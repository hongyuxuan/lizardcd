<template>
<el-card>
  <template #header>
    <div class="card-header" style="display:block;">
      <div>
        <b class="card-header-text" style="margin-right:20px">{{ route.params.name }}</b>
        <span style="color:#a0a0a0;font-size:15px" v-if="pipelineRunInfo.status">
          流水线启动于 {{moment.duration(moment(pipelineRunInfo.startTime)-moment()).humanize(true)}}
        </span>
        <el-dropdown @command="handleMore" v-if="pipelineRunInfo.conditions" class="pull-right">
          <el-button size="large">更多操作 <el-icon class="el-icon--right">
            <arrow-down />
          </el-icon></el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item :command="{action:'cancelled'}" :disabled="!['Running','Pending'].includes(pipelineRunInfo.conditions[0]?.reason)">终止运行</el-dropdown-item>
              <el-dropdown-item :command="{action:'rerun'}">重新运行</el-dropdown-item>
              <el-dropdown-item :command="{action:'params'}">查看启动参数</el-dropdown-item>
              <el-dropdown-item :command="{action:'yaml'}">查看YAML</el-dropdown-item>
              <el-dropdown-item :command="{action:'delete'}">删除</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
      <div style="margin-top:15px" v-if="pipelineRunInfo.conditions">
        <el-text :type="getColor(pipelineRunInfo.conditions[0]?.reason)">{{ pipelineRunInfo.conditions[0]?.reason }}</el-text>
        <el-text style="margin-left:20px">{{ pipelineRunInfo.conditions[0]?.message }}</el-text>
      </div>
      <div style="margin-top:10px"><el-tag v-for="(v,k,i) in pipelineRunInfo.annotations||{}" :key="i" size="large" type="warning">{{ k }}={{ v }}</el-tag></div>
      <div style="margin-top:10px"><div><el-tag v-for="(v,k,i) in pipelineRunInfo.labels||{}" :key="i" size="large">{{ k }}={{ v }}</el-tag></div></div>
    </div>
  </template>
  <div style="width:100%; height: calc(100vh - 375px)">
    <VueFlow 
      v-model="elements" 
      :node-types="nodeTypes" 
      :zoom-on-scroll="false"
      :zoom-on-pinch="false"
      :zoom-on-double-click="false"
      :pan-on-scroll="false" @node-click="handleNodeClick" />
  </div>
</el-card>
<el-drawer v-model="show.yaml" title="查看YAML" direction="rtl" size="700px">
  <div style="position:relative">
    <el-switch v-model="wrapLine" active-text="自动换行" inactive-text="不自动换行" inline-prompt class="el-switch-hover-right" size="large" />
    <v-ace-editor
      v-model:value="yamlContent"
      lang="yaml"
      theme="chrome"
      style="width:100%"
      :options="{
        wrap: wrapLine,
        enableBasicAutocompletion: true,
        enableSnippets: true,
        enableLiveAutocompletion: true,
        tabSize: 2,
        showPrintMargin: false,
        fontSize: 14,
        maxLines: 5000,
        minLines: 10,
        readOnly: true,
      }" />
  </div>
</el-drawer>
<el-drawer v-model="show.params" direction="rtl" title="查看启动参数"size="700px">
  <el-table :data="pipelineRunInfo.params" class="line-height25">
    <el-table-column label="名称" prop="name" min-width="150" />
    <el-table-column label="值" prop="value" min-width="350" />
  </el-table>
</el-drawer>
<el-dialog v-model="show.task" :title="currentTask?.displayName||''" @close="closeTask" width="80%">
  <template #header>
    <span style="font-size:18px">{{ currentTask?.displayName||'' }}</span>
    <span class="cell-comment">（{{ currentTask?.taskRef||'' }} ）</span>
  </template>
  <inner ref="refInner" />
</el-dialog>
</template>

<script setup>
import axios from 'axios';
import { onBeforeMount, onBeforeUnmount, ref, nextTick, markRaw, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getDuration, generateShortId, getColor } from '@/assets/util/common'
import inner from '../taskrun/result_inner.vue'
import { Duration } from 'luxon'
import moment from 'moment'
import yaml from 'js-yaml'
/* 引入v-ace-editor */
import { VAceEditor } from 'vue3-ace-editor'
import 'ace-builds/src-noconflict/mode-yaml'
import terminalUrl from 'ace-builds/src-noconflict/theme-terminal?url'
ace.config.setModuleUrl('ace/theme/terminal', terminalUrl)
import 'ace-builds/src-noconflict/ext-language_tools'
import 'ace-builds/src-noconflict/theme-chrome'
/* 引入vue-flow */
import { VueFlow } from '@vue-flow/core'
import '@vue-flow/core/dist/style.css'
import '@vue-flow/core/dist/theme-default.css'
import SpecialNode from '@/components/SpecialNode.vue'
import dagre from 'dagre'
/* 变量定义 */
const route = useRoute()
const router = useRouter()
const pipelineRunInfo = ref({})
const taskRunMap = ref({})
const timer = ref(null)
const currentTask = ref({})
const yamlContent = ref("")
const wrapLine = ref(true)
const show = ref({
  yaml: false,
  task: false
})
const refInner = ref(null)
const nodesAfter = ref({})
const nodes = ref([])
const edges = ref([])
const nodeTypes = {
  customTask: markRaw(SpecialNode),
}
const elements = ref([])
/* 生命周期函数 */
onBeforeMount(async () => {
  await getPipelineRun()
  timer.value = setInterval(async () => {
    try{
      await getPipelineRun()
      if(currentTask.value) {
        selectTask(currentTask.value.taskName)
      }
    }
    catch(e) {
      console.error(e)
    }
  }, 5000)
})
onBeforeUnmount(() => {
  if(timer) {
    clearInterval(timer.value)
    timer.value = null
  }
})
onMounted(() => {
  const layoutedNodes = applyDagreLayout(nodes.value, edges.value);
  elements.value = [...layoutedNodes, ...edges.value]
})
/* methods */
const getPipelineRun = async () => {
  nodes.value = []
  edges.value = []
  nodesAfter.value = {}
  pipelineRunInfo.value = await axios.get(`/tekton-pipelines/lizardcd/${route.query.namespace}/pipelineRuns/${route.params.name}`)
  for(let x of pipelineRunInfo.value.taskRuns) {
    x.expire = getDuration(x.endTime, x.startTime)
    let arr = x.taskRef?.split('/') || []
    if(arr.length > 1) {
      x.taskRefNamespace = arr[0]
      x.taskRefName = arr[1]
    } else {
      x.taskRefNamespace = undefined
      x.taskRefName = arr[0]
    }
    if(x.retriesStatus) {
      x.displayName += ` ( 重试 ${x.retriesStatus.length} 次 )`
    }
    for(let r of x.results || []) {
      if(!Array.isArray(r.value)) {
        if(r.value.startsWith("http://") || r.value.startsWith("https://")) {
          r.value = `<a style='color:var(--el-color-primary);text-decoration:underline;' href="${r.value}" target="_blank">${r.value}</a>`
        }
      }
    }
    taskRunMap.value[x.taskName] = x
    let status = null
    if(x.approvalStatus) {
      if(x.approvalStatus.status === 'pending' && x.conditions[0].reason !== 'Succeeded') {
        status = x.conditions[0].reason
      } else if(x.approvalStatus.status === 'completed') {
        status = x.labels.decision
      } else {
        status = x.approvalStatus.status
      }
    } else if(x.conditions) {
      status = x.conditions[0]?.reason
    }
    let node = {
      id: x.taskName,
      type: "customTask",
      data: {
        label: x.displayName||x.taskName,
        taskName: x.taskName,
        status: status,
      }
    }
    if(x.runAfter.length === 0) node.data.start = true
    nodes.value.push(node)
    if (x.runAfter) {
      for(let y of x.runAfter) {
        edges.value.push({
          id: `${y}-${x.taskName}`,
          source: y,
          target: x.taskName,
          type: 'default'
        })
        if(!nodesAfter.value.hasOwnProperty(y)) {
          nodesAfter.value[y] = []
        }
        nodesAfter.value[y].push(x.taskName)
      }
    }
  }
  for(let x of nodes.value) {
    if(!nodesAfter.value[x.data.taskName]) {
      x.data.end = true
    }
  }
  for(let x of edges.value) {
    if(taskRunMap.value[x.source]?.status === 'Running') {
      x.animated = true
    } else {
      x.animated = false
    }
  }
  const layoutedNodes = applyDagreLayout(nodes.value, edges.value);
  elements.value = [...layoutedNodes, ...edges.value]
  if(pipelineRunInfo.value.endTime) {
    clearInterval(timer.value)
    timer.value = null
  }
}
const selectTask = async (event) => {
  currentTask.value = taskRunMap.value[event]
  if(currentTask.value?.status) {
    nextTick(async () => {
      if(refInner.value) {
        await refInner.value.refreshTaskRun(currentTask.value)
        refInner.value.refreshCurrentStep()
      }
    })
  }
}
const handleMore = async (command) => {
  switch(command.action) {
    case "cancelled": {
      await axios.patch(`/lizardcd/tekton/cluster/${route.query.cluster}/namespace/${route.query.namespace}/pipelineruns/${route.params.name}`, [
        {
          "op": "replace",
          "path": "/spec/status",
          "value": "Cancelled"
        }
      ])
      ElMessage.success({message: `终止PipelineRun成功`})
      break
    }
    case "rerun": {
      let content = await axios.get(`/tekton-pipelines/lizardcd/${route.query.namespace}/pipelineRuns/${route.params.name}/yaml`)
      let params = yaml.load(content)
      params.metadata.name = pipelineRunInfo.value.pipelineRef + "-" + generateShortId(5)
      delete params.metadata.generateName
      delete params.spec.status
      try {
        await axios.post(`/lizardcd/tekton/cluster/${route.query.cluster}/namespace/${route.query.namespace}/apply?kind=PipelineRun`, {content: yaml.dump(params)})
        ElMessage.success({message: `重新执行PipelineRun成功`})
        setTimeout(() => {
          router.push({
            path: `/ci/tekton/pipelinerun/${params.metadata.name}`,
            query: route.query
          })
        }, 1000)
      } catch(e) {
        ElMessage.error({message: e})
      }
      break
    }
    case "yaml": {
      let response = await axios.get(`/tekton-pipelines/lizardcd/${route.query.namespace}/pipelineRuns/${route.params.name}/yaml`)
      let o = yaml.load(response)
      if(o.spec?.timeouts?.pipeline) {
        let timeout = parseInt(o.spec.timeouts.pipeline.slice(0, -2))
        const duration = Duration.fromMillis(timeout / 1e6)
        o.spec.timeouts.pipeline = duration.toFormat("h'h'm'm's's'")
      }
      yamlContent.value = yaml.dump(o)
      show.value.yaml = true
      break
    }
    case "params": {
      show.value.params = true
      break
    }
    case "delete": {
      await ElMessageBox.confirm('确定删除？','警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(async () => {
        await axios.delete(`/tekton-pipelines/lizardcd/${route.query.namespace}/pipelineRuns/${route.params.name}`)
        ElMessage.success({message: `删除PipelineRun成功`})
        setTimeout(async () => {
          router.push({
            path: `/ci/tekton`,
            query: {
              tab: "PipelineRun",
              namespace: route.query.namespace
            }
          })
        }, 1000)
      }).catch((e) => {
        console.warn(e)
      })
    }
  }
}
const applyDagreLayout = (nodes, edges) => {
  const dagreGraph = new dagre.graphlib.Graph();
  dagreGraph.setDefaultEdgeLabel(() => ({}));
  dagreGraph.setGraph({ rankdir: 'LR' }); // 方向：从左到右 (LR)

  nodes.forEach(node => {
    dagreGraph.setNode(node.id, { width: 150, height: 50 });
  });
  edges.forEach(edge => {
    dagreGraph.setEdge(edge.source, edge.target);
  });

  dagre.layout(dagreGraph);

  return nodes.map(node => {
    const dagreNode = dagreGraph.node(node.id);
    return { ...node, position: { x: dagreNode.x, y: dagreNode.y } };
  });
}
const handleNodeClick = ({event, node}) => {
  if(node.data.status) {
    show.value.task = true
  }
  selectTask(node.data.taskName)
}
const closeTask = () => {
  if(refInner.value) {
    refInner.value.close()
  }
}
</script>
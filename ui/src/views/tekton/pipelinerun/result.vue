<template>
<el-card>
  <template #header>
    <div class="card-header" style="display:block;">
        <div>
          <b class="card-header-text" style="margin-right:20px">{{ route.params.name }}</b>
          <span style="color:#a0a0a0;font-size:15px" v-if="pipelineRunInfo.status">
            流水线启动于 {{moment.duration(moment(pipelineRunInfo.status.startTime)-moment()).humanize(true)}}
          </span>
          <el-dropdown @command="handleMore" v-if="pipelineRunInfo.status?.conditions" class="pull-right">
            <el-button size="large">更多操作 <el-icon class="el-icon--right">
              <arrow-down />
            </el-icon></el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item :command="{action:'cancelled'}" :disabled="!['Running','Pending'].includes(pipelineRunInfo.status.conditions[0]?.reason)">终止运行</el-dropdown-item>
                <el-dropdown-item :command="{action:'rerun'}">重新运行</el-dropdown-item>
                <el-dropdown-item :command="{action:'params'}">查看启动参数</el-dropdown-item>
                <el-dropdown-item :command="{action:'yaml'}">查看YAML</el-dropdown-item>
                <el-dropdown-item :command="{action:'delete'}">删除</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
        <div style="margin-top:15px" v-if="pipelineRunInfo.status">
          <el-text :type="getColor(pipelineRunInfo.status.conditions[0]?.reason)">{{ pipelineRunInfo.status.conditions[0]?.reason }}</el-text>
          <el-text style="margin-left:20px">{{ pipelineRunInfo.status.conditions[0]?.message }}</el-text>
        </div>
        <div style="margin-top:10px"><el-tag v-for="(v,k,i) in pipelineRunInfo.metadata?.annotations||{}" :key="i" size="large" type="warning">{{ k }}={{ v }}</el-tag></div>
        <div style="margin-top:10px"><div><el-tag v-for="(v,k,i) in pipelineRunInfo.metadata?.labels||{}" :key="i" size="large">{{ k }}={{ v }}</el-tag></div></div>
    </div>
  </template>
  <el-row :gutter="15" v-if="pipelineRunInfo.status?.pipelineSpec">
    <el-col :span="6">
      <el-menu
        v-if="pipelineRunInfo.status"
        :default-active="activeName"
        background-color="#f0f0f0"
        unique-opened
        @select="selectStep"
        @open="selectTask"
        @close="selectTask">
        <el-sub-menu v-for="(item,i) in pipelineRunInfo.status.pipelineSpec.tasks" :key="i" :index="item.name">
          <template #title>
            <!-- approvetsk -->
            <el-icon v-if="taskRunMap[item.name]?.status.status==='pending'" :size="20" class="is-loading"><Refresh style="color:#e6a23c" /></el-icon>
            <el-icon v-else-if="taskRunMap[item.name]?.status.status==='completed'&&taskRunMap[item.name]?.status.results.decision==='true'" :size="20"><CircleCheckFilled style="color:green" /></el-icon>
            <el-icon v-else-if="taskRunMap[item.name]?.status.status==='completed'&&taskRunMap[item.name]?.status.results.decision==='false'" :size="20"><CircleCloseFilled style="color:red" /></el-icon>
            <!-- approvetsk end -->
            <el-icon v-else-if="taskRunMap[item.name]?.status.conditions&&taskRunMap[item.name]?.status.conditions[0].reason==='Succeeded'" :size="20"><CircleCheckFilled style="color:green" /></el-icon>
            <el-icon v-else-if="taskRunMap[item.name]?.status.conditions&&taskRunMap[item.name]?.status.conditions[0].reason==='Pending'" :size="20" class="is-loading"><Refresh style="color:#e6a23c" /></el-icon>
            <el-icon v-else-if="taskRunMap[item.name]?.status.conditions&&taskRunMap[item.name]?.status.conditions[0].reason==='Running'" :size="20" class="is-loading"><Refresh style="color:#409eff" /></el-icon>
            <el-icon v-else-if="taskRunMap[item.name]?.status.conditions&&taskRunMap[item.name]?.status.conditions[0].reason==='TaskRunCancelled'" :size="18"><font-awesome-icon icon="ban" style="color:#a0a0a0" /></el-icon>
            <el-icon v-else-if="!taskRunMap[item.name]||!taskRunMap[item.name]?.status?.conditions" :size="20"><RemoveFilled style="color:#a0a0a0" /></el-icon>
            <el-icon v-else :size="20"><CircleCloseFilled style="color:red" /></el-icon>
            <span>{{ item.displayName || item.name }}</span>
          </template>
          <div v-if="taskRunMap[item.name]"><!-- already run task -->
            <el-menu-item v-for="(y,i) in taskRunMap[item.name].status.steps" :key="i" :index="`${item.name}#${y.name}`">
              <el-icon v-if="y.terminated&&y.terminated.exitCode===0" :size="20"><CircleCheck style="color:green" /></el-icon>
              <el-icon v-else-if="y.terminated&&y.terminated.exitCode!==0&&y.terminated.reason==='TaskRunCancelled'" :size="18"><font-awesome-icon icon="ban" style="color:#a0a0a0" /></el-icon>
              <el-icon v-else-if="y.terminated&&y.terminated.exitCode!==0&&y.terminated.reason!=='TaskRunCancelled'" :size="20"><CircleClose style="color:red" /></el-icon>
              <el-icon v-else-if="y.waiting" :size="20" class="is-loading"><Refresh style="color:#e6a23c" /></el-icon>
              <el-icon v-else-if="(i==0&&y.running)||(i>0&&y.running&&taskRunMap[item.name].status.steps[i-1].terminated)" :size="20" class="is-loading"><Refresh style="color:#409eff" /></el-icon>
              <el-icon v-else :size="20"><RemoveFilled style="color:#a0a0a0" /></el-icon>
              <span>{{ y.name }}</span>
            </el-menu-item>
          </div>
          <div v-else><!-- task not run -->
            <el-menu-item v-for="(y,i) in taskMap[item.taskRef?.name||item.name]?.spec.steps||[]" :key="i" :index="`${item.name}#${y.name}`">
              <el-icon :size="20"><RemoveFilled style="color:#a0a0a0" /></el-icon>
              <span>{{ y.name }}</span>
            </el-menu-item>
          </div>
        </el-sub-menu>
      </el-menu>
    </el-col>
    <el-col :span="18">
      <el-card v-if="!showStep"> <!-- 点击task -->
        <template #header>
          <div class="card-header" style="display:block;" v-if="currentTask?.status">
            <span class="card-header-text" style="margin-right:20px"><b>{{ currentTask.metadata.name }}</b></span>
            <el-text v-if="currentTask.status.conditions" :type="getColor(currentTask.status.conditions[0]?.reason)">
              {{ currentTask.status.conditions[0]?.reason }}
            </el-text>
            <el-text v-if="currentTask.kind==='ApprovalTask'" :type="getColor2(currentTask.status)">{{ currentTask.metadata.labels.decision }}</el-text>
            <el-text v-if="currentTask.status.conditions" style="margin-left:20px" type="info">持续时间 {{ currentTask.expire }}</el-text>
            <el-dropdown @command="handleMoreTask" v-if="currentTask.status.conditions" class="pull-right" style="top:3px">
              <el-link underline="never">
                更多操作 <el-icon class="el-icon--right">
                  <arrow-down />
                </el-icon>
            </el-link>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item :command="{action:'cancelled',name:currentTask.metadata.name}" :disabled="!['Running','Pending'].includes(currentTask.status.conditions[0]?.reason)">终止运行</el-dropdown-item>
                  <el-dropdown-item :command="{action:'rerun',name:currentTask.metadata.name}">重新运行</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
          <div class="card-header" style="display:block;" v-else-if="currentTask?.taskRunName">
            <span class="card-header-text" style="margin-right:20px"><b>{{ currentTask.taskRunName }}</b></span>
            <el-text type="info">暂挂中</el-text>
          </div>
        </template>
        <el-tabs v-model="taskActive" v-if="currentTask?.kind!=='ApprovalTask'">
          <el-tab-pane label="参数" name="1" v-if="currentTask?.spec&&currentTask?.status">
            <el-table  :data="currentTask.spec.params" class="line-height25">
              <el-table-column label="名称" prop="name" min-width="150" />
              <el-table-column label="值" prop="value" min-width="350" />
            </el-table>
          </el-tab-pane>
          <el-tab-pane label="结果" name="2" v-if="currentTask?.status">
            <el-table :data="currentTask.status.results" class="line-height25">
              <el-table-column label="名称" prop="name" min-width="150" />
              <el-table-column label="值" min-width="350">
                <template #default="props"><div v-html="props.row.value" /></template>
              </el-table-column>
            </el-table>
          </el-tab-pane>
          <el-tab-pane label="状态" name="3" v-if="currentTask.status">
            <v-ace-editor
              v-model:value="statusManifest"
              lang="yaml"
              theme="terminal"
              style="width:100%"
              :options="{
                enableBasicAutocompletion: true,
                enableSnippets: true,
                enableLiveAutocompletion: true,
                tabSize: 2,
                showPrintMargin: false,
                fontSize: 14,
                maxLines: 5000,
                minLines: 1,
              }" />
          </el-tab-pane>
          <el-tab-pane label="Pod" name="4" v-if="currentTask.status">
            <v-ace-editor
              v-model:value="podManifest"
              lang="yaml"
              theme="terminal"
              style="width:100%"
              :options="{
                enableBasicAutocompletion: true,
                enableSnippets: true,
                enableLiveAutocompletion: true,
                tabSize: 2,
                showPrintMargin: false,
                fontSize: 14,
                maxLines: 5000,
                minLines: 1,
              }" />
          </el-tab-pane>
          <el-tab-pane label="事件" name="5" v-if="currentTask.status">
            <el-table :data="podEvents" style="width:100%" class="line-height25">
              <el-table-column prop="type" label="Type">
                <template #default="scope">
                  <el-tag v-if="scope.row.type==='Warning'" type="warning" size="large">{{ scope.row.type }}</el-tag>
                  <el-tag v-else-if="scope.row.type==='Normal'" type="primary" size="large">{{ scope.row.type }}</el-tag>
                  <el-tag v-else type="info" size="large">{{ scope.row.type }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="reason" label="Reason" />
              <el-table-column prop="age" label="Age" />
              <el-table-column prop="source.component" label="From" />
              <el-table-column prop="message" label="Message" min-width="300px" />
            </el-table>
          </el-tab-pane>
        </el-tabs>
        <el-descriptions title="" :column="1" size="large" border v-else-if="currentTask?.kind==='ApprovalTask'&&currentTask?.status.status==='completed'">
          <el-descriptions-item label="Description" label-width="150">{{ currentTask.spec.description }}</el-descriptions-item>
          <el-descriptions-item label="CreationTimestamp" label-width="150">{{ moment(currentTask.metadata.creationTimestamp).format('YYYY-MM-DD HH:mm:ss.SSS') }}</el-descriptions-item>
          <el-descriptions-item label="Labels" label-width="150">
            <el-tag v-for="(v,k,i) in currentTask.metadata.labels||{}" :key="i" size="large" type="warning">{{ k }}={{ v }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="Status" label-width="150">{{ currentTask.status.status }}</el-descriptions-item>
          <el-descriptions-item label="Results" label-width="150">{{ currentTask.status.results }}</el-descriptions-item>
        </el-descriptions>
        <iframe v-else-if="currentTask?.kind==='ApprovalTask'&&currentTask?.status.status==='pending'" :src="`${currentTask.status.approvalUrl}?user=${route.query.namespace}`" style="width:100%;min-height:650px;border:none" />
      </el-card>
      <el-card v-else> <!-- 点击step -->
        <template #header>
          <div class="card-header" style="display:block;" v-if="!currentStep.notRun">
            <span class="card-header-text" style="margin-right:20px"><b>{{ currentStep.name }}</b></span>
            <el-text v-if="currentStep.terminated&&currentStep.terminated.exitCode===0" type="success">{{ currentStep.terminated.reason }}</el-text>
            <el-text v-else-if="currentStep.terminated&&currentStep.terminated.exitCode!==0" type="danger">{{ currentStep.terminated.reason }}</el-text>
            <el-text v-else-if="currentStep.waiting" type="warning">{{ currentStep.waiting.reason }}</el-text>
            <el-text v-else-if="currentStep.running" type="primary">Running</el-text>
            <el-text style="margin-left:20px" type="info">持续时间 {{ currentStep.expire||'' }}</el-text>
          </div>
          <div class="card-header" style="display:block;" v-else>
            <span class="card-header-text" style="margin-right:20px"><b>{{ currentStep.name }}</b></span>
            <el-text type="info">未运行</el-text>
          </div>
        </template>
        <el-tabs v-model="stepActive" >
          <el-tab-pane label="日志" name="1">
            <pre class="dark fullscreen" ref="preRef">
              <p v-for="(item,i) in logs" :key="i" v-html="item" />
              <div v-if="currentStep.terminated?.exitCode===0" style="color:#24a148;font-weight:bold">步骤已完成</div>
              <div v-if="currentStep.terminated&&currentStep.terminated.exitCode!==0" style="color:#da1e28;font-weight:bold">步骤失败</div>
              <div v-if="currentStep.running"><el-icon class="is-loading" size="18"><Loading /></el-icon></div>
            </pre>
          </el-tab-pane>
          <el-tab-pane label="详情" name="2">
            <v-ace-editor
              v-model:value="stepSpec[currentStep.name]"
              lang="yaml"
              theme="terminal"
              style="width:100%"
              :options="{
                enableBasicAutocompletion: true,
                enableSnippets: true,
                enableLiveAutocompletion: true,
                tabSize: 2,
                showPrintMargin: false,
                fontSize: 14,
                maxLines: 5000,
                minLines: 1,
              }" />
          </el-tab-pane>
        </el-tabs>
      </el-card>
    </el-col>
  </el-row>
</el-card>
<el-drawer v-model="show.yaml" direction="rtl" size="700px">
  <template #header>
    <h4>查看YAML</h4>
  </template>
  <template #default>
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
  </template>
</el-drawer>
<el-drawer v-model="show.params" direction="rtl" size="700px">
  <template #header>
    <h4>查看启动参数</h4>
  </template>
  <template #default>
    <el-table :data="pipelineRunInfo.spec?.params" class="line-height25">
      <el-table-column label="名称" prop="name" min-width="150" />
      <el-table-column label="值" prop="value" min-width="350" />
    </el-table>
  </template>
</el-drawer>
</template>

<script setup>
import { Refresh, CircleCheckFilled, CircleCloseFilled, CircleClose, RemoveFilled, Loading, ArrowDown } from '@element-plus/icons-vue'
import axios from 'axios';
import { onBeforeMount, onBeforeUnmount, ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getDuration, generateShortId, getColor, getColor2 } from '@/assets/util/common'
import moment from 'moment'
import _ from 'lodash'
import yaml from 'js-yaml'
import { AnsiUp } from 'ansi_up'
/* 引入v-ace-editor */
import { VAceEditor } from 'vue3-ace-editor'
import 'ace-builds/src-noconflict/mode-yaml'
import terminalUrl from 'ace-builds/src-noconflict/theme-terminal?url'
ace.config.setModuleUrl('ace/theme/terminal', terminalUrl);
import 'ace-builds/src-noconflict/ext-language_tools'
/* 变量定义 */
const route = useRoute()
const router = useRouter()
const activeName = computed(() => {
  let taskName = pipelineRunInfo.value?.status.pipelineSpec.tasks[0].name
  if(taskName) {
    if(taskRunMap.value[taskName]?.status.steps)
      return taskRunMap.value[taskName]?.status.steps[0].name
    return ''
  }
  return ''
})
const pipelineRunInfo = ref({})
const taskRunMap = ref({})
const taskMap = ref({})
const taskActive = ref("1")
const stepActive = ref("1")
const currentTask = ref({})
const currentStep = ref({})
const showStep = ref(false)
const statusManifest = ref('')
const podManifest = ref('等待Pod资源')
const podEvents = ref([])
const logs = ref([])
const timer = ref(null)
const stepSpec = ref({})
const ansi_up = new AnsiUp()
const regWithHttp = /((https?:\/\/)?(www\.)?(?!\d+\.\d+)([a-zA-Z0-9-]+(\.[a-zA-Z]{2,})+)([\/?#][^\s]*)?)/g
const regWithOutHttp = /((www\.)?(?!\d+\.\d+)([a-zA-Z0-9-]+(\.[a-zA-Z]{2,})+)([\/?#][^\s]*)?)/g
const yamlContent = ref("")
const wrapLine = ref(true)
const show = ref({
  yaml: false
})
/* 生命周期函数 */
onBeforeMount(async () => {
  await getPipelineRun()
  await getTasks()
  await getTaskRun()
  selectTask(pipelineRunInfo.value.status.pipelineSpec.tasks[0].name)
  timer.value = setInterval(async () => {
    await getPipelineRun()
    await getTaskRun()
    await refreshCurrentTask(currentTask.value.taskRunName)
    if(currentStep.value) await refreshCurrentStep(currentTask.value.taskRunName + "#" + currentStep.value.name)
  }, 5000)
})
onBeforeUnmount(() => {
  if(timer) {
    clearInterval(timer.value)
    timer.value = null
  }
})
/* methods */
const getPipelineRun = async () => {
  pipelineRunInfo.value = await axios.get(`/lizardcd/tekton/cluster/${route.query.cluster}/namespace/${route.query.namespace}/pipelineruns/${route.params.name}`)
  if(pipelineRunInfo.value.status?.pipelineSpec?.tasks) {
    pipelineRunInfo.value.status.pipelineSpec.tasks = pipelineRunInfo.value.status.pipelineSpec.tasks.concat(pipelineRunInfo.value.status.pipelineSpec.finally||[])
    for(let x of pipelineRunInfo.value.status.pipelineSpec.tasks) {
      if(x.taskRef && !x.taskRef.hasOwnProperty("name") && x.taskRef.params) {
        x.taskRef.name = x.taskRef.params.find(n => n.name == "name")?.value
      }
    }
  }
  if(pipelineRunInfo.value.status.completionTime) {
    clearInterval(timer.value)
    timer.value = null
  }
}
const getTasks = async () => {
  let tasks = pipelineRunInfo.value.status.pipelineSpec.tasks.filter(n => n.hasOwnProperty("taskRef") && n.taskRef.kind == "Task")
  let response = await Promise.all(tasks.map(x => {
    let namespace = route.query.namespace
    if(x.taskRef.params) {
      namespace = x.taskRef.params.find(n => n.name === 'namespace')?.value || route.query.namespace
    }
    return axios.get(`/lizardcd/tekton/cluster/${route.query.cluster}/namespace/${namespace}/tasks/${x.taskRef.name}`)
  }))
  for(let x of response) {
    taskMap.value[x.metadata.name] = x
  }
  if(pipelineRunInfo.value.status.pipelineSpec.finally) {
    for(let x of pipelineRunInfo.value.status.pipelineSpec.finally) {
      taskMap.value[x.name] = {
        name: x.name,
        spec: x.taskSpec
      }
    }
  }
}
const getTaskRun = async () => {
  let response = await axios.get(`/lizardcd/tekton/cluster/${route.query.cluster}/namespace/${route.query.namespace}/taskruns?label_selector=tekton.dev/pipelineRun=${route.params.name}`)
  for(let x of _.sortBy(response.results, 'metadata.creationTimestamp')) {
    taskRunMap.value[x.metadata.labels['tekton.dev/pipelineTask']] = x
    x.expire = getDuration(x.status.completionTime, x.status.startTime)
    for(let y of x.status.taskSpec?.steps||[]) {
      stepSpec.value[y.name] = yaml.dump(y)
    }
    // if has retriesStatus
    if(x.status.retriesStatus) { // x = taskrun
      let name = x.spec.taskRef?.name || x.spec.taskSpec?.name
      if(x.spec.taskRef?.params)
        name = x.spec.taskRef.params.find(m => m.name === 'name')?.value
      let task = pipelineRunInfo.value.status.pipelineSpec.tasks.find(n => {
        return n.name === name
      })
      if(task)
        task.displayName += ` ( 重试 ${x.status.retriesStatus.length} 次 )`
    }
    for(let r of x.status.results || []) {
      if(!Array.isArray(r.value)) {
        if(r.value.startsWith("http://") || r.value.startsWith("https://")) {
          r.value = `<a style='color:#78a9ff;text-decoration:underline;' href="${r.value}" target="_blank">${r.value}</a>`
        }
      }
    }
  }
  // get approvaltasks
  let approvaltasks = pipelineRunInfo.value.status.childReferences?.filter(n => n.kind == "CustomRun")||[]
  for(let x of approvaltasks) {
    try {
      let response = await axios.get(`/lizardcd/tekton/cluster/${route.query.cluster}/namespace/${route.query.namespace}/approvaltasks/${x.name}`)
      let name = response.metadata.name.substring(route.params.name.length+1)
      taskRunMap.value[name] = response
    } catch(e) {}
  }
}
const selectTask = async (event) => {
  showStep.value = false
  podManifest.value = '等待Pod资源'
  podEvents.value = []
  await refreshCurrentTask(event)
}
const selectStep = async (event) => {
  showStep.value = true
  currentStep.value = null
  logs.value = []
  await refreshCurrentStep(event)
}
const refreshCurrentTask = async (event) => {
  currentTask.value = taskRunMap.value[event]
  if(currentTask.value?.kind === "ApprovalTask") {
    currentTask.value.taskRunName = event
    return
  }
  if(currentTask.value) {
    statusManifest.value = await axios.get(`/lizardcd/tekton/cluster/${route.query.cluster}/namespace/${route.query.namespace}/taskruns/${currentTask.value.metadata.name}/yaml?withStatus=true`)
    if(currentTask.value.status?.podName) {
      podManifest.value = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/pods/${currentTask.value.status.podName}/yaml`)
      let response = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/Pod/${currentTask.value.status.podName}/events`)
      podEvents.value = response.map(x => {
        x.age = moment.duration(moment(x.lastTimestamp)-moment()).humanize(true)
        return x
      })
    }
    currentTask.value.taskRunName = event
  } else {
    let found = pipelineRunInfo.value.status.pipelineSpec.tasks.find(x => x.name === event)
    if(found) currentTask.value = taskMap.value[found.taskRef?.name||found.name]
    if(currentTask.value)
      currentTask.value.taskRunName = event
    else {
      currentTask.value = {
        taskRunName: event
      }
    }
  }
}
const refreshCurrentStep = async (event) => {
  if(currentTask.value.status?.steps) {
    currentStep.value = currentTask.value.status.steps.find(n => currentTask.value.taskRunName + "#" + n.name === event)
    if(currentStep.value?.terminated) {
      currentStep.value.expire = getDuration(currentStep.value.terminated.finishedAt, currentStep.value.terminated.startedAt)
    } else if(currentStep.value?.running) {
      currentStep.value.expire = getDuration(moment(), currentStep.value.running.startedAt)
    }
    getLogs()
    if(currentStep.value) currentStep.value.notRun = false
  } else {
    currentStep.value = {
      notRun: true,
      name: event.split("#")[1]
    }
    logs.value = ['没有日志输出']
  }
}
const getLogs = async () => {
  if(currentStep.value?.container) {
    let response = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/pods/${currentTask.value.status.podName}/logs?container=${currentStep.value.container}&timestamps=true`)
    logs.value = response.split('\n').map(x => {
      if(regWithHttp.test(x))
        return ansi_up.ansi_to_html(x).replace(regWithHttp, `<a style='color:#78a9ff;text-decoration:underline;' href="$1" target="_blank">$1</a>`)
      else
        return ansi_up.ansi_to_html(x).replace(regWithOutHttp, `<a style='color:#78a9ff;text-decoration:underline;' href="http://$1" target="_blank">$1</a>`)
    })
    if(response === '') {
      logs.value = ['没有日志输出']
    }
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
      let content = await axios.get(`/lizardcd/tekton/cluster/${route.query.cluster}/namespace/${route.query.namespace}/pipelineruns/${route.params.name}/yaml`)
      let params = yaml.load(content)
      params.metadata.name = pipelineRunInfo.value.spec.pipelineRef.name + "-" + generateShortId(5)
      delete params.metadata.generateName
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
      yamlContent.value = await axios.get(`/lizardcd/tekton/cluster/${route.query.cluster}/namespace/${route.query.namespace}/pipelineruns/${route.params.name}/yaml`)
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
        await axios.delete(`/lizardcd/tekton/cluster/${route.query.cluster}/namespace/${route.query.namespace}/pipelineruns/${route.params.name}`)
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
const handleMoreTask = async (command) => {
  switch(command.action) {
    case "cancelled": {
      await axios.patch(`/lizardcd/tekton/cluster/${route.query.cluster}/namespace/${route.query.namespace}/taskruns/${command.name}`, [
        {
          "op": "replace",
          "path": "/spec/status",
          "value": "TaskRunCancelled"
        }
      ])
      ElMessage.success({message: `终止TaskRun成功`})
      break
    }
    case "rerun": {
      let content = await axios.get(`/lizardcd/tekton/cluster/${route.query.cluster}/namespace/${route.query.namespace}/taskruns/${command.name}/yaml`)
      let params = yaml.load(content)
      params.metadata.name = params.metadata.labels["tekton.dev/pipelineTask"] + "-" + generateShortId(5)
      delete params.metadata.generateName
      delete params.spec.status
      delete params.spec.statusMessage
      if(params.metadata.annotations)
        delete params.metadata.annotations["pipeline.tekton.dev/affinity-assistant"]
      if(params.spec.taskRef) {
        delete params.spec.taskRef.namespace
        if(params.spec.taskRef.resolver) {
          delete params.spec.taskRef.name
        }
      }
      try {
        await axios.post(`/lizardcd/tekton/cluster/${route.query.cluster}/namespace/${route.query.namespace}/apply?kind=TaskRun`, {content: yaml.dump(params)})
        ElMessage.success({message: `重新执行TaskRun成功`})
        setTimeout(() => {
          router.push({
            path: `/ci/tekton/taskrun/${params.metadata.name}`,
            query: route.query
          })
        }, 1000)
      } catch(e) {
        ElMessage.error({message: e})
      }
      break
    }
  }
}
</script>
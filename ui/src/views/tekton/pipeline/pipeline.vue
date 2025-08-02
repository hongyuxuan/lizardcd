<template>
<div class="box box-item">
  <div class="box-body" style="padding-top:20px;padding-bottom:0">
    <el-alert type="info" show-icon style="margin-bottom:15px">
      <template #title>
        关于 Pipeline 配置参考：<el-link href="https://tekton.dev/docs/pipelines/pipelines/" underline="never" type="primary" target="_blank">Pipelines</el-link>
      </template>
    </el-alert>
    <el-row>
      <div style="flex:1">
        <el-button-group>
          <el-button :icon="Refresh" size="large" @click="getList(true)" />
          <el-select v-model="namespace" placeholder="请选择命名空间" clearable filterable @change="getList(true)" style="width:250px;" size="large">
            <el-option v-for="(item) in props.namespaceList" :key="item" :label="item" :value="item" />
          </el-select>
          <el-select 
            v-model="inputLabels"
            multiple 
            clearable 
            allow-create 
            default-first-option 
            :reserve-keyword="false" 
            filterable 
            style="width:250px" 
            placeholder="输入标签 key=value 过滤……" 
            @change="getList(true);current=1;labelOptions=_.uniq(labelOptions.concat(inputLabels))"
            size="large">
            <el-option v-for="item in labelOptions" :key="item" :label="item" :value="item" />
          </el-select>
          <el-input v-model="searchKey" placeholder="输入名称进行搜索" size="large" :prefix-icon="Search" @change="getPage(1)" clearable style="width:250px;" />
          <el-radio-group v-model="pipelineType" size="large" style="vertical-align:top;margin-left:15px" @change="getList(true)">
            <el-radio value="pre-check">Precheck流水线</el-radio>
            <el-radio value="release">Release流水线</el-radio>
          </el-radio-group>
        </el-button-group>
      </div>
      <div>
        <el-dropdown @command="handleMore" class="pull-right">
          <el-button size="large">更多操作<el-icon class="el-icon--right"><arrow-down /></el-icon></el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item :command="{action:'startBatch'}">批量启动</el-dropdown-item>
              <el-dropdown-item :command="{action:'deleteBatch'}">批量删除</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-dropdown @command="handleMore" class="pull-right">
          <el-button size="large" type="primary">+ 创建Pipeline</el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item :command="{action:'create'}">模板创建</el-dropdown-item>
              <el-dropdown-item :command="{action:'createVisual'}">可视化创建</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </el-row>
    <el-table 
      :data="list" 
      v-loading="loading"
      element-loading-text="奋力加载中..."
      @selection-change="select"
      class="line-height40" 
      style="width:100%;margin-top:10px;min-height:150px">
      <el-table-column type="selection" width="45" />
      <el-table-column prop="metadata.name" label="名称" min-width="200" />
      <el-table-column label="注解" min-width="250">
        <template #default="scope">
          <el-tag v-for="(v,k,i) in scope.row.metadata.annotations||{}" :key="i" size="large" type="warning">{{ k }}={{ v }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="标签" min-width="300">
        <template #default="scope">
          <el-tag v-for="(v,k,i) in scope.row.metadata.labels||{}" :key="i" size="large">{{ k }}={{ v }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="最近一次运行" width="140">
        <template #default="scope">
          <div v-if="lastRun[scope.row.metadata.name]" @click="gotoLast(lastRun[scope.row.metadata.name])" class="pointer">
            <el-tooltip placement="top">
              <template #content>
                开始时间: {{ lastRun[scope.row.metadata.name].startTime }} <br>
                结束时间: {{ lastRun[scope.row.metadata.name].endTime}} <br>
              </template>
              <el-progress v-if="['Succeeded','Completed'].includes(lastRun[scope.row.metadata.name].conditions[0].reason)" :percentage="100" color="#5cb87a" :show-text="false" />
              <el-progress v-else-if="lastRun[scope.row.metadata.name]?.conditions[0].reason=='Cancelled'" :percentage="100" color="#bfbbbb" :show-text="false" />
              <el-progress v-else-if="lastRun[scope.row.metadata.name].conditions[0].status=='Unknown'" :percentage="50" color="#e6a23c" :show-text="false" /> <!-- 黄色 -->
              <el-progress v-else :percentage="100" color="#f56c6c" :show-text="false" /> <!-- 红色 -->
            </el-tooltip>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" width="170">
        <template #default="scope">
          {{ moment(scope.row.metadata.creationTimestamp).format('YYYY-MM-DD HH:mm:ss') }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="170">
        <template #default="scope">
          <el-button icon="EditPen" circle @click="editOne(scope.row)"></el-button>
          <el-tooltip content="启动流水线" placement="top">
            <el-button icon="ArrowRight" circle @click="batchStart=false;start(scope.row)"></el-button>
          </el-tooltip>
          <el-tooltip content="执行历史" placement="top">
            <el-button circle @click="gotoHistory(scope.row)">
              <font-awesome-icon icon="clock-rotate-left" />
            </el-button>
          </el-tooltip>
          <el-popconfirm title="确定删除?" confirm-button-text="确认" cancel-button-text="取消" @confirm="deleteOne(scope.row)">
            <template #reference>
              <el-button icon="Close" circle />
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination 
      class="pull-right"
      background 
      v-model:page-size="pageSize"
      :page-sizes="[10, 20, 50, 100]"
      layout="total, sizes, prev, pager, next, jumper" 
      :total="pageTotal"
      @size-change="handleSizeChange"
      @current-change="getPage"
      v-model:current-page="current" />
  </div>
</div>
<el-drawer v-model="show.yaml" direction="rtl" size="700px">
  <template #header>
    <h4>{{ edit === true ? '编辑YAML' : '创建Pipeline' }}</h4>
  </template>
  <template #default>
    <el-form ref="task" label-width="100px">
      <el-form-item label="从模板导入" v-if="edit===false">
        <el-select 
          v-model="form.templates" 
          placeholder="请选择模板" 
          value-key="id" 
          clearable 
          size="large"
          style="width:100%" 
          @change="selectTemplate">
          <el-option v-for="item in templateList" :key="item.id" :label="item.name" :value="item">
            <span style="float:left">{{item.name}}</span>
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="配置YAML">
        <el-switch v-model="wrapLine" active-text="自动换行" inactive-text="不自动换行" inline-prompt class="el-switch-hover-right" size="large" />
        <v-ace-editor
          v-model:value="form.content"
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
        }" />
      </el-form-item>
      <el-form-item label="模板变量" v-if="edit===false">
        <table class="table table-bordered" style="margin-bottom:0">
          <thead><tr><th>变量名</th><th>变量值</th></tr></thead>
          <tbody>
          <tr v-for="(item,index) in form.variables" :key="index" >
            <td style="vertical-align: top"><el-input v-model="item.key" size="large" /></td>
            <td><el-input v-model="item.value" size="large" /><myTips type="info">{{ item.memo }}</myTips></td>
            <td width="80">
              <el-button-group>
                <el-button icon="Plus" circle @click="addVar(index)"></el-button>
                <el-button icon="Close" circle @click="removeVar(index)"></el-button>
              </el-button-group>
            </td>
          </tr>
          </tbody>
        </table>
      </el-form-item>
    </el-form>
  </template>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show.yaml=false">取消</el-button>
      <el-button type="primary" @click="submitYaml">提交</el-button>
    </div>
  </template>
</el-drawer>
<el-drawer v-model="show.start" direction="rtl" size="700px">
  <template #header>
    <h4 v-if="batchStart===false">启动Pipeline</h4>
    <h4 v-else>批量启动Pipeline</h4>
  </template>
  <template #default>
    <el-alert v-if="batchStart===true" style="margin-bottom:20px" type="warning" show-icon title="批量启动要求所有流水线具有相同的必填参数，否则可能造成运行失败" />
    <el-form ref="task" label-width="100px">
      <el-form-item label="启动模式">
        <el-radio-group v-model="startMode">
          <el-radio value="1">表单模式</el-radio>
          <el-radio value="2" :disabled="batchStart">YAML模式</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="配置YAML" v-if="startMode==='2'">
        <v-ace-editor
          v-model:value="startContent"
          lang="yaml"
          theme="chrome"
          style="width:100%"
          :options="{
            enableBasicAutocompletion: true,
            enableSnippets: true,
            enableLiveAutocompletion: true,
            tabSize: 2,
            showPrintMargin: false,
            fontSize: 14,
            maxLines: 5000,
            minLines: 10,
        }" />
      </el-form-item>
      <el-form-item label="启动参数" v-if="startMode==='1'">
        <el-form label-width="120px" label-position="top" style="width:100%">
          <el-form-item v-for="(item,i) in formParams" :key="i" :label="item.name">
            <el-input v-model="item.value" size="large" v-if="!item.description?.startsWith('options:')" />
            <myTips type="info" v-if="!item.description?.startsWith('options:')"><div v-html="item.description" /></myTips>
            <el-select v-model="item.value" v-else size="large">
              <el-option v-for="y in item.options" :key="y" :label="y" :value="y" />
            </el-select>
          </el-form-item>
        </el-form>
      </el-form-item>
      <el-form-item label="超时时间" v-if="startMode==='1'">
        <el-input v-model="timeout" size="large" />
        <myTips type="info">0h0m0s表示永不超时</myTips>
      </el-form-item>
    </el-form>
  </template>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show.start=false">取消</el-button>
      <el-button type="primary" v-if="batchStart===false" @click="submitStart">提交</el-button>
      <el-button type="primary" v-else @click="submitStartBatch">提交</el-button>
    </div>
  </template>
</el-drawer>
</template>
<script setup>
import { Search,Refresh } from '@element-plus/icons-vue'
import { useRouter, useRoute } from 'vue-router'
import { onBeforeMount, onMounted, ref, computed } from 'vue'
import { useStore } from 'vuex'
import { ElMessage, ElMessageBox } from 'element-plus'
import { axios } from '/src/assets/util/axios'
import moment from 'moment'
import _ from 'lodash'
import yaml from 'js-yaml'
/* 引入v-ace-editor */
import { VAceEditor } from 'vue3-ace-editor'
import 'ace-builds/src-noconflict/mode-yaml'
import 'ace-builds/src-noconflict/theme-chrome'
import 'ace-builds/src-noconflict/ext-language_tools'
/* 变量定义 */
const props = defineProps({
  defaultTekton: { type: Object }, 
  namespaceList: { type: Array }, 
})
const store = useStore()
const role = computed(() => {
  return store.state.userInfo.role
})
const router = useRouter()
const route = useRoute()
const all = ref([])
const searchKey = ref("")
const list = ref([])
const pageSize = ref(10)
const pageTotal = ref(0)
const current = ref(1)
const namespace = ref("")
const loading = ref(false)
const show = ref({
  yaml: false
})
const edit = ref(false)
const form = ref({content:'',variables:[],templates:{}})
const templateList = ref([])
const startContent = ref("")
const startMode = ref("1")
const formParams = ref([])
const timeout = ref('0h0m0s')
const selected = ref([])
const batchStart = ref(false)
const pipelineRunObject = ref(null)
const pipelineRunObjects = ref([])
const lastRun = ref({})
const inputLabels = ref([])
const labelSelector = ref([])
const labelOptions = ref([])
const wrapLine = ref(true)
const pipelineType = ref("release")
/* 生命周期函数 */
onBeforeMount(async () => {
  getTemlates()
})
onMounted(async () => {
  if(route.query.namespace) {
    namespace.value = route.query.namespace
    setTimeout(async () => {
      await getList(true)
      if(route.query.name) {
        let row = all.value.find(n => n.metadata.name == route.query.name)
        if(row) 
          editOne(row)
        else  
          ElMessage.error({message: `未找到pipeline=${route.query.name}`})
      }
    }, 500)
  }
})
/* methods */
const getList = async (ifLoading) => {
  if(ifLoading) loading.value = true
  if(props.defaultTekton.cluster !== "" && namespace.value !== "") {
    let url = `/lizardcd/tekton/cluster/${props.defaultTekton.cluster}/namespace/${namespace.value}/pipelines`
    if(role.value !== 'admin') {
      labelSelector.value = [`project in (${localStorage.tenant})`]
    } else {
      labelSelector.value = []
    }
    if(inputLabels.value.length > 0) {
      labelSelector.value = _.uniq(labelSelector.value.concat(inputLabels.value))
    }
    labelSelector.value.push(`type=${pipelineType.value}`)
    if(labelSelector.value.length > 0) url += `?label_selector=${labelSelector.value.join(',')}`
    let response = await axios.get(url)
    all.value = _.sortBy(response.results, 'metadata.creationTimestamp').reverse()
    getPage(current.value)
  }
  if(ifLoading) loading.value = false
}
const getPage = async (page) => {
  let tmpList = all.value
  if(searchKey.value !== '') {
    tmpList = all.value.filter(n => n.metadata.name.includes(searchKey.value))
  }
  pageTotal.value = tmpList.length
  list.value = tmpList.slice((page-1)*pageSize.value, page*pageSize.value)
  getLastPipelineRun()
}
const getTemlates = async () => {
  let response = await axios.get(`/lizardcd/db/yaml_template?filter=type==tekton_pipeline&page=1&size=100&sort=update_at desc`)
  templateList.value = response.results.map(x => {
    x.variables = JSON.parse(x.variables)
    return x
  })
  response = await axios.get(`/lizardcd/db/yaml_template?filter=name==tekton_template_pipelinerun&page=1&size=1`)
  if(response.total > 0) {
    response = await axios.post(`/lizardcd/kubernetes/fetch/yaml`, {
      content: response.results[0].content,
      variables: {
        Annotations: {},
        Pipeline: "pipeline",
        Namespace: "tektoncd-default",
        Revision: "master"
      }
    })
    pipelineRunObject.value = yaml.load(response)
  }
}
const selectTemplate = (val) => {
  if(val) {
    form.value.content = val.content
    form.value.variables = val.variables
  }
  else {
    form.value.content = ""
  }
}
const editOne = async (row) => {
  show.value.yaml = true
  edit.value = true
  form.value.content = await axios.get(`/lizardcd/tekton/cluster/${props.defaultTekton.cluster}/namespace/${namespace.value}/pipelines/${row.metadata.name}/yaml`)
}
const deleteOne = async (row) => {
  await axios.delete(`/lizardcd/tekton/cluster/${props.defaultTekton.cluster}/namespace/${namespace.value}/pipelines/${row.metadata.name}`)
  ElMessage.success({message: `删除Pipeline成功`})
  setTimeout(async () => {
    await getList(true)
  }, 500)
}
const submitYaml = async () => {
  if(props.defaultTekton.cluster === "" || namespace.value == "") {
    ElMessage.warning({message: '请指定集群和命名空间'})
    return
  } 
  let params = Object.assign({}, form.value)
  delete params.templates
  let vars = {
    "Namespace": namespace.value,
    "Username": localStorage.username
  }
  for(let x of params.variables) {
    vars[x.key] = x.value
  }
  params.variables = vars
  try {
    if(edit.value === false) {
      await axios.post(`/lizardcd/tekton/cluster/${props.defaultTekton.cluster}/namespace/${namespace.value}/apply?kind=Pipeline`, params)
      ElMessage.success({message: `创建Pipeline成功`})
    }
    else {
      delete params.variables
      await axios.post(`/lizardcd/tekton/cluster/${props.defaultTekton.cluster}/namespace/${namespace.value}/apply?kind=Pipeline`, params)
        ElMessage.success({message: `更新Pipeline成功`})
    }
    show.value.yaml = false
    setTimeout(async () => {
      await getList(true)
    }, 500)
  } catch(e) {
    ElMessage.error({message: e})
  }
}
const start = async (row) => {
  if(!pipelineRunObject.value) {
    ElMessage.error({message: `Cannot find tekton_template_pipelinerun, please create it first`})
    return
  }
  pipelineRunObject.value.metadata = {
    annotations: Object.assign(row.metadata.annotations||{}, {
      'phecda.pipeline/triggered-by': localStorage.username
    }),
    labels: row.metadata.labels,
    generateName: row.metadata.name + "-",
    namespace: row.metadata.namespace,
  }
  pipelineRunObject.value.spec.params = row.spec.params?.map(x => {
    return {
      name: x.name,
      value: x.default || "",
    }
  })
  pipelineRunObject.value.spec.pipelineRef.name = row.metadata.name
  pipelineRunObject.value.spec.workspaces[0].volumeClaimTemplate.metadata.name = row.metadata.name
  startContent.value = yaml.dump(pipelineRunObject.value)
  formParams.value = _.cloneDeep(row.spec.params)
  for(let x of formParams.value||[]) {
    x.value = x.default || ""
    if(x.description) {
      x.description = x.description.trim().split("\n").join("<br>")
      if(x.description.startsWith("options:")) {
        x.options = x.description.slice("options:".length).split(",")
      }
    }
    delete x.type
  }
  show.value.start = true
  getLastPipelineRun()

}
const submitStart = async () => {
  if(startMode.value === "1") {
    if(formParams.value) {
      pipelineRunObject.value.spec.params = formParams.value.map(x => {
        delete x.description
        delete x.default
        delete x.options
        return x
      })
    }
    pipelineRunObject.value.spec.timeouts.pipeline = timeout.value
    startContent.value = yaml.dump(pipelineRunObject.value)
  }
  try {
    await axios.post(`/lizardcd/tekton/cluster/${props.defaultTekton.cluster}/namespace/${namespace.value}/apply?kind=PipelineRun`, {
      content: startContent.value
    })
    ElMessage.success({message: `启动Pipeline成功`})
    show.value.start = false
  } catch(e) {
    ElMessage.error({message: e})
  }
}
const submitStartBatch = async () => {
  try {
    await Promise.all(pipelineRunObjects.value.map(x => {
      x.metadata.annotations = Object.assign(x.metadata.annotations||{}, {
        'phecda.pipeline/triggered-by': localStorage.username
      })
      for(let [index, y] of x.spec.params.entries()) {
        let found = formParams.value.find(n => y.name === n.name)
        if(found) {
          x.spec.params[index] = _.cloneDeep(found)
          delete x.spec.params[index].description
          delete x.spec.params[index].default
          delete x.spec.params[index].type
        }
      }
      x.spec.timeouts.pipeline = timeout.value
      return axios.post(`/lizardcd/tekton/cluster/${props.defaultTekton.cluster}/namespace/${namespace.value}/apply?kind=PipelineRun`, {
        content: yaml.dump(x)
      })
    }))
    ElMessage.success({message: `批量启动Pipeline成功`})
    show.value.start = false
  } catch(e) {
    ElMessage.error({message: e})
  }
}
const handleMore = async (command) => {
  if(namespace.value === '') {
    ElMessage.warning({message: '请选择命名空间'})
    return
  }
  switch(command.action) {
    case "create": {
      show.value.yaml = true
      edit.value = false
      form.value = {content:'',variables:[]}
      break
    }
    case "createVisual": {
      router.push({
        path: '/ci/tekton/create',
        query: {
          namespace: namespace.value
        }
      })
      break
    }
    case "deleteBatch": {
      if(selected.value.length === 0) {
        ElMessage.warning({message: '请勾选记录'})
        return
      }
      await ElMessageBox.confirm('确定删除？','警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(async () => {
        await Promise.all(selected.value.map(x => {
          return axios.delete(`/lizardcd/tekton/cluster/${props.defaultTekton.cluster}/namespace/${namespace.value}/pipelines/${x.metadata.name}`)
        }))
        ElMessage.success({message: `删除Pipeline成功`})
        setTimeout(async () => {
          await getList(true)
        }, 500)
      }).catch(() =>{})
      break
    }
    case "startBatch": {
      if(selected.value.length === 0) {
        ElMessage.warning({message: '请勾选记录'})
        return
      }
      for(let row of selected.value) {
        let object = _.cloneDeep(pipelineRunObject.value)
        object.metadata = {
          annotations: row.metadata.annotations,
          labels: row.metadata.labels,
          generateName: row.metadata.name + "-",
          namespace: row.metadata.namespace,
        }
        object.spec.params = row.spec.params?.map(x => {
          return {
            name: x.name,
            value: x.default || "",
          }
        })
        object.spec.pipelineRef.name = row.metadata.name
        object.spec.workspaces[0].volumeClaimTemplate.metadata.name = row.metadata.name
        pipelineRunObjects.value.push(object)
      }
      formParams.value = _.cloneDeep(selected.value[0].spec.params.filter(n => !n.hasOwnProperty("default")))
      batchStart.value = true
      show.value.start = true
    }
  }
}
const handleSizeChange = async (size) => {
  pageSize.value = size
  await getList(true)
}
const addVar = (index) => {
  form.value.variables.splice(index+1, 0, {key:"", value:""})
}
const removeVar = (index) => {
  form.value.variables.splice(index, 1)
}
const select = (val) => {
  selected.value = val
}
const getLastPipelineRun = async () => {
  let names = list.value.map(x => x.metadata.name)
  lastRun.value = await axios.get(`/tekton-pipelines/lizardcd/${namespace.value}/pipelineRuns/lasts?pipelineRefs=${names.join(",")}`)
  for(let values of Object.values(lastRun.value)) {
    values.startTime = moment(values.startTime).format('YYYY-MM-DD HH:mm:ss')
    if(values.endTime)
      values.endTime = moment(values.endTime).format('YYYY-MM-DD HH:mm:ss')
  }
}
const gotoLast = (row) => {
  window.open(`/ci/tekton/pipelinerun/${row.name}?cluster=${props.defaultTekton.cluster}&namespace=${namespace.value}`, '_blank')
}
const gotoHistory = (row) => {
  window.open(`/ci/tekton?tab=PipelineRun&namespace=${namespace.value}&pipelineName=${row.metadata.name}`, '_blank')
}
const setNamespace = (ns) => {
  namespace.value = ns
  getList(true)
}
defineExpose({
  setNamespace
})
</script>
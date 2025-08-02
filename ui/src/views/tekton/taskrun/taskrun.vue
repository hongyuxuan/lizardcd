<template>
<div class="box box-item">
  <div class="box-body" style="padding-top:20px;padding-bottom:0">
    <el-alert type="info" show-icon style="margin-bottom:15px">
      <template #title>
        关于 TaskRun 配置参考：<el-link href="https://tekton.dev/docs/pipelines/taskruns/" underline="never" type="primary" target="_blank">TaskRuns</el-link>
      </template>
    </el-alert>
    <el-row>
      <el-col :span="18">
        <el-button-group>
          <el-button :icon="Refresh" size="large" @click="getList(true)" />
          <el-select v-model="namespace" placeholder="请选择命名空间" clearable filterable @change="getList(true)" style="width:250px;" size="large">
            <el-option v-for="(item) in props.namespaceList" :key="item" :label="item" :value="item" />
          </el-select>
          <el-input v-model="searchKey" placeholder="输入名称进行搜索" size="large" :prefix-icon="Search" @change="getPage(1)" clearable style="width:250px;" />
        </el-button-group>
      </el-col>
      <el-col :span="6">
        <el-dropdown @command="handleMore" class="pull-right">
          <el-button size="large">更多操作<el-icon class="el-icon--right"><arrow-down /></el-icon></el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item :command="{action:'deleteBatch'}">批量删除</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-button class="pull-right" size="large" type="primary" @click="show.yaml=true;edit=false;form={content:'',variables:[]}">+ 创建TaskRun</el-button>
      </el-col>
    </el-row>
    <el-table 
      :data="list" 
      v-loading="loading"
      element-loading-text="奋力加载中..."
      @selection-change="select"
      class="line-height40" 
      style="width:100%;margin-top:10px;min-height:150px">
      <el-table-column type="selection" width="45" />
      <el-table-column label="名称" min-width="200">
        <template #default="scope">
          <el-link underline="never" :href="`/ci/tekton/taskrun/${scope.row.metadata.name}?cluster=${defaultTekton.cluster}&namespace=${namespace}`" target="_blank">{{ scope.row.metadata.name }}</el-link>
        </template>
      </el-table-column>
      <el-table-column label="关联Task" min-width="160">
        <template #default="scope">
          {{ scope.row.spec.taskRef?.namespace || scope.row.metadata.namespace }}/{{ scope.row.spec.taskRef?.name }}
        </template>
      </el-table-column>
      <el-table-column label="执行状态" width="120">
        <template #default="scope">
          <el-progress v-if="scope.row.status?.conditions[0].reason=='Succeeded'" :percentage="100" color="#5cb87a" :show-text="false" />
          <el-tooltip v-else-if="scope.row.status?.conditions[0].reason=='TaskRunCancelled'" effect="dark" placement="top" content="TaskRunCancelled">
            <el-progress :percentage="100" color="#bfbbbb" :show-text="false" />
          </el-tooltip>
          <el-progress v-else-if="scope.row.status?.conditions[0].status=='Unknown'" :percentage="50" color="#e6a23c" :show-text="false" />
          <el-tooltip v-else effect="dark" placement="top" :content="scope.row.status.conditions[0].message">
            <el-progress :percentage="100" color="#f56c6c" :show-text="false" />
          </el-tooltip>
        </template>
      </el-table-column>
      <el-table-column label="开始时间" width="170">
        <template #default="scope">
          {{ scope.row.status?.startTime?moment(scope.row.status?.startTime).format('YYYY-MM-DD HH:mm:ss'):"" }}
        </template>
      </el-table-column>
      <el-table-column label="结束时间" width="170">
        <template #default="scope">
          {{ scope.row.status?.completionTime?moment(scope.row.status.completionTime).format('YYYY-MM-DD HH:mm:ss'):"" }}
        </template>
      </el-table-column>
      <el-table-column label="耗时" width="120">
        <template #default="scope">
          <span v-if="scope.row.status?.completionTime&&scope.row.status?.startTime">
            {{ moment.duration(moment(scope.row.status.completionTime)-moment(scope.row.status?.startTime)).humanize() }}
          </span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="170">
        <template #default="scope">
          <el-button icon="EditPen" circle @click="editOne(scope.row)"></el-button>
          <el-popconfirm title="确定重跑?" confirm-button-text="确认" cancel-button-text="取消" @confirm="reRun(scope.row)">
            <template #reference>
              <el-button icon="RefreshRight" circle />
            </template>
          </el-popconfirm>
          <el-popconfirm title="确定终止?" confirm-button-text="确认" cancel-button-text="取消" @confirm="cancel(scope.row)">
            <template #reference>
              <el-button circle><font-awesome-icon icon="ban" /></el-button>
            </template>
          </el-popconfirm>
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
    <h4>{{ edit === true ? '编辑YAML' : '创建TaskRun' }}</h4>
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
</template>
<script setup>
import { Search,Refresh } from '@element-plus/icons-vue'
import { onBeforeMount, ref } from 'vue'
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
const selected = ref([])
const edit = ref(false)
const form = ref({content:'',variables:[],templates:{}})
const templateList = ref([])
const wrapLine = ref(true)
/* 生命周期函数 */
onBeforeMount(async () => {
  getTemlates()
})
/* methods */
const getList = async (ifLoading) => {
  if(ifLoading) loading.value = true
  if(props.defaultTekton.cluster !== "" && namespace.value !== "") {
    let response = await axios.get(`/lizardcd/tekton/cluster/${props.defaultTekton.cluster}/namespace/${namespace.value}/taskruns`)
    for(let x of response.results) {
      if(x.spec.taskRef) { // 引用其它ns：taskRef
        x.spec.taskRef.name ||= x.spec.taskRef.params?.find(n => n.name === 'name')?.value
        x.spec.taskRef.namespace = x.spec.taskRef.params?.find(n => n.name === 'namespace')?.value
      } 
    }
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
}
const getTemlates = async () => {
  let response = await axios.get(`/lizardcd/db/yaml_template?filter=type==tekton_taskrun&page=1&size=100&sort=update_at desc`)
  templateList.value = response.results.map(x => {
    x.variables = JSON.parse(x.variables)
    return x
  })
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
  form.value.content = await axios.get(`/lizardcd/tekton/cluster/${props.defaultTekton.cluster}/namespace/${namespace.value}/taskruns/${row.metadata.name}/yaml`)
}
const deleteOne = async (row) => {
  await axios.delete(`/lizardcd/tekton/cluster/${props.defaultTekton.cluster}/namespace/${namespace.value}/taskruns/${row.metadata.name}`)
  ElMessage.success({message: `删除TaskRun成功`})
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
      await axios.post(`/lizardcd/tekton/cluster/${props.defaultTekton.cluster}/namespace/${namespace.value}/apply?kind=TaskRun`, params)
      ElMessage.success({message: `创建TaskRun成功`})
    }
    else {
      delete params.variables
      await axios.post(`/lizardcd/tekton/cluster/${props.defaultTekton.cluster}/namespace/${namespace.value}/apply?kind=TaskRun`, params)
      ElMessage.success({message: `更新TaskRun成功`})
    }
    show.value.yaml = false
    setTimeout(async () => {
      await getList(true)
    }, 500)
  } catch(e) {
    ElMessage.error({message: e})
  }
}
const handleSizeChange = async (size) => {
  pageSize.value = size
  await getList(true)
}
const reRun = async (row) => {
  let params = _.cloneDeep(row)
  params.metadata.generateName ||= params.metadata.labels["tekton.dev/task"] + "-"
  delete params.metadata.name
  delete params.metadata.resourceVersion
  delete params.metadata.uid
  delete params.metadata.generation
  delete params.spec.status
  delete params.status
  if(params.metadata.annotations)
    delete params.metadata.annotations["pipeline.tekton.dev/affinity-assistant"]
  if(params.spec.taskRef) {
    delete params.spec.taskRef.namespace
    if(params.spec.taskRef.resolver) {
      delete params.spec.taskRef.name
    }
  }
  try {
    await axios.post(`/lizardcd/tekton/cluster/${props.defaultTekton.cluster}/namespace/${namespace.value}/apply?kind=TaskRun`, {content: yaml.dump(params)})
    ElMessage.success({message: `重新执行TaskRun成功`})
    setTimeout(async () => {
      await getList(true)
    }, 500)
  } catch(e) {
    ElMessage.error({message: e})
  }
}
const handleMore = async (command) => {
  if(selected.value.length === 0) {
    ElMessage.warning({message: '请勾选TaskRun'})
    return
  }
  switch(command.action) {
    case "deleteBatch": {
      await ElMessageBox.confirm('确定删除？','警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(async () => {
        await Promise.all(selected.value.map(x => {
          return axios.delete(`/lizardcd/tekton/cluster/${props.defaultTekton.cluster}/namespace/${namespace.value}/taskruns/${x.metadata.name}`)
        }))
        ElMessage.success({message: `删除TaskRun成功`})
        setTimeout(async () => {
          await getList(true)
        }, 500)
      }).catch(() =>{})
      break
    }
  }
}
const select = (val) => {
  selected.value = val
}
const addVar = (index) => {
  form.value.variables.splice(index+1, 0, {key:"", value:""})
}
const removeVar = (index) => {
  form.value.variables.splice(index, 1)
}
const cancel = async (row) => {
  try {
    await axios.patch(`/lizardcd/tekton/cluster/${props.defaultTekton.cluster}/namespace/${namespace.value}/taskruns/${row.metadata.name}`, [
      {
        "op": "replace",
        "path": "/spec/status",
        "value": "TaskRunCancelled"
      }
    ])
    ElMessage.success({message: `终止TaskRun成功`})
    setTimeout(async () => {
      await getList(current.value)
    }, 1000)
  } catch(e) {
    ElMessage.error({message: e})
  }
}
</script>
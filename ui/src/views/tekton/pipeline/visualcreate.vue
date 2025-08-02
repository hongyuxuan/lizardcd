<template>
<el-breadcrumb :separator-icon="ArrowRight">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>持续集成</el-breadcrumb-item>
  <el-breadcrumb-item :to="{path:'/ci/tekton'}">流水线</el-breadcrumb-item>
  <el-breadcrumb-item>创建流水线</el-breadcrumb-item>
</el-breadcrumb>
<el-card>
  <template #header>
    <div class="card-header">
      <span class="card-header-text">创建流水线</span>
    </div>
  </template>
  <el-row :gutter="15">
    <el-col :span="12">
      <el-form ref="add" :model="form" label-width="100px">
        <el-form-item label="流水线名称" prop="name">
          <el-input v-model="form.metadata.name" size="large" />
        </el-form-item>
        <el-form-item label="标签" prop="labels">
          <el-row v-for="(item,i) in labels" :key="i" style="margin-bottom:5px;width:100%">
            <el-button-group>
              <el-input v-model="item.key" size="large" clearable style="width:200px;margin-right:5px" />
              <el-input v-model="item.value" size="large" clearable style="width:200px;" />
              <el-button circle :icon="Delete" size="large" style="float:right" @click="removeLabel(i)" />
            </el-button-group>
          </el-row>
          <el-button circle icon="Check" @click="okLabel()" v-if="labels.length>0" />
          <el-button circle icon="Plus" @click="addLabel()" />
        </el-form-item>
        <el-form-item label="注解" prop="annotations">
          <el-row v-for="(item,i) in annotations" :key="i" style="margin-bottom:5px;width:100%">
            <el-button-group>
              <el-input v-model="item.key" size="large" clearable style="width:200px;margin-right:5px" />
              <el-input v-model="item.value" size="large" clearable style="width:200px;" />
              <el-button circle :icon="Delete" size="large" style="float:right" @click="removeAnnotation(i)" />
            </el-button-group>
          </el-row>
          <el-button circle icon="Check" @click="okAnnotation()" v-if="annotations.length>0" />
          <el-button circle icon="Plus" @click="addAnnotation()" />
        </el-form-item>
        <el-form-item label="流水线参数">
          <table class="table table-bordered">
            <thead><tr><th>参数名</th><th>参数类型</th><th>默认值</th></tr></thead>
            <tbody>
            <tr v-for="(item,index) in form.spec.params" :key="index" >
              <td ><el-input v-model="item.name" size="large" /></td>
              <td><el-input v-model="item.type" size="large" /></td>
              <td width="200"><el-input v-model="item.default" size="large" /></td>
              <td width="80">
                <el-button-group>
                  <el-button icon="Plus" circle @click="addParam(index)"></el-button>
                  <el-button icon="Close" circle @click="removeParam(index)"></el-button>
                </el-button-group>
              </td>
            </tr>
            </tbody>
          </table>
        </el-form-item>
        <el-form-item label="任务步骤">
          <div style="min-height: 100px; max-width: 600px">
            <el-steps direction="vertical" :active="form.spec.tasks.length">
              <el-step v-for="(item,i) in form.spec.tasks" :key="i">
                <template #title>
                  <el-link v-if="!item.name" underline="never" type="primary" @click="currentTaskIndex=i;show.task=true">请选择任务</el-link>
                  <el-card v-else style="width:400px" shadow="always" @click="currentTaskIndex=i;show.task=true" class="pointer nobox">
                    <span style="font-size:16px;padding-bottom:10px">{{ item.name }}</span>
                  </el-card>
                </template>
                <template #description>
                  <el-button circle :icon="Plus" @click="addTask(i)" />
                  <el-button circle :icon="Minus" @click="removeTask(i)" />
                </template>
              </el-step>
            </el-steps>
          </div>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onSubmit">创建</el-button>
        </el-form-item>
      </el-form>
    </el-col>
    <el-col :span="12">
      <v-ace-editor
        v-model:value="yamlContent"
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
    </el-col>
  </el-row>
</el-card>
<el-drawer v-model="show.task" direction="rtl" size="700px">
  <template #header>
    <h4>设置Task</h4>
  </template>
  <template #default>
    <el-form ref="task" label-width="120px">
      <el-form-item label="任务名称">
        <el-input size="large" v-model="formTask.name" />
      </el-form-item>
      <el-form-item label="选择Namespace">
        <el-select 
          v-model="formTask.namespace" 
          placeholder="选择Namespace" 
          clearable 
          filterable
          size="large"
          @change="getTasks"
          style="width:100%">
          <el-option v-for="item in namespaceList" :key="item" :label="item" :value="item" />
        </el-select>
      </el-form-item>
      <el-form-item label="选择Task">
        <el-select 
          v-model="formTask.task" 
          placeholder="请选择Task" 
          value-key="metadata.name" 
          clearable 
          filterable
          size="large"
          @change="selectTask"
          style="width:100%">
          <el-option v-for="item in taskList" :key="item.metadata.name" :label="item.metadata.name" :value="item" />
        </el-select>
      </el-form-item>
      <el-form-item label="YAML">
        <v-ace-editor
          v-model:value="formTask.content"
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
      <el-form-item label="任务参数">
        <table class="table table-bordered">
          <thead><tr><th>参数名</th><th>参数值</th></tr></thead>
          <tbody>
          <tr v-for="(item,index) in formTask.params||[]" :key="index" >
            <td ><el-input v-model="item.name" size="large" /></td>
            <td><el-input v-model="item.value" size="large" /></td>
          </tr>
          </tbody>
        </table>
      </el-form-item>
      <el-form-item label="Workspace">
        <table class="table table-bordered">
          <thead><tr><th>名称</th><th>值</th></tr></thead>
          <tbody>
          <tr v-for="(item,index) in formTask.workspaces||[]" :key="index" >
            <td ><el-input v-model="item.name" size="large" /></td>
            <td><el-input v-model="item.workspace" size="large" /></td>
          </tr>
          </tbody>
        </table>
      </el-form-item>
    </el-form>
  </template>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show.task=false">取消</el-button>
      <el-button type="primary" @click="okTask">提交</el-button>
    </div>
  </template>
</el-drawer>
</template>

<script setup>
import { ArrowRight, Minus, Plus, Delete } from '@element-plus/icons-vue'
import { onBeforeMount, ref, computed,  } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { axios } from '/src/assets/util/axios'
import { ElMessage } from 'element-plus'
import _ from 'lodash'
import yaml from 'js-yaml'
/* 引入v-ace-editor */
import { VAceEditor } from 'vue3-ace-editor'
import 'ace-builds/src-noconflict/mode-yaml'
import 'ace-builds/src-noconflict/theme-chrome'
import 'ace-builds/src-noconflict/ext-language_tools'
/* 变量定义 */
const route = useRoute()
const router = useRouter()
const namespaceList = ref([])
const form = ref({
  "apiVersion": "tekton.dev/v1",
  "kind": "Pipeline",
  "metadata": {
    "namespace": route.query.namespace,
  },
  "spec": {
    "params": [{}],
    "tasks": [{}],
    "workspaces": [
      { name: "shared-workspace"}
    ]
  },
})
const yamlContent = computed(() => {
  return yaml.dump(form.value)
})
const labels = ref([])
const annotations = ref([])
const show = ref({
  task: false
})
const defaultTekton = ref({})
const formTask = ref({
  content: ''
})
const currentTaskIndex = ref({})
const taskList = ref([])
/* 生命周期函数 */
onBeforeMount(async () => {
  getClusterList()
})
/* methods */
const getClusterList = async () => {
  let response = await axios.get(`/lizardcd/db/settings?filter=setting_key==default_tekton`)
  if(response.total > 0) {
    defaultTekton.value = JSON.parse(response.results[0].setting_value)
    response = await axios.get(`/lizardcd/server/clusters`)
    if(response.hasOwnProperty(defaultTekton.value.cluster)) {
      namespaceList.value = response[defaultTekton.value.cluster].filter(n => n.startsWith("tektoncd-"))
    }
  }
}
const getTasks = async (val) => {
  if(defaultTekton.value.cluster !== "" && val !== "") {
    let response = await axios.get(`/lizardcd/tekton/cluster/${defaultTekton.value.cluster}/namespace/${val}/tasks`)
    taskList.value = _.sortBy(response.results, 'metadata.creationTimestamp').reverse()
  }
}
const selectTask = async (val) => {
  formTask.value.content = await axios.get(`/lizardcd/tekton/cluster/${defaultTekton.value.cluster}/namespace/${formTask.value.namespace}/tasks/${val.metadata.name}/yaml`)
  formTask.value.params = val.spec.params
  formTask.value.workspaces = val.spec.workspaces
}
const okTask = () => {
  form.value.spec.tasks[currentTaskIndex.value].name = formTask.value.name
  form.value.spec.tasks[currentTaskIndex.value].taskRef = {
    "kind": "Task",
    "params": [
      { name: 'kind', value: 'task'},
      { name: 'namespace', value: formTask.value.namespace},
      { name: 'name', value: formTask.value.task.metadata.name},
    ],
    "resolver": "cluster"
  }
  if(formTask.value.params) {
    form.value.spec.tasks[currentTaskIndex.value].params = formTask.value.params.filter(n => n.value).map(x => {
      return {
        name: x.name,
        value: x.value
      }
    })
  }
  if(formTask.value.workspaces) {
    form.value.spec.tasks[currentTaskIndex.value].workspaces = formTask.value.workspaces.filter(n => n.workspace).map(x => {
      return {
        name: x.name,
        workspace: x.workspace
      }
    })
  }
  if(currentTaskIndex.value > 0) {
    let lastTask = form.value.spec.tasks[currentTaskIndex.value-1]
    form.value.spec.tasks[currentTaskIndex.value].runAfter = lastTask.name
  }
  show.value.task = false
}
const addLabel = () => {
  labels.value.push({key: '', value: ''})
}
const removeLabel = (index) => {
  if(labels.value[index]?.key)
    delete form.value.metadata.labels[labels.value[index].key]
  labels.value.splice(index, 1)
}
const okLabel = () => {
  if(!form.value.metadata.labels) form.value.metadata.labels = {}
  for(let x of labels.value) {
    form.value.metadata.labels[x.key] = x.value
  }
}
const addAnnotation = () => {
  annotations.value.push({key: '', value: ''})
}
const removeAnnotation = (index) => {
  if(annotations.value[index]?.key)
    delete form.value.metadata.annotations[annotations.value[index].key]
  annotations.value.splice(index, 1)
}
const okAnnotation = () => {
  if(!form.value.metadata.annotations) form.value.metadata.annotations = {}
  for(let x of annotations.value) {
    form.value.metadata.annotations[x.key] = x.value
  }
}
const addParam = (index) => {
  form.value.spec.params.splice(index+1, 0, {name:"", type:""})
}
const removeParam = (index) => {
  form.value.spec.params.splice(index, 1)
}
const addTask = (index) => {
  form.value.spec.tasks.splice(index+1, 0, {})
}
const removeTask = (index) => {
  form.value.spec.tasks.splice(index, 1)
}
const onSubmit = async () => {
  try {
    await axios.post(`/lizardcd/tekton/cluster/${defaultTekton.value.cluster}/namespace/${route.query.namespace}/apply?kind=Pipeline`, {
      content: yamlContent.value
    })
    ElMessage.success({message: `创建Pipeline成功`})
    setTimeout(async () => {
      router.push({
        path: '/ci/tekton',
        query: {
          tab: 3
        }
      })
    }, 1000)
  }catch(e) {
    ElMessage.error({message: e})
  }
}
</script>
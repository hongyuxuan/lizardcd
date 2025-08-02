<template>
<el-form :model="form" label-width="120px">
  <el-form-item>
    <template #label><el-checkbox v-model="disabled.trigger_type">触发类型</el-checkbox></template>
    <el-radio-group v-model="form.trigger_type" :disabled="!disabled.trigger_type">
      <el-radio value="pipelinerun">Tekton PipelineRun</el-radio>
      <el-radio value="endpoint">Tekton Endpoint</el-radio>
      <el-radio value="自定义">自定义</el-radio>
    </el-radio-group>
  </el-form-item>
  <el-form-item v-if="form.trigger_type==='endpoint'">
    <template #label><el-checkbox v-model="disabled.cluster">选择EventListener</el-checkbox></template>
    <el-select 
      v-model="cluster" 
      placeholder="请选择集群" 
      clearable 
      filterable 
      :disabled="!disabled.cluster" 
      style="width:200px;margin:0 8px 8px 0" 
      size="large">
      <el-option v-for="(v,k) in props.clusterList" :key="k" :label="k" :value="k" />
    </el-select>
    <el-select 
      v-model="namespace" 
      placeholder="请选择命名空间" 
      clearable 
      filterable 
      @change="getEventListener" 
      :disabled="!disabled.cluster" 
      style="width:200px;margin:0 8px 8px 0" 
      size="large">
      <el-option v-for="(item) in clusterList[cluster]" :key="item" :label="item" :value="item" />
    </el-select>
    <el-select 
      v-model="eventListener" 
      placeholder="请选择EventListener" 
      clearable 
      filterable 
      value-key="metadata.name" 
      @change="selectEventListener"
      :disabled="!disabled.cluster"
      style="width:100%" 
      size="large">
      <el-option v-for="(item,index) in eventListenerList" :key="index" :label="item.metadata.name" :value="item" />
    </el-select>
  </el-form-item>
  <el-form-item label="触发地址" v-if="form.trigger_type!=='pipelinerun'">
    <el-input v-model="form.trigger_endpoint" size="large" :disabled="!disabled.cluster" clearable />
  </el-form-item>
  <el-form-item label="提交body" v-if="form.trigger_type!=='pipelinerun'">
    <v-ace-editor
      v-model:value="form.trigger_body"
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
        maxLines: 100,
      }" />
  </el-form-item>
  <el-form-item prop="ref_pattern">
    <template #label><el-checkbox v-model="disabled.ref_pattern">分支白名单</el-checkbox></template>
    <el-input v-model="form.ref_pattern" size="large" :disabled="!disabled.ref_pattern" clearable />
  </el-form-item>
  <el-form-item prop="trigger_event">
    <template #label><el-checkbox v-model="disabled.trigger_event">触发事件</el-checkbox></template>
    <el-checkbox-group v-model="form.trigger_event" :disabled="!disabled.trigger_event">
      <el-checkbox label="push" value="push" />
      <el-checkbox label="merge_request" value="merge_request" />
    </el-checkbox-group>
  </el-form-item>
  <el-form-item prop="uniq_instance">
    <template #label><el-checkbox v-model="disabled.uniq_instance">唯一实例</el-checkbox></template>
    <el-checkbox label="勾选后，后运行的流水线将终止前面的运行实例" v-model="form.uniq_instance" :disabled="!disabled.uniq_instance" />
  </el-form-item>
  <el-form-item label="所属租户" prop="tenant">
    <template #label><el-checkbox v-model="disabled.tenant">所属租户</el-checkbox></template>
    <el-select 
      v-model="form.tenant" 
      placeholder="请选择租户" 
      clearable 
      filterable 
      :disabled="userInfo.role!=='admin'||!disabled.tenant"
      style="width:100%;" 
      size="large">
      <el-option v-for="(item,index) in props.tenantList" :key="index" :label="item" :value="item" />
    </el-select>
  </el-form-item>
</el-form>
</template>

<script setup>
import { onBeforeMount, ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { useStore } from 'vuex'
import { axios } from '/src/assets/util/axios.js'
import moment from 'moment'
/* 引入v-ace-editor */
import { VAceEditor } from 'vue3-ace-editor'
import 'ace-builds/src-noconflict/mode-yaml'
import 'ace-builds/src-noconflict/theme-chrome'
import 'ace-builds/src-noconflict/ext-language_tools'
/* 变量定义 */
const props = defineProps({
  clusterList: { type: Object }, 
  tenantList: { type: Array },
})
const store = useStore()
const userInfo = computed(() => {
  return store.state.userInfo
})
const disabled = ref({
  trigger_type: false,
  cluster: false,
  ref_pattern: false,
  trigger_event: false,
  uniq_instance: false,
  tenant: false,
})
const form = ref({
  trigger_type: 'pipelinerun',
  trigger_body: ''
})
const cluster = ref("")
const namespace = ref("")
const eventListener = ref()
const eventListenerList = ref([])
/* methods */
const getEventListener = async () => {
  let response = await axios.get(`/lizardcd/tekton/cluster/${cluster.value}/namespace/${namespace.value}/eventlisteners`)
  eventListenerList.value = response.results
}
const selectEventListener = (val) => {
  if(!val.metadata.annotations?.endpoint)
    ElMessage.warning('EventListener未找到endpoint注解')
  else
    form.value.trigger_endpoint = val.metadata.annotations.endpoint
}
const submit = async (rows) => {
  let params = {
    update_at: moment()
  }
  for(let [k,v] of Object.entries(disabled.value)) {
    if(v === true) {
      params[k] = form.value[k]
      if(k === 'trigger_event') params[k] = params[k].join(',')
    }
  }
  await Promise.all(rows.map(x => {
    let body = Object.assign(x, params)
    delete body.application
    return axios.put(`/lizardcd/db/ci_trigger/${x.id}`, {body})
  }))
}
defineExpose({
  submit
})
</script>
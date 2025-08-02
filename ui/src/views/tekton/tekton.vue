<template>
<el-breadcrumb :separator-icon="ArrowRight">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>持续集成</el-breadcrumb-item>
  <el-breadcrumb-item>流水线</el-breadcrumb-item>
</el-breadcrumb>
<div class="box box-solid">
  <div class="box-header page-intro">
    <p>
      Tekton 是一个云原生的持续集成与持续构建平台，它在 Kubernetes 上以 CRD 的方式定义了构建所需的各种类型资源。参见 
      <el-link href="http://tekton.dev" underline="never" type="primary" target="_blank">Tekton</el-link><br>
      流水线页面为您提供 Tekton 的常见资源类型的增删改查等管理操作。<br>
      注意，本页面仅显示 Tekton 集群中以 <b>tektoncd-</b> 开头的命名空间，请将您的 CR 资源创建到以 <b>tektoncd-</b> 开头的命名空间中。
    </p>
  </div>
  <div class="box-body" style="padding:0">
    <el-tabs v-model="activeName" type="border-card" class="qingcloud-tab">
      <el-tab-pane name="Task" label="Task">
        <keep-alive>
          <task ref="refTask" :defaultTekton="defaultTekton" :namespaceList="namespaceList" />
        </keep-alive>
      </el-tab-pane>
      <el-tab-pane name="TaskRun" label="TaskRun">
        <keep-alive>
          <taskrun :defaultTekton="defaultTekton" :namespaceList="namespaceList" v-if="tektonSource==='crd'" />
          <taskrun_cloudevent ref="refTaskRun" :defaultTekton="defaultTekton" :namespaceList="namespaceList" v-else-if="tektonSource==='cloudevent'" />
        </keep-alive>
      </el-tab-pane>
      <el-tab-pane name="Pipeline" label="Pipeline">
        <keep-alive>
          <pipeline ref="refPipeline" :defaultTekton="defaultTekton" :namespaceList="namespaceList" />
        </keep-alive>
      </el-tab-pane>
      <el-tab-pane name="PipelineRun" label="PipelineRun">
        <keep-alive>
          <pipelinerun :defaultTekton="defaultTekton" :namespaceList="namespaceList" v-if="tektonSource==='crd'" />
          <pipelinerun_cloudevent ref="refPipelineRun" :defaultTekton="defaultTekton" :namespaceList="namespaceList" v-else-if="tektonSource==='cloudevent'" />
        </keep-alive>
      </el-tab-pane>
      <el-tab-pane name="TriggerBinding" label="TriggerBinding">
        <keep-alive>
          <triggerbinding :defaultTekton="defaultTekton" :namespaceList="namespaceList" />
        </keep-alive>
      </el-tab-pane>
      <el-tab-pane name="TriggerTemplate" label="TriggerTemplate">
        <keep-alive>
          <triggertemplate :defaultTekton="defaultTekton" :namespaceList="namespaceList" />
        </keep-alive>
      </el-tab-pane>
      <el-tab-pane name="EventListener" label="EventListener">
        <keep-alive>
          <eventlistener :defaultTekton="defaultTekton" :namespaceList="namespaceList" />
        </keep-alive>
      </el-tab-pane>
      <el-tab-pane name="ApprovalTask" label="ApprovalTask" v-if="userInfo.role==='admin'&&tektonSource==='crd'">
        <keep-alive>
          <approvaltask :defaultTekton="defaultTekton" :namespaceList="namespaceList" />
        </keep-alive>
      </el-tab-pane>
    </el-tabs>
  </div>
</div>
</template>

<script setup>
import { ArrowRight } from '@element-plus/icons-vue'
import { onBeforeMount, ref, computed, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { useStore } from 'vuex'
import task from './task.vue'
import pipeline from './pipeline/pipeline.vue'
import pipelinerun from './pipelinerun/pipelinerun.vue'
import pipelinerun_cloudevent from './pipelinerun/pipelinerun_cloudevent.vue'
import triggerbinding from './triggerbinding.vue'
import triggertemplate from './triggertemplate.vue'
import eventlistener from './eventlistener.vue'
import taskrun from './taskrun/taskrun.vue'
import taskrun_cloudevent from './taskrun/taskrun_cloudevent.vue'
import approvaltask from './approvaltask.vue'
import { axios } from '/src/assets/util/axios'
/* 变量定义 */
const route = useRoute()
const store = useStore()
const userInfo = computed(() => {
  return store.state.userInfo
})
const tektonSource = ref({})
const activeName = ref("Task")
const defaultTekton = ref({})
const namespaceList = ref([])
const refTask = ref(null)
const refTaskRun = ref(null)
const refPipeline = ref(null)
const refPipelineRun = ref(null)
/* 生命周期函数 */
onBeforeMount(async () => {
  // get defaultTekton
  let response = await axios.get(`/lizardcd/db/settings?filter=setting_key==default_tekton`)
  if(response.total > 0)
    defaultTekton.value = JSON.parse(response.results[0].setting_value)

  // get namespaceList
  response = await axios.get(`/lizardcd/server/clusters`)
  if(response.hasOwnProperty(defaultTekton.value.cluster)) {
    namespaceList.value = response[defaultTekton.value.cluster].filter(n => n.startsWith("tektoncd-"))
  }
  
  if(route.query.tab) {
    activeName.value = route.query.tab
  }

  response = await axios.get(`/lizardcd/db/settings?filter=setting_key==tekton_source`)
  if(response.total > 0) {
    tektonSource.value = response.results[0].setting_value
  }

  if(namespaceList.value.length > 0) {
    refTask.value.setNamespace(namespaceList.value[0])
    if(!route.query.namespace) {
      refPipeline.value.setNamespace(namespaceList.value[0])
    }
    nextTick(() => {
      if(tektonSource.value === 'cloudevent') {
        refTaskRun.value.setNamespace(namespaceList.value[0])
        // refPipelineRun.value.setNamespace(namespaceList.value[0])
      }
    })
  }
})
</script>
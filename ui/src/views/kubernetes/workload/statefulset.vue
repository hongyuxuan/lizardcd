<template>
<el-breadcrumb :separator-icon="ArrowRight" style="min-width:1366px">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>工作负载</el-breadcrumb-item>
  <el-breadcrumb-item :to="{ path: '/kubernetes/workload', query: {cluster: route.query.cluster, namespace: route.query.namespace, tab: 'statefulsets'} }">有状态副本集</el-breadcrumb-item>
  <el-breadcrumb-item>{{ statefulsetInfo.metadata?.name }}</el-breadcrumb-item>
</el-breadcrumb>
<el-row :gutter="15" style="min-width:1366px">
  <el-col :span="6">
    <el-card>
      <template #header>
        <div class="card-header">
          <span class="card-header-text"><b>{{ statefulsetInfo.metadata?.name }}</b></span>
            <el-dropdown @command="handleCommand" class="pull-right" style="top:3px">
              <el-link underline="never">
                更多操作
                <el-icon class="el-icon--right">
                  <arrow-down />
                </el-icon>
              </el-link>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="restart">重启</el-dropdown-item>
                  <el-dropdown-item command="yaml">编辑YAML</el-dropdown-item>
                  <el-dropdown-item command="delete">删除</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
        </div>
      </template>
      <el-descriptions :column="1" border class="no-color" :label-width="120">
        <el-descriptions-item label="集群">{{ route.query.cluster }}</el-descriptions-item>
        <el-descriptions-item label="命名空间">{{ route.query.namespace }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ statefulsetInfo.metadata?.creationTimestamp }}</el-descriptions-item>
      </el-descriptions>
    </el-card>
  </el-col>
  <el-col :span="18">
    <el-menu :default-active="activeIndex" mode="horizontal" @select="handleSelect">
      <el-menu-item index="status">资源状态</el-menu-item>
      <el-menu-item index="labels">标签</el-menu-item>
      <el-menu-item index="annotations">注解</el-menu-item>
      <el-menu-item index="events">事件</el-menu-item>
    </el-menu>
    <div class="box box-item" v-show="activeIndex==='status'">
      <div class="box-body" style="padding-top:20px">
        <el-row class="statistic">
          <el-col :span="8">
            <div class="statistic">
              <div class="statistic__head">Replicas</div>
              <div :class="`statistic__content ${statefulsetInfo.status?.readyReplicas<statefulsetInfo.spec?.replicas?'text-red':'text-green'}`" >
                {{  statefulsetInfo.status?.readyReplicas }} / {{ statefulsetInfo.spec?.replicas }}
              </div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="statistic">
              <div class="statistic__head">Resource.Limits ( cpu/memory )</div>
              <div class="statistic__content">
                {{  statefulsetInfo.spec?.template.spec.containers[0].resources?.limits?.cpu }} / {{ statefulsetInfo.spec?.template.spec.containers[0].resources?.limits?.memory }}
              </div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="statistic">
              <div class="statistic__head">Resource.Requests</div>
              <div class="statistic__content">
                {{  statefulsetInfo.spec?.template.spec.containers[0].resources?.requests?.cpu }} / {{ statefulsetInfo.spec?.template.spec.containers[0].resources?.requests?.memory }}
              </div>
            </div>
          </el-col>
        </el-row>
      </div>
    </div>
    <podList v-show="activeIndex==='status'" ref="refPods" workloadType="statefulsets" />
    <div class="box box-item" v-show="activeIndex==='labels'">
      <div class="box-body" style="padding-top:20px">
        <el-descriptions :column="1" border class="no-color" :label-width="150">
          <el-descriptions-item v-for="(v,k,i) in statefulsetInfo.metadata?.labels" :index="i" :label="k">{{ v }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </div>
    <div class="box box-item" v-show="activeIndex==='annotations'">
      <div class="box-body" style="padding-top:20px">
        <el-descriptions :column="1" border class="no-color" :label-width="180">
          <el-descriptions-item v-for="(v,k,i) in statefulsetInfo.metadata?.annotations" :index="i" :label="k">{{ v }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </div>
    <eventList v-show="activeIndex==='events'" ref="refEvents" resourceType="StatefulSet" :resourceName="route.params.workload_name" />
  </el-col>
</el-row>
<el-drawer v-model="show.yaml" direction="rtl" size="700px">
  <template #header>
    <h4>编辑YAML</h4>
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
          minLines: 10,
          maxLines: 5000,
        }" />
    </div>
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
import { ArrowRight } from '@element-plus/icons-vue'
import { onBeforeUnmount, onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRoute } from 'vue-router'
import podList from './podList.vue'
import eventList from '../eventList.vue'
import { axios } from '/src/assets/util/axios'
import moment from 'moment'
/* 引入v-ace-editor */
import { VAceEditor } from 'vue3-ace-editor'
import 'ace-builds/src-noconflict/mode-yaml'
import 'ace-builds/src-noconflict/theme-chrome'
import 'ace-builds/src-noconflict/ext-language_tools'
/* 变量定义 */
const route = useRoute()
const statefulsetInfo = ref({})
const activeIndex = ref("status")
const show = ref({
  event: false,
  yaml: false,
  log: false,
})
const yamlContent = ref("")
const timer = ref(null)
const wrapLine = ref(true)
const refPods = ref(null)
const refEvents = ref(null)
/* 生命周期函数 */
onMounted(async () => {
  doRequest()
  timer.value = setInterval(() => {
    doRequest()
  }, 15000)
})
onBeforeUnmount(() => {
  if(timer) {
    clearInterval(timer.value)
    timer.value = null
  }
})
/* methods */
const doRequest = () => {
  getStatefulset()
  refEvents.value.getEvents()
  refPods.value.getPods()
}
const getStatefulset = async () => {
  statefulsetInfo.value = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/statefulsets/${route.params.workload_name}`)
  statefulsetInfo.value.metadata.creationTimestamp = moment(statefulsetInfo.value.metadata.creationTimestamp).format('YYYY-MM-DD HH:mm:ss')
  statefulsetInfo.value.status.readyReplicas ||= 0
}
const handleSelect = async (key) => {
  activeIndex.value = key
}
const handleCommand = async (command) => {
  switch(command) {
    case "restart": {
      await ElMessageBox.confirm('确定重启？','警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(async () => {
        await axios.patch(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/statefulsets/${route.params.workload_name}/rollout`)
        ElMessage.success({message: '操作成功'})
      }).catch(() =>{})
      break
    }
    case "yaml": {
      yamlContent.value = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/statefulsets/${route.params.workload_name}/yaml`)
      show.value.yaml = true
      break
    }
    case "delete": {
      await ElMessageBox.confirm('确定删除？','警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(async () => {
        await axios.delete(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/statefulsets/${route.params.workload_name}`)
        ElMessage.success({message: '删除成功'})
      }).catch(() =>{})
      break
    }
  }
  setTimeout(async () => {
    await doRequest()
  }, 2000)
}
const submitYaml = async () => {
  await axios.patch(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/apply/yaml`, yamlContent.value, {
    headers: {
      'Content-Type': 'text/plain'
    }
  })
  show.value.yaml = false
  setTimeout(async () => {
    await doRequest()
  }, 2000)
}
</script>

<style scoped>
.my-header {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  gap: 16px;
}
</style>
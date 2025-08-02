<template>
<el-breadcrumb :separator-icon="ArrowRight" style="min-width:1366px">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>应用负载</el-breadcrumb-item>
  <el-breadcrumb-item :to="{ path: '/kubernetes/persistentvolumeclaims', query: {cluster: route.query.cluster, namespace: route.query.namespace} }">存储卷</el-breadcrumb-item>
  <el-breadcrumb-item>{{ pvcInfo.metadata?.name }}</el-breadcrumb-item>
</el-breadcrumb>
<el-row :gutter="15" style="min-width:1366px">
  <el-col :span="6">
    <el-card>
      <template #header>
        <div class="card-header">
          <el-text truncated size="large" style="width:250px;font-size:18px" class="text-black"><b>{{ pvcInfo.metadata?.name }}</b></el-text>
          <el-dropdown @command="handleCommand" class="pull-right" style="top:2px;min-width:75px;">
            <el-link underline="never">
              更多操作
              <el-icon class="el-icon--right">
                <arrow-down />
              </el-icon>
            </el-link>
            <template #dropdown>
              <el-dropdown-menu>
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
        <el-descriptions-item label="状态">{{ pvcInfo.status?.phase||'UnBound' }}</el-descriptions-item>
        <el-descriptions-item label="容量">{{ pvcInfo.spec?.resources.requests.storage }}</el-descriptions-item>
        <el-descriptions-item label="存储类">{{ pvcInfo.spec?.storageClassName }}</el-descriptions-item>
        <el-descriptions-item label="访问模式">{{ pvcInfo.spec?.accessModes.join(",") }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ pvcInfo.metadata?.creationTimestamp }}</el-descriptions-item>
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
    <podList v-show="activeIndex==='status'" ref="refPods" :pods="pods" style="margin-top:15px" />
    <div class="box box-item" v-show="activeIndex==='labels'">
      <div class="box-body" style="padding-top:20px">
        <el-descriptions :column="1" border class="no-color" :label-width="150">
          <el-descriptions-item v-for="(v,k,i) in pvcInfo.metadata?.labels" :index="i" :label="k">{{ v }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </div>
    <div class="box box-item" v-show="activeIndex==='annotations'">
      <div class="box-body" style="padding-top:20px">
        <el-descriptions :column="1" border class="no-color" :label-width="180">
          <el-descriptions-item v-for="(v,k,i) in pvcInfo.metadata?.annotations" :index="i" :label="k">{{ v }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </div>
    <eventList v-show="activeIndex==='events'" ref="refEvents" resourceType="PersistentVolumeClaim" :resourceName="route.params.pvc_name" />
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
import { onMounted, onBeforeUnmount, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRoute, useRouter } from 'vue-router'
import podList from '../workload/podList.vue'
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
const router = useRouter()
const pvcInfo = ref({})
const activeIndex = ref("status")
const timer = ref(null)
const refPods = ref(null)
const refEvents = ref(null)
const show = ref({
  yaml: false
})
const yamlContent = ref("")
const wrapLine = ref(true)
const pods = ref([])
/* 生命周期函数 */
onMounted(async () => {
  await getPersistentVolumeClaim()
  refPods.value.getPods()
  refEvents.value.getEvents()
  timer.value = setInterval(async () => {
    await getPods()
    refPods.value.getPods()
    refEvents.value.getEvents()
  }, 15000)
})
onBeforeUnmount(() => {
  if(timer) {
    clearInterval(timer.value)
    timer.value = null
  }
})
/* methods */
const getPersistentVolumeClaim = async () => {
  pvcInfo.value = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/persistentvolumeclaims/${route.params.pvc_name}`)
  if(!pvcInfo.value.metadata) {
    ElMessage.error({message: `未找到存储卷: ${route.params.pvc_name}`})
    return
  }
  pvcInfo.value.metadata.creationTimestamp = moment(pvcInfo.value.metadata.creationTimestamp).format('YYYY-MM-DD HH:mm:ss')
  await getPods()
}
const getPods = async () => {
  let response = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/pods?limit=5000`)
  pods.value = response.results.filter(n => {
    let found = false
    for(let volume of n.spec.volumes||[]) {
      if(volume.persistentVolumeClaim && volume.persistentVolumeClaim.claimName === pvcInfo.value.metadata.name) {
        found = true
      }
    }
    return found
  })
}
const handleCommand = async (command) => {
  switch(command) {
    case "yaml": {
      yamlContent.value = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/persistentvolumeclaims/${route.params.pvc_name}/yaml`)
      show.value.yaml = true
      break
    }
    case "delete": {
      await ElMessageBox.confirm('确定删除？','警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(async () => {
        await axios.delete(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/persistentvolumeclaims/${route.params.pvc_name}`)
        ElMessage.success({message: '删除成功'})
        router.push({
          path: '/kubernetes/persistentvolumeclaims',
          query: {
            cluster: route.query.cluster,
            namespace: route.query.namespace
          }
        })
      }).catch(() =>{})
      break
    }
  }
}
const handleSelect = async (key) => {
  activeIndex.value = key
}
const submitYaml = async () => {
  await axios.patch(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/apply/yaml`, yamlContent.value, {
    headers: {
      'Content-Type': 'text/plain'
    }
  })
  show.value.yaml = false
  await getPersistentVolumeClaim()
  refPods.value.getPods()
}
</script>
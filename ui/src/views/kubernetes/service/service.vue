<template>
<el-breadcrumb :separator-icon="ArrowRight" style="min-width:1366px">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>应用负载</el-breadcrumb-item>
  <el-breadcrumb-item :to="{ path: '/kubernetes/services', query: {cluster: route.query.cluster, namespace: route.query.namespace} }">服务</el-breadcrumb-item>
  <el-breadcrumb-item>{{ serviceInfo.metadata?.name }}</el-breadcrumb-item>
</el-breadcrumb>
<el-row :gutter="15" style="min-width:1366px">
  <el-col :span="6">
    <el-card>
      <template #header>
        <div class="card-header">
          <el-text truncated size="large" style="width:250px;font-size:18px" class="text-black"><b>{{ serviceInfo.metadata?.name }}</b></el-text>
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
        <el-descriptions-item label="服务类型">
          {{ serviceInfo.spec?.clusterIP==='None'?'Headless':'VirtualIP' }} ( {{ serviceInfo.spec?.type }} )
        </el-descriptions-item>
        <el-descriptions-item label="内部IP地址">
          <div v-for="item in serviceInfo.spec?.clusterIPs" :key="item">{{ item }}</div>
        </el-descriptions-item>
        <el-descriptions-item label="外部IP地址">
          <div v-if="serviceInfo.status?.loadBalancer&&Object.values(serviceInfo.status?.loadBalancer).length>0">
            {{ Object.values(serviceInfo.status?.loadBalancer)[0][0].ip }}
          </div>
        </el-descriptions-item>
        <el-descriptions-item label="会话亲和性">{{ serviceInfo.spec?.sessionAffinity }}</el-descriptions-item>
        <el-descriptions-item label="选择器">
          <div v-for="(v,k,i) in serviceInfo.spec?.selector" :key="i">
            <el-tooltip placement="top" :content="`${k}=${v}`">
              <el-text style="width:200px" truncated>{{ k }}={{ v }}</el-text>
            </el-tooltip>
          </div>
        </el-descriptions-item>
        <el-descriptions-item label="DNS">{{ serviceInfo.metadata?.name }}.{{ route.query.namespace }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ serviceInfo.metadata?.creationTimestamp }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ serviceInfo.metadata?.lastUpdateTime }}</el-descriptions-item>
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
    <el-card v-show="activeIndex==='status'" style="margin-top:15px">
      <template #header>
        <div class="card-header">
          <span class="card-header-text">端口</span>
        </div>
      </template>
      <el-table :data="serviceInfo.spec?.ports" :show-header="false">
        <el-table-column width="60">
          <template #default="scope">
            <font-awesome-icon icon="diagram-project" style="font-size:25px;vertical-align:middle;" />
          </template>
        </el-table-column>
        <el-table-column width="150" align="center">
          <template #default="scope">
            <b>{{ scope.row.port }}</b>
            <div class="text-gray cell-comment" style="font-size:14px;justify-content: center;">服务端口</div>
          </template>
        </el-table-column>
        <el-table-column width="100" align="center">
          <template #default="scope">
            <el-icon><Right /></el-icon> {{ scope.row.protocol }} <el-icon><Right /></el-icon>
          </template>
        </el-table-column>
        <el-table-column width="150" align="center">
          <template #default="scope">
            <b>{{ scope.row.targetPort }}</b>
            <div class="text-gray cell-comment" style="font-size:14px;justify-content: center;">容器端口</div>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    <el-card v-show="activeIndex==='status'" style="margin-top:15px" v-if="deploymentList.length>0">
      <template #header>
        <div class="card-header">
          <span class="card-header-text">部署</span>
        </div>
      </template>
      <el-table :data="deploymentList" :show-header="false">
        <el-table-column label="" width="45">
          <font-awesome-icon icon="layer-group" style="font-size:25px;vertical-align:middle;" />
        </el-table-column>
        <el-table-column prop="name" label="名称" min-width="200">
          <template #default="scope">
            <el-link underline="never" :href="`/kubernetes/workload/deployments/${scope.row.name}?cluster=${route.query.cluster}&namespace=${route.query.namespace}`">{{ scope.row.name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column label="状态" min-width="200">
          <template #default="scope">
            <font-awesome-icon icon="circle" v-if="scope.row.replicas==0" class="runningstatus text-gray" />
            <font-awesome-icon icon="circle" v-else-if="scope.row.available=='False'" class="runningstatus twinkling text-yellow" />
            <font-awesome-icon icon="circle" v-else-if="scope.row.available=='True'" class="runningstatus text-green" />
            <span v-if="scope.row.replicas === 0">停止 ( {{ scope.row.readyReplicas }} / {{ scope.row.replicas }} )</span>
            <span v-else>{{ scope.row.available === 'True' ? '运行中' : '更新中' }} ( {{ scope.row.readyReplicas }} / {{ scope.row.replicas }} )</span>
          </template>
        </el-table-column>
        <el-table-column label="更新时间">
          <template #default="scope">
            {{ moment(scope.row.lastUpdateTime).format('YYYY-MM-DD HH:mm:ss') }}
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    <el-card v-show="activeIndex==='status'" style="margin-top:15px" v-if="statefulsetList.length>0">
      <template #header>
        <div class="card-header">
          <span class="card-header-text">有状态副本集</span>
        </div>
      </template>
      <el-table :data="statefulsetList" :show-header="false">
        <el-table-column label="" width="45">
          <font-awesome-icon icon="layer-group" style="font-size:25px;vertical-align:middle;" />
        </el-table-column>
        <el-table-column prop="name" label="名称" min-width="200">
          <template #default="scope">
            <el-link underline="never" :href="`/kubernetes/workload/statefulsets/${scope.row.name}?cluster=${route.query.cluster}&namespace=${route.query.namespace}`">{{ scope.row.name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column label="状态" min-width="200">
          <template #default="scope">
            <font-awesome-icon icon="circle" v-if="scope.row.replicas==0" class="runningstatus text-gray" />
            <font-awesome-icon icon="circle" v-else-if="scope.row.available=='False'" class="runningstatus twinkling text-yellow" />
            <font-awesome-icon icon="circle" v-else-if="scope.row.available=='True'" class="runningstatus text-green" />
            <span v-if="scope.row.replicas === 0">停止 ( {{ scope.row.readyReplicas }} / {{ scope.row.replicas }} )</span>
            <span v-else>{{ scope.row.available === 'True' ? '运行中' : '更新中' }} ( {{ scope.row.readyReplicas }} / {{ scope.row.replicas }} )</span>
          </template>
        </el-table-column>
        <el-table-column label="更新时间">
          <template #default="scope">
            {{ moment(scope.row.lastUpdateTime).format('YYYY-MM-DD HH:mm:ss') }}
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    <podList v-show="activeIndex==='status'" ref="refPods" :labelSelector="labelSelector" />
    <div class="box box-item" v-show="activeIndex==='labels'">
      <div class="box-body" style="padding-top:20px">
        <el-descriptions :column="1" border class="no-color" :label-width="150">
          <el-descriptions-item v-for="(v,k,i) in serviceInfo.metadata?.labels" :index="i" :label="k">{{ v }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </div>
    <div class="box box-item" v-show="activeIndex==='annotations'">
      <div class="box-body" style="padding-top:20px">
        <el-descriptions :column="1" border class="no-color" :label-width="180">
          <el-descriptions-item v-for="(v,k,i) in serviceInfo.metadata?.annotations" :index="i" :label="k">{{ v }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </div>
    <eventList v-show="activeIndex==='events'" ref="refEvents" resourceType="Service" :resourceName="route.params.service_name" />
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
const serviceInfo = ref({})
const activeIndex = ref("status")
const labelSelector = ref("")
const timer = ref(null)
const refPods = ref(null)
const refEvents = ref(null)
const deploymentList = ref([])
const statefulsetList = ref([])
const show = ref({
  yaml: false
})
const yamlContent = ref("")
const wrapLine = ref(true)
/* 生命周期函数 */
onMounted(async () => {
  await getService()
  listWorkloads()
  refPods.value.getPods()
  refEvents.value.getEvents()
  timer.value = setInterval(() => {
    listWorkloads()
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
const getService = async () => {
  serviceInfo.value = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/services/${route.params.service_name}`)
  if(!serviceInfo.value.metadata) {
    ElMessage.error({message: `未找到服务: ${route.params.service_name}`})
    return
  }
  serviceInfo.value.metadata.creationTimestamp = moment(serviceInfo.value.metadata.creationTimestamp).format('YYYY-MM-DD HH:mm:ss')
  serviceInfo.value.metadata.lastUpdateTime = moment(serviceInfo.value.metadata.lastUpdateTime).format('YYYY-MM-DD HH:mm:ss')
  let ls = []
  for(let [k,v] of Object.entries(serviceInfo.value.spec.selector)) {
    ls.push(`${k}=${v}`)
  }
  labelSelector.value = ls.join(",")
}
const listWorkloads = async () => {
  // deployment
  let response = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/deployments?label_selector=${labelSelector.value}`)
  deploymentList.value = response.results.map(x => {
    let progress = x.status.conditions.find(n => n.type === 'Progressing')
    let r = {
      name: x.metadata.name,
      replicas: x.status.replicas||0,
      readyReplicas: x.status.readyReplicas||0,
      unavailableReplicas: x.status.unavailableReplicas||0,
      lastUpdateTime: progress?.lastUpdateTime
    }
    r.available = r.readyReplicas >= r.replicas ? 'True' : 'False'
    return r
  })
  // statefulset
  response = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/statefulsets?label_selector=${labelSelector.value}`)
  statefulsetList.value = response.results.map(x => {
    let r = {
      name: x.metadata.name,
      replicas: x.status.replicas||0,
      readyReplicas: x.status.readyReplicas||0,
      unavailableReplicas: x.status.unavailableReplicas||0,
      creationTimestamp: x.metadata.creationTimestamp
    }
    r.available = r.readyReplicas >= r.replicas ? 'True' : 'False'
    return r
  })
}
const handleCommand = async (command) => {
  switch(command) {
    case "yaml": {
      yamlContent.value = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/services/${route.params.service_name}/yaml`)
      show.value.yaml = true
      break
    }
    case "delete": {
      await ElMessageBox.confirm('确定删除？','警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(async () => {
        await axios.delete(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/services/${route.params.service_name}`)
        ElMessage.success({message: '删除成功'})
        router.push({
          path: '/kubernetes/services',
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
  await getService()
  listWorkloads()
  refPods.value.getPods()
}
</script>
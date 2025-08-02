<template>
<el-breadcrumb :separator-icon="ArrowRight" style="min-width:1366px">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>应用负载</el-breadcrumb-item>
  <el-breadcrumb-item :to="{ path: '/kubernetes/workload/pods', query:{cluster:route.query.cluster,namespace:route.query.namespace} }">容器组</el-breadcrumb-item>
  <el-breadcrumb-item>{{ route.params.pod_name }}</el-breadcrumb-item>
</el-breadcrumb>
<el-row :gutter="15" style="min-width:1366px">
  <el-col :span="6" >
    <el-card>
      <template #header>
        <div class="card-header">
          <el-text truncated size="large" style="width:250px;font-size:18px" class="text-black"><b>{{ podInfo.metadata?.name }}</b></el-text>
          <el-dropdown @command="handleCommand" class="pull-right" style="top:2px;min-width:75px;">
            <el-link underline="never">
              更多操作
              <el-icon class="el-icon--right">
                <arrow-down />
              </el-icon>
            </el-link>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="yaml">查看YAML</el-dropdown-item>
                <el-dropdown-item command="delete">删除</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </template>
      <el-descriptions :column="1" border class="no-color" :label-width="120">
        <el-descriptions-item label="集群">{{ route.query.cluster }}</el-descriptions-item>
        <el-descriptions-item label="命名空间">{{ route.query.namespace }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ podInfo.state }}</el-descriptions-item>
        <el-descriptions-item label="容器组IP">{{ podInfo.status?.podIP }}</el-descriptions-item>
        <el-descriptions-item label="节点名称">{{ podInfo.spec?.nodeName }}</el-descriptions-item>
        <el-descriptions-item label="节点IP">{{ podInfo.status?.hostIP }}</el-descriptions-item>
        <el-descriptions-item label="重启次数">{{ podInfo.status?.containerStatuses[0]?.restartCount }}</el-descriptions-item>
        <el-descriptions-item label="QoS类别">{{ podInfo.status?.qosClass }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ podInfo.metadata?.creationTimestamp }}</el-descriptions-item>
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
          <span class="card-header-text">容器</span>
        </div>
      </template>
      <el-table :data="podInfo.status?.containerStatuses||[]" :show-header="false">
        <el-table-column width="60">
          <template #default="props">
            <el-image style="width:30px;height:30px" src="/images/docker.svg" />
            <font-awesome-icon icon="circle" :class="`podstatus ${getPodClass(props.row.status, props.row.reason, props.row.ready)}`" />
          </template>
        </el-table-column>
        <el-table-column>
          <template #default="props">
            <div>
              <b style="margin-right:10px">{{props.row.name}}</b>
              <el-icon class="pointer text-blue" @click="container_name=props.row.name;getLogs(lines)"><Document /></el-icon>
              <el-tag v-if="props.row.initContainer===true" round type="warning" style="margin-left:10px">initContainer</el-tag>
            </div>
            <div class="text-gray cell-comment" v-if="props.row.ready==true">{{props.row.image}}</div>
            <div class="text-gray cell-comment" v-else>{{props.row.state_message}}</div>
          </template>
        </el-table-column>
        <el-table-column min-width="15%">
          <template #default="props">
            <div>{{props.row.status}}</div>
            <div class="text-gray cell-comment">Status</div>
          </template>
        </el-table-column>
        <el-table-column min-width="15%">
          <template #default="props">
            <div>{{props.row.restartCount}}</div>
            <div class="text-gray cell-comment">Restart Count</div>
          </template>
        </el-table-column>
        <el-table-column min-width="20%">
          <template #default="props">
            <div>{{ portsMap[props.row.name]?.map(n => `${n.containerPort}/${n.protocol}`).join(',')||'-' }}</div>
            <div class="text-gray cell-comment">Ports</div>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    <el-card v-show="activeIndex==='status'">
      <template #header>
        <div class="card-header">
          <span class="card-header-text">卷</span>
        </div>
      </template>
      <el-table :data="volumeList" :show-header="false">
        <el-table-column width="60">
          <template #default="props">
            <font-awesome-icon icon="floppy-disk" style="font-size:25px;vertical-align:middle;" />
          </template>
        </el-table-column>
        <el-table-column min-width="130">
          <template #default="props">
            <div><b>{{props.row.name}}</b></div>
            <div class="text-gray cell-comment">storageClass: {{props.row.storageClass||'-'}}</div>
          </template>
        </el-table-column>
        <el-table-column min-width="150">
          <template #default="props">
            <div><b>{{props.row.claimName}}</b></div>
            <div class="text-gray cell-comment">{{props.row.type}}</div>
          </template>
        </el-table-column>
        <el-table-column min-width="200">
          <template #default="props">
            <div><b>{{props.row.volumeMounts?.mountPath}}</b></div>
            <div class="text-gray cell-comment">mountPath</div>
          </template>
        </el-table-column>
        <el-table-column width="120">
          <template #default="props">
            <div><b>{{props.row.storage||'-'}}</b></div>
            <div class="text-gray cell-comment">storageSize</div>
          </template>
        </el-table-column>
        <el-table-column min-width="100">
          <template #default="props">
            <div><b>{{props.row.accessMode||'-'}}</b></div>
            <div class="text-gray cell-comment">accessMode</div>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    <div class="box box-item" v-show="activeIndex==='labels'">
      <div class="box-body" style="padding-top:20px">
        <el-descriptions :column="1" border class="no-color" :label-width="150">
          <el-descriptions-item v-for="(v,k,i) in podInfo.metadata?.labels" :index="i" :label="k">{{ v }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </div>
    <div class="box box-item" v-show="activeIndex==='annotations'">
      <div class="box-body" style="padding-top:20px">
        <el-descriptions :column="1" border class="no-color" :label-width="180">
          <el-descriptions-item v-for="(v,k,i) in podInfo.metadata?.annotations" :index="i" :label="k">{{ v }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </div>
    <div class="box box-item" v-show="activeIndex==='events'">
      <div class="box-body" style="padding-top:20px">
        <el-table :data="eventList" style="width:100%">
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
      </div>
    </div>
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
          readOnly: true,
        }" />
    </div>
  </template>
</el-drawer>
<el-dialog id="podLog" v-model="show.log" :fullscreen="true" @opened="openLog" @close="closeLog" :show-close="false" width="100%">
  <template #header="{close, titleId, titleClass}">
    <span class="el-dialog__title">容器日志</span>
    <el-button-group class="pull-right">
      <el-button circle @click="follow=!follow">
        <font-awesome-icon icon="play" v-if="!follow" />
        <font-awesome-icon icon="stop" v-else />
      </el-button>
      <el-button circle :icon="Refresh" @click="getLogs(lines)" />
      <el-button circle :icon="Close" @click="close" />
    </el-button-group>
  </template>
  <el-scrollbar style="height:calc(100vh - 90px);margin-top:5px" ref="scrollRef">
    <pre class="dark fullscreen" ref="preRef" style="min-height:calc(100vh - 90px)">
      <p class="pointer break-line" style="text-align: center" v-if="logs.length>=lines" @click="lines+=1000;getLogs(lines)">查看更多</p>
      <p v-for="(item,i) in logs" :key="i" class="break-line">{{ item }}</p>
    </pre>
  </el-scrollbar>
</el-dialog>
</template>

<script setup>
import { ArrowRight, Refresh, Close } from '@element-plus/icons-vue'
import { onBeforeMount, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRoute } from 'vue-router'
import { getPodClass } from '@/assets/util/common'
import { axios } from '/src/assets/util/axios'
import moment from 'moment'
/* 引入v-ace-editor */
import { VAceEditor } from 'vue3-ace-editor'
import 'ace-builds/src-noconflict/mode-yaml'
import 'ace-builds/src-noconflict/theme-chrome'
import 'ace-builds/src-noconflict/ext-language_tools'
/* 变量定义 */
const route = useRoute()
const podInfo = ref({})
const activeIndex = ref("status")
const container_name = ref("")
const portsMap = ref({})
const volumeList = ref([])
const eventList = ref([])
const show = ref({
  yaml: false,
  log: false,
})
const logs = ref([])
const lines = ref(1000)
const scrollRef = ref(null)
const follow = ref(false)
const preRef = ref(null)
const yamlContent = ref("")
const wrapLine = ref(true)
/* 生命周期函数 */
onBeforeMount(async () => {
  await getPodInfo()
  getVolumeList()
  getEvents()
})
/* methods */
const getPodInfo = async () => {
  podInfo.value = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/pods/${route.params.pod_name}`)
  let conditionReady = podInfo.value.status.conditions.find(n => {
    return n.type == 'Ready'
  })
  podInfo.value.metadata.creationTimestamp = moment(podInfo.value.metadata.creationTimestamp).format('YYYY-MM-DD HH:mm:ss')
  podInfo.value.status.ready = conditionReady?.status||'False'
  let containerStatuses = podInfo.value.status.containerStatuses?.map(y => {
    y.state_message = y.image
    y.status = Object.keys(y.state)[0]
    if(y.status!=='running') {
      y.state_message = y.state[y.status].reason
      y.reason = y.state[y.status].reason
    } else if(y.ready === false) { // 即使为running，也可能ready=False，需要显示ready为False的reason
      y.state_message = conditionReady.reason
    } 
    return y
  }) || podInfo.value.spec.containers.map(y => {  // 无法调度的pod没有containerStatuses字段
    return {
      name: y.name,
      state_message: podInfo.value.status.conditions[0].reason,
      reason: podInfo.value.status.conditions[0].reason,
      status: 'waiting'
    }
  })
  let initContainerStatuses = podInfo.value.status.initContainerStatuses?.map(y => {
    y.initContainer = true
    y.state_message = y.image
    y.status = Object.keys(y.state)[0]
    if(y.status!=='running') {
      y.state_message = y.state[y.status].reason
      y.reason = y.state[y.status].reason
    } else if(y.ready === false) { // 即使为running，也可能ready=False，需要显示ready为False的reason
      y.state_message = conditionReady.reason
    }
    return y
  }) || podInfo.value.spec.initContainers?.map(y => {  // 无法调度的pod没有containerStatuses字段
    return {
      name: y.name,
      state_message: podInfo.value.status.conditions[0].reason,
      reason: podInfo.value.status.conditions[0].reason,
      status: 'waiting'
    }
  }) || []
  podInfo.value.status.containerStatuses = containerStatuses.concat(initContainerStatuses)
  podInfo.value.state = podInfo.value.status.containerStatuses[0].status // running/waiting/terminated
  if(podInfo.value.state === 'running' && conditionReady.status === 'True') {
    podInfo.value.state_message = `Created ${moment.duration(moment(podInfo.value.status.containerStatuses[0].state.running.startedAt)-moment()).humanize(true)}`
  } else {
    podInfo.value.state_message = podInfo.value.status.containerStatuses[0].state_message
  }
  for(let x of podInfo.value.spec.containers) {
    portsMap.value[x.name] = x.ports
  }
}
const getVolumeList = async () => {
  let volumeMounts = []
  for(let x of podInfo.value.spec.containers) {
    for(let y of x.volumeMounts) {
      volumeMounts.push(y)
    }
  }
  volumeList.value = podInfo.value.spec.volumes.map(x => {
    if(x.hasOwnProperty('hostPath')) {
      x.type = 'hostPath'
      x.claimName = x.hostPath.path
    } else if(x.hasOwnProperty('persistentVolumeClaim')) {
      x.type = 'persistentVolumeClaim'
      x.claimName = x.persistentVolumeClaim.claimName
    } else {
      x.type = '-'
      x.claimName = '-'
    }
    x.volumeMounts = volumeMounts.find(n => n.name === x.name)
    return x
  })
  for(let x of volumeList.value) {
    if(x.type === 'persistentVolumeClaim') {
      let response = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/persistentvolumeclaims/${x.claimName}`)
      x.storage = response.spec.resources.requests.storage
      x.storageClass = response.spec.storageClassName
      x.accessMode = response.spec.accessModes?.join(',')
    }
  }
}
const handleSelect = async (key) => {
  activeIndex.value = key
}
const handleCommand = async (command) => {
  switch(command) {
    case "yaml": {
      yamlContent.value = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/pods/${route.params.pod_name}/yaml`)
      show.value.yaml = true
      break
    }
    case "delete": {
      await ElMessageBox.confirm('确定删除？','警告', {
        distinguishCancelAndClose: true,
        confirmButtonText: '正常删除',
        cancelButtonText: '强制删除',
        type: 'warning',
      }).then(async () => {
        await axios.delete(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/pods/${route.params.pod_name}`)
        ElMessage.success({message: '删除成功'})
      }).catch(async (action) =>{
        if(action === 'cancel') {
          await axios.delete(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/pods/${route.params.pod_name}?force=true`)
          ElMessage.success({message: '删除成功'})
        }
      })
      break
    }
  }
}
const getEvents = async () => {
  let response = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/Pod/${route.params.pod_name}/events`)
  eventList.value = response.map(x => {
    x.age = moment.duration(moment(x.lastTimestamp)-moment()).humanize(true)
    return x
  })
}
const getLogs = async (lines) => {
  let response = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/pods/${route.params.pod_name}/logs?container=${container_name.value}&lines=${lines}`)
  logs.value = response.split("\n")
  show.value.log = true
}
const openLog = () => {
  scrollRef.value.setScrollTop(preRef.value.clientHeight)
}
const closeLog = () => {}
</script>
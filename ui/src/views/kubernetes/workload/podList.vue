<template>
<div>
  <el-card>
    <template #header>
      <div class="card-header">
        <span class="card-header-text">容器组</span>
        <div class="box-tools pull-right">
          <span class="card-header-btn" @click="getPods()"><el-icon><Refresh /></el-icon></span>
        </div>
      </div>
    </template>
    <el-table :data="podList" :show-header="false" style="width:100%;">
      <el-table-column label="icon" width="60">
        <template #default="scope">
          <font-awesome-icon icon="cubes" style="font-size:25px" />
          <font-awesome-icon icon="circle" :class="`podstatus ${getPodClass(scope.row.state, scope.row.reason, scope.row.status?.ready)}`" />
        </template>
      </el-table-column>
      <el-table-column>
        <template #default="scope">
          <el-link underline="never" :href="`/kubernetes/workload/pods/${scope.row.pod_name}?cluster=${route.query.cluster}&namespace=${route.query.namespace}`"><b>{{scope.row.pod_name}}</b></el-link>
          <div v-if="scope.row.state==='waiting'" class="text-yellow cell-comment">
            <el-icon><WarningFilled /></el-icon>
            {{scope.row.state_message}}
          </div>
          <div v-else-if="scope.row.state==='terminated'" :class="`text-${scope.row.reason==='Error'?'red':'gray'} cell-comment`">
            <el-icon v-if="scope.row.reason==='Error'"><WarningFilled /></el-icon>
            {{scope.row.state_message}}
          </div>
          <div v-else-if="scope.row.state==='deleting'" class="text-yellow cell-comment">
            <el-icon><WarningFilled /></el-icon>
            {{scope.row.state}}
          </div>
          <div v-else :class="`text-${scope.row.status.ready==='False'?'yellow':'gray'} cell-comment`">{{scope.row.state_message}}</div>
        </template>
      </el-table-column>
      <el-table-column>
        <template #default="scope">
          <div>{{scope.row.node_name}} ( {{scope.row.hostip}} )</div>
          <div class="text-gray cell-comment">Worker Node</div>
        </template>
      </el-table-column>
      <el-table-column>
        <template #default="scope">
          <div>{{scope.row.podip}}</div>
          <div class="text-gray cell-comment">Pod IP</div>
        </template>
      </el-table-column>
      <el-table-column width="100">
        <template #default="scope">
          <el-link type="primary" underline="never" @click="getPodEvents(scope.row)">查看</el-link>
          <div class="text-gray cell-comment">Events</div>
        </template>
      </el-table-column>
      <el-table-column type="expand" width="45">
        <template #default="scope">
          <div style="padding-left:30px">
            <div class="text-gray" style="line-height:30px">Containers</div>
            <el-table :data="scope.row.status.containerStatuses" :show-header="false">
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
                    <el-icon class="pointer text-blue" @click="pod_name=scope.row.pod_name;container_name=props.row.name;getLogs(lines)">
                      <Document />
                    </el-icon>
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
            </el-table>
          </div>
        </template>
      </el-table-column>
    </el-table>
  </el-card>
  <el-dialog
    v-model="show.event"
    :title="`${currentPod.pod_name} 事件`"
    width="70%">
    <el-table :data="podEventList" style="width:100%">
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
  </el-dialog>
  <el-dialog id="podLog" v-model="show.log" :fullscreen="true" @opened="openLog" @close="closeLog" :show-close="false" width="100%">
    <template #header="{close, titleId, titleClass}">
      <span class="el-dialog__title">容器日志</span>
      <el-button-group class="pull-right">
        <el-button circle @click="followLogs(100)">
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
</div>
</template>

<script setup>
import { Refresh, Close } from '@element-plus/icons-vue'
import { onBeforeMount, onBeforeUnmount, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { forPodList, getPodClass } from '@/assets/util/common'
import { axios } from '/src/assets/util/axios'
import moment from 'moment'
import { v4 as uuidv4 } from 'uuid'
/* 变量定义 */
const props = defineProps({
  pods: { type: Array },
  workloadType: { type: String }, 
  labelSelector: { type: String },
})
const route = useRoute()
const podList = ref([])
const currentPod = ref({})
const podEventList = ref([])
const show = ref({
  event: false,
  log: false,
  terminal: false,
})
const wsto = ref("")
const logs = ref([])
const socket = ref(null)
const follow = ref(false)
const pod_name = ref("")
const container_name = ref("")
const lines = ref(1000)
const scrollRef = ref(null)
const preRef = ref(null)
/* 生命周期函数 */
onBeforeMount(async () => {
  // getPods()
})
onMounted(() => {
  wsto.value = uuidv4()
  // initLogSocket()
})
onBeforeUnmount(() => {
  if (socket.value) {
    socket.value.close()
  }
})
/* methods */
const getPods = async () => {
  if(props.pods) {
    podList.value = forPodList(props.pods)
    return
  }
  var url
  if(props.workloadType) {
    url = `/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/${props.workloadType}/${route.params.workload_name}/pods`
  } else if(props.labelSelector) {
    url = `/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/pods?label_selector=${props.labelSelector}`
  } else {
    return
  }
  let response = await axios.get(url)
  let results = response.results||response
  podList.value = forPodList(results)
}
const getPodEvents = async (row) => {
  currentPod.value = row
  let response = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/Pod/${row.pod_name}/events`)
  podEventList.value = response.map(x => {
    x.age = moment.duration(moment(x.lastTimestamp)-moment()).humanize(true)
    return x
  })
  show.value.event = true
}
const getLogs = async (lines) => {
  let response = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/pods/${pod_name.value}/logs?container=${container_name.value}&lines=${lines}`)
  logs.value = response.split("\n")
  show.value.log = true
}
const followLogs = async (lines) => {
  if(!follow.value) {
    logs.value = []
    initLogSocket(lines)
  } else {
    if (socket.value) {
      socket.value.close()
    }
    getLogs(1000)
  }
  follow.value = !follow.value
}
const openLog = () => {
  scrollRef.value.setScrollTop(preRef.value.clientHeight)
}
const closeLog = () => {
  if (socket.value) {
    socket.value.close()
  }
  follow.value = false
}
const initLogSocket = (lines) => {
  socket.value = new WebSocket(`/ws?id=${wsto.value}_log`)
  socket.value.onopen = async () => {
    console.log('successfully connected to websocket server')
    try {
      await socket.value.send(JSON.stringify({
        "message_type":"LIZARDCD_WS_PODLOG",
        "message_to": wsto.value + "_log",
        "message_data": {
          "cluster": route.query.cluster,
          "namespace": route.query.namespace,
          "pod_name": pod_name.value,
          "container": container_name.value,
          "lines": lines,
          "follow": true
        }
      }))
    }
    catch(e) {
      console.warn(e)
    }
  }
  socket.value.onclose = () => {
    console.warn('disconnected from websocket server')
  }
  socket.value.onerror = (error) => {
    console.error('websocket error:', error)
  }
  socket.value.onmessage = (event) => {
    for(let x of event.data.split('\n')) {
      try {
        let msg = JSON.parse(x)
        logs.value.push(msg.message_data.content)
      }
      catch(e) {
        console.warn(e)
      }
      if(scrollRef.value)
        scrollRef.value.setScrollTop(preRef.value.clientHeight)
    }
  }
}
defineExpose({ getPods })
</script>
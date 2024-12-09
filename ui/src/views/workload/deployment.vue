<template>
<el-breadcrumb :separator-icon="ArrowRight" style="min-width:1366px">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>工作负载</el-breadcrumb-item>
  <el-breadcrumb-item :to="{ path: '/workload/deployments', query: {cluster: route.query.cluster, namespace: route.query.namespace} }">部署</el-breadcrumb-item>
  <el-breadcrumb-item>{{ deploymentInfo.metadata?.name }}</el-breadcrumb-item>
</el-breadcrumb>
<el-row :gutter="15" style="min-width:1366px">
  <el-col :span="6">
    <el-card>
      <template #header>
        <div class="card-header">
          <span class="card-header-text"><b>{{ deploymentInfo.metadata?.name }}</b></span>
          <div class="box-tools pull-right" style="top:5px">
            <span class="card-header-btn">
              <el-dropdown @command="handleCommand">
                <el-link :underline="false" type="primary">
                  更多操作
                  <el-icon class="el-icon--right">
                    <arrow-down />
                  </el-icon>
                </el-link>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="restart">重启</el-dropdown-item>
                    <el-dropdown-item command="autoscaler">弹性伸缩</el-dropdown-item>
                    <el-dropdown-item command="yaml">编辑YAML</el-dropdown-item>
                    <el-dropdown-item command="delete">删除</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </span>
          </div>
        </div>
      </template>
      <el-descriptions :column="1">
        <el-descriptions-item label="集群">{{ route.query.cluster }}</el-descriptions-item>
        <el-descriptions-item label="命名空间">{{ route.query.namespace }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ deploymentInfo.metadata?.creationTimestamp }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ deploymentInfo.metadata?.lastUpdateTime }}</el-descriptions-item>
      </el-descriptions>
    </el-card>
  </el-col>
  <el-col :span="18">
    <el-menu :default-active="activeIndex" mode="horizontal" @select="handleSelect">
      <el-menu-item index="status">资源状态</el-menu-item>
      <el-menu-item index="labels">标签</el-menu-item>
      <el-menu-item index="annotations">注解</el-menu-item>
      <el-menu-item index="events">事件</el-menu-item>
      <el-menu-item index="autoscaler" v-if="hpa.found">弹性伸缩</el-menu-item>
    </el-menu>
    <div class="box box-item" v-show="activeIndex==='status'">
      <div class="box-body" style="padding-top:20px">
        <el-row class="statistic">
          <el-col :span="8">
            <div class="statistic">
              <div class="statistic__head">Replicas</div>
              <div :class="`statistic__content ${deploymentInfo.status?.readyReplicas<deploymentInfo.spec?.replicas?'text-red':'text-green'}`" >
                {{  deploymentInfo.status?.readyReplicas }} / {{ deploymentInfo.spec?.replicas }}
              </div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="statistic">
              <div class="statistic__head">Resource.Limits ( cpu/memory )</div>
              <div class="statistic__content">
                {{  deploymentInfo.spec?.template.spec.containers[0].resources?.limits?.cpu }} / {{ deploymentInfo.spec?.template.spec.containers[0].resources?.limits?.memory }}
              </div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="statistic">
              <div class="statistic__head">Resource.Requests</div>
              <div class="statistic__content">
                {{  deploymentInfo.spec?.template.spec.containers[0].resources?.requests?.cpu }} / {{ deploymentInfo.spec?.template.spec.containers[0].resources?.requests?.memory }}
              </div>
            </div>
          </el-col>
        </el-row>
      </div>
    </div>
    <el-card v-show="activeIndex==='status'">
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
            <font-awesome-icon icon="circle" v-if="scope.row.status.ready=='False'" class="podstatus twinkling text-yellow" />
            <font-awesome-icon icon="circle" v-if="scope.row.status.ready=='True'" class="podstatus text-green" />
          </template>
        </el-table-column>
        <el-table-column>
          <template #default="scope">
            <div><b>{{scope.row.pod_name}}</b></div>
            <div v-if="scope.row.state==='waiting'" class="text-yellow cell-comment">
              <el-icon><WarningFilled /></el-icon>
              {{scope.row.state_message}}
            </div>
            <div v-else-if="scope.row.state==='terminated'" class="text-red cell-comment">
              <el-icon><WarningFilled /></el-icon>
              {{scope.row.state}}
            </div>
            <div v-else class="text-gray cell-comment">{{scope.row.state_message}}</div>
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
            <el-link type="primary" :underline="false" @click="getPodEvents(scope.row)">查看</el-link>
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
                    <font-awesome-icon icon="circle" v-if="props.row.ready===false" class="twinkling podstatus text-yellow" />
                    <font-awesome-icon icon="circle" v-else-if="props.row.status==='terminated'" class="podstatus text-gray" />
                    <font-awesome-icon icon="circle" v-else class="podstatus text-green" />
                  </template>
                </el-table-column>
                <el-table-column>
                  <template #default="props">
                    <div><b>{{props.row.name}}</b> <el-icon class="pointer text-blue" @click="pod_name=scope.row.pod_name;container_name=props.row.name;getLogs(lines)"><Document /></el-icon> </div>
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
    <div class="box box-item" v-show="activeIndex==='labels'">
      <div class="box-body" style="padding-top:20px">
        <el-descriptions :column="1">
          <el-descriptions-item v-for="(v,k,i) in deploymentInfo.metadata?.labels" :index="i" :label="k">{{ v }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </div>
    <div class="box box-item" v-show="activeIndex==='annotations'">
      <div class="box-body" style="padding-top:20px">
        <el-descriptions :column="1">
          <el-descriptions-item v-for="(v,k,i) in deploymentInfo.metadata?.annotations" :index="i" :label="k">{{ v }}</el-descriptions-item>
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
    <div class="box box-item" v-if="hpa.found" v-show="activeIndex==='autoscaler'">
      <div class="box-body" style="padding-top:20px">
        <el-row class="statistic">
          <el-col :span="6">
            <div class="statistic">
              <div class="statistic__head">最大副本数量</div>
              <div class="statistic__content" >
                {{ hpa.max_pod }}
              </div>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="statistic">
              <div class="statistic__head">最小副本数量</div>
              <div class="statistic__content" >
                {{ hpa.min_pod }}
              </div>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="statistic">
              <div class="statistic__head">目标CPU使用率</div>
              <div class="statistic__content" >
                {{ hpa.cpu_usage }} %
              </div>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="statistic">
              <div class="statistic__head">目标内存使用率</div>
              <div class="statistic__content" >
                {{ hpa.mem_usage }} %
              </div>
            </div>
          </el-col>
        </el-row>
      </div>
    </div>
  </el-col>
</el-row>
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
<el-drawer v-model="show.yaml" direction="rtl" size="700px">
  <template #header>
    <h4>编辑YAML</h4>
  </template>
  <template #default>
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
        minLines: 10,
        maxLines: 5000,
        wrap: true
      }" />
  </template>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show.yaml=false">取消</el-button>
      <el-button type="primary" @click="submitYaml">提交</el-button>
    </div>
  </template>
</el-drawer>
<el-drawer v-model="show.autoscaler" direction="rtl" size="650px">
  <template #header>
    <h4>弹性伸缩</h4>
  </template>
  <template #default>
    <el-form :model="form" label-width="130px" >
      <el-form-item>
        <template #label><el-text>最小实例数量 
          <el-tooltip placement="top">
            <template #content>
              当Pod根据监控指标进行弹性伸缩时，保留的最小Pod数量
            </template>
            <el-icon class="text-yellow"><Warning /></el-icon>
          </el-tooltip></el-text>
        </template>
        <el-input-number v-model="form.min_pod" :max="100" :min="1" size="large" />
      </el-form-item>
      <el-form-item>
        <template #label><el-text>最大实例数量 
          <el-tooltip placement="top">
            <template #content>
              当Pod根据监控指标进行弹性伸缩时，扩充的最大Pod数量
            </template>
            <el-icon class="text-yellow"><Warning /></el-icon>
          </el-tooltip></el-text>
        </template>
        <el-input-number v-model="form.max_pod" :max="100" :min="0" size="large" />
      </el-form-item>
      <el-form-item>
        <template #label><el-text>目标CPU利用率 
          <el-tooltip placement="top">
            <template #content>
              当工作负载的所有Pod的平均CPU利用率低于此目标时，<br>将减少Pod数量，直至最小实例数量；<br>
              当工作负载的所有Pod的平均CPU利用率高于此目标时，<br>将增加Pod数量，直至最大实例数量
            </template>
            <el-icon class="text-yellow"><Warning /></el-icon>
          </el-tooltip></el-text>
        </template>
        <el-input-number v-model="form.cpu_usage" :max="100" :min="0" size="large" />
      </el-form-item>
      <el-form-item>
        <template #label><el-text>目标内存利用率 
          <el-tooltip placement="top">
            <template #content>
              当工作负载的所有Pod的平均内存利用率低于此目标时，<br>将减少Pod数量，直至最小实例数量；<br>
              当工作负载的所有Pod的平均内存利用率高于此目标时，<br>将增加Pod数量，直至最大实例数量
            </template>
            <el-icon class="text-yellow"><Warning /></el-icon>
          </el-tooltip></el-text>
        </template>
        <el-input-number v-model="form.mem_usage" :max="100" :min="0" size="large" />
      </el-form-item>          
    </el-form>
  </template>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show.add=false">取消</el-button>
      <el-button type="primary" @click="confirmScaler()">提交</el-button>
    </div>
  </template>
</el-drawer>
<el-dialog id="podLog" v-model="show.log" :fullscreen="true" @opened="openLog" @close="closeLog" :show-close="false" width="100%">
  <template #header="{close, titleId, titleClass}">
    <div class="my-header">
      <h4 :class="titleClass">容器日志</h4>
      <div class="box-tools pull-right" style="font-size:18px">
        <span class="card-header-btn" @click="followLogs(100)" v-if="!follow"><font-awesome-icon :icon="['far', 'circle-play']" /></span>
        <span class="card-header-btn" @click="followLogs(100)" v-else><font-awesome-icon :icon="['far', 'circle-stop']" /></span>
        <span class="card-header-btn" @click="getLogs(lines)"><el-icon><Refresh /></el-icon></span>
        <span class="card-header-btn" @click="show.log=false"><el-icon><Close /></el-icon></span>
      </div>
    </div>
  </template>
  <el-scrollbar style="height:calc(100vh - 90px)" ref="scrollRef">
    <pre class="dark fullscreen" ref="preRef">
      <p class="pointer" style="text-align: center" v-if="logs.length>=lines" @click="lines+=1000;getLogs(lines)">查看更多</p>
      <p v-for="(item,i) in logs" :key="i">{{ item }}</p>
    </pre>
  </el-scrollbar>
</el-dialog>
</template>

<script setup>
import { ArrowRight, Refresh, Close } from '@element-plus/icons-vue'
import { onBeforeMount, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRoute } from 'vue-router'
import { axios } from '/src/assets/util/axios'
import moment from 'moment'
import { v4 as uuidv4 } from 'uuid'
/* 引入v-ace-editor */
import { VAceEditor } from 'vue3-ace-editor'
import 'ace-builds/src-noconflict/mode-yaml'
import 'ace-builds/src-noconflict/theme-chrome'
import 'ace-builds/src-noconflict/ext-language_tools'
/* 变量定义 */
const route = useRoute()
const deploymentInfo = ref({})
const activeIndex = ref("status")
const podList = ref([])
const podEventList = ref([])
const eventList = ref([])
const currentPod = ref({})
const show = ref({
  event: false,
  yaml: false,
  log: false,
  autoscaler: false,
})
const yamlContent = ref("")
const form = ref({max_pod:10,min_pod:1,cpu_usage:50,mem_usage:50})
const timer = ref(null)
const logs = ref([])
const scrollRef = ref(null)
const preRef = ref(null)
const lines = ref(1000)
const pod_name = ref("")
const container_name = ref("")
const socket = ref(null)
const follow = ref(false)
const wsto = ref("")
const hpa = ref({found: false})
/* 生命周期函数 */
onBeforeMount(async () => {
  doRequest()
  timer.value = setInterval(() => {
    doRequest()
  }, 15000)
  getAutoScaler()
})
onBeforeUnmount(() => {
  if(timer) {
    clearInterval(timer.value)
    timer.value = null
  }
  if (socket.value) {
    socket.value.close()
  }
})
onMounted(() => {
  wsto.value = uuidv4()
  socket.value = new WebSocket(`/ws?id=${wsto.value}`)
  socket.value.onopen = () => {
    console.log('successfully connected to websocket server')
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
})
/* methods */
const doRequest = () => {
  getDeployment()
  getPods()
  getEvents()
}
const getDeployment = async () => {
  deploymentInfo.value = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/deployments/${route.params.workload_name}`)
  deploymentInfo.value.metadata.creationTimestamp = moment(deploymentInfo.value.metadata.creationTimestamp).format('YYYY-MM-DD HH:mm:ss')
  deploymentInfo.value.status.readyReplicas ||= 0
  for(let x of deploymentInfo.value.status.conditions) {
    if(x.type === "Available") {
      deploymentInfo.value.metadata.lastUpdateTime = x.lastUpdateTime
    }
  }
  deploymentInfo.value.metadata.lastUpdateTime = moment(deploymentInfo.value.metadata.lastUpdateTime).format('YYYY-MM-DD HH:mm:ss')
}
const handleSelect = async (key) => {
  activeIndex.value = key
}
const getPods = async () => {
  let response = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/deployments/${route.params.workload_name}/pods`)
  podList.value = response.map(x => {
    let conditionReady = x.status.conditions.find(n => {
      return n.type == 'Ready'
    })
    let m = {
      node_name: x.spec.nodeName,
      hostip: x.status.hostIP,
      podip: x.status.podIP,
      pod_name: x.metadata.name,
      status: {
        ready: conditionReady?.status||'False',
      },
    }
    let containerStatuses = x.status.containerStatuses?.map(y => {
      y.state_message = y.image
      y.status = Object.keys(y.state)[0]
      if(y.status!=='running') {
        y.state_message = y.state[y.status].reason
      } else if(y.ready === false) { // 即使为running，也可能ready=False，需要显示ready为False的reason
        y.state_message = conditionReady.reason
      } 
      return y
    }) || x.spec.containers.map(y => {  // 无法调度的pod没有containerStatuses字段
      return {
        name: y.name,
        state_message: x.status.conditions[0].reason,
        status: 'waiting'
      }
    })
    let initContainerStatuses = x.status.initContainerStatuses?.map(y => {
      y.state_message = y.image
      y.status = Object.keys(y.state)[0]
      if(y.status!=='running') {
        y.state_message = y.state[y.status].reason
      } else if(y.ready === false) { // 即使为running，也可能ready=False，需要显示ready为False的reason
        y.state_message = conditionReady.reason
      }
      return y
    }) || x.spec.initContainers?.map(y => {  // 无法调度的pod没有containerStatuses字段
      return {
        name: y.name,
        state_message: x.status.conditions[0].reason,
        status: 'waiting'
      }
    }) || []
    m.status.containerStatuses = containerStatuses.concat(initContainerStatuses)
    m.state = m.status.containerStatuses[0].status // running/waiting/terminated
    if(m.state === 'running' && conditionReady.status === 'True')
      m.state_message = `Created ${moment.duration(moment(m.status.containerStatuses[0].state.running.startedAt)-moment()).humanize(true)}`
    else
      m.state_message = m.status.containerStatuses[0].state_message
    return m
  })
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
const getEvents = async (row) => {
  let response = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/Deployment/${route.params.workload_name}/events`)
  eventList.value = response.map(x => {
    x.age = moment.duration(moment(x.lastTimestamp)-moment()).humanize(true)
    return x
  })
}
const handleCommand = async (command) => {
  switch(command) {
    case "restart": {
      await ElMessageBox.confirm('确定重启？','警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(async () => {
        await axios.patch(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/deployments/${route.params.workload_name}/rollout`)
        ElMessage.success({message: '操作成功'})
      }).catch(() =>{})
      break
    }
    case "autoscaler": {
      form.value = Object.assign({}, hpa.value)
      delete form.value.found
      show.value.autoscaler = true
      break
    }
    case "yaml": {
      yamlContent.value = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/deployments/${route.params.workload_name}/yaml`)
      show.value.yaml = true
      break
    }
    case "delete": {
      await ElMessageBox.confirm('确定删除？','警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(async () => {
        await axios.delete(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/deployments/${route.params.workload_name}`)
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
const getLogs = async (lines) => {
  let response = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/pods/${pod_name.value}/logs?container=${container_name.value}&lines=${lines}`)
  logs.value = response.split("\n")
  show.value.log = true
}
const followLogs = async (lines) => {
  logs.value = []
  try {
    await socket.value.send(JSON.stringify({
      "message_type":"LIZARDCD_WS_PODLOG",
      "message_to": wsto.value,
      "message_data": {
        "cluster": route.query.cluster,
        "namespace": route.query.namespace,
        "pod_name": pod_name.value,
        "container": container_name.value,
        "lines": lines,
        "follow": !follow.value
      }
    }))
  }
  catch(e) {
    console.warn(e)
  }
  follow.value = !follow.value
  if(!follow.value) {
    getLogs(1000)
  }
}
const openLog = () => {
  scrollRef.value.setScrollTop(preRef.value.clientHeight)
}
const closeLog = () => {
  if (socket.value) {
    socket.value.close()
  }
}
const confirmScaler = async () => {
  await axios.put(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/deployments/${route.params.workload_name}/hpa`, {
    "max": form.value.max_pod,
    "min": form.value.min_pod,
    "cpu": form.value.cpu_usage,
    "memory": form.value.mem_usage
  })
  show.value.autoscaler = false
}
const getAutoScaler = async () => {
  let response = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/deployments/${route.params.workload_name}/hpa`)
  if(response.metadata) {
    hpa.value = {
      found: true,
      max_pod: response.spec.maxReplicas,
      min_pod: response.spec.minReplicas,
      cpu_usage: response.spec.metrics.find(n => n.resource.name === 'cpu')?.resource.target.averageUtilization,
      mem_usage: response.spec.metrics.find(n => n.resource.name === 'memory')?.resource.target.averageUtilization
    }
  }
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
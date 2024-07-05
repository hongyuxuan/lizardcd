<template>
<el-breadcrumb :separator-icon="ArrowRight" style="min-width:1366px">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>工作负载</el-breadcrumb-item>
  <el-breadcrumb-item :to="{ path: '/workload/statefulsets', query: {cluster: route.query.cluster, namespace: route.query.namespace} }">有状态副本集</el-breadcrumb-item>
  <el-breadcrumb-item>{{ statefulsetInfo.metadata?.name }}</el-breadcrumb-item>
</el-breadcrumb>
<el-row :gutter="15" style="min-width:1366px">
  <el-col :span="6">
    <el-card>
      <template #header>
        <div class="card-header">
          <span class="card-header-text"><b>{{ statefulsetInfo.metadata?.name }}</b></span>
          <div class="box-tools pull-right">
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
            <div v-if="scope.row.state==='waiting'" class="text-yellow">
              <el-icon><WarningFilled /></el-icon>
              {{scope.row.state_message}}
            </div>
            <div v-else-if="scope.row.state==='terminated'" class="text-red">
              <el-icon><WarningFilled /></el-icon>
              {{scope.row.state}}
            </div>
            <div v-else class="text-gray">{{scope.row.state_message}}</div>
          </template>
        </el-table-column>
        <el-table-column>
          <template #default="scope">
            <div>{{scope.row.node_name}} ( {{scope.row.hostip}} )</div>
            <div class="text-gray">Worker Node</div>
          </template>
        </el-table-column>
        <el-table-column>
          <template #default="scope">
            <div>{{scope.row.podip}}</div>
            <div class="text-gray">Pod IP</div>
          </template>
        </el-table-column>
        <el-table-column width="100">
          <template #default="scope">
            <el-link type="primary" :underline="false" @click="getPodEvents(scope.row)">查看</el-link>
            <div class="text-gray">Events</div>
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
                    <div><b>{{props.row.name}}</b></div>
                    <div class="text-gray" v-if="props.row.ready==true">{{props.row.image}}</div>
                    <div class="text-gray" v-else>{{props.row.state_message}}</div>
                  </template>
                </el-table-column>
                <el-table-column min-width="15%">
                  <template #default="props">
                    <div>{{props.row.status}}</div>
                    <div class="text-gray">Status</div>
                  </template>
                </el-table-column>
                <el-table-column min-width="15%">
                  <template #default="props">
                    <div>{{props.row.restartCount}}</div>
                    <div class="text-gray">Restart Count</div>
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
          <el-descriptions-item v-for="(v,k,i) in statefulsetInfo.metadata?.labels" :index="i" :label="k">{{ v }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </div>
    <div class="box box-item" v-show="activeIndex==='annotations'">
      <div class="box-body" style="padding-top:20px">
        <el-descriptions :column="1">
          <el-descriptions-item v-for="(v,k,i) in statefulsetInfo.metadata?.annotations" :index="i" :label="k">{{ v }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </div>
    <div class="box box-item" v-show="activeIndex==='events'">
      <div class="box-body" style="padding-top:20px">
        <el-table :data="eventList" style="width:100%">
          <el-table-column prop="type" label="Type" />
          <el-table-column prop="reason" label="Reason" />
          <el-table-column prop="age" label="Age" />
          <el-table-column prop="source.component" label="From" />
          <el-table-column prop="message" label="Message" min-width="300px" />
        </el-table>
      </div>
    </div>
  </el-col>
</el-row>
<el-dialog
  v-model="show.event"
  :title="`${currentPod.pod_name} 事件`"
  width="70%">
  <el-table :data="podEventList" style="width:100%">
    <el-table-column prop="type" label="Type" />
    <el-table-column prop="reason" label="Reason" />
    <el-table-column prop="age" label="Age" />
    <el-table-column prop="source.component" label="From" />
    <el-table-column prop="message" label="Message" min-width="300px" />
  </el-table>
</el-dialog>
<el-drawer v-model="show.yaml" direction="rtl" size="800px">
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
</template>

<script setup>
import { ArrowRight } from '@element-plus/icons-vue'
import { onBeforeMount, onBeforeUnmount, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRoute } from 'vue-router'
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
const podList = ref([])
const podEventList = ref([])
const eventList = ref([])
const currentPod = ref({})
const show = ref({
  event: false,
  yaml: false,
})
const yamlContent = ref("")
const timer = ref(null)
/* 生命周期函数 */
onBeforeMount(async () => {
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
  getPods()
  getEvents()
}
const getStatefulset = async () => {
  statefulsetInfo.value = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/statefulsets/${route.params.workload_name}`)
  statefulsetInfo.value.metadata.creationTimestamp = moment(statefulsetInfo.value.metadata.creationTimestamp).format('YYYY-MM-DD HH:mm:ss')
  statefulsetInfo.value.status.readyReplicas ||= 0
}
const handleSelect = async (key) => {
  activeIndex.value = key
}
const getPods = async () => {
  let response = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/statefulsets/${route.params.workload_name}/pods`)
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
  let response = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/StatefulSet/${route.params.workload_name}/events`)
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
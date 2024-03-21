<template>
<el-breadcrumb :separator-icon="ArrowRight">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>工作负载</el-breadcrumb-item>
  <el-breadcrumb-item>部署</el-breadcrumb-item>
</el-breadcrumb>
<el-card>
  <template #header>
    <div class="card-header">
      <span class="card-header-text">部署</span>
      <span class="pull-right pointer" @click="getList(true)"><el-icon><Refresh /></el-icon></span>
    </div>
  </template>
  <el-row>
    <el-col :span="12">
      <el-button-group>
        <el-select v-model="cluster" placeholder="请选择集群" clearable filterable style="width:200px;margin-right:8px" size="large">
          <el-option v-for="(v,k) in clusterList" :key="k" :label="k" :value="k" />
        </el-select>
        <el-select v-model="namespace" placeholder="请选择命名空间" clearable filterable @change="getList(true)" style="width:200px;margin-right:8px" size="large">
          <el-option v-for="(item) in clusterList[cluster]" :key="item" :label="item" :value="item" />
        </el-select>
        <el-input v-model="searchKey" placeholder="输入名称进行搜索" size="large" :prefix-icon="Search" @change="getPage(1)" clearable style="width:200px;" />
      </el-button-group>
    </el-col>
    <el-col :span="12">
      <el-button class="pull-right" size="large" type="primary" @click="show.add=true;form={}">新建发布</el-button>
    </el-col>
  </el-row>
  <el-table 
    :data="list" 
    :show-header="false" 
    v-loading="loading"
    element-loading-text="奋力加载中..."
    class="line-height40" 
    style="width:100%;margin-top:10px;min-height:150px">
    <el-table-column label="" width="45">
      <font-awesome-icon icon="layer-group" style="font-size:25px;vertical-align:middle;" />
    </el-table-column>
    <el-table-column prop="name" label="名称" min-width="200">
      <template #default="scope">
        <el-link :underline="false" @click="openDeploy(scope.row)">{{ scope.row.name }}</el-link>
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
    <el-table-column label="操作" width="120">
      <template #default="scope">
        <el-dropdown @command="handleCommand" style="vertical-align:middle;">
          <el-button>更多操作<el-icon class="el-icon--right"><arrow-down /></el-icon></el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item :command="{action:'restart',row:scope.row}">重启</el-dropdown-item>
              <el-dropdown-item :command="{action:'setImage',row:scope.row}">设置镜像</el-dropdown-item>
              <el-dropdown-item :command="{action:'scale',row:scope.row}">设置副本</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </template>
    </el-table-column>
  </el-table>
  <el-pagination 
      class="pull-right"
      background 
      v-model:page-size="pageSize"
      :page-sizes="[20, 30, 50, 100]"
      layout="total, sizes, prev, pager, next, jumper" 
      :total="pageTotal"
      @current-change="getPage"
      v-model:current-page="current" />
</el-card>
<el-dialog v-model="show.pods" :title="currentDeploy.name" width="80%">
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
            {{scope.row.state_message}}
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
          <el-link type="primary" :underline="false" @click="getEvents(scope.row)">查看</el-link>
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
                  <font-awesome-icon icon="circle" v-if="props.row.ready==false" class="twinkling podstatus text-yellow" />
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
</el-dialog>
<el-dialog
  v-model="show.event"
  :title="`${currentPod.pod_name} 事件`"
  width="70%">
  <el-table :data="eventList" style="width:100%">
    <el-table-column prop="type" label="Type" />
    <el-table-column prop="reason" label="Reason" />
    <el-table-column prop="age" label="Age" />
    <el-table-column prop="source.component" label="From" />
    <el-table-column prop="message" label="Message" min-width="300px" />
  </el-table>
</el-dialog>
<el-drawer v-model="show.add" direction="rtl" size="600px">
  <template #header>
    <h4>新建发布</h4>
  </template>
  <template #default>
    <el-form ref="task" :model="form" :rules="rules" label-width="120px">
      <el-form-item label="集群" prop="cluster">
        <el-select v-model="form.cluster" placeholder="请选择集群" clearable filterable style="width:100%;margin-right:8px" size="large">
          <el-option v-for="(v,k) in clusterList" :key="k" :label="k" :value="k" />
        </el-select>
      </el-form-item>
      <el-form-item label="命名空间" prop="namespace">
        <el-select v-model="form.namespace" placeholder="请选择命名空间" clearable filterable style="width:100%" size="large">
          <el-option v-for="(item) in clusterList[form.cluster]" :key="item" :label="item" :value="item" />
        </el-select>
      </el-form-item>
      <el-form-item label="工作负载名称" prop="workload_name">
        <el-input v-model="form.workload_name" />
      </el-form-item>
      <el-form-item label="容器名称" prop="container">
        <el-input v-model="form.container" />
      </el-form-item>
      <el-form-item label="容器镜像" prop="image">
        <el-input v-model="form.image" />
      </el-form-item>
    </el-form>
  </template>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show.add=false">取消</el-button>
      <el-button type="primary" @click="confirmClick(task)">提交</el-button>
    </div>
  </template>
</el-drawer>
</template>

<script setup>
import { ArrowRight,Search,Refresh } from '@element-plus/icons-vue'
import { onBeforeMount, onBeforeUnmount, ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { axios } from '/src/assets/util/axios'
import moment from 'moment'
import _ from 'lodash'
/* 变量定义 */
const all = ref([])
const list = ref([])
const pageSize = ref(10)
const pageTotal = ref(0)
const current = ref(1)
const cluster = ref("")
const clusterList = ref({})
const namespace = ref("")
const show = ref({
  pods: false,
  event: false,
  add: false
})
const currentDeploy = ref({})
const currentPod = ref({})
const podList = ref([])
const eventList = ref([])
const loading = ref(false)
const form = ref({})
const task = ref(null)
const rules = reactive({
  cluster: [{required: true, message: '请选择集群', trigger: 'change'}],
  namespace: [{required: true, message: '请选择命名空间', trigger: 'change'}],
  workload_name: [{required: true, message: '请填写工作负载名称', trigger: 'blur'}],
  container: [{required: true, message: '请填写容器名称', trigger: 'blur'}],
  image: [{required: true, message: '请填写镜像', trigger: 'blur'}],
})
const searchKey = ref("")
const timer = ref(null)
/* 生命周期函数 */
onBeforeMount(async () => {
  getClusterList()
  timer.value = setInterval(() => {
    getList(false)
  }, 15000)
})
onBeforeUnmount(() => {
  if(timer) {
    clearInterval(timer.value)
    timer.value = null
  }
})
/* methods */
const getClusterList = async () => {
  clusterList.value = await axios.get(`/consul/clusters`)
}
const getList = async (ifLoading) => {
  if(ifLoading) loading.value = true
  if(cluster.value !== "" && namespace.value !== "") {
    let response = await axios.get(`/lizardcd/cluster/${cluster.value}/namespace/${namespace.value}/deployments`)
    response = response.map(x => {
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
    all.value = _.sortBy(response, 'lastUpdateTime').reverse()
    getPage(current.value)
  }
  if(ifLoading) loading.value = false
}
const getPage = async (page) => {
  let tmpList = all.value
  if(searchKey.value !== '') {
    tmpList = all.value.filter(n => n.name.includes(searchKey.value))
  }
  pageTotal.value = tmpList.length
  list.value = tmpList.slice((page-1)*pageSize.value, page*pageSize.value)
}
const openDeploy = async (row) => {
  currentDeploy.value = row
  let response = await axios.get(`/lizardcd/cluster/${cluster.value}/namespace/${namespace.value}/deployments/${row.name}/pods`)
  podList.value = response.map(x => {
    let status = x.status.conditions.find(n => {
      return n.type == 'Ready'
    })
    let m = {
      node_name: x.spec.nodeName,
      hostip: x.status.hostIP,
      podip: x.status.podIP,
      pod_name: x.metadata.name,
      status: {
        ready: status?.status||'False',
      },
    }
    let containerStatuses = x.status.containerStatuses?.map(y => {
      y.state_message = y.image
      y.status = Object.keys(y.state)[0]
      if(y.status!=='running') y.state_message = y.state[y.status].reason
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
      if(y.status!=='running') y.state_message = y.state[y.status].reason
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
    if(m.state === 'running')
      m.state_message = `Created ${moment.duration(moment(m.status.containerStatuses[0].state.running.startedAt)-moment()).humanize(true)}`
    else
      m.state_message = m.status.containerStatuses[0].state_message
    return m
  })
  show.value.pods = true
}
const getEvents = async (row) => {
  currentPod.value = row
  let response = await axios.get(`/lizardcd/cluster/${cluster.value}/namespace/${namespace.value}/pods/${row.pod_name}/events`)
  eventList.value = response.map(x => {
    x.age = moment.duration(moment(x.lastTimestamp)-moment()).humanize(true)
    return x
  })
  show.value.event = true
}
const handleCommand = async (command) => {
  switch(command.action) {
    case "restart": {
      await ElMessageBox.confirm('确定重启？','警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(async () => {
        await axios.patch(`/lizardcd/cluster/${cluster.value}/namespace/${namespace.value}/deployments/${command.row.name}/rollout`)
        ElMessage.success({message: '操作成功'})
      }).catch(() =>{})
      break
    }
    case "setImage": {
      let pods = await axios.get(`/lizardcd/cluster/${cluster.value}/namespace/${namespace.value}/deployments/${command.row.name}/pods`)
      let container = pods[0].spec.containers[0].name
      await ElMessageBox.prompt(`请输入容器 ${container} 的镜像`,'提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
      }).then(async ({value}) => {
        await axios.patch(`/lizardcd/cluster/${cluster.value}/namespace/${namespace.value}/deployments/${command.row.name}?container=${container}&image=${value}`)
        ElMessage.success({message: '操作成功'})
      }).catch(() =>{})
      break
    }
    case "scale": {
      await ElMessageBox.prompt(`请输入副本数`,'提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        inputPattern: /\d+/,
        inputErrorMessage: '必须输入数字'
      }).then(async ({value}) => {
        await axios.patch(`/lizardcd/cluster/${cluster.value}/namespace/${namespace.value}/deployments/scale`, {
          workloads: [
            {
              name: command.row.name,
              replicas: parseInt(value)
            }  
          ]
        })
        ElMessage.success({message: '操作成功'})
      }).catch(() =>{})
      break
    }
  }
  setTimeout(async () => {
    await getList(true)
  }, 2000)
}
const confirmClick = async (f) => {
  if(!f) return
  await f.validate(async (valid) => {
    if(valid) {
      await axios.patch(`/lizardcd/cluster/${form.value.cluster}/namespace/${form.value.namespace}/deployments/${form.value.workload_name}?container=${form.value.container}&image=${form.value.image}`)
      ElMessage.success({message: '发布成功'})
      show.value.add = false
    }
    else {
      ElMessage.warning('必填项未填完')
    }
  })
}
</script>
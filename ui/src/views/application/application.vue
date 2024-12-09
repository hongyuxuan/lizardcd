<template>
<el-breadcrumb :separator-icon="ArrowRight">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>应用管理</el-breadcrumb-item>
</el-breadcrumb>
<el-card>
  <template #header>
    <div class="card-header">
      <span class="card-header-text">应用管理</span>
    </div>
  </template>
  <el-row>
    <el-col :span="18">
      <el-button-group style="width:100%">
        <el-button :icon="Refresh" size="large" style="margin-right:5px" @click="getList(current)" />
        <el-input v-model="searchKey" clearable placeholder="输入应用名查询……" :prefix-icon="Search" @change="getList(1);current=1" style="width:25%;margin-right:5px" size="large" />
        <el-select 
          v-model="searchTags"
          multiple 
          clearable 
          allow-create 
          default-first-option 
          :reserve-keyword="false" 
          filterable 
          style="width:40%;margin-right:5px" 
          placeholder="输入标签过滤……" 
          @change="getList(1);current=1"
          size="large">
          <el-option v-for="item in tagOptions" :key="item" :label="item" :value="item" />
        </el-select>
      </el-button-group>
    </el-col>
    <el-col :span="6">
      <el-button-group class="pull-right">
        <el-button class="pull-right" size="large" type="primary" @click="edit=false;applicationInfo={};refAdd.open()">+ 新建应用</el-button>
        <el-button class="pull-right" size="large" type="primary" @click="show.deploy=true;formDeploy={policy:'same'}" style="margin-right:5px">发布应用</el-button>
      </el-button-group>
    </el-col>
  </el-row>
  <el-table 
    :data="list" 
    v-loading="loading.table"
    element-loading-text="奋力加载中..."
    class="line-height40" 
    style="width:100%;margin-top:10px">
    <el-table-column type="selection" width="45" />
    <el-table-column prop="app_name" label="应用名称" min-width="200">
      <template #default="scope">
        <el-link :underline="false" :href="`/application/${scope.row.id}`">{{ scope.row.app_name }}</el-link>
        <span v-if="syncList.hasOwnProperty(scope.row.app_name)" style="margin-left:5px">
          <el-popover placement="right-start" :width="300" trigger="hover">
            <template #reference>
              <el-icon v-if="syncList[scope.row.app_name].status==='已同步'" class="text-success" style="vertical-align:middle;"><CircleCheckFilled /></el-icon>
              <el-icon v-else class="text-red" style="vertical-align:middle;"><CircleCloseFilled /></el-icon>
            </template>
            <el-descriptions :column="1" size="small">
              <el-descriptions-item label="同步状态">
                <span v-if="syncList[scope.row.app_name].status==='已同步'" class="text-success">{{ syncList[scope.row.app_name].status }}</span>
                <span v-else class="text-red">{{ syncList[scope.row.app_name].status }}</span>
                </el-descriptions-item>
              <el-descriptions-item label="原因" v-if="syncList[scope.row.app_name].status!=='已同步'">{{ syncList[scope.row.app_name].reason }}</el-descriptions-item>
              <el-descriptions-item label="应用commitId">{{ syncList[scope.row.app_name].local_commitId }}</el-descriptions-item>
              <el-descriptions-item label="仓库commitId">{{ syncList[scope.row.app_name].remote_commitId }}</el-descriptions-item>
              <el-descriptions-item label="Last Comment">{{ syncList[scope.row.app_name].comment }}</el-descriptions-item>
              <el-descriptions-item label="Committer">{{ syncList[scope.row.app_name].author }}</el-descriptions-item>
            </el-descriptions>
          </el-popover>
        </span>
        <span v-if="taskList.hasOwnProperty(scope.row.app_name)" style="margin-left:20px">
          <el-tooltip v-for="(item,i) in taskList[scope.row.app_name]" :key="i" placement="top">
            <template #content>
              <span>FINISH_AT: {{ item.finish_at || "" }}</span>
            </template>
            <el-link  :underline="false" :href="`/task/history?id=${item.id}`" target="_blank" style="font-size:12px">
              <font-awesome-icon icon="circle" v-if="['finished','initialize','terminated'].includes(item.status)&&item.success===true" class="text-success" />
              <font-awesome-icon icon="circle" v-else-if="['finished','initialize','terminated'].includes(item.status)&&item.success===false" class="text-red" />
              <font-awesome-icon icon="circle" v-else class="text-yellow twinkling" />
            </el-link>
          </el-tooltip>
        </span>
      </template>
    </el-table-column>
    <el-table-column prop="deploy_type" label="部署方式" width="120" />
    <el-table-column prop="tags" label="标签" min-width="300">
      <template #default="scope">
        <el-tag v-for="item in scope.row.tags" :key="item" size="large">{{item}}</el-tag>
      </template>
    </el-table-column>
    <el-table-column prop="tenant" label="所属租户" min-width="80" />
    <el-table-column prop="update_at" label="更新时间" width="160">
      <template #default="scope">
        {{ moment(scope.row.update_at).format('YYYY-MM-DD HH:mm') }}
      </template>
    </el-table-column>
    <el-table-column prop="Option" label="操作" width="120">
      <template #default="scope">
        <el-dropdown @command="handleCommand" style="vertical-align:middle;">
          <el-button>更多操作<el-icon class="el-icon--right"><arrow-down /></el-icon></el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item :command="{action:'deploy',row:scope.row}" v-if="scope.row.deploy_type!=='GitOps'">发布</el-dropdown-item>
              <el-dropdown-item :command="{action:'sync',row:scope.row}" v-else>同步</el-dropdown-item>
              <el-dropdown-item :command="{action:'copy',row:scope.row}">复制</el-dropdown-item>
              <el-dropdown-item :command="{action:'edit',row:scope.row}">编辑</el-dropdown-item>
              <el-dropdown-item :command="{action:'restart',row:scope.row}" v-if="scope.row.deploy_type!=='GitOps'">重启</el-dropdown-item>
              <el-dropdown-item :command="{action:'delete',row:scope.row}">删除</el-dropdown-item>
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
    :page-sizes="[10, 30, 50, 100]"
    layout="total, sizes, prev, pager, next, jumper" 
    :total="pageTotal"
    @current-change="getList"
    @size-change="handleSizeChange"
    v-model:current-page="current" />
</el-card>
<addForm ref="refAdd" :form="applicationInfo" :k8scluster="k8scluster" :edit="edit" @submit="getList(current)" />
<deployForm ref="refDeploy" :applicationInfo="applicationInfo" />
</template>

<script setup>
import { ArrowRight,Search,Refresh } from '@element-plus/icons-vue'
import { onBeforeMount, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import addForm from './add.vue'
import deployForm from './deploy.vue'
import { axios } from '/src/assets/util/axios'
import moment from 'moment'
import _ from 'lodash'
/* 变量定义 */
const router = useRouter()
const route = useRoute()
const list = ref([])
const pageSize = ref(10)
const pageTotal = ref(0)
const current = ref(1)
const searchTags = ref([])
const searchKey = ref("")
const tagOptions = ref([])
const show = ref({
  deploy: false
})
const edit = ref(false)
const k8scluster = ref({})
const loading = ref({
  searchapp: false,
  table: false,
  artifact: false,
})
const formDeploy = ref({})
const applicationInfo = ref({})
const taskList = ref({})
const syncList = ref({})
const refAdd = ref(null)
const refDeploy = ref(null)
/* 生命周期函数 */
onBeforeMount(async () => {
  if(route.query.app_name) {
    searchKey.value = route.query.app_name
  }
  getClusterList()
  getList(1)
});
/* methods */
const getList = async (page) => {
  let url = `page=${page}&size=${pageSize.value}&sort=application.update_at desc`
  let tags = searchTags.value.map(x => `tags=="${x}"`)
  if(searchKey.value !== "") {
    url += `&search=app_name==${encodeURIComponent(searchKey.value)}`
    if(tags.length > 0) {
      url += `,${tags.join(",")}`
    }
  } else if(tags.length > 0) {
    url += `&search=${tags.join(",")}`
  }
  loading.value.table = true
  let response = await axios.get(`/lizardcd/db/application?${url}`)
  loading.value.table = false
  list.value = response.results||[]
  for(let x of list.value) {
    tagOptions.value = tagOptions.value.concat(x.tags)
  }
  tagOptions.value = _.uniq(tagOptions.value)
  pageTotal.value = response.total

  getApplicationTask()
  getApplicationSyncStatus()
}
const getApplicationTask = async () => {
  let response = await axios.get(`/lizardcd/task/history_by_application?size=5&sort=init_at%20desc&apps=${list.value.map(x => encodeURIComponent(x.app_name)).join(',')}`)
  for(let [k,v] of Object.entries(response)) {
    if(v) {
      taskList.value[k] = v.map(y => {
        return {
          id: y.id,
          status: y.status,
          success: y.success.Bool,
          finish_at: y.finish_at.Valid === true ? moment(y.finish_at.Time).format('YYYY-MM-DD HH:mm:ss') : undefined,
        }
      })
    }
  }
}
const getApplicationSyncStatus = async () => {
  syncList.value = {}
  syncList.value = await axios.get(`/lizardcd/git/sync/status?apps=${list.value.map(x => encodeURIComponent(x.app_name)).join(',')}`)
}
const getClusterList = async () => {
  k8scluster.value = await axios.get(`/lizardcd/server/clusters`)
}
const handleCommand = async (command) => {
  applicationInfo.value = Object.assign({}, command.row)
  switch(command.action) {
    case "deploy": {
      refDeploy.value.open()
      break
    }
    case "sync": {
      let response = await axios.get(`/lizardcd/db/application_resource?filter=application_id==${command.row.id}`)
      let txt = "确定同步以下资源？<br>"
      for(let x of response) {
        txt += `${x.cluster} / ${x.namespace} / ${x.resource_type} / ${x.resource_name}<br>`
      }
      await ElMessageBox.confirm(txt, '警告' , {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
        dangerouslyUseHTMLString: true,
      }).then(async () => {
        await axios.post(`/lizardcd/task/run`, {
          "app_name": command.row.app_name,
          "task_type": "synchronize",
          "trigger_type": "手动触发",
          "waiting": false
        })
      }).catch(() =>{})
      break
    }
    case "copy": {
      edit.value = false
      refAdd.value.open()
      break
    }
    case "edit": {
      edit.value = true
      refAdd.value.open()
      break
    }
    case "restart": {
      await ElMessageBox.confirm(`确定重启该应用的 ${command.row.workload.length} 个工作负载？`,'警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(async () => {
        let params = {
          "app_name": command.row.app_name,
          "task_type": "rollout",
          "trigger_type": "手动触发",
          "workloads": command.row.workload.map(x => {
            return {
              "cluster": x.cluster,
              "namespace": x.namespace,
              "workload_type": x.workload_type,
              "workload_name": x.workload_name,
            }
          })
        }
        let response = await axios.post(`/lizardcd/task/run`, params)
        router.push(`/task/history?id=${response.id}`)
      }).catch(() =>{})
      break
    }
    case "delete": {
      await ElMessageBox.confirm('确定删除？','警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(async () => {
        ElMessage.warning({message: `正在删除，请稍后`})
        await axios.delete(`/lizardcd/db/application/${command.row.id}`)
        getList(current.value)
      }).catch(() =>{})
      break
    }
  }
}
const handleSizeChange = async (size) => {
  pageSize.value = size
  await getList(current.value)
}
</script>
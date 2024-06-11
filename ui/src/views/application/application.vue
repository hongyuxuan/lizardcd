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
    <el-col :span="12">
      <el-button-group>
        <el-button :icon="Refresh" size="large" style="margin-right:5px" @click="getList(current)" />
        <el-input v-model="searchKey" placeholder="输入名称进行搜索" size="large" :prefix-icon="Search" @change="getList(1);current=1" clearable style="width:300px;" />
      </el-button-group>
    </el-col>
    <el-col :span="12">
      <el-button-group class="pull-right">
        <el-button class="pull-right" size="large" type="primary" @click="show.add=true;edit=false;form={workload:[],traffic_policy:'weight',enable_traffic_control:false,tenant:tenant,tags:[]}">新建应用</el-button>
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
    <el-table-column prop="app_name" label="应用名称" min-width="200" />
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
              <el-dropdown-item :command="{action:'deploy',row:scope.row}">发布</el-dropdown-item>
              <el-dropdown-item :command="{action:'copy',row:scope.row}">复制</el-dropdown-item>
              <el-dropdown-item :command="{action:'edit',row:scope.row}">编辑</el-dropdown-item>
              <el-dropdown-item :command="{action:'restart',row:scope.row}">重启</el-dropdown-item>
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
    v-model:current-page="current" />
</el-card>
<el-drawer v-model="show.add" direction="rtl" size="650px">
  <template #header>
    <h4 v-if="edit===false">新增应用</h4>
    <h4 v-if="edit===true">编辑应用</h4>
  </template>
  <template #default>
    <el-form ref="app" :model="form" :rules="rules" label-width="120px">
      <el-form-item label="应用名称" prop="app_name" >
        <el-input v-model="form.app_name" size="large" />
      </el-form-item>
      <el-form-item label="选择镜像仓库" prop="repo">
        <el-select v-model="form.repo" placeholder="请选择" value-key="id" clearable size="large" style="width:100%">
          <el-option v-for="item in repoList" :key="item.id" :label="item.repo_url" :value="item">
            <span style="float:left">{{item.repo_url}}</span>
            <span style="float:right;color:var(--el-text-color-secondary)">{{item.repo_account}}</span>
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="仓库/项目" prop="repo_name">
        <el-input v-model="form.repo_name" placeholder="请填写" size="large" />
        <el-alert type="info">Artifactory填写仓库名<br>Harbor填写项目名<br>DockerHub填写namespace</el-alert>
      </el-form-item>
      <el-form-item label="镜像名" prop="image_name">
        <el-input v-model="form.image_name" placeholder="请填写" size="large" />
      </el-form-item>
      <el-form-item label="所属租户" prop="tenant">
        <el-input v-model="form.tenant" disabled size="large" />
      </el-form-item>
      <el-form-item label="设置标签" prop="tags">
        <el-tag v-for="item in form.tags" :key="item" closable :disable-transitions="false" @close="handleClose(item)" size="large">{{item}}</el-tag>
        <el-input v-if="inputVisible" ref="inputRef" v-model="inputLabel" @keyup.enter="handleInputConfirm" @blur="handleInputConfirm" style="width:100px" />
        <el-button v-else @click="showInput">+ 添加标签</el-button>
      </el-form-item>
      <el-form-item label="开启灰度发布">
        <el-switch v-model="form.enable_traffic_control" />
      </el-form-item>
      <el-form-item label="灰度发布策略" v-if="form.enable_traffic_control">
        <el-radio-group v-model="form.traffic_policy">
          <el-radio value="weight" size="large">基于权重</el-radio>
          <el-radio value="header" size="large">基于头部字段</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="工作负载">
        <el-card v-for="(m,index) in form.workload" :key="index" style="width:100%">
          <template #header>
            <div class="card-header">
              <span>工作负载 {{ index+1 }}</span>
              <div class="box-tools pull-right">
                <span class="card-header-btn" @click="copyWorkload(index)"><el-icon><CopyDocument /></el-icon></span>
                <span class="card-header-btn" @click="removeWorkload(index)"><el-icon><Close /></el-icon></span>
              </div>
            </div>
          </template>
          <el-form label-width="100px">
            <el-form-item label="容器集群">
              <el-select v-model="m.cluster" clearable placeholder="请选择" size="large" style="width:100%">
                <el-option v-for="(v,k,i) in k8scluster" :key="i" :label="k" :value="k" />
              </el-select>
            </el-form-item>
            <el-form-item label="命名空间">
              <el-select v-model="m.namespace" clearable placeholder="请选择" size="large" style="width:100%">
                <el-option v-for="item in k8scluster[m.cluster]" :key="item" :label="item" :value="item" />
              </el-select>
            </el-form-item>
            <el-form-item label="工作负载类型">
              <el-radio-group v-model="m.workload_type">
                <el-radio label="deployments" value="deployments" />
                <el-radio label="statefulsets" value="statefulsets" />
              </el-radio-group>
            </el-form-item>
            <el-form-item label="工作负载名称">
              <el-input v-model="m.workload_name" size="large" />
            </el-form-item>
            <el-form-item label="容器名称">
              <el-input v-model="m.container_name" size="large" />
            </el-form-item>
            <el-form-item label="版本号" v-if="form.enable_traffic_control">
              <el-input v-model="m.version" size="large" />
              <el-alert title="版本号必须和POD的labels.version一致" />
            </el-form-item>
            <el-form-item label="权重" v-if="form.enable_traffic_control&&form.traffic_policy==='weight'">
              <el-input-number v-model="m.weight" :max="100" :min="0" size="large" />
            </el-form-item>
            <el-form-item label="匹配头部字段" v-if="form.enable_traffic_control&&form.traffic_policy==='header'">
              <!-- <el-input v-model="m.match_headers" size="large" placeholder="每行填写一个header，如userid:xxxx" type="textarea" :autosize="{minRows:2}" /> -->
              <table class="table table-bordered">
                <thead><tr><th width="30%">头部键</th><th width="40%">匹配方式</th><th width="30%">头部值</th></tr></thead>
                <tbody>
                  <tr v-for="(n,i) in m.headers" :key="i">
                    <td><el-input v-model="n.key" /></td>
                    <td>
                      <el-select v-model="n.match_type">
                        <el-option label="exact" value="exact" />
                        <el-option label="prefix" value="prefix" />
                        <el-option label="regex" value="regex" />
                      </el-select>
                    </td>
                    <td><el-input v-model="n.value" /></td>
                  </tr>
                </tbody>
              </table>
              <el-row>
                <el-button icon="Plus" circle @click="addMatchHeader(m.headers)"></el-button>
                <el-button icon="Close" circle @click="removeMatchHeader(m.headers)"></el-button>
              </el-row>
            </el-form-item>
          </el-form>
        </el-card>
        <el-row>
          <el-button icon="Plus" circle @click="addWorkload"></el-button>
        </el-row>
      </el-form-item>
    </el-form>
  </template>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show.add=false" size="large">取消</el-button>
      <el-button type="primary" @click="confirmClick(app)" size="large">提交</el-button>
    </div>
  </template>
</el-drawer>
<el-drawer v-model="show.deploy" direction="rtl" size="650px">
  <template #header>
    <h4>应用发布</h4>
  </template>
  <template #default>
    <el-form ref="release" :model="formDeploy" label-width="100px">
      <el-form-item label="选择应用">
        <el-select 
          v-model="formDeploy.app_name"
          placeholder="请选择" 
          value-key="id" 
          clearable 
          filterable 
          remote 
          :remote-method="searchApp" 
          :loading="loading.searchapp" 
          style="width:100%" 
          size="large">
          <el-option v-for="item in appList" :key="item.id" :label="item.app_name" :value="item" />
        </el-select>
      </el-form-item>
      <el-form-item label="发布策略">
        <el-radio-group v-model="formDeploy.policy">
          <el-radio value="same" size="large">所有工作负载使用相同镜像</el-radio>
          <el-radio value="different" size="large">不同工作负载使用不同镜像</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="选择制品" v-if="formDeploy.policy==='same'">
        <el-select 
          v-model="formDeploy.artifact_url" 
          placeholder="请选择" 
          value-key="tag" 
          filterable
          clearable
          allow-create
          default-first-option
          reserve-keyword
          remote
          :remote-method="getArtifacts"
          :loading="loading.artifact"
          style="width:100%" 
          size="large">
          <el-option v-for="item in artifactList" :key="item.tag" :label="item.tag" :value="item">
            <span style="float:left">{{item.tag}}</span>
            <span style="float:right;color:var(--el-text-color-secondary);font-size:12px">{{item.last_modified}}</span>
          </el-option>
        </el-select>
        <el-alert title="如因网络问题，无法在线获取DockerHub的镜像，可直接将镜像地址填写于此" />
      </el-form-item>
      <el-form-item label="工作负载" v-if="formDeploy.policy==='different'">
        <el-card v-for="(m,index) in formDeploy.app_name.workload" :key="index" style="width:100%">
          <template #header>
            <div class="card-header">
              <span>工作负载 {{ index+1 }}</span>
              <div class="box-tools pull-right">
                <span class="card-header-btn" @click="copyWorkload(index)"><el-icon><CopyDocument /></el-icon></span>
                <span class="card-header-btn" @click="removeWorkload(index)"><el-icon><Close /></el-icon></span>
              </div>
            </div>
          </template>
          <el-form label-width="100px">
            <el-form-item label="容器集群">
              <el-input v-model="m.cluster" size="large" disabled />
            </el-form-item>
            <el-form-item label="命名空间">
              <el-input v-model="m.namespace" size="large" disabled />
            </el-form-item>
            <el-form-item label="工作负载类型">
              <el-radio-group v-model="m.workload_type" disabled>
                <el-radio label="deployments" value="deployments" />
                <el-radio label="statefulsets" value="statefulsets" />
              </el-radio-group>
            </el-form-item>
            <el-form-item label="工作负载名称">
              <el-input v-model="m.workload_name" size="large" disabled />
            </el-form-item>
            <el-form-item label="容器名称">
              <el-input v-model="m.container_name" size="large" disabled />
            </el-form-item>
            <el-form-item label="选择制品">
              <el-select 
                v-model="formDeploy.artifact_url" 
                placeholder="请选择" 
                value-key="tag" 
                filterable
                clearable
                allow-create
                default-first-option
                reserve-keyword
                remote
                :remote-method="getArtifacts"
                :loading="loading.artifact"
                style="width:100%" 
                size="large">
                <el-option v-for="item in artifactList" :key="item.tag" :label="item.tag" :value="item">
                  <span style="float:left">{{item.tag}}</span>
                  <span style="float:right;color:var(--el-text-color-secondary);font-size:12px">{{item.last_modified}}</span>
                </el-option>
              </el-select>
            </el-form-item>
          </el-form>
        </el-card>
      </el-form-item>
    </el-form>
  </template>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show.deploy=false" size="large">取消</el-button>
      <el-button type="primary" @click="confirmDeploy()" size="large">提交</el-button>
    </div>
  </template>
</el-drawer>
</template>

<script setup>
import { ArrowRight,Search,Refresh,CopyDocument,Close } from '@element-plus/icons-vue'
import { onBeforeMount, ref, reactive, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { axios } from '/src/assets/util/axios'
import moment from 'moment'
import _ from 'lodash'
/* 变量定义 */
const router = useRouter()
const tenant = localStorage.tenant
const list = ref([])
const pageSize = ref(10)
const pageTotal = ref(0)
const current = ref(1)
const searchKey = ref("")
const show = ref({
  add: false,
  deploy: false
})
const edit = ref(false)
const form = ref({
  workload: [],
  traffic_policy: "weight",
  enable_traffic_control: false
})
const app = ref(null)
const tenants = ref([])
const repoList = ref([])
const k8scluster = ref({})
const rules = reactive({
  app_name: [{required: true, message: '请填写应用名称'}],
  repo_name: [{required: true, message: '请填写仓库/项目名称'}],
  image_name: [{required: true, message: '请填写镜像名称'}],
  repo: [{required: true, message: '请选择镜像仓库', trigger: 'change'}],
})
const loading = ref({
  searchapp: false,
  table: false,
  artifact: false,
})
const formDeploy = ref({})
const appList = ref([])
const artifactList = ref([])
// 添加标签的三个变量
const inputVisible = ref(false)
const inputLabel = ref('')
const inputRef = ref(null)
/* 生命周期函数 */
onBeforeMount(async () => {
  getRepoList()
  getClusterList()
  getList(1)
  getTenants()
});
/* methods */
const getList = async (page) => {
  let url = `page=${page}&size=${pageSize.value}&sort=update_at desc`
  if(searchKey.value != "") url += `&search=app_name==${searchKey.value}`
  loading.value.table = true
  let response = await axios.get(`/lizardcd/db/application?${url}`)
  loading.value.table = false
  list.value = response.results||[]
  pageTotal.value = response.total
}
const getClusterList = async () => {
  k8scluster.value = await axios.get(`/lizardcd/server/clusters`)
}
const getRepoList = async () => {
  let response = await axios.get(`/lizardcd/db/image_repository?size=100`)
  repoList.value = response.results
}
const getTenants = async () => {
  let response = await axios.get(`/lizardcd/db/tenant`)
  tenants.value = response.results
}
const addWorkload = () => {
  form.value.workload.push({
    cluster: '',
    namespace: '',
    workload_type: 'deployments',
    workload_name: '',
    container_name: '',
    weight: 50,
    headers: []
  })
}
const removeWorkload = (index) => {
  form.value.workload.splice(index, 1)
}
const copyWorkload = (index) => {
  let workload = Object.assign({}, form.value.workload[index])
  workload.headers = Object.assign([], workload.headers)
  form.value.workload.splice(index, 0, Object.assign({}, workload))
}
const confirmClick = async (f) => {
  if(!f) return
  await f.validate(async (valid) => {
    if(valid) {
      let params = Object.assign({}, form.value)
      params.update_at = moment()
      params.workload = JSON.stringify(params.workload)
      params.repo = JSON.stringify(params.repo)
      params.tags = JSON.stringify(params.tags)
      if(params.enable_traffic_control === true && params.traffic_policy === 'weight') {
        let weightTotal = 0
        for(let x of form.value.workload) {
          weightTotal += x.weight
        }
        if(weightTotal != 100) {
          ElMessage.warning('所有工作负载权重之和必须等于100')
          return
        }
      }
      if(edit.value === false) {
        await axios.post(`/lizardcd/db/application`, {body:params})
        getList(1)
        current.value = 1
        show.value.add = false
      }
      else {
        let id = params.id
        delete params.id
        await axios.put(`/lizardcd/db/application/${id}`, {body:params})
        getList(current.value)
        show.value.add = false
      }
    }
    else {
      ElMessage.warning('必填项未填完')
    }
  })
}
const searchApp = async (query) => {
  if(query) {
    loading.value.searchapp = true
    let response = await axios.get(`/lizardcd/db/application?search=app_name==${query}&size=20`)
    appList.value = response.results
    loading.value.searchapp = false
  }
  else {
    appList.value = []
  }
}
const getArtifacts = async (query) => {
  loading.value.artifact = true
  let url = `/lizardcd/server/repo/image/tags?app_name=${encodeURIComponent(formDeploy.value.app_name.app_name)}`
  if(query !== "") url += `&tag=${query}`
  let response = await axios.get(url)
  response ||= []
  artifactList.value = _.sortBy(response.map(x => {
    x.last_modified = moment(x.last_modified).format('YYYY-MM-DD HH:mm:ss')
    return x
  }), 'last_modified').reverse()
  loading.value.artifact = false
}
const confirmDeploy = async () => {
  ElMessageBox.confirm(
    '确认发布此版本？',
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    let params = {
      "app_name": formDeploy.value.app_name.app_name,
      "task_type": "deploy",
      "trigger_type": "手动触发",
      "workload": formDeploy.value.app_name.workload.map(x => {
        return {
          "cluster": x.cluster,
          "namespace": x.namespace,
          "workload_type": x.workload_type,
          "workload_name": x.workload_name,
          "container_name": x.container_name,
          "artifact_url": formDeploy.value.policy === 'same' ? (formDeploy.value.artifact_url.artifact_url || formDeploy.value.artifact_url) : (x.artifact_url.artifact_url || x.artifact_url)
        }
      })
    }
    let response = await axios.post(`/lizardcd/task/run`, params)
    router.push(`/task/history?id=${response.id}`)
  }).catch(() => {
    console.warn('cancel')
  })
}
const handleCommand = async (command) => {
  switch(command.action) {
    case "deploy": {
      let response = await axios.get(`/lizardcd/db/application?search=app_name==${command.row.app_name}&size=20`)
      appList.value = response.results
      formDeploy.value.app_name = command.row
      formDeploy.value.policy = 'same'
      show.value.deploy = true
      break
    }
    case "edit": {
      form.value = Object.assign({}, command.row)
      form.value.tenant ||= localStorage.tenant
      form.value.tags ||= []
      for(let w of form.value.workload) {
        w.headers = w.headers || []
      }
      edit.value = true
      show.value.add = true
      break
    }
    case "copy": {
      form.value = Object.assign({}, command.row)
      delete form.value.id
      edit.value = false
      show.value.add = true
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
          "workload": command.row.workload.map(x => {
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
        await axios.delete(`/lizardcd/db/application/${command.row.id}`)
        getList(current.value)
      }).catch(() =>{})
      break
    }
  }
}
const addMatchHeader = (headers) => {
  headers.push({
    key: "",
    match_type: "exact",
    value: ""
  })
}
const removeMatchHeader = (headers) => {
  headers.pop()
}
const handleClose = (tag) => {
  form.value.tags.splice(form.value.tags.indexOf(tag), 1)
}
const handleInputConfirm = () => {
  if(inputLabel.value) {
    form.value.tags.push(inputLabel.value)
  }
  inputVisible.value = false
  inputLabel.value = ''
}
const showInput = () => {
  inputVisible.value = true
  nextTick(() => {
    inputRef.value.input.focus()
  })
}
</script>
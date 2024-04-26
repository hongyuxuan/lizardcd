<template>
<el-breadcrumb :separator-icon="ArrowRight">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>应用管理</el-breadcrumb-item>
</el-breadcrumb>
<el-card>
  <template #header>
    <div class="card-header">
      <span class="card-header-text">应用管理</span>
      <span class="pull-right pointer" @click="getList(current)"><el-icon><Refresh /></el-icon></span>
    </div>
  </template>
  <el-row>
    <el-col :span="12">
      <el-button-group>
        <el-input v-model="searchKey" placeholder="输入名称进行搜索" size="large" :prefix-icon="Search" @change="getPage(1)" clearable style="width:300px;" />
      </el-button-group>
    </el-col>
    <el-col :span="12">
      <el-button-group class="pull-right">
        <el-button class="pull-right" size="large" type="primary" @click="show.add=true;edit=false;form={workload:[]}">新建应用</el-button>
        <el-button class="pull-right" size="large" type="primary" @click="show.deploy=true;formDeploy={}" style="margin-right:5px">发布应用</el-button>
      </el-button-group>
    </el-col>
  </el-row>
  <el-table :data="list" style="width:100%;margin-top:10px">
    <el-table-column type="selection" width="45" />
    <el-table-column prop="app_name" label="应用名称" min-width="200" />
    <el-table-column prop="repo" label="仓库地址" min-width="180">
      <template #default="scope">{{ scope.row.repo.repo_url }}</template>
    </el-table-column>
    <el-table-column prop="repo_name" label="仓库/项目" min-width="180" />
    <el-table-column prop="image_name" label="镜像名" min-width="180" />
    <el-table-column prop="update_at" label="更新时间" width="160">
      <template #default="scope">
        {{ moment(scope.row.update_at).format('YYYY-MM-DD HH:mm') }}
      </template>
    </el-table-column>
    <el-table-column prop="Option" label="操作" width="170">
      <template #default="scope">
        <el-tooltip effect="dark" content="发布" placement="top">
          <el-button circle @click="deployOne(scope.row)"><font-awesome-icon icon="rocket" /></el-button>
        </el-tooltip>
        <el-button :icon="EditPen" circle @click="editOne(scope.row)" />
        <el-tooltip effect="dark" content="复制" placement="top">
          <el-button :icon="CopyDocument" circle @click="copyOne(scope.row)" />
        </el-tooltip>
        <el-popconfirm title="确认删除？" @confirm="deleteOne(scope.row)">
          <template #reference>
            <el-button :icon="Delete" circle />
          </template>
        </el-popconfirm>
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
        <el-alert title="Artifactory填写仓库名，Harbor填写项目名" type="info" />
      </el-form-item>
      <el-form-item label="镜像名" prop="image_name">
        <el-input v-model="form.image_name" placeholder="请填写" size="large" />
      </el-form-item>
      <el-form-item label="工作负载">
        <table v-for="(m,index) in form.workload" :key="index" class="table table-bordered">
          <tbody>
          <tr>
            <td width=120>容器集群</td>
            <td>
              <el-select v-model="m.cluster" clearable placeholder="请选择" size="large" style="width:100%">
                <el-option v-for="(v,k,i) in k8scluster" :key="i" :label="k" :value="k" />
              </el-select>
            </td>
          </tr>
          <tr v-if="m.workload_type!='YAML'">
            <td width=120>命名空间</td>
            <td>
              <el-select v-model="m.namespace" clearable placeholder="请选择" size="large" style="width:100%">
                <el-option v-for="item in k8scluster[m.cluster]" :key="item" :label="item" :value="item" />
              </el-select>
            </td>
          </tr>
            <tr>
            <td width=120>工作负载类型</td>
            <td>
              <el-radio-group v-model="m.workload_type">
                <el-radio label="deployments" value="deployments" />
                <el-radio label="statefulsets" value="statefulsets" />
              </el-radio-group>
            </td>
          </tr>
          <tr v-if="m.workload_type!='YAML'">
            <td width=120>工作负载名称</td>
            <td>
              <el-input v-model="m.workload_name" size="large" />
            </td>
          </tr>
          <tr v-if="m.workload_type!='YAML'">
            <td width=120>容器名称</td>
            <td>
              <el-input v-model="m.container_name" size="large" />
            </td>
          </tr>
          <tr>
            <td width=120>操作</td>
            <td>
              <el-link type="primary" :underline="false" @click="removeWorkload(index)">删除本项</el-link>
            </td>
          </tr>
          </tbody>
        </table>
        <el-row>
          <el-button icon="Plus" circle @click="addWorkload"></el-button>
        </el-row>
      </el-form-item>
    </el-form>
  </template>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show.add=false">取消</el-button>
      <el-button type="primary" @click="confirmClick(app)">提交</el-button>
    </div>
  </template>
</el-drawer>
<el-drawer v-model="show.deploy" direction="rtl" size="600px">
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
          :loading="loading" 
          style="width:100%" 
          size="large">
          <el-option v-for="item in appList" :key="item.id" :label="item.app_name" :value="item" />
        </el-select>
      </el-form-item>
      <el-form-item label="选择制品">
        <el-select 
          v-model="formDeploy.artifact_url" 
          placeholder="请选择" 
          value-key="tag" 
          style="width:100%" 
          @focus="getArtifacts" 
          size="large">
          <el-option v-for="item in artifactList" :key="item.tag" :label="item.tag" :value="item">
            <span style="float:left">{{item.tag}}</span>
            <span style="float:right;color:var(--el-text-color-secondary);font-size:12px">{{item.last_modified}}</span>
          </el-option>
        </el-select>
      </el-form-item>
    </el-form>
  </template>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show.deploy=false">取消</el-button>
      <el-button type="primary" @click="confirmDeploy()">提交</el-button>
    </div>
  </template>
</el-drawer>
</template>

<script setup>
import { ArrowRight,Search,Refresh,EditPen,CopyDocument,Delete } from '@element-plus/icons-vue'
import { onBeforeMount, ref, reactive } from 'vue'
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
const searchKey = ref("")
const show = ref({
  add: false,
  deploy: false
})
const edit = ref(false)
const form = ref({
  workload: []
})
const app = ref(null)
const repoList = ref([])
const k8scluster = ref({})
const rules = reactive({
  app_name: [{required: true, message: '请填写应用名称'}],
  repo_name: [{required: true, message: '请填写仓库/项目名称'}],
  image_name: [{required: true, message: '请填写镜像名称'}],
  repo: [{required: true, message: '请选择镜像仓库', trigger: 'change'}],
})
const loading = ref(false)
const formDeploy = ref({})
const appList = ref([])
const artifactList = ref([])
/* 生命周期函数 */
onBeforeMount(async () => {
  getRepoList()
  getClusterList()
  getList(1)
});
/* methods */
const getList = async (page) => {
  let url = `page=${page}&size=${pageSize.value}&sort=update_at desc`
  if(searchKey.value != "") url += `&search=app_name==${searchKey.value}`
  let response = await axios.get(`/db/application?${url}`)
  list.value = response.results
  pageTotal.value = response.total
}
const getClusterList = async () => {
  k8scluster.value = await axios.get(`/lizardcd/clusters`)
}
const getRepoList = async () => {
  let response = await axios.get(`/db/image_repository?size=100`)
  repoList.value = response.results
}
const editOne = async (row) => {
  form.value = Object.assign({}, row)
  edit.value = true
  show.value.add = true
}
const copyOne = async (row) => {
  form.value = Object.assign({}, row)
  delete form.value.id
  edit.value = false
  show.value.add = true
}
const addWorkload = () => {
  form.value.workload.push({
    cluster: '',
    namespace: '',
    workload_type: 'deployments',
    workload_name: '',
    container_name: ''
  })
}
const removeWorkload = (index) => {
  form.value.workload.splice(index, 1)
}
const confirmClick = async (f) => {
  if(!f) return
  await f.validate(async (valid) => {
    if(valid) {
      let params = Object.assign({}, form.value)
      params.update_at = moment()
      params.workload = JSON.stringify(params.workload)
      params.repo = JSON.stringify(params.repo)
      if(edit.value === false) {
        await axios.post(`/db/application`, {body:params})
        getList(1)
        current.value = 1
        show.value.add = false
      }
      else {
        let id = params.id
        delete params.id
        await axios.put(`/db/application/${id}`, {body:params})
        getList(current.value)
        show.value.add = false
      }
    }
    else {
      ElMessage.warning('必填项未填完')
    }
  })
}
const deleteOne = async (row) => {
  await axios.delete(`/db/application/${row.id}`)
  getList(current.value)
}
const searchApp = async (query) => {
  if(query) {
    loading.value = true
    let response = await axios.get(`/db/application?search=app_name==${query}&size=20`)
    appList.value = response.results
    loading.value = false
  }
  else {
    appList.value = []
  }
}
const getArtifacts = async () => {
  let response = await axios.get(`/lizardcd/repo/image/tags?app_name=${formDeploy.value.app_name.app_name}`)
  response ||= []
  artifactList.value = _.sortBy(response.map(x => {
    x.last_modified = moment(x.last_modified).format('YYYY-MM-DD HH:mm:ss')
    return x
  }), 'last_modified').reverse()
}
const confirmDeploy = async () => {
  let params = Object.assign({}, formDeploy.value)
  ElMessageBox.confirm(
    '确认发布此版本？',
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    for(let x of formDeploy.value.app_name.workload) {
      await axios.patch(`/kubernetes/cluster/${x.cluster}/namespace/${x.namespace}/${x.workload_type}/${x.workload_name}?container=${x.container_name}&image=${formDeploy.value.artifact_url.artifact_url}`)
    }
    ElMessage.success('提交成功')
    getList(1)
    current.value = 1
    show.value.deploy = false
  }).catch(() => {
    console.warn('cancel')
  })
}
const deployOne = async (row) => {
  let response = await axios.get(`/db/application?search=app_name==${row.app_name}&size=20`)
  appList.value = response.results
  formDeploy.value.app_name = row
  show.value.deploy = true
} 
</script>
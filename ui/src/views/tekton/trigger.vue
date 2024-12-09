<template>
<el-breadcrumb :separator-icon="ArrowRight">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>持续集成</el-breadcrumb-item>
  <el-breadcrumb-item>精准触发</el-breadcrumb-item>
</el-breadcrumb>
<el-card>
  <template #header>
    <div class="card-header">
      <span class="card-header-text">精准触发</span>
    </div>
  </template>
  <el-row>
    <el-col :span="12">
      <el-button-group>
        <el-button :icon="Refresh" size="large" style="margin-right:5px" @click="getList(current)" />
        <el-input v-model="searchKey" placeholder="输入名称进行搜索" size="large" :prefix-icon="Search" @change="current=1;getList(1)" clearable style="width:300px;" />
      </el-button-group>
    </el-col>
    <el-col :span="12">
      <el-button class="pull-right" size="large" type="primary" @click="show=true;edit=false;form={trigger_path:[''],trigger_type:'Tekton',tenant}">+ 新建规则</el-button>
    </el-col>
  </el-row>
  <el-table :data="list" class="line-height40" style="width:100%;margin-top:10px" @selection-change="select">
    <el-table-column type="selection" width="45" />
    <el-table-column prop="trigger_name" label="应用名称" min-width="150" />
    <el-table-column prop="git_http_url" label="Git地址" min-width="250" />
    <el-table-column prop="trigger_path" label="触发路径" min-width="200">
      <template #default="scope">
        {{ scope.row.trigger_path.join(', ') }}
      </template>
    </el-table-column>
    <el-table-column prop="trigger_endpoint" label="触发地址" min-width="180" />
    <el-table-column prop="update_at" label="更新时间" width="150">
      <template #default="scope">
        {{ moment(scope.row.update_at).format('YYYY-MM-DD HH:mm') }}
      </template>
    </el-table-column>
    <el-table-column prop="Option" label="操作" width="140">
      <template #default="scope">
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
<el-drawer v-model="show" direction="rtl" size="700px">
  <template #header>
    <h4 v-if="edit===false">新建规则</h4>
    <h4 v-if="edit===true">编辑规则</h4>
  </template>
  <template #default>
    <el-form ref="templates" :model="form" :rules="rules" label-width="130px">
      <el-form-item label="选择应用" prop="app_name" >
        <el-select 
          style="width:100%"
          v-model="form.app_name" 
          filterable
          remote
          clearable
          reserve-keyword
          placeholder="输入名称进行查询"
          value-key="id"
          :remote-method="searchApp"
          @change="selectApp"
          size="large">
          <el-option v-for="item in appList" :key="item.id" :label="item.app_name" :value="item" />
        </el-select>
      </el-form-item>
      <el-form-item label="代码仓库地址" prop="git_http_url" >
        <el-input v-model="form.git_http_url" size="large" />
      </el-form-item>
      <el-form-item label="文件路径">
        <table v-for="(m,index) in form.trigger_path" :key="index" class="table table-bordered" style="margin-bottom:0">
          <tbody>
          <tr>
            <td><el-input v-model="form.trigger_path[index]" /></td>
            <td width="80">
              <el-button-group>
                <el-button icon="Plus" circle @click="addPath(index)"></el-button>
                <el-button icon="Close" circle @click="removePath(index)"></el-button>
              </el-button-group>
            </td>
          </tr>
          </tbody>
        </table>
      </el-form-item>
      <el-form-item label="触发类型">
        <el-radio-group v-model="form.trigger_type">
          <el-radio value="Tekton">Tekton</el-radio>
          <el-radio value="自定义">自定义</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="选择EventListener" v-if="form.trigger_type==='Tekton'">
        <el-select v-model="cluster" placeholder="请选择集群" clearable filterable style="width:200px;margin:0 8px 8px 0" size="large">
          <el-option v-for="(v,k) in clusterList" :key="k" :label="k" :value="k" />
        </el-select>
        <el-select v-model="namespace" placeholder="请选择命名空间" clearable filterable @change="getEventListener" style="width:200px;margin:0 8px 8px 0" size="large">
          <el-option v-for="(item) in clusterList[cluster]" :key="item" :label="item" :value="item" />
        </el-select>
        <el-select 
          v-model="eventListener" 
          placeholder="请选择EventListener" 
          clearable 
          filterable 
          value-key="metadata.name" 
          @change="selectEventListener"
          style="width:100%;" 
          size="large">
          <el-option v-for="(item,index) in eventListenerList" :key="index" :label="item.metadata.name" :value="item" />
        </el-select>
      </el-form-item>
      <el-form-item label="触发地址">
        <el-input v-model="form.trigger_endpoint" size="large" />
      </el-form-item>
      <el-form-item label="提交body">
        <v-ace-editor
          v-model:value="form.trigger_body"
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
            maxLines: 100,
          }" />
        <myTips type="info">go-template格式<br>支持以下变量：.Appname, .Ref, .GitSSHUrl, .GitHttpUrl</myTips>
      </el-form-item>
      <el-form-item label="Secret" prop="secret">
        <el-input v-model="form.secret" size="large" type="password" show-password />
      </el-form-item>
      <el-form-item label="所属租户" prop="tenant">
        <el-input v-model="form.tenant" disabled size="large" />
      </el-form-item>
    </el-form>
  </template>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show=false">取消</el-button>
      <el-button type="primary" @click="confirmClick(templates)">提交</el-button>
    </div>
  </template>
</el-drawer>
</template>
<script setup>
import { axios } from '/src/assets/util/axios.js'
import { ArrowRight,Search,Refresh,EditPen,Delete,Plus,CopyDocument } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { onBeforeMount, ref, reactive } from 'vue'
import MyTips from '/src/components/myTips/myTips.vue'
import moment from "moment"
import _ from 'lodash'
/* 引入v-ace-editor */
import { VAceEditor } from 'vue3-ace-editor'
import 'ace-builds/src-noconflict/mode-yaml'
import 'ace-builds/src-noconflict/theme-chrome'
import 'ace-builds/src-noconflict/ext-language_tools'
import myTips from '@/components/myTips'
/* 变量定义 */
const tenant = localStorage.tenant.split(',')[0]
const list = ref([])
const searchKey = ref("")
const pageSize = ref(20)
const pageTotal = ref(0)
const current = ref(1)
const selected = ref([])
const show = ref(false)
const edit = ref(false)
const form = ref({app_name:"",git_http_url:'',trigger_path:[""]})
const rules = reactive({
  app_name: [{required: true, message: '请填写应用名称', trigger: "blur"}],
  git_http_url: [{required: true, message: '请填写Git地址', trigger: "blur"}],
})
const templates = ref(null)
const appList = ref()
const cluster = ref("")
const clusterList = ref({})
const namespace = ref("")
const eventListener = ref()
const eventListenerList = ref([])
/* 生命周期函数 */
onBeforeMount(async () => {
  getClusterList()
  getList(1)
});
/* methods */
const getList = async (page) => {
  let url = `/lizardcd/db/ci_trigger?page=${page}&size=${pageSize.value}&sort=ci_trigger.update_at desc`
  if(searchKey.value != "")
    url += `&search=trigger_name==${searchKey.value}`
  let response = await axios.get(url)
  list.value = response.results
  pageTotal.value = response.total
}
const getClusterList = async () => {
  clusterList.value = await axios.get(`/lizardcd/server/clusters`)
}
const getEventListener = async () => {
  eventListenerList.value = await axios.get(`/lizardcd/tekton/cluster/${cluster.value}/namespace/${namespace.value}/eventlisteners`)
}
const searchApp = async (query) => {
  let response = await axios.get(`/lizardcd/db/application?page=1&size=20&sort=app_name&search=app_name==${encodeURIComponent(query)}`)
  appList.value = _.uniqBy(response.results, 'app_name')
}
const selectApp = async (val) => {
  if(!val) {
    form.value.git_http_url = ""
    form.value.trigger_name = ""
  }
  else {
    form.value.git_http_url = val.git_http_url
    form.value.trigger_name = val.app_name
  }
}
const editOne = async (row) => {
  form.value = Object.assign({}, row)
  form.value.tenant = tenant
  form.value.app_name = row.application
  form.value.trigger_name ||= row.application.app_name
  searchApp(row.application.app_name)
  delete form.value.application
  edit.value = true
  show.value = true
}
const select = (val) => {
  selected.value = val
}
const copyOne = async (row) => {
  let params = Object.assign({}, row)
  params.app_name += ` @${moment().format('YYYYMMDDHHmmss')}`
  delete params.id
  await editOne(params)
  edit.value = false
}
const deleteOne = async (row) => {
  await axios.delete(`/lizardcd/db/ci_trigger/${row.id}`)
  getList(current.value)
}
const addPath = (index) => {
  form.value.trigger_path.splice(index+1, 0, "")
}
const removePath = (index) => {
  form.value.trigger_path.splice(index, 1)
}
const confirmClick = async (f) => {
  if(!f) return
  await f.validate(async (valid) => {
    if(valid) {
      let params = Object.assign({}, form.value)
      params.app_id = params.app_name.id
      params.update_at = moment()
      delete params.app_name
      if(edit.value === false) {
        await axios.post(`/lizardcd/db/ci_trigger`, {body:params})
        getList(1)
        current.value = 1
        show.value = false
      }
      else {
        await axios.put(`/lizardcd/db/ci_trigger/${params.id}`, {body:params})
        getList(current.value)
        show.value = false
      }
    }
    else {
      ElMessage.warning('必填项未填完')
    }
  })
}
const selectEventListener = (val) => {
  if(!val.metadata.annotations?.endpoint)
    ElMessage.warning('EventListener未找到endpoint注解')
  else
    form.value.trigger_endpoint = val.metadata.annotations.endpoint
}
</script>
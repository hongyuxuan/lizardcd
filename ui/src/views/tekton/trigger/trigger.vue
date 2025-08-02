<template>
<el-breadcrumb :separator-icon="ArrowRight">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>持续集成</el-breadcrumb-item>
  <el-breadcrumb-item>触发配置</el-breadcrumb-item>
</el-breadcrumb>
<el-card>
  <template #header>
    <div class="card-header">
      <span class="card-header-text">触发配置</span>
    </div>
  </template>
  <el-row>
    <el-col :span="12">
      <el-button-group>
        <el-button :icon="Refresh" size="large" @click="getList(current)" />
        <el-input v-model="searchKey" placeholder="输入名称进行搜索" size="large" :prefix-icon="Search" @change="current=1;getList(1)" clearable style="width:300px;" />
      </el-button-group>
    </el-col>
    <el-col :span="12">
      <el-dropdown @command="handleMore" class="pull-right">
        <el-button size="large">更多操作<el-icon class="el-icon--right"><arrow-down /></el-icon></el-button>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item :command="{action:'modifyBatch'}">批量更新</el-dropdown-item>
            <el-dropdown-item :command="{action:'deleteBatch'}">批量删除</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
      <el-button class="pull-right" size="large" type="primary" @click="show.add=true;edit=false;form={trigger_path:[''],trigger_type:'pipelinerun',trigger_event:['push','merge_request'],tenant,pre_check:{}}">+ 新建规则</el-button>
    </el-col>
  </el-row>
  <el-table :data="list" class="line-height40" style="width:100%;margin-top:10px" @selection-change="select">
    <el-table-column type="selection" width="45" />
    <el-table-column prop="trigger_name" label="应用名称" min-width="150" />
    <el-table-column prop="git_http_url" label="Git地址" min-width="250" />
    <el-table-column prop="trigger_endpoint" label="标签" min-width="180">
      <template #default="scope">
        <el-tag v-for="(item,i) in scope.row.match_labels||[]" :key="i" size="large" type="primary">{{ item }}</el-tag>
      </template>
    </el-table-column>
    <el-table-column prop="ref_pattern" label="分支白名单" min-width="100" />
    <el-table-column prop="tenant" label="所属租户" min-width="80" />
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
<el-drawer v-model="show.add" direction="rtl" size="700px">
  <template #header>
    <h4 v-if="edit===false">新建规则</h4>
    <h4 v-if="edit===true">编辑规则</h4>
  </template>
  <template #default>
    <el-form ref="templates" :model="form" :rules="rules" label-width="130px">
      <el-form-item label="选择应用" prop="app_name">
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
          <el-radio value="pipelinerun">Tekton PipelineRun</el-radio>
          <el-radio value="endpoint">Tekton Endpoint</el-radio>
          <el-radio value="自定义">自定义</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="选择EventListener" v-if="form.trigger_type==='endpoint'">
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
      <el-form-item label="触发地址" v-if="form.trigger_type!=='pipelinerun'">
        <el-input v-model="form.trigger_endpoint" size="large" clearable />
      </el-form-item>
      <el-form-item label="提交body" v-if="form.trigger_type!=='pipelinerun'">
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
      <el-form-item v-if="form.trigger_type==='pipelinerun'">
        <template #label><el-text>标签匹配 
          <el-tooltip placement="top">
            <template #content>
              触发webhook后，将会从设置的默认tekton集群中查询含有选中labels的Pipeline，并执行
            </template>
            <el-icon class="text-yellow"><Warning /></el-icon>
          </el-tooltip></el-text>
        </template>
        <el-tag v-for="(item, i) in form.match_labels" :key="i" closable type="primary" size="large" @close="closeLabel(i)">{{ item }}</el-tag>
      </el-form-item>
      <el-form-item v-if="form.trigger_type==='pipelinerun'">
        <template #label><el-text>YAML扫描路径 
          <el-tooltip placement="top">
            <template #content>
              触发webhook后，将会从gitlab中的扫描路径搜索后缀为.tekton.(yaml|yml)的YAML文件，<br>并根据此文件执行PipelineRun
            </template>
            <el-icon class="text-yellow"><Warning /></el-icon>
          </el-tooltip></el-text>
        </template>
        <el-input v-model="form.scan_path" size="large" clearable />
        <myTips type="info">如配置了YAML扫描路径，则标签匹配失效。如不需要扫描，则此处留空</myTips>
      </el-form-item>
      <el-form-item label="分支白名单" prop="ref_pattern">
        <el-input v-model="form.ref_pattern" size="large" clearable />
        <myTips type="info">需要触发 webhook 的分支名称正则表达式</myTips>
      </el-form-item>
      <el-form-item label="触发事件" prop="trigger_event">
        <el-checkbox-group v-model="form.trigger_event">
          <el-checkbox label="push" value="push" />
          <el-checkbox label="merge_request" value="merge_request" />
        </el-checkbox-group>
      </el-form-item>
      <el-form-item label="开启预检查">
        <el-checkbox v-model="form.pre_check.enable" style="margin-top:-5px" @click="enablePreCheck()" />
      </el-form-item>
      <el-form-item label="预检查分支白名单" v-if="form.pre_check.enable===true">
        <el-input v-model="form.pre_check.ref_pattern" size="large" clearable />
      </el-form-item>
      <el-form-item label="预检查触发" v-if="form.pre_check.enable===true">
        <el-checkbox-group v-model="form.pre_check.trigger_event">
          <el-checkbox label="push" value="push" />
          <el-checkbox label="merge_request" value="merge_request" />
        </el-checkbox-group>
      </el-form-item>
      <el-form-item label="唯一实例" prop="uniq_instance">
        <el-checkbox label="勾选后，后运行的流水线将终止前面的运行实例" v-model="form.uniq_instance" />
      </el-form-item>
      <el-form-item label="所属租户" prop="tenant">
        <el-select 
          v-model="form.tenant" 
          placeholder="请选择租户" 
          clearable 
          filterable 
          :disabled="userInfo.role!=='admin'"
          style="width:100%;" 
          size="large">
          <el-option v-for="(item,index) in tenantList" :key="index" :label="item" :value="item" />
        </el-select>
      </el-form-item>
    </el-form>
  </template>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show.add=false">取消</el-button>
      <el-button type="primary" @click="confirmClick(templates)">提交</el-button>
    </div>
  </template>
</el-drawer>
<el-drawer v-model="show.batchedit" direction="rtl" size="700px">
  <template #header>
    <h4>批量更新</h4>
  </template>
  <template #default>
    <triggerbatchedit ref="batcheditRef" :clusterList="clusterList" :tenantList="tenantList" />
  </template>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show.batchedit=false">取消</el-button>
      <el-button type="primary" @click="submitBatchedit()">提交</el-button>
    </div>
  </template>
</el-drawer>
</template>
<script setup>
import { axios } from '/src/assets/util/axios.js'
import { ArrowRight,Search,Refresh,EditPen,Delete,Plus,CopyDocument } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { onBeforeMount, ref, reactive, computed } from 'vue'
import { useStore } from 'vuex'
import triggerbatchedit from './batchedit.vue'
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
const store = useStore()
const userInfo = computed(() => {
  return store.state.userInfo
})
const tenant = localStorage.tenant.split(',')[0]
const tenantList = ref([])
const list = ref([])
const searchKey = ref("")
const pageSize = ref(20)
const pageTotal = ref(0)
const current = ref(1)
const selected = ref([])
const show = ref({
  add: false,
  batchedit: false,
})
const edit = ref(false)
const form = ref({app_name:"",git_http_url:'',trigger_path:[""]})
const rules = reactive({
  app_name: [{required: true, message: '请填写应用名称', trigger: "blur"}],
  git_http_url: [{required: true, message: '请填写Git地址', trigger: "blur"}],
  ref_pattern: [{required: true, message: '请填写分支白名单', trigger: "blur"}],
})
const templates = ref(null)
const appList = ref()
const cluster = ref("")
const clusterList = ref({})
const namespace = ref("")
const eventListener = ref()
const eventListenerList = ref([])
const batcheditRef = ref(null)
/* 生命周期函数 */
onBeforeMount(async () => {
  getClusterList()
  getTenantList()
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
const getTenantList = async () => {
  let response = await axios.get(`/lizardcd/db/tenant`)
  tenantList.value = response.results.map(x => x.tenant_name)
}
const getClusterList = async () => {
  clusterList.value = await axios.get(`/lizardcd/server/clusters`)
}
const getEventListener = async () => {
  let response = await axios.get(`/lizardcd/tekton/cluster/${cluster.value}/namespace/${namespace.value}/eventlisteners`)
  eventListenerList.value = response.results
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
    form.value.match_labels = val.tags.map(x => x.replace(":", "="))
  }
}
const editOne = async (row) => {
  form.value = _.cloneDeep(row)
  form.value.tenant ||= tenant
  form.value.app_name = row.application
  form.value.trigger_name ||= row.application.app_name
  form.value.trigger_event = row.trigger_event?.split(',')||[]
  form.value.pre_check ||= {}
  if(form.value.pre_check.enable) {
    form.value.pre_check.trigger_event = form.value.pre_check.trigger_event.split(',')||[]
  }
  searchApp(row.application.app_name)
  delete form.value.application
  edit.value = true
  show.value.add = true
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
      let params = _.cloneDeep(form.value)
      params.app_id = params.app_name.id
      params.trigger_event = params.trigger_event.join(',')
      if(params.pre_check?.enable) {
        params.pre_check.trigger_event = params.pre_check.trigger_event.join(",")
      }
      params.update_at = moment()
      delete params.app_name
      if(edit.value === false) {
        await axios.post(`/lizardcd/db/ci_trigger`, {body:params})
        getList(1)
        current.value = 1
      }
      else {
        await axios.put(`/lizardcd/db/ci_trigger/${params.id}`, {body:params})
        getList(current.value)
      }
      show.value.add = false
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
const closeLabel = (index) => {
  form.value.match_labels.splice(index, 1)
}
const handleMore = async (command) => {
  if(selected.value.length === 0) {
    ElMessage.warning({message: '请勾选记录'})
    return
  }
  switch(command.action) {
    case "deleteBatch": {
      await ElMessageBox.confirm('确定删除？','警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(async () => {
        await Promise.all(selected.value.map(x => {
          return axios.delete(`/lizardcd/db/ci_trigger/${x.id}`)
        }))
        getList(current.value)
      }).catch(() =>{})
      break
    }
    case "modifyBatch": {
      show.value.batchedit = true
      break
    }
  }
}
const submitBatchedit = async () => {
  await batcheditRef.value.submit(selected.value)
  getList(current.value)
  show.value.batchedit = false
}
const enablePreCheck = () => {
  if(form.value.pre_check.enable === true) {
    form.value.pre_check.ref_pattern = "feature.*"
    form.value.pre_check.trigger_event = ["push", "merge_request"]
  } else {
    form.value.pre_check = {}
  }
}
</script>
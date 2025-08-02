<template>
<el-breadcrumb :separator-icon="ArrowRight">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>连接管理</el-breadcrumb-item>
</el-breadcrumb>
<el-card>
  <template #header>
    <div class="card-header">
      <span class="card-header-text">在线服务</span>
    </div>
  </template>
  <el-row>
    <el-col :span="12">
      <el-button-group style="width:100%">
        <el-button :icon="Refresh" size="large" @click="getList()" />
        <el-input v-model="searchKey" clearable placeholder="输入关键词查询……" :prefix-icon="Search" @change="current=1;filterList=Object.assign({},all);getPage(1)" style="width:50%" size="large" />
      </el-button-group>
    </el-col>
    <el-col :span="12">
        <el-button-group class="pull-right">
          <el-button class="pull-right" size="large" @click="show.list=true">连接配置</el-button>
          <el-button class="pull-right" size="large" type="primary" @click="show.add=true;form={kubeconfig:'',labels:[]}">+ 手工注册</el-button>
        </el-button-group>
      </el-col>
  </el-row>
  <el-table :data="list" class="line-height40" :show-header="true" @expand-change="getServiceMeta" style="width:100%;">
    <el-table-column type="expand">
      <template #default="scope">
        <el-table :data="serviceMeta[scope.row.service_name]||[]" style="width:100%;margin-left:50px;">
          <el-table-column prop="ServiceID" label="ServiceID" />
          <el-table-column prop="ServiceMeta" label="ServiceMeta">
            <template #default="props">
              <el-tag v-for="(v,k) of props.row.ServiceMeta" size="large">{{ k }}={{ v }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="Labels" label="Labels">
            <template #default="props">
              <el-tag v-for="(item,i) of props.row.Labels" :key="i" size="large">{{ item }}</el-tag>
            </template>
          </el-table-column>
        </el-table>
      </template>  
    </el-table-column>
    <el-table-column prop="service_name" label="Service Key" min-width="160" />
    <el-table-column prop="service_type" label="Service Type" min-width="160" />
    <el-table-column prop="service_source" label="Service Source" min-width="160" />
  </el-table>
  <el-pagination 
    class="pull-right"
    background 
    v-model:page-size="pageSize"
    :page-sizes="[20, 30, 50, 100]"
    layout="total, sizes, prev, pager, next, jumper" 
    :total="pageTotal"
    @current-change="getPage"
    @size-change="handleSizeChange"
    v-model:current-page="current" />
</el-card>
<el-drawer v-model="show.add" direction="rtl" size="700px">
  <template #header>
    <h4>手工注册</h4>
  </template>
  <template #default>
    <el-form ref="refAdd" :model="form" :rules="rules" label-width="100px">
      <el-form-item label="ServiceKey" prop="service_key">
        <el-input v-model="form.service_key" size="large"clearable  />
      </el-form-item>
      <el-form-item label="Endpoint" prop="endpoint">
        <el-input v-model="form.endpoint" size="large" clearable />
      </el-form-item>
      <el-form-item label="Proxy">
        <el-input v-model="form.proxy" size="large" clearable />
        <myTips type="info">如果 lizardcd-server 无法直接访问 lizardcd-agent，可通过设置代理访问</myTips>
      </el-form-item>
      <el-form-item label="Kubeconfig">
          <el-switch v-model="wrapLine" active-text="自动换行" inactive-text="不自动换行" inline-prompt class="el-switch-hover-right" size="large" />
          <v-ace-editor
            v-model:value="form.kubeconfig"
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
            }" />
        <myTips type="info">如果 lizardcd-server 可直接访问 Kubernetes 集群<br>且有 Kubeconfig，可在此配置</myTips>
      </el-form-item>
      <el-form-item label="设置标签" prop="labels">
        <el-row v-for="(item,i) in form.labels" :key="i" style="margin-bottom:5px;width:100%">
          <el-button-group>
            <el-input v-model="item.key" size="large" clearable style="width:200px;margin-right:5px" />
            <el-input v-model="item.value" size="large" clearable style="width:200px;" />
            <el-button circle :icon="Delete" size="large" style="float:right" @click="removeTag(i)" />
          </el-button-group>
        </el-row>
        <el-button circle icon="Plus" @click="addTag" />
      </el-form-item>
    </el-form>
  </template>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show.add=false">取消</el-button>
      <el-button type="primary" @click="confirmClick(refAdd)" :loading="loading">提交</el-button>
    </div>
  </template>
</el-drawer>
<el-drawer v-model="show.list" direction="rtl" size="700px">
  <template #header>
    <h4>连接配置</h4>
  </template>
  <template #default>
    <el-row>
      <el-col :span="24">
        <el-button-group style="width:100%">
          <el-button :icon="Refresh" size="large" style="margin-right:5px" @click="getAgents(current2)" />
          <el-input v-model="searchKey2" clearable placeholder="输入关键词查询……" :prefix-icon="Search" @change="current2=1;getAgents(1)" style="width:300px;margin-right:5px" size="large" />
        </el-button-group>
      </el-col>
    </el-row>
    <el-table :data="agents" class="line-height25" style="width:100%;">
      <el-table-column prop="service_key" label="Service Key" min-width="180" />
      <el-table-column label="Service Type" min-width="120">
        <template #default="scope">{{ scope.row.endpoint ? 'agent' : 'kubeconfig' }}</template>
      </el-table-column>
      <el-table-column prop="endpoint" label="Endpoint" min-width="120" />
      <el-table-column label="操作" width="60">
        <template #default="scope">
          <el-popconfirm title="确认删除？" @confirm="deleteOne(scope.row)">
            <template #reference>
              <el-button :icon="Delete" circle :disabled="userInfo.role!=='admin'" />
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination 
      class="pull-right"
      background 
      v-model:page-size="pageSize"
      :page-sizes="[10, 30, 50, 100]"
      layout="total, prev, pager, next, jumper" 
      :total="pageTotal2"
      @current-change="getAgents"
      v-model:current-page="current2" />
  </template>
</el-drawer>
</template>

<script setup>
import { ArrowRight,Refresh,Search,Delete } from '@element-plus/icons-vue'
import { onBeforeMount, ref, reactive, computed } from 'vue'
import { useStore } from 'vuex'
import { axios } from '/src/assets/util/axios.js'
import MyTips from '/src/components/myTips/myTips.vue'
/* 引入v-ace-editor */
import { VAceEditor } from 'vue3-ace-editor'
import 'ace-builds/src-noconflict/mode-yaml'
import 'ace-builds/src-noconflict/theme-chrome'
import 'ace-builds/src-noconflict/ext-language_tools'
import _ from 'lodash'
/* 变量定义 */
const store = useStore()
const userInfo = computed(() => {
  return store.state.userInfo
})
const list = ref([])
const filterList = ref([])
const all = ref([])
const pageSize = ref(20)
const pageTotal = ref(0)
const pageTotal2 = ref(0)
const current = ref(1)
const current2 = ref(1)
const serviceMeta = ref({})
const searchKey = ref("")
const searchKey2 = ref("")
const show = ref({
  add: false,
  list: false
})
const form = ref({kubeconfig:'',labels:[]})
const rules = reactive({
  service_key: [{required: true, message: '请填写ServiceKey'}],
})
const refAdd = ref(null)
const agents = ref([])
const loading = ref(false)
const wrapLine = ref(true)
/* 生命周期函数 */
onBeforeMount(async () => {
  getList()
  getAgents(1)
})
/* methods */
const getList = async () => {
  all.value = await axios.get(`/lizardcd/server/services`)
  getPage(current.value)
}
const getAgents = async (page) => {
  let url = `page=${page}&size=${pageSize.value}`
  if(searchKey2.value != "") url += `&search=service_key==${searchKey2.value}`
  let response = await axios.get(`/lizardcd/db/agent?${url}`)
  agents.value = response.results
  pageTotal2.value = response.total
}
const getPage = async (page) => {
  filterList.value = all.value.filter(n => n.service_name.includes(searchKey.value))
  pageTotal.value = filterList.value.length
  list.value = filterList.value.slice((page-1)*pageSize.value, page*pageSize.value)
}
const getServiceMeta = async (row) => {
  let serviceIds = await axios.get(`/lizardcd/server/services/${row.service_name}`)
  serviceMeta.value[row.service_name] = serviceIds.map(y => {
    return {
      ServiceName: y.ServiceName,
      ServiceID: y.ServiceID,
      ServiceMeta: y.ServiceMeta,
      Labels: y.Labels,
    }
  })
}
const handleSizeChange = async (size) => {
  pageSize.value = size
  await getList()
}
const confirmClick = async (f) => {
  if(!f) return
  await f.validate(async (valid) => {
    if(valid) {
      loading.value = true
      let params = _.cloneDeep(form.value)
      params.labels = form.value.labels.map(x => `${x.key}=${x.value}`)
      try {
        await axios.post(`/lizardcd/server/service`, params)
      } finally {
        loading.value = false
      }
      show.value.add = false
      getList()
      getAgents(current2.value)
    }
  })
}
const deleteOne = async (row) => {
  await axios.delete(`/lizardcd/db/agent/${row.id}`)
  getAgents(current2.value)
}
const addTag = () => {
  form.value.labels.push({key: '', value: ''})
}
const removeTag = (index) => {
  form.value.labels.splice(index, 1)
}
</script>
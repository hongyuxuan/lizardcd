<template>
<div class="box box-item">
  <div class="box-body" style="padding-top:20px;padding-bottom:0">
    <el-alert title="关于DestinationRule配置参考：https://istio.io/latest/zh/docs/reference/config/networking/destination-rule/" type="warning" style="margin-bottom:15px" />
    <el-row>
      <el-col :span="12">
        <el-button-group>
          <el-button :icon="Refresh" size="large" style="margin-right:5px" @click="getList(true)" />
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
        <el-button-group class="pull-right">
          <el-button class="pull-right" size="large" type="primary" @click="show.yaml=true;edit=false;yamlContent=''">创建目标规则</el-button>
        </el-button-group>
      </el-col>
    </el-row>
    <el-table 
      :data="list" 
      v-loading="loading"
      element-loading-text="奋力加载中..."
      class="line-height40" 
      style="width:100%;margin-top:10px;min-height:150px">
      <el-table-column label="" width="45">
        <font-awesome-icon icon="circle-nodes" style="font-size:25px;vertical-align:middle;" />
      </el-table-column>
      <el-table-column prop="name" label="名称" min-width="200">
        <template #default="scope">
          <el-link :underline="false">{{ scope.row.metadata.name }}</el-link>
        </template>
      </el-table-column>
      <el-table-column label="创建时间">
        <template #default="scope">
          {{ moment(scope.row.lastUpdateTime).format('YYYY-MM-DD HH:mm:ss') }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="100">
        <template #default="scope">
          <el-button icon="EditPen" circle @click="editOne(scope.row)"></el-button>
          <el-popconfirm title="确定删除?" confirm-button-text="确认" cancel-button-text="取消" @confirm="deleteOne(scope.row)">
            <template #reference>
              <el-button icon="Close" circle />
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination 
      class="pull-right"
      background 
      v-model:page-size="pageSize"
      :page-sizes="[10, 20, 50, 100]"
      layout="total, sizes, prev, pager, next, jumper" 
      :total="pageTotal"
      @size-change="handleSizeChange"
      @current-change="getPage"
      v-model:current-page="current" />
  </div>
</div>
<el-drawer v-model="show.yaml" direction="rtl" size="800px">
  <template #header>
    <h4>{{ edit === true ? '编辑YAML' : '创建目标规则' }}</h4>
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
        maxLines: 5000,
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
import { Search,Refresh } from '@element-plus/icons-vue'
import { onBeforeMount, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { axios } from '/src/assets/util/axios'
import moment from 'moment'
import _ from 'lodash'
/* 引入v-ace-editor */
import { VAceEditor } from 'vue3-ace-editor'
import 'ace-builds/src-noconflict/mode-yaml'
import 'ace-builds/src-noconflict/theme-chrome'
import 'ace-builds/src-noconflict/ext-language_tools'
/* 变量定义 */
const all = ref([])
const searchKey = ref("")
const list = ref([])
const pageSize = ref(10)
const pageTotal = ref(0)
const current = ref(1)
const cluster = ref("")
const clusterList = ref({})
const namespace = ref("")
const loading = ref(false)
const show = ref({
  yaml: false
})
const yamlContent = ref("")
const edit = ref(false)
/* 生命周期函数 */
onBeforeMount(async () => {
  getClusterList()
})
/* methods */
const getClusterList = async () => {
  clusterList.value = await axios.get(`/lizardcd/server/clusters`)
}
const getList = async (ifLoading) => {
  if(ifLoading) loading.value = true
  if(cluster.value !== "" && namespace.value !== "") {
    let response = await axios.get(`/lizardcd/istio/cluster/${cluster.value}/namespace/${namespace.value}/destinationrules`)
    all.value = _.sortBy(response, 'creationTimestamp').reverse()
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
const editOne = async (row) => {
  show.value.yaml = true
  edit.value = true
  yamlContent.value = await axios.get(`/lizardcd/istio/cluster/${cluster.value}/namespace/${namespace.value}/destinationrules/${row.metadata.name}/yaml`)
}
const deleteOne = async (row) => {
  await axios.delete(`/lizardcd/kubernetes/cluster/${cluster.value}/namespace/${namespace.value}/destinationrules/${row.metadata.name}`)
  ElMessage.success({message: '删除成功'})
  setTimeout(async () => {
    await getList(true)
  }, 2000)
}
const submitYaml = async () => {
  if(cluster.value === "" || namespace.value == "") {
    ElMessage.warning({message: '请指定集群和命名空间'})
    return
  } 
  await axios.patch(`/lizardcd/kubernetes/cluster/${cluster.value}/namespace/${namespace.value}/apply/yaml?kind=DestinationRule`, yamlContent.value, {
    headers: {
      'Content-Type': 'text/plain'
    }
  })
  show.value.yaml = false
  setTimeout(async () => {
    await getList(true)
  }, 2000)
}
const handleSizeChange = async (size) => {
  pageSize.value = size
  await getList(true)
}
</script>
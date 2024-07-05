<template>
<div class="box box-item">
  <div class="box-body" style="padding-top:20px;padding-bottom:0">
    <el-row>
      <el-col :span="12">
        <el-button-group>
          <el-button :icon="Refresh" size="large" style="margin-right:5px" @click="getList" />
          <el-select v-model="cluster" placeholder="请选择集群" clearable filterable style="width:200px;margin-right:8px" size="large">
            <el-option v-for="(v,k) in clusterList" :key="k" :label="k" :value="k" />
          </el-select>
          <el-select v-model="namespace" placeholder="请选择命名空间" clearable filterable @change="getList" style="width:200px;margin-right:8px" size="large">
            <el-option v-for="(item) in clusterList[cluster]" :key="item" :label="item" :value="item" />
          </el-select>
          <el-input v-model="searchKey" placeholder="输入名称进行搜索" size="large" :prefix-icon="Search" @change="getList" clearable style="width:200px;" />
        </el-button-group>
      </el-col>
    </el-row>
    <el-table 
      :data="list" 
      v-loading="loading.table"
      element-loading-text="奋力加载中..."
      class="line-height40" 
      style="width:100%;margin-top:10px;min-height:150px">
      <el-table-column prop="name" label="Release名称" min-width="150" />
      <el-table-column prop="chart" label="Chart版本" min-width="150" />
      <el-table-column prop="app_version" label="App版本" width="150" />
      <el-table-column prop="revision" label="Revision" width="150" />
      <el-table-column prop="status" label="状态" width="180">
        <template #default="scope">
          <font-awesome-icon icon="circle" v-if="['deployed','uninstalled'].includes(scope.row.status)" class="runningstatus text-green" />
          <font-awesome-icon icon="circle" v-else-if="scope.row.status=='failed'" class="runningstatus text-red" />
          <font-awesome-icon icon="circle" v-else class="runningstatus text-yellow" />
          {{ scope.row.status }}
        </template>
      </el-table-column>
      <el-table-column prop="updated" label="更新时间" width="210" />
      <el-table-column label="操作" fixed="right" width="130">
        <template #default="scope">
          <el-tooltip effect="dark" placement="top" content="重装">
            <el-button icon="Refresh" circle @click="upgrade(scope.row)"></el-button>
          </el-tooltip>
          <el-tooltip effect="dark" placement="top" content="回滚">
            <el-button icon="RefreshLeft" circle @click="rollback(scope.row)"></el-button>
          </el-tooltip>
          <el-popconfirm title="确定卸载?" confirm-button-text="确认" cancel-button-text="取消" @confirm="uninstall(scope.row)">
            <template #reference>
              <el-button icon="Delete" circle />
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
    <h4>重装Chart</h4>
  </template>
  <template #default>
    <el-form ref="release" :model="form" :rules="rules" label-width="120px">
      <el-form-item label="选择Helm仓库" prop="repo">
        <el-select 
          v-model="form.repo" 
          placeholder="请选择Helm仓库" 
          clearable 
          filterable 
          value-key="name"
          style="width:100%" 
          size="large"
          @change="getChartVersions">
          <el-option v-for="(item,i) in repoList" :key="i" :label="item.name" :value="item" />
        </el-select>
      </el-form-item>
      <el-form-item label="Chart名称" prop="chart_name">
        <el-input v-model="form.chart_name" disabled size="large" />
      </el-form-item>
      <el-form-item label="选择Chart版本" prop="chart_version">
        <el-select 
          v-model="form.chart_version" 
          placeholder="请选择版本" 
          clearable 
          filterable 
          value-key="ChartVersion"
          size="large">
          <el-option v-for="(item,i) in versions" :key="i" :label="item.ChartVersion" :value="item" />
        </el-select>
      </el-form-item>
      <el-form-item label="当前values.yaml">
        <v-ace-editor
          v-model:value="form.values"
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
      </el-form-item>
    </el-form>
  </template>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show.yaml=false" size="large">取消</el-button>
      <el-button type="primary" @click="submitValues(release)" :loading="loading.submit" size="large">提交</el-button>
    </div>
  </template>
</el-drawer>
<el-dialog
  v-model="show.rollback"
  title="版本回滚"
  width="70%">
  <el-table :data="history" class="line-height40" style="width:100%">
    <el-table-column prop="revision" label="版本" width="70" />
    <el-table-column prop="chart" label="Chart版本" min-width="150" />
    <el-table-column prop="app_version" label="APP版本" width="120" />
    <el-table-column prop="status" label="状态" width="180" />
    <el-table-column prop="description" label="信息" min-width="250" />
    <el-table-column prop="updated" label="更新时间" width="170">
      <template #default="scope">
        {{ moment(scope.row.updated).format('YYYY-MM-DD HH:mm:ss') }}
      </template>
    </el-table-column>
    <el-table-column label="操作" fixed="right" width="120">
      <template #default="scope">
        <el-popconfirm title="确定?" confirm-button-text="确认" cancel-button-text="取消" @confirm="submitRollback(scope.row)">
          <template #reference>
            <el-link type="primary" :underline="false" :disabled="loading.rollback">回滚到此版本</el-link>
          </template>
        </el-popconfirm>
      </template>
    </el-table-column>
  </el-table>
</el-dialog>
</template>
<script setup>
import { Refresh, Search, RefreshLeft } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { onBeforeMount, ref, reactive, onMounted } from 'vue'
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
const repoList = ref([])
const all = ref([])
const searchKey = ref("")
const list = ref([])
const pageSize = ref(10)
const pageTotal = ref(0)
const current = ref(1)
const loading = ref({
  table: false,
  submit: false,
  rollback: false
})
const cluster = ref("")
const clusterList = ref({})
const namespace = ref("")
const versions = ref([])
const show = ref({
  yaml: false,
  rollback: false
})
const form = ref({})
const rules = reactive({
  repo: [{required: true, message: '请选择仓库', trigger: 'change'}],
  chart: [{required: true, message: '请选择Chart', trigger: 'change'}],
  chart_version: [{required: true, message: '请选择版本', trigger: 'change'}],
})
const release = ref(null)
const history = ref([])
const currentRelease = ref({})
/* 生命周期函数 */
onBeforeMount(async () => {
  await getClusterList()
  getRepos()
  getList()
})
onMounted(() => {
  if(route.query.cluster !== undefined)
    cluster.value = route.query.cluster
  if(route.query.namespace !== undefined)
    namespace.value = route.query.namespace
})
/* methods */
const getClusterList = async () => {
  clusterList.value = await axios.get(`/lizardcd/server/clusters`)
}
const getList = async () => {
  let url = `/lizardcd/helm/cluster/${cluster.value}/namespace/${namespace.value}/releases`
  if(searchKey.value !== "") url += `?release_name=${searchKey.value}`
  loading.value.table = true
  if(cluster.value !== "" && namespace.value !== "") {
    all.value = await axios.get(url)
    pageTotal.value = all.value.length
    getPage(current.value)
  }
  loading.value.table = false
}
const getPage = async (page) => {
  list.value = all.value.slice((page-1)*pageSize.value, page*pageSize.value)
}
const getRepos = async () => {
  repoList.value = await axios.get(`/lizardcd/helm/repos`)
}
const getChartVersions = async (val) => {
  versions.value = await axios.get(`/lizardcd/helm/repo/${val.name}/${form.value.chart_name}`)
}
const uninstall = async (row) => {
  await axios.post(`/lizardcd/helm/cluster/${cluster.value}/namespace/${namespace.value}/charts/uninstall?release_name=${row.name}`)
  getList()
}
const upgrade = async (row) => {
  form.value = {
    release_name: row.name,
    chart_name: row.chart_name,
    chart_version: row.chart_version,
    revision: row.revision,
    values: await axios.get(`/lizardcd/helm/cluster/${cluster.value}/namespace/${namespace.value}/release/values?release_name=${row.name}&revision=${row.revision}`),
  }
  show.value.yaml = true
}
const submitValues = async (f) => {
  if(!f) return
  await f.validate(async (valid) => {
    if(valid) {
      let params = {
        "repo_url": form.value.repo.url,
        "chart_name": form.value.chart_name,
        "chart_version": form.value.chart_version.ChartVersion||form.value.chart_version,
        "revision": parseInt(form.value.revision),
        "release_name": form.value.release_name,
        "values": form.value.values
      }
      loading.value.submit = true
      try {
        await axios.post(`/lizardcd/helm/cluster/${cluster.value}/namespace/${namespace.value}/charts/upgrade`, params)
        getList()
      }
      catch(e) {}
      finally {
        loading.value.submit = false
        show.value.yaml = false
      }
    }
    else {
      ElMessage.warning('必填项未填完')
    }
  })
}
const rollback = async (row) => {
  history.value = await axios.get(`/lizardcd/helm/cluster/${cluster.value}/namespace/${namespace.value}/release/history?release_name=${row.name}`)
  currentRelease.value = row
  show.value.rollback = true
}
const submitRollback = async (row) => {
  loading.value.rollback = true
  try {
    await axios.post(`/lizardcd/helm/cluster/${cluster.value}/namespace/${namespace.value}/release/rollback`, {
      "release_name": currentRelease.value.name,
      "revision": row.revision
    })
    getList()
  }
  catch(e){}
  finally {
    loading.value.rollback = false
    show.value.rollback = false
  }
}
const handleSizeChange = async (size) => {
  pageSize.value = size
  await getList()
}
</script>
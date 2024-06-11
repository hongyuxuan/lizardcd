<template>
<el-breadcrumb :separator-icon="ArrowRight">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item :to="{ path: '/helm' }">Helm管理</el-breadcrumb-item>
  <el-breadcrumb-item :to="{ path: `/helm/${route.params.repo_name}`}">{{ route.params.repo_name }}</el-breadcrumb-item>
  <el-breadcrumb-item>{{ route.params.chart_name }}</el-breadcrumb-item>
</el-breadcrumb>
<el-menu :default-active="activeIndex" mode="horizontal" @select="handleSelect">
  <el-menu-item index="deploy">Helm包部署</el-menu-item>
  <el-menu-item index="readme">Helm包信息</el-menu-item>
</el-menu>
<keep-alive>
  <el-card v-if="activeIndex==='deploy'">
    <template #header>
      <div class="card-header">
        <span class="card-header-text">{{ route.params.chart_name }}</span>
        <div class="box-tools pull-right" style="top:-5px">
          <el-button type="primary" @click="download()" :loading="loading.download">下载</el-button>
          <el-popconfirm title="确定部署?" confirm-button-text="确认" cancel-button-text="取消" @confirm="onSubmit(deploy)">
            <template #reference>
              <el-button type="primary" :loading="loading.deploy">部署</el-button>
            </template>
          </el-popconfirm>
        </div>
      </div>
    </template>
    <el-form ref="deploy" :model="form" :rules="rules" label-width="120px">
      <el-form-item label="选择Chart版本" prop="chart_version">
        <el-select 
          v-model="form.chart_version" 
          placeholder="请选择版本" 
          clearable 
          filterable 
          value-key="ChartVersion"
          style="width:400px" 
          size="large" 
          @change="selectVersion">
          <el-option v-for="(item,i) in versions" :key="i" :label="item.ChartVersion" :value="item" />
        </el-select>
      </el-form-item>
      <el-form-item label="集群" prop="cluster">
        <el-select 
          v-model="form.cluster" 
          placeholder="请选择集群" 
          clearable 
          filterable 
          style="width:400px" 
          size="large">
          <el-option v-for="(v,k) in clusterList" :key="k" :label="k" :value="k" />
        </el-select>
      </el-form-item>
      <el-form-item label="命名空间" prop="namespace">
        <el-select v-model="form.namespace" placeholder="请选择命名空间" clearable filterable style="width:400px" size="large">
          <el-option v-for="(item) in clusterList[form.cluster]" :key="item" :label="item" :value="item" />
        </el-select>
      </el-form-item>
      <el-form-item label="发布名称">
        <el-input v-model="form.release_name" size="large" placeholder="不填默认等于Chart名" style="width:400px" />
      </el-form-item>
      <el-form-item label="values.yaml">
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
            minLines: 10,
            maxLines: 5000,
          }" />
      </el-form-item>
    </el-form>
  </el-card>
  <v-md-preview :text="readme" v-else-if="activeIndex==='readme' && loading.values===false" />
  <el-empty v-else-if="activeIndex==='readme' && loading.values===true" description="加载中，请稍后……" />
</keep-alive>
<el-backtop :right="100" :bottom="100" />
</template>
<script setup>
import { ArrowRight } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { onBeforeMount, ref, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { axios } from '/src/assets/util/axios'
/* 引入v-ace-editor */
import { VAceEditor } from 'vue3-ace-editor'
import 'ace-builds/src-noconflict/mode-yaml'
import 'ace-builds/src-noconflict/theme-chrome'
import 'ace-builds/src-noconflict/ext-language_tools'
/* 变量定义 */
const route = useRoute()
const router = useRouter()
const clusterList = ref({})
const activeIndex = ref("deploy")
const form = ref({
  values: ""
})
const rules = reactive({
  cluster: [{required: true, message: '请选择集群', trigger: 'change'}],
  namespace: [{required: true, message: '请选择命名空间', trigger: 'change'}],
  chart_version: [{required: true, message: '请选择版本', trigger: 'change'}],
})
const chartInfo = ref({})
const versions = ref([])
const deploy = ref(null)
const readme = ref("")
const loading = ref({
  deploy: false,
  download: false,
  values: false
})
/* 生命周期函数 */
onBeforeMount(async () => {
  getClusterList()
  await getChart()
  getChartVersions()
})
/* methods */
const getClusterList = async () => {
  clusterList.value = await axios.get(`/lizardcd/server/clusters`)
}
const getChart = async () => {
  let response = await axios.get(`/lizardcd/helm/repo/${route.params.repo_name}?chart_name=${route.params.chart_name}`)
  chartInfo.value = response.find(n => n.ChartName === route.params.chart_name)
}
const getChartVersions = async () => {
  versions.value = await axios.get(`/lizardcd/helm/repo/${route.params.repo_name}/${route.params.chart_name}`)
}
const selectVersion = async (val) => {
  loading.value.values = true
  let response = await Promise.all([
    axios.get(`/lizardcd/helm/repo/charts/values?repo_url=${chartInfo.value.RepoUrl}&chart_name=${chartInfo.value.ChartName}&chart_version=${val.ChartVersion}`),
    axios.get(`/lizardcd/helm/repo/charts/readme?repo_url=${chartInfo.value.RepoUrl}&chart_name=${chartInfo.value.ChartName}&chart_version=${val.ChartVersion}`)
  ])
  loading.value.values = false
  form.value.values = response[0]
  readme.value = response[1]
}
const onSubmit = async (f) => {
  if(!f) return
  await f.validate(async (valid) => {
    if(valid) {
      let params = {
        "repo_url": chartInfo.value.RepoUrl,
        "chart_name": route.params.chart_name,
        "chart_version": form.value.chart_version.ChartVersion,
        "release_name": form.value.release_name||route.params.chart_name,
        "values": form.value.values
      }
      loading.value.deploy = true
      try {
        await axios.post(`/lizardcd/helm/cluster/${form.value.cluster}/namespace/${form.value.namespace}/charts/install`, params)
        router.push(`/helm?tab=release&cluster=${form.value.cluster}&namespace=${form.value.namespace}`)
      }
      catch(e) {}
      finally {
        loading.value.deploy = false
      }
    }
    else {
      ElMessage.warning('必填项未填完')
    }
  })
}
const handleSelect = async (key) => {
  activeIndex.value = key
  if(key === 'readme' && readme.value === '') {
    readme.value = '请先选择Chart版本'
  }
}
const download = async () => {
  if(!form.value.chart_version?.ChartVersion) {
    ElMessage.warning('请选择Chart版本')
    return
  }
  loading.value.download = true
  let res = await axios.get(`/lizardcd/helm/repo/charts/download?repo_url=${chartInfo.value.RepoUrl}&chart_name=${route.params.chart_name}&chart_version=${form.value.chart_version.ChartVersion}`, {responseType: 'blob'})
  const reader = new FileReader()
  reader.readAsDataURL(res)
  reader.onload = (e) => {
    let a = document.createElement('a')
    a.download = `${route.params.chart_name}-${form.value.chart_version.ChartVersion}.tgz`
    a.href = e.target.result
    document.body.appendChild(a)
    a.click()
    document.body.removeChild(a)
    ElMessage.success({message: '下载文件成功'})
    loading.value.download = false
  }
  reader.onerror = (e) => {
    ElMessage.error({message: '下载文件失败: '+e})
  }
  reader.onabort = (e) => {
    ElMessage.error({message: '下载文件中断: '+e})
  }
}
</script>
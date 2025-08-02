<template>
<el-breadcrumb :separator-icon="ArrowRight">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item :to="{ path: '/application' }">应用管理</el-breadcrumb-item>
  <el-breadcrumb-item>应用发布</el-breadcrumb-item>
</el-breadcrumb>
<el-row :gutter="15">
  <el-col :span="12">
    <el-card>
      <template #header>
        <div class="card-header">
          <span class="card-header-text">发布信息</span>
        </div>
      </template>
      <el-form ref="release" :model="form" label-width="100px">
        <el-form-item label="选择应用">
          <el-select 
            v-model="form.app_name"
            placeholder="请输入关键词搜索" 
            value-key="id" 
            clearable 
            filterable 
            remote 
            :remote-method="searchApp" 
            @change="selectApp"
            :loading="loading.searchapp" 
            style="width:100%" 
            size="large">
            <el-option v-for="item in appList" :key="item.id" :label="item.app_name" :value="item" />
          </el-select>
        </el-form-item>
        <el-form-item label="发布策略">
          <el-radio-group v-model="form.policy">
            <el-radio value="same" size="large">所有工作负载使用相同镜像</el-radio>
            <el-radio value="different" size="large">不同工作负载使用不同镜像</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="选择制品" v-if="form.policy==='same'">
          <el-select 
            v-model="form.artifact_url" 
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
          <myTips type="info">如因网络问题，无法在线获取DockerHub的镜像，可直接将镜像地址填写于此</myTips>
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
    </el-card>
  </el-col>
  <el-col :span="12">
    <el-card>
      <template #header>
        <div class="card-header">
          <span class="card-header-text">部署目标</span>
        </div>
      </template>
      <el-form-item :label="form.app_name.deploy_type==='容器'?'工作负载':'目标地址'" v-if="form.app_name&&form.app_name.deploy_type!=='HTTP'">
        <el-row style="margin-bottom:10px;width:100%;">
          <el-check-tag v-for="item in labels" :key="item" :checked="checkedLabels[item]" @change="checkedLabels[item]=!checkedLabels[item];filterLabels()" style="margin-right:5px;">{{ item }}</el-check-tag>
        </el-row>
        <el-card v-for="(m,index) in form.app_name.workload" :key="index" style="width:100%">
          <template #header>
            <div class="card-header">
              <span v-if="form.app_name.deploy_type==='容器'">工作负载 {{ index+1 }}</span>
              <span v-else-if="form.app_name.deploy_type==='虚拟机'">目标服务器 {{ index+1 }}</span>
              <span v-else>目标服务器 {{ index+1 }}</span>
            </div>
          </template>
          <el-form label-width="100px">
            <el-form-item label="容器集群" v-if="form.app_name.deploy_type==='容器'">
              <el-input v-model="m.cluster" size="large" disabled />
            </el-form-item>
            <el-form-item label="命名空间" v-if="form.app_name.deploy_type==='容器'">
              <el-input v-model="m.namespace" size="large" disabled />
            </el-form-item>
            <el-form-item label="工作负载类型" v-if="form.app_name.deploy_type==='容器'">
              <el-radio-group v-model="m.workload_type" disabled>
                <el-radio-button label="deployments" value="deployments" />
                <el-radio-button label="statefulsets" value="statefulsets" />
                <el-radio-button label="jobs" value="jobs" />
                <el-radio-button label="cronjobs" value="cronjobs" />
                <el-radio-button label="yaml" value="yaml" />
              </el-radio-group>
            </el-form-item>
            <el-form-item :label="form.app_name.deploy_type==='容器'?'工作负载名称':'IP'" v-if="m.workload_type!=='yaml'">
              <el-input v-model="m.workload_name" size="large" disabled />
            </el-form-item>
            <el-form-item label="容器名称" v-if="form.app_name.deploy_type==='容器'&&m.workload_type!=='yaml'">
              <el-input v-model="m.container_name" size="large" disabled />
            </el-form-item>
            <el-form-item label="YAML配置" v-if="form.app_name.deploy_type==='容器'&&m.workload_type==='yaml'">
              <v-ace-editor
                v-model:value="fileContent"
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
                  wrap: true
                }" />
              <el-upload
                ref="upload"
                action="#"
                :auto-upload="false"
                :on-change="selectFile"
                :on-exceed="handleExceed"
                :limit="1"
                :show-file-list="false"
                accept=".yml,.yaml"
                style="z-index:1"
                ><el-button type="primary" size="small" style="margin-top:5px">选择文件</el-button>
              </el-upload>
            </el-form-item>
            <el-form-item label="标签">
              <el-tag v-for="item in m.labels" size="large">{{ item }}</el-tag>
            </el-form-item>
            <el-form-item label="选择制品" v-if="form.policy==='different'">
              <el-select 
                v-model="m.artifact_url" 
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
            <el-form-item label="是否启用">
              <el-switch v-model="m.enable" />
            </el-form-item>
          </el-form>
        </el-card>
        <el-row>
          <el-tooltip content="全部启用">
            <el-button circle @click="setWorkloadEnable(form,true)"><font-awesome-icon icon="toggle-on" /></el-button>
          </el-tooltip>
          <el-tooltip content="全部停用">
            <el-button circle @click="setWorkloadEnable(form,false)"><font-awesome-icon icon="toggle-off" /></el-button>
          </el-tooltip>
        </el-row>
      </el-form-item>
    </el-card>
  </el-col>
</el-row>
<el-footer style="position:fixed;right:0;bottom:0;width:100%;background-color: #fff;margin-left:-15px;box-shadow: var(--el-box-shadow-light);text-align: right;z-index:5">
  <el-button @click="confirmDeploy(true)">提交</el-button>
  <el-button @click="confirmDeploy(false)" type="primary">立即执行</el-button>
</el-footer>
</template>

<script setup>
import { ArrowRight, Delete } from '@element-plus/icons-vue'
import { onBeforeMount, ref } from 'vue'
import { ElMessageBox, ElMessage, genFileId } from 'element-plus'
import { useRoute, useRouter } from 'vue-router'
import { axios } from '/src/assets/util/axios'
import _ from 'lodash'
import moment from 'moment'
/* 引入v-ace-editor */
import { VAceEditor } from 'vue3-ace-editor'
import 'ace-builds/src-noconflict/mode-yaml'
import 'ace-builds/src-noconflict/theme-chrome'
import 'ace-builds/src-noconflict/ext-language_tools'
/* 变量定义 */
const route = useRoute()
const router = useRouter()
const form = ref({
  labels: [],
  policy: 'same'
})
const appList = ref([])
const artifactList = ref([])
const loading = ref({
  searchapp: false,
  artifact: false,
})
const labels = ref(['all'])
const checkedLabels = ref({})
const allWorkloads = ref([])
const fileContent = ref('')
const upload = ref(null)
/* 生命周期函数 */
onBeforeMount(async () => {
  if(route.query.app_name) {
    await searchApp(route.query.app_name)
    form.value.app_name = appList.value.find(n => n.app_name === route.query.app_name)
    selectApp()
    if(route.query.artifact_url) {
      form.value.artifact_url = route.query.artifact_url
    }
  }
})
/* methods */
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
  if(!form.value.app_name) {
    ElMessage.warning({message: '请先选择应用'})
    return
  }
  loading.value.artifact = true
  let url = `/lizardcd/server/repo/image/tags?app_name=${encodeURIComponent(form.value.app_name.app_name)}`
  if(query !== "") url += `&tag=${query}`
  let response = await axios.get(url)
  response ||= []
  artifactList.value = _.sortBy(response.map(x => {
    x.last_modified = moment(x.last_modified).format('YYYY-MM-DD HH:mm:ss')
    return x
  }), 'last_modified').reverse()
  loading.value.artifact = false
}
const selectApp = () => {
  labels.value = ['all']
  allWorkloads.value = _.cloneDeep(form.value.app_name.workload)
  for(let x of allWorkloads.value) {
    labels.value = labels.value.concat(x.labels||[])
  }
  labels.value = _.uniq(labels.value)
  if(route.query.labels) {
    for(let x of route.query.labels.split(',')) {
      checkedLabels.value[x] = true
    }
  } else {
    checkedLabels.value.all = true
  }
  filterLabels()
}
const selectFile = (file) => {
  const reader = new FileReader()
  reader.onload = (e) => {
    fileContent.value = e.target.result
  }
  reader.readAsText(file.raw)
}
const handleExceed = (files) => {
  upload.value[0]?.clearFiles()
  const file = files[0]
  file.uid = genFileId
  upload.value[0]?.handleStart(file)
}
const filterLabels = () => {
  if(checkedLabels.value.all === true) {
    form.value.app_name.workload = allWorkloads.value
  } else {
    let selected = []
    for(let [k,v] of Object.entries(checkedLabels.value)) {
      if(v === true) selected.push(k)
    }
    form.value.app_name.workload = allWorkloads.value.filter(n => {
      return _.intersection(n.labels, selected).length > 0
    })
  }
}
const setWorkloadEnable = (form, enable) => {
  for(let x of form.workload||form.app_name.workload) {
    x.enable = enable
  }
}
const confirmDeploy = async (waiting) => {
  ElMessageBox.confirm(
    `确认${waiting?'提交':'发布'}此版本？`,
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    let workloads = form.value.app_name.workload.filter(n => n.enable === true)
    let params = {
      "app_name": form.value.app_name.app_name,
      "task_type": "deploy",
      "trigger_type": "手动触发",
      "labels": form.value.labels.map(x => `${x.key}=${x.value}`),
      "waiting": waiting
    }
    if(form.value.app_name.deploy_type === 'HTTP') {
      params.artifact_url = form.value.artifact_url.artifact_url || form.value.artifact_url
    } else {
      params.workloads = workloads.map(x => {
        return {
          "cluster": x.cluster,
          "namespace": x.namespace,
          "workload_type": x.workload_type,
          "workload_name": x.workload_name || fileContent.value,
          "container_name": x.container_name,
          "artifact_url": form.value.policy === 'same' ? (form.value.artifact_url?.artifact_url || form.value.artifact_url) : (x.artifact_url.artifact_url || x.artifact_url),
          "enable": x.enable
        }
      })
    }
    let response = await axios.post(`/lizardcd/task/run`, params)
    router.push(`/task/history/${response.id}`)
  }).catch((e) => {
    console.warn(e)
  })
}
const addTag = () => {
  form.value.labels.push({key: '', value: ''})
}
const removeTag = (index) => {
  form.value.labels.splice(index, 1)
}
</script>
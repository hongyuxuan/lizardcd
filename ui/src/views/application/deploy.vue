<template>
<el-drawer v-model="show" @opened="afterOpen" direction="rtl" size="700px">
  <template #header>
    <h4>应用发布</h4>
  </template>
  <template #default>
    <el-form ref="release" :model="form" label-width="100px">
      <el-form-item label="选择应用">
        <el-select 
          v-model="form.app_name"
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
      <el-form-item :label="form.app_name.deploy_type==='容器'?'工作负载':'目标地址'" v-if="form.app_name&&form.app_name.deploy_type!=='HTTP'">
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
                <el-radio label="deployments" value="deployments" />
                <el-radio label="statefulsets" value="statefulsets" />
              </el-radio-group>
            </el-form-item>
            <el-form-item :label="form.app_name.deploy_type==='容器'?'工作负载名称':'IP'">
              <el-input v-model="m.workload_name" size="large" disabled />
            </el-form-item>
            <el-form-item label="容器名称" v-if="form.app_name.deploy_type==='容器'">
              <el-input v-model="m.container_name" size="large" disabled />
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
    </el-form>
  </template>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show.deploy=false">取消</el-button>
      <el-button @click="confirmDeploy(true)">提交</el-button>
      <el-button type="primary" @click="confirmDeploy(false)">立即执行</el-button>
    </div>
  </template>
</el-drawer>
</template>
<script setup>
import { ref } from 'vue'
import { ElMessageBox } from 'element-plus'
import { useRouter } from 'vue-router'
import { axios } from '/src/assets/util/axios'
import _ from 'lodash'
import moment from 'moment'
/* 变量定义 */
const props = defineProps({
  applicationInfo: { type: Object }, 
})
const router = useRouter()
const show = ref(false)
const form = ref({})
const appList = ref([])
const artifactList = ref([])
const loading = ref({
  searchapp: false,
  artifact: false,
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
const afterOpen = async () => {
  form.value.app_name = Object.assign({}, props.applicationInfo)
  form.value.policy = 'same'
  let response = await axios.get(`/lizardcd/db/application?search=app_name==${props.applicationInfo.app_name}&size=20`)
  appList.value = response.results
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
          "workload_name": x.workload_name,
          "container_name": x.container_name,
          "artifact_url": form.value.policy === 'same' ? (form.value.artifact_url.artifact_url || form.value.artifact_url) : (x.artifact_url.artifact_url || x.artifact_url)
        }
      })
    }
    let response = await axios.post(`/lizardcd/task/run`, params)
    router.push(`/task/history?id=${response.id}`)
  }).catch((e) => {
    console.warn(e)
  })
}
const setWorkloadEnable = (form, enable) => {
  for(let x of form.workload||form.app_name.workload) {
    x.enable = enable
  }
}
const open = () => {
  show.value = true
}
defineExpose({ open })
</script>
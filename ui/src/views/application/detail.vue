<template>
<el-breadcrumb :separator-icon="ArrowRight" style="min-width:1366px">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item :to="{ path: '/application' }">应用管理</el-breadcrumb-item>
  <el-breadcrumb-item>{{ applicationInfo.app_name }}</el-breadcrumb-item>
</el-breadcrumb>
<el-row :gutter="15" style="min-width:1366px">
  <el-col :span="6">
    <el-card>
      <template #header>
        <div class="card-header">
          <span class="card-header-text"><b>{{ applicationInfo.app_name }}</b></span>
          <div class="box-tools pull-right" style="top:5px">
            <span class="card-header-btn">
              <el-dropdown @command="handleCommand">
                <el-link :underline="false" type="primary">
                  更多操作
                  <el-icon class="el-icon--right">
                    <arrow-down />
                  </el-icon>
                </el-link>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="release" v-if="applicationInfo.deploy_type!=='GitOps'">发布</el-dropdown-item>
                    <el-dropdown-item command="sync" v-else>同步</el-dropdown-item>
                    <el-dropdown-item command="restart" v-if="applicationInfo.deploy_type!=='GitOps'">重启</el-dropdown-item>
                    <el-dropdown-item command="edit">编辑</el-dropdown-item>
                    <el-dropdown-item command="refresh">刷新</el-dropdown-item>
                    <el-dropdown-item command="delete">删除</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </span>
          </div>
        </div>
      </template>
      <el-descriptions :column="1" :label-width="200">
        <el-descriptions-item label="代码仓库">{{ applicationInfo.git_http_url }}</el-descriptions-item>
        <el-descriptions-item label="制品库" v-if="applicationInfo.deploy_type!=='GitOps'">{{ repoInfo.repo_url }}</el-descriptions-item>
        <el-descriptions-item label="仓库/项目" v-if="applicationInfo.deploy_type!=='GitOps'">{{ applicationInfo.repo_name }}</el-descriptions-item>
        <el-descriptions-item label="制品名称" v-if="applicationInfo.deploy_type!=='GitOps'">{{ applicationInfo.image_name }}</el-descriptions-item>
        <el-descriptions-item label="所属租户">{{ applicationInfo.tenant }}</el-descriptions-item>
        <el-descriptions-item label="标签">
          <el-tag v-for="item in applicationInfo.tags" :key="item" size="large" style="margin-bottom:5px">{{item}}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="部署方式">{{ applicationInfo.deploy_type }}</el-descriptions-item>
        <el-descriptions-item label="开启流量控制" v-if="applicationInfo.deploy_type==='容器'">{{ applicationInfo.enable_traffic_control?'是':'否' }}</el-descriptions-item>
        <el-descriptions-item label="开启自动构建">{{ applicationInfo.enable_build?'是':'否' }}</el-descriptions-item>
        <el-descriptions-item label="构建模板" v-if="applicationInfo.enable_build">
          <el-link :underline="false" type="primary" @click="show.template=true">点我查看</el-link>
        </el-descriptions-item>
        <el-descriptions-item label="构建脚本" v-if="applicationInfo.enable_build">
          <el-link :underline="false" type="primary" @click="show.build_script=true">点我查看</el-link>
        </el-descriptions-item>
        <el-descriptions-item label="版本获取脚本" v-if="applicationInfo.enable_build">
          <el-link :underline="false" type="primary" @click="show.version_script=true">点我查看</el-link>
        </el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ applicationInfo.update_at }}</el-descriptions-item>
      </el-descriptions>
    </el-card>
    <el-card v-if="applicationInfo.gitops">
      <template #header>
        <div class="card-header">
          <span class="card-header-text"><b>GitOps配置</b></span>
        </div>
      </template>
      <el-descriptions :column="1" :label-width="200">
        <el-descriptions-item label="Git Revision">{{ applicationInfo.gitops.git_revision }}</el-descriptions-item>
        <el-descriptions-item label="CommitId">{{ applicationInfo.gitops.commit_id }}</el-descriptions-item>
        <el-descriptions-item label="Path">{{ applicationInfo.gitops.path }}</el-descriptions-item>
        <el-descriptions-item label="Includes">{{ applicationInfo.gitops.include }}</el-descriptions-item>
        <el-descriptions-item label="Excludes">{{ applicationInfo.gitops.exclude }}</el-descriptions-item>
        <el-descriptions-item label="同步方式">{{ applicationInfo.gitops.sync_type }}</el-descriptions-item>
        <el-descriptions-item label="同步周期" v-if="applicationInfo.gitops.sync_type==='自动'">{{ applicationInfo.gitops.cron }}</el-descriptions-item>
        <el-descriptions-item label="同步参数">
          <el-checkbox v-model="applicationInfo.gitops.prune_on_delete" label="删除应用时同步删除 K8S 资源" disabled />
        </el-descriptions-item>
        <el-descriptions-item label="同步状态">
          <el-popover placement="right-start" :width="300" trigger="hover">
            <template #reference>
              <el-icon v-if="syncInfo.status==='已同步'" class="text-success" style="vertical-align:middle;font-size:18px"><CircleCheckFilled /></el-icon>
              <el-icon v-if="syncInfo.status==='未同步'" class="text-red" style="vertical-align:middle;font-size:18px"><CircleCloseFilled /></el-icon>
            </template>
            <el-descriptions :column="1" size="small">
              <el-descriptions-item label="同步状态">
                <span v-if="syncInfo.status==='已同步'" class="text-success">{{ syncInfo.status }}</span>
                <span v-else class="text-red">{{ syncInfo.status }}</span>
                </el-descriptions-item>
              <el-descriptions-item label="原因" v-if="syncInfo.status!=='已同步'">{{ syncInfo.reason }}</el-descriptions-item>
              <el-descriptions-item label="应用commitId">{{ syncInfo.local_commitId }}</el-descriptions-item>
              <el-descriptions-item label="仓库commitId">{{ syncInfo.remote_commitId }}</el-descriptions-item>
              <el-descriptions-item label="Last Comment">{{ syncInfo.comment }}</el-descriptions-item>
              <el-descriptions-item label="Committer">{{ syncInfo.author }}</el-descriptions-item>
            </el-descriptions>
          </el-popover>
        </el-descriptions-item>
      </el-descriptions>
    </el-card>
  </el-col>
  <el-col :span="18">
    <el-card v-for="(v,k,i) in workloads" :key="i">
      <template #header>
        <div class="card-header">
          <span class="card-header-text"><b>{{ k }}</b></span>
        </div>
      </template>
      <el-row :gutter="15" v-for="(rs, workload_type, j) in v" :key="j">
        <el-col :span="6" v-for="(item, m) in rs" :key="m">
          <el-card shadow="always" body-class="application-workload-body">
            <table class="table" style="margin-bottom:0">
              <tr>
                <td width="45"><font-awesome-icon :icon="resourceIcon[item.workload_type||item.resource_type]" style="font-size:25px" /></td>
                <td>
                  <div style="font-size:15px;line-height:20px">{{ item.workload_name||item.resource_name }}</div>
                  <div style="font-size:12px;color:#8d8c8c;line-height:20px">{{ item.workload_type||item.resource_type }}</div>
                </td>
                <td width="45" style="text-align: right">
                  <el-link type="primary" :underline="false" v-if="applicationInfo.gitops" @click="openManifest(item)"><font-awesome-icon icon="magnifying-glass"/></el-link>
                  <el-link type="primary" :underline="false" :href="`/workload/${item.workload_type||item.resource_type}/${item.workload_name||item.resource_name}?cluster=${item.cluster}&namespace=${item.namespace}`"><font-awesome-icon icon="share" /></el-link>
                </td>
              </tr>
            </table>
          </el-card>
        </el-col>
      </el-row>
    </el-card>
  </el-col>
</el-row>
<addForm ref="refAdd" :form="applicationInfo" :k8scluster="k8scluster" :edit="true" @submit="getInfo" />
<deployForm ref="refDeploy" :applicationInfo="applicationInfo" />
<el-drawer v-model="show.manifest" direction="rtl" size="1000px">
  <template #header>
    <h4>Manifests</h4>
  </template>
  <template #default>
    <el-menu :default-active="activeIndex" mode="horizontal" @select="handleSelect" v-if="applicationInfo.gitops">
      <el-menu-item index="live_manifest">Live Manifest</el-menu-item>
      <el-menu-item index="diff">Diff</el-menu-item>
      <el-menu-item index="desired_manifest">Desired Manifest</el-menu-item>
    </el-menu>
    <div class="box box-item" v-show="activeIndex==='live_manifest'">
      <div class="box-body" style="padding-top:20px">
        <v-ace-editor
          v-model:value="resource.live_manifest"
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
            minLines: 5,
            maxLines: 5000,
            wrap: true
          }" />
      </div>
    </div>
    <div class="box box-item" v-show="activeIndex==='diff'">
      <div class="box-body" style="padding-top:20px">
        <CodeDiff
          :old-string="resource.last_manifest"
          :new-string="resource.live_manifest"
          output-format="side-by-side"
          style="margin:0"
        />
      </div>
    </div>
    <div class="box box-item" v-show="activeIndex==='desired_manifest'">
      <div class="box-body" style="padding-top:20px">
        <v-ace-editor
          v-model:value="resource.desired_manifest"
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
            minLines: 5,
            maxLines: 5000,
            wrap: true
          }" />
      </div>
    </div>
  </template>
</el-drawer>
<el-drawer v-model="show.template" direction="rtl" size="700px">
  <template #header>
    <h4>{{ applicationInfo.template?.name }}</h4>
  </template>
  <v-ace-editor
    v-model:value="applicationInfo.template.content"
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
      minLines: 5,
      maxLines: 5000,
      wrap: true
    }" />
</el-drawer>
<el-drawer v-model="show.build_script" direction="rtl" size="700px">
  <template #header>
    <h4>构建脚本</h4>
  </template>
  <v-ace-editor
    v-model:value="applicationInfo.build_script"
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
      minLines: 5,
      maxLines: 5000,
      wrap: true
    }" />
</el-drawer>
<el-drawer v-model="show.version_script" direction="rtl" size="700px">
  <template #header>
    <h4>版本获取脚本</h4>
  </template>
  <v-ace-editor
    v-model:value="applicationInfo.version_script"
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
      minLines: 5,
      maxLines: 5000,
      wrap: true
    }" />
</el-drawer>
</template>
<script setup>
import { ArrowRight } from '@element-plus/icons-vue'
import axios from 'axios';
import { onBeforeMount, ref } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import addForm from './add.vue'
import deployForm from './deploy.vue'
import { CodeDiff } from 'v-code-diff'
/* 引入v-ace-editor */
import { VAceEditor } from 'vue3-ace-editor'
import 'ace-builds/src-noconflict/mode-yaml'
import 'ace-builds/src-noconflict/mode-json'
import 'ace-builds/src-noconflict/theme-chrome'
import 'ace-builds/src-noconflict/ext-language_tools'
import router from '@/router';
/* 变量定义 */
const route = useRoute()
const applicationInfo = ref({})
const repoInfo = ref({})
const syncInfo = ref({})
const activeIndex = ref('live_manifest')
const workloads = ref({})
const show = ref({
  add: false,
  manifest: false,
  build_script: false,
  version_script: false,
})
const resources = ref([])
const resource = ref({})
const resourceIcon = ref({
  "deployments": "cubes",
  "statefulsets": "cubes",
  "services": "network-wired",
  "vm": "laptop"
})
const refAdd = ref(null)
const refDeploy = ref(null)
const k8scluster = ref({})
/* 生命周期函数 */
onBeforeMount(async () => {
  getInfo()
  getClusterList()
})
/* methods */
const getInfo = async () => {
  applicationInfo.value = await axios.get(`/lizardcd/db/application/${route.params.id}`)
  if(applicationInfo.value.repo_id !== 0) {
    repoInfo.value = await axios.get(`/lizardcd/db/image_repository/${applicationInfo.value.repo_id}`)
  }
  workloads.value = {}
  if(applicationInfo.value.gitops) {
    syncInfo.value = {}
    let response = await axios.get(`/lizardcd/git/sync/status?apps=${applicationInfo.value.app_name}&withResource=true`)
    syncInfo.value = response[applicationInfo.value.app_name]
    resources.value = await axios.get(`/lizardcd/db/application_resource?filter=application_id==${route.params.id}`)
    for(let x of resources.value) {
      if(!workloads.value.hasOwnProperty(`${x.cluster} / ${x.namespace}`)) {
        workloads.value[`${x.cluster} / ${x.namespace}`] = {}
      }
      if(!workloads.value[`${x.cluster} / ${x.namespace}`].hasOwnProperty(x.resource_type)) {
        workloads.value[`${x.cluster} / ${x.namespace}`][x.resource_type] = []
      }
      workloads.value[`${x.cluster} / ${x.namespace}`][x.resource_type].push(x)
    }
  }
  else {
    for(let x of applicationInfo.value.workload) {
      let key = `${x.cluster} / ${x.namespace}`
      if(!x.cluster) key = x.workload_type
      if(!workloads.value.hasOwnProperty(key)) {
        workloads.value[key] = {}
      }
      if(!workloads.value[key].hasOwnProperty(x.workload_type)) {
        workloads.value[key][x.workload_type] = []
      }
      workloads.value[key][x.workload_type].push(x)
    }
  }
}
const getClusterList = async () => {
  k8scluster.value = await axios.get(`/lizardcd/server/clusters`)
}
const handleCommand = async (command) => {
  switch(command) {
    case "release": {
      refDeploy.value.open()
      break
    }
    case "edit": {
      refAdd.value.open()
      break
    }
    case "restart": {
      await ElMessageBox.confirm(`确定重启该应用的 ${applicationInfo.value.workload.length} 个工作负载？`,'警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(async () => {
        let params = {
          "app_name": applicationInfo.value.app_name,
          "task_type": "rollout",
          "trigger_type": "手动触发",
          "workloads": applicationInfo.value.workload.map(x => {
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
    case "sync": {
      let txt = "确定同步以下资源？<br>"
      for(let x of resources.value) {
        txt += `${x.cluster} / ${x.namespace} / ${x.resource_type} / ${x.resource_name}<br>`
      }
      await ElMessageBox.confirm(txt, '警告' , {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
        dangerouslyUseHTMLString: true,
      }).then(async () => {
        await axios.post(`/lizardcd/task/run`, {
          "app_name": applicationInfo.value.app_name,
          "task_type": "synchronize",
          "trigger_type": "手动触发",
          "waiting": false
        })
      }).catch(() =>{})
      break
    }
    case "refresh": {
      getInfo()
      break
    }
    case "delete": {
      await ElMessageBox.confirm('确定删除？','警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(async () => {
        ElMessage.warning({message: `正在删除，请稍后`})
        await axios.delete(`/lizardcd/db/application/${route.params.id}`)
        router.push(`application`)
      }).catch(() =>{})
      break
    }
  }
}
const handleSelect = async (key) => {
  activeIndex.value = key
}
const openManifest = (item) => {
  resource.value = item
  show.value.manifest = true
}
</script>
<style lang="css">
.application-workload-body {
  padding: 10px 10px 10px 20px;
  background-color: rgb(225, 225, 240);
}
</style>
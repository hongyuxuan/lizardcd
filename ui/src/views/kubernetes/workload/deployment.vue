<template>
<el-breadcrumb :separator-icon="ArrowRight" style="min-width:1366px">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>工作负载</el-breadcrumb-item>
  <el-breadcrumb-item :to="{ path: '/kubernetes/workload', query: {cluster: route.query.cluster, namespace: route.query.namespace, tab: 'deployments'} }">部署</el-breadcrumb-item>
  <el-breadcrumb-item>{{ deploymentInfo.metadata?.name }}</el-breadcrumb-item>
</el-breadcrumb>
<el-row :gutter="15" style="min-width:1366px">
  <el-col :span="6">
    <el-card>
      <template #header>
        <div class="card-header">
          <el-text truncated size="large" style="width:250px;font-size:18px" class="text-black"><b>{{ deploymentInfo.metadata?.name }}</b></el-text>
          <el-dropdown @command="handleCommand" class="pull-right" style="top:2px;min-width:75px;">
            <el-link underline="never">
              更多操作
              <el-icon class="el-icon--right">
                <arrow-down />
              </el-icon>
            </el-link>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="restart">重启</el-dropdown-item>
                <el-dropdown-item command="autoscaler">弹性伸缩</el-dropdown-item>
                <el-dropdown-item command="yaml">编辑YAML</el-dropdown-item>
                <el-dropdown-item command="delete">删除</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </template>
      <el-descriptions :column="1" border class="no-color" :label-width="120">
        <el-descriptions-item label="集群">{{ route.query.cluster }}</el-descriptions-item>
        <el-descriptions-item label="命名空间">{{ route.query.namespace }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ deploymentInfo.metadata?.creationTimestamp }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ deploymentInfo.metadata?.lastUpdateTime }}</el-descriptions-item>
      </el-descriptions>
    </el-card>
  </el-col>
  <el-col :span="18">
    <el-menu :default-active="activeIndex" mode="horizontal" @select="handleSelect">
      <el-menu-item index="status">资源状态</el-menu-item>
      <el-menu-item index="labels">标签</el-menu-item>
      <el-menu-item index="annotations">注解</el-menu-item>
      <el-menu-item index="events">事件</el-menu-item>
      <el-menu-item index="autoscaler" v-if="hpa.found">弹性伸缩</el-menu-item>
    </el-menu>
    <div class="box box-item" v-show="activeIndex==='status'">
      <div class="box-body" style="padding-top:20px">
        <el-row class="statistic">
          <el-col :span="8">
            <div class="statistic">
              <div class="statistic__head">Replicas</div>
              <div :class="`statistic__content ${deploymentInfo.status?.readyReplicas<deploymentInfo.spec?.replicas?'text-red':'text-green'}`" >
                {{  deploymentInfo.status?.readyReplicas }} / {{ deploymentInfo.spec?.replicas }}
              </div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="statistic">
              <div class="statistic__head">Resource.Limits ( cpu/memory )</div>
              <div class="statistic__content">
                {{  deploymentInfo.spec?.template.spec.containers[0].resources?.limits?.cpu }} / {{ deploymentInfo.spec?.template.spec.containers[0].resources?.limits?.memory }}
              </div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="statistic">
              <div class="statistic__head">Resource.Requests</div>
              <div class="statistic__content">
                {{  deploymentInfo.spec?.template.spec.containers[0].resources?.requests?.cpu }} / {{ deploymentInfo.spec?.template.spec.containers[0].resources?.requests?.memory }}
              </div>
            </div>
          </el-col>
        </el-row>
      </div>
    </div>
    <podList v-show="activeIndex==='status'" ref="refPods" workloadType="deployments" />
    <div class="box box-item" v-show="activeIndex==='labels'">
      <div class="box-body" style="padding-top:20px">
        <el-descriptions :column="1" border class="no-color" :label-width="150">
          <el-descriptions-item v-for="(v,k,i) in deploymentInfo.metadata?.labels" :index="i" :label="k">{{ v }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </div>
    <div class="box box-item" v-show="activeIndex==='annotations'">
      <div class="box-body" style="padding-top:20px">
        <el-descriptions :column="1" border class="no-color" :label-width="180">
          <el-descriptions-item v-for="(v,k,i) in deploymentInfo.metadata?.annotations" :index="i" :label="k">{{ v }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </div>
    <eventList v-show="activeIndex==='events'" ref="refEvents" resourceType="Deployment" :resourceName="route.params.workload_name" />
    <div class="box box-item" v-if="hpa.found" v-show="activeIndex==='autoscaler'">
      <div class="box-body" style="padding-top:20px">
        <el-row class="statistic">
          <el-col :span="6">
            <div class="statistic">
              <div class="statistic__head">最大副本数量</div>
              <div class="statistic__content" >
                {{ hpa.max_pod }}
              </div>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="statistic">
              <div class="statistic__head">最小副本数量</div>
              <div class="statistic__content" >
                {{ hpa.min_pod }}
              </div>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="statistic">
              <div class="statistic__head">目标CPU使用率</div>
              <div class="statistic__content" >
                {{ hpa.cpu_usage }} %
              </div>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="statistic">
              <div class="statistic__head">目标内存使用率</div>
              <div class="statistic__content" >
                {{ hpa.mem_usage }} %
              </div>
            </div>
          </el-col>
        </el-row>
      </div>
    </div>
  </el-col>
</el-row>
<el-drawer v-model="show.yaml" direction="rtl" size="700px">
  <template #header>
    <h4>编辑YAML</h4>
  </template>
  <template #default>
    <div style="position:relative">
      <el-switch v-model="wrapLine" active-text="自动换行" inactive-text="不自动换行" inline-prompt class="el-switch-hover-right" size="large" />
      <v-ace-editor
        v-model:value="yamlContent"
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
    </div>
  </template>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show.yaml=false">取消</el-button>
      <el-button type="primary" @click="submitYaml">提交</el-button>
    </div>
  </template>
</el-drawer>
<el-drawer v-model="show.autoscaler" direction="rtl" size="650px">
  <template #header>
    <h4>弹性伸缩</h4>
  </template>
  <template #default>
    <el-form :model="form" label-width="130px" >
      <el-form-item>
        <template #label><el-text>最小实例数量 
          <el-tooltip placement="top">
            <template #content>
              当Pod根据监控指标进行弹性伸缩时，保留的最小Pod数量
            </template>
            <el-icon class="text-yellow"><Warning /></el-icon>
          </el-tooltip></el-text>
        </template>
        <el-input-number v-model="form.min_pod" :max="100" :min="1" size="large" />
      </el-form-item>
      <el-form-item>
        <template #label><el-text>最大实例数量 
          <el-tooltip placement="top">
            <template #content>
              当Pod根据监控指标进行弹性伸缩时，扩充的最大Pod数量
            </template>
            <el-icon class="text-yellow"><Warning /></el-icon>
          </el-tooltip></el-text>
        </template>
        <el-input-number v-model="form.max_pod" :max="100" :min="0" size="large" />
      </el-form-item>
      <el-form-item>
        <template #label><el-text>目标CPU利用率 
          <el-tooltip placement="top">
            <template #content>
              当工作负载的所有Pod的平均CPU利用率低于此目标时，<br>将减少Pod数量，直至最小实例数量；<br>
              当工作负载的所有Pod的平均CPU利用率高于此目标时，<br>将增加Pod数量，直至最大实例数量
            </template>
            <el-icon class="text-yellow"><Warning /></el-icon>
          </el-tooltip></el-text>
        </template>
        <el-input-number v-model="form.cpu_usage" :max="100" :min="0" size="large" />
      </el-form-item>
      <el-form-item>
        <template #label><el-text>目标内存利用率 
          <el-tooltip placement="top">
            <template #content>
              当工作负载的所有Pod的平均内存利用率低于此目标时，<br>将减少Pod数量，直至最小实例数量；<br>
              当工作负载的所有Pod的平均内存利用率高于此目标时，<br>将增加Pod数量，直至最大实例数量
            </template>
            <el-icon class="text-yellow"><Warning /></el-icon>
          </el-tooltip></el-text>
        </template>
        <el-input-number v-model="form.mem_usage" :max="100" :min="0" size="large" />
      </el-form-item>          
    </el-form>
  </template>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show.add=false">取消</el-button>
      <el-button type="primary" @click="confirmScaler()">提交</el-button>
    </div>
  </template>
</el-drawer>
</template>

<script setup>
import { ArrowRight } from '@element-plus/icons-vue'
import { onBeforeUnmount, onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRoute } from 'vue-router'
import podList from './podList.vue'
import eventList from '../eventList.vue'
import { axios } from '/src/assets/util/axios'
import moment from 'moment'
/* 引入v-ace-editor */
import { VAceEditor } from 'vue3-ace-editor'
import 'ace-builds/src-noconflict/mode-yaml'
import 'ace-builds/src-noconflict/theme-chrome'
import 'ace-builds/src-noconflict/ext-language_tools'
import router from '@/router'
/* 变量定义 */
const route = useRoute()
const deploymentInfo = ref({})
const activeIndex = ref("status")
const show = ref({
  yaml: false,
  autoscaler: false,
})
const yamlContent = ref("")
const form = ref({max_pod:10,min_pod:1,cpu_usage:50,mem_usage:50})
const timer = ref(null)
const hpa = ref({found: false})
const wrapLine = ref(true)
const refPods = ref(null)
const refEvents = ref(null)
/* 生命周期函数 */
onMounted(async () => {
  doRequest()
  timer.value = setInterval(() => {
    doRequest()
  }, 15000)
  getAutoScaler()
})
onBeforeUnmount(() => {
  if(timer) {
    clearInterval(timer.value)
    timer.value = null
  }
})
/* methods */
const doRequest = () => {
  getDeployment()
  refEvents.value.getEvents()
  refPods.value.getPods()
}
const getDeployment = async () => {
  deploymentInfo.value = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/deployments/${route.params.workload_name}`)
  deploymentInfo.value.metadata.creationTimestamp = moment(deploymentInfo.value.metadata.creationTimestamp).format('YYYY-MM-DD HH:mm:ss')
  deploymentInfo.value.status.readyReplicas ||= 0
  for(let x of deploymentInfo.value.status.conditions) {
    if(x.type === "Available") {
      deploymentInfo.value.metadata.lastUpdateTime = x.lastUpdateTime
    }
  }
  deploymentInfo.value.metadata.lastUpdateTime = moment(deploymentInfo.value.metadata.lastUpdateTime).format('YYYY-MM-DD HH:mm:ss')
}
const handleSelect = async (key) => {
  activeIndex.value = key
}
const handleCommand = async (command) => {
  switch(command) {
    case "restart": {
      await ElMessageBox.confirm('确定重启？','警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(async () => {
        await axios.patch(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/deployments/${route.params.workload_name}/rollout`)
        ElMessage.success({message: '操作成功'})
      }).catch(() =>{})
      break
    }
    case "autoscaler": {
      form.value = Object.assign({}, hpa.value)
      delete form.value.found
      show.value.autoscaler = true
      break
    }
    case "yaml": {
      yamlContent.value = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/deployments/${route.params.workload_name}/yaml`)
      show.value.yaml = true
      break
    }
    case "delete": {
      await ElMessageBox.confirm('确定删除？','警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(async () => {
        await axios.delete(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/deployments/${route.params.workload_name}`)
        ElMessage.success({message: '删除成功'})
        router.push({
          path: `/kubernetes/workload`,
          query: {
            cluster: route.query.cluster,
            namespace: route.query.namespace,
            tab: 'deployments'
          }
        })
      }).catch(() =>{})
      break
    }
  }
}
const submitYaml = async () => {
  await axios.patch(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/apply/yaml`, yamlContent.value, {
    headers: {
      'Content-Type': 'text/plain'
    }
  })
  show.value.yaml = false
  setTimeout(async () => {
    await doRequest()
  }, 1000)
}
const confirmScaler = async () => {
  await axios.put(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/deployments/${route.params.workload_name}/hpa`, {
    "max": form.value.max_pod,
    "min": form.value.min_pod,
    "cpu": form.value.cpu_usage,
    "memory": form.value.mem_usage
  })
  show.value.autoscaler = false
}
const getAutoScaler = async () => {
  let response = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/deployments/${route.params.workload_name}/hpa`)
  if(response.metadata) {
    hpa.value = {
      found: true,
      max_pod: response.spec.maxReplicas,
      min_pod: response.spec.minReplicas,
      cpu_usage: response.spec.metrics.find(n => n.resource.name === 'cpu')?.resource.target.averageUtilization,
      mem_usage: response.spec.metrics.find(n => n.resource.name === 'memory')?.resource.target.averageUtilization
    }
  }
}
</script>

<style scoped>
.my-header {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  gap: 16px;
}
</style>
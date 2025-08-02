<template>
<el-card>
  <template #header>
    <div class="card-header" style="display:block;">
      <div>
        <b class="card-header-text" style="margin-right:20px">{{ route.params.name }}</b>
        <span style="color:#a0a0a0;font-size:15px" v-if="taskRunInfo.status">
          最后更新于 {{moment.duration(moment(taskRunInfo.status.completionTime)-moment()).humanize(true)}}
        </span>
      </div>
      <div style="margin-top:15px" v-if="taskRunInfo.status?.conditions">
        <el-text :type="getColor(taskRunInfo.status.conditions[0]?.reason)">{{ taskRunInfo.status.conditions[0]?.reason }}</el-text>
        <el-text style="margin-left:20px">{{ taskRunInfo.status.conditions[0]?.message }}</el-text>
      </div>
    </div>
  </template>
  <el-row :gutter="15">
    <el-col :span="6">
      <el-menu 
        :default-active="activeName" 
        background-color="#f0f0f0"
        unique-opened
        @select="selectStep"
        @open="selectTask"
        @close="selectTask">
        <el-sub-menu index="1" v-if="taskRunInfo.status?.conditions">
          <template #title>
            <el-icon v-if="taskRunInfo.status.conditions&&taskRunInfo.status.conditions[0].reason==='Succeeded'" :size="20"><CircleCheckFilled style="color:green" /></el-icon>
            <el-icon v-else-if="taskRunInfo.status.conditions&&taskRunInfo.status.conditions[0].reason==='Pending'" :size="20" class="is-loading"><Refresh style="color:#e6a23c" /></el-icon>
            <el-icon v-else-if="taskRunInfo.status.conditions&&taskRunInfo.status.conditions[0].reason==='Running'" :size="20" class="is-loading"><Refresh style="color:#409eff" /></el-icon>
            <el-icon v-else-if="taskRunInfo.status.conditions&&taskRunInfo.status.conditions[0].reason==='TaskRunCancelled'" :size="18"><font-awesome-icon icon="ban" style="color:#a0a0a0" /></el-icon>
            <el-icon v-else :size="20"><CircleCloseFilled style="color:red" /></el-icon>
            <span>{{ taskRunInfo.pipelineTask || taskRunInfo.metadata.name }}</span>
          </template>
          <el-menu-item v-for="(y,i) in taskRunInfo.status.steps" :key="i" :index="y.name">
            <el-icon v-if="y.terminated&&y.terminated.exitCode===0" :size="20"><CircleCheck style="color:green" /></el-icon>
            <el-icon v-else-if="y.terminated&&y.terminated.exitCode!==0&&y.terminated.reason==='TaskRunCancelled'" :size="18"><font-awesome-icon icon="ban" style="color:#a0a0a0" /></el-icon>
            <el-icon v-else-if="y.terminated&&y.terminated.exitCode!==0&&y.terminated.reason!=='TaskRunCancelled'" :size="20"><CircleClose style="color:red" /></el-icon>
            <el-icon v-else-if="y.waiting" :size="20" class="is-loading"><Refresh style="color:#e6a23c" /></el-icon>
            <el-icon v-else-if="(i==0&&y.running)||(i>0&&y.running&&taskRunInfo.status.steps[i-1].terminated)" :size="20" class="is-loading"><Refresh style="color:#409eff" /></el-icon>
            <el-icon v-else :size="20"><RemoveFilled style="color:#a0a0a0" /></el-icon>
            <span>{{ y.name }}</span>
          </el-menu-item>
        </el-sub-menu>
      </el-menu>
    </el-col>
    <el-col :span="18">
      <el-card v-if="!showStep"> <!-- 点击task -->
        <template #header>
          <div class="card-header" style="display:block;" v-if="taskRunInfo.status">
            <span class="card-header-text" style="margin-right:20px"><b>{{ taskRunInfo.metadata.name }}</b></span>
            <el-text v-if="taskRunInfo.status.conditions" :type="getColor(taskRunInfo.status.conditions[0]?.reason)">
              {{ taskRunInfo.status.conditions[0]?.reason }}
            </el-text>
            <el-text v-if="taskRunInfo.status.conditions" style="margin-left:20px" type="info">持续时间 {{ taskRunInfo.expire }}</el-text>
          </div>
        </template>
        <el-tabs v-model="taskActive">
          <el-tab-pane label="参数" name="1" v-if="taskRunInfo.spec&&taskRunInfo.status">
            <el-table  :data="taskRunInfo.spec.params" class="line-height25">
              <el-table-column label="名称" prop="name" min-width="150" />
              <el-table-column label="值" prop="value" min-width="350" />
            </el-table>
          </el-tab-pane>
          <el-tab-pane label="结果" name="2" v-if="taskRunInfo.status">
            <el-table :data="taskRunInfo.status.results" class="line-height25">
              <el-table-column label="名称" prop="name" min-width="150" />
              <el-table-column label="值" prop="value" min-width="350" />
            </el-table>
          </el-tab-pane>
          <el-tab-pane label="状态" name="3" v-if="taskRunInfo.status">
            <v-ace-editor
              v-model:value="statusManifest"
              lang="yaml"
              theme="terminal"
              style="width:100%"
              :options="{
                enableBasicAutocompletion: true,
                enableSnippets: true,
                enableLiveAutocompletion: true,
                tabSize: 2,
                showPrintMargin: false,
                fontSize: 14,
                maxLines: 5000,
                minLines: 1,
              }" />
          </el-tab-pane>
          <el-tab-pane label="Pod" name="4" v-if="taskRunInfo.status">
            <v-ace-editor
              v-model:value="podManifest"
              lang="yaml"
              theme="terminal"
              style="width:100%"
              :options="{
                enableBasicAutocompletion: true,
                enableSnippets: true,
                enableLiveAutocompletion: true,
                tabSize: 2,
                showPrintMargin: false,
                fontSize: 14,
                maxLines: 5000,
                minLines: 1,
              }" />
          </el-tab-pane>
          <el-tab-pane label="事件" name="5" v-if="taskRunInfo.status">
            <el-table :data="podEvents" style="width:100%" class="line-height25">
              <el-table-column prop="type" label="Type">
                <template #default="scope">
                  <el-tag v-if="scope.row.type==='Warning'" type="warning" size="large">{{ scope.row.type }}</el-tag>
                  <el-tag v-else-if="scope.row.type==='Normal'" type="primary" size="large">{{ scope.row.type }}</el-tag>
                  <el-tag v-else type="info" size="large">{{ scope.row.type }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="reason" label="Reason" />
              <el-table-column prop="age" label="Age" />
              <el-table-column prop="source.component" label="From" />
              <el-table-column prop="message" label="Message" min-width="300px" />
            </el-table>
          </el-tab-pane>
        </el-tabs>
      </el-card>
      <el-card v-else> <!-- 点击step -->
        <template #header>
          <div class="card-header" style="display:block;" v-if="!currentStep.notRun">
            <span class="card-header-text" style="margin-right:20px"><b>{{ currentStep.name }}</b></span>
            <el-text v-if="currentStep.terminated&&currentStep.terminated.exitCode===0" type="success">{{ currentStep.terminated.reason }}</el-text>
            <el-text v-else-if="currentStep.terminated&&currentStep.terminated.exitCode!==0" type="danger">{{ currentStep.terminated.reason }}</el-text>
            <el-text v-else-if="currentStep.waiting" type="warning">{{ currentStep.waiting.reason }}</el-text>
            <el-text v-else-if="currentStep.running" type="primary">Running</el-text>
            <el-text style="margin-left:20px" type="info">持续时间 {{ currentStep.expire||'' }}</el-text>
          </div>
          <div class="card-header" style="display:block;" v-else>
            <span class="card-header-text" style="margin-right:20px"><b>{{ currentStep.name }}</b></span>
            <el-text type="info">未运行</el-text>
          </div>
        </template>
        <el-tabs v-model="stepActive" >
          <el-tab-pane label="日志" name="1">
            <pre class="dark fullscreen" ref="preRef">
              <p v-for="(item,i) in logs" :key="i" v-html="item" />
              <div v-if="currentStep.terminated?.exitCode===0" style="color:#24a148;font-weight:bold">步骤已完成</div>
              <div v-if="currentStep.terminated&&currentStep.terminated.exitCode!==0" style="color:#da1e28;font-weight:bold">步骤失败</div>
              <div v-if="currentStep.running"><el-icon class="is-loading" size="18"><Loading /></el-icon></div>
            </pre>
          </el-tab-pane>
          <el-tab-pane label="详情" name="2">
            <v-ace-editor
              v-model:value="stepSpec[currentStep.name]"
              lang="yaml"
              theme="terminal"
              style="width:100%"
              :options="{
                enableBasicAutocompletion: true,
                enableSnippets: true,
                enableLiveAutocompletion: true,
                tabSize: 2,
                showPrintMargin: false,
                fontSize: 14,
                maxLines: 5000,
                minLines: 1,
              }" />
          </el-tab-pane>
        </el-tabs>
      </el-card>
    </el-col>
  </el-row>
</el-card>
<el-backtop :right="70" :bottom="50" />
</template>

<script setup>
import { ArrowRight,Search,Refresh, CircleCheckFilled, CircleCloseFilled, CircleClose, RemoveFilled, Loading } from '@element-plus/icons-vue'
import axios from 'axios';
import { onBeforeMount, onBeforeUnmount, ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useStore } from 'vuex'
import { getDuration } from '@/assets/util/common';
import moment from 'moment'
import _ from 'lodash'
import yaml from 'js-yaml'
import { AnsiUp } from 'ansi_up';
/* 引入v-ace-editor */
import { VAceEditor } from 'vue3-ace-editor'
import ace from 'ace-builds';
import 'ace-builds/src-noconflict/mode-yaml'
import terminalUrl from 'ace-builds/src-noconflict/theme-terminal?url'
ace.config.setModuleUrl('ace/theme/terminal', terminalUrl);
import 'ace-builds/src-noconflict/ext-language_tools'
/* 变量定义 */
const store = useStore()
const activeName = ref("1")
const taskRunInfo = ref({})
const route = useRoute()
const taskActive = ref("1")
const stepActive = ref("1")
const currentStep = ref({})
const showStep = ref(false)
const statusManifest = ref('')
const podManifest = ref('等待Pod资源')
const podEvents = ref([])
const logs = ref([])
const lines = ref(1000)
const lastTotal = ref(0)
const timer = ref(null)
const stepSpec = ref({})
const ansi_up = new AnsiUp()
const regWithHttp = /((https?:\/\/)?(www\.)?(?!\d+\.\d+)([a-zA-Z0-9-]+(\.[a-zA-Z]{2,})+)([\/?#][^\s]*)?)/g
const regWithOutHttp = /((www\.)?(?!\d+\.\d+)([a-zA-Z0-9-]+(\.[a-zA-Z]{2,})+)([\/?#][^\s]*)?)/g
/* 生命周期函数 */
onBeforeMount(async () => {
  await getTaskRun()
  timer.value = setInterval(async () => {
    await getTaskRun()
    if(currentStep.value) await refreshCurrentStep(currentStep.value.name)
  }, 5000)
})
onBeforeUnmount(() => {
  if(timer) {
    clearInterval(timer.value)
    timer.value = null
  }
})
/* methods */
const getTaskRun = async () => {
  let response = await axios.get(`/lizardcd/tekton/cluster/${route.query.cluster}/namespace/${route.query.namespace}/taskruns/${route.params.name}`)
  response.expire = getDuration(response.status.completionTime, response.status.startTime)
  for(let y of response.status.taskSpec?.steps||[]) {
    stepSpec.value[y.name] = yaml.dump(y)
  }
  if(response.metadata.labels && response.metadata.labels.hasOwnProperty('tekton.dev/pipelineTask')) {
    response.pipelineTask = response.metadata.labels['tekton.dev/pipelineTask']
  }
  if(response.status.conditions)
    taskRunInfo.value = response
  if(taskRunInfo.value.status?.podName) {
    podManifest.value = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/pods/${taskRunInfo.value.status.podName}/yaml`)
    let response = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/Pod/${taskRunInfo.value.status.podName}/events`)
    podEvents.value = response.map(x => {
      x.age = moment.duration(moment(x.lastTimestamp)-moment()).humanize(true)
      return x
    })
  }
  statusManifest.value = await axios.get(`/lizardcd/tekton/cluster/${route.query.cluster}/namespace/${route.query.namespace}/taskruns/${route.params.name}/yaml?withStatus=true`)
  if(taskRunInfo.value.status?.completionTime) {
    clearInterval(timer.value)
    timer.value = null
  }
}
const selectTask = async (event) => {
  showStep.value = false
}
const selectStep = async (event) => {
  showStep.value = true
  currentStep.value = null
  logs.value = []
  await refreshCurrentStep(event)
}
const refreshCurrentStep = async (event) => {
  if(taskRunInfo.value.status?.steps) {
    currentStep.value = taskRunInfo.value.status.steps.find(n => n.name === event)
    if(currentStep.value?.terminated) {
      currentStep.value.expire = getDuration(currentStep.value.terminated.finishedAt, currentStep.value.terminated.startedAt)
    } else if(currentStep.value?.running) {
      currentStep.value.expire = getDuration(moment(), currentStep.value.running.startedAt)
    }
    getLogs(lines.value)
    if(currentStep.value) currentStep.value.notRun = false
  } else {
    currentStep.value = {
      notRun: true,
      name: event
    }
    logs.value = ['没有日志输出']
  }
}
const getLogs = async (lines) => {
  if(currentStep.value?.container) {
    let response = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/pods/${taskRunInfo.value.status.podName}/logs?container=${currentStep.value.container}&timestamps=true`)
    logs.value = response.split('\n').map(x => {
      if(regWithHttp.test(x))
        return ansi_up.ansi_to_html(x).replace(regWithHttp, `<a style='color:#78a9ff;text-decoration:underline;' href="$1" target="_blank">$1</a>`)
      else
        return ansi_up.ansi_to_html(x).replace(regWithOutHttp, `<a style='color:#78a9ff;text-decoration:underline;' href="http://$1" target="_blank">$1</a>`)
    })
    lastTotal.value = logs.value.length
    if(response === '') {
      logs.value = ['没有日志输出']
    }
  }
}
const getColor = (reason) => {
  switch(reason) {
    case 'Succeeded': return 'success'
    case 'Pending': return 'warning'
    case 'Running': return 'primary'
    default: return 'danger'
  }
}
</script>
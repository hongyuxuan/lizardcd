<template>
<el-row :gutter="15">
  <el-col :span="6">
    <el-menu 
      :default-active="activeName" 
      background-color="#f0f0f0"
      unique-opened
      :default-openeds="['1']"
      @select="selectStep"
      @open="selectTask"
      @close="selectTask">
      <el-sub-menu index="1" v-if="taskRunInfo.conditions">
        <template #title>
          <!-- approvetsk -->
          <el-icon v-if="taskRunInfo.approvalStatus?.status==='pending'&&taskRunInfo.conditions[0].reason==='Succeeded'" :size="20" class="is-loading"><Refresh style="color:#e6a23c" /></el-icon>
          <el-icon v-else-if="taskRunInfo.approvalStatus?.status==='completed'&&taskRunInfo.approvalStatus?.results.decision==='true'" :size="20"><CircleCheckFilled style="color:green" /></el-icon>
          <el-icon v-else-if="taskRunInfo.approvalStatus?.status==='completed'&&taskRunInfo.approvalStatus?.results.decision==='false'" :size="20"><CircleCloseFilled style="color:red" /></el-icon>
          <!-- approvetsk end -->
          <el-icon v-else-if="taskRunInfo.conditions&&['Succeeded','Completed'].includes(taskRunInfo.conditions[0]?.reason)" :size="20"><CircleCheckFilled style="color:green" /></el-icon>
          <el-icon v-else-if="taskRunInfo.conditions&&['Pending','Started'].includes(taskRunInfo.conditions[0]?.reason)" :size="20" class="is-loading"><Refresh style="color:#e6a23c" /></el-icon>
          <el-icon v-else-if="taskRunInfo.conditions&&taskRunInfo.conditions[0]?.reason==='Running'" :size="20" class="is-loading"><Refresh style="color:#409eff" /></el-icon>
          <el-icon v-else-if="taskRunInfo.conditions&&taskRunInfo.conditions[0]?.reason==='RunTimedOut'" :size="20"><CircleCloseFilled style="color:red" /></el-icon>
          <el-icon v-else-if="taskRunInfo.conditions&&taskRunInfo.conditions[0]?.reason==='TaskRunCancelled'" :size="18"><font-awesome-icon icon="ban" style="color:#a0a0a0" /></el-icon>
          <el-icon v-else-if="!taskRunInfo.conditions" :size="20"><RemoveFilled style="color:#a0a0a0" /></el-icon>
          <el-icon v-else :size="20"><CircleCloseFilled style="color:red" /></el-icon>
          <span>{{ taskRunInfo.displayName || taskRunInfo.taskName }}</span>
        </template>
        <el-menu-item v-for="(y,i) in taskRunInfo.steps" :key="i" :index="y.name">
          <el-icon v-if="y.terminated&&y.terminated.exitCode===0" :size="20"><CircleCheck style="color:green" /></el-icon>
          <el-icon v-else-if="y.terminated&&y.terminated.exitCode!==0&&y.terminated.reason==='TaskRunCancelled'" :size="18"><font-awesome-icon icon="ban" style="color:#a0a0a0" /></el-icon>
          <el-icon v-else-if="y.terminated&&y.terminated.exitCode!==0&&y.terminated.reason!=='TaskRunCancelled'" :size="20"><CircleClose style="color:red" /></el-icon>
          <el-icon v-else-if="y.waiting" :size="20" class="is-loading"><Refresh style="color:#e6a23c" /></el-icon>
          <el-icon v-else-if="(i==0&&y.running)||(i>0&&y.running&&taskRunInfo.steps[i-1].terminated)" :size="20" class="is-loading"><Refresh style="color:#409eff" /></el-icon>
          <el-icon v-else :size="20"><RemoveFilled style="color:#a0a0a0" /></el-icon>
          <span>{{ y.name }}</span>
        </el-menu-item>
      </el-sub-menu>
    </el-menu>
  </el-col>
  <el-col :span="18">
    <el-card v-if="!showStep"> <!-- 点击task -->
      <template #header>
        <div class="card-header" style="display:block;" v-if="taskRunInfo?.uid">
          <span class="card-header-text" style="margin-right:20px"><b>{{ taskRunInfo.name }}</b></span>
          <el-text v-if="taskRunInfo.taskKind==='ApprovalTask'&&taskRunInfo.labels?.decision" :type="getColor2(taskRunInfo.approvalStatus)">{{ taskRunInfo.labels?.decision }}</el-text>
          <el-text v-else-if="taskRunInfo.conditions.length>0" :type="getColor(taskRunInfo.conditions[0]?.reason)">
            {{ taskRunInfo.conditions[0]?.reason }}
          </el-text>
          <el-text v-if="taskRunInfo.conditions" style="margin-left:20px" type="info">持续时间 {{ taskRunInfo.expire }}</el-text>
          <!-- handleMoreTask -->
          <el-dropdown @command="handleMoreTask" v-if="taskRunInfo.conditions" class="pull-right" style="top:3px">
            <el-link underline="never">
              更多操作 <el-icon class="el-icon--right">
                <arrow-down />
              </el-icon>
            </el-link>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item :command="{action:'cancelled',name:taskRunInfo.name}" :disabled="!['Running','Pending'].includes(taskRunInfo.conditions[0]?.reason)">终止运行</el-dropdown-item>
                <el-dropdown-item :command="{action:'rerun',name:taskRunInfo.name}">重新运行</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
        <div class="card-header" style="display:block;" v-else>
          <span class="card-header-text" style="margin-right:20px"><b>{{ taskRunInfo?.taskName }}</b></span>
          <el-text type="info">暂挂中</el-text>
        </div>
      </template>
      <el-tabs v-model="taskActive" v-if="taskRunInfo&&taskRunInfo.taskKind!=='ApprovalTask'">
        <el-tab-pane label="参数" name="1">
          <el-table  :data="taskRunInfo.params" class="line-height25">
            <el-table-column label="名称" prop="name" min-width="150" />
            <el-table-column label="值" prop="value" min-width="350" />
          </el-table>
        </el-tab-pane>
        <el-tab-pane label="结果" name="2" v-if="taskRunInfo.status">
          <el-table :data="taskRunInfo.results" class="line-height25">
            <el-table-column label="名称" prop="name" min-width="150" />
            <el-table-column label="值" prop="value" min-width="350">
              <template #default="scope">
                <div v-html="scope.row.value" />
              </template>
            </el-table-column>
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
      <el-descriptions title="" :column="1" size="large" border v-else-if="taskRunInfo?.taskKind==='ApprovalTask'&&taskRunInfo?.approvalStatus.status==='completed'">
        <el-descriptions-item label="Description" label-width="150">{{ taskRunInfo.approvalSpec?.description }}</el-descriptions-item>
        <el-descriptions-item label="CreationTimestamp" label-width="150">{{ moment(taskRunInfo.startTime).format('YYYY-MM-DD HH:mm:ss.SSS') }}</el-descriptions-item>
        <el-descriptions-item label="Labels" label-width="150">
          <el-tag v-for="(v,k,i) in taskRunInfo.labels||{}" :key="i" size="large" type="warning">{{ k }}={{ v }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="Status" label-width="150">{{ taskRunInfo.approvalStatus.status }}</el-descriptions-item>
        <el-descriptions-item label="Results" label-width="150">{{ taskRunInfo.approvalStatus.results }}</el-descriptions-item>
      </el-descriptions>
      <iframe v-else-if="taskRunInfo?.taskKind==='ApprovalTask'&&taskRunInfo?.approvalStatus?.status==='pending'&&taskRunInfo.conditions[0].reason==='Succeeded'" :src="`${taskRunInfo.approvalStatus.approvalUrl}?user=${route.query.namespace}`" style="width:100%;min-height:650px;border:none" />
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
          <div class="markdown-dark">
            <v-md-preview v-if="isMarkdown" :text="logs" :theme="githubTheme" />
            <pre v-else style="padding:15px 20px" v-html="logs" class="tekton-logs" />
            <div class="logs-footer" v-if="currentStep.terminated?.exitCode===0" style="color:#3eaf7c">步骤已完成</div>
            <div class="logs-footer" v-if="currentStep.terminated&&currentStep.terminated.exitCode!==0" style="color:#da1e28">步骤失败</div>
            <div class="logs-footer" v-if="currentStep.running"><el-icon class="is-loading" size="18"><Loading /></el-icon></div>
          </div>
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
</template>

<script setup>
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import axios from 'axios'
import yaml from 'js-yaml'
import moment from 'moment'
import { getDuration, generateShortId, getColor, getColor2 } from '@/assets/util/common'
import { AnsiUp } from 'ansi_up'
import githubTheme from '@kangc/v-md-editor/lib/theme/github'
/* 引入v-ace-editor */
import { VAceEditor } from 'vue3-ace-editor'
import ace from 'ace-builds';
import 'ace-builds/src-noconflict/mode-yaml'
import terminalUrl from 'ace-builds/src-noconflict/theme-terminal?url'
ace.config.setModuleUrl('ace/theme/terminal', terminalUrl);
import 'ace-builds/src-noconflict/ext-language_tools'
/* 变量定义 */
const route = useRoute()
const router = useRouter()
const taskRunInfo = ref({})
const activeName = ref("1")
const taskActive = ref("1")
const stepActive = ref("1")
const currentStep = ref({})
const showStep = ref(false)
const statusManifest = ref('')
const podManifest = ref('等待Pod资源')
const podEvents = ref([])
const logs = ref("")
const stepSpec = ref({})
const ansi_up = new AnsiUp()
const isMarkdown = ref(true)
/* 生命周期函数 */

/* methods */
const refreshTaskRun = async (info) => {
  taskRunInfo.value = info
  if(taskRunInfo.value.podName) {
    podManifest.value = await axios.get(`/tekton-pipelines/lizardcd/${route.query.namespace}/pods/${taskRunInfo.value.podName}/yaml`)
    let response = await axios.get(`/tekton-pipelines/lizardcd/${route.query.namespace}/pods/${taskRunInfo.value.podName}/events`)
    podEvents.value = response.map(x => {
      x.age = moment.duration(moment(x.lastTimestamp)-moment()).humanize(true)
      return x
    })
  }
  if(taskRunInfo.value.taskKind !== 'ApprovalTask' && taskRunInfo.value.taskKind !== 'PhecdaApproval') {
    let response = await axios.get(`/tekton-pipelines/lizardcd/${route.query.namespace}/taskRuns/${taskRunInfo.value.name}/yaml`)
    if(!response) return
    let taskrun = yaml.load(response)
    statusManifest.value = yaml.dump(taskrun.status)
    for(let x of taskrun.status?.taskSpec?.steps||[]) {
      stepSpec.value[x.name] = yaml.dump(x)
    }
  }
}
const selectTask = async (event) => {
  showStep.value = false
}
const selectStep = async (event) => {
  showStep.value = true
  currentStep.value = null
  logs.value = ""
  await refreshCurrentStep(event)
}
const refreshCurrentStep = async (event) => {
  if(taskRunInfo.value.steps) {
    event ||= taskRunInfo.value.steps[0]?.name
    if(currentStep.value) {
      event = currentStep.value.name
    }
    // currentStep.value = taskRunInfo.value.steps.find(n => n.name === event)
    let stepIndex = taskRunInfo.value.steps.findIndex(n => n.name === event)
    currentStep.value = taskRunInfo.value.steps[stepIndex]
    if(currentStep.value?.terminated) {
      currentStep.value.expire = getDuration(currentStep.value.terminated.finishedAt, currentStep.value.terminated.startedAt)
    } else if(currentStep.value?.running && stepIndex == 0 || (currentStep.value?.running && stepIndex > 0 && taskRunInfo.value.steps[stepIndex-1]?.terminated)) { // 因为一个task执行时所有step都是running，但实际只有第一个step在跑，所以只有stepIndex=0显示expire，或者上一个terminated了下一个才显示expire
      currentStep.value.expire = getDuration(moment(), currentStep.value.running.startedAt)
    }
    getLogs()
    if(currentStep.value) currentStep.value.notRun = false
  } else {
    currentStep.value = {
      notRun: true,
      name: event
    }
    logs.value = '没有日志输出'
  }
}
const getLogs = async () => {
  if(currentStep.value?.container) {
    let response = await axios.get(`/tekton-pipelines/lizardcd/${route.query.namespace}/pods/${taskRunInfo.value.podName}/${currentStep.value.container}/logs`)
    isMarkdown.value = response.includes("#enableMarkdown")
    response = response.replace(/^.*#enableMarkdown.*$\n?/gm, '')
    logs.value = ansi_up.ansi_to_html(response)
    if(response === '') {
      logs.value = '没有日志输出'
    }
  }
}
const handleMoreTask = async (command) => {
  switch(command.action) {
    case "cancelled": {
      await axios.patch(`/lizardcd/tekton/cluster/${route.query.cluster}/namespace/${route.query.namespace}/taskruns/${command.name}`, [
        {
          "op": "replace",
          "path": "/spec/status",
          "value": "TaskRunCancelled"
        }
      ])
      ElMessage.success({message: `终止TaskRun成功`})
      break
    }
    case "rerun": {
      let content = await axios.get(`/tekton-pipelines/lizardcd/${route.query.namespace}/taskRuns/${command.name}/yaml`)
      let params = yaml.load(content)
      params.metadata.name = params.metadata.labels["tekton.dev/pipelineTask"] || params.metadata.labels["tekton.dev/task"] + "-" + generateShortId(5)
      delete params.metadata.generateName
      if(params.metadata.annotations)
        delete params.metadata.annotations["pipeline.tekton.dev/affinity-assistant"]
      if(params.spec.taskRef) {
        delete params.spec.taskRef.namespace
        if(params.spec.taskRef.resolver) {
          delete params.spec.taskRef.name
        }
      }
      delete params.spec.status
      delete params.status
      try {
        await axios.post(`/lizardcd/tekton/cluster/${route.query.cluster}/namespace/${route.query.namespace}/apply?kind=TaskRun`, {content: yaml.dump(params)})
        ElMessage.success({message: `重新执行TaskRun成功`})
        setTimeout(() => {
          router.push({
            path: `/ci/tekton/taskrun/${params.metadata.name}`,
            query: route.query
          })
        }, 1000)
      } catch(e) {
        ElMessage.error({message: e})
      }
      break
    }
  }
}
const setActiveName = (name) => {
  activeName.value = name
}
const close = () => {
  showStep.value = false
  currentStep.value = null
}
defineExpose({ refreshTaskRun, refreshCurrentStep, setActiveName, close })
</script>
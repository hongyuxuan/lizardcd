<template>
  <el-card>
    <template #header>
      <div class="card-header" style="display:block;">
        <div>
          <b class="card-header-text" style="margin-right:20px">{{ route.params.name }}</b>
          <span style="color:#a0a0a0;font-size:15px" v-if="pipelineRunInfo.status">
            流水线启动于 {{moment.duration(moment(pipelineRunInfo.startTime)-moment()).humanize(true)}}
          </span>
          <el-dropdown @command="handleMore" v-if="pipelineRunInfo.conditions" class="pull-right">
            <el-button size="large">更多操作 <el-icon class="el-icon--right">
              <arrow-down />
            </el-icon></el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item :command="{action:'cancelled'}" :disabled="!['Running','Pending'].includes(pipelineRunInfo.conditions[0]?.reason)">终止运行</el-dropdown-item>
                <el-dropdown-item :command="{action:'rerun'}">重新运行</el-dropdown-item>
                <el-dropdown-item :command="{action:'params'}">查看启动参数</el-dropdown-item>
                <el-dropdown-item :command="{action:'yaml'}">查看YAML</el-dropdown-item>
                <el-dropdown-item :command="{action:'delete'}">删除</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
        <div style="margin-top:15px" v-if="pipelineRunInfo.conditions">
          <el-text :type="getColor(pipelineRunInfo.conditions[0]?.reason)">{{ pipelineRunInfo.conditions[0]?.reason }}</el-text>
          <el-text style="margin-left:20px">{{ pipelineRunInfo.conditions[0]?.message }}</el-text>
        </div>
        <div style="margin-top:10px"><el-tag v-for="(v,k,i) in pipelineRunInfo.annotations||{}" :key="i" size="large" type="warning">{{ k }}={{ v }}</el-tag></div>
        <div style="margin-top:10px"><div><el-tag v-for="(v,k,i) in pipelineRunInfo.labels||{}" :key="i" size="large">{{ k }}={{ v }}</el-tag></div></div>
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
          <el-sub-menu v-for="(item,i) in pipelineRunInfo.taskRuns||[]" :key="i" :index="item.taskName">
            <template #title>
              <!-- approvetsk -->
              <el-icon v-if="item.approvalStatus?.status==='pending'" :size="20" class="is-loading"><Refresh style="color:#e6a23c" /></el-icon>
              <el-icon v-else-if="item.approvalStatus?.status==='completed'&&item.approvalStatus?.results.decision==='true'" :size="20"><CircleCheckFilled style="color:green" /></el-icon>
              <el-icon v-else-if="item.approvalStatus?.status==='completed'&&item.approvalStatus?.results.decision==='false'" :size="20"><CircleCloseFilled style="color:red" /></el-icon>
              <!-- approvetsk end -->
              <el-icon v-else-if="['Succeeded','Completed'].includes(item.conditions&&item.conditions[0]?.reason)" :size="20"><CircleCheckFilled style="color:green" /></el-icon>
              <el-icon v-else-if="item.conditions&&['Pending','Started'].includes(item.conditions[0]?.reason)" :size="20" class="is-loading"><Refresh style="color:#e6a23c" /></el-icon>
              <el-icon v-else-if="item.conditions&&item.conditions[0]?.reason==='Running'" :size="20" class="is-loading"><Refresh style="color:#409eff" /></el-icon>
              <el-icon v-else-if="item.conditions&&item.conditions[0]?.reason==='RunTimedOut'" :size="20"><CircleCloseFilled style="color:red" /></el-icon>
              <el-icon v-else-if="item.conditions&&item.conditions[0]?.reason==='TaskRunCancelled'" :size="18"><font-awesome-icon icon="ban" style="color:#a0a0a0" /></el-icon>
              <el-icon v-else-if="!item.conditions" :size="20"><RemoveFilled style="color:#a0a0a0" /></el-icon>
              <el-icon v-else :size="20"><CircleCloseFilled style="color:red" /></el-icon>
              <span>{{ item.displayName || item.taskName }}</span>
            </template>
            <div v-if="item.steps"><!-- already run task -->
              <el-menu-item v-for="(y,i) in item.steps" :key="i" :index="`${item.taskName}#${y.name}`">
                <el-icon v-if="y.terminated&&y.terminated.exitCode===0" :size="20"><CircleCheck style="color:green" /></el-icon>
                <el-icon v-else-if="y.terminated&&y.terminated.exitCode!==0&&y.terminated.reason==='TaskRunCancelled'" :size="18"><font-awesome-icon icon="ban" style="color:#a0a0a0" /></el-icon>
                <el-icon v-else-if="y.terminated&&y.terminated.exitCode!==0&&y.terminated.reason!=='TaskRunCancelled'" :size="20"><CircleClose style="color:red" /></el-icon>
                <el-icon v-else-if="y.waiting" :size="20" class="is-loading"><Refresh style="color:#e6a23c" /></el-icon>
                <el-icon v-else-if="(i==0&&y.running)||(i>0&&y.running&&item.steps[i-1].terminated)" :size="20" class="is-loading"><Refresh style="color:#409eff" /></el-icon>
                <el-icon v-else :size="20"><RemoveFilled style="color:#a0a0a0" /></el-icon>
                <span>{{ y.name }}</span>
              </el-menu-item>
            </div>
            <div v-else><!-- task not run -->
              <el-menu-item v-for="(y,i) in taskMap[item.taskRefName||item.name]?.spec.steps||[]" :key="i" :index="`${item.taskName}#${y.name}`">
                <el-icon :size="20"><RemoveFilled style="color:#a0a0a0" /></el-icon>
                <span>{{ y.name }}</span>
              </el-menu-item>
            </div>
          </el-sub-menu>
        </el-menu>
      </el-col>
      <el-col :span="18">
        <el-card v-if="!showStep"> <!-- 点击task -->
          <template #header>
            <div class="card-header" style="display:block;" v-if="currentTask?.uid">
              <span class="card-header-text" style="margin-right:20px"><b>{{ currentTask.name }}</b></span>
              <el-text v-if="currentTask.taskKind==='ApprovalTask'" :type="getColor2(currentTask.approvalStatus)">{{ currentTask.labels?.decision }}</el-text>
              <el-text v-else-if="currentTask.conditions.length>0" :type="getColor(currentTask.conditions[0]?.reason)">
                {{ currentTask.conditions[0]?.reason }}
              </el-text>
              <el-text v-if="currentTask.conditions" style="margin-left:20px" type="info">持续时间 {{ currentTask.expire }}</el-text>
              <!-- handleMoreTask -->
              <el-dropdown @command="handleMoreTask" v-if="currentTask.conditions" class="pull-right" style="top:3px">
                <el-link underline="never">
                  更多操作 <el-icon class="el-icon--right">
                    <arrow-down />
                  </el-icon>
                </el-link>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item :command="{action:'cancelled',name:currentTask.name}" :disabled="!['Running','Pending'].includes(currentTask.conditions[0]?.reason)">终止运行</el-dropdown-item>
                    <el-dropdown-item :command="{action:'rerun',name:currentTask.name}">重新运行</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
            <div class="card-header" style="display:block;" v-else>
              <span class="card-header-text" style="margin-right:20px"><b>{{ currentTask?.taskName }}</b></span>
              <el-text type="info">暂挂中</el-text>
            </div>
          </template>
          <el-tabs v-model="taskActive" v-if="currentTask&&currentTask.taskKind!=='ApprovalTask'">
            <el-tab-pane label="参数" name="1" v-if="currentTask.conditions">
              <el-table  :data="currentTask.params" class="line-height25">
                <el-table-column label="名称" prop="name" min-width="150" />
                <el-table-column label="值" prop="value" min-width="350" />
              </el-table>
            </el-tab-pane>
            <el-tab-pane label="结果" name="2" v-if="currentTask.conditions">
              <el-table :data="currentTask.results" class="line-height25">
                <el-table-column label="名称" prop="name" min-width="150" />
                <el-table-column label="值"min-width="350">
                  <template #default="props"><div v-html="props.row.value" /></template>
                </el-table-column>
              </el-table>
            </el-tab-pane>
            <el-tab-pane label="状态" name="3" v-if="currentTask.conditions">
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
            <el-tab-pane label="Pod" name="4" v-if="currentTask.conditions">
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
            <el-tab-pane label="事件" name="5" v-if="currentTask.conditions">
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
          <el-descriptions title="" :column="1" size="large" border v-else-if="currentTask?.taskKind==='ApprovalTask'&&currentTask?.approvalStatus.status==='completed'">
            <el-descriptions-item label="Description" label-width="150">{{ currentTask.approvalSpec?.description }}</el-descriptions-item>
            <el-descriptions-item label="CreationTimestamp" label-width="150">{{ moment(currentTask.startTime).format('YYYY-MM-DD HH:mm:ss.SSS') }}</el-descriptions-item>
            <el-descriptions-item label="Labels" label-width="150">
              <el-tag v-for="(v,k,i) in currentTask.labels||{}" :key="i" size="large" type="warning">{{ k }}={{ v }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="Status" label-width="150">{{ currentTask.approvalStatus.status }}</el-descriptions-item>
            <el-descriptions-item label="Results" label-width="150">{{ currentTask.approvalStatus.results }}</el-descriptions-item>
          </el-descriptions>
          <iframe v-else-if="currentTask?.taskKind==='ApprovalTask'&&currentTask?.approvalStatus?.status==='pending'" :src="`${currentTask.approvalStatus.approvalUrl}?user=${route.query.namespace}`" style="width:100%;min-height:650px;border:none" />
        </el-card>
        <el-card v-else> <!-- 点击step -->
          <template #header>
            <div class="card-header" style="display:block;" v-if="!currentStep?.notRun">
              <span class="card-header-text" style="margin-right:20px"><b>{{ currentStep?.name||"" }}</b></span>
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
  <el-drawer v-model="show.yaml" direction="rtl" size="700px">
    <template #header>
      <h4>查看YAML</h4>
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
            maxLines: 5000,
            minLines: 10,
            readOnly: true,
          }" />
      </div>
    </template>
  </el-drawer>
  <el-drawer v-model="show.params" direction="rtl" size="700px">
    <template #header>
      <h4>查看启动参数</h4>
    </template>
    <template #default>
      <el-table :data="pipelineRunInfo.params" class="line-height25">
        <el-table-column label="名称" prop="name" min-width="150" />
        <el-table-column label="值" prop="value" min-width="350" />
      </el-table>
    </template>
  </el-drawer>
  </template>
  
  <script setup>
  import axios from 'axios';
  import { onBeforeMount, onBeforeUnmount, ref, computed } from 'vue'
  import { useRoute, useRouter } from 'vue-router'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import { getDuration, generateShortId, getColor, getColor2 } from '@/assets/util/common'
  import { Duration } from 'luxon'
  import moment from 'moment'
  import yaml from 'js-yaml'
  import { AnsiUp } from 'ansi_up'
  /* 引入v-ace-editor */
  import { VAceEditor } from 'vue3-ace-editor'
  import 'ace-builds/src-noconflict/mode-yaml'
  import terminalUrl from 'ace-builds/src-noconflict/theme-terminal?url'
  ace.config.setModuleUrl('ace/theme/terminal', terminalUrl)
  import 'ace-builds/src-noconflict/ext-language_tools'
  import 'ace-builds/src-noconflict/theme-chrome'
  /* 变量定义 */
  const route = useRoute()
  const router = useRouter()
  const pipelineRunInfo = ref({})
  const taskRunMap = ref({})
  const taskMap = ref({})
  const activeName = computed(() => {
    if(pipelineRunInfo.value.taskRuns && pipelineRunInfo.value.taskRuns.length > 0)
      return pipelineRunInfo.value.taskRuns[0].taskName
    return ''
  })
  const showStep = ref(false)
  const statusManifest = ref('')
  const podManifest = ref('等待Pod资源')
  const podEvents = ref([])
  const timer = ref(null)
  const logs = ref([])
  const taskActive = ref("1")
  const stepActive = ref("1")
  const currentTask = ref({})
  const currentStep = ref({})
  const stepSpec = ref({})
  const ansi_up = new AnsiUp()
  const regWithHttp = /((https?:\/\/)?(www\.)?(?!\d+\.\d+)([a-zA-Z0-9-]+(\.[a-zA-Z]{2,})+)([\/?#][^\s]*)?)/g
  const regWithOutHttp = /((www\.)?(?!\d+\.\d+)([a-zA-Z0-9-]+(\.[a-zA-Z]{2,})+)([\/?#][^\s]*)?)/g
  const yamlContent = ref("")
  const wrapLine = ref(true)
  const show = ref({
    yaml: false
  })
  /* 生命周期函数 */
  onBeforeMount(async () => {
    await getPipelineRun()
    await getTasks()
    selectTask(pipelineRunInfo.value.taskRuns[0].taskName)
    timer.value = setInterval(async () => {
      try{
        await getPipelineRun()
        await refreshCurrentTask(currentTask.value.taskName)
        if(currentStep.value) await refreshCurrentStep(currentTask.value.taskName + "#" + currentStep.value.name)
      }
      catch(e) {
        console.error(e)
      }
    }, 5000)
  })
  onBeforeUnmount(() => {
    if(timer) {
      clearInterval(timer.value)
      timer.value = null
    }
  })
  /* methods */
  const getPipelineRun = async () => {
    pipelineRunInfo.value = await axios.get(`/tekton-pipelines/lizardcd/${route.query.namespace}/pipelineRuns/${route.params.name}`)
    for(let x of pipelineRunInfo.value.taskRuns) {
      x.expire = getDuration(x.endTime, x.startTime)
      let arr = x.taskRef.split('/')
      if(arr.length > 1) {
        x.taskRefNamespace = arr[0]
        x.taskRefName = arr[1]
      } else {
        x.taskRefNamespace = undefined
        x.taskRefName = arr[0]
      }
      if(x.retriesStatus) {
        x.displayName += ` ( 重试 ${x.retriesStatus.length} 次 )`
      }
      for(let r of x.results || []) {
        if(!Array.isArray(r.value)) {
          if(r.value.startsWith("http://") || r.value.startsWith("https://")) {
            r.value = `<a style='color:#78a9ff;text-decoration:underline;' href="${r.value}" target="_blank">${r.value}</a>`
          }
        }
      }
      taskRunMap.value[x.taskName] = x
    }
    if(pipelineRunInfo.value.endTime) {
      clearInterval(timer.value)
      timer.value = null
    }
  }
  const getTasks = async () => {
    let tasks = pipelineRunInfo.value.taskRuns.filter(n => n.hasOwnProperty("taskRef") && n.taskKind !== "ApprovalTask")
    let response = await Promise.all(tasks.map(x => {
      return axios.get(`/lizardcd/tekton/cluster/${route.query.cluster}/namespace/${x.taskRefNamespace}/tasks/${x.taskRefName}`)
    }))
    for(let x of response) {
      taskMap.value[x.metadata.name] = x
    }
  }
  const selectTask = async (event) => {
    showStep.value = false
    podManifest.value = '等待Pod资源'
    podEvents.value = []
    await refreshCurrentTask(event)
  }
  const selectStep = async (event) => {
    showStep.value = true
    currentStep.value = null
    logs.value = []
    await refreshCurrentStep(event)
  }
  const refreshCurrentTask = async (event) => {
    currentTask.value = taskRunMap.value[event]
    // get status yaml
    if(currentTask.value&&currentTask.value.name&&currentTask.value.taskKind!=='ApprovalTask') {
      let response = await axios.get(`/tekton-pipelines/lizardcd/${route.query.namespace}/taskRuns/${currentTask.value.name}/yaml`)
      if(!response) return
      let taskrun = yaml.load(response)
      statusManifest.value = yaml.dump(taskrun.status)
      for(let x of taskrun.status?.taskSpec?.steps||[]) {
        stepSpec.value[x.name] = yaml.dump(x)
      }
      // get pod yaml & events
      if(currentTask.value.podName) {
        podManifest.value = await axios.get(`/tekton-pipelines/lizardcd/${route.query.namespace}/pods/${currentTask.value.podName}/yaml`)
        let response = await axios.get(`/tekton-pipelines/lizardcd/${route.query.namespace}/pods/${currentTask.value.podName}/events`)
        podEvents.value = response.map(x => {
          x.age = moment.duration(moment(x.lastTimestamp)-moment()).humanize(true)
          return x
        })
      }
    }
  }
  const refreshCurrentStep = async (event) => {
    if(currentTask.value.steps) {
      currentStep.value = currentTask.value.steps.find(n => currentTask.value.taskName + "#" + n.name === event)
      if(currentStep.value?.terminated) {
        currentStep.value.expire = getDuration(currentStep.value.terminated.finishedAt, currentStep.value.terminated.startedAt)
      } else if(currentStep.value?.running) {
        currentStep.value.expire = getDuration(moment(), currentStep.value.running.startedAt)
      }
      getLogs()
      if(currentStep.value) currentStep.value.notRun = false
    } else {
      currentStep.value = {
        notRun: true,
        name: event.split("#")[1]
      }
      logs.value = ['没有日志输出']
    }
  }
  const getLogs = async () => {
    if(currentStep.value?.container) {
      let response = await axios.get(`/tekton-pipelines/lizardcd/${route.query.namespace}/pods/${currentTask.value.podName}/${currentStep.value.container}/logs`)
      logs.value = response.split('\n').map(x => {
        if(regWithHttp.test(x))
          return ansi_up.ansi_to_html(x).replace(regWithHttp, `<a style='color:#78a9ff;text-decoration:underline;' href="$1" target="_blank">$1</a>`)
        else
          return ansi_up.ansi_to_html(x).replace(regWithOutHttp, `<a style='color:#78a9ff;text-decoration:underline;' href="http://$1" target="_blank">$1</a>`)
      })
      if(response === '') {
        logs.value = ['没有日志输出']
      }
    }
  }
  const handleMore = async (command) => {
    switch(command.action) {
      case "cancelled": {
        await axios.patch(`/lizardcd/tekton/cluster/${route.query.cluster}/namespace/${route.query.namespace}/pipelineruns/${route.params.name}`, [
          {
            "op": "replace",
            "path": "/spec/status",
            "value": "Cancelled"
          }
        ])
        ElMessage.success({message: `终止PipelineRun成功`})
        break
      }
      case "rerun": {
        let content = await axios.get(`/tekton-pipelines/lizardcd/${route.query.namespace}/pipelineRuns/${route.params.name}/yaml`)
        let params = yaml.load(content)
        params.metadata.name = pipelineRunInfo.value.pipelineRef + "-" + generateShortId(5)
        delete params.metadata.generateName
        delete params.spec.status
        try {
          await axios.post(`/lizardcd/tekton/cluster/${route.query.cluster}/namespace/${route.query.namespace}/apply?kind=PipelineRun`, {content: yaml.dump(params)})
          ElMessage.success({message: `重新执行PipelineRun成功`})
          setTimeout(() => {
            router.push({
              path: `/ci/tekton/pipelinerun/${params.metadata.name}`,
              query: route.query
            })
          }, 1000)
        } catch(e) {
          ElMessage.error({message: e})
        }
        break
      }
      case "yaml": {
        let response = await axios.get(`/tekton-pipelines/lizardcd/${route.query.namespace}/pipelineRuns/${route.params.name}/yaml`)
        let o = yaml.load(response)
        if(o.spec?.timeouts?.pipeline) {
          let timeout = parseInt(o.spec.timeouts.pipeline.slice(0, -2))
          const duration = Duration.fromMillis(timeout / 1e6)
          o.spec.timeouts.pipeline = duration.toFormat("h'h'm'm's's'")
        }
        yamlContent.value = yaml.dump(o)
        show.value.yaml = true
        break
      }
      case "params": {
        show.value.params = true
        break
      }
      case "delete": {
        await ElMessageBox.confirm('确定删除？','警告', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning',
        }).then(async () => {
          await axios.delete(`/tekton-pipelines/lizardcd/${route.query.namespace}/pipelineRuns/${route.params.name}`)
          ElMessage.success({message: `删除PipelineRun成功`})
          setTimeout(async () => {
            router.push({
              path: `/ci/tekton`,
              query: {
                tab: "PipelineRun",
                namespace: route.query.namespace
              }
            })
          }, 1000)
        }).catch((e) => {
          console.warn(e)
        })
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
        params.metadata.name = params.metadata.labels["tekton.dev/pipelineTask"] + "-" + generateShortId(5)
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
  </script>
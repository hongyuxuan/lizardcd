<template>
<el-breadcrumb :separator-icon="ArrowRight" style="min-width:1366px">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>应用负载</el-breadcrumb-item>
  <el-breadcrumb-item :to="{ path: '/kubernetes/jobs', query: {cluster: route.query.cluster, namespace: route.query.namespace} }">任务</el-breadcrumb-item>
  <el-breadcrumb-item>{{ jobInfo.metadata?.name }}</el-breadcrumb-item>
</el-breadcrumb>
<el-row :gutter="15" style="min-width:1366px">
  <el-col :span="6">
    <el-card>
      <template #header>
        <div class="card-header">
          <el-text truncated size="large" style="width:250px;font-size:18px" class="text-black"><b>{{ jobInfo.metadata?.name }}</b></el-text>
          <el-dropdown @command="handleCommand" class="pull-right" style="top:2px;min-width:75px;">
            <el-link underline="never">
              更多操作
              <el-icon class="el-icon--right">
                <arrow-down />
              </el-icon>
            </el-link>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="rerun">重新运行</el-dropdown-item>
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
        <el-descriptions-item label="状态">
          <div v-if="jobInfo.status?.active===1"><font-awesome-icon icon="circle" class="text-success" style="font-size:12px" /> 运行中</div>
          <div v-else-if="jobInfo.status?.succeeded===1"><font-awesome-icon icon="circle" class="text-gray" style="font-size:12px" /> 已完成</div>
          <div v-else-if="jobInfo.status?.failed"><font-awesome-icon icon="circle" class="text-red" style="font-size:12px" /> 失败</div>
        </el-descriptions-item>
        <el-descriptions-item label="最大重试次数">{{ jobInfo.spec?.backoffLimit }}</el-descriptions-item>
        <el-descriptions-item label="容器组完成数量">{{ jobInfo.spec?.completions }}</el-descriptions-item>
        <el-descriptions-item label="并行容器组数量">{{ jobInfo.spec?.parallelism }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ jobInfo.metadata?.creationTimestamp }}</el-descriptions-item>
      </el-descriptions>
    </el-card>
  </el-col>
  <el-col :span="18">
    <el-menu :default-active="activeIndex" mode="horizontal" @select="handleSelect">
      <el-menu-item index="history">运行记录</el-menu-item>
      <el-menu-item index="status">资源状态</el-menu-item>
      <el-menu-item index="labels">标签</el-menu-item>
      <el-menu-item index="annotations">注解</el-menu-item>
    </el-menu>
    <div class="box box-item" v-show="activeIndex==='history'">
      <div class="box-body" style="padding-top:20px">
        <el-table :data="history" class="line-height40">
          <el-table-column label="状态">
            <template #default="scope">
              <div v-if="scope.row.active===1"><font-awesome-icon icon="circle" class="text-success" style="font-size:12px" /> 运行中</div>
              <div v-else-if="scope.row.succeeded===1"><font-awesome-icon icon="circle" class="text-gray" style="font-size:12px" /> 已完成</div>
              <div v-else-if="scope.row.failed"><font-awesome-icon icon="circle" class="text-red" style="font-size:12px" /> 失败（{{ scope.row.failed }}）</div>
            </template>
          </el-table-column>
          <el-table-column label="消息" min-width="200">
            <template #default="scope">
              {{ scope.row.conditions?.length>0?scope.row.conditions[0].message||'-':'-' }}
            </template>
          </el-table-column>
          <el-table-column label="开始时间">
            <template #default="scope">{{ moment(scope.row.startTime).format('YYYY-MM-DD HH:mm:ss') }}</template>
          </el-table-column>
          <el-table-column label="结束时间">
            <template #default="scope">
              {{ scope.row.conditions?.length>0?moment(scope.row.conditions[0].lastProbeTime).format('YYYY-MM-DD HH:mm:ss'):'-' }}
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>
    <podList v-show="activeIndex==='status'" ref="refPods" :labelSelector="labelSelector" style="margin-top: 15px;" />
    <div class="box box-item" v-show="activeIndex==='labels'">
      <div class="box-body" style="padding-top:20px">
        <el-descriptions :column="1" border class="no-color" :label-width="150">
          <el-descriptions-item v-for="(v,k,i) in jobInfo.metadata?.labels" :index="i" :label="k">{{ v }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </div>
    <div class="box box-item" v-show="activeIndex==='annotations'">
      <div class="box-body" style="padding-top:20px">
        <el-descriptions :column="1" border class="no-color" :label-width="180">
          <el-descriptions-item v-for="(v,k,i) in jobInfo.metadata?.annotations" :index="i" :label="k">{{ v }}</el-descriptions-item>
        </el-descriptions>
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
</template>

<script setup>
import { ArrowRight } from '@element-plus/icons-vue'
import { onMounted, onBeforeUnmount, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRoute,useRouter } from 'vue-router'
import podList from '../workload/podList.vue'
import { axios } from '/src/assets/util/axios'
import moment from 'moment'
/* 引入v-ace-editor */
import { VAceEditor } from 'vue3-ace-editor'
import 'ace-builds/src-noconflict/mode-yaml'
import 'ace-builds/src-noconflict/theme-chrome'
import 'ace-builds/src-noconflict/ext-language_tools'
/* 变量定义 */
const route = useRoute()
const router = useRouter()
const jobInfo = ref({})
const activeIndex = ref("history")
const show = ref({
  yaml: false
})
const yamlContent = ref("")
const wrapLine = ref(true)
const history = ref([])
const labelSelector = ref("")
const timer = ref(null)
const refPods = ref(null)
/* 生命周期函数 */
onMounted(async () => {
  await getJob()
  refPods.value.getPods()
  timer.value = setInterval(async () => {
    await getJob()
    refPods.value.getPods()
  }, 15000)
})
onBeforeUnmount(() => {
  if(timer) {
    clearInterval(timer.value)
    timer.value = null
  }
})
/* methods */
const getJob = async () => {
  jobInfo.value = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/jobs/${route.params.job_name}`)
  if(!jobInfo.value.metadata) {
    ElMessage.error({message: `未找到任务: ${route.params.job_name}`})
    return
  }
  jobInfo.value.metadata.creationTimestamp = moment(jobInfo.value.metadata.creationTimestamp).format('YYYY-MM-DD HH:mm:ss')
  history.value = [
    jobInfo.value.status
  ]
  let ls = []
  for(let [k,v] of Object.entries(jobInfo.value.spec.selector.matchLabels)) {
    ls.push(`${k}=${v}`)
  }
  labelSelector.value = ls.join(",")
}
const handleSelect = async (key) => {
  activeIndex.value = key
}
const handleCommand = async (command) => {
  switch(command) {
    case "yaml": {
      yamlContent.value = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/jobs/${route.params.job_name}/yaml`)
      show.value.yaml = true
      break
    }
    case "delete": {
      await ElMessageBox.confirm('确定删除？','警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(async () => {
        await axios.delete(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/jobs/${route.params.job_name}`)
        ElMessage.success({message: '删除成功'})
        router.push({
          path: '/kubernetes/jobs',
          query: {
            cluster: route.query.cluster,
            namespace: route.query.namespace
          }
        })
      }).catch(() =>{})
      break
    }
    case "rerun": {
      await axios.patch(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/jobs/${route.params.job_name}/rollout`)
      ElMessage.success({message: '重新运行成功'})
      setTimeout(async () => {
        await getJob()
        refPods.value.getPods()
      }, 1000)
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
  await getJob()
}
</script>
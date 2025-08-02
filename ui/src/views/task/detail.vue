<template>
<el-breadcrumb :separator-icon="ArrowRight">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item :to="{ path: '/task/history' }">任务管理</el-breadcrumb-item>
  <el-breadcrumb-item>任务详情</el-breadcrumb-item>
</el-breadcrumb>
<el-card>
  <template #header>
    <div class="card-header">
      <span class="card-header-text">任务详情</span>
    </div>
  </template>
  <el-descriptions :column="3" border class="no-color" :label-width="120">
    <el-descriptions-item label="任务ID" width="33%">{{ route.params.id }}</el-descriptions-item>
    <el-descriptions-item label="任务类型" width="33%">{{ historyInfo.task_type }}</el-descriptions-item>
    <el-descriptions-item label="初始化时间">{{ historyInfo.init_at?.Valid?moment(historyInfo.init_at.Time).format('YYYY-MM-DD HH:mm:ss'):'-' }}</el-descriptions-item>
    <el-descriptions-item label="应用名称">{{ historyInfo.app_name }}</el-descriptions-item>
    <el-descriptions-item label="触发类型">{{ historyInfo.trigger_type }}</el-descriptions-item>
    <el-descriptions-item label="开始时间">{{ historyInfo.start_at?.Valid?moment(historyInfo.start_at.Time).format('YYYY-MM-DD HH:mm:ss'):'-' }}</el-descriptions-item>
    <el-descriptions-item label="任务状态">
      <div style="width:200px">
        <el-tooltip :content="historyInfo.status" placement="top">
          <el-progress v-if="historyInfo.status=='finished'" :percentage="100" :status="historyInfo.success.Bool?'success':'exception'" />
          <el-progress v-else-if="historyInfo.status=='running'" :percentage="50" color="#e6a23c" :show-text="false" :indeterminate="true" :duration="2" />
          <el-progress v-else :percentage="0" color="#e6a23c" :show-text="false" />
        </el-tooltip>
      </div>
    </el-descriptions-item>
    <el-descriptions-item label="租户">{{ historyInfo.tenant }}</el-descriptions-item>
    <el-descriptions-item label="结束时间">{{ historyInfo.finish_at?.Valid?moment(historyInfo.finish_at.Time).format('YYYY-MM-DD HH:mm:ss'):'-' }}</el-descriptions-item>
    <el-descriptions-item label="标签">
      <el-tag v-for="(item,i) in historyInfo.labels" :key="i">{{ item }}</el-tag>
    </el-descriptions-item>
    <el-descriptions-item label="输出信息">
      <el-text :type="historyInfo.success?.Bool===false?'danger':''">{{ historyInfo.err_message||'-' }}</el-text>
    </el-descriptions-item>
    <el-descriptions-item label="任务耗时">{{ historyInfo.expire||'-' }}</el-descriptions-item>
  </el-descriptions>
</el-card>
<el-card>
  <template #header>
    <div class="card-header">
      <span class="card-header-text">部署目标</span>
      <div class="box-tools pull-right">
        <span class="card-header-btn" @click="getTaskHistory()"><el-icon><Refresh /></el-icon></span>
      </div>
    </div>
  </template>
  <!-- 容器 -->
  <el-table :data="historyWorkload" :cell-style="{'line-height':'23px'}" v-if="appInfo.deploy_type==='容器'">
    <el-table-column type="index" label="#" width="45" />
    <el-table-column prop="workload.cluster" label="集群" width="150px" />
    <el-table-column prop="workload.namespace" label="命名空间" width="150px" />
    <el-table-column prop="workload.workload_type" label="资源类型" width="120px" />
    <el-table-column prop="workload.workload_name" label="资源名称" min-width="180px">
      <template #default="scope">
        <el-link :href="`/kubernetes/${['deployments','statefulsets','daemonsets'].includes(scope.row.workload.workload_type)?'workload/':''}${scope.row.workload.workload_type}/${scope.row.workload.workload_name}?cluster=${scope.row.workload.cluster}&namespace=${scope.row.workload.namespace}`" type="primary" underline="never" target="_blank">{{ scope.row.workload.workload_name }}</el-link>
      </template>
    </el-table-column>
    <el-table-column prop="workload.container_name" label="容器名称" min-width="150px">
      <template #default="scope">{{ scope.row.workload.deploy_type === 'yaml' ? '-' : scope.row.workload.container_name }}</template>
    </el-table-column>
    <el-table-column prop="workload.artifact_url" label="镜像" min-width="300px" />
    <el-table-column label="是否成功" width="90">
      <template #default="scope">
        <el-icon v-if="scope.row.success.Valid===true&&scope.row.success.Bool===true" class="text-success"><Check /></el-icon>
        <el-icon v-else-if="scope.row.success.Valid===true&&scope.row.success.Bool===false" class="text-danger"><Close /></el-icon>
      </template>
    </el-table-column>
    <el-table-column prop="status" label="状态" width="100px">
      <template #default="scope">
        <el-popover v-if="scope.row.status" placement="top" :width="800" trigger="hover">
          <template #reference>
            <el-link  underline="never" type="primary">查看</el-link>
          </template>
          <el-table :data="scope.row.status" v-if="scope.row.status[0]?.pod_name">
            <el-table-column label="容器组" prop="pod_name" width="250" />
            <el-table-column label="Ready" width="70">
              <template #default="props">
                <el-text :type="props.row.ready==='True'?'success':'warning'">{{ props.row.ready }}</el-text>
              </template>
            </el-table-column>
            <el-table-column label="镜像">
              <template #default="props">
                <div v-for="(v,k,i) in props.row.image" :key="i">
                  <span class="text-gray">{{ k }}</span>: <span>{{ v }}</span>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </el-popover>
        <span v-else>-</span>
      </template>
    </el-table-column>
    <el-table-column prop="err_message" label="输出信息" min-width="250px">
      <template #default="props">
        <el-text :type="historyInfo.success?.Bool===false?'danger':''">
          {{ props.row.err_message||'-' }}
        </el-text>
      </template>
    </el-table-column>
    <el-table-column prop="init_at" label="更新时间" width="160">
      <template #default="props">
        {{ moment(props.row.update_at).format('YYYY-MM-DD HH:mm') }}
      </template>
    </el-table-column>
  </el-table>
  <!-- 虚拟机/SSH -->
  <el-table :data="historyWorkload" :cell-style="{'line-height':'23px'}" v-else>
    <el-table-column type="index" label="#" width="45" />
    <el-table-column prop="workload.workload_type" label="目标类型" width="120px" />
    <el-table-column prop="workload.workload_name" label="目标名称" min-width="150px" />
    <el-table-column prop="workload.artifact_url" label="制品" min-width="300px" />
    <el-table-column label="是否成功" width="90">
      <template #default="scope">
        <el-icon v-if="scope.row.success.Valid===true&&scope.row.success.Bool===true" class="text-success"><Check /></el-icon>
        <el-icon v-else-if="scope.row.success.Valid===true&&scope.row.success.Bool===false" class="text-danger"><Close /></el-icon>
      </template>
    </el-table-column>
    <el-table-column prop="status" label="状态" min-width="200px">
      <template #default="scope">
        <div v-if="!scope.row.status">-</div>
        <div v-for="(item,i) in scope.row.status" :key="i" style="line-height:20px;">{{ item }}</div>
      </template>
    </el-table-column>
    <el-table-column prop="err_message" label="输出信息" min-width="250px">
      <template #default="scope">
        <el-text :type="historyInfo.success?.Bool===false?'danger':''">
          {{ scope.row.err_message||'-' }}
        </el-text>
      </template>
    </el-table-column>
    <el-table-column prop="init_at" label="更新时间" width="160">
      <template #default="scope">
        {{ moment(scope.row.update_at).format('YYYY-MM-DD HH:mm') }}
      </template>
    </el-table-column>
  </el-table>
</el-card>
</template>
<script setup>
import { ArrowRight,Refresh } from '@element-plus/icons-vue'
import { axios } from '/src/assets/util/axios'
import { onBeforeMount, onBeforeUnmount, ref } from 'vue'
import { useRoute } from 'vue-router'
import moment from 'moment'
import { faStethoscope } from '@fortawesome/free-solid-svg-icons'
/* 变量定义 */
const route = useRoute()
const historyInfo = ref({})
const historyWorkload = ref([])
const workloadType = ref("deployments")
const timer = ref(null)
const appInfo = ref({})
/* 生命周期函数 */
onBeforeMount(async () => {
  getTaskHistory()
  timer.value = setInterval(() => {
    getTaskHistory()
  }, 5000)
})
onBeforeUnmount(() => {
  if(timer) {
    clearInterval(timer.value)
    timer.value = null
  }
})
/* methods */
const getTaskHistory = async () => {
  historyInfo.value = await axios.get(`/lizardcd/db/task_history/${route.params.id}`)
  historyWorkload.value = historyInfo.value.workloads.map(x => {
    if(x.status === "") return x
    try {
      x.status = JSON.parse(x.status)
    }
    catch {
      x.status = [x.status]
    }
    return x
  })
  if(['terminated','finished'].includes(historyInfo.value.status)) {
    clearInterval(timer.value)
    timer.value = null
  }
  getApplication()
}
const getApplication = async () => {
  let response = await axios.get(`/lizardcd/db/application?filter=app_name==${historyInfo.value.app_name}`)
  if(response.total > 0) {
    appInfo.value = response.results[0]
  }
}
</script>
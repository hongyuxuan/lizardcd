<template>
<el-breadcrumb :separator-icon="ArrowRight">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>任务管理</el-breadcrumb-item>
</el-breadcrumb>
<el-card>
  <template #header>
    <div class="card-header">
      <span class="card-header-text">任务管理</span>
    </div>
  </template>
  <el-row>
    <el-col :span="18">
      <el-button-group>
        <el-button :icon="Refresh" size="large" style="margin-right:5px" @click="getList(current)" />
        <el-input v-model="searchKey" clearable placeholder="输入关键词查询……" @change="getList(1);current=1" style="width:350px;margin-right:5px" size="large">
          <template #prepend>
            <el-select v-model="searchField" placeholder="选择字段" style="width: 115px" size="large">
              <el-option label="应用名" value="app_name" />
              <el-option label="任务ID" value="id" />
            </el-select>
          </template>
        </el-input>
        <el-select 
          v-model="searchLabels"
          multiple 
          clearable 
          allow-create 
          default-first-option 
          :reserve-keyword="false" 
          filterable 
          style="width:400px;float:right;margin-right:5px" 
          placeholder="输入标签过滤……" 
          @change="getList(1);current=1"
          size="large">
          <el-option v-for="item in tagOptions" :key="item" :label="item" :value="item" />
        </el-select>
      </el-button-group>
    </el-col>
    <el-col :span="6">
      <el-date-picker
          style="float:right"
          v-model="timerange"
          type="datetimerange"
          :shortcuts="shortcuts"
          range-separator="To"
          start-placeholder="开始时间"
          end-placeholder="结束时间"
          size="large"
          @change="getList(1);current=1" />
    </el-col>
  </el-row>
  <el-table 
    :data="list" 
    v-loading="loading.table"
    @expand-change="getTaskHistoryWorkload"
    element-loading-text="奋力加载中..."
    class="line-height40" 
    style="width:100%;margin-top:10px">
    <el-table-column type="selection" width="40" />
    <el-table-column type="expand" width="40">
      <template #default="scope">
        <el-table :data="taskHistoryWorkload[scope.row.id]" style="margin-left:80px;" :cell-style="{'line-height':'23px'}">
          <el-table-column prop="workload.cluster" label="集群" width="110px" />
          <el-table-column prop="workload.namespace" label="命名空间" width="130px" />
          <el-table-column prop="workload.workload_type" label="负载/目标类型" width="120px"></el-table-column>
          <el-table-column prop="workload.workload_name" label="负载/目标名称" min-width="150px">
            <template #default="props">
              <el-link :href="`/workload/${props.row.workload.workload_type}/${props.row.workload.workload_name}?cluster=${props.row.workload.cluster}&namespace=${props.row.workload.namespace}`" type="primary" :underline="false" target="_blank">{{ props.row.workload.workload_name }}</el-link>
            </template>
          </el-table-column>
          <el-table-column prop="workload.container_name" label="容器名称" min-width="150px" />
          <el-table-column prop="workload.artifact_url" label="镜像/制品" min-width="200px" />
          <el-table-column prop="status" label="状态" min-width="300px">
            <template #default="props">
              <div v-for="(item,i) in props.row.status" :key="i" style="line-height:20px;">{{ item }}</div>
            </template>
          </el-table-column>
          <el-table-column prop="err_message" label="输出信息" min-width="180px">
            <template #default="props"><span class="text-red">{{ props.row.err_message }}</span></template>
          </el-table-column>
          <el-table-column prop="init_at" label="更新时间" width="160">
            <template #default="props">
              {{ moment(props.row.update_at).format('YYYY-MM-DD HH:mm') }}
            </template>
          </el-table-column>
        </el-table>
      </template>
    </el-table-column>
    <el-table-column prop="app_name" label="应用名称" min-width="150" />
    <el-table-column prop="task_type" label="任务类型" width="100" />
    <el-table-column prop="trigger_type" label="触发类型" width="100" />
    <el-table-column label="标签" min-width="180">
      <template #default="scope">
        <el-tag v-for="item in scope.row.labels" :key="item" size="large">{{item}}</el-tag>
      </template>
    </el-table-column>
    <el-table-column label="执行结果" width="100">
      <template #default="scope">
        <span v-if="scope.row.success.Bool===true&&scope.row.success.Valid===true" style="color:#5cb87a">
          <el-icon><Check /></el-icon>
        </span>
        <span v-else-if="scope.row.success.Bool===false&&scope.row.success.Valid===true" style="color:#f56c6c">
          <el-icon><Close /></el-icon>
        </span>
        <span v-else style="color:#e6a23c">
          <font-awesome-icon icon="circle" class="twinkling" style="font-size:12px " />
        </span>
      </template>
    </el-table-column>
    <el-table-column label="状态" width="150">
      <template #default="scope">
        <el-tooltip effect="dark" placement="top" :content="scope.row.err_message">
          <el-progress v-if="scope.row.status=='initialize'" :percentage="0" color="#e6a23c" :show-text="false" />
          <el-progress v-else-if="scope.row.status=='running'" :percentage="50" color="#e6a23c" :show-text="false" />
          <el-progress v-else-if="scope.row.status=='finished'&&scope.row.success.Bool===true" :percentage="100" color="#5cb87a" :show-text="false" />
          <el-progress :percentage="100" color="#f56c6c" :show-text="false" />
        </el-tooltip>
      </template>
    </el-table-column>
    <el-table-column prop="tenant" label="所属租户" width="120" />
    <el-table-column prop="init_at" label="初始时间" width="160">
      <template #default="scope">
        {{ moment(scope.row.init_at.Time).format('YYYY-MM-DD HH:mm') }}
      </template>
    </el-table-column>
    <el-table-column prop="expire" label="耗时" width="130" />
    <el-table-column prop="Option" label="操作" width="100">
      <template #default="scope">
        <el-tooltip effect="dark" content="回滚到此版本">
          <el-button circle icon="RefreshLeft" @click="redo(scope.row)" />
        </el-tooltip>
        <el-popconfirm title="确认删除？" @confirm="deleteOne(scope.row)">
          <template #reference>
            <el-button :icon="Close" circle :disabled="role!=='admin'" />
          </template>
        </el-popconfirm>
      </template>
    </el-table-column>
  </el-table>
  <el-pagination 
    class="pull-right"
    background 
    v-model:page-size="pageSize"
    :page-sizes="[10, 30, 50, 100]"
    layout="total, sizes, prev, pager, next, jumper" 
    :total="pageTotal"
    @current-change="getList"
    @size-change="handleSizeChange"
    v-model:current-page="current" />
</el-card>
</template>
<script setup>
import { ArrowRight,Refresh,Close } from '@element-plus/icons-vue'
import { onBeforeMount, ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useStore } from 'vuex'
import { ElMessage, ElMessageBox } from 'element-plus'
import { axios } from '/src/assets/util/axios'
import moment from 'moment'
import _ from 'lodash'
/* 变量定义 */
const store = useStore()
const role = computed(() => {
  return store.state.role
})
const route = useRoute()
const tenant = localStorage.tenant
const list = ref([])
const pageSize = ref(10)
const pageTotal = ref(0)
const current = ref(1)
const searchKey = ref("")
const searchLabels = ref([])
const tagOptions = ref([])
const loading = ref({
  table: false
})
const searchField = ref("app_name")
const timerange = ref([])
const shortcuts = [
  {
    text: '最近6小时',
    value: () => {
      return [moment().subtract(6,'hours'), moment()]
    }
  },
  {
    text: '最近1天',
    value: () => {
      return [moment().subtract(1,'days'), moment()]
    }
  },
  {
    text: '最近3天',
    value: () => {
      return [moment().subtract(3,'days'), moment()]
    }
  },
  {
    text: '最近1周',
    value: () => {
      return [moment().subtract(1,'weeks'), moment()]
    }
  },
]
const taskHistoryWorkload = ref({})
/* 生命周期函数 */
onBeforeMount(async () => {
  if(route.query.id) {
    searchField.value = 'id'
    searchKey.value = route.query.id
  }
  getList(1)
});
/* methods */
const getList = async (page) => {
  let url = `page=${page}&size=${pageSize.value}&sort=init_at desc`
  if(searchKey.value === "" || (searchKey.value !== "" && searchField.value !== "id")) {
    if(timerange.value?.length == 2)
      url += `&range=init_at==${moment(timerange.value[0]).format('YYYY-MM-DD HH:mm:ss')},${moment(timerange.value[1]).format('YYYY-MM-DD HH:mm:ss')}`
  }
  let labels = searchLabels.value.map(x => `labels==${x}`)
  if(searchKey.value !== "") {
    url += `&search=${searchField.value}==${searchKey.value}`
    if(labels.length > 0) {
      url += `,${labels.join(",")}`
    }
  } else if(labels.length > 0) {
    url += `&search=${labels.join(",")}`
  }
  loading.value.table = true
  let response = await axios.get(`/lizardcd/db/task_history?${url}`)
  loading.value.table = false
  list.value = response.results?.map(x => {
    tagOptions.value = tagOptions.value.concat(x.labels)
    try {
      x.err_message = JSON.parse(x.err_message).join(",")
    }
    catch {}
    return x
  })
  tagOptions.value = _.uniq(tagOptions.value)
  pageTotal.value = response.total
}
const redo = async (row) => {
  if(row.task_type === 'rollout') {
    ElMessage.warning('回滚只适用于deploy任务类型')
    return
  }
  ElMessageBox.confirm(
    '确认回滚到此版本？',
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    let response = await axios.get(`/lizardcd/db/task_history/${row.id}`)
    let params = {
      "app_name": row.app_name,
      "task_type": row.task_type,
      "labels": row.labels,
      "trigger_type": row.trigger_type,
    }
    if(response.workloads.length > 0 && response.workloads[0].workload.workload_type === 'HTTP') {
      params.artifact_url = response.workloads[0].workload.artifact_url
    } else {
      params.workloads = response.workloads.map(x => {
        return {
          "cluster": x.workload.cluster,
          "namespace": x.workload.namespace,
          "workload_type": x.workload.workload_type,
          "workload_name": x.workload.workload_name,
          "container_name": x.workload.container_name,
          "artifact_url": x.workload.artifact_url
        }
      })
    }
    await axios.post(`/lizardcd/task/run`, params)
    getList(current.value)
  }).catch((e) => {
    console.warn(e)
  })
}
const getTaskHistoryWorkload = async (row) => {
  let response = await axios.get(`/lizardcd/db/task_history/${row.id}`)
  taskHistoryWorkload.value[row.id] = response.workloads.map(x => {
    try {
      x.status = JSON.parse(x.status)
    }
    catch {
      x.status = [x.status]
    }
    return x
  })
}
const deleteOne = async (row) => {
  await axios.delete(`/lizardcd/task/history/${row.id}`)
  getList(current.value)
}
const handleSizeChange = async (size) => {
  pageSize.value = size
  await getList(current.value)
}
</script>
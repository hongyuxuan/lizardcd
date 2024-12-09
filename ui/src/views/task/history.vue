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
    <el-col :span="16">
      <el-button-group style="width:100%">
        <el-button :icon="Refresh" size="large" style="margin-right:5px" @click="getList(current)" />
        <el-input v-model="searchKey" clearable placeholder="输入关键词查询……" @change="getList(1);current=1" style="width:30%;margin-right:5px" size="large">
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
          style="width:35%;margin-right:5px" 
          placeholder="输入标签过滤……" 
          @change="getList(1);current=1"
          size="large">
          <el-option v-for="item in tagOptions" :key="item" :label="item" :value="item" />
        </el-select>
      </el-button-group>
    </el-col>
    <el-col :span="8">
      <el-dropdown @command="handleMore" class="pull-right">
        <el-button size="large">更多操作<el-icon class="el-icon--right"><arrow-down /></el-icon></el-button>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item :command="{action:'deleteBatch'}">批量删除</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
      <el-date-picker
          style="float:right;margin-right:5px;"
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
    @filter-change="filterTable"
    @sort-change="sortTable"
    @selection-change="select"
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
    <el-table-column prop="app_name" label="应用名称" min-width="150">
      <template #default="scope">
        <el-link :underline="false" :href="`/application?app_name=${encodeURIComponent(scope.row.app_name)}`" target="_blank">{{ scope.row.app_name }}</el-link>
      </template>
    </el-table-column>
    <el-table-column prop="task_type" label="任务类型" :filters="taskTypeFilters" column-key="task_type" :filter-multiple="false" width="120" />
    <el-table-column prop="trigger_type" label="触发类型" width="120" />
    <el-table-column label="执行结果" :filters="successFilters" column-key="success" :filter-multiple="false" width="120">
      <template #default="scope">
        <span v-if="scope.row.status==='waiting'" style="color:#5cb87a" />
        <span v-else-if="scope.row.success.Bool===true&&scope.row.success.Valid===true" style="color:#5cb87a">
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
    <el-table-column label="状态" :filters="statusFilters" column-key="status" :filter-multiple="false" width="150">
      <template #default="scope">
        <el-tooltip effect="dark" placement="top" :content="scope.row.err_message||scope.row.status">
          <el-progress v-if="['initialize','waiting','terminated'].includes(scope.row.status)" :percentage="0" color="#e6a23c" :show-text="false" />
          <el-progress v-else-if="scope.row.status=='running'" :percentage="50" color="#e6a23c" :show-text="false" />
          <el-progress v-else-if="scope.row.status=='finished'&&scope.row.success.Bool===true" :percentage="100" color="#5cb87a" :show-text="false" />
          <el-progress :percentage="100" color="#f56c6c" :show-text="false" />
        </el-tooltip>
      </template>
    </el-table-column>
    <el-table-column prop="tenant" label="所属租户" width="120" />
    <el-table-column prop="init_at" label="初始时间" sortable="custom" width="150">
      <template #default="scope">
        {{ scope.row.init_at.Valid ? moment(scope.row.init_at.Time).format('YYYY-MM-DD HH:mm') : '' }}
      </template>
    </el-table-column>
    <el-table-column prop="start_at" label="开始时间" sortable="custom" width="150">
      <template #default="scope">
        {{ scope.row.start_at.Valid ? moment(scope.row.start_at.Time).format('YYYY-MM-DD HH:mm') : '' }}
      </template>
    </el-table-column>
    <el-table-column prop="expire" label="耗时" width="120">
      <template #default="scope">
        {{ scope.row.start_at.Valid ? scope.row.expire : '' }}
      </template>
    </el-table-column>
    <el-table-column prop="Option" label="操作" width="130">
      <template #default="scope">
        <el-popover placement="left" :width="800" trigger="click">
          <template #reference>
            <el-button :icon="Search" circle />
          </template>
          <el-descriptions title="" :column="2">
            <el-descriptions-item label="任务ID">{{ scope.row.id }}</el-descriptions-item>
            <el-descriptions-item label="初始化时间">{{ moment(scope.row.init_at.Time).format('YYYY-MM-DD HH:mm:ss') }}</el-descriptions-item>
            <el-descriptions-item label="应用名称">{{ scope.row.app_name }}</el-descriptions-item>
            <el-descriptions-item label="开始时间">{{ moment(scope.row.start_at.Time).format('YYYY-MM-DD HH:mm:ss') }}</el-descriptions-item>
            <el-descriptions-item label="输出信息">{{ scope.row.err_message }}</el-descriptions-item>
            <el-descriptions-item label="结束时间">{{ moment(scope.row.finish_at.Time).format('YYYY-MM-DD HH:mm:ss') }}</el-descriptions-item>
            <el-descriptions-item label="标签">
              <el-tag v-for="item in scope.row.labels" :key="item" size="large">{{item}}</el-tag>
            </el-descriptions-item>
          </el-descriptions>
        </el-popover>
        <el-tooltip effect="dark" content="执行 / 回滚" placement="top">
          <el-button circle @click="execute(scope.row)" :disabled="scope.row.tenant!==tenant&&role!=='admin'" :icon="ArrowRight" />
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
import { ArrowRight,Refresh,Close,Search } from '@element-plus/icons-vue'
import { onBeforeMount, ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useStore } from 'vuex'
import { ElMessage, ElMessageBox } from 'element-plus'
import { axios } from '/src/assets/util/axios'
import moment from 'moment'
import _ from 'lodash'
/* 变量定义 */
const tenant = localStorage.tenant.split(",")[0]
const store = useStore()
const role = computed(() => {
  return store.state.role
})
const route = useRoute()
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
const statusFilters = ref([
  { text: 'initialize', value: 'initialize'},
  { text: 'waiting', value: 'waiting'},
  { text: 'terminated', value: 'terminated'},
  { text: 'finished', value: 'finished'},
  { text: 'running', value: 'running'},
])
const taskTypeFilters = ref([
  { text: 'rollout', value: 'rollout'},
  { text: 'deploy', value: 'deploy'},
  { text: 'synchronize', value: 'synchronize'},
])
const successFilters = ref([
  { text: '成功', value: '1'},
  { text: '失败', value: '0'},
])
const filterOptions = ref({})
const sort = ref({prop: 'init_at', order: 'descending'})
const selected = ref([])
/* 生命周期函数 */
onBeforeMount(async () => {
  if(route.query.id) {
    searchField.value = 'id'
    searchKey.value = route.query.id
  }
  if(route.query.tag) {
    searchLabels.value = [route.query.tag]
  }
  getList(1)
});
/* methods */
const getList = async (page) => {
  let url = `page=${page}&size=${pageSize.value}&sort=${sort.value.prop} ${sort.value.order==='descending'?'desc':'asc'}`
  if(searchKey.value === "" || (searchKey.value !== "" && searchField.value !== "id")) {
    if(timerange.value?.length == 2)
      url += `&range=${sort.value.prop}==${moment(timerange.value[0]).format('YYYY-MM-DD HH:mm:ss')},${moment(timerange.value[1]).format('YYYY-MM-DD HH:mm:ss')}`
  }
  let labels = searchLabels.value.map(x => `labels==${x}`)
  if(searchKey.value !== "") {
    url += `&search=${searchField.value}==${encodeURIComponent(searchKey.value)}`
    if(labels.length > 0) {
      url += `,${labels.join(",")}`
    }
  } else if(labels.length > 0) {
    url += `&search=${labels.join(",")}`
  }
  let filterParams = []
  if(Object.keys(filterOptions.value).length > 0) {
    for(let [k,v] of Object.entries(filterOptions.value)) {
      filterParams.push(`${k}==${v[0]}`)
    }
  }
  if(filterParams.length > 0) url += `&filter=${filterParams.join(",")}`
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
const execute = async (row) => {
  ElMessageBox.confirm(
    '确认执行此任务？',
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    }
  ).then(async () => {
    await axios.post(`/lizardcd/task/execute/${row.id}`)
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
  await axios.delete(`/lizardcd/db/task_history/${row.id}`)
  ElMessage.success({message: '删除成功'})
  getList(current.value)
}
const handleSizeChange = async (size) => {
  pageSize.value = size
  await getList(current.value)
}
const filterTable = async (filter) => {
  filterOptions.value = Object.assign(filterOptions.value, filter)
  for(let [k,v] of Object.entries(filterOptions.value)) {
    if(v.length === 0)
      delete filterOptions.value[k]
  }
  getList(current.value)
}
const sortTable = async (data) => {
  sort.value = data
  getList(current.value)
}
const handleMore = async (command) => {
  if(selected.value.length === 0) {
    ElMessage.warning({message: '请勾选任务'})
    return
  }
  switch(command.action) {
    case "deleteBatch": {
      await ElMessageBox.confirm('确定删除？','警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(async () => {
        ElMessage.warning({message: '正在删除，请稍后'})
        await Promise.all(selected.value.map(x => {
          return axios.delete(`/lizardcd/db/task_history/${x.id}`)
        }))
        ElMessage.success({message: '删除成功'})
        setTimeout(async () => {
          await getList(current.value)
        }, 500)
      }).catch(() =>{})
      break
    }
  }
}
const select = (val) => {
  selected.value = val
}
</script>
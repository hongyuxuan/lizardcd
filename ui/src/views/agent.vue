<template>
<el-breadcrumb :separator-icon="ArrowRight">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>Agent管理</el-breadcrumb-item>
</el-breadcrumb>
<el-card>
  <template #header>
    <div class="card-header">
      <span class="card-header-text">Agent管理</span>
      <div class="box-tools pull-right">
        <span class="card-header-btn" @click="getList(current)"><el-icon><Refresh /></el-icon></span>
      </div>
    </div>
  </template>
  <el-table :data="list" class="line-height40" :show-header="true" @expand-change="getServiceMeta" style="width:100%;">
    <el-table-column type="expand">
      <template #default="scope">
        <el-table :data="serviceMeta[scope.row.service_name]||[]" style="width:100%;margin-left:50px;">
          <el-table-column prop="ServiceID" label="ServiceID" />
          <el-table-column prop="ServiceMeta" label="ServiceMeta">
            <template #default="props">
              <el-tag v-for="(v,k) of props.row.ServiceMeta">{{ k }}={{ v }}</el-tag>
            </template>
          </el-table-column>
        </el-table>
      </template>  
    </el-table-column>
    <el-table-column prop="service_name" label="Service Key" min-width="160" />
    <el-table-column prop="service_source" label="Service Source" min-width="160" />
  </el-table>
  <el-pagination 
    class="pull-right"
    background 
    v-model:page-size="pageSize"
    :page-sizes="[20, 30, 50, 100]"
    layout="total, sizes, prev, pager, next, jumper" 
    :total="pageTotal"
    @current-change="getPage"
    @size-change="handleSizeChange"
    v-model:current-page="current" />
</el-card>
</template>

<script setup>
import { ArrowRight,Refresh } from '@element-plus/icons-vue'
import { onBeforeMount, ref } from 'vue'
import { axios } from '/src/assets/util/axios.js'
/* 变量定义 */
const list = ref([])
const all = ref([])
const pageSize = ref(20)
const pageTotal = ref(0)
const current = ref(1)
const serviceMeta = ref({})
/* 生命周期函数 */
onBeforeMount(async () => {
  getList()
})
/* methods */
const getList = async () => {
  let response = await axios.get(`/lizardcd/server/services`)
  // for(let x of response) {
    
  // }
  all.value = response
  getPage(current.value)
}
const getPage = async (page) => {
  pageTotal.value = all.value.length
  list.value = all.value.slice((page-1)*pageSize.value, page*pageSize.value)
}
const getServiceMeta = async (row) => {
  let serviceIds = await axios.get(`/lizardcd/server/services/${row.service_name}`)
  serviceMeta.value[row.service_name] = serviceIds.map(y => {
    return {
      ServiceName: y.ServiceName,
      ServiceID: y.ServiceID,
      ServiceMeta: y.ServiceMeta
    }
  })
}
const handleSizeChange = async (size) => {
  pageSize.value = size
  await getList()
}
</script>
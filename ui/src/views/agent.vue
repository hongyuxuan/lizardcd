<template>
<el-breadcrumb :separator-icon="ArrowRight">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>Agent管理</el-breadcrumb-item>
</el-breadcrumb>
<el-card>
  <template #header>
    <div class="card-header">
      <span class="card-header-text">Agent管理</span>
      <span class="pull-right pointer" @click="getList"><el-icon><Refresh /></el-icon></span>
    </div>
  </template>
  <el-table :data="list" class="line-height40" style="width:100%;" :show-header="true">
    <el-table-column type="expand">
      <template #default="scope">
        <el-table :data="scope.row.serviceIds||[]" style="width:100%;margin-left:50px;">
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
</el-card>
</template>

<script setup>
import { ArrowRight,Search,Refresh,EditPen,Delete,Plus,CopyDocument } from '@element-plus/icons-vue'
import { onBeforeMount, ref } from 'vue'
import { axios } from '/src/assets/util/axios.js'
/* 变量定义 */
const list = ref([])
const pageSize = ref(20)
const pageTotal = ref(0)
const current = ref(1)
/* 生命周期函数 */
onBeforeMount(async () => {
  getList()
})
/* methods */
const getList = async () => {
  let response = await axios.get(`/lizardcd/services`)
  for(let x of response) {
    let serviceIds = await axios.get(`/lizardcd/services/${x.service_name}`)
    x.serviceIds = serviceIds.map(y => {
      return {
        ServiceName: y.ServiceName,
        ServiceID: y.ServiceID,
        ServiceMeta: y.ServiceMeta
      }
    })
  }
  list.value = response
}
</script>
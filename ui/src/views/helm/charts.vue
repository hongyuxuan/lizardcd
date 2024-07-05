<template>
<el-breadcrumb :separator-icon="ArrowRight">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item :to="{ path: '/helm' }">Helm管理</el-breadcrumb-item>
  <el-breadcrumb-item>{{ route.params.repo_name }}</el-breadcrumb-item>
</el-breadcrumb>
<el-card>
  <template #header>
    <div class="card-header">
      <span class="card-header-text">包管理 - {{ route.params.repo_name }}</span>
    </div>
  </template>
  <el-row style="margin-bottom:15px">
    <el-col :span="12">
      <el-button-group>
        <el-button :icon="Refresh" size="large" style="margin-right:5px" @click="getList" />
        <el-input v-model="searchKey" placeholder="输入名称进行搜索" size="large" :prefix-icon="Search" @change="getList" clearable style="width:200px;" />
      </el-button-group>
    </el-col>
  </el-row>
  <el-empty v-if="loading" description="加载中，请稍后……" />
  <el-row :gutter="15" v-if="!loading">
    <el-col :span="4" v-for="(item,i) in list" :key="i">
      <el-card shadow="hover" :body-style="{'cursor':'pointer'}" class="card" @click="goto(item)" style="margin-bottom:15px">
        <div class="card-icon">
          <el-image style="width:48px;height:48px" :src="item.Icon">
            <template #error><div class="image-slot"><el-icon><Picture /></el-icon></div></template>
          </el-image>
        </div>
        <div class="card-text">{{item.ChartName}}</div>
        <div class="card-text" style="font-size:12px;color: #a7a4a4;">Chart版本：{{item.ChartVersion}}</div>
        <div class="card-text" style="font-size:12px;color: #a7a4a4;">APP版本：{{item.ChartVersion}}</div>
      </el-card>
    </el-col>
  </el-row>
  <el-pagination 
    v-if="!loading"
    class="pull-right"
    background 
    v-model:page-size="pageSize"
    :page-sizes="[18, 36, 72]"
    layout="total, sizes, prev, pager, next, jumper" 
    :total="pageTotal"
    @size-change="handleSizeChange"
    @current-change="getPage"
    v-model:current-page="current" />
</el-card>
</template>
  
<script setup>
import { Refresh, ArrowRight, Search, Picture } from '@element-plus/icons-vue'
import { onBeforeMount, ref, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { axios } from '/src/assets/util/axios'
/* 变量定义 */
const route = useRoute()
const router = useRouter()
const searchKey = ref("")
const pageSize = ref(18)
const pageTotal = ref(0)
const current = ref(1)
const all = ref([])
const list = ref([])
const loading = ref(false)
/* 生命周期函数 */
onBeforeMount(async () => {
  getList()
})
/* methods */
const getList = async () => {
  loading.value = true
  let url = `/lizardcd/helm/repo/${route.params.repo_name}`
  if(searchKey.value !== "") url += `?chart_name=${searchKey.value}`
  all.value = await axios.get(url)
  loading.value = false
  pageTotal.value = all.value.length
  getPage(current.value)
}
const getPage = async (page) => {
  list.value = all.value.slice((page-1)*pageSize.value, page*pageSize.value)
}
const goto = async (item) => {
  router.push(`${window.location.pathname}/${item.ChartName}`)
}
const handleSizeChange = async (size) => {
  pageSize.value = size
  await getList()
}
</script>
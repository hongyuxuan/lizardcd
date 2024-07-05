<template>
<el-breadcrumb :separator-icon="ArrowRight">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>Helm管理</el-breadcrumb-item>
</el-breadcrumb>
<el-menu :default-active="activeIndex" mode="horizontal" @select="handleSelect">
  <el-menu-item index="repo">Helm仓库</el-menu-item>
  <el-menu-item index="release">Helm包发布</el-menu-item>
</el-menu>
<keep-alive>
  <repo v-if="activeIndex==='repo'" />
  <release v-else-if="activeIndex==='release'" />
</keep-alive>
</template>

<script setup>
import { ArrowRight } from '@element-plus/icons-vue'
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import repo from './repo.vue'
import release from './release.vue'
/* 变量定义 */
const route = useRoute()
const activeIndex = ref("repo")
/* 生命周期函数 */
onMounted(() => {
  if(route.query.tab !== undefined)
    activeIndex.value = route.query.tab
})
/* methods */
const handleSelect = async (key) => {
  activeIndex.value = key
}
</script>
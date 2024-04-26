<template>
<el-breadcrumb :separator-icon="ArrowRight">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>配置</el-breadcrumb-item>
</el-breadcrumb>
<el-menu :default-active="activeIndex" mode="horizontal" @select="handleSelect">
  <el-menu-item index="repo">镜像仓库管理</el-menu-item>
</el-menu>
<keep-alive>
  <repo v-if="activeIndex==='repo'" />
</keep-alive>
</template>

<script setup>
import { ArrowRight,Search,Refresh,EditPen,Delete,Plus,CopyDocument } from '@element-plus/icons-vue'
import { onBeforeMount, ref } from 'vue'
import repo from './repo.vue'
/* 变量定义 */
const activeIndex = ref("repo")
/* methods */
const handleSelect = async (key) => {
  activeIndex.value = key
}
</script>
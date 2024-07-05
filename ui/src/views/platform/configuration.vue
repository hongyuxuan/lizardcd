<template>
<el-breadcrumb :separator-icon="ArrowRight">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>配置</el-breadcrumb-item>
</el-breadcrumb>
<el-menu :default-active="activeIndex" mode="horizontal" @select="handleSelect">
  <el-menu-item index="settings">设置</el-menu-item>
  <el-menu-item index="tenant" v-if="role==='admin'">租户管理</el-menu-item>
  <el-menu-item index="user" v-if="role==='admin'">用户管理</el-menu-item>
  <el-menu-item index="repo">镜像仓库管理</el-menu-item>
</el-menu>
<keep-alive>
  <repo v-if="activeIndex==='repo'" />
  <settings v-else-if="activeIndex==='settings'" />
  <tenant v-else-if="activeIndex==='tenant'" />
  <user v-else-if="activeIndex==='user'" />
</keep-alive>
</template>

<script setup>
import { ArrowRight } from '@element-plus/icons-vue'
import { ref, computed } from 'vue'
import { useStore } from 'vuex'
import repo from './repo.vue'
import settings from './settings.vue'
import user from './user.vue'
import tenant from './teanant.vue'
/* 变量定义 */
const store = useStore()
const role = computed(() => {
  return store.state.role
})
const activeIndex = ref("settings")
/* methods */
const handleSelect = async (key) => {
  activeIndex.value = key
}
</script>
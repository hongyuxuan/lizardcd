<template>
  <el-container>
    <el-aside width="240px">
      <Sidebar />
    </el-aside>
    <el-container>
      <el-header height="56px">
        <HeadBar />
      </el-header>
      <el-main>
        <router-view :key="viewKey" />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import _ from 'lodash'
import { ref, computed, provide, onBeforeMount, watch, nextTick } from 'vue'
import { useRouter, useRoute, onBeforeRouteUpdate } from 'vue-router'
import Sidebar from './components/sidebar.vue'
import HeadBar from './components/header.vue'
/* 变量定义 */
const route = useRoute()
const router = useRouter()
const viewKey = computed(() => {
  return router.currentRoute.value.fullPath
})
/* 生命周期函数 */
onBeforeMount(async () => {
})
/* methods */
</script>

<style>
body {
  margin: 0;
  background-color: #f0f0f0;
}
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  /* text-align: center; */
  color: #2c3e50;
}

.el-header {
  line-height: 50px;
  padding-left: 15px;
}

.el-footer {
  background-color: #fff;
  color: var(--el-text-color-primary);
  text-align: center;
  line-height: 60px;
}

.el-aside {
  color: var(--el-text-color-primary);
  text-align: center;
}

.el-main {
  padding: 15px;
  color: var(--el-text-color-primary);
  min-height: calc(100vh - 56px);
}

.el-container:nth-child(5) .el-aside,
.el-container:nth-child(6) .el-aside {
  line-height: 260px;
}

.el-container:nth-child(7) .el-aside {
  line-height: 320px;
}
</style>

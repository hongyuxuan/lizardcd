<template>
  <el-container>
    <el-aside width="240px" style="z-index:6">
      <el-container>
        <el-main style="padding:0 0 60px">
          <Sidebar />
        </el-main>
        <el-footer style="position: fixed;bottom:0;">
          当前版本: {{ version }}
        </el-footer>
      </el-container>
    </el-aside>
    <el-container>
      <el-header height="56px">
        <HeadBar />
      </el-header>
      <el-main >
        <router-view :key="viewKey" />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import _ from 'lodash'
import { ref, computed, onBeforeMount, watch, nextTick, provide } from 'vue'
import { useRouter } from 'vue-router'
import { useStore } from 'vuex'
import Sidebar from './components/sidebar.vue'
import HeadBar from './components/header.vue'
import { axios } from '/src/assets/util/axios'
import { getSettings } from '/src/assets/util/common'
import * as echarts from 'echarts'
provide('ec', echarts)
/* 变量定义 */
const store = useStore()
const router = useRouter()
const viewKey = computed(() => {
  return router.currentRoute.value.fullPath
})
const version = ref("")
/* watch */
watch(
  () => router.currentRoute.value,
  (newVal, oldVal) => {
    if(newVal.meta?.description) {
      document.title = newVal.meta.description
    }
  }
)
/* 生命周期函数 */
onBeforeMount(async () => {
  await checkLogin()
  getVersion()
  store.state.settings = await getSettings()
})
/* methods */
const checkLogin = async () => {
  let response = await axios.get(`/lizardcd/auth/user/info`, {
    headers: {
      'Authorization': `Bearer ${localStorage.access_token}`
    }
  })
  localStorage.username = response.username
  localStorage.tenant = response.tenant
  delete response.password
  response.profile ||= {}
  store.state.userInfo = response
}
const getVersion = async () => {
  version.value = await axios.get(`/lizardcd/server/version`)
}
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
  text-align: center;
  line-height: 55px;
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

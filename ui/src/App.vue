<template>
  <el-container>
    <el-aside width="240px">
      <el-container>
        <el-main style="padding:0">
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
import { ref, computed, onBeforeMount, watch, nextTick } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useStore } from 'vuex'
import Sidebar from './components/sidebar.vue'
import HeadBar from './components/header.vue'
import { axios } from '/src/assets/util/axios'
/* 变量定义 */
const store = useStore()
const router = useRouter()
const viewKey = computed(() => {
  return router.currentRoute.value.fullPath
})
const version = ref("")
/* 生命周期函数 */
onBeforeMount(async () => {
  checkLogin()
  getVersion()
  getSettings()
})
/* methods */
const checkLogin = async () => {
  let response = await axios.get(`/lizardcd/auth/user/info`, {
    headers: {
      'Authorization': `Bearer ${localStorage.access_token}`
    }
  })
  localStorage.username = response.username
  store.state.username = response.username
  store.state.role = response.role
  localStorage.tenant = response.tenant
}
const getVersion = async () => {
  version.value = await axios.get(`/lizardcd/server/version`)
}
const getSettings = async () => {
  let response = await axios.get(`/lizardcd/db/settings?size=1000&filter=tenant==${localStorage.tenant}`)
  let settings = {}
  for(let x of response.results) {
    if(x.setting_value === 'true' || x.setting_value === 'false')
      x.setting_value = JSON.parse(x.setting_value)
    settings[x.setting_key] = x.setting_value
  }
  store.state.settings = settings
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

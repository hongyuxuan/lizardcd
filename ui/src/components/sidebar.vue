<template>
  <el-scrollbar>
    <div class="sidebar-logo">
      <div v-if="isCollapse" class="logo-mini"><b>A</b>-</div>
      <el-image v-else style="height:45px;margin-top:15px" src="/images/lizardcd-logo.png" />
    </div>
    <div class="user-panel">
      <div class="pull-left image">
        <el-avatar :size="40" :src="avatar" />
      </div>
      <div class="pull-left info sidebar-userinfo">
        <p style="text-align:left">{{username}}</p>
        <a><font-awesome-icon icon="circle" style="color:green" /> Online</a>
      </div>
    </div>
    <div class="user-panel sidenav" style="text-align:left;padding:8px 10px;margin-top:5px;">
      导航
    </div>
    <el-menu
        active-text-color="#ffd04b"
        background-color="#141f29"
        class="el-menu-vertical-demo"
        :default-active="activeIndex"
        text-color="#fff"
        :unique-opened="true"
        :router="true">
        <el-menu-item index="/"><font-awesome-icon icon="home" />首页</el-menu-item>
        <el-menu-item index="/agent"><font-awesome-icon icon="sliders" />连接管理</el-menu-item>
        <el-menu-item index="/application"><font-awesome-icon icon="laptop-code" />应用管理</el-menu-item>
        <el-menu-item index="/task/history"><font-awesome-icon icon="list-check" />任务管理</el-menu-item>
        <el-menu-item index="/template"><font-awesome-icon icon="ghost" />模板管理</el-menu-item>
        <el-sub-menu index="2">
          <template #title>
            <font-awesome-icon icon="layer-group" />
            <span>容器集群资源</span>
          </template>
          <el-sub-menu index="2-2">
            <template #title>
              <span>应用负载</span>
            </template>
            <el-menu-item index="/kubernetes/services">服务</el-menu-item>
            <el-menu-item index="/kubernetes/workload">工作负载</el-menu-item>
            <el-menu-item index="/kubernetes/ingresses">路由与网关</el-menu-item>
            <el-menu-item index="/kubernetes/jobs">任务</el-menu-item>
            <el-menu-item index="/kubernetes/workload/pods">容器组</el-menu-item>
          </el-sub-menu>
          <el-sub-menu index="2-3">
            <template #title>
              <span>存储</span>
            </template>
            <el-menu-item index="/kubernetes/persistentvolumeclaims">存储卷</el-menu-item>
          </el-sub-menu>
          <el-sub-menu index="2-4">
            <template #title>
              <span>配置</span>
            </template>
            <el-menu-item index="/kubernetes/secrets">保密字典</el-menu-item>
            <el-menu-item index="/kubernetes/configmaps">配置字典</el-menu-item>
            <el-menu-item index="/kubernetes/serviceaccounts">服务账户</el-menu-item>
          </el-sub-menu>
        </el-sub-menu>
        <el-menu-item index="/mesh/istio" v-if="enableIstio===true"><span class="iconmoon icon-istio svg-inline--fa text-white"></span>Istio管理</el-menu-item>
        <el-menu-item index="/helm" v-if="enableHelm===true"><span class="iconmoon icon-helm svg-inline--fa"></span>Helm管理</el-menu-item>
        <el-sub-menu index="3" v-if="enableTekton===true">
          <template #title>
            <span class="iconmoon icon-tekton svg-inline--fa" />
            <span>持续集成</span>
          </template>
          <el-menu-item index="/ci/tekton">流水线</el-menu-item>
          <el-menu-item index="/ci/trigger">触发配置</el-menu-item>
        </el-sub-menu>
    </el-menu>
  </el-scrollbar>
</template>

<script setup>
import { onBeforeMount, ref, computed, watch } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'
/* 变量定义 */
const store = useStore()
const router = useRouter()
const avatar = computed(() => {
  return store.state.userInfo.profile?.avatar || "/images/avator.png"
})
const username = computed(() => {
  return store.state.userInfo.username
})
const enableIstio = computed(() => {
  return store.state.settings?.enable_istio||false
})
const enableTekton = computed(() => {
  return store.state.settings?.enable_tekton||false
})
const enableHelm = computed(() => {
  return store.state.settings?.enable_helm||false
})
const isCollapse = ref({})
const activeIndex = ref("/")
/* watch */
watch(
  () => router.currentRoute.value,
  (newVal, oldVal) => {
    activeIndex.value = newVal.path
  }
)
/* 生命周期函数 */
onBeforeMount(async () => {
  isCollapse.value = false
})
</script>
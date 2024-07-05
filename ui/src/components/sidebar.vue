<template>
  <el-scrollbar>
    <div class="sidebar-logo">
      <div v-if="isCollapse" class="logo-mini"><b>A</b>-</div>
      <el-image v-else style="height:45px;margin-top:15px" src="/images/lizardcd-logo.png" />
    </div>
    <div class="user-panel">
      <div class="pull-left image">
        <el-avatar :size="45" :src="avator" />
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
        default-active="2"
        text-color="#fff"
        :unique-opened="true"
        :router="true">
        <el-menu-item index="/agent"><font-awesome-icon icon="sliders" />Agent管理</el-menu-item>
        <el-menu-item index="/application"><font-awesome-icon icon="laptop-code" />应用管理</el-menu-item>
        <el-menu-item index="/task/history"><font-awesome-icon icon="list-check" />任务管理</el-menu-item>
        <el-menu-item index="/template"><font-awesome-icon icon="ghost" />YAML模板</el-menu-item>
        <el-sub-menu index="2">
          <template #title>
            <font-awesome-icon icon="layer-group" />
            <span>工作负载</span>
          </template>
          <el-menu-item index="/workload/deployments">部署</el-menu-item>
          <el-menu-item index="/workload/statefulsets">有状态副本集</el-menu-item>
        </el-sub-menu>
        <el-menu-item index="/mesh/istio" v-if="enableIstio===true"><span class="iconmoon icon-istio svg-inline--fa text-white"></span>Istio管理</el-menu-item>
        <el-menu-item index="/ci/tekton" v-if="enableTekton===true"><span class="iconmoon icon-tekton svg-inline--fa"></span>Tekton管理</el-menu-item>
        <el-menu-item index="/helm" v-if="enableHelm===true"><span class="iconmoon icon-helm svg-inline--fa"></span>Helm管理</el-menu-item>
    </el-menu>
  </el-scrollbar>
</template>

<script setup>
import { onBeforeMount, reactive, toRefs, ref, computed } from 'vue'
import { useStore } from 'vuex'
/* 变量定义 */
const store = useStore()
const avator = ref("/images/avator.png")
const username = computed(() => {
  return store.state.username
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
/* 生命周期函数 */
onBeforeMount(async () => {
  isCollapse.value = false
})
</script>
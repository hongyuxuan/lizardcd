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
        <el-sub-menu index="2">
          <template #title>
            <font-awesome-icon icon="layer-group" />
            <span>工作负载</span>
          </template>
          <el-menu-item index="/workload/deployment">部署</el-menu-item>
          <el-menu-item index="/workload/statefulset">有状态副本集</el-menu-item>
        </el-sub-menu>
    </el-menu>
  </el-scrollbar>
</template>

<script setup>
import { onBeforeMount, reactive, toRefs, ref, computed } from 'vue'
const avator = ref("/images/avator.png")
const username = ref("游客")
const isCollapse = ref({})
/* 生命周期函数 */
onBeforeMount(async () => {
  isCollapse.value = false
})
</script>
<template>
  <el-menu
    mode="horizontal"
    :ellipsis="false"
    @select="handleSelect">
    <el-menu-item index="0">
      <span class="portal-title">Lizardcd Cloud Native Continuous Delivery for Kubernetes</span>
    </el-menu-item>
    <div class="flex-grow" />
    <el-menu-item index="1">
      <div>
        <el-avatar :size="45" :src="avator" style="vertical-align:middle" />
        <span style="margin-left:15px">{{username}}</span>
      </div>
    </el-menu-item>
    <el-menu-item index="2"><font-awesome-icon icon="gears" style="font-size:20px" /></el-menu-item>
  </el-menu>
</template>

<script setup lang="ts">
import { ref } from 'vue'
/* 变量定义 */
const avator = ref("/images/avator.png")
const username = ref("游客")
/* methods */
const handleSelect = (index: string) => {
  switch(index) {
    case '2-1': console.log("平台管理")
  }
}
</script>

<style>
.el-header {
  padding: 0 !important;
}
.sidebar-toggle {
  float: left;
  padding: 0 15px;
  color: inherit;
}
.portal-title {
  float:left;
  padding-left: 10px;
  font-size: 16px;
}
.flex-grow {
  flex-grow: 1;
}
.el-header .el-menu--horizontal .el-sub-menu__icon-arrow {
  display: none;
}
.el-header .el-menu--horizontal>.el-menu-item.is-active,
.el-header .el-menu--horizontal>.el-submenu.is-active>.el-submenu__title {
  border-bottom: 0px solid #ffffff !important;
}
</style>
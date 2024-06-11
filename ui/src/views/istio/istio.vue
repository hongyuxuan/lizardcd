<template>
<el-breadcrumb :separator-icon="ArrowRight">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>Istio管理</el-breadcrumb-item>
</el-breadcrumb>
<el-menu :default-active="activeIndex" mode="horizontal" @select="handleSelect">
  <el-menu-item index="destinationrule">目标规则</el-menu-item>
  <el-menu-item index="virtualservice">虚拟服务</el-menu-item>
</el-menu>
<keep-alive>
  <destinationrule v-if="activeIndex==='destinationrule'" />
  <virtualservice v-else-if="activeIndex==='virtualservice'" />
</keep-alive>
</template>

<script setup>
import { ArrowRight } from '@element-plus/icons-vue'
import { ref } from 'vue'
import destinationrule from './destinationrule.vue'
import virtualservice from './virtualservice.vue'
/* 变量定义 */
const activeIndex = ref("destinationrule")
/* methods */
const handleSelect = async (key) => {
  activeIndex.value = key
}
</script>
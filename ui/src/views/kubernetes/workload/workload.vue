<template>
<el-breadcrumb :separator-icon="ArrowRight">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>应用负载</el-breadcrumb-item>
  <el-breadcrumb-item>工作负载</el-breadcrumb-item>
</el-breadcrumb>
<div class="box box-solid">
  <div class="box-header page-intro">
    <p>
      工作负载 (Workload) 通常是访问服务的实际载体, 也是对节点日志收集、监控等系统应用的实际运行载体，是对一组容器组 (Pod) 的抽象模型。
    </p>
  </div>
  <div class="box-body" style="padding:0">
    <el-tabs v-model="activeName" type="border-card" class="qingcloud-tab">
      <el-tab-pane name="deployments" label="部署">
        <keep-alive>
          <deployments />
        </keep-alive>
      </el-tab-pane>
      <el-tab-pane name="statefulsets" label="有状态副本集">
        <keep-alive>
          <statefulsets />
        </keep-alive>
      </el-tab-pane>
      <el-tab-pane name="daemonsets" label="守护进程集">
        <keep-alive>
          
        </keep-alive>
      </el-tab-pane>
    </el-tabs>
  </div>
</div>
</template>

<script setup>
import { ArrowRight } from '@element-plus/icons-vue'
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import deployments from './deployments.vue'
import statefulsets from './statefulsets.vue'
/* 变量定义 */
const route = useRoute()
const activeName = ref("deployments")
/* 生命周期函数 */
onMounted(() => {
  if(route.query.tab) {
    activeName.value = route.query.tab
  }
})
</script>
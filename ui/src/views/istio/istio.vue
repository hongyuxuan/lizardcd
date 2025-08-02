<template>
<el-breadcrumb :separator-icon="ArrowRight">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>Istio管理</el-breadcrumb-item>
</el-breadcrumb>
<div class="box box-solid">
  <div class="box-header page-intro">
    <p>
      Istio 是一个开源的服务网格平台，可以在 Kubernetes 和传统工作负载上扩展功能，提供流量管理、可观测性、安全和策略等功能。参见 
      <el-link href="http://istio.io" underline="never" type="primary" target="_blank">Istio</el-link><br>
      Istio 管理页面为您提供 Istio 的常见资源类型的增删改查等管理操作。
    </p>
  </div>
  <div class="box-body" style="padding:0">
    <el-tabs v-model="activeName" type="border-card" class="qingcloud-tab">
      <el-tab-pane name="1" label="目标规则">
        <keep-alive>
          <destinationrule />
        </keep-alive>
      </el-tab-pane>
      <el-tab-pane name="2" label="虚拟服务">
        <keep-alive>
          <virtualservice />
        </keep-alive>
      </el-tab-pane>
      <el-tab-pane name="3" label="网关">
        <keep-alive>
          <gateway />
        </keep-alive>
      </el-tab-pane>
    </el-tabs>
  </div>
</div>
</template>

<script setup>
import { ArrowRight } from '@element-plus/icons-vue'
import { ref } from 'vue'
import destinationrule from './destinationrule.vue'
import virtualservice from './virtualservice.vue'
import gateway from './gateway.vue'
/* 变量定义 */
const activeName = ref("1")
/* methods */
const handleSelect = async (key) => {
  activeIndex.value = key
}
</script>
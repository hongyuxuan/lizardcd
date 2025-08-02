<template>
<el-breadcrumb :separator-icon="ArrowRight">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>Helm管理</el-breadcrumb-item>
</el-breadcrumb>
<div class="box box-solid">
  <div class="box-header page-intro">
    <p>
      Helm 是 Kubernetes 环境下的包管理工具，可以帮助您在 Kubernetes 上快速部署和管理应用程序。参见 
      <el-link href="https://helm.sh" underline="never" type="primary" target="_blank">Helm</el-link><br>
      Helm 管理页面为您提供 Helm 仓库的增删改查、Charts 包搜索、查看、安装；已安装 Charts 包的更新、重装、卸载等管理操作。
    </p>
  </div>
  <div class="box-body" style="padding:0">
    <el-tabs v-model="activeName" type="border-card" class="qingcloud-tab">
      <el-tab-pane name="repo" label="Helm仓库">
        <keep-alive>
          <repo />
        </keep-alive>
      </el-tab-pane>
      <el-tab-pane name="release" label="Helm包发布">
          <keep-alive>
          <release />
        </keep-alive>
      </el-tab-pane>
    </el-tabs>
  </div>
</div>
</template>

<script setup>
import { ArrowRight } from '@element-plus/icons-vue'
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import repo from './repo.vue'
import release from './release.vue'
/* 变量定义 */
const route = useRoute()
const activeName = ref("repo")
/* 生命周期函数 */
onMounted(() => {
  if(route.query.tab !== undefined)
  activeName.value = route.query.tab
})
/* methods */
const handleSelect = async (key) => {
  activeName.value = key
}
</script>
<template>
<el-breadcrumb :separator-icon="ArrowRight">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>配置</el-breadcrumb-item>
</el-breadcrumb>
<div class="box box-solid">
  <div class="box-header page-intro">
    <p v-html="info" />
  </div>
  <div class="box-body" style="padding:0">
    <el-tabs v-model="activeName" type="border-card" @tab-change="changeTab" class="qingcloud-tab">
      <el-tab-pane name="1" label="设置">
        <keep-alive>
          <settings />
        </keep-alive>
      </el-tab-pane>
      <el-tab-pane name="2" label="租户管理" v-if="role==='admin'">
          <keep-alive>
          <tenant />
        </keep-alive>
      </el-tab-pane>
      <el-tab-pane name="3" label="用户管理" v-if="role==='admin'">
          <keep-alive>
          <user />
        </keep-alive>
      </el-tab-pane>
      <el-tab-pane name="4" label="登录认证" v-if="role==='admin'">
          <keep-alive>
          <authorized />
        </keep-alive>
      </el-tab-pane>
      <el-tab-pane name="5" label="镜像仓库管理">
          <keep-alive>
          <repo />
        </keep-alive>
      </el-tab-pane>
      <el-tab-pane name="6" label="代码仓库管理">
          <keep-alive>
          <git />
        </keep-alive>
      </el-tab-pane>
      <el-tab-pane name="7" label="令牌管理" v-if="role==='admin'">
          <keep-alive>
          <token />
        </keep-alive>
      </el-tab-pane>
    </el-tabs>
  </div>
</div>
</template>

<script setup>
import { ArrowRight } from '@element-plus/icons-vue'
import { ref, computed } from 'vue'
import { useStore } from 'vuex'
import repo from './repo.vue'
import settings from './settings.vue'
import user from './user.vue'
import tenant from './teanant.vue'
import authorized from './authorized.vue'
import git from './git.vue'
import token from './token.vue'
/* 变量定义 */
const store = useStore()
const role = computed(() => {
  return store.state.userInfo.role
})
const activeName = ref("1")
const info = ref("配置页面为您提供 Lizardcd 平台的通用设置，包括参数设置、镜像仓库、用户设置、租户设置等。<br>部分功能仅管理员有权限设置。如有问题，请联系您的管理员。")
/* methods */
const changeTab = (tabName) => {
  if(tabName === "7") {
    info.value = "您可以在此页面创建访问 Lizardcd 的 JWT 令牌，并为令牌设置过期时间。<br>该页面功能只有管理员可操作。"
  } else {
    info.value = "配置页面为您提供 Lizardcd 平台的通用设置，包括参数设置、镜像仓库、用户设置、租户设置等。<br>部分功能仅管理员有权限设置。如有问题，请联系您的管理员。"
  }
}
</script>
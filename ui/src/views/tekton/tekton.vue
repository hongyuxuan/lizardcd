<template>
<el-breadcrumb :separator-icon="ArrowRight">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>持续集成</el-breadcrumb-item>
  <el-breadcrumb-item>Tekton管理</el-breadcrumb-item>
</el-breadcrumb>
<div class="box box-solid">
  <div class="box-header page-intro">
    <p>
      Tekton 是一个云原生的持续集成与持续构建平台，它在 Kubernetes 上以 CRD 的方式定义了构建所需的各种类型资源。参见 
      <el-link href="http://tekton.dev" :underline="false" type="primary" target="_blank">Tekton</el-link><br>
      Tekton 管理页面为您提供 Tekton 的常见资源类型的增删改查等管理操作。
    </p>
  </div>
  <div class="box-body" style="padding:0">
    <el-tabs v-model="activeName" type="border-card" class="qingcloud-tab">
      <el-tab-pane name="1" label="Task">
        <keep-alive>
          <task />
        </keep-alive>
      </el-tab-pane>
      <el-tab-pane name="2" label="Pipeline">
        <keep-alive>
          <pipeline />
        </keep-alive>
      </el-tab-pane>
      <el-tab-pane name="3" label="PipelineRun">
        <keep-alive>
          <pipelinerun />
        </keep-alive>
      </el-tab-pane>
      <el-tab-pane name="4" label="TriggerBinding">
        <keep-alive>
          <triggerbinding />
        </keep-alive>
      </el-tab-pane>
      <el-tab-pane name="5" label="TriggerTemplate">
        <keep-alive>
          <triggertemplate />
        </keep-alive>
      </el-tab-pane>
      <el-tab-pane name="6" label="EventListener">
        <keep-alive>
          <eventlistener />
        </keep-alive>
      </el-tab-pane>
    </el-tabs>
  </div>
</div>
</template>

<script setup>
import { ArrowRight } from '@element-plus/icons-vue'
import { ref } from 'vue'
import task from './task.vue'
import pipeline from './pipeline.vue'
import pipelinerun from './pipelinerun.vue'
import triggerbinding from './triggerbinding.vue'
import triggertemplate from './triggertemplate.vue'
import eventlistener from './eventlistener.vue'
/* 变量定义 */
const activeName = ref("1")
</script>
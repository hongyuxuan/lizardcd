<template>
<div class="box box-item">
  <div class="box-body" style="padding-top:20px;padding-bottom:0">
    <el-collapse v-model="activeNames">
      <el-collapse-item name="1">
        <template #title><h4><b>开启对接Istio</b></h4></template>
        <el-row>
          <el-col :span="18">Istio是一个云原生的服务网格解决方案，提供微服务调用、流量治理、网关、服务跟踪监控等组件功能。<br>如需开启对接Istio，lizardcd-agent所在（或对接）集群需要已经部署好Istio，具备<code>istio.io</code> Group下的相关apiVersion，同时agent启动时使用的<code>kubeconfig</code>或<code>serviceaccount</code>必须具备对<code>istio.io</code> Group下资源的相关操作权限。</el-col>
          <el-col :span="6" ><el-switch v-model="settings.enable_istio.setting_value" size="large" @change="setEnable('enable_istio')" /> </el-col>
        </el-row>
      </el-collapse-item>
      <el-collapse-item name="2">
        <template #title><h4><b>开启对接Tekton</b></h4></template>
        <el-row>
          <el-col :span="18">Tekton是一个云原生的持续集成工具。<br>如需开启对接Tekton，lizardcd-agent所在（或对接）集群需要已经部署好Tekton，具备<code>tekton.dev</code> Group下的相关apiVersion，同时agent启动时使用的<code>kubeconfig</code>或<code>serviceaccount</code>必须具备对<code>tekton.dev</code> Group下资源的相关操作权限。</el-col>
          <el-col :span="6" ><el-switch v-model="settings.enable_tekton.setting_value" size="large" @change="setEnable('enable_tekton')" /> </el-col>
        </el-row>
      </el-collapse-item>
      <el-collapse-item name="3">
        <template #title><h4><b>开启Helm包管理功能</b></h4></template>
        <el-row>
          <el-col :span="18">Helm是Kubernetes环境下的包管理工具。<br>开启Helm包管理功能仅需agent启动时使用的<code>kubeconfig</code>或<code>serviceaccount</code>具备对相关namespace的操作权限。</el-col>
          <el-col :span="6" ><el-switch v-model="settings.enable_helm.setting_value" size="large" @change="setEnable('enable_helm')" /> </el-col>
        </el-row>
        <el-divider />
        <el-row>
          <el-col :span="18">Helm包安装升级开启等待（--wait）。<br>如不开启等待，Helm将values参数及内置参数按照模板渲染成Kubernetes资源并提交到Apiserver，当Apiserver接收到资源请求后，Helm即认为安装升级成功（此时资源可能并没有真正Ready）。如开启等待，则Helm会等待所有的Pods, PVCs, Services, 和minimum number of Pods of a Deployment, statefulSet, or ReplicaSet处于Ready状态，才认为本次安装升级成功。</el-col>
          <el-col :span="6" ><el-switch v-model="settings.helm_wait.setting_value" size="large" @change="setEnable('helm_wait')" /> </el-col>
        </el-row>
        <el-divider />
        <el-row>
          <el-col :span="18">Helm包安装升级等待超时时间（--timeout）。<br>在Helm开启等待条件下，Helm等待资源处于Ready状态的超时时间，单位为<code>秒</code></el-col>
          <el-col :span="6" style="text-align:right;"><el-input-number v-model="settings.helm_timeout.setting_value" size="large" @change="setValue('helm_timeout')" /> </el-col>
        </el-row>
      </el-collapse-item>
    </el-collapse>
  </div>
</div>
</template>
<script setup>
import axios from 'axios';
import { onBeforeMount, ref, reactive } from 'vue'
/* 变量定义 */
const activeNames = ref(["1","2","3"])
const settings = ref({
  enable_istio: {},
  enable_tekton: {},
  enable_helm: {},
  helm_wait: {},
  helm_timeout: {},
})
/* 生命周期函数 */
onBeforeMount(async () => {
  getSettings()
})
/* methods */
const getSettings = async () => {
  let response = await axios.get(`/lizardcd/db/settings?size=1000&filter=tenant==${localStorage.tenant}`)
  for(let x of response.results) {
    if(x.setting_value === 'true' || x.setting_value === 'false')
      x.setting_value = JSON.parse(x.setting_value)
    if(x.setting_key == 'helm_timeout')
      x.setting_value = parseInt(x.setting_value)
    settings.value[x.setting_key] = x
  }
}
const setEnable = async (setting_key) => {
  await axios.put(`/lizardcd/db/settings/${settings.value[setting_key].id}`, {
    body: {
      setting_value: settings.value[setting_key].setting_value.toString()
    }
  })
}
const setValue = async (setting_key) => {
  await axios.put(`/lizardcd/db/settings/${settings.value[setting_key].id}`, {
    body: {
      setting_value: settings.value[setting_key].setting_value
    }
  })
}
</script>
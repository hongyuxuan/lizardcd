<template>
  <div class="box box-item">
    <div class="box-body" style="padding-top:20px;padding-bottom:0">
      <el-collapse v-model="activeNames">
        <el-collapse-item name="1" v-if="role==='admin'">
          <template #title><h4><b>设置服务端URL</b></h4></template>
          <el-row>
            <el-col :span="15">您可能将Lizardcd部署于您环境的网关之内、或者反向代理服务器之后，<br>因部分对接第三方平台功能需要从第三方平台回调Lizardcd的接口，因此您需要在此配置Lizardcd的接口URL，也即服务端URL。</el-col>
            <el-col :span="9" ><el-input v-model="settings.server_url.setting_value" size="large" clearable style="width:400px" placeholder="按回车提交" @change="setValue('server_url')" /></el-col>
          </el-row>
        </el-collapse-item>
        <el-collapse-item name="2">
          <template #title><h4><b>开启对接Istio</b></h4></template>
          <el-row>
            <el-col :span="18">Istio是一个云原生的服务网格解决方案，提供微服务调用、流量治理、网关、服务跟踪监控等组件功能。<br>如需开启对接Istio，lizardcd-agent所在（或对接）集群需要已经部署好Istio，具备<code>istio.io</code> Group下的相关apiVersion，同时agent启动时使用的<code>kubeconfig</code>或<code>serviceaccount</code>必须具备对<code>istio.io</code> Group下资源的相关操作权限。</el-col>
            <el-col :span="6" ><el-switch v-model="settings.enable_istio.setting_value" size="large" @change="setEnable('enable_istio')" /> </el-col>
          </el-row>
        </el-collapse-item>
        <el-collapse-item name="3">
          <template #title><h4><b>开启持续集成功能</b></h4></template>
          <el-row>
            <el-col :span="18">Tekton是一个云原生的持续集成工具。Lizardcd支持对接Tekton，并提供对目标Tekton资源进行管理。同时Lizardcd的Serverless功能必须开启持续集成。<br>如需开启对接Tekton，lizardcd-agent所在（或对接）集群需要已经部署好Tekton，具备<code>tekton.dev</code> Group下的相关apiVersion，同时agent启动时使用的<code>kubeconfig</code>或<code>serviceaccount</code>必须具备对<code>tekton.dev</code> Group下资源的相关操作权限。<br>开启持续集成后，Lizardcd还提供对代码仓库的工程进行精准触发的功能。</el-col>
            <el-col :span="6" ><el-switch v-model="settings.enable_tekton.setting_value" size="large" @change="setEnable('enable_tekton')" /> </el-col>
          </el-row>
          <el-divider v-if="role==='admin'" />
          <el-row v-if="role==='admin'">
            <el-col :span="12">设置默认Tekton用于应用自动构建</el-col>
            <el-col :span="12">
              <el-button-group>
                <el-select v-model="cluster" placeholder="请选择集群" clearable filterable @change="setTekton()" style="width:200px;margin-right:5px" size="large">
                  <el-option v-for="(v,k) in clusterList" :key="k" :label="k" :value="k" />
                </el-select>
                <el-select v-model="namespace" placeholder="请选择命名空间" clearable filterable @change="setTekton()" style="width:200px;margin-right:5px" size="large">
                  <el-option v-for="(item) in clusterList[cluster]" :key="item" :label="item" :value="item" />
                </el-select>
                <el-input v-model="defaultTekton" size="large" clearable style="width:200px" />
              </el-button-group>
            </el-col>
          </el-row>
          <el-divider v-if="role==='admin'" />
          <el-row v-if="role==='admin'">
            <el-col :span="12">设置Tekton运行时数据获取方式</el-col>
            <el-col :span="12">
              <el-radio-group v-model="settings.tekton_source.setting_value" @change="setValue('tekton_source')" class="pull-right">
                <el-radio value="crd">From Kubernetes Tekton CRD</el-radio>
                <el-radio value="cloudevent">From CloudEvents Database</el-radio>
              </el-radio-group>
            </el-col>
          </el-row>
          <el-divider />
          <el-row>
            <el-col :span="12">设置Tekton流水线详情页展示方式</el-col>
            <el-col :span="12">
              <el-radio-group v-model="settings.pipelinerun_result.setting_value" @change="setValue('pipelinerun_result')" class="pull-right">
                <el-radio value="flow">流程图</el-radio>
                <el-radio value="collapse">折叠面板</el-radio>
              </el-radio-group>
            </el-col>
          </el-row>
        </el-collapse-item>
        <el-collapse-item name="4">
          <template #title><h4><b>开启对接DolphinScheduler</b></h4></template>
          <el-row>
            <el-col :span="18">DolphinScheduler是一个apache开源的分布式任务调度平台。Lizardcd支持对接DolphinScheduler实现更高层次的任务调度功能。<br>应用的GitOps如果选择使用DolphinScheduler进行资源定时同步，则必须在此设置默认DolphinScheduler。<br>同时，您还必须在此设置DolphinScheduler回调Lizardcd时所用的JwtToken，否则接口将返回401错误。</el-col>
          </el-row>
          <el-divider />
          <el-row>
            <el-col :span="6">设置DolphinScheduler</el-col>
            <el-col :span="18">
              <el-button circle :icon="Check" size="large" style="float:right" @click="setDS" />
              <el-input v-model="dolphinscheduler.project" size="large" clearable style="width:180px;margin-right:5px;float:right" placeholder="请输入Project Code" />
              <el-input v-model="dolphinscheduler.token" size="large" clearable style="width:300px;margin-right:5px;float:right" placeholder="请输入dolphinscheduler Token" />
              <el-input v-model="dolphinscheduler.base_url" size="large" clearable style="width:400px;margin-right:5px;float:right" placeholder="请输入dolphinscheduler API地址" />
            </el-col>
          </el-row>
          <el-divider />
          <el-row>
            <el-col :span="6">设置回调Lizardcd的JwtToken</el-col>
              <el-col :span="18">
                <el-button circle :icon="Check" size="large" style="float:right" @click="setDS" />
                <el-input v-model="dolphinscheduler.lizardcd_jwt_token" size="large" clearable style="width:890px;margin-right:5px;float:right" placeholder="请输入Lizardcd的JwtToken" type="textarea" :autosize="{ minRows: 2}" />
            </el-col>
          </el-row>
        </el-collapse-item>
        <el-collapse-item name="5">
          <template #title><h4><b>开启Helm包管理功能</b></h4></template>
          <el-row>
            <el-col :span="18">Helm是Kubernetes环境下的包管理工具。<br>开启Helm包管理功能仅需 server（直连K8S）或 agent（间接连 K8S）启动时使用的<code>kubeconfig</code>或<code>serviceaccount</code>具备对相关namespace的操作权限。</el-col>
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
            <el-col :span="6"><el-input-number v-model="settings.helm_timeout.setting_value" size="large" @change="setValue('helm_timeout')" /> </el-col>
          </el-row>
        </el-collapse-item>
      </el-collapse>
    </div>
  </div>
  </template>
  <script setup>
  import { onBeforeMount, ref, computed } from 'vue'
  import { Check } from '@element-plus/icons-vue'
  import { useStore } from 'vuex'
  import axios from 'axios'
  /* 变量定义 */
  const store = useStore()
  const role = computed(() => {
    return store.state.userInfo.role
  })
  const activeNames = ref(["1","2","3","4","5"])
  const settings = ref({
    enable_istio: {},
    enable_tekton: {},
    enable_helm: {},
    helm_wait: {},
    helm_timeout: {},
    default_tekton: {},
    server_url: {},
    tekton_source: {},
    pipelinerun_result: {},
  })
  const cluster = ref("")
  const clusterList = ref({})
  const namespace = ref("")
  const defaultTekton = ref("")
  const dolphinscheduler = ref({})
  /* 生命周期函数 */
  onBeforeMount(async () => {
    getSettings()
    getClusterList()
  })
  /* methods */
  const getSettings = async () => {
    let tenants = localStorage.tenant?.split(",") || []
    let response = await axios.get(`/lizardcd/db/settings?size=1000&filter=tenant==${tenants[0]}`)
    for(let x of response.results) {
      if(x.setting_value === 'true' || x.setting_value === 'false')
        x.setting_value = JSON.parse(x.setting_value)
      if(x.setting_key == 'helm_timeout')
        x.setting_value = parseInt(x.setting_value)
      if(x.setting_key == 'default_tekton') {
        x.setting_value = JSON.parse(x.setting_value)
        defaultTekton.value = x.setting_value.endpoint
        cluster.value = x.setting_value.cluster
        namespace.value = x.setting_value.namespace
      }
      if(x.setting_key == 'dolphinscheduler') {
        dolphinscheduler.value = x.setting_value = JSON.parse(x.setting_value)
      }
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
  const setTekton = async () => {
    defaultTekton.value = ""
    let setting_value = {
      cluster: cluster.value,
      namespace: namespace.value,
    }
    if(cluster.value&&namespace.value) {
      setting_value.endpoint = await axios.get(`/lizardcd/tekton/cluster/${cluster.value}/namespace/${namespace.value}/endpoint`)
      defaultTekton.value = setting_value.endpoint
    }
    settings.value.default_tekton.setting_value = JSON.stringify(setting_value)
    await setValue("default_tekton")
  }
  const getClusterList = async () => {
    clusterList.value = await axios.get(`/lizardcd/server/clusters`)
  }
  const setDS = async () => {
    settings.value.dolphinscheduler.setting_value = JSON.stringify(dolphinscheduler.value)
    await setValue("dolphinscheduler")
  }
  </script>
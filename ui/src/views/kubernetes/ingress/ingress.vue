<template>
<el-breadcrumb :separator-icon="ArrowRight" style="min-width:1366px">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>应用负载</el-breadcrumb-item>
  <el-breadcrumb-item :to="{ path: '/kubernetes/ingresses', query: {cluster: route.query.cluster, namespace: route.query.namespace} }">路由</el-breadcrumb-item>
  <el-breadcrumb-item>{{ ingressInfo.metadata?.name }}</el-breadcrumb-item>
</el-breadcrumb>
<el-row :gutter="15" style="min-width:1366px">
  <el-col :span="6">
    <el-card>
      <template #header>
        <div class="card-header">
          <el-text truncated size="large" style="width:250px;font-size:18px" class="text-black"><b>{{ ingressInfo.metadata?.name }}</b></el-text>
          <el-dropdown @command="handleCommand" class="pull-right" style="top:2px;min-width:75px;">
            <el-link underline="never">
              更多操作
              <el-icon class="el-icon--right">
                <arrow-down />
              </el-icon>
            </el-link>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="yaml">编辑YAML</el-dropdown-item>
                <el-dropdown-item command="delete">删除</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </template>
      <el-descriptions :column="1" border class="no-color" :label-width="120">
        <el-descriptions-item label="集群">{{ route.query.cluster }}</el-descriptions-item>
        <el-descriptions-item label="命名空间">{{ route.query.namespace }}</el-descriptions-item>
        <el-descriptions-item label="网关地址">
          <div v-if="ingressInfo.status?.loadBalancer&&Object.values(ingressInfo.status?.loadBalancer).length>0">
            {{ Object.values(ingressInfo.status?.loadBalancer)[0][0].ip }}
          </div>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ ingressInfo.metadata?.creationTimestamp }}</el-descriptions-item>
      </el-descriptions>
    </el-card>
  </el-col>
  <el-col :span="18">
    <el-menu :default-active="activeIndex" mode="horizontal" @select="handleSelect">
      <el-menu-item index="status">资源状态</el-menu-item>
      <el-menu-item index="labels">标签</el-menu-item>
      <el-menu-item index="annotations">注解</el-menu-item>
      <el-menu-item index="events">事件</el-menu-item>
    </el-menu>
    <el-card v-show="activeIndex==='status'" style="margin-top:15px">
      <template #header>
        <div class="card-header">
          <span class="card-header-text">规则</span>
        </div>
      </template>
      <div class="box box-item" v-for="(item,i) in ingressInfo.spec?.rules||[]" :key="i">
        <div class="box-header" style="line-height:25px;font-size:15px">
          <b>{{ item.host }}</b>
          <div class="box-tools">{{ item.http ? 'http' : 'https' }}</div>
        </div>
        <div class="box-body" style="padding:5px 0 0 0">
          <el-table :data="item.http?.paths||item.https.paths" :show-header="false">
            <el-table-column>
              <template #default="scope">
                <b>{{ scope.row.path }}</b>
                <div class="text-gray cell-comment" style="font-size:14px">Path</div>
              </template>
            </el-table-column>
            <el-table-column>
              <template #default="scope">
                <el-link underline="never" :href="`/kubernetes/services/${scope.row.backend[Object.keys(scope.row.backend)[0]].name}?cluster=${route.query.cluster}&namespace=${route.query.namespace}`"><b>{{ scope.row.backend[Object.keys(scope.row.backend)[0]].name }}</b></el-link>
                <div class="text-gray cell-comment" style="font-size:14px">Service</div>
              </template>
            </el-table-column>
            <el-table-column>
              <template #default="scope">
                <b>{{ scope.row.backend[Object.keys(scope.row.backend)[0]].port.number }}</b>
                <div class="text-gray cell-comment" style="font-size:14px">Port</div>
              </template>
            </el-table-column>
            <el-table-column align="right" width="100">
              <template #default="scope">
                <el-button round @click="goto(item.http ? 'http' : 'https', item.host, scope.row)">点击访问</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </div>
    </el-card>
    <div class="box box-item" v-show="activeIndex==='labels'">
      <div class="box-body" style="padding-top:20px">
        <el-descriptions :column="1" border class="no-color" :label-width="150">
          <el-descriptions-item v-for="(v,k,i) in ingressInfo.metadata?.labels" :index="i" :label="k">{{ v }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </div>
    <div class="box box-item" v-show="activeIndex==='annotations'">
      <div class="box-body" style="padding-top:20px">
        <el-descriptions :column="1" border class="no-color" :label-width="180">
          <el-descriptions-item v-for="(v,k,i) in ingressInfo.metadata?.annotations" :index="i" :label="k">{{ v }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </div>
    <eventList v-show="activeIndex==='events'" ref="refEvents" resourceType="Ingress" :resourceName="route.params.ingress_name" />
  </el-col>
</el-row>
<el-drawer v-model="show.yaml" direction="rtl" size="700px">
  <template #header>
    <h4>编辑YAML</h4>
  </template>
  <template #default>
    <div style="position:relative">
      <el-switch v-model="wrapLine" active-text="自动换行" inactive-text="不自动换行" inline-prompt class="el-switch-hover-right" size="large" />
      <v-ace-editor
        v-model:value="yamlContent"
        lang="yaml"
        theme="chrome"
        style="width:100%"
        :options="{
          wrap: wrapLine,
          enableBasicAutocompletion: true,
          enableSnippets: true,
          enableLiveAutocompletion: true,
          tabSize: 2,
          showPrintMargin: false,
          fontSize: 14,
          minLines: 10,
          maxLines: 5000,
        }" />
    </div>
  </template>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show.yaml=false">取消</el-button>
      <el-button type="primary" @click="submitYaml">提交</el-button>
    </div>
  </template>
</el-drawer>
</template>

<script setup>
import { ArrowRight } from '@element-plus/icons-vue'
import { onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRoute, useRouter } from 'vue-router'
import eventList from '../eventList.vue'
import { axios } from '/src/assets/util/axios'
import moment from 'moment'
/* 引入v-ace-editor */
import { VAceEditor } from 'vue3-ace-editor'
import 'ace-builds/src-noconflict/mode-yaml'
import 'ace-builds/src-noconflict/theme-chrome'
import 'ace-builds/src-noconflict/ext-language_tools'
/* 变量定义 */
const route = useRoute()
const router = useRouter()
const ingressInfo = ref({})
const activeIndex = ref("status")
const show = ref({
  yaml: false
})
const yamlContent = ref("")
const wrapLine = ref(true)
const refEvents = ref(null)
/* 生命周期函数 */
onMounted(async () => {
  await getIngress()
  refEvents.value.getEvents()
})
/* methods */
const getIngress = async () => {
  ingressInfo.value = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/ingresses/${route.params.ingress_name}`)
  if(!ingressInfo.value.metadata) {
    ElMessage.error({message: `未找到路由: ${route.params.ingress_name}`})
    return
  }
  ingressInfo.value.metadata.creationTimestamp = moment(ingressInfo.value.metadata.creationTimestamp).format('YYYY-MM-DD HH:mm:ss')
}
const handleCommand = async (command) => {
  switch(command) {
    case "yaml": {
      yamlContent.value = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/ingresses/${route.params.ingress_name}/yaml`)
      show.value.yaml = true
      break
    }
    case "delete": {
      await ElMessageBox.confirm('确定删除？','警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(async () => {
        await axios.delete(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/ingresses/${route.params.ingress_name}`)
        ElMessage.success({message: '删除成功'})
        router.push({
          path: '/kubernetes/workload/ingresses',
          query: {
            cluster: route.query.cluster,
            namespace: route.query.namespace
          }
        })
      }).catch(() =>{})
      break
    }
  }
}
const handleSelect = async (key) => {
  activeIndex.value = key
}
const submitYaml = async () => {
  await axios.patch(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/apply/yaml`, yamlContent.value, {
    headers: {
      'Content-Type': 'text/plain'
    }
  })
  show.value.yaml = false
  await getIngress()
  listWorkloads()
  refPods.value.getPods()
}
const goto = async (protocol, host, row) => {
  window.open(`${protocol}://${host}${row.path}`)
}
</script>
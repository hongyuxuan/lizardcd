<template>
<el-breadcrumb :separator-icon="ArrowRight">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>配置</el-breadcrumb-item>
  <el-breadcrumb-item>保密字典</el-breadcrumb-item>
</el-breadcrumb>
<div class="box box-solid">
  <div class="box-header page-intro">
    <p>
      保密字典（Secret）常用于存储工作负载所需的配置信息，许多应用程序会从配置文件、命令行参数或环境变量中读取配置信息。
    </p>
  </div>
  <div class="box-body" style="padding:0">
    <el-tabs v-model="activeName" type="border-card" class="qingcloud-tab">
      <el-tab-pane name="secrets" label="保密字典">
        <div class="box box-item">
          <div class="box-body" style="padding-top:20px;padding-bottom:0">
            <el-row>
              <el-col :span="12">
                <el-button-group>
                  <el-button :icon="Refresh" size="large" @click="getList(current)" />
                  <el-select v-model="cluster" placeholder="请选择集群" clearable filterable @change="namespace=''" style="width:200px;" size="large">
                    <el-option v-for="(v,k) in clusterList" :key="k" :label="k" :value="k" />
                  </el-select>
                  <el-select v-model="namespace" placeholder="请选择命名空间" clearable filterable @change="getList(true)" style="width:200px;" size="large">
                    <el-option v-for="(item) in clusterList[cluster]" :key="item" :label="item" :value="item" />
                  </el-select>
                  <el-input v-model="searchKey" placeholder="输入名称进行搜索" size="large" :prefix-icon="Search" @change="getPage(1)" clearable style="width:200px;" />
                </el-button-group>
              </el-col>
              <el-col :span="12">
                <el-dropdown @command="handleMore" class="pull-right">
                  <el-button size="large">更多操作<el-icon class="el-icon--right"><arrow-down /></el-icon></el-button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item :command="{action:'deleteBatch'}">批量删除</el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
                <el-button class="pull-right" size="large" type="primary" @click="show.new=true;">+ 新建保密字典</el-button>
              </el-col>
            </el-row>
            <el-table 
              :data="list" 
              v-loading="loading"
              element-loading-text="奋力加载中..."
              class="line-height40" 
              @selection-change="select"
              style="width:100%;margin-top:10px;min-height:150px">
              <el-table-column type="selection" width="45" />
              <el-table-column label="" width="45">
                <font-awesome-icon icon="key" style="font-size:25px;vertical-align:middle;" />
              </el-table-column>
              <el-table-column prop="name" label="名称" min-width="200">
                <template #default="scope">
                  <el-link underline="never" :href="`/kubernetes/secrets/${scope.row.metadata.name}?cluster=${cluster}&namespace=${namespace}`">{{ scope.row.metadata.name }}</el-link>
                </template>
              </el-table-column>
              <el-table-column prop="type" label="类型" min-width="250" />
              <el-table-column prop="keys" label="字段" min-width="350" />
              <el-table-column label="创建时间" width="170">
                <template #default="scope">
                  {{ moment(scope.row.metadata.creationTimestamp).format('YYYY-MM-DD HH:mm:ss') }}
                </template>
              </el-table-column>
              <el-table-column label="操作" width="120">
                <template #default="scope">
                  <el-dropdown @command="handleCommand" style="vertical-align:middle;">
                    <el-button>更多操作<el-icon class="el-icon--right"><arrow-down /></el-icon></el-button>
                    <template #dropdown>
                      <el-dropdown-menu>
                        <el-dropdown-item :command="{action:'yaml',row:scope.row}">编辑YAML</el-dropdown-item>
                        <el-dropdown-item :command="{action:'delete',row:scope.row}">删除</el-dropdown-item>
                      </el-dropdown-menu>
                    </template>
                  </el-dropdown>
                </template>
              </el-table-column>
            </el-table>
            <el-pagination 
                class="pull-right"
                background 
                v-model:page-size="pageSize"
                :page-sizes="[10, 20, 50, 100]"
                layout="total, sizes, prev, pager, next, jumper" 
                :total="pageTotal"
                @size-change="handleSizeChange"
                @current-change="getPage"
                v-model:current-page="current" />
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>
  </div>
</div>
<el-drawer v-model="show.new" direction="rtl" size="700px">
  <template #header>
    <h4>新建保密字典</h4>
  </template>
  <template #default>
    <newWorkload ref="refWorkload" />
  </template>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show.new=false">取消</el-button>
      <el-button type="primary" @click="submitNew()">提交</el-button>
    </div>
  </template>
</el-drawer>
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
import { ArrowRight,Search,Refresh } from '@element-plus/icons-vue'
import { onBeforeMount, onBeforeUnmount, ref } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import newWorkload from '../new.vue'
import { axios } from '/src/assets/util/axios'
import moment from 'moment'
import _ from 'lodash'
/* 引入v-ace-editor */
import { VAceEditor } from 'vue3-ace-editor'
import 'ace-builds/src-noconflict/mode-yaml'
import 'ace-builds/src-noconflict/theme-chrome'
import 'ace-builds/src-noconflict/ext-language_tools'
/* 变量定义 */
const route = useRoute()
const activeName = ref("secrets")
const all = ref([])
const searchKey = ref("")
const list = ref([])
const pageSize = ref(10)
const pageTotal = ref(0)
const current = ref(1)
const cluster = ref("")
const clusterList = ref({})
const namespace = ref("")
const show = ref({
  new: false,
  yaml: false,
})
const loading = ref(false)
const selected = ref([])
const refWorkload = ref(null)
const wrapLine = ref(true)
const yamlContent = ref("")
const timer = ref(null)
/* 生命周期函数 */
onBeforeMount(async () => {
  await getClusterList()
  if(route.query.cluster && route.query.namespace) {
    cluster.value = route.query.cluster
    namespace.value = route.query.namespace
    getList(true)
  }
  timer.value = setInterval(() => {
    getList(false)
  }, 15000)
})
onBeforeUnmount(() => {
  if(timer) {
    clearInterval(timer.value)
    timer.value = null
  }
})
/* methods */
const getClusterList = async () => {
  clusterList.value = await axios.get(`/lizardcd/server/clusters`)
}
const getList = async (ifLoading) => {
  if(ifLoading) loading.value = true
  if(cluster.value && namespace.value) {
    let response = await axios.get(`/lizardcd/kubernetes/cluster/${cluster.value}/namespace/${namespace.value}/secrets`)
    all.value = _.sortBy(response.results, 'metadata.creationTimestamp').reverse()
    for(let x of all.value) {
      x.keys = Object.keys(x.data||{})
    }
    getPage(current.value)
  }
  if(ifLoading) loading.value = false
}
const getPage = async (page) => {
  let tmpList = all.value
  if(searchKey.value !== '') {
    tmpList = all.value.filter(n => n.name.includes(searchKey.value))
  }
  pageTotal.value = tmpList.length
  list.value = tmpList.slice((page-1)*pageSize.value, page*pageSize.value)
}
const handleCommand = async (command) => {
  switch(command.action) {
      case "yaml": {
      yamlContent.value = await axios.get(`/lizardcd/kubernetes/cluster/${cluster.value}/namespace/${namespace.value}/secrets/${command.row.metadata.name}/yaml`)
      show.value.yaml = true
      break
    }
    case "delete": {
      await ElMessageBox.confirm('确定删除？','警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(async () => {
        await axios.delete(`/lizardcd/kubernetes/cluster/${cluster.value}/namespace/${namespace.value}/secrets/${command.row.metadata.name}`)
        ElMessage.success({message: '删除成功'})
        setTimeout(async () => {
          await getList(true)
        }, 1000)
      }).catch(() =>{})
      break
    }
  }
}
const submitNew = async () => {
  refWorkload.value.getFormData(async (params) => {
    let cluster = params.cluster
    let namespace = params.namespace
    delete params.cluster
    delete params.namespace
    delete params.versions
    await axios.patch(`/lizardcd/kubernetes/cluster/${cluster}/namespace/${namespace}/apply/variable?kind=Secret`, params)
    ElMessage.success({message: '提交成功'})
    setTimeout(async () => {
      await getList(true)
    }, 1000)
    show.value.new = false
  })
}
const submitYaml = async () => {
  await axios.patch(`/lizardcd/kubernetes/cluster/${cluster.value}/namespace/${namespace.value}/apply/yaml`, yamlContent.value, {
    headers: {
      'Content-Type': 'text/plain'
    }
  })
  show.value.yaml = false
}
const handleMore = async (command) => {
  if(selected.value.length === 0) {
    ElMessage.warning({message: '请勾选保密字典'})
    return
  }
  switch(command.action) {
    case "deleteBatch": {
      await ElMessageBox.confirm('确定删除？','警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(async () => {
        for(let x of selected.value) {
          await axios.delete(`/lizardcd/kubernetes/cluster/${cluster.value}/namespace/${namespace.value}/secrets/${x.metadata.name}`)
        }
        ElMessage.success({message: '删除成功'})
      }).catch(() =>{})
      break
    }
  }
  setTimeout(async () => {
    await getList(true)
  }, 1000)
}
const select = (val) => {
  selected.value = val
}
const handleSizeChange = async (size) => {
  pageSize.value = size
  await getList(true)
}
</script>
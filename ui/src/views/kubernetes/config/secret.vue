<template>
<el-breadcrumb :separator-icon="ArrowRight" style="min-width:1366px">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>配置</el-breadcrumb-item>
  <el-breadcrumb-item :to="{ path: '/kubernetes/secrets', query: {cluster: route.query.cluster, namespace: route.query.namespace} }">保密字典</el-breadcrumb-item>
  <el-breadcrumb-item>{{ secretInfo.metadata?.name }}</el-breadcrumb-item>
</el-breadcrumb>
<el-row :gutter="15" style="min-width:1366px">
  <el-col :span="6">
    <el-card>
      <template #header>
        <div class="card-header">
          <el-text truncated size="large" style="width:250px;font-size:18px" class="text-black"><b>{{ secretInfo.metadata?.name }}</b></el-text>
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
                <el-dropdown-item command="edit">编辑设置</el-dropdown-item>
                <el-dropdown-item command="delete">删除</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </template>
      <el-descriptions :column="1" border class="no-color" :label-width="120">
        <el-descriptions-item label="集群">{{ route.query.cluster }}</el-descriptions-item>
        <el-descriptions-item label="命名空间">{{ route.query.namespace }}</el-descriptions-item>
        <el-descriptions-item label="类型">{{ secretInfo.type }}</el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ moment(secretInfo.metadata?.creationTimestamp).format('YYYY-MM-DD HH:mm:ss') }}</el-descriptions-item>
      </el-descriptions>
    </el-card>
  </el-col>
  <el-col :span="18">
    <el-card header-class="header-height-60">
      <template #header>
        <div class="card-header">
          <span class="card-header-text">数据</span>
          <span class="card-header-btn" @click="hideData">
            <el-icon v-if="!hide"><View /></el-icon>
            <el-icon v-else><Hide /></el-icon>
          </span>
        </div>
      </template>
      <el-descriptions :column="1" border class="no-color" :label-width="150">
        <el-descriptions-item v-for="(v,k,i) in secretInfo.data" :key="i" :label="k" class-name="no-wrap">
          <div v-html="data[k]" />
        </el-descriptions-item>
      </el-descriptions>
    </el-card>
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
<el-dialog v-model="show.edit" title="编辑保密字典" @closed="afterEdit()" width="70%">
  <el-card v-for="(item,i) in editdata" :key="i" header-class="header-height-50">
    <template #header>
      <div class="card-header">
        <span class="card-header-text">{{ item.key }}</span>
        <span class="card-header-btn">
          <el-icon @click="removeOne(i)"><Close /></el-icon>
        </span>
      </div>
    </template>
    <el-form :model="item" label-position="top">
      <el-form-item label="键">
        <el-input v-model="item.key" size="large" />
      </el-form-item>
      <el-form-item label="值">
        <v-ace-editor
          v-model:value="item.value"
          lang="text"
          theme="chrome"
          style="width:100%;"
          :options="{
            enableBasicAutocompletion: true,
            enableSnippets: true,
            enableLiveAutocompletion: true,
            tabSize: 2,
            showPrintMargin: false,
            fontSize: 14,
            minLines: 5,
            maxLines: 5000
          }" />
      </el-form-item>
    </el-form>
  </el-card>
  <el-button circle icon="Plus" @click="addOne()" />
  <template #footer>
    <div class="dialog-footer">
      <el-button @click="show.edit=false">取消</el-button>
      <el-button type="primary" @click="submitEdit">提交</el-button>
    </div>
  </template>
</el-dialog>
</template>

<script setup>
import { ArrowRight, Check, EditPen, View, Hide } from '@element-plus/icons-vue'
import { onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRoute, useRouter } from 'vue-router'
import { axios } from '/src/assets/util/axios'
import moment from 'moment'
import yaml from 'js-yaml'
import { codeToHtml } from 'shiki'
/* 引入v-ace-editor */
import { VAceEditor } from 'vue3-ace-editor'
import 'ace-builds/src-noconflict/mode-yaml'
import 'ace-builds/src-noconflict/mode-text'
import 'ace-builds/src-noconflict/theme-chrome'
import 'ace-builds/src-noconflict/ext-language_tools'
import _ from 'lodash'
/* 变量定义 */
const route = useRoute()
const router = useRouter()
const secretInfo = ref({})
const show = ref({
  yaml: false,
  edit: false,
})
const yamlContent = ref("")
const wrapLine = ref(true)
const edit = ref({})
const hide = ref(true)
const data = ref({})
const editdata = ref([])
/* 生命周期函数 */
onMounted(async () => {
  await getConfigMap()
})
/* methods */
const getConfigMap = async () => {
  secretInfo.value = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/secrets/${route.params.secret_name}`)
  if(!secretInfo.value.metadata) {
    ElMessage.error({message: `未找到保密字典: ${route.params.secret_name}`})
    return
  }
  editdata.value = []
  for(let key of Object.keys(secretInfo.value.data)) {
    editdata.value.push({
      key,
      value: secretInfo.value.data[key]
    })
    edit.value[key] = false
    data.value[key] = await codeToHtml(secretInfo.value.data[key], { lang: 'text', theme: 'github-dark' })
  }
}
const handleCommand = async (command) => {
  switch(command) {
    case "yaml": {
      yamlContent.value = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/secrets/${route.params.secret_name}/yaml`)
      show.value.yaml = true
      break
    }
    case "delete": {
      await ElMessageBox.confirm('确定删除？','警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(async () => {
        await axios.delete(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/secrets/${route.params.secret_name}`)
        ElMessage.success({message: '删除成功'})
        router.push({
          path: '/kubernetes/secrets',
          query: {
            cluster: route.query.cluster,
            namespace: route.query.namespace
          }
        })
      }).catch(() =>{})
      break
    }
    case "edit": {
      for(let x of editdata.value) {
        x.value = atob(x.value)
      }
      show.value.edit = true
      break
    }
  }
}
const submitYaml = async () => {
  await axios.patch(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/apply/yaml`, yamlContent.value, {
    headers: {
      'Content-Type': 'text/plain'
    }
  })
  show.value.yaml = false
  await getConfigMap()
}
const hideData = async () => {
  hide.value = !hide.value
  for(let k of Object.keys(secretInfo.value.data)) {
    if(!hide.value) {
      secretInfo.value.data[k] = atob(secretInfo.value.data[k])
    } else {
      secretInfo.value.data[k] = btoa(secretInfo.value.data[k])
    }
    data.value[k] = await codeToHtml(secretInfo.value.data[k], { lang: 'text', theme: 'github-dark' })
  }
}
const removeOne = (index) => {
  editdata.value.splice(index, 1)
}
const addOne = () => {
  editdata.value.push({
    key: '',
    value: ''
  })
}
const submitEdit = async () => {
  let data = {}
  for(let x of editdata.value) {
    data[x.key] = btoa(x.value)
  }
  let params = _.cloneDeep(secretInfo.value)
  params.data = data
  params.apiVersion = "v1"
  params.kind = "Secret"
  yamlContent.value = yaml.dump(params)
  show.value.edit = false
  await submitYaml()
  for(let x of editdata.value) {
    x.value = atob(x.value)
  }
}
const afterEdit = () => {
  for(let x of editdata.value) {
    x.value = btoa(x.value)
  }
}
</script>
<template>
<el-breadcrumb :separator-icon="ArrowRight">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>YAML模板</el-breadcrumb-item>
</el-breadcrumb>
<el-card>
  <template #header>
    <div class="card-header">
      <span class="card-header-text">YAML模板</span>
    </div>
  </template>
  <el-row>
    <el-alert title="关于模板用法参考Go-template：https://pkg.go.dev/text/template" type="warning" style="margin-bottom:15px" />
    <el-col :span="12">
      <el-button-group>
        <el-button :icon="Refresh" size="large" style="margin-right:5px" @click="getList(current)" />
        <el-input v-model="searchKey" placeholder="输入名称进行搜索" size="large" :prefix-icon="Search" @change="getPage(1)" clearable style="width:300px;" />
      </el-button-group>
    </el-col>
    <el-col :span="12">
      <el-button class="pull-right" size="large" type="primary" @click="show=true;edit=false;form={content:'',variables:[{key:'',value:''}],tenant:tenant}">新建模板</el-button>
    </el-col>
  </el-row>
  <el-table 
    :data="list" 
    v-loading="loading.table"
    element-loading-text="奋力加载中..."
    class="line-height40" 
    style="width:100%;margin-top:10px">
    <el-table-column prop="name" label="模板名称" min-width="160" />
    <el-table-column prop="tenant" label="所属租户" min-width="80" />
    <el-table-column prop="update_at" label="更新时间" width="160">
      <template #default="scope">
        {{ moment(scope.row.update_at).format('YYYY-MM-DD HH:mm') }}
      </template>
    </el-table-column>
    <el-table-column prop="Option" label="操作" width="140">
      <template #default="scope">
        <el-button :icon="EditPen" circle @click="editOne(scope.row)" />
        <el-tooltip effect="dark" content="复制" placement="top">
          <el-button :icon="CopyDocument" circle @click="copyOne(scope.row)" />
        </el-tooltip>
        <el-popconfirm title="确认删除？" @confirm="deleteOne(scope.row)">
          <template #reference>
            <el-button :icon="Delete" circle />
          </template>
        </el-popconfirm>
      </template>
    </el-table-column>
  </el-table>
  <el-pagination 
    class="pull-right"
    background 
    v-model:page-size="pageSize"
    :page-sizes="[10, 30, 50, 100]"
    layout="total, sizes, prev, pager, fnext, jumper" 
    :total="pageTotal"
    @size-change="handleSizeChange"
    @current-change="getList"
    v-model:current-page="current" />
</el-card>
<el-drawer v-model="show" direction="rtl" size="800px">
  <template #header>
    <h4 v-if="edit===false">新建模板</h4>
    <h4 v-if="edit===true">编辑模板</h4>
  </template>
  <template #default>
    <el-form ref="template" :model="form" :rules="rules" label-width="120px">
      <el-form-item label="模板名称" prop="name">
        <el-input v-model="form.name" size="large" />
      </el-form-item>
      <el-form-item label="模板定义" prop="content">
        <v-ace-editor
          v-model:value="form.content"
          lang="yaml"
          theme="chrome"
          style="width:100%;height:700px"
          :options="{
            enableBasicAutocompletion: true,
            enableSnippets: true,
            enableLiveAutocompletion: true,
            tabSize: 2,
            showPrintMargin: false,
            fontSize: 14
          }" />
      </el-form-item>
      <el-form-item label="模板变量">
        <table class="table table-bordered">
          <thead><tr><th>变量名</th><th>默认变量值</th></tr></thead>
          <tbody>
          <tr v-for="(item,index) in form.variables" :key="index" >
            <td><el-input v-model="item.key" size="large" /></td>
            <td><el-input v-model="item.value" size="large" /></td>
            <td width="100">
              <el-button icon="Plus" circle @click="addVar(index)"></el-button>
              <el-button icon="Close" circle @click="removeVar(index)"></el-button>
            </td>
          </tr>
          </tbody>
        </table>
      </el-form-item>
      <el-form-item label="所属租户" prop="tenant">
        <el-input v-model="form.tenant" disabled size="large" />
      </el-form-item>
    </el-form>
  </template>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show=false" size="large">取消</el-button>
      <el-button type="primary" @click="confirmClick(template)" size="large">提交</el-button>
    </div>
  </template>
</el-drawer>
</template>

<script setup>
import { ArrowRight,Search,Refresh,EditPen,CopyDocument,Delete } from '@element-plus/icons-vue'
import { onBeforeMount, ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { axios } from '/src/assets/util/axios'
import moment from 'moment'
import _ from 'lodash'
/* 引入v-ace-editor */
import { VAceEditor } from 'vue3-ace-editor'
import 'ace-builds/src-noconflict/mode-yaml'
import 'ace-builds/src-noconflict/theme-chrome'
import 'ace-builds/src-noconflict/ext-language_tools'
/* 变量定义 */
const tenant = localStorage.tenant
const list = ref([])
const pageSize = ref(10)
const pageTotal = ref(0)
const current = ref(1)
const searchKey = ref("")
const loading = ref({
  table: false
})
const show = ref(false)
const edit = ref(false)
const form = ref({content:'',variables:[{key:'',value:''}]})
const rules = reactive({
  name: [{required: true, message: '请填写模板名称'}],
  content: [{required: true, message: '请填写模板定义'}],
})
const template = ref(null)
/* 生命周期函数 */
onBeforeMount(async () => {
  getList(1)
})
/* methods */
const getList = async (page) => {
  let url = `page=${page}&size=${pageSize.value}&sort=update_at desc`
  if(searchKey.value != "") url += `&search=name==${searchKey.value}`
  let response = await axios.get(`/lizardcd/db/application_template?${url}`)
  list.value = response.results.map(x => {
    x.variables = JSON.parse(x.variables)
    return x
  })
  pageTotal.value = response.total
}
const editOne = async (row) => {
  form.value = row
  form.value.tenant ||= localStorage.tenant
  edit.value = true
  show.value = true
}
const copyOne = async (row) => {
  form.value = Object.assign({}, row)
  delete form.value.id
  edit.value = false
  show.value = true
}
const addVar = (index) => {
  form.value.variables.splice(index+1, 0, {key:"", value:""})
}
const removeVar = (index) => {
  form.value.variables.splice(index, 1)
}
const confirmClick = async (f) => {
  if(!f) return
  await f.validate(async (valid) => {
    if(valid) {
      let params = Object.assign({}, form.value)
      params.variables = JSON.stringify(params.variables)
      params.update_at = moment()
      if(edit.value === false) {
        await axios.post(`/lizardcd/db/application_template`, {body:params})
        getList(1)
        current.value = 1
        show.value = false
      }
      else {
        let id = params.id
        delete params.id
        await axios.put(`/lizardcd/db/application_template/${id}`, {body:params})
        getList(current.value)
        show.value = false
      }
    }
    else {
      ElMessage.warning('必填项未填完')
    }
  })
}
const deleteOne = async (row) => {
  await axios.delete(`/lizardcd/db/application_template/${row.id}`)
  getList(current.value)
}
const handleSizeChange = async (size) => {
  pageSize.value = size
  await getList(current.value)
}
</script>
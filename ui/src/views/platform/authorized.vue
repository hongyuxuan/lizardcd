<template>
<div class="box box-item">
  <div class="box-body" style="padding-top:20px;padding-bottom:0">
    <el-row>
      <el-col :span="12">
        <el-button-group>
          <el-button icon="refresh" size="large" style="margin-right:5px" @click="getList(1)" />
          <el-input v-model="searchKey" placeholder="输入名称进行搜索" size="large" :prefix-icon="Search" @change="getPage(1)" clearable style="width:300px;" />
        </el-button-group>
      </el-col>
      <el-col :span="12">
        <el-button class="pull-right" size="large" type="primary" @click="show=true;form={}">新建认证</el-button>
      </el-col>
    </el-row>
    <el-table 
      :data="list"
      class="line-height40" 
      style="width:100%;margin-top:10px">
      <el-table-column prop="name" label="名称" min-width="150" />
      <el-table-column prop="authorize_url" label="authorize_url" min-width="350" />
      <el-table-column prop="Option" label="操作" width="130">
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
  </div>
</div>
<el-drawer v-model="show" direction="rtl" size="600px">
  <template #header>
    <h4 v-if="edit===false">新建认证</h4>
    <h4 v-if="edit===true">编辑认证</h4>
  </template>
  <template #default>
    <el-form ref="authorize" :model="form" :rules="rules" label-width="100px">
      <el-form-item label="认证名称" prop="name">
        <el-input v-model="form.name" size="large" />
      </el-form-item>
      <el-divider><span style="color:#b4b4b4">Oauth2配置</span></el-divider>
      <el-form-item label="client_id" prop="client_id">
        <el-input v-model="form.client_id" size="large" />
        <myTips type="info">第三方认证平台颁发的ClientID</myTips>
      </el-form-item>
      <el-form-item label="client_secret" prop="client_secret">
        <el-input v-model="form.client_secret" size="large" />
        <myTips type="info">第三方认证平台颁发的ClientSecret</myTips>
      </el-form-item>
      <el-form-item label="authorize_url" prop="authorize_url">
        <el-input v-model="form.authorize_url" size="large" />
        <myTips type="info">第三方认证平台授权接口地址</myTips>
      </el-form-item>
      <el-form-item label="redirect_url" prop="redirect_url">
        <el-input v-model="form.redirect_url" size="large" />
        <myTips type="info">Lizardcd重定向地址，固定写法：<br>http://{SERVER_URL:PORT}/lizardcd/auth/callback/{认证名称}</myTips>
      </el-form-item>
      <el-form-item label="token_url" prop="token_url">
        <el-input v-model="form.token_url" size="large" />
        <myTips type="info">第三方认证平台获取token接口地址</myTips>
      </el-form-item>
      <el-form-item label="userinfo_url" prop="userinfo_url">
        <el-input v-model="form.userinfo_url" size="large" />
      </el-form-item>
      <el-form-item label="callback_url" prop="callback_url">
        <el-input v-model="form.callback_url" size="large" />
        <myTips type="info">Lizardcd回调地址，固定写法：<br>http://{UI_URL:PORT}/callback.html</myTips>
      </el-form-item>
      <el-divider><span style="color:#b4b4b4">属性映射</span></el-divider>
      <el-form-item label="用户账户" prop="user_jsonpath">
        <el-input v-model="form.user_jsonpath" size="large" placeholder="$.username" />
        <myTips type="info">通过jsonpath从第三方平台提取用户名</myTips>
      </el-form-item>
      <el-form-item label="用户头像" prop="avatar_jsonpath">
        <el-input v-model="form.avatar_jsonpath" size="large" placeholder="$.avatar_url" />
        <myTips type="info">通过jsonpath从第三方平台提取头像</myTips>
      </el-form-item>
      <el-divider><span style="color:#b4b4b4">代理设置</span></el-divider>
      <el-form-item label="HTTP_PROXY" prop="http_proxy">
        <el-input v-model="form.http_proxy" size="large" />
        <myTips type="info">如果Lizardcd服务端无法直接连通第三方平台，可在此配置代理</myTips>
      </el-form-item>
    </el-form>
  </template>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show=false">取消</el-button>
      <el-button type="primary" @click="confirmClick(authorize)">提交</el-button>
    </div>
  </template>
</el-drawer>
</template>
<script setup>
import { Search,EditPen,CopyDocument,Delete } from '@element-plus/icons-vue'
import { onBeforeMount, ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { axios } from '/src/assets/util/axios'
import _ from 'lodash'
/* 变量定义 */
const list = ref([])
const searchKey = ref("")
const show = ref(false)
const form = ref({})
const rules = reactive({
  name: [{required: true, message: '请填写名称'}],
  client_id: [{required: true, message: '请填写client_id'}],
  client_secret: [{required: true, message: '请填写client_secret'}],
  authorize_url: [{required: true, message: '请填写authorize_url'}],
  redirect_url: [{required: true, message: '请填写redirect_url'}],
  token_url: [{required: true, message: '请填写token_url'}],
  userinfo_url: [{required: true, message: '请填写userinfo_url'}],
  callback_url: [{required: true, message: '请填写callback_url'}],
  user_jsonpath: [{required: true, message: '请填写用户账户'}],
})
const edit = ref(false)
const authorize = ref(null)
/* 生命周期函数 */
onBeforeMount(async () => {
  getList()
})
/* methods */
const getList = async () => {
  let url = ""
  if(searchKey.value != "") url += `?search=name==${searchKey.value}`
  let response = await axios.get(`/lizardcd/db/oauth2?${url}`)
  list.value = response.results
}
const confirmClick = async (f) => {
  if(!f) return
  await f.validate(async (valid) => {
    if(valid) {
      let params = Object.assign({}, form.value)
      if(edit.value === false) {
        await axios.post(`/lizardcd/db/oauth2`, {body:params})
        getList()
        show.value = false
      }
      else {
        let id = params.id
        delete params.id
        await axios.put(`/lizardcd/db/oauth2/${id}`, {body:params})
        getList()
        show.value = false
      }
    }
    else {
      ElMessage.warning('必填项未填完')
    }
  })
}
const editOne = async (row) => {
  form.value = Object.assign({}, row)
  edit.value = true
  show.value = true
}
const deleteOne = async (row) => {
  await axios.delete(`/lizardcd/db/oauth2/${row.id}`)
  getList()
}
const copyOne = async (row) => {
  form.value = Object.assign({}, row)
  delete form.value.id
  edit.value = false
  show.value = true
}
</script>
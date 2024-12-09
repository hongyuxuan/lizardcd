<template>
<div class="box box-item">
  <div class="box-body" style="padding-top:20px;padding-bottom:0">
    <el-row>
      <el-col :span="12">
        <el-button-group>
          <el-button icon="refresh" size="large" style="margin-right:5px" @click="getList(1)" />
          <el-input v-model="searchKey" placeholder="输入名称进行搜索" size="large" :prefix-icon="Search" @change="current=1;getList(1)" clearable style="width:300px;" />
        </el-button-group>
      </el-col>
      <el-col :span="12">
        <el-button class="pull-right" size="large" type="primary" @click="show=true;edit=false;form={tenant}">新建仓库</el-button>
      </el-col>
    </el-row>
    <el-table 
      :data="list"
      class="line-height40" 
      style="width:100%;margin-top:10px">
      <el-table-column prop="git_http_url" label="仓库地址" min-width="300" />
      <el-table-column prop="tenant" label="所属租户" min-width="80" />
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
      layout="total, sizes, prev, pager, next, jumper" 
      :total="pageTotal"
      @size-change="handleSizeChange"
      @current-change="getList"
      v-model:current-page="current" />
  </div>
</div>
<el-drawer v-model="show" direction="rtl" size="600px">
  <template #header>
    <h4 v-if="edit===false">新建仓库</h4>
    <h4 v-if="edit===true">编辑仓库</h4>
  </template>
  <template #default>
    <el-form ref="repo" :model="form" :rules="rules" label-width="150px">
      <el-form-item label="仓库地址" prop="git_http_url">
        <el-input v-model="form.git_http_url" size="large" clearable placeholder="http(s)://" />
      </el-form-item>
      <el-form-item>
        <template #label><el-text>AccessToken 
          <el-tooltip placement="top">
            <template #content>
              关于如何获取AccessToken参见：http(s)://&lt;your_git_url&gt;/-/user_settings/personal_access_tokens
            </template>
            <el-icon class="text-yellow"><Warning /></el-icon>
          </el-tooltip></el-text>
        </template>
        <el-input v-model="form.access_token" type="password" size="large" clearable />
      </el-form-item>
      <el-form-item>
        <template #label><el-text>WebhookBaseURL 
          <el-tooltip placement="top">
            <template #content>
              该地址用于在代码仓库自动配置触发构建流水线的基准URL<br>通常设置为Lizardcd Server的对外服务地址（需从您的gitlab能够访问）
            </template>
            <el-icon class="text-yellow"><Warning /></el-icon>
          </el-tooltip></el-text>
        </template>
        <el-input v-model="form.webhook_base_url" size="large" clearable />
      </el-form-item>
      <el-form-item label="DefaultSecret" prop="secret">
        <el-input v-model="form.secret" size="large" type="password" show-password />
      </el-form-item>
      <el-form-item label="所属租户" prop="tenant">
        <el-input v-model="form.tenant" disabled size="large" />
      </el-form-item>
    </el-form>
  </template>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show=false">取消</el-button>
      <el-button type="primary" @click="confirmClick(repo)">提交</el-button>
    </div>
  </template>
</el-drawer>
</template>

<script setup>
import { Search,EditPen,CopyDocument,Delete } from '@element-plus/icons-vue'
import { onBeforeMount, ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { axios } from '/src/assets/util/axios'
/* 变量定义 */
const tenant = localStorage.tenant.split(",")[0]
const list = ref([])
const pageSize = ref(10)
const pageTotal = ref(0)
const current = ref(1)
const searchKey = ref("")
const show = ref(false)
const form = ref({})
const rules = reactive({
  git_http_url: [{required: true, message: '请填写仓库地址'}],
  secret: [{required: true, message: '请填写默认secret'}],
})
const edit = ref(false)
const repo = ref(null)
/* 生命周期函数 */
onBeforeMount(async () => {
  getList(1)
})
/* methods */
const getList = async (page) => {
  let url = `page=${page}&size=${pageSize.value}`
  if(searchKey.value != "") url += `&search=git_http_url==${searchKey.value}`
  let response = await axios.get(`/lizardcd/db/git_repository?${url}`)
  list.value = response.results
  pageTotal.value = response.total
}
const confirmClick = async (f) => {
  if(!f) return
  await f.validate(async (valid) => {
    if(valid) {
      let params = Object.assign({}, form.value)
      if(edit.value === false) {
        await axios.post(`/lizardcd/db/git_repository`, {body:params})
        getList(1)
        current.value = 1
      }
      else {
        await axios.put(`/lizardcd/db/git_repository/${params.id}`, {body:params})
        getList(current.value)
      }
      show.value = false
    }
    else {
      ElMessage.warning('必填项未填完')
    }
  })
}
const editOne = async (row) => {
  form.value = Object.assign({}, row)
  form.value = row
  form.value.tenant ||= localStorage.tenant
  edit.value = true
  show.value = true
}
const deleteOne = async (row) => {
  await axios.delete(`/lizardcd/db/git_repository/${row.id}`)
  getList(current.value)
}
const copyOne = async (row) => {
  form.value = Object.assign({}, row)
  delete form.value.id
  edit.value = false
  show.value = true
}
const handleSizeChange = async (size) => {
  pageSize.value = size
  await getList(current.value)
}
</script>
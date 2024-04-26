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
        <el-button class="pull-right" size="large" type="primary" @click="show=true;form={}">新建仓库</el-button>
      </el-col>
    </el-row>
    <el-table 
      :data="list"
      class="line-height40" 
      style="width:100%;margin-top:10px">
      <el-table-column prop="repo_url" label="仓库地址" min-width="300" />
      <el-table-column prop="repo_account" label="仓库账户" width="150" />
      <el-table-column prop="repo_type" label="仓库类型" width="100" />
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
      :page-sizes="[20, 30, 50, 100]"
      layout="total, sizes, prev, pager, next, jumper" 
      :total="pageTotal"
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
    <el-form ref="repo" :model="form" :rules="rules" label-width="100px">
      <el-form-item label="仓库地址" prop="repo_url">
        <el-input v-model="form.repo_url" size="large" />
      </el-form-item>
      <el-form-item label="仓库账号" prop="repo_account">
        <el-input v-model="form.repo_account" size="large" />
      </el-form-item>
      <el-form-item label="仓库密码" prop="repo_password">
        <el-input v-model="form.repo_password" type="password" size="large" />
      </el-form-item>
      <el-form-item label="仓库类型" prop="repo_type">
        <el-select v-model="form.repo_type" size="large">
          <el-option label="Artifactory" value="Artifactory" />
          <el-option label="Harbor" value="Harbor" />
        </el-select>
      </el-form-item>
    </el-form>
  </template>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show=false" size="large">取消</el-button>
      <el-button type="primary" @click="confirmClick(repo)" size="large">提交</el-button>
    </div>
  </template>
</el-drawer>
</template>

<script setup>
import { ArrowRight,Search,Refresh,EditPen,CopyDocument,Delete } from '@element-plus/icons-vue'
import { onBeforeMount, ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { axios } from '/src/assets/util/axios'
/* 变量定义 */
const all = ref([])
const list = ref([])
const pageSize = ref(10)
const pageTotal = ref(0)
const current = ref(1)
const searchKey = ref("")
const show = ref(false)
const form = ref({})
const rules = reactive({
  repo_url: [{required: true, message: '请填写仓库地址'}],
  repo_account: [{required: true, message: '请填写仓库账户'}],
  repo_type: [{required: true, message: '请选择仓库类型', trigger: 'change'}],
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
  if(searchKey.value != "") url += `&search=repo_url==${searchKey.value}`
  let response = await axios.get(`/db/image_repository?${url}`)
  list.value = response.results
  pageTotal.value = response.total
}
const confirmClick = async (f) => {
  if(!f) return
  await f.validate(async (valid) => {
    if(valid) {
      let params = Object.assign({}, form.value)
      if(edit.value === false) {
        await axios.post(`/db/image_repository`, {body:params})
        getList(1)
        current.value = 1
        show.value = false
      }
      else {
        let id = params.id
        delete params.id
        await axios.put(`/db/image_repository/${id}`, {body:params})
        getList(current.value)
        show.value = false
      }
    }
    else {
      ElMessage.warning('必填项未填完')
    }
  })
}
const editOne = async (row) => {
  form.value = row
  edit.value = true
  show.value = true
}
const deleteOne = async (row) => {
  await axios.delete(`/db/image_repository/${row.id}`)
  getList(current.value)
}
const copyOne = async (row) => {
  form.value = Object.assign({}, row)
  delete form.value.id
  edit.value = false
  show.value = true
}
</script>
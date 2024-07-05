<template>
<div class="box box-item">
  <div class="box-body" style="padding-top:20px;padding-bottom:0">
    <el-row>
      <el-col :span="12">
        <el-button-group>
          <el-button icon="refresh" size="large" style="margin-right:5px" @click="getList(1)" />
          <el-input v-model="searchKey" placeholder="输入名称进行搜索" size="large" :prefix-icon="Search" @change="getList(1);current=1" clearable style="width:300px;" />
        </el-button-group>
      </el-col>
      <el-col :span="12">
        <el-button class="pull-right" size="large" type="primary" @click="show=true;edit=false;form={role:'admin'}">新建用户</el-button>
      </el-col>
    </el-row>
    <el-table 
      :data="list"
      class="line-height40" 
      style="width:100%;margin-top:10px">
      <el-table-column prop="username" label="用户名" min-width="150" />
      <el-table-column prop="tenant" label="所属租户" min-width="150" />
      <el-table-column prop="role" label="用户权限" min-width="150" />
      <el-table-column prop="update_at" label="更新时间" width="170">
        <template #default="scope">
        {{ moment(scope.row.update_at).format('YYYY-MM-DD HH:mm:ss') }}
      </template>
      </el-table-column>
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
    <el-pagination 
      class="pull-right"
      background 
      v-model:page-size="pageSize"
      :page-sizes="[20, 30, 50, 100]"
      layout="total, sizes, prev, pager, next, jumper" 
      :total="pageTotal"
      @size-change="handleSizeChange"
      @current-change="getList"
      v-model:current-page="current" />
  </div>
</div>
<el-drawer v-model="show" direction="rtl" size="600px">
  <template #header>
    <h4 v-if="edit===false">新建用户</h4>
    <h4 v-if="edit===true">编辑用户</h4>
  </template>
  <template #default>
    <el-form ref="user" :model="form" :rules="rules" label-width="100px">
      <el-form-item label="用户名" prop="username">
        <el-input v-model="form.username" size="large" :disabled="edit" />
      </el-form-item>
      <el-form-item label="所属租户" prop="tenant">
        <el-select 
          v-model="form.tenant" 
          placeholder="请选择租户" 
          clearable
          filterable
          value-key="id"
          style="width:100%" 
          size="large">
          <el-option v-for="(item,i) in tenants" :key="i" :label="item.tenant_name" :value="item" />
        </el-select>
      </el-form-item>
      <el-form-item label="权限" prop="role">
        <el-radio-group v-model="form.role" style="">
          <el-radio value="admin" size="large">admin</el-radio>
          <el-radio value="readonly" size="large">readonly</el-radio>
          <el-radio value="readwrite" size="large">readwrite</el-radio>
        </el-radio-group>
      </el-form-item>
    </el-form>
  </template>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show=false" size="large">取消</el-button>
      <el-button type="primary" @click="confirmClick(user)" size="large">提交</el-button>
    </div>
  </template>
</el-drawer>
</template>

<script setup>
import { Search,EditPen,CopyDocument,Delete } from '@element-plus/icons-vue'
import { onBeforeMount, ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { axios } from '/src/assets/util/axios'
import moment from 'moment'
/* 变量定义 */
const list = ref([])
const pageSize = ref(10)
const pageTotal = ref(0)
const current = ref(1)
const searchKey = ref("")
const show = ref(false)
const form = ref({
  role: 'admin'
})
const rules = reactive({
  username: [{required: true, message: '请填写用户名'}],
  password: [{required: true, message: '请填写密码'}],
})
const tenants = ref([])
const edit = ref(false)
const user = ref(null)
/* 生命周期函数 */
onBeforeMount(async () => {
  getList(1)
  getTenants()
})
/* methods */
const getList = async (page) => {
  let url = `page=${page}&size=${pageSize.value}`
  if(searchKey.value != "") url += `&search=username==${searchKey.value}`
  let response = await axios.get(`/lizardcd/db/user?${url}`)
  list.value = response.results
  pageTotal.value = response.total
}
const getTenants = async () => {
  let response = await axios.get(`/lizardcd/db/tenant`)
  tenants.value = response.results
}
const confirmClick = async (f) => {
  if(!f) return
  await f.validate(async (valid) => {
    if(valid) {
      let params = Object.assign({}, form.value)
      params.tenant = form.value.tenant.tenant_name
      if(edit.value === false) {
        let password = await axios.post(`/lizardcd/auth/adduser`, params)
        getList(1)
        current.value = 1
        show.value = false
        ElMessageBox.alert(`这是为您新建用户生成的临时密码：<code>${password}</code><br>仅出现在此对话框一次，请牢记并及时修改密码`, '密码告知', {
          confirmButtonText: '知道了',
          dangerouslyUseHTMLString: true,
        })
      }
      else {
        let id = params.id
        delete params.id
        delete params.password
        await axios.put(`/lizardcd/db/user/${id}`, {body:params})
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
  form.value = Object.assign({}, row)
  form.value.tenant = tenants.value.find(n => n.tenant_name === row.tenant)
  edit.value = true
  show.value = true
}
const deleteOne = async (row) => {
  await axios.delete(`/lizardcd/db/user/${row.id}`)
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
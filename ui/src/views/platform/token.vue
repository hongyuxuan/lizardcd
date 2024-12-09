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
        <el-button class="pull-right" size="large" type="primary" @click="show=true;edit=false;form={}">新建令牌</el-button>
      </el-col>
    </el-row>
    <el-table 
      :data="list"
      class="line-height25" 
      style="width:100%;margin-top:10px">
      <el-table-column prop="username" label="用户" min-width="100" />
      <el-table-column prop="jwt_token" label="令牌" min-width="200" />
      <el-table-column prop="create_at" label="创建时间" width="180">
        <template #default="props">
          {{ moment(props.row.create_at).format('YYYY-MM-DD HH:mm:SS') }}
        </template>
      </el-table-column>
      <el-table-column prop="expire_at" label="过期时间" width="180">
        <template #default="props">
          {{ moment(props.row.expire_at).format('YYYY-MM-DD HH:mm:SS') }}
        </template>
      </el-table-column>
      <el-table-column prop="Option" label="操作" width="100">
        <template #default="scope">
          <el-button class="copy" :icon="DocumentCopy" :data-clipboard-text="scope.row.jwt_token" circle @click="copyDocument" />
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
    <h4>新建令牌</h4>
  </template>
  <template #default>
    <el-form ref="token" :model="form" :rules="rules" label-width="100px">
      <el-form-item label="选择用户" prop="user_id">
        <el-select v-model="form.user_id" size="large" clearable value-key="id">
          <el-option v-for="(item,i) in userList" :key="i" :label="item.username" :value="item" />
        </el-select>
      </el-form-item>
      <el-form-item label="过期时间" prop="expire_at">
        <el-date-picker v-model="form.expire_at" type="date" placeholder="请选择时间" size="large" clearable style="width:100%" />
      </el-form-item>
      <el-form-item label="令牌" prop="jwt_token">
          <el-input v-model="form.jwt_token" disabled size="large" placeholder="点击生成令牌" style="width:410px;float:left;margin-right:5px" />
          <el-button circle :icon="RefreshRight" size="large" @click="generateToken(token)" />
      </el-form-item>
    </el-form>
  </template>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show=false">取消</el-button>
      <el-button type="primary" @click="confirmClick(token)">提交</el-button>
    </div>
  </template>
</el-drawer>
</template>

<script setup>
import { Search,RefreshRight,Delete,DocumentCopy } from '@element-plus/icons-vue'
import { onBeforeMount, ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { axios } from '/src/assets/util/axios'
import moment from 'moment'
import { copyDocument } from '/src/assets/util/clipboard.js'
/* 变量定义 */
const props = defineProps({
  count: {
    type: String,
  }
})
const list = ref([])
const pageSize = ref(10)
const pageTotal = ref(0)
const current = ref(1)
const searchKey = ref("")
const show = ref(false)
const form = ref({})
const rules = reactive({
  user_id: [{required: true, message: '请选择用户', trigger: 'change'}],
  expire_at: [{required: true, message: '请选择过期时间', trigger: 'blur'}],
})
const edit = ref(false)
const token = ref(null)
const userList = ref([])
/* 生命周期函数 */
onBeforeMount(async () => {
  getList(1)
  getUserList()
})
/* methods */
const getList = async (page) => {
  let url = `page=${page}&size=${pageSize.value}`
  if(searchKey.value != "") url += `&search=username==${searchKey.value}`
  let response = await axios.get(`/lizardcd/db/tokens?${url}`)
  list.value = response.results
  pageTotal.value = response.total
}
const getUserList = async () => {
  let response = await axios.get(`/lizardcd/db/user`)
  userList.value = response.results
}
const confirmClick = async (f) => {
  await f.validate(async (valid) => {
    if(valid) {
      let params = Object.assign({}, form.value)
      params.username = params.user_id.username
      params.user_id = params.user_id.id
      params.create_at = moment()
      await axios.post(`/lizardcd/db/tokens`, {body:params})
      getList(current.value)
      show.value = false
    }
    else {
      ElMessage.warning('必填项未填完')
    }
  })
}
const deleteOne = async (row) => {
  await axios.delete(`/lizardcd/db/tokens/${row.id}`)
  getList(current.value)
}
const handleSizeChange = async (size) => {
  pageSize.value = size
  await getList(current.value)
}
const generateToken = async (f) => {
  await f.validate(async (valid) => {
    if(valid) {
      let token = await axios.post(`/lizardcd/auth/token/generate`, {
        user_id: form.value.user_id.id,
        expire_seconds: moment(form.value.expire_at).diff(moment(), 'seconds')
      })
      form.value.jwt_token = token
    }
    else {
      ElMessage.warning('必填项未填完')
    }
  })
}
</script>
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
        <el-button class="pull-right" size="large" type="primary" @click="show=true;form={}">新建租户</el-button>
      </el-col>
    </el-row>
    <el-table 
      :data="list"
      class="line-height40" 
      style="width:100%;margin-top:10px">
      <el-table-column prop="tenant_name" label="租户名" min-width="150" />
      <el-table-column prop="namespace" label="命名空间" min-width="350">
        <template #default="scope">
          <el-tag v-for="item in scope.row.namespaces.map(x => x.namespace)" size="large" type="warning">{{ item }}</el-tag>
        </template>
      </el-table-column>
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
    <h4 v-if="edit===false">新建租户</h4>
    <h4 v-if="edit===true">编辑租户</h4>
  </template>
  <template #default>
    <el-form ref="tenant" :model="form" :rules="rules" label-width="100px">
      <el-form-item label="租户名" prop="tenant_name">
        <el-input v-model="form.tenant_name" size="large" :disabled="edit" />
      </el-form-item>
      <el-form-item label="设置命名空间" prop="namespaces">
        <el-select 
          v-model="form.namespaces" 
          placeholder="请选择命名空间" 
          clearable
          filterable
          multiple
          value-key="namespace"
          style="width:100%" 
          size="large">
          <el-option v-for="(item,i) in namespaces" :key="i" :label="item.namespace" :value="item">
            <span style="float:left">{{item.namespace}}</span>
            <span style="float:right;color:var(--el-text-color-secondary)">{{item.cluster}}</span>
          </el-option>
        </el-select>
      </el-form-item>
    </el-form>
  </template>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show=false" size="large">取消</el-button>
      <el-button type="primary" @click="confirmClick(tenant)" size="large">提交</el-button>
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
import _ from 'lodash'
/* 变量定义 */
const list = ref([])
const pageSize = ref(10)
const pageTotal = ref(0)
const current = ref(1)
const searchKey = ref("")
const show = ref(false)
const form = ref({})
const rules = reactive({
  tenant_name: [{required: true, message: '请填写租户名'}],
})
const edit = ref(false)
const tenant = ref(null)
const namespaces = ref([])
/* 生命周期函数 */
onBeforeMount(async () => {
  getList(1)
  getNamespaces()
})
/* methods */
const getList = async (page) => {
  let url = `page=${page}&size=${pageSize.value}`
  if(searchKey.value != "") url += `&search=tenant_name==${searchKey.value}`
  let response = await axios.get(`/lizardcd/db/tenant?${url}`)
  list.value = response.results.map(x => {
    try {
      x.namespaces = JSON.parse(x.namespaces)
    }
    catch(e) {
      x.namespaces = []
    }
    return x
  })
  pageTotal.value = response.total
}
const getNamespaces = async () => {
  let response = await axios.get(`/lizardcd/server/clusters?with_vm=true`)
  for(let [k,v] of Object.entries(response)) {
    for(let ns of v) {
      namespaces.value.push({
        namespace: ns,
        cluster: k,
      })
    }
  }
  namespaces.value = _.uniqBy(namespaces.value, "namespace")
}
const confirmClick = async (f) => {
  if(!f) return
  await f.validate(async (valid) => {
    if(valid) {
      let params = Object.assign({}, form.value)
      params.namespaces = JSON.stringify(params.namespaces)
      params.update_at = moment()
      if(edit.value === false) {
        await axios.post(`/lizardcd/db/tenant`, {body:params})
        getList(1)
        current.value = 1
        show.value = false
      }
      else {
        let id = params.id
        delete params.id
        await axios.put(`/lizardcd/db/tenant/${id}`, {body:params})
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
  edit.value = true
  show.value = true
}
const deleteOne = async (row) => {
  await axios.delete(`/lizardcd/db/tenant/${row.id}`)
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
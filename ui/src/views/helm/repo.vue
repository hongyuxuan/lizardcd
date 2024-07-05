<template>
<div class="box box-item">
  <div class="box-body" style="padding-top:20px;padding-bottom:0">
    <el-row>
      <el-col :span="12">
        <el-button-group>
          <el-button :icon="Refresh" size="large" style="margin-right:5px" @click="getList" />
          <el-input v-model="searchKey" placeholder="输入名称进行搜索" size="large" :prefix-icon="Search" @change="getPage(1)" clearable style="width:200px;" />
        </el-button-group>
      </el-col>
      <el-col :span="12">
        <el-button-group class="pull-right">
          <el-button class="pull-right" size="large" type="primary" @click="show.update=true">更新仓库</el-button>
          <el-button class="pull-right" size="large" type="primary" @click="show.add=true;edit=false;form={tenant:tenant}" style="margin-right:5px">添加仓库</el-button>
        </el-button-group>
      </el-col>
    </el-row>
    <el-table 
      :data="list" 
      v-loading="loading.table"
      element-loading-text="奋力加载中..."
      class="line-height40" 
      style="width:100%;margin-top:10px;min-height:150px">
      <el-table-column label="" width="45">
        <span class="iconmoon icon-helm svg-inline--fa" style="font-size: 25px;vertical-align: middle;"></span>
      </el-table-column>
      <el-table-column prop="name" label="名称" min-width="150">
        <template #default="scope">
          <el-link :href="`/helm/${scope.row.name}`" :underline="false">{{ scope.row.name }}</el-link>
        </template>
      </el-table-column>
      <el-table-column prop="url" label="URL" min-width="350" />
      <el-table-column prop="tenant" label="所属租户" min-width="80" />
      <el-table-column label="操作" width="100">
        <template #default="scope">
          <el-button :icon="EditPen" circle @click="editOne(scope.row)" />
          <el-popconfirm title="确定删除?" confirm-button-text="确认" cancel-button-text="取消" @confirm="deleteOne(scope.row)">
            <template #reference>
              <el-button icon="Close" circle />
            </template>
          </el-popconfirm>
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
      @current-change="getPage"
      v-model:current-page="current" />
  </div>
</div>
<el-drawer v-model="show.add" direction="rtl" size="650px">
  <template #header>
    <h4 v-if="edit===false">添加仓库</h4>
    <h4 v-if="edit===true">编辑仓库</h4>
  </template>
  <template #default>
    <el-form ref="repo" :model="form" :rules="rules" label-width="100px">
      <el-form-item label="仓库名称" prop="name" >
        <el-input v-model="form.name" size="large" />
      </el-form-item>
      <el-form-item label="仓库URL" prop="url" >
        <el-input v-model="form.url" size="large" />
      </el-form-item>
      <el-form-item label="仓库用户名" prop="app_name" >
        <el-input v-model="form.username" size="large" clearable />
      </el-form-item>
      <el-form-item label="仓库密码" prop="app_name" >
        <el-input v-model="form.password" size="large" type="password" clearable />
      </el-form-item>
      <el-form-item label="所属租户" prop="tenant">
        <el-input v-model="form.tenant" disabled size="large" />
      </el-form-item>
    </el-form>
  </template>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show.add=false">取消</el-button>
      <el-button type="primary" @click="confirmClick(repo)">提交</el-button>
    </div>
  </template>
</el-drawer>
<el-drawer v-model="show.update" direction="rtl" size="650px">
  <template #header>
    <h4>更新仓库</h4>
  </template>
  <el-table 
    :data="namespaces" 
    class="line-height40" 
    style="width:100%;margin-top:10px;min-height:150px">
    <el-table-column prop="cluster" label="集群" min-width="150" />
    <el-table-column prop="namespace" label="命名空间" min-width="150" />
    <el-table-column label="操作" width="80">
      <template #default="scope">
        <el-link type="primary" :underline="false" :disabled="scope.row.loading" @click="updateRepo(scope.row)">更新仓库</el-link>
      </template>
    </el-table-column>
  </el-table>
</el-drawer>
</template>

<script setup>
import { Refresh, Search, EditPen } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { onBeforeMount, ref, reactive } from 'vue'
import { axios } from '/src/assets/util/axios'
/* 变量定义 */
const tenant = localStorage.tenant
const all = ref([])
const searchKey = ref("")
const list = ref([])
const pageSize = ref(10)
const pageTotal = ref(0)
const current = ref(1)
const loading = ref({
  table: false,
  update: false
})
const show = ref({
  add: false,
  update: false,
})
const edit = ref(false)
const form = ref({})
const repo = ref(null)
const namespaces = ref([])
const rules = reactive({
  name: [{required: true, message: '请填写仓库名称'}],
  url: [{required: true, message: '请填写仓库URL'}],
})
/* 生命周期函数 */
onBeforeMount(async () => {
  getList()
  getNamespaces()
})
/* methods */
const getList = async () => {
  loading.value.table = true
  all.value = await axios.get(`/lizardcd/helm/repos`)
  loading.value.table = false
  getPage(current.value)
}
const getPage = async (page) => {
  let tmpList = all.value
  if(searchKey.value !== '') {
    tmpList = all.value.filter(n => n.name.includes(searchKey.value))
  }
  pageTotal.value = tmpList.length
  list.value = tmpList.slice((page-1)*pageSize.value, page*pageSize.value)
}
const getNamespaces = async () => {
  let response = await axios.get(`/lizardcd/server/clusters`)
  for(let [k,v] of Object.entries(response)) {
    for(let ns of v) {
      namespaces.value.push({
        namespace: ns,
        cluster: k,
        loading: false
      })
    }
  }
}
const editOne = (row) => {
  form.value = Object.assign({}, row)
  form.value.tenant ||= localStorage.tenant
  edit.value = true
  show.value.add = true
}
const deleteOne = async (row) => {
  await axios.delete(`/lizardcd/helm/repo/${row.name}`)
  getList()
}
const confirmClick = async (f) => {
  if(!f) return
  await f.validate(async (valid) => {
    if(valid) {
      let params = Object.assign({}, form.value)
      if(edit.value === false) {
        await axios.post(`/lizardcd/helm/repo`, params)
        getList()
        current.value = 1
      }
      else {
        let id = params.id
        delete params.id
        await axios.put(`/lizardcd/db/helm_repositories/${id}`, {body:params})
        getList(current.value)
      }
      show.value.add = false
    }
    else {
      ElMessage.warning('必填项未填完')
    }
  })
}
const updateRepo = async (row) => {
  row.loading = true
  ElMessage.success('正在更新，请稍后')
  try {
    await axios.post(`/lizardcd/helm/cluster/${row.cluster}/namespace/${row.namespace}/repo/update`)
  }
  catch(e){}
  finally {
    row.loading = false
  }
}
</script>
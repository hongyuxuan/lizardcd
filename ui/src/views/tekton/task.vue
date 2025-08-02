<template>
<div class="box box-item">
  <div class="box-body" style="padding-top:20px;padding-bottom:0">
    <el-alert type="info" show-icon style="margin-bottom:15px">
      <template #title>
        关于 Task 配置参考：<el-link href="https://tekton.dev/docs/pipelines/tasks/" underline="never" type="primary" target="_blank">Tasks</el-link>
      </template>
    </el-alert>
    <el-row>
      <el-col :span="18">
        <el-button-group>
          <el-button :icon="Refresh" size="large" @click="getList(true)" />
          <el-select v-model="namespace" placeholder="请选择命名空间" clearable filterable @change="getList(true)" style="width:250px;" size="large">
            <el-option v-for="(item) in props.namespaceList" :key="item" :label="item" :value="item" />
          </el-select>
          <el-input v-model="searchKey" placeholder="输入名称进行搜索" size="large" :prefix-icon="Search" @change="getPage(1)" clearable style="width:250px;" />
        </el-button-group>
      </el-col>
      <el-col :span="6">
        <el-button-group class="pull-right">
          <el-button class="pull-right" size="large" type="primary" @click="show.yaml=true;edit=false;form={content:'',variables:[]}">+ 创建Task</el-button>
        </el-button-group>
      </el-col>
    </el-row>
    <el-table 
      :data="list" 
      v-loading="loading"
      element-loading-text="奋力加载中..."
      class="line-height40" 
      style="width:100%;margin-top:10px;min-height:150px">
      <el-table-column type="selection" width="45" />
      <el-table-column prop="metadata.name" label="名称" min-width="150" />
      <el-table-column label="注解" min-width="300">
        <template #default="scope">
          <el-tag v-for="(v,k,i) in scope.row.metadata.annotations||{}" :key="i" size="large" type="warning">{{ k }}={{ v }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="标签" min-width="200">
        <template #default="scope">
          <el-tag v-for="(v,k,i) in scope.row.metadata.labels||{}" :key="i" size="large" type="primary">{{ k }}={{ v }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" width="170">
        <template #default="scope">
          {{ moment(scope.row.metadata.creationTimestamp).format('YYYY-MM-DD HH:mm:ss') }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="100">
        <template #default="scope">
          <el-button icon="EditPen" circle @click="editOne(scope.row)"></el-button>
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
      @size-change="handleSizeChange"
      @current-change="getPage"
      v-model:current-page="current" />
  </div>
</div>
<el-drawer v-model="show.yaml" direction="rtl" size="700px">
  <template #header>
    <h4>{{ edit === true ? '编辑YAML' : '创建Task' }}</h4>
  </template>
  <template #default>
    <el-form ref="task" label-width="100px">
      <el-form-item label="从模板导入" v-if="edit===false">
        <el-select 
          v-model="form.templates" 
          placeholder="请选择模板" 
          value-key="id" 
          clearable 
          size="large"
          style="width:100%" 
          @change="selectTemplate">
          <el-option v-for="item in templateList" :key="item.id" :label="item.name" :value="item">
            <span style="float:left">{{item.name}}</span>
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="配置YAML">
        <el-switch v-model="wrapLine" active-text="自动换行" inactive-text="不自动换行" inline-prompt class="el-switch-hover-right" size="large" />
        <v-ace-editor
          v-model:value="form.content"
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
            maxLines: 5000,
            minLines: 10,
          }" />
      </el-form-item>
      <el-form-item label="模板变量" v-if="edit===false">
        <table class="table table-bordered" style="margin-bottom:0">
          <thead><tr><th>变量名</th><th>默认变量值</th></tr></thead>
          <tbody>
          <tr v-for="(item,index) in form.variables" :key="index" >
            <td style="vertical-align: top"><el-input v-model="item.key" size="large" /></td>
            <td><el-input v-model="item.value" size="large" /><myTips type="info">{{ item.memo }}</myTips></td>
            <td width="80">
              <el-button-group>
                <el-button icon="Plus" circle @click="addVar(index)"></el-button>
                <el-button icon="Close" circle @click="removeVar(index)"></el-button>
              </el-button-group>
            </td>
          </tr>
          </tbody>
        </table>
      </el-form-item>
    </el-form>
  </template>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show.yaml=false">取消</el-button>
      <el-button type="primary" @click="submitYaml" :disabled="namespace.includes('default')&&userInfo.role!=='admin'">提交</el-button>
    </div>
  </template>
</el-drawer>
</template>
<script setup>
import { Search,Refresh } from '@element-plus/icons-vue'
import { onBeforeMount, ref, computed } from 'vue'
import { useStore } from 'vuex'
import { ElMessage } from 'element-plus'
import { axios } from '/src/assets/util/axios'
import moment from 'moment'
import _ from 'lodash'
/* 引入v-ace-editor */
import { VAceEditor } from 'vue3-ace-editor'
import 'ace-builds/src-noconflict/mode-yaml'
import 'ace-builds/src-noconflict/theme-chrome'
import 'ace-builds/src-noconflict/ext-language_tools'
/* 变量定义 */
const props = defineProps({
  defaultTekton: { type: Object }, 
  namespaceList: { type: Array }, 
})
const store = useStore()
const userInfo = computed(() => {
  return store.state.userInfo
})
const all = ref([])
const searchKey = ref("")
const list = ref([])
const pageSize = ref(10)
const pageTotal = ref(0)
const current = ref(1)
const namespace = ref("")
const loading = ref(false)
const show = ref({
  yaml: false
})
const edit = ref(false)
const form = ref({content:'',variables:[],templates:{}})
const templateList = ref([])
const wrapLine = ref(true)
/* 生命周期函数 */
onBeforeMount(async () => {
  getTemlates()
})
/* methods */
const getList = async (ifLoading) => {
  if(ifLoading) loading.value = true
  if(props.defaultTekton.cluster !== "" && namespace.value !== "") {
    let response = await axios.get(`/lizardcd/tekton/cluster/${props.defaultTekton.cluster}/namespace/${namespace.value}/tasks`)
    all.value = _.sortBy(response.results, 'metadata.creationTimestamp').reverse()
    getPage(current.value)
  }
  if(ifLoading) loading.value = false
}
const getPage = async (page) => {
  let tmpList = all.value
  if(searchKey.value !== '') {
    tmpList = all.value.filter(n => n.metadata.name.includes(searchKey.value))
  }
  pageTotal.value = tmpList.length
  list.value = tmpList.slice((page-1)*pageSize.value, page*pageSize.value)
}
const getTemlates = async () => {
  let response = await axios.get(`/lizardcd/db/yaml_template?filter=type==tekton_task&page=1&size=100&sort=update_at desc`)
  templateList.value = response.results.map(x => {
    x.variables = JSON.parse(x.variables)
    return x
  })
}
const selectTemplate = (val) => {
  if(val) {
    form.value.content = val.content
    form.value.variables = val.variables
  }
  else {
    form.value.content = ""
  }
}
const editOne = async (row) => {
  show.value.yaml = true
  edit.value = true
  form.value.content = await axios.get(`/lizardcd/tekton/cluster/${props.defaultTekton.cluster}/namespace/${namespace.value}/tasks/${row.metadata.name}/yaml`)
}
const deleteOne = async (row) => {
  await axios.delete(`/lizardcd/tekton/cluster/${props.defaultTekton.cluster}/namespace/${namespace.value}/tasks/${row.metadata.name}`)
  ElMessage.success({message: `删除Task成功`})
  setTimeout(async () => {
    await getList(true)
  }, 500)
}
const submitYaml = async () => {
  if(props.defaultTekton.cluster === "" || namespace.value == "") {
    ElMessage.warning({message: '请指定集群和命名空间'})
    return
  } 
  let params = Object.assign({}, form.value)
  delete params.templates
  let vars = {
    "Namespace": namespace.value,
    "Username": localStorage.username
  }
  for(let x of params.variables) {
    vars[x.key] = x.value
  }
  params.variables = vars
  try {
    if(edit.value === false) {
      await axios.post(`/lizardcd/tekton/cluster/${props.defaultTekton.cluster}/namespace/${namespace.value}/apply?kind=Task`, params)
      ElMessage.success({message: `创建Task成功`})
    }
    else {
      delete params.variables
      await axios.post(`/lizardcd/tekton/cluster/${props.defaultTekton.cluster}/namespace/${namespace.value}/apply?kind=Task`, params)
      ElMessage.success({message: `更新Task成功`})
    }
    show.value.yaml = false
    setTimeout(async () => {
      await getList(true)
    }, 500)
  } catch(e) {
    ElMessage.error({message: e})
  }
}
const handleSizeChange = async (size) => {
  pageSize.value = size
  await getList(true)
}
const addVar = (index) => {
  form.value.variables.splice(index+1, 0, {key:"", value:""})
}
const removeVar = (index) => {
  form.value.variables.splice(index, 1)
}
const setNamespace = (ns) => {
  namespace.value = ns
  getList(true)
}
defineExpose({
  setNamespace
})
</script>
<template>
<el-breadcrumb :separator-icon="ArrowRight">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>工作负载</el-breadcrumb-item>
  <el-breadcrumb-item>部署</el-breadcrumb-item>
</el-breadcrumb>
<el-card>
  <template #header>
    <div class="card-header">
      <span class="card-header-text">部署</span>
    </div>
  </template>
  <el-row>
    <el-col :span="12">
      <el-button-group>
        <el-button :icon="Refresh" size="large" style="margin-right:5px" @click="getList(current)" />
        <el-select v-model="cluster" placeholder="请选择集群" clearable filterable style="width:200px;margin-right:8px" size="large">
          <el-option v-for="(v,k) in clusterList" :key="k" :label="k" :value="k" />
        </el-select>
        <el-select v-model="namespace" placeholder="请选择命名空间" clearable filterable @change="getList(true)" style="width:200px;margin-right:8px" size="large">
          <el-option v-for="(item) in clusterList[cluster]" :key="item" :label="item" :value="item" />
        </el-select>
        <el-input v-model="searchKey" placeholder="输入名称进行搜索" size="large" :prefix-icon="Search" @change="getPage(1)" clearable style="width:200px;" />
      </el-button-group>
    </el-col>
    <el-col :span="12">
      <el-dropdown @command="handleMore" class="pull-right">
        <el-button size="large">更多操作<el-icon class="el-icon--right"><arrow-down /></el-icon></el-button>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item :command="{action:'restartBatch'}">批量重启</el-dropdown-item>
            <el-dropdown-item :command="{action:'deleteBatch'}">批量删除</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
      <el-button class="pull-right" size="large" type="primary" @click="show.new=true;formNew={content:'',variables:[],cluster:cluster,namespace:namespace}" style="margin-right:5px">新建工作负载</el-button>
    </el-col>
  </el-row>
  <el-table 
    :data="list" 
    v-loading="loading"
    element-loading-text="奋力加载中..."
    class="line-height40" 
    @selection-change="select"
    style="width:100%;margin-top:10px;min-height:150px">
    <el-table-column type="selection" width="45" />
    <el-table-column label="" width="45">
      <font-awesome-icon icon="layer-group" style="font-size:25px;vertical-align:middle;" />
    </el-table-column>
    <el-table-column prop="name" label="名称" min-width="200">
      <template #default="scope">
        <el-link :underline="false" :href="`/workload/deployments/${scope.row.name}?cluster=${cluster}&namespace=${namespace}`">{{ scope.row.name }}</el-link>
      </template>
    </el-table-column>
    <el-table-column label="状态" min-width="200">
      <template #default="scope">
        <font-awesome-icon icon="circle" v-if="scope.row.replicas==0" class="runningstatus text-gray" />
        <font-awesome-icon icon="circle" v-else-if="scope.row.available=='False'" class="runningstatus twinkling text-yellow" />
        <font-awesome-icon icon="circle" v-else-if="scope.row.available=='True'" class="runningstatus text-green" />
        <span v-if="scope.row.replicas === 0">停止 ( {{ scope.row.readyReplicas }} / {{ scope.row.replicas }} )</span>
        <span v-else>{{ scope.row.available === 'True' ? '运行中' : '更新中' }} ( {{ scope.row.readyReplicas }} / {{ scope.row.replicas }} )</span>
      </template>
    </el-table-column>
    <el-table-column label="更新时间">
      <template #default="scope">
        {{ moment(scope.row.lastUpdateTime).format('YYYY-MM-DD HH:mm:ss') }}
      </template>
    </el-table-column>
    <el-table-column label="操作" width="120">
      <template #default="scope">
        <el-dropdown @command="handleCommand" style="vertical-align:middle;">
          <el-button>更多操作<el-icon class="el-icon--right"><arrow-down /></el-icon></el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item :command="{action:'restart',row:scope.row}">重启</el-dropdown-item>
              <el-dropdown-item :command="{action:'setImage',row:scope.row}">设置镜像</el-dropdown-item>
              <el-dropdown-item :command="{action:'scale',row:scope.row}">设置副本</el-dropdown-item>
              <el-dropdown-item :command="{action:'yaml',row:scope.row}">编辑YAML</el-dropdown-item>
              <el-dropdown-item :command="{action:'delete',row:scope.row}">删除</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
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
</el-card>
<el-drawer v-model="show.new" direction="rtl" size="800px">
  <template #header>
    <h4>新建工作负载</h4>
  </template>
  <template #default>
    <el-form ref="refNew" :rules="rules" :model="formNew" label-width="120px">
      <el-form-item label="集群" prop="cluster">
        <el-input v-model="formNew.cluster" disabled size="large" />
      </el-form-item>
      <el-form-item label="命名空间" prop="namespace">
        <el-input v-model="formNew.namespace" disabled size="large" />
      </el-form-item>
      <el-form-item label="从模板导入" prop="templates">
        <el-select 
          v-model="formNew.templates" 
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
      <el-form-item label="配置YAML" prop="content">
        <v-ace-editor
          v-model:value="formNew.content"
          lang="yaml"
          theme="chrome"
          style="width:100%;"
          :options="{
            enableBasicAutocompletion: true,
            enableSnippets: true,
            enableLiveAutocompletion: true,
            tabSize: 2,
            showPrintMargin: false,
            fontSize: 14,
            minLines: 10,
            maxLines: 5000
          }" />
      </el-form-item>
      <el-form-item label="模板变量">
        <table class="table table-bordered">
          <thead><tr><th>变量名</th><th>默认变量值</th></tr></thead>
          <tbody>
          <tr v-for="(item,index) in formNew.variables" :key="index" >
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
    </el-form>
  </template>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show.new=false">取消</el-button>
      <el-button type="primary" @click="submitNew(refNew)">提交</el-button>
    </div>
  </template>
</el-drawer>
<el-drawer v-model="show.yaml" direction="rtl" size="800px">
  <template #header>
    <h4>编辑YAML</h4>
  </template>
  <template #default>
    <v-ace-editor
      v-model:value="yamlContent"
      lang="yaml"
      theme="chrome"
      style="width:100%"
      :options="{
        enableBasicAutocompletion: true,
        enableSnippets: true,
        enableLiveAutocompletion: true,
        tabSize: 2,
        showPrintMargin: false,
        fontSize: 14,
        minLines: 10,
        maxLines: 5000,
        wrap: true
      }" />
  </template>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show.yaml=false" size="large">取消</el-button>
      <el-button type="primary" @click="submitYaml" size="large">提交</el-button>
    </div>
  </template>
</el-drawer>
</template>

<script setup>
import { ArrowRight,Search,Refresh } from '@element-plus/icons-vue'
import { onBeforeMount, onBeforeUnmount, ref, reactive } from 'vue'
import { useRoute } from 'vue-router'
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
const route = useRoute()
const all = ref([])
const searchKey = ref("")
const list = ref([])
const pageSize = ref(10)
const pageTotal = ref(0)
const current = ref(1)
const cluster = ref("")
const clusterList = ref({})
const namespace = ref("")
const show = ref({
  pods: false,
  event: false,
  new: false,
  yaml: false
})
const rules = reactive({
  cluster: [{required: true, message: '请选择集群'}],
  namespace: [{required: true, message: '请选择命名空间'}],
})
const loading = ref(false)
const selected = ref([])
/* 新建发布 */
const formNew = ref({content:'',variables:[]})
const templateList = ref([])
const refNew = ref(null)
const timer = ref(null)
/* YAML配置 */
const yamlContent = ref("")
/* 生命周期函数 */
onBeforeMount(async () => {
  await getClusterList()
  timer.value = setInterval(() => {
    getList(false)
  }, 15000)
  getTemlates()
  if(route.query.cluster && route.query.namespace) {
    cluster.value = route.query.cluster
    namespace.value = route.query.namespace
    getList(true)
  }
})
onBeforeUnmount(() => {
  if(timer) {
    clearInterval(timer.value)
    timer.value = null
  }
})
/* methods */
const getClusterList = async () => {
  clusterList.value = await axios.get(`/lizardcd/server/clusters`)
}
const getList = async (ifLoading) => {
  if(ifLoading) loading.value = true
  if(cluster.value !== "" && namespace.value !== "") {
    let response = await axios.get(`/lizardcd/kubernetes/cluster/${cluster.value}/namespace/${namespace.value}/deployments`)
    response = response.map(x => {
      let progress = x.status.conditions.find(n => n.type === 'Progressing')
      let r = {
        name: x.metadata.name,
        replicas: x.status.replicas||0,
        readyReplicas: x.status.readyReplicas||0,
        unavailableReplicas: x.status.unavailableReplicas||0,
        lastUpdateTime: progress?.lastUpdateTime
      }
      r.available = r.readyReplicas >= r.replicas ? 'True' : 'False'
      return r
    })
    all.value = _.sortBy(response, 'lastUpdateTime').reverse()
    getPage(current.value)
  }
  if(ifLoading) loading.value = false
}
const getPage = async (page) => {
  let tmpList = all.value
  if(searchKey.value !== '') {
    tmpList = all.value.filter(n => n.name.includes(searchKey.value))
  }
  pageTotal.value = tmpList.length
  list.value = tmpList.slice((page-1)*pageSize.value, page*pageSize.value)
}
const handleCommand = async (command) => {
  switch(command.action) {
    case "restart": {
      await ElMessageBox.confirm('确定重启？','警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(async () => {
        await axios.patch(`/lizardcd/kubernetes/cluster/${cluster.value}/namespace/${namespace.value}/deployments/${command.row.name}/rollout`)
        ElMessage.success({message: '操作成功'})
      }).catch(() =>{})
      break
    }
    case "setImage": {
      let pods = await axios.get(`/lizardcd/kubernetes/cluster/${cluster.value}/namespace/${namespace.value}/deployments/${command.row.name}/pods`)
      let container = pods[0].spec.containers[0].name
      await ElMessageBox.prompt(`请输入容器 ${container} 的镜像`,'提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
      }).then(async ({value}) => {
        await axios.patch(`/lizardcd/kubernetes/cluster/${cluster.value}/namespace/${namespace.value}/deployments/${command.row.name}?container=${container}&image=${value}`)
        ElMessage.success({message: '操作成功'})
      }).catch(() =>{})
      break
    }
    case "scale": {
      await ElMessageBox.prompt(`请输入副本数`,'提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        inputPattern: /\d+/,
        inputErrorMessage: '必须输入数字'
      }).then(async ({value}) => {
        await axios.patch(`/lizardcd/kubernetes/cluster/${cluster.value}/namespace/${namespace.value}/deployments/scale`, {
          workloads: [
            {
              name: command.row.name,
              replicas: parseInt(value)
            }  
          ]
        })
        ElMessage.success({message: '操作成功'})
      }).catch(() =>{})
      break
    }
    case "yaml": {
      yamlContent.value = await axios.get(`/lizardcd/kubernetes/cluster/${cluster.value}/namespace/${namespace.value}/deployments/${command.row.name}/yaml`)
      show.value.yaml = true
      break
    }
    case "delete": {
      await ElMessageBox.confirm('确定删除？','警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(async () => {
        await axios.delete(`/lizardcd/kubernetes/cluster/${cluster.value}/namespace/${namespace.value}/deployments/${command.row.name}`)
        ElMessage.success({message: '删除成功'})
      }).catch(() =>{})
      break
    }
  }
  setTimeout(async () => {
    await getList(true)
  }, 2000)
}
const getTemlates = async () => {
  let response = await axios.get(`/lizardcd/db/application_template?page=1&size=100&sort=update_at desc`)
  templateList.value = response.results.map(x => {
    x.variables = JSON.parse(x.variables)
    return x
  })
}
const selectTemplate = (val) => {
  if(val) {
    formNew.value.content = val.content
    formNew.value.variables = val.variables
  }
  else {
    formNew.value.content = ""
  }
}
const submitNew = async (f) => {
  if(!f) return
  await f.validate(async (valid) => {
    if(valid) {
      let params = Object.assign({}, formNew.value)
      delete params.cluster
      delete params.namespace
      delete params.templates
      let vars = {}
      for(let x of params.variables) {
        vars[x.key] = x.value
      }
      params.variables = vars
      await axios.patch(`/lizardcd/kubernetes/cluster/${formNew.value.cluster}/namespace/${formNew.value.namespace}/apply/variable?kind=Deployment`, params)
      ElMessage.success({message: '发布成功'})
      show.value.new = false
    }
    else {
      ElMessage.warning('必填项未填完')
    }
  })
}
const submitYaml = async () => {
  await axios.patch(`/lizardcd/kubernetes/cluster/${cluster.value}/namespace/${namespace.value}/apply/yaml`, yamlContent.value, {
    headers: {
      'Content-Type': 'text/plain'
    }
  })
  show.value.yaml = false
}
const addVar = (index) => {
  formNew.value.variables.splice(index+1, 0, {key:"", value:""})
}
const removeVar = (index) => {
  formNew.value.variables.splice(index, 1)
}
const handleMore = async (command) => {
  if(selected.value.length === 0) {
    ElMessage.warning({message: '请勾选工作负载'})
    return
  }
  switch(command.action) {
    case "restartBatch": {
      await ElMessageBox.confirm('确定重启？','警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(async () => {
        for(let x of selected.value) {
          await axios.patch(`/lizardcd/kubernetes/cluster/${cluster.value}/namespace/${namespace.value}/deployments/${x.name}/rollout`)
        }
        ElMessage.success({message: '操作成功'})
      }).catch(() =>{})
      break
    }
    case "deleteBatch": {
      await ElMessageBox.confirm('确定删除？','警告', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }).then(async () => {
        for(let x of selected.value) {
          await axios.delete(`/lizardcd/kubernetes/cluster/${cluster.value}/namespace/${namespace.value}/deployments/${x.name}`)
        }
        ElMessage.success({message: '删除成功'})
      }).catch(() =>{})
      break
    }
  }
}
const select = (val) => {
  selected.value = val
}
const handleSizeChange = async (size) => {
  pageSize.value = size
  await getList(true)
}
</script>
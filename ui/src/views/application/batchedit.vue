<template>
<el-drawer v-model="show" direction="rtl" size="700px">
  <template #header>
    <h4>批量更新</h4>
  </template>
  <template #default>
    <el-form ref="refForm" :model="form" label-width="135px">
      <el-form-item>
        <template #label><el-checkbox v-model="disabled.git_http_url">代码仓库</el-checkbox></template>
        <el-input v-model="form.git_http_url" size="large" clearable :disabled="!disabled.git_http_url" />
        <myTips type="info">仅支持http(s)开头</myTips>
      </el-form-item>
      <el-form-item>
        <template #label><el-checkbox v-model="disabled.repo">镜像仓库/制品库</el-checkbox></template>
        <el-select 
          v-model="form.repo" 
          placeholder="请选择" 
          value-key="id" 
          clearable 
          :disabled="!disabled.repo"
          size="large" 
          style="width:100%">
          <el-option v-for="item in repoList" :key="item.id" :label="item.repo_url" :value="item">
            <span style="float:left">{{item.repo_url}}</span>
            <span style="float:right;color:var(--el-text-color-secondary)">{{item.repo_account}}</span>
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item>
        <template #label><el-checkbox v-model="disabled.repo_name">仓库/项目</el-checkbox></template>
        <el-input v-model="form.repo_name" size="large" clearable :disabled="!disabled.repo_name" />
        <myTips type="info">Artifactory填写仓库名<br>Harbor填写项目名<br>DockerHub填写namespace<br>S3填写桶名</myTips>
      </el-form-item>
      <el-form-item>
        <template #label><el-checkbox v-model="disabled.timeout">超时时间</el-checkbox></template>
        <el-input v-model="form.timeout" type="number" size="large" clearable :disabled="!disabled.timeout" />
      </el-form-item>
      <el-form-item label="所属租户" prop="tenant">
        <template #label><el-checkbox v-model="disabled.tenant">所属租户</el-checkbox></template>
        <el-select 
          v-model="form.tenant" 
          placeholder="请选择租户" 
          clearable 
          filterable 
          :disabled="userInfo.role!=='admin'||!disabled.tenant"
          style="width:100%;" 
          size="large">
          <el-option v-for="(item,index) in tenantList" :key="index" :label="item" :value="item" />
        </el-select>
      </el-form-item>
    </el-form>
  </template>
  <template #footer>
    <div style="flex: auto">
      <el-button @click="show=false">取消</el-button>
      <el-button type="primary" @click="confirmClick(refForm)" :loading="loading">提交</el-button>
    </div>
  </template>
</el-drawer>
</template>

<script setup>
import { onBeforeMount, ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { useStore } from 'vuex'
import { axios } from '/src/assets/util/axios.js'
import moment from 'moment'
/* 引入v-ace-editor */
import 'ace-builds/src-noconflict/mode-yaml'
import 'ace-builds/src-noconflict/theme-chrome'
import 'ace-builds/src-noconflict/ext-language_tools'
/* 变量定义 */
const props = defineProps({
  clusterList: { type: Object }, 
  tenantList: { type: Array },
  repoList: { type: Array },
  rows: { type: Array },
})
const store = useStore()
const userInfo = computed(() => {
  return store.state.userInfo
})
const show = ref(false)
const disabled = ref({
  git_http_url: false,
  tenant: false,
  repo: false,
  repo_name: false,
  timeout: false,
})
const form = ref({
  timeout: 300,
})
const refForm = ref(null)
const loading = ref(false)
/* methods */
const confirmClick = async (f) => {
  if(!f) return
  await f.validate(async (valid) => {
    if(valid) {
      let params = {
        update_at: moment()
      }
      if(params.repo) {
        params.repo_id = params.repo.id
        delete params.repo
      }
      for(let [k,v] of Object.entries(disabled.value)) {
        if(v === true) {
          params[k] = form.value[k]
        }
      }
      await Promise.all(props.rows.map(x => {
        let body = Object.assign(x, params)
        body.timeout = parseInt(body.timeout)
        return axios.put(`/lizardcd/db/application/${x.id}`, {body})
      }))
      show.value = false
      loading.value = false
      emit('submit')
    }
  })
}
const open = () => {
  show.value = true
}
defineExpose({ open })
const emit = defineEmits(['submit'])
</script>
<template>
<el-breadcrumb :separator-icon="ArrowRight">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
  <el-breadcrumb-item>持续集成</el-breadcrumb-item>
  <el-breadcrumb-item :to="{path:'/ci/tekton'}">流水线</el-breadcrumb-item>
  <el-breadcrumb-item>{{ route.params.name }}</el-breadcrumb-item>
</el-breadcrumb>
<result v-if="tektonSource === 'crd'" />
<resultflow v-if="tektonSource === 'cloudevent' && pipelinerunShow === 'flow'" />
<resultcollapse v-if="tektonSource === 'cloudevent' && pipelinerunShow === 'collapse'" />
<el-backtop :right="70" :bottom="50" />
</template>

<script setup>
import { ArrowRight} from '@element-plus/icons-vue'
import { ref, onBeforeMount } from 'vue'
import { useRoute } from 'vue-router'
import result from './result.vue'
import resultflow from './result_flow.vue'
import resultcollapse from './result_collapse.vue'
import { axios } from '/src/assets/util/axios'
/* 变量定义 */
const route = useRoute()
const tektonSource = ref('')
const pipelinerunShow = ref('')
/* 生命周期函数 */
onBeforeMount(async () => {
  let response = await axios.get(`/lizardcd/db/settings?filter=setting_key==tekton_source`)
  if(response.total > 0)
    tektonSource.value = response.results[0].setting_value

  response = await axios.get(`/lizardcd/db/settings?filter=setting_key==pipelinerun_result`)
  if(response.total > 0)
  pipelinerunShow.value = response.results[0].setting_value
})
</script>
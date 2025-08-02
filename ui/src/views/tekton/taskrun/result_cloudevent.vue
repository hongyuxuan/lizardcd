<template>
<el-card>
  <template #header>
    <div class="card-header" style="display:block;">
      <div>
        <b class="card-header-text" style="margin-right:20px">{{ route.params.name }}</b>
        <span style="color:#a0a0a0;font-size:15px">
          任务启动于 {{moment.duration(moment(taskRunInfo.startTime)-moment()).humanize(true)}}
        </span>
      </div>
      <div style="margin-top:15px" v-if="taskRunInfo.conditions">
        <el-text :type="getColor(taskRunInfo.conditions[0]?.reason)">{{ taskRunInfo.conditions[0]?.reason }}</el-text>
        <el-text style="margin-left:20px">{{ taskRunInfo.conditions[0]?.message }}</el-text>
      </div>
    </div>
  </template>
  <inner ref="refInner" />
</el-card>
<el-backtop :right="70" :bottom="50" />
</template>

<script setup>
import axios from 'axios'
import { onBeforeMount, onBeforeUnmount, ref } from 'vue'
import { useRoute } from 'vue-router'
import { getDuration, getColor } from '@/assets/util/common'
import moment from 'moment'
import _ from 'lodash'
import inner from './result_inner.vue'
/* 变量定义 */
const taskRunInfo = ref({})
const route = useRoute()
const timer = ref(null)
const refInner = ref(null)
/* 生命周期函数 */
onBeforeMount(async () => {
  await getTaskRun()
  timer.value = setInterval(async () => {
    await getTaskRun()
  }, 5000)
})
onBeforeUnmount(() => {
  if(timer) {
    clearInterval(timer.value)
    timer.value = null
  }
})
/* methods */
const getTaskRun = async () => {
  let response = await axios.get(`/tekton-pipelines/lizardcd/${route.query.namespace}/taskRuns/${route.params.name}`)
  response.expire = getDuration(response.status.endTime, response.status.startTime)
  taskRunInfo.value = response
  refInner.value.refreshTaskRun(response)
  refInner.value.refreshCurrentStep()
  if(taskRunInfo.value.endTime) {
    clearInterval(timer.value)
    timer.value = null
  }
}
</script>
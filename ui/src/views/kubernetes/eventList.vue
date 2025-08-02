<template>
<div class="box box-item">
  <div class="box-body" style="padding-top:20px">
    <el-table :data="eventList" style="width:100%">
      <el-table-column prop="type" label="Type">
        <template #default="scope">
          <el-tag v-if="scope.row.type==='Warning'" type="warning" size="large">{{ scope.row.type }}</el-tag>
          <el-tag v-else-if="scope.row.type==='Normal'" type="primary" size="large">{{ scope.row.type }}</el-tag>
          <el-tag v-else type="info" size="large">{{ scope.row.type }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="reason" label="Reason" />
      <el-table-column prop="age" label="Age" />
      <el-table-column prop="source.component" label="From" />
      <el-table-column prop="message" label="Message" min-width="300px" />
    </el-table>
  </div>
</div>
</template>

<script setup>
import { ref } from 'vue'
import { useRoute } from 'vue-router'
import { axios } from '/src/assets/util/axios'
import moment from 'moment'
/* 变量定义 */
const props = defineProps({
  resourceType: { type: String },
  resourceName: { type: String }, 
})
const route = useRoute()
const eventList = ref([])
/* methods */
const getEvents = async (row) => {
  let response = await axios.get(`/lizardcd/kubernetes/cluster/${route.query.cluster}/namespace/${route.query.namespace}/${props.resourceType}/${props.resourceName}/events`)
  eventList.value = response.map(x => {
    x.age = moment.duration(moment(x.lastTimestamp)-moment()).humanize(true)
    return x
  })
}
defineExpose({ getEvents })
</script>
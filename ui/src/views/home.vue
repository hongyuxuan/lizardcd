<template>
<el-breadcrumb :separator-icon="ArrowRight">
  <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
</el-breadcrumb>
<el-row :gutter="15">
  <el-col :span="12">
    <el-card>
      <template #header>
        <div class="card-header">
          <span class="card-header-text">应用部署月度统计</span>
          <div class="box-tools pull-right">
            <el-date-picker v-model="date1" type="month" @change="getStatsByAppname();getStatsBySuccess()" />
          </div>
        </div>
      </template>
      <div id="task-by-app_name" style="height:350px" />
    </el-card>
  </el-col>
  <el-col :span="12"><el-card>
      <template #header>
        <div class="card-header">
          <span class="card-header-text">应用部署成功率统计</span>
          <div class="box-tools pull-right">
            <el-date-picker v-model="date1" type="month" @change="getStatsByAppname();getStatsBySuccess()" />
          </div>
        </div>
      </template>
      <div id="task-by-success" style="height:350px" />
    </el-card></el-col>
</el-row>
<el-row :gutter="15">
  <el-col :col="24">
    <el-card>
      <template #header>
        <div class="card-header">
          <span class="card-header-text">近一年应用部署月度分布</span>
        </div>
      </template>
      <div id="task-by-month" style="height:350px" />
    </el-card>
  </el-col>
</el-row>
</template>
<script setup>
import { ref, inject, onMounted, onBeforeUnmount } from "vue";
import { ArrowRight,Search,RefreshLeft,EditPen,Delete,Plus } from '@element-plus/icons-vue'
import moment from "moment"
import { axios } from '/src/assets/util/axios.js'
let echarts = inject("ec")
const date1 = ref(new Date())
const chart1 = ref({})
const chart2 = ref({})
const chart3 = ref({})
onMounted(async () => {
  getStatsByAppname()
  getStatsBySuccess()
  getStatsByMonth()
})
onBeforeUnmount(() => {
  if(chart1.value) {
    chart1.value.clear()
    chart1.value.dispose()
    chart1.value = null
  }
  if(chart2.value) {
    chart2.value.clear()
    chart2.value.dispose()
    chart2.value = null
  }
  if(chart3.value) {
    chart3.value.clear()
    chart3.value.dispose()
    chart3.value = null
  }
})
/* methods */
const getStatsByAppname = async () => {
  let startOfmonth = moment(date1.value).startOf('month')
  let response = await axios.get(`/lizardcd/task/stats?by=app_name&time_from=${startOfmonth.format('YYYY-MM-DD')}&time_till=${startOfmonth.add(1,'month').format('YYYY-MM-DD')}`)
  chart1.value = echarts.init(document.getElementById("task-by-app_name"))
  chart1.value.setOption({
    title: {
      show: false,
    },
    series: [
      {
        name: 'stats',
        type: 'pie',
        radius: ['30%', '60%'],
        left: 'center',
        width: '85%',
        emphasis: {
          label: {
            show: true,
          }
        },
        label: {
          alignTo: 'edge',
          formatter: '{b}\n{c}',
          minMargin: 5,
          edgeDistance: 10,
          lineHeight: 15,
        },
        labelLine: {
          length: 15,
          length2: 0,
          maxSurfaceAngle: 80
        },
        data: response.map(x => {
          return { value: x.count, name: x.app_name}
        })
      }
    ]
  })
}
const getStatsBySuccess = async () => {
  let startOfmonth = moment(date1.value).startOf('month')
  let response = await axios.get(`/lizardcd/task/stats?by=success&time_from=${startOfmonth.format('YYYY-MM-DD')}&time_till=${startOfmonth.add(1,'month').format('YYYY-MM-DD')}`)
  chart2.value = echarts.init(document.getElementById("task-by-success"))
  chart2.value.setOption({
    title: {
      show: false,
    },
    series: [
      {
        name: 'stats',
        type: 'pie',
        radius: ['30%', '60%'],
        left: 'center',
        width: '85%',
        emphasis: {
          label: {
            show: true,
          }
        },
        label: {
          alignTo: 'edge',
          formatter: '{b}\n{c}',
          minMargin: 5,
          edgeDistance: 10,
          lineHeight: 15,
        },
        labelLine: {
          length: 15,
          length2: 0,
          maxSurfaceAngle: 80
        },
        data: response.map(x => {
          return { value: x.count, name: x.success === 1 ? '成功' : '失败'}
        })
      }
    ]
  })
}
const getStatsByMonth = async () => {
  let startMonth = moment().add(1, 'month').subtract(1, 'year')
  let endMonth = moment().add(1, 'month')
  let response = await axios.get(`/lizardcd/task/stats?by=month&time_from=${startMonth.format('YYYY-MM-DD')}&time_till=${endMonth.format('YYYY-MM-DD')}`)
  chart3.value = echarts.init(document.getElementById("task-by-month"))
  chart3.value.setOption({
    title: {
      show: false,
    },
    xAxis: {
      type: 'category',
      data: response.map(x => x.month),
    },
    yAxis: {
      type: 'value'
    },
    series: [
      {
        name: '任务数量',
        data: response.map(x => x.count),
        type: 'bar',
        label: {
          show: true,
          position: 'outside'
        },
        showBackground: true,
        backgroundStyle: {
          color: 'rgba(180, 180, 180, 0.2)'
        }
      }
    ]
  })
}
</script>
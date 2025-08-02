<template>
  <el-breadcrumb :separator-icon="ArrowRight">
    <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
    <el-breadcrumb-item>个人设置</el-breadcrumb-item>
  </el-breadcrumb>
  <el-card>
    <template #header>
      <div class="card-header">
        <span class="card-header-text">个人设置</span>
      </div>
    </template>
    <el-descriptions title="" :column="1" size="large" border>
      <el-descriptions-item label="头像" label-width="150">
        <div style="margin-bottom:10px"><el-avatar :size="150" :src="userInfo.profile.avatar||'/images/avator.png'" /></div>
        <el-link v-if="!edit.avatar" type="primary" underline="never" @click="edit.avatar=true">修改头像</el-link>
        <div v-else>
            <el-input v-model="form.avatar" placeholder="请输入头像链接" size="large" style="width:400px;margin-right:5px" />
            <el-button :icon="Check" circle size="large" @click="updateAvatar()" />
            <el-button :icon="Close" circle size="large" @click="edit.avatar=false" />
        </div>
      </el-descriptions-item>
      <el-descriptions-item label="用户ID" label-width="150">{{ userInfo.userid }}</el-descriptions-item>
      <el-descriptions-item label="用户名" label-width="150">{{ userInfo.username }}</el-descriptions-item>
      <el-descriptions-item label="邮箱" label-width="150">{{ userInfo.email }}</el-descriptions-item>
      <el-descriptions-item label="所属租户" label-width="150">{{ userInfo.tenant }}</el-descriptions-item>
      <el-descriptions-item label="授权命名空间" label-width="150">
        <el-tag v-for="item in namespaces.map(x => x.namespace)" size="large" type="warning" syl>{{ item }}</el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="用户角色" label-width="150">{{ userInfo.role }}</el-descriptions-item>
    </el-descriptions>
  </el-card>
</template>
<script setup>
import { axios } from '/src/assets/util/axios.js'
import { ArrowRight, Check, Close } from '@element-plus/icons-vue'
import { onBeforeMount, computed, ref } from 'vue'
import { useStore } from 'vuex'
/* 变量定义 */
const store = useStore()
const userInfo = computed(() => {
  return store.state.userInfo
})
const namespaces = ref([])
const edit = ref({
  avatar: false
})
const form = ref({})
/* 生命周期函数 */
onBeforeMount(async () => {
  getTenant()
  form.value.avatar = userInfo.value.profile.avatar || ''
})
/* methods */
const getTenant = async () => {
  namespaces.value = []
  for(let x of userInfo.value.tenant.split(',')) {
    let response = await axios.get(`/lizardcd/db/tenant?filter=tenant_name==${x}&page=1&size=1`)
    namespaces.value = namespaces.value.concat(JSON.parse(response.results[0].namespaces))
  }
}
const updateAvatar = async () => {
  userInfo.value.profile.avatar = form.value.avatar
  await axios.put(`/lizardcd/db/user/${userInfo.value.id}`, {body: {profile: JSON.stringify(userInfo.value.profile)}})
}
</script>
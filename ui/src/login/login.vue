<template>
  <el-container>
    <el-aside class="login-aside" width="40%">
      <el-image style="height:180px;margin-top:calc(50vh - 90px)" src="/images/lizardcd-logo.png" />
    </el-aside>
    <el-container>
      <el-main style="padding:0">
        <el-row justify="center" style="flex-direction:column;align-items:center;height:calc(100vh)">
          <div style="color:#141f29;font-size:25px;">WELCOME TO LIZARDCD</div>
          <el-form
            ref="login"
            label-position="top"
            label-width="auto"
            :model="form"
            :rules="rules"
            style="width:400px;margin-top:30px">
            <el-form-item label="Username" prop="username">
              <el-input v-model="form.username" size="large" />
            </el-form-item>
            <el-form-item label="Password" prop="password">
              <el-input v-model="form.password" type="password" show-password size="large" @keyup.enter="submit(login)" />
            </el-form-item>
            <el-form-item>
              <el-button @click="submit(login)" style="width:100%" type="primary" size="large">Login</el-button>
            </el-form-item>
            <el-divider v-if="oauth2"><span style="color:#b4b4b4">SIGN IN WITH OAUTH2</span></el-divider>
          </el-form>
          <el-row justify="center" v-if="oauth2" style="align-items:center;">
            <el-button v-for="(item,i) in oauth2List" :key="i" @click="loginWithOauth2(item)">{{ item.name }}</el-button>
          </el-row>
        </el-row>
      </el-main>
    </el-container>
  </el-container>
</template>
<script setup>
import { onBeforeMount, ref, reactive } from 'vue'
import { useRoute } from 'vue-router'
import { axios } from '/src/assets/util/axios'
import moment from 'moment'
const route = useRoute()
const form = ref({})
const rules = reactive({
  username: [{required: true, message: '请输入用户名'}],
  password: [{required: true, message: '请输入密码'}],
})
const login = ref(null)
const oauth2 = ref(false)
const oauth2List = ref([])
/* 生命周期函数 */
onBeforeMount(async () => {
  await getServerConfig()
  getOauth2()
})
/* methods */
const getServerConfig = async () => {
  let response = await axios.get(`/lizardcd/server/config`)
  oauth2.value = response.Auth.Oauth2
}
const getOauth2 = async () => {
  oauth2List.value = await axios.get(`/lizardcd/server/oauth2`)
}
const submit = async (f) => {
  if(!f) return
  await f.validate(async (valid) => {
    if(valid) {
      let response = await axios.post(`/lizardcd/auth/login`, form.value)
      localStorage.access_token = response.access_token
      if(window.location.search.startsWith('?redirect=')) {
        window.location.href = window.location.search.substring('?redirect='.length, window.location.search.length)
      } else {
        window.location.href = "/"
      }
    } else {
      ElMessage.warning('必填项未填完')
    }
  })
}
const loginWithOauth2 = async (item) => {
  window.location.href = `${item.authorized_url}?client_id=${item.client_id}&response_type=code&redirect_uri=${item.redirect_url}&oauth_timestamp=${moment().valueOf()}`
}
</script>
<style>
.login-aside {
  background-color: #141f29;
  height: calc(100vh);
  text-align:center;
}
.login-panel {
  text-align:center;
}
</style>
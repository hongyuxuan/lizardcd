import axios from 'axios'
import moment from 'moment'
import { ElMessage } from 'element-plus'

axios.defaults.timeout = 60000;

axios.interceptors.request.use((config) => {
  if(config.url != '/lizardcd/auth/login') {
    config.headers.Authorization = `Bearer ${localStorage.access_token}`
  }
	return config
});

axios.interceptors.response.use(
	(response) => {
    if(response.data?.code && response.data?.code !== 200 && response.data?.code !== 0) { // 有code，但code既不等于200也不等于0，则视为错误
      return Promise.reject(response.data?.message)
    }
    if(response.data?.message) {
      ElMessage.success({message: response.data.message})
    }
		if(response.data?.data !== undefined ){
			return response.data.data;
		}
		else {
			return response.data;
		}
	},
	async (err) => {
		if(err.response) {
      if(err.response.status === 401 && !['/lizardcd/auth/login','/lizardcd/auth/chpasswd'].includes(err.response.config.url)) { // 登录失效
        window.location.href = `/login/?redirect=${window.location.href}`
      }
      let err_message = err.response.data ? (err.response.data.message || err.response.data) : err.response.statusText
      ElMessage.error({message: err_message})
    }
		return Promise.reject(err.response ? err.response.data : err)
	}
);

export {axios};

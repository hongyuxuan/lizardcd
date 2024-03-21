import axios from 'axios'
import moment from 'moment'
import { ElMessage } from 'element-plus'

axios.defaults.timeout = 60000;

axios.interceptors.request.use((config) => {
	return config
});

axios.interceptors.response.use(
	(response) => {
    if(response.data.message) {
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
      let err_message = err.response.data ? (err.response.data.message || err.response.data) : err.response.statusText
      ElMessage.error({message: err_message})
    }
		return Promise.reject(err.response ? err.response.data : err)
	}
);

export {axios};

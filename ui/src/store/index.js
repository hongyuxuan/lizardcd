import { createStore } from 'vuex'

export default createStore({
  state: {
    menuCollpase: false,
    username: '游客',
    avatar: '',
    role: '',
    settings: {},
  },
  modules: {
  }
})
import { createStore } from 'vuex'

export default createStore({
  state: {
    menuCollpase: false,
    userInfo: {
      username: '游客',
      profile: {}
    },
    settings: {},
  },
  modules: {
  }
})
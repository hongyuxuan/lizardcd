import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import 'element-plus/dist/index.css'

import './assets/css/bootstrap.css'
import './assets/css/AdminLTE.css'
import './assets/css/theme.css'
import './assets/css/app.css'

import App from './App.vue'
import router from './router'

/* 加载 font-awesome */
import { library } from '@fortawesome/fontawesome-svg-core'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import { faGears } from '@fortawesome/free-solid-svg-icons/faGears'
import {faCircle,faHouse,faLayerGroup,faSliders,faPlus,faMinus,faCubes} from '@fortawesome/free-solid-svg-icons'
library.add(faGears,faCircle,faHouse,faLayerGroup,faSliders,faPlus,faMinus,faCubes)

/* 加载 element-plus */
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

const app = createApp(App)
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app
  .use(createPinia())
  .use(ElementPlus, {locale: zhCn})
  .use(router)
  .component("font-awesome-icon", FontAwesomeIcon)
  .mount('#app')

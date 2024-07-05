import { createApp } from 'vue'
import ElementPlus from 'element-plus'
import zhCn from 'element-plus/es/locale/lang/zh-cn'
import 'element-plus/dist/index.css'
import './assets/css/bootstrap.css'
import './assets/css/AdminLTE.css'
import './assets/css/theme.css'
import './assets/css/app.css'
import './assets/css/style.css'
import App from './App.vue'
import router from './router'
/*** v-md-editor **/
import VMdPreview from '@kangc/v-md-editor/lib/preview';
import '@kangc/v-md-editor/lib/style/preview.css';
import vuepressTheme from '@kangc/v-md-editor/lib/theme/vuepress.js';
import '@kangc/v-md-editor/lib/theme/style/vuepress.css';
// Prism
import Prism from 'prismjs';
// highlight code
import 'prismjs/components/prism-json';

/* 加载 font-awesome */
import { library } from '@fortawesome/fontawesome-svg-core'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import {faGears,faCircle,faHouse,faLayerGroup,faSliders,faPlus,faMinus,faCubes,faGhost,faLaptopCode,faRocket,faCircleNodes,faTimeline,faFileCode,faListCheck} from '@fortawesome/free-solid-svg-icons'
import {faCircleQuestion as farCircleQuestion} from '@fortawesome/free-regular-svg-icons'
library.add(faGears,faCircle,faHouse,faLayerGroup,faSliders,faPlus,faMinus,faCubes,faGhost,faLaptopCode,faRocket,farCircleQuestion,faCircleNodes,faTimeline,faFileCode,faListCheck)

/* 加载 element-plus */
import * as ElementPlusIconsVue from '@element-plus/icons-vue'

/* 加载自定义组件 */
import myTips from './components/myTips'

/* 引入store */
import store from './store'

/* 引入VMdEditor */
VMdPreview.use(vuepressTheme, {
  Prism,
});

const app = createApp(App)
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component)
}

app
  .use(ElementPlus, {locale: zhCn})
  .use(router)
  .use(store)
  .use(VMdPreview)
  .use(myTips)
  .component("font-awesome-icon", FontAwesomeIcon)
  .mount('#app')

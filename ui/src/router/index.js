import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: () => import("../views/agent.vue")
    },
    {
      path: '/agent',
      name: 'agent',
      component: () => import("../views/agent.vue")
    },
    {
      path: '/template',
      name: 'template',
      component: () => import("../views/template.vue")
    },
    {
      path: '/workload/deployment',
      name: '',
      component: () => import("../views/workload/deployment.vue")
    },
    {
      path: '/workload/statefulset',
      name: '',
      component: () => import("../views/workload/statefulset.vue")
    },
  ]
})

export default router

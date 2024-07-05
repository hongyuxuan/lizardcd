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
      path: '/workload/deployments',
      name: '',
      component: () => import("../views/workload/deployments.vue")
    },
    {
      path: '/workload/statefulsets',
      name: '',
      component: () => import("../views/workload/statefulsets.vue")
    },
    {
      path: '/workload/deployments/:workload_name',
      name: '',
      component: () => import("../views/workload/deployment.vue")
    },
    {
      path: '/workload/statefulsets/:workload_name',
      name: '',
      component: () => import("../views/workload/statefulset.vue")
    },
    {
      path: '/platform/settings',
      name: '',
      component: () => import("../views/platform/configuration.vue")
    },
    {
      path: '/application',
      name: '',
      component: () => import("../views/application/application.vue")
    },
    {
      path: '/task/history',
      name: '',
      component: () => import("../views/task/history.vue")
    },
    {
      path: '/mesh/istio',
      name: '',
      component: () => import("../views/istio/istio.vue")
    },
    {
      path: '/helm',
      name: '',
      component: () => import("../views/helm/helm.vue")
    },
    {
      path: '/helm/:repo_name',
      name: '',
      component: () => import("../views/helm/charts.vue")
    },
    {
      path: '/helm/:repo_name/:chart_name',
      name: '',
      component: () => import("../views/helm/chart.vue")
    },
  ]
})

export default router

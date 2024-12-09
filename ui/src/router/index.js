import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      meta: {
        description: "首页",
      },
      component: () => import("../views/home.vue")
    },
    {
      path: '/agent',
      name: 'agent',
      meta: {
        description: "Agent管理",
      },
      component: () => import("../views/agent.vue")
    },
    {
      path: '/template',
      name: 'template',
      meta: {
        description: "模板管理",
      },
      component: () => import("../views/template.vue")
    },
    {
      path: '/workload/deployments',
      name: '',
      meta: {
        description: "部署",
      },
      component: () => import("../views/workload/deployments.vue")
    },
    {
      path: '/workload/statefulsets',
      name: '',
      meta: {
        description: "有状态副本集",
      },
      component: () => import("../views/workload/statefulsets.vue")
    },
    {
      path: '/workload/deployments/:workload_name',
      name: '',
      meta: {
        description: "部署详情",
      },
      component: () => import("../views/workload/deployment.vue")
    },
    {
      path: '/workload/statefulsets/:workload_name',
      name: '',
      meta: {
        description: "有状态副本集详情",
      },
      component: () => import("../views/workload/statefulset.vue")
    },
    {
      path: '/platform/settings',
      name: '',
      meta: {
        description: "配置",
      },
      component: () => import("../views/platform/configuration.vue")
    },
    {
      path: '/application',
      name: '',
      meta: {
        description: "应用管理",
      },
      component: () => import("../views/application/application.vue")
    },
    {
      path: '/application/:id',
      name: '',
      meta: {
        description: "应用详情",
      },
      component: () => import("../views/application/detail.vue")
    },
    {
      path: '/task/history',
      name: '',
      meta: {
        description: "任务管理",
      },
      component: () => import("../views/task/history.vue")
    },
    {
      path: '/mesh/istio',
      name: '',
      meta: {
        description: "Istio管理",
      },
      component: () => import("../views/istio/istio.vue")
    },
    {
      path: '/ci/tekton',
      name: '',
      meta: {
        description: "Tekton管理",
      },
      component: () => import("../views/tekton/tekton.vue")
    },
    {
      path: '/ci/trigger',
      name: '',
      meta: {
        description: "精准触发",
      },
      component: () => import("../views/tekton/trigger.vue")
    },
    {
      path: '/helm',
      name: '',
      meta: {
        description: "Helm管理",
      },
      component: () => import("../views/helm/helm.vue")
    },
    {
      path: '/helm/:repo_name',
      name: '',
      meta: {
        description: "Helm包列表",
      },
      component: () => import("../views/helm/charts.vue")
    },
    {
      path: '/helm/:repo_name/:chart_name',
      name: '',
      meta: {
        description: "Helm包安装",
      },
      component: () => import("../views/helm/chart.vue")
    },
    { 
      path: '/:path(.*)', 
      meta: {
        description: "页面未找到",
      },
      component: () => import("../views/404.vue"),
    }
  ]
})

export default router

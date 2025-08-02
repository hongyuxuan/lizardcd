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
      path: '/profile',
      name: 'profile',
      meta: {
        description: "个人设置",
      },
      component: () => import("../views/platform/profile.vue")
    },
    {
      path: '/agent',
      name: 'agent',
      meta: {
        description: "连接管理",
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
      path: '/kubernetes/workload',
      name: '',
      meta: {
        description: "工作负载",
      },
      component: () => import("../views/kubernetes/workload/workload.vue")
    },
    {
      path: '/kubernetes/services',
      name: '',
      meta: {
        description: "服务",
      },
      component: () => import("../views/kubernetes/service/services.vue")
    },
    {
      path: '/kubernetes/services/:service_name',
      name: '',
      meta: {
        description: "服务详情",
      },
      component: () => import("../views/kubernetes/service/service.vue")
    },
    {
      path: '/kubernetes/workload/deployments/:workload_name',
      name: '',
      meta: {
        description: "部署详情",
      },
      component: () => import("../views/kubernetes/workload/deployment.vue")
    },
    {
      path: '/kubernetes/workload/statefulsets/:workload_name',
      name: '',
      meta: {
        description: "有状态副本集详情",
      },
      component: () => import("../views/kubernetes/workload/statefulset.vue")
    },
    {
      path: '/kubernetes/ingresses',
      name: '',
      meta: {
        description: "路由",
      },
      component: () => import("../views/kubernetes/ingress/ingresses.vue")
    },
    {
      path: '/kubernetes/ingresses/:ingress_name',
      name: '',
      meta: {
        description: "路由详情",
      },
      component: () => import("../views/kubernetes/ingress/ingress.vue")
    },
    {
      path: '/kubernetes/workload/pods',
      name: '',
      meta: {
        description: "容器组",
      },
      component: () => import("../views/kubernetes/workload/pods.vue")
    },
    {
      path: '/kubernetes/workload/pods/:pod_name',
      name: '',
      meta: {
        description: "容器组详情",
      },
      component: () => import("../views/kubernetes/workload/pod.vue")
    },
    {
      path: '/kubernetes/jobs',
      name: '',
      meta: {
        description: "任务",
      },
      component: () => import("../views/kubernetes/job/jobs.vue")
    },
    {
      path: '/kubernetes/jobs/:job_name',
      name: '',
      meta: {
        description: "任务详情",
      },
      component: () => import("../views/kubernetes/job/job.vue")
    },
    {
      path: '/kubernetes/cronjobs/:job_name',
      name: '',
      meta: {
        description: "定时任务详情",
      },
      component: () => import("../views/kubernetes/job/cronjob.vue")
    },
    {
      path: '/kubernetes/persistentvolumeclaims',
      name: '',
      meta: {
        description: "存储卷",
      },
      component: () => import("../views/kubernetes/pvc/pvcs.vue")
    },
    {
      path: '/kubernetes/persistentvolumeclaims/:pvc_name',
      name: '',
      meta: {
        description: "存储卷详情",
      },
      component: () => import("../views/kubernetes/pvc/pvc.vue")
    },
    {
      path: '/kubernetes/configmaps',
      name: '',
      meta: {
        description: "配置字典",
      },
      component: () => import("../views/kubernetes/config/configmaps.vue")
    },
    {
      path: '/kubernetes/configmaps/:configmap_name',
      name: '',
      meta: {
        description: "配置字典详情",
      },
      component: () => import("../views/kubernetes/config/configmap.vue")
    },
    {
      path: '/kubernetes/secrets',
      name: '',
      meta: {
        description: "保密字典",
      },
      component: () => import("../views/kubernetes/config/secrets.vue")
    },
    {
      path: '/kubernetes/secrets/:secret_name',
      name: '',
      meta: {
        description: "保密字典详情",
      },
      component: () => import("../views/kubernetes/config/secret.vue")
    },
    {
      path: '/kubernetes/serviceaccounts',
      name: '',
      meta: {
        description: "服务账户",
      },
      component: () => import("../views/kubernetes/config/serviceaccounts.vue")
    },
    // {
    //   path: '/kubernetes/serviceaccounts/:serviceaccount_name',
    //   name: '',
    //   meta: {
    //     description: "保密字典详情",
    //   },
    //   component: () => import("../views/kubernetes/config/serviceaccount.vue")
    // },
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
      path: '/task/history/:id',
      name: '',
      meta: {
        description: "任务详情",
      },
      component: () => import("../views/task/detail.vue")
    },
    {
      path: '/application/release',
      name: '',
      meta: {
        description: "应用发布",
      },
      component: () => import("../views/application/release.vue")
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
        description: "流水线",
      },
      component: () => import("../views/tekton/tekton.vue")
    },
    {
      path: '/ci/tekton/pipelinerun/:name',
      name: '',
      meta: {
        description: "PipelineRun",
      },
      component: () => import("../views/tekton/pipelinerun/pipelinerunresult.vue")
    },
    {
      path: '/ci/tekton/taskrun/:name',
      name: '',
      meta: {
        description: "TaskRun",
      },
      component: () => import("../views/tekton/taskrun/taskrunresult.vue")
    },
    {
      path: '/ci/tekton/create',
      name: '',
      meta: {
        description: "创建流水线",
      },
      component: () => import("../views/tekton/pipeline/visualcreate.vue")
    },
    {
      path: '/ci/trigger',
      name: '',
      meta: {
        description: "触发配置",
      },
      component: () => import("../views/tekton/trigger/trigger.vue")
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

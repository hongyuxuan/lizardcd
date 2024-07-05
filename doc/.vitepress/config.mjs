import { defineConfig } from 'vitepress'

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "Lizardcd - 文档",
  description: "Lizardcd - Cloud Native Continuous Delivery Platform",
  base: "/lizardcd-doc/",
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: '主页', link: '/' },
      { text: '文档', link: '/docs/tintroduce' },
    ],

    head: [['link', { rel: 'icon', href: '/public/favicon.ico' }]],

    sidebar: [
      {
        text: '简介',
        items: [
          { text: 'Lizardcd简介', link: '/docs/introduce' }
        ]
      },
      {
        text: '快速开始',
        items: [
          { text: '快速部署', link: '/docs/quickstart/deployment' },
          { text: '新建应用', link: '/docs/quickstart/newapp' },
          { text: '查看任务', link: '/docs/quickstart/task' },
        ]
      },
      {
        text: '架构设计',
        items: [
          { text: '部署架构', link: '/docs/architecture/design' },
        ],
        collapsed: true,
      },
      {
        text: '功能介绍',
        items: [
          { text: '用户与权限管理', link: '/docs/features/rbac' },
          { text: '平台设置', link: '/docs/features/platform' },
          { text: '应用管理', 
            link: '/docs/features/application',
            items: [
              { text: '容器部署', link: '/docs/features/application/k8s' },
              { text: '虚拟机部署', link: '/docs/features/application/vm' },
              { text: 'HTTP部署', link: '/docs/features/application/http' }
            ]
          },
          { text: '任务管理', link: '/docs/features/task' },
          { text: '工作负载', link: '/docs/features/workload' },
          { text: 'Istio管理', link: '/docs/features/istio' },
          { text: 'Helm管理', link: '/docs/features/helm' },
        ],
        collapsed: true,
      },
      {
        text: '部署与配置',
        items: [
          { text: 'server部署', link: '/docs/deploy/server' },
          { text: 'agent部署', link: '/docs/deploy/agent' },
          { text: 'ui部署', link: '/docs/deploy/ui' },
          { text: '配置',
            items: [
              { text: 'agent配置', link: '/docs/deploy/configure/agent' },
              { text: 'server配置', link: '/docs/deploy/configure/server' },
              { text: 'RBAC权限设置', link: '/docs/deploy/configure/rbac' },
            ]
          },
        ],
        collapsed: true,
      },
      {
        text: '客户端用法',
        items: [
          { text: 'cli用法', link: '/docs/cli/cli' },
        ],
        collapsed: true,
      }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/hongyuxuan/Lizardcd' }
    ]
  }
})

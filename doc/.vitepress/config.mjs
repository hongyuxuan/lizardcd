import { defineConfig } from 'vitepress'
import mdItCustomAttrs from "markdown-it-custom-attrs";

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: "Lizardcd: Cloud Native Continuous Delivery Platform",
  description: "Lizardcd - Cloud Native Continuous Delivery Platform",
  base: "/lizardcd-doc/",
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    siteTitle: 'Lizardcd - 文档',
    nav: [
      { text: '主页', link: '/' },
      { text: '文档', link: '/docs/tintroduce' },
    ],

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
        ],
        collapsed: false,
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
              { text: 'HTTP 部署', link: '/docs/features/application/http' }
            ]
          },
          { text: '任务管理', link: '/docs/features/task' },
          { text: '工作负载', link: '/docs/features/workload' },
          { text: 'Istio 管理', link: '/docs/features/istio' },
          { text: 'Helm 管理', link: '/docs/features/helm' },
        ],
        collapsed: true,
      },
      {
        text: '部署与配置',
        items: [
          { text: 'agent 部署', link: '/docs/deploy/agent' },
          { text: 'server 部署', link: '/docs/deploy/server' },
          { text: 'ui 部署', link: '/docs/deploy/ui' },
          { text: '配置', link: '/docs/deploy/configure',
            items: [
              { text: 'agent 配置', link: '/docs/deploy/configure/agent' },
              { text: 'server 配置', link: '/docs/deploy/configure/server' },
            ]
          },
        ],
        collapsed: true,
      },
      {
        text: '客户端用法',
        items: [
          { text: 'cli 用法', link: '/docs/cli/cli' },
        ],
        collapsed: true,
      },
      {
        text: '运维',
        items: [
          { text: '链路跟踪', link: '/docs/operate/trace' },
          { text: '指标监控', link: '/docs/operate/prometheus' },
        ],
        collapsed: true,
      }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/hongyuxuan/Lizardcd' }
    ],
    outline:{
      label: '页面导航',
      level: [2,6]
    },
    lastUpdated: {
      text: '最后更新于'
    },
    search: {
      provider: 'local'
    }
  },
  markdown: {
    config: (md) => {
      // use more markdown-it plugins!
      md.use(mdItCustomAttrs, "image", {
        "data-fancybox": "gallery",
      });
    }
  },
  head: [
    ['link', { rel: 'icon', href: 'https://project-1255547500.cos.ap-beijing.myqcloud.com/lizardcd%2Ffavicon.ico' }],
    [
      "link",
      {
        rel: "stylesheet",
        href: "https://cdn.jsdelivr.net/npm/@fancyapps/ui/dist/fancybox.css",
      },
    ],
    ["script", { src: "https://cdn.jsdelivr.net/npm/@fancyapps/ui@4.0/dist/fancybox.umd.js" }],
  ]
})

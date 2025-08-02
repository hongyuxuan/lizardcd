import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import VueDevTools from 'vite-plugin-vue-devtools'
import { resolve } from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    VueDevTools(),
  ],
  build: {
    rollupOptions: {
      input: {
        main: resolve(__dirname, 'index.html'),
        login: resolve(__dirname, 'login/index.html')
      }
    }
  },
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    host: '0.0.0.0',
    port: '5173',
    proxy: {
      '/lizardcd': {
        target: 'http://localhost:5117',
        changeOrigin: true,
      },
      '/tekton-pipelines': {
        target: 'http://tekton-pipelines',
        changeOrigin: true,
      },
      '/ws': {
        target: 'ws://localhost:5117',
        rewriteWsOrigin: true,
        ws: true,
      },
    }
  }
})

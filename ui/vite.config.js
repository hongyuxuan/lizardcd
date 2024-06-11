import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import VueDevTools from 'vite-plugin-vue-devtools'
import { resolve } from 'path';

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
    proxy: {
      '/lizardcd': {
        target: 'http://localhost:5117',
        changeOrigin: true,
      },
      // '/kubernetes': {
      //   target: 'http://localhost:5117',
      //   changeOrigin: true,
      // },
      // '/istio': {
      //   target: 'http://localhost:5117',
      //   changeOrigin: true,
      // },
      // '/helm': {
      //   target: 'http://localhost:5117',
      //   changeOrigin: true,
      // },
      // '/db': {
      //   target: 'http://localhost:5117',
      //   changeOrigin: true,
      // },
      // '/auth': {
      //   target: 'http://localhost:5117',
      //   changeOrigin: true,
      // },
      // '/swagger': {
      //   target: 'http://localhost:9088/docs',
      //   changeOrigin: true,
      //   rewrite: (path) => path.replace(/^\/swagger/, '')
      // },
      // '/server-static': {
      //   target: 'http://localhost:5117',
      //   changeOrigin: true,
      // },
    }
  }
})

import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import VueDevTools from 'vite-plugin-vue-devtools'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    VueDevTools(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  server: {
    proxy: {
      '/consul': {
        // target: 'http://lizardcd-server.fiofa-ofa-bdev.cicc.io',
        target: 'http://localhost:5117',
        changeOrigin: true,
      },
      '/lizardcd': {
        // target: 'http://lizardcd-server.fiofa-ofa-bdev.cicc.io',
        target: 'http://localhost:5117',
        changeOrigin: true,
      },
    }
  }
})

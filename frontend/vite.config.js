import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';

export default defineConfig({
  base: '/messenger/',

  plugins: [vue()],

  resolve: {
    alias: {
      '@': '/src',
    },
  },

  server: {
    proxy: {
      '/messenger/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
});

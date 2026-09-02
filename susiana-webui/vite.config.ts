import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import path from 'node:path'
import { VitePWA } from 'vite-plugin-pwa'

const basePath = process.env.VITE_BASE_PATH || '/'
const apiProxyTarget = process.env.VITE_API_PROXY_TARGET || 'http://127.0.0.1:8080'
const baseNoSlash = basePath.replace(/\/$/, '') || ''

export default defineConfig({
  base: basePath,
  plugins: [
    vue(),
    VitePWA({
      registerType: 'autoUpdate',
      includeAssets: ['pwa-192.png', 'pwa-512.png', 'apple-touch-icon.png'],
      manifest: {
        name: 'Susiana',
        short_name: 'Susiana',
        description: 'Susiana market WebUI',
        start_url: basePath,
        display: 'standalone',
        background_color: '#0b1220',
        theme_color: '#0f766e',
        icons: [
          { src: 'pwa-192.png', sizes: '192x192', type: 'image/png', purpose: 'any' },
          { src: 'pwa-512.png', sizes: '512x512', type: 'image/png', purpose: 'any maskable' },
        ],
      },
      workbox: {
        globPatterns: ['**/*.{js,css,html,ico,png,svg,woff,woff2,webmanifest}'],
        navigateFallback: basePath === '/' ? '/index.html' : `${baseNoSlash}/index.html`,
        runtimeCaching: [
          {
            urlPattern: ({ request, url }) =>
              request.method === 'GET' && url.pathname.includes('/api/'),
            handler: 'NetworkOnly',
            options: { cacheName: 'susiana-api-get' },
          },
        ],
      },
      devOptions: { enabled: false },
    }),
  ],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
  server: {
    host: true,
    port: Number(process.env.VITE_DEV_PORT || 5178),
    allowedHosts: true,
    proxy: {
      '/api': {
        target: apiProxyTarget,
        changeOrigin: true,
      },
      '/health': {
        target: apiProxyTarget,
        changeOrigin: true,
      },
    },
  },
})

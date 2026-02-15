import tailwindcss from '@tailwindcss/vite'

export default defineNuxtConfig({
  compatibilityDate: '2024-11-01',
  devtools: { enabled: false },
  future: { compatibilityVersion: 4 },
  ssr: false,
  modules: ['shadcn-nuxt', '@pinia/nuxt'],
  shadcn: {
    prefix: '',
    componentDir: './app/components/ui',
  },
  css: ['~/assets/css/tailwind.css'],
  vite: {
    plugins: [tailwindcss()],
  },
  runtimeConfig: {
    public: {
      apiBase: '/api',
      wsBase: '/ws',
      firebaseApiKey: '',
      firebaseAuthDomain: 'cinema-25e75.firebaseapp.com',
      firebaseProjectId: 'cinema-25e75',
    },
  },
  nitro: {
    devProxy: {
      '/api': { target: 'http://localhost:8080/api', changeOrigin: true },
      '/ws': { target: 'http://localhost:8080/ws', ws: true },
    },
  },
})

import { defineNuxtConfig } from 'nuxt/config'
import tailwindcss from "@tailwindcss/vite";

export default defineNuxtConfig({
  css: [
    '@/main.css', // create this file
  ],
  components: true,
  modules: ['shadcn-nuxt'],
  vite: {
    plugins: [
      tailwindcss(),
    ],
  },
  shadcn: {
    /**
     * Prefix for all the imported component.
     * @default "Ui"
     */
    prefix: '',
    /**
     * Directory that the component lives in.
     * Will respect the Nuxt aliases.
     * @link https://nuxt.com/docs/api/nuxt-config#alias
     * @default "@/components/ui"
     */
    componentDir: '@/components/ui'
  },
  
  runtimeConfig: {
    public: {
      backendUrl: 'http://localhost:3000/v1',
    },
  },

  nitro: {
    devProxy: {
      //needed to do this workaround to call /api because its been acting up (redirect to nuxt pages instead of api calls)
      '/v1': {
        target: 'http://localhost:8080',
        changeOrigin: true,
        rewrite: (path) => path.replace(/^\/v1/, '')
      }
    }
  },

  compatibilityDate: "2026-08-01",
})
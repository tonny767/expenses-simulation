import { useAuth } from '~/composables/useAuth'

export default defineNuxtRouteMiddleware((to) => {
  const { isLoggedIn, isManager } = useAuth()
  const publicPages = ['/login']

  if (!isLoggedIn.value && !publicPages.includes(to.path)) {
    return navigateTo('/login')
  }

  if (isLoggedIn.value && to.path === '/login') {
    return navigateTo('/dashboard')
  }
})
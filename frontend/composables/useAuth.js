export const useAuth = () => {
  const role = useCookie('role', {
    maxAge: 60 * 60 * 24,
    sameSite: 'lax',
    path: '/',
  })

  const userName = useCookie('user_name', {
    maxAge: 60 * 60 * 24,
    sameSite: 'lax',
    path: '/',
  })

  const userId = useCookie('user_id', {
    maxAge: 60 * 60 * 24,
    sameSite: 'lax',
    path: '/',
  })

  const isLoggedIn = computed(() => !!userId.value)

  const roleValue = computed(() => role.value || 'user')
  const isManager = computed(() => role.value === 'manager')

  const requireAuth = () => {
    if (!isLoggedIn.value) {
      return navigateTo('/login')
    }
  }

  const requireRole = (requiredRole) => {
    if (roleValue.value !== requiredRole) {
      return navigateTo('/dashboard')
    }
  }

  const login = (roleVal, name, id) => {
    role.value = roleVal
    userName.value = name
    userId.value = id
  }

  const logout = async () => {
    role.value = null
    userName.value = null
    userId.value = null

    await navigateTo('/login')
  }

  return {
    role,
    userName,
    userId,

    isLoggedIn,
    roleValue,
    isManager,

    requireAuth,
    requireRole,
    login,
    logout,
  }
}

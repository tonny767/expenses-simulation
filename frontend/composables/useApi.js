export const useApi = (role) => {
  const config = useRuntimeConfig()

  const base =
    role === 'manager'
      ? '/v1/api/manager'  // Changed from /api to /backend/api
      : role === 'user'
      ? '/v1/api/user'
      : '/v1/api'

  const request = async (path, options = {}) => {
    // const isBrowser = import.meta.client
    const url = `${base}${path}`
    
    console.log('Request URL:', url)
    
    return await $fetch(url, {
      ...options,
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
        ...options.headers,
      },
    })
  }

  return {
    get: (path, options = {}) => request(path, { ...options, method: 'GET' }),
    post: (path, body, options = {}) => request(path, { ...options, method: 'POST', body }),
    put: (path, body, options = {}) => request(path, { ...options, method: 'PUT', body }),
    delete: (path, options = {}) => request(path, { ...options, method: 'DELETE' }),
  }
}
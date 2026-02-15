import axios from 'axios'

export function useApi() {
  const auth = useAuthStore()
  const api = axios.create({ baseURL: '/api' })

  api.interceptors.request.use((config) => {
    if (auth.token) {
      config.headers.Authorization = `Bearer ${auth.token}`
    }
    return config
  })

  api.interceptors.response.use(
    (res) => res,
    (err) => {
      if (err.response?.status === 401) {
        auth.logout()
        navigateTo('/login')
      }
      return Promise.reject(err)
    },
  )

  return api
}

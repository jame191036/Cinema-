import { defineStore } from 'pinia'

interface User {
  id: string
  email: string
  name: string
  role: string
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null as User | null,
    token: '' as string,
  }),
  getters: {
    isAuthenticated: (state) => !!state.token,
    isAdmin: (state) => state.user?.role === 'ADMIN',
  },
  actions: {
    setAuth(data: { user: User; token: string }) {
      this.user = data.user
      this.token = data.token
      if (import.meta.client) {
        localStorage.setItem('cinema_token', data.token)
        localStorage.setItem('cinema_user', JSON.stringify(data.user))
      }
    },
    logout() {
      this.user = null
      this.token = ''
      if (import.meta.client) {
        localStorage.removeItem('cinema_token')
        localStorage.removeItem('cinema_user')
      }
    },
    restore() {
      if (import.meta.client) {
        const token = localStorage.getItem('cinema_token')
        const userStr = localStorage.getItem('cinema_user')
        if (token && userStr) {
          this.token = token
          try { this.user = JSON.parse(userStr) } catch {}
        }
      }
    },
  },
})

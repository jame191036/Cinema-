<script setup lang="ts">
import { GoogleAuthProvider, signInWithPopup } from 'firebase/auth'

definePageMeta({ layout: 'default' })

const { $firebaseAuth } = useNuxtApp()
const api = useApi()
const auth = useAuthStore()
const email = ref('user@demo.com')
const name = ref('Demo User')
const loading = ref(false)
const googleLoading = ref(false)
const error = ref('')

async function login(role: string) {
  loading.value = true
  error.value = ''
  try {
    const { data } = await api.post('/auth/login', {
      email: email.value,
      name: name.value,
      demo: true,
      role,
    })
    auth.setAuth({ user: data.user, token: data.token })
    navigateTo('/')
  } catch (e: any) {
    error.value = e.response?.data?.error || 'Login failed'
  } finally {
    loading.value = false
  }
}

async function loginWithGoogle() {
  googleLoading.value = true
  error.value = ''
  try {
    const provider = new GoogleAuthProvider()
    const result = await signInWithPopup($firebaseAuth, provider)
    const idToken = await result.user.getIdToken()

    const { data } = await api.post('/auth/google', { id_token: idToken })
    auth.setAuth({ user: data.user, token: data.token })
    navigateTo('/')
  } catch (e: any) {
    if (e.code === 'auth/popup-closed-by-user') {
      error.value = 'Sign-in popup was closed'
    } else if (e.code === 'auth/cancelled-popup-request') {
      // Ignore duplicate popup
    } else {
      error.value = e.response?.data?.error || e.message || 'Google sign-in failed'
    }
  } finally {
    googleLoading.value = false
  }
}
</script>

<template>
  <div class="flex min-h-[60vh] items-center justify-center">
    <Card class="w-full max-w-md">
      <CardHeader class="text-center">
        <div class="mb-2 text-4xl">🎬</div>
        <CardTitle class="text-2xl">Welcome to Cinema</CardTitle>
        <p class="text-sm text-muted-foreground">Sign in to book your seats</p>
      </CardHeader>
      <CardContent class="space-y-4">
        <div v-if="error" class="rounded-lg bg-destructive/10 p-3 text-sm text-destructive">{{ error }}</div>

        <!-- Google Sign-In -->
        <Button
          class="w-full gap-2"
          variant="outline"
          :disabled="googleLoading || loading"
          @click="loginWithGoogle"
        >
          <svg class="h-5 w-5" viewBox="0 0 24 24">
            <path
              d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92a5.06 5.06 0 0 1-2.2 3.32v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.1z"
              fill="#4285F4"
            />
            <path
              d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
              fill="#34A853"
            />
            <path
              d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
              fill="#FBBC05"
            />
            <path
              d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
              fill="#EA4335"
            />
          </svg>
          {{ googleLoading ? 'Signing in...' : 'Sign in with Google' }}
        </Button>

        <div class="relative">
          <div class="absolute inset-0 flex items-center">
            <Separator class="w-full" />
          </div>
          <div class="relative flex justify-center text-xs uppercase">
            <span class="bg-card px-2 text-muted-foreground">Or demo login</span>
          </div>
        </div>

        <!-- Demo Login -->
        <div class="space-y-2">
          <label class="text-sm font-medium">Email</label>
          <Input v-model="email" type="email" placeholder="your@email.com" />
        </div>
        <div class="space-y-2">
          <label class="text-sm font-medium">Display Name</label>
          <Input v-model="name" type="text" placeholder="Your name" />
        </div>
      </CardContent>
      <CardFooter class="flex gap-3">
        <Button class="flex-1" :disabled="loading || googleLoading" @click="login('USER')">
          {{ loading ? 'Signing in...' : 'Demo User' }}
        </Button>
        <Button class="flex-1" variant="secondary" :disabled="loading || googleLoading" @click="login('ADMIN')">
          Demo Admin
        </Button>
      </CardFooter>
    </Card>
  </div>
</template>

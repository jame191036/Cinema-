<script setup lang="ts">
const auth = useAuthStore()
const router = useRouter()

function handleLogout() {
  auth.logout()
  navigateTo('/login')
}
</script>

<template>
  <div class="min-h-screen bg-background text-foreground">
    <nav class="sticky top-0 z-50 border-b border-border bg-card/95 backdrop-blur supports-[backdrop-filter]:bg-card/60">
      <div class="mx-auto flex h-14 max-w-7xl items-center justify-between px-4 sm:px-6">
        <div class="flex items-center gap-6">
          <NuxtLink to="/" class="flex items-center gap-2 text-lg font-bold tracking-tight">
            <span class="text-xl">🎬</span>
            <span>Cinema</span>
          </NuxtLink>
          <Separator orientation="vertical" class="h-6" />
          <nav class="flex items-center gap-4 text-sm">
            <NuxtLink to="/" class="text-muted-foreground transition-colors hover:text-foreground" active-class="text-foreground font-medium">Movies</NuxtLink>
            <template v-if="auth.isAdmin">
              <NuxtLink to="/admin/bookings" class="text-muted-foreground transition-colors hover:text-foreground" active-class="text-foreground font-medium">Bookings</NuxtLink>
              <NuxtLink to="/admin/audit-logs" class="text-muted-foreground transition-colors hover:text-foreground" active-class="text-foreground font-medium">Audit Logs</NuxtLink>
            </template>
          </nav>
        </div>
        <div class="flex items-center gap-3">
          <template v-if="auth.isAuthenticated">
            <span class="text-sm text-muted-foreground">{{ auth.user?.name }}</span>
            <Badge v-if="auth.isAdmin" variant="secondary" class="text-xs">Admin</Badge>
            <Button variant="ghost" size="sm" @click="handleLogout">Logout</Button>
          </template>
          <template v-else>
            <Button size="sm" as-child>
              <NuxtLink to="/login">Login</NuxtLink>
            </Button>
          </template>
        </div>
      </div>
    </nav>
    <main class="mx-auto max-w-7xl px-4 py-8 sm:px-6">
      <slot />
    </main>
  </div>
</template>

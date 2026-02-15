<script setup lang="ts">
definePageMeta({ middleware: ['auth'] })

const route = useRoute()
const api = useApi()
const movieId = route.params.movieId as string
const movie = ref<any>(null)
const showtimes = ref<any[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    const [mRes, stRes] = await Promise.all([
      api.get(`/movies/${movieId}`),
      api.get('/showtimes', { params: { movie_id: movieId } }),
    ])
    movie.value = mRes.data
    showtimes.value = stRes.data || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
})

function formatTime(dt: string) {
  return new Date(dt).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}
function formatDate(dt: string) {
  return new Date(dt).toLocaleDateString([], { weekday: 'short', month: 'short', day: 'numeric' })
}
</script>

<template>
  <div>
    <NuxtLink to="/" class="mb-6 inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground">
      ← Back to Movies
    </NuxtLink>

    <div v-if="loading" class="space-y-4">
      <div class="h-8 w-48 animate-pulse rounded bg-muted" />
      <div class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <Card v-for="i in 4" :key="i" class="animate-pulse"><CardContent class="p-6"><div class="h-12 rounded bg-muted" /></CardContent></Card>
      </div>
    </div>

    <template v-else>
      <div class="mb-6">
        <h1 class="text-2xl font-bold tracking-tight">{{ movie?.title }}</h1>
        <div class="mt-1 flex items-center gap-2">
          <Badge variant="secondary">{{ movie?.durationMin }} min</Badge>
          <Badge v-if="movie?.rating" variant="outline">{{ movie?.rating }}</Badge>
        </div>
      </div>

      <h2 class="mb-4 text-lg font-semibold">Available Showtimes</h2>

      <div v-if="showtimes.length === 0" class="py-16 text-center text-muted-foreground">No showtimes available</div>

      <div v-else class="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <NuxtLink v-for="st in showtimes" :key="st.id" :to="`/seats/${st.id}`">
          <Card class="transition-colors hover:border-primary/50 hover:bg-accent/50 cursor-pointer">
            <CardContent class="p-6">
              <div class="text-2xl font-bold">{{ formatTime(st.startTime) }}</div>
              <div class="mt-1 text-sm text-muted-foreground">{{ formatDate(st.startTime) }}</div>
              <Separator class="my-3" />
              <div class="text-xs text-muted-foreground">Hall: {{ st.auditoriumId }}</div>
            </CardContent>
          </Card>
        </NuxtLink>
      </div>
    </template>
  </div>
</template>

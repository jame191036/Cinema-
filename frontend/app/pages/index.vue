<script setup lang="ts">
const api = useApi()
const movies = ref<any[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    const { data } = await api.get('/movies')
    movies.value = data
  } catch (e) {
    console.error('Failed to load movies', e)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div>
    <div class="mb-8">
      <h1 class="text-3xl font-bold tracking-tight">Now Showing</h1>
      <p class="mt-1 text-muted-foreground">Select a movie to view available showtimes</p>
    </div>

    <div v-if="loading" class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
      <Card v-for="i in 6" :key="i" class="animate-pulse">
        <CardHeader><div class="h-5 w-3/4 rounded bg-muted" /></CardHeader>
        <CardContent><div class="h-4 w-1/2 rounded bg-muted" /></CardContent>
      </Card>
    </div>

    <div v-else-if="movies.length === 0" class="flex flex-col items-center justify-center py-20 text-muted-foreground">
      <span class="text-4xl mb-4">🎬</span>
      <p>No movies available</p>
    </div>

    <div v-else class="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
      <NuxtLink v-for="movie in movies" :key="movie.id" :to="`/showtimes/${movie.id}`">
        <Card class="h-full transition-colors hover:border-primary/50 hover:bg-accent/50 cursor-pointer">
          <CardHeader>
            <CardTitle class="text-lg">{{ movie.title }}</CardTitle>
          </CardHeader>
          <CardContent>
            <div class="flex items-center gap-3">
              <Badge variant="secondary">{{ movie.durationMin }} min</Badge>
              <Badge v-if="movie.rating" variant="outline">{{ movie.rating }}</Badge>
            </div>
          </CardContent>
        </Card>
      </NuxtLink>
    </div>
  </div>
</template>

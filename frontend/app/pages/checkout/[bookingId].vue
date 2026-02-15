<script setup lang="ts">
definePageMeta({ middleware: ['auth'] })

const route = useRoute()
const api = useApi()
const bid = route.params.bookingId as string

const booking = ref<any>(null)
const loading = ref(true)
const paying = ref(false)
const success = ref(false)
const error = ref('')
const countdown = ref('')
let timer: ReturnType<typeof setInterval> | null = null

onMounted(async () => {
  try {
    const { data } = await api.get(`/bookings/${bid}`)
    booking.value = data
    if (data.lockExpiresAt) {
      const exp = new Date(data.lockExpiresAt)
      timer = setInterval(() => {
        const diff = exp.getTime() - Date.now()
        if (diff <= 0) { countdown.value = 'EXPIRED'; clearInterval(timer!); return }
        const m = Math.floor(diff / 60000)
        const s = Math.floor((diff % 60000) / 1000)
        countdown.value = `${m}:${s.toString().padStart(2, '0')}`
      }, 1000)
    }
  } catch { error.value = 'Failed to load booking' }
  finally { loading.value = false }
})

async function payAndConfirm() {
  paying.value = true; error.value = ''
  try {
    await api.post(`/bookings/${bid}/pay`)
    await api.post(`/bookings/${bid}/confirm`)
    success.value = true
    if (timer) clearInterval(timer)
  } catch (e: any) { error.value = e.response?.data?.error || 'Payment failed' }
  finally { paying.value = false }
}

async function cancel() {
  try { await api.post(`/bookings/${bid}/cancel`); navigateTo('/') }
  catch (e: any) { error.value = e.response?.data?.error || 'Cancel failed' }
}

onUnmounted(() => { if (timer) clearInterval(timer) })
</script>

<template>
  <div class="flex min-h-[60vh] items-center justify-center">
    <div v-if="loading" class="h-64 w-full max-w-md animate-pulse rounded-xl bg-muted" />

    <Card v-else-if="success" class="w-full max-w-md border-emerald-500/30 bg-emerald-500/5 text-center">
      <CardHeader>
        <div class="text-5xl mb-2">🎉</div>
        <CardTitle class="text-2xl text-emerald-400">Booking Confirmed!</CardTitle>
      </CardHeader>
      <CardContent>
        <p class="text-muted-foreground">Seats: <span class="font-semibold text-foreground">{{ booking?.seats?.join(', ') }}</span></p>
      </CardContent>
      <CardFooter class="justify-center">
        <Button as-child><NuxtLink to="/">Back to Movies</NuxtLink></Button>
      </CardFooter>
    </Card>

    <Card v-else class="w-full max-w-md">
      <CardHeader class="text-center">
        <CardTitle class="text-2xl">Checkout</CardTitle>
      </CardHeader>
      <CardContent class="space-y-4">
        <div v-if="error" class="rounded-lg border border-destructive/50 bg-destructive/10 p-3 text-sm text-destructive">{{ error }}</div>

        <div v-if="booking" class="space-y-3">
          <div class="rounded-lg bg-muted p-4">
            <div class="text-xs text-muted-foreground">Seats</div>
            <div class="mt-1 text-lg font-bold">{{ booking.seats?.join(', ') }}</div>
          </div>
          <div class="flex gap-3">
            <div class="flex-1 rounded-lg bg-muted p-4">
              <div class="text-xs text-muted-foreground">Status</div>
              <Badge class="mt-1" :variant="booking.status === 'LOCKED' ? 'secondary' : 'default'">{{ booking.status }}</Badge>
            </div>
            <div v-if="countdown" class="flex-1 rounded-lg bg-muted p-4">
              <div class="text-xs text-muted-foreground">Time Left</div>
              <div class="mt-1 font-mono text-lg font-bold" :class="countdown === 'EXPIRED' ? 'text-red-400' : 'text-blue-300'">{{ countdown }}</div>
            </div>
          </div>
        </div>
      </CardContent>
      <CardFooter class="gap-3">
        <Button variant="ghost" class="flex-1" @click="cancel">Cancel</Button>
        <Button class="flex-1" :disabled="paying || countdown === 'EXPIRED'" @click="payAndConfirm">
          {{ paying ? 'Processing...' : 'Pay & Confirm' }}
        </Button>
      </CardFooter>
    </Card>
  </div>
</template>

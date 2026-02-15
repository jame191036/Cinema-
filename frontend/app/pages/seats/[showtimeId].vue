<script setup lang="ts">
definePageMeta({ middleware: ['auth'] })

const route = useRoute()
const api = useApi()
const auth = useAuthStore()
const showtimeId = route.params.showtimeId as string

const seats = ref<any[]>([])
const selectedSeats = ref<Set<string>>(new Set())
const loading = ref(true)
const locking = ref(false)
const bookingId = ref('')
const lockExpiresAt = ref<Date | null>(null)
const countdown = ref('')
const wsConnected = ref(false)
const error = ref('')

let ws: WebSocket | null = null
let countdownTimer: ReturnType<typeof setInterval> | null = null

const seatsByRow = computed(() => {
  const map = new Map<string, any[]>()
  for (const s of seats.value) {
    const row = s.seatCode.charAt(0)
    if (!map.has(row)) map.set(row, [])
    map.get(row)!.push(s)
  }
  for (const [, rs] of map) {
    rs.sort((a: any, b: any) => parseInt(a.seatCode.substring(1)) - parseInt(b.seatCode.substring(1)))
  }
  return map
})

function seatColor(seat: any) {
  if (seat.state === 'BOOKED') return 'bg-red-600/80 cursor-not-allowed'
  if (seat.state === 'LOCKED' && seat.lockedByUserId === auth.user?.id) return 'bg-blue-500'
  if (seat.state === 'LOCKED') return 'bg-amber-500 cursor-not-allowed'
  if (selectedSeats.value.has(seat.seatCode)) return 'bg-indigo-500 ring-2 ring-white/50'
  return 'bg-emerald-600 hover:bg-emerald-500 cursor-pointer'
}

function toggleSeat(seat: any) {
  if (seat.state !== 'AVAILABLE' || bookingId.value) return
  const s = new Set(selectedSeats.value)
  if (s.has(seat.seatCode)) s.delete(seat.seatCode)
  else s.add(seat.seatCode)
  selectedSeats.value = s
}

async function lockSeats() {
  if (selectedSeats.value.size === 0) return
  locking.value = true
  error.value = ''
  try {
    const { data } = await api.post(`/showtimes/${showtimeId}/seats/lock`, {
      seats: Array.from(selectedSeats.value),
    })
    bookingId.value = data.bookingId
    lockExpiresAt.value = new Date(data.lockExpiresAt)
    selectedSeats.value = new Set()
    startCountdown()
  } catch (e: any) {
    error.value = e.response?.data?.error || 'Failed to lock seats'
  } finally {
    locking.value = false
  }
}

async function cancelBooking() {
  if (!bookingId.value) return
  try {
    await api.post(`/bookings/${bookingId.value}/cancel`)
    bookingId.value = ''
    lockExpiresAt.value = null
    countdown.value = ''
    if (countdownTimer) clearInterval(countdownTimer)
  } catch (e: any) {
    error.value = e.response?.data?.error || 'Cancel failed'
  }
}

function startCountdown() {
  if (countdownTimer) clearInterval(countdownTimer)
  countdownTimer = setInterval(() => {
    if (!lockExpiresAt.value) return
    const diff = lockExpiresAt.value.getTime() - Date.now()
    if (diff <= 0) {
      countdown.value = 'EXPIRED'
      bookingId.value = ''
      lockExpiresAt.value = null
      clearInterval(countdownTimer!)
      return
    }
    const m = Math.floor(diff / 60000)
    const s = Math.floor((diff % 60000) / 1000)
    countdown.value = `${m}:${s.toString().padStart(2, '0')}`
  }, 1000)
}

function connectWs() {
  const proto = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  ws = new WebSocket(`${proto}//${window.location.host}/ws/showtimes/${showtimeId}?token=${auth.token}`)
  ws.onopen = () => { wsConnected.value = true }
  ws.onclose = () => {
    wsConnected.value = false
    setTimeout(connectWs, 3000)
  }
  ws.onmessage = (e) => {
    try { handleMsg(JSON.parse(e.data)) } catch {}
  }
}

function handleMsg(msg: any) {
  if (msg.type === 'SYNC_SNAPSHOT' && Array.isArray(msg.seats)) {
    seats.value = msg.seats
    return
  }
  const idx = seats.value.findIndex((s) => s.seatCode === msg.seatCode)
  if (idx === -1) return
  const updated = { ...seats.value[idx] }
  if (msg.type === 'SEAT_LOCKED') {
    updated.state = 'LOCKED'
    updated.lockedByUserId = msg.lockedByUserId
    updated.lockExpiresAt = msg.lockExpiresAt
  } else if (msg.type === 'SEAT_RELEASED') {
    updated.state = 'AVAILABLE'
    updated.lockedByUserId = null
    updated.lockExpiresAt = null
  } else if (msg.type === 'SEAT_BOOKED') {
    updated.state = 'BOOKED'
  }
  seats.value[idx] = updated
}

onMounted(async () => {
  try {
    const { data } = await api.get(`/showtimes/${showtimeId}/seats`)
    seats.value = data
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
  connectWs()
})

onUnmounted(() => {
  if (ws) { ws.onclose = null; ws.close() }
  if (countdownTimer) clearInterval(countdownTimer)
})
</script>

<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <div>
        <NuxtLink to="/" class="mb-2 inline-flex items-center text-sm text-muted-foreground hover:text-foreground">← Back</NuxtLink>
        <h1 class="text-2xl font-bold tracking-tight">Select Your Seats</h1>
      </div>
      <div class="flex items-center gap-2 text-sm text-muted-foreground">
        <span :class="wsConnected ? 'bg-emerald-500' : 'bg-red-500'" class="h-2 w-2 rounded-full" />
        {{ wsConnected ? 'Live' : 'Reconnecting...' }}
      </div>
    </div>

    <div v-if="error" class="mb-4 rounded-lg border border-destructive/50 bg-destructive/10 p-3 text-sm text-destructive">{{ error }}</div>

    <!-- Lock info bar -->
    <Card v-if="bookingId" class="mb-6 border-blue-500/30 bg-blue-500/5">
      <CardContent class="flex items-center justify-between p-4">
        <div class="flex items-center gap-3">
          <Badge variant="secondary" class="bg-blue-500/20 text-blue-400">Locked</Badge>
          <span class="text-sm text-muted-foreground">Time remaining:</span>
          <span class="font-mono text-lg font-bold" :class="countdown === 'EXPIRED' ? 'text-red-400' : 'text-blue-300'">{{ countdown }}</span>
        </div>
        <div class="flex gap-2">
          <Button variant="ghost" size="sm" @click="cancelBooking">Cancel</Button>
          <Button size="sm" as-child>
            <NuxtLink :to="`/checkout/${bookingId}`">Pay & Confirm →</NuxtLink>
          </Button>
        </div>
      </CardContent>
    </Card>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center py-20">
      <div class="h-8 w-8 animate-spin rounded-full border-2 border-primary border-t-transparent" />
    </div>

    <template v-else>
      <!-- Screen -->
      <div class="mx-auto mb-10 max-w-lg text-center">
        <div class="mx-auto h-1 w-3/4 rounded-t-full bg-muted-foreground/30" />
        <p class="mt-2 text-xs font-medium tracking-widest text-muted-foreground">SCREEN</p>
      </div>

      <!-- Seat Grid -->
      <div class="mx-auto max-w-2xl space-y-2">
        <div v-for="[rowLabel, rowSeats] in seatsByRow" :key="rowLabel" class="flex items-center gap-3">
          <span class="w-5 text-center text-xs font-bold text-muted-foreground">{{ rowLabel }}</span>
          <div class="flex flex-1 justify-center gap-1.5">
            <button
              v-for="seat in rowSeats" :key="seat.seatCode"
              @click="toggleSeat(seat)"
              :disabled="seat.state === 'BOOKED' || (seat.state === 'LOCKED' && seat.lockedByUserId !== auth.user?.id)"
              :class="seatColor(seat)"
              class="flex h-8 w-8 items-center justify-center rounded-t-lg text-[10px] font-semibold text-white transition-all disabled:opacity-70"
              :title="seat.seatCode + ' - ' + seat.state"
            >
              {{ seat.seatCode.substring(1) }}
            </button>
          </div>
          <span class="w-5 text-center text-xs font-bold text-muted-foreground">{{ rowLabel }}</span>
        </div>
      </div>

      <!-- Legend -->
      <div class="mx-auto mt-8 flex max-w-lg flex-wrap justify-center gap-x-5 gap-y-2 text-xs text-muted-foreground">
        <div class="flex items-center gap-1.5"><span class="h-3 w-3 rounded-sm bg-emerald-600" /> Available</div>
        <div class="flex items-center gap-1.5"><span class="h-3 w-3 rounded-sm bg-indigo-500 ring-1 ring-white/50" /> Selected</div>
        <div class="flex items-center gap-1.5"><span class="h-3 w-3 rounded-sm bg-blue-500" /> My Lock</div>
        <div class="flex items-center gap-1.5"><span class="h-3 w-3 rounded-sm bg-amber-500" /> Locked</div>
        <div class="flex items-center gap-1.5"><span class="h-3 w-3 rounded-sm bg-red-600/80" /> Booked</div>
      </div>

      <!-- Lock button -->
      <div v-if="!bookingId && selectedSeats.size > 0" class="mt-8 flex justify-center">
        <Button size="lg" :disabled="locking" @click="lockSeats">
          {{ locking ? 'Locking...' : `Lock ${selectedSeats.size} Seat(s)` }}
        </Button>
      </div>
    </template>
  </div>
</template>

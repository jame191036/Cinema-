<script setup lang="ts">
definePageMeta({ middleware: ['auth', 'admin'] })

const api = useApi()
const bookings = ref<any[]>([])
const loading = ref(true)
const statusFilter = ref('')
const dateFilter = ref('')

async function fetchBookings() {
  loading.value = true
  try {
    const params: Record<string, string> = {}
    if (statusFilter.value) params.status = statusFilter.value
    if (dateFilter.value) params.date = dateFilter.value
    const { data } = await api.get('/admin/bookings', { params })
    bookings.value = data || []
  } catch (e) { console.error(e) }
  finally { loading.value = false }
}

onMounted(fetchBookings)

function statusVariant(s: string) {
  if (s === 'BOOKED') return 'default' as const
  if (s === 'LOCKED') return 'secondary' as const
  return 'outline' as const
}

function fmtDt(dt: string) { return new Date(dt).toLocaleString() }
function shortId(id: string) { return id?.substring(0, 8) + '...' }
</script>

<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <h1 class="text-2xl font-bold tracking-tight">Bookings</h1>
      <Button variant="outline" size="sm" @click="fetchBookings">Refresh</Button>
    </div>

    <div class="mb-4 flex flex-wrap gap-3">
      <select v-model="statusFilter" @change="fetchBookings" class="rounded-md border border-input bg-background px-3 py-2 text-sm">
        <option value="">All Statuses</option>
        <option v-for="s in ['LOCKED','BOOKED','CANCELLED','EXPIRED']" :key="s" :value="s">{{ s }}</option>
      </select>
      <Input v-model="dateFilter" type="date" class="w-auto" @change="fetchBookings" />
    </div>

    <Card>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>ID</TableHead>
            <TableHead>User</TableHead>
            <TableHead>Showtime</TableHead>
            <TableHead>Seats</TableHead>
            <TableHead>Status</TableHead>
            <TableHead>Created</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-if="loading">
            <TableCell colspan="6" class="py-10 text-center text-muted-foreground">Loading...</TableCell>
          </TableRow>
          <TableRow v-else-if="bookings.length === 0">
            <TableCell colspan="6" class="py-10 text-center text-muted-foreground">No bookings found</TableCell>
          </TableRow>
          <TableRow v-for="b in bookings" :key="b.id" v-else>
            <TableCell class="font-mono text-xs">{{ shortId(b.id) }}</TableCell>
            <TableCell>{{ shortId(b.userId) }}</TableCell>
            <TableCell>{{ shortId(b.showtimeId) }}</TableCell>
            <TableCell>{{ b.seats?.join(', ') }}</TableCell>
            <TableCell><Badge :variant="statusVariant(b.status)">{{ b.status }}</Badge></TableCell>
            <TableCell class="text-muted-foreground">{{ fmtDt(b.createdAt) }}</TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </Card>
  </div>
</template>

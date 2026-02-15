<script setup lang="ts">
definePageMeta({ middleware: ['auth', 'admin'] })

const api = useApi()
const logs = ref<any[]>([])
const loading = ref(true)
const eventType = ref('')
const dateFrom = ref('')
const dateTo = ref('')

async function fetchLogs() {
  loading.value = true
  try {
    const params: Record<string, string> = {}
    if (eventType.value) params.event_type = eventType.value
    if (dateFrom.value) params.date_from = dateFrom.value
    if (dateTo.value) params.date_to = dateTo.value
    const { data } = await api.get('/admin/audit-logs', { params })
    logs.value = data || []
  } catch (e) { console.error(e) }
  finally { loading.value = false }
}

onMounted(fetchLogs)

function eventVariant(t: string) {
  if (t.includes('SUCCESS')) return 'default' as const
  if (t.includes('TIMEOUT') || t.includes('ERROR')) return 'destructive' as const
  return 'secondary' as const
}

function fmtDt(dt: string) { return new Date(dt).toLocaleString() }
function shortId(id: string | undefined) { return id ? id.substring(0, 8) + '...' : '-' }
</script>

<template>
  <div>
    <div class="mb-6 flex items-center justify-between">
      <h1 class="text-2xl font-bold tracking-tight">Audit Logs</h1>
      <Button variant="outline" size="sm" @click="fetchLogs">Refresh</Button>
    </div>

    <div class="mb-4 flex flex-wrap gap-3">
      <select v-model="eventType" @change="fetchLogs" class="rounded-md border border-input bg-background px-3 py-2 text-sm">
        <option value="">All Events</option>
        <option v-for="e in ['BOOKING_SUCCESS','BOOKING_TIMEOUT','SEAT_LOCKED','SEAT_RELEASED','SYSTEM_ERROR']" :key="e" :value="e">{{ e }}</option>
      </select>
      <Input v-model="dateFrom" type="date" class="w-auto" placeholder="From" @change="fetchLogs" />
      <Input v-model="dateTo" type="date" class="w-auto" placeholder="To" @change="fetchLogs" />
    </div>

    <Card>
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Time</TableHead>
            <TableHead>Event</TableHead>
            <TableHead>User</TableHead>
            <TableHead>Seat</TableHead>
            <TableHead>Booking</TableHead>
            <TableHead>Details</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow v-if="loading">
            <TableCell colspan="6" class="py-10 text-center text-muted-foreground">Loading...</TableCell>
          </TableRow>
          <TableRow v-else-if="logs.length === 0">
            <TableCell colspan="6" class="py-10 text-center text-muted-foreground">No audit logs found</TableCell>
          </TableRow>
          <TableRow v-for="log in logs" :key="log.id" v-else>
            <TableCell class="whitespace-nowrap text-muted-foreground">{{ fmtDt(log.createdAt) }}</TableCell>
            <TableCell><Badge :variant="eventVariant(log.eventType)">{{ log.eventType }}</Badge></TableCell>
            <TableCell class="font-mono text-xs">{{ shortId(log.userId) }}</TableCell>
            <TableCell>{{ log.seatCode || '-' }}</TableCell>
            <TableCell class="font-mono text-xs">{{ shortId(log.bookingId) }}</TableCell>
            <TableCell class="max-w-[200px] truncate text-xs text-muted-foreground">{{ JSON.stringify(log.payload) }}</TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </Card>
  </div>
</template>

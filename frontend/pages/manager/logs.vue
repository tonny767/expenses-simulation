<template>
  <Container class="p-6 space-y-6">
    <div>
      <h2 class="text-xl font-bold">Audit Logs</h2>
      <p class="text-muted-foreground">
        History of all expense status changes
      </p>
    </div>

    <div class="rounded-lg border p-4 space-y-4">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Expense ID</TableHead>
            <TableHead>Actor</TableHead>
            <TableHead>From Status</TableHead>
            <TableHead>To Status</TableHead>
            <TableHead>Reason</TableHead>
            <TableHead>Timestamp</TableHead>
          </TableRow>
        </TableHeader>

        <TableBody>
          <TableRow v-for="log in auditLogs" :key="log.id">
            <TableCell>
              <NuxtLink
                :to="`/expenses/${log.expense_id}`"
                class="font-medium text-blue-600 hover:underline"
              >
                #{{ log.expense_id }}
              </NuxtLink>
            </TableCell>
            <TableCell>
                <span class="text-sm">{{ log.actor_id || 'system' }}</span>
            </TableCell>

            <TableCell>
              <span
                :class="[
                  'px-2 py-1 text-sm rounded',
                  getStatusBadgeClass(log.from_status)
                ]"
              >
                {{ log.from_status.toUpperCase() }}
              </span>              
            </TableCell>

            <TableCell>
              <span
                :class="[
                  'px-2 py-1 text-sm rounded',
                  getStatusBadgeClass(log.to_status)
                ]"
              >
                {{ log.to_status.toUpperCase() }}
              </span>
            </TableCell>

            <TableCell>
              <span v-if="log.reason" class="text-sm">{{ log.reason }}</span>
              <span v-else class="text-sm text-muted-foreground">-</span>
            </TableCell>

            <TableCell>{{ formatDateTime(log.created_at) }}</TableCell>
          </TableRow>
        </TableBody>
      </Table>

      <div v-if="!auditLogs.length" class="p-6 text-center text-muted-foreground">
        No audit logs available
      </div>

      <div
        v-if="totalPages > 1"
        class="flex items-center justify-between mt-4 pt-4 border-t"
      >
        <p class="text-sm text-muted-foreground">
          Showing {{ startIndex }} to {{ endIndex }} of {{ total }} results
        </p>

        <div class="flex gap-2">
          <Button
            size="sm"
            variant="outline"
            :disabled="page === 1"
            @click="handlePageChange(page - 1)"
          >
            Previous
          </Button>

          <Button
            size="sm"
            variant="outline"
            :disabled="page >= totalPages"
            @click="handlePageChange(page + 1)"
          >
            Next
          </Button>
        </div>
      </div>
    </div>
  </Container>
</template>

<script setup>
useHead({
  title: 'Expense Audit Logs'
})

import { ref, computed, onMounted } from 'vue'
import { Container } from '~/components/ui/container'
import { Button } from '~/components/ui/button'
import { useApi } from '~/composables/useApi'
import { useHead } from 'nuxt/app'

const auditLogs = ref([])
const total = ref(0)
const page = ref(1)
const limit = ref(10)

const totalPages = computed(() => Math.ceil(total.value / limit.value))
const startIndex = computed(() => (page.value - 1) * limit.value + 1)
const endIndex = computed(() => Math.min(page.value * limit.value, total.value))

onMounted(fetchAuditLogs)

async function fetchAuditLogs() {
  const { get } = useApi('manager')

  const params = new URLSearchParams({
    page: page.value.toString(),
    limit: limit.value.toString(),
  })

  const res = await get(`/expense-logs?${params}`)
  auditLogs.value = res?.data || []
  total.value = res?.meta?.total || 0
}

const handlePageChange = (newPage) => {
  if (newPage < 1 || newPage > totalPages.value) return
  page.value = newPage
  fetchAuditLogs()
}

const getStatusBadgeClass = (status) => {
  const statusMap = {
    approved: 'bg-green-100 text-green-700',
    rejected: 'bg-red-100 text-red-700',
    pending: 'bg-yellow-100 text-yellow-700',
  }
  return statusMap[status] || 'bg-gray-100 text-gray-700'
}

const formatDateTime = (date) => {
  if (!date) return '-'
  const d = new Date(date)
  return d.toLocaleString('id-ID')
}
</script>
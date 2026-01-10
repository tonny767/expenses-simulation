<template>
  <Container class="p-6 space-y-6">
    <div>
      <h2 class="text-xl font-bold">Pending Approvals</h2>
      <p class="text-muted-foreground">
        Expenses waiting for your approval
      </p>
    </div>

    <div class="rounded-lg border p-4 space-y-4">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Description</TableHead>
            <TableHead>Amount</TableHead>
            <TableHead>Status</TableHead>
            <TableHead>Submitted</TableHead>
            <TableHead>User</TableHead>
            <TableHead class="text-right">Action</TableHead>
          </TableRow>
        </TableHeader>

        <TableBody>
          <TableRow v-for="expense in expenses" :key="expense.id">
            <TableCell>
              <NuxtLink
                :to="`/expenses/${expense.id}`"
                class="font-medium text-blue-600 hover:underline"
              >
                {{ expense.description }}
              </NuxtLink>
            </TableCell>
            <TableCell>Rp {{ formatAmount(expense.amount_idr) }}</TableCell>

            <!-- Status + hover notes -->
            <TableCell>
              <div class="flex items-center gap-2">
                <span
                  class="px-2 py-1 text-sm rounded bg-yellow-100 text-yellow-700"
                >
                  {{ expense.status.toUpperCase() }}
                </span>
              </div>
            </TableCell>

            <TableCell>{{ formatDateTime(expense.submitted_at) }}</TableCell>
            <TableCell>{{ expense.user?.name }}</TableCell>

            <TableCell class="text-right">
              <div class="flex justify-end gap-2">
                <ButtonAlert
                  :expense-id="expense.id"
                  title="Approve Expense"
                  message="Approve this expense?"
                  confirm-label="Approve"
                  action-type="approve"
                  :on-action-complete="fetchExpenses"
                >
                  <template #trigger>
                    <Button size="sm" variant="green">Approve</Button>
                  </template>
                </ButtonAlert>

                <ButtonAlert
                  :expense-id="expense.id"
                  title="Reject Expense"
                  message="Reject this expense?"
                  confirm-label="Reject"
                  action-type="reject"
                  :on-action-complete="fetchExpenses"
                >
                  <template #trigger>
                    <Button size="sm" variant="destructive">Reject</Button>
                  </template>
                </ButtonAlert>
              </div>
            </TableCell>
          </TableRow>
        </TableBody>
      </Table>

      <div v-if="!expenses.length" class="p-6 text-center text-muted-foreground">
        No pending approvals 🎉
      </div>

      <!-- Pagination (reused) -->
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
  title: 'Pending Approvals'
})

import { ref, computed, onMounted } from 'vue'
import Container from '~/components/ui/container/Container.vue'
import { Button } from '~/components/ui/button'
import { useApi } from '~/composables/useApi'
import ButtonAlert from '~/components/ButtonAlert.vue'
import { useHead } from 'nuxt/app'

const expenses = ref([])
const total = ref(0)
const page = ref(1)
const limit = ref(10)

const totalPages = computed(() => Math.ceil(total.value / limit.value))
const startIndex = computed(() => (page.value - 1) * limit.value + 1)
const endIndex = computed(() => Math.min(page.value * limit.value, total.value))

onMounted(fetchExpenses)

async function fetchExpenses() {
  const { get } = useApi('manager')

  const params = new URLSearchParams({
    page: page.value.toString(),
    limit: limit.value.toString(),
    status: 'pending',
  })

  const res = await get(`/expenses?${params}`)
  expenses.value = res?.data || []
  total.value = res?.meta?.total || 0
}

const handlePageChange = (newPage) => {
  if (newPage < 1 || newPage > totalPages.value) return
  page.value = newPage
  fetchExpenses()
}

const formatAmount = (amount) =>
  amount != null ? amount.toLocaleString('id-ID') : '0'

const formatDateTime = (date) => {
  if (!date) return '-'
  const d = new Date(date)
  return d.toLocaleString('id-ID')
}
</script>


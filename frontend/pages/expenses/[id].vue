<template>
  <Container class="p-6 space-y-6 max-w-3xl">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-xl font-bold">Expense Details</h2>
        <p class="text-sm text-muted-foreground">
          Expense ID #{{ expense?.id }}
        </p>
      </div>

      <Button variant="outline" @click="$router.back()">
        Back
      </Button>
    </div>

    <!-- Summary Card -->
    <div class="rounded-xl border p-6 space-y-4">
      <div>
        <p class="text-sm text-muted-foreground">Description</p>
        <p class="text-lg font-medium">
          {{ expense?.description || '-' }}
        </p>
      </div>

      <div class="grid grid-cols-2 gap-6">
        <div>
          <p class="text-sm text-muted-foreground">Amount</p>
          <p class="text-lg font-semibold">
            Rp {{ formatAmount(expense?.amount_idr) }}
          </p>
        </div>

        <div>
          <p class="text-sm text-muted-foreground">Submitted</p>
          <p>
            {{ formatDateTime(expense?.submitted_at) }}
          </p>
        </div>
      </div>

      <!-- Status + Notes -->
      <div class="flex items-center gap-2">
        <span
          class="px-3 py-1 rounded text-sm font-medium"
          :class="statusClass(expense?.status)"
        >
          {{ expense?.status?.toUpperCase() }}
        </span>

        <HoverCard v-if="expense?.approval?.notes">
          <HoverCardTrigger as-child>
            <button
              class="text-muted-foreground hover:text-foreground"
              aria-label="View notes"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M13 16h-1v-4h-1m1-4h.01M12 2a10 10 0 100 20 10 10 0 000-20z"
                />
              </svg>
            </button>
          </HoverCardTrigger>

          <HoverCardContent class="w-64">
            <p class="text-sm whitespace-pre-wrap">
              {{ expense.approval.notes }}
            </p>
          </HoverCardContent>
        </HoverCard>
      </div>

      <!-- Manager-only -->
      <div v-if="isManager">
        <p class="text-sm text-muted-foreground">Submitted by</p>
        <p>{{ expense?.user.name }}</p>
      </div>
    </div>

    <div class="rounded-xl border p-6 space-y-3">
      <h3 class="font-semibold">Receipt</h3>

      <div v-if="expense?.receipt_url">
        <NuxtLink
          v-if="isImage(expense.receipt_url)"
          :to="expense.receipt_url"
          target="_blank"
          class="inline-block mt-3 cursor-pointer" 
        >
          <img
            :src="expense.receipt_url"
            class="rounded-lg border max-h-100"
          />
        </NuxtLink>
      </div>

      <p v-else class="text-sm text-muted-foreground">
        No receipt uploaded
      </p>
    </div>

    <!-- Manager Actions -->
    <div
      v-if="canManage(expense)"
      class="flex justify-end gap-3"
    >
      <ButtonAlert
        :expense-id="expense.id"
        title="Approve Expense"
        message="Do you want to approve this expense?"
        confirm-label="Approve"
        cancel-label="Cancel"
        action-type="approve"
        :on-action-complete="fetchExpense"
      >
        <template #trigger>
          <Button variant="green">Approve</Button>
        </template>
      </ButtonAlert>

      <ButtonAlert
        :expense-id="expense.id"
        title="Reject Expense"
        message="Do you want to reject this expense?"
        confirm-label="Reject"
        cancel-label="Cancel"
        action-type="reject"
        :on-action-complete="fetchExpense"
      >
        <template #trigger>
          <Button variant="destructive">Reject</Button>
        </template>
      </ButtonAlert>
    </div>
  </Container>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import Container from '~/components/ui/container/Container.vue'
import { Button } from '~/components/ui/button'
import ButtonAlert from '~/components/ButtonAlert.vue'
import { useAuth } from '~/composables/useAuth'
import { useApi } from '~/composables/useApi'
import {
  HoverCard,
  HoverCardTrigger,
  HoverCardContent,
} from '~/components/ui/hover-card'

const route = useRoute()
const { role } = useAuth()
const expense = ref(null)
const loading = ref(true)

const isManager = computed(() => role.value === 'manager')

onMounted(fetchExpense)

async function fetchExpense() {
  try {
    loading.value = true
    const { get } = useApi(role.value)

    expense.value = await get(`/expenses/${route.params.id}`)
  } finally {
    loading.value = false
  }
}

const canManage = (expense) =>
  isManager.value && expense?.status === 'pending'

const formatAmount = (amount) =>
  amount != null ? amount.toLocaleString('id-ID') : '0'

const formatDateTime = (date) => {
  if (!date) return '-'
  const d = new Date(date)
  return d.toLocaleString()
}

const statusClass = (status) => {
  switch (status) {
    case 'approved':
      return 'bg-green-100 text-green-700'
    case 'rejected':
      return 'bg-red-100 text-red-700'
    case 'pending':
      return 'bg-yellow-100 text-yellow-700'
    default:
      return 'bg-gray-100 text-gray-700'
  }
}

const isImage = (url) =>
  /\.(jpg|jpeg|png|gif|webp)$/i.test(url)
</script>

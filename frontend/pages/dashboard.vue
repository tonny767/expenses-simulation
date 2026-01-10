<template>
  <Container class="p-6 space-y-6">
    <div>
      <h2 class="text-xl font-bold">Dashboard</h2>
      <p class="text-muted-foreground">
        Welcome back, {{ userName }}!
      </p>
    </div>

    <div class="rounded-lg border p-4 space-y-4">
      <div class="flex justify-between items-center mb-4">
        <h2 class="text-xl font-bold">Expenses</h2>
        <Button v-if="!isManager" size="md" @click="openModal = true">New Expense</Button>
      </div>

      <div class="flex items-center gap-3">
        <select
          v-model="statusFilter"
          @change="() => { page = 1; fetchExpenses()}"
          class="border rounded px-3 py-1 text-sm"
        >
          <option value="">All Status</option>
          <option value="pending">Pending</option>
          <option value="approved">Approved</option>
          <option value="rejected">Rejected</option>
          <option value="completed">Completed</option>
        </select>
      </div>
      <!-- New Expense Dialog -->
      <Dialog v-model:open="openModal">
        <DialogContent class="sm:max-w-lg">
          <DialogHeader>
            <DialogTitle>New Expense</DialogTitle>
            <DialogDescription>Fill out the details below</DialogDescription>
          </DialogHeader>

          <form @submit.prevent="submitExpense" class="space-y-4 mt-4">
            <!-- Description -->
            <div class="space-y-1">
              <Label for="description">
                Description <span class="text-red-500">*</span>
              </Label>
              <Input
                id="description"
                v-model="newExpense.description"
                type="text"
                placeholder="Description"
                required
              />
            </div>

            <!-- Amount -->
            <div class="space-y-1">
              <Label for="amount">
                Amount (IDR) <span class="text-red-500">*</span>
              </Label>
              <MoneyInput v-model="newExpense.amount_idr" />
              
              <!-- Validation Messages -->
              <div v-if="newExpense.amount_idr">
                <!-- Too Low -->
                <p v-if="newExpense.amount_idr < 10000" class="text-xs text-red-600 mt-1 flex items-center gap-1">
                  <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  Minimum amount is Rp 10.000
                </p>
                
                <!-- > 50mil -->
                <p v-else-if="newExpense.amount_idr > 50000000" class="text-xs text-red-600 mt-1 flex items-center gap-1">
                  <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  Maximum amount is Rp 50.000.000
                </p>
                
                <!-- >= 1mil - <= 50mil -->
                <p v-else-if="newExpense.amount_idr >= 1000000" class="text-xs text-amber-600 mt-1 flex items-center gap-1">
                  <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                  </svg>
                  This amount requires manager approval
                </p>
                
                <!-- >= 10k - < 1mil -->
                <p v-else class="text-xs text-green-600 mt-1 flex items-center gap-1">
                  <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  This expense will be auto-approved
                </p>
              </div>
              
              <p class="text-xs text-muted-foreground mt-1">
                Min: Rp 10.000 | Max: Rp 50.000.000
              </p>
            </div>

            <!-- Receipt Upload -->
            <div class="space-y-1">
              <Label for="receipt">
                Receipt <span class="text-red-500">*</span>
              </Label>
              <div class="space-y-2">
                <Input
                  id="receipt"
                  type="file"
                  accept="image/*,.pdf"
                  @change="handleFileUpload"
                  required
                  class="cursor-pointer"
                />
                <p class="text-xs text-muted-foreground">
                  Upload receipt image or PDF (max 5MB)
                </p>
                
                <!-- File Preview -->
                <div v-if="uploadedFile" class="flex items-center gap-2 p-2 bg-gray-50 rounded-md border">
                  <svg class="w-4 h-4 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <span class="text-sm text-gray-700">{{ uploadedFile.name }}</span>
                  <span class="text-xs text-gray-500">({{ formatFileSize(uploadedFile.size) }})</span>
                  <button
                    type="button"
                    @click="removeFile"
                    class="ml-auto text-red-500 hover:text-red-700"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                  </button>
                </div>
              </div>
            </div>

            <!-- Buttons -->
            <div class="flex justify-end gap-2 mt-6">
              <Button 
                type="submit" 
                variant="green"
                :disabled="!isAmountValid"
              >
                Submit
              </Button>
              <Button
                type="button"
                variant="destructive"
                @click="closeModal"
              >
                Cancel
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>

      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Description</TableHead>
            <TableHead>Amount</TableHead>
            <TableHead>Status</TableHead>
            <TableHead>Submitted</TableHead>
            <TableHead>Updated At</TableHead>

            <TableHead v-if="isManager">User</TableHead>
            <TableHead v-if="isManager" class="text-right">Action</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow
            v-for="expense in expenses"
            :key="expense.id"
          >
          <TableCell>
            <NuxtLink
              :to="`/expenses/${expense.id}`"
              class="font-medium text-blue-600 hover:underline"
            >
              {{ expense.description }}
            </NuxtLink>
          </TableCell>            
            <TableCell>Rp {{ formatAmount(expense.amount_idr) }}</TableCell>
            <TableCell>
              <div class="flex items-center gap-2">

                <span
                  class="px-2 py-1 text-sm rounded"
                  :class="statusClass(expense.status)"
                >
                  {{ expense.status.toUpperCase() }}
                </span>

                <!-- Notes hover card -->
                <HoverCard v-if="expense.approval?.notes">
                  <HoverCardTrigger as-child>
                    <button
                      class="text-muted-foreground hover:text-foreground transition"
                      aria-label="View notes"
                    >
                      <svg
                        class="w-4 h-4"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
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
                    <div class="space-y-1">
                      <p class="text-sm text-muted-foreground whitespace-pre-wrap">
                        {{ expense.approval.notes }}
                      </p>
                    </div>
                  </HoverCardContent>
                </HoverCard>     
              </div>       
            </TableCell>
            <TableCell>{{ formatDateTime(expense.submitted_at) }}</TableCell>
            <TableCell>{{ formatDateTime(expense.updated_at) }}</TableCell>
            <TableCell v-if="isManager">{{ expense.user?.name || '-' }}</TableCell>
            <TableCell v-if="isManager" class="text-right">
              <div v-if="canManage(expense)" class="flex justify-end gap-2">
                <ButtonAlert
                  :expense-id="expense.id"
                  title="Approve Expense"
                  message="Do you want to approve this expense?"
                  confirm-label="Approve"
                  cancel-label="Cancel"
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
                  message="Do you want to reject this expense?"
                  confirm-label="Reject"
                  cancel-label="Cancel"
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
        No expenses found
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

          <!-- Page Numbers -->
          <div class="flex gap-1">
            <Button
              v-for="pageNum in visiblePages"
              :key="pageNum"
              size="sm"
              :variant="pageNum === page ? 'default' : 'outline'"
              @click="handlePageChange(pageNum)"
            >
              {{ pageNum }}
            </Button>
          </div>

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
  title: 'Dashboard'
})

import { ref, computed } from 'vue'
import Container from '~/components/ui/container/Container.vue'
import { Button } from '~/components/ui/button'
import { useAuth } from '~/composables/useAuth'
import { useApi } from '~/composables/useApi'
import ButtonAlert from '~/components/ButtonAlert.vue'
import { Label } from '~/components/ui/label'
import { Input } from '~/components/ui/input'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from '~/components/ui/dialog'
import MoneyInput from '~/components/MoneyInput.vue'
import {
  HoverCard,
  HoverCardTrigger,
  HoverCardContent,
} from '~/components/ui/hover-card'
import { useHead } from 'nuxt/app'

const { userName, role, userId } = useAuth()

const expenses = ref([])
const loading = ref(true)
const error = ref(null)
const openModal = ref(false)
const uploadedFile = ref(null)
const total = ref(0)
const page = ref(1)
const limit = ref(10)
const statusFilter = ref('')

const newExpense = ref({
  description: '',
  amount_idr: null,
  receipt_url: '/receipt-placeholder.png', // only fake mock image
})

const isManager = ref(role.value === 'manager')

// pagination computed
const totalPages = computed(() => Math.ceil(total.value / limit.value))

const startIndex = computed(() => {
  if (total.value === 0) return 0
  return (page.value - 1) * limit.value + 1
})

const endIndex = computed(() => {
  const end = page.value * limit.value
  return end > total.value ? total.value : end
})

const visiblePages = computed(() => {
  const pages = []
  const maxVisible = 5
  let start = Math.max(1, page.value - Math.floor(maxVisible / 2))
  let end = Math.min(totalPages.value, start + maxVisible - 1)
  
  if (end - start < maxVisible - 1) {
    start = Math.max(1, end - maxVisible + 1)
  }
  
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  
  return pages
})

// Validation for amount
const isAmountValid = computed(() => {
  const amount = newExpense.value.amount_idr
  if (!amount) return false
  return amount >= 10000 && amount <= 50000000
})

onMounted(async () => {
  if (!role?.value) return
  await fetchExpenses()
})

async function fetchExpenses() {
  if (!role?.value) {
    loading.value = false
    return
  }

  try {
    loading.value = true
    error.value = null

    const { get } = useApi(role.value)

    const params = new URLSearchParams({
      page: page.value.toString(),
      limit: limit.value.toString(),
    })
    
    if (statusFilter.value) {
      params.append('status', statusFilter.value)
    }

    const path = `/expenses?${params.toString()}`
    const res = await get(path)

    expenses.value = res?.data || []
    total.value = res?.meta?.total || 0
    page.value = res?.meta?.page || 1
    limit.value = res?.meta?.limit || 10

  } catch (err) {
    error.value = err?.error || 'Failed to load expenses'
    expenses.value = []
    total.value = 0
  } finally {
    loading.value = false
  }
}

const handlePageChange = (newPage) => {
  if (newPage < 1 || newPage > totalPages.value) return
  page.value = newPage
  fetchExpenses()
}

const handleFileUpload = (event) => {
  const file = event.target.files[0]

  if (!file) return

  const maxSize = 5 * 1024 * 1024
  if (file.size > maxSize) {
    alert('File size must be less than 5MB')
    event.target.value = ''
    return
  }

  const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'application/pdf']
  if (!allowedTypes.includes(file.type)) {
    alert('Only images (JPG, PNG, GIF) and PDF files are allowed')
    event.target.value = ''
    return
  }

  uploadedFile.value = file
}

const removeFile = () => {
  uploadedFile.value = null
  const fileInput = document.getElementById('receipt')
  if (fileInput) fileInput.value = ''
}

const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

const closeModal = () => {
  openModal.value = false
  newExpense.value.description = ''
  newExpense.value.amount_idr = null
  uploadedFile.value = null
}

const submitExpense = async () => {
  // Validate amount
  if (!newExpense.value.amount_idr) {
    alert('Amount is required')
    return
  }

  if (newExpense.value.amount_idr < 10000) {
    alert('Minimum amount is Rp 10.000')
    return
  }

  if (newExpense.value.amount_idr > 50000000) {
    alert('Maximum amount is Rp 50.000.000')
    return
  }

  if (!newExpense.value.description) {
    alert('Description is required')
    return
  }

  if (!uploadedFile.value) {
    alert('Receipt is required')
    return
  }

  try {
    const { post } = useApi(role.value)
    const path = '/expenses'

    const res =await post(path, {
      user_id: userId.value,
      description: newExpense.value.description,
      amount_idr: newExpense.value.amount_idr,
      receipt_url: '/receipt-placeholder.png' // only fake mock image
    })

    alert(res?.message)

    closeModal()
    await fetchExpenses()
  } catch (err) {
    alert(res?.error || 'Failed to submit expense')
  }
}

const canManage = (expense) => {
  return (
    isManager &&
    expense?.status === 'pending'
  )
}

const formatAmount = (amount) => {
  return amount != null ? amount.toLocaleString('id-ID') : '0'
}

const formatDateTime = (date) => {
  if (!date) return '-'
  try {
    const d = new Date(date).toLocaleDateString('en-GB', { 
      timeZone: 'Asia/Jakarta',
      day: '2-digit',
      month: '2-digit',
      year: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: false
    })
    
    return d
  } catch {
    return 'Invalid date'
  }
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
</script>
<template>
    <Card class="p-6">
      <div class="flex justify-between mb-4">
        <h2 class="text-xl font-bold">Expenses</h2>
        <Button @click="openModal">New Expense</Button>
      </div>

      <Table class="w-full">
        <thead>
          <tr>
            <th>ID</th>
            <th>Description</th>
            <th>Amount</th>
            <th>Status</th>
            <th>Submitted</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="exp in expenses" :key="exp.id">
            <td>{{ exp.id }}</td>
            <td>{{ exp.description }}</td>
            <td>{{ exp.amountIDR }}</td>
            <td>{{ exp.status }}</td>
            <td>{{ exp.submittedAt }}</td>
          </tr>
        </tbody>
      </Table>
    </Card>

    <Popover open="showModal">
      <template #header>New Expense</template>
      <template #body>
        <form @submit.prevent="submitExpense">
          <Label>Description</Label>
          <Input v-model="desc" placeholder="Expense description" class="mb-4" />

          <Label>Amount (IDR)</Label>
          <Input v-model="amount" type="number" placeholder="100000" class="mb-4" />

          <Label>Receipt URL</Label>
          <Input v-model="receipt" placeholder="https://..." class="mb-4" />
        </form>
      </template>
      <template #footer>
        <Button @click="submitExpense">Submit</Button>
        <Button variant="secondary" @click="showModal = false">Cancel</Button>
      </template>
    </Popover>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useApi } from '@/composables/useApi'

const api = useApi()
const expenses = ref([])
const showModal = ref(false)

const desc = ref('')
const amount = ref('')
const receipt = ref('')

const fetchExpenses = async () => {
  expenses.value = await api.get('/user/expenses')
}

const submitExpense = async () => {
  const newExp = await api.post('/user/expenses', {
    description: desc.value,
    amountIDR: parseInt(amount.value),
    receiptURL: receipt.value,
  })
  expenses.value.push(newExp.expense)
  showModal.value = false
}

onMounted(() => fetchExpenses())
</script>

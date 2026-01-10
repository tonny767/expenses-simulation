<template>
  <AlertDialog v-model:open="open">
    <AlertDialogTrigger asChild>
      <slot name="trigger"></slot>
    </AlertDialogTrigger>

    <AlertDialogContent class="sm:max-w-lg">
      <AlertDialogHeader>
        <AlertDialogTitle>{{ title }}</AlertDialogTitle>
        <AlertDialogDescription class="mt-2 text-sm text-muted-foreground">
          {{ message }}
        </AlertDialogDescription>
      </AlertDialogHeader>

      <div class="mt-4">
        <label class="block text-sm font-medium">Notes (optional)</label>
        <Textarea
          v-model="notes"
          :rows="5"
          class="mt-1 w-full border rounded p-2 text-sm"
          placeholder="Add a note..."
        />
      </div>

      <AlertDialogFooter class="mt-4 flex justify-end gap-2">
        <Button
          size="sm"
          variant="outline"
          @click="handleCancel"
        >
          {{ cancelLabel }}
        </Button>
        <Button
          size="sm"

          :variant="actionType == 'approve' ? 'green' : 'destructive'"
          @click="handleConfirm"
        >
          {{ confirmLabel }}
        </Button>
      </AlertDialogFooter>
    </AlertDialogContent>
  </AlertDialog>
</template>

<script setup>
import { ref } from 'vue'
import { AlertDialog, AlertDialogTrigger, AlertDialogContent, AlertDialogHeader, AlertDialogTitle, AlertDialogDescription, AlertDialogFooter } from '~/components/ui/alert-dialog'
import { Button } from '~/components/ui/button'
import { useApi } from '~/composables/useApi'
import { useAuth } from '~/composables/useAuth'
import { Textarea } from '~/components/ui/textarea'

const props = defineProps({
  expenseId: { type: Number, required: true },
  title: { type: String, default: 'Confirm Action' },
  message: { type: String, default: 'Are you sure?' },
  confirmLabel: { type: String, default: 'Approve' },
  cancelLabel: { type: String, default: 'Cancel' },
  actionType: { type: String, default: 'approve' }, // approve | reject
  onActionComplete: { type: Function }
})

const { role, userId } = useAuth()
const api = useApi(role.value)
const open = ref(false)
const notes = ref('')

const handleConfirm = async () => {
  try {
    await api.put(`/expenses/${props.expenseId}/${props.actionType}`, {
      approver_id: userId.value,
      notes: notes.value
    })

    props.actionType === 'approve'
      ? alert('Expense has been approved!')
      : alert('Expense has been rejected!')

    notes.value = ''
    open.value = false
    props.onActionComplete?.()
  } catch (err) {
    alert(`Failed to ${props.actionType} expense:`, err)
  }
}

const handleCancel = () => {
  notes.value = ''
  open.value = false
}
</script>

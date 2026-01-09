<template>
  <div class="relative flex items-center">
    <span class="absolute left-3 text-gray-500 font-medium text-sm pointer-events-none">
      Rp
    </span>
    <Input
      v-model="localValue"
      @input="handleInput"
      @blur="handleBlur"
      placeholder="0"
      class="pl-9"
    />
  </div>

</template>

<script setup>
import { ref, watch } from "vue"
import { Input } from "~/components/ui/input"

const props = defineProps({
  modelValue: {
    type: [Number, String],
    default: 0,
  },
})

const emit = defineEmits(["update:modelValue"])

const localValue = ref("")
const isTyping = ref(false)

const formatNumber = (val) => {
  if (!val && val !== 0) return ""
  const n = Number(val) || 0
  return n.toLocaleString("id-ID")
}

// Only update from props when NOT typing
watch(
  () => props.modelValue,
  (val) => {
    if (!isTyping.value) {
      localValue.value = formatNumber(val)
    }
  },
  { immediate: true }
)

const handleInput = (event) => {
  isTyping.value = true
  const raw = event.target.value.replace(/\D/g, "")
  localValue.value = raw ? formatNumber(raw) : ""
  emit("update:modelValue", Number(raw) || 0)
}

const handleBlur = () => {
  isTyping.value = false
  // Re-format on blur
  localValue.value = formatNumber(props.modelValue)
}
</script>
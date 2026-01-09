<script setup>
defineProps({
  class: {
    type: String,
    default: '',
  },
})
</script>

<template>
  <div :class="['mx-auto max-w-7xl px-4', $props.class]">
    <slot />
  </div>
</template>
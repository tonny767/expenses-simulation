<!-- pages/index.vue -->
<template>
  <div class="flex items-center justify-center min-h-screen">
    <div class="text-center">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-gray-900 mx-auto"></div>
      <p class="mt-4 text-gray-600">Redirecting...</p>
    </div>
  </div>
</template>

<script setup>
const { isLoggedIn } = useAuth()

// Redirect based on auth status
if (isLoggedIn.value) {
  navigateTo('/dashboard')
} else {
  navigateTo('/login')
}
</script>
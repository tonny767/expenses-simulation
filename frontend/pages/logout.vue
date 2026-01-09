<template>
  <div class="flex items-center justify-center min-h-screen">
    <Card class="p-6">
      <p class="text-lg">Logging out...</p>
    </Card>
  </div>
</template>

<script setup>
const { logout } = useAuth()

// Run logout immediately when page loads
onMounted(async () => {
  await logout() // This will clear cookies and redirect to /login
})
</script>
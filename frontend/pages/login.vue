<template>
  <Card class="max-w-md mx-auto mt-20 p-6">
    <h2 class="text-2xl font-bold mb-6 text-center">Login</h2>

    <form @submit.prevent="submit">
      <div class="mb-4">
        <Label for="email">Email</Label>
        <Input 
          id="email" 
          v-model="email" 
          type="email" 
          placeholder="you@example.com" 
          required 
        />
      </div>

      <div class="mb-4">
        <Label for="password">Password</Label>
        <Input 
          id="password" 
          v-model="password" 
          type="password" 
          placeholder="********" 
          required 
        />
      </div>

      <Button 
        class="w-full mt-4 cursor-pointer" 
        :disabled="loading"
      >
        {{ loading ? 'Logging in...' : 'Login' }}
      </Button>

      <p v-if="error" class="text-red-500 mt-2 text-sm">{{ error }}</p>
    </form>
  </Card>
</template>

<script setup>
const { login } = useAuth()
const api = useApi()

const email = ref('') 
const password = ref('')
const loading = ref(false)
const error = ref('')

const submit = async () => {
  loading.value = true
  error.value = ''
  
  try {
    const res = await api.post('/auth/login', { 
      email: email.value, 
      password: password.value 
    })
    
    login(res.role, res.name, res.id)
    
    // Redirect based on role
    await navigateTo('/dashboard')
} catch (err) {
    error.value = err?.data?.message || err?.message || 'Login failed'
  } finally {
    loading.value = false
  }
}
</script>
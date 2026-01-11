<template>
  <nav class="bg-white shadow-md">
    <div v-if="isLoggedIn" class="max-w-7xl mx-auto px-4 py-4">
      <div class="flex justify-between items-center">
        <NuxtLink to="/" class="text-xl font-bold text-gray-800">
          Expense Management
        </NuxtLink>

        <!-- Right side -->
        <div  class="flex items-center gap-6">
          <NuxtLink 
            v-if="isManager" 
            to="/manager/logs"
            class="text-gray-700 hover:text-gray-900 font-medium"
          >
            Expense Logs
          </NuxtLink>

          <NuxtLink 
            v-if="isManager" 
            to="/manager/approvals"
            class="text-gray-700 hover:text-gray-900 font-medium"
          >
            Pending Approvals
          </NuxtLink>

          <div class="flex items-center gap-4 border-l pl-6">
            <div class="text-sm">
              <div class="font-medium text-gray-900">{{ userName }}</div>
              <div class="text-gray-500">{{ role }}</div>
            </div>
            
            <Button
              variant="destructive"
              @click="logout"
              class="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600 text-sm"
            >
              Logout
          </Button>
          </div>
        </div>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { useAuth } from '~/composables/useAuth';
const { isLoggedIn, userName, role, isManager, logout } = useAuth()
</script>
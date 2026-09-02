<script setup lang="ts">
import { ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { login } from '@/lib/auth'
import { Button } from '@/components/ui/button'

const username = ref('armin')
const password = ref('')
const error = ref<string | null>(null)
const loading = ref(false)
const router = useRouter()
const route = useRoute()

async function onSubmit() {
  error.value = null
  loading.value = true
  try {
    login(username.value, password.value)
    const redirect = typeof route.query.redirect === 'string' ? route.query.redirect : '/dashboard'
    await router.replace(redirect)
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Login failed'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="flex min-h-full items-center justify-center bg-muted/40 p-6">
    <form class="w-full max-w-sm space-y-4 rounded-lg border bg-card p-6 shadow-sm" @submit.prevent="onSubmit">
      <div>
        <h1 class="text-xl font-semibold">Susiana</h1>
        <p class="text-sm text-muted-foreground">Sign in to continue</p>
      </div>
      <label class="block space-y-1 text-sm">
        <span>Username</span>
        <input
          v-model="username"
          autocomplete="username"
          class="h-9 w-full rounded-md border border-input bg-background px-3 outline-none focus-visible:ring-2 focus-visible:ring-ring"
        />
      </label>
      <label class="block space-y-1 text-sm">
        <span>Password</span>
        <input
          v-model="password"
          type="password"
          autocomplete="current-password"
          class="h-9 w-full rounded-md border border-input bg-background px-3 outline-none focus-visible:ring-2 focus-visible:ring-ring"
        />
      </label>
      <p v-if="error" class="text-sm text-ask">{{ error }}</p>
      <Button class="w-full" type="submit" :disabled="loading">
        {{ loading ? 'Signing in…' : 'Sign in' }}
      </Button>
    </form>
  </div>
</template>

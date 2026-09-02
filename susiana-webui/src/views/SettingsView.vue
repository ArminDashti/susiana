<script setup lang="ts">
import { onMounted, ref } from 'vue'

const apiBase = ref(localStorage.getItem('susiana_api_base') ?? '')
const theme = ref(localStorage.getItem('susiana_theme') ?? 'system')
const locale = ref(localStorage.getItem('susiana_locale') ?? 'en')
const saved = ref(false)

function applyTheme(value: string) {
  const root = document.documentElement
  if (value === 'dark') root.classList.add('dark')
  else if (value === 'light') root.classList.remove('dark')
  else if (window.matchMedia('(prefers-color-scheme: dark)').matches) root.classList.add('dark')
  else root.classList.remove('dark')
}

onMounted(() => applyTheme(theme.value))

function save() {
  localStorage.setItem('susiana_api_base', apiBase.value)
  localStorage.setItem('susiana_theme', theme.value)
  localStorage.setItem('susiana_locale', locale.value)
  applyTheme(theme.value)
  saved.value = true
  window.setTimeout(() => {
    saved.value = false
  }, 1500)
}
</script>

<template>
  <form class="max-w-lg space-y-4" @submit.prevent="save">
    <label class="block space-y-1 text-sm">
      <span>api_base_url (optional override; leave empty to use Vite proxy)</span>
      <input
        v-model="apiBase"
        placeholder="http://127.0.0.1:8080"
        class="h-9 w-full rounded-md border border-input bg-background px-3 font-mono text-sm outline-none focus-visible:ring-2 focus-visible:ring-ring"
      />
    </label>
    <label class="block space-y-1 text-sm">
      <span>theme</span>
      <select
        v-model="theme"
        class="h-9 w-full rounded-md border border-input bg-background px-3 text-sm outline-none focus-visible:ring-2 focus-visible:ring-ring"
      >
        <option value="system">system</option>
        <option value="light">light</option>
        <option value="dark">dark</option>
      </select>
    </label>
    <label class="block space-y-1 text-sm">
      <span>locale</span>
      <select
        v-model="locale"
        class="h-9 w-full rounded-md border border-input bg-background px-3 text-sm outline-none focus-visible:ring-2 focus-visible:ring-ring"
      >
        <option value="en">en</option>
        <option value="fa">fa</option>
      </select>
    </label>
    <button
      type="submit"
      class="inline-flex h-9 items-center rounded-md bg-primary px-4 text-sm font-medium text-primary-foreground hover:bg-primary/90"
    >
      Save
    </button>
    <p v-if="saved" class="text-sm text-primary">Saved.</p>
  </form>
</template>

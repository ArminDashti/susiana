<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter, RouterLink } from 'vue-router'
import { LogOut, Settings } from 'lucide-vue-next'
import { getStoredUser, logout } from '@/lib/auth'
import { Button } from '@/components/ui/button'

type NavItem = { label: string; to: string }
type NavGroup = { label: string; items: NavItem[] }

const groups: NavGroup[] = [
  {
    label: 'Overview',
    items: [
      { label: 'Dashboard', to: '/dashboard' },
      { label: 'Settings', to: '/settings' },
    ],
  },
  {
    label: 'Stock market',
    items: [
      { label: 'TSE', to: '/stock-market/tse' },
      { label: 'TSE indexes', to: '/stock-market/tse/indexes' },
    ],
  },
  {
    label: 'Commodity',
    items: [
      { label: 'Crude oil', to: '/commodity/crude-oil' },
      { label: 'Natural gas', to: '/commodity/natural-gas' },
      { label: 'Gold', to: '/commodity/gold' },
      { label: 'Silver', to: '/commodity/silver' },
      { label: 'Copper', to: '/commodity/copper' },
      { label: 'Steel', to: '/commodity/steel' },
    ],
  },
  {
    label: 'Indexes',
    items: [
      { label: 'Asia', to: '/indexes/asia' },
      { label: 'Top 30', to: '/indexes/top30' },
      { label: 'Europe', to: '/indexes/europe' },
      { label: 'North America', to: '/indexes/north-america' },
      { label: 'South America', to: '/indexes/south-america' },
      { label: 'Middle East', to: '/indexes/middle-east' },
      { label: 'Australia', to: '/indexes/australia' },
      { label: 'Africa', to: '/indexes/africa' },
    ],
  },
  {
    label: 'FX & crypto',
    items: [
      { label: 'Currencies', to: '/currencies' },
      { label: 'Crypto', to: '/crypto-currencies' },
    ],
  },
]

const route = useRoute()
const router = useRouter()
const user = computed(() => getStoredUser())

function onLogout() {
  logout()
  router.push({ name: 'login' })
}

function isActive(to: string): boolean {
  return route.path === to || route.path.startsWith(`${to}/`)
}
</script>

<template>
  <div class="flex h-full min-h-0">
    <aside class="flex w-60 shrink-0 flex-col border-r bg-card">
      <div class="border-b px-4 py-4">
        <RouterLink to="/dashboard" class="text-lg font-semibold tracking-tight">Susiana</RouterLink>
        <p class="text-xs text-muted-foreground">Market command center</p>
      </div>
      <nav class="flex-1 overflow-y-auto px-2 py-3">
        <div v-for="group in groups" :key="group.label" class="mb-4">
          <p class="px-2 pb-1 text-[11px] font-semibold uppercase tracking-wide text-muted-foreground">
            {{ group.label }}
          </p>
          <RouterLink
            v-for="item in group.items"
            :key="item.to"
            :to="item.to"
            class="mb-0.5 block rounded-md px-2 py-1.5 text-sm transition-colors"
            :class="
              isActive(item.to)
                ? 'bg-primary/10 font-medium text-primary'
                : 'text-foreground/80 hover:bg-accent'
            "
          >
            {{ item.label }}
          </RouterLink>
        </div>
      </nav>
      <div class="flex items-center justify-between gap-2 border-t px-3 py-3">
        <div class="min-w-0">
          <p class="truncate text-sm font-medium">{{ user?.display_name ?? 'User' }}</p>
          <p class="truncate text-xs text-muted-foreground">{{ user?.username }}</p>
        </div>
        <div class="flex gap-1">
          <Button variant="ghost" size="sm" class="px-2" @click="router.push('/settings')">
            <Settings class="h-4 w-4" />
          </Button>
          <Button variant="ghost" size="sm" class="px-2" @click="onLogout">
            <LogOut class="h-4 w-4" />
          </Button>
        </div>
      </div>
    </aside>
    <main class="min-w-0 flex-1 overflow-auto">
      <header class="sticky top-0 z-10 border-b bg-background/90 px-6 py-3 backdrop-blur">
        <h1 class="text-base font-semibold">{{ route.meta.title ?? 'Susiana' }}</h1>
      </header>
      <div class="p-6">
        <slot />
      </div>
    </main>
  </div>
</template>

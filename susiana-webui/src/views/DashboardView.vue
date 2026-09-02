<script setup lang="ts">
import { ref } from 'vue'
import DataGrid from '@/components/DataGrid.vue'
import { INDEX_COLUMNS, TSE_LIST_COLUMNS } from '@/mocks/columns'
import { mockIndexes, mockTseStocks } from '@/mocks/market'
import { Button } from '@/components/ui/button'

const movers = mockTseStocks
const asia = mockIndexes.asia.slice(0, 5)
const apiHint = ref(import.meta.env.VITE_API_PROXY_TARGET || 'http://127.0.0.1:8080')
</script>

<template>
  <div class="space-y-6">
    <div class="grid gap-3 sm:grid-cols-3">
      <div class="rounded-md border bg-card p-4">
        <p class="text-xs uppercase text-muted-foreground">API proxy</p>
        <p class="mt-1 font-mono text-sm">{{ apiHint }}</p>
      </div>
      <div class="rounded-md border bg-card p-4">
        <p class="text-xs uppercase text-muted-foreground">TSE symbols (mock/API)</p>
        <p class="mt-1 text-2xl font-semibold tabular-nums">{{ movers.length }}+</p>
      </div>
      <div class="rounded-md border bg-card p-4">
        <p class="text-xs uppercase text-muted-foreground">Quick open</p>
        <Button class="mt-2" size="sm" @click="$router.push('/stock-market/tse')">Open TSE grid</Button>
      </div>
    </div>

    <section>
      <h2 class="mb-2 text-sm font-semibold">Watchlist</h2>
      <DataGrid :columns="TSE_LIST_COLUMNS" :rows="movers" row-key="ticker" @row-click="(r) => $router.push(`/stock-market/tse/stock/${r.ticker}`)" />
    </section>

    <section>
      <h2 class="mb-2 text-sm font-semibold">Asia indexes</h2>
      <DataGrid :columns="INDEX_COLUMNS" :rows="asia" />
    </section>
  </div>
</template>

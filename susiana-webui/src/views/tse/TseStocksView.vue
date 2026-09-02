<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import DataGrid from '@/components/DataGrid.vue'
import { TSE_LIST_COLUMNS } from '@/mocks/columns'
import type { MarketRow } from '@/mocks/columns'
import { mockTseStocks } from '@/mocks/market'
import { listStocks } from '@/api/stocks'

const rows = ref<MarketRow[]>([...mockTseStocks])
const loading = ref(false)
const error = ref<string | null>(null)
const router = useRouter()

onMounted(async () => {
  loading.value = true
  error.value = null
  try {
    const res = await listStocks({ page: 1, limit: 100 })
    if (res.data?.length) {
      rows.value = res.data.map((s) => ({
        ticker: s.ticker || s.local_ticker || s.isin,
        isin: s.isin,
        last: s.last ?? null,
        close: s.close ?? null,
        open: s.open ?? null,
        volume: s.volume ?? null,
        turnover: s.turnover ?? null,
        trades: s.trades ?? null,
        pe_ratio: s.pe_ratio ?? null,
        eps: s.eps ?? null,
      }))
    }
  } catch (e) {
    error.value = e instanceof Error ? `${e.message} — showing mock rows` : 'API unavailable — showing mock rows'
  } finally {
    loading.value = false
  }
})

function onRowClick(row: MarketRow) {
  const symbol = String(row.ticker ?? row.isin ?? '')
  if (symbol) router.push(`/stock-market/tse/stock/${encodeURIComponent(symbol)}`)
}
</script>

<template>
  <div class="space-y-3">
    <p class="text-sm text-muted-foreground">
      Click a row to open the stock dashboard. History:
      <code class="rounded bg-muted px-1">/stock-market/tse/stock/history/&lt;symbol&gt;</code>
    </p>
    <DataGrid
      :columns="TSE_LIST_COLUMNS"
      :rows="rows"
      :loading="loading"
      :error="error"
      row-key="ticker"
      @row-click="onRowClick"
    />
  </div>
</template>

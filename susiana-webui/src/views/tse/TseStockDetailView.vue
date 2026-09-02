<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import StockDetailDashboard from '@/components/StockDetailDashboard.vue'
import { mockStockDetail, type StockDetailMock } from '@/mocks/market'
import { findStockBySymbol } from '@/api/stocks'
import { Button } from '@/components/ui/button'

const route = useRoute()
const symbol = computed(() => String(route.params.symbol ?? ''))
const detail = ref<StockDetailMock>(mockStockDetail(symbol.value))
const source = ref<'api' | 'mock'>('mock')
const loading = ref(false)
const error = ref<string | null>(null)

async function load() {
  loading.value = true
  error.value = null
  const base = mockStockDetail(symbol.value)
  try {
    const stock = await findStockBySymbol(symbol.value)
    if (stock) {
      const last = stock.last ?? stock.close ?? base.last
      const close = stock.close ?? last
      const open = stock.open ?? last
      const prev = close * 0.97
      detail.value = {
        ...base,
        company_name: stock.security_name || base.company_name,
        symbol: stock.ticker || stock.local_ticker || symbol.value,
        market_segment: stock.market_segment || base.market_segment,
        trading_status: stock.trading_status || base.trading_status,
        last: last ?? base.last,
        close: close ?? base.close,
        open: open ?? base.open,
        prev_close: prev,
        change: (last ?? 0) - prev,
        change_pct: prev ? (((last ?? 0) - prev) / prev) * 100 : 0,
        trades: stock.trades ?? base.trades,
        volume: stock.volume ?? base.volume,
        turnover: stock.turnover ?? base.turnover,
        lower_limit: stock.lower_limit ?? base.lower_limit,
        upper_limit: stock.upper_limit ?? base.upper_limit,
        shares_outstanding: stock.shares_outstanding ?? base.shares_outstanding,
        base_volume: stock.base_volume ?? base.base_volume,
        free_float: stock.free_float ?? base.free_float,
        eps: stock.eps ?? base.eps,
        pe_ratio: stock.pe_ratio ?? base.pe_ratio,
        retail_buy_volume: stock.retail_buy_volume ?? base.retail_buy_volume,
        institutional_buy_volume: stock.institutional_buy_volume ?? base.institutional_buy_volume,
        retail_buy_count: stock.retail_buy_count ?? base.retail_buy_count,
        institutional_buy_count: stock.institutional_buy_count ?? base.institutional_buy_count,
      }
      source.value = 'api'
    } else {
      detail.value = base
      source.value = 'mock'
      error.value = 'Stock not found in API — showing mock dashboard'
    }
  } catch (e) {
    detail.value = base
    source.value = 'mock'
    error.value = e instanceof Error ? `${e.message} — showing mock dashboard` : 'API error'
  } finally {
    loading.value = false
  }
}

onMounted(load)
watch(symbol, load)
</script>

<template>
  <div class="space-y-3">
    <div class="flex flex-wrap items-center gap-2 text-sm">
      <RouterLink class="text-primary hover:underline" to="/stock-market/tse">← TSE list</RouterLink>
      <span class="text-muted-foreground">·</span>
      <RouterLink
        class="text-primary hover:underline"
        :to="`/stock-market/tse/stock/history/${encodeURIComponent(symbol)}`"
      >
        History
      </RouterLink>
      <span class="rounded bg-muted px-2 py-0.5 text-xs">source: {{ source }}</span>
      <Button size="sm" variant="outline" :disabled="loading" @click="load">Refresh</Button>
    </div>
    <p v-if="error" class="text-sm text-muted-foreground">{{ error }}</p>
    <StockDetailDashboard :detail="detail" />
  </div>
</template>

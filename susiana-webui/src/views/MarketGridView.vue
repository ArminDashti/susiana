<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import DataGrid from '@/components/DataGrid.vue'
import type { GridColumn, MarketRow } from '@/mocks/columns'
import {
  COMMODITY_COLUMNS,
  CRYPTO_COLUMNS,
  FX_COLUMNS,
  HISTORY_COLUMNS,
  INDEX_COLUMNS,
} from '@/mocks/columns'
import {
  mockCommodities,
  mockCrypto,
  mockCurrencies,
  mockHistory,
  mockIndexes,
} from '@/mocks/market'

const route = useRoute()
const router = useRouter()

const grid = computed(() => {
  const kind = String(route.meta.gridKind ?? '')
  const key = String(route.meta.gridKey ?? route.params.symbol ?? '')

  let columns: GridColumn[] = INDEX_COLUMNS
  let rows: MarketRow[] = []
  let rowKey = 'symbol'

  switch (kind) {
    case 'commodity':
      columns = COMMODITY_COLUMNS
      rows = mockCommodities[key] ?? []
      break
    case 'index':
      columns = INDEX_COLUMNS
      rows = mockIndexes[key as keyof typeof mockIndexes] ?? []
      break
    case 'tse-indexes':
      columns = INDEX_COLUMNS
      rows = mockIndexes.tse
      break
    case 'currencies':
      columns = FX_COLUMNS
      rows = mockCurrencies
      rowKey = 'pair'
      break
    case 'crypto':
      columns = CRYPTO_COLUMNS
      rows = mockCrypto
      break
    case 'history':
      columns = HISTORY_COLUMNS
      rows = mockHistory(String(route.params.symbol ?? ''))
      rowKey = 'as_of'
      break
    default:
      rows = []
  }

  return { columns, rows, rowKey }
})

function onRowClick(row: MarketRow) {
  if (route.meta.gridKind === 'history') return
  const symbol = row.symbol ?? row.pair
  if (typeof symbol === 'string' && route.path.includes('/stock-market/tse')) {
    router.push(`/stock-market/tse/stock/${symbol}`)
  }
}
</script>

<template>
  <DataGrid
    :columns="grid.columns"
    :rows="grid.rows"
    :row-key="grid.rowKey"
    @row-click="onRowClick"
  />
</template>

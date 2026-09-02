<script setup lang="ts">
import { computed, ref } from 'vue'
import type { GridColumn, MarketRow } from '@/mocks/columns'
import { formatCompact, formatNumber, formatPct } from '@/lib/utils'

const props = withDefaults(
  defineProps<{
    columns: GridColumn[]
    rows: MarketRow[]
    loading?: boolean
    error?: string | null
    searchable?: boolean
    rowKey?: string
  }>(),
  {
    loading: false,
    error: null,
    searchable: true,
    rowKey: 'symbol',
  },
)

const emit = defineEmits<{
  rowClick: [row: MarketRow]
}>()

const query = ref('')
const sortKey = ref<string | null>(null)
const sortDir = ref<'asc' | 'desc'>('asc')

function cellValue(row: MarketRow, key: string): string | number | null | undefined {
  return row[key]
}

function formatCell(col: GridColumn, value: string | number | null | undefined): string {
  if (value == null || value === '') return '—'
  if (typeof value === 'string') return value
  switch (col.format) {
    case 'number':
      return formatNumber(value, Number.isInteger(value) ? 0 : 2)
    case 'compact':
      return formatCompact(value)
    case 'pct':
      return formatPct(value)
    default:
      return String(value)
  }
}

const filtered = computed(() => {
  let list = [...props.rows]
  const q = query.value.trim().toLowerCase()
  if (q) {
    list = list.filter((row) =>
      props.columns.some((col) => String(cellValue(row, col.key) ?? '').toLowerCase().includes(q)),
    )
  }
  if (sortKey.value) {
    const key = sortKey.value
    list.sort((a, b) => {
      const av = cellValue(a, key)
      const bv = cellValue(b, key)
      if (av == null && bv == null) return 0
      if (av == null) return 1
      if (bv == null) return -1
      if (typeof av === 'number' && typeof bv === 'number') {
        return sortDir.value === 'asc' ? av - bv : bv - av
      }
      const as = String(av)
      const bs = String(bv)
      return sortDir.value === 'asc' ? as.localeCompare(bs) : bs.localeCompare(as)
    })
  }
  return list
})

function toggleSort(key: string) {
  if (sortKey.value === key) {
    sortDir.value = sortDir.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortKey.value = key
    sortDir.value = 'asc'
  }
}

function headerLabel(col: GridColumn): string {
  return col.label ?? col.key
}
</script>

<template>
  <div class="space-y-3">
    <div v-if="searchable" class="flex items-center gap-2">
      <input
        v-model="query"
        type="search"
        placeholder="Filter rows…"
        class="h-9 w-full max-w-sm rounded-md border border-input bg-background px-3 text-sm outline-none ring-offset-background focus-visible:ring-2 focus-visible:ring-ring"
      />
      <span class="text-xs text-muted-foreground">{{ filtered.length }} rows</span>
    </div>

    <p v-if="error" class="rounded-md border border-ask/40 bg-ask/10 px-3 py-2 text-sm text-ask">
      {{ error }}
    </p>
    <p v-else-if="loading" class="text-sm text-muted-foreground">Loading…</p>

    <div class="overflow-auto rounded-md border">
      <table class="w-full min-w-[640px] border-collapse text-sm">
        <thead class="bg-muted/60">
          <tr>
            <th
              v-for="col in columns"
              :key="col.key"
              class="cursor-pointer select-none whitespace-nowrap border-b px-3 py-2 font-medium text-muted-foreground"
              :class="{
                'text-left': (col.align ?? 'left') === 'left',
                'text-right': col.align === 'right' || col.format === 'number' || col.format === 'compact' || col.format === 'pct',
                'text-center': col.align === 'center',
              }"
              @click="toggleSort(col.key)"
            >
              {{ headerLabel(col) }}
              <span v-if="sortKey === col.key" class="ml-1 text-xs">{{ sortDir === 'asc' ? '↑' : '↓' }}</span>
            </th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="(row, idx) in filtered"
            :key="String(cellValue(row, rowKey) ?? idx)"
            class="border-b last:border-0 hover:bg-accent/50 cursor-pointer"
            @click="emit('rowClick', row)"
          >
            <td
              v-for="col in columns"
              :key="col.key"
              class="whitespace-nowrap px-3 py-2 font-mono text-[13px] tabular-nums"
              :class="{
                'text-left font-sans': col.format === 'text' || !col.format,
                'text-right': col.format === 'number' || col.format === 'compact' || col.format === 'pct',
              }"
            >
              {{ formatCell(col, cellValue(row, col.key)) }}
            </td>
          </tr>
          <tr v-if="!loading && filtered.length === 0">
            <td :colspan="columns.length" class="px-3 py-8 text-center text-muted-foreground">No rows</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

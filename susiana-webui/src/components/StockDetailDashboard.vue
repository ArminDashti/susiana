<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import type { StockDetailMock } from '@/mocks/market'
import { formatCompact, formatNumber, formatPct } from '@/lib/utils'
import DataGrid from '@/components/DataGrid.vue'
import type { GridColumn } from '@/mocks/columns'

const props = defineProps<{
  detail: StockDetailMock
}>()

const now = ref(new Date())
let timer: number | undefined

onMounted(() => {
  timer = window.setInterval(() => {
    now.value = new Date()
  }, 1000)
})
onUnmounted(() => {
  if (timer) window.clearInterval(timer)
})

const clock = computed(() =>
  now.value.toLocaleTimeString('en-GB', { hour12: false }),
)

const changePositive = computed(() => props.detail.change >= 0)

const peerColumns: GridColumn[] = [
  { key: 'symbol' },
  { key: 'close', format: 'number' },
  { key: 'last', format: 'number' },
  { key: 'count', format: 'number' },
  { key: 'volume', format: 'compact' },
  { key: 'value', format: 'compact' },
]

const announcementColumns: GridColumn[] = [
  { key: 'date' },
  { key: 'title', format: 'text' },
]

const chartPoints = computed(() => {
  const prices = props.detail.chart.map((p) => p.price)
  const min = Math.min(...prices)
  const max = Math.max(...prices)
  const w = 520
  const h = 180
  const pad = 8
  return props.detail.chart
    .map((p, i) => {
      const x = pad + (i / Math.max(props.detail.chart.length - 1, 1)) * (w - pad * 2)
      const y = pad + (1 - (p.price - min) / Math.max(max - min, 1)) * (h - pad * 2)
      return `${x},${y}`
    })
    .join(' ')
})

function pctOf(part: number, total: number): string {
  if (!total) return '0%'
  return `${((part / total) * 100).toFixed(1)}%`
}
</script>

<template>
  <div class="space-y-4" dir="ltr">
    <!-- Header -->
    <div class="flex flex-wrap items-start justify-between gap-3 rounded-md border bg-card px-4 py-3">
      <div>
        <h2 class="text-lg font-semibold">
          {{ detail.company_name }}
          <span class="text-muted-foreground">({{ detail.symbol }})</span>
        </h2>
        <p class="text-sm text-muted-foreground">{{ detail.market_segment }}</p>
      </div>
      <div class="text-right text-sm">
        <p>
          Status:
          <span class="font-medium text-primary">{{ detail.trading_status }}</span>
        </p>
        <p class="font-mono tabular-nums text-muted-foreground">{{ clock }}</p>
      </div>
    </div>

    <!-- Stats + retail/institutional -->
    <div class="grid gap-4 xl:grid-cols-[1fr_280px]">
      <div class="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
        <div class="rounded-md border bg-card p-3 text-sm">
          <p class="mb-2 text-xs font-semibold uppercase text-muted-foreground">Price</p>
          <dl class="space-y-1 font-mono tabular-nums">
            <div class="flex justify-between"><dt>last</dt><dd>{{ formatNumber(detail.last) }}</dd></div>
            <div class="flex justify-between"><dt>close</dt><dd>{{ formatNumber(detail.close) }}</dd></div>
            <div class="flex justify-between"><dt>open</dt><dd>{{ formatNumber(detail.open) }}</dd></div>
            <div class="flex justify-between"><dt>prev_close</dt><dd>{{ formatNumber(detail.prev_close) }}</dd></div>
            <div class="flex justify-between font-semibold" :class="changePositive ? 'text-emerald-600' : 'text-ask'">
              <dt>change</dt>
              <dd>{{ formatNumber(detail.change) }} ({{ formatPct(detail.change_pct) }})</dd>
            </div>
          </dl>
        </div>
        <div class="rounded-md border bg-card p-3 text-sm">
          <p class="mb-2 text-xs font-semibold uppercase text-muted-foreground">Activity</p>
          <dl class="space-y-1 font-mono tabular-nums">
            <div class="flex justify-between"><dt>trades</dt><dd>{{ formatNumber(detail.trades) }}</dd></div>
            <div class="flex justify-between"><dt>volume</dt><dd>{{ formatCompact(detail.volume) }}</dd></div>
            <div class="flex justify-between"><dt>turnover</dt><dd>{{ formatCompact(detail.turnover) }}</dd></div>
            <div class="flex justify-between"><dt>market_cap</dt><dd>{{ formatCompact(detail.market_cap) }}</dd></div>
          </dl>
        </div>
        <div class="rounded-md border bg-card p-3 text-sm">
          <p class="mb-2 text-xs font-semibold uppercase text-muted-foreground">Range</p>
          <dl class="space-y-1 font-mono tabular-nums">
            <div class="flex justify-between"><dt>day</dt><dd>{{ formatNumber(detail.day_low) }} – {{ formatNumber(detail.day_high) }}</dd></div>
            <div class="flex justify-between"><dt>limits</dt><dd>{{ formatNumber(detail.lower_limit) }} – {{ formatNumber(detail.upper_limit) }}</dd></div>
            <div class="flex justify-between"><dt>week</dt><dd>{{ formatNumber(detail.week_low) }} – {{ formatNumber(detail.week_high) }}</dd></div>
            <div class="flex justify-between"><dt>year</dt><dd>{{ formatNumber(detail.year_low) }} – {{ formatNumber(detail.year_high) }}</dd></div>
          </dl>
        </div>
        <div class="rounded-md border bg-card p-3 text-sm">
          <p class="mb-2 text-xs font-semibold uppercase text-muted-foreground">Company</p>
          <dl class="space-y-1 font-mono tabular-nums">
            <div class="flex justify-between"><dt>shares_outstanding</dt><dd>{{ formatCompact(detail.shares_outstanding) }}</dd></div>
            <div class="flex justify-between"><dt>base_volume</dt><dd>{{ formatNumber(detail.base_volume) }}</dd></div>
            <div class="flex justify-between"><dt>free_float</dt><dd>{{ formatNumber(detail.free_float, 0) }}%</dd></div>
            <div class="flex justify-between"><dt>avg_monthly_volume</dt><dd>{{ formatCompact(detail.avg_monthly_volume) }}</dd></div>
            <div class="flex justify-between"><dt>eps</dt><dd>{{ formatNumber(detail.eps) }}</dd></div>
            <div class="flex justify-between"><dt>pe_ratio</dt><dd>{{ formatNumber(detail.pe_ratio, 2) }}</dd></div>
            <div class="flex justify-between"><dt>group_pe</dt><dd>{{ formatNumber(detail.group_pe, 2) }}</dd></div>
            <div class="flex justify-between"><dt>ps_ratio</dt><dd>{{ formatNumber(detail.ps_ratio, 2) }}</dd></div>
          </dl>
        </div>
      </div>

      <div class="rounded-md border bg-card p-3 text-sm">
        <p class="mb-2 text-xs font-semibold uppercase text-muted-foreground">Retail / Institutional</p>
        <table class="w-full text-xs">
          <thead>
            <tr class="text-muted-foreground">
              <th class="py-1 text-left"></th>
              <th class="py-1 text-right">buy_volume</th>
              <th class="py-1 text-right">%</th>
              <th class="py-1 text-right">buy_count</th>
            </tr>
          </thead>
          <tbody class="font-mono tabular-nums">
            <tr>
              <td class="py-1">retail</td>
              <td class="text-right">{{ formatCompact(detail.retail_buy_volume) }}</td>
              <td class="text-right">{{ pctOf(detail.retail_buy_volume, detail.retail_buy_volume + detail.institutional_buy_volume) }}</td>
              <td class="text-right">{{ formatNumber(detail.retail_buy_count) }}</td>
            </tr>
            <tr>
              <td class="py-1">institutional</td>
              <td class="text-right">{{ formatCompact(detail.institutional_buy_volume) }}</td>
              <td class="text-right">{{ pctOf(detail.institutional_buy_volume, detail.retail_buy_volume + detail.institutional_buy_volume) }}</td>
              <td class="text-right">{{ formatNumber(detail.institutional_buy_count) }}</td>
            </tr>
          </tbody>
        </table>
        <table class="mt-3 w-full text-xs">
          <thead>
            <tr class="text-muted-foreground">
              <th class="py-1 text-left"></th>
              <th class="py-1 text-right">sell_volume</th>
              <th class="py-1 text-right">%</th>
              <th class="py-1 text-right">sell_count</th>
            </tr>
          </thead>
          <tbody class="font-mono tabular-nums">
            <tr>
              <td class="py-1">retail</td>
              <td class="text-right">{{ formatCompact(detail.retail_sell_volume) }}</td>
              <td class="text-right">{{ pctOf(detail.retail_sell_volume, detail.retail_sell_volume + detail.institutional_sell_volume) }}</td>
              <td class="text-right">{{ formatNumber(detail.retail_sell_count) }}</td>
            </tr>
            <tr>
              <td class="py-1">institutional</td>
              <td class="text-right">{{ formatCompact(detail.institutional_sell_volume) }}</td>
              <td class="text-right">{{ pctOf(detail.institutional_sell_volume, detail.retail_sell_volume + detail.institutional_sell_volume) }}</td>
              <td class="text-right">{{ formatNumber(detail.institutional_sell_count) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Chart + order book -->
    <div class="grid gap-4 lg:grid-cols-2">
      <div class="rounded-md border bg-card p-3">
        <p class="mb-2 text-xs font-semibold uppercase text-muted-foreground">Intraday (price)</p>
        <svg viewBox="0 0 520 180" class="h-48 w-full">
          <polyline
            fill="none"
            stroke="hsl(var(--primary))"
            stroke-width="2"
            :points="chartPoints"
          />
        </svg>
        <div class="mt-1 flex justify-between text-[10px] text-muted-foreground">
          <span v-for="p in detail.chart" :key="p.t">{{ p.t }}</span>
        </div>
      </div>

      <div class="rounded-md border bg-card p-3">
        <p class="mb-2 text-xs font-semibold uppercase text-muted-foreground">Order book</p>
        <div class="grid grid-cols-2 gap-2 text-xs">
          <table class="w-full">
            <thead>
              <tr class="bg-bid text-bid-foreground">
                <th class="px-2 py-1 text-right">count</th>
                <th class="px-2 py-1 text-right">volume</th>
                <th class="px-2 py-1 text-right">bid</th>
              </tr>
            </thead>
            <tbody class="font-mono tabular-nums">
              <tr v-for="(b, i) in detail.bids" :key="'b' + i" class="border-b">
                <td class="px-2 py-1 text-right">{{ formatNumber(b.count) }}</td>
                <td class="px-2 py-1 text-right">{{ formatCompact(b.volume) }}</td>
                <td class="px-2 py-1 text-right text-bid">{{ formatNumber(b.price) }}</td>
              </tr>
            </tbody>
          </table>
          <table class="w-full">
            <thead>
              <tr class="bg-ask text-ask-foreground">
                <th class="px-2 py-1 text-left">ask</th>
                <th class="px-2 py-1 text-right">volume</th>
                <th class="px-2 py-1 text-right">count</th>
              </tr>
            </thead>
            <tbody class="font-mono tabular-nums">
              <tr v-for="(a, i) in detail.asks" :key="'a' + i" class="border-b">
                <td class="px-2 py-1 text-left text-ask">{{ formatNumber(a.price) }}</td>
                <td class="px-2 py-1 text-right">{{ formatCompact(a.volume) }}</td>
                <td class="px-2 py-1 text-right">{{ formatNumber(a.count) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Announcements + peers -->
    <div class="grid gap-4 lg:grid-cols-2">
      <div class="rounded-md border bg-card p-3">
        <p class="mb-2 text-xs font-semibold uppercase text-muted-foreground">Announcements</p>
        <DataGrid :columns="announcementColumns" :rows="detail.announcements" :searchable="false" row-key="title" />
      </div>
      <div class="rounded-md border bg-card p-3">
        <p class="mb-2 text-xs font-semibold uppercase text-muted-foreground">Peer group</p>
        <DataGrid :columns="peerColumns" :rows="detail.peers" :searchable="false" row-key="symbol" />
      </div>
    </div>
  </div>
</template>

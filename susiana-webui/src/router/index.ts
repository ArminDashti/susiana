import { createRouter, createWebHistory } from 'vue-router'
import { getToken } from '@/lib/auth'
import AppLayout from '@/views/AppLayout.vue'
import LoginView from '@/views/LoginView.vue'
import DashboardView from '@/views/DashboardView.vue'
import SettingsView from '@/views/SettingsView.vue'
import MarketGridView from '@/views/MarketGridView.vue'
import TseStocksView from '@/views/tse/TseStocksView.vue'
import TseStockDetailView from '@/views/tse/TseStockDetailView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/login',
      name: 'login',
      component: LoginView,
      meta: { guest: true, title: 'Login' },
    },
    {
      path: '/',
      component: AppLayout,
      meta: { requiresAuth: true },
      children: [
        { path: '', redirect: '/dashboard' },
        { path: 'dashboard', name: 'dashboard', component: DashboardView, meta: { title: 'Dashboard' } },
        { path: 'settings', name: 'settings', component: SettingsView, meta: { title: 'Settings' } },
        {
          path: 'stock-market/tse',
          name: 'tse',
          component: TseStocksView,
          meta: { title: 'Stock market / TSE' },
        },
        {
          path: 'stock-market/tse/stock/history/:symbol',
          name: 'tse-stock-history',
          component: MarketGridView,
          meta: { title: 'TSE stock history', gridKind: 'history' },
        },
        {
          path: 'stock-market/tse/stock/:symbol',
          name: 'tse-stock',
          component: TseStockDetailView,
          meta: { title: 'TSE stock' },
        },
        {
          path: 'stock-market/tse/indexes',
          name: 'tse-indexes',
          component: MarketGridView,
          meta: { title: 'TSE indexes', gridKind: 'tse-indexes' },
        },
        {
          path: 'commodity/crude-oil',
          component: MarketGridView,
          meta: { title: 'Commodity / Crude oil', gridKind: 'commodity', gridKey: 'crude-oil' },
        },
        {
          path: 'commodity/natural-gas',
          component: MarketGridView,
          meta: { title: 'Commodity / Natural gas', gridKind: 'commodity', gridKey: 'natural-gas' },
        },
        {
          path: 'commodity/gold',
          component: MarketGridView,
          meta: { title: 'Commodity / Gold', gridKind: 'commodity', gridKey: 'gold' },
        },
        {
          path: 'commodity/silver',
          component: MarketGridView,
          meta: { title: 'Commodity / Silver', gridKind: 'commodity', gridKey: 'silver' },
        },
        {
          path: 'commodity/copper',
          component: MarketGridView,
          meta: { title: 'Commodity / Copper', gridKind: 'commodity', gridKey: 'copper' },
        },
        {
          path: 'commodity/steel',
          component: MarketGridView,
          meta: { title: 'Commodity / Steel', gridKind: 'commodity', gridKey: 'steel' },
        },
        {
          path: 'indexes/asia',
          component: MarketGridView,
          meta: { title: 'Indexes / Asia', gridKind: 'index', gridKey: 'asia' },
        },
        {
          path: 'indexes/top30',
          component: MarketGridView,
          meta: { title: 'Indexes / Top 30', gridKind: 'index', gridKey: 'top30' },
        },
        {
          path: 'indexes/europe',
          component: MarketGridView,
          meta: { title: 'Indexes / Europe', gridKind: 'index', gridKey: 'europe' },
        },
        {
          path: 'indexes/north-america',
          component: MarketGridView,
          meta: { title: 'Indexes / North America', gridKind: 'index', gridKey: 'north-america' },
        },
        {
          path: 'indexes/south-america',
          component: MarketGridView,
          meta: { title: 'Indexes / South America', gridKind: 'index', gridKey: 'south-america' },
        },
        {
          path: 'indexes/middle-east',
          component: MarketGridView,
          meta: { title: 'Indexes / Middle East', gridKind: 'index', gridKey: 'middle-east' },
        },
        {
          path: 'indexes/australia',
          component: MarketGridView,
          meta: { title: 'Indexes / Australia', gridKind: 'index', gridKey: 'australia' },
        },
        {
          path: 'indexes/africa',
          component: MarketGridView,
          meta: { title: 'Indexes / Africa', gridKind: 'index', gridKey: 'africa' },
        },
        {
          path: 'currencies',
          component: MarketGridView,
          meta: { title: 'Currencies', gridKind: 'currencies' },
        },
        {
          path: 'crypto-currencies',
          component: MarketGridView,
          meta: { title: 'Crypto currencies', gridKind: 'crypto' },
        },
      ],
    },
  ],
})

router.beforeEach((to) => {
  const token = getToken()
  if (to.meta.requiresAuth && !token) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }
  if (to.matched.some((r) => r.meta.requiresAuth) && !token && to.name !== 'login') {
    return { name: 'login', query: { redirect: to.fullPath } }
  }
  if (to.meta.guest && token) {
    return { name: 'dashboard' }
  }
  return true
})

export default router

import type { MarketRow } from './columns'

const asOf = '2026-08-26 12:30'

function indexRows(prefix: string, names: string[]): MarketRow[] {
  return names.map((name, i) => ({
    symbol: `${prefix}${i + 1}`,
    name,
    last: 1000 + i * 37.5,
    change: (i % 2 === 0 ? 1 : -1) * (2.5 + i),
    change_pct: (i % 2 === 0 ? 1 : -1) * (0.35 + i * 0.08),
    volume: 1_200_000 + i * 85_000,
  }))
}

function commodityRows(symbol: string, name: string): MarketRow[] {
  return [
    {
      symbol,
      name,
      last: 72.4,
      change_pct: 0.85,
      open: 71.9,
      high: 73.1,
      low: 71.2,
      volume: 4_250_000,
      as_of: asOf,
    },
    {
      symbol: `${symbol}-F`,
      name: `${name} Futures`,
      last: 73.0,
      change_pct: -0.22,
      open: 73.2,
      high: 73.8,
      low: 72.6,
      volume: 2_100_000,
      as_of: asOf,
    },
  ]
}

export const mockIndexes = {
  asia: indexRows('ASIA', ['Nikkei 225', 'Hang Seng', 'Shanghai Comp', 'KOSPI', 'Sensex']),
  top30: indexRows('TOP', ['Top 30 A', 'Top 30 B', 'Top 30 C', 'Top 30 D', 'Top 30 E']),
  europe: indexRows('EU', ['FTSE 100', 'DAX', 'CAC 40', 'IBEX 35', 'SMI']),
  'north-america': indexRows('NA', ['S&P 500', 'Dow Jones', 'Nasdaq', 'TSX', 'IPC']),
  'south-america': indexRows('SA', ['Bovespa', 'IPSA', 'MERVAL', 'COLCAP', 'IGBVL']),
  'middle-east': indexRows('ME', ['TASI', 'DFMGI', 'QSI', 'ADI', 'KWSE']),
  australia: indexRows('AU', ['ASX 200', 'NZX 50', 'All Ords']),
  africa: indexRows('AF', ['JSE Top 40', 'EGX 30', 'NSE ASI']),
  tse: indexRows('TSE', ['TEDPIX', 'Equal Weight', 'Industry', 'Financial', 'Free Float']),
}

export const mockCommodities: Record<string, MarketRow[]> = {
  'crude-oil': commodityRows('CL', 'Crude Oil'),
  'natural-gas': commodityRows('NG', 'Natural Gas'),
  gold: commodityRows('XAU', 'Gold'),
  silver: commodityRows('XAG', 'Silver'),
  copper: commodityRows('HG', 'Copper'),
  steel: commodityRows('STL', 'Steel'),
}

export const mockCurrencies: MarketRow[] = [
  { pair: 'USD/IRR', last: 920_000, bid: 919_500, ask: 920_500, change_pct: 0.12, volume_24h: 1.2e9, as_of: asOf },
  { pair: 'EUR/USD', last: 1.0842, bid: 1.084, ask: 1.0844, change_pct: -0.08, volume_24h: 8.4e10, as_of: asOf },
  { pair: 'GBP/USD', last: 1.271, bid: 1.2708, ask: 1.2712, change_pct: 0.05, volume_24h: 3.1e10, as_of: asOf },
  { pair: 'USD/JPY', last: 149.32, bid: 149.3, ask: 149.34, change_pct: 0.21, volume_24h: 4.5e10, as_of: asOf },
]

export const mockCrypto: MarketRow[] = [
  { symbol: 'BTC/USDT', last: 64_250, bid: 64_240, ask: 64_260, change_pct: 1.45, volume_24h: 2.8e10, as_of: asOf },
  { symbol: 'ETH/USDT', last: 3_420, bid: 3_418, ask: 3_422, change_pct: 0.92, volume_24h: 1.1e10, as_of: asOf },
  { symbol: 'SOL/USDT', last: 148.5, bid: 148.4, ask: 148.6, change_pct: -0.55, volume_24h: 2.4e9, as_of: asOf },
]

export function mockHistory(symbol: string): MarketRow[] {
  const base = 12_500
  return Array.from({ length: 12 }, (_, i) => {
    const close = base + i * 35 + (i % 3) * 10
    return {
      as_of: `2026-08-${String(26 - i).padStart(2, '0')}`,
      open: close - 40,
      high: close + 80,
      low: close - 90,
      close,
      volume: 80_000_000 + i * 2_500_000,
      turnover: (close * (80_000_000 + i * 2_500_000)) / 10,
      symbol,
    }
  })
}

export const mockTseStocks: MarketRow[] = [
  {
    ticker: 'Shabandar',
    isin: 'IRO1BNDR0001',
    last: 12950,
    close: 12950,
    open: 12950,
    volume: 120_666_000,
    turnover: 1_562_628_000_000,
    trades: 3757,
    pe_ratio: 5.07,
    eps: 2555,
  },
  {
    ticker: 'Shepna',
    isin: 'IRO1PNZA0001',
    last: 4210,
    close: 4180,
    open: 4150,
    volume: 45_200_000,
    turnover: 190_000_000_000,
    trades: 2100,
    pe_ratio: 6.2,
    eps: 680,
  },
  {
    ticker: 'Shotran',
    isin: 'IRO1TRAN0001',
    last: 890,
    close: 875,
    open: 870,
    volume: 12_400_000,
    turnover: 11_000_000_000,
    trades: 980,
    pe_ratio: 8.1,
    eps: 110,
  },
]

export type OrderBookLevel = { count: number; volume: number; price: number }

export type StockDetailMock = {
  company_name: string
  symbol: string
  market_segment: string
  trading_status: string
  last: number
  close: number
  open: number
  prev_close: number
  change: number
  change_pct: number
  trades: number
  volume: number
  turnover: number
  market_cap: number
  day_low: number
  day_high: number
  lower_limit: number
  upper_limit: number
  week_low: number
  week_high: number
  year_low: number
  year_high: number
  shares_outstanding: number
  base_volume: number
  free_float: number
  avg_monthly_volume: number
  eps: number
  pe_ratio: number
  group_pe: number
  ps_ratio: number
  retail_buy_volume: number
  institutional_buy_volume: number
  retail_buy_count: number
  institutional_buy_count: number
  retail_sell_volume: number
  institutional_sell_volume: number
  retail_sell_count: number
  institutional_sell_count: number
  chart: { t: string; price: number }[]
  bids: OrderBookLevel[]
  asks: OrderBookLevel[]
  announcements: { date: string; title: string }[]
  peers: MarketRow[]
}

export function mockStockDetail(symbol: string): StockDetailMock {
  const last = 12950
  const prev = 12580
  return {
    company_name: 'Bandar Abbas Oil Refining Co.',
    symbol: symbol || 'Shabandar',
    market_segment: 'First Market (Main Board)',
    trading_status: 'Open',
    last,
    close: last,
    open: last,
    prev_close: prev,
    change: last - prev,
    change_pct: ((last - prev) / prev) * 100,
    trades: 3757,
    volume: 120_666_000,
    turnover: 1_562_628_000_000_000,
    market_cap: 5_306_170_064_000_000,
    day_low: 12690,
    day_high: 12950,
    lower_limit: 11960,
    upper_limit: 13200,
    week_low: 11680,
    week_high: 12950,
    year_low: 8550,
    year_high: 15580,
    shares_outstanding: 409_743_000_000,
    base_volume: 1,
    free_float: 26,
    avg_monthly_volume: 492_752_000,
    eps: 2555,
    pe_ratio: 5.07,
    group_pe: 5.68,
    ps_ratio: 0.62,
    retail_buy_volume: 68_400_000,
    institutional_buy_volume: 52_266_000,
    retail_buy_count: 2140,
    institutional_buy_count: 85,
    retail_sell_volume: 71_100_000,
    institutional_sell_volume: 49_566_000,
    retail_sell_count: 1980,
    institutional_sell_count: 62,
    chart: [
      { t: '09:00', price: 12690 },
      { t: '09:30', price: 12740 },
      { t: '10:00', price: 12810 },
      { t: '10:30', price: 12780 },
      { t: '11:00', price: 12860 },
      { t: '11:30', price: 12910 },
      { t: '12:00', price: 12940 },
      { t: '12:30', price: 12950 },
    ],
    bids: [
      { count: 42, volume: 1_250_000, price: 12940 },
      { count: 18, volume: 890_000, price: 12930 },
      { count: 27, volume: 1_100_000, price: 12920 },
      { count: 11, volume: 640_000, price: 12910 },
      { count: 9, volume: 510_000, price: 12900 },
    ],
    asks: [
      { count: 35, volume: 980_000, price: 12950 },
      { count: 22, volume: 720_000, price: 12960 },
      { count: 14, volume: 540_000, price: 12970 },
      { count: 19, volume: 810_000, price: 12980 },
      { count: 8, volume: 430_000, price: 12990 },
    ],
    announcements: [
      { date: '1404/05/28', title: 'Board decision on capital increase proposal' },
      { date: '1404/05/20', title: 'Audited financial statements — Q1' },
      { date: '1404/05/12', title: 'Material information disclosure' },
      { date: '1404/04/30', title: 'Monthly production and sales report' },
    ],
    peers: [
      { symbol: 'Shepna', close: 4180, last: 4210, count: 2100, volume: 45_200_000, value: 190e9 },
      { symbol: 'Shotran', close: 875, last: 890, count: 980, volume: 12_400_000, value: 11e9 },
      { symbol: 'Shavan', close: 1520, last: 1535, count: 1400, volume: 8_200_000, value: 12.5e9 },
    ],
  }
}

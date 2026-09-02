export type GridColumn = {
  key: string
  label?: string
  align?: 'left' | 'right' | 'center'
  format?: 'number' | 'compact' | 'pct' | 'text'
}

export type MarketRow = Record<string, string | number | null | undefined>

export const INDEX_COLUMNS: GridColumn[] = [
  { key: 'symbol' },
  { key: 'name' },
  { key: 'last', format: 'number' },
  { key: 'change', format: 'number' },
  { key: 'change_pct', format: 'pct' },
  { key: 'volume', format: 'compact' },
]

export const COMMODITY_COLUMNS: GridColumn[] = [
  { key: 'symbol' },
  { key: 'name' },
  { key: 'last', format: 'number' },
  { key: 'change_pct', format: 'pct' },
  { key: 'open', format: 'number' },
  { key: 'high', format: 'number' },
  { key: 'low', format: 'number' },
  { key: 'volume', format: 'compact' },
  { key: 'as_of' },
]

export const FX_COLUMNS: GridColumn[] = [
  { key: 'pair' },
  { key: 'last', format: 'number' },
  { key: 'bid', format: 'number' },
  { key: 'ask', format: 'number' },
  { key: 'change_pct', format: 'pct' },
  { key: 'volume_24h', format: 'compact' },
  { key: 'as_of' },
]

export const CRYPTO_COLUMNS: GridColumn[] = [
  { key: 'symbol' },
  { key: 'last', format: 'number' },
  { key: 'bid', format: 'number' },
  { key: 'ask', format: 'number' },
  { key: 'change_pct', format: 'pct' },
  { key: 'volume_24h', format: 'compact' },
  { key: 'as_of' },
]

export const HISTORY_COLUMNS: GridColumn[] = [
  { key: 'as_of' },
  { key: 'open', format: 'number' },
  { key: 'high', format: 'number' },
  { key: 'low', format: 'number' },
  { key: 'close', format: 'number' },
  { key: 'volume', format: 'compact' },
  { key: 'turnover', format: 'compact' },
]

export const TSE_LIST_COLUMNS: GridColumn[] = [
  { key: 'ticker' },
  { key: 'isin' },
  { key: 'last', format: 'number' },
  { key: 'close', format: 'number' },
  { key: 'open', format: 'number' },
  { key: 'volume', format: 'compact' },
  { key: 'turnover', format: 'compact' },
  { key: 'trades', format: 'number' },
  { key: 'pe_ratio', format: 'number' },
  { key: 'eps', format: 'number' },
]

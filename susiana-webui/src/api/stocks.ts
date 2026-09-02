import { apiFetch, type ListResponse } from './client'

export type Stock = {
  last?: number | null
  close?: number | null
  open?: number | null
  lower_limit?: number | null
  upper_limit?: number | null
  base_volume?: number | null
  volume?: number | null
  turnover?: number | null
  trades?: number | null
  retail_buy_volume?: number | null
  institutional_buy_volume?: number | null
  retail_buy_count?: number | null
  institutional_buy_count?: number | null
  sector_code?: string
  board_code?: string
  security_name?: string
  local_ticker?: string
  ticker?: string
  isin: string
  tsetmc_code?: string
  market_segment?: string
  trading_status?: string
  shares_outstanding?: number | null
  free_float?: number | null
  pe_ratio?: number | null
  eps?: number | null
  as_of?: string | null
}

export async function listStocks(params?: {
  page?: number
  limit?: number
  ticker?: string
}): Promise<ListResponse<Stock>> {
  const q = new URLSearchParams()
  if (params?.page) q.set('page', String(params.page))
  if (params?.limit) q.set('limit', String(params.limit))
  if (params?.ticker) q.set('ticker', params.ticker)
  const qs = q.toString()
  return apiFetch<ListResponse<Stock>>(`/api/v1/stocks${qs ? `?${qs}` : ''}`)
}

export async function getStockByIsin(isin: string): Promise<Stock> {
  return apiFetch<Stock>(`/api/v1/stocks/${encodeURIComponent(isin)}`)
}

export async function findStockBySymbol(symbol: string): Promise<Stock | null> {
  const needle = symbol.trim().toLowerCase()
  const page = await listStocks({ page: 1, limit: 200, ticker: symbol })
  const hit =
    page.data.find(
      (s) =>
        s.ticker?.toLowerCase() === needle ||
        s.local_ticker?.toLowerCase() === needle ||
        s.isin.toLowerCase() === needle,
    ) ?? null
  if (hit) return hit
  try {
    return await getStockByIsin(symbol)
  } catch {
    return null
  }
}

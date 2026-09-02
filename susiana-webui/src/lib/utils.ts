import type { ClassValue } from 'clsx'
import { clsx } from 'clsx'
import { twMerge } from 'tailwind-merge'

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export function formatNumber(value: number | null | undefined, digits = 0): string {
  if (value == null || Number.isNaN(value)) return '—'
  return new Intl.NumberFormat('en-US', {
    maximumFractionDigits: digits,
    minimumFractionDigits: digits,
  }).format(value)
}

export function formatCompact(value: number | null | undefined): string {
  if (value == null || Number.isNaN(value)) return '—'
  const abs = Math.abs(value)
  if (abs >= 1e9) return `${formatNumber(value / 1e9, 3)} B`
  if (abs >= 1e6) return `${formatNumber(value / 1e6, 3)} M`
  if (abs >= 1e3) return `${formatNumber(value / 1e3, 3)} K`
  return formatNumber(value, 0)
}

export function formatPct(value: number | null | undefined): string {
  if (value == null || Number.isNaN(value)) return '—'
  const sign = value > 0 ? '+' : ''
  return `${sign}${formatNumber(value, 2)}%`
}

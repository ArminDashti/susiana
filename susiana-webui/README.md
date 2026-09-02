# susiana-webui

Vue 3 + Vite + Tailwind + shadcn-style UI for Susiana market dashboards (TSE, commodities, indexes, currencies, crypto).

## Stack

- Vue 3 + TypeScript + Vue Router
- Tailwind CSS + Inter
- PWA via `vite-plugin-pwa` (Workbox)

## Setup

```bash
cp .env.example .env
npm install
npm run dev
```

Dev URL: `http://localhost:5178`  
API proxy target: `VITE_API_PROXY_TARGET` (default `http://127.0.0.1:8080` → [susiana-api](https://github.com/ArminDashti/susiana-api))

## Login

| Field | Value |
|-------|-------|
| Username | `armin` |
| Password | `dopadopa123` |

Auth is client-side until the API exposes login endpoints.

## Routes

- `/login`, `/settings`, `/dashboard`
- `/stock-market/tse` (+ stock detail, history, indexes)
- `/commodity/*`, `/indexes/*`, `/currencies`, `/crypto-currencies`

TSE stock list/detail prefer live `susiana-api` data; other market pages use typed mock fixtures.

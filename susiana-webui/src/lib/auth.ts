const TOKEN_KEY = 'susiana_token'
const USER_KEY = 'susiana_user'

export type AuthUser = {
  username: string
  display_name: string
}

const DEFAULT_USER: AuthUser = {
  username: 'armin',
  display_name: 'Armin',
}

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY)
}

export function getStoredUser(): AuthUser | null {
  const raw = localStorage.getItem(USER_KEY)
  if (!raw) return null
  try {
    return JSON.parse(raw) as AuthUser
  } catch {
    return null
  }
}

export function login(username: string, password: string): AuthUser {
  if (username.trim() === 'armin' && password === 'dopadopa123') {
    const user = { ...DEFAULT_USER }
    localStorage.setItem(TOKEN_KEY, `local.${Date.now()}`)
    localStorage.setItem(USER_KEY, JSON.stringify(user))
    return user
  }
  throw new Error('Invalid username or password')
}

export function logout(): void {
  localStorage.removeItem(TOKEN_KEY)
  localStorage.removeItem(USER_KEY)
}

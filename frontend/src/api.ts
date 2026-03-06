import type { Box, Collection, DashboardStats, Game, Miniature, MiniatureImage, MiniaturePaint, Paint, User } from './types'

const BASE = ''

function token(): string | null {
  return localStorage.getItem('token')
}

function authHeaders(): HeadersInit {
  return { 'Content-Type': 'application/json', Authorization: `Bearer ${token()}` }
}

async function req<T>(method: string, path: string, body?: unknown): Promise<T> {
  const res = await fetch(BASE + path, {
    method,
    headers: authHeaders(),
    body: body !== undefined ? JSON.stringify(body) : undefined,
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error(err.error ?? res.statusText)
  }
  if (res.status === 204) return undefined as T
  return res.json()
}

// Auth
export async function register(username: string, email: string, password: string) {
  const data = await req<{ token: string; user: User }>('POST', '/api/auth/register', { username, email, password })
  localStorage.setItem('token', data.token)
  return data.user
}

export async function login(email: string, password: string) {
  const data = await req<{ token: string; user: User }>('POST', '/api/auth/login', { email, password })
  localStorage.setItem('token', data.token)
  return data.user
}

export function logout() {
  localStorage.removeItem('token')
}

// Dashboard
export const getDashboard = () => req<DashboardStats>('GET', '/api/dashboard')

// Collections
export const getCollections = () => req<Collection[]>('GET', '/api/collections')
export const getCollection = (id: string) => req<Collection>('GET', `/api/collections/${id}`)
export const createCollection = (body: { name: string; notes?: string }) =>
  req<Collection>('POST', '/api/collections', body)
export const updateCollection = (id: string, body: Partial<Collection>) =>
  req<Collection>('PUT', `/api/collections/${id}`, body)
export const deleteCollection = (id: string) => req<void>('DELETE', `/api/collections/${id}`)

// Games
export const getGames = (collectionID: string) =>
  req<Game[]>('GET', `/api/collections/${collectionID}/games`)
export const getGame = (id: string) => req<Game>('GET', `/api/games/${id}`)
export const createGame = (collectionID: string, body: { name: string; publisher?: string; year?: number }) =>
  req<Game>('POST', `/api/collections/${collectionID}/games`, body)
export const updateGame = (id: string, body: { name: string; publisher?: string; year?: number }) =>
  req<Game>('PUT', `/api/games/${id}`, body)
export const deleteGame = (id: string) => req<void>('DELETE', `/api/games/${id}`)

// Boxes
export const getBoxes = (gameID: string) =>
  req<Box[]>('GET', `/api/games/${gameID}/boxes`)
export const getBox = (id: string) => req<Box>('GET', `/api/boxes/${id}`)
export interface MiniSuggestion {
  name: string
  unit_type: string
  quantity: number
}

export const createBox = (gameID: string, body: { name: string }) =>
  req<{ box: Box; suggestions: MiniSuggestion[]; claude_error?: string; source: 'catalog' | 'claude' | 'none' }>('POST', `/api/games/${gameID}/boxes`, body)
export const confirmBox = (boxID: string, miniatures: MiniSuggestion[]) =>
  req<{ miniatures: Miniature[] }>('POST', `/api/boxes/${boxID}/confirm`, { miniatures })
export const updateBox = (id: string, body: { name: string }) =>
  req<Box>('PUT', `/api/boxes/${id}`, body)
export const deleteBox = (id: string) => req<void>('DELETE', `/api/boxes/${id}`)

// Miniatures
export const getMiniatures = (boxID: string) =>
  req<Miniature[]>('GET', `/api/boxes/${boxID}/miniatures`)
export const getMiniature = (id: string) => req<Miniature>('GET', `/api/miniatures/${id}`)
export const createMiniature = (boxID: string, body: { name: string; unit_type?: string; quantity?: number; notes?: string }) =>
  req<Miniature>('POST', `/api/boxes/${boxID}/miniatures`, body)
export const updateMiniature = (id: string, body: Partial<Miniature>) =>
  req<Miniature>('PATCH', `/api/miniatures/${id}`, body)
export const deleteMiniature = (id: string) => req<void>('DELETE', `/api/miniatures/${id}`)

// Miniature paints
export const addMiniaturePaint = (miniID: string, body: { paint_id: string; purpose?: string; notes?: string }) =>
  req<MiniaturePaint>('POST', `/api/miniatures/${miniID}/paints`, body)
export const removeMiniaturePaint = (miniID: string, linkID: string) =>
  req<void>('DELETE', `/api/miniatures/${miniID}/paints/${linkID}`)

// Images
export async function uploadImage(miniID: string, file: File, stage: string, caption: string): Promise<MiniatureImage> {
  const form = new FormData()
  form.append('image', file)
  form.append('stage', stage)
  form.append('caption', caption)
  const res = await fetch(`${BASE}/api/miniatures/${miniID}/images`, {
    method: 'POST',
    headers: { Authorization: `Bearer ${token()}` },
    body: form,
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }))
    throw new Error(err.error ?? res.statusText)
  }
  return res.json()
}
export const deleteImage = (miniID: string, imageID: string) =>
  req<void>('DELETE', `/api/miniatures/${miniID}/images/${imageID}`)

// Paints
export const getPaints = () => req<Paint[]>('GET', '/api/paints')
export const createPaint = (body: { brand: string; name: string; color?: string; type?: string }) =>
  req<Paint>('POST', '/api/paints', body)
export const updatePaint = (id: string, body: Partial<Paint>) => req<Paint>('PUT', `/api/paints/${id}`, body)
export const deletePaint = (id: string) => req<void>('DELETE', `/api/paints/${id}`)

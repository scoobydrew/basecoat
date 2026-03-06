export type PaintingStatus =
  | 'unpainted'
  | 'primed'
  | 'basecoated'
  | 'shaded'
  | 'detailed'
  | 'finished'

export interface User {
  id: string
  username: string
  email: string
  created_at: string
}

export interface Collection {
  id: string
  user_id: string
  name: string
  game: string
  notes: string
  created_at: string
}

export interface Miniature {
  id: string
  collection_id: string
  user_id: string
  name: string
  unit_type: string
  quantity: number
  status: PaintingStatus
  notes: string
  created_at: string
  updated_at: string
  paints?: MiniaturePaint[]
  images?: MiniatureImage[]
}

export interface Paint {
  id: string
  user_id: string
  brand: string
  name: string
  color: string
  type: string
  created_at: string
}

export interface MiniaturePaint {
  id: string
  miniature_id: string
  paint_id: string
  purpose: string
  notes: string
  created_at: string
  paint?: Paint
}

export interface MiniatureImage {
  id: string
  miniature_id: string
  stage: PaintingStatus
  url: string
  caption: string
  created_at: string
}

export interface DashboardStats {
  total_minis: number
  finished_minis: number
  in_progress_minis: number
  unpainted_minis: number
  shame_percent: number
  by_status: Record<string, number>
  recent_activity: Miniature[]
}

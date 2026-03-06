import { useEffect, useState } from 'react'
import { confirmBox, createBox, createMiniature, deleteBox, deleteMiniature, getBoxes, getMiniatures, type MiniSuggestion } from '../api'
import type { Box, Game, Miniature } from '../types'
import type { Page } from '../App'
import { StatusBadge } from './DashboardPage'

interface ReviewState {
  box: Box
  rows: MiniSuggestion[]
  claudeError: string
  source: 'catalog' | 'claude' | 'none'
  saving: boolean
}

const SOURCE_LABEL: Record<ReviewState['source'], string> = {
  catalog: 'From shared catalog',
  claude: 'Suggested by Claude',
  none: 'No suggestions found',
}

export default function GameDetailPage({ game, nav }: { game: Game; nav: (p: Page) => void }) {
  const [boxes, setBoxes] = useState<Box[]>([])
  const [error, setError] = useState('')
  const [showForm, setShowForm] = useState(false)
  const [boxName, setBoxName] = useState('')
  const [creating, setCreating] = useState(false)
  const [review, setReview] = useState<ReviewState | null>(null)
  const [refreshKeys, setRefreshKeys] = useState<Record<string, number>>({})

  useEffect(() => {
    getBoxes(game.id).then(setBoxes).catch(e => setError(e.message))
  }, [game.id])

  async function handleCreateBox(e: React.FormEvent) {
    e.preventDefault()
    setCreating(true)
    try {
      const res = await createBox(game.id, { name: boxName })
      setBoxes(prev => [...prev, res.box].sort((a, b) => a.name.localeCompare(b.name)))
      setShowForm(false)
      setBoxName('')
      setReview({
        box: res.box,
        rows: (res.suggestions ?? []).map(s => ({ ...s, quantity: s.quantity || 1 })),
        claudeError: res.claude_error ?? '',
        source: res.source,
        saving: false,
      })
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Failed')
    } finally {
      setCreating(false)
    }
  }

  async function confirmReview() {
    if (!review) return
    setReview(r => r && ({ ...r, saving: true }))
    try {
      const boxId = review.box.id
      await confirmBox(
        boxId,
        review.rows.filter(r => r.name.trim()).map(r => ({ ...r, name: r.name.trim() })),
      )
      setRefreshKeys(prev => ({ ...prev, [boxId]: (prev[boxId] ?? 0) + 1 }))
      setReview(null)
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Failed to save minis')
      setReview(r => r && ({ ...r, saving: false }))
    }
  }

  async function handleDeleteBox(boxId: string) {
    if (!confirm('Delete this box and all its miniatures?')) return
    try {
      await deleteBox(boxId)
      setBoxes(prev => prev.filter(b => b.id !== boxId))
      if (review?.box.id === boxId) setReview(null)
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Failed')
    }
  }

  return (
    <div>
      <button onClick={() => nav({ name: 'collection', id: game.collection_id })} style={linkBtn}>← Back to collection</button>

      <div style={{ display: 'flex', alignItems: 'center', margin: '1rem 0' }}>
        <div>
          <h2 style={{ margin: 0 }}>{game.name}</h2>
          {(game.publisher || game.year) && (
            <div style={{ fontSize: '0.85rem', color: '#888', marginTop: 2 }}>
              {[game.publisher, game.year].filter(Boolean).join(' · ')}
            </div>
          )}
        </div>
        <button
          onClick={() => { setShowForm(s => !s); setReview(null) }}
          style={{ marginLeft: 'auto', ...btnStyle }}
        >
          {showForm ? 'Cancel' : '+ Add Box'}
        </button>
      </div>

      {error && <p style={{ color: 'red' }}>{error}</p>}

      {showForm && (
        <form onSubmit={handleCreateBox} style={formStyle}>
          <h3 style={{ margin: '0 0 0.5rem' }}>Add Box</h3>
          <input
            placeholder="Box name * (e.g. Core Set, Mystics of Midgard)"
            value={boxName}
            onChange={e => setBoxName(e.target.value)}
            required
            style={inputStyle}
          />
          <button type="submit" disabled={creating} style={btnStyle}>
            {creating ? 'Looking up minis…' : 'Add Box'}
          </button>
        </form>
      )}

      {review && (
        <MiniReview
          review={review}
          onChange={setReview}
          onConfirm={confirmReview}
          onSkip={() => setReview(null)}
        />
      )}

      {boxes.length === 0 && !showForm && !review && (
        <p style={{ color: '#666' }}>No boxes yet. Add a box or expansion to get started.</p>
      )}

      <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem', marginTop: '0.5rem' }}>
        {boxes.map(b => (
          <BoxSection
            key={`${b.id}-${refreshKeys[b.id] ?? 0}`}
            box={b}
            isReviewing={review?.box.id === b.id}
            nav={nav}
            onDelete={() => handleDeleteBox(b.id)}
          />
        ))}
      </div>
    </div>
  )
}

// ── BoxSection ────────────────────────────────────────────────────────────────

function BoxSection({ box, isReviewing, nav, onDelete }: {
  box: Box
  isReviewing: boolean
  nav: (p: Page) => void
  onDelete: () => void
}) {
  const [minis, setMinis] = useState<Miniature[]>([])
  const [loading, setLoading] = useState(true)
  const [showAddForm, setShowAddForm] = useState(false)
  const [name, setName] = useState('')
  const [unitType, setUnitType] = useState('')
  const [quantity, setQuantity] = useState(1)
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    getMiniatures(box.id)
      .then(ms => setMinis(ms ?? []))
      .finally(() => setLoading(false))
  }, [box.id])

  async function handleAddMini(e: React.FormEvent) {
    e.preventDefault()
    setSaving(true)
    try {
      const m = await createMiniature(box.id, { name, unit_type: unitType || undefined, quantity })
      setMinis(prev => [...prev, m].sort((a, b) => a.name.localeCompare(b.name)))
      setShowAddForm(false)
      setName(''); setUnitType(''); setQuantity(1)
    } finally {
      setSaving(false)
    }
  }

  async function handleDeleteMini(miniId: string, e: React.MouseEvent) {
    e.stopPropagation()
    if (!confirm('Delete this miniature?')) return
    await deleteMiniature(miniId)
    setMinis(prev => prev.filter(m => m.id !== miniId))
  }

  return (
    <div style={{
      border: isReviewing ? '2px solid #2980b9' : '1px solid #ddd',
      borderRadius: 8,
      overflow: 'hidden',
    }}>
      {/* Box header */}
      <div style={{ display: 'flex', alignItems: 'center', padding: '0.6rem 1rem', background: '#f7f7f7', borderBottom: '1px solid #eee', gap: '0.5rem' }}>
        <strong style={{ flex: 1 }}>{box.name}</strong>
        {isReviewing && <span style={{ fontSize: '0.75rem', color: '#2980b9' }}>reviewing…</span>}
        <button
          onClick={() => setShowAddForm(s => !s)}
          style={{ ...smallBtn, marginLeft: 'auto' }}
        >
          {showAddForm ? 'Cancel' : '+ Mini'}
        </button>
        <button
          onClick={() => nav({ name: 'box', id: box.id })}
          style={{ ...smallBtn, background: 'none', color: '#2980b9', border: '1px solid #2980b9' }}
          title="Open box detail (photos & paints)"
        >
          Manage
        </button>
        <button
          onClick={e => { e.stopPropagation(); onDelete() }}
          style={{ background: 'none', border: 'none', cursor: 'pointer', color: '#c0392b', padding: '0 0.25rem' }}
        >
          ✕
        </button>
      </div>

      {/* Add mini form */}
      {showAddForm && (
        <form onSubmit={handleAddMini} style={{ display: 'flex', gap: '0.5rem', padding: '0.75rem 1rem', background: '#fafafa', borderBottom: '1px solid #eee', flexWrap: 'wrap' }}>
          <input placeholder="Name *" value={name} onChange={e => setName(e.target.value)} required style={{ ...inputStyle, flex: 2, minWidth: 120 }} />
          <input placeholder="Unit type" value={unitType} onChange={e => setUnitType(e.target.value)} style={{ ...inputStyle, flex: 1, minWidth: 100 }} />
          <input type="number" min={1} value={quantity} onChange={e => setQuantity(Number(e.target.value))} style={{ ...inputStyle, width: 60 }} />
          <button type="submit" disabled={saving} style={btnStyle}>{saving ? '…' : 'Add'}</button>
        </form>
      )}

      {/* Mini list */}
      {loading ? (
        <p style={{ padding: '0.75rem 1rem', margin: 0, color: '#888', fontSize: '0.9rem' }}>Loading…</p>
      ) : minis.length === 0 ? (
        <p style={{ padding: '0.75rem 1rem', margin: 0, color: '#aaa', fontSize: '0.9rem' }}>No minis yet.</p>
      ) : (
        <table style={{ width: '100%', borderCollapse: 'collapse' }}>
          <thead>
            <tr style={{ textAlign: 'left', borderBottom: '1px solid #eee', background: '#fafafa' }}>
              <th style={th}>Name</th>
              <th style={th}>Type</th>
              <th style={{ ...th, width: 40 }}>Qty</th>
              <th style={th}>Status</th>
              <th style={{ ...th, width: 30 }}></th>
            </tr>
          </thead>
          <tbody>
            {minis.map(m => (
              <tr
                key={m.id}
                onClick={() => nav({ name: 'miniature', id: m.id })}
                style={{ cursor: 'pointer', borderBottom: '1px solid #f0f0f0' }}
              >
                <td style={td}>{m.name}</td>
                <td style={{ ...td, color: '#888' }}>{m.unit_type}</td>
                <td style={td}>{m.quantity}</td>
                <td style={td}><StatusBadge status={m.status} /></td>
                <td style={td}>
                  <button
                    onClick={e => handleDeleteMini(m.id, e)}
                    style={{ background: 'none', border: 'none', cursor: 'pointer', color: '#c0392b' }}
                  >
                    ✕
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  )
}

// ── MiniReview ────────────────────────────────────────────────────────────────

function MiniReview({ review, onChange, onConfirm, onSkip }: {
  review: ReviewState
  onChange: (r: ReviewState) => void
  onConfirm: () => void
  onSkip: () => void
}) {
  function updateRow(i: number, field: keyof MiniSuggestion, value: string | number) {
    const rows = review.rows.map((r, idx) => idx === i ? { ...r, [field]: value } : r)
    onChange({ ...review, rows })
  }

  return (
    <div style={{ border: '2px solid #2980b9', borderRadius: 8, padding: '1rem', marginBottom: '1rem' }}>
      <div style={{ display: 'flex', alignItems: 'center', marginBottom: '0.75rem' }}>
        <h3 style={{ margin: 0 }}>Review minis for "{review.box.name}"</h3>
        <span style={{ marginLeft: '0.75rem', fontSize: '0.8rem', color: review.source === 'catalog' ? '#27ae60' : '#888', background: '#f5f5f5', padding: '0.2rem 0.5rem', borderRadius: 4 }}>
          {SOURCE_LABEL[review.source]}
        </span>
      </div>

      {review.claudeError && <p style={{ color: '#c0392b', margin: '0 0 0.75rem' }}>Claude lookup failed: {review.claudeError}</p>}
      {!review.claudeError && review.rows.length === 0 && <p style={{ color: '#666', margin: '0 0 0.75rem' }}>No suggestions. Add minis manually below.</p>}

      {review.rows.length > 0 && (
        <table style={{ width: '100%', borderCollapse: 'collapse', marginBottom: '0.75rem' }}>
          <thead>
            <tr style={{ textAlign: 'left', borderBottom: '1px solid #ccc' }}>
              <th style={th}>Name</th><th style={th}>Unit Type</th><th style={{ ...th, width: 70 }}>Qty</th><th style={{ ...th, width: 30 }}></th>
            </tr>
          </thead>
          <tbody>
            {review.rows.map((row, i) => (
              <tr key={i} style={{ borderBottom: '1px solid #eee' }}>
                <td style={td}><input value={row.name} onChange={e => updateRow(i, 'name', e.target.value)} style={{ ...inputStyle, width: '100%', boxSizing: 'border-box' }} /></td>
                <td style={td}><input value={row.unit_type} onChange={e => updateRow(i, 'unit_type', e.target.value)} style={{ ...inputStyle, width: '100%', boxSizing: 'border-box' }} /></td>
                <td style={td}><input type="number" min={1} value={row.quantity} onChange={e => updateRow(i, 'quantity', Number(e.target.value))} style={{ ...inputStyle, width: 60 }} /></td>
                <td style={td}><button onClick={() => onChange({ ...review, rows: review.rows.filter((_, idx) => idx !== i) })} style={{ background: 'none', border: 'none', cursor: 'pointer', color: '#c0392b' }}>✕</button></td>
              </tr>
            ))}
          </tbody>
        </table>
      )}

      <div style={{ display: 'flex', gap: '0.5rem', alignItems: 'center' }}>
        <button onClick={onConfirm} disabled={review.saving} style={btnStyle}>
          {review.saving ? 'Saving…' : `Confirm & add ${review.rows.length} mini${review.rows.length !== 1 ? 's' : ''}`}
        </button>
        <button onClick={() => onChange({ ...review, rows: [...review.rows, { name: '', unit_type: '', quantity: 1 }] })} style={outlineBtn}>+ Row</button>
        <button onClick={onSkip} style={{ ...outlineBtn, marginLeft: 'auto' }}>Skip</button>
      </div>
    </div>
  )
}

const inputStyle: React.CSSProperties = { padding: '0.4rem 0.5rem', fontSize: '0.95rem', border: '1px solid #ccc', borderRadius: 4 }
const btnStyle: React.CSSProperties = { padding: '0.5rem 1rem', cursor: 'pointer', background: '#333', color: '#fff', border: 'none', borderRadius: 4, fontSize: '0.95rem' }
const smallBtn: React.CSSProperties = { padding: '0.25rem 0.6rem', cursor: 'pointer', background: '#555', color: '#fff', border: 'none', borderRadius: 4, fontSize: '0.8rem' }
const outlineBtn: React.CSSProperties = { padding: '0.4rem 0.75rem', cursor: 'pointer', background: '#fff', color: '#333', border: '1px solid #ccc', borderRadius: 4, fontSize: '0.9rem' }
const linkBtn: React.CSSProperties = { background: 'none', border: 'none', cursor: 'pointer', color: '#2980b9', fontSize: '0.95rem', padding: 0 }
const formStyle: React.CSSProperties = { border: '1px solid #ccc', borderRadius: 8, padding: '1rem', marginBottom: '1.5rem', display: 'flex', flexDirection: 'column', gap: '0.5rem' }
const th: React.CSSProperties = { padding: '0.35rem 0.75rem 0.35rem 0.75rem', fontWeight: 600, fontSize: '0.85rem', color: '#555' }
const td: React.CSSProperties = { padding: '0.4rem 0.75rem' }

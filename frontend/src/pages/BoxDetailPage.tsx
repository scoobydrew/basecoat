import { useEffect, useState } from 'react'
import { createMiniature, deleteMiniature, getMiniatures } from '../api'
import type { Box, Miniature } from '../types'
import type { Page } from '../App'
import { StatusBadge } from './DashboardPage'

export default function BoxDetailPage({ box, nav }: { box: Box; nav: (p: Page) => void }) {
  const [minis, setMinis] = useState<Miniature[]>([])
  const [error, setError] = useState('')
  const [showForm, setShowForm] = useState(false)
  const [name, setName] = useState('')
  const [unitType, setUnitType] = useState('')
  const [quantity, setQuantity] = useState(1)
  const [notes, setNotes] = useState('')
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    getMiniatures(box.id).then(setMinis).catch(e => setError(e.message))
  }, [box.id])

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault()
    setSaving(true)
    try {
      const m = await createMiniature(box.id, { name, unit_type: unitType || undefined, quantity, notes: notes || undefined })
      setMinis(prev => [...prev, m].sort((a, b) => a.name.localeCompare(b.name)))
      setShowForm(false)
      setName(''); setUnitType(''); setQuantity(1); setNotes('')
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Failed')
    } finally {
      setSaving(false)
    }
  }

  async function handleDelete(miniId: string, e: React.MouseEvent) {
    e.stopPropagation()
    if (!confirm('Delete this miniature?')) return
    try {
      await deleteMiniature(miniId)
      setMinis(prev => prev.filter(m => m.id !== miniId))
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Failed')
    }
  }

  return (
    <div>
      <button onClick={() => nav({ name: 'game', id: box.game_id })} style={linkBtn}>← Back to game</button>

      <div style={{ display: 'flex', alignItems: 'center', margin: '1rem 0' }}>
        <h2 style={{ margin: 0 }}>{box.name}</h2>
        <button onClick={() => setShowForm(s => !s)} style={{ marginLeft: 'auto', ...btnStyle }}>
          {showForm ? 'Cancel' : '+ Add Mini'}
        </button>
      </div>

      {error && <p style={{ color: 'red' }}>{error}</p>}

      {showForm && (
        <form onSubmit={handleCreate} style={formStyle}>
          <input placeholder="Name *" value={name} onChange={e => setName(e.target.value)} required style={inputStyle} />
          <input placeholder="Unit type (e.g. infantry, hero)" value={unitType} onChange={e => setUnitType(e.target.value)} style={inputStyle} />
          <label style={{ fontSize: '0.9rem' }}>
            Quantity
            <input type="number" min={1} value={quantity} onChange={e => setQuantity(Number(e.target.value))} style={{ ...inputStyle, marginLeft: '0.5rem', width: 60 }} />
          </label>
          <textarea placeholder="Notes" value={notes} onChange={e => setNotes(e.target.value)} rows={2} style={inputStyle} />
          <button type="submit" disabled={saving} style={btnStyle}>{saving ? 'Adding...' : 'Add'}</button>
        </form>
      )}

      {minis.length === 0 && !showForm && (
        <p style={{ color: '#666' }}>No miniatures yet.</p>
      )}

      <table style={{ width: '100%', borderCollapse: 'collapse' }}>
        {minis.length > 0 && (
          <thead>
            <tr style={{ textAlign: 'left', borderBottom: '1px solid #ccc' }}>
              <th style={th}>Name</th>
              <th style={th}>Type</th>
              <th style={th}>Qty</th>
              <th style={th}>Status</th>
              <th style={th}></th>
            </tr>
          </thead>
        )}
        <tbody>
          {minis.map(m => (
            <tr key={m.id} onClick={() => nav({ name: 'miniature', id: m.id })} style={{ cursor: 'pointer', borderBottom: '1px solid #eee' }}>
              <td style={td}>{m.name}</td>
              <td style={td}>{m.unit_type}</td>
              <td style={td}>{m.quantity}</td>
              <td style={td}><StatusBadge status={m.status} /></td>
              <td style={td}>
                <button onClick={e => handleDelete(m.id, e)} style={{ background: 'none', border: 'none', cursor: 'pointer', color: '#c0392b' }}>✕</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}

const inputStyle: React.CSSProperties = { padding: '0.5rem', fontSize: '1rem', border: '1px solid #ccc', borderRadius: 4 }
const btnStyle: React.CSSProperties = { padding: '0.5rem 1rem', cursor: 'pointer', background: '#333', color: '#fff', border: 'none', borderRadius: 4, fontSize: '0.95rem' }
const linkBtn: React.CSSProperties = { background: 'none', border: 'none', cursor: 'pointer', color: '#2980b9', fontSize: '0.95rem', padding: 0 }
const formStyle: React.CSSProperties = { padding: '0.75rem', background: '#f9f9f9', borderRadius: 6, marginBottom: '0.75rem', display: 'flex', flexDirection: 'column', gap: '0.5rem' }
const th: React.CSSProperties = { padding: '0.4rem 0.75rem 0.4rem 0', fontWeight: 'bold' }
const td: React.CSSProperties = { padding: '0.4rem 0.75rem 0.4rem 0' }

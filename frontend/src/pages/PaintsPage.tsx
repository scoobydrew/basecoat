import { useEffect, useState } from 'react'
import { createPaint, deletePaint, getPaints } from '../api'
import type { Paint } from '../types'

export default function PaintsPage() {
  const [paints, setPaints] = useState<Paint[]>([])
  const [error, setError] = useState('')
  const [showForm, setShowForm] = useState(false)
  const [brand, setBrand] = useState('')
  const [name, setName] = useState('')
  const [color, setColor] = useState('')
  const [type, setType] = useState('')
  const [saving, setSaving] = useState(false)
  const [filter, setFilter] = useState('')

  useEffect(() => {
    getPaints().then(setPaints).catch(e => setError(e.message))
  }, [])

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault()
    setSaving(true)
    try {
      const p = await createPaint({ brand, name, color: color || undefined, type: type || undefined })
      setPaints(prev => [...prev, p].sort((a, b) => a.brand.localeCompare(b.brand) || a.name.localeCompare(b.name)))
      setShowForm(false)
      setBrand(''); setName(''); setColor(''); setType('')
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Failed')
    } finally {
      setSaving(false)
    }
  }

  async function handleDelete(id: string) {
    if (!confirm('Remove this paint from your library?')) return
    try {
      await deletePaint(id)
      setPaints(prev => prev.filter(p => p.id !== id))
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Failed')
    }
  }

  const filtered = filter
    ? paints.filter(p => `${p.brand} ${p.name} ${p.type}`.toLowerCase().includes(filter.toLowerCase()))
    : paints

  return (
    <div>
      <div style={{ display: 'flex', alignItems: 'center', marginBottom: '1rem' }}>
        <h2 style={{ margin: 0 }}>Paint Library</h2>
        <button onClick={() => setShowForm(s => !s)} style={{ marginLeft: 'auto', ...btnStyle }}>
          {showForm ? 'Cancel' : '+ Add Paint'}
        </button>
      </div>

      {error && <p style={{ color: 'red' }}>{error}</p>}

      {showForm && (
        <form onSubmit={handleCreate} style={{ border: '1px solid #ccc', borderRadius: 8, padding: '1rem', marginBottom: '1.5rem', display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
          <h3 style={{ margin: '0 0 0.5rem' }}>New Paint</h3>
          <input placeholder="Brand * (e.g. Citadel, Vallejo)" value={brand} onChange={e => setBrand(e.target.value)} required style={inputStyle} />
          <input placeholder="Name * (e.g. Nuln Oil)" value={name} onChange={e => setName(e.target.value)} required style={inputStyle} />
          <input placeholder="Color (e.g. black)" value={color} onChange={e => setColor(e.target.value)} style={inputStyle} />
          <select value={type} onChange={e => setType(e.target.value)} style={inputStyle}>
            <option value="">Type (optional)</option>
            <option value="base">Base</option>
            <option value="shade">Shade</option>
            <option value="layer">Layer</option>
            <option value="highlight">Highlight</option>
            <option value="contrast">Contrast</option>
            <option value="technical">Technical</option>
            <option value="texture">Texture</option>
            <option value="dry">Dry</option>
            <option value="air">Air</option>
          </select>
          <button type="submit" disabled={saving} style={btnStyle}>{saving ? 'Saving...' : 'Add Paint'}</button>
        </form>
      )}

      <input
        placeholder="Filter paints…"
        value={filter}
        onChange={e => setFilter(e.target.value)}
        style={{ ...inputStyle, width: '100%', marginBottom: '1rem', boxSizing: 'border-box' }}
      />

      {filtered.length === 0 && <p style={{ color: '#666' }}>No paints found.</p>}

      <table style={{ width: '100%', borderCollapse: 'collapse' }}>
        {filtered.length > 0 && (
          <thead>
            <tr style={{ textAlign: 'left', borderBottom: '1px solid #ccc' }}>
              <th style={th}>Brand</th>
              <th style={th}>Name</th>
              <th style={th}>Color</th>
              <th style={th}>Type</th>
              <th style={th}></th>
            </tr>
          </thead>
        )}
        <tbody>
          {filtered.map(p => (
            <tr key={p.id} style={{ borderBottom: '1px solid #eee' }}>
              <td style={td}>{p.brand}</td>
              <td style={td}>{p.name}</td>
              <td style={td}>{p.color}</td>
              <td style={td}>{p.type}</td>
              <td style={td}>
                <button onClick={() => handleDelete(p.id)} style={{ background: 'none', border: 'none', cursor: 'pointer', color: '#c0392b' }}>✕</button>
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
const th: React.CSSProperties = { padding: '0.4rem 0.75rem 0.4rem 0', fontWeight: 'bold' }
const td: React.CSSProperties = { padding: '0.4rem 0.75rem 0.4rem 0' }

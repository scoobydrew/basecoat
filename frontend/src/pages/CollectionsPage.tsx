import { useEffect, useState } from 'react'
import { createCollection, deleteCollection, getCollections } from '../api'
import type { Collection } from '../types'
import type { Page } from '../App'

export default function CollectionsPage({ nav }: { nav: (p: Page) => void }) {
  const [collections, setCollections] = useState<Collection[]>([])
  const [error, setError] = useState('')
  const [showForm, setShowForm] = useState(false)
  const [name, setName] = useState('')
  const [notes, setNotes] = useState('')
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    getCollections().then(setCollections).catch(e => setError(e.message))
  }, [])

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault()
    setSaving(true)
    try {
      const col = await createCollection({ name, notes: notes || undefined })
      setCollections(prev => [...prev, col].sort((a, b) => a.name.localeCompare(b.name)))
      setShowForm(false)
      setName(''); setNotes('')
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Failed')
    } finally {
      setSaving(false)
    }
  }

  async function handleDelete(id: string, e: React.MouseEvent) {
    e.stopPropagation()
    if (!confirm('Delete this collection and everything in it?')) return
    try {
      await deleteCollection(id)
      setCollections(prev => prev.filter(c => c.id !== id))
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Failed')
    }
  }

  return (
    <div>
      <div style={{ display: 'flex', alignItems: 'center', marginBottom: '1.5rem' }}>
        <h2 style={{ margin: 0 }}>Collections</h2>
        <button onClick={() => setShowForm(s => !s)} style={{ marginLeft: 'auto', ...btnStyle }}>
          {showForm ? 'Cancel' : '+ New Collection'}
        </button>
      </div>

      {error && <p style={{ color: 'red' }}>{error}</p>}

      {showForm && (
        <form onSubmit={handleCreate} style={formStyle}>
          <h3 style={{ margin: '0 0 0.5rem' }}>New Collection</h3>
          <input placeholder="Name *" value={name} onChange={e => setName(e.target.value)} required style={inputStyle} />
          <textarea placeholder="Notes (optional)" value={notes} onChange={e => setNotes(e.target.value)} rows={2} style={inputStyle} />
          <button type="submit" disabled={saving} style={btnStyle}>{saving ? 'Creating...' : 'Create'}</button>
        </form>
      )}

      {collections.length === 0 && !showForm && (
        <p style={{ color: '#666' }}>No collections yet. Create one to get started.</p>
      )}

      <div style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
        {collections.map(c => (
          <div
            key={c.id}
            onClick={() => nav({ name: 'collection', id: c.id })}
            style={{ display: 'flex', alignItems: 'center', padding: '0.75rem 1rem', border: '1px solid #ddd', borderRadius: 8, cursor: 'pointer' }}
          >
            <div>
              <strong>{c.name}</strong>
              {c.notes && <p style={{ margin: '0.2rem 0 0', fontSize: '0.85rem', color: '#666' }}>{c.notes}</p>}
            </div>
            <button onClick={e => handleDelete(c.id, e)} style={{ marginLeft: 'auto', background: 'none', border: 'none', cursor: 'pointer', color: '#c0392b' }}>✕</button>
          </div>
        ))}
      </div>
    </div>
  )
}

const inputStyle: React.CSSProperties = { padding: '0.5rem', fontSize: '1rem', border: '1px solid #ccc', borderRadius: 4 }
const btnStyle: React.CSSProperties = { padding: '0.5rem 1rem', cursor: 'pointer', background: '#333', color: '#fff', border: 'none', borderRadius: 4, fontSize: '0.95rem' }
const formStyle: React.CSSProperties = { border: '1px solid #ccc', borderRadius: 8, padding: '1rem', marginBottom: '1.5rem', display: 'flex', flexDirection: 'column', gap: '0.5rem' }

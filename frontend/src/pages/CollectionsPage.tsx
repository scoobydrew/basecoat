import { useEffect, useState } from 'react'
import { createCollection, deleteCollection, getCollections } from '../api'
import type { Collection } from '../types'
import type { Page } from '../App'

export default function CollectionsPage({ nav }: { nav: (p: Page) => void }) {
  const [collections, setCollections] = useState<Collection[]>([])
  const [error, setError] = useState('')
  const [showForm, setShowForm] = useState(false)
  const [name, setName] = useState('')
  const [game, setGame] = useState('')
  const [set, setSet] = useState('')
  const [notes, setNotes] = useState('')
  const [creating, setCreating] = useState(false)

  useEffect(() => {
    getCollections().then(setCollections).catch(e => setError(e.message))
  }, [])

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault()
    setCreating(true)
    try {
      const res = await createCollection({ name, game, set: set || undefined, notes: notes || undefined })
      setCollections(prev => [res.collection, ...prev])
      setShowForm(false)
      setName(''); setGame(''); setSet(''); setNotes('')
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Failed to create')
    } finally {
      setCreating(false)
    }
  }

  async function handleDelete(id: string, e: React.MouseEvent) {
    e.stopPropagation()
    if (!confirm('Delete this collection and all its miniatures?')) return
    try {
      await deleteCollection(id)
      setCollections(prev => prev.filter(c => c.id !== id))
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Failed to delete')
    }
  }

  return (
    <div>
      <div style={{ display: 'flex', alignItems: 'center', marginBottom: '1rem' }}>
        <h2 style={{ margin: 0 }}>Collections</h2>
        <button onClick={() => setShowForm(s => !s)} style={{ marginLeft: 'auto', ...btnStyle }}>
          {showForm ? 'Cancel' : '+ New Collection'}
        </button>
      </div>

      {error && <p style={{ color: 'red' }}>{error}</p>}

      {showForm && (
        <form onSubmit={handleCreate} style={{ border: '1px solid #ccc', borderRadius: 8, padding: '1rem', marginBottom: '1.5rem', display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
          <h3 style={{ margin: '0 0 0.5rem' }}>New Collection</h3>
          <input placeholder="Collection name *" value={name} onChange={e => setName(e.target.value)} required style={inputStyle} />
          <input placeholder="Game / system *" value={game} onChange={e => setGame(e.target.value)} required style={inputStyle} />
          <input placeholder="Set or box (optional — Claude will populate minis)" value={set} onChange={e => setSet(e.target.value)} style={inputStyle} />
          <textarea placeholder="Notes" value={notes} onChange={e => setNotes(e.target.value)} rows={2} style={inputStyle} />
          <button type="submit" disabled={creating} style={btnStyle}>
            {creating ? 'Creating...' : 'Create'}
          </button>
          {creating && set && <small style={{ color: '#666' }}>Asking Claude to populate minis from "{set}"…</small>}
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
            style={{ border: '1px solid #ccc', borderRadius: 8, padding: '0.75rem 1rem', cursor: 'pointer', display: 'flex', alignItems: 'center' }}
          >
            <div>
              <strong>{c.name}</strong>
              <span style={{ color: '#666', marginLeft: '0.5rem', fontSize: '0.9rem' }}>{c.game}</span>
              {c.notes && <p style={{ margin: '0.25rem 0 0', fontSize: '0.85rem', color: '#555' }}>{c.notes}</p>}
            </div>
            <button
              onClick={e => handleDelete(c.id, e)}
              style={{ marginLeft: 'auto', background: 'none', border: 'none', cursor: 'pointer', color: '#c0392b', fontSize: '1.1rem' }}
              title="Delete"
            >
              ✕
            </button>
          </div>
        ))}
      </div>
    </div>
  )
}

const inputStyle: React.CSSProperties = { padding: '0.5rem', fontSize: '1rem', border: '1px solid #ccc', borderRadius: 4 }
const btnStyle: React.CSSProperties = { padding: '0.5rem 1rem', cursor: 'pointer', background: '#333', color: '#fff', border: 'none', borderRadius: 4, fontSize: '0.95rem' }

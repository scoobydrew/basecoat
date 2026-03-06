import { useEffect, useState } from 'react'
import { createGame, deleteGame, getBoxes, getCollection, getGames } from '../api'
import type { Box, Collection, Game } from '../types'
import type { Page } from '../App'

export default function CollectionDetailPage({ id, nav }: { id: string; nav: (p: Page) => void }) {
  const [collection, setCollection] = useState<Collection | null>(null)
  const [games, setGames] = useState<Game[]>([])
  const [error, setError] = useState('')
  const [showForm, setShowForm] = useState(false)
  const [gameName, setGameName] = useState('')
  const [publisher, setPublisher] = useState('')
  const [year, setYear] = useState('')
  const [saving, setSaving] = useState(false)
  const [expanded, setExpanded] = useState<Set<string>>(new Set())
  const [boxCache, setBoxCache] = useState<Record<string, Box[]>>({})

  useEffect(() => {
    Promise.all([getCollection(id), getGames(id)])
      .then(([col, gs]) => { setCollection(col); setGames(gs) })
      .catch(e => setError(e.message))
  }, [id])

  async function toggleGame(gameId: string) {
    const next = new Set(expanded)
    if (next.has(gameId)) {
      next.delete(gameId)
    } else {
      next.add(gameId)
      if (!boxCache[gameId]) {
        const boxes = await getBoxes(gameId).catch(() => [] as Box[])
        setBoxCache(c => ({ ...c, [gameId]: boxes }))
      }
    }
    setExpanded(next)
  }

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault()
    setSaving(true)
    try {
      const body: { name: string; publisher?: string; year?: number } = { name: gameName }
      if (publisher.trim()) body.publisher = publisher.trim()
      const y = parseInt(year)
      if (!isNaN(y) && y > 0) body.year = y
      const g = await createGame(id, body)
      setGames(prev => [...prev, g].sort((a, b) => a.name.localeCompare(b.name)))
      setShowForm(false)
      setGameName(''); setPublisher(''); setYear('')
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Failed')
    } finally {
      setSaving(false)
    }
  }

  async function handleDelete(gameId: string, e: React.MouseEvent) {
    e.stopPropagation()
    if (!confirm('Delete this game and all its boxes and miniatures?')) return
    try {
      await deleteGame(gameId)
      setGames(prev => prev.filter(g => g.id !== gameId))
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Failed')
    }
  }

  if (error) return <p style={{ color: 'red' }}>{error}</p>
  if (!collection) return <p>Loading...</p>

  return (
    <div>
      <button onClick={() => nav({ name: 'collections' })} style={linkBtn}>← Collections</button>

      <div style={{ display: 'flex', alignItems: 'center', margin: '1rem 0' }}>
        <h2 style={{ margin: 0 }}>{collection.name}</h2>
        <button onClick={() => setShowForm(s => !s)} style={{ marginLeft: 'auto', ...btnStyle }}>
          {showForm ? 'Cancel' : '+ Add Game'}
        </button>
      </div>

      {collection.notes && <p style={{ color: '#666', marginTop: 0 }}>{collection.notes}</p>}

      {showForm && (
        <form onSubmit={handleCreate} style={formStyle}>
          <h3 style={{ margin: '0 0 0.5rem' }}>Add Game</h3>
          <input placeholder="Game name * (e.g. Blood Rage, Warhammer 40K)" value={gameName} onChange={e => setGameName(e.target.value)} required style={inputStyle} />
          <div style={{ display: 'flex', gap: '0.5rem' }}>
            <input placeholder="Publisher (optional)" value={publisher} onChange={e => setPublisher(e.target.value)} style={{ ...inputStyle, flex: 1 }} />
            <input placeholder="Year (optional)" value={year} onChange={e => setYear(e.target.value)} type="number" min={1970} max={2100} style={{ ...inputStyle, width: 120 }} />
          </div>
          <small style={{ color: '#666' }}>Publisher and year help Claude suggest the right minis when you add boxes.</small>
          <button type="submit" disabled={saving} style={btnStyle}>{saving ? 'Adding...' : 'Add Game'}</button>
        </form>
      )}

      {games.length === 0 && !showForm && (
        <p style={{ color: '#666' }}>No games in this collection yet.</p>
      )}

      <div style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
        {games.map(g => {
          const open = expanded.has(g.id)
          const boxes = boxCache[g.id]
          return (
            <div key={g.id} style={{ border: '1px solid #ddd', borderRadius: 8, overflow: 'hidden' }}>
              {/* Game row */}
              <div style={{ display: 'flex', alignItems: 'center', padding: '0.75rem 1rem', cursor: 'pointer', background: open ? '#f7f7f7' : '#fff' }}
                onClick={() => toggleGame(g.id)}>
                <span style={{ marginRight: '0.5rem', color: '#888', fontSize: '0.8rem' }}>{open ? '▼' : '▶'}</span>
                <div style={{ flex: 1 }}>
                  <strong>{g.name}</strong>
                  {(g.publisher || g.year) && (
                    <span style={{ marginLeft: '0.5rem', fontSize: '0.8rem', color: '#888' }}>
                      {[g.publisher, g.year].filter(Boolean).join(' · ')}
                    </span>
                  )}
                </div>
                <button
                  onClick={e => { e.stopPropagation(); nav({ name: 'game', id: g.id }) }}
                  style={{ ...smallBtn, marginRight: '0.5rem' }}
                >
                  Open
                </button>
                <button onClick={e => handleDelete(g.id, e)} style={{ background: 'none', border: 'none', cursor: 'pointer', color: '#c0392b' }}>✕</button>
              </div>

              {/* Expanded boxes */}
              {open && (
                <div style={{ borderTop: '1px solid #eee' }}>
                  {!boxes ? (
                    <p style={{ margin: 0, padding: '0.6rem 1rem', color: '#aaa', fontSize: '0.9rem' }}>Loading…</p>
                  ) : boxes.length === 0 ? (
                    <p style={{ margin: 0, padding: '0.6rem 1rem', color: '#aaa', fontSize: '0.9rem' }}>No boxes yet.</p>
                  ) : (
                    boxes.map(b => (
                      <div
                        key={b.id}
                        onClick={() => nav({ name: 'game', id: g.id })}
                        style={{ display: 'flex', alignItems: 'center', padding: '0.5rem 1rem 0.5rem 2.5rem', borderBottom: '1px solid #f5f5f5', cursor: 'pointer' }}
                      >
                        <span style={{ flex: 1, fontSize: '0.95rem' }}>{b.name}</span>
                        <span style={{ fontSize: '0.75rem', color: '#aaa' }}>→</span>
                      </div>
                    ))
                  )}
                </div>
              )}
            </div>
          )
        })}
      </div>
    </div>
  )
}

const inputStyle: React.CSSProperties = { padding: '0.5rem', fontSize: '1rem', border: '1px solid #ccc', borderRadius: 4 }
const btnStyle: React.CSSProperties = { padding: '0.5rem 1rem', cursor: 'pointer', background: '#333', color: '#fff', border: 'none', borderRadius: 4, fontSize: '0.95rem' }
const smallBtn: React.CSSProperties = { padding: '0.25rem 0.6rem', cursor: 'pointer', background: '#555', color: '#fff', border: 'none', borderRadius: 4, fontSize: '0.8rem' }
const linkBtn: React.CSSProperties = { background: 'none', border: 'none', cursor: 'pointer', color: '#2980b9', fontSize: '0.95rem', padding: 0 }
const formStyle: React.CSSProperties = { border: '1px solid #ccc', borderRadius: 8, padding: '1rem', marginBottom: '1.5rem', display: 'flex', flexDirection: 'column', gap: '0.5rem' }

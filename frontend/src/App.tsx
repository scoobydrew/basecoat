import { useEffect, useState } from 'react'
import AuthPage from './pages/AuthPage'
import DashboardPage from './pages/DashboardPage'
import CollectionsPage from './pages/CollectionsPage'
import CollectionDetailPage from './pages/CollectionDetailPage'
import GameDetailPage from './pages/GameDetailPage'
import BoxDetailPage from './pages/BoxDetailPage'
import MiniatureDetailPage from './pages/MiniatureDetailPage'
import PaintsPage from './pages/PaintsPage'
import { getBox, getGame, logout } from './api'
import type { Box, Game } from './types'

export type Page =
  | { name: 'dashboard' }
  | { name: 'collections' }
  | { name: 'collection'; id: string }
  | { name: 'game'; id: string }
  | { name: 'box'; id: string }
  | { name: 'miniature'; id: string }
  | { name: 'paints' }

export default function App() {
  const [authed, setAuthed] = useState(!!localStorage.getItem('token'))
  const [page, setPage] = useState<Page>({ name: 'dashboard' })

  useEffect(() => {
    const handler = () => setAuthed(!!localStorage.getItem('token'))
    window.addEventListener('auth', handler)
    return () => window.removeEventListener('auth', handler)
  }, [])

  if (!authed) return <AuthPage onAuth={() => setAuthed(true)} />

  function nav(p: Page) { setPage(p) }

  function handleLogout() {
    logout()
    setAuthed(false)
  }

  const collectionsActive = ['collections', 'collection', 'game', 'miniature'].includes(page.name)

  return (
    <div style={{ fontFamily: 'sans-serif', maxWidth: 900, margin: '0 auto', padding: '0 1rem' }}>
      <nav style={{ display: 'flex', gap: '1rem', padding: '1rem 0', borderBottom: '1px solid #ccc', marginBottom: '1.5rem', alignItems: 'center' }}>
        <strong style={{ marginRight: '0.5rem' }}>🎨 Basecoat</strong>
        <NavLink label="Dashboard" active={page.name === 'dashboard'} onClick={() => nav({ name: 'dashboard' })} />
        <NavLink label="Collections" active={collectionsActive} onClick={() => nav({ name: 'collections' })} />
        <NavLink label="Paints" active={page.name === 'paints'} onClick={() => nav({ name: 'paints' })} />
        <button onClick={handleLogout} style={{ marginLeft: 'auto', cursor: 'pointer' }}>Logout</button>
      </nav>

      {page.name === 'dashboard' && <DashboardPage nav={nav} />}
      {page.name === 'collections' && <CollectionsPage nav={nav} />}
      {page.name === 'collection' && <CollectionDetailPage id={page.id} nav={nav} />}
      {page.name === 'game' && <GameLoader id={page.id} nav={nav} />}
      {page.name === 'box' && <BoxLoader id={page.id} nav={nav} />}
      {page.name === 'miniature' && <MiniatureDetailPage id={page.id} nav={nav} />}
      {page.name === 'paints' && <PaintsPage />}
    </div>
  )
}

function BoxLoader({ id, nav }: { id: string; nav: (p: Page) => void }) {
  const [box, setBox] = useState<Box | null>(null)
  const [error, setError] = useState('')

  useEffect(() => {
    getBox(id).then(setBox).catch(e => setError(e.message))
  }, [id])

  if (error) return <p style={{ color: 'red' }}>{error}</p>
  if (!box) return <p>Loading...</p>
  return <BoxDetailPage box={box} nav={nav} />
}

function GameLoader({ id, nav }: { id: string; nav: (p: Page) => void }) {
  const [game, setGame] = useState<Game | null>(null)
  const [error, setError] = useState('')

  useEffect(() => {
    getGame(id).then(setGame).catch(e => setError(e.message))
  }, [id])

  if (error) return <p style={{ color: 'red' }}>{error}</p>
  if (!game) return <p>Loading...</p>
  return <GameDetailPage game={game} nav={nav} />
}

function NavLink({ label, active, onClick }: { label: string; active: boolean; onClick: () => void }) {
  return (
    <button
      onClick={onClick}
      style={{ background: 'none', border: 'none', cursor: 'pointer', fontWeight: active ? 'bold' : 'normal', textDecoration: active ? 'underline' : 'none', fontSize: '1rem' }}
    >
      {label}
    </button>
  )
}

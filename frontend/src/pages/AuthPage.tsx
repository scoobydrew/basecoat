import { useState } from 'react'
import { login, register } from '../api'

export default function AuthPage({ onAuth }: { onAuth: () => void }) {
  const [mode, setMode] = useState<'login' | 'register'>('login')
  const [username, setUsername] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  async function submit(e: React.FormEvent) {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      if (mode === 'login') {
        await login(email, password)
      } else {
        await register(username, email, password)
      }
      onAuth()
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Something went wrong')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div style={{ maxWidth: 360, margin: '4rem auto', fontFamily: 'sans-serif' }}>
      <h1 style={{ textAlign: 'center' }}>🎨 Basecoat</h1>
      <div style={{ display: 'flex', marginBottom: '1rem' }}>
        <TabButton label="Login" active={mode === 'login'} onClick={() => setMode('login')} />
        <TabButton label="Register" active={mode === 'register'} onClick={() => setMode('register')} />
      </div>
      <form onSubmit={submit} style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem' }}>
        {mode === 'register' && (
          <input placeholder="Username" value={username} onChange={e => setUsername(e.target.value)} required style={inputStyle} />
        )}
        <input placeholder="Email" type="email" value={email} onChange={e => setEmail(e.target.value)} required style={inputStyle} />
        <input placeholder="Password" type="password" value={password} onChange={e => setPassword(e.target.value)} required style={inputStyle} />
        {error && <p style={{ color: 'red', margin: 0 }}>{error}</p>}
        <button type="submit" disabled={loading} style={btnStyle}>
          {loading ? '...' : mode === 'login' ? 'Login' : 'Create account'}
        </button>
      </form>
    </div>
  )
}

function TabButton({ label, active, onClick }: { label: string; active: boolean; onClick: () => void }) {
  return (
    <button onClick={onClick} style={{ flex: 1, padding: '0.5rem', cursor: 'pointer', background: active ? '#333' : '#eee', color: active ? '#fff' : '#333', border: 'none', fontWeight: active ? 'bold' : 'normal' }}>
      {label}
    </button>
  )
}

const inputStyle: React.CSSProperties = { padding: '0.5rem', fontSize: '1rem', border: '1px solid #ccc', borderRadius: 4 }
const btnStyle: React.CSSProperties = { padding: '0.6rem', fontSize: '1rem', cursor: 'pointer', background: '#333', color: '#fff', border: 'none', borderRadius: 4 }

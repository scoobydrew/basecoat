import { useEffect, useState } from 'react'
import { getDashboard } from '../api'
import type { DashboardStats } from '../types'
import type { Page } from '../App'

const STATUS_LABELS: Record<string, string> = {
  unpainted: 'Unpainted',
  primed: 'Primed',
  basecoated: 'Basecoated',
  shaded: 'Shaded',
  detailed: 'Detailed',
  finished: 'Finished',
}

export function StatusBadge({ status }: { status: string }) {
  const colors: Record<string, string> = {
    unpainted: '#95a5a6',
    primed: '#7f8c8d',
    basecoated: '#2980b9',
    shaded: '#8e44ad',
    detailed: '#e67e22',
    finished: '#27ae60',
  }
  return (
    <span style={{ background: colors[status] ?? '#ccc', color: '#fff', borderRadius: 4, padding: '0.1rem 0.4rem', fontSize: '0.8rem' }}>
      {STATUS_LABELS[status] ?? status}
    </span>
  )
}

export default function DashboardPage({ nav }: { nav: (p: Page) => void }) {
  const [stats, setStats] = useState<DashboardStats | null>(null)
  const [error, setError] = useState('')

  useEffect(() => {
    getDashboard().then(setStats).catch(e => setError(e.message))
  }, [])

  if (error) return <p style={{ color: 'red' }}>{error}</p>
  if (!stats) return <p>Loading...</p>

  const shameColor = stats.shame_percent > 66 ? '#c0392b' : stats.shame_percent > 33 ? '#e67e22' : '#27ae60'

  return (
    <div>
      <h2>Dashboard</h2>

      <div style={{ display: 'flex', gap: '1rem', flexWrap: 'wrap', marginBottom: '2rem' }}>
        <StatCard label="Total Minis" value={stats.total_minis} />
        <StatCard label="Finished" value={stats.finished_minis} />
        <StatCard label="In Progress" value={stats.in_progress_minis} />
        <StatCard label="Unpainted" value={stats.unpainted_minis} />
      </div>

      <div style={{ marginBottom: '2rem' }}>
        <h3 style={{ marginBottom: '0.5rem' }}>
          Grey Pile of Shame — {stats.shame_percent.toFixed(1)}%
        </h3>
        <div style={{ background: '#eee', borderRadius: 8, height: 24, overflow: 'hidden' }}>
          <div style={{ width: `${stats.shame_percent}%`, height: '100%', background: shameColor, transition: 'width 0.4s' }} />
        </div>
        <small style={{ color: '#666' }}>{stats.unpainted_minis} unpainted or primed out of {stats.total_minis} total</small>
      </div>

      <div style={{ marginBottom: '2rem' }}>
        <h3>By Status</h3>
        <table style={{ borderCollapse: 'collapse' }}>
          <tbody>
            {Object.entries(STATUS_LABELS).map(([key, label]) => (
              <tr key={key}>
                <td style={{ padding: '0.25rem 1rem 0.25rem 0' }}>{label}</td>
                <td style={{ padding: '0.25rem 0', fontWeight: 'bold' }}>{stats.by_status[key] ?? 0}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {stats.recent_activity?.length > 0 && (
        <div>
          <h3>Recent Activity</h3>
          <table style={{ width: '100%', borderCollapse: 'collapse' }}>
            <thead>
              <tr style={{ textAlign: 'left', borderBottom: '1px solid #ccc' }}>
                <th style={th}>Name</th>
                <th style={th}>Status</th>
                <th style={th}>Updated</th>
              </tr>
            </thead>
            <tbody>
              {stats.recent_activity.map(m => (
                <tr key={m.id} onClick={() => nav({ name: 'miniature', id: m.id })} style={{ cursor: 'pointer', borderBottom: '1px solid #eee' }}>
                  <td style={td}>{m.name}</td>
                  <td style={td}><StatusBadge status={m.status} /></td>
                  <td style={td}>{new Date(m.updated_at).toLocaleDateString()}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}

function StatCard({ label, value }: { label: string; value: number }) {
  return (
    <div style={{ border: '1px solid #ccc', borderRadius: 8, padding: '1rem 1.5rem', minWidth: 120, textAlign: 'center' }}>
      <div style={{ fontSize: '2rem', fontWeight: 'bold' }}>{value}</div>
      <div style={{ color: '#666', fontSize: '0.85rem' }}>{label}</div>
    </div>
  )
}

const th: React.CSSProperties = { padding: '0.4rem 0.75rem 0.4rem 0', fontWeight: 'bold' }
const td: React.CSSProperties = { padding: '0.4rem 0.75rem 0.4rem 0' }

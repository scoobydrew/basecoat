import { useEffect, useRef, useState } from 'react'
import { addMiniaturePaint, deleteImage, getMiniature, getPaints, removeMiniaturePaint, updateMiniature, uploadImage } from '../api'
import type { Miniature, MiniaturePaint, Paint, PaintingStatus } from '../types'
import type { Page } from '../App'
import { StatusBadge } from './DashboardPage'

const STATUSES: PaintingStatus[] = ['unpainted', 'primed', 'basecoated', 'shaded', 'detailed', 'finished']

export default function MiniatureDetailPage({ id, nav }: { id: string; nav: (p: Page) => void }) {
  const [mini, setMini] = useState<Miniature | null>(null)
  const [paints, setPaints] = useState<Paint[]>([])
  const [error, setError] = useState('')
  const [editingNotes, setEditingNotes] = useState(false)
  const [notes, setNotes] = useState('')
  const [savingNotes, setSavingNotes] = useState(false)

  // Add paint form
  const [showPaintForm, setShowPaintForm] = useState(false)
  const [selectedPaint, setSelectedPaint] = useState('')
  const [paintPurpose, setPaintPurpose] = useState('')
  const [paintNotes, setPaintNotes] = useState('')
  const [addingPaint, setAddingPaint] = useState(false)

  // Image upload
  const fileRef = useRef<HTMLInputElement>(null)
  const [uploadStage, setUploadStage] = useState<PaintingStatus>('unpainted')
  const [uploadCaption, setUploadCaption] = useState('')
  const [uploading, setUploading] = useState(false)

  useEffect(() => {
    Promise.all([getMiniature(id), getPaints()])
      .then(([m, ps]) => { setMini(m); setNotes(m.notes); setPaints(ps) })
      .catch(e => setError(e.message))
  }, [id])

  async function handleStatusChange(status: PaintingStatus) {
    if (!mini) return
    const updated = await updateMiniature(mini.id, { status }).catch(e => { setError(e.message); return null })
    if (updated) setMini({ ...mini, ...updated })
  }

  async function saveNotes() {
    if (!mini) return
    setSavingNotes(true)
    const updated = await updateMiniature(mini.id, { notes }).catch(e => { setError(e.message); return null })
    if (updated) { setMini({ ...mini, ...updated }); setEditingNotes(false) }
    setSavingNotes(false)
  }

  async function handleAddPaint(e: React.FormEvent) {
    e.preventDefault()
    if (!mini || !selectedPaint) return
    setAddingPaint(true)
    try {
      const mp = await addMiniaturePaint(mini.id, { paint_id: selectedPaint, purpose: paintPurpose, notes: paintNotes })
      const paint = paints.find(p => p.id === selectedPaint)
      const enriched: MiniaturePaint = { ...mp, paint }
      setMini({ ...mini, paints: [...(mini.paints ?? []), enriched] })
      setShowPaintForm(false); setSelectedPaint(''); setPaintPurpose(''); setPaintNotes('')
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Failed')
    } finally {
      setAddingPaint(false)
    }
  }

  async function handleRemovePaint(linkID: string) {
    if (!mini) return
    await removeMiniaturePaint(mini.id, linkID).catch(e => setError(e.message))
    setMini({ ...mini, paints: (mini.paints ?? []).filter(p => p.id !== linkID) })
  }

  async function handleUpload() {
    if (!mini || !fileRef.current?.files?.[0]) return
    setUploading(true)
    try {
      const img = await uploadImage(mini.id, fileRef.current.files[0], uploadStage, uploadCaption)
      setMini({ ...mini, images: [...(mini.images ?? []), img] })
      setUploadCaption('')
      if (fileRef.current) fileRef.current.value = ''
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Upload failed')
    } finally {
      setUploading(false)
    }
  }

  async function handleDeleteImage(imageID: string) {
    if (!mini) return
    await deleteImage(mini.id, imageID).catch(e => setError(e.message))
    setMini({ ...mini, images: (mini.images ?? []).filter(i => i.id !== imageID) })
  }

  if (error) return <p style={{ color: 'red' }}>{error}</p>
  if (!mini) return <p>Loading...</p>

  return (
    <div>
      <button onClick={() => nav({ name: 'collection', id: mini.collection_id })} style={linkBtn}>← Back to collection</button>

      <div style={{ display: 'flex', alignItems: 'baseline', gap: '1rem', margin: '1rem 0' }}>
        <h2 style={{ margin: 0 }}>{mini.name}</h2>
        {mini.unit_type && <span style={{ color: '#666' }}>{mini.unit_type}</span>}
        {mini.quantity > 1 && <span style={{ color: '#666' }}>×{mini.quantity}</span>}
      </div>

      {/* Status */}
      <section style={section}>
        <h3 style={sectionTitle}>Status</h3>
        <div style={{ display: 'flex', gap: '0.5rem', flexWrap: 'wrap' }}>
          {STATUSES.map(s => (
            <button
              key={s}
              onClick={() => handleStatusChange(s)}
              style={{
                padding: '0.35rem 0.75rem',
                borderRadius: 4,
                border: mini.status === s ? '2px solid #333' : '1px solid #ccc',
                cursor: 'pointer',
                background: mini.status === s ? '#333' : '#fff',
                color: mini.status === s ? '#fff' : '#333',
                fontWeight: mini.status === s ? 'bold' : 'normal',
              }}
            >
              {s}
            </button>
          ))}
        </div>
        <div style={{ marginTop: '0.5rem' }}><StatusBadge status={mini.status} /></div>
      </section>

      {/* Notes */}
      <section style={section}>
        <h3 style={sectionTitle}>Notes</h3>
        {editingNotes ? (
          <div style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem' }}>
            <textarea value={notes} onChange={e => setNotes(e.target.value)} rows={4} style={{ ...inputStyle, resize: 'vertical' }} />
            <div style={{ display: 'flex', gap: '0.5rem' }}>
              <button onClick={saveNotes} disabled={savingNotes} style={btnStyle}>{savingNotes ? 'Saving...' : 'Save'}</button>
              <button onClick={() => { setEditingNotes(false); setNotes(mini.notes) }} style={outlineBtn}>Cancel</button>
            </div>
          </div>
        ) : (
          <div>
            <p style={{ margin: '0 0 0.5rem', whiteSpace: 'pre-wrap', color: mini.notes ? '#333' : '#999' }}>
              {mini.notes || 'No notes yet.'}
            </p>
            <button onClick={() => setEditingNotes(true)} style={outlineBtn}>Edit notes</button>
          </div>
        )}
      </section>

      {/* Paints */}
      <section style={section}>
        <div style={{ display: 'flex', alignItems: 'center', marginBottom: '0.5rem' }}>
          <h3 style={{ ...sectionTitle, margin: 0 }}>Paints Used</h3>
          <button onClick={() => setShowPaintForm(s => !s)} style={{ marginLeft: 'auto', ...outlineBtn }}>
            {showPaintForm ? 'Cancel' : '+ Add paint'}
          </button>
        </div>

        {showPaintForm && (
          <form onSubmit={handleAddPaint} style={{ display: 'flex', flexDirection: 'column', gap: '0.5rem', marginBottom: '0.75rem', padding: '0.75rem', background: '#f9f9f9', borderRadius: 6 }}>
            <select value={selectedPaint} onChange={e => setSelectedPaint(e.target.value)} required style={inputStyle}>
              <option value="">Select paint…</option>
              {paints.map(p => (
                <option key={p.id} value={p.id}>{p.brand} — {p.name}</option>
              ))}
            </select>
            <input placeholder="Purpose (e.g. base coat, shade)" value={paintPurpose} onChange={e => setPaintPurpose(e.target.value)} style={inputStyle} />
            <input placeholder="Notes" value={paintNotes} onChange={e => setPaintNotes(e.target.value)} style={inputStyle} />
            <button type="submit" disabled={addingPaint} style={btnStyle}>{addingPaint ? 'Adding...' : 'Add'}</button>
            {paints.length === 0 && <small style={{ color: '#c0392b' }}>No paints in your library yet — add some in the Paints tab first.</small>}
          </form>
        )}

        {(mini.paints ?? []).length === 0 && !showPaintForm && <p style={{ color: '#999' }}>No paints logged yet.</p>}

        <div style={{ display: 'flex', flexDirection: 'column', gap: '0.4rem' }}>
          {(mini.paints ?? []).map(mp => (
            <div key={mp.id} style={{ display: 'flex', alignItems: 'center', gap: '0.5rem', padding: '0.4rem 0.5rem', background: '#f5f5f5', borderRadius: 4 }}>
              <span style={{ fontWeight: 'bold' }}>{mp.paint?.brand} — {mp.paint?.name}</span>
              {mp.purpose && <span style={{ color: '#666', fontSize: '0.85rem' }}>({mp.purpose})</span>}
              {mp.notes && <span style={{ color: '#888', fontSize: '0.85rem' }}>· {mp.notes}</span>}
              <button onClick={() => handleRemovePaint(mp.id)} style={{ marginLeft: 'auto', background: 'none', border: 'none', cursor: 'pointer', color: '#c0392b' }}>✕</button>
            </div>
          ))}
        </div>
      </section>

      {/* Images */}
      <section style={section}>
        <h3 style={sectionTitle}>Photos</h3>
        <div style={{ display: 'flex', gap: '0.5rem', alignItems: 'center', flexWrap: 'wrap', marginBottom: '0.75rem' }}>
          <input type="file" accept="image/*" ref={fileRef} style={{ fontSize: '0.9rem' }} />
          <select value={uploadStage} onChange={e => setUploadStage(e.target.value as PaintingStatus)} style={inputStyle}>
            {STATUSES.map(s => <option key={s} value={s}>{s}</option>)}
          </select>
          <input placeholder="Caption" value={uploadCaption} onChange={e => setUploadCaption(e.target.value)} style={{ ...inputStyle, flex: 1 }} />
          <button onClick={handleUpload} disabled={uploading} style={btnStyle}>{uploading ? 'Uploading...' : 'Upload'}</button>
        </div>

        {(mini.images ?? []).length === 0 && <p style={{ color: '#999' }}>No photos yet.</p>}

        <div style={{ display: 'flex', gap: '0.75rem', flexWrap: 'wrap' }}>
          {(mini.images ?? []).map(img => (
            <div key={img.id} style={{ position: 'relative', width: 150 }}>
              <img src={img.url} alt={img.caption} style={{ width: '100%', height: 120, objectFit: 'cover', borderRadius: 6, display: 'block' }} />
              <div style={{ fontSize: '0.75rem', color: '#666', marginTop: 2 }}>
                <StatusBadge status={img.stage} /> {img.caption}
              </div>
              <button
                onClick={() => handleDeleteImage(img.id)}
                style={{ position: 'absolute', top: 4, right: 4, background: 'rgba(0,0,0,0.5)', color: '#fff', border: 'none', borderRadius: '50%', width: 20, height: 20, cursor: 'pointer', lineHeight: '20px', padding: 0, fontSize: '0.75rem' }}
              >✕</button>
            </div>
          ))}
        </div>
      </section>
    </div>
  )
}

const section: React.CSSProperties = { marginBottom: '2rem', paddingBottom: '1.5rem', borderBottom: '1px solid #eee' }
const sectionTitle: React.CSSProperties = { marginTop: 0, marginBottom: '0.75rem' }
const inputStyle: React.CSSProperties = { padding: '0.45rem 0.5rem', fontSize: '0.95rem', border: '1px solid #ccc', borderRadius: 4 }
const btnStyle: React.CSSProperties = { padding: '0.45rem 1rem', cursor: 'pointer', background: '#333', color: '#fff', border: 'none', borderRadius: 4, fontSize: '0.95rem' }
const outlineBtn: React.CSSProperties = { padding: '0.4rem 0.75rem', cursor: 'pointer', background: '#fff', color: '#333', border: '1px solid #ccc', borderRadius: 4, fontSize: '0.9rem' }
const linkBtn: React.CSSProperties = { background: 'none', border: 'none', cursor: 'pointer', color: '#2980b9', fontSize: '0.95rem', padding: 0 }

import { useRef, useState, type DragEvent } from 'react'

interface FileUploadProps {
  onUpload: (file: File) => void
  disabled?: boolean
}

export function FileUpload({ onUpload, disabled = false }: FileUploadProps) {
  const inputRef = useRef<HTMLInputElement>(null)
  const [dragging, setDragging] = useState(false)
  const [error, setError] = useState<string | null>(null)

  function select(file: File | undefined) {
    if (!file) {
      return
    }
    if (file.type !== 'application/pdf' && !file.name.toLowerCase().endsWith('.pdf')) {
      setError('Please choose a PDF file.')
      return
    }
    setError(null)
    onUpload(file)
  }

  function handleDrop(event: DragEvent<HTMLDivElement>) {
    event.preventDefault()
    setDragging(false)
    select(event.dataTransfer.files[0])
  }

  return (
    <div className="file-upload">
      <div
        className={`dropzone${dragging ? ' dropzone--active' : ''}`}
        onDragOver={(event) => {
          event.preventDefault()
          setDragging(true)
        }}
        onDragLeave={() => setDragging(false)}
        onDrop={handleDrop}
      >
        <p className="dropzone__title">Drop your resume here</p>
        <p className="dropzone__hint">PDF only, up to 5 MB</p>
        <button
          type="button"
          className="button"
          disabled={disabled}
          onClick={() => inputRef.current?.click()}
        >
          {disabled ? 'Uploading…' : 'Choose PDF'}
        </button>
        <input
          ref={inputRef}
          type="file"
          accept="application/pdf,.pdf"
          hidden
          onChange={(event) => {
            select(event.target.files?.[0])
            event.target.value = ''
          }}
        />
      </div>
      {error && <p className="message message--error">{error}</p>}
    </div>
  )
}

import { FileUpload } from '../components/FileUpload'
import { useUpload } from '../hooks/useUpload'
import { formatBytes } from '../utils/formatBytes'

export function UploadPage() {
  const { status, result, error, upload, reset } = useUpload()
  const uploading = status === 'uploading'

  return (
    <section className="upload-page">
      <h1>Upload your resume</h1>
      <p className="upload-page__lead">
        We will read your resume and help you build a sharper version of it.
      </p>

      <FileUpload onUpload={upload} disabled={uploading} />

      {status === 'success' && result && (
        <div className="message message--success">
          <span>
            Uploaded <strong>{result.filename}</strong> ({formatBytes(result.size)})
          </span>
          <button type="button" className="button button--ghost" onClick={reset}>
            Upload another
          </button>
        </div>
      )}

      {status === 'error' && (
        <div className="message message--error">
          <span>{error}</span>
          <button type="button" className="button button--ghost" onClick={reset}>
            Try again
          </button>
        </div>
      )}
    </section>
  )
}

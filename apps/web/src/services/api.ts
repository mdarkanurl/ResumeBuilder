import type { UploadResult } from '../types'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? '/api'

export async function uploadResume(file: File): Promise<UploadResult> {
  const body = new FormData()
  body.append('file', file)

  const response = await fetch(`${API_BASE_URL}/upload`, {
    method: 'POST',
    body,
  })

  if (!response.ok) {
    const message = (await response.text()).trim()
    throw new Error(message || `Upload failed (${response.status})`)
  }

  return (await response.json()) as UploadResult
}

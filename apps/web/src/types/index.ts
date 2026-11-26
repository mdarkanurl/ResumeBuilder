export interface UploadResult {
  filename: string
  size: number
}

export type UploadStatus = 'idle' | 'uploading' | 'success' | 'error'

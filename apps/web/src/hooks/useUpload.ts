import { useCallback, useState } from 'react'
import { uploadResume } from '../services/api'
import type { UploadResult, UploadStatus } from '../types'

interface UploadState {
  status: UploadStatus
  result?: UploadResult
  error?: string
}

export function useUpload() {
  const [state, setState] = useState<UploadState>({ status: 'idle' })

  const upload = useCallback(async (file: File) => {
    setState({ status: 'uploading' })
    try {
      const result = await uploadResume(file)
      setState({ status: 'success', result })
    } catch (error) {
      setState({
        status: 'error',
        error: error instanceof Error ? error.message : 'Upload failed',
      })
    }
  }, [])

  const reset = useCallback(() => setState({ status: 'idle' }), [])

  return { ...state, upload, reset }
}

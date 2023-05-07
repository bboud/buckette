type FileDB = {
  FileName: string
  UUID: ArrayBuffer
  URL: string
  Size: number
  Uploaded: string
  UserUploaded: string
  DownloadCount: number
}

type FileStatus = {
  name: string
  type: string
  status: string
  progress: number
  length: number
  url: string
  error?: string
}

type FileResponse = {
  URL: string
  Duplicate: boolean
  File: FileDB
}

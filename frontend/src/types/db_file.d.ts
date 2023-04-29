type FileDB = {
  FileName: string
  UUID: string
  Size: number
  Uploaded: string
  UserUploaded: string
  DownloadCount: number
}

type FileStatus = {
  name: string
  status: string
  progress: number
  length: number
  url: string
}

type FileResponse = {
  FileName: string
  RecordHash: string
  FileHash: string
  Size: number
  ContentType: string
  Uploaded: string
  UserUploaded: string
  DownloadCount: number
}

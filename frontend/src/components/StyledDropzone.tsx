import { useMutation } from '@tanstack/react-query'
import { useCallback, useEffect, useMemo, useState } from 'react'
import { FileWithPath, useDropzone, DropzoneRootProps } from 'react-dropzone'
import { Link, useNavigate } from 'react-router-dom'

import ClimbingBoxLoader from 'react-spinners/ClimbingBoxLoader'

import axios from 'axios'
import { fileStatusAtom } from '../utils/atoms'

import { animated as a, useSpring } from '@react-spring/web'

import { useAtom } from 'jotai'

const baseStyle: DropzoneRootProps = {
  flex: 1,
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  padding: '20px',
  borderWidth: 2,
  borderRadius: 2,
  borderColor: '#eeeeee',
  borderStyle: 'dashed',
  backgroundColor: '#fafafa',
  color: '#bdbdbd',
  outline: 'none',
  transition: 'border .24s ease-in-out',
}

const focusedStyle = {
  borderColor: '#2196f3',
}

const acceptStyle = {
  borderColor: '#00e676',
}

const rejectStyle = {
  borderColor: '#ff1744',
}

export default function StyledDropzone() {
  const [isUploading, setUploading] = useState(false)

  const [fileData, setFileData] = useState<ArrayBuffer[]>()

  const [fileStatus, setFileStatus] = useAtom(fileStatusAtom)

  const [loadingSpring, setLoadingSpring] = useSpring(() => ({
    from: {
      opacity: 1,
    },
    to: {
      opacity: 1,
    },
  }))

  useEffect(() => {
    if (fileStatus) {
      const finishedFiles = fileStatus.filter((f) => f.status === 'success')

      if (finishedFiles.length === fileStatus.length) {
        setTimeout(() => {
          setLoadingSpring.start({
            opacity: 0,
          })
        }, 1000)
      } else {
        setLoadingSpring.start({
          opacity: 1,
        })
      }
    }
  }, [fileStatus])

  const onDrop = useCallback((acceptedFiles: File[]) => {
    acceptedFiles.forEach((file) => {
      const reader = new FileReader()

      reader.onabort = () => console.log('file reading was aborted')
      reader.onerror = () => console.log('file reading has failed')
      reader.onload = () => {
        const binaryStr = reader.result as ArrayBuffer
        setFileData((prev) => {
          if (prev) {
            return [...prev, binaryStr]
          } else {
            return [binaryStr]
          }
        })
      }
      reader.readAsArrayBuffer(file)
    })
  }, [])

  const {
    getRootProps,
    getInputProps,
    isFocused,
    isDragAccept,
    isDragReject,
    acceptedFiles,
  } = useDropzone({ onDrop })

  const style = useMemo(
    () => ({
      ...baseStyle,
      ...(isFocused ? focusedStyle : {}),
      ...(isDragAccept ? acceptStyle : {}),
      ...(isDragReject ? rejectStyle : {}),
    }),
    [isFocused, isDragAccept, isDragReject],
  )

  const fileEndpoint = useMutation({
    mutationFn: async (file: File[]) => {
      // boundary is file.name
      file.forEach(async (f, i) => {
        const formData = new FormData()
        if (fileData) {
          formData.append(
            'file',
            new Blob([fileData[i]], { type: f.type }),
            f.name,
          )
        }
        // const res = await fetch('http://localhost:8080/upl', {
        //   method: 'POST',
        //   mode: 'no-cors',
        //   body: formData,
        // })

        const res = await axios.post('http://localhost:3030/upload', formData, {
          onUploadProgress: (progressEvent) => {
            setFileStatus((prev) => {
              if (prev) {
                prev[i].progress =
                  (progressEvent.loaded / (progressEvent.total || 1)) * 100
                return [...prev]
              } else {
                return []
              }
            })
          },
        })

        if (res.status === 200) {
          setFileStatus((prev) => {
            if (prev) {
              prev[i].status = 'success'
              return [...prev]
            } else {
              return []
            }
          })
        }
      })
    },
  })

  const acceptedFileItems = acceptedFiles.map((file: FileWithPath) => (
    <li key={file.path}>
      {file.path} - {file.size} bytes
    </li>
  ))

  return (
    <>
      {!isUploading && (
        <div className='container'>
          <div {...getRootProps({ style })}>
            <input {...getInputProps()} />
            <p>Drag 'n' drop some files here, or click to select files</p>
            <ul>{acceptedFileItems}</ul>
          </div>
        </div>
      )}
      {isUploading && (
        <div>
          {fileEndpoint.error &&
          typeof fileEndpoint.error === 'object' &&
          'message' in fileEndpoint.error &&
          typeof fileEndpoint.error.message === 'string' ? (
            <div className='text-red-500'>
              <p>Oops, something went wrong :(</p>
              <p>{fileEndpoint.error.message}</p>
              <Link className='text-blue-500 hover:text-blue-700' to='/'>
                Return home
              </Link>
            </div>
          ) : (
            <>
              <a.div style={loadingSpring}>
                <ClimbingBoxLoader
                  color={
                    window.matchMedia('(prefers-color-scheme: dark)').matches
                      ? '#ffffff'
                      : '#000000'
                  }
                  loading={isUploading}
                  size={20}
                />
              </a.div>
            </>
          )}
        </div>
      )}
      {acceptedFiles.length > 0 && !isUploading && (
        <button
          className='p-2 rounded-lg top-4 relative bg-zinc-200 dark:bg-zinc-800 hover:bg-zinc-300 dark:hover:bg-zinc-700 transition-colors'
          onClick={() => {
            setUploading(true)

            setFileStatus(
              acceptedFiles.map((file) => ({
                name: file.name,
                status: 'Uploading',
                length: file.size,
                progress: 0,
                url: '',
              })),
            )

            fileEndpoint.mutate(acceptedFiles)
          }}>
          Upload
        </button>
      )}
    </>
  )
}

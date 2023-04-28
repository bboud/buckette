import { useMutation } from '@tanstack/react-query'
import { useCallback, useMemo, useState } from 'react'
import { useDropzone } from 'react-dropzone'
import { Link, useNavigate } from 'react-router-dom'

import ClimbingBoxLoader from 'react-spinners/ClimbingBoxLoader'

const baseStyle = {
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

  const onDrop = useCallback((acceptedFiles: File[]) => {
    acceptedFiles.forEach((file) => {
      const reader = new FileReader()

      reader.onload = () => {
        const binaryStr = reader.result
        console.log('ya')
        console.log(binaryStr)
      }
    })

    console.log(acceptedFiles)
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
    mutationFn: async (file: File) => {
      const formData = new FormData()

      formData.append('file', acceptedFiles[0])

      const res = await fetch('http://localhost:8080/upl', {
        method: 'POST',
        headers: {
          'Content-Type': 'multipart/form-data',
        },
        body: formData,
      })

      if (!res.ok) {
        throw new Error(res.statusText)
      }

      return res.json()
    },
  })

  const acceptedFileItems = acceptedFiles.map((file) => (
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
          {fileEndpoint.error || fileEndpoint.status === 'error' ? (
            <div className='text-red-500'>
              <p>Oops, something went wrong :(</p>
              <p>{fileEndpoint.error.message}</p>
              <Link className='text-blue-500 hover:text-blue-700' to='/'>
                Return home
              </Link>
            </div>
          ) : (
            <ClimbingBoxLoader color={'#ffffff'} loading={true} size={20} />
          )}
        </div>
      )}
      {acceptedFiles.length > 0 && !isUploading && (
        <button
          className='p-2 rounded-lg top-4 relative bg-zinc-200 dark:bg-zinc-800 hover:bg-zinc-300 dark:hover:bg-zinc-700 transition-colors'
          onClick={() => {
            setUploading(true)
            fileEndpoint.mutate(acceptedFiles[0])
          }}>
          Upload
        </button>
      )}
    </>
  )
}

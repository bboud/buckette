import { useMutation } from '@tanstack/react-query'
import { useCallback, useEffect, useMemo, useState } from 'react'
import { FileWithPath, useDropzone, DropzoneRootProps } from 'react-dropzone'
import { Link, useNavigate } from 'react-router-dom'

import ClimbingBoxLoader from 'react-spinners/ClimbingBoxLoader'

import axios from 'axios'
import {
  fileStatusAtom,
  filesAtom,
  inputPropsAtom,
  isDragAcceptAtom,
  isDragRejectAtom,
  isFocusedAtom,
  isUploadingAtom,
} from '../utils/atoms'

import { animated as a, useSpring } from '@react-spring/web'

import { useAtom, useAtomValue } from 'jotai'

import { Icon } from '@iconify/react'

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
  backgroundColor: '#f0f0f0',
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
  const [isUploading, setUploading] = useAtom(isUploadingAtom)

  const [files, setFiles] = useAtom(filesAtom)

  const [fileStatus, setFileStatus] = useAtom(fileStatusAtom)

  const isFocused = useAtomValue(isFocusedAtom)
  const isDragAccept = useAtomValue(isDragAcceptAtom)
  const isDragReject = useAtomValue(isDragRejectAtom)

  const acceptedFiles = useAtomValue(filesAtom)

  const [loadingSpring, setLoadingSpring] = useSpring(() => ({
    from: {
      opacity: 1,
    },
    to: {
      opacity: 1,
    },
  }))

  const [resetSpring, setResetSpring] = useSpring(() => ({
    opacity: 0,
    display: 'none',
  }))

  useEffect(() => {
    if (fileStatus) {
      const finishedFiles = fileStatus.filter(
        (f) => f.status === 'success' || f.status === 'error',
      )

      if (finishedFiles.length < 1 || fileStatus.length < 1) {
        return
      }

      if (finishedFiles.length === fileStatus.length) {
        console.log('yup')
        setTimeout(() => {
          setLoadingSpring.start({
            opacity: 0,
            onRest: () => {
              console.log('hi2')
              setResetSpring.start({
                opacity: 1,
                display: 'block',
              })
            },
          })
        }, 1000)
      } else {
        setLoadingSpring.start({
          opacity: 1,
        })
      }
    }
  }, [fileStatus])

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
        if (files) {
          formData.append('file', f, f.name)
          // const res = await fetch('http://localhost:8080/upl', {
          //   method: 'POST',
          //   mode: 'no-cors',
          //   body: formData,
          // })

          // const digest = await window.crypto.subtle.digest(
          //   'SHA-256',
          //   await blob.arrayBuffer(),
          // )

          // const hashHex = Array.from(new Uint8Array(digest))
          //   .map((b) => b.toString(16).padStart(2, '0'))
          //   .join('')

          // const hashBase64 = Buffer.from(hashHex).toString('base64')

          try {
            const res = await axios.post<FileResponse>(
              'http://localhost:8080/upl',
              formData,
              {
                onUploadProgress: (progressEvent) => {
                  setFileStatus((prev) => {
                    if (prev) {
                      prev[i].progress =
                        (progressEvent.loaded / (progressEvent.total || 1)) *
                        100
                      return [...prev]
                    } else {
                      return []
                    }
                  })
                },

                headers: {
                  'File-Size': f.size,
                  'File-Name': f.name,
                  'File-Type':
                    f.type === '' ? 'application/octet-stream' : f.type,
                },

                maxRedirects: 0,
              },
            )

            const data = res.data

            if (res.status === 200) {
              console.log('RESPONSE:')
              console.log(data)

              setFileStatus((prev) => {
                if (prev) {
                  prev[i].status = 'success'
                  prev[i].url = data.URL
                  data.Duplicate && (prev[i].error = 'File is a duplicate!')
                  return [...prev]
                } else {
                  return []
                }
              })
            } else {
              setFileStatus((prev) => {
                if (prev) {
                  prev[i].status = 'error'
                  prev[i].error = res.statusText
                  return [...prev]
                } else {
                  return []
                }
              })
            }
          } catch (e) {
            setFileStatus((prev) => {
              if (prev) {
                prev[i].status = 'error'
                prev[i].error = `an error has occurred: ${e}`
                return [...prev]
              } else {
                return []
              }
            })
          }
        }
      })
    },
  })

  return (
    <>
      {isUploading && (
        <div>
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
        </div>
      )}
      <a.div className='absolute top-2 w-40' style={resetSpring}>
        <button
          className='p-2 rounded-lg top-4 relative bg-zinc-200 dark:bg-zinc-800 hover:bg-zinc-300 dark:hover:bg-zinc-700 transition-colors'
          onClick={() => {
            setFiles([])
            setFileStatus([])
            setResetSpring.start({
              opacity: 0,
              onRest: (_, ctrl) => {
                ctrl.stop()
                setResetSpring.start({
                  display: 'none',
                })

                setUploading(false)
              },
            })
          }}>
          Upload more
        </button>
      </a.div>
      {fileStatus.length > 0 && !isUploading && (
        <button
          className='p-2 rounded-lg top-4 relative bg-zinc-200 dark:bg-zinc-800 hover:bg-zinc-300 dark:hover:bg-zinc-700 transition-colors'
          onClick={() => {
            setUploading(true)

            setFileStatus((prev) => {
              if (prev) {
                prev.forEach((f) => {
                  f.status = 'uploading'
                })
                return [...prev]
              } else {
                return []
              }
            })

            fileEndpoint.mutate(acceptedFiles)
          }}>
          Upload
        </button>
      )}
    </>
  )
}

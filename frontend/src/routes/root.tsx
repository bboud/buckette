import { Link, Outlet, useLocation, useNavigate } from 'react-router-dom'

import CountUp from 'react-countup'
import { SignInButton, useClerk, useUser } from '@clerk/clerk-react'
import {
  fileStatusAtom,
  filesAtom,
  inputPropsAtom,
  isDragAcceptAtom,
  isDragRejectAtom,
  isFocusedAtom,
  isUploadingAtom,
} from '../utils/atoms'

import {
  animated as a,
  useSpring,
  useSprings,
  useTransition,
} from '@react-spring/web'

import { useAtom } from 'jotai'
import { useCallback, useEffect } from 'react'
import { Icon } from '@iconify/react'
import { FileWithPath, useDropzone } from 'react-dropzone'

export default function Root() {
  const { isSignedIn } = useUser()
  const { signOut } = useClerk()
  const location = useLocation()
  const navigate = useNavigate()

  const [fileStatus, setFileStatus] = useAtom(fileStatusAtom)

  const [isUploading, setIsUploading] = useAtom(isUploadingAtom)

  const [files, setFiles] = useAtom(filesAtom)

  const [focused, setFocused] = useAtom(isFocusedAtom)
  const [dragAccept, setDragAccept] = useAtom(isDragAcceptAtom)
  const [dragReject, setDragReject] = useAtom(isDragRejectAtom)
  const [inputProps, setInputProps] = useAtom(inputPropsAtom)

  const onDrop = useCallback(
    (acceptedFiles: File[]) => {
      if (isUploading) return

      setIconSpring.start({
        scale: dragAccept ? 0.2 : 1,
        config: {
          tension: 200,
          friction: 10,
        },
      })

      console.log(files)

      const filtered = acceptedFiles.filter(
        (file) => !files.some((f) => f.name === file.name),
      )
      console.log(filtered)

      setFileStatus((prev) => {
        return [
          ...prev,
          ...filtered.map((file) => ({
            name: file.name,
            type: file.type === '' ? 'application/octet-stream' : file.type,
            status: 'idle',
            length: file.size,
            progress: 0,
            url: '',
          })),
        ]
      })

      setFiles((prev) => [...prev, ...filtered])
      navigate('/upload')
    },
    [setFiles, files, isUploading, navigate, setFileStatus, dragAccept],
  )

  const {
    getRootProps,
    getInputProps,
    isFocused,
    isDragAccept,
    isDragReject,
    acceptedFiles,
  } = useDropzone({
    onDrop,
  })

  const [fileSprings, setFileSprings] = useSprings(fileStatus.length, () => ({
    width: '0%',
  }))

  const [successSprings, setSuccessSprings] = useSprings(
    fileStatus.length,
    () => ({
      opacity: 0,
      display: 'none',
    }),
  )

  const [progressSprings, setProgressSprings] = useSprings(
    fileStatus.length,
    (i) => ({
      opacity: 1,
      display: 'block',
    }),
  )

  const [uploadWindowSpring, setUploadWindowSpring] = useSpring(() => ({
    opacity: 0,
    display: 'none',
  }))

  const [iconSpring, setIconSpring] = useSpring(() => ({
    scale: 1,
  }))

  useEffect(() => {
    setUploadWindowSpring.start({
      opacity: isDragAccept || isDragReject ? 1 : 0,
      display: isDragAccept || isDragReject ? 'block' : 'none',
    })

    setIconSpring.start({
      scale: isDragAccept || isDragReject ? 1.2 : 1,
      delay: 100,
    })

    setDragAccept(isDragAccept)
    setDragReject(isDragReject)
    setFocused(isFocused)
  }, [isDragAccept, isDragReject, isFocused, getInputProps])

  useEffect(() => {
    if (fileStatus) {
      setFileSprings.start((i) => ({
        width: `${fileStatus[i].progress}%`,
      }))

      fileStatus.forEach((f, i) => {
        if (f.status === 'success') {
          setSuccessSprings.start((i) => ({
            opacity: 1,
            display: 'inline',
          }))

          setTimeout(() => {
            setProgressSprings.start((i) => ({
              opacity: fileStatus[i].status === 'success' ? 0 : 1,
              onRest: (i, ctrl) => {
                ctrl.start({
                  display: 'none',
                })
              },
            }))
          }, 3000)
        }
      })
    }
  }, [fileStatus])

  useEffect(() => {
    if (location.pathname.includes('welcome')) {
      document.title = 'Welcome to Buckette'
    } else if (location.pathname.includes('upload')) {
      document.title = 'Upload to Buckette'
    } else {
      document.title = 'Buckette'
    }
  }, [location.pathname])

  return (
    <>
      <a.div
        style={uploadWindowSpring}
        className={`w-full h-full absolute z-30 ${
          isDragAccept
            ? 'bg-emerald-500/20'
            : isDragReject
            ? 'bg-red-500/20'
            : 'bg-gray-400/20'
        }`}
        {...getRootProps()}>
        <input {...getInputProps()} />
      </a.div>
      <div
        onDragEnter={() => {
          setUploadWindowSpring.start({
            opacity: 1,
            display: 'block',
          })
        }}
        className='absolute z-10 left-0 right-0 w-32 h-24 m-auto my-3 grid space-y-2 justify-center'>
        {location.pathname.includes('upload') && (
          <div className='bg-zinc-300/20 p-5 rounded-md'>
            <input
              className='absolute opacity-0 bg-zinc-300 w-24 h-24 left-0 right-1 m-auto top-1'
              onInput={(e) => {
                const fileList = e.currentTarget.files
                if (fileList) {
                  const inputFiles = [...fileList]

                  // filter any files that are already in the list
                  const filtered = inputFiles.filter(
                    (file) => !files.some((f) => f.name === file.name),
                  )

                  setFileStatus((prev) => {
                    return [
                      ...prev,
                      ...filtered.map((file) => ({
                        name: file.name,
                        type:
                          file.type === ''
                            ? 'application/octet-stream'
                            : file.type,
                        status: 'idle',
                        length: file.size,
                        progress: 0,
                        url: '',
                      })),
                    ]
                  })

                  setFiles((prev) => [...prev, ...filtered])
                }
              }}
              disabled={isUploading}
              type='file'
              multiple={true}
            />
            <a.div style={iconSpring}>
              <Icon
                className={`w-full transition-colors ${
                  isDragAccept
                    ? 'text-green-500'
                    : isDragReject
                    ? 'text-red-500'
                    : isFocused
                    ? 'text-blue-500'
                    : 'text-gray-200'
                }`}
                icon='ic:baseline-download-for-offline'
                width={64}
                height={64}
              />
            </a.div>
          </div>
        )}
      </div>
      <div
        onDragEnter={() => {
          setUploadWindowSpring.start({
            opacity: 1,
            display: 'block',
          })
        }}
        className='absolute w-full z-20 grid space-y-2 justify-center h-auto top-32'>
        {location.pathname.includes('upload') &&
          fileStatus.map((f, i) => (
            <div
              className='bg-zinc-200/75 rounded-lg dark:bg-zinc-800/75 dark:text-zinc-300 text-zinc-800 p-2 w-96'
              key={i}>
              <div>
                <span>
                  <Icon
                    className={`transition-colors inline mr-2 ${
                      f.status === 'uploading'
                        ? 'animate-spin text-zinc-400'
                        : f.status === 'success'
                        ? f.error
                          ? 'text-yellow-300'
                          : 'text-green-500'
                        : f.status === 'error'
                        ? 'text-red-500'
                        : ''
                    }`}
                    inline={true}
                    icon={
                      f.status == 'uploading'
                        ? 'gg:spinner'
                        : f.status === 'success'
                        ? f.error
                          ? 'material-symbols:warning'
                          : 'akar-icons:check'
                        : f.status === 'error' || f.error
                        ? 'akar-icons:close'
                        : 'ph:dot-outline-light'
                    }
                  />
                </span>
                <b>{f.name}</b>
              </div>
              <a.div style={progressSprings[i]}>
                <CountUp
                  end={f.progress}
                  decimals={2}
                  suffix={'%'}
                  preserveValue={true}
                />
                <a.div
                  className='dark:bg-zinc-200 rounded-lg bg-zinc-800 h-2'
                  style={fileSprings[i]}></a.div>
              </a.div>
              {f.url && (
                <a
                  className='text-blue-500 hover:underline transition-all'
                  href={f.url}>
                  {f.url}
                </a>
              )}
              {f.error && <p className='text-red-500'>{f.error}</p>}
            </div>
          ))}
      </div>
      <div
        onDragEnter={() => {
          setUploadWindowSpring.start({
            opacity: 1,
            display: 'block',
          })
        }}
        className={`w-full h-full ${
          location.pathname !== '/' ? 'fixed' : 'absolute'
        } text-zinc-900 dark:text-zinc-300 bg-zinc-100 dark:bg-zinc-900`}>
        <div className='m-4 grid justify-start grid-flow-col space-x-4'>
          {isSignedIn && (
            <>
              <Link
                to={'upload'}
                className='p-2 rounded-lg bg-zinc-200 dark:bg-zinc-800 hover:bg-zinc-300 dark:hover:bg-zinc-700 transition-colors'>
                Upload
              </Link>
              <button
                className='p-2 rounded-lg bg-zinc-200 dark:bg-zinc-800 hover:bg-zinc-300 dark:hover:bg-zinc-700 transition-colors'
                onClick={() => {
                  signOut()
                }}>
                Sign out
              </button>
            </>
          )}
          {!isSignedIn && <SignInButton mode='modal' />}
        </div>
        <table className='table-fixed text-center w-[95%] left-0 right-0 m-auto'>
          <tbody>
            <tr>
              <th className='hover:bg-zinc-200 hover:dark:bg-zinc-800 rounded-sm transition-colors'>
                Name
              </th>
              <th className='hover:bg-zinc-200 hover:dark:bg-zinc-800 rounded-sm transition-colors'>
                Uploaded at
              </th>
              <th className='hover:bg-zinc-200 hover:dark:bg-zinc-800 rounded-sm transition-colors'>
                Size
              </th>
            </tr>
            {location.pathname.includes('/#/upload') &&
              location.pathname !== '/' && (
                <tr className='hover:bg-zinc-200 transition-colors rounded-lg hover:dark:bg-zinc-800'>
                  <td>
                    <Link
                      onClick={(e) => {
                        if (location.pathname.split('/').length === 2) {
                          e.preventDefault()
                          navigate('/')
                        }
                      }}
                      to={location.pathname.split('/').slice(0, -1).join('/')}>
                      ..
                    </Link>
                  </td>
                  <td></td>
                  <td></td>
                </tr>
              )}
            {/* {test_data.map((file) => (
              <tr
                key={file.FileName}
                className='hover:bg-zinc-200 transition-colors rounded-lg hover:dark:bg-zinc-800'>
                <td>{file.FileName}</td>
                <td>{file.Uploaded}</td>
                <td>
                  <CountUp end={file.Size} decimals={2} /> MB
                </td>
              </tr>
            ))} */}
          </tbody>
        </table>
      </div>
      {location.pathname !== '/' && (
        <div className='w-full h-full absolute text-zinc-900 dark:text-zinc-300'>
          <Outlet />
        </div>
      )}
    </>
  )
}

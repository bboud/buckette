import { Link, Outlet, useLocation, useNavigate } from 'react-router-dom'

import CountUp from 'react-countup'
import { SignInButton, useClerk, useUser } from '@clerk/clerk-react'
import { fileStatusAtom } from '../utils/atoms'

import { animated as a, useSprings } from '@react-spring/web'

import { useAtom } from 'jotai'
import { useEffect } from 'react'

export default function Root() {
  const { isSignedIn } = useUser()
  const { signOut } = useClerk()
  const location = useLocation()
  const navigate = useNavigate()

  console.log(location.pathname)

  const [fileStatus, setFileStatus] = useAtom(fileStatusAtom)

  const [fileSprings, setFileSprings] = useSprings(fileStatus.length, () => ({
    width: '0%',
  }))

  const [progressSprings, setProgressSprings] = useSprings(
    fileStatus.length,
    (i) => ({
      opacity: 1,
      display: 'block',
    }),
  )

  useEffect(() => {
    if (fileStatus) {
      setFileSprings.start((i) => ({
        width: `${fileStatus[i].progress}%`,
      }))

      fileStatus.forEach((f, i) => {
        if (f.status === 'success') {
          setTimeout(() => {
            setProgressSprings.start((i, a) => ({
              opacity: fileStatus[i].status === 'success' ? 0 : 1,
              onRest: (_, ctrl) => {
                ctrl.start({
                  display:
                    fileStatus[i].status === 'success' ? 'none' : 'block',
                })
              },
            }))
          }, 3000)
        }
      })
    }
  }, [fileStatus])

  return (
    <>
      <div className='absolute z-50 w-full top-5 grid space-y-2 justify-center'>
        {fileStatus &&
          fileStatus.map((f, i) => (
            <div
              className='bg-zinc-200/75 rounded-lg dark:bg-zinc-800/75 p-2 w-96'
              key={i}>
              <p>
                <b>{f.name}</b>
              </p>
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
              {!f.error && f.status === 'success' && <p>Success! </p>}
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
        className={`w-full h-full ${
          location.pathname !== '/' ? 'fixed' : 'absolute'
        } text-zinc-900 dark:text-zinc-300 bg-zinc-100 dark:bg-zinc-900`}>
        <Outlet />
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
    </>
  )
}

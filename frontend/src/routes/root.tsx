import { Link, Outlet, useLocation, useNavigate } from 'react-router-dom'

import CountUp from 'react-countup'
import { useClerk, useUser } from '@clerk/clerk-react'

const test_data: FileDB[] = [
  {
    FileName: 'test',
    UUID: 'blah',
    Size: 100,
    Uploaded: new Date().toISOString(),
    UserUploaded: '',
    DownloadCount: 0,
  },
  {
    FileName: 'test2',
    UUID: 'blah2',
    Size: 100,
    Uploaded: new Date().toISOString(),
    UserUploaded: '',
    DownloadCount: 0,
  },
]

export default function Root() {
  const { isSignedIn } = useUser()
  const { signOut } = useClerk()
  const location = useLocation()
  const navigate = useNavigate()
  return (
    <>
      <div className='w-full h-full absolute text-zinc-900 dark:text-zinc-300 bg-zinc-100 dark:bg-zinc-900'>
        <Outlet />
        <div className='m-4 grid justify-start grid-flow-col space-x-4'>
          {isSignedIn && (
            <button
              className='p-2 rounded-lg bg-zinc-200 dark:bg-zinc-800 hover:bg-zinc-300 dark:hover:bg-zinc-700 transition-colors'
              onClick={() => {
                signOut()
              }}>
              Sign out
            </button>
          )}
          <Link
            to={'upload'}
            className='p-2 rounded-lg bg-zinc-200 dark:bg-zinc-800 hover:bg-zinc-300 dark:hover:bg-zinc-700 transition-colors'>
            {isSignedIn ? 'Upload' : 'Sign in'}
          </Link>
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
            {location.pathname !== '/' && (
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
            {test_data.map((file) => (
              <tr
                key={file.FileName}
                className='hover:bg-zinc-200 transition-colors rounded-lg hover:dark:bg-zinc-800'>
                <td>{file.FileName}</td>
                <td>{file.Uploaded}</td>
                <td>
                  <CountUp end={file.Size} decimals={2} /> MB
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </>
  )
}

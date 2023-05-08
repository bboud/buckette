import { SignIn, useUser } from '@clerk/clerk-react'
import { animated as a, useSpring } from '@react-spring/web'
import { useNavigate } from 'react-router-dom'

import UploadComponent from '../components/UploadComponent'

import Dropzone from 'react-dropzone'
import { useMutation } from '@tanstack/react-query'

function Upload() {
  const [uploadSpring, setUploadSpring] = useSpring(() => ({
    from: {
      opacity: 0,
    },
    to: {
      opacity: 1,
    },
  }))

  const { isSignedIn } = useUser()

  const navigate = useNavigate()

  return (
    <a.div
      style={uploadSpring}
      className='absolute bg-zinc-900/20 backdrop-blur-sm w-full h-full grid justify-center'>
      {isSignedIn && (
        <div className='left-28 relative'>
          <UploadComponent />
        </div>
      )}
    </a.div>
  )
}

export default Upload

import { SignIn, useUser } from '@clerk/clerk-react'
import { animated as a, useSpring } from '@react-spring/web'
import { useNavigate } from 'react-router-dom'

import StyledDropzone from '../components/StyledDropzone'

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
      className='absolute z-10 bg-zinc-900/20 backdrop-blur-sm w-full h-full'>
      {!isSignedIn && (
        <div className='w-full h-full flex flex-col justify-center items-center'>
          <SignIn afterSignInUrl={'/upload'} />
        </div>
      )}
      {isSignedIn && (
        <div className='w-full h-full flex flex-col justify-center items-center'>
          <StyledDropzone />
        </div>
      )}
    </a.div>
  )
}

export default Upload

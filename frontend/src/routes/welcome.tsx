import { Icon } from '@iconify/react'
import { useSpring, animated as a } from '@react-spring/web'
import { useNavigate } from 'react-router-dom'

function Welcome() {
  const [welcomeSpring, setWelcomeSpring] = useSpring(() => ({
    opacity: 1,
    display: 'block',
  }))

  const navigate = useNavigate()

  return (
    <>
      <a.div
        style={welcomeSpring}
        className='w-full h-full bg-zinc-300/10 dark:bg-zinc-800/20 backdrop-blur-sm absolute grid justify-center'>
        <div className='bg-zinc-300/50 dark:bg-zinc-800/50 relative h-36 rounded-lg top-1/4 text-center p-4 ml-4 mr-4'>
          <p className='text-3xl'>Welcome to Buckette</p>
          <div className='top-12 relative'>
            <button
              onClick={() =>
                setWelcomeSpring.start({
                  opacity: 0,
                  onRest: (_, ctrl) => {
                    ctrl.start({ display: 'none' })
                    navigate('/')
                  },
                })
              }
              className='text-sm'>
              <Icon
                icon='material-symbols:home'
                className='inline text-lg'
                inline={true}
              />
              Click to enter the main menu
            </button>
          </div>
        </div>
      </a.div>
    </>
  )
}

export default Welcome

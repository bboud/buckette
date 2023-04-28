import { useClerk, useUser } from '@clerk/clerk-react'
import { useQuery } from '@tanstack/react-query'
import { useState } from 'react'

function TestRoute() {
  const user = useUser()
  const clerk = useClerk()

  const [sessionToken, setSessionToken] = useState<string | null>(null)

  clerk.session?.getToken().then((token) => {
    if (token) {
      setSessionToken(token)
    }
  })

  const { data, isLoading, error } = useQuery({
    queryKey: ['test'],
    queryFn: async () => {
      // Fetch with clerk auth
      const res = await fetch('http://localhost:3000', {
        headers: {
          Authorization: `Bearer ${sessionToken}`,
        },
      })

      return res.json()
    },
  })

  if (isLoading) {
    return <div>Loading...</div>
  }

  return <>hello, {data.hello}</>
}

export default TestRoute

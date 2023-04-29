import React from 'react'
import ReactDOM from 'react-dom/client'

import { dark } from '@clerk/themes'

import './index.css'

import { createHashRouter, RouterProvider } from 'react-router-dom'
import Root from './routes/root'
import Upload from './routes/upload'
import { ClerkProvider } from '@clerk/clerk-react'
import TestRoute from './routes/test'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'

const clerkPubKey = 'pk_test_ZXhwZXJ0LXdvcm0tMS5jbGVyay5hY2NvdW50cy5kZXYk'

console.log(clerkPubKey)

const queryClient = new QueryClient()

const router = createHashRouter([
  {
    path: '/',
    element: <Root />,
    children: [
      {
        path: 'upload',
        element: <Upload />,
      },
      {
        // directory listing
        path: ':directory/*',
        element: <Root />,
      },
    ],
  },
  {
    path: 'test',
    element: <TestRoute />,
  },
])

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <ClerkProvider
      appearance={{
        baseTheme: window.matchMedia('(prefers-color-scheme: dark)').matches
          ? dark
          : undefined,
      }}
      publishableKey={clerkPubKey}>
      <QueryClientProvider client={queryClient}>
        <RouterProvider router={router} />
      </QueryClientProvider>
    </ClerkProvider>
  </React.StrictMode>,
)

import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import axios from 'axios'
import './index.css'
import App from './App.tsx'

axios.interceptors.response.use(
  (response) => response,
  (error) => {
    const status = (error as { response?: { status?: number } })?.response?.status

    if (status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('username')

      if (window.location.pathname !== '/login') {
        window.location.replace('/login')
      }
    }

    return Promise.reject(error)
  },
)

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <BrowserRouter>
      <App />
    </BrowserRouter>
  </StrictMode>,
)

import type { Metadata } from 'next'
import './globals.css'
import { LangProvider } from '@/context/LangContext'

export const metadata: Metadata = {
  title: 'Masaar CRM',
  description: 'The open-source CRM built for the UAE market',
}

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en" dir="ltr">
      <head>
        <link rel="preconnect" href="https://fonts.googleapis.com" />
        <link rel="preconnect" href="https://fonts.gstatic.com" crossOrigin="anonymous" />
        <link
          href="https://fonts.googleapis.com/css2?family=Cairo:wght@400;500;600;700&family=Inter:wght@400;500;600;700&display=swap"
          rel="stylesheet"
        />
      </head>
      <body className="bg-gray-50 text-gray-900 antialiased">
        <LangProvider>
          {children}
        </LangProvider>
      </body>
    </html>
  )
}

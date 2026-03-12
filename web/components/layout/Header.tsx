'use client'
import { useRouter } from 'next/navigation'
import { useAuthStore } from '@/store/auth'
import { useLang } from '@/context/LangContext'
import { NotificationBell } from './NotificationBell'
import { api } from '@/lib/api'
import { getRefreshToken } from '@/lib/auth'

export function Header({ title }: { title: string }) {
  const { user, logout } = useAuthStore()
  const { lang, setLang, t } = useLang()
  const router = useRouter()

  const handleLogout = async () => {
    const rt = getRefreshToken()
    if (rt) await api.auth.logout(rt).catch(() => {})
    logout()
    router.push('/login')
  }

  return (
    <header className="h-14 flex items-center justify-between px-6 bg-white border-b border-gray-200 shrink-0">
      <h1 className="font-semibold text-gray-800 text-base">{title}</h1>

      <div className="flex items-center gap-3">
        {/* RTL/LTR toggle */}
        <button
          onClick={() => setLang(lang === 'ar' ? 'en' : 'ar')}
          className="text-xs px-2.5 py-1 rounded-md border border-gray-200 text-gray-600 hover:bg-gray-50 transition-colors font-medium"
        >
          {lang === 'ar' ? 'EN' : 'عربي'}
        </button>

        <NotificationBell />

        {/* Avatar + logout */}
        <div className="flex items-center gap-2">
          <div className="w-8 h-8 rounded-full bg-brand-600 text-white text-sm font-semibold flex items-center justify-center">
            {user?.name?.[0]?.toUpperCase() ?? 'U'}
          </div>
          <button
            onClick={handleLogout}
            className="text-xs text-gray-500 hover:text-red-500 transition-colors"
          >
            {t('خروج', 'Sign out')}
          </button>
        </div>
      </div>
    </header>
  )
}

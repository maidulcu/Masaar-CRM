'use client'
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { useLang } from '@/context/LangContext'
import clsx from 'clsx'

const nav = [
  {
    href: '/pipeline',
    icon: '⬛',
    label: { en: 'Pipeline', ar: 'خط الأنابيب' },
  },
  {
    href: '/inbox',
    icon: '💬',
    label: { en: 'Inbox', ar: 'الرسائل' },
  },
  {
    href: '/contacts',
    icon: '👥',
    label: { en: 'Contacts', ar: 'جهات الاتصال' },
  },
  {
    href: '/deals',
    icon: '💼',
    label: { en: 'Deals', ar: 'الصفقات' },
  },
]

export function Sidebar() {
  const pathname = usePathname()
  const { lang, t } = useLang()

  return (
    <aside className="flex flex-col w-56 min-h-screen bg-gray-900 text-white shrink-0">
      {/* Logo */}
      <div className="flex items-center gap-2 px-5 py-5 border-b border-gray-700">
        <span className="text-brand-500 text-2xl font-bold">M</span>
        <span className="font-semibold text-lg tracking-tight">
          {lang === 'ar' ? 'مسار' : 'Masaar'}
        </span>
      </div>

      {/* Nav */}
      <nav className="flex-1 px-3 py-4 space-y-1">
        {nav.map((item) => {
          const active = pathname.startsWith(item.href)
          return (
            <Link
              key={item.href}
              href={item.href}
              className={clsx(
                'flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors',
                active
                  ? 'bg-brand-600 text-white'
                  : 'text-gray-400 hover:bg-gray-800 hover:text-white'
              )}
            >
              <span className="text-base">{item.icon}</span>
              <span>{lang === 'ar' ? item.label.ar : item.label.en}</span>
            </Link>
          )
        })}
      </nav>

      {/* Footer */}
      <div className="px-4 py-3 border-t border-gray-700 text-xs text-gray-500 text-center">
        {t('مسار CRM', 'Masaar CRM')} · MIT
      </div>
    </aside>
  )
}

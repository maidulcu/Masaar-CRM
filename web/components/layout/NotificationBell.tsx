'use client'
import { useState } from 'react'
import { useNotifications } from '@/hooks/useNotifications'
import { useAuthStore } from '@/store/auth'
import { useLang } from '@/context/LangContext'
import clsx from 'clsx'

export function NotificationBell() {
  const user = useAuthStore((s) => s.user)
  const { notifications, unread, markRead } = useNotifications(user?.id ?? null)
  const [open, setOpen] = useState(false)
  const { t } = useLang()

  return (
    <div className="relative">
      <button
        onClick={() => setOpen((v) => !v)}
        className="relative p-2 rounded-full hover:bg-gray-100 transition-colors"
        aria-label="Notifications"
      >
        <svg className="w-5 h-5 text-gray-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
            d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6 6 0 10-12 0v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
        </svg>
        {unread > 0 && (
          <span className="absolute top-1 end-1 w-4 h-4 bg-red-500 text-white text-[10px] font-bold rounded-full flex items-center justify-center">
            {unread > 9 ? '9+' : unread}
          </span>
        )}
      </button>

      {open && (
        <>
          <div className="fixed inset-0 z-10" onClick={() => setOpen(false)} />
          <div className="absolute end-0 top-10 z-20 w-80 bg-white rounded-xl shadow-xl border border-gray-100 overflow-hidden">
            <div className="px-4 py-3 border-b border-gray-100 font-semibold text-sm text-gray-700">
              {t('الإشعارات', 'Notifications')}
            </div>
            <div className="max-h-80 overflow-y-auto divide-y divide-gray-50">
              {notifications.length === 0 ? (
                <p className="px-4 py-6 text-sm text-gray-400 text-center">
                  {t('لا توجد إشعارات', 'No notifications')}
                </p>
              ) : (
                notifications.slice(0, 15).map((n) => (
                  <button
                    key={n.id}
                    onClick={() => markRead(n.id)}
                    className={clsx(
                      'w-full text-start px-4 py-3 text-sm hover:bg-gray-50 transition-colors',
                      !n.read && 'bg-blue-50'
                    )}
                  >
                    <p className="font-medium text-gray-800">{n.title}</p>
                    {n.body && <p className="text-gray-500 text-xs mt-0.5">{n.body}</p>}
                  </button>
                ))
              )}
            </div>
          </div>
        </>
      )}
    </div>
  )
}

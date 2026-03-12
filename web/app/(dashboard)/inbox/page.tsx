'use client'
import { useEffect, useState } from 'react'
import Link from 'next/link'
import { Header } from '@/components/layout/Header'
import { api } from '@/lib/api'
import { useLang } from '@/context/LangContext'
import type { WhatsAppThread } from '@/types'
import clsx from 'clsx'

const statusTabs = [
  { key: '', label: { en: 'All', ar: 'الكل' } },
  { key: 'open', label: { en: 'Open', ar: 'مفتوح' } },
  { key: 'pending', label: { en: 'Pending', ar: 'معلق' } },
  { key: 'closed', label: { en: 'Closed', ar: 'مغلق' } },
]

const statusColor: Record<string, string> = {
  open:    'bg-green-100 text-green-700',
  pending: 'bg-yellow-100 text-yellow-700',
  closed:  'bg-gray-100 text-gray-500',
}

export default function InboxPage() {
  const [threads, setThreads] = useState<WhatsAppThread[]>([])
  const [status, setStatus] = useState('')
  const [loading, setLoading] = useState(true)
  const { lang, t } = useLang()

  useEffect(() => {
    setLoading(true)
    api.threads.list({ status, limit: 40 }).then((res: any) => {
      setThreads(Array.isArray(res) ? res : [])
    }).catch(() => setThreads([])).finally(() => setLoading(false))
  }, [status])

  return (
    <div className="flex flex-col flex-1 overflow-hidden">
      <Header title={t('صندوق الرسائل', 'Inbox')} />

      {/* Status tabs */}
      <div className="flex gap-1 px-6 pt-4 pb-0 border-b border-gray-100 bg-white">
        {statusTabs.map((tab) => (
          <button
            key={tab.key}
            onClick={() => setStatus(tab.key)}
            className={clsx(
              'px-4 py-2 text-sm font-medium rounded-t-lg border-b-2 transition-colors',
              status === tab.key
                ? 'border-brand-600 text-brand-600'
                : 'border-transparent text-gray-500 hover:text-gray-700'
            )}
          >
            {lang === 'ar' ? tab.label.ar : tab.label.en}
          </button>
        ))}
      </div>

      {/* Thread list */}
      <div className="flex-1 overflow-y-auto bg-white">
        {loading ? (
          <div className="flex items-center justify-center h-40 text-sm text-gray-400">
            {t('جاري التحميل...', 'Loading...')}
          </div>
        ) : threads.length === 0 ? (
          <div className="flex items-center justify-center h-40 text-sm text-gray-400">
            {t('لا توجد محادثات', 'No threads')}
          </div>
        ) : (
          <ul className="divide-y divide-gray-50">
            {threads.map((thread) => (
              <li key={thread.id}>
                <Link
                  href={`/inbox/${thread.id}`}
                  className="flex items-center gap-4 px-6 py-4 hover:bg-gray-50 transition-colors"
                >
                  {/* Avatar */}
                  <div className="w-10 h-10 rounded-full bg-brand-100 text-brand-700 font-semibold text-sm flex items-center justify-center shrink-0">
                    {thread.contact?.full_name?.[0]?.toUpperCase() ?? '?'}
                  </div>

                  {/* Info */}
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center justify-between gap-2">
                      <p className="font-medium text-sm text-gray-900 truncate">
                        {thread.contact?.full_name ?? thread.contact_id}
                      </p>
                      <span className={clsx('text-[10px] font-medium px-2 py-0.5 rounded-full shrink-0', statusColor[thread.thread_status])}>
                        {thread.thread_status}
                      </span>
                    </div>
                    <p className="text-xs text-gray-400 mt-0.5 truncate">
                      {thread.contact?.phone_wa} · {thread.message_count} {t('رسالة', 'messages')}
                    </p>
                    {thread.ai_summary && (
                      <p className="text-xs text-gray-500 mt-1 truncate italic">{thread.ai_summary}</p>
                    )}
                  </div>

                  {/* Time */}
                  {thread.last_message_at && (
                    <span className="text-xs text-gray-300 shrink-0">
                      {new Date(thread.last_message_at).toLocaleDateString(lang === 'ar' ? 'ar-AE' : 'en-AE', {
                        month: 'short', day: 'numeric'
                      })}
                    </span>
                  )}
                </Link>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  )
}

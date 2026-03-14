'use client'
import { useEffect, useState } from 'react'
import { useParams, useRouter } from 'next/navigation'
import { Header } from '@/components/layout/Header'
import { api } from '@/lib/api'
import { useLang } from '@/context/LangContext'
import type { WhatsAppThread, WhatsAppMessage } from '@/types'
import clsx from 'clsx'

export default function ThreadPage() {
  const { id } = useParams<{ id: string }>()
  const [thread, setThread] = useState<WhatsAppThread | null>(null)
  const [messages, setMessages] = useState<WhatsAppMessage[]>([])
  const [summary, setSummary] = useState('')
  const [summarizing, setSummarizing] = useState(false)
  const [loading, setLoading] = useState(true)
  const { lang, t } = useLang()
  const router = useRouter()

  useEffect(() => {
    if (!id) return
    Promise.all([
      api.threads.get(id) as Promise<WhatsAppThread>,
      api.threads.messages(id) as Promise<WhatsAppMessage[]>,
    ]).then(([t, msgs]) => {
      setThread(t)
      setMessages(Array.isArray(msgs) ? msgs : [])
    }).finally(() => setLoading(false))
  }, [id])

  const handleSummarize = async () => {
    setSummarizing(true)
    try {
      const res = await api.ai.summarize(id) as { summary: string }
      setSummary(res.summary)
    } catch {
      setSummary(t('حدث خطأ أثناء التلخيص', 'Summarization failed'))
    } finally {
      setSummarizing(false)
    }
  }

  const handleClose = async () => {
    await api.threads.close(id).catch(() => {})
    router.push('/inbox')
  }

  if (loading) {
    return (
      <div className="flex-1 flex items-center justify-center text-gray-400 text-sm">
        {t('جاري التحميل...', 'Loading...')}
      </div>
    )
  }

  return (
    <div className="flex flex-col flex-1 overflow-hidden">
      <Header title={thread?.contact?.full_name ?? t('المحادثة', 'Thread')} />

      {/* Thread meta bar */}
      <div className="bg-white border-b border-gray-100 px-6 py-3 flex items-center justify-between gap-4">
        <div className="flex items-center gap-3">
          <div className="w-9 h-9 rounded-full bg-brand-100 text-brand-700 font-semibold text-sm flex items-center justify-center">
            {thread?.contact?.full_name?.[0]?.toUpperCase() ?? '?'}
          </div>
          <div>
            <p className="font-medium text-sm text-gray-900">{thread?.contact?.full_name}</p>
            <p className="text-xs text-gray-400">{thread?.contact?.phone_wa}</p>
          </div>
        </div>

        <div className="flex items-center gap-2">
          <button
            onClick={handleSummarize}
            disabled={summarizing || messages.length === 0}
            className="flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium bg-indigo-50 text-indigo-700 rounded-lg hover:bg-indigo-100 transition-colors disabled:opacity-50"
          >
            {summarizing ? '⏳' : '✨'}
            {summarizing
              ? t('جاري التلخيص...', 'Summarizing...')
              : t('تلخيص بالذكاء الاصطناعي', 'AI Summarize')}
          </button>

          {thread?.thread_status !== 'closed' && (
            <button
              onClick={handleClose}
              className="px-3 py-1.5 text-xs font-medium bg-gray-100 text-gray-600 rounded-lg hover:bg-gray-200 transition-colors"
            >
              {t('إغلاق', 'Close thread')}
            </button>
          )}
        </div>
      </div>

      {/* AI Summary */}
      {summary && (
        <div className="mx-6 mt-4 p-4 bg-indigo-50 border border-indigo-100 rounded-xl text-sm text-indigo-800">
          <p className="font-semibold text-xs text-indigo-500 mb-1">{t('ملخص الذكاء الاصطناعي', 'AI Summary')}</p>
          <p>{summary}</p>
        </div>
      )}

      {/* Messages */}
      <div className="flex-1 overflow-y-auto px-6 py-4 space-y-3">
        {messages.length === 0 ? (
          <p className="text-center text-sm text-gray-400 mt-10">
            {t('لا توجد رسائل', 'No messages yet')}
          </p>
        ) : (
          messages.map((msg) => {
            const isInbound = msg.direction === 'inbound'
            return (
              <div
                key={msg.id}
                className={clsx('flex', isInbound ? 'justify-start' : 'justify-end')}
              >
                <div
                  className={clsx(
                    'max-w-xs md:max-w-md px-4 py-2.5 rounded-2xl text-sm',
                    isInbound
                      ? 'bg-white border border-gray-100 text-gray-800 rounded-tl-sm'
                      : 'bg-brand-600 text-white rounded-tr-sm'
                  )}
                >
                  <p className="leading-relaxed">{msg.body}</p>
                  <p className={clsx('text-[10px] mt-1', isInbound ? 'text-gray-400' : 'text-blue-100')}>
                    {new Date(msg.sent_at).toLocaleTimeString(lang === 'ar' ? 'ar-AE' : 'en-AE', {
                      hour: '2-digit', minute: '2-digit'
                    })}
                  </p>
                </div>
              </div>
            )
          })
        )}
      </div>

      {/* Read-only notice (WhatsApp Sender is Enterprise) */}
      <div className="px-6 py-3 bg-gray-50 border-t border-gray-100 text-center">
        <p className="text-xs text-gray-400">
          {t(
            'إرسال الرسائل متوفر في الإصدار Enterprise',
            'Sending messages is available in the Enterprise edition'
          )}
        </p>
      </div>
    </div>
  )
}

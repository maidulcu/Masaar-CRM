'use client'
import { useEffect, useState, useCallback } from 'react'
import { Header } from '@/components/layout/Header'
import { api } from '@/lib/api'
import { useLang } from '@/context/LangContext'
import type { Contact, PaginatedResult } from '@/types'
import clsx from 'clsx'

export default function ContactsPage() {
  const [result, setResult] = useState<PaginatedResult<Contact> | null>(null)
  const [search, setSearch] = useState('')
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(true)
  const { lang, t } = useLang()

  const load = useCallback(async () => {
    setLoading(true)
    try {
      const data = await api.contacts.list({ search, page, limit: 20 }) as PaginatedResult<Contact>
      setResult(data)
    } catch {
      setResult(null)
    } finally {
      setLoading(false)
    }
  }, [search, page])

  useEffect(() => { load() }, [load])

  const handleSearch = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setPage(1)
    load()
  }

  const scoreColor = (score: number) => {
    if (score >= 80) return 'text-green-600 bg-green-50'
    if (score >= 50) return 'text-yellow-700 bg-yellow-50'
    return 'text-gray-500 bg-gray-50'
  }

  const contacts = result?.data ?? []
  const total = result?.total ?? 0
  const totalPages = Math.ceil(total / 20)

  return (
    <div className="flex flex-col flex-1 overflow-hidden">
      <Header title={t('جهات الاتصال', 'Contacts')} />

      <div className="flex-1 overflow-auto">
        <div className="max-w-5xl mx-auto px-6 py-6">

          {/* Search bar */}
          <form onSubmit={handleSearch} className="flex gap-2 mb-6">
            <input
              type="search"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              placeholder={t('ابحث بالاسم أو الهاتف أو البريد...', 'Search by name, phone, or email...')}
              className="flex-1 px-4 py-2 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-brand-500 focus:border-transparent"
            />
            <button
              type="submit"
              className="px-4 py-2 bg-brand-600 text-white text-sm font-medium rounded-lg hover:bg-brand-700 transition-colors"
            >
              {t('بحث', 'Search')}
            </button>
          </form>

          {/* Stats */}
          <p className="text-xs text-gray-400 mb-4">
            {total} {t('جهة اتصال', 'contacts')}
          </p>

          {/* Table */}
          {loading ? (
            <div className="text-center py-16 text-sm text-gray-400">
              {t('جاري التحميل...', 'Loading...')}
            </div>
          ) : contacts.length === 0 ? (
            <div className="text-center py-16 text-sm text-gray-400">
              {t('لا توجد نتائج', 'No results found')}
            </div>
          ) : (
            <div className="bg-white rounded-xl border border-gray-100 overflow-hidden">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b border-gray-100 bg-gray-50 text-gray-500 text-xs font-medium">
                    <th className="text-start px-4 py-3">{t('الاسم', 'Name')}</th>
                    <th className="text-start px-4 py-3">{t('واتساب', 'WhatsApp')}</th>
                    <th className="text-start px-4 py-3">{t('البريد', 'Email')}</th>
                    <th className="text-start px-4 py-3">{t('اللغة', 'Lang')}</th>
                    <th className="text-start px-4 py-3">{t('النقاط', 'Score')}</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-gray-50">
                  {contacts.map((c) => (
                    <tr key={c.id} className="hover:bg-gray-50 transition-colors">
                      <td className="px-4 py-3">
                        <div className="flex items-center gap-2">
                          <div className="w-7 h-7 rounded-full bg-brand-100 text-brand-700 text-xs font-semibold flex items-center justify-center shrink-0">
                            {c.full_name[0]?.toUpperCase()}
                          </div>
                          <span className="font-medium text-gray-900 truncate max-w-[160px]">{c.full_name}</span>
                        </div>
                      </td>
                      <td className="px-4 py-3 text-gray-500 font-mono text-xs">{c.phone_wa}</td>
                      <td className="px-4 py-3 text-gray-500 truncate max-w-[180px]">{c.email || '—'}</td>
                      <td className="px-4 py-3">
                        <span className="text-xs px-1.5 py-0.5 bg-gray-100 text-gray-600 rounded">
                          {c.language.toUpperCase()}
                        </span>
                      </td>
                      <td className="px-4 py-3">
                        <span className={clsx('text-xs font-semibold px-2 py-0.5 rounded-full', scoreColor(c.lead_score))}>
                          {c.lead_score}
                        </span>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}

          {/* Pagination */}
          {totalPages > 1 && (
            <div className="flex justify-center gap-2 mt-6">
              <button
                onClick={() => setPage((p) => Math.max(1, p - 1))}
                disabled={page === 1}
                className="px-3 py-1.5 text-xs rounded-lg border border-gray-200 disabled:opacity-40 hover:bg-gray-50"
              >
                {t('السابق', 'Previous')}
              </button>
              <span className="px-3 py-1.5 text-xs text-gray-500">
                {page} / {totalPages}
              </span>
              <button
                onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                disabled={page === totalPages}
                className="px-3 py-1.5 text-xs rounded-lg border border-gray-200 disabled:opacity-40 hover:bg-gray-50"
              >
                {t('التالي', 'Next')}
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

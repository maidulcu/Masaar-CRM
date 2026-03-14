'use client'
import { useEffect, useState, useCallback } from 'react'
import { Header } from '@/components/layout/Header'
import { Modal, FormField, FormError } from '@/components/ui/Modal'
import { Pagination } from '@/components/ui/Pagination'
import { api } from '@/lib/api'
import { useLang } from '@/context/LangContext'
import type { Contact, PaginatedResult } from '@/types'
import clsx from 'clsx'

export default function ContactsPage() {
  const [result, setResult] = useState<PaginatedResult<Contact> | null>(null)
  const [search, setSearch] = useState('')
  const [page, setPage] = useState(1)
  const [loading, setLoading] = useState(true)
  const [showAddModal, setShowAddModal] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [submitError, setSubmitError] = useState('')
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

  const handleAddContact = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setSubmitError('')
    setSubmitting(true)
    
    const formData = new FormData(e.currentTarget)
    const data = {
      phone_wa: formData.get('phone_wa'),
      full_name: formData.get('full_name'),
      email: formData.get('email') || undefined,
      language: formData.get('language') || 'en',
    }

    try {
      await api.contacts.create(data)
      setShowAddModal(false)
      load()
    } catch (err: any) {
      setSubmitError(err.message || t('حدث خطأ', 'Something went wrong'))
    } finally {
      setSubmitting(false)
    }
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

          {/* Search bar + Add button */}
          <div className="flex gap-2 mb-6">
            <form onSubmit={handleSearch} className="flex-1 flex gap-2">
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
            <button
              onClick={() => setShowAddModal(true)}
              className="px-4 py-2 bg-green-600 text-white text-sm font-medium rounded-lg hover:bg-green-700 transition-colors"
            >
              + {t('إضافة', 'Add')}
            </button>
          </div>

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

          <Pagination
            page={page}
            totalPages={totalPages}
            onPageChange={setPage}
          />
        </div>
      </div>

      {/* Add Contact Modal */}
      <Modal
        open={showAddModal}
        onClose={() => setShowAddModal(false)}
        title={t('إضافة جهة اتصال', 'Add Contact')}
      >
        <form onSubmit={handleAddContact} className="space-y-4">
          <FormField label={t('رقم الواتساب *', 'WhatsApp Number *')}>
            <input
              type="tel"
              name="phone_wa"
              required
              placeholder={t('971501234567', '971501234567')}
              className="w-full px-3 py-2 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-brand-500 focus:border-transparent"
            />
          </FormField>
          
          <FormField label={t('الاسم الكامل *', 'Full Name *')}>
            <input
              type="text"
              name="full_name"
              required
              placeholder={t('أدخل الاسم الكامل', 'Enter full name')}
              className="w-full px-3 py-2 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-brand-500 focus:border-transparent"
            />
          </FormField>
          
          <FormField label={t('البريد الإلكتروني', 'Email')}>
            <input
              type="email"
              name="email"
              placeholder={t('example@email.com', 'example@email.com')}
              className="w-full px-3 py-2 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-brand-500 focus:border-transparent"
            />
          </FormField>
          
          <FormField label={t('اللغة', 'Language')}>
            <select
              name="language"
              defaultValue="en"
              className="w-full px-3 py-2 border border-gray-200 rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-brand-500 focus:border-transparent"
            >
              <option value="en">English</option>
              <option value="ar">العربية</option>
            </select>
          </FormField>

          {submitError && <FormError error={submitError} />}

          <div className="flex gap-2 pt-2">
            <button
              type="button"
              onClick={() => setShowAddModal(false)}
              className="flex-1 px-4 py-2 border border-gray-200 text-gray-600 text-sm font-medium rounded-lg hover:bg-gray-50 transition-colors"
            >
              {t('إلغاء', 'Cancel')}
            </button>
            <button
              type="submit"
              disabled={submitting}
              className="flex-1 px-4 py-2 bg-brand-600 text-white text-sm font-medium rounded-lg hover:bg-brand-700 transition-colors disabled:opacity-60"
            >
              {submitting ? t('جاري الإضافة...', 'Adding...') : t('إضافة', 'Add')}
            </button>
          </div>
        </form>
      </Modal>
    </div>
  )
}

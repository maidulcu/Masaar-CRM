'use client'
import { useEffect, useState } from 'react'
import Link from 'next/link'
import { Header } from '@/components/layout/Header'
import { api } from '@/lib/api'
import { useLang } from '@/context/LangContext'
import type { Deal, PaginatedResult } from '@/types'
import clsx from 'clsx'

const stageColor: Record<string, string> = {
  open:   'bg-blue-100 text-blue-700',
  won:    'bg-green-100 text-green-700',
  lost:   'bg-red-100 text-red-700',
}

const stageLabel: Record<string, { en: string; ar: string }> = {
  open:   { en: 'Open',   ar: 'مفتوحة' },
  won:    { en: 'Won',    ar: 'مكسوبة' },
  lost:   { en: 'Lost',   ar: 'خاسرة' },
}

export default function DealsPage() {
  const [deals, setDeals] = useState<Deal[]>([])
  const [loading, setLoading] = useState(true)
  const { lang, t } = useLang()

  useEffect(() => {
    api.deals.list({ limit: 50 }).then((res: any) => {
      const data = res?.data ?? res
      setDeals(Array.isArray(data) ? data : [])
    }).catch(() => setDeals([])).finally(() => setLoading(false))
  }, [])

  const formatAmount = (amount: number, currency: string) =>
    new Intl.NumberFormat(lang === 'ar' ? 'ar-AE' : 'en-AE', {
      style: 'currency', currency: currency || 'AED', maximumFractionDigits: 0,
    }).format(amount)

  return (
    <div className="flex flex-col flex-1 overflow-hidden">
      <Header title={t('الصفقات', 'Deals')} />

      <div className="flex-1 overflow-y-auto p-6">
        {loading ? (
          <div className="flex items-center justify-center h-40 text-sm text-gray-400">
            {t('جاري التحميل...', 'Loading...')}
          </div>
        ) : deals.length === 0 ? (
          <div className="flex flex-col items-center justify-center h-40 gap-2">
            <p className="text-sm text-gray-400">{t('لا توجد صفقات بعد', 'No deals yet')}</p>
            <p className="text-xs text-gray-300">{t('أنشئ صفقة من صفحة خط الأنابيب', 'Create a deal from the Pipeline page')}</p>
          </div>
        ) : (
          <div className="bg-white rounded-xl border border-gray-100 overflow-hidden">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-gray-100 text-xs text-gray-400 uppercase tracking-wide">
                  <th className="text-start px-5 py-3 font-medium">{t('العنوان', 'Title')}</th>
                  <th className="text-start px-5 py-3 font-medium">{t('المرحلة', 'Stage')}</th>
                  <th className="text-start px-5 py-3 font-medium">{t('القيمة', 'Amount')}</th>
                  <th className="text-start px-5 py-3 font-medium">{t('الاحتمالية', 'Probability')}</th>
                  <th className="text-start px-5 py-3 font-medium">{t('تاريخ الإغلاق', 'Close Date')}</th>
                  <th className="px-5 py-3" />
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-50">
                {deals.map((deal) => (
                  <tr key={deal.id} className="hover:bg-gray-50 transition-colors">
                    <td className="px-5 py-3.5 font-medium text-gray-900">{deal.title}</td>
                    <td className="px-5 py-3.5">
                      <span className={clsx('text-xs font-medium px-2 py-0.5 rounded-full', stageColor[deal.stage])}>
                        {lang === 'ar' ? stageLabel[deal.stage]?.ar : stageLabel[deal.stage]?.en}
                      </span>
                    </td>
                    <td className="px-5 py-3.5 text-gray-700 font-medium">
                      {deal.amount ? formatAmount(deal.amount, deal.currency) : '—'}
                    </td>
                    <td className="px-5 py-3.5">
                      <div className="flex items-center gap-2">
                        <div className="w-16 h-1.5 rounded-full bg-gray-100 overflow-hidden">
                          <div
                            className="h-full bg-brand-500 rounded-full"
                            style={{ width: `${deal.probability}%` }}
                          />
                        </div>
                        <span className="text-xs text-gray-400">{deal.probability}%</span>
                      </div>
                    </td>
                    <td className="px-5 py-3.5 text-gray-400 text-xs">
                      {deal.close_date
                        ? new Date(deal.close_date).toLocaleDateString(
                            lang === 'ar' ? 'ar-AE' : 'en-AE',
                            { year: 'numeric', month: 'short', day: 'numeric' }
                          )
                        : '—'}
                    </td>
                    <td className="px-5 py-3.5 text-end">
                      <Link
                        href={`/deals/${deal.id}`}
                        className="text-xs text-brand-600 hover:underline font-medium"
                      >
                        {t('عرض', 'View')}
                      </Link>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  )
}

'use client'
import { useEffect, useState } from 'react'
import Link from 'next/link'
import { Header } from '@/components/layout/Header'
import { api } from '@/lib/api'
import { useLang } from '@/context/LangContext'

interface Stats {
  total_contacts: number
  active_leads: number
  new_leads_week: number
  open_threads: number
  open_deals: number
  open_deals_value: number
  won_deals: number
  won_deals_value: number
}

function StatCard({
  label,
  value,
  sub,
  href,
  color,
}: {
  label: string
  value: string | number
  sub?: string
  href: string
  color: string
}) {
  return (
    <Link href={href} className="block group">
      <div className="bg-white rounded-xl border border-gray-200 p-5 hover:border-brand-300 hover:shadow-sm transition-all">
        <p className="text-xs font-medium text-gray-500 mb-2">{label}</p>
        <p className={`text-2xl font-bold ${color}`}>{value}</p>
        {sub && <p className="text-xs text-gray-400 mt-1">{sub}</p>}
      </div>
    </Link>
  )
}

function SkeletonCard() {
  return (
    <div className="bg-white rounded-xl border border-gray-200 p-5 animate-pulse">
      <div className="h-3 w-24 bg-gray-100 rounded mb-3" />
      <div className="h-7 w-16 bg-gray-100 rounded" />
    </div>
  )
}

export default function DashboardPage() {
  const [stats, setStats] = useState<Stats | null>(null)
  const [loading, setLoading] = useState(true)
  const { t } = useLang()

  useEffect(() => {
    api.stats.overview()
      .then((data) => setStats(data as Stats))
      .catch(() => {})
      .finally(() => setLoading(false))
  }, [])

  const fmt = (n: number) =>
    new Intl.NumberFormat('en-AE', { maximumFractionDigits: 0 }).format(n)

  return (
    <div className="flex flex-col flex-1 min-h-0">
      <Header title={t('لوحة التحكم', 'Dashboard')} />

      <main className="flex-1 overflow-y-auto p-6 bg-gray-50">
        <div className="max-w-4xl space-y-6">

          {/* Top row */}
          <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
            {loading ? (
              Array.from({ length: 4 }).map((_, i) => <SkeletonCard key={i} />)
            ) : (
              <>
                <StatCard
                  label={t('إجمالي جهات الاتصال', 'Total Contacts')}
                  value={fmt(stats?.total_contacts ?? 0)}
                  href="/contacts"
                  color="text-gray-800"
                />
                <StatCard
                  label={t('Leads النشطة', 'Active Leads')}
                  value={fmt(stats?.active_leads ?? 0)}
                  sub={t(`${stats?.new_leads_week ?? 0} هذا الأسبوع`, `${stats?.new_leads_week ?? 0} this week`)}
                  href="/pipeline"
                  color="text-brand-600"
                />
                <StatCard
                  label={t('محادثات مفتوحة', 'Open Threads')}
                  value={fmt(stats?.open_threads ?? 0)}
                  href="/inbox"
                  color="text-blue-600"
                />
                <StatCard
                  label={t('صفقات مفتوحة', 'Open Deals')}
                  value={fmt(stats?.open_deals ?? 0)}
                  sub={stats?.open_deals_value ? `AED ${fmt(stats.open_deals_value)}` : undefined}
                  href="/deals"
                  color="text-yellow-600"
                />
              </>
            )}
          </div>

          {/* Won deals highlight */}
          {!loading && (stats?.won_deals ?? 0) > 0 && (
            <Link href="/deals" className="block group">
              <div className="bg-green-50 border border-green-200 rounded-xl p-5 flex items-center justify-between hover:border-green-400 transition-all">
                <div>
                  <p className="text-xs font-medium text-green-600 mb-1">
                    {t('الصفقات المكتسبة', 'Won Deals')}
                  </p>
                  <p className="text-2xl font-bold text-green-700">{fmt(stats?.won_deals ?? 0)}</p>
                </div>
                <div className="text-right">
                  <p className="text-xs text-green-500 mb-1">{t('إجمالي القيمة', 'Total Value')}</p>
                  <p className="text-xl font-bold text-green-700">AED {fmt(stats?.won_deals_value ?? 0)}</p>
                </div>
                <div className="w-10 h-10 rounded-full bg-green-100 flex items-center justify-center ml-4 shrink-0">
                  <svg className="w-5 h-5 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </div>
              </div>
            </Link>
          )}

          {/* Quick links */}
          <div>
            <p className="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-3">
              {t('روابط سريعة', 'Quick Links')}
            </p>
            <div className="grid grid-cols-2 sm:grid-cols-4 gap-3">
              {[
                { href: '/pipeline', label: { en: 'Pipeline', ar: 'خط الأنابيب' }, icon: '📊' },
                { href: '/inbox',    label: { en: 'Inbox',    ar: 'الرسائل'     }, icon: '💬' },
                { href: '/contacts', label: { en: 'Contacts', ar: 'جهات الاتصال' }, icon: '👤' },
                { href: '/deals',    label: { en: 'Deals',    ar: 'الصفقات'     }, icon: '🤝' },
              ].map((item) => (
                <Link
                  key={item.href}
                  href={item.href}
                  className="flex items-center gap-2 bg-white border border-gray-200 rounded-xl px-4 py-3 text-sm font-medium text-gray-700 hover:border-brand-300 hover:text-brand-700 transition-all"
                >
                  <span>{item.icon}</span>
                  <span>{t(item.label.ar, item.label.en)}</span>
                </Link>
              ))}
            </div>
          </div>

        </div>
      </main>
    </div>
  )
}

'use client'
import { useDroppable } from '@dnd-kit/core'
import { SortableContext, verticalListSortingStrategy } from '@dnd-kit/sortable'
import type { Lead, LeadStage } from '@/types'
import { KanbanCard } from './Card'
import { useLang } from '@/context/LangContext'
import clsx from 'clsx'

const stageConfig: Record<LeadStage, { label: { en: string; ar: string }; color: string }> = {
  new:       { label: { en: 'New',       ar: 'جديد'      }, color: 'bg-gray-100 text-gray-600' },
  contacted: { label: { en: 'Contacted', ar: 'تم التواصل' }, color: 'bg-blue-100 text-blue-600' },
  qualified: { label: { en: 'Qualified', ar: 'مؤهل'      }, color: 'bg-indigo-100 text-indigo-600' },
  proposal:  { label: { en: 'Proposal',  ar: 'عرض'       }, color: 'bg-yellow-100 text-yellow-700' },
  won:       { label: { en: 'Won',       ar: 'مكسب'      }, color: 'bg-green-100 text-green-700' },
  lost:      { label: { en: 'Lost',      ar: 'خسارة'     }, color: 'bg-red-100 text-red-600' },
}

interface Props {
  stage: LeadStage
  leads: Lead[]
}

export function KanbanColumn({ stage, leads }: Props) {
  const { setNodeRef, isOver } = useDroppable({ id: stage })
  const { lang, t } = useLang()
  const config = stageConfig[stage]

  const totalValue = leads.reduce((sum, l) => sum + l.deal_value, 0)
  const currency = leads[0]?.currency ?? 'AED'

  return (
    <div className="flex flex-col w-64 shrink-0">
      {/* Column header */}
      <div className="flex items-center justify-between mb-3 px-1">
        <div className="flex items-center gap-2">
          <span className={clsx('text-xs font-semibold px-2 py-0.5 rounded-full', config.color)}>
            {lang === 'ar' ? config.label.ar : config.label.en}
          </span>
          <span className="text-xs text-gray-400 font-medium">{leads.length}</span>
        </div>
        {totalValue > 0 && (
          <span className="text-xs text-gray-400">{currency} {totalValue.toLocaleString()}</span>
        )}
      </div>

      {/* Drop zone */}
      <div
        ref={setNodeRef}
        className={clsx(
          'flex-1 min-h-32 rounded-xl p-2 space-y-2 transition-colors',
          isOver ? 'bg-brand-50 ring-2 ring-brand-200' : 'bg-gray-100/60'
        )}
      >
        <SortableContext items={leads.map((l) => l.id)} strategy={verticalListSortingStrategy}>
          {leads.map((lead) => (
            <KanbanCard key={lead.id} lead={lead} />
          ))}
        </SortableContext>

        {leads.length === 0 && (
          <div className="flex items-center justify-center h-20 text-xs text-gray-400">
            {t('أفلت هنا', 'Drop here')}
          </div>
        )}
      </div>
    </div>
  )
}

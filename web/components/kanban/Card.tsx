'use client'
import { useSortable } from '@dnd-kit/sortable'
import { CSS } from '@dnd-kit/utilities'
import type { Lead } from '@/types'
import { useLang } from '@/context/LangContext'
import clsx from 'clsx'

const sourceColors: Record<string, string> = {
  whatsapp: 'bg-green-100 text-green-700',
  web:       'bg-blue-100 text-blue-700',
  referral:  'bg-purple-100 text-purple-700',
  event:     'bg-orange-100 text-orange-700',
}

interface Props {
  lead: Lead
  onOpen?: (lead: Lead) => void
}

export function KanbanCard({ lead, onOpen }: Props) {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } =
    useSortable({ id: lead.id })
  const { t } = useLang()

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  }

  return (
    <div
      ref={setNodeRef}
      style={style}
      {...attributes}
      className={clsx(
        'bg-white rounded-xl border border-gray-100 shadow-sm select-none',
        isDragging && 'opacity-50 shadow-lg ring-2 ring-brand-400'
      )}
    >
      {/* Clickable body — opens notes modal */}
      <div
        className="p-3 cursor-pointer"
        onClick={() => onOpen?.(lead)}
      >
        {/* Contact name */}
        <p className="font-medium text-sm text-gray-900 truncate">
          {lead.contact?.full_name ?? t('جهة اتصال غير معروفة', 'Unknown contact')}
        </p>

        {/* Phone */}
        {lead.contact?.phone_wa && (
          <p className="text-xs text-gray-400 mt-0.5 truncate">{lead.contact.phone_wa}</p>
        )}

        {/* Deal value */}
        <div className="flex items-center justify-between mt-2">
          <span className="text-sm font-semibold text-gray-700">
            {lead.currency} {lead.deal_value.toLocaleString()}
          </span>
          {lead.source && (
            <span className={clsx('text-[10px] font-medium px-1.5 py-0.5 rounded-full', sourceColors[lead.source] ?? 'bg-gray-100 text-gray-600')}>
              {lead.source}
            </span>
          )}
        </div>

        {/* Lead score */}
        {lead.contact?.lead_score != null && lead.contact.lead_score > 0 && (
          <div className="mt-2 flex items-center gap-1">
            <div className="flex-1 h-1 bg-gray-100 rounded-full overflow-hidden">
              <div
                className="h-full bg-brand-500 rounded-full"
                style={{ width: `${lead.contact.lead_score}%` }}
              />
            </div>
            <span className="text-[10px] text-gray-400">{lead.contact.lead_score}</span>
          </div>
        )}

        {/* Notes preview */}
        {lead.notes && (
          <p className="mt-2 text-[10px] text-gray-400 truncate italic">
            {lead.notes}
          </p>
        )}
      </div>

      {/* Drag handle — only this triggers DnD */}
      <div
        {...listeners}
        className="flex items-center justify-center py-1 border-t border-gray-50 cursor-grab active:cursor-grabbing"
        title={t('اسحب للنقل', 'Drag to move')}
      >
        <svg className="w-4 h-4 text-gray-300" fill="currentColor" viewBox="0 0 24 24">
          <path d="M8 6a1.5 1.5 0 110-3 1.5 1.5 0 010 3zm8 0a1.5 1.5 0 110-3 1.5 1.5 0 010 3zM8 13.5a1.5 1.5 0 110-3 1.5 1.5 0 010 3zm8 0a1.5 1.5 0 110-3 1.5 1.5 0 010 3zM8 21a1.5 1.5 0 110-3 1.5 1.5 0 010 3zm8 0a1.5 1.5 0 110-3 1.5 1.5 0 010 3z"/>
        </svg>
      </div>
    </div>
  )
}

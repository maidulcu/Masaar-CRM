'use client'
import { useLang } from '@/context/LangContext'
import clsx from 'clsx'

interface PaginationProps {
  page: number
  totalPages: number
  onPageChange: (page: number) => void
}

export function Pagination({ page, totalPages, onPageChange }: PaginationProps) {
  const { t } = useLang()

  if (totalPages <= 1) return null

  return (
    <div className="flex justify-center items-center gap-2 mt-6">
      <button
        onClick={() => onPageChange(page - 1)}
        disabled={page === 1}
        className={clsx(
          "px-3 py-1.5 text-xs font-medium rounded-lg border transition-colors",
          page === 1
            ? "border-gray-100 text-gray-300 cursor-not-allowed"
            : "border-gray-200 text-gray-600 hover:bg-gray-50"
        )}
      >
        {t('السابق', 'Previous')}
      </button>
      
      <span className="px-3 py-1.5 text-xs text-gray-500">
        {page} / {totalPages}
      </span>
      
      <button
        onClick={() => onPageChange(page + 1)}
        disabled={page === totalPages}
        className={clsx(
          "px-3 py-1.5 text-xs font-medium rounded-lg border transition-colors",
          page === totalPages
            ? "border-gray-100 text-gray-300 cursor-not-allowed"
            : "border-gray-200 text-gray-600 hover:bg-gray-50"
        )}
      >
        {t('التالي', 'Next')}
      </button>
    </div>
  )
}

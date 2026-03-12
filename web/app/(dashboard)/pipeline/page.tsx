'use client'
import { useEffect, useState, useCallback } from 'react'
import {
  DndContext, DragEndEvent, DragOverlay, DragStartEvent,
  PointerSensor, useSensor, useSensors, closestCorners,
} from '@dnd-kit/core'
import { Header } from '@/components/layout/Header'
import { KanbanColumn } from '@/components/kanban/Column'
import { KanbanCard } from '@/components/kanban/Card'
import { api } from '@/lib/api'
import { useLang } from '@/context/LangContext'
import type { KanbanBoard, Lead, LeadStage } from '@/types'

const STAGES: LeadStage[] = ['new', 'contacted', 'qualified', 'proposal', 'won', 'lost']

export default function PipelinePage() {
  const [board, setBoard] = useState<KanbanBoard>({})
  const [activeCard, setActiveCard] = useState<Lead | null>(null)
  const [loading, setLoading] = useState(true)
  const { t } = useLang()

  const sensors = useSensors(
    useSensor(PointerSensor, { activationConstraint: { distance: 8 } })
  )

  const load = useCallback(async () => {
    try {
      const data = await api.leads.kanban() as KanbanBoard
      setBoard(data ?? {})
    } catch {
      // handle error silently
    } finally {
      setLoading(false)
    }
  }, [])

  useEffect(() => { load() }, [load])

  const findCard = (id: string): Lead | null => {
    for (const leads of Object.values(board)) {
      const found = leads?.find((l) => l.id === id)
      if (found) return found
    }
    return null
  }

  const findStage = (id: string): LeadStage | null => {
    for (const [stage, leads] of Object.entries(board)) {
      if (leads?.some((l) => l.id === id)) return stage as LeadStage
    }
    return null
  }

  const handleDragStart = (e: DragStartEvent) => {
    setActiveCard(findCard(String(e.active.id)))
  }

  const handleDragEnd = async (e: DragEndEvent) => {
    setActiveCard(null)
    const { active, over } = e
    if (!over) return

    const leadId = String(active.id)
    const targetStage = (STAGES.includes(over.id as LeadStage)
      ? over.id
      : findStage(String(over.id))) as LeadStage | null

    const currentStage = findStage(leadId)
    if (!targetStage || targetStage === currentStage) return

    // Optimistic update
    setBoard((prev) => {
      const next = { ...prev }
      const card = (next[currentStage!] ?? []).find((l) => l.id === leadId)
      if (!card) return prev
      next[currentStage!] = (next[currentStage!] ?? []).filter((l) => l.id !== leadId)
      next[targetStage] = [{ ...card, stage: targetStage }, ...(next[targetStage] ?? [])]
      return next
    })

    await api.leads.updateStage(leadId, targetStage).catch(() => load())
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
      <Header title={t('خط الأنابيب', 'Sales Pipeline')} />

      <div className="flex-1 overflow-x-auto p-6">
        <DndContext
          sensors={sensors}
          collisionDetection={closestCorners}
          onDragStart={handleDragStart}
          onDragEnd={handleDragEnd}
        >
          <div className="flex gap-4 min-w-max pb-4">
            {STAGES.map((stage) => (
              <KanbanColumn
                key={stage}
                stage={stage}
                leads={board[stage] ?? []}
              />
            ))}
          </div>

          <DragOverlay>
            {activeCard && <KanbanCard lead={activeCard} />}
          </DragOverlay>
        </DndContext>
      </div>
    </div>
  )
}

'use client'
import { create } from 'zustand'
import type { AuthUser } from '@/types'
import { getUser, getToken, saveSession, clearSession, updateUser } from '@/lib/auth'

interface AuthState {
  user: AuthUser | null
  token: string | null
  setSession: (token: string, refreshToken: string, user: AuthUser) => void
  updateUser: (updates: Partial<AuthUser>) => void
  logout: () => void
  init: () => void
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  token: null,

  init: () => {
    set({ user: getUser(), token: getToken() })
  },

  setSession: (token, refreshToken, user) => {
    saveSession(token, refreshToken, user)
    set({ user, token })
  },

  updateUser: (updates) => {
    updateUser(updates)
    set((state) => ({ user: state.user ? { ...state.user, ...updates } : null }))
  },

  logout: () => {
    clearSession()
    set({ user: null, token: null })
  },
}))

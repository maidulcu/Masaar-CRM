'use client'
import { create } from 'zustand'
import type { AuthUser } from '@/types'
import { getUser, getToken, saveSession, clearSession } from '@/lib/auth'

interface AuthState {
  user: AuthUser | null
  token: string | null
  setSession: (token: string, refreshToken: string, user: AuthUser) => void
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

  logout: () => {
    clearSession()
    set({ user: null, token: null })
  },
}))

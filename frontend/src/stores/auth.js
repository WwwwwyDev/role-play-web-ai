import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/services/api'

export const useAuthStore = defineStore('auth', () => {
  const user = ref(null)
  const token = ref(localStorage.getItem('token'))
  const isInitialized = ref(false)

  const isAuthenticated = computed(() => !!token.value)

  const login = async (email, password) => {
    try {
      const response = await api.post('/auth/login', { email, password })
      const { token: newToken, user: userData } = response.data
      
      token.value = newToken
      user.value = userData
      localStorage.setItem('token', newToken)
      
      return { success: true }
    } catch (error) {
      return { 
        success: false, 
        error: error.response?.data?.error || '登录失败' 
      }
    }
  }

  const register = async (username, email, password) => {
    try {
      const response = await api.post('/auth/register', { username, email, password })
      const { token: newToken, user: userData } = response.data
      
      token.value = newToken
      user.value = userData
      localStorage.setItem('token', newToken)
      
      return { success: true }
    } catch (error) {
      return { 
        success: false, 
        error: error.response?.data?.error || '注册失败' 
      }
    }
  }

  const logout = () => {
    user.value = null
    token.value = null
    localStorage.removeItem('token')
  }

  const fetchProfile = async () => {
    try {
      const response = await api.get('/auth/me')
      user.value = response.data.user
      isInitialized.value = true
      console.log('用户信息获取成功:', user.value)
      return { success: true }
    } catch (error) {
      console.log('用户信息获取失败:', error)
      logout()
      isInitialized.value = true
      return { success: false }
    }
  }

  const initializeAuth = async () => {
    if (token.value && !isInitialized.value) {
      await fetchProfile()
    } else if (!token.value) {
      isInitialized.value = true
    }
  }

  // 初始化时检查token
  if (token.value) {
    fetchProfile()
  } else {
    isInitialized.value = true
  }

  return {
    user,
    token,
    isAuthenticated,
    isInitialized,
    login,
    register,
    logout,
    fetchProfile,
    initializeAuth
  }
})

import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'Home',
      component: () => import('@/views/Home.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/chat/:conversationId?',
      name: 'Chat',
      component: () => import('@/views/Chat.vue'),
      meta: { requiresAuth: true }
    }
  ]
})

// 路由守卫
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  
  // 如果认证状态还未初始化，等待初始化完成
  if (!authStore.isInitialized) {
    try {
      await authStore.initializeAuth()
    } catch (error) {
      console.error('认证初始化失败:', error)
      // 即使初始化失败，也继续路由，避免卡住
    }
  }
  
  // 确保状态更新完成
  await new Promise(resolve => setTimeout(resolve, 10))
  
  // 路由逻辑
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/')
  } else if (to.name === 'Home' && authStore.isAuthenticated) {
    next('/chat')
  } else {
    next()
  }
})

export default router

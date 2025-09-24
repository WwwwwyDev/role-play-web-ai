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
  
  // 如果认证状态还未初始化，先初始化
  if (!authStore.isInitialized) {
    await authStore.initializeAuth()
  }
  
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next('/')
  } else if (to.name === 'Home' && authStore.isAuthenticated) {
    next('/chat')
  } else {
    next()
  }
})

export default router

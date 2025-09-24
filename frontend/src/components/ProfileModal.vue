<template>
  <div
    v-if="isVisible"
    class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
    @click="closeModal"
  >
    <div
      class="bg-white rounded-2xl p-6 max-w-2xl w-full mx-4 max-h-[80vh] overflow-hidden flex flex-col"
      @click.stop
    >
      <!-- Header -->
      <div class="flex items-center justify-between mb-6">
        <h3 class="text-2xl font-bold text-gray-900">个人资料</h3>
        <button
          @click="closeModal"
          class="p-2 text-gray-400 hover:bg-gray-100 rounded-lg transition-colors duration-200"
        >
          <X class="w-6 h-6" />
        </button>
      </div>

      <!-- Content -->
      <div class="flex-1 overflow-y-auto">
        <div class="grid lg:grid-cols-3 gap-6">
          <!-- Profile Info -->
          <div class="lg:col-span-1">
            <div class="bg-gray-50 rounded-xl p-6">
              <div class="text-center">
                <div class="w-20 h-20 bg-primary-100 rounded-full flex items-center justify-center mx-auto mb-4">
                  <User class="w-10 h-10 text-primary-600" />
                </div>
                <h2 class="text-xl font-semibold text-gray-900 mb-1">
                  {{ authStore.user?.username }}
                </h2>
                <p class="text-gray-600 mb-4">{{ authStore.user?.email }}</p>
                <p class="text-sm text-gray-500">
                  注册时间：{{ formatDate(authStore.user?.created_at) }}
                </p>
              </div>
            </div>
          </div>

          <!-- Account Settings -->
          <div class="lg:col-span-2">
            <div class="space-y-6">
              <div>
                <h4 class="text-lg font-semibold text-gray-900 mb-4">账户信息</h4>
                <div class="space-y-4">
                  <div class="bg-gray-50 rounded-lg p-4">
                    <label class="block text-sm font-medium text-gray-700 mb-1">用户名</label>
                    <p class="text-sm text-gray-900">{{ authStore.user?.username }}</p>
                  </div>
                  <div class="bg-gray-50 rounded-lg p-4">
                    <label class="block text-sm font-medium text-gray-700 mb-1">邮箱地址</label>
                    <p class="text-sm text-gray-900">{{ authStore.user?.email }}</p>
                  </div>
                </div>
              </div>
              
              <div class="border-t border-gray-200 pt-6">
                <h4 class="text-lg font-semibold text-gray-900 mb-4">对话管理</h4>
                <p class="text-sm text-gray-600 mb-4">查看和管理您的对话历史</p>
                <button
                  @click="goToChat"
                  class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-lg text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 transition-colors duration-200"
                >
                  <MessageCircle class="w-4 h-4 mr-2" />
                  查看对话历史
                </button>
              </div>

              <div class="border-t border-gray-200 pt-6">
                <h4 class="text-lg font-semibold text-gray-900 mb-4">账户操作</h4>
                <div class="space-y-3">
                  <button
                    @click="handleLogout"
                    class="w-full inline-flex items-center justify-center px-4 py-2 border border-red-300 text-sm font-medium rounded-lg text-red-700 bg-red-50 hover:bg-red-100 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 transition-colors duration-200"
                  >
                    <LogOut class="w-4 h-4 mr-2" />
                    退出登录
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { User, X, MessageCircle, LogOut } from 'lucide-vue-next'

const props = defineProps({
  isVisible: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close'])

const router = useRouter()
const authStore = useAuthStore()

const closeModal = () => {
  emit('close')
}

const goToChat = () => {
  closeModal()
  router.push('/chat')
}

const handleLogout = () => {
  authStore.logout()
  closeModal()
  router.push('/')
}

const formatDate = (timestamp) => {
  if (!timestamp) return ''
  
  // 确保正确处理时间戳
  let date
  if (typeof timestamp === 'string') {
    date = new Date(timestamp)
  } else {
    date = new Date(timestamp)
  }
  
  // 检查日期是否有效
  if (isNaN(date.getTime())) {
    return '时间格式错误'
  }
  
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}
</script>

<style scoped>
/* 确保弹出框在移动设备上也能正常显示 */
@media (max-width: 768px) {
  .max-w-2xl {
    max-width: 95vw;
  }
}
</style>

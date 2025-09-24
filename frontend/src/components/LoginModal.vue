<template>
  <div
    v-if="isVisible"
    class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
    @click="closeModal"
  >
    <div
      class="bg-white rounded-2xl p-8 max-w-md w-full mx-4 relative"
      @click.stop
    >
      <!-- 关闭按钮 -->
      <button
        @click="closeModal"
        class="absolute top-4 right-4 p-2 text-gray-400 hover:bg-gray-100 rounded-lg transition-colors duration-200"
      >
        <X class="w-6 h-6" />
      </button>

      <!-- 头部 -->
      <div class="text-center mb-8">
        <div class="mx-auto h-16 w-16 bg-primary-600 rounded-2xl flex items-center justify-center mb-4">
          <MessageCircle class="h-10 w-10 text-white" />
        </div>
        <h2 class="text-2xl font-bold text-gray-900 mb-2">
          登录到您的账户
        </h2>
        <p class="text-gray-600">
          还没有账户？
          <button
            @click="switchToRegister"
            class="font-medium text-primary-600 hover:text-primary-500 transition-colors duration-200"
          >
            立即注册
          </button>
        </p>
      </div>
      
      <!-- 登录表单 -->
      <form class="space-y-6" @submit.prevent="handleLogin">
        <div class="space-y-4">
          <div>
            <label for="email" class="block text-sm font-medium text-gray-700 mb-2">
              邮箱地址
            </label>
            <input
              id="email"
              v-model="form.email"
              type="email"
              required
              class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-colors duration-200"
              placeholder="请输入您的邮箱"
            />
          </div>
          
          <div>
            <label for="password" class="block text-sm font-medium text-gray-700 mb-2">
              密码
            </label>
            <input
              id="password"
              v-model="form.password"
              type="password"
              required
              class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-colors duration-200"
              placeholder="请输入您的密码"
            />
          </div>
        </div>

        <!-- 错误信息 -->
        <div v-if="error" class="bg-red-50 border border-red-200 rounded-xl p-4 animate-fade-in">
          <div class="flex">
            <AlertCircle class="h-5 w-5 text-red-400 flex-shrink-0 mt-0.5" />
            <div class="ml-3">
              <p class="text-sm text-red-800 font-medium">{{ error }}</p>
              <p class="text-xs text-red-600 mt-1">请检查您的邮箱和密码，然后重试</p>
            </div>
          </div>
        </div>

        <!-- 登录按钮 -->
        <button
          type="submit"
          :disabled="isLoading"
          class="w-full flex justify-center items-center py-3 px-4 border border-transparent text-sm font-medium rounded-xl text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200"
        >
          <Loader2 v-if="isLoading" class="h-5 w-5 text-white animate-spin mr-2" />
          {{ isLoading ? '登录中...' : '登录' }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { MessageCircle, AlertCircle, Loader2, X } from 'lucide-vue-next'

const props = defineProps({
  isVisible: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close', 'switch-to-register'])

const router = useRouter()
const authStore = useAuthStore()

const form = ref({
  email: '',
  password: ''
})

const error = ref('')
const isLoading = ref(false)

const closeModal = () => {
  emit('close')
  resetForm()
}

const switchToRegister = () => {
  emit('switch-to-register')
  resetForm()
}

const resetForm = () => {
  form.value = {
    email: '',
    password: ''
  }
  error.value = ''
  isLoading.value = false
}

const handleLogin = async () => {
  error.value = ''
  isLoading.value = true

  try {
    const result = await authStore.login(form.value.email, form.value.password)
    
    if (result.success) {
      // 登录成功，关闭弹窗并跳转
      closeModal()
      router.push('/chat')
    } else {
      // 登录失败，显示错误信息
      error.value = result.error || '登录失败，请重试'
    }
  } catch (err) {
    // 网络错误或其他异常
    error.value = '网络连接失败，请检查网络后重试'
  } finally {
    isLoading.value = false
  }
}

// 监听弹窗显示状态，重置表单
watch(() => props.isVisible, (newVal) => {
  if (newVal) {
    resetForm()
  }
})
</script>

<style scoped>
@keyframes fade-in {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-fade-in {
  animation: fade-in 0.3s ease-out;
}
</style>

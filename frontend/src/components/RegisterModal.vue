<template>
  <div
    v-if="isVisible"
    class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
    @click="closeModal"
  >
    <div
      class="bg-white rounded-2xl p-8 max-w-md w-full mx-4 relative max-h-[90vh] overflow-y-auto"
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
          创建新账户
        </h2>
        <p class="text-gray-600">
          已有账户？
          <button
            @click="switchToLogin"
            class="font-medium text-primary-600 hover:text-primary-500 transition-colors duration-200"
          >
            立即登录
          </button>
        </p>
      </div>
      
      <!-- 注册表单 -->
      <form class="space-y-6" @submit.prevent="handleRegister">
        <div class="space-y-4">
          <div>
            <label for="username" class="block text-sm font-medium text-gray-700 mb-2">
              用户名
            </label>
            <input
              id="username"
              v-model="form.username"
              type="text"
              required
              class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-colors duration-200"
              placeholder="请输入用户名"
            />
          </div>
          
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
              placeholder="请输入密码（至少6位）"
            />
          </div>
          
          <div>
            <label for="confirmPassword" class="block text-sm font-medium text-gray-700 mb-2">
              确认密码
            </label>
            <input
              id="confirmPassword"
              v-model="form.confirmPassword"
              type="password"
              required
              class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent transition-colors duration-200"
              placeholder="请再次输入密码"
            />
          </div>
        </div>

        <!-- 错误信息 -->
        <div v-if="error" class="bg-red-50 border border-red-200 rounded-xl p-4 animate-fade-in">
          <div class="flex">
            <AlertCircle class="h-5 w-5 text-red-400 flex-shrink-0 mt-0.5" />
            <div class="ml-3">
              <p class="text-sm text-red-800 font-medium">{{ error }}</p>
              <p class="text-xs text-red-600 mt-1">请检查您的信息并重试</p>
            </div>
          </div>
        </div>

        <!-- 注册按钮 -->
        <button
          type="submit"
          :disabled="isLoading || !isFormValid"
          class="w-full flex justify-center items-center py-3 px-4 border border-transparent text-sm font-medium rounded-xl text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200"
        >
          <Loader2 v-if="isLoading" class="h-5 w-5 text-white animate-spin mr-2" />
          {{ isLoading ? '注册中...' : '注册' }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { MessageCircle, AlertCircle, Loader2, X } from 'lucide-vue-next'

const props = defineProps({
  isVisible: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close', 'switch-to-login'])

const router = useRouter()
const authStore = useAuthStore()

const form = ref({
  username: '',
  email: '',
  password: '',
  confirmPassword: ''
})

const error = ref('')
const isLoading = ref(false)

const isFormValid = computed(() => {
  return form.value.username.length >= 3 &&
         form.value.email.includes('@') &&
         form.value.password.length >= 6 &&
         form.value.password === form.value.confirmPassword
})

const closeModal = () => {
  emit('close')
  resetForm()
}

const switchToLogin = () => {
  emit('switch-to-login')
  resetForm()
}

const resetForm = () => {
  form.value = {
    username: '',
    email: '',
    password: '',
    confirmPassword: ''
  }
  error.value = ''
  isLoading.value = false
}

const handleRegister = async () => {
  error.value = ''
  
  // 客户端验证
  if (form.value.password !== form.value.confirmPassword) {
    error.value = '两次输入的密码不一致'
    return
  }
  
  if (form.value.password.length < 6) {
    error.value = '密码长度至少为6位'
    return
  }
  
  if (form.value.username.length < 3) {
    error.value = '用户名长度至少为3位'
    return
  }
  
  isLoading.value = true

  try {
    const result = await authStore.register(
      form.value.username,
      form.value.email,
      form.value.password
    )
    
    if (result.success) {
      // 注册成功，关闭弹窗并跳转
      closeModal()
      router.push('/chat')
    } else {
      // 注册失败，显示错误信息
      error.value = result.error || '注册失败，请重试'
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

<template>
  <nav class="fixed top-0 left-0 right-0 bg-white shadow-sm border-b border-gray-200 z-50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center h-16">
        <!-- Logo -->
        <div class="flex items-center">
          <router-link to="/" class="flex items-center space-x-2">
            <div class="w-8 h-8 bg-primary-600 rounded-lg flex items-center justify-center">
              <MessageCircle class="w-5 h-5 text-white" />
            </div>
            <span class="text-xl font-bold text-gray-900">AI角色扮演</span>
          </router-link>
        </div>


        <!-- 用户菜单 -->
        <div class="flex items-center space-x-4">
          <div v-if="authStore.isAuthenticated" class="flex items-center space-x-4">
            <span class="text-sm text-gray-600">
              欢迎，{{ authStore.user?.username }}
            </span>
            <div class="relative">
              <button 
                @click="showUserMenu = !showUserMenu"
                class="flex items-center space-x-2 text-gray-600 hover:text-gray-900 transition-colors duration-200"
              >
                <User class="w-5 h-5" />
                <ChevronDown class="w-4 h-4" />
              </button>
              
              <!-- 下拉菜单 -->
              <div 
                v-if="showUserMenu"
                class="absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-lg border border-gray-200 py-1 z-50"
              >
                <button
                  @click="showProfile"
                  class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                >
                  个人资料
                </button>
                <button 
                  @click="handleLogout"
                  class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                >
                  退出登录
                </button>
              </div>
            </div>
          </div>
          
          <div v-else class="flex items-center space-x-2">
            <button
              @click="$emit('show-login')"
              class="text-gray-600 hover:text-primary-600 transition-colors duration-200"
            >
              登录
            </button>
            <span class="text-gray-300">|</span>
            <button
              @click="$emit('show-register')"
              class="btn-primary text-sm"
            >
              注册
            </button>
          </div>
        </div>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { MessageCircle, User, ChevronDown } from 'lucide-vue-next'

const emit = defineEmits(['show-login', 'show-register', 'show-profile'])

const router = useRouter()
const authStore = useAuthStore()
const showUserMenu = ref(false)

const showProfile = () => {
  showUserMenu.value = false
  emit('show-profile')
}

const handleLogout = () => {
  authStore.logout()
  showUserMenu.value = false
  router.push('/')
}

// 点击外部关闭菜单
import { onClickOutside } from '@vueuse/core'
const menuRef = ref(null)
onClickOutside(menuRef, () => {
  showUserMenu.value = false
})
</script>

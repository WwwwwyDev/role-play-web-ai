<template>
  <div id="app" class="min-h-screen bg-gray-50">
    <Navbar 
      v-if="showNavbar" 
      @show-login="handleShowLogin"
      @show-register="handleShowRegister"
      @show-profile="handleShowProfile"
    />
    <main :class="showNavbar ? 'pt-16' : ''">
      <router-view />
    </main>

    <!-- 全局登录弹出框 -->
    <LoginModal
      :is-visible="showLoginModal"
      @close="showLoginModal = false"
      @switch-to-register="switchToRegister"
    />

    <!-- 全局注册弹出框 -->
    <RegisterModal
      :is-visible="showRegisterModal"
      @close="showRegisterModal = false"
      @switch-to-login="switchToLogin"
    />

    <!-- 全局个人资料弹出框 -->
    <ProfileModal
      :is-visible="showProfileModal"
      @close="showProfileModal = false"
    />
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import Navbar from '@/components/Navbar.vue'
import LoginModal from '@/components/LoginModal.vue'
import RegisterModal from '@/components/RegisterModal.vue'
import ProfileModal from '@/components/ProfileModal.vue'

const route = useRoute()
const authStore = useAuthStore()

// 弹出框状态
const showLoginModal = ref(false)
const showRegisterModal = ref(false)
const showProfileModal = ref(false)

const showNavbar = computed(() => {
  // 只在非聊天页面显示导航栏
  return route.name !== 'Chat'
})

const handleShowLogin = () => {
  showLoginModal.value = true
}

const handleShowRegister = () => {
  showRegisterModal.value = true
}

const handleShowProfile = () => {
  showProfileModal.value = true
}

const switchToRegister = () => {
  showLoginModal.value = false
  showRegisterModal.value = true
}

const switchToLogin = () => {
  showRegisterModal.value = false
  showLoginModal.value = true
}
</script>

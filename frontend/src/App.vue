<template>
  <div id="app" v-cloak class="min-h-screen bg-gray-50">
    <!-- 全局加载指示器 -->
    <transition name="loader" appear>
      <div v-if="!authStore.isInitialized" class="fixed inset-0 bg-white/90 backdrop-blur-sm flex items-center justify-center z-50">
        <div class="text-center">
          <div class="w-16 h-16 bg-gradient-to-r from-blue-500 to-indigo-600 rounded-full flex items-center justify-center mx-auto mb-4 animate-pulse">
            <div class="w-8 h-8 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
          </div>
          <p class="text-gray-600 font-medium">正在加载...</p>
        </div>
      </div>
    </transition>

     <transition name="navbar" appear>
       <Navbar 
         v-if="showNavbar && authStore.isInitialized" 
         @show-login="handleShowLogin"
         @show-register="handleShowRegister"
         @show-profile="handleShowProfile"
       />
     </transition>
     <main :class="showNavbar && authStore.isInitialized ? 'pt-16' : ''">
       <router-view v-if="authStore.isInitialized" v-slot="{ Component, route }">
         <component :is="Component" :key="route.path" />
       </router-view>
     </main>

    <!-- 全局登录弹出框 -->
    <transition name="modal" appear>
      <LoginModal
        v-if="showLoginModal"
        :is-visible="showLoginModal"
        @close="showLoginModal = false"
        @switch-to-register="switchToRegister"
      />
    </transition>

    <!-- 全局注册弹出框 -->
    <transition name="modal" appear>
      <RegisterModal
        v-if="showRegisterModal"
        :is-visible="showRegisterModal"
        @close="showRegisterModal = false"
        @switch-to-login="switchToLogin"
      />
    </transition>

    <!-- 全局个人资料弹出框 -->
    <transition name="modal" appear>
      <ProfileModal
        v-if="showProfileModal"
        :is-visible="showProfileModal"
        @close="showProfileModal = false"
      />
    </transition>
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

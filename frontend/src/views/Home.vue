<template>
  <div class="min-h-screen overflow-hidden">
    <!-- Hero Section -->
    <section class="relative bg-gradient-to-br from-primary-600 via-primary-700 to-primary-800 text-white overflow-hidden">
      <!-- 背景动画元素 -->
      <div class="absolute inset-0">
        <div class="floating-shapes">
          <div class="shape shape-1"></div>
          <div class="shape shape-2"></div>
          <div class="shape shape-3"></div>
          <div class="shape shape-4"></div>
          <div class="shape shape-5"></div>
        </div>
      </div>
      
      <div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-24">
        <div class="text-center">
          <div class="hero-content">
            <h1 class="text-5xl md:text-7xl font-bold mb-6 hero-title">
              <span class="gradient-text">AI角色扮演</span>
              <br>
              <span class="text-4xl md:text-6xl">聊天</span>
            </h1>
            <p class="text-xl md:text-2xl mb-8 text-primary-100 hero-subtitle">
              与历史人物、文学角色、科学家等AI角色进行深度对话
            </p>
            <div class="flex flex-col sm:flex-row gap-4 justify-center hero-buttons">
              <button
                @click="handleStartChat"
                class="btn-hero btn-primary-hero"
              >
                <span>开始聊天</span>
                <ArrowRight class="w-5 h-5 ml-2" />
              </button>
              <button
                @click="showRegisterModal = true"
                class="btn-hero btn-secondary-hero"
              >
                <span>免费注册</span>
                <UserPlus class="w-5 h-5 ml-2" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section class="py-24 bg-gradient-to-b from-white to-gray-50 relative">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-20">
          <h2 class="text-4xl md:text-5xl font-bold text-gray-900 mb-6 section-title">
            为什么选择我们？
          </h2>
          <p class="text-xl text-gray-600 max-w-2xl mx-auto">
            体验最真实的AI角色扮演对话，开启跨越时空的智慧交流
          </p>
        </div>

        <div class="grid md:grid-cols-3 gap-8">
          <div class="feature-card group">
            <div class="feature-icon">
              <Users class="w-8 h-8 text-primary-600" />
            </div>
            <h3 class="text-xl font-semibold text-gray-900 mb-3">丰富角色库</h3>
            <p class="text-gray-600 leading-relaxed">
              从哈利·波特到苏格拉底，从爱因斯坦到达芬奇，与各种历史人物和文学角色对话
            </p>
          </div>

          <div class="feature-card group">
            <div class="feature-icon">
              <MessageCircle class="w-8 h-8 text-primary-600" />
            </div>
            <h3 class="text-xl font-semibold text-gray-900 mb-3">智能对话</h3>
            <p class="text-gray-600 leading-relaxed">
              基于先进的大语言模型，每个角色都有独特的性格和知识背景
            </p>
          </div>

          <div class="feature-card group">
            <div class="feature-icon">
              <Mic class="w-8 h-8 text-primary-600" />
            </div>
            <h3 class="text-xl font-semibold text-gray-900 mb-3">语音聊天</h3>
            <p class="text-gray-600 leading-relaxed">
              支持语音输入和输出，让对话更加自然流畅
            </p>
          </div>
        </div>
      </div>
    </section>

    <!-- Characters Preview -->
    <section class="py-24 bg-gradient-to-b from-gray-50 to-white relative">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-20">
          <h2 class="text-4xl md:text-5xl font-bold text-gray-900 mb-6 section-title">
            热门角色
          </h2>
          <p class="text-xl text-gray-600 max-w-2xl mx-auto">
            选择您感兴趣的角色开始对话，体验跨越时空的智慧交流
          </p>
        </div>

        <div class="grid md:grid-cols-2 lg:grid-cols-4 gap-8">
          <div 
            v-for="(character, index) in featuredCharacters" 
            :key="character.id"
            class="character-card group"
            :style="{ animationDelay: `${index * 0.1}s` }"
            @click="startChat(character)"
          >
            <div class="text-center">
              <div class="character-avatar">
                <img 
                  v-if="character.avatar_url" 
                  :src="character.avatar_url" 
                  :alt="character.name"
                  class="w-full h-full object-cover"
                />
                <div v-else class="w-full h-full bg-primary-100 flex items-center justify-center">
                  <User class="w-10 h-10 text-primary-600" />
                </div>
              </div>
              <h3 class="text-lg font-semibold text-gray-900 mb-3">{{ character.name }}</h3>
              <p class="text-sm text-gray-600 mb-4 leading-relaxed">{{ character.description }}</p>
              <span class="character-category">
                {{ character.category }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="py-24 bg-gradient-to-br from-primary-600 via-primary-700 to-primary-800 text-white relative overflow-hidden">
      <!-- 背景装饰 -->
      <div class="absolute inset-0">
        <div class="cta-shapes">
          <div class="cta-shape cta-shape-1"></div>
          <div class="cta-shape cta-shape-2"></div>
          <div class="cta-shape cta-shape-3"></div>
        </div>
      </div>
      
      <div class="relative max-w-4xl mx-auto text-center px-4 sm:px-6 lg:px-8">
        <h2 class="text-4xl md:text-5xl font-bold mb-6 cta-title">
          准备好开始您的AI角色扮演之旅了吗？
        </h2>
        <p class="text-xl text-primary-100 mb-10 max-w-2xl mx-auto">
          立即注册，与您喜爱的角色进行深度对话，开启跨越时空的智慧交流
        </p>
        <button
          @click="showRegisterModal = true"
          class="btn-cta"
        >
          <span>免费注册</span>
          <ArrowRight class="w-5 h-5 ml-2" />
        </button>
      </div>
    </section>

    <!-- 登录弹出框 -->
    <LoginModal
      :is-visible="showLoginModal"
      @close="showLoginModal = false"
      @switch-to-register="switchToRegister"
    />

    <!-- 注册弹出框 -->
    <RegisterModal
      :is-visible="showRegisterModal"
      @close="showRegisterModal = false"
      @switch-to-login="switchToLogin"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useChatStore } from '@/stores/chat'
import { useAuthStore } from '@/stores/auth'
import { Users, MessageCircle, Mic, User, ArrowRight, UserPlus } from 'lucide-vue-next'
import LoginModal from '@/components/LoginModal.vue'
import RegisterModal from '@/components/RegisterModal.vue'

const router = useRouter()
const chatStore = useChatStore()
const authStore = useAuthStore()

// 弹出框状态
const showLoginModal = ref(false)
const showRegisterModal = ref(false)

const featuredCharacters = ref([
  {
    id: 1,
    name: '哈利·波特',
    description: '来自霍格沃茨魔法学校的年轻巫师',
    category: '文学人物',
    avatar_url: '/avatars/harry_potter.svg'
  },
  {
    id: 2,
    name: '苏格拉底',
    description: '古希腊哲学家，以苏格拉底式问答法闻名',
    category: '历史人物',
    avatar_url: '/avatars/socrates.svg'
  },
  {
    id: 3,
    name: '爱因斯坦',
    description: '理论物理学家，相对论的创立者',
    category: '科学家',
    avatar_url: '/avatars/einstein.svg'
  },
  {
    id: 4,
    name: '达芬奇',
    description: '文艺复兴时期的博学者，艺术家、发明家',
    category: '艺术家',
    avatar_url: '/avatars/da_vinci.svg'
  }
])

const startChat = (character) => {
  if (authStore.isAuthenticated) {
    router.push(`/chat?character=${character.id}`)
  } else {
    showRegisterModal.value = true
  }
}

const handleStartChat = () => {
  if (authStore.isAuthenticated) {
    router.push('/chat')
  } else {
    showLoginModal.value = true
  }
}

const switchToRegister = () => {
  showLoginModal.value = false
  showRegisterModal.value = true
}

const switchToLogin = () => {
  showRegisterModal.value = false
  showLoginModal.value = true
}

onMounted(async () => {
  if (authStore.isAuthenticated) {
    await chatStore.fetchCharacters()
    if (chatStore.characters.length > 0) {
      featuredCharacters.value = chatStore.characters.slice(0, 4)
    }
  }
})
</script>

<style scoped>
/* Hero Section Styles */
.hero-content {
  animation: fadeInUp 1s ease-out;
}

.hero-title {
  animation: slideInDown 1s ease-out 0.2s both;
}

.hero-subtitle {
  animation: slideInUp 1s ease-out 0.4s both;
}

.hero-buttons {
  animation: slideInUp 1s ease-out 0.6s both;
}

.gradient-text {
  background: linear-gradient(135deg, #ffffff 0%, #e0e7ff 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

/* Floating Shapes Animation */
.floating-shapes {
  position: absolute;
  width: 100%;
  height: 100%;
  overflow: hidden;
}

.shape {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  animation: float 6s ease-in-out infinite;
}

.shape-1 {
  width: 80px;
  height: 80px;
  top: 20%;
  left: 10%;
  animation-delay: 0s;
}

.shape-2 {
  width: 120px;
  height: 120px;
  top: 60%;
  right: 10%;
  animation-delay: 2s;
}

.shape-3 {
  width: 60px;
  height: 60px;
  top: 30%;
  right: 20%;
  animation-delay: 4s;
}

.shape-4 {
  width: 100px;
  height: 100px;
  bottom: 20%;
  left: 20%;
  animation-delay: 1s;
}

.shape-5 {
  width: 40px;
  height: 40px;
  top: 10%;
  right: 30%;
  animation-delay: 3s;
}

/* Hero Buttons */
.btn-hero {
  @apply inline-flex items-center justify-center px-8 py-4 text-lg font-semibold rounded-xl transition-all duration-300 transform hover:scale-105 focus:outline-none focus:ring-4 focus:ring-opacity-50;
}

.btn-primary-hero {
  @apply bg-white text-primary-600 hover:bg-gray-100 focus:ring-white shadow-lg hover:shadow-xl;
}

.btn-secondary-hero {
  @apply border-2 border-white text-white hover:bg-white hover:text-primary-600 focus:ring-white;
}

/* Section Titles */
.section-title {
  animation: fadeInUp 0.8s ease-out;
}

/* Feature Cards */
.feature-card {
  @apply bg-white p-8 rounded-2xl shadow-lg hover:shadow-2xl transition-all duration-300 transform hover:-translate-y-2 border border-gray-100;
  animation: fadeInUp 0.8s ease-out;
}

.feature-card:nth-child(1) { animation-delay: 0.1s; }
.feature-card:nth-child(2) { animation-delay: 0.2s; }
.feature-card:nth-child(3) { animation-delay: 0.3s; }

.feature-icon {
  @apply w-20 h-20 bg-gradient-to-br from-primary-100 to-primary-200 rounded-2xl flex items-center justify-center mx-auto mb-6 transition-all duration-300 group-hover:scale-110 group-hover:rotate-3;
}

/* Character Cards */
.character-card {
  @apply bg-white p-6 rounded-2xl shadow-lg hover:shadow-2xl transition-all duration-300 transform hover:-translate-y-3 cursor-pointer border border-gray-100;
  animation: fadeInUp 0.8s ease-out both;
}

.character-avatar {
  @apply w-24 h-24 rounded-full mx-auto mb-6 overflow-hidden border-4 border-gray-200 group-hover:border-primary-300 transition-all duration-300 group-hover:scale-110;
}

.character-category {
  @apply inline-block bg-gradient-to-r from-primary-100 to-primary-200 text-primary-800 text-sm px-4 py-2 rounded-full font-medium;
}

/* CTA Section */
.cta-title {
  animation: fadeInUp 0.8s ease-out;
}

.cta-shapes {
  position: absolute;
  width: 100%;
  height: 100%;
  overflow: hidden;
}

.cta-shape {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  animation: float 8s ease-in-out infinite;
}

.cta-shape-1 {
  width: 200px;
  height: 200px;
  top: -100px;
  left: -100px;
  animation-delay: 0s;
}

.cta-shape-2 {
  width: 150px;
  height: 150px;
  bottom: -75px;
  right: -75px;
  animation-delay: 3s;
}

.cta-shape-3 {
  width: 100px;
  height: 100px;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  animation-delay: 6s;
}

.btn-cta {
  @apply inline-flex items-center justify-center px-10 py-5 text-xl font-bold bg-white text-primary-600 rounded-2xl transition-all duration-300 transform hover:scale-105 hover:shadow-2xl focus:outline-none focus:ring-4 focus:ring-white focus:ring-opacity-50;
}

/* Animations */
@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes slideInDown {
  from {
    opacity: 0;
    transform: translateY(-50px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes slideInUp {
  from {
    opacity: 0;
    transform: translateY(50px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px) rotate(0deg);
  }
  50% {
    transform: translateY(-20px) rotate(180deg);
  }
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .hero-title {
    @apply text-4xl;
  }
  
  .section-title {
    @apply text-3xl;
  }
  
  .cta-title {
    @apply text-3xl;
  }
  
  .feature-card {
    @apply p-6;
  }
  
  .character-card {
    @apply p-4;
  }
}
</style>

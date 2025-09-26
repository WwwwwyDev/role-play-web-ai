<template>
  <div
    v-if="isVisible"
    class="fixed inset-0 bg-gradient-to-br from-black/30 via-purple-900/20 to-blue-900/30 backdrop-blur-md flex items-center justify-center z-50 animate-fade-in"
    @click="handleBackdropClick"
  >
    <div
      class="bg-white/90 backdrop-blur-2xl rounded-3xl p-8 max-w-lg w-full mx-4 shadow-2xl border border-white/40 animate-slide-up relative overflow-hidden"
      @click.stop
    >
      <!-- 背景装饰 -->
      <div class="absolute top-0 left-0 w-full h-32 bg-gradient-to-br from-blue-500/10 via-indigo-500/10 to-purple-500/10 rounded-t-3xl"></div>
      <div class="absolute top-4 right-4 w-20 h-20 bg-gradient-to-br from-blue-400/20 to-indigo-400/20 rounded-full blur-xl"></div>
      <div class="absolute bottom-4 left-4 w-16 h-16 bg-gradient-to-br from-purple-400/20 to-pink-400/20 rounded-full blur-lg"></div>
      
      <!-- 关闭按钮 -->
      <button 
        @click="handleClose" 
        class="absolute top-6 right-6 p-2 text-gray-400 hover:text-gray-600 hover:bg-white/80 rounded-xl transition-all duration-300 z-50"
        title="关闭"
      >
        <X class="w-5 h-5" />
      </button>
      
      <!-- 主要内容 -->
      <div class="relative z-10">
        <!-- 标题 -->
        <div class="text-center mb-6">
          <div class="w-16 h-16 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-full flex items-center justify-center mx-auto mb-4 shadow-xl">
            <Mic class="w-8 h-8 text-white" />
          </div>
          <h3 class="text-2xl font-bold bg-gradient-to-r from-blue-600 to-indigo-600 bg-clip-text text-transparent mb-2">
            语音输入
          </h3>
          <p class="text-gray-600 text-sm">
            点击开始录音，清晰地说出您想发送的消息
          </p>
        </div>

        <!-- 录音状态区域 -->
        <div class="text-center mb-8">
          <!-- 录音按钮 -->
          <div class="flex justify-center mb-4">
            <button
              @click="toggleRecording"
              :disabled="isProcessing"
              class="w-24 h-24 rounded-full transition-all duration-300 shadow-2xl transform hover:scale-105 disabled:scale-100 disabled:opacity-50 relative overflow-hidden flex items-center justify-center"
              :class="isRecording 
                ? 'bg-gradient-to-r from-red-500 to-red-600 text-white animate-pulse' 
                : 'bg-gradient-to-r from-blue-500 to-indigo-600 text-white hover:from-blue-600 hover:to-indigo-700'"
            >
              <!-- 录音时的脉冲效果 -->
              <div v-if="isRecording" class="absolute inset-0 bg-red-400 rounded-full animate-ping opacity-30"></div>
              <Mic v-if="!isRecording" class="w-8 h-8 relative z-10" />
              <Square v-else class="w-8 h-8 relative z-10" />
            </button>
          </div>

          <!-- 状态文字 -->
          <div>
            <p class="text-sm font-medium" :class="{
              'text-red-600': isRecording,
              'text-blue-600': !isRecording && !isProcessing,
              'text-gray-500': isProcessing
            }">
              {{ getStatusText() }}
            </p>
          </div>
        </div>

        <!-- 音频质量指示器 -->
        <div v-if="isRecording" class="mb-6">
          <div class="flex items-center justify-center space-x-3 mb-3">
            <span class="text-xs text-gray-600">音频质量</span>
            <div class="flex items-center space-x-1">
              <div class="w-2 h-2 rounded-full" :class="{
                'bg-green-500': audioQuality.quality === 'excellent',
                'bg-yellow-500': audioQuality.quality === 'good',
                'bg-orange-500': audioQuality.quality === 'fair',
                'bg-red-500': audioQuality.quality === 'poor'
              }"></div>
              <span class="text-xs font-medium" :class="{
                'text-green-600': audioQuality.quality === 'excellent',
                'text-yellow-600': audioQuality.quality === 'good',
                'text-orange-600': audioQuality.quality === 'fair',
                'text-red-600': audioQuality.quality === 'poor'
              }">
                {{ getQualityText(audioQuality.quality) }}
              </span>
            </div>
          </div>
          
          <!-- 动态波形动画 -->
          <div class="flex justify-center space-x-1">
            <div 
              v-for="i in 7" 
              :key="i"
              class="w-1 rounded-full animate-pulse"
              :style="{
                height: `${3 + (audioQuality.level * 4)}px`,
                animationDelay: `${i * 0.1}s`,
                backgroundColor: audioQuality.quality === 'excellent' ? '#10b981' : 
                               audioQuality.quality === 'good' ? '#f59e0b' :
                               audioQuality.quality === 'fair' ? '#f97316' : '#ef4444'
              }"
            ></div>
          </div>
        </div>

        <!-- 识别结果显示 -->
        <div v-if="recognizedText" class="mb-6">
          <div class="p-4 bg-gradient-to-r from-green-50 to-emerald-50 border border-green-200 rounded-xl">
            <div class="flex items-center space-x-2 mb-2">
              <div class="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
              <span class="text-sm font-medium text-green-800">识别成功</span>
            </div>
            <p class="text-sm text-gray-700 bg-white/70 p-3 rounded-lg mb-2">
              {{ recognizedText }}
            </p>
            <div class="flex items-center justify-center space-x-2 text-xs text-green-600">
              <div class="w-1.5 h-1.5 bg-green-500 rounded-full animate-pulse"></div>
              <span>正在发送消息...</span>
            </div>
          </div>
        </div>

        <!-- 错误提示 -->
        <div v-if="error" class="mb-6">
          <div class="p-4 bg-gradient-to-r from-red-50 to-orange-50 border border-red-200 rounded-xl">
            <div class="flex items-center space-x-2 mb-2">
              <AlertCircle class="w-4 h-4 text-red-500" />
              <span class="text-sm font-medium text-red-800">识别错误</span>
            </div>
            <p class="text-sm text-red-700 mb-3">{{ error }}</p>
            <div class="flex space-x-2">
              <button
                v-if="!error.includes('权限被拒绝')"
                @click="retryRecording"
                class="px-3 py-1.5 text-xs bg-blue-500 text-white rounded-full hover:bg-blue-600 transition-all duration-200"
              >
                重试
              </button>
              <button
                @click="error = ''"
                class="px-3 py-1.5 text-xs bg-red-500 text-white rounded-full hover:bg-red-600 transition-all duration-200"
              >
                关闭
              </button>
            </div>
          </div>
        </div>


        <!-- 使用提示 -->
        <div class="mt-6 text-center">
          <div class="flex items-center justify-center space-x-4 text-xs text-gray-500">
            <div class="flex items-center space-x-1">
              <div class="w-1.5 h-1.5 bg-blue-500 rounded-full"></div>
              <span>清晰发音</span>
            </div>
            <div class="flex items-center space-x-1">
              <div class="w-1.5 h-1.5 bg-green-500 rounded-full"></div>
              <span>安静环境</span>
            </div>
            <div class="flex items-center space-x-1">
              <div class="w-1.5 h-1.5 bg-purple-500 rounded-full"></div>
              <span>网络连接</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { Mic, Square, X, AlertCircle } from 'lucide-vue-next'
import voiceService from '@/services/voice'

const props = defineProps({
  isVisible: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close', 'message'])

const isRecording = ref(false)
const isProcessing = ref(false)
const recognizedText = ref('')
const error = ref('')
const audioQuality = ref({ level: 0, quality: 'unknown' })

const getStatusText = () => {
  if (isProcessing.value) return '处理中...'
  if (isRecording.value) return '正在听取您的语音...'
  return '点击开始录音'
}

const getQualityText = (quality) => {
  const qualityMap = {
    'excellent': '优秀',
    'good': '良好', 
    'fair': '一般',
    'poor': '较差',
    'unknown': '未知'
  }
  return qualityMap[quality] || '未知'
}

const toggleRecording = async () => {
  if (isRecording.value) {
    stopRecording()
  } else {
    await startRecording()
  }
}

const startRecording = async () => {
  try {
    error.value = ''
    recognizedText.value = ''
    isRecording.value = true
    isProcessing.value = false
    
    const text = await voiceService.startRecognition({
      onInterimResult: (transcript) => {
        recognizedText.value = transcript
      },
      onFinalResult: (transcript) => {
        recognizedText.value = transcript
        isRecording.value = false
        isProcessing.value = false
        
        // 自动发送消息
        if (transcript.trim()) {
          setTimeout(() => {
            sendMessage()
          }, 500) // 延迟500ms让用户看到识别结果
        }
      },
      onQualityUpdate: (quality) => {
        audioQuality.value = quality
      },
      onError: (errorMsg, shouldRetry) => {
        if (shouldRetry) {
          error.value = `${errorMsg} (正在重试...)`
        } else {
          error.value = errorMsg
          isRecording.value = false
          isProcessing.value = false
        }
      }
    })
    
  } catch (err) {
    error.value = getErrorMessage(err.message)
    isRecording.value = false
    isProcessing.value = false
  }
}

const stopRecording = () => {
  voiceService.stopRecognition()
  isRecording.value = false
  isProcessing.value = false
}

const retryRecording = async () => {
  error.value = ''
  await startRecording()
}

const sendMessage = () => {
  if (recognizedText.value.trim()) {
    emit('message', recognizedText.value.trim())
    // 延迟关闭，让用户看到发送状态
    setTimeout(() => {
      handleClose()
    }, 1000)
  }
}

const editMessage = () => {
  // 这里可以打开一个编辑对话框
  // 暂时直接发送
  sendMessage()
}

const handleClose = () => {
  if (isRecording.value) {
    stopRecording()
  }
  emit('close')
}

const handleBackdropClick = () => {
  handleClose()
}

const getErrorMessage = (error) => {
  const errorMap = {
    'no-speech': '没有检测到语音，请重试',
    'audio-capture': '无法访问麦克风，请检查权限',
    'not-allowed': '麦克风权限被拒绝，请在浏览器设置中允许',
    'network': '网络错误，请检查网络连接',
    'aborted': '语音识别被中断',
    'language-not-supported': '不支持当前语言设置'
  }
  return errorMap[error] || `语音识别错误：${error}`
}

// 监听气泡显示状态
watch(() => props.isVisible, (newValue) => {
  if (!newValue) {
    // 气泡关闭时重置状态
    isRecording.value = false
    isProcessing.value = false
    recognizedText.value = ''
    error.value = ''
    audioQuality.value = { level: 0, quality: 'unknown' }
  }
})
</script>

<style scoped>
@keyframes fade-in {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

@keyframes slide-up {
  from {
    opacity: 0;
    transform: translateY(30px) scale(0.95);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px) rotate(0deg);
  }
  50% {
    transform: translateY(-10px) rotate(2deg);
  }
}

@keyframes float-delayed {
  0%, 100% {
    transform: translateY(0px) rotate(0deg);
  }
  50% {
    transform: translateY(-8px) rotate(-1deg);
  }
}

@keyframes pulse-slow {
  0%, 100% {
    opacity: 0.3;
    transform: scale(1);
  }
  50% {
    opacity: 0.6;
    transform: scale(1.1);
  }
}

@keyframes wave {
  0%, 100% {
    transform: scaleY(1);
  }
  50% {
    transform: scaleY(1.5);
  }
}

@keyframes glow {
  0%, 100% {
    box-shadow: 0 0 20px rgba(59, 130, 246, 0.3);
  }
  50% {
    box-shadow: 0 0 30px rgba(59, 130, 246, 0.6);
  }
}

.animate-fade-in {
  animation: fade-in 0.4s ease-out;
}

.animate-slide-up {
  animation: slide-up 0.5s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.animate-float {
  animation: float 6s ease-in-out infinite;
}

.animate-float-delayed {
  animation: float-delayed 8s ease-in-out infinite;
}

.animate-pulse-slow {
  animation: pulse-slow 4s ease-in-out infinite;
}

.animate-wave {
  animation: wave 1.5s ease-in-out infinite;
}

.animate-glow {
  animation: glow 2s ease-in-out infinite;
}

/* 录音按钮脉冲动画 */
@keyframes pulse-recording {
  0%, 100% {
    transform: scale(1);
    box-shadow: 0 0 0 0 rgba(239, 68, 68, 0.7);
  }
  50% {
    transform: scale(1.05);
    box-shadow: 0 0 0 15px rgba(239, 68, 68, 0);
  }
}

.animate-pulse-recording {
  animation: pulse-recording 2s infinite;
}

/* 渐变背景动画 */
@keyframes gradient-shift {
  0% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
  100% {
    background-position: 0% 50%;
  }
}

.animate-gradient {
  background-size: 200% 200%;
  animation: gradient-shift 3s ease infinite;
}
</style>

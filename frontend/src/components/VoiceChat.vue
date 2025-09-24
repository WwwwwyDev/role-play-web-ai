<template>
  <div class="voice-chat">
    <!-- 语音控制按钮 -->
    <div class="flex items-center justify-center space-x-4">
      <button
        @click="toggleRecording"
        :disabled="isSpeaking"
        class="p-4 rounded-full transition-all duration-300 shadow-lg hover:shadow-xl transform hover:scale-105"
        :class="isRecording 
          ? 'bg-gradient-to-r from-red-500 to-red-600 text-white animate-pulse' 
          : 'bg-gradient-to-r from-blue-500 to-indigo-600 text-white hover:from-blue-600 hover:to-indigo-700'"
        :title="isRecording ? '停止录音' : '开始语音输入'"
      >
        <Mic v-if="!isRecording" class="w-6 h-6" />
        <Square v-else class="w-6 h-6" />
      </button>

      <button
        @click="toggleSpeaking"
        :disabled="isRecording"
        class="p-4 rounded-full transition-all duration-300 shadow-lg hover:shadow-xl transform hover:scale-105"
        :class="isSpeaking 
          ? 'bg-gradient-to-r from-green-500 to-green-600 text-white' 
          : 'bg-gradient-to-r from-gray-500 to-gray-600 text-white hover:from-gray-600 hover:to-gray-700'"
        :title="isSpeaking ? '停止播放' : '播放AI回复'"
      >
        <Volume2 v-if="!isSpeaking" class="w-6 h-6" />
        <VolumeX v-else class="w-6 h-6" />
      </button>
    </div>

    <!-- 录音状态指示器 -->
    <div v-if="isRecording" class="mt-4 text-center">
      <div class="flex items-center justify-center space-x-3 text-red-600">
        <div class="w-3 h-3 bg-red-600 rounded-full animate-pulse"></div>
        <span class="text-sm font-medium">正在听取您的语音...</span>
      </div>
      <div class="mt-2 text-xs text-gray-500">
        请清晰地说话，识别结果将自动填入输入框
      </div>
    </div>

    <!-- 语音识别结果 -->
    <div v-if="recognizedText" class="mt-4 p-3 bg-gradient-to-r from-blue-50 to-indigo-50 border border-blue-200 rounded-xl">
      <div class="flex items-center space-x-2 mb-2">
        <div class="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
        <span class="text-sm font-medium text-blue-800">识别结果：</span>
      </div>
      <p class="text-sm text-gray-700 bg-white/70 p-2 rounded-lg">
        {{ recognizedText }}
      </p>
      <div class="mt-2 text-xs text-gray-500">
        已自动填入输入框，您可以编辑后发送
      </div>
    </div>

    <!-- 错误提示 -->
    <div v-if="error" class="mt-4 p-3 bg-red-50 border border-red-200 rounded-xl">
      <div class="flex items-center space-x-2">
        <AlertCircle class="w-4 h-4 text-red-500" />
        <p class="text-sm text-red-800 font-medium">{{ error }}</p>
      </div>
      <div class="mt-1 text-xs text-red-600">
        请检查麦克风权限或重试
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { Mic, Square, Volume2, VolumeX, AlertCircle } from 'lucide-vue-next'
import voiceService from '@/services/voice'

const props = defineProps({
  onMessage: {
    type: Function,
    required: true
  },
  lastMessage: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['message', 'error'])

const isRecording = ref(false)
const isSpeaking = ref(false)
const recognizedText = ref('')
const interimText = ref('')
const error = ref('')
const isSupported = ref({ recognition: false, synthesis: false })

const toggleRecording = async () => {
  if (isRecording.value) {
    stopRecording()
  } else {
    await startRecording()
  }
}

const toggleSpeaking = () => {
  if (isSpeaking.value) {
    stopSpeaking()
  } else if (props.lastMessage) {
    speakMessage(props.lastMessage)
  }
}

const startRecording = async () => {
  try {
    error.value = ''
    recognizedText.value = ''
    interimText.value = ''
    
    const text = await voiceService.startRecognition({
      onInterimResult: (transcript) => {
        interimText.value = transcript
        // 实时更新输入框内容
        emit('message', transcript)
      },
      onFinalResult: (transcript) => {
        recognizedText.value = transcript
        interimText.value = ''
        // 最终结果写入输入框
        emit('message', transcript)
      },
      onError: (errorMsg) => {
        error.value = getErrorMessage(errorMsg)
        emit('error', errorMsg)
      }
    })
    
  } catch (err) {
    error.value = getErrorMessage(err.message)
    emit('error', err.message)
  }
}

const stopRecording = () => {
  voiceService.stopRecognition()
  isRecording.value = false
  interimText.value = ''
}

const speakMessage = async (text) => {
  try {
    error.value = ''
    await voiceService.speak(text, {
      lang: 'zh-CN',
      rate: 0.9,
      pitch: 1,
      volume: 0.8
    })
  } catch (err) {
    error.value = getErrorMessage(err.message)
    emit('error', err.message)
  }
}

const stopSpeaking = () => {
  voiceService.stopSpeaking()
  isSpeaking.value = false
}

// 获取友好的错误信息
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

// 监听语音服务状态
const updateStatus = () => {
  isRecording.value = voiceService.isRecording
  isSpeaking.value = voiceService.isSpeaking
}

onMounted(() => {
  isSupported.value = voiceService.isSupported()
  
  if (!isSupported.value.recognition) {
    error.value = '您的浏览器不支持语音识别功能，请使用Chrome或Edge浏览器'
  }
  
  if (!isSupported.value.synthesis) {
    console.warn('浏览器不支持语音合成功能')
  }

  // 定期更新状态
  const statusInterval = setInterval(updateStatus, 100)
  
  onUnmounted(() => {
    clearInterval(statusInterval)
    stopRecording()
    stopSpeaking()
  })
})
</script>

<style scoped>
.voice-chat {
  @apply p-6 bg-white/95 backdrop-blur-sm rounded-2xl border border-gray-200/50 shadow-lg;
}

/* 语音按钮动画 */
@keyframes pulse-recording {
  0%, 100% {
    transform: scale(1);
    box-shadow: 0 0 0 0 rgba(239, 68, 68, 0.7);
  }
  50% {
    transform: scale(1.05);
    box-shadow: 0 0 0 10px rgba(239, 68, 68, 0);
  }
}

.animate-pulse-recording {
  animation: pulse-recording 2s infinite;
}

/* 识别结果动画 */
@keyframes slideInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-slide-in-up {
  animation: slideInUp 0.3s ease-out;
}
</style>

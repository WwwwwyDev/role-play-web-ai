class VoiceService {
  constructor() {
    this.recognition = null
    this.synthesis = window.speechSynthesis
    this.isRecording = false
    this.isSpeaking = false
    this.onInterimResult = null
    this.onFinalResult = null
    this.onError = null
    this.onQualityUpdate = null
    this.retryCount = 0
    this.maxRetries = 3
    this.audioContext = null
    this.analyser = null
    this.microphone = null
    this.initSpeechRecognition()
  }

  initSpeechRecognition() {
    if ('webkitSpeechRecognition' in window || 'SpeechRecognition' in window) {
      const SpeechRecognition = window.SpeechRecognition || window.webkitSpeechRecognition
      this.recognition = new SpeechRecognition()
      
      this.recognition.continuous = false
      this.recognition.interimResults = true // 启用临时结果
      this.recognition.lang = 'zh-CN' // 中文识别
      this.recognition.maxAlternatives = 3 // 增加备选结果
      this.recognition.serviceURI = 'wss://www.google.com/speech-api/full-duplex/v1/up' // 使用Google语音API
      
      this.recognition.onstart = () => {
        this.isRecording = true
        console.log('语音识别开始')
      }
      
      this.recognition.onend = () => {
        this.isRecording = false
        console.log('语音识别结束')
      }
      
      this.recognition.onresult = (event) => {
        let interimTranscript = ''
        let finalTranscript = ''
        
        for (let i = event.resultIndex; i < event.results.length; i++) {
          const result = event.results[i]
          const transcript = result[0].transcript
          const confidence = result[0].confidence || 0
          
          // 只处理置信度较高的结果
          if (confidence > 0.3) {
            if (result.isFinal) {
              finalTranscript += transcript
            } else {
              interimTranscript += transcript
            }
          }
        }
        
        // 清理和格式化文本
        const cleanInterim = this.cleanTranscript(interimTranscript)
        const cleanFinal = this.cleanTranscript(finalTranscript)
        
        // 触发临时结果回调
        if (cleanInterim && this.onInterimResult) {
          this.onInterimResult(cleanInterim)
        }
        
        // 触发最终结果回调
        if (cleanFinal && this.onFinalResult) {
          this.onFinalResult(cleanFinal)
        }
      }
      
      this.recognition.onerror = (event) => {
        console.error('语音识别错误:', event.error)
        this.isRecording = false
        
        // 处理不同类型的错误
        let errorMessage = event.error
        let shouldRetry = false
        
        switch (event.error) {
          case 'network':
            errorMessage = '网络连接错误，请检查网络连接后重试'
            shouldRetry = true
            break
          case 'no-speech':
            errorMessage = '没有检测到语音，请大声一点说话'
            shouldRetry = true
            break
          case 'audio-capture':
            errorMessage = '无法访问麦克风，请检查麦克风连接和权限设置'
            break
          case 'not-allowed':
            errorMessage = '麦克风权限被拒绝，请点击地址栏的麦克风图标允许访问'
            break
          case 'aborted':
            errorMessage = '语音识别被中断，请重新开始'
            shouldRetry = true
            break
          case 'language-not-supported':
            errorMessage = '当前浏览器不支持中文语音识别，建议使用Chrome或Edge浏览器'
            break
          case 'service-not-allowed':
            errorMessage = '语音服务被拒绝，请检查网络连接'
            shouldRetry = true
            break
          default:
            errorMessage = `语音识别错误：${event.error}，请重试`
            shouldRetry = true
        }
        
        if (this.onError) {
          this.onError(errorMessage, shouldRetry)
        }
      }
    }
  }

  // 清理和格式化识别文本
  cleanTranscript(text) {
    if (!text) return ''
    
    // 移除多余空格
    let cleaned = text.trim()
    
    // 处理常见的中文识别错误
    const corrections = {
      '你好': '你好',
      '谢谢': '谢谢',
      '再见': '再见',
      '请问': '请问',
      '可以': '可以',
      '什么': '什么',
      '怎么': '怎么',
      '为什么': '为什么',
      '哪里': '哪里',
      '什么时候': '什么时候'
    }
    
    // 应用常见错误修正
    Object.keys(corrections).forEach(wrong => {
      const correct = corrections[wrong]
      cleaned = cleaned.replace(new RegExp(wrong, 'g'), correct)
    })
    
    // 移除标点符号错误
    cleaned = cleaned.replace(/[，。！？；：""''（）【】]/g, '')
    
    return cleaned
  }

  // 获取智能重试延迟
  getRetryDelay(errorType, retryCount) {
    const baseDelays = {
      'network': 2000,      // 网络错误延迟更长
      'no-speech': 1000,    // 无语音错误延迟较短
      'aborted': 500,       // 中断错误延迟最短
      'service-not-allowed': 3000 // 服务错误延迟最长
    }
    
    const baseDelay = baseDelays[errorType] || 1000
    // 指数退避：每次重试延迟翻倍
    return baseDelay * Math.pow(2, retryCount - 1)
  }

  // 初始化音频质量检测
  async initAudioQualityDetection() {
    try {
      this.audioContext = new (window.AudioContext || window.webkitAudioContext)()
      this.analyser = this.audioContext.createAnalyser()
      this.analyser.fftSize = 256
      
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
      this.microphone = this.audioContext.createMediaStreamSource(stream)
      this.microphone.connect(this.analyser)
      
      return true
    } catch (error) {
      console.warn('音频质量检测初始化失败:', error)
      return false
    }
  }

  // 获取音频质量数据
  getAudioQuality() {
    if (!this.analyser) return { level: 0, quality: 'unknown' }
    
    const dataArray = new Uint8Array(this.analyser.frequencyBinCount)
    this.analyser.getByteFrequencyData(dataArray)
    
    // 计算平均音量
    const average = dataArray.reduce((sum, value) => sum + value, 0) / dataArray.length
    const level = Math.min(average / 128, 1) // 归一化到0-1
    
    let quality = 'poor'
    if (level > 0.7) quality = 'excellent'
    else if (level > 0.4) quality = 'good'
    else if (level > 0.2) quality = 'fair'
    
    return { level, quality }
  }

  // 开始音频质量监控
  startQualityMonitoring() {
    if (!this.analyser) return
    
    const monitor = () => {
      if (this.isRecording) {
        const quality = this.getAudioQuality()
        if (this.onQualityUpdate) {
          this.onQualityUpdate(quality)
        }
        requestAnimationFrame(monitor)
      }
    }
    monitor()
  }

  // 开始语音识别
  startRecognition(options = {}) {
    return new Promise((resolve, reject) => {
      if (!this.recognition) {
        reject(new Error('浏览器不支持语音识别'))
        return
      }

      // 重置重试计数
      this.retryCount = 0

      // 设置回调函数
      this.onInterimResult = options.onInterimResult || null
      this.onFinalResult = options.onFinalResult || null
      this.onQualityUpdate = options.onQualityUpdate || null
      this.onError = (errorMessage, shouldRetry) => {
        if (shouldRetry && this.retryCount < this.maxRetries) {
          this.retryCount++
          console.log(`语音识别重试 ${this.retryCount}/${this.maxRetries}`)
          
          // 智能延迟重试：根据错误类型和重试次数调整延迟
          const delay = this.getRetryDelay(event.error, this.retryCount)
          
          setTimeout(() => {
            try {
              this.recognition.start()
            } catch (err) {
              if (options.onError) {
                options.onError(errorMessage)
              }
              reject(new Error(errorMessage))
            }
          }, delay)
        } else {
          if (options.onError) {
            options.onError(errorMessage)
          }
          reject(new Error(errorMessage))
        }
      }

      // 设置最终结果处理
      const originalOnFinalResult = this.onFinalResult
      this.onFinalResult = (transcript) => {
        if (originalOnFinalResult) {
          originalOnFinalResult(transcript)
        }
        this.retryCount = 0 // 重置重试计数
        resolve(transcript)
      }

      try {
        // 初始化音频质量检测
        this.initAudioQualityDetection().then(success => {
          if (success) {
            this.startQualityMonitoring()
          }
        })
        
        this.recognition.start()
      } catch (error) {
        reject(error)
      }
    })
  }

  // 停止语音识别
  stopRecognition() {
    if (this.recognition && this.isRecording) {
      this.recognition.stop()
    }
  }

  // 语音合成
  speak(text, options = {}) {
    return new Promise((resolve, reject) => {
      if (!this.synthesis) {
        reject(new Error('浏览器不支持语音合成'))
        return
      }

      // 停止当前播放
      this.synthesis.cancel()

      const utterance = new SpeechSynthesisUtterance(text)
      
      // 设置语音参数
      utterance.lang = options.lang || 'zh-CN'
      utterance.rate = options.rate || 1
      utterance.pitch = options.pitch || 1
      utterance.volume = options.volume || 1

      // 选择最佳中文语音
      const voices = this.synthesis.getVoices()
      const chineseVoices = voices.filter(voice => 
        voice.lang.includes('zh') || voice.lang.includes('CN')
      )
      
      // 优先选择中文语音，按质量排序
      const preferredVoices = [
        'Microsoft Huihui - Chinese (Simplified, PRC)',
        'Microsoft Yaoyao - Chinese (Simplified, PRC)', 
        'Microsoft Kangkang - Chinese (Simplified, PRC)',
        'Google 普通话（中国大陆）',
        'Google 中文（中国）'
      ]
      
      let selectedVoice = null
      for (const preferred of preferredVoices) {
        selectedVoice = chineseVoices.find(voice => voice.name === preferred)
        if (selectedVoice) break
      }
      
      // 如果没有找到首选语音，选择第一个中文语音
      if (!selectedVoice && chineseVoices.length > 0) {
        selectedVoice = chineseVoices[0]
      }
      
      if (selectedVoice) {
        utterance.voice = selectedVoice
      }

      utterance.onstart = () => {
        this.isSpeaking = true
      }

      utterance.onend = () => {
        this.isSpeaking = false
        resolve()
      }

      utterance.onerror = (event) => {
        this.isSpeaking = false
        reject(new Error(event.error))
      }

      this.synthesis.speak(utterance)
    })
  }

  // 停止语音合成
  stopSpeaking() {
    if (this.synthesis && this.isSpeaking) {
      this.synthesis.cancel()
      this.isSpeaking = false
    }
  }

  // 检查浏览器支持
  isSupported() {
    return {
      recognition: !!(window.SpeechRecognition || window.webkitSpeechRecognition),
      synthesis: !!window.speechSynthesis,
      audioContext: !!(window.AudioContext || window.webkitAudioContext),
      mediaDevices: !!navigator.mediaDevices
    }
  }

  // 获取可用的语音识别引擎
  getAvailableEngines() {
    const engines = []
    
    if (window.SpeechRecognition || window.webkitSpeechRecognition) {
      engines.push({
        name: 'Web Speech API',
        type: 'native',
        supported: true,
        description: '浏览器原生语音识别'
      })
    }
    
    // 检查是否支持其他语音识别服务
    if (navigator.mediaDevices && navigator.mediaDevices.getUserMedia) {
      engines.push({
        name: 'MediaRecorder API',
        type: 'fallback',
        supported: true,
        description: '媒体录制降级方案'
      })
    }
    
    return engines
  }

  // 降级语音识别方案
  async startFallbackRecognition(options = {}) {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
      const mediaRecorder = new MediaRecorder(stream)
      const chunks = []
      
      mediaRecorder.ondataavailable = (event) => {
        chunks.push(event.data)
      }
      
      mediaRecorder.onstop = async () => {
        const blob = new Blob(chunks, { type: 'audio/wav' })
        // 这里可以调用后端API进行语音识别
        // 或者使用其他第三方语音识别服务
        console.log('降级方案：录音完成，需要后端处理')
      }
      
      mediaRecorder.start()
      
      // 5秒后自动停止
      setTimeout(() => {
        mediaRecorder.stop()
        stream.getTracks().forEach(track => track.stop())
      }, 5000)
      
    } catch (error) {
      throw new Error('降级语音识别失败：' + error.message)
    }
  }

  // 获取可用语音列表
  getVoices() {
    return new Promise((resolve) => {
      if (this.synthesis.getVoices().length > 0) {
        resolve(this.synthesis.getVoices())
      } else {
        this.synthesis.onvoiceschanged = () => {
          resolve(this.synthesis.getVoices())
        }
      }
    })
  }
}

export default new VoiceService()

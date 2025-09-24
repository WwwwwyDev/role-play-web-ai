class VoiceService {
  constructor() {
    this.recognition = null
    this.synthesis = window.speechSynthesis
    this.isRecording = false
    this.isSpeaking = false
    this.onInterimResult = null
    this.onFinalResult = null
    this.onError = null
    this.retryCount = 0
    this.maxRetries = 3
    this.initSpeechRecognition()
  }

  initSpeechRecognition() {
    if ('webkitSpeechRecognition' in window || 'SpeechRecognition' in window) {
      const SpeechRecognition = window.SpeechRecognition || window.webkitSpeechRecognition
      this.recognition = new SpeechRecognition()
      
      this.recognition.continuous = false
      this.recognition.interimResults = true // 启用临时结果
      this.recognition.lang = 'zh-CN'
      this.recognition.maxAlternatives = 1
      
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
          const transcript = event.results[i][0].transcript
          if (event.results[i].isFinal) {
            finalTranscript += transcript
          } else {
            interimTranscript += transcript
          }
        }
        
        // 触发临时结果回调
        if (interimTranscript && this.onInterimResult) {
          this.onInterimResult(interimTranscript)
        }
        
        // 触发最终结果回调
        if (finalTranscript && this.onFinalResult) {
          this.onFinalResult(finalTranscript)
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
            errorMessage = '网络连接错误，请检查网络连接'
            shouldRetry = true
            break
          case 'no-speech':
            errorMessage = '没有检测到语音，请重试'
            shouldRetry = true
            break
          case 'audio-capture':
            errorMessage = '无法访问麦克风，请检查权限'
            break
          case 'not-allowed':
            errorMessage = '麦克风权限被拒绝，请在浏览器设置中允许'
            break
          case 'aborted':
            errorMessage = '语音识别被中断'
            shouldRetry = true
            break
          case 'language-not-supported':
            errorMessage = '不支持当前语言设置'
            break
          default:
            errorMessage = `语音识别错误：${event.error}`
            shouldRetry = true
        }
        
        if (this.onError) {
          this.onError(errorMessage, shouldRetry)
        }
      }
    }
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
      this.onError = (errorMessage, shouldRetry) => {
        if (shouldRetry && this.retryCount < this.maxRetries) {
          this.retryCount++
          console.log(`语音识别重试 ${this.retryCount}/${this.maxRetries}`)
          
          // 延迟重试
          setTimeout(() => {
            try {
              this.recognition.start()
            } catch (err) {
              if (options.onError) {
                options.onError(errorMessage)
              }
              reject(new Error(errorMessage))
            }
          }, 1000)
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

      // 选择中文语音
      const voices = this.synthesis.getVoices()
      const chineseVoice = voices.find(voice => 
        voice.lang.includes('zh') || voice.lang.includes('CN')
      )
      if (chineseVoice) {
        utterance.voice = chineseVoice
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
      synthesis: !!window.speechSynthesis
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

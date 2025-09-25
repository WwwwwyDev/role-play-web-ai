import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/services/api'

export const useChatStore = defineStore('chat', () => {
  const characters = ref([])
  const conversations = ref([])
  const currentConversation = ref(null)
  const messages = ref([])
  const isLoading = ref(false)

  const fetchCharacters = async () => {
    try {
      const response = await api.get('/characters')
      characters.value = response.data.characters
      return { success: true }
    } catch (error) {
      return { 
        success: false, 
        error: error.response?.data?.error || '获取角色列表失败' 
      }
    }
  }

  const searchCharacters = async (query) => {
    try {
      const response = await api.get(`/characters/search?q=${encodeURIComponent(query)}`)
      return { 
        success: true, 
        characters: response.data.characters 
      }
    } catch (error) {
      return { 
        success: false, 
        error: error.response?.data?.error || '搜索角色失败' 
      }
    }
  }

  const fetchConversations = async () => {
    try {
      const response = await api.get('/conversations')
      conversations.value = response.data.conversations || []
      return { success: true }
    } catch (error) {
      console.error('获取对话列表失败:', error)
      conversations.value = []
      return { 
        success: false, 
        error: error.response?.data?.error || '获取对话列表失败' 
      }
    }
  }

  const createConversation = async (characterId) => {
    try {
      const response = await api.post('/conversations', { character_id: characterId })
      const conversation = response.data.conversation
      conversations.value.unshift(conversation)
      return { success: true, conversation }
    } catch (error) {
      return { 
        success: false, 
        error: error.response?.data?.error || '创建对话失败' 
      }
    }
  }

  const fetchConversation = async (conversationId) => {
    try {
      console.log('Fetching conversation:', conversationId)
      const response = await api.get(`/conversations/${conversationId}`)
      currentConversation.value = response.data.conversation
      messages.value = response.data.messages || []
      console.log('Conversation loaded:', response.data.conversation?.title, 'Messages:', messages.value.length)
      return { success: true }
    } catch (error) {
      console.error('Failed to fetch conversation:', error)
      currentConversation.value = null
      messages.value = []
      return { 
        success: false, 
        error: error.response?.data?.error || '获取对话失败' 
      }
    }
  }

  const sendMessage = async (conversationId, content, audioUrl = null) => {
    isLoading.value = true
    try {
      console.log('Sending message:', { conversationId, content, audioUrl })
      const response = await api.post(`/conversations/${conversationId}/messages`, {
        content,
        audio_url: audioUrl
      })
      
      console.log('Message response:', response.data)
      
      // 更新用户消息的ID（从临时ID更新为真实ID）
      const userMessageIndex = messages.value.findIndex(msg => 
        msg.role === 'user' && msg.content === content && msg.id > 1000000000000 // 临时ID通常是时间戳
      )
      if (userMessageIndex !== -1 && response.data.user_message) {
        messages.value[userMessageIndex] = response.data.user_message
        console.log('Updated user message:', response.data.user_message)
      }
      
      // 添加AI消息
      if (response.data.ai_message) {
        messages.value.push(response.data.ai_message)
        console.log('Added AI message:', response.data.ai_message)
      } else {
        console.warn('No AI message in response')
      }
      
      return { success: true }
    } catch (error) {
      console.error('Send message error:', error)
      return { 
        success: false, 
        error: error.response?.data?.error || '发送消息失败' 
      }
    } finally {
      isLoading.value = false
    }
  }

  const sendMessageStream = async (conversationId, content, audioUrl = null, onStreamData = null) => {
    isLoading.value = true
    let hasReceivedFirstData = false
    
    try {
      console.log('Sending streaming message:', { conversationId, content, audioUrl })
      
      const token = localStorage.getItem('token')
      const response = await fetch(`/api/v1/conversations/${conversationId}/messages/stream`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({
          content,
          audio_url: audioUrl
        })
      })

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const reader = response.body.getReader()
      const decoder = new TextDecoder()
      let aiMessageIndex = -1

      while (true) {
        const { done, value } = await reader.read()
        if (done) break

        const chunk = decoder.decode(value)
        const lines = chunk.split('\n')

        for (const line of lines) {
          if (line.startsWith('data: ')) {
            const data = line.slice(6)
            if (data === '[DONE]') {
              console.log('Stream completed')
              break
            }

            try {
              const parsedData = JSON.parse(data)
              
              // 处理用户消息
              if (parsedData.role === 'user') {
                const userMessageIndex = messages.value.findIndex(msg => 
                  msg.role === 'user' && msg.content === content && msg.id > 1000000000000
                )
                if (userMessageIndex !== -1) {
                  messages.value[userMessageIndex] = parsedData
                  console.log('Updated user message:', parsedData)
                }
              }
              
              // 处理AI消息
              if (parsedData.role === 'assistant') {
                // 第一次收到AI数据时，隐藏加载指示器
                if (!hasReceivedFirstData) {
                  hasReceivedFirstData = true
                  isLoading.value = false
                }
                
                if (aiMessageIndex === -1) {
                  // 只有在有内容时才添加AI消息
                  if (parsedData.content && parsedData.content.trim()) {
                    messages.value.push(parsedData)
                    aiMessageIndex = messages.value.length - 1
                    console.log('Added AI message:', parsedData)
                  }
                } else {
                  // 更新AI消息内容
                  messages.value[aiMessageIndex] = parsedData
                }
                
                // 调用流式数据回调
                if (onStreamData) {
                  onStreamData(parsedData)
                }
              }
            } catch (e) {
              // JSON解析失败，跳过这行数据
              console.warn('Failed to parse JSON:', data, e)
              continue
            }
          }
        }
      }

      return { success: true }
    } catch (error) {
      console.error('Send streaming message error:', error)
      return { 
        success: false, 
        error: error.message || '发送消息失败' 
      }
    } finally {
      // 确保在最后也设置加载状态为false
      isLoading.value = false
    }
  }

  const deleteConversation = async (conversationId) => {
    try {
      await api.delete(`/conversations/${conversationId}`)
      conversations.value = conversations.value.filter(c => c.id !== conversationId)
      if (currentConversation.value?.id === conversationId) {
        currentConversation.value = null
        messages.value = []
      }
      return { success: true }
    } catch (error) {
      return { 
        success: false, 
        error: error.response?.data?.error || '删除对话失败' 
      }
    }
  }

  const batchDeleteConversations = async (conversationIds) => {
    try {
      const response = await api.delete('/conversations/batch', {
        data: { ids: conversationIds }
      })
      
      // 从对话列表中移除已删除的对话
      conversations.value = conversations.value.filter(c => !conversationIds.includes(c.id))
      
      // 如果当前对话被删除，清空当前对话
      if (currentConversation.value && conversationIds.includes(currentConversation.value.id)) {
        currentConversation.value = null
        messages.value = []
      }
      
      return { 
        success: true, 
        deletedCount: response.data.deleted_count 
      }
    } catch (error) {
      return { 
        success: false, 
        error: error.response?.data?.error || '批量删除对话失败' 
      }
    }
  }

  const clearCurrentConversation = () => {
    currentConversation.value = null
    messages.value = []
  }

  return {
    characters,
    conversations,
    currentConversation,
    messages,
    isLoading,
    fetchCharacters,
    searchCharacters,
    fetchConversations,
    createConversation,
    fetchConversation,
    sendMessage,
    sendMessageStream,
    deleteConversation,
    batchDeleteConversations,
    clearCurrentConversation
  }
})

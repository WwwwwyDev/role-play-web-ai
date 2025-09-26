<template>
  <div v-cloak class="h-screen flex bg-gradient-to-br from-slate-50 to-blue-50">
    <!-- 左侧边栏 - 历史对话 -->
    <div 
      class="w-80 md:w-80 w-full md:relative fixed inset-0 md:inset-auto bg-white/95 md:bg-white/80 backdrop-blur-sm border-r border-gray-200/50 flex flex-col shadow-xl transition-all duration-300 ease-in-out z-40"
      :class="{ 'hidden': sidebarCollapsed }"
    >
      <!-- 侧边栏头部 -->
      <div class="p-6 border-b border-gray-200/50 bg-gradient-to-r from-blue-50 to-indigo-50">
        <div class="flex items-center justify-between mb-6">
          <div class="flex items-center space-x-3">
            <div class="w-8 h-8 bg-gradient-to-r from-blue-500 to-indigo-600 rounded-lg flex items-center justify-center">
              <MessageCircle class="w-5 h-5 text-white" />
            </div>
            <h2 class="text-xl font-bold bg-gradient-to-r from-gray-800 to-gray-600 bg-clip-text text-transparent">对话历史</h2>
          </div>
          <!-- 移动端关闭按钮 -->
          <button
            @click="sidebarCollapsed = true"
            class="md:hidden p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg transition-all duration-300"
            title="关闭"
          >
            <X class="w-5 h-5" />
          </button>
          <div class="flex items-center space-x-2">
            <!-- 批量操作按钮 -->
            <button
              v-if="!isBatchMode"
              @click="toggleBatchMode"
              class="p-2 text-gray-600 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-all duration-300"
              title="批量操作"
            >
              <CheckSquare class="w-5 h-5" />
            </button>
            <button
              v-if="isBatchMode"
              @click="toggleBatchMode"
              class="p-2 text-blue-600 bg-blue-50 rounded-lg transition-all duration-300"
              title="取消批量操作"
            >
              <X class="w-5 h-5" />
            </button>
            <!-- 批量删除按钮 -->
            <button
              v-if="isBatchMode && selectedConversations.size > 0"
              @click="showDeleteDialog()"
              class="p-2 text-red-600 hover:text-red-700 hover:bg-red-50 rounded-lg transition-all duration-300"
              title="批量删除"
            >
              <Trash2 class="w-5 h-5" />
            </button>
            <!-- 新建对话按钮 -->
            <button
              @click="showCharacterSelector = true"
              class="p-3 text-white bg-gradient-to-r from-blue-500 to-indigo-600 hover:from-blue-600 hover:to-indigo-700 rounded-xl shadow-lg"
              title="新建对话"
            >
              <Plus class="w-5 h-5" />
            </button>
          </div>
        </div>
        
        <!-- 搜索框 -->
        <div class="relative flex items-center space-x-3">
          <div class="relative flex-1">
            <Search class="absolute left-4 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
            <input
              v-model="searchQuery"
              type="text"
              placeholder="搜索对话..."
              class="w-full pl-12 pr-4 py-3 bg-white/70 backdrop-blur-sm border border-gray-200/50 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500/50 focus:border-blue-300 text-sm shadow-sm hover:shadow-md transition-all duration-300"
            />
          </div>
          <button
            @click="fetchConversations"
            :disabled="isLoadingConversations"
            class="p-3 text-gray-500 bg-white/70 backdrop-blur-sm hover:text-blue-600 hover:bg-blue-50 rounded-xl transition-all duration-300 disabled:opacity-50 disabled:cursor-not-allowed shadow-sm hover:shadow-md"
            title="刷新对话列表"
          >
            <Loader2 v-if="isLoadingConversations" class="w-5 h-5 animate-spin" />
            <RefreshCw v-else class="w-5 h-5" />
          </button>
        </div>
      </div>

      <!-- 对话列表 -->
      <div class="flex-1 overflow-y-auto">
        <!-- 骨架屏加载状态 -->
        <div v-if="!isInitialized && isLoadingConversations" class="p-3 space-y-2">
          <div v-for="i in 3" :key="i" class="p-4 rounded-xl bg-gray-100/50 skeleton-pulse">
            <div class="flex items-center space-x-3">
              <div class="w-10 h-10 bg-gray-200 rounded-full"></div>
              <div class="flex-1">
                <div class="h-4 bg-gray-200 rounded w-3/4 mb-2"></div>
                <div class="h-3 bg-gray-200 rounded w-1/2"></div>
              </div>
            </div>
          </div>
        </div>
        
        <div v-else-if="conversationsError" class="p-6 text-center">
          <div class="w-16 h-16 bg-red-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <AlertCircle class="w-8 h-8 text-red-500" />
          </div>
          <p class="text-red-600 text-sm mb-4">{{ conversationsError }}</p>
          <button
            @click="fetchConversations"
            class="px-4 py-2 bg-red-600 text-white text-sm font-medium rounded-lg hover:bg-red-700 transition-all duration-300"
          >
            重试
          </button>
        </div>
        
        <div v-else-if="filteredConversations.length === 0" class="p-6 text-center">
          <div class="w-16 h-16 bg-gradient-to-r from-blue-100 to-indigo-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <MessageCircle class="w-8 h-8 text-blue-500" />
          </div>
          <p class="text-gray-600 text-sm mb-4">还没有对话记录</p>
          <button
            @click="showCharacterSelector = true"
            class="px-4 py-2 bg-gradient-to-r from-blue-500 to-indigo-600 text-white text-sm font-medium rounded-lg hover:from-blue-600 hover:to-indigo-700 shadow-sm"
          >
            开始第一个对话
          </button>
        </div>

        <div v-else class="p-3 space-y-2">
          <!-- 批量模式下的全选按钮 -->
          <div v-if="isBatchMode && filteredConversations.length > 0" class="p-3 bg-blue-50 rounded-xl border border-blue-200">
            <div class="flex items-center justify-between">
              <button
                @click="selectAllConversations"
                class="flex items-center space-x-2 text-sm font-medium text-blue-700 hover:text-blue-800"
              >
                <CheckSquare v-if="selectedConversations.size === filteredConversations.length" class="w-4 h-4" />
                <Square v-else class="w-4 h-4" />
                <span>{{ selectedConversations.size === filteredConversations.length ? '取消全选' : '全选' }}</span>
              </button>
              <span class="text-xs text-blue-600">
                已选择 {{ selectedConversations.size }} / {{ filteredConversations.length }}
              </span>
            </div>
          </div>
          
          <!-- 对话列表 -->
          <div
            v-for="conversation in filteredConversations"
            :key="conversation.id"
            class="p-4 rounded-xl transition-all duration-300 group relative"
            :class="[
              chatStore.currentConversation?.id === conversation.id 
                ? 'bg-gradient-to-r from-blue-50 to-indigo-50 border border-blue-200 shadow-md' 
                : 'hover:bg-white/80 hover:shadow-md border border-transparent hover:border-gray-200/50',
              isBatchMode ? 'cursor-default' : 'cursor-pointer',
              selectedConversations.has(conversation.id) ? 'ring-2 ring-blue-500 ring-opacity-50' : ''
            ]"
            @click="isBatchMode ? toggleConversationSelection(conversation.id) : selectConversation(conversation)"
          >
            <!-- 批量模式选择框 -->
            <div v-if="isBatchMode" class="absolute left-3 top-1/2 transform -translate-y-1/2">
              <button
                @click.stop="toggleConversationSelection(conversation.id)"
                class="w-5 h-5 rounded border-2 flex items-center justify-center transition-colors duration-200"
                :class="selectedConversations.has(conversation.id) 
                  ? 'bg-blue-500 border-blue-500 text-white' 
                  : 'border-gray-300 hover:border-blue-400'"
              >
                <CheckSquare v-if="selectedConversations.has(conversation.id)" class="w-3 h-3" />
              </button>
            </div>
            
            <div class="flex items-center space-x-3" :class="{ 'ml-8': isBatchMode }">
              <div class="w-10 h-10 rounded-full overflow-hidden border-2 border-white shadow-sm flex-shrink-0 group-hover:border-blue-200 transition-colors duration-300">
                <img 
                  v-if="conversation.character?.avatar_url" 
                  :src="conversation.character.avatar_url" 
                  :alt="conversation.character.name"
                  class="w-full h-full object-cover"
                  @error="handleAvatarError"
                />
                <div v-else class="w-full h-full bg-gradient-to-r from-blue-100 to-indigo-100 flex items-center justify-center fallback-avatar">
                  <User class="w-5 h-5 text-blue-600" />
                </div>
              </div>
              <div class="flex-1 min-w-0">
                <h4 class="font-semibold text-gray-900 text-sm truncate group-hover:text-blue-700 transition-colors duration-300">
                  {{ conversation.character?.name || '未知角色' }}
                </h4>
                <p class="text-xs text-gray-500 truncate">
                  {{ formatTime(conversation.updated_at) }}
                </p>
              </div>
            </div>
            
            <!-- 删除按钮 -->
            <div class="absolute right-3 top-1/2 transform -translate-y-1/2 opacity-0 group-hover:opacity-100 transition-opacity duration-200">
              <button
                @click.stop="showDeleteDialog(conversation.id)"
                class="p-1.5 text-red-500 hover:text-red-700 hover:bg-red-50 rounded-lg transition-all duration-200"
                title="删除对话"
              >
                <Trash2 class="w-4 h-4" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 右侧聊天区域 -->
    <div class="flex-1 flex flex-col pb-16 transition-all duration-300 ease-in-out"
         :class="sidebarCollapsed ? 'w-full' : 'md:w-auto'">
      <!-- Chat Header -->
      <div class="bg-white/80 backdrop-blur-sm border-b border-gray-200/50 px-6 py-4">
        <div class="flex items-center justify-between">
          <div class="flex items-center space-x-4">
            <!-- 侧边栏切换按钮 -->
            <button
              @click="toggleSidebar"
              class="p-2 text-gray-600 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-all duration-300"
              :title="sidebarCollapsed ? '显示对话历史' : '隐藏对话历史'"
            >
              <PanelLeftOpen v-if="sidebarCollapsed" class="w-5 h-5" />
              <PanelLeftClose v-else class="w-5 h-5" />
            </button>
            
            <!-- 角色信息 -->
            <div v-if="chatStore.currentConversation?.character" class="flex items-center space-x-3">
              <div class="w-10 h-10 rounded-full overflow-hidden border-2 border-white shadow-sm">
                <img 
                  v-if="chatStore.currentConversation.character.avatar_url" 
                  :src="chatStore.currentConversation.character.avatar_url" 
                  :alt="chatStore.currentConversation.character.name"
                  class="w-full h-full object-cover"
                  @error="handleAvatarError"
                />
                <div v-else class="w-full h-full bg-gradient-to-r from-blue-100 to-indigo-100 flex items-center justify-center">
                  <User class="w-6 h-6 text-blue-600" />
                </div>
              </div>
              <div>
                <h1 class="text-lg font-semibold text-gray-900">
                  {{ chatStore.currentConversation.character.name }}
                </h1>
                <p class="text-sm text-gray-500">
                  {{ chatStore.currentConversation.character.description }}
                </p>
              </div>
            </div>
            
            <!-- 未选择角色时的提示 -->
            <div v-else class="flex items-center space-x-3">
              <div class="w-10 h-10 rounded-full bg-gradient-to-r from-blue-100 to-indigo-100 flex items-center justify-center">
                <MessageCircle class="w-6 h-6 text-blue-600" />
              </div>
              <div>
                <h1 class="text-lg font-semibold text-gray-900">选择角色开始对话</h1>
                <p class="text-sm text-gray-500">请选择一个角色开始新的对话</p>
              </div>
            </div>
          </div>
          
          <!-- 右上角按钮组 -->
          <div class="flex items-center space-x-2">
            <!-- 个人资料按钮 -->
            <button
              @click="showProfileModal = true"
              class="p-2 text-gray-600 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-all duration-300"
              title="个人资料"
            >
              <User class="w-5 h-5" />
            </button>
          </div>
        </div>
      </div>

      <!-- Messages Area -->
      <div ref="messagesContainerRef" class="flex-1 overflow-y-auto p-6 space-y-6">
        <div v-if="!chatStore.currentConversation" class="flex items-center justify-center h-full">
          <div class="text-center">
            <div class="w-24 h-24 bg-gradient-to-r from-blue-100 to-indigo-100 rounded-full flex items-center justify-center mx-auto mb-6">
              <MessageCircle class="w-12 h-12 text-blue-500" />
            </div>
            <h2 class="text-2xl font-bold text-gray-900 mb-4">欢迎使用AI角色扮演聊天</h2>
            <p class="text-gray-600 mb-8 max-w-md mx-auto">
              选择一个角色开始对话，体验与不同角色的有趣交流
            </p>
            <button
              @click="showCharacterSelector = true"
              class="px-6 py-3 bg-gradient-to-r from-blue-500 to-indigo-600 text-white font-semibold rounded-xl hover:from-blue-600 hover:to-indigo-700 shadow-lg"
            >
              选择角色开始对话
            </button>
          </div>
        </div>

        <div v-else-if="chatStore.messages.length === 0" class="flex items-center justify-center h-full">
          <div class="text-center">
            <div class="w-20 h-20 bg-gradient-to-r from-blue-100 to-indigo-100 rounded-full flex items-center justify-center mx-auto mb-6">
              <MessageCircle class="w-10 h-10 text-blue-500" />
            </div>
            <h3 class="text-xl font-semibold text-gray-900 mb-2">
              与 {{ chatStore.currentConversation.character.name }} 的对话
            </h3>
            <p class="text-gray-600 mb-6">
              {{ chatStore.currentConversation.character.description }}
            </p>
            <p class="text-sm text-gray-500">开始发送消息来开始对话</p>
          </div>
        </div>

        <div v-else class="space-y-6">
          <div
            v-for="message in chatStore.messages"
            :key="message.id"
            class="flex"
            :class="message.role === 'user' ? 'justify-end' : 'justify-start'"
          >
            <div
              class="max-w-3xl px-6 py-4 rounded-2xl shadow-md"
              :class="message.role === 'user' 
                ? 'bg-gradient-to-r from-blue-500 to-indigo-600 text-white' 
                : 'bg-white/90 backdrop-blur-sm text-gray-900'"
            >
              <div class="flex items-start space-x-3">
                <div v-if="message.role === 'assistant'" class="w-8 h-8 rounded-full overflow-hidden border-2 border-white shadow-sm flex-shrink-0">
                  <img 
                    v-if="chatStore.currentConversation?.character?.avatar_url" 
                    :src="chatStore.currentConversation.character.avatar_url" 
                    :alt="chatStore.currentConversation.character.name"
                    class="w-full h-full object-cover"
                    @error="handleAvatarError"
                  />
                  <div v-else class="w-full h-full bg-gradient-to-r from-blue-100 to-indigo-100 flex items-center justify-center">
                    <User class="w-4 h-4 text-blue-600" />
                  </div>
                </div>
                <div class="flex-1">
                  <MarkdownRenderer 
                    v-if="message.role === 'assistant'" 
                    :content="message.content" 
                  />
                  <p v-else class="whitespace-pre-wrap">{{ message.content }}</p>
                  <p class="text-xs mt-2 opacity-70">
                    {{ formatTime(message.created_at) }}
                  </p>
                </div>
              </div>
            </div>
          </div>
          
          <!-- AI正在思考的加载指示器 -->
          <div v-if="chatStore.isLoading && chatStore.messages.length > 0" class="flex justify-start">
            <div class="max-w-3xl px-6 py-5 rounded-2xl shadow-xl bg-gradient-to-r from-white/95 to-blue-50/95 backdrop-blur-sm text-gray-900 border border-blue-100/50 relative overflow-hidden">
              <!-- 背景装饰 -->
              <div class="absolute inset-0 bg-gradient-to-r from-blue-50/30 to-indigo-50/30"></div>
              <div class="absolute top-0 right-0 w-20 h-20 bg-blue-100/20 rounded-full -translate-y-10 translate-x-10"></div>
              
              <div class="relative flex items-start space-x-4">
                <!-- 角色头像 -->
                <div class="w-12 h-12 rounded-full overflow-hidden border-2 border-blue-200 shadow-lg flex-shrink-0 relative">
                  <img 
                    v-if="chatStore.currentConversation?.character?.avatar_url" 
                    :src="chatStore.currentConversation.character.avatar_url" 
                    :alt="chatStore.currentConversation.character.name"
                    class="w-full h-full object-cover"
                    @error="handleAvatarError"
                  />
                  <div v-else class="w-full h-full bg-gradient-to-r from-blue-100 to-indigo-100 flex items-center justify-center">
                    <User class="w-6 h-6 text-blue-600" />
                  </div>
                  <!-- 头像周围的脉冲效果 -->
                  <div class="absolute inset-0 rounded-full border-2 border-blue-300 opacity-30"></div>
                  <div class="absolute inset-0 rounded-full border border-blue-400 opacity-20"></div>
                </div>
                
                <!-- 加载动画区域 -->
                <div class="flex-1">
                  <div class="flex items-center space-x-3 mb-3">
                    <span class="text-sm font-semibold text-gray-700">正在思考中</span>
                    <div class="flex items-center space-x-1">
                      <!-- 美化的三个点动画 -->
                      <div class="w-2.5 h-2.5 bg-gradient-to-r from-blue-400 to-blue-600 rounded-full"></div>
                      <div class="w-2.5 h-2.5 bg-gradient-to-r from-blue-400 to-blue-600 rounded-full"></div>
                      <div class="w-2.5 h-2.5 bg-gradient-to-r from-blue-400 to-blue-600 rounded-full"></div>
                    </div>
                  </div>
                  
                  <!-- 思考进度条 -->
                  <div class="w-full bg-gray-200/60 rounded-full h-2 overflow-hidden mb-3">
                    <div class="h-full bg-gradient-to-r from-blue-400 via-blue-500 to-blue-600 rounded-full relative">
                      <div class="absolute inset-0 bg-gradient-to-r from-transparent via-white/30 to-transparent"></div>
                    </div>
                  </div>
                  
                  <!-- 思考提示文字 -->
                  <div class="flex items-center space-x-2 text-xs text-gray-600">
                    <div class="w-1.5 h-1.5 bg-blue-400 rounded-full"></div>
                    <span class="font-medium">{{ chatStore.currentConversation?.character?.name || 'AI' }} 正在组织语言...</span>
                    <div class="flex space-x-1">
                      <div class="w-1 h-1 bg-blue-300 rounded-full"></div>
                      <div class="w-1 h-1 bg-blue-300 rounded-full"></div>
                      <div class="w-1 h-1 bg-blue-300 rounded-full"></div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          
        </div>
      </div>

    </div>

    <!-- 固定输入条 -->
    <div class="fixed bottom-0 bg-white/95 backdrop-blur-md border-t border-gray-200/50 shadow-lg z-30 transition-all duration-300 ease-in-out"
         :class="sidebarCollapsed ? 'left-0 right-0' : 'left-0 md:left-80 right-0'">
      <div class="max-w-7xl mx-auto px-4 py-3">
        <div class="flex items-center space-x-3">
          <!-- 输入框容器 -->
          <div class="flex-1 relative">
            <textarea
              v-model="messageInput"
              ref="messageInputRef"
              placeholder="输入消息..."
              class="w-full px-3 py-2 bg-white/80 backdrop-blur-sm border border-gray-200/50 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500/50 focus:border-blue-300 resize-none shadow-sm hover:shadow-md transition-all duration-300 text-sm min-h-[40px] max-h-24 leading-5"
              rows="1"
              @keydown="handleKeyDown"
              @input="autoResize"
            ></textarea>
            <!-- 输入框装饰 -->
            <div class="absolute inset-0 rounded-xl bg-gradient-to-r from-blue-500/5 to-indigo-500/5 pointer-events-none opacity-0 transition-opacity duration-300"
                 :class="{ 'opacity-100': messageInput.trim() }"></div>
          </div>
          
          <!-- 语音识别按钮 -->
          <button
            @click="toggleVoiceRecording"
            :disabled="chatStore.isLoading"
            class="flex-shrink-0 w-10 h-10 rounded-lg transition-all duration-300 shadow-sm hover:shadow-md flex items-center justify-center relative overflow-hidden"
            :class="isVoiceRecording 
              ? 'bg-gradient-to-r from-red-500 to-red-600 text-white hover:from-red-600 hover:to-red-700' 
              : 'bg-white/80 backdrop-blur-sm text-gray-600 hover:text-red-500 hover:bg-red-50 border border-gray-200/50 hover:border-red-200'"
            :title="isVoiceRecording ? '停止语音识别' : '开始语音输入'"
          >
            <!-- 录音时的脉冲效果 -->
            <div v-if="isVoiceRecording" class="absolute inset-0 bg-red-400 rounded-lg opacity-30"></div>
            <Mic v-if="!isVoiceRecording" class="w-4 h-4 relative z-10" />
            <Square v-else class="w-4 h-4 relative z-10" />
          </button>
          
          <!-- 发送按钮 -->
          <button
            @click="sendMessage"
            :disabled="!messageInput.trim() || chatStore.isLoading"
            class="flex-shrink-0 w-10 h-10 bg-gradient-to-r from-blue-500 to-indigo-600 text-white rounded-lg hover:from-blue-600 hover:to-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed shadow-lg flex items-center justify-center"
            title="发送消息"
          >
            <Send v-if="!chatStore.isLoading" class="w-4 h-4" />
            <div v-else class="w-4 h-4 border-2 border-white border-t-transparent rounded-full"></div>
          </button>
        </div>
        
        <!-- 底部装饰线 -->
        <div class="mt-2 h-px bg-gradient-to-r from-transparent via-gray-200/50 to-transparent"></div>
      </div>
    </div>

    <!-- Character Selector Modal -->
    <div
      v-if="showCharacterSelector"
      class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50"
      @click="showCharacterSelector = false"
    >
      <div
        class="bg-white/95 backdrop-blur-md rounded-3xl p-8 max-w-5xl w-full mx-4 max-h-[85vh] overflow-hidden flex flex-col shadow-2xl border border-white/20"
        @click.stop
      >
        <div class="flex items-center justify-between mb-8">
          <div class="flex items-center space-x-3">
            <div class="w-10 h-10 bg-gradient-to-r from-blue-500 to-indigo-600 rounded-xl flex items-center justify-center">
              <Users class="w-6 h-6 text-white" />
            </div>
            <h3 class="text-3xl font-bold bg-gradient-to-r from-gray-800 to-gray-600 bg-clip-text text-transparent">选择角色</h3>
          </div>
          <button
            @click="showCharacterSelector = false"
            class="p-3 text-gray-400 hover:bg-gray-100 rounded-xl"
          >
            <X class="w-6 h-6" />
          </button>
        </div>
        
        <!-- 搜索和分类 -->
        <div class="mb-8">
          <div class="relative mb-6">
            <Search class="absolute left-4 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" />
            <input
              v-model="characterSearchQuery"
              type="text"
              placeholder="搜索角色..."
              class="w-full pl-12 pr-4 py-4 bg-white/70 backdrop-blur-sm border border-gray-200/50 rounded-2xl focus:outline-none focus:ring-2 focus:ring-blue-500/50 focus:border-blue-300 shadow-sm hover:shadow-md transition-all duration-300 text-lg"
            />
          </div>
          
          <!-- 分类筛选 -->
          <div class="flex flex-wrap gap-3">
            <button
              v-for="category in categories"
              :key="category"
              @click="selectedCategory = selectedCategory === category ? null : category"
              class="px-6 py-3 rounded-full text-sm font-semibold"
              :class="selectedCategory === category 
                ? 'bg-gradient-to-r from-blue-500 to-indigo-600 text-white shadow-lg' 
                : 'bg-white/70 backdrop-blur-sm text-gray-700 hover:bg-blue-50 border border-gray-200/50 shadow-sm'"
            >
              {{ category }}
            </button>
          </div>
        </div>
        
        <!-- 角色列表 -->
        <div class="flex-1 overflow-y-auto">
          <div class="grid md:grid-cols-2 lg:grid-cols-3 gap-4">
            <div
              v-for="character in filteredCharacters"
              :key="character.id"
              class="p-6 bg-white/70 backdrop-blur-sm border border-gray-200/50 rounded-2xl hover:bg-white hover:border-blue-300 cursor-pointer group shadow-sm"
              @click="selectCharacter(character)"
            >
              <div class="text-center">
                <div class="w-20 h-20 rounded-full mx-auto mb-4 overflow-hidden border-3 border-white shadow-lg group-hover:border-blue-200">
                  <img 
                    v-if="character.avatar_url" 
                    :src="character.avatar_url" 
                    :alt="character.name"
                    class="w-full h-full object-cover"
                    @error="handleAvatarError"
                  />
                  <div v-else class="w-full h-full bg-gradient-to-r from-blue-100 to-indigo-100 flex items-center justify-center fallback-avatar">
                    <User class="w-10 h-10 text-blue-600" />
                  </div>
                </div>
                <h4 class="font-bold text-gray-900 mb-2 text-lg group-hover:text-blue-700">{{ character.name }}</h4>
                <p class="text-sm text-gray-600 mb-3 leading-relaxed">{{ character.description }}</p>
                <span class="inline-block bg-gradient-to-r from-blue-100 to-indigo-100 text-blue-800 text-xs px-3 py-1.5 rounded-full font-medium">
                  {{ character.category }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 删除确认弹窗 -->
    <div
      v-if="showDeleteConfirm"
      class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-[9998]"
      @click="cancelDelete"
    >
      <div
        class="bg-white/95 backdrop-blur-md rounded-2xl p-6 max-w-md w-full mx-4 shadow-2xl border border-white/20"
        @click.stop
      >
        <div class="flex items-center space-x-3 mb-4">
          <div class="w-10 h-10 bg-red-100 rounded-full flex items-center justify-center">
            <Trash2 class="w-6 h-6 text-red-600" />
          </div>
          <h3 class="text-xl font-bold text-gray-900">确认删除</h3>
        </div>
        
        <p class="text-gray-600 mb-6">
          {{ Array.isArray(deleteTarget) 
            ? `确定要删除选中的 ${deleteTarget.length} 个对话吗？此操作无法撤销。`
            : '确定要删除这个对话吗？此操作无法撤销。'
          }}
        </p>
        
        <div class="flex space-x-3">
          <button
            @click="cancelDelete"
            class="flex-1 px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors duration-200"
            :disabled="isDeleting"
          >
            取消
          </button>
          <button
            @click="confirmDelete"
            class="flex-1 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
            :disabled="isDeleting"
          >
            <span v-if="isDeleting">删除中...</span>
            <span v-else>确认删除</span>
          </button>
        </div>
      </div>
    </div>

    <!-- 个人资料弹窗 -->
    <div
      v-if="showProfileModal"
      class="fixed inset-0 bg-black/50 backdrop-blur-md flex items-center justify-center z-[9999]"
      @click="showProfileModal = false"
    >
      <div
        class="bg-white/95 backdrop-blur-xl rounded-3xl p-8 max-w-lg w-full mx-4 shadow-2xl border border-white/30 relative overflow-hidden"
        @click.stop
      >
        <!-- 背景装饰 -->
        <div class="absolute top-0 left-0 w-full h-32 bg-gradient-to-br from-blue-500/10 via-indigo-500/10 to-purple-500/10 rounded-t-3xl"></div>
        <div class="absolute top-4 right-4 w-20 h-20 bg-gradient-to-br from-blue-400/20 to-indigo-400/20 rounded-full blur-xl"></div>
        <div class="absolute bottom-4 left-4 w-16 h-16 bg-gradient-to-br from-purple-400/20 to-pink-400/20 rounded-full blur-lg"></div>
        
        <!-- 关闭按钮 -->
        <button 
          @click="() => { console.log('Close button clicked'); showProfileModal = false; }" 
          class="absolute top-6 right-6 p-2 text-gray-400 hover:text-gray-600 hover:bg-white/80 rounded-xl z-50 cursor-pointer"
          title="关闭"
        >
          <X class="w-5 h-5" />
        </button>
        
        <!-- 头部 -->
        <div class="relative z-10 text-center mb-8">
          <div class="w-28 h-28 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-full flex items-center justify-center mx-auto mb-6 shadow-xl border-4 border-white/50 relative">
            <User class="w-14 h-14 text-white" />
            <div class="absolute -bottom-1 -right-1 w-8 h-8 bg-green-500 rounded-full border-4 border-white flex items-center justify-center">
              <div class="w-3 h-3 bg-white rounded-full"></div>
            </div>
          </div>
          <h2 class="text-3xl font-bold bg-gradient-to-r from-gray-800 to-gray-600 bg-clip-text text-transparent mb-2">
            {{ authStore.user?.username }}
          </h2>
          <p class="text-gray-600 text-lg mb-3">{{ authStore.user?.email }}</p>
          <div class="inline-flex items-center px-4 py-2 bg-gradient-to-r from-blue-50 to-indigo-50 rounded-full border border-blue-200/50">
            <div class="w-2 h-2 bg-green-500 rounded-full mr-2"></div>
            <span class="text-sm font-medium text-blue-700">在线</span>
          </div>
        </div>
        
        <!-- 用户信息卡片 -->
        <div class="relative z-10 bg-white/70 backdrop-blur-sm rounded-2xl p-6 mb-6 border border-gray-200/50 shadow-lg">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-semibold text-gray-800">账户信息</h3>
            <div class="w-8 h-8 bg-gradient-to-r from-blue-100 to-indigo-100 rounded-lg flex items-center justify-center">
              <User class="w-4 h-4 text-blue-600" />
            </div>
          </div>
          <div class="space-y-3">
            <div class="flex items-center justify-between py-2">
              <span class="text-sm text-gray-600">注册时间</span>
              <span class="text-sm font-medium text-gray-800">{{ formatDate(authStore.user?.created_at) }}</span>
            </div>
            <div class="flex items-center justify-between py-2">
              <span class="text-sm text-gray-600">用户ID</span>
              <span class="text-sm font-medium text-gray-800 font-mono">#{{ authStore.user?.id || 'N/A' }}</span>
            </div>
          </div>
        </div>
        
        <!-- 对话统计卡片 -->
        <div class="relative z-10 bg-gradient-to-r from-blue-50 to-indigo-50 rounded-2xl p-6 mb-6 border border-blue-200/50">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-semibold text-gray-800">对话统计</h3>
            <div class="w-8 h-8 bg-gradient-to-r from-blue-500 to-indigo-600 rounded-lg flex items-center justify-center">
              <MessageCircle class="w-4 h-4 text-white" />
            </div>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div class="text-center">
              <div class="text-2xl font-bold text-blue-600 mb-1">{{ conversations.length }}</div>
              <div class="text-sm text-gray-600">总对话数</div>
            </div>
            <div class="text-center">
              <div class="text-2xl font-bold text-indigo-600 mb-1">{{ characters.length }}</div>
              <div class="text-sm text-gray-600">可用角色</div>
            </div>
          </div>
        </div>
        
        <!-- 操作按钮 -->
        <div class="relative z-10 space-y-3">
          <button 
            @click="handleLogout"
            class="w-full inline-flex items-center justify-center px-6 py-3 bg-gradient-to-r from-red-500 to-red-600 text-white font-semibold rounded-xl hover:from-red-600 hover:to-red-700 focus:outline-none focus:ring-2 focus:ring-red-500/50 focus:ring-offset-2 shadow-lg"
          >
            <div class="w-5 h-5 mr-2">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/>
                <polyline points="16,17 21,12 16,7"/>
                <line x1="21" y1="12" x2="9" y2="12"/>
              </svg>
            </div>
            退出登录
          </button>
        </div>
      </div>
    </div>

    <!-- 语音气泡 -->
    <VoiceBubble
      :is-visible="showVoiceBubble"
      @close="handleVoiceBubbleClose"
      @message="handleVoiceMessage"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useChatStore } from '@/stores/chat'
import { useAuthStore } from '@/stores/auth'
import MarkdownRenderer from '@/components/MarkdownRenderer.vue'
import { 
  MessageCircle, 
  Send, 
  Mic, 
  MicOff, 
  Search, 
  Users, 
  User, 
  Plus, 
  X, 
  Loader2, 
  RefreshCw,
  CheckSquare,
  Square,
  Trash2,
  AlertCircle,
  PanelLeftOpen,
  PanelLeftClose
} from 'lucide-vue-next'
import voiceService from '@/services/voice'
import VoiceBubble from '@/components/VoiceBubble.vue'

const route = useRoute()
const router = useRouter()
const chatStore = useChatStore()
const authStore = useAuthStore()

// 响应式数据
const messageInput = ref('')
const messageInputRef = ref(null)
const messagesContainerRef = ref(null)
const isVoiceRecording = ref(false)
const voiceError = ref('')
const voiceSuccess = ref('')
const audioQuality = ref({ level: 0, quality: 'unknown' })
const showVoiceBubble = ref(false)
const showCharacterSelector = ref(false)
const characters = ref([])
const conversations = ref([])
const searchQuery = ref('')
const characterSearchQuery = ref('')
const selectedCategory = ref(null)
const isLoadingConversations = ref(false)
const conversationsError = ref('')
const isInitialized = ref(false)

// 批量删除相关
const selectedConversations = ref(new Set())
const isBatchMode = ref(false)
const showDeleteConfirm = ref(false)
const deleteTarget = ref(null)
const isDeleting = ref(false)

// 侧边栏折叠状态 - 从localStorage读取，默认隐藏
const getInitialSidebarState = () => {
  const saved = localStorage.getItem('sidebarCollapsed')
  return saved !== null ? JSON.parse(saved) : true
}
const sidebarCollapsed = ref(getInitialSidebarState())

// 个人资料相关
const showProfileModal = ref(false)

// 计算属性
const categories = computed(() => {
  const cats = [...new Set(characters.value.map(c => c.category))]
  return cats.sort()
})

const filteredConversations = computed(() => {
  let filtered = conversations.value
  
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(conv => 
      conv.character?.name?.toLowerCase().includes(query) ||
      conv.title?.toLowerCase().includes(query)
    )
  }
  
  return filtered.sort((a, b) => new Date(b.updated_at) - new Date(a.updated_at))
})

const filteredCharacters = computed(() => {
  let filtered = characters.value
  
  if (characterSearchQuery.value) {
    const query = characterSearchQuery.value.toLowerCase()
    filtered = filtered.filter(char => 
      char.name.toLowerCase().includes(query) ||
      char.description.toLowerCase().includes(query) ||
      char.category.toLowerCase().includes(query)
    )
  }
  
  if (selectedCategory.value) {
    filtered = filtered.filter(char => char.category === selectedCategory.value)
  }
  
  return filtered
})

const lastAIMessage = computed(() => {
  const messages = chatStore.messages
  for (let i = messages.length - 1; i >= 0; i--) {
    if (messages[i].role === 'assistant') {
      return messages[i].content
    }
  }
  return ''
})

// 方法
const toggleSidebar = () => {
  sidebarCollapsed.value = !sidebarCollapsed.value
}

const toggleVoiceRecording = () => {
  showVoiceBubble.value = true
}

const handleVoiceBubbleClose = () => {
  showVoiceBubble.value = false
  // 停止任何正在进行的录音
  if (isVoiceRecording.value) {
    voiceService.stopRecognition()
    isVoiceRecording.value = false
  }
}

const handleVoiceMessage = (text) => {
  messageInput.value = text
  autoResize()
  showVoiceBubble.value = false
  
  // 聚焦到输入框
  if (messageInputRef.value) {
    messageInputRef.value.focus()
  }
}

const startVoiceRecording = async () => {
  try {
    voiceError.value = ''
    isVoiceRecording.value = true
    
    // 检查浏览器支持情况
    const support = voiceService.isSupported()
    if (!support.recognition) {
      // 尝试降级方案
      if (support.mediaDevices) {
        await voiceService.startFallbackRecognition()
        return
      } else {
        throw new Error('浏览器不支持语音识别功能')
      }
    }
    
    const text = await voiceService.startRecognition({
      onInterimResult: (transcript) => {
        // 实时更新输入框内容
        messageInput.value = transcript
        autoResize()
      },
      onFinalResult: (transcript) => {
        // 最终结果写入输入框
        messageInput.value = transcript
        autoResize()
        isVoiceRecording.value = false
        
        // 显示成功提示
        voiceSuccess.value = '语音识别成功！'
        setTimeout(() => {
          voiceSuccess.value = ''
        }, 2000)
        
        // 聚焦到输入框
        if (messageInputRef.value) {
          messageInputRef.value.focus()
        }
      },
      onQualityUpdate: (quality) => {
        // 更新音频质量数据
        audioQuality.value = quality
      },
      onError: (errorMsg, shouldRetry) => {
        if (shouldRetry) {
          // 如果是可重试的错误，显示重试提示
          voiceError.value = `${errorMsg} (正在重试...)`
        } else {
          voiceError.value = errorMsg
          isVoiceRecording.value = false
        }
      }
    })
    
  } catch (err) {
    voiceError.value = getVoiceErrorMessage(err.message)
    isVoiceRecording.value = false
  }
}

const stopVoiceRecording = () => {
  voiceService.stopRecognition()
  isVoiceRecording.value = false
}

const retryVoiceRecording = async () => {
  voiceError.value = ''
  await startVoiceRecording()
}

const getVoiceErrorMessage = (error) => {
  // 如果错误信息已经包含中文，直接返回
  if (error.includes('网络连接错误') || error.includes('没有检测到语音') || 
      error.includes('无法访问麦克风') || error.includes('麦克风权限被拒绝') ||
      error.includes('语音识别被中断') || error.includes('不支持当前语言设置')) {
    return error
  }
  
  const errorMap = {
    'no-speech': '没有检测到语音，请重试',
    'audio-capture': '无法访问麦克风，请检查权限',
    'not-allowed': '麦克风权限被拒绝，请在浏览器设置中允许',
    'network': '网络连接错误，请检查网络连接',
    'aborted': '语音识别被中断',
    'language-not-supported': '不支持当前语言设置'
  }
  return errorMap[error] || `语音识别错误：${error}`
}

const autoResize = () => {
  nextTick(() => {
    if (messageInputRef.value) {
      messageInputRef.value.style.height = 'auto'
      messageInputRef.value.style.height = messageInputRef.value.scrollHeight + 'px'
    }
  })
}

// 自动滚动到消息底部
const scrollToBottom = (smooth = true) => {
  nextTick(() => {
    if (messagesContainerRef.value) {
      if (smooth) {
        messagesContainerRef.value.scrollTo({
          top: messagesContainerRef.value.scrollHeight,
          behavior: 'smooth'
        })
      } else {
        messagesContainerRef.value.scrollTop = messagesContainerRef.value.scrollHeight
      }
    }
  })
}

const handleKeyDown = (event) => {
  if (event.key === 'Enter' && !event.shiftKey) {
    event.preventDefault()
    sendMessage()
  }
}

const sendMessage = async () => {
  if (!messageInput.value.trim() || chatStore.isLoading) return
  
  const content = messageInput.value.trim()
  messageInput.value = ''
  autoResize()
  
  if (!chatStore.currentConversation?.id) {
    console.error('No current conversation')
    return
  }
  
  // 立即显示用户消息
  const userMessage = {
    id: Date.now(), // 临时ID
    role: 'user',
    content: content,
    created_at: new Date().toISOString(),
    conversation_id: chatStore.currentConversation.id
  }
  chatStore.messages.push(userMessage)
  
  // 滚动到底部显示用户消息
  scrollToBottom()
  
  // 发送消息到后端（使用流式响应）
  const result = await chatStore.sendMessageStream(
    chatStore.currentConversation.id, 
    content, 
    null, 
    (streamData) => {
      // 流式数据回调，每次收到新内容时滚动到底部
      scrollToBottom()
    }
  )
  
  // 如果发送失败，显示错误消息
  if (!result.success) {
    console.error('Send message failed:', result.error)
    const errorMessage = {
      id: Date.now() + 1,
      role: 'assistant',
      content: `抱歉，消息发送失败：${result.error || '请稍后重试'}`,
      created_at: new Date().toISOString(),
      conversation_id: chatStore.currentConversation.id
    }
    chatStore.messages.push(errorMessage)
    // 滚动到底部显示错误消息
    scrollToBottom()
  } else {
    console.log('Message sent successfully')
    // 流式响应完成后滚动到底部
    scrollToBottom()
  }
}


const handleAvatarError = (event) => {
  event.target.style.display = 'none'
}

const selectCharacter = async (character) => {
  showCharacterSelector.value = false
  
  // 创建新对话
  const result = await chatStore.createConversation(character.id)
  
  if (result.success) {
    // 立即切换到新创建的对话
    await chatStore.fetchConversation(result.conversation.id)
    router.push(`/chat/${result.conversation.id}`)
    
    // 刷新对话列表以显示新对话
    await fetchConversations()
  } else {
    console.error('创建对话失败:', result.error)
    // 可以在这里显示错误提示
  }
}

const selectConversation = async (conversation) => {
  if (isBatchMode.value) return
  
  await chatStore.fetchConversation(conversation.id)
  router.push(`/chat/${conversation.id}`)
}

const fetchConversations = async () => {
  isLoadingConversations.value = true
  conversationsError.value = ''
  
  try {
    const result = await chatStore.fetchConversations()
    if (result.success) {
      conversations.value = chatStore.conversations
    } else {
      conversationsError.value = result.error || '获取对话列表失败'
    }
  } catch (error) {
    conversationsError.value = '网络连接失败，请重试'
  } finally {
    isLoadingConversations.value = false
    isInitialized.value = true
  }
}

const fetchCharacters = async () => {
  try {
    const response = await fetch('/api/v1/characters')
    if (response.ok) {
      const data = await response.json()
      characters.value = data.characters || []
    }
  } catch (error) {
    console.error('获取角色列表失败:', error)
  }
}

const formatTime = (timestamp) => {
  if (!timestamp) return ''
  
  // 确保正确处理时间戳
  let date
  if (typeof timestamp === 'string') {
    // 如果是字符串，直接创建Date对象
    date = new Date(timestamp)
  } else {
    // 如果是数字或其他类型，也创建Date对象
    date = new Date(timestamp)
  }
  
  // 检查日期是否有效
  if (isNaN(date.getTime())) {
    return '时间格式错误'
  }
  
  const now = new Date()
  const diff = now - date
  
  if (diff < 60000) return '刚刚'
  if (diff < 3600000) return `${Math.floor(diff / 60000)}分钟前`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)}小时前`
  if (diff < 604800000) return `${Math.floor(diff / 86400000)}天前`
  
  return date.toLocaleDateString('zh-CN', {
    month: 'short',
    day: 'numeric'
  })
}

const formatDate = (timestamp) => {
  if (!timestamp) return ''
  
  // 确保正确处理时间戳
  let date
  if (typeof timestamp === 'string') {
    date = new Date(timestamp)
  } else {
    date = new Date(timestamp)
  }
  
  // 检查日期是否有效
  if (isNaN(date.getTime())) {
    return '时间格式错误'
  }
  
  return date.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 批量删除相关方法
const toggleBatchMode = () => {
  isBatchMode.value = !isBatchMode.value
  selectedConversations.value.clear()
}

const toggleConversationSelection = (conversationId) => {
  if (selectedConversations.value.has(conversationId)) {
    selectedConversations.value.delete(conversationId)
  } else {
    selectedConversations.value.add(conversationId)
  }
}

const selectAllConversations = () => {
  if (selectedConversations.value.size === filteredConversations.value.length) {
    selectedConversations.value.clear()
  } else {
    filteredConversations.value.forEach(conv => {
      selectedConversations.value.add(conv.id)
    })
  }
}

const showDeleteDialog = (conversationId = null) => {
  if (conversationId) {
    deleteTarget.value = conversationId
  } else {
    deleteTarget.value = Array.from(selectedConversations.value)
  }
  showDeleteConfirm.value = true
}

const confirmDelete = async () => {
  if (!deleteTarget.value) return
  
  isDeleting.value = true
  
  try {
    let result
    if (Array.isArray(deleteTarget.value)) {
      // 批量删除
      result = await chatStore.batchDeleteConversations(deleteTarget.value)
    } else {
      // 单个删除
      result = await chatStore.deleteConversation(deleteTarget.value)
    }
    
    if (result.success) {
      // 重新获取对话列表
      await fetchConversations()
      // 如果当前对话被删除，清空当前对话
      if (chatStore.currentConversation && 
          (Array.isArray(deleteTarget.value) ? deleteTarget.value.includes(chatStore.currentConversation.id) : deleteTarget.value === chatStore.currentConversation.id)) {
        chatStore.clearCurrentConversation()
        router.push('/chat')
      }
    }
  } catch (error) {
    console.error('删除对话失败:', error)
  } finally {
    isDeleting.value = false
    showDeleteConfirm.value = false
    deleteTarget.value = null
    if (isBatchMode.value) {
      toggleBatchMode()
    }
  }
}

const cancelDelete = () => {
  showDeleteConfirm.value = false
  deleteTarget.value = null
}

// 个人资料相关方法
const handleLogout = () => {
  authStore.logout()
  showProfileModal.value = false
  router.push('/')
}

// 监听路由变化
watch(() => route.params.conversationId, async (newConversationId) => {
  if (newConversationId && newConversationId !== chatStore.currentConversation?.id) {
    await chatStore.fetchConversation(newConversationId)
    // 切换对话后滚动到底部
    scrollToBottom(false) // 不使用平滑滚动，立即跳转
  }
})

// 监听消息变化，自动滚动到底部
watch(() => chatStore.messages, () => {
  scrollToBottom()
}, { deep: true })

// 监听侧边栏状态变化，自动保存到localStorage
watch(sidebarCollapsed, (newValue) => {
  localStorage.setItem('sidebarCollapsed', JSON.stringify(newValue))
})

// 生命周期
onMounted(async () => {
  // 确保侧边栏状态从localStorage正确加载
  const savedState = localStorage.getItem('sidebarCollapsed')
  if (savedState !== null) {
    sidebarCollapsed.value = JSON.parse(savedState)
  }
  
  await Promise.all([
    fetchCharacters(),
    fetchConversations()
  ])
  
  // 如果有对话ID参数，加载对应对话
  if (route.params.conversationId) {
    await chatStore.fetchConversation(route.params.conversationId)
  }
  
  // 如果有角色ID参数，创建新对话
  const characterId = route.query.character
  if (characterId && !route.params.conversationId) {
    const character = characters.value.find(c => c.id == characterId)
    if (character) {
      await selectCharacter(character)
    }
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

@keyframes slide-up {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes float {
  0%, 100% {
    transform: translateY(0px);
  }
  50% {
    transform: translateY(-10px);
  }
}

@keyframes thinking-pulse {
  0%, 100% {
    opacity: 0.4;
    transform: scale(1);
  }
  50% {
    opacity: 1;
    transform: scale(1.1);
  }
}

@keyframes progress-wave {
  0% {
    transform: translateX(-100%);
  }
  100% {
    transform: translateX(100%);
  }
}

@keyframes thinking-dots {
  0%, 20% {
    transform: scale(1);
    opacity: 0.7;
  }
  50% {
    transform: scale(1.3);
    opacity: 1;
  }
  80%, 100% {
    transform: scale(1);
    opacity: 0.7;
  }
}

.animate-fade-in {
  animation: fade-in 0.3s ease-out;
}

.animate-slide-up {
  animation: slide-up 0.3s ease-out;
}

.animate-float {
  animation: float 3s ease-in-out infinite;
}

.animate-thinking-pulse {
  animation: thinking-pulse 2s ease-in-out infinite;
}

.animate-progress-wave {
  animation: progress-wave 2s ease-in-out infinite;
}

.animate-thinking-dots {
  animation: thinking-dots 1.5s ease-in-out infinite;
}

/* 自定义滚动条 */
.overflow-y-auto::-webkit-scrollbar {
  width: 6px;
}

.overflow-y-auto::-webkit-scrollbar-track {
  background: transparent;
}

.overflow-y-auto::-webkit-scrollbar-thumb {
  background: rgba(156, 163, 175, 0.3);
  border-radius: 3px;
}

.overflow-y-auto::-webkit-scrollbar-thumb:hover {
  background: rgba(156, 163, 175, 0.5);
}

/* 头像错误处理 */
.fallback-avatar {
  background: linear-gradient(135deg, #dbeafe 0%, #e0e7ff 100%);
}
</style>
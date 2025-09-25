<template>
  <div class="markdown-content" v-html="renderedMarkdown"></div>
</template>

<script setup>
import { computed } from 'vue'
import { marked } from 'marked'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'

// 配置marked
marked.setOptions({
  highlight: function(code, lang) {
    if (lang && hljs.getLanguage(lang)) {
      try {
        return hljs.highlight(code, { language: lang }).value
      } catch (err) {
        console.warn('Code highlighting error:', err)
      }
    }
    try {
      return hljs.highlightAuto(code).value
    } catch (err) {
      console.warn('Auto code highlighting error:', err)
      return code
    }
  },
  breaks: true, // 支持换行
  gfm: true, // GitHub风格markdown
  pedantic: false, // 不严格模式
  sanitize: false, // 允许HTML（但要注意XSS）
  smartLists: true, // 智能列表
  smartypants: true, // 智能标点
})

const props = defineProps({
  content: {
    type: String,
    required: true
  }
})

const renderedMarkdown = computed(() => {
  if (!props.content) return ''
  return marked(props.content)
})
</script>

<style scoped>
.markdown-content {
  @apply max-w-none;
}

.markdown-content :deep(h1) {
  @apply text-xl font-bold text-gray-900 mb-4 mt-6;
}

.markdown-content :deep(h2) {
  @apply text-lg font-semibold text-gray-800 mb-3 mt-5;
}

.markdown-content :deep(h3) {
  @apply text-base font-semibold text-gray-700 mb-2 mt-4;
}

.markdown-content :deep(h4) {
  @apply text-sm font-semibold text-gray-700 mb-2 mt-3;
}

.markdown-content :deep(p) {
  @apply text-gray-700 mb-3 leading-relaxed;
}

.markdown-content :deep(strong) {
  @apply font-semibold text-gray-900;
}

.markdown-content :deep(em) {
  @apply italic text-gray-800;
}

.markdown-content :deep(code) {
  @apply bg-gray-100 text-gray-800 px-1.5 py-0.5 rounded text-sm font-mono;
}

.markdown-content :deep(pre) {
  @apply bg-gray-900 text-gray-100 p-4 rounded-lg overflow-x-auto my-4;
}

.markdown-content :deep(pre code) {
  @apply bg-transparent text-gray-100 p-0;
}

.markdown-content :deep(blockquote) {
  @apply border-l-4 border-blue-500 pl-4 py-2 bg-blue-50 text-gray-700 italic my-4;
}

.markdown-content :deep(ul) {
  @apply list-disc list-inside mb-3 space-y-1;
}

.markdown-content :deep(ol) {
  @apply list-decimal list-inside mb-3 space-y-1;
}

.markdown-content :deep(li) {
  @apply text-gray-700;
}

.markdown-content :deep(a) {
  @apply text-blue-600 hover:text-blue-800 underline;
}

.markdown-content :deep(table) {
  @apply w-full border-collapse border border-gray-300 my-4;
}

.markdown-content :deep(th) {
  @apply border border-gray-300 bg-gray-100 px-3 py-2 text-left font-semibold;
}

.markdown-content :deep(td) {
  @apply border border-gray-300 px-3 py-2;
}

.markdown-content :deep(hr) {
  @apply border-t border-gray-300 my-6;
}

.markdown-content :deep(img) {
  @apply max-w-full h-auto rounded-lg shadow-sm my-4;
}

.markdown-content :deep(del) {
  @apply line-through text-gray-500;
}

.markdown-content :deep(ins) {
  @apply underline;
}

.markdown-content :deep(mark) {
  @apply bg-yellow-200 px-1 rounded;
}

.markdown-content :deep(small) {
  @apply text-xs text-gray-500;
}

.markdown-content :deep(sub) {
  @apply text-xs;
}

.markdown-content :deep(sup) {
  @apply text-xs;
}

/* 代码高亮样式调整 */
.markdown-content :deep(.hljs) {
  @apply bg-gray-900 text-gray-100;
}

.markdown-content :deep(.hljs-comment) {
  @apply text-gray-500;
}

.markdown-content :deep(.hljs-keyword) {
  @apply text-blue-400;
}

.markdown-content :deep(.hljs-string) {
  @apply text-green-400;
}

.markdown-content :deep(.hljs-number) {
  @apply text-yellow-400;
}

.markdown-content :deep(.hljs-function) {
  @apply text-purple-400;
}

.markdown-content :deep(.hljs-variable) {
  @apply text-orange-400;
}
</style>

<template>
  <div class="chat-container" :data-mode="currentMode">
    <!-- 左侧功能选择区 -->
    <div class="feature-sidebar">
      <div class="feature-options">
        <div
          v-for="(option, index) in features"
          :key="index"
          :class="['feature-item', { active: currentMode === option.mode }]"
          @click="switchMode(option.mode)"
        >
          <el-icon>
            <component :is="option.icon" />
          </el-icon>
          {{ option.label }}
        </div>
      </div>
    </div>

    <!-- 右侧主对话区 -->
    <el-card class="chat-card" :style="{ backgroundColor: currentBackground }">
      <div class="chat-messages" ref="messagesContainer">
        <template v-if="chatStore.messages.length === 0">
          <div class="welcome-message">
            <el-empty description="开始一段新的对话吧！" />
          </div>
        </template>
        <template v-else>
          <div
            v-for="(message, index) in chatStore.messages"
            :key="index"
            :class="['message', message.role]"
          >
            <div class="message-avatar">
              <el-avatar 
                :size="40"
                :src="message.role === 'user' ? '/avatar-user.png' : '/avatar-ai.png'"
                :icon="message.role === 'user' ? 'User' : 'Service'" 
              />
            </div>
            <div class="message-content">
              <div class="message-text markdown-body" v-html="renderMarkdown(message.content)"></div>
              <div class="message-time">
                {{ new Date(message.timestamp).toLocaleTimeString() }}
              </div>
            </div>
          </div>
        </template>
        <div v-if="chatStore.loading" class="loading-indicator">
          <el-icon class="is-loading"><Loading /></el-icon>
          AI正在思考中...
        </div>
      </div>
      <div class="chat-input">
        <el-input
          v-model="messageInput"
          type="textarea"
          :rows="3"
          :placeholder="inputPlaceholder"
          :disabled="chatStore.loading"
          @keydown.enter.exact.prevent="sendMessage"
        >
          <template #append>
            <el-button
              type="primary"
              :icon="currentFeature.icon"
              :loading="chatStore.loading"
              @click="sendMessage"
            >
              {{ currentFeature.buttonText }}
            </el-button>
          </template>
        </el-input>
        <div class="input-tips">
          按 Enter 发送，Shift + Enter 换行
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick, watch } from 'vue'
import { 
  ChatLineRound, 
  Picture, 
  Microphone,
  Loading
} from '@element-plus/icons-vue'
import { useChatStore } from '../stores/chat'
import { ElMessage } from 'element-plus'
import { marked } from 'marked'
import DOMPurify from 'dompurify' // 用于防止XSS攻击

const chatStore = useChatStore()
const messageInput = ref('')
const messagesContainer = ref(null)

// 添加markdown渲染函数
const renderMarkdown = (content) => {
  const html = marked(content)
  return DOMPurify.sanitize(html)
}

const scrollToBottom = async () => {
  await nextTick()
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

const isImageMode = ref(false)

const features = [
  { 
    mode: 'chat',
    label: '智能对话',
    icon: 'ChatLineRound',
    buttonText: '发送',
    placeholder: '输入您的问题...',
    background: '#fff8e6'  // 淡雅米黄色背景
  },
  { 
    mode: 'image',
    label: '图像生成',
    icon: 'Picture',
    buttonText: '生成图片',
    placeholder: '请描述要生成的图片...',
    background: '#f5e6ff'  // 淡紫色背景
  },
  { 
    mode: 'voice',
    label: '语音合成',
    icon: 'Microphone',
    buttonText: '生成语音',
    placeholder: '输入要转换的文字...',
    background: '#e6fff0'  // 淡绿色背景
  }
]

const currentMode = ref('chat')
const currentFeature = computed(() => features.find(f => f.mode === currentMode.value))
const currentBackground = computed(() => currentFeature.value.background)
const inputPlaceholder = computed(() => currentFeature.value.placeholder)

const switchMode = (mode) => {
  currentMode.value = mode
}

const sendMessage = async () => {
  const message = messageInput.value.trim()
  if (!message || chatStore.loading) return

  try {
    messageInput.value = ''
    switch (currentMode.value) {
      case 'image':
        await chatStore.generateImage(message)
        break
      case 'voice':
        await chatStore.generateVoice(message)
        break
      default:
        await chatStore.sendMessage(message)
    }
  } catch (error) {
    ElMessage.error('操作失败，请重试')
  }
}

watch(
  () => chatStore.messages,
  () => scrollToBottom(),
  { deep: true }
)

onMounted(() => {
  scrollToBottom()
})
</script>

<style scoped>
.chat-container {
  height: calc(100vh - 80px);
  padding: 20px;
  padding-top: 0;
  box-sizing: border-box;
  display: flex;
  gap: 20px;
}

.feature-sidebar {
  width: 200px;
  background: #fff;
  border-radius: 8px;
  padding: 20px 0;
  box-shadow: 0 2px 12px 0 rgba(0,0,0,0.1);
}

.feature-item {
  padding: 16px 24px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: all 0.3s ease;
}

.feature-item:hover {
  background: #f5f7fa;
}

.feature-item.active {
  background: #40ffb6;
  color: #fff;
}

.chat-card {
  flex: 1;
  height: 100%;
  transition: background-color 0.3s ease;
}

.chat-input-wrapper {
  position: sticky;
  bottom: 0;
  background: #fff;
  z-index: 10;
  box-shadow: 0 -2px 12px rgba(0,0,0,0.05);
}

.chat-input {
  padding: 20px;
  border-top: 1px solid #dcdfe6;
}


.chat-messages {
  flex: 1 1 0; /* 关键：允许收缩/扩展 */
  min-height:440px;
  max-height: 440px; /* 修复Flex容器高度坍塌 */
  overflow-y: auto;
  padding: 20px;
}

/* 自定义滚动条样式 */
.chat-messages::-webkit-scrollbar {
  width: 6px;
}

.chat-messages::-webkit-scrollbar-thumb {
  background-color: #909399;
  border-radius: 3px;
}

.chat-messages::-webkit-scrollbar-track {
  background-color: #f5f7fa;
}

/* Markdown 样式 */
:deep(.message-text) {
  line-height: 1.4; /* 原值通常为1.5-1.8 */
  margin: 8px 0; /* 减少段落间距 */
}

/* 清除Markdown默认间距 */
:deep(.message-text p) {
  margin: 4px 0;
}

/* 代码块间距优化 */
:deep(.message-text pre) {
  margin: 8px 0;
}

:deep(.message-text code) {
  font-family: monospace;
  padding: 2px 4px;
  background-color: rgba(0, 0, 0, 0.1);
  border-radius: 4px;
}

:deep(.user .message-text pre),
:deep(.user .message-text code) {
  background-color: rgba(255, 255, 255, 0.1);
}

.welcome-message {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.message {
  display: flex;
  margin-bottom: 24px;
  gap: 16px;
  position: relative;
  z-index: 2;
  max-width: 85%;
}

.message.assistant {
  flex-direction: row;
  margin-right: auto;
}

.message.user {
  flex-direction: row-reverse;
  margin-left: auto;
}
.message-time {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
  text-align: right;
}
.message-content {
  background: #f5f7fa !important;  /* 固定为浅灰色背景 */
  padding: 12px 16px;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
  position: relative;
  display: inline-block;
  max-width: 100%;
  border: 1px solid #e4e7ed !important;  /* 固定边框颜色 */
}

/* 移除之前的动态背景色设置 */
.message.assistant .message-content,
.message.user .message-content {
  background: #f5f7fa !important;
  border: 1px solid #e4e7ed !important;
}

/* 修改小三角的颜色为固定的灰色 */
.message.assistant .message-content::before {
  border-color: transparent #f5f7fa transparent transparent !important;
}

.message.user .message-content::before {
  content: '';
  position: absolute;
  right: -8px;
  top: 15px;
  border-style: solid;
  border-width: 8px 0 8px 8px;
  border-color: transparent transparent transparent #ecf5ff;
}

.feature-sidebar {
  width: 220px;
  background: #fff;
  border-radius: 12px;
  padding: 24px 0;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.feature-item {
  padding: 16px 24px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 12px;
  transition: all 0.3s ease;
  border-left: 4px solid transparent;
}

.feature-item:hover {
  background: #f5f7fa;
  border-left-color: #a0cfff;
}

.feature-item.active {
  background: linear-gradient(to right, #ecf5ff, #ffffff);
  border-left-color: #409eff;
  color: #409eff;
}

.chat-card {
  flex: 1;
  height: 100%;
  transition: all 0.3s ease;
  border-radius: 12px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.chat-input {
  padding: 24px;
  border-top: 1px solid #ebeef5;
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(10px);
  border-radius: 0 0 12px 12px;
}

:deep(.el-input__wrapper) {
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

:deep(.el-button) {
  border-radius: 8px;
  padding: 12px 24px;
}

/* Markdown 内容样式优化 */
:deep(.message-text) {
  line-height: 1.6;
  color: #2c3e50;
}

:deep(.message-text pre) {
  background: rgba(0, 0, 0, 0.02);
  border-radius: 8px;
  padding: 16px;
  margin: 12px 0;
}

:deep(.message-text code) {
  background: rgba(0, 0, 0, 0.04);
  padding: 2px 6px;
  border-radius: 4px;
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', monospace;
}

/* 加载动画美化 */
.loading-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  color: #909399;
  margin: 20px 0;
}


.input-tips {
  font-size: 12px;
  color: #909399;
  margin-top: 8px;
  text-align: right;
}

/* 为图片消息添加样式 */
:deep(.message-text img) {
  max-width: 40%;
  border-radius: 4px;
  background-color: #f6f8fa;
  overflow-x: auto;
}

:deep(.markdown-body code) {
  font-family: monospace;
  padding: 2px 4px;
  background-color: rgba(0, 0, 0, 0.1);
  border-radius: 4px;
  word-break: break-all;
  white-space: pre-wrap;
}

:deep(.user .markdown-body pre),
:deep(.user .markdown-body code) {
  background-color: rgba(255, 255, 255, 0.1);
}

.chat-messages::-webkit-scrollbar-track {
  background-color: transparent;
}

/* 欢迎消息美化 */
.welcome-message {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #909399;
  font-size: 16px;
}

:deep(.el-empty__description) {
  margin-top: 16px;
  font-size: 16px;
  color: #909399;
}
</style>

<style scoped>
.chat-container {
  height: calc(100vh - 80px);
  padding: 20px;
  padding-top: 0;
  box-sizing: border-box;
  display: flex;
  gap: 20px;
}

.feature-sidebar {
  width: 200px;
  background: #fff;
  border-radius: 8px;
  padding: 20px 0;
  box-shadow: 0 2px 12px 0 rgba(0,0,0,0.1);
}

.feature-item {
  padding: 16px 24px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: all 0.3s ease;
}

.feature-item:hover {
  background: #f5f7fa;
}

.feature-item.active {
  background: #40ffb6;
  color: #fff;
}

.chat-card {
  flex: 1;
  height: 100%;
  transition: background-color 0.3s ease;
}

.chat-input-wrapper {
  position: sticky;
  bottom: 0;
  background: #fff;
  z-index: 10;
  box-shadow: 0 -2px 12px rgba(0,0,0,0.05);
}

.chat-input {
  padding: 20px;
  border-top: 1px solid #dcdfe6;
}


.chat-messages {
  flex: 1 1 0; /* 关键：允许收缩/扩展 */
  min-height:470px;
  max-height: 470px; /* 修复Flex容器高度坍塌 */
  overflow-y: auto;
  padding: 20px;
}

/* 自定义滚动条样式 */
.chat-messages::-webkit-scrollbar {
  width: 6px;
}

.chat-messages::-webkit-scrollbar-thumb {
  background-color: #909399;
  border-radius: 3px;
}

.chat-messages::-webkit-scrollbar-track {
  background-color: #f5f7fa;
}

/* Markdown 样式 */
:deep(.message-text) {
  line-height: 1.4; /* 原值通常为1.5-1.8 */
  margin: 8px 0; /* 减少段落间距 */
}

/* 清除Markdown默认间距 */
:deep(.message-text p) {
  margin: 4px 0;
}

/* 代码块间距优化 */
:deep(.message-text pre) {
  margin: 8px 0;
}

:deep(.message-text code) {
  font-family: monospace;
  padding: 2px 4px;
  background-color: rgba(0, 0, 0, 0.1);
  border-radius: 4px;
}

:deep(.user .message-text pre),
:deep(.user .message-text code) {
  background-color: rgba(255, 255, 255, 0.1);
}

.welcome-message {
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.message {
  display: flex;
  margin-bottom: 24px;
  gap: 16px;
  position: relative;
  z-index: 2;
  max-width: 85%;
}

.message.assistant {
  flex-direction: row;
  margin-right: auto;
}

.message.user {
  flex-direction: row-reverse;
  margin-left: auto;
}
.message-time {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
  text-align: right;
}
.message-content {
  background: #fff;
  padding: 12px 16px;
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
  position: relative;
  display: inline-block; /* 改为行内块级元素 */
  max-width: 100%;
}
/* 修改图片消息的特殊样式 */
:deep(.message-text img) {
  width: 40%;
  border-radius: 4px;
  background-color: #f6f8fa;
  display: block; /* 让图片独占一行 */
}
.message.assistant .message-content {
  background: var(--message-bg, #fff8e6);
  border: 1px solid var(--message-border, #ffd591);
}

.message.user .message-content {
  background: var(--message-bg, #fff8e6);
  border: 1px solid var(--message-border, #ffd591);
}

/* 添加小三角颜色适配 */
.message.assistant .message-content::before {
  border-color: transparent var(--message-bg, #fff8e6) transparent transparent;
}

.message.user .message-content::before {
  border-color: transparent transparent transparent var(--message-bg, #fff8e6);
}

/* 不同模式下的颜色变化 */
.chat-container[data-mode='chat'] {
  --message-bg: #fff8e6;
  --message-border: #ffd591;
}

.chat-container[data-mode='image'] {
  --message-bg: #f5e6ff;
  --message-border: #d3adf7;
}

.chat-container[data-mode='voice'] {
  --message-bg: #e6fff0;
  --message-border: #b7eb8f;
}

.feature-item.active {
  background: var(--active-bg, #ffd591);
  color: #fff;
}

/* 不同模式下的激活颜色 */
.chat-container[data-mode='chat'] .feature-item.active {
  --active-bg: #ffa940;
}

.chat-container[data-mode='image'] .feature-item.active {
  --active-bg: #9254de;
}

.chat-container[data-mode='voice'] .feature-item.active {
  --active-bg: #52c41a;
}
</style>
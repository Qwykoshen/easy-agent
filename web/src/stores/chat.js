import { defineStore } from 'pinia'
import { ref } from 'vue'
import axios from 'axios'
import { useUserStore } from './user'

export const useChatStore = defineStore('chat', () => {
  const messages = ref([])
  const loading = ref(false)
  const userStore = useUserStore()

  const addMessage = (role, content) => {
    messages.value.push({
      role,
      content,
      timestamp: new Date().toISOString()
    })
  }

  const sendMessage = async (content) => {
    if (!content.trim()) return

    try {
      loading.value = true
      addMessage('user', content)

      const response = await axios.post('/api/v1/chat', {
        messages: messages.value.map(msg => ({
          role: msg.role,
          content: msg.content
        })),
        model: 'deepseek-chat',
        temperature: 0.7,
        max_tokens: 2000
      }, {
        headers: {
          'Authorization': `Bearer ${userStore.token}`
        }
      })

      addMessage('assistant', response.data.choices[0].message.content)
    } catch (error) {
      console.error('发送消息失败:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  const clearMessages = () => {
    messages.value = []
  }

  const generateImage = async (prompt) => {
    if (!prompt.trim()) return

    try {
      loading.value = true
      addMessage('user', `生成图片：${prompt}`)

      const response = await axios.post('/api/v1/generate-image', {
        prompt
      }, {
        headers: {
          'Authorization': `Bearer ${userStore.token}`
        }
      })

      // 直接获取图片URL并显示
      addMessage('assistant', `![生成的图片](${response.data.image_url})`)
    } catch (error) {
      console.error('生成图片失败:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  const generateVoice = async (text) => {
    if (!text.trim()) return

    try {
      loading.value = true
      addMessage('user', `生成语音：${text}`)

      const response = await axios.post('/api/v1/generate-voice', {
        text
      }, {
        headers: {
          'Authorization': `Bearer ${userStore.token}`
        }
      })

      // 添加音频消息
      addMessage('assistant', `
<audio controls>
  <source src="${response.data.audio_url}" type="audio/mp3">
  您的浏览器不支持音频播放
</audio>
      `)
    } catch (error) {
      console.error('生成语音失败:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  return {
    messages,
    loading,
    sendMessage,
    clearMessages,
    generateImage,
    generateVoice
  }
})
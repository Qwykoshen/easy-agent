import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import axios from 'axios'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const username = ref(localStorage.getItem('username') || '')

  const isLoggedIn = computed(() => !!token.value)

  const setToken = (newToken) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  const setUsername = (name) => {
    username.value = name
    localStorage.setItem('username', name)
  }

  const login = async (credentials) => {
    try {
      const response = await axios.post('/api/v1/login', credentials)
      setToken(response.data.token)
      setUsername(credentials.username)
      return response
    } catch (error) {
      throw error
    }
  }

  const register = async (userData) => {
    try {
      const response = await axios.post('/api/v1/register', userData)
      setToken(response.data.token)
      setUsername(userData.username)
      return response
    } catch (error) {
      throw error
    }
  }

  const logout = () => {
    token.value = ''
    username.value = ''
    localStorage.removeItem('token')
    localStorage.removeItem('username')
  }

  return {
    token,
    username,
    isLoggedIn,
    login,
    register,
    logout
  }
})
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import type { User } from '@/types/user'
import type { LoginRequest, RegisterRequest } from '@/types/auth'
import * as authApi from '@/api/auth'
import router from '@/router'

export const useAuthStore = defineStore('auth', () => {
  // 状态
  const token = ref(localStorage.getItem('token') || '')
  const user = ref<User | null>(null)
  const loading = ref(false)

  // 登录
  async function login(data: LoginRequest) {
    loading.value = true
    try {
      const res = await authApi.login(data)
      token.value = res.token
      user.value = res.user
      localStorage.setItem('token', res.token)
      ElMessage.success('登录成功')
      router.push('/')
      return true
    } catch (error) {
      return false
    } finally {
      loading.value = false
    }
  }

  // 注册
  async function register(data: RegisterRequest) {
    loading.value = true
    try {
      const res = await authApi.register(data)
      token.value = res.token
      user.value = res.user
      localStorage.setItem('token', res.token)
      ElMessage.success('注册成功')
      router.push('/')
      return true
    } catch (error) {
      return false
    } finally {
      loading.value = false
    }
  }

  // 登出
  async function logout() {
    try {
      await authApi.logout()
      token.value = ''
      user.value = null
      localStorage.removeItem('token')
      router.push('/login')
      return true
    } catch (error) {
      return false
    }
  }

  // 获取用户信息
  async function getInfo() {
    try {
      const res = await authApi.getCurrentUser()
      user.value = res
      return true
    } catch (error) {
      token.value = ''
      user.value = null
      localStorage.removeItem('token')
      return false
    }
  }

  // 检查登录状态
  async function checkAuth() {
    if (!token.value) return false
    if (!user.value) {
      return await getInfo()
    }
    return true
  }

  return {
    // 状态
    token,
    user,
    loading,

    // 方法
    login,
    register,
    logout,
    getInfo,
    checkAuth
  }
}) 
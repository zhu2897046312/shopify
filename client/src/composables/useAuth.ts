import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'
import type { LoginRequest, RegisterRequest } from '@/types/auth'

export function useAuth() {
  const router = useRouter()
  const authStore = useAuthStore()
  const loading = ref(false)

  // 登录
  async function handleLogin(email: string, password: string) {
    loading.value = true
    try {
      const success = await authStore.login({ email, password })
      if (!success) {
        ElMessage.error('登录失败，请检查邮箱和密码')
      }
      return success
    } finally {
      loading.value = false
    }
  }

  // 注册
  async function handleRegister(email: string, password: string, nickname?: string) {
    loading.value = true
    try {
      const success = await authStore.register({ 
        email, 
        password, 
        nickname,
        role: 'admin'
      })
      if (!success) {
        ElMessage.error('注册失败，请稍后重试')
      }
      return success
    } finally {
      loading.value = false
    }
  }

  // 登出
  async function handleLogout() {
    const success = await authStore.logout()
    if (!success) {
      ElMessage.error('登出失败，请稍后重试')
    }
    return success
  }

  // 检查登录状态
  async function checkAuth() {
    return await authStore.checkAuth()
  }

  return {
    loading,
    handleLogin,
    handleRegister,
    handleLogout,
    checkAuth
  }
}

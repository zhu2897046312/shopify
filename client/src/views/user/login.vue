<template>
  <div class="login-container">
    <el-card class="login-card">
      <template #header>
        <h2>管理员登录</h2>
      </template>
      <el-form :model="loginForm" :rules="rules" ref="loginFormRef">
        <el-form-item prop="email">
          <el-input v-model="loginForm.email" placeholder="邮箱">
            <template #prefix>
              <el-icon><User /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item prop="password">
          <el-input v-model="loginForm.password" type="password" placeholder="密码">
            <template #prefix>
              <el-icon><Lock /></el-icon>
            </template>
          </el-input>
        </el-form-item>
        <div class="form-actions">
          <el-checkbox v-model="rememberMe">记住我</el-checkbox>
          <el-link type="primary" @click="router.push('/forgot-password')">忘记密码？</el-link>
        </div>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleSubmit" block>
            登录
          </el-button>
        </el-form-item>
        <div class="register-link">
          没有账号？<el-link type="primary" @click="router.push('/register')">去注册</el-link>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { User, Lock } from '@element-plus/icons-vue'
import type { FormInstance } from 'element-plus'
import { useAuth } from '@/composables/useAuth'

const router = useRouter()
const { loading, handleLogin } = useAuth()
const loginFormRef = ref<FormInstance>()
const rememberMe = ref(false)

const loginForm = reactive({
  email: '',
  password: ''
})

const rules = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能小于6位', trigger: 'blur' }
  ]
}

const handleSubmit = async () => {
  if (!loginFormRef.value) return
  
  await loginFormRef.value.validate(async (valid) => {
    if (valid) {
      // 如果记住我，保存邮箱到本地存储
      if (rememberMe.value) {
        localStorage.setItem('remembered_email', loginForm.email)
      } else {
        localStorage.removeItem('remembered_email')
      }
      
      const success = await handleLogin(loginForm.email, loginForm.password)
      if (success) {
        router.push('/dashboard')
      }
    }
  })
}

// 初始化时，如果有记住的邮箱，自动填充
const rememberedEmail = localStorage.getItem('remembered_email')
if (rememberedEmail) {
  loginForm.email = rememberedEmail
  rememberMe.value = true
}
</script>

<style scoped lang="scss">
.login-container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f5f7fa;
  
  .login-card {
    width: 400px;
    
    :deep(.el-card__header) {
      text-align: center;
      
      h2 {
        margin: 0;
      }
    }
  }

  .form-actions {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
  }

  .register-link {
    text-align: right;
    margin-top: 10px;
  }
}
</style> 
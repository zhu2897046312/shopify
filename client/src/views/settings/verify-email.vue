<template>
  <div class="verify-email-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>邮箱验证</span>
        </div>
      </template>

      <div v-if="!sent" class="send-section">
        <p class="tip">我们将向您的邮箱 {{ email }} 发送验证邮件</p>
        <el-button 
          type="primary" 
          :loading="loading" 
          @click="handleSendVerification"
        >
          发送验证邮件
        </el-button>
      </div>

      <div v-else class="verify-section">
        <el-result
          icon="success"
          title="验证邮件已发送"
          sub-title="请查看您的邮箱，点击邮件中的验证链接完成验证"
        >
          <template #extra>
            <el-button type="primary" @click="$router.push('/settings/profile')">
              返回个人信息
            </el-button>
            <el-button @click="handleResend">重新发送</el-button>
          </template>
        </el-result>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { ElMessage } from 'element-plus'

const authStore = useAuthStore()
const loading = ref(false)
const sent = ref(false)

const email = computed(() => authStore.user?.email)

const handleSendVerification = async () => {
  loading.value = true
  try {
    await authStore.sendVerificationEmail()
    ElMessage.success('验证邮件已发送')
    sent.value = true
  } catch (error: any) {
    ElMessage.error(error.message || '发送失败')
  } finally {
    loading.value = false
  }
}

const handleResend = () => {
  sent.value = false
}
</script>

<style scoped lang="scss">
.verify-email-container {
  max-width: 600px;
  margin: 0 auto;

  .send-section {
    text-align: center;
    padding: 40px 0;

    .tip {
      margin-bottom: 20px;
      color: #666;
    }
  }

  .verify-section {
    padding: 20px 0;
  }
}
</style> 
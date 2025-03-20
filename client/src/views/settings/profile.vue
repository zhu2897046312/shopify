<template>
  <div class="profile-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>个人信息</span>
        </div>
      </template>

      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="头像">
          <el-upload
            class="avatar-uploader"
            action="/api/v1/upload"
            :show-file-list="false"
            :on-success="handleAvatarSuccess"
          >
            <el-avatar v-if="form.avatar" :src="form.avatar" :size="100" />
            <el-icon v-else class="avatar-uploader-icon"><Plus /></el-icon>
          </el-upload>
        </el-form-item>

        <el-form-item label="邮箱">
          <el-input v-model="form.email" disabled />
          <template #append>
            <el-tag v-if="form.email_verified" type="success">已验证</el-tag>
            <el-button v-else type="primary" link @click="handleVerifyEmail">
              去验证
            </el-button>
          </template>
        </el-form-item>

        <el-form-item label="昵称" prop="nickname">
          <el-input v-model="form.nickname" />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleSubmit">
            保存
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import type { FormInstance } from 'element-plus'
import { useAuthStore } from '@/stores/auth'
import { ElMessage } from 'element-plus'

const authStore = useAuthStore()
const formRef = ref<FormInstance>()
const loading = ref(false)

const form = reactive({
  email: authStore.user?.email || '',
  nickname: authStore.user?.nickname || '',
  avatar: authStore.user?.avatar || '',
  email_verified: authStore.user?.email_verified || false
})

const rules = {
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' }
  ]
}

const handleAvatarSuccess = (res: any) => {
  form.avatar = res.data.url
}

const handleVerifyEmail = () => {
  ElMessage.success('验证邮件已发送，请查收')
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        await authStore.updateProfile({
          nickname: form.nickname,
          avatar: form.avatar
        })
        ElMessage.success('保存成功')
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped lang="scss">
.profile-container {
  .avatar-uploader {
    :deep(.el-upload) {
      border: 1px dashed #d9d9d9;
      border-radius: 50%;
      cursor: pointer;
      position: relative;
      overflow: hidden;
      transition: var(--el-transition-duration-fast);
      
      &:hover {
        border-color: var(--el-color-primary);
      }
    }
    
    .avatar-uploader-icon {
      font-size: 28px;
      color: #8c939d;
      width: 100px;
      height: 100px;
      line-height: 100px;
      text-align: center;
    }
  }
}
</style> 
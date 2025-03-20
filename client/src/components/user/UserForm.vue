<template>
  <el-form
    ref="formRef"
    :model="form"
    :rules="rules"
    label-width="100px"
  >
    <el-form-item label="邮箱" prop="email">
      <el-input v-model="form.email" :disabled="!!user" />
    </el-form-item>

    <el-form-item label="昵称" prop="nickname">
      <el-input v-model="form.nickname" />
    </el-form-item>

    <!-- <el-form-item label="头像">
      <el-upload
        class="avatar-uploader"
        action="/api/v1/upload"
        :show-file-list="false"
        :on-success="handleAvatarSuccess"
      >
        <el-avatar v-if="form.avatar" :src="form.avatar" :size="100" />
        <el-icon v-else class="avatar-uploader-icon"><Plus /></el-icon>
      </el-upload>
    </el-form-item> -->

    <el-form-item label="角色" prop="role">
      <el-select v-model="form.role">
        <el-option label="管理员" value="admin" />
        <el-option label="用户" value="user" />
      </el-select>
    </el-form-item>

    <el-form-item label="状态" prop="status">
      <el-select v-model="form.status">
        <el-option label="正常" value="active" />
        <el-option label="禁用" value="inactive" />
      </el-select>
    </el-form-item>

    <el-form-item>
      <el-button type="primary" :loading="loading" @click="handleSubmit">
        保存
      </el-button>
      <el-button @click="$emit('cancel')">取消</el-button>
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import type { FormInstance } from 'element-plus'
import type { User } from '@/types/user'
import { UserRole, UserStatus } from '@/types/user'

const props = defineProps<{
  user?: User
}>()

const emit = defineEmits<{
  (e: 'success', user: User): void
  (e: 'cancel'): void
}>()

const formRef = ref<FormInstance>()
const loading = ref(false)

const form = reactive({
  email: props.user?.email || '',
  nickname: props.user?.nickname || '',
  avatar: props.user?.avatar || '',
  role: props.user?.role || UserRole.User,
  status: props.user?.status || UserStatus.Active,
  email_verified: props.user?.email_verified || false
})

const rules = {
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' }
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'change' }
  ],
  status: [
    { required: true, message: '请选择状态', trigger: 'change' }
  ]
}

const handleAvatarSuccess = (res: any) => {
  form.avatar = res.data.url
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        emit('success', {
          id: props.user?.id,
          ...form
        } as User)
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped lang="scss">
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
</style>
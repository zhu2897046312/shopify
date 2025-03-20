<template>
  <div class="user-page">
    <!-- 搜索栏 -->
    <el-card class="search-card">
      <el-form :model="searchForm" inline>
        <el-form-item label="关键词">
          <el-input
            v-model="searchForm.keyword"
            placeholder="用户名/邮箱"
            clearable
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item label="角色">
          <el-select v-model="searchForm.role" placeholder="全部" clearable>
            <el-option label="管理员" :value="UserRole.Admin" />
            <el-option label="用户" :value="UserRole.User" />
          </el-select>
        </el-form-item>
        <el-form-item label="邮箱验证">
          <el-select v-model="searchForm.email_verified" placeholder="全部" clearable>
            <el-option label="已验证" :value="true" />
            <el-option label="未验证" :value="false" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>搜索
          </el-button>
          <el-button @click="handleReset">
            <el-icon><Refresh /></el-icon>重置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 操作栏 -->
    <el-card class="action-card">
      <template #header>
        <div class="card-header">
          <span>用户列表</span>
          <div class="actions">
            <el-button type="primary" @click="handleAdd">
              <el-icon><Plus /></el-icon>新增用户
            </el-button>
            <el-button type="danger" :disabled="!selectedUsers.length" @click="handleBatchDelete">
              <el-icon><Delete /></el-icon>批量删除
            </el-button>
            <el-button @click="handleExport">
              <el-icon><Download /></el-icon>导出数据
            </el-button>
          </div>
        </div>
      </template>

      <!-- 用户表格 -->
      <user-table
        ref="tableRef"
        v-model:selected="selectedUsers"
        @edit="handleEdit"
        @refresh="loadUsers"
      />
    </el-card>

    <!-- 用户表单对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="currentUser ? '编辑用户' : '新增用户'"
      width="600px"
      destroy-on-close
    >
      <user-form
        :user="currentUser"
        @success="handleFormSuccess"
        @cancel="dialogVisible = false"
      />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessageBox, ElMessage } from 'element-plus'
import { Search, Refresh, Plus, Delete, Download } from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import { UserRole, UserStatus } from '@/types/user'
import type { User } from '@/types/user'
import UserTable from '@/components/user/UserTable.vue'
import UserForm from '@/components/user/UserForm.vue'

const userStore = useUserStore()
const tableRef = ref()
const dialogVisible = ref(false)
const currentUser = ref<User | undefined>(undefined)
const selectedUsers = ref<number[]>([])

// 搜索表单
const searchForm = reactive({
  keyword: '',
  role: undefined as UserRole | undefined,
  status: undefined as UserStatus | undefined,
  email_verified: undefined as boolean | undefined,
  page: 1,
  page_size: 10
})

// 加载用户列表
const loadUsers = async () => {
  await userStore.getUsers({
    page: searchForm.page,
    page_size: searchForm.page_size,
    keyword: searchForm.keyword,
    role: searchForm.role,
    status: searchForm.status,
    email_verified: searchForm.email_verified
  })
}

// 搜索
const handleSearch = () => {
  searchForm.page = 1
  loadUsers()
}

// 重置搜索
const handleReset = () => {
  Object.assign(searchForm, {
    keyword: '',
    role: undefined,
    status: undefined,
    email_verified: undefined,
    page: 1
  })
  loadUsers()
}

// 新增用户
const handleAdd = () => {
  currentUser.value = undefined
  dialogVisible.value = true
}

// 编辑用户
const handleEdit = (user: User) => {
  currentUser.value = user
  dialogVisible.value = true
}

// 批量删除
const handleBatchDelete = async () => {
  if (!selectedUsers.value.length) return
  
  try {
    await ElMessageBox.confirm('确定要删除选中的用户吗？', '提示', {
      type: 'warning'
    })
    await userStore.batchDeleteUsers(selectedUsers.value)
    ElMessage.success('删除成功')
    selectedUsers.value = []
    loadUsers()
  } catch {
    // 用户取消操作
  }
}

// 导出数据
const handleExport = async () => {
  try {
    await userStore.exportUsers(searchForm)
    ElMessage.success('导出成功')
  } catch (error: any) {
    ElMessage.error(error.message || '导出失败')
  }
}

// 表单提交成功
const handleFormSuccess = async (user: User) => {
  dialogVisible.value = false
  if (currentUser.value) {
    await userStore.updateUser(user.id, user)
  } else {
    await userStore.createUser(user)
  }
  ElMessage.success(`${currentUser.value ? '更新' : '创建'}成功`)
  loadUsers()
}

// 确保组件挂载时加载数据
onMounted(() => {
  loadUsers()
})
</script>

<style scoped lang="scss">
.user-page {
  .search-card {
    margin-bottom: 20px;
  }

  .action-card {
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
  }
}
</style> 
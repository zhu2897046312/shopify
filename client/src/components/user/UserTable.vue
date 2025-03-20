<template>
  <div class="user-table">
    <el-table
      v-loading="loading"
      :data="users"
      border
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column label="用户信息" min-width="200">
        <template #default="{ row }">
          <div class="user-info">
            <el-avatar :size="40" :src="row.avatar">
              {{ row.nickname?.charAt(0)?.toUpperCase() }}
            </el-avatar>
            <div class="user-detail">
              <div class="nickname">{{ row.nickname }}</div>
              <div class="email">{{ row.email }}</div>
            </div>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="角色" width="100">
        <template #default="{ row }">
          <el-tag :type="row.role === 'admin' ? 'danger' : 'info'">
            {{ row.role === 'admin' ? '管理员' : '用户' }}
          </el-tag>
        </template>
      </el-table-column>
      <!-- <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 'active' ? 'success' : 'info'">
            {{ row.status === 'active' ? '正常' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column> -->
      <el-table-column label="邮箱验证" width="100">
        <template #default="{ row }">
          <el-tag :type="row.email_verified ? 'success' : 'warning'">
            {{ row.email_verified ? '已验证' : '未验证' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="注册时间" width="180" />
      <el-table-column label="操作" width="250" fixed="right">
        <template #default="{ row }">
          <el-button-group>
            <el-button type="primary" link @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button type="primary" link @click="handleDetail(row)">
              详情
            </el-button>
            <el-button 
              type="primary" 
              link 
              @click="handleStatusChange(row)"
            >
              {{ row.status === 'active' ? '禁用' : '启用' }}
            </el-button>
            <el-button 
              v-if="!row.email_verified" 
              type="warning" 
              link 
              @click="handleVerifyEmail(row)"
            >
              验证邮箱
            </el-button>
            <el-button type="danger" link @click="handleDelete(row)">
              删除
            </el-button>
          </el-button-group>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessageBox } from 'element-plus'
import { useUserStore } from '@/stores/user'
import { storeToRefs } from 'pinia'
import type { User } from '@/types/user'
import { UserStatus } from '@/types/user'

const router = useRouter()
const userStore = useUserStore()
const { users, total, loading } = storeToRefs(userStore)

const currentPage = ref(1)
const pageSize = ref(10)
const selectedUsers = ref<number[]>([])

const emit = defineEmits<{
  (e: 'edit', user: User): void
  (e: 'refresh'): void
}>()

// 加载用户列表
async function loadUsers() {
  await userStore.getUsers({
    page: currentPage.value,
    page_size: pageSize.value
  })
}

// 选择变化
const handleSelectionChange = (selection: User[]) => {
  selectedUsers.value = selection.map(item => item.id)
}

// 编辑用户
const handleEdit = (user: User) => {
  emit('edit', user)
}

// 查看详情
const handleDetail = (user: User) => {
  router.push(`/users/${user.id}`)
}

// 更改状态
const handleStatusChange = async (user: User) => {
  try {
    await ElMessageBox.confirm(
      `确定要${user.status === UserStatus.Active ? '禁用' : '启用'}该用户吗？`,
      '提示',
      { type: 'warning' }
    )
    await userStore.updateUserStatus(user.id, 
      user.status === UserStatus.Active ? UserStatus.Inactive : UserStatus.Active
    )
    emit('refresh')
  } catch {
    // 用户取消操作
  }
}

// 验证邮箱
const handleVerifyEmail = async (user: User) => {
  console.log(user)
}

// 删除用户
const handleDelete = async (user: User) => {
  try {
    await ElMessageBox.confirm('确定要删除该用户吗？', '提示', {
      type: 'warning'
    })
    await userStore.deleteUser(user.id)
    emit('refresh')
  } catch {
    // 用户取消操作
  }
}

// 分页相关
const handleSizeChange = () => {
  currentPage.value = 1
  loadUsers()
}

const handleCurrentChange = () => {
  loadUsers()
}

// 初始加载
loadUsers()
</script>

<style scoped lang="scss">
.user-table {
  .user-info {
    display: flex;
    align-items: center;
    gap: 12px;

    .user-detail {
      .nickname {
        font-weight: bold;
      }
      .email {
        font-size: 12px;
        color: #666;
      }
    }
  }

  .pagination {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>
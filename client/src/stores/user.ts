import { defineStore } from 'pinia'
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import type { User, ListUserParams, ListUserResponse, UpdateUserRequest } from '@/types/user'
import { UserRole, UserStatus } from '@/types/user'
import * as userApi from '@/api/user'

// 添加类型定义
interface UserStats {
  total_count: number
  active_count: number
  inactive_count: number
  admin_count: number
  user_count: number
  verified_count: number
  today_count: number
  total_growth: number
}

export const useUserStore = defineStore('admin-user', () => {
  // 状态
  const users = ref<User[]>([])
  const total = ref(0)
  const loading = ref(false)
  const currentUser = ref<User | null>(null)
  const userStats = ref<UserStats>({
    total_count: 0,
    active_count: 0,
    inactive_count: 0,
    admin_count: 0,
    user_count: 0,
    verified_count: 0,
    today_count: 0,
    total_growth: 0
  })

  // 获取用户列表
  async function getUsers(params: ListUserParams): Promise<ListUserResponse> {
    loading.value = true
    try {
      const res = await userApi.getUsers(params)
      users.value = res.users
      total.value = res.total
      return res
    } catch (error) {
      console.error('Failed to get users:', error)
      throw error
    } finally {
      loading.value = false
    }
  }

  // 获取用户详情
  async function getUser(id: number) {
    try {
      const user = await userApi.getUser(id)
      currentUser.value = user
      return user
    } catch (error) {
      console.error('Failed to get user:', error)
      return null
    }
  }

  // 更新用户
  async function updateUser(id: number, data: UpdateUserRequest) {
    try {
      const user = await userApi.updateUser(id, data)
      // 更新列表中的用户数据
      const index = users.value.findIndex(u => u.id === id)
      if (index > -1) {
        users.value[index] = user
      }
      ElMessage.success('更新成功')
      return true
    } catch (error) {
      console.error('Failed to update user:', error)
      return false
    }
  }

  // 删除用户
  async function deleteUser(id: number) {
    try {
      await userApi.deleteUser(id)
      // 从列表中移除用户
      users.value = users.value.filter(u => u.id !== id)
      total.value--
      ElMessage.success('删除成功')
      return true
    } catch (error) {
      console.error('Failed to delete user:', error)
      return false
    }
  }

  // 批量删除用户
  async function batchDeleteUsers(ids: number[]) {
    try {
      await userApi.batchDeleteUsers(ids)
      // 从列表中移除用户
      users.value = users.value.filter(u => !ids.includes(u.id))
      total.value -= ids.length
      ElMessage.success('批量删除成功')
      return true
    } catch (error) {
      console.error('Failed to batch delete users:', error)
      return false
    }
  }

  // 更新用户状态
  async function updateUserStatus(id: number, status: string) {
    try {
      const user = await userApi.updateUserStatus(id, status)
      // 更新列表中的用户状态
      const index = users.value.findIndex(u => u.id === id)
      if (index > -1) {
        users.value[index] = user
      }
      ElMessage.success('状态更新成功')
      return true
    } catch (error) {
      console.error('Failed to update user status:', error)
      return false
    }
  }

  // 更新用户角色
  async function updateUserRole(id: number, role: string) {
    try {
      const user = await userApi.updateUserRole(id, role)
      // 更新列表中的用户角色
      const index = users.value.findIndex(u => u.id === id)
      if (index > -1) {
        users.value[index] = user
      }
      ElMessage.success('角色更新成功')
      return true
    } catch (error) {
      console.error('Failed to update user role:', error)
      return false
    }
  }

  // 验证用户邮箱
  async function verifyUserEmail(id: number) {
    try {
      const user = await userApi.verifyUserEmail(id)
      // 更新列表中的用户邮箱验证状态
      const index = users.value.findIndex(u => u.id === id)
      if (index > -1) {
        users.value[index] = user
      }
      ElMessage.success('邮箱验证成功')
      return true
    } catch (error) {
      console.error('Failed to verify user email:', error)
      return false
    }
  }

  // 重置用户密码
  async function resetUserPassword(id: number) {
    try {
      const res = await userApi.resetUserPassword(id)
      ElMessage.success('密码重置成功')
      return res.password
    } catch (error) {
      console.error('Failed to reset user password:', error)
      return null
    }
  }

  // 添加获取统计数据的方法
  const getUserStats = async (): Promise<UserStats> => {
    try {
      // 获取所有用户数据来计算统计信息
      const response = await getUsers({ page: 1, page_size: 1000 })
      const allUsers = response.users
      const now = new Date()
      const today = new Date(now.getFullYear(), now.getMonth(), now.getDate()).getTime()

      // 计算各种统计数据
      const stats: UserStats = {
        total_count: allUsers.length,
        admin_count: allUsers.filter(user => user.role === UserRole.Admin).length,
        user_count: allUsers.filter(user => user.role === UserRole.User).length,
        active_count: allUsers.filter(user => user.status === UserStatus.Active).length,
        inactive_count: allUsers.filter(user => user.status === UserStatus.Inactive).length,
        verified_count: allUsers.filter(user => user.email_verified).length,
        today_count: allUsers.filter(user => new Date(user.created_at).getTime() >= today).length,
        total_growth: calculateGrowthRate(allUsers.length, allUsers.filter(user => {
          const userDate = new Date(user.created_at).getTime()
          const yesterdayStart = today - 24 * 60 * 60 * 1000
          return userDate < today
        }).length)
      }

      // 更新store中的统计数据
      userStats.value = stats
      return stats
    } catch (error) {
      console.error('Failed to calculate user stats:', error)
      throw error
    }
  }

  // 计算增长率
  const calculateGrowthRate = (current: number, previous: number): number => {
    if (previous === 0) return current > 0 ? 100 : 0
    return Number(((current - previous) / previous * 100).toFixed(2))
  }

  // 导出用户数据
  const exportUsers = async (params: ListUserParams) => {
    return await userApi.exportUsers(params)
  }

  // 创建用户
  const createUser = async (data: User) => {
    return await userApi.createUser(data)
  }

  // 计算用户趋势数据
  async function getUserTrend(timeRange: string) {
    try {
      // 获取所有用户列表，包含创建时间
      const response: ListUserResponse = await getUsers({ page: 1, page_size: 1000 })
      const allUsers = response.users
      
      // 根据时间范围计算日期范围
      const now = new Date()
      let startDate = new Date()
      switch (timeRange) {
        case 'week':
          startDate.setDate(now.getDate() - 7)
          break
        case 'month':
          startDate.setMonth(now.getMonth() - 1)
          break
        case 'year':
          startDate.setFullYear(now.getFullYear() - 1)
          break
        default:
          startDate.setDate(now.getDate() - 7)
      }

      // 生成日期数组
      const dates: string[] = []
      const currentDate = new Date(startDate)
      while (currentDate <= now) {
        dates.push(currentDate.toISOString().split('T')[0])
        currentDate.setDate(currentDate.getDate() + 1)
      }

      // 计算每天的新增用户数
      const dailyCounts = dates.map(date => {
        const count = allUsers.filter((user: User) => {
          const userDate = new Date(user.created_at).toISOString().split('T')[0]
          return userDate === date
        }).length

        return {
          date,
          count,
          growth: 0 // 初始化增长率为0
        }
      })

      // 计算环比增长率
      for (let i = 1; i < dailyCounts.length; i++) {
        const prevCount = dailyCounts[i - 1].count
        const currentCount = dailyCounts[i].count
        if (prevCount > 0) {
          dailyCounts[i].growth = ((currentCount - prevCount) / prevCount) * 100
        }
      }

      return { data: dailyCounts }
    } catch (error) {
      console.error('Failed to calculate user trend:', error)
      return { data: [] }
    }
  }

  return {
    // 状态
    users,
    total,
    loading,
    currentUser,
    userStats,

    // 方法
    getUsers,
    getUser,
    updateUser,
    deleteUser,
    batchDeleteUsers,
    updateUserStatus,
    updateUserRole,
    verifyUserEmail,
    resetUserPassword,
    getUserStats,
    exportUsers,
    createUser,
    getUserTrend
  }
})
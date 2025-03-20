import request from '@/utils/request'
import type { 
  User, 
  ListUserParams, 
  ListUserResponse, 
  UpdateUserRequest,
  UserStatus,
  UserStats
} from '@/types/user'

// 获取用户列表
export function getUsers(params: ListUserParams): Promise<ListUserResponse> {
  return request.get('/admin/users', { params })
}

// 获取用户详情
export function getUser(id: number): Promise<User> {
  return request.get(`/admin/users/${id}`)
}

// 更新用户
export function updateUser(id: number, data: UpdateUserRequest): Promise<User> {
  return request.put(`/admin/users/${id}`, data)
}

// 删除用户
export function deleteUser(id: number) {
  return request.delete(`/admin/users/${id}`)
}

// 获取用户统计信息
export function getUserStats(): Promise<UserStats> {
  return request.get('/admin/users/stats')
}

// 更新用户状态
export function updateUserStatus(id: number, status: string): Promise<User> {
  return request.put(`/admin/users/${id}/status`, { status })
}

// 更新用户角色
export function updateUserRole(id: number, role: string): Promise<User> {
  return request.put(`/admin/users/${id}/role`, { role })
}

// 验证用户邮箱
export function verifyUserEmail(id: number): Promise<User> {
  return request.put(`/admin/users/${id}/verify-email`)
}

// 重置用户密码
export function resetUserPassword(id: number): Promise<{ password: string }> {
  return request.put(`/admin/users/${id}/reset-password`)
}

// 批量删除用户
export function batchDeleteUsers(ids: number[]) {
  return request.delete('/admin/users/batch', { data: { ids } })
}

// 导出用户数据
export function exportUsers(params: ListUserParams) {
  return request.get('/admin/users/export', { 
    params,
    responseType: 'blob'
  })
}

// 创建用户
export function createUser(data: User): Promise<User> {
  return request.post('users/register', data)
}
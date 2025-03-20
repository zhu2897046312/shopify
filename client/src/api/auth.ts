import request from '@/utils/request'
import type { LoginRequest, LoginResponse, RegisterRequest, RegisterResponse } from '@/types/auth'
import type { User } from '@/types/user'

// 登录
export function login(data: LoginRequest): Promise<LoginResponse> {
  return request.post('/users/login', data)
}

// 注册
export function register(data: RegisterRequest): Promise<RegisterResponse> {
  return request.post('/users/register', data)
}

// 登出
export function logout() {
  return request.post('/users/logout')
}

// 获取当前用户信息
export function getCurrentUser(): Promise<User> {
  return request.get('/users/profile')
}

// 修改密码
export function changePassword(oldPassword: string, newPassword: string) {
  return request.put('/users/password', {
    old_password: oldPassword,
    new_password: newPassword
  })
}

// 重置密码
export function resetPassword(email: string) {
  return request.post('/users/password/reset', { email })
}

// 验证邮箱
export function verifyEmail(token: string) {
  return request.post('/users/email/verify', { token })
}
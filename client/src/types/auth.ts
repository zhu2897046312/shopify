import { User } from './user'

// 登录请求参数
export interface LoginRequest {
  email: string
  password: string
}

// 登录响应
export interface LoginResponse {
  user: User
  token: string
}

// 注册请求参数
export interface RegisterRequest {
  email: string
  password: string
  nickname?: string
  role?: string
}

// 注册响应
export interface RegisterResponse {
  user: User
  token: string
} 
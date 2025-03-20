import { PaginationParams, PaginationResponse } from './base'

// 用户角色枚举
export enum UserRole {
  User = 'user',
  Admin = 'admin'
}

// 用户状态枚举
export enum UserStatus {
  Active = 'active',
  Inactive = 'inactive'
}

// 用户数据模型
export interface User {
  id: number
  email: string
  nickname: string
  avatar: string
  role: UserRole
  status: UserStatus
  email_verified: boolean
  created_at: string
  updated_at: string
}

// 地址数据模型
export interface Address {
  id: number
  user_id: number
  name: string
  phone: string
  province: string
  city: string
  district: string
  street: string
  post_code: string
  is_default: boolean
  created_at: string
  updated_at: string
}

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

// 用户列表查询参数
export interface ListUserParams extends PaginationParams {
  role?: UserRole
  status?: UserStatus
  email_verified?: boolean
  keyword?: string
}

// 用户列表响应
export interface ListUserResponse {
  page: number
  page_size: number
  total: number
  users: User[]
}

// 更新用户请求参数
export interface UpdateUserRequest {
  nickname?: string
  avatar?: string
  role?: UserRole
  status?: UserStatus
  email_verified?: boolean
}

// 设置默认地址请求参数
export interface SetDefaultAddressRequest {
  address_id: number
}

// 地址查询参数
export interface ListAddressParams {
  user_id?: number
  is_default?: boolean
}

// 地址验证响应
export interface AddressValidationResponse {
  valid: boolean
  message?: string
}

// 地址统计信息
export interface AddressStats {
  total_count: number
  default_count: number
  user_count: number  // 有地址的用户数
}

// 省市区数据结构
export interface Region {
  code: string
  name: string
  children?: Region[]
}

// 地址标签
export enum AddressTag {
  Home = 'home',
  Company = 'company',
  School = 'school',
  Other = 'other'
}

// 扩展地址接口(包含标签)
export interface AddressWithTags extends Address {
  tags?: AddressTag[]
}

// 注册请求参数
export interface RegisterRequest {
  email: string
  password: string
  nickname?: string
  role?: UserRole
}

// 注册响应
export interface RegisterResponse {
  user: User
  token: string
}

// 用户统计信息
export interface UserStats {
  total_count: number
  total_growth: number
  admin_count: number
  user_count: number
  verified_count: number
  today_count: number
  today_growth: number
  active_count: number
  inactive_count: number
}

// 用户趋势数据
export interface UserTrend {
  date: string
  count: number
  growth: number
}

export interface UserTrendResponse {
  data: UserTrend[]
}
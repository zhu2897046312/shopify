// 基础响应类型
export interface BaseResponse<T = any> {
  code: number
  message: string
  data: T
}

// 分页请求参数
export interface PaginationParams {
  page: number
  page_size: number
}

// 分页响应数据
export interface PaginationResponse<T> {
  items: T[]
  total: number
  page: number
  page_size: number
}

// 统计数据
export interface StatsData {
  total: number
  today: number
  week: number
  month: number
} 
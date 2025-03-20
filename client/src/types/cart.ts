import { Product } from './product'
import { User } from './user'
import { PaginationParams, PaginationResponse } from './base'

// 购物车项数据模型
export interface CartItem {
  id: number
  user_id: number
  user?: User
  product_id: number
  product: Product
  quantity: number
  selected: boolean
  created_at: string
  updated_at: string
}

// 添加购物车请求参数
export interface AddCartItemRequest {
  product_id: number
  quantity: number
}

// 更新购物车项数量请求参数
export interface UpdateCartItemQuantityRequest {
  quantity: number
}

// 更新购物车项选中状态请求参数
export interface UpdateCartItemSelectedRequest {
  selected: boolean
}

// 批量更新购物车项选中状态请求参数
export interface BatchUpdateSelectedRequest {
  selected: boolean
}

// 购物车列表查询参数
export interface ListCartParams extends PaginationParams {
  user_id?: number
  selected?: boolean
}

// 购物车列表响应
export type ListCartResponse = PaginationResponse<CartItem>

// 购物车统计信息
export interface CartStats {
  total_count: number        // 购物车商品总数
  selected_count: number     // 已选商品数量
  total_amount: string       // 商品总金额
  selected_amount: string    // 已选商品总金额
}

// 购物车完整信息
export interface CartInfo {
  items: CartItem[]
  stats: CartStats
} 
import { PaginationParams, PaginationResponse, StatsData } from './base'
import { User } from './user'
// 商品数据模型
export interface Product {
  id: number
  name: string
  description: string
  category: ProductCategory
  price: number
  stock: number
  sales: number
  images: string[]
  status: ProductStatus
  rating: number
  tags: string[]
  created_at: string
  updated_at: string
}

// Review 类型
export interface Review {
  id: number; // 评价的唯一标识符
  user: User; // 关联的用户信息
  product_id: number; // 关联的产品ID
  product: Product; // 关联的产品对象
  order_id: number; // 关联的订单ID
  rating: number; // 评分（1-5星）
  content: string; // 评价内容
  images: string[]; // 评价图片，URL列表
  created_at: string; // 创建时间，ISO 格式的日期字符串
  updated_at: string; // 更新时间，ISO 格式的日期字符串
}

// Response 类型
export interface ReviewResponse {
  reviews: Review[]; // 评论列表
  total: number; // 总评论数
}

// 商品列表查询参数
export interface ListProductParams extends PaginationParams {
  category_id?: number
  status?: ProductStatus
  keyword?: string
  min_price?: number
  max_price?: number
  sort_by?: string
  sort_order?: 'asc' | 'desc'
}

export interface ListReviewsParams extends PaginationParams { }

// 商品列表响应
export type ListProductResponse = PaginationResponse<Product>

// 创建商品请求参数
export interface CreateProductRequest {
  name: string
  description: string
  category: string
  price: number
  stock: number
  rating: number
  sales: number
  images: string[]
  tags?: string[]
  status: ProductStatus
}

// 更新商品请求参数
export interface UpdateProductRequest {
  name?: string
  description?: string
  category?: string
  price?: number
  rating: number
  stock?: number
  sales: number
  images?: string[]
  tags?: string[]
  status?: ProductStatus
}

// 商品统计信息
export interface ProductStats {
  total_count: number
  active_count: number
  inactive_count: number
  out_of_stock_count: number
  category_count: number
  total_sales: number
  total_amount: number
}

// 商品状态
export enum ProductStatus {
  Active = 'active',     // 上架
  Inactive = 'inactive', // 下架
  OutOfStock = 'out_of_stock' // 缺货
}

// 商品分类
// 商品分类枚举
export enum ProductCategory {
  Clothing = 'clothing',      // 服装
  Electronics = 'electronics', // 电子产品
  Food = 'food',              // 食品
  Books = 'books',            // 图书
  Beauty = 'beauty',          // 美妆
  Sports = 'sports',          // 运动
  Home = 'home',              // 家居
  Toys = 'toys'               // 玩具
}

// 分类显示名称映射
export const ProductCategoryNames: Record<ProductCategory, string> = {
  [ProductCategory.Clothing]: '服装',
  [ProductCategory.Electronics]: '电子产品',
  [ProductCategory.Food]: '食品',
  [ProductCategory.Books]: '图书',
  [ProductCategory.Beauty]: '美妆',
  [ProductCategory.Sports]: '运动',
  [ProductCategory.Home]: '家居',
  [ProductCategory.Toys]: '玩具'
}

// 分类列表查询参数
export interface ListCategoryParams extends PaginationParams {
  parent_id?: number
  keyword?: string
}

// 分类列表响应
export type ListCategoryResponse = PaginationResponse<ProductCategory>

// 创建分类请求
export interface CreateCategoryRequest {
  name: string
  description?: string
  parent_id?: number
  sort_order?: number
  image?: string
}

// 更新分类请求
export interface UpdateCategoryRequest {
  name?: string
  description?: string
  parent_id?: number
  sort_order?: number
  image?: string
}

// 更新商品状态请求
export interface UpdateProductStatusRequest {
  status: ProductStatus
} 
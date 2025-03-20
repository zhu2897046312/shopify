import type { User } from './user'
import type { Product } from './product'

export enum OrderStatus {
  Pending = 'pending',       // 待付款
  Paid = 'paid',            // 已付款
  Shipped = 'shipped',      // 已发货
  Delivered = 'delivered',  // 已送达
  Completed = 'completed',  // 已完成
  Cancelled = 'cancelled',  // 已取消
  Refunding = 'refunding',  // 退款中
  Refunded = 'refunded'     // 已退款
}

export enum PaymentStatus {
  Paid = 'paid',
  Unpaid = 'unpaid'
}

export enum PaymentMethod {
  Alipay = 'alipay',    // 支付宝
  WechatPay = 'wechat', // 微信支付
  CreditCard = 'credit_card' // 信用卡
}

// 地址信息
export interface Address {
  id: number
  name: string
  phone: string
  province: string
  city: string
  district: string
  street: string
  post_code: string
  is_default: boolean
}

// 物流信息
export interface Logistics {
  id: number
  order_id: number
  tracking_number: string
  carrier: string
  status: string
  traces: LogisticsTrace[]
  created_at: string
  updated_at: string
}

// 物流跟踪
export interface LogisticsTrace {
  id: number
  logistics_id: number
  content: string
  location: string
  created_at: string
}

// 订单商品
// 订单商品
export interface OrderItem {
  id: number
  order_id: number
  product_id: number
  product: Product
  quantity: number
  price: string         // 改为 string 类型
  created_at: string    // 新增
  updated_at: string    // 新增
}

// 订单信息
export interface Order {
  id: number
  user_id: number
  order_number: string
  status: OrderStatus
  total_amount: string  // 改为 string 类型，因为后端返回的是字符串
  address_id: number    // 新增
  address: Address
  payment_method: PaymentMethod
  payment_status: PaymentStatus
  payment_time: string | null  // 新增
  order_items: OrderItem[]     // 改名从 items 到 order_items
  logistics?: Logistics
  created_at: string
  updated_at: string
}

// 订单统计
export interface OrderStats {
  total_count: number
  total_amount: number
  today_count: number
  today_amount: number
  pending_count: number
  shipped_count: number
  completed_count: number
  cancelled_count: number
}

// 订单列表查询参数
export interface ListOrderParams {
  page: number
  page_size: number
  status?: OrderStatus
  payment_status?: PaymentStatus
  start_time?: string
  end_time?: string
  keyword?: string
}

// 订单列表响应
export interface ListOrderResponse {
  orders: Order[]
  total: number
  page: number
  page_size: number
}

// 更新订单状态请求
export interface UpdateOrderStatusRequest {
  status: OrderStatus
}

// 更新物流信息请求
export interface UpdateLogisticsRequest {
  tracking_number: string
  carrier: string
}

// 添加物流跟踪请求
export interface AddLogisticsTraceRequest {
  content: string
  location: string
} 
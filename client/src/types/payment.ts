import { Order } from './order'
import { PaginationParams, PaginationResponse } from './base'

// 支付方式枚举
export enum PaymentMethod {
  Wechat = 'wechat',  // 微信支付
  Alipay = 'alipay'   // 支付宝支付
}

// 支付状态枚举
export enum PaymentStatus {
  Pending = 'pending',    // 待支付
  Paid = 'paid',         // 已支付
  Failed = 'failed',     // 支付失败
  Refunded = 'refunded'  // 已退款
}

// 支付记录数据模型
export interface Payment {
  id: number
  order_id: number
  order: Order
  payment_method: PaymentMethod
  amount: string
  trade_no: string
  status: PaymentStatus
  pay_time: string | null
  created_at: string
  updated_at: string
}

// 支付回调记录
export interface PaymentCallback {
  id: number
  payment_id: number
  payment: Payment
  trade_no: string
  status: PaymentStatus
  raw_data: string
  created_at: string
  updated_at: string
}

// 创建支付请求参数
export interface CreatePaymentRequest {
  order_id: number
  method: PaymentMethod
}

// 创建支付响应
export interface CreatePaymentResponse {
  payment_id: number
  pay_url: string  // 支付链接或二维码链接
}

// 支付回调请求参数
export interface PaymentCallbackRequest {
  [key: string]: string  // 不同支付方式的回调参数不同
}

// 支付列表查询参数
export interface ListPaymentParams extends PaginationParams {
  order_id?: number
  method?: PaymentMethod
  status?: PaymentStatus
  start_time?: string
  end_time?: string
}

// 支付列表响应
export type ListPaymentResponse = PaginationResponse<Payment>

// 支付统计信息
export interface PaymentStats {
  total_amount: string           // 总支付金额
  today_amount: string          // 今日支付金额
  total_count: number           // 总支付笔数
  today_count: number          // 今日支付笔数
  method_stats: {              // 支付方式统计
    [key in PaymentMethod]: {
      amount: string
      count: number
    }
  }
  status_stats: {              // 支付状态统计
    [key in PaymentStatus]: {
      amount: string
      count: number
    }
  }
  hourly_stats: {             // 每小时支付统计
    hour: number
    amount: string
    count: number
  }[]
}

// 退款请求参数
export interface RefundRequest {
  payment_id: number
  amount: string
  reason: string
}

// 退款响应
export interface RefundResponse {
  refund_id: string
  status: PaymentStatus
  refund_time: string
} 
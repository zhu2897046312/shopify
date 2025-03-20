import request from '@/utils/request'
import type {
  Order,
  OrderStats,
  OrderStatus,
  ListOrderParams,
  ListOrderResponse,
  UpdateOrderStatusRequest,
  UpdateLogisticsRequest,
  AddLogisticsTraceRequest,
  Logistics
} from '@/types/order'

// 获取订单列表
export function getOrders(params: ListOrderParams): Promise<ListOrderResponse> {
  return request.get('/admin/orders', { params })
}

// 获取订单详情
export function getOrder(id: number): Promise<Order> {
  return request.get(`/admin/orders/${id}`)
}

// 更新订单状态
export function updateOrderStatus(id: number, data: UpdateOrderStatusRequest): Promise<Order> {
  return request.put(`/admin/orders/${id}/status`, data)
}

// 更新物流信息
export function updateLogistics(id: number, data: UpdateLogisticsRequest): Promise<Logistics> {
  return request.post(`/admin/orders/${id}/logistics`, data)
}

// 添加物流跟踪
export function addLogisticsTrace(id: number, data: AddLogisticsTraceRequest): Promise<Logistics> {
  return request.post(`/admin/orders/${id}/logistics/trace`, data)
}

// 获取订单统计信息
export function getOrderStats(): Promise<OrderStats> {
  return request.get('/admin/orders/stats')
}

// 导出订单数据
export function exportOrders(params: ListOrderParams) {
  return request.get('/admin/orders/export', {
    params,
    responseType: 'blob'
  })
}

// 批量更新订单状态
export function batchUpdateOrderStatus(ids: number[], status: OrderStatus): Promise<void> {
  return request.put('/admin/orders/batch/status', {
    ids,
    status
  })
}

// 批量删除订单
export function batchDeleteOrders(ids: number[]): Promise<void> {
  return request.delete('/admin/orders/batch', {
    data: { ids }
  })
}
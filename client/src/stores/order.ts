import { defineStore } from 'pinia'
import { ref } from 'vue'
import { storeToRefs } from 'pinia'
import type {
  Order,
  OrderStats,
  ListOrderParams,
  UpdateLogisticsRequest,
  AddLogisticsTraceRequest
} from '@/types/order'
import { OrderStatus } from '@/types/order'
import * as orderApi from '@/api/order'
import { ElMessage } from 'element-plus'

export const useOrderStore = defineStore('admin-order', () => {
  // 状态
  const orders = ref<Order[]>([])
  const total = ref(0)
  const loading = ref(false)
  const currentOrder = ref<Order | null>(null)
  const stats = ref<OrderStats>({
    total_count: 0,
    total_amount: 0,
    today_count: 0,
    today_amount: 0,
    pending_count: 0,
    shipped_count: 0,
    completed_count: 0,
    cancelled_count: 0
  })

  // 获取订单列表
  const getOrders = async (params: ListOrderParams) => {
    loading.value = true
    try {
      const res = await orderApi.getOrders(params)
      orders.value = res.orders
      total.value = res.total
      console.log(res.orders)
      return res
    } finally {
      loading.value = false
    }
  }

  // 获取订单详情
  const getOrder = async (id: number) => {
    loading.value = true
    try {
      const order = await orderApi.getOrder(id)
      currentOrder.value = order
      return order
    } finally {
      loading.value = false
    }
  }

  // 更新订单状态
  const updateOrderStatus = async (id: number, status: OrderStatus) => {
    loading.value = true
    try {
      await orderApi.updateOrderStatus(id, { status })
      ElMessage.success('状态更新成功')
    } finally {
      loading.value = false
    }
  }

  // 批量更新订单状态
  const batchUpdateOrderStatus = async (ids: number[], status: OrderStatus) => {
    loading.value = true
    try {
      await orderApi.batchUpdateOrderStatus(ids, status)
      ElMessage.success('批量更新成功')
    } finally {
      loading.value = false
    }
  }

  // 批量删除订单
  const batchDeleteOrders = async (ids: number[]) => {
    loading.value = true
    try {
      await orderApi.batchDeleteOrders(ids)
    } finally {
      loading.value = false
    }
  }

  // 更新物流信息
  const updateLogistics = async (id: number, data: UpdateLogisticsRequest) => {
    loading.value = true
    try {
      await orderApi.updateLogistics(id, data)
      ElMessage.success('物流信息更新成功')
    } finally {
      loading.value = false
    }
  }

  // 添加物流跟踪
  const addLogisticsTrace = async (id: number, data: AddLogisticsTraceRequest) => {
    loading.value = true
    try {
      await orderApi.addLogisticsTrace(id, data)
    } finally {
      loading.value = false
    }
  }

  // 获取订单统计
  const getOrderStats = async (): Promise<OrderStats> => {
    try {
      // 获取所有订单数据来计算统计信息
      const { orders: allOrders } = await getOrders({ page: 1, page_size: 1000 })

      const now = new Date()
      const today = new Date(now.getFullYear(), now.getMonth(), now.getDate()).getTime()

      // 计算今日订单
      const todayOrders = allOrders.filter(order => 
        new Date(order.created_at).getTime() >= today
      )

      // 计算各种状态的订单数量
      const order_stats: OrderStats = {
        total_count: allOrders.length,
        total_amount: calculateTotalAmount(allOrders),
        today_count: todayOrders.length,
        today_amount: calculateTotalAmount(todayOrders),
        pending_count: allOrders.filter(order => order.status === OrderStatus.Pending).length,
        shipped_count: allOrders.filter(order => order.status === OrderStatus.Shipped).length,
        completed_count: allOrders.filter(order => order.status === OrderStatus.Completed).length,
        cancelled_count: allOrders.filter(order => order.status === OrderStatus.Cancelled).length
      }

      // 更新store中的统计数据
      stats.value = order_stats
      return order_stats
    } catch (error) {
      console.error('Failed to calculate order stats:', error)
      throw error
    }
  }

  // 计算订单总金额
  const calculateTotalAmount = (orders: Order[]): number => {
    return Number(orders.reduce((sum, order) => 
      sum + Number(order.total_amount), 0
    ).toFixed(2))
  }

  // 获取订单趋势数据
  const getOrderTrend = async (timeRange: string) => {
    try {
      // 获取所有订单
      const { orders: allOrders } = await getOrders({ page: 1, page_size: 1000 })

      // 根据时间范围计算日期范围
      const now = new Date()
      let startDate = new Date()
      let dateFormat: string
      let step: number

      switch (timeRange) {
        case 'week':
          startDate.setDate(now.getDate() - 7)
          dateFormat = 'MM-dd'
          step = 24 * 60 * 60 * 1000 // 1 day
          break
        case 'month':
          startDate.setMonth(now.getMonth() - 1)
          dateFormat = 'MM-dd'
          step = 24 * 60 * 60 * 1000 // 1 day
          break
        case 'year':
          startDate.setFullYear(now.getFullYear() - 1)
          dateFormat = 'yyyy-MM'
          step = 30 * 24 * 60 * 60 * 1000 // 30 days
          break
        default:
          throw new Error('Invalid time range')
      }

      // 生成日期数组
      const dates: string[] = []
      const data: { date: string; count: number; amount: number }[] = []
      
      for (let date = startDate; date <= now; date = new Date(date.getTime() + step)) {
        const dateStr = formatDate(date, dateFormat)
        dates.push(dateStr)
        
        // 统计当前日期的订单
        const periodStart = date.getTime()
        const periodEnd = new Date(date.getTime() + step).getTime()
        
        const periodOrders = allOrders.filter(order => {
          const orderTime = new Date(order.created_at).getTime()
          return orderTime >= periodStart && orderTime < periodEnd
        })

        data.push({
          date: dateStr,
          count: periodOrders.length,
          amount: calculateTotalAmount(periodOrders)
        })
      }

      return {
        data
      }
    } catch (error) {
      console.error('Failed to get order trend:', error)
      throw error
    }
  }

  // 格式化日期
  const formatDate = (date: Date, format: string): string => {
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')

    return format
      .replace('yyyy', String(year))
      .replace('MM', month)
      .replace('dd', day)
  }

  return {
    // 状态
    orders,
    total,
    loading,
    currentOrder,
    stats,

    // 方法
    getOrders,
    getOrder,
    updateOrderStatus,
    batchUpdateOrderStatus,
    batchDeleteOrders,
    updateLogistics,
    addLogisticsTrace,
    getOrderStats,
    getOrderTrend
  }
})

// 导出组合式函数，用于在组件中获取响应式状态
export function useOrderStoreRefs() {
  const store = useOrderStore()
  return {
    ...storeToRefs(store),
    // 非响应式方法直接返回
    getOrders: store.getOrders,
    getOrder: store.getOrder,
    updateOrderStatus: store.updateOrderStatus,
    batchUpdateOrderStatus: store.batchUpdateOrderStatus,
    batchDeleteOrders: store.batchDeleteOrders,
    updateLogistics: store.updateLogistics,
    addLogisticsTrace: store.addLogisticsTrace,
    getOrderStats: store.getOrderStats,
    getOrderTrend: store.getOrderTrend
  }
}
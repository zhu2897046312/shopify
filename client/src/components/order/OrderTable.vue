<template>
  <div class="order-table">
    <el-table
      v-loading="loading"
      :data="orders"
      border
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column prop="order_number" label="订单号" width="180" />
      <el-table-column label="用户信息" min-width="200">
        <template #default="{ row }">
          <div class="user-info">
            <el-avatar :size="40">
              {{ row.user?.nickname?.charAt(0)?.toUpperCase() || 'U' }}
            </el-avatar>
            <div class="user-detail">
              <div class="nickname">{{ row.user?.nickname || '未知用户' }}</div>
              <div class="email">{{ row.user?.email || '-' }}</div>
            </div>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="total_amount" label="订单金额" width="120">
        <template #default="{ row }">
          ¥{{ row.total_amount }}
        </template>
      </el-table-column>
      <el-table-column label="订单状态" width="120">
        <template #default="{ row }">
          <el-tag :type="getOrderStatusType(row.status)">
            {{ getOrderStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="支付状态" width="120">
        <template #default="{ row }">
          <el-tag :type="row.payment_status === 'paid' ? 'success' : 'warning'">
            {{ row.payment_status === 'paid' ? '已支付' : '未支付' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="下单时间" width="180" />
      <el-table-column label="商品信息" min-width="200">
        <template #default="{ row }">
          <div v-for="item in row.order_items" :key="item.id">
            {{ item.product.name }} x {{ item.quantity }}
          </div>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="300" fixed="right">
        <template #default="{ row }">
          <el-button-group>
            <el-button type="primary" link @click="handleDetail(row)">
              详情
            </el-button>
            <el-button 
              type="primary" 
              link 
              @click="handleStatusChange(row)"
              :disabled="!canUpdateStatus(row.status)"
            >
              {{ getNextStatusText(row.status) }}
            </el-button>
            <el-button 
              type="primary" 
              link 
              @click="handleViewLogistics(row)"
              v-if="row.logistics?.tracking_number"
            >
              查看物流
            </el-button>
            <el-button 
              type="primary" 
              link 
              @click="handleLogistics(row)"
              v-if="canUpdateLogistics(row.status)"
            >
              更新物流
            </el-button>
            
          </el-button-group>
        </template>
      </el-table-column>
    </el-table>

    <div class="pagination">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    <!-- 物流信息对话框 -->
    <el-dialog
      v-model="logisticsDialogVisible"
      title="物流信息"
      width="600px"
    >
      <logistics-form
        v-if="currentOrder"
        :order="currentOrder"
        @success="handleLogisticsSuccess"
        @cancel="logisticsDialogVisible = false"
      />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useOrderStore } from '@/stores/order'
import { useRouter } from 'vue-router'
import type { Order, LogisticsTrace } from '@/types/order'
import { OrderStatus } from '@/types/order'
import LogisticsForm from './LogisticsForm.vue'

const emit = defineEmits(['refresh'])
const orderStore = useOrderStore()
const router = useRouter()

const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const currentOrder = ref<Order | null>(null)
const logisticsDialogVisible = ref(false)
const orders = ref<Order[]>([])

// 加载订单列表
const loadOrders = async () => {
  loading.value = true
  try {
    const res = await orderStore.getOrders({
      page: currentPage.value,
      page_size: pageSize.value
    })
    orders.value = res.orders
    total.value = res.total
    emit('refresh')
  } catch (error: any) {
    ElMessage.error(error.message || '加载订单失败')
  } finally {
    loading.value = false
  }
}

// 选择变化
const handleSelectionChange = (selection: Order[]) => {
  // 处理表格选择变化
}

// 查看详情
const handleDetail = (row: Order) => {
  router.push(`/orders/${row.id}`)
}

// 更改状态
const handleStatusChange = async (row: Order) => {
  const nextStatus = getNextStatus(row.status)
  if (!nextStatus) {
    ElMessage.warning('当前状态无法更改')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确定要将订单状态更改为"${getOrderStatusText(nextStatus)}"吗？`,
      '提示',
      { 
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning' 
      }
    )
    
    loading.value = true
    await orderStore.updateOrderStatus(row.id, nextStatus)
    ElMessage.success(`订单状态已更新为"${getOrderStatusText(nextStatus)}"`)
    await loadOrders()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '更新状态失败')
    }
  } finally {
    loading.value = false
  }
}

// 更新物流
const handleLogistics = (row: Order) => {
  currentOrder.value = row
  logisticsDialogVisible.value = true
}

// 物流更新成功
const handleLogisticsSuccess = () => {
  logisticsDialogVisible.value = false
  loadOrders()
}

// 查看物流信息
const handleViewLogistics = (row: Order) => {
  if (!row.logistics) return
  router.push(`/orders/${row.id}/logistics`)
}

const formatDate = (date: Date) => {
  return new Intl.DateTimeFormat('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false
  }).format(date)
}

// 获取物流轨迹类型
const getTraceType = (trace: LogisticsTrace): '' | 'primary' | 'success' => {
  const content = trace.content.toLowerCase()
  if (content.includes('签收') || content.includes('已送达')) {
    return 'success'
  }
  if (content.includes('派送') || content.includes('运输')) {
    return 'primary'
  }
  return ''
}

// 订单状态相关方法
const getOrderStatusType = (status: OrderStatus): string => {
  const map: Record<OrderStatus, string> = {
    [OrderStatus.Pending]: 'warning',
    [OrderStatus.Paid]: 'primary',
    [OrderStatus.Shipped]: 'info',
    [OrderStatus.Delivered]: 'success',
    [OrderStatus.Completed]: 'success',
    [OrderStatus.Cancelled]: 'danger',
    [OrderStatus.Refunding]: 'warning',
    [OrderStatus.Refunded]: 'info'
  }
  return map[status]
}

const getOrderStatusText = (status: OrderStatus): string => {
  const map: Record<OrderStatus, string> = {
    [OrderStatus.Pending]: '待付款',
    [OrderStatus.Paid]: '待发货',
    [OrderStatus.Shipped]: '已发货',
    [OrderStatus.Delivered]: '已送达',
    [OrderStatus.Completed]: '已完成',
    [OrderStatus.Cancelled]: '已取消',
    [OrderStatus.Refunding]: '退款中',
    [OrderStatus.Refunded]: '已退款'
  }
  return map[status]
}

// 分页相关
const handleSizeChange = (size: number) => {
  pageSize.value = size
  currentPage.value = 1
  loadOrders()
}

const handleCurrentChange = (page: number) => {
  currentPage.value = page
  loadOrders()
}

// 订单状态相关方法
const getNextStatus = (status: OrderStatus): OrderStatus | null => {
  const map: Partial<Record<OrderStatus, OrderStatus>> = {
    [OrderStatus.Pending]: OrderStatus.Paid,
    [OrderStatus.Paid]: OrderStatus.Shipped,
    [OrderStatus.Shipped]: OrderStatus.Delivered,
    [OrderStatus.Delivered]: OrderStatus.Completed
  }
  return map[status] || null
}

const getNextStatusText = (status: OrderStatus): string => {
  const nextStatus = getNextStatus(status)
  if (!nextStatus) return '无法更改'
  
  const actionMap: Partial<Record<OrderStatus, string>> = {
    [OrderStatus.Paid]: '发货',
    [OrderStatus.Shipped]: '送达',
    [OrderStatus.Delivered]: '完成',
  }
  return actionMap[status] || '更新状态'
}

const canUpdateStatus = (status: OrderStatus): boolean => {
  return getNextStatus(status) !== null && 
    ![OrderStatus.Completed, OrderStatus.Cancelled, OrderStatus.Refunded].includes(status)
}

const canUpdateLogistics = (status: OrderStatus): boolean => {
  return [OrderStatus.Paid, OrderStatus.Shipped, OrderStatus.Delivered, OrderStatus.Completed].includes(status)
}

// 初始加载
loadOrders()
</script>

<style scoped lang="scss">
.order-table {
  .user-info {
    display: flex;
    align-items: center;
    gap: 12px;

    .user-detail {
      .nickname {
        font-weight: bold;
      }
      .email {
        font-size: 12px;
        color: #666;
      }
    }
  }

  .pagination {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}

:deep(.logistics-dialog) {
  .el-message-box__content {
    max-height: 70vh;
    overflow-y: auto;
  }
  
  .logistics-info {
    text-align: left;
    
    .logistics-header {
      margin-bottom: 20px;
      padding-bottom: 15px;
      border-bottom: 1px solid #eee;
      
      p {
        margin: 8px 0;
      }
    }
    
    .logistics-traces {
      .el-timeline-item {
        padding-bottom: 20px;
        
        .el-timeline-item__timestamp {
          font-size: 13px;
        }
        
        h4 {
          margin: 0;
          font-size: 14px;
          color: #303133;
        }
        
        .location {
          margin: 5px 0 0;
          color: #666;
          font-size: 13px;
        }
      }
    }
  }
}
</style>
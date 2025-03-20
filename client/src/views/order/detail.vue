<template>
  <div class="order-detail-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>订单详情</span>
          <div class="actions">
            <el-button 
              type="primary"
              :disabled="!canUpdateStatus(order?.status)"
              @click="handleStatusChange"
            >
              {{ getNextStatusText(order?.status) }}
            </el-button>
            <el-button
              v-if="canUpdateLogistics(order?.status)"
              @click="handleLogistics"
            >
              更新物流
            </el-button>
            <el-button
              v-if="order?.logistics"
              @click="handleViewLogistics"
            >
              查看物流
            </el-button>
          </div>
        </div>
      </template>

      <order-detail v-if="order" :order="order" />
    </el-card>

    <!-- 物流信息对话框 -->
    <el-dialog
      v-model="logisticsDialogVisible"
      title="物流信息"
      width="600px"
    >
      <logistics-form
        v-if="!viewMode"
        :order="order!"
        @success="handleLogisticsSuccess"
        @cancel="logisticsDialogVisible = false"
      />
      <div v-else class="logistics-details">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="物流单号">{{ order?.logistics?.tracking_number }}</el-descriptions-item>
          <el-descriptions-item label="承运商">{{ order?.logistics?.carrier }}</el-descriptions-item>
        </el-descriptions>
        
        <!-- 物流轨迹 -->
        <div class="logistics-traces">
          <div class="title">物流轨迹</div>
          <el-timeline>
            <el-timeline-item
              v-for="trace in order?.logistics?.traces"
              :key="trace.id"
              :type="getTraceType(trace)"
              :timestamp="formatDate(new Date(trace.created_at))"
            >
              <div class="trace-content">
                <div>{{ trace.content }}</div>
                <div class="location">{{ trace.location }}</div>
              </div>
            </el-timeline-item>
          </el-timeline>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import { useOrderStore } from '@/stores/order'
import { OrderStatus } from '@/types/order'
import type { Order, LogisticsTrace } from '@/types/order'
import OrderDetail from '@/components/order/OrderDetail.vue'
import LogisticsForm from '@/components/order/LogisticsForm.vue'

const route = useRoute()
const router = useRouter()
const orderStore = useOrderStore()

const order = ref<Order>()
const logisticsDialogVisible = ref(false)
const viewMode = ref(false)
const loading = ref(false)

// 加载订单详情
const loadOrder = async () => {
  loading.value = true
  try {
    const id = Number(route.params.id)
    if (!id) {
      router.push('/orders')
      return
    }

    order.value = await orderStore.getOrder(id)
    if (!order.value) {
      ElMessage.error('订单不存在')
      router.push('/orders')
    }
  } catch (error: any) {
    ElMessage.error(error.message || '加载订单详情失败')
    router.push('/orders')
  } finally {
    loading.value = false
  }
}

// 更改状态
const handleStatusChange = async () => {
  if (!order.value) return
  
  const nextStatus = getNextStatus(order.value.status)
  if (!nextStatus) return

  try {
    await ElMessageBox.confirm(
      `确定要将订单状态更改为"${getNextStatusText(nextStatus)}"吗？`,
      '提示',
      { type: 'warning' }
    )
    await orderStore.updateOrderStatus(order.value.id, nextStatus)
    loadOrder()
  } catch {
    // 用户取消操作
  }
}

// 查看物流信息
const handleViewLogistics = () => {
  if (!order.value?.logistics) {
    ElMessage.warning('暂无物流信息')
    return
  }
  viewMode.value = true
  logisticsDialogVisible.value = true
}

// 更新物流信息
const handleLogistics = () => {
  viewMode.value = false
  logisticsDialogVisible.value = true
}

// 物流信息更新成功
const handleLogisticsSuccess = () => {
  logisticsDialogVisible.value = false
  loadOrder()
}

// 订单状态相关方法
const getNextStatus = (status?: OrderStatus): OrderStatus | null => {
  if (!status) return null
  
  const map: Partial<Record<OrderStatus, OrderStatus>> = {
    [OrderStatus.Pending]: OrderStatus.Paid,
    [OrderStatus.Paid]: OrderStatus.Shipped,
    [OrderStatus.Shipped]: OrderStatus.Completed
  }
  return map[status] || null
}

const getNextStatusText = (status?: OrderStatus): string => {
  if (!status) return ''
  const nextStatus = getNextStatus(status)
  if (!nextStatus) return ''

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
  return map[nextStatus]
}

const canUpdateStatus = (status?: OrderStatus): boolean => {
  if (!status) return false
  return status !== OrderStatus.Completed && status !== OrderStatus.Cancelled
}

const canUpdateLogistics = (status?: OrderStatus): boolean => {
  if (!status) return false
  return status === OrderStatus.Shipped || status === OrderStatus.Completed
}

// 格式化日期
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

// 获取轨迹类型
const getTraceType = (trace: LogisticsTrace): '' | 'primary' | 'success' => {
  if (trace === order.value?.logistics?.traces[0]) return 'success'
  if (trace === order.value?.logistics?.traces[order.value?.logistics?.traces.length - 1]) return 'primary'
  return ''
}

// 初始加载
onMounted(() => {
  loadOrder()
})
</script>

<style scoped lang="scss">
.order-detail-page {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .logistics-details {
    .title {
      font-size: 16px;
      font-weight: 500;
      margin: 20px 0;
    }

    .logistics-traces {
      margin-top: 20px;
    }

    .trace-content {
      .location {
        font-size: 12px;
        color: #999;
        margin-top: 4px;
      }
    }
  }
}
</style>
<template>
  <div class="order-detail">
    <!-- 基本信息 -->
    <el-descriptions title="订单信息" :column="2" border>
      <el-descriptions-item label="订单号">{{ order.order_number }}</el-descriptions-item>
      <el-descriptions-item label="下单时间">{{ order.created_at }}</el-descriptions-item>
      <el-descriptions-item label="订单状态">
        <el-tag :type="getOrderStatusType(order.status)">
          {{ getOrderStatusText(order.status) }}
        </el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="支付状态">
        <el-tag :type="order.payment_status === PaymentStatus.Paid ? 'success' : 'warning'">
          {{ order.payment_status === PaymentStatus.Paid ? '已支付' : '未支付' }}
        </el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="订单金额">¥{{ order.total_amount }}</el-descriptions-item>
      <el-descriptions-item label="支付方式">{{ getPaymentMethodText(order.payment_method) }}</el-descriptions-item>
      <el-descriptions-item label="支付时间" v-if="order.payment_time">
        {{ order.payment_time }}
      </el-descriptions-item>
    </el-descriptions>

    <!-- 收货信息 -->
    <el-descriptions title="收货信息" :column="1" border class="mt-4">
      <el-descriptions-item label="收货人">{{ order.address.name }}</el-descriptions-item>
      <el-descriptions-item label="联系电话">{{ order.address.phone }}</el-descriptions-item>
      <el-descriptions-item label="收货地址">
        {{ formatAddress(order.address) }}
      </el-descriptions-item>
      <el-descriptions-item label="邮政编码" v-if="order.address.post_code">
        {{ order.address.post_code }}
      </el-descriptions-item>
    </el-descriptions>

    <!-- 商品信息 -->
    <div class="order-items mt-4">
      <div class="title">商品信息</div>
      <el-table :data="order.order_items" border>
        <el-table-column label="商品信息">
          <template #default="{ row }">
            <div class="product-info">
              <el-image
                v-if="row.product.image"
                :src="row.product.image"
                :preview-src-list="[row.product.image]"
                fit="cover"
                class="product-image"
              />
              <div class="product-detail">
                <div class="name">{{ row.product.name }}</div>
                <div class="sku" v-if="row.product.sku">SKU: {{ row.product.sku }}</div>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="price" label="单价" width="120">
          <template #default="{ row }">¥{{ row.price }}</template>
        </el-table-column>
        <el-table-column prop="quantity" label="数量" width="120" />
        <el-table-column label="小计" width="120">
          <template #default="{ row }">
            ¥{{ calculateSubtotal(row) }}
          </template>
        </el-table-column>
      </el-table>
      <div class="order-total">
        总计：<span class="amount">¥{{ order.total_amount }}</span>
      </div>
    </div>

    <!-- 物流信息 -->
    <div v-if="showLogistics" class="logistics-info mt-4">
      <div class="title">物流信息</div>
      <!-- 添加物流信息按钮 -->
      <div v-if="!order.logistics" class="add-logistics">
        <el-empty description="暂无物流信息">
          <template #extra>
            <el-button type="primary" @click="$emit('update-logistics', order)">
              添加物流信息
            </el-button>
          </template>
        </el-empty>
      </div>
      
      <!-- 已有物流信息 -->
      <template v-else>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="物流公司">{{ order.logistics.carrier }}</el-descriptions-item>
          <el-descriptions-item label="物流单号">
            <el-link type="primary" @click="handleTrackOrder(order.logistics.tracking_number)">
              {{ order.logistics.tracking_number }}
            </el-link>
          </el-descriptions-item>
          <el-descriptions-item label="物流状态">
            <el-tag type="info">{{ order.logistics.status }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="更新时间">{{ order.logistics.updated_at }}</el-descriptions-item>
        </el-descriptions>

        <!-- 物流跟踪记录 -->
        <div v-if="order.logistics.traces?.length" class="logistics-traces mt-4">
          <div class="subtitle">物流跟踪</div>
          <el-timeline>
            <el-timeline-item
              v-for="trace in order.logistics.traces"
              :key="trace.id"
              :timestamp="trace.created_at"
              :type="getTraceType(trace)"
            >
              <div class="trace-content">
                <div class="location" v-if="trace.location">{{ trace.location }}</div>
                <div class="content">{{ trace.content }}</div>
              </div>
            </el-timeline-item>
          </el-timeline>
        </div>

        <!-- 更新物流按钮 -->
        <div class="actions mt-4">
          <el-button-group>
            <el-button type="primary" @click="$emit('update-logistics', order)">
              更新物流信息
            </el-button>
            <el-button type="primary" @click="$emit('add-trace', order)">
              添加物流跟踪
            </el-button>
          </el-button-group>
        </div>
      </template>
    </div>

    <!-- 操作按钮 -->
    <div class="order-actions mt-4">
      <el-button-group>
        <el-button
          v-if="canUpdateStatus"
          type="primary"
          @click="$emit('update-status', order)"
        >
          {{ getNextStatusText(order.status) }}
        </el-button>
      </el-button-group>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Order, OrderStatus, PaymentMethod } from '@/types/order'
import { OrderStatus as OrderStatusEnum, PaymentStatus, PaymentMethod as PaymentMethodEnum } from '@/types/order'
import type { Address, OrderItem, LogisticsTrace } from '@/types/order'

const props = defineProps<{
  order: Order
}>()

defineEmits<{
  (e: 'update-status', order: Order): void
  (e: 'update-logistics', order: Order): void
  (e: 'add-trace', order: Order): void
}>()

// 订单状态相关方法
const getOrderStatusType = (status: OrderStatus): string => {
  const map: Record<OrderStatus, string> = {
    [OrderStatusEnum.Pending]: 'warning',
    [OrderStatusEnum.Paid]: 'primary',
    [OrderStatusEnum.Shipped]: 'info',
    [OrderStatusEnum.Delivered]: 'success',
    [OrderStatusEnum.Completed]: 'success',
    [OrderStatusEnum.Cancelled]: 'danger',
    [OrderStatusEnum.Refunding]: 'warning',
    [OrderStatusEnum.Refunded]: 'info'
  }
  return map[status]
}

const getOrderStatusText = (status: OrderStatus): string => {
  const map: Record<OrderStatus, string> = {
    [OrderStatusEnum.Pending]: '待付款',
    [OrderStatusEnum.Paid]: '待发货',
    [OrderStatusEnum.Shipped]: '已发货',
    [OrderStatusEnum.Delivered]: '已送达',
    [OrderStatusEnum.Completed]: '已完成',
    [OrderStatusEnum.Cancelled]: '已取消',
    [OrderStatusEnum.Refunding]: '退款中',
    [OrderStatusEnum.Refunded]: '已退款'
  }
  return map[status]
}

const getPaymentMethodText = (method: PaymentMethod): string => {
  const map: Record<PaymentMethod, string> = {
    [PaymentMethodEnum.Alipay]: '支付宝',
    [PaymentMethodEnum.WechatPay]: '微信支付',
    [PaymentMethodEnum.CreditCard]: '信用卡'
  }
  return map[method] || '未知'
}

// 格式化地址
const formatAddress = (address: Address): string => {
  return `${address.province}${address.city}${address.district}${address.street}`
}

// 计算小计金额
const calculateSubtotal = (item: OrderItem): string => {
  const price = parseFloat(item.price)
  const quantity = item.quantity
  return (price * quantity).toFixed(2)
}

// 物流相关
const showLogistics = computed(() => {
  const status = props.order.status
  return [
    OrderStatusEnum.Shipped,
    OrderStatusEnum.Delivered,
    OrderStatusEnum.Completed
  ].includes(status)
})

const getTraceType = (trace: LogisticsTrace): '' | 'primary' | 'success' => {
  if (trace.content.includes('已签收')) return 'success'
  if (trace.content.includes('派送中')) return 'primary'
  return ''
}

// 状态更新相关
const canUpdateStatus = computed(() => {
  const status = props.order.status
  return ![
    OrderStatusEnum.Completed,
    OrderStatusEnum.Cancelled,
    OrderStatusEnum.Refunded
  ].includes(status)
})

const getNextStatusText = (status: OrderStatus): string => {
  const map: Partial<Record<OrderStatus, string>> = {
    [OrderStatusEnum.Pending]: '确认支付',
    [OrderStatusEnum.Paid]: '发货',
    [OrderStatusEnum.Shipped]: '确认送达',
    [OrderStatusEnum.Delivered]: '完成订单'
  }
  return map[status] || '更新状态'
}

// 跟踪订单
const handleTrackOrder = (trackingNumber: string) => {
  // 这里可以实现跳转到物流跟踪页面或打开物流跟踪弹窗
  console.log('跟踪订单:', trackingNumber)
}
</script>

<style scoped lang="scss">
.order-detail {
  .mt-4 {
    margin-top: 16px;
  }

  .title {
    font-size: 16px;
    font-weight: bold;
    margin-bottom: 16px;
  }

  .subtitle {
    font-size: 14px;
    font-weight: bold;
    margin: 16px 0;
    color: #666;
  }

  .product-info {
    display: flex;
    align-items: center;
    gap: 12px;

    .product-image {
      width: 60px;
      height: 60px;
      border-radius: 4px;
    }

    .product-detail {
      .name {
        font-weight: bold;
      }
      .sku {
        font-size: 12px;
        color: #666;
        margin-top: 4px;
      }
    }
  }

  .order-total {
    margin-top: 16px;
    text-align: right;
    font-size: 14px;

    .amount {
      font-size: 18px;
      font-weight: bold;
      color: #f56c6c;
      margin-left: 8px;
    }
  }

  .add-logistics {
    padding: 24px;
    background-color: #f8f9fa;
    border-radius: 4px;
  }

  .logistics-traces {
    .trace-content {
      .location {
        font-size: 12px;
        color: #666;
        margin-bottom: 4px;
      }
      .content {
        color: #333;
      }
    }
  }

  .actions {
    display: flex;
    justify-content: flex-end;
  }

  .order-actions {
    display: flex;
    justify-content: flex-end;
    padding-top: 16px;
    border-top: 1px solid #ebeef5;
  }
}
</style>
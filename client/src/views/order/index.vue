<template>
  <div class="order-page">
    <!-- 搜索栏 -->
    <el-card class="search-card">
      <el-form :model="searchForm" inline>
        <el-form-item label="订单号">
          <el-input
            v-model="searchForm.order_id"
            placeholder="请输入订单号"
            clearable
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item label="订单状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option
              v-for="status in orderStatusOptions"
              :key="status.value"
              :label="status.label"
              :value="status.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="支付状态">
          <el-select v-model="searchForm.payment_status" placeholder="全部" clearable>
            <el-option label="已支付" :value="PaymentStatus.Paid" />
            <el-option label="未支付" :value="PaymentStatus.Unpaid" />
          </el-select>
        </el-form-item>
        <el-form-item label="下单时间">
          <el-date-picker
            v-model="searchForm.date_range"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>搜索
          </el-button>
          <el-button @click="handleReset">
            <el-icon><Refresh /></el-icon>重置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 操作栏 -->
    <el-card class="action-card">
      <template #header>
        <div class="card-header">
          <span>订单列表</span>
          <div class="actions">
            <el-button @click="handleExport">
              <el-icon><Download /></el-icon>导出数据
            </el-button>
          </div>
        </div>
      </template>

      <!-- 订单表格 -->
      <order-table
        ref="tableRef"
        @refresh="loadOrders"
      >
        <template #default>
          <el-table-column prop="status" label="订单状态" width="120">
            <template #default="{ row }">
              <el-tag :type="getOrderStatusType(row.status)">
                {{ getOrderStatusLabel(row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="250">
            <template #default="{ row }">
              <el-button-group>
                <el-button
                  v-if="row.status === OrderStatus.Pending"
                  size="small"
                  type="primary"
                  @click="handleUpdateStatus(row.id, OrderStatus.Shipped)"
                >
                  发货
                </el-button>
                <el-button
                  v-if="row.status === OrderStatus.Shipped"
                  size="small"
                  type="success"
                  @click="handleUpdateStatus(row.id, OrderStatus.Completed)"
                >
                  完成
                </el-button>
                <el-button
                  v-if="row.status === OrderStatus.Pending"
                  size="small"
                  type="danger"
                  @click="handleUpdateStatus(row.id, OrderStatus.Cancelled)"
                >
                  取消
                </el-button>
                <el-button size="small" @click="handleViewDetails(row.id)">
                  查看
                </el-button>
              </el-button-group>
            </template>
          </el-table-column>
        </template>
      </order-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Download } from '@element-plus/icons-vue'
import { useRouter } from 'vue-router'
import { useOrderStore } from '@/stores/order'
import { OrderStatus, PaymentStatus } from '@/types/order'
import OrderTable from '@/components/order/OrderTable.vue'

const orderStore = useOrderStore()
const tableRef = ref()
const loading = ref(false)
const router = useRouter()

// 搜索表单
const searchForm = reactive({
  order_id: '',
  status: undefined as OrderStatus | undefined,
  payment_status: undefined as PaymentStatus | undefined,
  date_range: [] as string[],
  page: 1,
  page_size: 10
})

// 订单状态选项
const orderStatusOptions = [
  { label: '待付款', value: OrderStatus.Pending },
  { label: '待发货', value: OrderStatus.Paid },
  { label: '已发货', value: OrderStatus.Shipped },
  { label: '已完成', value: OrderStatus.Completed },
  { label: '已取消', value: OrderStatus.Cancelled }
]

// 加载订单列表
const loadOrders = async () => {
  await orderStore.getOrders({
    ...searchForm,
    start_time: searchForm.date_range[0],
    end_time: searchForm.date_range[1]
  })
}

// 搜索
const handleSearch = () => {
  searchForm.page = 1
  loadOrders()
}

// 重置搜索
const handleReset = () => {
  Object.assign(searchForm, {
    order_id: '',
    status: undefined,
    payment_status: undefined,
    date_range: [],
    page: 1
  })
  loadOrders()
}

// 导出数据
const handleExport = async () => {
  // try {
  //   await orderStore.exportOrders({
  //     ...searchForm,
  //     start_time: searchForm.date_range[0],
  //     end_time: searchForm.date_range[1]
  //   })
  //   ElMessage.success('导出成功')
  // } catch (error: any) {
  //   ElMessage.error(error.message || '导出失败')
  // }
}

// 确保组件挂载时加载数据
onMounted(() => {
  loadOrders()
})

// 获取订单状态类型
const getOrderStatusType = (status: OrderStatus) => {
  switch (status) {
    case OrderStatus.Pending:
      return 'warning'
    case OrderStatus.Paid:
      return 'success'
    case OrderStatus.Shipped:
      return 'primary'
    case OrderStatus.Completed:
      return 'success'
    case OrderStatus.Cancelled:
      return 'danger'
    default:
      return 'info'
  }
}

// 获取订单状态标签
const getOrderStatusLabel = (status: OrderStatus) => {
  switch (status) {
    case OrderStatus.Pending:
      return '待付款'
    case OrderStatus.Paid:
      return '待发货'
    case OrderStatus.Shipped:
      return '已发货'
    case OrderStatus.Completed:
      return '已完成'
    case OrderStatus.Cancelled:
      return '已取消'
    default:
      return '未知'
  }
}

// 更新订单状态
const handleUpdateStatus = async (id: number, status: OrderStatus) => {
  const statusLabel = getOrderStatusLabel(status)
  try {
    await ElMessageBox.confirm(
      `确定要将订单状态更新为"${statusLabel}"吗？`,
      '确认操作',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning',
      }
    )
    
    loading.value = true
    await orderStore.updateOrderStatus(id, status)
    await loadOrders() // 刷新订单列表
    ElMessage.success(`订单状态已更新为"${statusLabel}"`)
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '更新失败')
    }
  } finally {
    loading.value = false
  }
}

// 查看订单详情
const handleViewDetails = (id: number) => {
  router.push(`/order/detail/${id}`)
}
</script>

<style scoped lang="scss">
.order-page {
  .search-card {
    margin-bottom: 20px;
  }

  .action-card {
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
  }
}
</style>
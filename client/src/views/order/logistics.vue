&lt;template>
  &lt;div class="logistics-page">
    &lt;el-card>
      &lt;template #header>
        &lt;div class="card-header">
          &lt;span>物流信息&lt;/span>
          &lt;div class="actions">
            &lt;el-button @click="$router.back()">返回&lt;/el-button>
          &lt;/div>
        &lt;/div>
      &lt;/template>

      &lt;div v-if="logistics" class="logistics-info">
        &lt;div class="logistics-header">
          &lt;el-descriptions :column="2" border>
            &lt;el-descriptions-item label="物流单号">{{ logistics.tracking_number }}&lt;/el-descriptions-item>
            &lt;el-descriptions-item label="承运商">{{ logistics.carrier }}&lt;/el-descriptions-item>
            &lt;el-descriptions-item label="物流状态">{{ logistics.status || '运输中' }}&lt;/el-descriptions-item>
            &lt;el-descriptions-item label="更新时间">{{ formatDate(new Date(logistics.updated_at)) }}&lt;/el-descriptions-item>
          &lt;/el-descriptions>
        &lt;/div>

        &lt;div class="logistics-traces">
          &lt;h3>物流轨迹&lt;/h3>
          &lt;el-timeline>
            &lt;el-timeline-item
              v-for="trace in logistics.traces"
              :key="trace.id"
              :timestamp="formatDate(new Date(trace.created_at))"
              :type="getTraceType(trace)"
            >
              &lt;h4>{{ trace.content }}&lt;/h4>
              &lt;p class="location">{{ trace.location }}&lt;/p>
            &lt;/el-timeline-item>
          &lt;/el-timeline>
        &lt;/div>
      &lt;/div>

      &lt;div v-else class="no-logistics">
        暂无物流信息
      &lt;/div>
    &lt;/el-card>
  &lt;/div>
&lt;/template>

&lt;script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useOrderStore } from '@/stores/order'
import type { Logistics, LogisticsTrace } from '@/types/order'
import { ElMessage } from 'element-plus'

const route = useRoute()
const orderStore = useOrderStore()
const logistics = ref&lt;Logistics | null>(null)

// 加载物流信息
const loadLogistics = async () => {
  try {
    const orderId = parseInt(route.params.id as string)
    if (isNaN(orderId)) {
      ElMessage.error('无效的订单ID')
      return
    }

    const order = await orderStore.getOrder(orderId)
    if (!order.logistics) {
      ElMessage.warning('该订单暂无物流信息')
      return
    }

    logistics.value = order.logistics
  } catch (error: any) {
    ElMessage.error(error.message || '加载物流信息失败')
  }
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

// 初始加载
onMounted(() => {
  loadLogistics()
})
&lt;/script>

&lt;style scoped lang="scss">
.logistics-page {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .logistics-info {
    .logistics-header {
      margin-bottom: 30px;
    }

    .logistics-traces {
      h3 {
        margin: 20px 0;
        font-weight: 500;
      }

      .el-timeline-item {
        padding-bottom: 20px;

        :deep(.el-timeline-item__timestamp) {
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

  .no-logistics {
    text-align: center;
    color: #909399;
    padding: 30px 0;
  }
}
&lt;/style>

<template>
  <div class="order-stats-container">
    <el-row :gutter="20">
      <!-- 总览数据卡片 -->
      <el-col :span="24">
        <el-card class="overview-card">
          <template #header>
            <div class="card-header">
              <span>订单总览</span>
              <el-button type="primary" link @click="refreshStats">
                <el-icon><Refresh /></el-icon> 刷新
              </el-button>
            </div>
          </template>
          <el-row :gutter="20" v-loading="loading">
            <el-col :span="6">
              <div class="stat-item">
                <div class="label">总订单数</div>
                <div class="value">{{ stats.total_count }}</div>
              </div>
            </el-col>
            <el-col :span="6">
              <div class="stat-item">
                <div class="label">总金额</div>
                <div class="value">¥{{ formatAmount(stats.total_amount) }}</div>
              </div>
            </el-col>
            <el-col :span="6">
              <div class="stat-item">
                <div class="label">今日订单数</div>
                <div class="value">{{ stats.today_count }}</div>
              </div>
            </el-col>
            <el-col :span="6">
              <div class="stat-item">
                <div class="label">今日金额</div>
                <div class="value">¥{{ formatAmount(stats.today_amount) }}</div>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="mt-20">
      <!-- 订单状态统计 -->
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>订单状态统计</span>
            </div>
          </template>
          <div class="status-stats" v-loading="loading">
            <div class="status-item">
              <el-tag>待付款</el-tag>
              <span class="count">{{ stats.pending_count }}</span>
            </div>
            <div class="status-item">
              <el-tag type="warning">已发货</el-tag>
              <span class="count">{{ stats.shipped_count }}</span>
            </div>
            <div class="status-item">
              <el-tag type="success">已完成</el-tag>
              <span class="count">{{ stats.completed_count }}</span>
            </div>
            <div class="status-item">
              <el-tag type="danger">已取消</el-tag>
              <span class="count">{{ stats.cancelled_count }}</span>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 订单趋势图 -->
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>订单趋势</span>
              <el-radio-group v-model="timeRange" size="small" @change="getOrderTrendData">
                <el-radio-button label="week">最近7天</el-radio-button>
                <el-radio-button label="month">最近30天</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div ref="chartRef" class="chart-container"></div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useOrderStore } from '@/stores/order'
import type { OrderStats } from '@/types/order'
import * as echarts from 'echarts'
import { Refresh } from '@element-plus/icons-vue'
import { storeToRefs } from 'pinia'

const orderStore = useOrderStore()
const { loading, stats } = storeToRefs(orderStore)
const timeRange = ref('week')
const chartRef = ref<HTMLElement>()
let chart: echarts.ECharts | null = null

// 格式化金额
const formatAmount = (amount: number) => {
  return amount.toFixed(2)
}

// 刷新统计数据
const refreshStats = async () => {
  await orderStore.getOrderStats()
  getOrderTrendData()
}

// 获取订单趋势数据
const getOrderTrendData = async () => {
  if (!chart) return
  const data = await orderStore.getOrderTrend(timeRange.value)
  
  const option = {
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: data.data.map(item => item.date),
      axisLabel: {
        rotate: 45
      }
    },
    yAxis: {
      type: 'value'
    },
    series: [
      {
        name: '订单数',
        type: 'bar',
        data: data.data.map(item => item.count),
        itemStyle: {
          color: '#409EFF'
        }
      },
      {
        name: '金额',
        type: 'line',
        yAxisIndex: 0,
        data: data.data.map(item => item.amount),
        itemStyle: {
          color: '#67C23A'
        }
      }
    ]
  }
  
  chart.setOption(option)
}

// 初始化图表
const initChart = () => {
  if (chartRef.value) {
    chart = echarts.init(chartRef.value)
    getOrderTrendData()
  }
}

// 监听窗口大小变化
const handleResize = () => {
  chart?.resize()
}

onMounted(() => {
  refreshStats()
  initChart()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  chart?.dispose()
})
</script>

<style scoped lang="scss">
.order-stats-container {
  .mt-20 {
    margin-top: 20px;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .overview-card {
    .stat-item {
      text-align: center;
      padding: 20px 0;

      .label {
        color: #909399;
        font-size: 14px;
        margin-bottom: 8px;
      }

      .value {
        font-size: 24px;
        font-weight: bold;
        color: #303133;
      }
    }
  }

  .status-stats {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 20px;
    padding: 10px;

    .status-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 10px;
      background-color: #f5f7fa;
      border-radius: 4px;

      .count {
        font-size: 18px;
        font-weight: bold;
        color: #303133;
      }
    }
  }

  .chart-container {
    height: 300px;
  }
}
</style>
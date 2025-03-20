<template>
  <div class="product-stats-container">
    <el-row :gutter="20">
      <!-- 总览数据卡片 -->
      <el-col :span="24">
        <el-card class="overview-card">
          <template #header>
            <div class="card-header">
              <span>商品总览</span>
              <el-button type="primary" link @click="refreshStats">
                <el-icon><Refresh /></el-icon> 刷新
              </el-button>
            </div>
          </template>
          <el-row :gutter="20" v-loading="loading">
            <el-col :span="6">
              <div class="stat-item">
                <div class="label">总商品数</div>
                <div class="value">{{ stats.total_count }}</div>
              </div>
            </el-col>
            <el-col :span="6">
              <div class="stat-item">
                <div class="label">在售商品</div>
                <div class="value">{{ stats.active_count }}</div>
              </div>
            </el-col>
            <el-col :span="6">
              <div class="stat-item">
                <div class="label">总销量</div>
                <div class="value">{{ stats.total_sales }}</div>
              </div>
            </el-col>
            <el-col :span="6">
              <div class="stat-item">
                <div class="label">总销售额</div>
                <div class="value">¥{{ formatAmount(stats.total_amount) }}</div>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="mt-20">
      <!-- 商品状态统计 -->
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>商品状态统计</span>
            </div>
          </template>
          <div class="status-stats" v-loading="loading">
            <div class="status-item">
              <el-tag type="success">在售商品</el-tag>
              <span class="count">{{ stats.active_count }}</span>
            </div>
            <div class="status-item">
              <el-tag type="warning">下架商品</el-tag>
              <span class="count">{{ stats.inactive_count }}</span>
            </div>
            <div class="status-item">
              <el-tag type="danger">缺货商品</el-tag>
              <span class="count">{{ stats.out_of_stock_count }}</span>
            </div>
            <div class="status-item">
              <el-tag>商品分类数</el-tag>
              <span class="count">{{ stats.category_count }}</span>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 商品趋势图 -->
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>销售趋势</span>
              <el-radio-group v-model="timeRange" size="small" @change="getProductTrendData">
                <el-radio-button label="week">最近7天</el-radio-button>
                <el-radio-button label="month">最近30天</el-radio-button>
                <el-radio-button label="year">最近一年</el-radio-button>
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
import { useProductStore } from '@/stores/product'
import * as echarts from 'echarts'
import { Refresh } from '@element-plus/icons-vue'
import { storeToRefs } from 'pinia'

const productStore = useProductStore()
const { stats, loading } = storeToRefs(productStore)
const timeRange = ref('week')
const chartRef = ref<HTMLElement>()
let chart: echarts.ECharts | null = null

// 格式化金额
const formatAmount = (amount: number) => {
  return (amount / 100).toFixed(2)
}

// 刷新统计数据
const refreshStats = async () => {
  await productStore.getProductStats()
}

// 获取商品趋势数据
const getProductTrendData = async () => {
  try {
    const data = await productStore.getProductTrend(timeRange.value)
    updateChart(data)
  } catch (error) {
    console.error('获取商品趋势数据失败:', error)
  }
}

// 更新图表数据
const updateChart = (data: { dates: string[], sales: number[], amounts: number[] }) => {
  if (!chart) return

  chart.setOption({
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    legend: {
      data: ['销量', '销售额']
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      data: data.dates
    },
    yAxis: [
      {
        type: 'value',
        name: '销量',
        position: 'left'
      },
      {
        type: 'value',
        name: '销售额',
        position: 'right',
        axisLabel: {
          formatter: '{value} 元'
        }
      }
    ],
    series: [
      {
        name: '销量',
        type: 'bar',
        data: data.sales
      },
      {
        name: '销售额',
        type: 'line',
        yAxisIndex: 1,
        data: data.amounts
      }
    ]
  })
}

// 初始化图表
const initChart = () => {
  if (chartRef.value) {
    chart = echarts.init(chartRef.value)
    getProductTrendData()
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

<style scoped>
.product-stats-container {
  padding: 20px;
}

.overview-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stat-item {
  text-align: center;
  padding: 20px 0;
}

.stat-item .label {
  color: #909399;
  font-size: 14px;
  margin-bottom: 10px;
}

.stat-item .value {
  font-size: 24px;
  font-weight: bold;
  color: #409EFF;
}

.status-stats {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
  padding: 20px;
}

.status-item {
  display: flex;
  align-items: center;
  gap: 10px;
}

.status-item .count {
  font-size: 20px;
  font-weight: bold;
  color: #606266;
}

.chart-container {
  height: 400px;
  width: 100%;
}

.mt-20 {
  margin-top: 20px;
}
</style>
<template>
  <div class="advertisement-stats-container">
    <el-row :gutter="20">
      <!-- 总览数据卡片 -->
      <el-col :span="24">
        <el-card class="overview-card">
          <template #header>
            <div class="card-header">
              <span>广告总览</span>
              <el-button type="primary" link @click="refreshStats">
                <el-icon><Refresh /></el-icon> 刷新
              </el-button>
            </div>
          </template>
          <el-row :gutter="20" v-loading="loading">
            <el-col :span="6">
              <div class="stat-item">
                <div class="label">总广告数</div>
                <div class="value">{{ stats.total_count }}</div>
              </div>
            </el-col>
            <el-col :span="6">
              <div class="stat-item">
                <div class="label">活跃广告</div>
                <div class="value">{{ stats.active_count }}</div>
              </div>
            </el-col>
            <el-col :span="6">
              <div class="stat-item">
                <div class="label">总点击量</div>
                <div class="value">{{ stats.click_count }}</div>
              </div>
            </el-col>
            <el-col :span="6">
              <div class="stat-item">
                <div class="label">点击率</div>
                <div class="value">{{ stats.ctr.toFixed(2) }}%</div>
              </div>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="mt-20">
      <!-- 广告位置统计 -->
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>广告位置统计</span>
            </div>
          </template>
          <div class="position-stats" v-loading="loading">
            <div class="position-item">
              <el-tag type="success">顶部横幅</el-tag>
              <span class="count">{{ stats.position_counts.banner_top }}</span>
            </div>
            <div class="position-item">
              <el-tag type="warning">底部横幅</el-tag>
              <span class="count">{{ stats.position_counts.banner_bottom }}</span>
            </div>
            <div class="position-item">
              <el-tag type="info">侧边栏顶部</el-tag>
              <span class="count">{{ stats.position_counts.sidebar_top }}</span>
            </div>
            <div class="position-item">
              <el-tag>侧边栏底部</el-tag>
              <span class="count">{{ stats.position_counts.sidebar_bottom }}</span>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 广告状态统计 -->
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>广告状态统计</span>
            </div>
          </template>
          <div class="status-stats" v-loading="loading">
            <div class="status-item">
              <el-tag type="success">活跃广告</el-tag>
              <span class="count">{{ stats.active_count }}</span>
            </div>
            <div class="status-item">
              <el-tag type="danger">非活跃广告</el-tag>
              <span class="count">{{ stats.inactive_count }}</span>
            </div>
            <div class="status-item">
              <el-tag type="info">展示次数</el-tag>
              <span class="count">{{ stats.impression_count }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="mt-20">
      <!-- 广告趋势图 -->
      <el-col :span="24">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>广告趋势</span>
              <el-radio-group v-model="timeRange" size="small" @change="getAdvertisementTrendData">
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
import { useAdvertisementStore } from '@/stores/advertisement'
import * as echarts from 'echarts'
import { Refresh } from '@element-plus/icons-vue'
import { storeToRefs } from 'pinia'

const advertisementStore = useAdvertisementStore()
const { stats, loading } = storeToRefs(advertisementStore)
const timeRange = ref('week')
const chartRef = ref<HTMLElement>()
let chart: echarts.ECharts | null = null

// 刷新统计数据
const refreshStats = async () => {
  await advertisementStore.getAdvertisementStats()
}

// 获取广告趋势数据
const getAdvertisementTrendData = async () => {
  try {
    const data = await advertisementStore.getAdvertisementTrend(timeRange.value)
    updateChart(data)
  } catch (error) {
    console.error('获取广告趋势数据失败:', error)
  }
}

// 更新图表数据
const updateChart = (data: { dates: string[], active_ads: number[], clicks: number[], impressions: number[] }) => {
  if (!chart) return

  chart.setOption({
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    legend: {
      data: ['活跃广告', '点击量', '展示量']
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
        name: '广告数',
        position: 'left'
      },
      {
        type: 'value',
        name: '点击/展示量',
        position: 'right'
      }
    ],
    series: [
      {
        name: '活跃广告',
        type: 'bar',
        data: data.active_ads
      },
      {
        name: '点击量',
        type: 'line',
        yAxisIndex: 1,
        data: data.clicks
      },
      {
        name: '展示量',
        type: 'line',
        yAxisIndex: 1,
        data: data.impressions
      }
    ]
  })
}

// 初始化图表
const initChart = () => {
  if (chartRef.value) {
    chart = echarts.init(chartRef.value)
    getAdvertisementTrendData()
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
.advertisement-stats-container {
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

.position-stats,
.status-stats {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
  padding: 20px;
}

.position-item,
.status-item {
  display: flex;
  align-items: center;
  gap: 10px;
}

.position-item .count,
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

<template>
  <div class="user-stats-container">

    <!-- 用户趋势图 -->
    <el-card class="chart-card">
      <template #header>
        <div class="card-header">
          <span>用户增长趋势</span>
          <el-radio-group v-model="timeRange" size="small">
            <el-radio-button label="week">最近一周</el-radio-button>
            <el-radio-button label="month">最近一月</el-radio-button>
            <el-radio-button label="year">最近一年</el-radio-button>
          </el-radio-group>
        </div>
      </template>
      <div class="chart-container">
        <el-empty v-if="!chartData.length" description="暂无数据" />
        <div v-else ref="chartRef" class="chart"></div>
      </div>
    </el-card>

    <!-- 用户分布 -->
    <el-row :gutter="20" class="distribution-row">
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>用户角色分布</span>
            </div>
          </template>
          <div class="chart-container">
            <el-empty v-if="!roleData.length" description="暂无数据" />
            <div v-else ref="roleChartRef" class="chart"></div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>用户状态分布</span>
            </div>
          </template>
          <div class="chart-container">
            <el-empty v-if="!statusData.length" description="暂无数据" />
            <div v-else ref="statusChartRef" class="chart"></div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import * as echarts from 'echarts'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import type { UserTrend } from '@/types/user'

const userStore = useUserStore()
const timeRange = ref('week')
const chartRef = ref<HTMLElement>()
const roleChartRef = ref<HTMLElement>()
const statusChartRef = ref<HTMLElement>()

const chartData = ref<UserTrend[]>([])
const roleData = ref<{ name: string; value: number }[]>([])
const statusData = ref<{ name: string; value: number }[]>([])

let chart: echarts.ECharts | null = null
let roleChart: echarts.ECharts | null = null
let statusChart: echarts.ECharts | null = null

// 初始化趋势图
const initChart = () => {
  if (!chartRef.value) return
  chart = echarts.init(chartRef.value)
  
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
      data: chartData.value.map(item => item.date),
      axisLabel: {
        rotate: 45
      }
    },
    yAxis: [
      {
        type: 'value',
        name: '新增用户数',
        position: 'left'
      },
      {
        type: 'value',
        name: '增长率(%)',
        position: 'right',
        axisLabel: {
          formatter: '{value}%'
        }
      }
    ],
    series: [
      {
        name: '新增用户',
        type: 'bar',
        data: chartData.value.map(item => item.count)
      },
      {
        name: '增长率',
        type: 'line',
        yAxisIndex: 1,
        data: chartData.value.map(item => item.growth),
        lineStyle: {
          width: 2
        },
        symbol: 'circle',
        symbolSize: 8
      }
    ]
  }
  
  chart.setOption(option)
}

// 初始化饼图
const initPieChart = (chart: echarts.ECharts, data: { name: string; value: number }[]) => {
  chart.setOption({
    tooltip: {
      trigger: 'item',
      formatter: '{b}: {c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      right: 10,
      top: 'center'
    },
    series: [
      {
        type: 'pie',
        radius: ['40%', '70%'],
        avoidLabelOverlap: false,
        itemStyle: {
          borderRadius: 10,
          borderColor: '#fff',
          borderWidth: 2
        },
        label: {
          show: false
        },
        emphasis: {
          label: {
            show: true,
            fontSize: '16',
            fontWeight: 'bold'
          }
        },
        labelLine: {
          show: false
        },
        data
      }
    ]
  })
}

// 加载数据
const loadData = async () => {
  try {
    // 获取用户统计数据
    await userStore.getUserStats()
    
    // 获取用户趋势数据
    const trendResponse = await userStore.getUserTrend(timeRange.value)
    chartData.value = trendResponse.data

    // 角色分布数据
    roleData.value = [
      { name: '管理员', value: userStore.userStats.admin_count },
      { name: '普通用户', value: userStore.userStats.user_count }
    ]

    // 状态分布数据
    statusData.value = [
      { name: '活跃用户', value: userStore.userStats.active_count },
      { name: '非活跃用户', value: userStore.userStats.inactive_count }
    ]

    // 初始化图表
    await nextTick()
    initChart()
    
    if (roleChartRef.value) {
      roleChart = echarts.init(roleChartRef.value)
      initPieChart(roleChart, roleData.value)
    }
    
    if (statusChartRef.value) {
      statusChart = echarts.init(statusChartRef.value)
      initPieChart(statusChart, statusData.value)
    }
  } catch (error) {
    console.error('Failed to load data:', error)
    ElMessage.error('加载数据失败')
  }
}

// 监听时间范围变化
watch(timeRange, () => {
  loadData()
})

// 处理窗口大小变化
const handleResize = () => {
  chart?.resize()
  roleChart?.resize()
  statusChart?.resize()
}

onMounted(() => {
  loadData()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  chart?.dispose()
  roleChart?.dispose()
  statusChart?.dispose()
})
</script>

<style scoped>
.user-stats-container {
  padding: 20px;
}

.chart-card {
  margin-top: 20px;
}

.distribution-row {
  margin-top: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chart-container {
  height: 400px;
  width: 100%;
}

.chart {
  height: 100%;
  width: 100%;
}
</style>
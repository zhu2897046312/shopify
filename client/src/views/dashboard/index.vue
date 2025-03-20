<template>
  <div class="dashboard-container">
    <el-row :gutter="20" class="mb-4">
      <el-col :span="6">
        <el-card shadow="hover" v-loading="userStore.loading">
          <template #header>
            <div class="card-header">
              <span>总用户数</span>
              <el-button type="primary" link @click="refreshStats">
                <el-icon><Refresh /></el-icon>
              </el-button>
            </div>
          </template>
          <div class="card-body">
            <h2>{{ userStore.userStats?.total_count }}</h2>
            <div class="card-footer">
              <div class="footer-item">
                <span>今日新增:</span>
                <span class="value">{{ userStore.userStats?.today_count }}</span>
              </div>
              <div class="footer-item" v-if="userStore.userStats?.total_growth > 0">
                <span>增长率:</span>
                <span class="value success">+{{ userStore.userStats?.total_growth }}%</span>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" v-loading="loading">
          <template #header>
            <div class="card-header">
              <span>总订单数</span>
            </div>
          </template>
          <div class="card-body">
            <h2>{{ stats.total_count }}</h2>
            <div class="card-footer">
              <div class="footer-item">
                <span>今日订单:</span>
                <span class="value">{{ stats.today_count }}</span>
              </div>
              <div class="footer-item">
                <span>待处理:</span>
                <span class="value warning">{{ stats.pending_count }}</span>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" v-loading="productStore.loading">
          <template #header>
            <div class="card-header">
              <span>商品统计</span>
            </div>
          </template>
          <div class="card-body">
            <h2>{{ productStore.stats.total_count }}</h2>
            <div class="card-footer">
              <div class="footer-item">
                <span>在售商品:</span>
                <span class="value">{{ productStore.stats.active_count }}</span>
              </div>
              <div class="footer-item">
                <span>总销量:</span>
                <span class="value success">{{ productStore.stats.total_sales }}</span>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" v-loading="loading">
          <template #header>
            <div class="card-header">
              <span>总销售额</span>
            </div>
          </template>
          <div class="card-body">
            <h2>¥{{ formatAmount(stats.total_amount) }}</h2>
            <div class="card-footer">
              <div class="footer-item">
                <span>今日销售:</span>
                <span class="value">¥{{ formatAmount(stats.today_amount) }}</span>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20">
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>订单趋势</span>
              <el-radio-group v-model="orderTimeRange" size="small" @change="refreshOrderTrend">
                <el-radio-button label="week">最近7天</el-radio-button>
                <el-radio-button label="month">最近30天</el-radio-button>
                <el-radio-button label="year">最近一年</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div class="chart-container">
            <el-empty v-if="!orderTrendData.length" description="暂无数据" />
            <div v-else ref="orderChartRef" style="width: 100%; height: 350px"></div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>用户趋势</span>
              <el-radio-group v-model="userTimeRange" size="small" @change="refreshUserTrend">
                <el-radio-button label="week">最近7天</el-radio-button>
                <el-radio-button label="month">最近30天</el-radio-button>
                <el-radio-button label="year">最近一年</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div class="chart-container">
            <el-empty v-if="!userTrendData.length" description="暂无数据" />
            <div v-else ref="userChartRef" style="width: 100%; height: 350px"></div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="mt-4">
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>商品销售趋势</span>
              <el-radio-group v-model="productTimeRange" size="small" @change="refreshProductTrend">
                <el-radio-button label="week">最近7天</el-radio-button>
                <el-radio-button label="month">最近30天</el-radio-button>
                <el-radio-button label="year">最近一年</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div class="chart-container">
            <el-empty v-if="!productTrendData.dates?.length" description="暂无数据" />
            <div v-else ref="productChartRef" style="width: 100%; height: 350px"></div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>广告效果趋势</span>
              <el-radio-group v-model="adTimeRange" size="small" @change="refreshAdTrend">
                <el-radio-button label="week">最近7天</el-radio-button>
                <el-radio-button label="month">最近30天</el-radio-button>
                <el-radio-button label="year">最近一年</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div class="chart-container">
            <el-empty v-if="!adTrendData.dates?.length" description="暂无数据" />
            <div v-else ref="adChartRef" style="width: 100%; height: 350px"></div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { useOrderStore } from '@/stores/order'
import { useProductStore } from '@/stores/product'
import { useAdvertisementStore } from '@/stores/advertisement'
import type { OrderStats } from '@/types/order'
import * as echarts from 'echarts'
import { storeToRefs } from 'pinia'
import { Refresh } from '@element-plus/icons-vue'

const userStore = useUserStore()
const orderStore = useOrderStore()
const productStore = useProductStore()
const advertisementStore = useAdvertisementStore()
const { loading, stats } = storeToRefs(orderStore)

// 图表相关
const orderChartRef = ref<HTMLElement>()
const userChartRef = ref<HTMLElement>()
const productChartRef = ref<HTMLElement>()
const adChartRef = ref<HTMLElement>()

const orderTimeRange = ref('week')
const userTimeRange = ref('week')
const productTimeRange = ref('week')
const adTimeRange = ref('week')

const orderTrendData = ref<any[]>([])
const userTrendData = ref<any[]>([])
const productTrendData = ref<any>({})
const adTrendData = ref<any>({})

let orderChart: echarts.ECharts | null = null
let userChart: echarts.ECharts | null = null
let productChart: echarts.ECharts | null = null
let adChart: echarts.ECharts | null = null

// 格式化金额
const formatAmount = (amount: number) => {
  return amount.toFixed(2)
}

// 刷新统计数据
const refreshStats = async () => {
  await Promise.all([
    userStore.getUserStats(),
    orderStore.getOrderStats(),
    productStore.getProductStats(),
    advertisementStore.getAdvertisementStats()
  ])
  refreshAllTrends()
}

// 刷新所有趋势数据
const refreshAllTrends = () => {
  refreshOrderTrend()
  refreshUserTrend()
  refreshProductTrend()
  refreshAdTrend()
}

// 刷新订单趋势
const refreshOrderTrend = async () => {
  try {
    const { data } = await orderStore.getOrderTrend(orderTimeRange.value)
    orderTrendData.value = data || []
    if (!orderChartRef.value) return
    
    if (!orderChart) {
      orderChart = echarts.init(orderChartRef.value)
    }
    
    orderChart.setOption({
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'shadow'
        }
      },
      legend: {
        data: ['订单数', '金额']
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: data.map(item => item.date),
        axisLabel: {
          rotate: 45
        }
      },
      yAxis: [
        {
          type: 'value',
          name: '订单数',
          position: 'left'
        },
        {
          type: 'value',
          name: '金额',
          position: 'right',
          axisLabel: {
            formatter: '¥{value}'
          }
        }
      ],
      series: [
        {
          name: '订单数',
          type: 'bar',
          data: data.map(item => item.count),
          itemStyle: {
            color: '#409EFF'
          }
        },
        {
          name: '金额',
          type: 'line',
          yAxisIndex: 1,
          data: data.map(item => item.amount),
          itemStyle: {
            color: '#67C23A'
          }
        }
      ]
    })
  } catch (error) {
    console.error('Failed to get order trend:', error)
  }
}

// 刷新用户趋势
const refreshUserTrend = async () => {
  try {
    const { data } = await userStore.getUserTrend(userTimeRange.value)
    userTrendData.value = data || []
    if (!userChartRef.value) return
    
    if (!userChart) {
      userChart = echarts.init(userChartRef.value)
    }

    userChart.setOption({
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'shadow'
        }
      },
      legend: {
        data: ['新增用户', '增长率']
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        data: data.map(item => item.date),
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
          data: data.map(item => item.count),
          itemStyle: {
            color: '#409EFF'
          }
        },
        {
          name: '增长率',
          type: 'line',
          yAxisIndex: 1,
          data: data.map(item => item.growth),
          itemStyle: {
            color: '#67C23A'
          },
          lineStyle: {
            width: 2
          },
          symbol: 'circle',
          symbolSize: 8
        }
      ]
    })
  } catch (error) {
    console.error('Failed to get user trend:', error)
  }
}

// 刷新商品趋势
const refreshProductTrend = async () => {
  try {
    const data = await productStore.getProductTrend(productTimeRange.value)
    productTrendData.value = data || {}
    if (!productChartRef.value) return
    
    if (!productChart) {
      productChart = echarts.init(productChartRef.value)
    }

    const dates = data.dates || []
    const sales = data.sales || []
    const amounts = data.amounts || []

    productChart.setOption({
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
        data: dates,
        axisLabel: {
          rotate: 45
        }
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
            formatter: '¥{value}'
          }
        }
      ],
      series: [
        {
          name: '销量',
          type: 'bar',
          data: sales,
          itemStyle: {
            color: '#409EFF'
          }
        },
        {
          name: '销售额',
          type: 'line',
          yAxisIndex: 1,
          data: amounts,
          itemStyle: {
            color: '#67C23A'
          }
        }
      ]
    })
  } catch (error) {
    console.error('Failed to get product trend:', error)
  }
}

// 刷新广告趋势
const refreshAdTrend = async () => {
  try {
    const data = await advertisementStore.getAdvertisementTrend(adTimeRange.value)
    adTrendData.value = data || {}
    if (!adChartRef.value) return
    
    if (!adChart) {
      adChart = echarts.init(adChartRef.value)
    }

    const dates = data.dates || []
    const active_ads = data.active_ads || []
    const clicks = data.clicks || []
    const impressions = data.impressions || []

    adChart.setOption({
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
        data: dates,
        axisLabel: {
          rotate: 45
        }
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
          data: active_ads,
          itemStyle: {
            color: '#409EFF'
          }
        },
        {
          name: '点击量',
          type: 'line',
          yAxisIndex: 1,
          data: clicks,
          itemStyle: {
            color: '#67C23A'
          }
        },
        {
          name: '展示量',
          type: 'line',
          yAxisIndex: 1,
          data: impressions,
          itemStyle: {
            color: '#E6A23C'
          }
        }
      ]
    })
  } catch (error) {
    console.error('Failed to get advertisement trend:', error)
  }
}

// 处理窗口大小变化
const handleResize = () => {
  orderChart?.resize()
  userChart?.resize()
  productChart?.resize()
  adChart?.resize()
}

// 初始化
onMounted(async () => {
  await refreshStats()
  await refreshAllTrends()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  orderChart?.dispose()
  userChart?.dispose()
  productChart?.dispose()
  adChart?.dispose()
})
</script>

<style scoped lang="scss">
.dashboard-container {
  padding: 20px;

  .mb-4 {
    margin-bottom: 20px;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .card-body {
    text-align: center;
    
    h2 {
      font-size: 28px;
      margin: 10px 0;
      color: #303133;
    }

    .card-footer {
      margin-top: 10px;
      
      .footer-item {
        display: flex;
        justify-content: space-between;
        margin-top: 5px;
        color: #909399;

        .value {
          font-weight: bold;
          color: #606266;

          &.success {
            color: #67C23A;
          }

          &.warning {
            color: #E6A23C;
          }

          &.danger {
            color: #F56C6C;
          }
        }
      }
    }
  }

  .chart-container {
    width: 100%;
    height: 350px;
  }
}
</style>
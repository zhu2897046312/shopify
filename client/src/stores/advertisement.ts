import { defineStore } from 'pinia'
import { ref } from 'vue'
import { storeToRefs } from 'pinia'
import type {
  Advertisement,
  ListAdvertisementParams,
  CreateAdvertisementRequest,
  UpdateAdvertisementRequest,
  UpdateAdvertisementStatusRequest
} from '@/types/advertisement'
import { AdvertisementStatus, AdvertisementPosition } from '@/types/advertisement'
import * as advertisementApi from '@/api/advertisement'
import { ElMessage } from 'element-plus'

// 广告统计数据接口
interface AdvertisementStats {
  total_count: number
  active_count: number
  inactive_count: number
  position_counts: Record<AdvertisementPosition, number>
  click_count: number
  impression_count: number
  ctr: number // Click-Through Rate
}

export const useAdvertisementStore = defineStore('admin-advertisement', () => {
  // 状态
  const advertisements = ref<Advertisement[]>([])
  const total = ref(0)
  const loading = ref(false)
  const currentAdvertisement = ref<Advertisement | null>(null)
  const stats = ref<AdvertisementStats>({
    total_count: 0,
    active_count: 0,
    inactive_count: 0,
    position_counts: {
      [AdvertisementPosition.BannerTop]: 0,
      [AdvertisementPosition.BannerBottom]: 0,
      [AdvertisementPosition.SidebarTop]: 0,
      [AdvertisementPosition.SidebarBottom]: 0
    },
    click_count: 0,
    impression_count: 0,
    ctr: 0
  })

  // 获取广告列表
  const getAdvertisements = async (params: ListAdvertisementParams) => {
    loading.value = true
    try {
      const res = await advertisementApi.getAdvertisements(params)
      advertisements.value = res.items
      total.value = res.total
      return res
    } finally {
      loading.value = false
    }
  }

  // 获取广告详情
  const getAdvertisement = async (id: number) => {
    loading.value = true
    try {
      const advertisement = await advertisementApi.getAdvertisement(id)
      currentAdvertisement.value = advertisement
      return advertisement
    } finally {
      loading.value = false
    }
  }

  // 创建广告
  const createAdvertisement = async (data: CreateAdvertisementRequest) => {
    loading.value = true
    try {
      const advertisement = await advertisementApi.createAdvertisement(data)
      ElMessage.success('创建广告成功')
      return advertisement
    } finally {
      loading.value = false
    }
  }

  // 更新广告
  const updateAdvertisement = async (id: number, data: UpdateAdvertisementRequest) => {
    loading.value = true
    try {
      const advertisement = await advertisementApi.updateAdvertisement(id, data)
      ElMessage.success('更新广告成功')
      return advertisement
    } finally {
      loading.value = false
    }
  }

  // 删除广告
  const deleteAdvertisement = async (id: number) => {
    loading.value = true
    try {
      await advertisementApi.deleteAdvertisement(id)
      ElMessage.success('删除广告成功')
    } finally {
      loading.value = false
    }
  }

  // 更新广告状态
  const updateAdvertisementStatus = async (id: number, status: AdvertisementStatus) => {
    loading.value = true
    try {
      const data: UpdateAdvertisementStatusRequest = { status }
      const advertisement = await advertisementApi.updateAdvertisementStatus(id, data)
      ElMessage.success('更新广告状态成功')
      return advertisement
    } finally {
      loading.value = false
    }
  }

  // 获取广告统计数据
  const getAdvertisementStats = async () => {
    try {
      // 获取所有广告数据
      const allAds = await advertisementApi.getAdvertisements({ page: 1, page_size: 1000 })
      
      // 计算各种状态的广告数量
      const activeAds = allAds.items.filter(ad => ad.status === AdvertisementStatus.Active)
      const inactiveAds = allAds.items.filter(ad => ad.status === AdvertisementStatus.Inactive)
      
      // 计算各个位置的广告数量
      const positionCounts = {
        [AdvertisementPosition.BannerTop]: 0,
        [AdvertisementPosition.BannerBottom]: 0,
        [AdvertisementPosition.SidebarTop]: 0,
        [AdvertisementPosition.SidebarBottom]: 0
      }
      
      allAds.items.forEach(ad => {
        positionCounts[ad.position]++
      })

      // 模拟点击和展示数据（实际应该从后端获取）
      const clickCount = activeAds.length * 100 // 假设每个活跃广告平均100次点击
      const impressionCount = activeAds.length * 1000 // 假设每个活跃广告平均1000次展示
      
      // 更新统计数据
      stats.value = {
        total_count: allAds.total,
        active_count: activeAds.length,
        inactive_count: inactiveAds.length,
        position_counts: positionCounts,
        click_count: clickCount,
        impression_count: impressionCount,
        ctr: impressionCount > 0 ? (clickCount / impressionCount) * 100 : 0
      }
      
      return stats.value
    } catch (error) {
      console.error('获取广告统计数据失败:', error)
      ElMessage.error('获取广告统计数据失败')
      throw error
    }
  }

  // 获取广告趋势数据
  const getAdvertisementTrend = async (timeRange: string) => {
    try {
      // 获取所有广告
      const { items: allAds } = await advertisementApi.getAdvertisements({ page: 1, page_size: 1000 })

      // 根据时间范围计算日期范围
      const now = new Date()
      let startDate = new Date()
      let dateFormat: string
      let step: number

      switch (timeRange) {
        case 'week':
          startDate.setDate(now.getDate() - 7)
          dateFormat = 'MM-DD'
          step = 24 * 60 * 60 * 1000 // 1 day
          break
        case 'month':
          startDate.setMonth(now.getMonth() - 1)
          dateFormat = 'MM-DD'
          step = 24 * 60 * 60 * 1000 // 1 day
          break
        case 'year':
          startDate.setFullYear(now.getFullYear() - 1)
          dateFormat = 'YYYY-MM'
          step = 30 * 24 * 60 * 60 * 1000 // 30 days
          break
        default:
          throw new Error('Invalid time range')
      }

      // 生成日期数组和数据
      const dates: string[] = []
      const activeAds: number[] = []
      const clicks: number[] = []
      const impressions: number[] = []

      for (let date = startDate; date <= now; date = new Date(date.getTime() + step)) {
        const dateStr = formatDate(date, dateFormat)
        dates.push(dateStr)

        // 统计当前日期的活跃广告数
        const periodStart = date.getTime()
        const periodEnd = new Date(date.getTime() + step).getTime()

        const periodAds = allAds.filter(ad => {
          const startTime = new Date(ad.start_time).getTime()
          const endTime = new Date(ad.end_time).getTime()
          return startTime <= periodEnd && endTime >= periodStart && ad.status === AdvertisementStatus.Active
        })

        // 计算当期数据
        const activeCount = periodAds.length
        // 模拟点击和展示数据
        const clickCount = activeCount * 100
        const impressionCount = activeCount * 1000

        activeAds.push(activeCount)
        clicks.push(clickCount)
        impressions.push(impressionCount)
      }

      return {
        dates,
        active_ads: activeAds,
        clicks,
        impressions
      }
    } catch (error) {
      console.error('获取广告趋势数据失败:', error)
      ElMessage.error('获取广告趋势数据失败')
      throw error
    }
  }

  // 格式化日期
  const formatDate = (date: Date, format: string): string => {
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    return format
      .replace('YYYY', String(year))
      .replace('MM', month)
      .replace('DD', day)
  }

  return {
    // 状态
    advertisements,
    total,
    loading,
    currentAdvertisement,
    stats,

    // 方法
    getAdvertisements,
    getAdvertisement,
    createAdvertisement,
    updateAdvertisement,
    deleteAdvertisement,
    updateAdvertisementStatus,
    getAdvertisementStats,
    getAdvertisementTrend,
    formatDate
  }
})

// 导出组合式函数，用于在组件中获取响应式状态
export const useAdvertisementStoreRefs = () => {
  const store = useAdvertisementStore()
  return storeToRefs(store)
}
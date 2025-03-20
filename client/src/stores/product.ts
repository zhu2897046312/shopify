import { defineStore } from 'pinia'
import { ref } from 'vue'
import { storeToRefs } from 'pinia'
import type {
  Product,
  ProductStats,
  ListProductParams,
  CreateProductRequest,
  UpdateProductRequest,
  ListReviewsParams
} from '@/types/product'
import { ProductStatus } from '@/types/product'
import * as productApi from '@/api/product'
import { ElMessage } from 'element-plus'

export const useProductStore = defineStore('admin-product', () => {
  // 状态
  const products = ref<Product[]>([])
  const total = ref(0)
  const loading = ref(false)
  const currentProduct = ref<Product | null>(null)
  const stats = ref<ProductStats>({
    total_count: 0,
    active_count: 0,
    inactive_count: 0,
    out_of_stock_count: 0,
    category_count: 0,
    total_sales: 0,
    total_amount: 0
  })

  // 获取商品列表
  const getProducts = async (params: ListProductParams) => {
    loading.value = true
    try {
      const res = await productApi.getProducts(params)
      products.value = res.items
      total.value = res.total
      return res
    } finally {
      loading.value = false
    }
  }

  // 获取商品评论
  const getProductReviewsList = async (product_id: number, params: ListReviewsParams) => {
    loading.value = true
    try {
      const res = await productApi.getProductReviews(product_id, params)
      return res
    } finally {
      loading.value = false
    }
  }

  // 删除商品评论
  const deleteProductReview = async (review_id: number) => {
    loading.value = true
    try {
      await productApi.deleteProductReview(review_id)
      ElMessage.success('删除评论成功')
    } finally {
      loading.value = false
    }
  }

  // 获取商品详情
  const getProduct = async (id: number) => {
    loading.value = true
    try {
      const product = await productApi.getProduct(id)
      currentProduct.value = product
      return product
    } finally {
      loading.value = false
    }
  }

  // 创建商品
  const createProduct = async (data: CreateProductRequest) => {
    loading.value = true
    try {
      const product = await productApi.createProduct(data)
      ElMessage.success('创建商品成功')
      return product
    } finally {
      loading.value = false
    }
  }

  // 更新商品
  const updateProduct = async (id: number, data: UpdateProductRequest) => {
    loading.value = true
    try {
      const product = await productApi.updateProduct(id, data)
      ElMessage.success('更新商品成功')
      return product
    } finally {
      loading.value = false
    }
  }

  // 删除商品
  const deleteProduct = async (id: number) => {
    loading.value = true
    try {
      await productApi.deleteProduct(id)
      ElMessage.success('删除商品成功')
    } finally {
      loading.value = false
    }
  }

  // 获取商品统计数据
  const getProductStats = async () => {
    try {
      // 获取所有商品数据
      const allProducts = await productApi.getProducts({ page: 1, page_size: 1000 })
      
      // 计算各种状态的商品数量
      const activeProducts = allProducts.items.filter(p => p.status === ProductStatus.Active)
      const inactiveProducts = allProducts.items.filter(p => p.status === ProductStatus.Inactive)
      const outOfStockProducts = allProducts.items.filter(p => p.stock === 0)
      
      // 计算分类数量（通过Set去重）
      const categories = new Set(allProducts.items.map(p => p.category))
      
      // 计算总销售额
      const totalAmount = allProducts.items.reduce((sum, p) => sum + (p.price * p.sales), 0)
      
      // 更新统计数据
      stats.value = {
        total_count: allProducts.total,
        active_count: activeProducts.length,
        inactive_count: inactiveProducts.length,
        out_of_stock_count: outOfStockProducts.length,
        category_count: categories.size,
        total_sales: allProducts.items.reduce((sum, p) => sum + p.sales, 0),
        total_amount: totalAmount
      }
      
      return stats.value
    } catch (error) {
      console.error('获取商品统计数据失败:', error)
      ElMessage.error('获取商品统计数据失败')
      throw error
    }
  }

  // 获取商品趋势数据（按时间范围）
  const getProductTrend = async (timeRange: string) => {
    try {
      // 获取所有商品
      const { items: allProducts } = await productApi.getProducts({ page: 1, page_size: 1000 })

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
      const sales: number[] = []
      const amounts: number[] = []

      for (let date = startDate; date <= now; date = new Date(date.getTime() + step)) {
        const dateStr = formatDate(date, dateFormat)
        dates.push(dateStr)

        // 统计当前日期的销量和销售额
        const periodStart = date.getTime()
        const periodEnd = new Date(date.getTime() + step).getTime()

        // 假设商品的 updated_at 表示最后一次销售时间
        const periodProducts = allProducts.filter(product => {
          const updateTime = new Date(product.updated_at).getTime()
          return updateTime >= periodStart && updateTime < periodEnd
        })

        // 计算当期销量和销售额
        const periodSales = periodProducts.reduce((sum, product) => sum + product.sales, 0)
        const periodAmount = periodProducts.reduce((sum, product) => sum + (product.price * product.sales), 0)

        sales.push(periodSales)
        amounts.push(periodAmount)
      }

      return {
        dates,
        sales,
        amounts
      }
    } catch (error) {
      console.error('获取商品趋势数据失败:', error)
      ElMessage.error('获取商品趋势数据失败')
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
    products,
    total,
    loading,
    currentProduct,
    stats,
    
    // 方法
    getProducts,
    getProduct,
    createProduct,
    updateProduct,
    deleteProduct,
    getProductStats,
    getProductTrend,
    formatDate,
    getProductReviewsList,
    deleteProductReview
  }
})

// 导出组合式函数，用于在组件中获取响应式状态
export const useProductStoreRefs = () => {
  const store = useProductStore()
  return storeToRefs(store)
}
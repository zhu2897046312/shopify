<template>
  <div class="product-page">
    <!-- 搜索栏 -->
    <el-card class="search-card">
      <el-form :model="searchForm" inline>
        <el-form-item label="关键词">
          <el-input
            v-model="searchForm.keyword"
            placeholder="商品名称/描述"
            clearable
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item label="商品分类">
          <el-select v-model="searchForm.category" placeholder="全部" clearable>
            <el-option
              v-for="(value, key) in ProductCategory"
              :key="key"
              :label="ProductCategoryNames[value as ProductCategory]"
              :value="value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="上架" :value="ProductStatus.Active" />
            <el-option label="下架" :value="ProductStatus.Inactive" />
          </el-select>
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
          <span>商品列表</span>
          <div class="actions">
            <el-button type="primary" @click="handleAdd">
              <el-icon><Plus /></el-icon>新增商品
            </el-button>
            <el-button type="danger" :disabled="!selectedProducts.length" @click="handleBatchDelete">
              批量删除
            </el-button>
            <el-button @click="handleExport">
              <el-icon><Download /></el-icon>导出数据
            </el-button>
          </div>
        </div>
      </template>

      <!-- 商品表格 -->
      <product-table
        ref="tableRef"
        @selection-change="handleSelectionChange"
        @view-reviews="handleViewReviews"
      />
    </el-card>

    <!-- 评论对话框 -->
    <el-dialog
      v-model="reviewDialogVisible"
      title="商品评论"
      width="800px"
    >
      <el-table
        v-loading="loading"
        :data="reviews"
        style="width: 100%"
      >
        <el-table-column prop="user.username" label="用户" width="120" />
        <el-table-column prop="rating" label="评分" width="100">
          <template #default="{ row }">
            <el-rate
              v-model="row.rating"
              disabled
              show-score
            />
          </template>
        </el-table-column>
        <el-table-column prop="content" label="评论内容" />
        <el-table-column prop="created_at" label="评论时间" width="180">
          <template #default="{ row }">
            {{ formatDate(new Date(row.created_at)) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button
              type="danger"
              link
              @click="handleDeleteReview(row.id)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-model:current-page="reviewPage"
        v-model:page-size="reviewPageSize"
        :total="reviewTotal"
        class="mt-4"
        background
        layout="total, sizes, prev, pager, next"
        @size-change="loadReviews"
        @current-change="loadReviews"
      />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessageBox, ElMessage } from 'element-plus'
import { Search, Refresh, Plus, Download } from '@element-plus/icons-vue'
import { useProductStore } from '@/stores/product'
import type { Product, Review } from '@/types/product'
import { ProductStatus, ProductCategoryNames, ProductCategory } from '@/types/product'
import ProductTable from '@/components/product/ProductTable.vue'

const router = useRouter()
const productStore = useProductStore()
const tableRef = ref()

const searchForm = ref({
  keyword: '',
  category: undefined as ProductCategory | undefined,
  category_id: undefined as number | undefined,
  status: undefined as ProductStatus | undefined,
  page: 1,
  page_size: 10
})

const categories = ref<ProductCategory[]>([])
const selectedProducts = ref<number[]>([])

// 评论相关
const reviewDialogVisible = ref(false)
const reviews = ref<Review[]>([])
const reviewTotal = ref(0)
const reviewPage = ref(1)
const reviewPageSize = ref(10)
const currentProductId = ref<number | null>(null)
const loading = ref(false)

// 格式化日期
const formatDate = (date: Date) => {
  return new Intl.DateTimeFormat('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit'
  }).format(date)
}

// 加载评论数据
const loadReviews = async () => {
  if (!currentProductId.value) return
  
  try {
    const res = await productStore.getProductReviewsList(currentProductId.value, {
      page: reviewPage.value,
      page_size: reviewPageSize.value
    })
    reviews.value = res.reviews
    reviewTotal.value = res.total
  } catch (error) {
    ElMessage.error('获取评论失败')
  }
}

// 查看评论
const handleViewReviews = async (productId: number) => {
  currentProductId.value = productId
  reviewDialogVisible.value = true
  reviewPage.value = 1
  await loadReviews()
}

// 删除评论
const handleDeleteReview = async (reviewId: number) => {
  try {
    await ElMessageBox.confirm('确定要删除这条评论吗？', '提示', {
      type: 'warning'
    })
    
    // TODO: 实现删除评论的 API 调用
    await productStore.deleteProductReview(reviewId)
    ElMessage.success('删除评论成功')
    await loadReviews()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除评论失败')
    }
  }
}

// 搜索 
const handleSearch = () => {
  searchForm.value.page = 1
  loadProducts()
}

// 重置
const handleReset = () => {
  searchForm.value = {
    keyword: '',
    category: undefined,
    category_id: undefined,
    status: undefined,
    page: 1,
    page_size: 10
  }
  loadProducts()
}

// 加载商品数据
const loadProducts = () => {
  productStore.getProducts(searchForm.value)
}

// 新增商品
const handleAdd = () => {
  router.push('/products/create')
}

// 选择变化
const handleSelectionChange = (selection: Product[]) => {
  selectedProducts.value = selection.map(item => item.id)
}

// 批量删除
const handleBatchDelete = async () => {
  // try {
  //   await ElMessageBox.confirm('确定要删除选中的商品吗？', '提示', {
  //     type: 'warning'
  //   })
  //   await productStore.deleteProduct(selectedProducts.value)
  //   ElMessage.success('删除成功')
  //   loadProducts()
  // } catch {
  //   // 用户取消操作
  // }
  console.log(selectedProducts.value)
}

// 导出数据
const handleExport = async () => {
  // try {
  //   await productStore.exportProducts(searchForm.value)
  //   ElMessage.success('导出成功')
  // } catch (error: any) {
  //   ElMessage.error(error.message || '导出失败')
  // }
  console.log(searchForm.value)
}

onMounted(() => {
  loadProducts()
})
</script>

<style scoped lang="scss">
.product-page {
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
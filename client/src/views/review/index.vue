<template>
  <div class="review-list-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>评论管理</span>
          <el-button-group>
            <el-button 
              type="danger" 
              :disabled="!selectedReviews.length"
              @click="handleBatchDelete"
            >
              批量删除
            </el-button>
          </el-button-group>
        </div>
      </template>

      <el-table
        v-loading="loading"
        :data="reviews"
        border
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="用户" width="200">
          <template #default="{ row }">
            <div class="user-info">
              <el-avatar :size="32" :src="row.user.avatar">
                {{ row.user.nickname?.charAt(0)?.toUpperCase() }}
              </el-avatar>
              <span class="nickname">{{ row.user.nickname }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="商品" min-width="200">
          <template #default="{ row }">
            <el-link @click="$router.push(`/products/${row.product_id}`)">
              {{ row.product.name }}
            </el-link>
          </template>
        </el-table-column>
        <el-table-column label="评分" width="120">
          <template #default="{ row }">
            <el-rate v-model="row.rating" disabled />
          </template>
        </el-table-column>
        <el-table-column prop="content" label="评价内容" show-overflow-tooltip />
        <el-table-column label="图片" width="120">
          <template #default="{ row }">
            <el-image
              v-if="row.images?.length"
              :src="row.images[0]"
              :preview-src-list="row.images"
              fit="cover"
              style="width: 60px; height: 60px"
            />
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="评价时间" width="180" />
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button type="danger" link @click="handleDelete(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessageBox } from 'element-plus'
import { useReviewStore } from '@/stores/review'

const store = useReviewStore()
const { reviews, total, loading } = storeToRefs(store)

const currentPage = ref(1)
const pageSize = ref(10)
const selectedReviews = ref<number[]>([])

// 加载数据
const loadData = () => {
  store.getReviews({
    page: currentPage.value,
    page_size: pageSize.value
  })
}

// 选择变化
const handleSelectionChange = (selection: any[]) => {
  selectedReviews.value = selection.map(item => item.id)
}

// 删除评论
const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定要删除该评论吗？', '提示', {
      type: 'warning'
    })
    await store.deleteReview(row.id)
    loadData()
  } catch {
    // 用户取消操作
  }
}

// 批量删除
const handleBatchDelete = async () => {
  try {
    await ElMessageBox.confirm('确定要删除选中的评论吗？', '提示', {
      type: 'warning'
    })
    await store.batchDeleteReviews(selectedReviews.value)
    selectedReviews.value = []
    loadData()
  } catch {
    // 用户取消操作
  }
}

// 分页相关
const handleSizeChange = () => {
  currentPage.value = 1
  loadData()
}

const handleCurrentChange = () => {
  loadData()
}

onMounted(() => {
  loadData()
})
</script>

<style scoped lang="scss">
.review-list-container {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .user-info {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .pagination {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}
</style> 
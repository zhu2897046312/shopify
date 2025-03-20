<template>
  <div class="product-table">
    <el-table
      v-loading="loading"
      :data="products"
      border
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column label="商品信息" min-width="300">
        <template #default="{ row }">
          <div class="product-info">
            <el-image
              :src="row.images[0]"
              :preview-src-list="row.images"
              fit="cover"
              class="product-image"
            />
            <div class="product-detail">
              <div class="name">{{ row.name }}</div>
              <div class="category">{{ row.category?.name }}</div>
            </div>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="price" label="价格" width="120">
        <template #default="{ row }">
          ¥{{ row.price }}
        </template>
      </el-table-column>
      <el-table-column prop="stock" label="库存" width="100" />
      <el-table-column prop="sales" label="销量" width="100" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 'active' ? 'success' : 'info'">
            {{ row.status === 'active' ? '上架' : '下架' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button-group>
            <el-button type="primary" link @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button 
              type="primary" 
              link 
              @click="handleStatusChange(row)"
            >
              {{ row.status === 'active' ? '下架' : '上架' }}
            </el-button>
            <el-button type="primary" link @click="handleViewReviews(row)">
              评论
            </el-button>
            <el-button type="danger" link @click="handleDelete(row)">
              删除
            </el-button>
          </el-button-group>
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
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessageBox } from 'element-plus'
import { useProductStore } from '@/stores/product'
import { storeToRefs } from 'pinia'
import type { Product } from '@/types/product'
import type { ListProductParams } from '@/types/product'
import { ProductStatus } from '@/types/product'

const router = useRouter()
const productStore = useProductStore()
const { products, total, loading } = storeToRefs(productStore)

const emit = defineEmits<{
  (e: 'selection-change', selection: Product[]): void
  (e: 'view-reviews', productId: number): void
}>()

const currentPage = ref(1)
const pageSize = ref(10)
const selectedProducts = ref<number[]>([])

// 加载商品列表
const loadProducts = (params: ListProductParams) => {
  productStore.getProducts({
  ...params,  // 先展开params
  page: currentPage.value,  // 后面的值会覆盖前面的
  page_size: pageSize.value
})
}

// 选择变化
const handleSelectionChange = (selection: Product[]) => {
  selectedProducts.value = selection.map(item => item.id)
  emit('selection-change', selection)
}

// 编辑商品
const handleEdit = (row: Product) => {
  router.push(`/products/${row.id}/edit`)
}

// 更改状态
const handleStatusChange = async (row: Product) => {
  try {
    await ElMessageBox.confirm(
      `确定要${row.status === 'active' ? '下架' : '上架'}该商品吗？`,
      '提示',
      { type: 'warning' }
    )
    await productStore.updateProduct(row.id, { 
      status: row.status === ProductStatus.Active ? ProductStatus.Inactive : ProductStatus.Active 
    })
    loadProducts({
      page: currentPage.value,
      page_size: pageSize.value
    })
  } catch {
    // 用户取消操作
  }
}

// 查看评论
const handleViewReviews = (row: Product) => {
  emit('view-reviews', row.id)
}

// 删除商品
const handleDelete = async (row: Product) => {
  try {
    await ElMessageBox.confirm('确定要删除该商品吗？', '提示', {
      type: 'warning'
    })
    await productStore.deleteProduct(row.id)
    loadProducts({
      page: currentPage.value,
      page_size: pageSize.value
    })
  } catch {
    // 用户取消操作
  }
}

// 分页相关
const handleSizeChange = () => {
  currentPage.value = 1
  loadProducts({
    page: currentPage.value,
    page_size: pageSize.value
  })
}

const handleCurrentChange = () => {
  loadProducts({
    page: currentPage.value,
    page_size: pageSize.value
  })
}

// 初始加载
loadProducts({
  page: currentPage.value,
  page_size: pageSize.value
})

// 暴露方法给父组件
defineExpose({
  loadProducts
})
</script>

<style scoped lang="scss">
.product-table {
  .product-info {
    display: flex;
    align-items: center;
    gap: 12px;

    .product-image {
      width: 60px;
      height: 60px;
      border-radius: 4px;
    }

    .product-detail {
      .name {
        font-weight: bold;
      }
      .category {
        font-size: 12px;
        color: #666;
      }
    }
  }

  .pagination {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>
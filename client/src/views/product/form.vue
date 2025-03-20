<template>
  <div class="product-form-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ isEdit ? '编辑商品' : '新增商品' }}</span>
          <el-button @click="$router.back()">返回</el-button>
        </div>
      </template>

      <product-form
        :product="product"
        @success="handleSuccess"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useProductStore } from '@/stores/product'
import type { Product } from '@/types/product'
import ProductForm from '@/components/product/ProductForm.vue'

const route = useRoute()
const router = useRouter()
const productStore = useProductStore()

const isEdit = computed(() => route.name === 'EditProduct')
const product = ref<Product>()

// 加载商品数据
const loadProduct = async () => {
  if (!isEdit.value) return
  
  const id = Number(route.params.id)
  if (!id) {
    router.push('/products')
    return
  }

  product.value = await productStore.getProduct(id)
  if (!product.value) {
    ElMessage.error('商品不存在')
    router.push('/products')
  }
}

// 表单提交成功
const handleSuccess = async (formData: Product) => {
  if (isEdit.value) {
    await productStore.updateProduct(formData.id!, formData)
    ElMessage.success('更新成功')
  } else {
    await productStore.createProduct(formData)
    ElMessage.success('创建成功')
  }
  router.push('/products')
}

onMounted(() => {
  loadProduct()
})
</script>

<style scoped lang="scss">
.product-form-page {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
}
</style> 
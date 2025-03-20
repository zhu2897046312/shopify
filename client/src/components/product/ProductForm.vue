<template>
  <el-form
    ref="formRef"
    :model="form"
    :rules="rules"
    label-width="100px"
  >
    <el-form-item label="商品名称" prop="name">
      <el-input v-model="form.name" placeholder="请输入商品名称"/>
    </el-form-item>

    <el-form-item label="商品分类" prop="category">
    <el-select v-model="form.category"
      placeholder="请选择商品分类"
      aria-label="商品分类"
    >
      <el-option
        v-for="category in Object.entries(ProductCategoryNames)"
        :key="category[0]"
        :label="category[1]"
        :value="category[0]"
      />
    </el-select>
  </el-form-item>

  <el-form-item label="商品价格" prop="price">
  <el-input-number 
    v-model="form.price" 
    :precision="2" 
    :step="0.1" 
    :min="0"
    aria-label="商品价格"
    placeholder="请输入商品价格"
    />
  </el-form-item>

  <el-form-item label="商品库存" prop="stock">
    <el-input-number 
      v-model="form.stock" 
      :min="0" 
      :step="1"
      aria-label="商品库存"
      placeholder="请输入商品库存"
    />
  </el-form-item>

    <!-- <el-form-item label="商品图片" prop="images">
      <el-upload
        class="image-uploader"
        action="/api/v1/upload/product"
        list-type="picture-card"
        :on-success="handleImageSuccess"
        :on-remove="handleImageRemove"
        multiple
      >
        <el-icon><Plus /></el-icon>
      </el-upload>
    </el-form-item> -->
    <el-form-item label="商品图片" prop="images">
      <el-input
        v-model="imageInput"
        placeholder="请输入图片 URL，多个URL请用逗号分隔"
        suffix-icon="el-icon-link"
        @input="handleImageInput"
      ></el-input>
    </el-form-item>

    <el-form-item label="商品描述" prop="description">
      <el-input
        v-model="form.description"
        type="textarea"
        :rows="4"
        placeholder="请输入商品描述"
        aria-label="商品描述"
      />
    </el-form-item>

    <el-form-item label="商品状态" prop="status">
      <el-radio-group v-model="form.status"  aria-label="商品状态">
        <el-radio :label="ProductStatus.Active" data-test="status-active">上架</el-radio>
        <el-radio :label="ProductStatus.Inactive" data-test="status-inactive">下架</el-radio>
      </el-radio-group>
    </el-form-item>

    <el-form-item label="商品评分" prop="rating">
      <el-input-number 
        v-model="form.rating" 
        :min="0" 
        :step="1"
        aria-label="商品评分"
        placeholder="请输入商品评分"
      />
    </el-form-item>

    <el-form-item label="商品销量" prop="sales">
      <el-input-number 
        v-model="form.sales" 
        :min="0" 
        :step="1"
        aria-label="商品销量"
        placeholder="请输入商品销量"
      />
    </el-form-item>

    <el-form-item>
      <el-button type="primary" :loading="loading" @click="handleSubmit">
        保存
      </el-button>
      <el-button @click="$router.back()">取消</el-button>
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
import { useProductStore } from '@/stores/product'
import type { Product, } from '@/types/product'
import { ProductStatus } from '@/types/product'
import type { FormInstance } from 'element-plus'
import { ProductCategory, ProductCategoryNames,UpdateProductRequest} from '@/types/product'
import { ElMessage } from 'element-plus'

const props = defineProps<{
  product?: Product
}>()

const emit = defineEmits<{
  (e: 'success', product: Product): void
}>()

const productStore = useProductStore()
const formRef = ref<FormInstance>()
const loading = ref(false)
const categories = ref<ProductCategory[]>([])
const imageInput = ref('')

const form = reactive<UpdateProductRequest>({
  name: props.product?.name || '',
  description: props.product?.description || '',
  category: props.product?.category || ProductCategory.Clothing,
  price: props.product?.price || 0,
  stock: props.product?.stock || 0,
  images: props.product?.images || [],
  status: props.product?.status || ProductStatus.Active,
  tags: props.product?.tags || [],
  rating: props.product?.rating || 0,
  sales: props.product?.sales || 0
})

const rules = {
  name: [
    { required: true, message: '请输入商品名称', trigger: 'blur' }
  ],
  category: [
    { required: true, message: '请选择商品分类', trigger: 'change' }
  ]
}

// 监听产品数据变化
watch(() => props.product, (newProduct) => {
  if (newProduct) {
    // 更新表单数据
    Object.assign(form, {
      name: newProduct.name,
      description: newProduct.description,
      category: newProduct.category,
      price: newProduct.price,
      stock: newProduct.stock,
      status: newProduct.status,
      rating: newProduct.rating,
      sales: newProduct.sales,
      images: newProduct.images,
      tags: newProduct.tags
    })
    // 初始化图片输入
    imageInput.value = newProduct.images?.join(', ') || ''
  }
}, { immediate: true })

const handleImageInput = () => {
  const images = imageInput.value.split(',').map((url) => url.trim())
  form.images = images
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    loading.value = true
    
    const formData = {
      ...form,
      category: form.category?.toString() || ProductCategory.Clothing.toString(),
      images: form.images?.filter(url => url !== '') || [],
      tags: form.tags || []
    }

    if (props.product) {
      // Update existing product
      const updatedProduct = await productStore.updateProduct(props.product.id, formData)
      ElMessage.success('商品更新成功')
      emit('success', updatedProduct)
    } else {
      // Create new product
      const newProduct = await productStore.createProduct(formData as any)
      ElMessage.success('商品创建成功')
      emit('success', newProduct)
    }
    
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  // loadCategories()
})
</script>

<style scoped lang="scss">
.image-uploader {
  :deep(.el-upload--picture-card) {
    width: 100px;
    height: 100px;
    line-height: 100px;
  }
}
</style>
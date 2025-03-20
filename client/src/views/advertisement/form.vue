<template>
  <div class="advertisement-form-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>{{ isEdit ? '编辑广告' : '添加广告' }}</span>
          <el-button @click="$router.back()">返回</el-button>
        </div>
      </template>

      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" />
        </el-form-item>

        <!-- <el-form-item label="图片" prop="image">
          <el-upload
            class="image-uploader"
            action="/api/v1/upload"
            :show-file-list="false"
            :on-success="handleImageSuccess"
          >
            <el-image
              v-if="form.image"
              :src="form.image"
              fit="cover"
              class="uploaded-image"
            />
            <el-icon v-else class="image-uploader-icon"><Plus /></el-icon>
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


        <el-form-item label="链接" prop="url">
          <el-input v-model="form.url" />
        </el-form-item>

        <el-form-item label="位置" prop="position">
          <el-select v-model="form.position">
            <el-option label="首页" :value="AdvertisementPosition.BannerTop" />
            <el-option label="分类页" :value="AdvertisementPosition.BannerBottom" />
            <el-option label="详情页" :value="AdvertisementPosition.SidebarTop" />
            <el-option label="横幅广告" :value="AdvertisementPosition.SidebarBottom" />
          </el-select>
        </el-form-item>

        <el-form-item label="有效期" prop="date_range">
          <el-date-picker
            v-model="form.date_range"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
          />
        </el-form-item>

        <el-form-item label="状态" prop="status">
          <el-radio-group v-model="form.status">
            <el-radio label="active">启用</el-radio>
            <el-radio label="inactive">禁用</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleSubmit">
            保存
          </el-button>
          <el-button @click="$router.back()">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Plus } from '@element-plus/icons-vue'
import type { FormInstance } from 'element-plus'
import { useAdvertisementStore } from '@/stores/advertisement'
import { AdvertisementPosition, AdvertisementStatus } from '@/types/advertisement'

const route = useRoute()
const router = useRouter()
const store = useAdvertisementStore()
const loading = ref(false)
const formRef = ref<FormInstance>()
const imageInput = ref('')

const isEdit = computed(() => route.name === 'EditAdvertisement')

const form = reactive({
  title: '',
  image: '',
  url: '',
  position: AdvertisementPosition.BannerTop,
  date_range: [] as string[],
  status: AdvertisementStatus.Active
})

const rules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  image: [{ required: true, message: '请上传图片', trigger: 'change' }],
  url: [{ required: true, message: '请输入链接', trigger: 'blur' }],
  position: [{ required: true, message: '请选择位置', trigger: 'change' }],
  date_range: [{ required: true, message: '请选择有效期', trigger: 'change' }]
}

const handleImageSuccess = (res: any) => {
  form.image = res.data.url
}

const handleImageInput = () => {
  const images = imageInput.value.split(',').map((url) => url.trim())
  form.image = images[0]
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const data = {
          ...form,
          start_time: form.date_range[0],
          end_time: form.date_range[1]
        }
        
        if (isEdit.value) {
          await store.updateAdvertisement(Number(route.params.id), data)
        } else {
          await store.createAdvertisement(data)
        }
        
        router.push('/advertisements')
      } finally {
        loading.value = false
      }
    }
  })
}

// 如果是编辑模式，加载广告数据
onMounted(async () => {
  if (isEdit.value) {
    loading.value = true
    try {
      const ad = await store.getAdvertisement(Number(route.params.id))
      if (ad) {
        // 确保所有字段都被正确赋值
        form.title = ad.title
        form.image = ad.image
        form.url = ad.url
        form.position = ad.position
        form.status = ad.status
        form.date_range = [ad.start_time, ad.end_time]
        
        console.log('Loaded advertisement data:', {
          title: form.title,
          image: form.image,
          url: form.url,
          position: form.position,
          status: form.status,
          date_range: form.date_range
        })
      }
    } catch (error) {
      console.error('Failed to load advertisement:', error)
    } finally {
      loading.value = false
    }
  }
})
</script>

<style scoped lang="scss">
.advertisement-form-container {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .image-uploader {
    :deep(.el-upload) {
      border: 1px dashed #d9d9d9;
      border-radius: 6px;
      cursor: pointer;
      position: relative;
      overflow: hidden;
      transition: var(--el-transition-duration-fast);
      
      &:hover {
        border-color: var(--el-color-primary);
      }
    }
    
    .image-uploader-icon {
      font-size: 28px;
      color: #8c939d;
      width: 178px;
      height: 178px;
      text-align: center;
      line-height: 178px;
    }
    
    .uploaded-image {
      width: 178px;
      height: 178px;
      display: block;
    }
  }
}
</style>
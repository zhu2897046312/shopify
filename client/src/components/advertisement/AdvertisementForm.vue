<template>
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
        action="/api/v1/admin/upload"
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
        placeholder="请输入图片 URL"
        suffix-icon="el-icon-link"
        @input="handleImageInput"
      ></el-input>
    </el-form-item>

    <el-form-item label="链接" prop="url">
      <el-input v-model="form.url" />
    </el-form-item>

    <el-form-item label="位置" prop="position">
      <el-select v-model="form.position">
        <el-option label="顶部横幅" :value="AdvertisementPosition.BannerTop" />
        <el-option label="底部横幅" :value="AdvertisementPosition.BannerBottom" />
        <el-option label="首页轮播图" :value="AdvertisementPosition.HomeBanner" />
        <el-option label="首页中部广告位" :value="AdvertisementPosition.HomeMiddle" />
        <el-option label="侧边栏顶部" :value="AdvertisementPosition.SidebarTop" />
        <el-option label="侧边栏底部" :value="AdvertisementPosition.SidebarBottom" />
      </el-select>
    </el-form-item>

    <el-form-item label="有效期" required>
      <el-date-picker
        v-model="form.time_range"
        type="datetimerange"
        range-separator="至"
        start-placeholder="开始时间"
        end-placeholder="结束时间"
        value-format="YYYY-MM-DDTHH:mm:ssZ"
        :default-time="[
          new Date(2000, 1, 1, 0, 0, 0),
          new Date(2000, 1, 1, 23, 59, 59),
        ]"
      />
    </el-form-item>

    <el-form-item label="状态" prop="status">
      <el-radio-group v-model="form.status">
        <el-radio :value="AdvertisementStatus.Active">启用</el-radio>
        <el-radio :value="AdvertisementStatus.Inactive">禁用</el-radio>
      </el-radio-group>
    </el-form-item>

    <el-form-item>
      <el-button type="primary" :loading="loading" @click="handleSubmit">
        {{ advertisement ? '更新' : '创建' }}
      </el-button>
      <el-button @click="$emit('cancel')">取消</el-button>
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import type { FormInstance } from 'element-plus'
import type { Advertisement } from '@/types/advertisement'
import { AdvertisementPosition, AdvertisementStatus } from '@/types/advertisement'

const props = defineProps<{
  advertisement?: Advertisement
}>()

const emit = defineEmits<{
  (e: 'success', data: Advertisement): void
  (e: 'cancel'): void
}>()

const formRef = ref<FormInstance>()
const loading = ref(false)
const imageInput = ref('')

// Initialize form data when advertisement prop changes
watch(() => props.advertisement, (newAd) => {
  if (newAd) {
    imageInput.value = newAd.image || ''
  }
}, { immediate: true })

const form = reactive({
  title: props.advertisement?.title || '',
  image: props.advertisement?.image || 'http://dummyimage.com/400x400',
  url: props.advertisement?.url || '',
  position: props.advertisement?.position || AdvertisementPosition.BannerTop,
  time_range: props.advertisement 
    ? [props.advertisement.start_time, props.advertisement.end_time]
    : [],
  status: props.advertisement?.status || AdvertisementStatus.Active
})

// Watch for changes in advertisement prop
watch(
  () => props.advertisement,
  (newVal) => {
    if (newVal) {
      form.title = newVal.title
      form.image = newVal.image
      form.url = newVal.url
      form.position = newVal.position
      form.time_range = [newVal.start_time, newVal.end_time]
      form.status = newVal.status
    } else {
      form.title = ''
      form.image = 'http://dummyimage.com/400x400'
      form.url = ''
      form.position = AdvertisementPosition.BannerTop
      form.time_range = []
      form.status = AdvertisementStatus.Active
    }
  },
  { deep: true }
)

const rules = {
  title: [
    { required: true, message: '请输入广告标题', trigger: 'blur' }
  ],
  image: [
    { required: true, message: '请上传广告图片', trigger: 'change' }
  ],
  url: [
    { required: true, message: '请输入广告链接', trigger: 'blur' }
  ],
  position: [
    { required: true, message: '请选择广告位置', trigger: 'change' }
  ],
  time_range: [
    { required: true, message: '请选择有效期', trigger: 'change' }
  ]
}

const handleImageSuccess = (response: any) => {
  form.image = response.url
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const { time_range, ...formData } = form
        const submitData = {
          title: formData.title,
          image: formData.image,
          url: formData.url,
          position: formData.position,
          start_time: time_range[0],
          end_time: time_range[1],
          status: formData.status
        }
        emit('success', submitData as Advertisement)
      } finally {
        loading.value = false
      }
    }
  })
}

const handleImageInput = () => {
  const images = imageInput.value.split(',').map((url) => url.trim())
  form.image = images[0]
}

</script>

<style scoped lang="scss">
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
</style>
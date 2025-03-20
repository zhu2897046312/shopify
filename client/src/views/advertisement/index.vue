<template>
  <div class="advertisement-page">
    <!-- 操作栏 -->
    <el-card class="action-card">
      <template #header>
        <div class="card-header">
          <span>广告列表</span>
          <el-button type="primary" @click="handleCreate">
            <el-icon><Plus /></el-icon>新增广告
          </el-button>
        </div>
      </template>

      <advertisement-table
        ref="tableRef"
        @edit="handleEdit"
      />
    </el-card>

    <!-- 表单对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="currentAd ? '编辑广告' : '新增广告'"
      width="600px"
    >
      <advertisement-form
        :advertisement="currentAd"
        @success="handleFormSuccess"
        @cancel="dialogVisible = false"
      />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { useAdvertisementStore } from '@/stores/advertisement'
import type { Advertisement } from '@/types/advertisement'
import  { AdvertisementStatus } from '@/types/advertisement'
import AdvertisementTable from '@/components/advertisement/AdvertisementTable.vue'
import AdvertisementForm from '@/components/advertisement/AdvertisementForm.vue'

const store = useAdvertisementStore()
const tableRef = ref()
const dialogVisible = ref(false)
const currentAd = ref<Advertisement | undefined>(undefined)

// 新增广告
const handleCreate = () => {
  currentAd.value = undefined
  dialogVisible.value = true
}

// 编辑广告
const handleEdit = (ad: Advertisement) => {
  currentAd.value = ad
  dialogVisible.value = true
}

// 表单提交成功
const handleFormSuccess = async (ad: Advertisement) => {
  dialogVisible.value = false
  try {
    const submitData = {
      title: ad.title,
      image: ad.image,
      url: ad.url,
      position: ad.position,
      start_time: ad.start_time,
      end_time: ad.end_time,
      status: currentAd.value ? ad.status : AdvertisementStatus.Active
    }

    if (currentAd.value) {
      await store.updateAdvertisement(currentAd.value.id, submitData)
    } else {
      await store.createAdvertisement(submitData)
    }
    ElMessage.success(currentAd.value ? '更新成功' : '创建成功')
    tableRef.value?.loadData()
  } catch (error) {
    const message = error instanceof Error ? error.message : '未知错误'
    ElMessage.error('操作失败：' + message)
  }
}

// 组件挂载时加载数据
onMounted(() => {
  tableRef.value?.loadData()
})
</script>

<style scoped lang="scss">
.advertisement-page {
  .action-card {
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
  }
}
</style> 
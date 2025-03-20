<template>
  <div class="advertisement-table">
    <el-table
      v-loading="loading"
      :data="advertisements"
      border
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="title" label="标题" min-width="150" />
      <el-table-column label="图片" width="120">
        <template #default="{ row }">
          <el-image 
            :src="row.image" 
            :preview-src-list="[row.image]"
            fit="cover"
            class="ad-image"
          />
        </template>
      </el-table-column>
      <el-table-column prop="position" label="位置" width="120">
        <template #default="{ row }">
          <el-tag>{{ row.position }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 'active' ? 'success' : 'info'">
            {{ row.status === 'active' ? '启用' : '禁用' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="有效期" width="340">
        <template #default="{ row }">
          {{ row.start_time }} 至 {{ row.end_time }}
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
              {{ row.status === 'active' ? '禁用' : '启用' }}
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
import { ref, onMounted } from 'vue'
import { ElMessageBox } from 'element-plus'
import { useAdvertisementStore } from '@/stores/advertisement'
import { storeToRefs } from 'pinia'
import type { Advertisement} from '@/types/advertisement'
import  { AdvertisementStatus } from '@/types/advertisement'

const store = useAdvertisementStore()
const { advertisements, total, loading } = storeToRefs(store)
const currentPage = ref(1)
const pageSize = ref(10)
const selectedAds = ref<Advertisement[]>([])

// 加载数据
const loadData = () => {
  store.getAdvertisements({
    page: currentPage.value,
    page_size: pageSize.value
  })
}

// 选择变化
const handleSelectionChange = (selection: Advertisement[]) => {
  selectedAds.value = selection
}

// 编辑广告
const handleEdit = (row: Advertisement) => {
  emit('edit', row)
}

// 更新状态
const handleStatusChange = async (row: Advertisement) => {
  try {
    await ElMessageBox.confirm(
      `确定要${row.status === 'active' ? '禁用' : '启用'}该广告吗？`,
      '提示',
      { type: 'warning' }
    )
    await store.updateAdvertisementStatus(row.id, {
      status: row.status === AdvertisementStatus.Active 
        ? AdvertisementStatus.Inactive 
        : AdvertisementStatus.Active
    })
    loadData()
  } catch {
    // 用户取消操作
  }
}

// 删除广告
const handleDelete = async (row: Advertisement) => {
  try {
    await ElMessageBox.confirm('确定要删除该广告吗？', '提示', {
      type: 'warning'
    })
    await store.deleteAdvertisement(row.id)
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

// 暴露方法给父组件
defineExpose({
  loadData
})

const emit = defineEmits<{
  (e: 'edit', advertisement: Advertisement): void
}>()

// 组件挂载时加载数据
onMounted(() => {
  loadData()
})
</script>

<style scoped lang="scss">
.advertisement-table {
  .ad-image {
    width: 80px;
    height: 80px;
    border-radius: 4px;
  }

  .pagination {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}
</style> 
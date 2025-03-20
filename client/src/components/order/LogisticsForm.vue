<template>
  <el-form
    ref="formRef"
    :model="form"
    :rules="rules"
    label-width="100px"
  >
    <el-form-item label="物流公司" prop="carrier">
      <el-select v-model="form.carrier" placeholder="请选择物流公司">
        <el-option label="顺丰快递" value="SF" />
        <el-option label="中通快递" value="ZTO" />
        <el-option label="圆通快递" value="YTO" />
        <el-option label="韵达快递" value="YD" />
      </el-select>
    </el-form-item>

    <el-form-item label="物流单号" prop="tracking_number">
      <el-input v-model="form.tracking_number" />
    </el-form-item>

    <el-form-item label="物流备注">
      <el-input
        v-model="form.remark"
        type="textarea"
        :rows="3"
        placeholder="请输入物流备注信息"
      />
    </el-form-item>

    <el-form-item>
      <el-button type="primary" :loading="loading" @click="handleSubmit">
        保存
      </el-button>
      <el-button @click="$emit('cancel')">取消</el-button>
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import type { FormInstance } from 'element-plus'
import type { Order } from '@/types/order'
import { useOrderStore } from '@/stores/order'

const props = defineProps<{
  order: Order
}>()

const emit = defineEmits<{
  (e: 'success'): void
  (e: 'cancel'): void
}>()

const orderStore = useOrderStore()
const formRef = ref<FormInstance>()
const loading = ref(false)

const form = reactive({
  carrier: '',
  tracking_number: '',
  remark: ''
})

const rules = {
  carrier: [
    { required: true, message: '请选择物流公司', trigger: 'change' }
  ],
  tracking_number: [
    { required: true, message: '请输入物流单号', trigger: 'blur' }
  ]
}

const handleSubmit = async () => {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        await orderStore.updateLogistics(props.order.id, form)
        emit('success')
      } finally {
        loading.value = false
      }
    }
  })
}
</script> 
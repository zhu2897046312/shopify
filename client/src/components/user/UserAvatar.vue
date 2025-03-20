<template>
  <el-avatar :size="size" :src="src">
    {{ nickname?.charAt(0)?.toUpperCase() }}
  </el-avatar>
</template>

<script setup lang="ts">
defineProps<{
  src?: string
  nickname?: string
  size?: number
}>()
</script>

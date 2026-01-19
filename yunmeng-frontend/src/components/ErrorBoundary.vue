<template>
  <div v-if="hasError" class="error-boundary">
    <el-result
      icon="error"
      title="页面加载失败"
      sub-title="抱歉，页面出现了问题，请刷新重试"
    >
      <template #extra>
        <el-button type="primary" @click="handleRetry">重试</el-button>
        <el-button @click="handleBack">返回</el-button>
      </template>
    </el-result>
  </div>
  <div v-else>
    <slot />
  </div>
</template>

<script setup>
import { ref, onErrorCaptured } from 'vue'
import { useRouter } from 'vue-router'
import { handleVueError } from '@/utils/errorHandler'

const router = useRouter()
const hasError = ref(false)

onErrorCaptured((error, instance, info) => {
  // 使用统一的错误处理函数
  const shouldPreventPropagation = handleVueError(error, instance, info)
  
  // 如果是数据迭代错误，尝试恢复而不是显示错误页面
  if (error.message && error.message.includes('is not iterable')) {
    console.warn('检测到数据迭代错误，尝试自动恢复...')
    return true // 阻止错误继续传播，尝试恢复
  }
  
  hasError.value = true
  return shouldPreventPropagation
})

const handleRetry = () => {
  hasError.value = false
  window.location.reload()
}

const handleBack = () => {
  router.push('/login')
}
</script>

<style scoped>
.error-boundary {
  padding: 40px;
  text-align: center;
}
</style> 
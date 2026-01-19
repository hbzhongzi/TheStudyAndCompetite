<template>
  <div class="teacher-selector">
    <el-select
      v-model="selectedTeacherId"
      placeholder="请选择指导老师"
      style="width: 100%"
      :loading="loading"
      :disabled="disabled"
      @change="handleTeacherChange"
      filterable
      clearable
    >
      <el-option
        v-for="teacher in teachers"
        :key="teacher.id"
        :label="teacher.name"
        :value="teacher.id"
      >
        <div class="teacher-option">
          <div class="teacher-info">
            <span class="teacher-name">{{ teacher.name }}</span>
            <span class="teacher-title">{{ teacher.title || '教师' }}</span>
          </div>
          <div class="teacher-details">
            <span class="department">{{ teacher.department }}</span>
            <span class="specialty" v-if="teacher.specialty">{{ teacher.specialty }}</span>
          </div>
        </div>
      </el-option>
    </el-select>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { teacherService } from '../services/teacherService'

// Props
const props = defineProps({
  modelValue: {
    type: [String, Number],
    default: ''
  },
  disabled: {
    type: Boolean,
    default: false
  },
  placeholder: {
    type: String,
    default: '请选择指导老师'
  }
})

// Emits
const emit = defineEmits(['update:modelValue', 'change'])

// 响应式数据
const loading = ref(false)
const teachers = ref([])
const selectedTeacherId = ref(props.modelValue)

// 监听modelValue变化
watch(() => props.modelValue, (newValue) => {
  selectedTeacherId.value = newValue
})

// 监听selectedTeacherId变化
watch(selectedTeacherId, (newValue) => {
  emit('update:modelValue', newValue)
})

// 获取教师列表
const loadTeachers = async () => {
  try {
    loading.value = true
    const response = await teacherService.getTeacherList()
    if (response.code === 200) {
      teachers.value = response.data || []
    } else {
      ElMessage.error(response.message || '获取教师列表失败')
    }
  } catch (error) {
    console.error('获取教师列表失败:', error)
    ElMessage.error('获取教师列表失败')
  } finally {
    loading.value = false
  }
}

// 处理教师选择变化
const handleTeacherChange = (teacherId) => {
  const selectedTeacher = teachers.value.find(t => t.id === teacherId)
  emit('change', {
    teacherId,
    teacher: selectedTeacher
  })
}

// 组件挂载时加载教师列表
onMounted(() => {
  loadTeachers()
})
</script>

<style scoped>
.teacher-selector {
  width: 100%;
}

.teacher-option {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.teacher-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.teacher-name {
  font-weight: 500;
  color: #303133;
}

.teacher-title {
  font-size: 12px;
  color: #909399;
  background-color: #f5f7fa;
  padding: 2px 6px;
  border-radius: 4px;
}

.teacher-details {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.department {
  font-size: 12px;
  color: #606266;
}

.specialty {
  font-size: 11px;
  color: #909399;
}
</style> 
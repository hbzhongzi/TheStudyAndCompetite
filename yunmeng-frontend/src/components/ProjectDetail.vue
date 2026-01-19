<template>
  <el-dialog
    v-model="dialogVisible"
    title="项目详情"
    width="900px"
    :before-close="handleClose"
  >
    <div v-loading="loading" class="project-detail">
      <!-- 项目基本信息 -->
      <el-card class="info-card">
        <template #header>
          <div class="card-header">
            <span>基本信息</span>
            <el-tag :type="getStatusType(project.status)">
              {{ getStatusText(project.status) }}
            </el-tag>
          </div>
        </template>
        
        <el-descriptions :column="2" border>
          <el-descriptions-item label="项目ID">{{ project.id }}</el-descriptions-item>
          <el-descriptions-item label="项目标题">{{ project.title }}</el-descriptions-item>
          <el-descriptions-item label="项目类型">
            <el-tag :type="project.type === '科研' ? 'primary' : 'success'">
              {{ project.type }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ formatDate(project.createdAt) }}</el-descriptions-item>
          <el-descriptions-item label="更新时间">{{ formatDate(project.updatedAt) }}</el-descriptions-item>
          <el-descriptions-item label="项目状态">
            <el-tag :type="getStatusType(project.status)">
              {{ getStatusText(project.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="项目描述" :span="2">
            <div class="description-content">{{ project.description || '暂无描述' }}</div>
          </el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 项目负责人信息 -->
      <el-card class="info-card" v-if="project.student">
        <template #header>
          <span>项目负责人</span>
        </template>
        
        <el-descriptions :column="2" border>
          <el-descriptions-item label="姓名">{{ project.student.realName || project.student.username }}</el-descriptions-item>
          <el-descriptions-item label="学号">{{ project.student.studentId || '暂无' }}</el-descriptions-item>
          <el-descriptions-item label="邮箱">{{ project.student.email || '暂无' }}</el-descriptions-item>
          <el-descriptions-item label="电话">{{ project.student.phone || '暂无' }}</el-descriptions-item>
          <el-descriptions-item label="学院">{{ project.student.department || '暂无' }}</el-descriptions-item>
        </el-descriptions>
      </el-card>

      <!-- 项目成员 -->
      <el-card class="info-card" v-if="project.members && project.members.length > 0">
        <template #header>
          <span>项目成员 ({{ project.members.length }}人)</span>
        </template>
        
        <el-table :data="project.members" style="width: 100%">
          <el-table-column prop="name" label="姓名" width="120" />
          <el-table-column prop="studentNumber" label="学号" width="120" />
          <el-table-column prop="role" label="角色" width="120" />
        </el-table>
      </el-card>

      <!-- 项目附件 -->
      <el-card class="info-card" v-if="project.files && project.files.length > 0">
        <template #header>
          <span>项目附件 ({{ project.files.length }}个)</span>
        </template>
        
        <div class="attachments-list">
          <div 
            v-for="file in project.files" 
            :key="file.id" 
            class="attachment-item"
          >
            <div class="attachment-info">
              <el-icon class="file-icon"><Document /></el-icon>
              <span class="file-name">{{ file.fileName }}</span>
              <span class="file-time">{{ formatDate(file.uploadTime) }}</span>
            </div>
            <div class="attachment-actions">
              <el-button 
                type="primary" 
                size="small" 
                @click="downloadFile(file)"
              >
                下载
              </el-button>
            </div>
          </div>
        </div>
      </el-card>

      <!-- 审核记录 -->
      <el-card class="info-card" v-if="project.reviews && project.reviews.length > 0">
        <template #header>
          <span>审核记录 ({{ project.reviews.length }}条)</span>
        </template>
        
        <el-timeline>
          <el-timeline-item
            v-for="review in project.reviews"
            :key="review.id"
            :timestamp="formatDate(review.reviewTime)"
            :type="review.status === 'approved' ? 'success' : 'danger'"
          >
            <el-card class="review-card">
              <div class="review-header">
                <span class="reviewer">{{ review.reviewer.realName || review.reviewer.username }}</span>
                <el-tag :type="review.status === 'approved' ? 'success' : 'danger'">
                  {{ review.status === 'approved' ? '通过' : '拒绝' }}
                </el-tag>
              </div>
              <div class="review-content">
                {{ review.comments || '暂无审核意见' }}
              </div>
            </el-card>
          </el-timeline-item>
        </el-timeline>
      </el-card>

      <!-- 空状态 -->
      <el-empty 
        v-if="!loading && !project.id" 
        description="项目信息加载失败"
        style="margin-top: 40px;"
      />
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">关闭</el-button>
        <el-button 
          v-if="project.status === 'draft'" 
          type="primary" 
          @click="handleEdit"
        >
          编辑项目
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Document } from '@element-plus/icons-vue'
import { projectService } from '../services/projectService'

// Props
const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  projectId: {
    type: [String, Number],
    default: null
  }
})

// Emits
const emit = defineEmits(['update:visible', 'edit'])

// 响应式数据
const loading = ref(false)
const project = ref({})

// 计算属性
const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

// 监听项目ID变化
watch(() => props.projectId, (newId) => {
  if (newId && dialogVisible.value) {
    loadProjectDetail()
  }
}, { immediate: true })

// 监听对话框显示状态
watch(() => dialogVisible.value, (visible) => {
  if (visible && props.projectId) {
    loadProjectDetail()
  }
})

// 加载项目详情
const loadProjectDetail = async () => {
  if (!props.projectId) return
  
  loading.value = true
  try {
    const response = await projectService.getProjectDetail(props.projectId)
    project.value = response.data || {}
  } catch (error) {
    ElMessage.error(error.message || '加载项目详情失败')
    project.value = {}
  } finally {
    loading.value = false
  }
}

// 状态类型映射
const getStatusType = (status) => {
  const statusMap = {
    draft: 'info',
    pending: 'warning',
    approved: 'success',
    rejected: 'danger'
  }
  return statusMap[status] || 'info'
}

// 状态文本映射
const getStatusText = (status) => {
  const statusMap = {
    draft: '草稿',
    pending: '待审核',
    approved: '已通过',
    rejected: '已拒绝'
  }
  return statusMap[status] || status
}

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return '暂无'
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

// 下载文件
const downloadFile = (file) => {
  const link = document.createElement('a')
  link.href = `http://localhost:8080${file.fileUrl}`
  link.download = file.fileName
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

// 编辑项目
const handleEdit = () => {
  emit('edit', project.value)
  handleClose()
}

// 关闭对话框
const handleClose = () => {
  dialogVisible.value = false
  project.value = {}
}
</script>

<style scoped>
.project-detail {
  max-height: 70vh;
  overflow-y: auto;
}

.info-card {
  margin-bottom: 20px;
}

.info-card:last-child {
  margin-bottom: 0;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.description-content {
  white-space: pre-wrap;
  line-height: 1.6;
  color: #606266;
}

.attachments-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.attachment-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px;
  background: #f8f9fa;
  border-radius: 4px;
  border: 1px solid #e4e7ed;
}

.attachment-info {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
}

.file-icon {
  color: #409eff;
  font-size: 16px;
}

.file-name {
  color: #303133;
  font-size: 14px;
  font-weight: 500;
}

.file-time {
  color: #909399;
  font-size: 12px;
}

.attachment-actions {
  display: flex;
  gap: 8px;
}

.review-card {
  border: none;
  box-shadow: none;
  background: transparent;
}

.review-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.reviewer {
  font-weight: 500;
  color: #303133;
}

.review-content {
  color: #606266;
  line-height: 1.6;
  white-space: pre-wrap;
}

.dialog-footer {
  text-align: right;
}

:deep(.el-descriptions__label) {
  font-weight: 500;
  color: #303133;
}

:deep(.el-timeline-item__node) {
  background-color: #409eff;
}

:deep(.el-timeline-item__node--success) {
  background-color: #67c23a;
}

:deep(.el-timeline-item__node--danger) {
  background-color: #f56c6c;
}
</style> 
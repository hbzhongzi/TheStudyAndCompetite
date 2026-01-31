<template>
  <div class="project-detail">

    <!-- 项目基本信息 -->
    <el-card class="mb-20">
      <template #header>
        <span>项目详情</span>
      </template>

      <el-descriptions border :column="2">
        <el-descriptions-item label="项目名称">
          {{ project.title }}
        </el-descriptions-item>

        <el-descriptions-item label="项目类型">
          {{ project.type }}
        </el-descriptions-item>

        <el-descriptions-item label="指导教师">
          {{ project.teacherName || '—' }}
        </el-descriptions-item>

        <el-descriptions-item label="项目状态">
          <el-tag :type="statusType(project.status)">
            {{ statusText(project.status) }}
          </el-tag>
        </el-descriptions-item>

        <el-descriptions-item label="创建时间">
          {{ formatDateTime(project.createdAt) }}
        </el-descriptions-item>

        <el-descriptions-item label="预计完成时间">
          {{ project.finishedAt || '—' }}
        </el-descriptions-item>
      </el-descriptions>

      <el-alert
        v-if="!canEdit"
        class="mt-10"
        type="warning"
        show-icon
        title="项目尚未通过审核，暂不可上传或删除文件"
      />
    </el-card>

    <!-- 文件管理 -->
    <el-card>
      <template #header>
        <div class="header-content">
          <span>项目文件</span>
          <el-button
            v-if="canEdit"
            type="primary"
            @click="openUpload"
          >
            上传文件
          </el-button>
        </div>
      </template>

      <el-table
        :data="files"
        v-loading="loading"
        border
      >
        <el-table-column prop="name" label="文件名" min-width="200" />

        <el-table-column prop="category" label="类型" width="120">
          <template #default="{ row }">
            <el-tag size="small">
              {{ categoryLabel(row.category) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="大小" width="120">
          <template #default="{ row }">
            {{ formatFileSize(row.size) }}
          </template>
        </el-table-column>

        <el-table-column prop="createdAt" label="上传时间" width="160" />

        <el-table-column label="操作" width="240">
          <template #default="{ row }">
            <el-button size="small" @click="previewFile(row)">
              预览
            </el-button>
            <el-button size="small" type="primary" @click="downloadFile(row)">
              下载
            </el-button>
            <el-button
              v-if="canEdit"
              size="small"
              type="danger"
              @click="deleteFile(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-empty
        v-if="!loading && files.length === 0"
        description="暂无项目文件"
      />
    </el-card>

    <!-- 上传文件弹窗 -->
    <el-dialog
      v-model="uploadDialogVisible"
      title="上传文件"
      width="500px"
    >
      <el-upload
        drag
        multiple
        :auto-upload="false"
        :on-change="handleFileChange"
      >
        <el-icon><Upload /></el-icon>
        <div class="el-upload__text">
          点击或拖拽文件到此处上传
        </div>
      </el-upload>

      <template #footer>
        <el-button @click="uploadDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitUpload">
          上传
        </el-button>
      </template>
    </el-dialog>

  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { projectService } from '@/services/projectService'

const route = useRoute()
const projectId = Number(route.params.id)

const project = ref({})
const files = ref([])
const loading = ref(false)

const uploadDialogVisible = ref(false)
const uploadFiles = ref([])

/* ========= 权限 ========= */
const canEdit = computed(() => project.value.status === 'approved')

/* ========= 状态显示 ========= */
const statusType = (status) => {
  switch (status) {
    case 'approved': return 'success'
    case 'pending': return 'warning'
    case 'rejected': return 'danger'
    default: return 'info'
  }
}

const statusText = (status) => {
  switch (status) {
    case 'approved': return '已通过'
    case 'pending': return '审核中'
    case 'rejected': return '已驳回'
    default: return '未知'
  }
}

/* ========= 工具方法 ========= */
const formatDateTime = (time) => {
  if (!time) return '—'
  return new Date(time).toLocaleString()
}

const formatFileSize = (size) => {
  if (!size) return '0 KB'
  if (size < 1024) return size + ' B'
  if (size < 1024 * 1024) return (size / 1024).toFixed(1) + ' KB'
  return (size / 1024 / 1024).toFixed(1) + ' MB'
}

const categoryLabel = (val) => {
  const map = {
    document: '文档',
    image: '图片',
    video: '视频',
    code: '代码',
    other: '其他'
  }
  return map[val] || '未知'
}

/* ========= 数据加载 ========= */
const loadProject = async () => {
  const res = await projectService.getProjectDetail(projectId)
  project.value = res.data || {}
}

const loadFiles = async () => {
  loading.value = true
  const res = await projectService.getProjectFiles(projectId)
  files.value = res.data || []
  loading.value = false
}

/* ========= 上传 ========= */
const openUpload = () => {
  uploadFiles.value = []
  uploadDialogVisible.value = true
}

const handleFileChange = (file) => {
  uploadFiles.value = [file.raw]
}

const submitUpload = async () => {
  if (uploadFiles.value.length === 0) {
    ElMessage.warning('请选择文件')
    return
  }
  await projectService.uploadProjectFiles(projectId, uploadFiles.value)
  ElMessage.success('上传成功')
  uploadDialogVisible.value = false
  loadFiles()
}

/* ========= 删除 ========= */
const deleteFile = (file) => {
  ElMessageBox.confirm('确认删除该文件？', '提示', { type: 'warning' })
    .then(async () => {
      await projectService.deleteProjectFile(projectId, [file.id])
      ElMessage.success('删除成功')
      loadFiles()
    })
}

/* ========= 占位操作 ========= */
const previewFile = (file) => {
  ElMessage.info('预览功能可后续扩展')
}

const downloadFile = (file) => {
  window.open(file.url, '_blank')
}

onMounted(() => {
  loadProject()
  loadFiles()
})
</script>

<style scoped>
.project-detail {
  padding: 20px;
}

.mb-20 {
  margin-bottom: 20px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.mt-10 {
  margin-top: 10px;
}
</style>
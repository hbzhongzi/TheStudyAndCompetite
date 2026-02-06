<template>
  <div class="project-detail">

    <!-- 项目选择 -->
    <el-card class="mb-20">
      <div class="project-select">
        <span class="label">选择项目：</span>
        <el-select
          v-model="projectId"
          placeholder="请选择项目"
          style="width: 300px"
          @change="handleProjectChange"
        >
          <el-option
            v-for="item in projects"
            :key="item.id"
            :label="item.title"
            :value="item.id"
          />
        </el-select>
      </div>
    </el-card>

<!-- 项目基本信息 -->
<el-card class="mb-20" v-if="project.id">
  <template #header>
    <div class="header-content">
      <span>项目详情</span>
      <el-tag
        :type="project.isApproved ? 'success' : 'warning'"
        size="small"
      >
        {{ project.isApproved ? '已立项' : '未立项' }}
      </el-tag>
    </div>
  </template>

  <el-descriptions border :column="2">
    <!-- 项目 -->
    <el-descriptions-item label="项目名称">
      {{ project.title }}
    </el-descriptions-item>

    <el-descriptions-item label="项目类型">
      {{ project.type }}
    </el-descriptions-item>

    <el-descriptions-item label="项目状态">
      <el-tag :type="statusType(project.status)">
        {{ statusText(project.status) }}
      </el-tag>
    </el-descriptions-item>

    <el-descriptions-item label="创建时间">
      {{ formatDateTime(project.createdAt) }}
    </el-descriptions-item>

    <!-- 学生 -->
    <el-descriptions-item label="学生姓名">
      {{ project.student?.realName || '—' }}
    </el-descriptions-item>

    <el-descriptions-item label="学号">
      {{ project.student?.studentId || '—' }}
    </el-descriptions-item>

    <el-descriptions-item label="学生学院">
      {{ project.student?.department || '—' }}
    </el-descriptions-item>

    <el-descriptions-item label="学生邮箱">
      {{ project.student?.email || '—' }}
    </el-descriptions-item>

    <!-- 教师（仅显示指导教师名称） -->
    <el-descriptions-item label="指导教师">
      {{ project.teacher?.realName || '—' }}
    </el-descriptions-item>
  </el-descriptions>
</el-card>

    <!-- 项目成果文件 -->
    <el-card>
      <template #header>
        <div class="header-content">
          <span>项目成果文件</span>
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
        <el-table-column
          prop="originalName"
          label="文件名"
          min-width="220"
        />

        <el-table-column label="类型" width="120">
          <template #default="{ row }">
            <el-tag size="small">
              {{ row.ext}}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="大小" width="120">
          <template #default="{ row }">
            {{ formatFileSize(row.size) }}
          </template>
        </el-table-column>

        <el-table-column label="上传时间" width="170">
          <template #default="{ row }">
            {{ formatDateTime(row.createdAt) }}
          </template>
        </el-table-column>

        <el-table-column label="操作" width="220">
          <template #default="{ row }">
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
        description="暂无项目成果文件"
      />
    </el-card>

    <!-- 上传文件 -->
    <el-dialog v-model="uploadDialogVisible" title="上传文件" width="500px">
      <el-upload
        drag
        multiple
        :auto-upload="false"
        :on-change="handleFileChange"
      >
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
import { ElMessage, ElMessageBox } from 'element-plus'
import { projectService } from '@/services/projectService'

/* ========= state ========= */
const projects = ref([])
const projectId = ref(null)
const project = ref({})
const files = ref([])
const loading = ref(false)

const uploadDialogVisible = ref(false)
const uploadFiles = ref([])

/* ========= 权限 ========= */
const canEdit = computed(() => project.value.status === 'approved')

/* ========= 状态 ========= */
const statusText = (status)  =>
  ({ draft: '草稿', submitted: '待审核', approved: '已通过', rejected: '已拒绝' }[status] || status)

const statusType = (status)  =>
  ({ draft: '', submitted: 'warning', approved: 'success', rejected: 'danger' }[status] || '')


/* ========= 工具 ========= */
const formatDateTime = (time) =>
  time ? new Date(time).toLocaleString() : '—'

const formatFileSize = (size) => {
  if (!size) return '0 KB'
  if (size < 1024) return size + ' B'
  if (size < 1024 * 1024) return (size / 1024).toFixed(1) + ' KB'
  return (size / 1024 / 1024).toFixed(1) + ' MB'
}


/* ========= 数据加载 ========= */

// 1️⃣ 获取我的项目
const loadMyProjects = async () => {
  const res = await projectService.getMyProjects()
  const list = res.data?.list || []

  if (list.length === 0) {
    ElMessage.warning('暂无项目')
    return
  }

  projects.value = list
  projectId.value = list[0].id // ✅ 默认第一个
}

// 2️⃣ 加载详情
const loadProjectDetail = async () => {
  if (!projectId.value) return
  const res = await projectService.getProjectDetail(projectId.value)
  project.value = res.data || {}
}

// 3️⃣ 加载文件
const loadFiles = async () => {
  if (!projectId.value) return
  loading.value = true
  const res = await projectService.getProjectFiles(projectId.value)
  files.value = res.data || []
  loading.value = false
}

/* ========= 项目切换 ========= */
const handleProjectChange = async () => {
  project.value = {}
  files.value = []
  await loadProjectDetail()
  await loadFiles()
}

/* ========= 上传 ========= */
const openUpload = () => {
  uploadFiles.value = []
  uploadDialogVisible.value = true
}

const handleFileChange = (file) => {
  uploadFiles.value.push(file.raw)
}

const submitUpload = async () => {
  if (uploadFiles.value.length === 0) {
    ElMessage.warning('请选择文件')
    return
  }
  await projectService.uploadProjectFiles(projectId.value, uploadFiles.value)
  ElMessage.success('上传成功')
  uploadDialogVisible.value = false
  loadFiles()
}

/* ========= 删除 ========= */
const deleteFile = (file) => {
  ElMessageBox.confirm('确认删除该文件？', '提示', { type: 'warning' })
    .then(async () => {
      await projectService.deleteProjectFile(projectId.value, file.id)
      ElMessage.success('删除成功')
      loadFiles()
    })
}

/* ========= 下载 ========= */
const downloadFile = (file) => {
  window.open('/' + file.filePath)
}

/* ========= 生命周期 ========= */
onMounted(async () => {
  await loadMyProjects()
  await loadProjectDetail()
  await loadFiles()
})
</script>

<style scoped>
.project-detail {
  padding: 20px;
}

.mb-20 {
  margin-bottom: 20px;
}

.project-select {
  display: flex;
  align-items: center;
}

.project-select .label {
  margin-right: 10px;
  font-weight: 500;
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
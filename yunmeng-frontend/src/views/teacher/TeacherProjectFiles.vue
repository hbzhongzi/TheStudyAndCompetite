<template>
  <div class="teacher-file-review">

    <!-- 学生项目选择 -->
    <el-card class="mb-16">
      <template #header>
        <span>学生项目选择</span>
      </template>

      <el-table
        :data="projects"
        v-loading="projectLoading"
        highlight-current-row
        @current-change="selectProject"
        row-key="id"
      >
        <el-table-column prop="title" label="项目名称" min-width="220" />
        <el-table-column prop="studentName" label="学生" width="120" />
        <el-table-column prop="type" label="类型" width="120" />
        <el-table-column prop="status" label="状态" width="120">
          <template #default="{ row }">
            <el-tag>{{ row.status }}</el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 文件列表 -->
    <el-card>
      <template #header>
        <span>
          项目成果文件
          <span v-if="currentProject">
            （{{ currentProject.title }} - {{ currentProject.studentName }}）
          </span>
        </span>
      </template>

      <el-table
        :data="files"
        v-loading="fileLoading"
        border
      >
        <el-table-column prop="originalName" label="文件名" min-width="240" />

        <el-table-column label="类型" width="120">
          <template #default="{ row }">
            <el-tag size="small">{{ row.category }}</el-tag>
          </template>
        </el-table-column>

        <el-table-column label="大小" width="120">
          <template #default="{ row }">
            {{ formatSize(row.size) }}
          </template>
        </el-table-column>

        <el-table-column label="审核状态" width="120">
          <template #default="{ row }">
            <el-tag :type="reviewType(row.reviewStatus)">
              {{ reviewText(row.reviewStatus) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="上传时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.createdAt) }}
          </template>
        </el-table-column>

        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button size="small" @click="preview(row)">预览</el-button>
            <el-button
              v-if="row.reviewStatus === 'pending'"
              size="small"
              type="success"
              @click="review(row, 'approved')"
            >
              通过
            </el-button>
            <el-button
              v-if="row.reviewStatus === 'pending'"
              size="small"
              type="danger"
              @click="review(row, 'rejected')"
            >
              驳回
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-empty
        v-if="!fileLoading && files.length === 0"
        description="暂无成果文件"
      />
    </el-card>

  </div>
</template>
<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { teacherService } from '@/services/teacherService'

/* ================= 项目列表 ================= */
const projects = ref([])
const projectLoading = ref(false)
const currentProject = ref(null)

/* ================= 文件列表 ================= */
const files = ref([])
const fileLoading = ref(false)

/* ================= 加载教师项目 ================= */
const fetchProjects = async () => {
  projectLoading.value = true
  try {
    const res = await teacherService.getTeacherProjects()
    projects.value = res.data.list || []

    // ✅ 默认选中第一个项目
    if (projects.value.length > 0) {
      selectProject(projects.value[0])
    }
  } finally {
    projectLoading.value = false
  }
}

/* ================= 选择项目 ================= */
const selectProject = (row) => {
  if (!row) return
  currentProject.value = row
  fetchFiles(row.id, row.studentId)
}

/* ================= 加载文件 ================= */
const fetchFiles = async (projectId, studentId) => {
  fileLoading.value = true
  try {
    const res = await teacherService.getStudentProjectsFiles({
      projectId,
      studentId
    })
    files.value = res.data.list || []
  } finally {
    fileLoading.value = false
  }
}

/* ================= 文件审核 ================= */
const review = async (file, status) => {
  await ElMessageBox.confirm(
    `确认${status === 'approved' ? '通过' : '驳回'}该文件？`,
    '提示',
    { type: 'warning' }
  )

  await teacherService.reviewProjectFile({
    fileId: file.id,
    status
  })

  ElMessage.success('审核完成')
  fetchFiles(currentProject.value.id, currentProject.value.studentId)
}

/* ================= 工具函数 ================= */
const formatSize = (s) =>
  s < 1024 ? `${s} B` :
  s < 1024 * 1024 ? `${(s / 1024).toFixed(1)} KB` :
  `${(s / 1024 / 1024).toFixed(1)} MB`

const formatTime = (t) =>
  t ? t.replace('T', ' ').substring(0, 19) : '-'

const reviewText = (s) =>
  s === 'approved' ? '已通过' :
  s === 'rejected' ? '已驳回' : '待审核'

const reviewType = (s) =>
  s === 'approved' ? 'success' :
  s === 'rejected' ? 'danger' : 'warning'

const preview = (file) => {
  window.open('/' + file.filePath)
}

onMounted(fetchProjects)
</script>
<style scoped>
.teacher-file-review {
  padding: 20px;
}

.mb-16 {
  margin-bottom: 16px;
}
</style>
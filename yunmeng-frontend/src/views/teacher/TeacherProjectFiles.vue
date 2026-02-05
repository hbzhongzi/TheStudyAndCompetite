<template>
  <div class="teacher-file-review">

    <!-- ================= 学生项目选择 ================= -->
    <el-card class="mb-16">
      <template #header>
        <span>学生项目选择</span>
      </template>

      <el-table
        ref="projectTableRef"
        :data="projects"
        v-loading="projectLoading"
        highlight-current-row
        row-key="id"
        @current-change="selectProject"
      >
        <!-- 项目名称 -->
        <el-table-column
          prop="title"
          label="项目名称"
          min-width="240"
          show-overflow-tooltip
        />

        <!-- 学生 -->
        <el-table-column label="学生" width="160">
          <template #default="{ row }">
            <div class="student-cell">
              <div class="name">{{ row.student?.name || '-' }}</div>
              <div class="sid">{{ row.student?.studentId || '' }}</div>
            </div>
          </template>
        </el-table-column>

        <!-- 类型 -->
        <el-table-column label="类型" width="120">
          <template #default="{ row }">
            <el-tag v-if="row.type" size="small">{{ row.type }}</el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>

        <!-- 状态 -->
        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="projectStatusTag(row.status)">
              {{ projectStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- ================= 文件列表 ================= -->
    <el-card>
      <template #header>
        <span>
          项目成果文件
          <span v-if="currentProject">
            （{{ currentProject.title }} - {{ currentProject.student?.name }}）
          </span>
        </span>
      </template>

      <el-table
        :data="files"
        v-loading="fileLoading"
        border
      >
        <!-- 文件名 -->
        <el-table-column
          prop="originalName"
          label="文件名"
          min-width="260"
          show-overflow-tooltip
        />

        <!-- 文件类型 -->
        <el-table-column label="类型" width="120">
          <template #default="{ row }">
            <el-tag size="small">
              {{ fileCategory(row) }}
            </el-tag>
          </template>
        </el-table-column>

        <!-- 大小 -->
        <el-table-column label="大小" width="120">
          <template #default="{ row }">
            {{ formatSize(row.size) }}
          </template>
        </el-table-column>

        <!-- 审核状态 -->
        <el-table-column label="审核状态" width="120">
          <template #default="{ row }">
            <el-tag :type="reviewTag(row)">
              {{ reviewText(row) }}
            </el-tag>
          </template>
        </el-table-column>

        <!-- 上传时间 -->
        <el-table-column label="上传时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.createdAt) }}
          </template>
        </el-table-column>

        <!-- 操作 -->
        <el-table-column label="操作" width="220">
          <template #default="{ row }">
            <el-button size="small" @click="preview(row)">
              预览
            </el-button>

            <el-button
              v-if="canReview(row)"
              size="small"
              type="success"
              @click="review(row, 'approved')"
            >
              通过
            </el-button>

            <el-button
              v-if="canReview(row)"
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
import { ref, onMounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { teacherService } from '@/services/teacherService'

/* ================= 项目 ================= */
const projects = ref([])
const projectLoading = ref(false)
const currentProject = ref(null)
const projectTableRef = ref(null)

/* ================= 文件 ================= */
const files = ref([])
const fileLoading = ref(false)

/* ================= 获取教师项目 ================= */
const fetchProjects = async () => {
  projectLoading.value = true
  try {
    const res = await teacherService.getTeacherProjects()
    projects.value = res.data.list || []

    // ✅ 默认选中第一个项目
    if (projects.value.length > 0) {
      nextTick(() => {
        projectTableRef.value?.setCurrentRow(projects.value[0])
        selectProject(projects.value[0])
      })
    }
  } finally {
    projectLoading.value = false
  }
}

/* ================= 选择项目 ================= */
const selectProject = (row) => {
  if (!row) return
  currentProject.value = row
  fetchFiles(row.id, row.student?.user_id)
}

/* ================= 获取文件 ================= */
const fetchFiles = async (projectId, user_id) => {
  fileLoading.value = true
  try {
// 调用时使用正确的参数名
    const res = await teacherService.getStudentProjectsFiles({
          id: projectId,          // 对应后端的 id
          student_id: user_id   // 对应后端的 user_id
    })
    files.value = res.data || []
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
  fetchFiles(currentProject.value.id, currentProject.value.student?.studentId)
}

/* ================= 工具函数 ================= */

/* 项目状态 */
const projectStatusText = s =>
  ({ submitted: '已提交', approved: '已通过', reviewing: '审核中', rejected: '已驳回' }[s] || s)

const projectStatusTag = s =>
  ({ submitted: 'warning', approved: 'success', reviewing: 'info', rejected: 'danger' }[s] || '')

/* 文件审核状态（核心修复点） */
const reviewText = row => {
  if (row.reviewStatus === 'approved') return '已通过'
  if (row.reviewStatus === 'rejected') return '已驳回'
  return '待审核'
}

const reviewTag = row => {
  if (row.reviewStatus === 'approved') return 'success'
  if (row.reviewStatus === 'rejected') return 'danger'
  return 'warning'
}

const canReview = row =>
  !row.reviewStatus || row.reviewStatus === 'pending' || row.status === 'draft'

/* 文件类型 */
const fileCategory = row => {
  if (row.category) return row.category
  if (row.ext?.includes('jpg') || row.ext?.includes('png')) return '图片'
  if (row.ext?.includes('pdf') || row.ext?.includes('doc')) return '文档'
  return '其他'
}

const formatSize = s =>
  s < 1024 ? `${s} B` :
  s < 1024 * 1024 ? `${(s / 1024).toFixed(1)} KB` :
  `${(s / 1024 / 1024).toFixed(1)} MB`

const formatTime = t =>
  t ? t.replace('T', ' ').substring(0, 19) : '-'

const preview = file => {
  window.open('/' + file.filePath)
}

/* ================= 生命周期 ================= */
onMounted(fetchProjects)
</script>

<style scoped>
.teacher-file-review {
  padding: 20px;
}

.mb-16 {
  margin-bottom: 16px;
}

.student-cell {
  line-height: 1.2;
}

.student-cell .name {
  font-weight: 500;
}

.student-cell .sid {
  font-size: 12px;
  color: #999;
}
</style>
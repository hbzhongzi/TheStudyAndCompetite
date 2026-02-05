<template>
  <div class="teacher-project-review">

    <!-- 左：教师项目列表 -->
    <el-card class="left-card">
      <template #header>
        <span>我的指导项目</span>
      </template>

      <el-table
        :data="projects"
        v-loading="listLoading"
        highlight-current-row
        row-key="id"
        @current-change="selectProject"
      >
        <el-table-column
          prop="title"
          label="项目名称"
          min-width="220"
          show-overflow-tooltip
        />

        <el-table-column label="学生" width="120">
          <template #default="{ row }">
            {{ row.studentName }}
          </template>
        </el-table-column>

        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusTag(row.status)">
              {{ statusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 右：项目详情 -->
    <el-card class="right-card" v-loading="detailLoading">
      <template #header>
        <span>项目详情</span>
      </template>

      <template v-if="project">

        <!-- 基本信息 -->
        <el-descriptions :column="2" border class="mb-16">
          <el-descriptions-item label="项目名称">
            {{ project.title }}
          </el-descriptions-item>
          <el-descriptions-item label="项目类型">
            {{ project.type }}
          </el-descriptions-item>
          <el-descriptions-item label="项目状态">
            <el-tag :type="statusTag(project.status)">
              {{ statusText(project.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="是否通过">
            <el-tag :type="project.isApproved ? 'success' : 'warning'">
              {{ project.isApproved ? '已通过' : '未通过' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">
            {{ formatTime(project.createdAt) }}
          </el-descriptions-item>
          <el-descriptions-item label="更新时间">
            {{ formatTime(project.updatedAt) }}
          </el-descriptions-item>
          <el-descriptions-item label="项目计划" :span="2">
            {{ project.plan || '-' }}
          </el-descriptions-item>
          <el-descriptions-item label="项目描述" :span="2">
            {{ project.description || '-' }}
          </el-descriptions-item>
        </el-descriptions>

        <!-- 学生 / 教师 -->
        <el-row :gutter="16" class="mb-16">
          <el-col :span="12">
            <el-descriptions title="学生信息" border>
              <el-descriptions-item label="姓名">
                {{ project.student.realName }}
              </el-descriptions-item>
              <el-descriptions-item label="学号">
                {{ project.student.studentId }}
              </el-descriptions-item>
              <el-descriptions-item label="学院">
                {{ project.student.department }}
              </el-descriptions-item>
              <el-descriptions-item label="邮箱">
                {{ project.student.email }}
              </el-descriptions-item>
            </el-descriptions>
          </el-col>

          <el-col :span="12">
            <el-descriptions title="指导教师" border>
              <el-descriptions-item label="姓名">
                {{ project.teacher.realName }}
              </el-descriptions-item>
              <el-descriptions-item label="学院">
                {{ project.teacher.department }}
              </el-descriptions-item>
              <el-descriptions-item label="邮箱">
                {{ project.teacher.email }}
              </el-descriptions-item>
            </el-descriptions>
          </el-col>
        </el-row>

        <!-- 项目审核 -->
        <el-card
          v-if="project.status === 'submitted'"
          class="mb-16"
          shadow="never"
        >
          <template #header>
            <span>项目审核</span>
          </template>

          <el-input
            v-model="reviewReason"
            type="textarea"
            rows="3"
            placeholder="如驳回，请填写审核意见（可选）"
          />

          <div class="mt-12">
            <el-button type="success" @click="reviewProject('approved')">
              通过审核
            </el-button>
            <el-button type="danger" @click="reviewProject('rejected')">
              驳回项目
            </el-button>
          </div>
        </el-card>

        <!-- 成果文件 -->
        <el-card shadow="never">
          <template #header>
            <span>项目成果文件</span>
          </template>

          <el-table :data="project.files" border>
            <el-table-column
              prop="originalName"
              label="文件名"
              min-width="240"
            />
            <el-table-column label="大小" width="120">
              <template #default="{ row }">
                {{ formatSize(row.size) }}
              </template>
            </el-table-column>
            <el-table-column label="审核状态" width="120">
              <template #default="{ row }">
                <el-tag :type="reviewTag(row.reviewStatus)">
                  {{ reviewText(row.reviewStatus) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="上传时间" width="180">
              <template #default="{ row }">
                {{ formatTime(row.uploadTime || row.createdAt) }}
              </template>
            </el-table-column>
            <el-table-column label="操作" width="220">
              <template #default="{ row }">
                <el-button size="small" @click="preview(row)">
                  预览
                </el-button>
                <el-button
                  size="small"
                  type="success"
                  v-if="row.reviewStatus !== 'approved'"
                  @click="reviewFile(row.id, 'approved')"
                >
                  通过
                </el-button>
                <el-button
                  size="small"
                  type="danger"
                  v-if="row.reviewStatus !== 'rejected'"
                  @click="reviewFile(row.id, 'rejected')"
                >
                  驳回
                </el-button>
              </template>
            </el-table-column>
          </el-table>

          <el-empty
            v-if="project.files.length === 0"
            description="暂无成果文件"
          />
        </el-card>

      </template>
    </el-card>
  </div>
</template>
<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { teacherService } from '@/services/teacherService'

/* ================= 项目列表 ================= */

const projects = ref([])
const listLoading = ref(false)
const detailLoading = ref(false)
const currentProject = ref(null)
const project = ref(null)
const reviewReason = ref('')

/* ================= 加载教师项目列表 ================= */

const fetchProjects = async () => {
  listLoading.value = true
  try {
    const res = await teacherService.getTeacherProjects()
    projects.value = res.data.list || []

    // ✅ 默认选中第一个
    if (projects.value.length > 0) {
      selectProject(projects.value[0])
    }
  } finally {
    listLoading.value = false
  }
}

/* ================= 选择项目 ================= */

const selectProject = async (row) => {
  if (!row) return
  currentProject.value = row
  await fetchProjectDetail(row.id)
}

/* ================= 项目详情 ================= */

const fetchProjectDetail = async (projectId) => {
  detailLoading.value = true
  try {
    const res = await teacherService.getProjectDetail(projectId)
    project.value = res.data
  } finally {
    detailLoading.value = false
  }
}

/* ================= 项目审核 ================= */

const reviewProject = async (status) => {
  await ElMessageBox.confirm(
    `确认${status === 'approved' ? '通过' : '驳回'}该项目？`,
    '提示',
    { type: 'warning' }
  )

  await teacherService.reviewProject({
    projectId: project.value.id,
    status,
    reason: reviewReason.value
  })

  ElMessage.success('项目审核完成')
  fetchProjects()
}

/* ================= 文件审核 ================= */

const reviewFile = async (fileId, status) => {
  await teacherService.reviewProjectFile({ fileId, status })
  ElMessage.success('文件审核完成')
  fetchProjectDetail(project.value.id)
}

/* ================= 工具函数 ================= */

const statusText = s =>
  ({ submitted: '待审核', approved: '已通过', rejected: '已驳回' }[s] || s)

const statusTag = s =>
  ({ submitted: 'warning', approved: 'success', rejected: 'danger' }[s] || '')

const reviewText = s =>
  ({ approved: '已通过', rejected: '已驳回' }[s] || '待审核')

const reviewTag = s =>
  ({ approved: 'success', rejected: 'danger' }[s] || 'warning')

const formatTime = t =>
  t ? t.replace('T', ' ').substring(0, 19) : '-'

const formatSize = s =>
  s < 1024 ? `${s} B`
    : s < 1024 * 1024 ? `${(s / 1024).toFixed(1)} KB`
    : `${(s / 1024 / 1024).toFixed(1)} MB`

const preview = file => {
  window.open(file.fileUrl)
}

onMounted(fetchProjects)
</script>
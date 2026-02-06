<template>
  <div class="project-extensions">
    <el-card>
      <template #header>
        <div class="header-content">
          <span>项目延期申请管理</span>
          <el-button type="primary" @click="openApplyDialog">
            申请延期
          </el-button>
        </div>
      </template>

      <!-- 延期申请列表 -->
      <el-table
        :data="extensions"
        style="width: 100%"
        v-loading="loading"
      >
        <el-table-column prop="projectTitle" label="项目名称" min-width="180" />

        <el-table-column label="原完成时间" width="150">
          <template #default="scope">
            {{ formatDate(scope.row.originalFinishTime) }}
          </template>
        </el-table-column>

        <el-table-column label="申请完成时间" width="150">
          <template #default="scope">
            {{ formatDate(scope.row.requestedFinishTime) }}
          </template>
        </el-table-column>

        <el-table-column label="延期天数" width="120">
          <template #default="scope">
            <el-tag type="warning">
              +{{ calcDays(scope.row) }} 天
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column
          prop="applyReason"
          label="申请原因"
          min-width="200"
          show-overflow-tooltip
        />

        <el-table-column label="状态" width="100">
          <template #default="scope">
            <el-tag :type="statusType(scope.row.status)">
              {{ statusLabel(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="申请时间" width="160">
          <template #default="scope">
            {{ formatDateTime(scope.row.createdAt) }}
          </template>
        </el-table-column>

        <el-table-column label="指导教师" width="120">
          <template #default="scope">
            {{ scope.row.teacherName || '—' }}
          </template>
        </el-table-column>

        <el-table-column label="操作" width="120">
          <template #default="scope">
            <el-button size="small" @click="viewDetail(scope.row)">
              查看
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-empty
        v-if="!loading && extensions.length === 0"
        description="暂无延期申请记录"
      />
    </el-card>

    <!-- 申请延期弹窗 -->
    <el-dialog
      v-model="applyDialogVisible"
      title="申请项目延期"
      width="40%"
    >
      <el-form :model="form" label-width="120px">
        <!-- 项目选择 -->
        <el-form-item label="选择项目">
          <el-select
            v-model="form.projectId"
            placeholder="请选择项目"
            style="width: 100%"
          >
            <el-option
              v-for="item in projects"
              :key="item.id"
              :label="item.title"
              :value="item.id"
            />
          </el-select>
        </el-form-item>

        <!-- 原完成时间 -->
<el-form-item label="原完成时间">
  <el-date-picker
    v-model="form.originalFinishTime"
    type="date"
    disabled
    format="YYYY-MM-DD"
    value-format="YYYY-MM-DD"
  />
</el-form-item>

        <!-- 申请完成时间 -->
        <el-form-item label="申请完成时间">
<el-date-picker
  v-model="form.requestedFinishTime"
  type="date"
  format="YYYY-MM-DD"
  value-format="YYYY-MM-DDT00:00:00+08:00"
/>
        </el-form-item>

        <!-- 原因 -->
        <el-form-item label="延期原因">
          <el-input
            v-model="form.applyReason"
            type="textarea"
            :rows="3"
            placeholder="请输入延期原因"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="applyDialogVisible = false">取消</el-button>
        <el-button
          type="primary"
          :loading="submitting"
          @click="submitApply"
        >
          提交申请
        </el-button>
      </template>
    </el-dialog>

    <!-- 详情弹窗 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="延期申请详情"
      width="50%"
    >
      <el-descriptions v-if="current" :column="2" border>
        <el-descriptions-item label="项目名称">
          {{ current.projectTitle }}
        </el-descriptions-item>

        <el-descriptions-item label="状态">
          <el-tag :type="statusType(current.status)">
            {{ statusLabel(current.status) }}
          </el-tag>
        </el-descriptions-item>

        <el-descriptions-item label="原完成时间">
          {{ formatDate(current.originalFinishTime) }}
        </el-descriptions-item>

        <el-descriptions-item label="申请完成时间">
          {{ formatDate(current.requestedFinishTime) }}
        </el-descriptions-item>

        <el-descriptions-item label="延期原因" :span="2">
          {{ current.applyReason }}
        </el-descriptions-item>

        <el-descriptions-item label="审核意见" :span="2">
          {{ current.reviewReason || '暂无' }}
        </el-descriptions-item>
      </el-descriptions>

      <template #footer>
        <el-button @click="detailDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import studentService from '../../services/studentService'

/* 列表 */
const extensions = ref([])
const loading = ref(false)

/* 学生项目 */
const projects = ref([])

/* 弹窗 */
const applyDialogVisible = ref(false)
const detailDialogVisible = ref(false)
const submitting = ref(false)

/* 当前详情 */
const current = ref(null)

/* 表单 */
const form = ref({
  projectId: null,
  originalFinishTime: '',
  requestedFinishTime: '',
  applyReason: ''
})

/* 加载延期列表 */
const loadList = async () => {
  loading.value = true
  try {
    const res = await studentService.getMyExtensionApplications()
    extensions.value = res.data.list || []
  } finally {
    loading.value = false
  }
}

/* 加载学生项目 */
const loadMyProjects = async () => {
  const res = await studentService.getMyProjects()
  projects.value = res.data.list || []

  if (projects.value.length > 0) {
    const first = projects.value[0]
    form.value.projectId = first.id
    form.value.originalFinishTime = first.finishTime
  }
}

onMounted(loadList)

/* 监听项目切换 */
watch(() => form.value.projectId, (val) => {
  const project = projects.value.find(p => p.id === val)
  if (project) {
    form.value.originalFinishTime = project.finishTime
  }
})

/* 工具方法 */
const formatDate = v => v ? v.slice(0, 10) : ''
const formatDateTime = v => v ? v.replace('T', ' ').slice(0, 19) : ''

const calcDays = row => {
  const a = new Date(row.originalFinishTime)
  const b = new Date(row.requestedFinishTime)
  return Math.max(0, Math.ceil((b - a) / 86400000))
}

const statusLabel = s => ({
  pending: '待审核',
  approved: '已通过',
  rejected: '已驳回'
}[s] || s)

const statusType = s => ({
  pending: 'warning',
  approved: 'success',
  rejected: 'danger'
}[s] || 'info')

/* 打开弹窗 */
const openApplyDialog = async () => {
  form.value = {
    projectId: null,
    originalFinishTime: '',
    requestedFinishTime: '',
    applyReason: ''
  }

  await loadMyProjects()
  applyDialogVisible.value = true
}

/* 提交 */
const submitApply = async () => {
  if (!form.value.projectId) {
    return ElMessage.warning('请选择项目')
  }
  if (!form.value.requestedFinishTime) {
    return ElMessage.warning('请选择申请完成时间')
  }

  submitting.value = true
  try {
    await studentService.createExtensionApplication(form.value)
    ElMessage.success('延期申请提交成功')
    applyDialogVisible.value = false
    loadList()
  } finally {
    submitting.value = false
  }
}

/* 查看详情 */
const viewDetail = row => {
  current.value = row
  detailDialogVisible.value = true
}
</script>

<style scoped>
.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
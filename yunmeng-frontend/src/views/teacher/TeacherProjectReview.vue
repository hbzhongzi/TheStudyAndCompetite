<template>
  <div class="project-review-page">
    <!-- 顶部 -->
    <div class="page-header">
      <div>
        <h2>项目延期审核</h2>
        <p class="sub">学生项目延期申请审核与处理</p >
      </div>
      <el-button type="primary" @click="fetchList">
        刷新
      </el-button>
    </div>

    <!-- 统计看板 -->
    <el-row :gutter="16" class="stats">
      <el-col :span="8" v-for="stat in stats" :key="stat.label">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <span class="stat-label">{{ stat.label }}</span>
            <span class="stat-num">{{ stat.value }}</span>
          </div>
        </el-card>
      </el-col>
    </el-row>


    <!-- 表格 -->
    <el-card>
      <el-table
        :data="list"
        v-loading="loading"
        stripe
        highlight-current-row
      >
        <el-table-column
          prop="projectTitle"
          label="项目名称"
          min-width="220"
        />

        <el-table-column
          prop="studentName"
          label="学生"
          width="120"
        />

        <el-table-column label="延期时间">
          <template #default="{ row }">
            <div class="date-col">
              <span class="old">
                {{ formatTime(row.originalFinishTime) }}
              </span>
              <span class="new">
                → {{ formatTime(row.requestedFinishTime) }}
              </span>
            </div>
          </template>
        </el-table-column>

        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusType(row.status)">
              {{ statusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="申请时间" width="170">
          <template #default="{ row }">
            {{ formatTime(row.createdAt) }}
          </template>
        </el-table-column>

        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="view(row)">
              详情
            </el-button>

            <el-button
              v-if="row.status === 'pending'"
              size="small"
              type="success"
              @click="openReview(row, 'approved')"
            >
              通过
            </el-button>

            <el-button
              v-if="row.status === 'pending'"
              size="small"
              type="danger"
              @click="openReview(row, 'rejected')"
            >
              驳回
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination
        class="pagination"
        v-model:current-page="query.page"
        v-model:page-size="query.size"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="fetchList"
      />
    </el-card>

    <!-- 审核弹窗 -->
    <el-dialog v-model="reviewVisible" width="480px">
      <template #header>
        <strong>确认审核结果</strong>
      </template>

      <el-alert
        :title="
          reviewStatus === 'approved'
            ? '确认通过延期申请？'
            : '确认驳回延期申请？'
        "
        :type="reviewStatus === 'approved' ? 'success' : 'error'"
        show-icon
        class="mb-12"
      />

      <el-input
        v-model="reviewReason"
        type="textarea"
        rows="4"
        placeholder="请输入审核意见（可选）"
      />

      <template #footer>
        <el-button @click="reviewVisible = false">取消</el-button>
        <el-button type="primary" @click="submitReview">
          确认提交
        </el-button>
      </template>
    </el-dialog>

    <!-- 查看详情 -->
    <el-dialog v-model="viewVisible" title="延期申请详情" width="600px">
      <el-descriptions border :column="2">
        <el-descriptions-item label="项目名称">
          {{ viewRow.projectTitle }}
        </el-descriptions-item>

        <el-descriptions-item label="学生">
          {{ viewRow.studentName }}
        </el-descriptions-item>

        <el-descriptions-item label="原截止时间">
          {{ formatTime(viewRow.originalFinishTime) }}
        </el-descriptions-item>

        <el-descriptions-item label="申请截止时间">
          {{ formatTime(viewRow.requestedFinishTime) }}
        </el-descriptions-item>

        <el-descriptions-item label="申请时间">
          {{ formatTime(viewRow.createdAt) }}
        </el-descriptions-item>

        <el-descriptions-item label="状态">
          <el-tag :type="statusType(viewRow.status)">
            {{ statusLabel(viewRow.status) }}
          </el-tag>
        </el-descriptions-item>

        <el-descriptions-item label="申请理由" :span="2">
          {{ viewRow.applyReason || '—' }}
        </el-descriptions-item>

        <el-descriptions-item
          v-if="viewRow.status !== 'pending'"
          label="审核意见"
          :span="2"
        >
          {{ viewRow.reviewReason || '—' }}
        </el-descriptions-item>
      </el-descriptions>

      <template #footer>
        <el-button type="primary" @click="viewVisible = false">
          确定
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import teacherService from '@/services/teacherService'

/* ========= 数据 ========= */
const list = ref([])
const total = ref(0)
const loading = ref(false)

const query = ref({
  page: 1,
  size: 10,
  status: '',
  keyword: ''
})

const stats = ref([
  { label: '待审核', value: 0 },
  { label: '已通过', value: 0 },
  { label: '已拒绝', value: 0 }
])

/* ========= 审核 ========= */
const reviewVisible = ref(false)
const currentRow = ref(null)
const reviewStatus = ref('')
const reviewReason = ref('')

/* ========= 查看 ========= */
const viewVisible = ref(false)
const viewRow = ref({})

/* ========= 方法 ========= */
const fetchList = async () => {
  loading.value = true
  try {
    const res = await teacherService.getProjectReviews(query.value)
    list.value = res.data.list || []
    total.value = res.data.total || 0

    stats.value[0].value = list.value.filter(i => i.status === 'pending').length
    stats.value[1].value = list.value.filter(i => i.status === 'approved').length
    stats.value[2].value = list.value.filter(i => i.status === 'rejected').length
  } finally {
    loading.value = false
  }
}

const statusLabel = s =>
  s === 'pending' ? '待审核' : s === 'approved' ? '已通过' : '已拒绝'

const statusType = s =>
  s === 'pending' ? 'warning' : s === 'approved' ? 'success' : 'danger'

const formatTime = t => {
  if (!t) return '-'
  return t.replace('T', ' ').substring(0, 19)
}

const view = row => {
  viewRow.value = row
  viewVisible.value = true
}

const openReview = (row, status) => {
  currentRow.value = row
  reviewStatus.value = status
  reviewReason.value = ''
  reviewVisible.value = true
}

const submitReview = async () => {
  await teacherService.updateProjectReviews({
    ApplicationID: currentRow.value.id,
    Action: reviewStatus.value,
    Reason: reviewReason.value
  })
  ElMessage.success('审核完成')
  reviewVisible.value = false
  fetchList()
}

onMounted(fetchList)
</script>

<style scoped>
.project-review-page {
  padding: 16px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.page-header h2 {
  margin: 0;
}

.sub {
  font-size: 13px;
  color: #888;
}

.stats {
  margin-bottom: 16px;
}

.stat-card {
  text-align: center;
}

.stat-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stat-label {
  font-size: 14px;
  color: #666;
}

.stat-num {
  font-size: 26px;
  font-weight: bold;
}

.search-card {
  margin-bottom: 12px;
}

.date-col {
  display: flex;
  flex-direction: column;
  font-size: 12px;
}

.date-col .old {
  color: #999;
  text-decoration: line-through;
}

.date-col .new {
  color: #409eff;
}

.pagination {
  margin-top: 16px;
  text-align: right;
}
</style>
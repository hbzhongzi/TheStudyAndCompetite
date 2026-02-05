<template>
  <div class="project-review">
    <el-card>
      <!-- 头部 -->
      <template #header>
        <div class="header-content">
          <span>项目延期审核</span>
          <el-button type="primary" @click="fetchList">刷新</el-button>
        </div>
      </template>

      <!-- 统计 -->
      <el-row :gutter="20" class="stats">
        <el-col :span="6" v-for="stat in stats" :key="stat.label">
          <el-card class="stat-card">
            <h4>{{ stat.label }}</h4>
            <p class="stat-num">{{ stat.value }}</p >
          </el-card>
        </el-col>
      </el-row>

      <!-- 搜索 -->
      <el-form inline class="search-bar">
        <el-form-item>
          <el-input
            v-model="query.keyword"
            placeholder="项目名 / 学生姓名"
            clearable
            @input="fetchList"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-form-item>

        <el-form-item>
          <el-select
            v-model="query.status"
            placeholder="状态"
            clearable
            @change="fetchList"
          >
            <el-option label="待审核" value="pending" />
            <el-option label="已通过" value="approved" />
            <el-option label="已拒绝" value="rejected" />
          </el-select>
        </el-form-item>
      </el-form>

      <!-- 表格 -->
      <el-table :data="list" v-loading="loading" style="width: 100%">
        <el-table-column prop="projectTitle" label="项目名称" min-width="200" />
        <el-table-column prop="studentName" label="学生" width="120" />

        <el-table-column label="原截止时间" width="170">
          <template #default="scope">
            {{ formatTime(scope.row.originalFinishTime) }}
          </template>
        </el-table-column>

        <el-table-column label="申请截止时间" width="170">
          <template #default="scope">
            {{ formatTime(scope.row.requestedFinishTime) }}
          </template>
        </el-table-column>

        <el-table-column label="状态" width="100">
          <template #default="scope">
            <el-tag :type="statusType(scope.row.status)">
              {{ statusLabel(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="申请时间" width="170">
          <template #default="scope">
            {{ formatTime(scope.row.createdAt) }}
          </template>
        </el-table-column>

        <el-table-column label="操作" width="220" fixed="right">
          <template #default="scope">
            <el-button size="small" @click="view(scope.row)">查看</el-button>

            <el-button
              v-if="scope.row.status === 'pending'"
              size="small"
              type="success"
              @click="openReview(scope.row, 'approved')"
            >
              通过
            </el-button>

            <el-button
              v-if="scope.row.status === 'pending'"
              size="small"
              type="danger"
              @click="openReview(scope.row, 'rejected')"
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
    <el-dialog v-model="reviewVisible" title="延期审核" width="500px">
      <el-form label-width="90px">
        <el-form-item label="审核结果">
          <el-tag :type="reviewStatus === 'approved' ? 'success' : 'danger'">
            {{ statusLabel(reviewStatus) }}
          </el-tag>
        </el-form-item>

        <el-form-item label="审核意见">
          <el-input
            v-model="reviewReason"
            type="textarea"
            rows="4"
            placeholder="请输入审核意见"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="reviewVisible = false">取消</el-button>
        <el-button type="primary" @click="submitReview">提交</el-button>
      </template>
    </el-dialog>
  </div>

<!-- 查看延期申请 -->
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

</template>


<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import { teacherService } from '@/services/teacherService'

/* ========= 表格 ========= */
const list = ref([])
const total = ref(0)
const loading = ref(false)

const query = ref({
  page: 1,
  size: 10,
  status: '',
  keyword: ''
})

/* ========= 统计 ========= */
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

/* ========= 获取列表 ========= */
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

/* ========= 状态 ========= */
const statusLabel = s =>
  s === 'pending' ? '待审核' : s === 'approved' ? '已通过' : '已拒绝'

const statusType = s =>
  s === 'pending' ? 'warning' : s === 'approved' ? 'success' : 'danger'

/* ========= 时间 ========= */
const formatTime = t => {
  if (!t) return '-'
  return t.replace('T', ' ').substring(0, 19)
}

/* ========= 查看 ========= */
const view = row => {
  viewRow.value = row
  viewVisible.value = true
}

/* ========= 审核 ========= */
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
.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stats {
  margin-bottom: 16px;
}

.stat-card {
  text-align: center;
}

.stat-num {
  font-size: 22px;
  font-weight: bold;
}

.search-bar {
  margin-bottom: 12px;
}

.pagination {
  margin-top: 16px;
  text-align: right;
}
</style>
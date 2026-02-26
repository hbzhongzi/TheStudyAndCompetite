<template>
  <div class="competition-judging">

    <!-- 左侧：已分配竞赛 -->
    <div class="left-panel">
      <h3>已分配竞赛</h3>

      <el-menu
        :default-active="activeCompetitionId?.toString()"
        @select="handleSelectCompetition"
      >
        <el-menu-item
          v-for="item in competitionList"
          :key="item.competition.id"
          :index="item.competition.id.toString()"
        >
          {{ item.competition.title }}
        </el-menu-item>
      </el-menu>
    </div>

    <!-- 右侧：作品评审 -->
    <div class="right-panel">

      <div class="page-header">
        <h2>竞赛作品评审</h2>
      </div>

      <el-table
        :data="submissionList"
        border
        stripe
        v-loading="loading"
      >

        <el-table-column label="学生" width="120">
          <template #default="{ row }">
            {{ row.student?.realName }}
          </template>
        </el-table-column>

        <el-table-column label="学院" width="120">
          <template #default="{ row }">
            {{ row.student?.department }}
          </template>
        </el-table-column>

        <el-table-column prop="version" label="版本" width="80" />

        <el-table-column label="提交时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.submit_time) }}
          </template>
        </el-table-column>

        <el-table-column label="是否已查看" width="100">
          <template #default="{ row }">
            <el-tag :type="row.teacher_viewed ? 'success' : 'warning'">
              {{ row.teacher_viewed ? '已查看' : '未查看' }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="340">
          <template #default="{ row }">

            <el-button size="small" @click="openDetailDialog(row)">
              查看
            </el-button>

            <el-button size="small" type="primary" @click="downloadFile(row)">
              下载
            </el-button>

            <el-button size="small" type="success" @click="openReviewDialog(row)">
              评审
            </el-button>

            <el-button size="small" type="warning" @click="openScoreHistory(row)">
              评审记录
            </el-button>

          </template>
        </el-table-column>

      </el-table>

      <el-pagination
        background
        layout="prev, pager, next"
        :total="total"
        :page-size="10"
        style="margin-top:20px; text-align:right"
        @current-change="handlePageChange"
      />

    </div>

    <!-- 评审弹窗 -->
    <el-dialog
      v-model="reviewDialogVisible"
      title="作品评审"
      width="500px"
    >
      <el-form :model="reviewForm">

        <el-form-item label="学生">
          {{ currentRow?.student?.realName }}
        </el-form-item>

        <el-form-item label="评分">
          <el-input-number
            v-model="reviewForm.score"
            :min="0"
            :max="100"
          />
        </el-form-item>

        <el-form-item label="评审意见">
          <el-input
            type="textarea"
            v-model="reviewForm.comment"
            :rows="4"
          />
        </el-form-item>

      </el-form>

      <template #footer>
        <el-button @click="reviewDialogVisible=false">取消</el-button>
        <el-button type="primary" @click="submitReview">
          提交评审
        </el-button>
      </template>
    </el-dialog>

  </div>

  <!-- 作品详情弹窗 -->
  <el-dialog
    v-model="detailDialogVisible"
    title="作品详情"
    width="600px"
  >
    <div v-if="currentDetail">

      <el-descriptions :column="2" border>

        <el-descriptions-item label="学生姓名">
          {{ currentDetail.student?.realName }}
        </el-descriptions-item>

        <el-descriptions-item label="学院">
          {{ currentDetail.student?.department }}
        </el-descriptions-item>

        <el-descriptions-item label="邮箱">
          {{ currentDetail.student?.email }}
        </el-descriptions-item>

        <el-descriptions-item label="电话">
          {{ currentDetail.student?.phone }}
        </el-descriptions-item>

        <el-descriptions-item label="版本">
          {{ currentDetail.version }}
        </el-descriptions-item>

        <el-descriptions-item label="提交时间">
          {{ formatDate(currentDetail.submit_time) }}
        </el-descriptions-item>

        <el-descriptions-item label="文件大小">
          {{ (currentDetail.file_size / 1024).toFixed(2) }} KB
        </el-descriptions-item>

      </el-descriptions>

      <el-divider />

      <h4>作品描述</h4>
      <p>{{ currentDetail.description }}</p >

      <el-divider />

      <el-button type="primary" @click="downloadFile(currentDetail)">
        下载文件
      </el-button>

    </div>
  </el-dialog>

  <!-- 评分记录弹窗 -->
  <el-dialog
    v-model="scoreHistoryVisible"
    title="评审记录"
    width="650px"
  >
    <el-table
      :data="scoreHistoryList"
      border
      stripe
    >
      <el-table-column label="评审人" width="120">
        <template #default="{ row }">
          {{ row.judge?.realName }}
        </template>
      </el-table-column>

      <el-table-column label="学院" width="120">
        <template #default="{ row }">
          {{ row.judge?.department }}
        </template>
      </el-table-column>

      <el-table-column prop="score" label="评分" width="80" />

      <el-table-column prop="comment" label="评语" />

      <el-table-column label="评分时间" width="160">
        <template #default="{ row }">
          {{ formatDate(row.scored_at) }}
        </template>
      </el-table-column>

    </el-table>
  </el-dialog>

</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import teacherService from '@/services/teacherService'

const loading = ref(false)

const competitionList = ref([])
const activeCompetitionId = ref(null)

const submissionList = ref([])
const total = ref(0)
const currentPage = ref(1)

const detailDialogVisible = ref(false)
const currentDetail = ref(null)

const reviewDialogVisible = ref(false)
const reviewForm = ref({
  score: null,
  comment: ''
})
const currentRow = ref(null)

const scoreHistoryVisible = ref(false)
const scoreHistoryList = ref([])

const openDetailDialog = (row) => {
  currentDetail.value = row
  detailDialogVisible.value = true
}

const openScoreHistory = async (row) => {
  try {
    const res = await teacherService.getSubmissionScores(row.id)
    scoreHistoryList.value = res.data
    scoreHistoryVisible.value = true
  } catch (error) {
    ElMessage.error('获取评分记录失败')
  }
}

const loadCompetitions = async () => {
  try {
    const res = await teacherService.getMyTasks()
    if (res.code === 200 && res.data.list.length > 0) {
      competitionList.value = res.data.list
      activeCompetitionId.value = res.data.list[0].competition.id
      loadSubmissions()
    }
  } catch (error) {
    ElMessage.error('加载竞赛任务失败')
  }
}

const loadSubmissions = async () => {
  if (!activeCompetitionId.value) return
  loading.value = true
  try {
    const res = await teacherService.getCompetitionSubmissions(
      activeCompetitionId.value,
      { page: currentPage.value, size: 10 }
    )
    submissionList.value = res.data.list
    total.value = res.data.total
  } catch (error) {
    ElMessage.error('加载作品失败')
  }
  loading.value = false
}

const handleSelectCompetition = (id) => {
  activeCompetitionId.value = id
  currentPage.value = 1
  loadSubmissions()
}

const handlePageChange = (page) => {
  currentPage.value = page
  loadSubmissions()
}

const formatDate = (time) => {
  if (!time) return '-'
  return time.replace('T', ' ').substring(0, 16)
}

const downloadFile = (row) => {
  window.open(
    `http://localhost:8080/${row.file_url.replace(/\\/g, '/')}`
  )
}

const openReviewDialog = (row) => {
  currentRow.value = row
  reviewForm.value.score = null
  reviewForm.value.comment = ''
  reviewDialogVisible.value = true
}

const submitReview = async () => {
  try {
    await teacherService.submitSubmissionReview(
      currentRow.value.id,
      reviewForm.value
    )
    ElMessage.success('提交评分成功')
    reviewDialogVisible.value = false
    loadSubmissions()
  } catch (error) {
    ElMessage.error('评分失败')
  }
}

onMounted(loadCompetitions)
</script>

<style scoped>
.competition-judging {
  display: flex;
  gap: 20px;
}

.left-panel {
  width: 260px;
  background: #fff;
  padding: 15px;
  border-radius: 8px;
}

.right-panel {
  flex: 1;
  background: #fff;
  padding: 20px;
  border-radius: 8px;
}
</style>
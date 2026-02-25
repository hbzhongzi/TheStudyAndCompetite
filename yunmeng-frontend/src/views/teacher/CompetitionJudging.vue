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

        <el-table-column label="操作" width="240">
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

      <el-descriptions-item label="是否锁定">
        <el-tag :type="currentDetail.locked ? 'danger' : 'success'">
          {{ currentDetail.locked ? '已锁定' : '未锁定' }}
        </el-tag>
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

</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import  teacherService from '@/services/teacherService'

const loading = ref(false)
 
// 左侧任务
const competitionList = ref([])
const activeCompetitionId = ref(null)

// 右侧提交列表
const submissionList = ref([])
const total = ref(0)
const currentPage = ref(1)

// 详情弹窗
const detailDialogVisible = ref(false)
const currentDetail = ref(null)

const openDetailDialog = (row) => {
  currentDetail.value = row
  detailDialogVisible.value = true
}

// 评审
const reviewDialogVisible = ref(false)
const reviewForm = ref({
  score: null,
  comment: ''
})
const currentRow = ref(null)


// 加载已分配竞赛
const loadCompetitions = async () => {
  try {
    const res = await teacherService.getMyTasks()

    if (res.code === 200 && res.data.list.length > 0) {
      competitionList.value = res.data.list

      // ✅ 默认选中第一个
      const first = res.data.list[0]
      activeCompetitionId.value = first.competition.id

      loadSubmissions()
    }
  } catch (error) {
    ElMessage.error('加载竞赛任务失败')
  }
}


// 加载作品列表
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


// 切换竞赛
const handleSelectCompetition = (id) => {
  activeCompetitionId.value = id
  currentPage.value = 1
  loadSubmissions()
}


// 分页
const handlePageChange = (page) => {
  currentPage.value = page
  loadSubmissions()
}


// 时间格式
const formatDate = (time) => {
  if (!time) return '-'
  return time.replace('T', ' ').substring(0, 16)
}


// 文件预览
const previewFile = (row) => {
  window.open(
    `http://localhost:8080/${row.file_url.replace(/\\/g, '/')}`
  )
}


// 下载
const downloadFile = (row) => {
  window.open(
    `http://localhost:8080/${row.file_url.replace(/\\/g, '/')}`
  )
}


// 打开评审
const openReviewDialog = (row) => {
  currentRow.value = row
  reviewForm.value.score = row.scores || 0
  reviewForm.value.comment = row.teacher_feedback || ''
  reviewDialogVisible.value = true
}


// 提交评审
const submitReview = async () => {
  try {
    await teacherService.submitSubmissionReview(
      currentRow.value.id,
      reviewForm.value
    )
    ElMessage.success('评审成功')
    reviewDialogVisible.value = false
    loadSubmissions()
  } catch (error) {
    ElMessage.error('评审失败')
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
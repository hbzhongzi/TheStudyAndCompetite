<template>
  <div class="competition-review">

    <!-- 左侧竞赛列表 -->
    <div class="left-panel">
      <h3>竞赛列表</h3>

      <el-menu
        :default-active="activeCompetitionId?.toString()"
        @select="handleSelectCompetition"
      >
        <el-menu-item
          v-for="item in competitionList"
          :key="item.id"
          :index="item.id.toString()"
        >
          {{ item.title }}
        </el-menu-item>
      </el-menu>
    </div>

    <!-- 右侧报名审核 -->
    <div class="right-panel">

      <h2>报名审核</h2>

      <el-table
        :data="registerList"
        border
        stripe
        v-loading="loading"
      >
        <el-table-column prop="id" label="ID" width="80" />

        <el-table-column label="队伍名称">
          <template #default="{ row }">
            {{ row.team_name }}
          </template>
        </el-table-column>

        <el-table-column prop="team_leader" label="队长ID" width="100" />

        <el-table-column label="报名时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.register_time) }}
          </template>
        </el-table-column>

        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <el-tag
              :type="row.status === 'approved'
                ? 'success'
                : row.status === 'pending'
                ? 'warning'
                : 'danger'"
            >
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="200">
          <template #default="{ row }">

            <template v-if="row.status === 'pending'">

              <el-button
                size="small"
                type="success"
                @click="openReviewDialog(row.id, 'approved')"
              >
                通过
              </el-button>

              <el-button
                size="small"
                type="danger"
                @click="openReviewDialog(row.id, 'rejected')"
              >
                拒绝
              </el-button>

            </template>

            <span v-else>-</span>

          </template>
        </el-table-column>

      </el-table>

      <el-pagination
        background
        layout="prev, pager, next"
        :total="total"
        :page-size="10"
        style="margin-top:20px;text-align:right"
        @current-change="handlePageChange"
      />

    </div>
    <!-- 审核弹窗 -->
  <el-dialog
    v-model="reviewDialogVisible"
    title="报名审核"
    width="420px"
  >
    <el-form :model="reviewForm" label-width="90px">

      <el-form-item label="审核结果">
        <el-select v-model="reviewForm.status" style="width:100%">
          <el-option label="通过" value="approved" />
          <el-option label="拒绝" value="rejected" />
        </el-select>
      </el-form-item>

      <el-form-item label="审核意见">
        <el-input
          type="textarea"
          v-model="reviewForm.common"
          :rows="4"
          placeholder="请输入审核意见"
        />
      </el-form-item>

    </el-form>

    <template #footer>
      <el-button @click="reviewDialogVisible = false">
        取消
      </el-button>
      <el-button type="primary" @click="submitReview">
        提交审核
      </el-button>
    </template>

  </el-dialog>
  </div>

</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import competitionService from '@/services/competitionService'


const competitionList = ref([])
const activeCompetitionId = ref(null)

const registerList = ref([])
const total = ref(0)
const currentPage = ref(1)
const loading = ref(false)

// 审核弹窗
const reviewDialogVisible = ref(false)
const reviewForm = ref({
  id: null,
  status: '',
  common: ''
})


// ========================
// 加载竞赛列表
// ========================
const loadCompetitions = async () => {
  try {
    const res = await competitionService.getCompetitions({
      page: 1,
      size: 50,
      is_open: 1
    })

    competitionList.value = res.data.list

    if (res.data.list.length > 0) {
      activeCompetitionId.value = res.data.list[0].id
      loadRegisters()
    }

  } catch (error) {
    ElMessage.error('加载竞赛失败')
  }
}


// ========================
// 加载报名列表
// ========================
const loadRegisters = async () => {
  if (!activeCompetitionId.value) return

  loading.value = true

  try {
    const res = await competitionService.getCompetitionRegisters(
      activeCompetitionId.value,
      { page: currentPage.value, size: 10 }
    )

    registerList.value = res.data.list
    total.value = res.data.total

  } catch (error) {
    ElMessage.error('加载报名记录失败')
  }

  loading.value = false
}


// ========================
// 切换竞赛
// ========================
const handleSelectCompetition = (id) => {
  activeCompetitionId.value = id
  currentPage.value = 1
  loadRegisters()
}


// ========================
// 分页
// ========================
const handlePageChange = (page) => {
  currentPage.value = page
  loadRegisters()
}


// ========================
// 打开审核弹窗
// ========================
const openReviewDialog = (id, status) => {
  reviewForm.value = {
    id,
    status,
    common: ''
  }
  reviewDialogVisible.value = true
}


// ========================
// 提交审核
// ========================
const submitReview = async () => {

  console.log('提交审核被点击')

  if (!reviewForm.value.common) {
    ElMessage.warning('请输入审核意见')
    return
  }

  try {
    console.log('开始调用接口')

    await competitionService.reviewCompetitionRegister(
      reviewForm.value.id,
      {
        Status: reviewForm.value.status,
        Common: reviewForm.value.common
      }
    )

    console.log('接口调用成功')

    ElMessage.success('审核成功')
    reviewDialogVisible.value = false
    loadRegisters()

  } catch (error) {
    console.error(error)
    ElMessage.error('审核失败')
  }
}


// ========================
// 状态文本
// ========================
const getStatusText = (status) => {
  if (status === 'approved') return '已通过'
  if (status === 'pending') return '待审核'
  return '已拒绝'
}


// ========================
// 时间格式
// ========================
const formatDate = (time) => {
  if (!time) return '-'
  return time.replace('T', ' ').substring(0, 16)
}


onMounted(loadCompetitions)
</script>

<style scoped>
.competition-review {
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
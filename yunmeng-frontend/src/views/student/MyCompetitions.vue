<template>
  <div class="my-competitions">

    <!-- 标题 -->
    <div class="page-header">
      <h2>我的竞赛</h2>
      <p>查看您已报名的竞赛及审核状态</p >
    </div>

    <!-- 卡片列表 -->
    <el-row :gutter="20" v-loading="loading" style="margin-top:20px">

      <el-col
        v-for="item in competitionList"
        :key="item.id"
        :xs="24"
        :sm="12"
        :md="8"
        style="margin-bottom:20px"
      >
        <el-card shadow="hover" class="competition-card">

          <!-- 状态 -->
          <div class="status-tag">
            <el-tag :type="getStatusType(item.status)">
              {{ getStatusText(item.status) }}
            </el-tag>
          </div>

          <!-- 标题 -->
          <h3>{{ item.competition?.title }}</h3>

          <!-- 队伍名称 -->
          <p><strong>队伍名称：</strong>{{ item.team_name }}</p >

          <!-- 报名时间 -->
          <p>
            <strong>报名时间：</strong>
            {{ formatDate(item.register_time) }}
          </p >

          <!-- 操作按钮 -->
          <div class="actions">
            <el-button size="small" @click="viewDetail(item)">
              查看详情
            </el-button>
          </div>

        </el-card>
      </el-col>

    </el-row>

    <!-- 空状态 -->
    <el-empty
      v-if="!loading && competitionList.length === 0"
      description="暂无报名记录"
    />

    <!-- 详情弹窗 -->
    <el-dialog v-model="detailVisible" title="报名详情" width="500px">
      <div v-if="selectedItem">

        <h3>{{ selectedItem.competition?.title }}</h3>

        <el-divider />

        <el-descriptions :column="1" border>

          <el-descriptions-item label="队伍名称">
            {{ selectedItem.team_name }}
          </el-descriptions-item>

          <el-descriptions-item label="审核状态">
            {{ getStatusText(selectedItem.status) }}
          </el-descriptions-item>

          <el-descriptions-item label="报名时间">
            {{ formatDate(selectedItem.register_time) }}
          </el-descriptions-item>

        </el-descriptions>

      </div>
    </el-dialog>

  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import CompetitionService from '@/services/CompetitionService'

export default {
  name: 'MyCompetitions',

  setup() {

    const loading = ref(false)
    const competitionList = ref([])

    const detailVisible = ref(false)
    const selectedItem = ref(null)

    // 获取我的报名记录
    const loadMyCompetitions = async () => {
      try {
        loading.value = true
        const res = await CompetitionService.getMyRegistrations()
        competitionList.value = res.data.list
      } catch (error) {
        ElMessage.error('获取报名记录失败')
      } finally {
        loading.value = false
      }
    }

    const formatDate = (time) => {
      if (!time) return '-'
      return time.replace('T', ' ').substring(0, 16)
    }

    const getStatusText = (status) => {
      const map = {
        pending: '待审核',
        approved: '已通过',
        rejected: '已拒绝'
      }
      return map[status] || status
    }

    const getStatusType = (status) => {
      const map = {
        pending: 'warning',
        approved: 'success',
        rejected: 'danger'
      }
      return map[status] || 'info'
    }

    const viewDetail = (item) => {
      selectedItem.value = item
      detailVisible.value = true
    }

    onMounted(loadMyCompetitions)

    return {
      loading,
      competitionList,
      detailVisible,
      selectedItem,
      formatDate,
      getStatusText,
      getStatusType,
      viewDetail
    }
  }
}
</script>

<style scoped>
.my-competitions {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
}

.competition-card {
  position: relative;
  border-radius: 12px;
  transition: 0.3s;
}

.competition-card:hover {
  transform: translateY(-5px);
}

.status-tag {
  position: absolute;
  top: 10px;
  right: 10px;
}

.actions {
  margin-top: 10px;
  text-align: right;
}
</style>
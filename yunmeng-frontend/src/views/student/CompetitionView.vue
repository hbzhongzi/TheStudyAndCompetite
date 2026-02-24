<template>
  <div class="competition-hall">

    <!-- 页面标题 -->
    <div class="page-header">
      <h2>竞赛大厅</h2>
      <p>浏览当前可参与的竞赛信息并进行报名</p >
    </div>

    <!-- 筛选区域 -->
    <div class="filter-section">
      <el-row :gutter="20">

        <el-col :span="6">
          <el-select v-model="filterLevel" placeholder="竞赛级别">
            <el-option label="全部" value="" />
            <el-option label="校级" value="school" />
            <el-option label="国家级" value="national" />
          </el-select>
        </el-col>

        <el-col :span="6">
          <el-select v-model="filterStatus" placeholder="竞赛状态">
            <el-option label="全部" value="" />
            <el-option label="报名中" value="registration" />
            <el-option label="已结束" value="finished" />
          </el-select>
        </el-col>

        <el-col :span="8">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索竞赛名称"
            clearable
          />
        </el-col>

      </el-row>
    </div>

    <!-- 卡片列表 -->
    <el-row :gutter="20" style="margin-top: 20px">

      <el-col
        v-for="item in filteredCompetitions"
        :key="item.id"
        :xs="24"
        :sm="12"
        :md="8"
        style="margin-bottom: 20px"
      >
        <el-card shadow="hover" class="competition-card">

          <!-- 状态角标 -->
          <div class="status-tag">
            <el-tag :type="getStatusType(item.status)">
              {{ getStatusText(item.status) }}
            </el-tag>
          </div>

          <!-- 标题 -->
          <h3 class="title">{{ item.title }}</h3>

          <!-- 分类 + 级别 -->
          <div class="meta">
            <el-tag size="small">{{ item.category }}</el-tag>
            <el-tag
              size="small"
              :type="item.level === 'national' ? 'danger' : 'success'"
            >
              {{ item.level === 'national' ? '国家级' : '校级' }}
            </el-tag>
          </div>

          <!-- 描述 -->
          <p class="description">
            {{ item.description }}
          </p >

          <!-- 报名时间 -->
          <div class="time">
            🗓 {{ formatDate(item.registrationStart) }}
            -
            {{ formatDate(item.registrationEnd) }}
          </div>

          <!-- 报名进度 -->
          <div class="progress-section">
            <span>报名人数：{{ item.currentParticipants }}/{{ item.maxParticipants }}</span>
            <el-progress
              :percentage="getProgress(item)"
              :stroke-width="8"
            />
          </div>

          <!-- 操作按钮 -->
          <div class="actions">
            <el-button size="small" @click="viewDetail(item)">
              详情
            </el-button>

            <el-button
              v-if="item.isOpen && item.status === 'registration'"
              type="primary"
              size="small"
              @click="registerCompetition(item)"
            >
              立即报名
            </el-button>
          </div>

        </el-card>
      </el-col>

    </el-row>

    <!-- 详情弹窗 -->
    <el-dialog v-model="detailVisible" title="竞赛详情" width="50%">
      <div v-if="selectedCompetition">
        <h3>{{ selectedCompetition.title }}</h3>
        <p>{{ selectedCompetition.description }}</p >

        <el-divider />

        <el-descriptions :column="1" border>
          <el-descriptions-item label="类别">
            {{ selectedCompetition.category }}
          </el-descriptions-item>
          <el-descriptions-item label="级别">
            {{ selectedCompetition.level }}
          </el-descriptions-item>
          <el-descriptions-item label="报名时间">
            {{ formatDate(selectedCompetition.registrationStart) }}
            -
            {{ formatDate(selectedCompetition.registrationEnd) }}
          </el-descriptions-item>
          <el-descriptions-item label="提交时间">
            {{ formatDate(selectedCompetition.submissionStart) }}
            -
            {{ formatDate(selectedCompetition.submissionEnd) }}
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </el-dialog>

  </div>
</template>
<script>
import { ref, computed, onMounted } from 'vue'
import competitionService from '@/services/competitionService'

export default {
  name: 'CompetitionHall',

  setup() {

    const competitions = ref([])
    const loading = ref(false)

    const filterLevel = ref('')
    const filterStatus = ref('')
    const searchKeyword = ref('')

    const detailVisible = ref(false)
    const selectedCompetition = ref(null)

    const loadCompetitions = async () => {
      loading.value = true
      const res = await competitionService.getCompetitions()
      competitions.value = res.data.list
      loading.value = false
    }

const filteredCompetitions = computed(() => {
  return competitions.value
    .filter(c => c.status !== 'draft') // ❌ 过滤草稿
    .filter(c => {
      return (
        (!filterLevel.value || c.level === filterLevel.value) &&
        (!filterStatus.value || c.status === filterStatus.value) &&
        (!searchKeyword.value || c.title.includes(searchKeyword.value))
      )
    })
})

    const formatDate = (time) => {
      if (!time) return '-'
      return time.replace('T', ' ').substring(0, 16)
    }

    const getStatusType = (status) => {
      const map = {
        draft: 'info',
        registration: 'success',
        ongoing: 'warning',
        finished: 'danger'
      }
      return map[status] || 'info'
    }

const getProgress = (item) => {
  if (!item.maxParticipants) return 0
  return Math.round((item.currentParticipants / item.maxParticipants) * 100)
}

    const getStatusText = (status) => {
      const map = {
        draft: '草稿',
        registration: '报名中',
        ongoing: '进行中',
        finished: '已结束'
      }
      return map[status] || status
    }

    const viewDetail = (row) => {
      selectedCompetition.value = row
      detailVisible.value = true
    }

    const registerCompetition = (row) => {
      console.log("报名竞赛:", row.id)
    }

    onMounted(loadCompetitions)

    return {
      loading,
      filterLevel,
      filterStatus,
      searchKeyword,
      filteredCompetitions,
      detailVisible,
      selectedCompetition,
      loadCompetitions,
      formatDate,
      getStatusType,
      getStatusText,
      viewDetail,
      getProgress,
      registerCompetition
    }
  }
}
</script>

<style scoped>
.competition-hall {
  padding: 20px;
}
.page-header {
  margin-bottom: 20px;
}
.filter-section {
  margin-bottom: 20px;
}

.competition-card {
  position: relative;
  border-radius: 12px;
  transition: all 0.3s;
}

.competition-card:hover {
  transform: translateY(-5px);
}

.status-tag {
  position: absolute;
  top: 10px;
  right: 10px;
}

.title {
  font-size: 18px;
  margin-bottom: 10px;
}

.meta {
  margin-bottom: 10px;
}

.description {
  font-size: 14px;
  color: #666;
  height: 40px;
  overflow: hidden;
}

.time {
  font-size: 13px;
  margin: 10px 0;
  color: #999;
}

.progress-section {
  margin: 10px 0;
}

.actions {
  display: flex;
  justify-content: space-between;
  margin-top: 10px;
}

</style>
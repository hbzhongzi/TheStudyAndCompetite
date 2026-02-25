<template>
  <div class="competition-guidance">

    <!-- 页面标题 -->
    <div class="page-header">
      <h2>我的竞赛任务</h2>
      <p>查看已分配的竞赛评审任务</p >
    </div>

    <!-- 任务卡片 -->
    <el-row :gutter="20" v-loading="loading" style="margin-top:20px">

      <el-col
        v-for="item in taskList"
        :key="item.id"
        :xs="24"
        :sm="12"
        :md="8"
        style="margin-bottom:20px"
      >
        <el-card shadow="hover" class="task-card">

          <div class="status-tag">
            <el-tag :type="getStatusType(item.status)">
              {{ getStatusText(item.status) }}
            </el-tag>
          </div>

          <h3>{{ item.competition?.title }}</h3>

          <p class="description">
            {{ item.competition?.description }}
          </p >

          <p class="time">
            <strong>分配时间：</strong>
            {{ formatDate(item.assigned_at) }}
          </p >

          <div class="actions">
            <el-button
              type="primary"
              size="small"
              @click="openDetail(item.competition_id)"
            >
              查看详情
            </el-button>
          </div>

        </el-card>
      </el-col>

    </el-row>

    <!-- 详情弹窗 -->
    <el-dialog
      v-model="detailVisible"
      title="竞赛详情"
      width="700px"
    >
      <el-descriptions
        v-if="competitionDetail"
        :column="2"
        border
      >

        <el-descriptions-item label="竞赛名称">
          {{ competitionDetail.title }}
        </el-descriptions-item>

        <el-descriptions-item label="竞赛级别">
          {{ getLevelText(competitionDetail.level) }}
        </el-descriptions-item>

        <el-descriptions-item label="竞赛类别">
          {{ competitionDetail.category }}
        </el-descriptions-item>

        <el-descriptions-item label="状态">
          {{ competitionDetail.status }}
        </el-descriptions-item>

        <el-descriptions-item label="报名时间">
          {{ formatDate(competitionDetail.registrationStart) }}
          -
          {{ formatDate(competitionDetail.registrationEnd) }}
        </el-descriptions-item>

        <el-descriptions-item label="提交时间">
          {{ formatDate(competitionDetail.submissionStart) }}
          -
          {{ formatDate(competitionDetail.submissionEnd) }}
        </el-descriptions-item>

        <el-descriptions-item label="人数上限">
          {{ competitionDetail.maxParticipants }}
        </el-descriptions-item>

        <el-descriptions-item label="当前人数">
          {{ competitionDetail.currentParticipants }}
        </el-descriptions-item>

        <el-descriptions-item label="是否开放">
          <el-tag :type="competitionDetail.isOpen ? 'success' : 'danger'">
            {{ competitionDetail.isOpen ? '开放' : '关闭' }}
          </el-tag>
        </el-descriptions-item>

        <el-descriptions-item label="创建时间">
          {{ formatDate(competitionDetail.createdAt) }}
        </el-descriptions-item>

      </el-descriptions>

      <!-- 奖项配置 -->
      <div v-if="competitionDetail?.awardConfig" style="margin-top:20px">
        <h4>奖项配置</h4>
        <el-tag type="success">
          一等奖 {{ competitionDetail.awardConfig.first_prize }} 名
        </el-tag>
        <el-tag type="warning" style="margin-left:10px">
          二等奖 {{ competitionDetail.awardConfig.second_prize }} 名
        </el-tag>
        <el-tag type="info" style="margin-left:10px">
          三等奖 {{ competitionDetail.awardConfig.third_prize }} 名
        </el-tag>
      </div>

    </el-dialog>

    <el-empty
      v-if="!loading && taskList.length === 0"
      description="暂无竞赛任务"
    />

  </div>
</template>

<script>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import teacherService from '@/services/teacherService'
import adminService from '@/services/adminService'

export default {
  name: 'TeacherTasks',

  setup() {

    const loading = ref(false)
    const taskList = ref([])

    const detailVisible = ref(false)
    const competitionDetail = ref(null)

    const loadTasks = async () => {
      try {
        loading.value = true
        const res = await teacherService.getMyTasks()
        taskList.value = res.data.list
      } catch (error) {
        ElMessage.error('获取任务失败')
      } finally {
        loading.value = false
      }
    }

    const openDetail = async (competitionId) => {
      try {
        const res = await adminService.getCompetitionDetail(competitionId)
        competitionDetail.value = res.data
        detailVisible.value = true
      } catch (error) {
        ElMessage.error('获取详情失败')
      }
    }

    const formatDate = (time) => {
      if (!time) return '-'
      return time.replace('T', ' ').substring(0, 16)
    }

    const getStatusText = (status) => {
      const map = {
        active: '进行中',
        completed: '已完成'
      }
      return map[status] || status
    }

    const getStatusType = (status) => {
      const map = {
        active: 'warning',
        completed: 'success'
      }
      return map[status] || 'info'
    }

    const getLevelText = (level) => {
      const map = {
        school: '校级',
        provincial: '省级',
        national: '国家级'
      }
      return map[level] || level
    }

    onMounted(loadTasks)

    return {
      loading,
      taskList,
      formatDate,
      getStatusText,
      getStatusType,
      openDetail,
      detailVisible,
      competitionDetail,
      getLevelText
    }
  }
}
</script>
<template>
  <div class="competition-tasks">

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

          <!-- 状态标签 -->
          <div class="status-tag">
            <el-tag :type="getStatusType(item.status)">
              {{ getStatusText(item.status) }}
            </el-tag>
          </div>

          <!-- 竞赛标题 -->
          <h3>{{ item.competition?.title }}</h3>

          <!-- 竞赛简介 -->
          <p class="description">
            {{ item.competition?.description }}
          </p >

          <!-- 分配时间 -->
          <p class="time">
            <strong>分配时间：</strong>
            {{ formatDate(item.assigned_at) }}
          </p >

          <!-- 操作按钮 -->
          <div class="actions">
            <el-button
              type="primary"
              size="small"
              @click="enterReview(item)"
            >
              进入评审
            </el-button>
          </div>

        </el-card>
      </el-col>

    </el-row>

    <!-- 空状态 -->
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
import { useRouter } from 'vue-router'

export default {
  name: 'TeacherTasks',

  setup() {

    const loading = ref(false)
    const taskList = ref([])
    const router = useRouter()

    const loadTasks = async () => {
      try {
        //loading.value = true
      // const res = await teacherService.getCompetitionRegistrations()
       // taskList.value = res.data.list
      } catch (error) {
        ElMessage.error('获取任务失败')
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

    const enterReview = (item) => {
      router.push(`/teacher/review/${item.competition_id}`)
    }

    onMounted(loadTasks)

    return {
      loading,
      taskList,
      formatDate,
      getStatusText,
      getStatusType,
      enterReview
    }
  }
}
</script>

<style scoped>
.teacher-tasks {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
}

.task-card {
  position: relative;
  border-radius: 12px;
  transition: 0.3s;
  min-height: 220px;
}

.task-card:hover {
  transform: translateY(-5px);
}

.status-tag {
  position: absolute;
  top: 12px;
  right: 12px;
}

.description {
  color: #666;
  margin: 10px 0;
  min-height: 60px;
}

.time {
  font-size: 13px;
  color: #888;
}

.actions {
  margin-top: 10px;
  text-align: right;
}
</style>
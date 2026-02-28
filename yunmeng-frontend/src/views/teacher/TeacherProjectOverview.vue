<template>
  <div class="project-overview">

    <h2>我的指导项目</h2>

    <!-- 统计卡片 -->
    <div class="stat-cards">
      <el-card class="stat-item">
        <div class="stat-title">项目总数</div>
        <div class="stat-value">{{ total }}</div>
      </el-card>

      <el-card class="stat-item">
        <div class="stat-title">已通过</div>
        <div class="stat-value success">{{ approvedCount }}</div>
      </el-card>

      <el-card class="stat-item">
        <div class="stat-title">已拒绝</div>
        <div class="stat-value danger">{{ rejectedCount }}</div>
      </el-card>

      <el-card class="stat-item">
        <div class="stat-title">待审核</div>
        <div class="stat-value warning">{{ pendingCount }}</div>
      </el-card>
    </div>

    <!-- 表格 -->
    <el-table
      :data="projectList"
      border
      stripe
      v-loading="loading"
    >

      <el-table-column prop="id" label="ID" width="80" />

      <el-table-column prop="title" label="项目名称" />

      <el-table-column prop="type" label="类型" width="120" />

      <el-table-column prop="studentName" label="学生" width="120" />

      <el-table-column label="状态" width="120">
        <template #default="{ row }">
          <el-tag
            :type="row.status === 'approved'
              ? 'success'
              : row.status === 'rejected'
              ? 'danger'
              : 'warning'"
          >
            {{ getStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column label="成员数" width="100">
        <template #default="{ row }">
          {{ row.memberCount }}
        </template>
      </el-table-column>

      <el-table-column label="文件数" width="100">
        <template #default="{ row }">
          {{ row.fileCount }}
        </template>
      </el-table-column>

      <el-table-column label="操作" width="120">
        <template #default="{ row }">
          <el-button
            size="small"
            type="primary"
            @click="openDetail(row)"
          >
            查看详情
          </el-button>
        </template>
      </el-table-column>

    </el-table>

    <el-pagination
      background
      layout="prev, pager, next"
      :total="total"
      :page-size="20"
      style="margin-top:20px;text-align:right"
      @current-change="handlePageChange"
    />

    <!-- 详情弹窗 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="项目详情"
      width="600px"
    >
      <el-descriptions :column="2" border>

        <el-descriptions-item label="项目名称">
          {{ currentProject.title }}
        </el-descriptions-item>

        <el-descriptions-item label="项目类型">
          {{ currentProject.type }}
        </el-descriptions-item>

        <el-descriptions-item label="学生">
          {{ currentProject.studentName }}
        </el-descriptions-item>

        <el-descriptions-item label="状态">
          {{ getStatusText(currentProject.status) }}
        </el-descriptions-item>

        <el-descriptions-item label="提交时间">
          {{ formatDate(currentProject.submittedAt) }}
        </el-descriptions-item>

        <el-descriptions-item label="创建时间">
          {{ formatDate(currentProject.createdAt) }}
        </el-descriptions-item>

        <el-descriptions-item label="成员数量">
          {{ currentProject.memberCount }}
        </el-descriptions-item>

        <el-descriptions-item label="文件数量">
          {{ currentProject.fileCount }}
        </el-descriptions-item>

        <el-descriptions-item label="评审次数">
          {{ currentProject.reviewCount }}
        </el-descriptions-item>

        <el-descriptions-item label="项目描述" :span="2">
          {{ currentProject.description }}
        </el-descriptions-item>

      </el-descriptions>

      <template #footer>
        <el-button @click="detailDialogVisible=false">
          关闭
        </el-button>
      </template>

    </el-dialog>

  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import project from '@/services/projectService'

const projectList = ref([])
const total = ref(0)
const loading = ref(false)
const currentPage = ref(1)

const detailDialogVisible = ref(false)
const currentProject = ref({})


// 加载项目
const loadProjects = async () => {
  loading.value = true

  try {
    const res = await project.getTeacherProjects({
      page: currentPage.value,
      size: 20
    })

    projectList.value = res.data.list
    total.value = res.data.total

  } catch (error) {
    ElMessage.error('加载项目失败')
  }

  loading.value = false
}


// 分页
const handlePageChange = (page) => {
  currentPage.value = page
  loadProjects()
}


// 打开详情
const openDetail = (row) => {
  currentProject.value = row
  detailDialogVisible.value = true
}


// 状态文本
const getStatusText = (status) => {
  if (status === 'approved') return '已通过'
  if (status === 'rejected') return '已拒绝'
  return '待审核'
}


// 时间格式
const formatDate = (time) => {
  if (!time) return '-'
  return time.replace('T', ' ').substring(0, 16)
}


// 统计
const approvedCount = computed(() =>
  projectList.value.filter(p => p.status === 'approved').length
)

const rejectedCount = computed(() =>
  projectList.value.filter(p => p.status === 'rejected').length
)

const pendingCount = computed(() =>
  projectList.value.filter(p => p.status !== 'approved' && p.status !== 'rejected').length
)


onMounted(loadProjects)
</script>

<style scoped>
.project-overview {
  background: #fff;
  padding: 20px;
  border-radius: 10px;
}

.stat-cards {
  display: flex;
  gap: 20px;
  margin-bottom: 20px;
}

.stat-item {
  flex: 1;
  text-align: center;
}

.stat-title {
  font-size: 14px;
  color: #666;
}

.stat-value {
  font-size: 26px;
  font-weight: bold;
  margin-top: 5px;
}

.success {
  color: #67c23a;
}

.danger {
  color: #f56c6c;
}

.warning {
  color: #e6a23c;
}
</style>
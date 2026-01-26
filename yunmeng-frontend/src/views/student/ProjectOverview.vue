<template>
  <div class="project-overview">

    <!-- 搜索与筛选 -->
    <el-card class="filter-card">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-input
            v-model="query.keyword"
            placeholder="搜索项目名称"
            clearable
            @input="reload"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-col>

        <el-col :span="4">
          <el-select v-model="query.status" placeholder="项目状态" clearable @change="reload">
            <el-option label="进行中" value="ongoing" />
            <el-option label="已完成" value="completed" />
            <el-option label="待审核" value="pending" />
            <el-option label="已拒绝" value="rejected" />
          </el-select>
        </el-col>

        <el-col :span="4">
          <el-select v-model="query.type" placeholder="项目类型" clearable @change="reload">
            <el-option label="科研项目" value="科研项目" />
            <el-option label="创新项目" value="创新项目" />
            <el-option label="竞赛项目" value="竞赛项目" />
          </el-select>
        </el-col>

        <el-col :span="6">
          <el-button type="primary" @click="createProject">创建项目</el-button>
          <el-button @click="reload">刷新</el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stat-row">
      <el-col :span="6" v-for="item in stats" :key="item.title">
        <el-card class="stat-card">
          <div class="stat-value">{{ item.value }}</div>
          <div class="stat-title">{{ item.title }}</div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 项目列表 -->
    <el-card>
      <template #header>
        <span>我的项目</span>
      </template>

      <el-table :data="projects" v-loading="loading">
        <el-table-column label="项目名称" min-width="200">
          <template #default="{ row }">
            <el-link type="primary" @click="openDetail(row)">
              {{ row.name }}
            </el-link>
          </template>
        </el-table-column>

        <el-table-column prop="type" label="类型" width="120" />

        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusTag(row.status)">
              {{ statusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="进度" width="160">
          <template #default="{ row }">
            <el-progress :percentage="row.progress" />
          </template>
        </el-table-column>
          <el-table-column prop="createdAt" label="创建时间" width="140">
              <template #default="scope">
                    <span>{{ formatDateTime(scope.row.createdAt) }}</span>
              </template>
          </el-table-column>
        <el-table-column label="操作" width="260" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="openDetail(row)">查看</el-button>
            <el-button size="small" type="success" @click="openProgress(row)">进度</el-button>
            <el-button size="small" type="danger" @click="remove(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="page.page"
          v-model:page-size="page.size"
          :total="page.total"
          layout="total, prev, pager, next"
          @current-change="reload"
        />
      </div>
    </el-card>

    <!-- 项目详情 -->
    <el-dialog v-model="detailVisible" title="项目详情" width="70%">
      <el-descriptions border :column="2">
        <el-descriptions-item label="项目名称">{{ current?.name }}</el-descriptions-item>
        <el-descriptions-item label="类型">{{ current?.type }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ statusText(current?.status) }}</el-descriptions-item>
        <el-descriptions-item label="进度">{{ current?.progress }}%</el-descriptions-item>
      </el-descriptions>

      <el-divider />
      <p><strong>项目描述：</strong>{{ current?.description }}</p >
      <p><strong>项目计划：</strong>{{ current?.plan }}</p >
    </el-dialog>

    <!-- 进度更新 -->
    <el-dialog v-model="progressVisible" title="更新进度" width="40%">
      <el-form>
        <el-form-item label="进度">
          <el-slider v-model="progressForm.progress" />
        </el-form-item>
        <el-form-item label="说明">
          <el-input v-model="progressForm.description" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="progressVisible=false">取消</el-button>
        <el-button type="primary" @click="submitProgress">提交</el-button>
      </template>
    </el-dialog>

  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import studentService from '../../services/studentService'

const loading = ref(false)
const projects = ref([])
const stats = ref([])
const current = ref(null)

const detailVisible = ref(false)
const progressVisible = ref(false)

const progressForm = ref({ progress: 0, description: '' })

const query = ref({ keyword: '', status: '', type: '' })
const page = ref({ page: 1, size: 10, total: 0 })

const reload = async () => {
  loading.value = true
  const res = await studentService.getMyProjects({
    ...query.value,
    page: page.value.page,
    size: page.value.size
  })
  projects.value = (res.data.list || []).map(p => ({
    ...p,
    name: p.title
  }))
  page.value.total = res.data.total
  loading.value = false
}

const loadStats = async () => {
  const res = await studentService.getProjectStats()
  const d = res.data
  stats.value = [
    { title: '项目总数', value: d.totalProjects },
    { title: '进行中', value: d.ongoingProjects },
    { title: '已完成', value: d.completedProjects },
    { title: '待审核', value: d.pendingProjects }
  ]
}


const formatDateTime = (dateString) => {
  if (!dateString) return '--'
  
  try {
    const date = new Date(dateString)
    
    // 格式：YYYY-MM-DD HH:mm:ss
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    const hours = String(date.getHours()).padStart(2, '0')
    const minutes = String(date.getMinutes()).padStart(2, '0')
    const seconds = String(date.getSeconds()).padStart(2, '0')
    
    return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`
  } catch (error) {
    console.error('日期格式化错误:', error)
    return dateString
  }
}

const openDetail = row => {
  current.value = row
  detailVisible.value = true
}

const openProgress = row => {
  current.value = row
  progressForm.value.progress = row.progress
  progressVisible.value = true
}

const submitProgress = async () => {
  await studentService.updateProjectProgress(current.value.id, progressForm.value)
  ElMessage.success('更新成功')
  progressVisible.value = false
  reload()
}

const remove = row => {
  ElMessageBox.confirm('确认删除该项目？', '提示', { type: 'warning' })
    .then(async () => {
      await studentService.deleteProject(row.id)
      ElMessage.success('删除成功')
      reload()
    })
}

const statusText = s =>
  ({ ongoing: '进行中', completed: '已完成', pending: '待审核', rejected: '已拒绝' }[s] || s)

const statusTag = s =>
  ({ ongoing: 'success', completed: 'success', pending: 'warning', rejected: 'danger' }[s] || '')

const createProject = () => {
  ElMessage.info('跳转到项目创建页面')
}

onMounted(() => {
  reload()
  loadStats()
})
</script>

<style scoped>
.project-overview {
  padding: 20px;
}
.filter-card {
  margin-bottom: 20px;
}
.stat-row {
  margin-bottom: 20px;
}
.stat-card {
  text-align: center;
}
.stat-value {
  font-size: 26px;
  font-weight: bold;
}
.stat-title {
  color: #666;
  margin-top: 6px;
}
.pagination {
  margin-top: 20px;
  text-align: right;
}
</style>
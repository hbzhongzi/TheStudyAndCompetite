<template>
  <div class="project-list">
    <!-- 顶部操作栏 -->
    <div class="action-bar">
      <div class="left-actions">
        <el-button type="primary" @click="showCreateDialog = true">
          <el-icon><Plus /></el-icon>
          新建项目
        </el-button>
        <el-select v-model="filterStatus" placeholder="筛选状态" @change="loadProjects" style="margin-left: 10px;">
          <el-option label="全部" value=""></el-option>
          <el-option label="草稿" value="draft"></el-option>
          <el-option label="待审核" value="pending"></el-option>
          <el-option label="已通过" value="approved"></el-option>
          <el-option label="已驳回" value="rejected"></el-option>
        </el-select>
      </div>
      <div class="right-actions">
        <el-button @click="loadProjects">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>
    </div>

    <!-- 项目列表 -->
    <el-table 
      :data="projects" 
      v-loading="loading"
      style="width: 100%"
      @row-click="handleRowClick"
    >
      <el-table-column prop="id" label="项目ID" width="80" />
      <el-table-column prop="title" label="项目标题" min-width="200" />
      <el-table-column prop="type" label="项目类型" width="100">
        <template #default="{ row }">
          <el-tag :type="row.type === '科研' ? 'primary' : 'success'">
            {{ row.type }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="120">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">
            {{ getStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="createdAt" label="创建时间" width="180">
        <template #default="{ row }">
          {{ formatDate(row.createdAt) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="250" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click.stop="viewProject(row)">查看</el-button>
          <el-button 
            v-if="row.status === 'draft'" 
            size="small" 
            type="primary" 
            @click.stop="editProject(row)"
          >
            编辑
          </el-button>
          <el-button 
            v-if="row.status === 'draft'" 
            size="small" 
            type="success" 
            @click.stop="submitForReview(row)"
          >
            提交审核
          </el-button>
          <el-button 
            v-if="row.status === 'rejected'" 
            size="small" 
            type="warning" 
            @click.stop="editProject(row)"
          >
            修改
          </el-button>
          <el-button 
            v-if="row.status === 'rejected'" 
            size="small" 
            type="danger" 
            @click.stop="deleteProject(row)"
          >
            删除
          </el-button>
          <el-button 
            v-if="row.status === 'draft'" 
            size="small" 
            type="danger" 
            @click.stop="deleteProject(row)"
          >
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- 空状态 -->
    <el-empty 
      v-if="!loading && projects.length === 0" 
      description="暂无项目数据"
      style="margin-top: 40px;"
    >
      <el-button type="primary" @click="showCreateDialog = true">创建第一个项目</el-button>
    </el-empty>

    <!-- 创建项目对话框 -->
    <ProjectForm 
      v-model:visible="showCreateDialog"
      @success="handleCreateSuccess"
    />

    <!-- 项目详情对话框 -->
    <ProjectDetail 
      v-model:visible="showDetailDialog"
      :project-id="selectedProjectId"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Refresh } from '@element-plus/icons-vue'
import { projectService } from '../services/projectService'
import ProjectForm from './ProjectForm.vue'
import ProjectDetail from './ProjectDetail.vue'

const router = useRouter()

// 响应式数据
const projects = ref([])
const loading = ref(false)
const filterStatus = ref('')
const showCreateDialog = ref(false)
const showDetailDialog = ref(false)
const selectedProjectId = ref(null)

// 获取项目列表
const loadProjects = async () => {
  loading.value = true
  try {
    const response = await projectService.getMyProjects(filterStatus.value)
    
    // 检查响应格式
    if (response && response.code === 200) {
      projects.value = response.data || []
    } else {
      projects.value = []
      ElMessage.warning('获取项目列表失败，请稍后重试')
    }
  } catch (error) {
    console.error('获取项目列表失败:', error)
    projects.value = []
    
    // 如果是网络错误，显示友好提示
    if (error.code === 'ECONNABORTED' || error.message.includes('Network Error')) {
      ElMessage.warning('网络连接失败，请检查网络后重试')
    } else {
      ElMessage.error(error.message || '获取项目列表失败')
    }
  } finally {
    loading.value = false
  }
}

// 状态类型映射
const getStatusType = (status) => {
  const statusMap = {
    draft: 'info',
    pending: 'warning',
    approved: 'success',
    rejected: 'danger'
  }
  return statusMap[status] || 'info'
}

// 状态文本映射
const getStatusText = (status) => {
  const statusMap = {
    draft: '草稿',
    pending: '待审核',
    approved: '已通过',
    rejected: '已驳回'
  }
  return statusMap[status] || status
}

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

// 查看项目详情
const viewProject = (project) => {
  selectedProjectId.value = project.id
  showDetailDialog.value = true
}

// 编辑项目
const editProject = (project) => {
  // 跳转到编辑页面，传递项目信息
  router.push({
    path: `/project/edit/${project.id}`,
    query: {
      fromStatus: project.status // 传递原始状态
    }
  })
}

// 提交审核
const submitForReview = async (project) => {
  try {
    await ElMessageBox.confirm(
      `确定要提交项目"${project.title}"进行审核吗？提交后将无法修改。`,
      '确认提交',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await projectService.updateProject(project.id, { ...project, status: 'pending' })
    ElMessage.success('项目已提交审核')
    loadProjects()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '提交审核失败')
    }
  }
}

// 删除项目
const deleteProject = async (project) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除项目"${project.title}"吗？此操作不可恢复。`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await projectService.deleteProject(project.id)
    ElMessage.success('删除成功')
    loadProjects()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(error.message)
    }
  }
}

// 行点击事件
const handleRowClick = (row) => {
  viewProject(row)
}

// 创建项目成功回调
const handleCreateSuccess = () => {
  showCreateDialog.value = false
  loadProjects()
  ElMessage.success('项目创建成功')
}

// 组件挂载时加载数据
onMounted(() => {
  loadProjects()
})
</script>

<style scoped>
.project-list {
  padding: 20px;
}

.action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 15px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.left-actions {
  display: flex;
  align-items: center;
}

.right-actions {
  display: flex;
  align-items: center;
}

.el-table {
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.el-table :deep(.el-table__row) {
  cursor: pointer;
}

.el-table :deep(.el-table__row:hover) {
  background-color: #f5f7fa;
}
</style> 
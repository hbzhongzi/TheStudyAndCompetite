<template>
  <div class="student-project-management">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h2>我的项目管理</h2>
        <p class="header-desc">管理您参与的所有项目，包括创建、编辑、查看进度等操作</p>
      </div>
      <div class="header-right">
        <el-button type="primary" @click="createNewProject">
          <i class="el-icon-plus"></i>
          创建新项目
        </el-button>
        <el-button type="success" @click="showProjectTemplates">
          <i class="el-icon-document"></i>
          项目模板
        </el-button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon draft">
              <i class="el-icon-document"></i>
            </div>
            <div class="stat-info">
              <h4>草稿项目</h4>
              <p class="stat-number">{{ stats.draftCount }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon pending">
              <i class="el-icon-time"></i>
            </div>
            <div class="stat-info">
              <h4>待审核</h4>
              <p class="stat-number">{{ stats.pendingCount }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon approved">
              <i class="el-icon-check"></i>
            </div>
            <div class="stat-info">
              <h4>进行中</h4>
              <p class="stat-number">{{ stats.approvedCount }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon completed">
              <i class="el-icon-trophy"></i>
            </div>
            <div class="stat-info">
              <h4>已完成</h4>
              <p class="stat-number">{{ stats.completedCount }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 高级筛选 -->
    <el-card class="filter-card">
      <template #header>
        <div class="filter-header">
          <span>项目筛选</span>
          <el-button link @click="resetFilters">重置筛选</el-button>
        </div>
      </template>
      
      <el-row :gutter="20">
        <el-col :span="6">
          <el-input
            v-model="filters.search"
            placeholder="搜索项目名称、描述"
            clearable
            @input="handleSearch"
          >
            <template #prefix>
              <i class="el-icon-search"></i>
            </template>
          </el-input>
        </el-col>
        <el-col :span="4">
          <el-select v-model="filters.status" placeholder="项目状态" clearable @change="handleSearch">
            <el-option label="全部状态" value="" />
            <el-option label="草稿" value="draft" />
            <el-option label="待审核" value="pending" />
            <el-option label="已通过" value="approved" />
            <el-option label="进行中" value="in_progress" />
            <el-option label="已完成" value="completed" />
            <el-option label="已驳回" value="rejected" />
            <el-option label="已暂停" value="suspended" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-select v-model="filters.type" placeholder="项目类型" clearable @change="handleSearch">
            <el-option label="全部类型" value="" />
            <el-option label="科研" value="科研" />
            <el-option label="竞赛" value="竞赛" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-select v-model="filters.role" placeholder="我的角色" clearable @change="handleSearch">
            <el-option label="全部角色" value="" />
            <el-option label="项目负责人" value="leader" />
            <el-option label="项目成员" value="member" />
            <el-option label="指导老师" value="teacher" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-select v-model="filters.level" placeholder="项目级别" clearable @change="handleSearch">
            <el-option label="全部级别" value="" />
            <el-option label="校级" value="校级" />
            <el-option label="省级" value="省级" />
            <el-option label="国家级" value="国家级" />
            <el-option label="国际级" value="国际级" />
          </el-select>
        </el-col>
        <el-col :span="2">
          <el-button type="primary" @click="handleSearch">搜索</el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- 项目列表 -->
    <el-card class="project-list-card">
      <template #header>
        <div class="list-header">
          <span>项目列表</span>
          <div class="list-actions">
            <el-button size="small" @click="refreshProjects">
              <i class="el-icon-refresh"></i>
              刷新
            </el-button>
            <el-button size="small" @click="exportMyProjects">
              <i class="el-icon-download"></i>
              导出
            </el-button>
          </div>
        </div>
      </template>

      <el-table 
        :data="filteredProjects" 
        style="width: 100%" 
        v-loading="loading"
        @row-click="viewProjectDetail"
        :row-class-name="getRowClassName"
      >
        <el-table-column prop="title" label="项目名称" min-width="200" show-overflow-tooltip>
          <template #default="{ row }">
            <div class="project-title">
              <span class="title-text">{{ row.title }}</span>
              <el-tag v-if="row.isExtended" size="small" type="warning">延期</el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="项目类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.type === '科研' ? 'primary' : 'success'">
              {{ row.type }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="level" label="项目级别" width="100">
          <template #default="{ row }">
            <el-tag :type="getLevelType(row.level)" size="small">
              {{ row.level }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="progress" label="进度" width="150">
          <template #default="{ row }">
            <el-progress 
              :percentage="row.progress || 0" 
              :status="getProgressStatus(row.progress)"
              :stroke-width="8"
            />
          </template>
        </el-table-column>
        <el-table-column prop="myRole" label="我的角色" width="100">
          <template #default="{ row }">
            <el-tag :type="getRoleType(row.myRole)" size="small">
              {{ getRoleText(row.myRole) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="expectedEndDate" label="预计完成" width="120">
          <template #default="{ row }">
            {{ formatDate(row.expectedEndDate) }}
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="120">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click.stop="viewProjectDetail(row)">
              <i class="el-icon-view"></i>
              查看
            </el-button>
            <el-button 
              v-if="canEdit(row)" 
              size="small" 
              type="primary" 
              @click.stop="editProject(row)"
            >
              <i class="el-icon-edit"></i>
              编辑
            </el-button>
            <el-button 
              v-if="canSubmit(row)" 
              size="small" 
              type="success" 
              @click.stop="submitProject(row)"
            >
              <i class="el-icon-upload"></i>
              提交
            </el-button>
            <el-button 
              v-if="canUpdateProgress(row)" 
              size="small" 
              type="warning" 
              @click.stop="updateProgress(row)"
            >
              <i class="el-icon-refresh"></i>
              进度
            </el-button>
            <el-dropdown @command="(command) => handleCommand(command, row)" @click.stop>
              <el-button size="small">
                更多<i class="el-icon-arrow-down el-icon--right"></i>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="files">管理附件</el-dropdown-item>
                  <el-dropdown-item command="members">管理成员</el-dropdown-item>
                  <el-dropdown-item command="timeline">项目时间线</el-dropdown-item>
                  <el-dropdown-item command="reviews">审核记录</el-dropdown-item>
                  <el-dropdown-item command="extend" v-if="canExtend(row)">申请延期</el-dropdown-item>
                  <el-dropdown-item command="archive" v-if="canArchive(row)" divided>归档项目</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="totalProjects"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 项目详情对话框 -->
    <el-dialog
      v-model="showDetailDialog"
      title="项目详情"
      width="80%"
      :close-on-click-modal="false"
    >
      <ProjectDetail 
        v-if="showDetailDialog"
        :project-id="selectedProjectId"
        :is-student="true"
        @refresh="loadProjects"
      />
    </el-dialog>

    <!-- 创建/编辑项目对话框 -->
    <el-dialog
      v-model="showProjectFormDialog"
      :title="isEditing ? '编辑项目' : '创建新项目'"
      width="70%"
      :close-on-click-modal="false"
    >
      <ProjectForm 
        v-if="showProjectFormDialog"
        :project="currentProject"
        :is-editing="isEditing"
        @submit="handleProjectSubmit"
        @cancel="showProjectFormDialog = false"
      />
    </el-dialog>

    <!-- 进度更新对话框 -->
    <el-dialog
      v-model="showProgressDialog"
      title="更新项目进度"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form :model="progressForm" label-width="100px" ref="progressFormRef">
        <el-form-item label="项目名称">
          <span>{{ currentProgressProject?.title }}</span>
        </el-form-item>
        <el-form-item label="当前进度" prop="progress">
          <el-slider
            v-model="progressForm.progress"
            :min="0"
            :max="100"
            :step="5"
            show-input
            input-size="small"
          />
        </el-form-item>
        <el-form-item label="进度说明" prop="description">
          <el-input
            v-model="progressForm.description"
            type="textarea"
            :rows="3"
            placeholder="请描述项目进展和下一步计划"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>
        <el-form-item label="预计完成" prop="expectedEndDate">
          <el-date-picker
            v-model="progressForm.expectedEndDate"
            type="date"
            placeholder="选择预计完成时间"
            :disabled-date="disabledDate"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showProgressDialog = false">取消</el-button>
        <el-button type="primary" @click="submitProgress" :loading="updatingProgress">
          提交进度
        </el-button>
      </template>
    </el-dialog>

    <!-- 项目模板对话框 -->
    <el-dialog
      v-model="showTemplatesDialog"
      title="项目模板库"
      width="80%"
      :close-on-click-modal="false"
    >
      <ProjectTemplates @select-template="handleTemplateSelect" />
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import ProjectDetail from '../../components/ProjectDetail.vue'
import ProjectForm from '../../components/ProjectForm.vue'
import ProjectTemplates from '../../components/ProjectTemplates.vue'
import { projectService } from '../../services/projectService'
import { studentService } from '../../services/studentService'

// 响应式数据
const loading = ref(false)
const updatingProgress = ref(false)
const currentPage = ref(1)
const pageSize = ref(20)
const totalProjects = ref(0)
const showDetailDialog = ref(false)
const showProjectFormDialog = ref(false)
const showProgressDialog = ref(false)
const showTemplatesDialog = ref(false)
const selectedProjectId = ref(null)
const currentProject = ref(null)
const currentProgressProject = ref(null)
const isEditing = ref(false)

// 统计数据
const stats = ref({
  draftCount: 0,
  pendingCount: 0,
  approvedCount: 0,
  completedCount: 0
})

// 筛选条件
const filters = reactive({
  search: '',
  status: '',
  type: '',
  role: '',
  level: ''
})

// 项目列表
const projects = ref([])

// 进度表单
const progressForm = reactive({
  progress: 0,
  description: '',
  expectedEndDate: null
})

// 表单引用
const progressFormRef = ref(null)

// 计算属性
const filteredProjects = computed(() => {
  return projects.value
})

const safeProjects = computed(() => {
  return Array.isArray(projects.value) ? projects.value : []
})

// 方法
const loadProjects = async () => {
  loading.value = true
  try {
    const response = await studentService.getMyProjects({
      page: currentPage.value,
      size: pageSize.value,
      ...filters
    })
    
    if (response && response.code === 200) {
      projects.value = response.data.list || []
      totalProjects.value = response.data.total || 0
    } else {
      projects.value = []
      totalProjects.value = 0
    }
  } catch (error) {
    console.error('加载项目列表失败:', error)
    ElMessage.error('加载项目列表失败')
    projects.value = []
    totalProjects.value = 0
  } finally {
    loading.value = false
  }
}

const loadStats = async () => {
  try {
    const response = await studentService.getProjectStats()
    if (response && response.code === 200) {
      stats.value = response.data
    }
  } catch (error) {
    console.error('加载统计数据失败:', error)
    stats.value = {
      draftCount: 0,
      pendingCount: 0,
      approvedCount: 0,
      completedCount: 0
    }
  }
}

const handleSearch = () => {
  currentPage.value = 1
  loadProjects()
}

const resetFilters = () => {
  Object.keys(filters).forEach(key => {
    filters[key] = ''
  })
  currentPage.value = 1
  loadProjects()
}

const handleSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
  loadProjects()
}

const handleCurrentChange = (page) => {
  currentPage.value = page
  loadProjects()
}

const refreshProjects = () => {
  loadProjects()
  loadStats()
}

const viewProjectDetail = (project) => {
  selectedProjectId.value = project.id
  showDetailDialog.value = true
}

const createNewProject = () => {
  currentProject.value = null
  isEditing.value = false
  showProjectFormDialog.value = true
}

const editProject = (project) => {
  currentProject.value = { ...project }
  isEditing.value = true
  showProjectFormDialog.value = true
}

const submitProject = async (project) => {
  try {
    await ElMessageBox.confirm(
      `确定要提交项目"${project.title}"吗？提交后将进入审核流程。`,
      '确认提交',
      { type: 'warning' }
    )
    
    await projectService.updateProject(project.id, { status: 'pending' })
    ElMessage.success('项目提交成功，等待审核')
    loadProjects()
    loadStats()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('项目提交失败')
    }
  }
}

const updateProgress = (project) => {
  currentProgressProject.value = project
  progressForm.progress = project.progress || 0
  progressForm.description = ''
  progressForm.expectedEndDate = project.expectedEndDate
  showProgressDialog.value = true
}

const submitProgress = async () => {
  if (!progressFormRef.value) return
  
  try {
    await progressFormRef.value.validate()
    
    updatingProgress.value = true
    await projectService.updateProjectProgress(currentProgressProject.value.id, {
      progress: progressForm.progress,
      description: progressForm.description,
      expectedEndDate: progressForm.expectedEndDate
    })
    
    ElMessage.success('进度更新成功')
    showProgressDialog.value = false
    loadProjects()
    loadStats()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('进度更新失败')
    }
  } finally {
    updatingProgress.value = false
  }
}

const showProjectTemplates = () => {
  showTemplatesDialog.value = true
}

const handleTemplateSelect = (template) => {
  currentProject.value = { ...template, id: null }
  isEditing.value = false
  showTemplatesDialog.value = false
  showProjectFormDialog.value = true
}

const handleProjectSubmit = async (projectData) => {
  try {
    if (isEditing.value) {
      await projectService.updateProject(currentProject.value.id, projectData)
      ElMessage.success('项目更新成功')
    } else {
      await projectService.createProject(projectData)
      ElMessage.success('项目创建成功')
    }
    
    showProjectFormDialog.value = false
    loadProjects()
    loadStats()
  } catch (error) {
    ElMessage.error(isEditing.value ? '项目更新失败' : '项目创建失败')
  }
}

const exportMyProjects = async () => {
  try {
    const response = await studentService.exportMyProjects(filters)
    if (response && response.data) {
      const blob = new Blob([response.data], { type: 'application/vnd.ms-excel' })
      const url = window.URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = url
      link.download = `我的项目_${new Date().toISOString().split('T')[0]}.xlsx`
      link.click()
      window.URL.revokeObjectURL(url)
      ElMessage.success('导出成功')
    }
  } catch (error) {
    console.error('导出失败:', error)
    ElMessage.error('导出失败')
  }
}

const handleCommand = (command, project) => {
  switch (command) {
    case 'files':
      ElMessage.info('管理附件功能开发中...')
      break
    case 'members':
      ElMessage.info('管理成员功能开发中...')
      break
    case 'timeline':
      ElMessage.info('项目时间线功能开发中...')
      break
    case 'reviews':
      ElMessage.info('审核记录功能开发中...')
      break
    case 'extend':
      ElMessage.info('申请延期功能开发中...')
      break
    case 'archive':
      ElMessage.info('归档项目功能开发中...')
      break
  }
}

// 权限检查方法
const canEdit = (project) => {
  return project.status === 'draft' && project.myRole === 'leader'
}

const canSubmit = (project) => {
  return project.status === 'draft' && project.myRole === 'leader'
}

const canUpdateProgress = (project) => {
  return ['approved', 'in_progress'].includes(project.status) && 
         ['leader', 'member'].includes(project.myRole)
}

const canExtend = (project) => {
  return ['approved', 'in_progress'].includes(project.status) && 
         project.myRole === 'leader'
}

const canArchive = (project) => {
  return project.status === 'completed' && project.myRole === 'leader'
}

// 工具方法
const getStatusType = (status) => {
  const statusMap = {
    draft: 'info',
    pending: 'warning',
    approved: 'success',
    in_progress: 'primary',
    completed: 'success',
    rejected: 'danger',
    suspended: 'warning'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status) => {
  const statusMap = {
    draft: '草稿',
    pending: '待审核',
    approved: '已通过',
    in_progress: '进行中',
    completed: '已完成',
    rejected: '已驳回',
    suspended: '已暂停'
  }
  return statusMap[status] || status
}

const getLevelType = (level) => {
  const levelMap = {
    '校级': 'info',
    '省级': 'success',
    '国家级': 'warning',
    '国际级': 'danger'
  }
  return levelMap[level] || 'info'
}

const getRoleType = (role) => {
  const roleMap = {
    leader: 'danger',
    member: 'primary',
    teacher: 'success'
  }
  return roleMap[role] || 'info'
}

const getRoleText = (role) => {
  const roleMap = {
    leader: '负责人',
    member: '成员',
    teacher: '指导老师'
  }
  return roleMap[role] || role
}

const getProgressStatus = (progress) => {
  if (progress >= 100) return 'success'
  if (progress >= 80) return 'warning'
  if (progress >= 50) return ''
  return 'exception'
}

const getRowClassName = ({ row }) => {
  if (row.status === 'pending') return 'pending-row'
  if (row.status === 'rejected') return 'rejected-row'
  if (row.isExtended) return 'extended-row'
  return ''
}

const formatDate = (dateString) => {
  if (!dateString) return '-'
  try {
    const date = new Date(dateString)
    return date.toLocaleDateString('zh-CN')
  } catch (error) {
    return dateString
  }
}

const disabledDate = (time) => {
  return time.getTime() < Date.now()
}

// 组件挂载和卸载
onMounted(() => {
  loadProjects()
  loadStats()
})

onUnmounted(() => {
  projects.value = []
})
</script>

<style scoped>
.student-project-management {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
}

.header-left h2 {
  margin: 0 0 8px 0;
  color: #2c3e50;
  font-size: 24px;
  font-weight: 600;
}

.header-desc {
  margin: 0;
  color: #7f8c8d;
  font-size: 14px;
}

.header-right {
  display: flex;
  gap: 10px;
}

.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 15px;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-icon.draft {
  background: linear-gradient(135deg, #3498db, #2980b9);
}

.stat-icon.pending {
  background: linear-gradient(135deg, #f39c12, #e67e22);
}

.stat-icon.approved {
  background: linear-gradient(135deg, #27ae60, #2ecc71);
}

.stat-icon.completed {
  background: linear-gradient(135deg, #9b59b6, #8e44ad);
}

.stat-icon i {
  font-size: 24px;
  color: white;
}

.stat-info h4 {
  margin: 0 0 5px 0;
  color: #7f8c8d;
  font-size: 14px;
}

.stat-number {
  margin: 0;
  font-size: 28px;
  font-weight: 600;
  color: #2c3e50;
}

.filter-card {
  margin-bottom: 20px;
  border-radius: 8px;
}

.filter-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.project-list-card {
  border-radius: 8px;
}

.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.list-actions {
  display: flex;
  gap: 10px;
}

.pagination-wrapper {
  margin-top: 20px;
  text-align: right;
}

.project-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.title-text {
  flex: 1;
}

.el-table {
  border-radius: 8px;
  overflow: hidden;
}

.el-table :deep(.el-table__row) {
  cursor: pointer;
}

.el-table :deep(.el-table__row:hover) {
  background-color: #f5f7fa;
}

.el-table :deep(.pending-row) {
  background-color: #fef9e7;
}

.el-table :deep(.rejected-row) {
  background-color: #fdf2f2;
}

.el-table :deep(.extended-row) {
  background-color: #f0f8ff;
}

.el-progress {
  margin-top: 4px;
}

.el-progress :deep(.el-progress__text) {
  font-size: 12px;
}
</style> 
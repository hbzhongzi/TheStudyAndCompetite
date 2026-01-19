<template>
  <div class="project-management">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h2>项目管理</h2>
        <p class="header-desc">管理系统内所有学生提交的项目，支持查看、审核、导出等操作</p>
      </div>
      <div class="header-right">
        <el-button type="success" @click="exportProjects">
          <i class="el-icon-download"></i>
          导出数据
        </el-button>
        <el-button type="primary" @click="showStatistics">
          <i class="el-icon-s-data"></i>
          统计报告
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
              <h4>已通过</h4>
              <p class="stat-number">{{ stats.approvedCount }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon total">
              <i class="el-icon-s-data"></i>
            </div>
            <div class="stat-info">
              <h4>总项目数</h4>
              <p class="stat-number">{{ stats.totalCount }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 高级筛选 -->
    <el-card class="filter-card">
      <template #header>
        <div class="filter-header">
          <span>高级筛选</span>
          <el-button link @click="resetFilters">重置筛选</el-button>
        </div>
      </template>
      
      <el-row :gutter="20">
        <el-col :span="6">
          <el-input
            v-model="filters.search"
            placeholder="搜索项目名称、学生姓名"
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
            <el-option label="已驳回" value="rejected" />
            <el-option label="已删除" value="deleted" />
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
          <el-select v-model="filters.department" placeholder="所属院系" clearable @change="handleSearch">
            <el-option label="全部院系" value="" />
            <el-option label="计算机学院" value="计算机学院" />
            <el-option label="机械学院" value="机械学院" />
            <el-option label="经管学院" value="经管学院" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-select v-model="filters.teacher" placeholder="指导老师" clearable @change="handleSearch">
            <el-option label="全部老师" value="" />
            <el-option label="张老师" value="张老师" />
            <el-option label="李老师" value="李老师" />
            <el-option label="王老师" value="王老师" />
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
            <el-button size="small" @click="batchApprove" :disabled="!hasSelectedProjects">
              批量通过
            </el-button>
            <el-button size="small" @click="batchReject" :disabled="!hasSelectedProjects">
              批量驳回
            </el-button>
            <el-button size="small" @click="cleanupProjects">
              清理无效项目
            </el-button>
          </div>
        </div>
      </template>

      <el-table 
        :data="filteredProjects" 
        style="width: 100%" 
        v-loading="loading"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="id" label="项目ID" width="80" />
        <el-table-column prop="title" label="项目名称" min-width="200" show-overflow-tooltip />
        <el-table-column prop="studentName" label="学生姓名" width="100" />
        <el-table-column prop="studentId" label="学号" width="120" />
        <el-table-column prop="department" label="院系" width="120" />
        <el-table-column prop="teacherName" label="指导老师" width="100" />
        <el-table-column prop="type" label="项目类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.type === '科研' ? 'primary' : 'success'">
              {{ row.type }}
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
        <el-table-column prop="memberCount" label="成员数" width="80" />
        <el-table-column prop="fileCount" label="附件数" width="80" />
        <el-table-column prop="createdAt" label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column prop="submittedAt" label="提交时间" width="160">
          <template #default="{ row }">
            {{ row.submittedAt ? formatDate(row.submittedAt) : '-' }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="viewProjectDetail(row)">查看详情</el-button>
            <el-button 
              v-if="row.status === 'pending'" 
              size="small" 
              type="success" 
              @click="reviewProject(row, 'approved')"
            >
              通过
            </el-button>
            <el-button 
              v-if="row.status === 'pending'" 
              size="small" 
              type="danger" 
              @click="reviewProject(row, 'rejected')"
            >
              驳回
            </el-button>
            <el-button 
              v-if="row.status === 'draft'" 
              size="small" 
              type="warning" 
              @click="forceSubmit(row)"
            >
              强制提交
            </el-button>
            <el-button 
              v-if="row.status === 'approved'" 
              size="small" 
              type="info" 
              @click="withdrawProject(row)"
            >
              撤回
            </el-button>
            <el-button 
              v-if="row.status === 'deleted'" 
              size="small" 
              type="success" 
              @click="restoreProject(row)"
            >
              恢复
            </el-button>
            <el-dropdown @command="(command) => handleCommand(command, row)">
              <el-button size="small">
                更多<i class="el-icon-arrow-down el-icon--right"></i>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="edit">编辑项目</el-dropdown-item>
                  <el-dropdown-item command="files">查看附件</el-dropdown-item>
                  <el-dropdown-item command="reviews">审核记录</el-dropdown-item>
                  <el-dropdown-item command="logs">操作日志</el-dropdown-item>
                  <el-dropdown-item command="force-status">强制更新状态</el-dropdown-item>
                  <el-dropdown-item command="soft-delete" divided>软删除</el-dropdown-item>
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
        :is-admin="true"
        @refresh="loadProjects"
      />
    </el-dialog>

    <!-- 审核对话框 -->
    <el-dialog
      v-model="showReviewDialog"
      title="项目审核"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form :model="reviewForm" label-width="100px">
        <el-form-item label="审核结果">
          <el-radio-group v-model="reviewForm.status">
            <el-radio label="approved">通过</el-radio>
            <el-radio label="rejected">驳回</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="审核意见">
          <el-input
            v-model="reviewForm.comments"
            type="textarea"
            :rows="4"
            placeholder="请输入审核意见"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showReviewDialog = false">取消</el-button>
        <el-button type="primary" @click="submitReview" :loading="reviewing">
          提交审核
        </el-button>
      </template>
    </el-dialog>

    <!-- 统计报告对话框 -->
    <el-dialog
      v-model="showStatsDialog"
      title="项目统计报告"
      width="70%"
      :close-on-click-modal="false"
    >
      <ProjectStatistics />
    </el-dialog>

    <!-- 批量操作确认对话框 -->
    <el-dialog
      v-model="showBatchDialog"
      title="批量操作确认"
      width="400px"
    >
      <p>{{ batchMessage }}</p>
      <template #footer>
        <el-button @click="showBatchDialog = false">取消</el-button>
        <el-button type="primary" @click="confirmBatchOperation" :loading="batchProcessing">
          确认
        </el-button>
      </template>
    </el-dialog>

    <!-- 强制更新状态对话框 -->
    <el-dialog
      v-model="showForceStatusDialog"
      title="强制更新项目状态"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form :model="forceStatusForm" label-width="100px">
        <el-form-item label="项目名称">
          <span>{{ currentForceProject?.title }}</span>
        </el-form-item>
        <el-form-item label="当前状态">
          <el-tag :type="getStatusType(currentForceProject?.status)">
            {{ getStatusText(currentForceProject?.status) }}
          </el-tag>
        </el-form-item>
        <el-form-item label="新状态">
          <el-select v-model="forceStatusForm.status" placeholder="选择新状态">
            <el-option label="草稿" value="draft" />
            <el-option label="待审核" value="pending" />
            <el-option label="已通过" value="approved" />
            <el-option label="已驳回" value="rejected" />
            <el-option label="已删除" value="deleted" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showForceStatusDialog = false">取消</el-button>
        <el-button type="primary" @click="submitForceStatus">
          确认更新
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import ProjectDetail from '../../components/ProjectDetail.vue'
import ProjectStatistics from './ProjectStatistics.vue'
import { projectService } from '../../services/projectService'
import { adminService } from '../../services/adminService'

// 响应式数据
const loading = ref(false)
const reviewing = ref(false)
const batchProcessing = ref(false)
const currentPage = ref(1)
const pageSize = ref(20)
const totalProjects = ref(0)
const selectedProjects = ref([])
const showDetailDialog = ref(false)
const showReviewDialog = ref(false)
const showStatsDialog = ref(false)
const showBatchDialog = ref(false)
const showForceStatusDialog = ref(false)
const selectedProjectId = ref(null)
const currentReviewProject = ref(null)
const currentForceProject = ref(null)
const batchMessage = ref('')
const batchOperation = ref('')
const forceStatusForm = reactive({
  status: 'pending'
})

// 统计数据
const stats = ref({
  draftCount: 0,
  pendingCount: 0,
  approvedCount: 0,
  totalCount: 0
})

// 筛选条件
const filters = reactive({
  search: '',
  status: '',
  type: '',
  department: '',
  teacher: ''
})

// 项目列表
const projects = ref([])

// 确保项目列表始终是一个数组
const safeProjects = computed(() => {
  return Array.isArray(projects.value) ? projects.value : []
})

// 审核表单
const reviewForm = reactive({
  status: 'approved',
  comments: ''
})

// 计算属性
const filteredProjects = computed(() => {
  return safeProjects.value
})

const hasSelectedProjects = computed(() => {
  return selectedProjects.value.length > 0
})

// 方法
const loadProjects = async () => {
  loading.value = true
  try {
    const response = await projectService.getProjectList({
      page: currentPage.value,
      size: pageSize.value,
      ...filters
    })
    
    if (response && response.code === 200) {
      projects.value = response.data.list || []
      totalProjects.value = response.data.total || 0
    }
  } catch (error) {
    console.error('加载项目列表失败:', error)
    ElMessage.error('加载项目列表失败')
  } finally {
    loading.value = false
  }
}

const loadStats = async () => {
  try {
    const response = await adminService.getProjectStats()
    if (response && response.code === 200) {
      stats.value = response.data
    }
  } catch (error) {
    console.error('加载统计数据失败:', error)
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
  handleSearch()
}

const handleSelectionChange = (selection) => {
  try {
    // 确保 selection 是一个数组
    if (Array.isArray(selection)) {
      selectedProjects.value = selection
    } else {
      selectedProjects.value = []
    }
  } catch (error) {
    console.error('选择变化处理错误:', error)
    selectedProjects.value = []
  }
}

const handleSizeChange = (size) => {
  pageSize.value = size
  loadProjects()
}

const handleCurrentChange = (page) => {
  currentPage.value = page
  loadProjects()
}

const viewProjectDetail = (project) => {
  selectedProjectId.value = project.id
  showDetailDialog.value = true
}

const reviewProject = (project, status) => {
  currentReviewProject.value = project
  reviewForm.status = status
  reviewForm.comments = ''
  showReviewDialog.value = true
}

const submitReview = async () => {
  if (!reviewForm.comments.trim()) {
    ElMessage.warning('请输入审核意见')
    return
  }

  reviewing.value = true
  try {
    await projectService.reviewProject(currentReviewProject.value.id, {
      status: reviewForm.status,
      comments: reviewForm.comments
    })
    
    ElMessage.success('审核提交成功')
    showReviewDialog.value = false
    loadProjects()
    loadStats()
  } catch (error) {
    ElMessage.error('审核提交失败')
  } finally {
    reviewing.value = false
  }
}

const forceSubmit = async (project) => {
  try {
    await ElMessageBox.confirm(
      `确定要强制提交项目"${project.title}"吗？`,
      '确认操作',
      { type: 'warning' }
    )
    
    await projectService.updateProject(project.id, { status: 'pending' })
    ElMessage.success('项目已强制提交')
    loadProjects()
    loadStats()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('操作失败')
    }
  }
}

const withdrawProject = async (project) => {
  try {
    await ElMessageBox.confirm(
      `确定要撤回项目"${project.title}"吗？`,
      '确认操作',
      { type: 'warning' }
    )
    
    await projectService.updateProject(project.id, { status: 'draft' })
    ElMessage.success('项目已撤回')
    loadProjects()
    loadStats()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('操作失败')
    }
  }
}

const openForceStatusDialog = (project) => {
  currentForceProject.value = project
  forceStatusForm.status = project.status
  showForceStatusDialog.value = true
}

const submitForceStatus = async () => {
  try {
    await adminService.forceUpdateProjectStatus(currentForceProject.value.id, forceStatusForm.status)
    ElMessage.success('项目状态强制更新成功')
    showForceStatusDialog.value = false
    loadProjects()
    loadStats()
  } catch (error) {
    ElMessage.error('强制更新状态失败')
  }
}

const softDeleteProject = async (project) => {
  try {
    await ElMessageBox.confirm(
      `确定要软删除项目"${project.title}"吗？`,
      '确认操作',
      { type: 'warning' }
    )
    
    await adminService.softDeleteProject(project.id)
    ElMessage.success('项目已软删除')
    loadProjects()
    loadStats()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('软删除失败')
    }
  }
}

const restoreProject = async (project) => {
  try {
    await ElMessageBox.confirm(
      `确定要恢复项目"${project.title}"吗？`,
      '确认操作',
      { type: 'warning' }
    )
    
    await adminService.restoreProject(project.id)
    ElMessage.success('项目已恢复')
    loadProjects()
    loadStats()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('恢复项目失败')
    }
  }
}

const batchApprove = () => {
  if (selectedProjects.value.length === 0) {
    ElMessage.warning('请先选择要操作的项目')
    return
  }
  batchMessage.value = `确定要批量通过选中的 ${selectedProjects.value.length} 个项目吗？`
  batchOperation.value = 'approve'
  showBatchDialog.value = true
}

const batchReject = () => {
  if (selectedProjects.value.length === 0) {
    ElMessage.warning('请先选择要操作的项目')
    return
  }
  batchMessage.value = `确定要批量驳回选中的 ${selectedProjects.value.length} 个项目吗？`
  batchOperation.value = 'reject'
  showBatchDialog.value = true
}

const confirmBatchOperation = async () => {
  batchProcessing.value = true
  try {
    const projectIds = selectedProjects.value.map(p => p.id)
    const status = batchOperation.value === 'approve' ? 'approved' : 'rejected'
    
    await projectService.batchReviewProjects(projectIds, {
      status,
      comments: `管理员批量${batchOperation.value === 'approve' ? '通过' : '驳回'}`
    })
    
    ElMessage.success('批量操作成功')
    showBatchDialog.value = false
    selectedProjects.value = []
    loadProjects()
    loadStats()
  } catch (error) {
    console.error('批量操作失败:', error)
    ElMessage.error('批量操作失败')
  } finally {
    batchProcessing.value = false
  }
}

const cleanupProjects = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要清理长时间未提交的草稿项目吗？此操作不可恢复。',
      '确认清理',
      { type: 'warning' }
    )
    
    await adminService.cleanupProjects()
    ElMessage.success('清理完成')
    loadProjects()
    loadStats()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('清理失败')
    }
  }
}

const exportProjects = async () => {
  try {
    const response = await adminService.exportProjects(filters)
    if (response && response.data) {
      const blob = new Blob([response.data], { type: 'application/vnd.ms-excel' })
      const url = window.URL.createObjectURL(blob)
      const link = document.createElement('a')
      link.href = url
      link.download = `项目列表_${new Date().toISOString().split('T')[0]}.xlsx`
      link.click()
      window.URL.revokeObjectURL(url)
      ElMessage.success('导出成功')
    }
  } catch (error) {
    console.error('导出失败:', error)
    ElMessage.error('导出失败')
  }
}

const showStatistics = () => {
  showStatsDialog.value = true
}

const handleCommand = (command, project) => {
  switch (command) {
    case 'edit':
      editProject(project)
      break
    case 'files':
      viewProjectFiles(project)
      break
    case 'reviews':
      viewProjectReviews(project)
      break
    case 'logs':
      viewProjectLogs(project)
      break
    case 'force-status':
      openForceStatusDialog(project)
      break
    case 'soft-delete':
      softDeleteProject(project)
      break
  }
}

// 编辑项目
const editProject = (project) => {
  ElMessage.info('编辑项目功能开发中...')
}

// 查看项目附件
const viewProjectFiles = (project) => {
  ElMessage.info('查看附件功能开发中...')
}

// 查看审核记录
const viewProjectReviews = (project) => {
  ElMessage.info('查看审核记录功能开发中...')
}

// 查看操作日志
const viewProjectLogs = (project) => {
  ElMessage.info('查看操作日志功能开发中...')
}

const getStatusType = (status) => {
  const statusMap = {
    draft: 'info',
    pending: 'warning',
    approved: 'success',
    rejected: 'danger',
    deleted: 'danger'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status) => {
  const statusMap = {
    draft: '草稿',
    pending: '待审核',
    approved: '已通过',
    rejected: '已驳回',
    deleted: '已删除'
  }
  return statusMap[status] || status
}

const formatDate = (dateString) => {
  if (!dateString) return '-'
  try {
    const date = new Date(dateString)
    return date.toLocaleString('zh-CN')
  } catch (error) {
    return dateString
  }
}

// 组件挂载时加载数据
onMounted(() => {
  loadProjects()
  loadStats()
})

onUnmounted(() => {
  // 清理资源
  selectedProjects.value = []
  projects.value = []
})
</script>

<style scoped>
.project-management {
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

.stat-icon.total {
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
</style> 
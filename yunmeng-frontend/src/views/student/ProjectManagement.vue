<template>
  <div class="student-project-management">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h2>我的项目管理</h2>
        <p class="header-desc">管理您参与的所有项目，包括查看、跟踪进度等操作</p>
      </div>
      <div class="header-right">
        <el-button type="primary" @click="createNewProject">
          <el-icon><Plus /></el-icon>
          创建新项目
        </el-button>
        <el-button type="success" @click="showProjectTemplates">
          <el-icon><Document /></el-icon>
          项目模板
        </el-button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="24" :sm="12" :md="6" :lg="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon draft">
              <el-icon><Document /></el-icon>
            </div>
            <div class="stat-info">
              <h4>草稿项目</h4>
              <p class="stat-number">{{ stats.draftCount || 0 }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6" :lg="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon pending">
              <el-icon><Clock /></el-icon>
            </div>
            <div class="stat-info">
              <h4>待审核</h4>
              <p class="stat-number">{{ stats.pendingCount || 0 }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6" :lg="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon approved">
              <el-icon><Check /></el-icon>
            </div>
            <div class="stat-info">
              <h4>进行中</h4>
              <p class="stat-number">{{ stats.approvedCount || 0 }}</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6" :lg="6">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon completed">
              <el-icon><Trophy /></el-icon>
            </div>
            <div class="stat-info">
              <h4>已完成</h4>
              <p class="stat-number">{{ stats.completedCount || 0 }}</p>
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
        <el-col :xs="24" :sm="12" :md="8" :lg="6">
          <el-input
            v-model="filters.search"
            placeholder="搜索项目名称"
            clearable
            @input="handleSearch"
            @clear="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-col>
        <el-col :xs="24" :sm="12" :md="6" :lg="4">
          <el-select 
            v-model="filters.status" 
            placeholder="项目状态" 
            clearable 
            filterable
            @change="handleSearch"
          >
            <el-option label="全部状态" value="" />
            <el-option label="待审核" value="pending" />
            <el-option label="已通过" value="approved" />
            <el-option label="进行中" value="in_progress" />
            <el-option label="已完成" value="completed" />
            <el-option label="已驳回" value="rejected" />
          </el-select>
        </el-col>
        <el-col :xs="24" :sm="12" :md="6" :lg="4">
          <el-select 
            v-model="filters.type" 
            placeholder="项目类型" 
            clearable 
            filterable
            @change="handleSearch"
          >
            <el-option label="全部类型" value="" />
            <el-option label="科研项目" value="科研项目" />
            <el-option label="创新项目" value="创新项目" />
            <el-option label="竞赛项目" value="竞赛项目" />
            <el-option label="软件开发" value="软件开发" />
          </el-select>
        </el-col>
        <el-col :xs="24" :sm="12" :md="6" :lg="4">
          <el-select 
            v-model="filters.level" 
            placeholder="项目级别" 
            clearable 
            filterable
            @change="handleSearch"
          >
            <el-option label="全部级别" value="" />
            <el-option label="院级" value="院级" />
            <el-option label="校级" value="校级" />
            <el-option label="省级" value="省级" />
            <el-option label="国家级" value="国家级" />
          </el-select>
        </el-col>
        <el-col :xs="24" :sm="24" :md="6" :lg="2">
          <el-button type="primary" @click="handleSearch" :loading="loading">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
        </el-col>
      </el-row>
    </el-card>

<!-- 项目列表 -->
<el-card class="project-list-card">
  <template #header>
    <div class="list-header">
      <span>项目列表 ({{ totalProjects }})</span>
      <div class="list-actions">
        <el-button size="small" @click="refreshProjects" :loading="loading">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
        <el-button size="small" @click="exportMyProjects">
          <el-icon><Download /></el-icon>
          导出
        </el-button>
      </div>
    </div>
  </template>

  <el-table 
    :data="projects" 
    style="width: 100%" 
    v-loading="loading"
    stripe
    border
    @row-click="viewProjectDetail"
  >
    <el-table-column prop="id" label="ID" width="80" align="center" sortable>
      <template #default="{ row }">
        <span class="project-id">#{{ row.id }}</span>
      </template>
    </el-table-column>
    
    <el-table-column prop="title" label="项目名称" min-width="220" show-overflow-tooltip>
      <template #default="{ row }">
        <div class="project-title">
          <span class="title-text">{{ row.title }}</span>
          <el-tag v-if="row.status === 'rejected'" size="small" type="danger">已驳回</el-tag>
          <el-tag v-else-if="row.status === 'pending'" size="small" type="warning">待审核</el-tag>
          <el-tag v-else-if="row.status === 'approved'" size="small" type="success">已通过</el-tag>
        </div>
        <div v-if="row.description" class="project-description">
          {{ row.description }}
        </div>
      </template>
    </el-table-column>
    
    <el-table-column prop="type" label="项目类型" width="100" align="center">
      <template #default="{ row }">
        <el-tag :type="getTypeTagType(row.type)" size="small">
          {{ row.type || '未设置' }}
        </el-tag>
      </template>
    </el-table-column>
    
    <el-table-column prop="status" label="状态" width="100" align="center">
      <template #default="{ row }">
        <el-tag :type="getStatusTagType(row.status)" size="small">
          {{ getStatusText(row.status) }}
        </el-tag>
      </template>
    </el-table-column>
    
    <el-table-column prop="progress" label="进度" width="130" align="center">
      <template #default="{ row }">
        <div class="progress-container">
          <el-progress 
            v-if="row.progress !== undefined && row.progress !== null"
            :percentage="row.progress" 
            :status="getProgressStatus(row.progress)"
            :stroke-width="10"
            :show-text="true"
            :text-inside="true"
          />
          <span v-else class="text-gray">-</span>
        </div>
      </template>
    </el-table-column>
    
    <el-table-column prop="level" label="项目级别" width="100" align="center">
      <template #default="{ row }">
        <el-tag v-if="row.level" :type="getLevelTagType(row.level)" size="small">
          {{ row.level }}
        </el-tag>
        <span v-else class="text-gray">-</span>
      </template>
    </el-table-column>
    
    <el-table-column prop="studentName" label="负责人" width="100" align="center">
      <template #default="{ row }">
        <span>{{ row.studentName || '未设置' }}</span>
      </template>
    </el-table-column>
    
    <el-table-column prop="teacherName" label="指导老师" width="100" align="center">
      <template #default="{ row }">
        <span>{{ row.teacherName || '未设置' }}</span>
      </template>
    </el-table-column>
    
    <el-table-column prop="createdAt" label="创建时间" width="120" align="center" sortable>
      <template #default="{ row }">
        <div class="date-cell">
          {{ formatDate(row.createdAt) }}
        </div>
      </template>
    </el-table-column>
    
    <el-table-column prop="updatedAt" label="更新时间" width="120" align="center" sortable>
      <template #default="{ row }">
        <div class="date-cell">
          {{ formatDate(row.updatedAt) }}
        </div>
      </template>
    </el-table-column>
    
    <el-table-column label="操作" width="220" fixed="right" align="center">
      <template #default="{ row }">
        <el-button-group>
          <el-button size="small" @click.stop="viewProjectDetail(row)" type="info">
            <el-icon><View /></el-icon>
            查看
          </el-button>
          
          <el-button 
            v-if="canEdit(row)" 
            size="small" 
            type="primary" 
            @click.stop="editProject(row)"
          >
            <el-icon><Edit /></el-icon>
            编辑
          </el-button>
          
          <el-dropdown @command="(command) => handleCommand(command, row)" @click.stop>
            <el-button size="small">
              <el-icon><More /></el-icon>
              更多
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item 
                  command="submit" 
                  v-if="canSubmit(row)"
                  :icon="Upload"
                >
                  提交审核
                </el-dropdown-item>
                
                <el-dropdown-item 
                  command="progress" 
                  v-if="canUpdateProgress(row)"
                  :icon="TrendCharts"
                >
                  更新进度
                </el-dropdown-item>
                
                <el-dropdown-item 
                  command="files" 
                  :icon="Folder"
                  divided
                >
                  附件管理
                </el-dropdown-item>
                
                <el-dropdown-item 
                  command="members" 
                  :icon="User"
                >
                  成员管理
                </el-dropdown-item>
                
                <el-dropdown-item 
                  command="timeline" 
                  :icon="Clock"
                >
                  时间线
                </el-dropdown-item>
                
                <el-dropdown-item 
                  command="extend" 
                  v-if="canExtend(row)"
                  :icon="Calendar"
                  divided
                >
                  申请延期
                </el-dropdown-item>
                
                <el-dropdown-item 
                  command="report" 
                  v-if="canGenerateReport(row)"
                  :icon="Document"
                >
                  生成报告
                </el-dropdown-item>
                
                <el-dropdown-item 
                  command="archive" 
                  v-if="canArchive(row)"
                  :icon="Box"
                  divided
                >
                  归档项目
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </el-button-group>
      </template>
    </el-table-column>
  </el-table>


      <!-- 分页 -->
      <div class="pagination-wrapper" v-if="totalProjects > 0">
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

      <!-- 空状态 -->
      <div class="empty-state" v-if="!loading && projects.length === 0">
        <el-empty description="暂无项目数据">
          <template #image>
            <el-icon size="60"><Document /></el-icon>
          </template>
          <el-button type="primary" @click="createNewProject">创建第一个项目</el-button>
        </el-empty>
      </div>
    </el-card>

    <!-- 项目详情对话框 -->
    <el-dialog
      v-model="showDetailDialog"
      :title="selectedProject?.title || '项目详情'"
      width="80%"
      :close-on-click-modal="false"
      :destroy-on-close="true"
    >
      <ProjectDetail 
        v-if="showDetailDialog"
        :project="selectedProject"
        :is-student="true"
        @refresh="refreshProjects"
      />
    </el-dialog>

    <!-- 创建/编辑项目对话框 -->
    <el-dialog
      v-model="showProjectFormDialog"
      :title="isEditing ? '编辑项目' : '创建新项目'"
      width="70%"
      :close-on-click-modal="false"
      :destroy-on-close="true"
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
      :title="`更新项目进度 - ${currentProgressProject?.title}`"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form :model="progressForm" label-width="100px">
        <el-form-item label="当前进度">
          <el-slider
            v-model="progressForm.progress"
            :min="0"
            :max="100"
            :step="5"
            show-input
            input-size="small"
          />
        </el-form-item>
        <el-form-item label="进度说明">
          <el-input
            v-model="progressForm.description"
            type="textarea"
            :rows="3"
            placeholder="请描述项目进展、遇到的问题和下一步计划"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>
        <el-form-item label="预计完成">
          <el-date-picker
            v-model="progressForm.expectedEndDate"
            type="date"
            placeholder="选择预计完成时间"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
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
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import ProjectDetail from '../../components/ProjectDetail.vue'
import ProjectForm from '../../components/ProjectForm.vue'
import ProjectTemplates from '../../components/ProjectTemplates.vue'
import  projectService from '@/services/projectService'
import { studentService } from '../../services/studentService'

// 响应式数据
const loading = ref(false)
const updatingProgress = ref(false)
const currentPage = ref(1)
const pageSize = ref(10) // 根据接口返回默认调整为10
const totalProjects = ref(0)
const showDetailDialog = ref(false)
const showProjectFormDialog = ref(false)
const showProgressDialog = ref(false)
const showTemplatesDialog = ref(false)
const selectedProject = ref(null)
const currentProject = ref(null)
const currentProgressProject = ref(null)
const isEditing = ref(false)

// 工具方法
const getStatusTagType = (status) => {
  const statusMap = {
    draft: 'info',          // 草稿
    pending: 'warning',     // 待审核
    approved: 'success',    // 已通过
    in_progress: 'primary', // 进行中
    completed: 'success',   // 已完成
    rejected: 'danger',     // 已驳回
    suspended: 'warning'    // 已暂停
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

const getTypeTagType = (type) => {
  const typeMap = {
    '科研项目': 'primary',
    '创新项目': 'success',
    '竞赛项目': 'warning',
    '软件开发': 'info',
    '毕业论文': 'danger',
    '课程设计': 'primary'
  }
  return typeMap[type] || 'info'
}

const getLevelTagType = (level) => {
  const levelMap = {
    '院级': 'info',
    '校级': 'success',
    '省级': 'warning',
    '国家级': 'danger',
    '国际级': 'danger'
  }
  return levelMap[level] || 'info'
}

const getProgressStatus = (progress) => {
  if (progress >= 100) return 'success'
  if (progress >= 80) return 'warning'
  if (progress >= 50) return ''
  return 'exception'
}

const formatDate = (dateString) => {
  if (!dateString) return '-'
  try {
    const date = new Date(dateString)
    // 格式化为 YYYY-MM-DD HH:mm
    const year = date.getFullYear()
    const month = String(date.getMonth() + 1).padStart(2, '0')
    const day = String(date.getDate()).padStart(2, '0')
    const hours = String(date.getHours()).padStart(2, '0')
    const minutes = String(date.getMinutes()).padStart(2, '0')
    return `${year}-${month}-${day} ${hours}:${minutes}`
  } catch (error) {
    return dateString
  }
}

// 权限检查方法（根据实际业务逻辑调整）
const canEdit = (project) => {
  return project.status === 'draft'
}

const canSubmit = (project) => {
  return project.status === 'draft'
}

const canUpdateProgress = (project) => {
  return ['approved', 'in_progress'].includes(project.status)
}

const canExtend = (project) => {
  return ['approved', 'in_progress'].includes(project.status)
}

const canGenerateReport = (project) => {
  return ['completed', 'approved', 'in_progress'].includes(project.status)
}

const canArchive = (project) => {
  return project.status === 'completed'
}

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

// 加载项目列表
const loadProjects = async () => {
  loading.value = true
  try {
    console.log('📡 请求项目数据，参数:', {
      page: currentPage.value,
      size: pageSize.value,
      ...filters
    })
    
    const response = await studentService.getMyProjects({
      page: currentPage.value,
      size: pageSize.value,
      ...filters
    })
    
    console.log('✅ 项目数据响应:', response)
    
    if (response && response.code === 200) {
      // 处理接口返回的数据结构
      let projectList = []
      let total = 0
      
      if (response.data && response.data.list) {
        // 如果返回的是接口格式
        projectList = response.data.list || []
        total = response.data.total || 0
        currentPage.value = response.data.page || 1
        pageSize.value = response.data.size || 10
      } else if (Array.isArray(response.data)) {
        // 如果返回的是模拟数据格式
        projectList = response.data
        total = projectList.length
      } else if (Array.isArray(response)) {
        // 如果直接返回数组
        projectList = response
        total = projectList.length
      }
      
      // 映射字段：将接口返回的字段映射到组件使用的字段
      projects.value = projectList.map(item => ({
        // 主键
        id: item.id,
        
        // 项目基本信息
        title: item.title || item.name || '未命名项目',
        description: item.description || '',
        type: item.type || '',
        
        // 项目状态和进度
        status: item.status || '',
        progress: item.progress || 0,
        
        // 人员信息（模拟数据，实际应从接口获取）
        studentName: item.studentName || '张三',
        teacherName: item.teacherName || '李老师',
        
        // 级别信息
        level: item.level || '院级',
        
        // 时间信息
        createdAt: item.createdAt || item.createTime || '',
        updatedAt: item.updated_at || item.updateTime || '',
        deadline: item.deadline || '',
        
        // 计划信息
        plan: item.plan || '',
        
        // 扩展信息
        expectedEndDate: item.expectedEndDate || '',
        isExtended: false
      }))
      
      totalProjects.value = total
      
      // 更新统计数据
      updateProjectStats(projectList)
      
      console.log(`✅ 成功加载 ${projects.value.length} 个项目，总计 ${totalProjects.value} 个`)
      
    } else {
      console.warn('⚠️ 响应格式异常:', response)
      projects.value = []
      totalProjects.value = 0
      resetStats()
      ElMessage.warning(response?.message || '获取项目数据失败')
    }
  } catch (error) {
    console.error('❌ 加载项目列表失败:', error)
    // 使用模拟数据作为后备
    projects.value = [
      {
        id: 1,
        title: '智能校园系统',
        type: '软件开发',
        status: 'in_progress',
        progress: 75,
        createdAt: '2024-01-15',
        studentName: '张三',
        teacherName: '李老师',
        level: '校级',
        description: '基于物联网技术的智能校园管理系统',
        plan: '预计6个月完成，分为需求分析、设计、开发、测试四个阶段'
      },
      {
        id: 2,
        title: '数据分析平台',
        type: '科研项目',
        status: 'pending',
        progress: 90,
        createdAt: '2024-01-14',
        studentName: '李四',
        teacherName: '王老师',
        level: '省级',
        description: '大数据分析平台，支持多种数据源和算法',
        plan: '预计8个月完成，包括数据采集、预处理、分析、可视化等模块'
      },
      {
        id: 3,
        title: '在线教育平台',
        type: '创新项目',
        status: 'completed',
        progress: 100,
        createdAt: '2024-01-10',
        studentName: '王五',
        teacherName: '赵老师',
        level: '国家级',
        description: '基于Web的在线教育学习平台',
        plan: '预计4个月完成，包括用户管理、课程管理、学习跟踪等模块'
      }
    ]
    totalProjects.value = projects.value.length
    updateProjectStats(projects.value)
    ElMessage.error('加载项目列表失败，显示模拟数据')
  } finally {
    loading.value = false
  }
}

// 更新项目统计数据
const updateProjectStats = (projectList) => {
  const statsData = {
    draftCount: 0,
    pendingCount: 0,
    approvedCount: 0,
    completedCount: 0,
    rejectedCount: 0,
    inProgressCount: 0
  }
  
  projectList.forEach(project => {
    const status = project.status || project.status
    if (status === 'draft') statsData.draftCount++
    else if (status === 'pending') statsData.pendingCount++
    else if (status === 'approved') statsData.approvedCount++
    else if (status === 'completed') statsData.completedCount++
    else if (status === 'rejected') statsData.rejectedCount++
    else if (status === 'in_progress') statsData.inProgressCount++
  })
  
  // 更新统计卡片数据
  stats.value = {
    draftCount: statsData.draftCount,
    pendingCount: statsData.pendingCount,
    // 将"已通过"和"进行中"合并为"进行中"
    approvedCount: statsData.approvedCount + statsData.inProgressCount,
    completedCount: statsData.completedCount
  }
  
  console.log('📊 更新统计信息:', stats.value)
}

// 重置统计数据
const resetStats = () => {
  stats.value = {
    draftCount: 0,
    pendingCount: 0,
    approvedCount: 0,
    completedCount: 0
  }
}

const handleSearch = () => {
  console.log('🔍 执行搜索，重置到第一页')
  currentPage.value = 1
  loadProjects()
}

const resetFilters = () => {
  console.log('🔄 重置所有筛选条件')
  Object.keys(filters).forEach(key => {
    filters[key] = ''
  })
  currentPage.value = 1
  loadProjects()
}

const handleSizeChange = (size) => {
  console.log(`📐 更改页大小为: ${size}`)
  pageSize.value = size
  currentPage.value = 1
  loadProjects()
}

const handleCurrentChange = (page) => {
  console.log(`📄 切换到页面: ${page}`)
  currentPage.value = page
  loadProjects()
}

const refreshProjects = () => {
  console.log('🔄 手动刷新项目列表')
  loadProjects()
}

const loadStats = async () => {
  try {
    // 如果接口有专门的统计接口，可以调用
    // const response = await studentService.getProjectStats()
    // 如果没有，就从项目列表中计算
    updateProjectStats(projects.value)
  } catch (error) {
    console.error('加载统计数据失败:', error)
    resetStats()
  }
}

const viewProjectDetail = (project) => {
  selectedProject.value = project
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
    
    // 调用服务层方法
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
  if (!currentProgressProject.value) return
  
  try {
    updatingProgress.value = true
    
    // 调用服务层方法
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
    console.error('进度更新失败:', error)
    ElMessage.error('进度更新失败')
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
    console.error('项目操作失败:', error)
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
    case 'submit':
      submitProject(project)
      break
    case 'progress':
      updateProgress(project)
      break
    case 'files':
      ElMessage.info('管理附件功能开发中...')
      break
    case 'members':
      ElMessage.info('管理成员功能开发中...')
      break
    case 'timeline':
      ElMessage.info('项目时间线功能开发中...')
      break
    case 'extend':
      ElMessage.info('申请延期功能开发中...')
      break
    case 'report':
      ElMessage.info('生成报告功能开发中...')
      break
    case 'archive':
      ElMessage.info('归档项目功能开发中...')
      break
  }
}

// 组件挂载和卸载
onMounted(() => {
  loadProjects()
})

onUnmounted(() => {
  projects.value = []
})
</script>
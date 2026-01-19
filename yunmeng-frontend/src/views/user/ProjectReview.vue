<template>
  <div class="project-review">
    <!-- 筛选条件 -->
    <el-card style="margin-bottom: 20px;">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-select v-model="filters.status" placeholder="项目状态" clearable>
            <el-option label="待审核" value="pending" />
            <el-option label="已通过" value="approved" />
            <el-option label="已拒绝" value="rejected" />
            <el-option label="需要修改" value="revision" />
          </el-select>
        </el-col>
        <el-col :span="6">
          <el-select v-model="filters.type" placeholder="项目类型" clearable>
            <el-option label="软件开发" value="software" />
            <el-option label="科研项目" value="research" />
            <el-option label="创新项目" value="innovation" />
            <el-option label="竞赛项目" value="competition" />
          </el-select>
        </el-col>
        <el-col :span="6">
          <el-date-picker
            v-model="filters.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            format="YYYY-MM-DD"
          />
        </el-col>
        <el-col :span="6">
          <el-button type="primary" @click="searchProjects">搜索</el-button>
          <el-button @click="resetFilters">重置</el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- 项目列表 -->
    <el-card>
      <template #header>
        <span>项目审核列表</span>
        <el-button link @click="refreshProjects">
          <el-icon><Refresh /></el-icon>
        </el-button>
      </template>
      
      <el-table :data="projectList" style="width: 100%" v-loading="loading">
        <el-table-column prop="name" label="项目名称" min-width="150" />
        <el-table-column prop="type" label="项目类型" width="120" />
        <el-table-column prop="applicant" label="申请人" width="100" />
        <el-table-column prop="createTime" label="申请时间" width="120" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">
              {{ scope.row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="priority" label="优先级" width="80">
          <template #default="scope">
            <el-tag :type="getPriorityType(scope.row.priority)">
              {{ scope.row.priority }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="scope">
            <el-button size="small" @click="viewProject(scope.row)">查看</el-button>
            <el-button 
              size="small" 
              type="success" 
              @click="approveProject(scope.row)"
              v-if="scope.row.status === 'pending'"
            >
              通过
            </el-button>
            <el-button 
              size="small" 
              type="danger" 
              @click="rejectProject(scope.row)"
              v-if="scope.row.status === 'pending'"
            >
              拒绝
            </el-button>
            <el-button 
              size="small" 
              type="warning" 
              @click="requestRevision(scope.row)"
              v-if="scope.row.status === 'pending'"
            >
              要求修改
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      
      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.currentPage"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 项目详情对话框 -->
    <el-dialog
      v-model="projectDetailVisible"
      title="项目详情"
      width="60%"
      :before-close="handleCloseDialog"
    >
      <div v-if="currentProject">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="项目名称">{{ currentProject.name }}</el-descriptions-item>
          <el-descriptions-item label="项目类型">{{ currentProject.type }}</el-descriptions-item>
          <el-descriptions-item label="申请人">{{ currentProject.applicant }}</el-descriptions-item>
          <el-descriptions-item label="申请时间">{{ currentProject.createTime }}</el-descriptions-item>
          <el-descriptions-item label="项目状态">{{ currentProject.status }}</el-descriptions-item>
          <el-descriptions-item label="优先级">{{ currentProject.priority }}</el-descriptions-item>
        </el-descriptions>
        
        <el-divider />
        
        <h4>项目描述</h4>
        <p>{{ currentProject.description }}</p>
        
        <el-divider />
        
        <h4>项目计划</h4>
        <p>{{ currentProject.plan }}</p>
      </div>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="projectDetailVisible = false">关闭</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 审核操作对话框 -->
    <el-dialog
      v-model="reviewDialogVisible"
      :title="reviewDialogTitle"
      width="40%"
    >
      <el-form :model="reviewForm" label-width="80px">
        <el-form-item label="审核结果">
          <el-select v-model="reviewForm.result" placeholder="请选择审核结果">
            <el-option label="通过" value="approved" />
            <el-option label="拒绝" value="rejected" />
            <el-option label="要求修改" value="revision" />
          </el-select>
        </el-form-item>
        <el-form-item label="审核意见">
          <el-input
            v-model="reviewForm.comment"
            type="textarea"
            :rows="4"
            placeholder="请输入审核意见"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="reviewDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitReview">提交审核</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'

// 响应式数据
const loading = ref(false)
const projectList = ref([])
const projectDetailVisible = ref(false)
const reviewDialogVisible = ref(false)
const currentProject = ref(null)
const reviewDialogTitle = ref('')

const filters = reactive({
  status: '',
  type: '',
  dateRange: []
})

const pagination = reactive({
  currentPage: 1,
  pageSize: 20,
  total: 0
})

const reviewForm = reactive({
  result: '',
  comment: ''
})

// 获取状态类型
const getStatusType = (status) => {
  const statusMap = {
    'pending': 'warning',
    'approved': 'success',
    'rejected': 'danger',
    'revision': 'info'
  }
  return statusMap[status] || 'info'
}

// 获取优先级类型
const getPriorityType = (priority) => {
  const priorityMap = {
    '高': 'danger',
    '中': 'warning',
    '低': 'info'
  }
  return priorityMap[priority] || 'info'
}

// 搜索项目
const searchProjects = () => {
  pagination.currentPage = 1
  loadProjects()
}

// 重置筛选条件
const resetFilters = () => {
  Object.assign(filters, {
    status: '',
    type: '',
    dateRange: []
  })
  searchProjects()
}

// 加载项目列表
const loadProjects = async () => {
  loading.value = true
  try {
    // 这里应该调用实际的API，暂时使用模拟数据
    const mockData = [
      {
        id: 1,
        name: '智能校园系统',
        type: '软件开发',
        applicant: '张三',
        createTime: '2024-01-15',
        status: 'pending',
        priority: '高',
        description: '基于物联网技术的智能校园管理系统',
        plan: '预计6个月完成，分为需求分析、设计、开发、测试四个阶段'
      },
      {
        id: 2,
        name: '数据分析平台',
        type: '科研项目',
        applicant: '李四',
        createTime: '2024-01-14',
        status: 'pending',
        priority: '中',
        description: '大数据分析平台，支持多种数据源和算法',
        plan: '预计8个月完成，包括数据采集、预处理、分析、可视化等模块'
      }
    ]
    
    projectList.value = mockData
    pagination.total = mockData.length
  } catch (error) {
    console.error('加载项目列表失败:', error)
    ElMessage.error('加载项目列表失败')
  } finally {
    loading.value = false
  }
}

// 刷新项目列表
const refreshProjects = () => {
  loadProjects()
  ElMessage.success('项目列表已刷新')
}

// 查看项目详情
const viewProject = (project) => {
  currentProject.value = project
  projectDetailVisible.value = true
}

// 通过项目
const approveProject = (project) => {
  reviewForm.result = 'approved'
  reviewForm.comment = ''
  currentProject.value = project
  reviewDialogTitle.value = '审核通过'
  reviewDialogVisible.value = true
}

// 拒绝项目
const rejectProject = (project) => {
  reviewForm.result = 'rejected'
  reviewForm.comment = ''
  currentProject.value = project
  reviewDialogTitle.value = '审核拒绝'
  reviewDialogVisible.value = true
}

// 要求修改
const requestRevision = (project) => {
  reviewForm.result = 'revision'
  reviewForm.comment = ''
  currentProject.value = project
  reviewDialogTitle.value = '要求修改'
  reviewDialogVisible.value = true
}

// 提交审核
const submitReview = async () => {
  if (!reviewForm.result) {
    ElMessage.warning('请选择审核结果')
    return
  }
  
  try {
    // 这里应该调用实际的API
    ElMessage.success('审核提交成功')
    reviewDialogVisible.value = false
    loadProjects() // 重新加载列表
  } catch (error) {
    console.error('提交审核失败:', error)
    ElMessage.error('提交审核失败')
  }
}

// 处理分页大小变化
const handleSizeChange = (val) => {
  pagination.pageSize = val
  loadProjects()
}

// 处理当前页变化
const handleCurrentChange = (val) => {
  pagination.currentPage = val
  loadProjects()
}

// 关闭对话框
const handleCloseDialog = () => {
  projectDetailVisible.value = false
  currentProject.value = null
}

// 组件挂载时加载数据
onMounted(() => {
  loadProjects()
})
</script>

<style scoped>
.project-review {
  padding: 20px;
}

.pagination-container {
  margin-top: 20px;
  text-align: right;
}

.dialog-footer {
  text-align: right;
}
</style> 
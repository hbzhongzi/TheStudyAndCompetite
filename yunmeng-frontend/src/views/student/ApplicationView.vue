<template>
  <div class="application-view">
    <!-- 页面标题 -->
    <div class="page-header">
      <h2>申请管理</h2>
      <p>管理您的项目申请和竞赛报名记录</p>
    </div>

    <!-- 筛选和搜索 -->
    <div class="filter-section">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-select v-model="filterType" placeholder="申请类型" @change="loadApplications">
            <el-option label="全部" value=""></el-option>
            <el-option label="项目申请" value="project"></el-option>
            <el-option label="竞赛报名" value="competition"></el-option>
          </el-select>
        </el-col>
        <el-col :span="6">
          <el-select v-model="filterStatus" placeholder="申请状态" @change="loadApplications">
            <el-option label="全部" value=""></el-option>
            <el-option label="待审核" value="pending"></el-option>
            <el-option label="已通过" value="approved"></el-option>
            <el-option label="已拒绝" value="rejected"></el-option>
          </el-select>
        </el-col>
        <el-col :span="8">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索申请标题"
            @input="handleSearch"
            clearable
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-col>
        <el-col :span="4">
          <el-button type="primary" @click="loadApplications">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </el-col>
      </el-row>
    </div>

    <!-- 申请列表 -->
    <div class="application-list">
      <el-table 
        :data="filteredApplications" 
        v-loading="loading"
        style="width: 100%"
        @row-click="handleRowClick"
      >
        <el-table-column prop="id" label="申请ID" width="80" />
        <el-table-column prop="title" label="申请标题" min-width="200" />
        <el-table-column prop="type" label="申请类型" width="120">
          <template #default="{ row }">
            <el-tag :type="row.type === '项目申请' ? 'primary' : 'success'">
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
        <el-table-column prop="submitTime" label="提交时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.submitTime) }}
          </template>
        </el-table-column>
        <el-table-column prop="reviewTime" label="审核时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.reviewTime) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click.stop="viewApplicationDetail(row)">查看详情</el-button>
            <el-button 
              v-if="row.status === 'pending'"
              size="small" 
              type="warning" 
              @click.stop="cancelApplication(row)"
            >
              取消申请
            </el-button>
            <el-button 
              v-if="row.status === 'rejected'"
              size="small" 
              type="primary" 
              @click.stop="resubmitApplication(row)"
            >
              重新申请
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 空状态 -->
    <el-empty 
      v-if="!loading && filteredApplications.length === 0" 
      description="暂无申请记录"
      style="margin-top: 40px;"
    />

    <!-- 申请详情对话框 -->
    <el-dialog
      v-model="showDetailDialog"
      title="申请详情"
      width="70%"
      :before-close="handleCloseDetail"
    >
      <div v-if="selectedApplication" class="application-detail">
        <div class="detail-header">
          <h2>{{ selectedApplication.title }}</h2>
          <div class="detail-tags">
            <el-tag :type="selectedApplication.type === '项目申请' ? 'primary' : 'success'">
              {{ selectedApplication.type }}
            </el-tag>
            <el-tag :type="getStatusType(selectedApplication.status)">
              {{ getStatusText(selectedApplication.status) }}
            </el-tag>
          </div>
        </div>
        
        <el-divider />
        
        <div class="detail-content">
          <el-row :gutter="20">
            <el-col :span="12">
              <div class="detail-section">
                <h3>基本信息</h3>
                <el-descriptions :column="1" border>
                  <el-descriptions-item label="申请ID">{{ selectedApplication.id }}</el-descriptions-item>
                  <el-descriptions-item label="申请类型">{{ selectedApplication.type }}</el-descriptions-item>
                  <el-descriptions-item label="申请状态">
                    <el-tag :type="getStatusType(selectedApplication.status)">
                      {{ getStatusText(selectedApplication.status) }}
                    </el-tag>
                  </el-descriptions-item>
                  <el-descriptions-item label="提交时间">{{ formatDate(selectedApplication.submitTime) }}</el-descriptions-item>
                  <el-descriptions-item label="审核时间">{{ formatDate(selectedApplication.reviewTime) }}</el-descriptions-item>
                </el-descriptions>
              </div>
            </el-col>
            
            <el-col :span="12">
              <div class="detail-section">
                <h3>申请内容</h3>
                <div class="content-preview">
                  <p><strong>描述：</strong>{{ selectedApplication.description }}</p>
                  <p v-if="selectedApplication.type === '项目申请'">
                    <strong>项目类型：</strong>{{ selectedApplication.projectType }}
                  </p>
                  <p v-if="selectedApplication.type === '竞赛报名'">
                    <strong>竞赛名称：</strong>{{ selectedApplication.competitionName }}
                  </p>
                  <p v-if="selectedApplication.teamName">
                    <strong>队伍名称：</strong>{{ selectedApplication.teamName }}
                  </p>
                </div>
              </div>
            </el-col>
          </el-row>
          
          <div class="detail-section">
            <h3>审核记录</h3>
            <el-timeline>
              <el-timeline-item 
                v-for="(review, index) in selectedApplication.reviews" 
                :key="index"
                :timestamp="formatDate(review.time)"
                :type="getReviewType(review.status)"
              >
                <el-card>
                  <h4>{{ review.reviewer }}</h4>
                  <p><strong>审核结果：</strong>
                    <el-tag :type="getStatusType(review.status)" size="small">
                      {{ getStatusText(review.status) }}
                    </el-tag>
                  </p>
                  <p v-if="review.comment"><strong>审核意见：</strong>{{ review.comment }}</p>
                </el-card>
              </el-timeline-item>
            </el-timeline>
          </div>
          
          <div v-if="selectedApplication.attachments && selectedApplication.attachments.length > 0" class="detail-section">
            <h3>附件文件</h3>
            <el-table :data="selectedApplication.attachments" border>
              <el-table-column prop="fileName" label="文件名" />
              <el-table-column prop="fileSize" label="文件大小" width="120" />
              <el-table-column prop="uploadTime" label="上传时间" width="180">
                <template #default="{ row }">
                  {{ formatDate(row.uploadTime) }}
                </template>
              </el-table-column>
              <el-table-column label="操作" width="120">
                <template #default="{ row }">
                  <el-button size="small" @click="downloadFile(row)">下载</el-button>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </div>
        
        <div class="detail-actions">
          <el-button @click="showDetailDialog = false">关闭</el-button>
          <el-button 
            v-if="selectedApplication.status === 'pending'"
            type="warning"
            @click="cancelApplication(selectedApplication)"
          >
            取消申请
          </el-button>
          <el-button 
            v-if="selectedApplication.status === 'rejected'"
            type="primary"
            @click="resubmitApplication(selectedApplication)"
          >
            重新申请
          </el-button>
        </div>
      </div>
    </el-dialog>

    <!-- 取消申请确认对话框 -->
    <el-dialog
      v-model="showCancelDialog"
      title="确认取消申请"
      width="40%"
    >
      <div class="cancel-confirm">
        <p>您确定要取消申请 <strong>{{ selectedApplication?.title }}</strong> 吗？</p>
        <p class="warning-text">取消后需要重新提交申请。</p>
        
        <el-form :model="cancelForm" label-width="100px">
          <el-form-item label="取消原因">
            <el-input 
              v-model="cancelForm.reason" 
              type="textarea" 
              placeholder="请输入取消原因（可选）"
              :rows="3"
            />
          </el-form-item>
        </el-form>
      </div>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showCancelDialog = false">取消</el-button>
          <el-button type="warning" @click="confirmCancel" :loading="cancelLoading">
            确认取消
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh } from '@element-plus/icons-vue'

// 响应式数据
const loading = ref(false)
const filterType = ref('')
const filterStatus = ref('')
const searchKeyword = ref('')
const applications = ref([])
const showDetailDialog = ref(false)
const showCancelDialog = ref(false)
const selectedApplication = ref(null)
const cancelLoading = ref(false)

// 取消表单
const cancelForm = ref({
  reason: ''
})

// 模拟申请数据
const mockApplications = [
  {
    id: 1,
    title: '智能校园管理系统项目申请',
    type: '项目申请',
    status: 'pending',
    submitTime: '2024-01-15T10:00:00Z',
    reviewTime: null,
    description: '基于物联网技术的智能校园管理系统，实现校园设施的智能化管理',
    projectType: '科研项目',
    teamName: '智能校园团队',
    reviews: [],
    attachments: [
      {
        fileName: '项目申请书.pdf',
        fileSize: '2.5MB',
        uploadTime: '2024-01-15T10:00:00Z'
      },
      {
        fileName: '技术方案.docx',
        fileSize: '1.8MB',
        uploadTime: '2024-01-15T10:00:00Z'
      }
    ]
  },
  {
    id: 2,
    title: '全国大学生程序设计竞赛报名',
    type: '竞赛报名',
    status: 'approved',
    submitTime: '2024-01-10T14:30:00Z',
    reviewTime: '2024-01-12T09:15:00Z',
    description: '报名参加全国大学生程序设计竞赛',
    competitionName: '全国大学生程序设计竞赛',
    teamName: '编程精英队',
    reviews: [
      {
        reviewer: '张老师',
        status: 'approved',
        time: '2024-01-12T09:15:00Z',
        comment: '申请材料完整，符合参赛条件，审核通过。'
      }
    ],
    attachments: [
      {
        fileName: '报名表.pdf',
        fileSize: '500KB',
        uploadTime: '2024-01-10T14:30:00Z'
      }
    ]
  },
  {
    id: 3,
    title: '在线学习平台开发项目申请',
    type: '项目申请',
    status: 'rejected',
    submitTime: '2024-01-05T09:15:00Z',
    reviewTime: '2024-01-08T16:30:00Z',
    description: '开发一个支持多种学习模式的在线教育平台',
    projectType: '竞赛项目',
    teamName: '在线教育团队',
    reviews: [
      {
        reviewer: '李老师',
        status: 'rejected',
        time: '2024-01-08T16:30:00Z',
        comment: '项目描述不够详细，技术方案需要进一步完善。'
      }
    ],
    attachments: [
      {
        fileName: '项目申请书.pdf',
        fileSize: '3.2MB',
        uploadTime: '2024-01-05T09:15:00Z'
      }
    ]
  }
]

// 过滤后的申请列表
const filteredApplications = computed(() => {
  let filtered = applications.value

  // 类型筛选
  if (filterType.value) {
    filtered = filtered.filter(a => a.type === filterType.value)
  }

  // 状态筛选
  if (filterStatus.value) {
    filtered = filtered.filter(a => a.status === filterStatus.value)
  }

  // 关键词搜索
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    filtered = filtered.filter(a => 
      a.title.toLowerCase().includes(keyword) || 
      a.description.toLowerCase().includes(keyword)
    )
  }

  return filtered
})

// 加载申请数据
const loadApplications = async () => {
  loading.value = true
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500))
    applications.value = mockApplications
  } catch (error) {
    console.error('加载申请数据失败:', error)
    ElMessage.error('加载申请数据失败')
  } finally {
    loading.value = false
  }
}

// 状态类型映射
const getStatusType = (status) => {
  const statusMap = {
    pending: 'warning',
    approved: 'success',
    rejected: 'danger'
  }
  return statusMap[status] || 'info'
}

// 状态文本映射
const getStatusText = (status) => {
  const statusMap = {
    pending: '待审核',
    approved: '已通过',
    rejected: '已拒绝'
  }
  return statusMap[status] || status
}

// 审核类型映射
const getReviewType = (status) => {
  const typeMap = {
    approved: 'success',
    rejected: 'danger'
  }
  return typeMap[status] || 'primary'
}

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return '暂无'
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

// 搜索处理
const handleSearch = () => {
  // 实时搜索，不需要额外处理，computed会自动更新
}

// 查看申请详情
const viewApplicationDetail = (application) => {
  selectedApplication.value = application
  showDetailDialog.value = true
}

// 关闭详情对话框
const handleCloseDetail = () => {
  showDetailDialog.value = false
  selectedApplication.value = null
}

// 行点击事件
const handleRowClick = (row) => {
  viewApplicationDetail(row)
}

// 取消申请
const cancelApplication = (application) => {
  selectedApplication.value = application
  cancelForm.value.reason = ''
  showCancelDialog.value = true
}

// 确认取消申请
const confirmCancel = async () => {
  cancelLoading.value = true
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // 更新申请状态
    const index = applications.value.findIndex(a => a.id === selectedApplication.value.id)
    if (index !== -1) {
      applications.value[index].status = 'cancelled'
      applications.value[index].reviewTime = new Date().toISOString()
    }
    
    showCancelDialog.value = false
    ElMessage.success('申请已取消')
  } catch (error) {
    console.error('取消申请失败:', error)
    ElMessage.error('取消申请失败，请稍后重试')
  } finally {
    cancelLoading.value = false
  }
}

// 重新申请
const resubmitApplication = (application) => {
  ElMessage.info('重新申请功能开发中...')
}

// 下载文件
const downloadFile = (file) => {
  ElMessage.info(`下载文件：${file.fileName}`)
}

// 组件挂载时加载数据
onMounted(() => {
  loadApplications()
})
</script>

<style scoped>
.application-view {
  padding: 20px;
}

.page-header {
  margin-bottom: 30px;
  text-align: center;
}

.page-header h2 {
  margin: 0 0 10px 0;
  color: #2c3e50;
  font-size: 28px;
  font-weight: 600;
}

.page-header p {
  margin: 0;
  color: #7f8c8d;
  font-size: 16px;
}

.filter-section {
  margin-bottom: 30px;
  padding: 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.application-list {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.el-table {
  border-radius: 8px;
}

.el-table :deep(.el-table__row) {
  cursor: pointer;
}

.el-table :deep(.el-table__row:hover) {
  background-color: #f5f7fa;
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
}

.detail-header h2 {
  margin: 0;
  color: #2c3e50;
}

.detail-tags {
  display: flex;
  gap: 10px;
}

.detail-section {
  margin-bottom: 25px;
}

.detail-section h3 {
  margin: 0 0 15px 0;
  color: #2c3e50;
  font-size: 18px;
  font-weight: 600;
}

.content-preview p {
  margin: 0 0 10px 0;
  line-height: 1.6;
  color: #5a6c7d;
}

.detail-actions {
  display: flex;
  justify-content: flex-end;
  gap: 15px;
  margin-top: 30px;
  padding-top: 20px;
  border-top: 1px solid #e9ecef;
}

.cancel-confirm {
  padding: 20px 0;
}

.warning-text {
  color: #e6a23c;
  font-size: 14px;
  margin: 10px 0 20px 0;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style> 
<template>
  <div class="notification-view">
    <!-- 页面标题 -->
    <div class="page-header">
      <h2>消息通知</h2>
      <p>查看系统通知和审核状态变更</p>
    </div>

    <!-- 筛选和操作 -->
    <div class="filter-section">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-select v-model="filterType" placeholder="通知类型" @change="loadNotifications">
            <el-option label="全部" value=""></el-option>
            <el-option label="系统通知" value="system"></el-option>
            <el-option label="审核通知" value="review"></el-option>
            <el-option label="竞赛通知" value="competition"></el-option>
          </el-select>
        </el-col>
        <el-col :span="6">
          <el-select v-model="filterStatus" placeholder="阅读状态" @change="loadNotifications">
            <el-option label="全部" value=""></el-option>
            <el-option label="未读" value="unread"></el-option>
            <el-option label="已读" value="read"></el-option>
          </el-select>
        </el-col>
        <el-col :span="8">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索通知内容"
            @input="handleSearch"
            clearable
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-col>
        <el-col :span="4">
          <el-button type="primary" @click="loadNotifications">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </el-col>
      </el-row>
      
      <div class="bulk-actions" style="margin-top: 15px;">
        <el-button size="small" @click="markAllAsRead" :disabled="!hasUnread">
          全部标记为已读
        </el-button>
        <el-button size="small" type="danger" @click="deleteSelected" :disabled="selectedNotifications.length === 0">
          删除选中
        </el-button>
      </div>
    </div>

    <!-- 通知列表 -->
    <div class="notification-list">
      <el-table 
        :data="filteredNotifications" 
        v-loading="loading"
        style="width: 100%"
        @selection-change="handleSelectionChange"
        @row-click="handleRowClick"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" label="通知标题" min-width="200">
          <template #default="{ row }">
            <div class="notification-title">
                              <el-icon v-if="!row.isRead" class="unread-icon"><Check /></el-icon>
              <span :class="{ 'unread': !row.isRead }">{{ row.title }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型" width="120">
          <template #default="{ row }">
            <el-tag :type="getTypeTag(row.type)" size="small">
              {{ getTypeText(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="priority" label="优先级" width="100">
          <template #default="{ row }">
            <el-tag :type="getPriorityType(row.priority)" size="small">
              {{ getPriorityText(row.priority) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createTime" label="发送时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.createTime) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click.stop="viewNotification(row)">查看</el-button>
            <el-button 
              v-if="!row.isRead"
              size="small" 
              type="success" 
              @click.stop="markAsRead(row)"
            >
              标记已读
            </el-button>
            <el-button 
              size="small" 
              type="danger" 
              @click.stop="deleteNotification(row)"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 空状态 -->
    <el-empty 
      v-if="!loading && filteredNotifications.length === 0" 
      description="暂无通知消息"
      style="margin-top: 40px;"
    />

    <!-- 通知详情对话框 -->
    <el-dialog
      v-model="showDetailDialog"
      title="通知详情"
      width="60%"
      :before-close="handleCloseDetail"
    >
      <div v-if="selectedNotification" class="notification-detail">
        <div class="detail-header">
          <h2>{{ selectedNotification.title }}</h2>
          <div class="detail-tags">
            <el-tag :type="getTypeTag(selectedNotification.type)">
              {{ getTypeText(selectedNotification.type) }}
            </el-tag>
            <el-tag :type="getPriorityType(selectedNotification.priority)">
              {{ getPriorityText(selectedNotification.priority) }}
            </el-tag>
            <el-tag v-if="!selectedNotification.isRead" type="warning">未读</el-tag>
          </div>
        </div>
        
        <el-divider />
        
        <div class="detail-content">
          <div class="detail-section">
            <h3>通知内容</h3>
            <div class="content-body" v-html="selectedNotification.content"></div>
          </div>
          
          <div class="detail-section">
            <h3>详细信息</h3>
            <el-descriptions :column="2" border>
              <el-descriptions-item label="发送时间">
                {{ formatDate(selectedNotification.createTime) }}
              </el-descriptions-item>
              <el-descriptions-item label="阅读时间">
                {{ formatDate(selectedNotification.readTime) }}
              </el-descriptions-item>
              <el-descriptions-item label="发送者">
                {{ selectedNotification.sender }}
              </el-descriptions-item>
              <el-descriptions-item label="优先级">
                <el-tag :type="getPriorityType(selectedNotification.priority)">
                  {{ getPriorityText(selectedNotification.priority) }}
                </el-tag>
              </el-descriptions-item>
            </el-descriptions>
          </div>
          
          <div v-if="selectedNotification.relatedData" class="detail-section">
            <h3>相关数据</h3>
            <el-card>
              <pre>{{ JSON.stringify(selectedNotification.relatedData, null, 2) }}</pre>
            </el-card>
          </div>
        </div>
        
        <div class="detail-actions">
          <el-button @click="showDetailDialog = false">关闭</el-button>
          <el-button 
            v-if="!selectedNotification.isRead"
            type="success"
            @click="markAsRead(selectedNotification)"
          >
            标记已读
          </el-button>
          <el-button 
            type="danger"
            @click="deleteNotification(selectedNotification)"
          >
            删除通知
          </el-button>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Check } from '@element-plus/icons-vue'

// 响应式数据
const loading = ref(false)
const filterType = ref('')
const filterStatus = ref('')
const searchKeyword = ref('')
const notifications = ref([])
const selectedNotifications = ref([])
const showDetailDialog = ref(false)
const selectedNotification = ref(null)

// 模拟通知数据
const mockNotifications = [
  {
    id: 1,
    title: '项目申请审核通过',
    content: '<p>恭喜！您提交的"智能校园管理系统"项目申请已通过审核。</p><p>项目编号：PRJ2024001</p><p>审核人：张老师</p><p>审核时间：2024-01-15 14:30:00</p>',
    type: 'review',
    priority: 'high',
    isRead: false,
    createTime: '2024-01-15T14:30:00Z',
    readTime: null,
    sender: '系统管理员',
    relatedData: {
      projectId: 1,
      projectTitle: '智能校园管理系统',
      reviewer: '张老师'
    }
  },
  {
    id: 2,
    title: '竞赛报名成功',
    content: '<p>您已成功报名参加"全国大学生程序设计竞赛"。</p><p>竞赛时间：2024-03-01 09:00:00</p><p>竞赛地点：计算机学院实验室</p><p>请提前做好准备，祝您取得好成绩！</p>',
    type: 'competition',
    priority: 'medium',
    isRead: true,
    createTime: '2024-01-10T10:00:00Z',
    readTime: '2024-01-10T10:15:00Z',
    sender: '竞赛管理员',
    relatedData: {
      competitionId: 1,
      competitionName: '全国大学生程序设计竞赛',
      teamName: '编程精英队'
    }
  },
  {
    id: 3,
    title: '系统维护通知',
    content: '<p>系统将于2024-01-20 02:00:00 至 06:00:00 进行维护升级。</p><p>维护期间系统将暂停服务，请提前保存相关工作。</p><p>给您带来的不便，敬请谅解。</p>',
    type: 'system',
    priority: 'low',
    isRead: false,
    createTime: '2024-01-18T09:00:00Z',
    readTime: null,
    sender: '系统管理员',
    relatedData: {
      maintenanceStart: '2024-01-20T02:00:00Z',
      maintenanceEnd: '2024-01-20T06:00:00Z'
    }
  },
  {
    id: 4,
    title: '项目申请被拒绝',
    content: '<p>很遗憾，您提交的"在线学习平台开发"项目申请未通过审核。</p><p>拒绝原因：项目描述不够详细，技术方案需要进一步完善。</p><p>建议：请完善项目描述和技术方案后重新提交。</p>',
    type: 'review',
    priority: 'high',
    isRead: true,
    createTime: '2024-01-08T16:30:00Z',
    readTime: '2024-01-08T17:00:00Z',
    sender: '李老师',
    relatedData: {
      projectId: 3,
      projectTitle: '在线学习平台开发',
      reason: '项目描述不够详细，技术方案需要进一步完善'
    }
  },
  {
    id: 5,
    title: '新功能上线通知',
    content: '<p>系统新增了文件管理功能，支持在线预览和下载。</p><p>新功能包括：</p><ul><li>支持多种文件格式预览</li><li>文件在线编辑</li><li>文件版本管理</li></ul><p>欢迎体验新功能！</p>',
    type: 'system',
    priority: 'medium',
    isRead: false,
    createTime: '2024-01-16T11:00:00Z',
    readTime: null,
    sender: '系统管理员',
    relatedData: {
      feature: '文件管理',
      version: '2.1.0'
    }
  }
]

// 过滤后的通知列表
const filteredNotifications = computed(() => {
  let filtered = notifications.value

  // 类型筛选
  if (filterType.value) {
    filtered = filtered.filter(n => n.type === filterType.value)
  }

  // 状态筛选
  if (filterStatus.value) {
    filtered = filtered.filter(n => {
      if (filterStatus.value === 'read') return n.isRead
      if (filterStatus.value === 'unread') return !n.isRead
      return true
    })
  }

  // 关键词搜索
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    filtered = filtered.filter(n => 
      n.title.toLowerCase().includes(keyword) || 
      n.content.toLowerCase().includes(keyword)
    )
  }

  return filtered
})

// 是否有未读通知
const hasUnread = computed(() => {
  return notifications.value.some(n => !n.isRead)
})

// 加载通知数据
const loadNotifications = async () => {
  loading.value = true
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500))
    notifications.value = mockNotifications
  } catch (error) {
    console.error('加载通知数据失败:', error)
    ElMessage.error('加载通知数据失败')
  } finally {
    loading.value = false
  }
}

// 类型标签映射
const getTypeTag = (type) => {
  const tagMap = {
    system: 'info',
    review: 'primary',
    competition: 'success'
  }
  return tagMap[type] || 'info'
}

// 类型文本映射
const getTypeText = (type) => {
  const textMap = {
    system: '系统通知',
    review: '审核通知',
    competition: '竞赛通知'
  }
  return textMap[type] || type
}

// 优先级类型映射
const getPriorityType = (priority) => {
  const typeMap = {
    high: 'danger',
    medium: 'warning',
    low: 'info'
  }
  return typeMap[priority] || 'info'
}

// 优先级文本映射
const getPriorityText = (priority) => {
  const textMap = {
    high: '高',
    medium: '中',
    low: '低'
  }
  return textMap[priority] || priority
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

// 查看通知详情
const viewNotification = (notification) => {
  selectedNotification.value = notification
  showDetailDialog.value = true
  
  // 如果未读，自动标记为已读
  if (!notification.isRead) {
    markAsRead(notification)
  }
}

// 关闭详情对话框
const handleCloseDetail = () => {
  showDetailDialog.value = false
  selectedNotification.value = null
}

// 行点击事件
const handleRowClick = (row) => {
  viewNotification(row)
}

// 标记为已读
const markAsRead = async (notification) => {
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 300))
    
    notification.isRead = true
    notification.readTime = new Date().toISOString()
    
    ElMessage.success('已标记为已读')
  } catch (error) {
    console.error('标记失败:', error)
    ElMessage.error('标记失败')
  }
}

// 删除通知
const deleteNotification = async (notification) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除通知"${notification.title}"吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 300))
    
    const index = notifications.value.findIndex(n => n.id === notification.id)
    if (index !== -1) {
      notifications.value.splice(index, 1)
      ElMessage.success('删除成功')
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 选择变化处理
const handleSelectionChange = (selection) => {
  selectedNotifications.value = selection
}

// 全部标记为已读
const markAllAsRead = async () => {
  try {
    await ElMessageBox.confirm(
      '确定要将所有未读通知标记为已读吗？',
      '确认操作',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'info'
      }
    )
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500))
    
    notifications.value.forEach(notification => {
      if (!notification.isRead) {
        notification.isRead = true
        notification.readTime = new Date().toISOString()
      }
    })
    
    ElMessage.success('全部标记为已读')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('操作失败')
    }
  }
}

// 删除选中通知
const deleteSelected = async () => {
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedNotifications.value.length} 条通知吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500))
    
    const selectedIds = selectedNotifications.value.map(n => n.id)
    notifications.value = notifications.value.filter(n => !selectedIds.includes(n.id))
    
    ElMessage.success('删除成功')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 组件挂载时加载数据
onMounted(() => {
  loadNotifications()
})
</script>

<style scoped>
.notification-view {
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

.bulk-actions {
  display: flex;
  gap: 10px;
}

.notification-list {
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

.notification-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.unread-icon {
  color: #409eff;
  font-size: 12px;
}

.unread {
  font-weight: 600;
  color: #2c3e50;
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

.content-body {
  line-height: 1.6;
  color: #5a6c7d;
}

.content-body p {
  margin: 0 0 10px 0;
}

.content-body ul {
  margin: 10px 0;
  padding-left: 20px;
}

.content-body li {
  margin-bottom: 5px;
}

.detail-actions {
  display: flex;
  justify-content: flex-end;
  gap: 15px;
  margin-top: 30px;
  padding-top: 20px;
  border-top: 1px solid #e9ecef;
}

pre {
  background: #f8f9fa;
  padding: 15px;
  border-radius: 4px;
  font-size: 12px;
  overflow-x: auto;
}
</style> 
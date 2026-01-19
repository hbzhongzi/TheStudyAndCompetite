<template>
  <div class="teacher-project-files">
    <el-card>
      <template #header>
        <div class="header-content">
          <span>项目文件管理</span>
          <el-button type="primary" @click="refreshFiles">刷新</el-button>
        </div>
      </template>

      <!-- 项目选择 -->
      <el-form :inline="true" class="project-selector">
        <el-form-item label="选择项目">
          <el-select v-model="selectedProjectId" placeholder="请选择项目" @change="loadFiles">
            <el-option
              v-for="project in projectList"
              :key="project.id"
              :label="project.name || '未命名项目'"
              :value="project.id || ''"
              v-if="project && project.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="success" @click="loadFiles" :disabled="!selectedProjectId">
            加载文件
          </el-button>
        </el-form-item>
      </el-form>

      <!-- 文件管理区域 -->
      <div v-if="selectedProjectId && files.length > 0" class="files-container">
        <!-- 文件统计 -->
        <el-row :gutter="20" class="file-stats">
          <el-col :span="6" v-for="stat in fileStats" :key="stat.label">
            <el-card class="stat-card" :class="stat.type">
              <div class="stat-content">
                <div class="stat-icon">
                  <el-icon><component :is="stat.icon" /></el-icon>
                </div>
                <div class="stat-info">
                  <h4>{{ stat.label }}</h4>
                  <p class="stat-number">{{ stat.value }}</p>
                  <p class="stat-desc">{{ stat.description }}</p>
                </div>
              </div>
            </el-card>
          </el-col>
        </el-row>

        <!-- 搜索和筛选 -->
        <el-card class="search-card">
          <el-form :inline="true" class="search-form">
            <el-form-item label="搜索">
              <el-input
                v-model="searchQuery"
                placeholder="搜索文件名或描述"
                clearable
                @input="handleSearch"
              >
                <template #prefix>
                  <el-icon><Search /></el-icon>
                </template>
              </el-input>
            </el-form-item>
            <el-form-item label="文件类型">
              <el-select v-model="selectedCategory" placeholder="所有类型" clearable @change="handleSearch">
                <el-option label="所有类型" :value="''" />
                <el-option label="文档" value="document" />
                <el-option label="代码" value="code" />
                <el-option label="图片" value="image" />
                <el-option label="视频" value="video" />
                <el-option label="其他" value="other" />
              </el-select>
            </el-form-item>
            <el-form-item label="排序">
              <el-select v-model="sortBy" @change="handleSearch">
                <el-option label="上传时间" value="uploadTime" />
                <el-option label="文件大小" value="size" />
                <el-option label="文件名" value="name" />
              </el-select>
            </el-form-item>
          </el-form>
        </el-card>

        <!-- 视图切换 -->
        <div class="view-controls">
          <el-radio-group v-model="viewMode" size="large">
            <el-radio-button label="table">表格视图</el-radio-button>
            <el-radio-button label="grid">网格视图</el-radio-button>
          </el-radio-group>
        </div>

        <!-- 文件列表 -->
        <div v-if="viewMode === 'table'" class="table-view">
          <el-table :data="filteredFiles" style="width: 100%" @selection-change="handleSelectionChange">
            <el-table-column type="selection" width="55" />
            <el-table-column prop="name" label="文件名" min-width="200">
              <template #default="scope">
                <div class="file-name">
                  <el-icon class="file-icon" :class="getFileIconClass(scope.row.type)">
                    <component :is="getFileIcon(scope.row.type)" />
                  </el-icon>
                  <span>{{ scope.row.name }}</span>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="type" label="类型" width="100">
              <template #default="scope">
                <el-tag :type="getFileTypeTag(scope.row.type)" size="small">
                  {{ getFileTypeLabel(scope.row.type) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="size" label="大小" width="100">
              <template #default="scope">
                {{ formatFileSize(scope.row.size) }}
              </template>
            </el-table-column>
            <el-table-column prop="uploadTime" label="上传时间" width="150" />
            <el-table-column prop="uploader" label="上传者" width="100" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="scope">
                <el-tag :type="getFileStatusType(scope.row.status)" size="small">
                  {{ getFileStatusLabel(scope.row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200" fixed="right">
              <template #default="scope">
                <el-button size="small" @click="previewFile(scope.row)">预览</el-button>
                <el-button size="small" type="primary" @click="downloadFile(scope.row)">下载</el-button>
                <el-button size="small" type="warning" @click="reviewFile(scope.row)">审核</el-button>
                <el-button size="small" type="info" @click="viewFileDetail(scope.row)">详情</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <!-- 网格视图 -->
        <div v-else class="grid-view">
          <el-row :gutter="20">
            <el-col :span="6" v-for="file in filteredFiles" :key="file.id">
              <el-card class="file-card" :class="{ selected: selectedFiles.includes(file.id) }">
                <div class="file-card-header">
                  <el-checkbox 
                    v-model="file.selected" 
                    @change="(checked) => handleFileSelect(file.id, checked)"
                  />
                  <el-dropdown @command="(command) => handleFileAction(command, file)">
                    <el-button type="text" size="small">
                      <el-icon><More /></el-icon>
                    </el-button>
                    <template #dropdown>
                      <el-dropdown-menu>
                        <el-dropdown-item command="preview">预览</el-dropdown-item>
                        <el-dropdown-item command="download">下载</el-dropdown-item>
                        <el-dropdown-item command="review">审核</el-dropdown-item>
                        <el-dropdown-item command="detail">详情</el-dropdown-item>
                      </el-dropdown-menu>
                    </template>
                  </el-dropdown>
                </div>
                
                <div class="file-card-content" @click="viewFileDetail(file)">
                  <div class="file-icon-large">
                    <el-icon :class="getFileIconClass(file.type)">
                      <component :is="getFileIcon(file.type)" />
                    </el-icon>
                  </div>
                  <h4 class="file-title">{{ file.name }}</h4>
                  <p class="file-info">{{ formatFileSize(file.size) }} • {{ getFileTypeLabel(file.type) }}</p>
                  <p class="file-uploader">上传者: {{ file.uploader }}</p>
                  <p class="file-time">{{ file.uploadTime }}</p>
                </div>
                
                <div class="file-card-footer">
                  <el-tag :type="getFileStatusType(file.status)" size="small">
                    {{ getFileStatusLabel(file.status) }}
                  </el-tag>
                </div>
              </el-card>
            </el-col>
          </el-row>
        </div>

        <!-- 分页 -->
        <div class="pagination-container">
          <el-pagination
            v-model:current-page="currentPage"
            v-model:page-size="pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="filteredFiles.length"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="handleSizeChange"
            @current-change="handleCurrentChange"
          />
        </div>
      </div>

      <!-- 空状态 -->
      <el-empty
        v-else-if="selectedProjectId && files.length === 0"
        description="该项目暂无文件"
      >
        <el-button type="primary" @click="uploadFile">上传第一个文件</el-button>
      </el-empty>

      <!-- 项目选择提示 -->
      <el-empty
        v-else
        description="请先选择一个项目"
      />
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Document, Search, Grid, List, Download, Delete, Edit, View } from '@element-plus/icons-vue'
import { teacherService } from '../../services/teacherService'
import { validateProjectList, validateApiResponse, getDefaultProjects } from '../../utils/dataValidator'

// 响应式数据
const selectedProjectId = ref('')
const projectList = ref([])
const files = ref([])
const loading = ref(false)
const searchQuery = ref('')
const selectedCategory = ref('')
const sortBy = ref('uploadTime')
const viewMode = ref('table')
const currentPage = ref(1)
const pageSize = ref(20)
const selectedFiles = ref([])

// 文件统计
const fileStats = computed(() => [
  {
    label: '总文件数',
    value: files.value.length,
    description: '所有文件',
    icon: Files,
    type: 'total'
  },
  {
    label: '待审核',
    value: files.value.filter(f => f.status === 'pending').length,
    description: '需要审核',
    icon: Document,
    type: 'pending'
  },
  {
    label: '已通过',
    value: files.value.filter(f => f.status === 'approved').length,
    description: '审核通过',
    icon: Folder,
    type: 'approved'
  },
  {
    label: '已拒绝',
    value: files.value.filter(f => f.status === 'rejected').length,
    description: '审核拒绝',
    icon: Files,
    type: 'rejected'
  }
])

// 过滤后的文件
const filteredFiles = computed(() => {
  let result = [...files.value]
  
  // 搜索过滤
  if (searchQuery.value) {
    result = result.filter(file => 
      file.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      (file.description && file.description.toLowerCase().includes(searchQuery.value.toLowerCase()))
    )
  }
  
  // 类型过滤
  if (selectedCategory.value) {
    result = result.filter(file => file.type === selectedCategory.value)
  }
  
  // 排序
  result.sort((a, b) => {
    switch (sortBy.value) {
      case 'uploadTime':
        return new Date(b.uploadTime) - new Date(a.uploadTime)
      case 'size':
        return b.size - a.size
      case 'name':
        return a.name.localeCompare(b.name)
      default:
        return 0
    }
  })
  
  return result
})

// 获取文件图标
const getFileIcon = (type) => {
  const iconMap = {
    document: Document,
    code: Files,
    image: Picture,
    video: VideoPlay,
    other: Files
  }
  return iconMap[type] || Files
}

// 获取文件图标样式类
const getFileIconClass = (type) => {
  const classMap = {
    document: 'file-icon-document',
    code: 'file-icon-code',
    image: 'file-icon-image',
    video: 'file-icon-video',
    other: 'file-icon-other'
  }
  return classMap[type] || 'file-icon-other'
}

// 获取文件类型标签
const getFileTypeLabel = (type) => {
  const labelMap = {
    document: '文档',
    code: '代码',
    image: '图片',
    video: '视频',
    other: '其他'
  }
  return labelMap[type] || '其他'
}

// 获取文件类型标签样式
const getFileTypeTag = (type) => {
  const tagMap = {
    document: 'primary',
    code: 'success',
    image: 'warning',
    video: 'danger',
    other: 'info'
  }
  return tagMap[type] || 'info'
}

// 获取文件状态类型
const getFileStatusType = (status) => {
  const typeMap = {
    pending: 'warning',
    approved: 'success',
    rejected: 'danger',
    reviewing: 'info'
  }
  return typeMap[status] || 'info'
}

// 获取文件状态标签
const getFileStatusLabel = (status) => {
  const labelMap = {
    pending: '待审核',
    approved: '已通过',
    rejected: '已拒绝',
    reviewing: '审核中'
  }
  return labelMap[status] || '未知'
}

// 格式化文件大小
const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 加载项目列表
const loadProjects = async () => {
  try {
    const response = await teacherService.getGuidedProjects()
    
    // 使用验证工具验证响应数据
    const validatedResponse = validateApiResponse(response)
    
    if (validatedResponse.code === 200) {
      // 使用验证工具验证项目列表数据
      projectList.value = validateProjectList(validatedResponse.data, getDefaultProjects())
    } else {
      // 使用默认数据
      projectList.value = getDefaultProjects()
    }
  } catch (error) {
    console.error('加载项目列表失败:', error)
    // 使用默认数据
    projectList.value = getDefaultProjects()
  }
  
  // 验证数据完整性
  console.log('项目列表数据:', projectList.value)
}

// 加载文件
const loadFiles = async () => {
  if (!selectedProjectId.value) return
  
  loading.value = true
  try {
    // 这里应该调用实际的API
    // const response = await teacherService.getProjectFiles(selectedProjectId.value)
    
    // 使用模拟数据
    files.value = [
      {
        id: 1,
        name: '需求分析报告.pdf',
        type: 'document',
        size: 2048576,
        uploadTime: '2024-01-25 15:30:00',
        uploader: '张三',
        status: 'approved',
        description: '项目需求分析详细报告',
        path: '/files/requirements.pdf'
      },
      {
        id: 2,
        name: '系统架构设计.docx',
        type: 'document',
        size: 1536000,
        uploadTime: '2024-02-20 14:20:00',
        uploader: '李四',
        status: 'approved',
        description: '系统整体架构设计文档',
        path: '/files/architecture.docx'
      },
      {
        id: 3,
        name: '用户界面设计.psd',
        type: 'image',
        size: 52428800,
        uploadTime: '2024-03-15 10:15:00',
        uploader: '王五',
        status: 'pending',
        description: '用户界面设计稿',
        path: '/files/ui-design.psd'
      },
      {
        id: 4,
        name: '核心功能代码.zip',
        type: 'code',
        size: 10485760,
        uploadTime: '2024-04-10 16:45:00',
        uploader: '赵六',
        status: 'reviewing',
        description: '核心功能模块源代码',
        path: '/files/core-code.zip'
      },
      {
        id: 5,
        name: '测试报告.xlsx',
        type: 'document',
        size: 512000,
        uploadTime: '2024-05-05 09:30:00',
        uploader: '张三',
        status: 'pending',
        description: '功能测试结果报告',
        path: '/files/test-report.xlsx'
      },
      {
        id: 6,
        name: '演示视频.mp4',
        type: 'video',
        size: 524288000,
        uploadTime: '2024-05-20 11:20:00',
        uploader: '李四',
        status: 'approved',
        description: '系统功能演示视频',
        path: '/files/demo-video.mp4'
      }
    ]
  } catch (error) {
    console.error('加载文件失败:', error)
    ElMessage.error('加载文件失败')
  } finally {
    loading.value = false
  }
}

// 搜索处理
const handleSearch = () => {
  currentPage.value = 1
}

// 视图切换处理
const handleViewModeChange = (mode) => {
  viewMode.value = mode
}

// 文件选择处理
const handleSelectionChange = (selection) => {
  selectedFiles.value = selection.map(item => item.id)
}

const handleFileSelect = (fileId, checked) => {
  if (checked) {
    if (!selectedFiles.value.includes(fileId)) {
      selectedFiles.value.push(fileId)
    }
  } else {
    selectedFiles.value = selectedFiles.value.filter(id => id !== fileId)
  }
}

// 文件操作处理
const handleFileAction = (command, file) => {
  switch (command) {
    case 'preview':
      previewFile(file)
      break
    case 'download':
      downloadFile(file)
      break
    case 'review':
      reviewFile(file)
      break
    case 'detail':
      viewFileDetail(file)
      break
  }
}

// 预览文件
const previewFile = (file) => {
  ElMessage.info(`预览文件: ${file.name}`)
  // 这里应该实现文件预览逻辑
}

// 下载文件
const downloadFile = (file) => {
  ElMessage.success(`开始下载文件: ${file.name}`)
  // 这里应该实现文件下载逻辑
}

// 审核文件
const reviewFile = (file) => {
  ElMessage.info(`审核文件: ${file.name}`)
  // 这里应该打开文件审核对话框
}

// 查看文件详情
const viewFileDetail = (file) => {
  ElMessage.info(`查看文件详情: ${file.name}`)
  // 这里应该打开文件详情对话框
}

// 上传文件
const uploadFile = () => {
  ElMessage.info('跳转到文件上传页面')
}

// 分页处理
const handleSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
}

const handleCurrentChange = (page) => {
  currentPage.value = page
}

// 刷新文件
const refreshFiles = () => {
  loadFiles()
  ElMessage.success('文件列表已刷新')
}

// 组件挂载时加载数据
onMounted(async () => {
  try {
    await loadProjects()
    
    // 验证数据完整性
    if (projectList.value.length === 0) {
      console.warn('项目列表为空，使用默认数据')
      // 确保有默认数据
      projectList.value = [
        { id: 1, name: '智能校园系统' },
        { id: 2, name: '数据分析平台' },
        { id: 3, name: '在线教育平台' }
      ]
    }
  } catch (error) {
    console.error('组件初始化失败:', error)
    ElMessage.error('组件初始化失败，请刷新页面重试')
  }
})
</script>

<style scoped>
.teacher-project-files {
  padding: 20px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.project-selector {
  margin-bottom: 20px;
  padding: 20px;
  background-color: #f5f7fa;
  border-radius: 8px;
}

.files-container {
  margin-top: 20px;
}

.file-stats {
  margin-bottom: 20px;
}

.stat-card {
  margin-bottom: 20px;
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}

.stat-card.total {
  border-left: 4px solid #667eea;
}

.stat-card.pending {
  border-left: 4px solid #f093fb;
}

.stat-card.approved {
  border-left: 4px solid #4facfe;
}

.stat-card.rejected {
  border-left: 4px solid #43e97b;
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

.stat-icon.total {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stat-icon.pending {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.stat-icon.approved {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.stat-icon.rejected {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
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
  margin: 0 0 5px 0;
  font-size: 28px;
  font-weight: 600;
  color: #2c3e50;
}

.stat-desc {
  margin: 0;
  color: #95a5a6;
  font-size: 12px;
}

.search-card {
  margin-bottom: 20px;
}

.search-form {
  margin: 0;
}

.view-controls {
  margin-bottom: 20px;
  text-align: center;
}

.table-view {
  margin-bottom: 20px;
}

.file-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.file-icon {
  font-size: 18px;
}

.file-icon-document {
  color: #409EFF;
}

.file-icon-code {
  color: #67C23A;
}

.file-icon-image {
  color: #E6A23C;
}

.file-icon-video {
  color: #F56C6C;
}

.file-icon-other {
  color: #909399;
}

.grid-view {
  margin-bottom: 20px;
}

.file-card {
  margin-bottom: 20px;
  transition: all 0.3s ease;
  cursor: pointer;
}

.file-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}

.file-card.selected {
  border: 2px solid #409EFF;
  background-color: #f0f9ff;
}

.file-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px;
  border-bottom: 1px solid #ebeef5;
}

.file-card-content {
  padding: 20px;
  text-align: center;
}

.file-icon-large {
  font-size: 48px;
  margin-bottom: 15px;
}

.file-icon-large.file-icon-document {
  color: #409EFF;
}

.file-icon-large.file-icon-code {
  color: #67C23A;
}

.file-icon-large.file-icon-image {
  color: #E6A23C;
}

.file-icon-large.file-icon-video {
  color: #F56C6C;
}

.file-icon-large.file-icon-other {
  color: #909399;
}

.file-title {
  margin: 0 0 10px 0;
  font-size: 16px;
  color: #2c3e50;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-info {
  margin: 0 0 8px 0;
  color: #606266;
  font-size: 14px;
}

.file-uploader {
  margin: 0 0 8px 0;
  color: #909399;
  font-size: 12px;
}

.file-time {
  margin: 0;
  color: #c0c4cc;
  font-size: 12px;
}

.file-card-footer {
  padding: 10px;
  border-top: 1px solid #ebeef5;
  text-align: center;
}

.pagination-container {
  text-align: center;
  margin-top: 20px;
}

:deep(.el-table .el-table__row:hover) {
  background-color: #f5f7fa;
}

:deep(.el-table .el-table__row.selected) {
  background-color: #f0f9ff;
}

:deep(.el-card__header) {
  padding: 15px 20px;
}

:deep(.el-card__body) {
  padding: 20px;
}

:deep(.el-form-item) {
  margin-bottom: 15px;
}

:deep(.el-button--text) {
  padding: 0;
  font-size: 16px;
}

:deep(.el-dropdown-menu__item) {
  padding: 8px 20px;
}

:deep(.el-pagination) {
  justify-content: center;
}
</style> 
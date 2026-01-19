<template>
  <div class="project-files">
    <el-card>
      <template #header>
        <div class="header-content">
          <span>项目文件管理</span>
          <el-button type="primary" @click="showUploadDialog">上传文件</el-button>
        </div>
      </template>

      <!-- 项目选择 -->
      <el-form :inline="true" class="project-selector">
        <el-form-item label="选择项目">
          <el-select v-model="selectedProjectId" placeholder="请选择项目" @change="loadFiles">
            <el-option
              v-for="project in projectList"
              :key="project.id"
              :label="project.name"
              :value="project.id"
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
      <div v-if="selectedProjectId" class="files-container">
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

        <!-- 文件操作栏 -->
        <div class="file-actions">
          <el-row :gutter="20">
            <el-col :span="6">
              <el-input
                v-model="searchQuery"
                placeholder="搜索文件名"
                clearable
                @input="handleSearch"
              >
                <template #prefix>
                  <el-icon><Search /></el-icon>
                </template>
              </el-input>
            </el-col>
            <el-col :span="4">
              <el-select v-model="categoryFilter" placeholder="文件类型" clearable @change="handleFilter">
                <el-option label="全部" value="" />
                <el-option label="文档" value="document" />
                <el-option label="图片" value="image" />
                <el-option label="视频" value="video" />
                <el-option label="代码" value="code" />
                <el-option label="其他" value="other" />
              </el-select>
            </el-col>
            <el-col :span="4">
              <el-select v-model="sortBy" placeholder="排序方式" @change="handleSort">
                <el-option label="上传时间" value="uploadTime" />
                <el-option label="文件名" value="name" />
                <el-option label="文件大小" value="size" />
                <el-option label="文件类型" value="type" />
              </el-select>
            </el-col>
            <el-col :span="10">
              <el-button-group>
                <el-button @click="viewMode = 'list'" :type="viewMode === 'list' ? 'primary' : ''">
                  <el-icon><List /></el-icon>
                  列表视图
                </el-button>
                <el-button @click="viewMode = 'grid'" :type="viewMode === 'grid' ? 'primary' : ''">
                  <el-icon><Grid /></el-icon>
                  网格视图
                </el-button>
              </el-button-group>
              <el-button type="warning" @click="batchDelete" :disabled="selectedFiles.length === 0">
                批量删除
              </el-button>
            </el-col>
          </el-row>
        </div>

        <!-- 文件列表 -->
        <div v-if="filteredFiles.length > 0" class="files-content">
          <!-- 列表视图 -->
          <div v-if="viewMode === 'list'">
            <el-table :data="filteredFiles" @selection-change="handleSelectionChange" v-loading="loading">
              <el-table-column type="selection" width="55" />
              <el-table-column prop="name" label="文件名" min-width="200">
                <template #default="scope">
                  <div class="file-name">
                    <el-icon class="file-icon" :class="getFileIconClass(scope.row.type)">
                      <component :is="getFileIcon(scope.row.type)" />
                    </el-icon>
                    <span @click="previewFile(scope.row)" class="file-link">
                      {{ scope.row.name }}
                    </span>
                  </div>
                </template>
              </el-table-column>
              <el-table-column prop="category" label="类型" width="100">
                <template #default="scope">
                  <el-tag :type="getCategoryType(scope.row.category)" size="small">
                    {{ getCategoryLabel(scope.row.category) }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="size" label="大小" width="120">
                <template #default="scope">
                  {{ formatFileSize(scope.row.size) }}
                </template>
              </el-table-column>
              <el-table-column prop="uploadTime" label="上传时间" width="150" />
              <el-table-column prop="uploader" label="上传者" width="100" />
              <el-table-column label="操作" width="200" fixed="right">
                <template #default="scope">
                  <el-button size="small" @click="previewFile(scope.row)">预览</el-button>
                  <el-button size="small" type="primary" @click="downloadFile(scope.row)">下载</el-button>
                  <el-button size="small" type="warning" @click="editFile(scope.row)">编辑</el-button>
                  <el-button size="small" type="danger" @click="deleteFile(scope.row)">删除</el-button>
                </template>
              </el-table-column>
            </el-table>
          </div>

          <!-- 网格视图 -->
          <div v-else class="files-grid">
            <el-row :gutter="20">
              <el-col :span="6" v-for="file in filteredFiles" :key="file.id">
                <el-card class="file-card" :class="{ selected: selectedFiles.includes(file.id) }">
                  <div class="file-card-header">
                    <el-checkbox 
                      :model-value="selectedFiles.includes(file.id)"
                      @change="(checked) => toggleFileSelection(file.id, checked)"
                    />
                    <el-dropdown @command="handleFileAction">
                      <el-button type="text" size="small">
                        <el-icon><More /></el-icon>
                      </el-button>
                      <template #dropdown>
                        <el-dropdown-menu>
                          <el-dropdown-item :command="{ action: 'preview', file }">预览</el-dropdown-item>
                          <el-dropdown-item :command="{ action: 'download', file }">下载</el-dropdown-item>
                          <el-dropdown-item :command="{ action: 'edit', file }">编辑</el-dropdown-item>
                          <el-dropdown-item :command="{ action: 'delete', file }" divided>删除</el-dropdown-item>
                        </el-dropdown-menu>
                      </template>
                    </el-dropdown>
                  </div>
                  
                  <div class="file-card-content" @click="previewFile(file)">
                    <div class="file-icon-large">
                      <el-icon :class="getFileIconClass(file.type)">
                        <component :is="getFileIcon(file.type)" />
                      </el-icon>
                    </div>
                    <h4 class="file-title">{{ file.name }}</h4>
                    <p class="file-info">{{ formatFileSize(file.size) }} • {{ getCategoryLabel(file.category) }}</p>
                    <p class="file-uploader">上传者: {{ file.uploader }}</p>
                    <p class="file-time">{{ file.uploadTime }}</p>
                  </div>
                </el-card>
              </el-col>
            </el-row>
          </div>
        </div>

        <!-- 空状态 -->
        <el-empty
          v-else
          description="暂无文件"
        >
          <el-button type="primary" @click="showUploadDialog">上传第一个文件</el-button>
        </el-empty>
      </div>

      <!-- 项目选择提示 -->
      <el-empty
        v-else
        description="请先选择一个项目"
      />
    </el-card>

    <!-- 文件上传对话框 -->
    <el-dialog
      v-model="uploadDialogVisible"
      title="上传文件"
      width="50%"
    >
      <el-form :model="uploadForm" :rules="uploadRules" ref="uploadFormRef" label-width="100px">
        <el-form-item label="文件" prop="files">
          <el-upload
            ref="uploadRef"
            :auto-upload="false"
            :on-change="handleFileChange"
            :on-remove="handleFileRemove"
            :file-list="uploadForm.files"
            multiple
            drag
          >
            <el-icon class="el-icon--upload"><Upload /></el-icon>
            <div class="el-upload__text">
              将文件拖到此处，或<em>点击上传</em>
            </div>
            <template #tip>
              <div class="el-upload__tip">
                支持任意类型文件，单个文件不超过100MB
              </div>
            </template>
          </el-upload>
        </el-form-item>
        
        <el-form-item label="文件分类" prop="category">
          <el-select v-model="uploadForm.category" placeholder="选择文件分类">
            <el-option label="文档" value="document" />
            <el-option label="图片" value="image" />
            <el-option label="视频" value="video" />
            <el-option label="代码" value="code" />
            <el-option label="其他" value="other" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="文件描述">
          <el-input
            v-model="uploadForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入文件描述"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="uploadDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitUpload" :loading="uploading">开始上传</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 文件预览对话框 -->
    <el-dialog
      v-model="previewDialogVisible"
      title="文件预览"
      width="70%"
    >
      <div v-if="currentFile" class="file-preview">
        <div class="file-info-header">
          <h3>{{ currentFile.name }}</h3>
          <div class="file-meta">
            <el-tag :type="getCategoryType(currentFile.category)" size="small">
              {{ getCategoryLabel(currentFile.category) }}
            </el-tag>
            <span class="file-size">{{ formatFileSize(currentFile.size) }}</span>
            <span class="file-uploader">上传者: {{ currentFile.uploader }}</span>
            <span class="file-time">{{ currentFile.uploadTime }}</span>
          </div>
        </div>
        
        <el-divider />
        
        <div v-if="currentFile.description" class="file-description">
          <h4>文件描述</h4>
          <p>{{ currentFile.description }}</p>
        </div>
        
        <!-- 文件预览内容 -->
        <div class="preview-content">
          <div v-if="isImageFile(currentFile.type)" class="image-preview">
            <img :src="currentFile.url || '/placeholder-image.jpg'" :alt="currentFile.name" />
          </div>
          <div v-else-if="isTextFile(currentFile.type)" class="text-preview">
            <pre>{{ currentFile.content || '文件内容预览不可用' }}</pre>
          </div>
          <div v-else class="no-preview">
            <el-icon><Document /></el-icon>
            <p>此文件类型不支持预览</p>
            <el-button type="primary" @click="downloadFile(currentFile)">下载文件</el-button>
          </div>
        </div>
      </div>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="previewDialogVisible = false">关闭</el-button>
          <el-button type="primary" @click="downloadFile(currentFile)">下载</el-button>
          <el-button type="warning" @click="editFile(currentFile)">编辑</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Document, Folder, Picture, VideoPlay, Files, 
  Search, List, Grid, More, Upload 
} from '@element-plus/icons-vue'
import { studentService } from '../../services/studentService'

// 响应式数据
const selectedProjectId = ref('')
const projectList = ref([])
const files = ref([])
const loading = ref(false)
const uploadDialogVisible = ref(false)
const previewDialogVisible = ref(false)
const currentFile = ref(null)
const viewMode = ref('list')
const searchQuery = ref('')
const categoryFilter = ref('')
const sortBy = ref('uploadTime')
const selectedFiles = ref([])
const uploading = ref(false)

const uploadForm = ref({
  files: [],
  category: 'document',
  description: ''
})

const uploadRules = {
  files: [{ required: true, message: '请选择要上传的文件', trigger: 'change' }],
  category: [{ required: true, message: '请选择文件分类', trigger: 'change' }]
}

const uploadFormRef = ref()
const uploadRef = ref()

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
    label: '文档文件',
    value: files.value.filter(f => f.category === 'document').length,
    description: '文档类型',
    icon: Document,
    type: 'document'
  },
  {
    label: '图片文件',
    value: files.value.filter(f => f.category === 'image').length,
    description: '图片类型',
    icon: Picture,
    type: 'image'
  },
  {
    label: '代码文件',
    value: files.value.filter(f => f.category === 'code').length,
    description: '代码类型',
    icon: Files,
    type: 'code'
  }
])

// 计算属性
const filteredFiles = computed(() => {
  let filtered = files.value

  // 搜索过滤
  if (searchQuery.value) {
    filtered = filtered.filter(file => 
      file.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      (file.description && file.description.toLowerCase().includes(searchQuery.value.toLowerCase()))
    )
  }

  // 分类过滤
  if (categoryFilter.value) {
    filtered = filtered.filter(file => file.category === categoryFilter.value)
  }

  // 排序
  filtered.sort((a, b) => {
    switch (sortBy.value) {
      case 'name':
        return a.name.localeCompare(b.name)
      case 'size':
        return b.size - a.size
      case 'type':
        return a.category.localeCompare(b.category)
      case 'uploadTime':
      default:
        return new Date(b.uploadTime) - new Date(a.uploadTime)
    }
  })

  return filtered
})

// 加载项目列表
const loadProjects = async () => {
  try {
    const response = await studentService.getMyProjects()
    if (response && response.code === 200) {
      projectList.value = response.data || []
    } else {
      // 使用模拟数据
      projectList.value = [
        { id: 1, name: '智能校园系统' },
        { id: 2, name: '数据分析平台' },
        { id: 3, name: '在线教育平台' }
      ]
    }
  } catch (error) {
    console.error('加载项目列表失败:', error)
    // 使用模拟数据
    projectList.value = [
      { id: 1, name: '智能校园系统' },
      { id: 2, name: '数据分析平台' },
      { id: 3, name: '在线教育平台' }
    ]
  }
}

// 加载文件
const loadFiles = async () => {
  if (!selectedProjectId.value) return
  
  loading.value = true
  try {
    // 这里应该调用实际的API
    // const response = await studentService.getProjectFiles(selectedProjectId.value)
    
    // 使用模拟数据
    files.value = [
      {
        id: 1,
        name: '需求分析文档.docx',
        category: 'document',
        size: 1024 * 1024 * 2.5, // 2.5MB
        uploadTime: '2024-01-20 10:30:00',
        uploader: '张三',
        description: '项目需求分析详细文档',
        url: '/api/files/1',
        type: 'docx'
      },
      {
        id: 2,
        name: '系统架构图.png',
        category: 'image',
        size: 1024 * 512, // 512KB
        uploadTime: '2024-01-22 14:20:00',
        uploader: '李四',
        description: '系统整体架构设计图',
        url: '/api/files/2',
        type: 'png'
      },
      {
        id: 3,
        name: '数据库设计.sql',
        category: 'code',
        size: 1024 * 15, // 15KB
        uploadTime: '2024-01-25 09:15:00',
        uploader: '王五',
        description: '数据库表结构设计SQL脚本',
        url: '/api/files/3',
        type: 'sql'
      },
      {
        id: 4,
        name: '项目计划书.pdf',
        category: 'document',
        size: 1024 * 1024 * 1.8, // 1.8MB
        uploadTime: '2024-01-28 16:45:00',
        uploader: '赵六',
        description: '项目整体实施计划书',
        url: '/api/files/4',
        type: 'pdf'
      }
    ]
  } catch (error) {
    console.error('加载文件失败:', error)
    ElMessage.error('加载文件失败')
  } finally {
    loading.value = false
  }
}

// 显示上传对话框
const showUploadDialog = () => {
  if (!selectedProjectId.value) {
    ElMessage.warning('请先选择一个项目')
    return
  }
  
  uploadForm.value = {
    files: [],
    category: 'document',
    description: ''
  }
  uploadDialogVisible.value = true
}

// 文件选择变化
const handleFileChange = (file, fileList) => {
  uploadForm.value.files = fileList
}

// 文件移除
const handleFileRemove = (file, fileList) => {
  uploadForm.value.files = fileList
}

// 提交上传
const submitUpload = async () => {
  try {
    await uploadFormRef.value.validate()
    
    if (uploadForm.value.files.length === 0) {
      ElMessage.warning('请选择要上传的文件')
      return
    }
    
    uploading.value = true
    
    // 这里应该调用实际的API
    // await studentService.uploadFiles(selectedProjectId.value, uploadForm.value)
    
    // 模拟上传过程
    await new Promise(resolve => setTimeout(resolve, 2000))
    
    // 添加新文件到列表
    uploadForm.value.files.forEach((file, index) => {
      const newFile = {
        id: Date.now() + index,
        name: file.name,
        category: uploadForm.value.category,
        size: file.size || 1024 * 1024,
        uploadTime: new Date().toLocaleString('zh-CN'),
        uploader: '当前用户',
        description: uploadForm.value.description,
        url: `/api/files/${Date.now() + index}`,
        type: file.name.split('.').pop()
      }
      files.value.unshift(newFile)
    })
    
    ElMessage.success('文件上传成功')
    uploadDialogVisible.value = false
  } catch (error) {
    console.error('上传失败:', error)
    ElMessage.error('上传失败')
  } finally {
    uploading.value = false
  }
}

// 预览文件
const previewFile = (file) => {
  currentFile.value = file
  previewDialogVisible.value = true
}

// 下载文件
const downloadFile = (file) => {
  // 这里应该调用实际的下载API
  ElMessage.success(`开始下载文件: ${file.name}`)
  
  // 模拟下载
  const link = document.createElement('a')
  link.href = file.url || '#'
  link.download = file.name
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

// 编辑文件
const editFile = (file) => {
  ElMessage.info(`编辑文件: ${file.name}`)
}

// 删除文件
const deleteFile = async (file) => {
  try {
    await ElMessageBox.confirm('确定要删除这个文件吗？', '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    // 这里应该调用实际的API
    // await studentService.deleteFile(file.id)
    
    const index = files.value.findIndex(f => f.id === file.id)
    if (index > -1) {
      files.value.splice(index, 1)
    }
    
    ElMessage.success('文件删除成功')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除文件失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 批量删除
const batchDelete = async () => {
  if (selectedFiles.value.length === 0) {
    ElMessage.warning('请选择要删除的文件')
    return
  }
  
  try {
    await ElMessageBox.confirm(`确定要删除选中的 ${selectedFiles.value.length} 个文件吗？`, '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    // 这里应该调用实际的API
    // await studentService.batchDeleteFiles(selectedFiles.value)
    
    // 从列表中移除选中的文件
    selectedFiles.value.forEach(fileId => {
      const index = files.value.findIndex(f => f.id === fileId)
      if (index > -1) {
        files.value.splice(index, 1)
      }
    })
    
    selectedFiles.value = []
    ElMessage.success('批量删除成功')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('批量删除失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 文件选择变化
const handleSelectionChange = (selection) => {
  selectedFiles.value = selection.map(item => item.id)
}

// 切换文件选择
const toggleFileSelection = (fileId, checked) => {
  if (checked) {
    if (!selectedFiles.value.includes(fileId)) {
      selectedFiles.value.push(fileId)
    }
  } else {
    const index = selectedFiles.value.indexOf(fileId)
    if (index > -1) {
      selectedFiles.value.splice(index, 1)
    }
  }
}

// 文件操作处理
const handleFileAction = (command) => {
  const { action, file } = command
  switch (action) {
    case 'preview':
      previewFile(file)
      break
    case 'download':
      downloadFile(file)
      break
    case 'edit':
      editFile(file)
      break
    case 'delete':
      deleteFile(file)
      break
  }
}

// 搜索处理
const handleSearch = () => {
  // 搜索逻辑已在计算属性中处理
}

// 筛选处理
const handleFilter = () => {
  // 筛选逻辑已在计算属性中处理
}

// 排序处理
const handleSort = () => {
  // 排序逻辑已在计算属性中处理
}

// 获取文件图标
const getFileIcon = (type) => {
  const iconMap = {
    'doc': Document,
    'docx': Document,
    'pdf': Document,
    'txt': Document,
    'png': Picture,
    'jpg': Picture,
    'jpeg': Picture,
    'gif': Picture,
    'mp4': VideoPlay,
    'avi': VideoPlay,
    'sql': Files,
    'js': Files,
    'py': Files,
    'java': Files
  }
  return iconMap[type] || Document
}

// 获取文件图标样式
const getFileIconClass = (type) => {
  const classMap = {
    'doc': 'file-doc',
    'docx': 'file-doc',
    'pdf': 'file-pdf',
    'png': 'file-image',
    'jpg': 'file-image',
    'jpeg': 'file-image',
    'gif': 'file-image',
    'mp4': 'file-video',
    'sql': 'file-code',
    'js': 'file-code',
    'py': 'file-code',
    'java': 'file-code'
  }
  return classMap[type] || 'file-default'
}

// 获取分类类型
const getCategoryType = (category) => {
  const typeMap = {
    'document': 'primary',
    'image': 'success',
    'video': 'warning',
    'code': 'info',
    'other': 'default'
  }
  return typeMap[category] || 'default'
}

// 获取分类标签
const getCategoryLabel = (category) => {
  const labelMap = {
    'document': '文档',
    'image': '图片',
    'video': '视频',
    'code': '代码',
    'other': '其他'
  }
  return labelMap[category] || '其他'
}

// 格式化文件大小
const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 判断是否为图片文件
const isImageFile = (type) => {
  return ['png', 'jpg', 'jpeg', 'gif', 'bmp', 'webp'].includes(type)
}

// 判断是否为文本文件
const isTextFile = (type) => {
  return ['txt', 'md', 'js', 'py', 'java', 'sql', 'html', 'css'].includes(type)
}

// 组件挂载时加载数据
onMounted(() => {
  loadProjects()
})
</script>

<style scoped>
.project-files {
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

.stat-card.document {
  border-left: 4px solid #f093fb;
}

.stat-card.image {
  border-left: 4px solid #4facfe;
}

.stat-card.code {
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

.stat-icon.document {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.stat-icon.image {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.stat-icon.code {
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

.file-actions {
  margin-bottom: 20px;
  padding: 20px;
  background-color: #f5f7fa;
  border-radius: 8px;
}

.files-content {
  margin-top: 20px;
}

.files-grid {
  margin-top: 20px;
}

.file-card {
  margin-bottom: 20px;
  transition: all 0.3s ease;
  cursor: pointer;
}

.file-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.1);
}

.file-card.selected {
  border: 2px solid #409EFF;
  background-color: #f0f9ff;
}

.file-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.file-card-content {
  text-align: center;
}

.file-icon-large {
  font-size: 48px;
  color: #409EFF;
  margin-bottom: 15px;
}

.file-title {
  margin: 0 0 10px 0;
  color: #2c3e50;
  font-size: 14px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-info {
  margin: 5px 0;
  color: #606266;
  font-size: 12px;
}

.file-uploader {
  margin: 5px 0;
  color: #909399;
  font-size: 12px;
}

.file-time {
  margin: 5px 0;
  color: #c0c4cc;
  font-size: 11px;
}

.file-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.file-icon {
  font-size: 18px;
}

.file-link {
  color: #409EFF;
  cursor: pointer;
}

.file-link:hover {
  text-decoration: underline;
}

.file-preview {
  padding: 20px;
}

.file-info-header h3 {
  margin: 0 0 15px 0;
  color: #2c3e50;
}

.file-meta {
  display: flex;
  gap: 15px;
  align-items: center;
  color: #606266;
  font-size: 14px;
}

.file-description {
  margin: 20px 0;
}

.file-description h4 {
  margin: 0 0 10px 0;
  color: #2c3e50;
}

.preview-content {
  margin-top: 20px;
  text-align: center;
}

.image-preview img {
  max-width: 100%;
  max-height: 500px;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.text-preview {
  text-align: left;
  background-color: #f5f7fa;
  padding: 20px;
  border-radius: 8px;
  max-height: 400px;
  overflow-y: auto;
}

.text-preview pre {
  margin: 0;
  white-space: pre-wrap;
  word-wrap: break-word;
}

.no-preview {
  padding: 40px;
  color: #909399;
}

.no-preview .el-icon {
  font-size: 48px;
  margin-bottom: 15px;
}

.dialog-footer {
  text-align: right;
}

/* 文件图标样式 */
.file-doc {
  color: #409EFF;
}

.file-pdf {
  color: #F56C6C;
}

.file-image {
  color: #67C23A;
}

.file-video {
  color: #E6A23C;
}

.file-code {
  color: #909399;
}

.file-default {
  color: #C0C4CC;
}
</style> 
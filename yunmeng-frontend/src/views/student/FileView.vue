<template>
  <div class="file-view">
    <!-- 页面标题 -->
    <div class="page-header">
      <h2>文件管理</h2>
      <p>管理您的个人文件和项目附件</p>
    </div>

    <!-- 操作栏 -->
    <div class="action-bar">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-select v-model="filterCategory" placeholder="文件分类" @change="loadFiles">
            <el-option label="全部" value=""></el-option>
            <el-option label="项目文件" value="project"></el-option>
            <el-option label="竞赛文件" value="competition"></el-option>
            <el-option label="个人文档" value="document"></el-option>
            <el-option label="图片资料" value="image"></el-option>
            <el-option label="其他" value="other"></el-option>
          </el-select>
        </el-col>
        <el-col :span="6">
          <el-select v-model="filterType" placeholder="文件类型" @change="loadFiles">
            <el-option label="全部" value=""></el-option>
            <el-option label="PDF" value="pdf"></el-option>
            <el-option label="Word" value="doc"></el-option>
            <el-option label="Excel" value="xls"></el-option>
            <el-option label="PPT" value="ppt"></el-option>
            <el-option label="图片" value="image"></el-option>
            <el-option label="压缩包" value="zip"></el-option>
          </el-select>
        </el-col>
        <el-col :span="8">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索文件名"
            @input="handleSearch"
            clearable
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-col>
        <el-col :span="4">
          <el-button type="primary" @click="loadFiles">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </el-col>
      </el-row>
      
      <div class="upload-section" style="margin-top: 15px;">
        <el-upload
          ref="uploadRef"
          :action="uploadUrl"
          :headers="uploadHeaders"
          :data="uploadData"
          :before-upload="beforeUpload"
          :on-success="onUploadSuccess"
          :on-error="onUploadError"
          :on-progress="onUploadProgress"
          :file-list="uploadFileList"
          :auto-upload="false"
          multiple
          drag
        >
                          <el-icon class="el-icon--upload"><Upload /></el-icon>
          <div class="el-upload__text">
            将文件拖到此处，或<em>点击上传</em>
          </div>
          <template #tip>
            <div class="el-upload__tip">
              支持 PDF、Word、Excel、PPT、图片等格式，单个文件不超过10MB
            </div>
          </template>
        </el-upload>
        
        <div class="upload-actions" style="margin-top: 10px;">
          <el-button type="primary" @click="submitUpload" :loading="uploading">
            开始上传
          </el-button>
          <el-button @click="clearUpload">清空列表</el-button>
        </div>
      </div>
    </div>

    <!-- 文件列表 -->
    <div class="file-list">
      <el-table 
        :data="filteredFiles" 
        v-loading="loading"
        style="width: 100%"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" />
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="fileName" label="文件名" min-width="200">
          <template #default="{ row }">
            <div class="file-name">
              <el-icon class="file-icon" :class="getFileIconClass(row.fileType)">
                <component :is="getFileIcon(row.fileType)" />
              </el-icon>
              <span>{{ row.fileName }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="category" label="分类" width="120">
          <template #default="{ row }">
            <el-tag :type="getCategoryType(row.category)" size="small">
              {{ getCategoryText(row.category) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="fileSize" label="大小" width="120" />
        <el-table-column prop="uploadTime" label="上传时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.uploadTime) }}
          </template>
        </el-table-column>
        <el-table-column prop="downloadCount" label="下载次数" width="100" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="downloadFile(row)">下载</el-button>
            <el-button size="small" @click="openFilePreview(row)">预览</el-button>
            <el-button size="small" @click="editFile(row)">编辑</el-button>
            <el-button size="small" type="danger" @click="deleteFile(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <!-- 空状态 -->
    <el-empty 
      v-if="!loading && filteredFiles.length === 0" 
      description="暂无文件"
      style="margin-top: 40px;"
    />

    <!-- 批量操作 -->
    <div v-if="selectedFiles.length > 0" class="bulk-actions">
      <el-card>
        <div class="bulk-actions-content">
          <span>已选择 {{ selectedFiles.length }} 个文件</span>
          <div class="bulk-buttons">
            <el-button size="small" @click="downloadSelected">批量下载</el-button>
            <el-button size="small" @click="moveSelected">移动到</el-button>
            <el-button size="small" type="danger" @click="deleteSelected">批量删除</el-button>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 文件编辑对话框 -->
    <el-dialog
      v-model="showEditDialog"
      title="编辑文件信息"
      width="50%"
    >
      <el-form :model="editForm" :rules="editRules" ref="editFormRef" label-width="100px">
        <el-form-item label="文件名" prop="fileName">
          <el-input v-model="editForm.fileName" placeholder="请输入文件名" />
        </el-form-item>
        
        <el-form-item label="文件分类" prop="category">
          <el-select v-model="editForm.category" placeholder="请选择分类" style="width: 100%">
            <el-option label="项目文件" value="project" />
            <el-option label="竞赛文件" value="competition" />
            <el-option label="个人文档" value="document" />
            <el-option label="图片资料" value="image" />
            <el-option label="其他" value="other" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="文件描述" prop="description">
          <el-input 
            v-model="editForm.description" 
            type="textarea" 
            placeholder="请输入文件描述"
            :rows="3"
          />
        </el-form-item>
        
        <el-form-item label="标签" prop="tags">
          <el-select
            v-model="editForm.tags"
            multiple
            filterable
            allow-create
            default-first-option
            placeholder="请选择或输入标签"
            style="width: 100%"
          >
            <el-option label="重要" value="important" />
            <el-option label="项目相关" value="project" />
            <el-option label="竞赛相关" value="competition" />
            <el-option label="参考资料" value="reference" />
          </el-select>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showEditDialog = false">取消</el-button>
          <el-button type="primary" @click="saveFileEdit" :loading="saving">
            保存
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 文件预览对话框 -->
    <el-dialog
      v-model="showPreviewDialog"
      title="文件预览"
      width="80%"
      :before-close="handleClosePreview"
    >
      <div v-if="previewFile" class="file-preview">
        <div class="preview-header">
          <h3>{{ previewFile.fileName }}</h3>
          <div class="preview-info">
            <span>大小：{{ previewFile.fileSize }}</span>
            <span>上传时间：{{ formatDate(previewFile.uploadTime) }}</span>
            <span>下载次数：{{ previewFile.downloadCount }}</span>
          </div>
        </div>
        
        <el-divider />
        
        <div class="preview-content">
          <!-- 图片预览 -->
          <div v-if="isImageFile(previewFile.fileType)" class="image-preview">
            <img :src="previewFile.fileUrl" :alt="previewFile.fileName" />
          </div>
          
          <!-- PDF预览 -->
          <div v-else-if="previewFile.fileType === 'pdf'" class="pdf-preview">
            <iframe :src="previewFile.fileUrl" width="100%" height="600"></iframe>
          </div>
          
          <!-- 其他文件 -->
          <div v-else class="other-preview">
            <el-empty description="该文件类型暂不支持预览">
              <el-button type="primary" @click="downloadFile(previewFile)">
                下载文件
              </el-button>
            </el-empty>
          </div>
        </div>
        
        <div class="preview-actions">
          <el-button @click="showPreviewDialog = false">关闭</el-button>
          <el-button type="primary" @click="downloadFile(previewFile)">
            下载文件
          </el-button>
        </div>
      </div>
    </el-dialog>

    <!-- 移动文件对话框 -->
    <el-dialog
      v-model="showMoveDialog"
      title="移动文件"
      width="40%"
    >
      <div class="move-dialog">
        <p>将选中的 {{ selectedFiles.length }} 个文件移动到：</p>
        
        <el-form :model="moveForm" label-width="100px">
          <el-form-item label="目标分类">
            <el-select v-model="moveForm.targetCategory" placeholder="请选择目标分类" style="width: 100%">
              <el-option label="项目文件" value="project" />
              <el-option label="竞赛文件" value="competition" />
              <el-option label="个人文档" value="document" />
              <el-option label="图片资料" value="image" />
              <el-option label="其他" value="other" />
            </el-select>
          </el-form-item>
        </el-form>
      </div>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showMoveDialog = false">取消</el-button>
          <el-button type="primary" @click="confirmMove" :loading="moving">
            确认移动
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Search, 
  Refresh, 
  Upload,
  Document,
  Picture,
  VideoPlay,
  Folder,
  Files
} from '@element-plus/icons-vue'

// 响应式数据
const loading = ref(false)
const uploading = ref(false)
const saving = ref(false)
const moving = ref(false)
const filterCategory = ref('')
const filterType = ref('')
const searchKeyword = ref('')
const files = ref([])
const selectedFiles = ref([])
const showEditDialog = ref(false)
const showPreviewDialog = ref(false)
const showMoveDialog = ref(false)
const previewFile = ref(null)
const uploadRef = ref()
const uploadFileList = ref([])

// 上传相关
const uploadUrl = 'http://localhost:8080/api/files/upload'
const uploadHeaders = {
  Authorization: `Bearer ${localStorage.getItem('token')}`
}
const uploadData = {
  category: 'document'
}

// 编辑表单
const editForm = ref({
  fileName: '',
  category: '',
  description: '',
  tags: []
})

// 移动表单
const moveForm = ref({
  targetCategory: ''
})

// 表单引用
const editFormRef = ref()

// 编辑验证规则
const editRules = {
  fileName: [
    { required: true, message: '请输入文件名', trigger: 'blur' }
  ],
  category: [
    { required: true, message: '请选择文件分类', trigger: 'change' }
  ]
}

// 模拟文件数据
const mockFiles = [
  {
    id: 1,
    fileName: '项目申请书.pdf',
    fileType: 'pdf',
    category: 'project',
    fileSize: '2.5MB',
    uploadTime: '2024-01-15T10:00:00Z',
    downloadCount: 5,
    fileUrl: '/uploads/project_application.pdf',
    description: '智能校园管理系统项目申请书',
    tags: ['重要', '项目相关']
  },
  {
    id: 2,
    fileName: '技术方案.docx',
    fileType: 'doc',
    category: 'project',
    fileSize: '1.8MB',
    uploadTime: '2024-01-15T10:00:00Z',
    downloadCount: 3,
    fileUrl: '/uploads/technical_proposal.docx',
    description: '项目技术实现方案',
    tags: ['项目相关']
  },
  {
    id: 3,
    fileName: '竞赛报名表.pdf',
    fileType: 'pdf',
    category: 'competition',
    fileSize: '500KB',
    uploadTime: '2024-01-10T14:30:00Z',
    downloadCount: 2,
    fileUrl: '/uploads/competition_form.pdf',
    description: '程序设计竞赛报名表',
    tags: ['竞赛相关']
  },
  {
    id: 4,
    fileName: '项目截图.png',
    fileType: 'png',
    category: 'image',
    fileSize: '800KB',
    uploadTime: '2024-01-12T16:20:00Z',
    downloadCount: 8,
    fileUrl: '/uploads/project_screenshot.png',
    description: '项目界面截图',
    tags: ['图片资料']
  },
  {
    id: 5,
    fileName: '参考资料.zip',
    fileType: 'zip',
    category: 'document',
    fileSize: '5.2MB',
    uploadTime: '2024-01-08T09:15:00Z',
    downloadCount: 12,
    fileUrl: '/uploads/reference_materials.zip',
    description: '项目参考资料压缩包',
    tags: ['参考资料']
  }
]

// 过滤后的文件列表
const filteredFiles = computed(() => {
  let filtered = files.value

  // 分类筛选
  if (filterCategory.value) {
    filtered = filtered.filter(f => f.category === filterCategory.value)
  }

  // 类型筛选
  if (filterType.value) {
    filtered = filtered.filter(f => f.fileType === filterType.value)
  }

  // 关键词搜索
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    filtered = filtered.filter(f => 
      f.fileName.toLowerCase().includes(keyword) || 
      f.description.toLowerCase().includes(keyword)
    )
  }

  return filtered
})

// 加载文件数据
const loadFiles = async () => {
  loading.value = true
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500))
    files.value = mockFiles
  } catch (error) {
    console.error('加载文件数据失败:', error)
    ElMessage.error('加载文件数据失败')
  } finally {
    loading.value = false
  }
}

// 获取文件图标
const getFileIcon = (fileType) => {
  const iconMap = {
    pdf: Document,
    doc: Document,
    docx: Document,
    xls: Document,
    xlsx: Document,
    ppt: Document,
    pptx: Document,
    png: Picture,
    jpg: Picture,
    jpeg: Picture,
    gif: Picture,
    mp4: VideoPlay,
    avi: VideoPlay,
    zip: Folder,
    rar: Folder
  }
  return iconMap[fileType] || Files
}

// 获取文件图标样式
const getFileIconClass = (fileType) => {
  const classMap = {
    pdf: 'file-pdf',
    doc: 'file-word',
    docx: 'file-word',
    xls: 'file-excel',
    xlsx: 'file-excel',
    ppt: 'file-ppt',
    pptx: 'file-ppt',
    png: 'file-image',
    jpg: 'file-image',
    jpeg: 'file-image',
    gif: 'file-image',
    mp4: 'file-video',
    avi: 'file-video',
    zip: 'file-archive',
    rar: 'file-archive'
  }
  return classMap[fileType] || 'file-default'
}

// 分类类型映射
const getCategoryType = (category) => {
  const typeMap = {
    project: 'primary',
    competition: 'success',
    document: 'info',
    image: 'warning',
    other: 'default'
  }
  return typeMap[category] || 'default'
}

// 分类文本映射
const getCategoryText = (category) => {
  const textMap = {
    project: '项目文件',
    competition: '竞赛文件',
    document: '个人文档',
    image: '图片资料',
    other: '其他'
  }
  return textMap[category] || category
}

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

// 搜索处理
const handleSearch = () => {
  // 实时搜索，不需要额外处理，computed会自动更新
}

// 选择变化处理
const handleSelectionChange = (selection) => {
  selectedFiles.value = selection
}

// 上传前检查
const beforeUpload = (file) => {
  const isValidType = [
    'application/pdf',
    'application/msword',
    'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
    'application/vnd.ms-excel',
    'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
    'application/vnd.ms-powerpoint',
    'application/vnd.openxmlformats-officedocument.presentationml.presentation',
    'image/jpeg',
    'image/png',
    'image/gif',
    'application/zip',
    'application/x-rar-compressed'
  ].includes(file.type)
  
  const isLt10M = file.size / 1024 / 1024 < 10

  if (!isValidType) {
    ElMessage.error('不支持的文件类型！')
    return false
  }
  if (!isLt10M) {
    ElMessage.error('文件大小不能超过 10MB！')
    return false
  }
  return true
}

// 上传成功
const onUploadSuccess = (response, file, fileList) => {
  ElMessage.success(`${file.name} 上传成功`)
  loadFiles() // 重新加载文件列表
}

// 上传失败
const onUploadError = (error, file, fileList) => {
  ElMessage.error(`${file.name} 上传失败`)
}

// 上传进度
const onUploadProgress = (event, file, fileList) => {
  console.log('上传进度:', event.percent)
}

// 提交上传
const submitUpload = () => {
  uploadRef.value.submit()
}

// 清空上传列表
const clearUpload = () => {
  uploadRef.value.clearFiles()
}

// 下载文件
const downloadFile = (file) => {
  ElMessage.info(`下载文件：${file.fileName}`)
  // 这里应该实现实际的文件下载逻辑
}

// 预览文件
const openFilePreview = (file) => {
  previewFile.value = file
  showPreviewDialog.value = true
}

// 关闭预览
const handleClosePreview = () => {
  showPreviewDialog.value = false
  previewFile.value = null
}

// 判断是否为图片文件
const isImageFile = (fileType) => {
  return ['png', 'jpg', 'jpeg', 'gif', 'bmp'].includes(fileType)
}

// 编辑文件
const editFile = (file) => {
  editForm.value = {
    fileName: file.fileName,
    category: file.category,
    description: file.description || '',
    tags: file.tags || []
  }
  showEditDialog.value = true
}

// 保存文件编辑
const saveFileEdit = async () => {
  try {
    await editFormRef.value.validate()
    saving.value = true
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 800))
    
    // 更新文件信息
    const index = files.value.findIndex(f => f.id === previewFile.value?.id)
    if (index !== -1) {
      Object.assign(files.value[index], editForm.value)
    }
    
    showEditDialog.value = false
    ElMessage.success('文件信息保存成功')
  } catch (error) {
    console.error('保存失败:', error)
    ElMessage.error('保存失败，请检查输入信息')
  } finally {
    saving.value = false
  }
}

// 删除文件
const deleteFile = async (file) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除文件"${file.fileName}"吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 300))
    
    const index = files.value.findIndex(f => f.id === file.id)
    if (index !== -1) {
      files.value.splice(index, 1)
      ElMessage.success('删除成功')
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 批量下载
const downloadSelected = () => {
  ElMessage.info(`批量下载 ${selectedFiles.value.length} 个文件`)
}

// 移动文件
const moveSelected = () => {
  moveForm.value.targetCategory = ''
  showMoveDialog.value = true
}

// 确认移动
const confirmMove = async () => {
  if (!moveForm.value.targetCategory) {
    ElMessage.warning('请选择目标分类')
    return
  }

  moving.value = true
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // 更新文件分类
    selectedFiles.value.forEach(file => {
      const index = files.value.findIndex(f => f.id === file.id)
      if (index !== -1) {
        files.value[index].category = moveForm.value.targetCategory
      }
    })
    
    showMoveDialog.value = false
    ElMessage.success('文件移动成功')
  } catch (error) {
    console.error('移动失败:', error)
    ElMessage.error('移动失败')
  } finally {
    moving.value = false
  }
}

// 批量删除
const deleteSelected = async () => {
  try {
    await ElMessageBox.confirm(
      `确定要删除选中的 ${selectedFiles.value.length} 个文件吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500))
    
    const selectedIds = selectedFiles.value.map(f => f.id)
    files.value = files.value.filter(f => !selectedIds.includes(f.id))
    
    ElMessage.success('删除成功')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 组件挂载时加载数据
onMounted(() => {
  loadFiles()
})
</script>

<style scoped>
.file-view {
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

.action-bar {
  margin-bottom: 30px;
  padding: 20px;
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.upload-section {
  border-top: 1px solid #e9ecef;
  padding-top: 20px;
}

.upload-actions {
  display: flex;
  gap: 10px;
}

.file-list {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  margin-bottom: 20px;
}

.el-table {
  border-radius: 8px;
}

.file-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.file-icon {
  font-size: 16px;
}

.file-pdf { color: #f56c6c; }
.file-word { color: #409eff; }
.file-excel { color: #67c23a; }
.file-ppt { color: #e6a23c; }
.file-image { color: #909399; }
.file-video { color: #f56c6c; }
.file-archive { color: #909399; }
.file-default { color: #909399; }

.bulk-actions {
  margin-top: 20px;
}

.bulk-actions-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.bulk-buttons {
  display: flex;
  gap: 10px;
}

.file-preview {
  padding: 20px 0;
}

.preview-header {
  margin-bottom: 20px;
}

.preview-header h3 {
  margin: 0 0 10px 0;
  color: #2c3e50;
}

.preview-info {
  display: flex;
  gap: 20px;
  color: #7f8c8d;
  font-size: 14px;
}

.preview-content {
  margin: 20px 0;
}

.image-preview {
  text-align: center;
}

.image-preview img {
  max-width: 100%;
  max-height: 500px;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.pdf-preview {
  border: 1px solid #e9ecef;
  border-radius: 8px;
  overflow: hidden;
}

.other-preview {
  text-align: center;
  padding: 40px 0;
}

.preview-actions {
  display: flex;
  justify-content: flex-end;
  gap: 15px;
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #e9ecef;
}

.move-dialog {
  padding: 20px 0;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style> 
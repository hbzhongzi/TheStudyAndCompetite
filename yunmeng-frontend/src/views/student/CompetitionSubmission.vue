<template>
  <div class="competition-submission">
    <!-- 页面标题 -->
    <div class="page-header">
      <h2>竞赛作品提交</h2>
      <p>提交和管理您的竞赛作品</p>
    </div>

    <!-- 竞赛信息卡片 -->
    <el-card v-if="currentCompetition" class="competition-info-card">
      <template #header>
        <div class="competition-header">
          <h3>{{ currentCompetition.title }}</h3>
          <el-tag :type="getStatusType(currentCompetition.status)">
            {{ getStatusText(currentCompetition.status) }}
          </el-tag>
        </div>
      </template>
      
      <el-descriptions :column="3" border>
        <el-descriptions-item label="竞赛类型">
          {{ currentCompetition.type }}
        </el-descriptions-item>
        <el-descriptions-item label="主办方">
          {{ currentCompetition.organizer }}
        </el-descriptions-item>
        <el-descriptions-item label="提交截止时间">
          {{ formatDate(currentCompetition.submissionEnd) }}
        </el-descriptions-item>
        <el-descriptions-item label="支持格式">
          {{ currentCompetition.fileFormats || 'PDF、DOC、ZIP等' }}
        </el-descriptions-item>
        <el-descriptions-item label="文件大小限制">
          {{ currentCompetition.fileSizeLimit || '单个文件不超过50MB' }}
        </el-descriptions-item>
        <el-descriptions-item label="当前状态">
          <el-tag :type="getSubmissionStatusType(submissionStatus)">
            {{ getSubmissionStatusText(submissionStatus) }}
          </el-tag>
        </el-descriptions-item>
      </el-descriptions>
    </el-card>

    <!-- 作品提交表单 -->
    <el-card class="submission-form-card">
      <template #header>
        <h3>提交作品</h3>
      </template>
      
      <el-form
        ref="submissionFormRef"
        :model="submissionForm"
        :rules="submissionRules"
        label-width="120px"
      >
        <el-form-item label="作品标题" prop="title">
          <el-input 
            v-model="submissionForm.title" 
            placeholder="请输入作品标题"
            maxlength="100"
            show-word-limit
          />
        </el-form-item>
        
        <el-form-item label="作品描述" prop="description">
          <el-input
            v-model="submissionForm.description"
            type="textarea"
            :rows="4"
            placeholder="请详细描述您的作品内容、创新点、技术特点等"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>
        
        <el-form-item label="作品文件" prop="file">
          <el-upload
            ref="uploadRef"
            :action="uploadAction"
            :before-upload="beforeUpload"
            :on-success="handleUploadSuccess"
            :on-error="handleUploadError"
            :on-remove="handleFileRemove"
            :file-list="fileList"
            :limit="1"
            :accept="acceptedFileTypes"
            drag
          >
            <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
            <div class="el-upload__text">
              将文件拖到此处，或<em>点击上传</em>
            </div>
            <template #tip>
              <div class="el-upload__tip">
                支持 {{ acceptedFileTypes }} 格式，单个文件不超过 {{ maxFileSize }}MB
              </div>
            </template>
          </el-upload>
        </el-form-item>
        
        <el-form-item label="版本说明" prop="versionNote">
          <el-input
            v-model="submissionForm.versionNote"
            placeholder="请说明本次提交的版本更新内容（可选）"
            maxlength="200"
            show-word-limit
          />
        </el-form-item>
        
        <el-form-item>
          <el-button 
            type="primary" 
            @click="submitWork"
            :loading="submitting"
            :disabled="!canSubmit"
          >
            提交作品
          </el-button>
          <el-button @click="saveAsDraft">保存草稿</el-button>
          <el-button @click="resetForm">重置表单</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 提交历史 -->
    <el-card class="submission-history-card">
      <template #header>
        <h3>提交历史</h3>
      </template>
      
      <el-table :data="submissionHistory" style="width: 100%">
        <el-table-column prop="version" label="版本" width="80" />
        <el-table-column prop="title" label="作品标题" width="200" />
        <el-table-column prop="submitTime" label="提交时间" width="150">
          <template #default="scope">
            {{ formatDate(scope.row.submitTime) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getSubmissionStatusType(scope.row.status)">
              {{ getSubmissionStatusText(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="fileSize" label="文件大小" width="100">
          <template #default="scope">
            {{ formatFileSize(scope.row.fileSize) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="scope">
            <div class="action-buttons">
              <el-button 
                size="small" 
                type="primary" 
                @click="viewSubmission(scope.row)"
              >
                查看
              </el-button>
              <el-button 
                v-if="scope.row.status === 'submitted' && !scope.row.locked"
                size="small" 
                type="warning" 
                @click="editSubmission(scope.row)"
              >
                编辑
              </el-button>
              <el-button 
                v-if="scope.row.status === 'submitted' && !scope.row.locked"
                size="small" 
                type="danger" 
                @click="deleteSubmission(scope.row)"
              >
                删除
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 作品详情对话框 -->
    <el-dialog
      v-model="showDetailDialog"
      title="作品详情"
      width="70%"
    >
      <div v-if="selectedSubmission" class="submission-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="作品标题">
            {{ selectedSubmission.title }}
          </el-descriptions-item>
          <el-descriptions-item label="版本号">
            {{ selectedSubmission.version }}
          </el-descriptions-item>
          <el-descriptions-item label="提交时间">
            {{ formatDate(selectedSubmission.submitTime) }}
          </el-descriptions-item>
          <el-descriptions-item label="文件大小">
            {{ formatFileSize(selectedSubmission.fileSize) }}
          </el-descriptions-item>
          <el-descriptions-item label="当前状态">
            <el-tag :type="getSubmissionStatusType(selectedSubmission.status)">
              {{ getSubmissionStatusText(selectedSubmission.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="是否锁定">
            <el-tag :type="selectedSubmission.locked ? 'danger' : 'success'">
              {{ selectedSubmission.locked ? '已锁定' : '未锁定' }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>
        
        <div class="description-section">
          <h4>作品描述</h4>
          <p>{{ selectedSubmission.description }}</p>
        </div>
        
        <div class="file-section">
          <h4>作品文件</h4>
          <el-button 
            type="primary" 
            @click="downloadFile(selectedSubmission)"
            :icon="Download"
          >
            下载文件
          </el-button>
        </div>
        
        <div v-if="selectedSubmission.reviewComments" class="review-section">
          <h4>评审意见</h4>
          <p>{{ selectedSubmission.reviewComments }}</p>
        </div>
      </div>
    </el-dialog>

    <!-- 编辑作品对话框 -->
    <el-dialog
      v-model="showEditDialog"
      title="编辑作品"
      width="60%"
    >
      <el-form
        ref="editFormRef"
        :model="editForm"
        :rules="submissionRules"
        label-width="120px"
      >
        <el-form-item label="作品标题" prop="title">
          <el-input v-model="editForm.title" />
        </el-form-item>
        
        <el-form-item label="作品描述" prop="description">
          <el-input
            v-model="editForm.description"
            type="textarea"
            :rows="4"
          />
        </el-form-item>
        
        <el-form-item label="版本说明" prop="versionNote">
          <el-input v-model="editForm.versionNote" />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showEditDialog = false">取消</el-button>
          <el-button type="primary" @click="confirmEdit" :loading="editing">
            确认修改
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { UploadFilled, Download } from '@element-plus/icons-vue'
import { formatDate } from '@/utils/dateUtils'

// 响应式数据
const currentCompetition = ref(null)
const submissionForm = ref({
  title: '',
  description: '',
  file: null,
  versionNote: ''
})
const editForm = ref({
  title: '',
  description: '',
  versionNote: ''
})
const fileList = ref([])
const submissionHistory = ref([])
const selectedSubmission = ref(null)
const showDetailDialog = ref(false)
const showEditDialog = ref(false)
const submitting = ref(false)
const editing = ref(false)

// 表单引用
const submissionFormRef = ref()
const editFormRef = ref()
const uploadRef = ref()

// 上传配置
const uploadAction = '/api/files/upload' // 实际的上传地址
const acceptedFileTypes = '.pdf,.doc,.docx,.zip,.rar,.ppt,.pptx'
const maxFileSize = 50 // MB

// 表单验证规则
const submissionRules = {
  title: [
    { required: true, message: '请输入作品标题', trigger: 'blur' },
    { min: 2, max: 100, message: '标题长度在 2 到 100 个字符', trigger: 'blur' }
  ],
  description: [
    { required: true, message: '请输入作品描述', trigger: 'blur' },
    { min: 10, max: 500, message: '描述长度在 10 到 500 个字符', trigger: 'blur' }
  ],
  file: [
    { required: true, message: '请选择要上传的文件', trigger: 'change' }
  ]
}

// 计算属性
const submissionStatus = computed(() => {
  if (submissionHistory.value.length === 0) return 'not_submitted'
  const latest = submissionHistory.value[0]
  return latest.status
})

const canSubmit = computed(() => {
  return currentCompetition.value && 
         currentCompetition.value.status === 'submission' &&
         submissionForm.value.title &&
         submissionForm.value.description &&
         fileList.value.length > 0
})

// 方法
const loadCompetitionInfo = async () => {
  try {
    // 模拟加载竞赛信息
    currentCompetition.value = {
      id: 1,
      title: '全国大学生程序设计竞赛',
      type: '程序设计',
      organizer: '计算机学院',
      status: 'submission',
      submissionEnd: new Date('2024-03-01'),
      fileFormats: 'PDF、DOC、ZIP、RAR',
      fileSizeLimit: '单个文件不超过50MB'
    }
  } catch (error) {
    console.error('加载竞赛信息失败:', error)
    ElMessage.error('加载竞赛信息失败')
  }
}

const loadSubmissionHistory = async () => {
  try {
    // 模拟加载提交历史
    submissionHistory.value = [
      {
        id: 1,
        version: '1.0',
        title: '程序设计竞赛作品v1.0',
        description: '基于C++的程序设计作品',
        submitTime: new Date('2024-02-15'),
        status: 'submitted',
        fileSize: 1024 * 1024 * 5, // 5MB
        locked: false,
        reviewComments: null
      },
      {
        id: 2,
        version: '1.1',
        title: '程序设计竞赛作品v1.1',
        description: '优化后的程序设计作品',
        submitTime: new Date('2024-02-20'),
        status: 'reviewing',
        fileSize: 1024 * 1024 * 6, // 6MB
        locked: true,
        reviewComments: '作品质量良好，建议继续优化'
      }
    ]
  } catch (error) {
    console.error('加载提交历史失败:', error)
    ElMessage.error('加载提交历史失败')
  }
}

const beforeUpload = (file) => {
  // 检查文件大小
  const isLt50M = file.size / 1024 / 1024 < maxFileSize
  if (!isLt50M) {
    ElMessage.error(`文件大小不能超过 ${maxFileSize}MB!`)
    return false
  }
  
  // 检查文件类型
  const isValidType = acceptedFileTypes.includes(file.name.substring(file.name.lastIndexOf('.')))
  if (!isValidType) {
    ElMessage.error('文件格式不支持!')
    return false
  }
  
  return true
}

const handleUploadSuccess = (response, file) => {
  submissionForm.value.file = file
  ElMessage.success('文件上传成功')
}

const handleUploadError = (error) => {
  ElMessage.error('文件上传失败')
}

const handleFileRemove = () => {
  submissionForm.value.file = null
}

const submitWork = async () => {
  try {
    await submissionFormRef.value.validate()
    
    submitting.value = true
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // 创建新的提交记录
    const newSubmission = {
      id: Date.now(),
      version: (submissionHistory.value.length + 1).toString(),
      title: submissionForm.value.title,
      description: submissionForm.value.description,
      submitTime: new Date(),
      status: 'submitted',
      fileSize: submissionForm.value.file?.size || 0,
      locked: false,
      reviewComments: null
    }
    
    submissionHistory.value.unshift(newSubmission)
    
    // 重置表单
    resetForm()
    
    ElMessage.success('作品提交成功！')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('提交失败：' + (error.message || '未知错误'))
    }
  } finally {
    submitting.value = false
  }
}

const saveAsDraft = () => {
  ElMessage.info('草稿保存功能开发中...')
}

const resetForm = () => {
  submissionForm.value = {
    title: '',
    description: '',
    file: null,
    versionNote: ''
  }
  fileList.value = []
  submissionFormRef.value?.resetFields()
}

const viewSubmission = (submission) => {
  selectedSubmission.value = submission
  showDetailDialog.value = true
}

const editSubmission = (submission) => {
  editForm.value = {
    title: submission.title,
    description: submission.description,
    versionNote: ''
  }
  showEditDialog.value = true
}

const confirmEdit = async () => {
  try {
    await editFormRef.value.validate()
    
    editing.value = true
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500))
    
    // 更新提交记录
    const index = submissionHistory.value.findIndex(s => s.id === selectedSubmission.value.id)
    if (index !== -1) {
      submissionHistory.value[index].title = editForm.value.title
      submissionHistory.value[index].description = editForm.value.description
    }
    
    showEditDialog.value = false
    ElMessage.success('作品修改成功')
  } catch (error) {
    ElMessage.error('修改失败：' + (error.message || '未知错误'))
  } finally {
    editing.value = false
  }
}

const deleteSubmission = async (submission) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除作品 "${submission.title}" 吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500))
    
    const index = submissionHistory.value.findIndex(s => s.id === submission.id)
    if (index !== -1) {
      submissionHistory.value.splice(index, 1)
    }
    
    ElMessage.success('作品删除成功')
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败：' + (error.message || '未知错误'))
    }
  }
}

const downloadFile = (submission) => {
  ElMessage.success('开始下载文件...')
  // 实际项目中这里会调用下载API
}

// 工具方法
const getStatusType = (status) => {
  const statusMap = {
    draft: 'info',
    registration: 'warning',
    submission: 'success',
    review: 'primary',
    completed: 'success'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status) => {
  const statusMap = {
    draft: '草稿',
    registration: '报名中',
    submission: '提交中',
    review: '评审中',
    completed: '已完成'
  }
  return statusMap[status] || status
}

const getSubmissionStatusType = (status) => {
  const statusMap = {
    not_submitted: 'info',
    submitted: 'success',
    reviewing: 'warning',
    approved: 'success',
    rejected: 'danger'
  }
  return statusMap[status] || 'info'
}

const getSubmissionStatusText = (status) => {
  const statusMap = {
    not_submitted: '未提交',
    submitted: '已提交',
    reviewing: '评审中',
    approved: '已通过',
    rejected: '已拒绝'
  }
  return statusMap[status] || status
}

const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 组件挂载时加载数据
onMounted(() => {
  loadCompetitionInfo()
  loadSubmissionHistory()
})
</script>

<style scoped>
.competition-submission {
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

.competition-info-card {
  margin-bottom: 30px;
}

.competition-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.competition-header h3 {
  margin: 0;
  color: #2c3e50;
}

.submission-form-card {
  margin-bottom: 30px;
}

.submission-form-card h3 {
  margin: 0;
  color: #2c3e50;
}

.submission-history-card {
  margin-bottom: 30px;
}

.submission-history-card h3 {
  margin: 0;
  color: #2c3e50;
}

.action-buttons {
  display: flex;
  gap: 5px;
}

.description-section,
.file-section,
.review-section {
  margin-top: 20px;
}

.description-section h4,
.file-section h4,
.review-section h4 {
  margin: 0 0 10px 0;
  color: #2c3e50;
  font-size: 16px;
}

.description-section p,
.review-section p {
  margin: 0;
  color: #606266;
  line-height: 1.6;
  background: #f8f9fa;
  padding: 15px;
  border-radius: 6px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .competition-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
  
  .action-buttons {
    flex-direction: column;
  }
}
</style> 
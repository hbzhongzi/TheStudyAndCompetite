<template>
  <el-dialog
    v-model="dialogVisible"
    :title="isEdit ? '编辑项目' : '创建项目'"
    width="800px"
    :before-close="handleClose"
  >
    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="100px"
      v-loading="loading"
    >
      <el-row :gutter="20">
        <el-col :span="24">
          <el-form-item label="项目标题" prop="title">
            <el-input 
              v-model="form.title" 
              placeholder="请输入项目标题"
              maxlength="100"
              show-word-limit
            />
          </el-form-item>
        </el-col>
        
        <el-col :span="24">
          <el-form-item label="项目描述" prop="description">
            <el-input
              v-model="form.description"
              type="textarea"
              :rows="4"
              placeholder="请输入项目描述"
              maxlength="1000"
              show-word-limit
            />
          </el-form-item>
        </el-col>
        
        <el-col :span="12">
          <el-form-item label="项目类型" prop="type">
            <el-select v-model="form.type" placeholder="请选择项目类型" style="width: 100%">
              <el-option label="科研" value="科研" />
              <el-option label="竞赛" value="竞赛" />
            </el-select>
          </el-form-item>
        </el-col>
        
        <el-col :span="12">
          <el-form-item label="指导老师" prop="teacherId">
            <TeacherSelector
              v-model="form.teacherId"
              :disabled="isEdit && form.status !== 'draft'"
              @change="handleTeacherChange"
            />
          </el-form-item>
        </el-col>
        
        <el-col :span="12">
          <el-form-item label="项目状态" prop="status">
            <el-select v-model="form.status" placeholder="请选择项目状态" style="width: 100%" :disabled="true">
              <el-option label="草稿" value="draft" />
              <el-option label="待审核" value="pending" />
              <el-option label="已通过" value="approved" />
              <el-option label="已驳回" value="rejected" />
            </el-select>
          </el-form-item>
        </el-col>
      </el-row>

      <!-- 项目成员 -->
      <el-form-item label="项目成员">
        <div class="members-section">
          <div class="members-header">
            <h4>项目成员</h4>
            <el-button type="primary" size="small" @click="addMember">
              <el-icon><Plus /></el-icon>
              添加成员
            </el-button>
          </div>
          
          <div v-if="form.members.length === 0" class="empty-members">
            <el-empty description="暂无项目成员" :image-size="60" />
          </div>
          
          <div v-else class="members-list">
            <div 
              v-for="(member, index) in form.members" 
              :key="index" 
              class="member-item"
            >
              <el-row :gutter="10">
                <el-col :span="8">
                  <el-input
                    v-model="member.name"
                    placeholder="成员姓名"
                    size="small"
                  />
                </el-col>
                <el-col :span="8">
                  <el-input
                    v-model="member.studentNumber"
                    placeholder="学号"
                    size="small"
                  />
                </el-col>
                <el-col :span="6">
                  <el-input
                    v-model="member.role"
                    placeholder="角色"
                    size="small"
                  />
                </el-col>
                <el-col :span="2">
                  <el-button 
                    type="danger" 
                    size="small" 
                    @click="removeMember(index)"
                    :icon="Delete"
                  />
                </el-col>
              </el-row>
            </div>
          </div>
        </div>
      </el-form-item>

      <!-- 项目附件 -->
      <el-form-item label="项目附件">
        <div class="attachments-section">
          <div class="attachments-header">
            <h4>项目附件</h4>
            <el-upload
              ref="uploadRef"
              :action="uploadAction"
              :headers="uploadHeaders"
              :before-upload="beforeUpload"
              :on-success="handleUploadSuccess"
              :on-error="handleUploadError"
              :show-file-list="false"
              accept=".pdf,.doc,.docx,.xls,.xlsx,.ppt,.pptx,.txt,.jpg,.jpeg,.png,.gif"
            >
              <el-button type="primary" size="small">
                <el-icon><Upload /></el-icon>
                上传文件
              </el-button>
            </el-upload>
          </div>
          
          <div v-if="form.attachments.length === 0" class="empty-attachments">
            <el-empty description="暂无附件" :image-size="60" />
          </div>
          
          <div v-else class="attachments-list">
            <div 
              v-for="(attachment, index) in form.attachments" 
              :key="index" 
              class="attachment-item"
            >
              <el-row :gutter="10" align="middle">
                <el-col :span="16">
                  <div class="attachment-info">
                    <el-icon class="file-icon"><Document /></el-icon>
                    <span class="file-name">{{ attachment.fileName }}</span>
                  </div>
                </el-col>
                <el-col :span="6">
                  <el-button 
                    type="primary" 
                    size="small" 
                    @click="downloadFile(attachment)"
                  >
                    下载
                  </el-button>
                </el-col>
                <el-col :span="2">
                  <el-button 
                    type="danger" 
                    size="small" 
                    @click="removeAttachment(index)"
                    :icon="Delete"
                  />
                </el-col>
              </el-row>
            </div>
          </div>
        </div>
      </el-form-item>
    </el-form>

    <template #footer>
      <div class="dialog-footer">
        <el-button @click="handleClose">取消</el-button>
        <el-button 
          v-if="form.status === 'draft' || form.status === 'rejected'"
          type="primary" 
          @click="handleSubmit" 
          :loading="submitting"
        >
          {{ isEdit ? '保存' : '创建' }}
        </el-button>
        <el-button 
          v-if="form.status === 'draft'"
          type="success" 
          @click="handleSubmitForReview" 
          :loading="submitting"
        >
          提交审核
        </el-button>
        <el-button 
          v-if="form.status === 'rejected'"
          type="warning" 
          @click="handleResubmit" 
          :loading="submitting"
        >
          重新提交
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus, Delete, Upload, Document } from '@element-plus/icons-vue'
import { teacherService } from '../services/teacherService'
import { getToken } from '../services/auth'
import TeacherSelector from './TeacherSelector.vue'

// Props
const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  projectData: {
    type: Object,
    default: null
  }
})

// Emits
const emit = defineEmits(['update:visible', 'success'])

// 响应式数据
const formRef = ref()
const loading = ref(false)
const submitting = ref(false)
const uploadRef = ref()

// 表单数据
const form = ref({
  title: '',
  description: '',
  type: '科研',
  status: 'draft',
  teacherId: '', // 指导老师ID
  members: [],
  attachments: []
})

// 表单验证规则
const rules = {
  title: [
    { required: true, message: '请输入项目标题', trigger: 'blur' },
    { min: 2, max: 100, message: '标题长度在 2 到 100 个字符', trigger: 'blur' }
  ],
  description: [
    { required: true, message: '请输入项目描述', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择项目类型', trigger: 'change' }
  ],
  teacherId: [
    { required: true, message: '请选择指导老师', trigger: 'change' }
  ],
  status: [
    { required: true, message: '请选择项目状态', trigger: 'change' }
  ]
}

// 计算属性
const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

const isEdit = computed(() => !!props.projectData)

// 上传相关
const uploadAction = 'http://localhost:8080/api/files/upload'
const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${getToken()}`
}))

// 重置表单
const resetForm = () => {
  form.value = {
    title: '',
    description: '',
    type: '科研',
    status: 'draft',
    teacherId: '', // 指导老师ID
    members: [],
    attachments: []
  }
  formRef.value?.clearValidate()
}

// 处理教师选择变化
const handleTeacherChange = (teacherInfo) => {
  console.log('选择的指导老师:', teacherInfo)
}

// 监听项目数据变化
watch(() => props.projectData, (newData) => {
  if (newData) {
    form.value = { ...newData }
  } else {
    resetForm()
  }
}, { immediate: true })

// 添加成员
const addMember = () => {
  form.value.members.push({
    name: '',
    studentNumber: '',
    role: ''
  })
}

// 移除成员
const removeMember = (index) => {
  form.value.members.splice(index, 1)
}

// 文件上传前检查
const beforeUpload = (file) => {
  const isValidType = [
    'application/pdf',
    'application/msword',
    'application/vnd.openxmlformats-officedocument.wordprocessingml.document',
    'application/vnd.ms-excel',
    'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
    'application/vnd.ms-powerpoint',
    'application/vnd.openxmlformats-officedocument.presentationml.presentation',
    'text/plain',
    'image/jpeg',
    'image/png',
    'image/gif'
  ].includes(file.type)

  if (!isValidType) {
    ElMessage.error('不支持的文件类型')
    return false
  }

  const isLt10M = file.size / 1024 / 1024 < 10
  if (!isLt10M) {
    ElMessage.error('文件大小不能超过 10MB')
    return false
  }

  return true
}

// 文件上传成功
const handleUploadSuccess = (response) => {
  if (response.code === 200) {
    form.value.attachments.push({
      fileName: response.data.fileName,
      fileUrl: response.data.fileUrl
    })
    ElMessage.success('文件上传成功')
  } else {
    ElMessage.error(response.message || '文件上传失败')
  }
}

// 文件上传失败
const handleUploadError = () => {
  ElMessage.error('文件上传失败')
}

// 下载文件
const downloadFile = (attachment) => {
  const link = document.createElement('a')
  link.href = `http://localhost:8080${attachment.fileUrl}`
  link.download = attachment.fileName
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

// 移除附件
const removeAttachment = (index) => {
  form.value.attachments.splice(index, 1)
}

// 提交表单（保存）
const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    submitting.value = true

    if (isEdit.value) {
      await projectService.updateProject(props.projectData.id, form.value)
      ElMessage.success('项目保存成功')
    } else {
      await projectService.createProject(form.value)
      ElMessage.success('项目创建成功')
    }

    emit('success')
    handleClose()
  } catch (error) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

// 提交审核
const handleSubmitForReview = async () => {
  try {
    await formRef.value.validate()
    submitting.value = true

    // 设置状态为待审核
    const submitData = { ...form.value, status: 'pending' }

    if (isEdit.value) {
      await projectService.updateProject(props.projectData.id, submitData)
      ElMessage.success('项目已提交审核')
    } else {
      await projectService.createProject(submitData)
      ElMessage.success('项目已提交审核')
    }

    emit('success')
    handleClose()
  } catch (error) {
    ElMessage.error(error.message || '提交审核失败')
  } finally {
    submitting.value = false
  }
}

// 重新提交
const handleResubmit = async () => {
  try {
    await formRef.value.validate()
    submitting.value = true

    // 设置状态为待审核
    const submitData = { ...form.value, status: 'pending' }

    await projectService.updateProject(props.projectData.id, submitData)
    ElMessage.success('项目已重新提交审核')

    emit('success')
    handleClose()
  } catch (error) {
    ElMessage.error(error.message || '重新提交失败')
  } finally {
    submitting.value = false
  }
}

// 关闭对话框
const handleClose = () => {
  dialogVisible.value = false
  resetForm()
}
</script>

<style scoped>
.members-section,
.attachments-section {
  border: 1px solid #e4e7ed;
  border-radius: 4px;
  padding: 15px;
  background: #fafafa;
}

.members-header,
.attachments-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.members-header h4,
.attachments-header h4 {
  margin: 0;
  color: #303133;
  font-size: 14px;
}

.empty-members,
.empty-attachments {
  text-align: center;
  padding: 20px 0;
}

.member-item,
.attachment-item {
  margin-bottom: 10px;
  padding: 10px;
  background: white;
  border-radius: 4px;
  border: 1px solid #e4e7ed;
}

.attachment-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.file-icon {
  color: #409eff;
  font-size: 16px;
}

.file-name {
  color: #303133;
  font-size: 14px;
}

.dialog-footer {
  text-align: right;
}
</style> 
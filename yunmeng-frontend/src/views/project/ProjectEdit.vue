<template>
  <div class="project-edit">
    <el-card>
      <template #header>
        <div class="page-header">
          <el-button @click="goBack" icon="ArrowLeft">返回</el-button>
          <span class="page-title">编辑项目</span>
        </div>
      </template>

      <div v-loading="loading" class="edit-content">
        <el-form
          ref="formRef"
          :model="form"
          :rules="rules"
          label-width="100px"
          class="edit-form"
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
              <el-form-item label="项目状态" prop="status">
                <el-select 
                  v-model="form.status" 
                  placeholder="请选择项目状态" 
                  style="width: 100%"
                  :disabled="route.query.fromStatus === 'rejected'"
                >
                  <el-option label="草稿" value="draft" />
                  <el-option 
                    v-if="route.query.fromStatus !== 'rejected'" 
                    label="提交审核" 
                    value="pending" 
                  />
                </el-select>
                <div v-if="route.query.fromStatus === 'rejected'" class="status-tip">
                  <el-text type="info" size="small">已驳回项目修改后自动转为草稿状态</el-text>
                </div>
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

          <!-- 操作按钮 -->
          <el-form-item>
            <div class="form-actions">
              <el-button @click="goBack">取消</el-button>
              <el-button type="primary" @click="handleSubmit" :loading="submitting">
                保存项目
              </el-button>
            </div>
          </el-form-item>
        </el-form>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Plus, Delete, Upload, Document, ArrowLeft } from '@element-plus/icons-vue'
import { projectService } from '../../services/projectService'
import { getToken } from '../../services/auth'

const route = useRoute()
const router = useRouter()

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
  teacherId: 0, // 添加teacherId字段
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
  status: [
    { required: true, message: '请选择项目状态', trigger: 'change' }
  ]
}

// 上传相关
const uploadAction = 'http://localhost:8080/api/files/upload'
const uploadHeaders = {
  Authorization: `Bearer ${getToken()}`
}

// 加载项目数据
const loadProject = async () => {
  const projectId = route.params.id
  if (!projectId) {
    ElMessage.error('项目ID不存在')
    goBack()
    return
  }

  loading.value = true
  try {
    const response = await projectService.getProjectDetail(projectId)
    const projectData = response.data
    
    // 检查项目状态，只有草稿或已驳回状态可以编辑
    if (projectData.status !== 'draft' && projectData.status !== 'rejected') {
      ElMessage.error('只有草稿或已驳回状态的项目可以编辑')
      goBack()
      return
    }

    // 获取原始状态（用于判断是否是从驳回状态编辑）
    const fromStatus = route.query.fromStatus || projectData.status

    // 填充表单数据
    form.value = {
      title: projectData.title || '',
      description: projectData.description || '',
      type: projectData.type || '科研',
      status: fromStatus === 'rejected' ? 'draft' : (projectData.status || 'draft'), // 已驳回项目编辑时转为草稿
      teacherId: projectData.teacher?.id || projectData.teacherId || 0, // 添加teacherId字段
      members: projectData.members || [],
      attachments: projectData.files || []
    }

    // 如果是从驳回状态编辑，显示提示信息
    if (fromStatus === 'rejected') {
      ElMessage.info('项目已从驳回状态转为草稿，修改完成后可重新提交审核')
    }
  } catch (error) {
    ElMessage.error(error.message || '加载项目数据失败')
    goBack()
  } finally {
    loading.value = false
  }
}

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

// 提交表单
const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    submitting.value = true

    const projectId = route.params.id
    const fromStatus = route.query.fromStatus
    
    // 如果是从驳回状态编辑，确保状态为草稿
    const updateData = {
      ...form.value,
      status: fromStatus === 'rejected' ? 'draft' : form.value.status
    }
    
    await projectService.updateProject(projectId, updateData)
    
    const successMessage = fromStatus === 'rejected' 
      ? '项目修改成功，已转为草稿状态' 
      : '项目更新成功'
    ElMessage.success(successMessage)
    goBack()
  } catch (error) {
    ElMessage.error(error.message || '更新项目失败')
  } finally {
    submitting.value = false
  }
}

// 返回上一页
const goBack = () => {
  router.go(-1)
}

// 组件挂载时加载数据
onMounted(() => {
  loadProject()
})
</script>

<style scoped>
.project-edit {
  padding: 20px;
}

.page-header {
  display: flex;
  align-items: center;
  gap: 15px;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.edit-content {
  padding: 20px 0;
}

.edit-form {
  max-width: 800px;
}

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

.form-actions {
  display: flex;
  justify-content: center;
  gap: 15px;
  margin-top: 30px;
}

.status-tip {
  margin-top: 5px;
}
</style> 
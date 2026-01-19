<template>
  <div class="project-extensions">
    <el-card>
      <template #header>
        <div class="header-content">
          <span>项目延期申请管理</span>
          <el-button type="primary" @click="showApplyDialog">申请延期</el-button>
        </div>
      </template>

      <!-- 项目选择 -->
      <el-form :inline="true" class="project-selector">
        <el-form-item label="选择项目">
          <el-select v-model="selectedProjectId" placeholder="请选择项目" @change="loadExtensions">
            <el-option
              v-for="project in projectList"
              :key="project.id"
              :label="project.name"
              :value="project.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="success" @click="loadExtensions" :disabled="!selectedProjectId">
            加载申请记录
          </el-button>
        </el-form-item>
      </el-form>

      <!-- 延期申请管理区域 -->
      <div v-if="selectedProjectId" class="extensions-container">
        <!-- 统计信息 -->
        <el-row :gutter="20" class="extension-stats">
          <el-col :span="6" v-for="stat in extensionStats" :key="stat.label">
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

        <!-- 申请列表 -->
        <div v-if="extensions.length > 0" class="extensions-content">
          <el-table :data="extensions" style="width: 100%" v-loading="loading">
            <el-table-column prop="projectName" label="项目名称" min-width="150" />
            <el-table-column prop="originalDeadline" label="原截止时间" width="120" />
            <el-table-column prop="requestedDeadline" label="申请截止时间" width="120" />
            <el-table-column prop="extensionDays" label="延期天数" width="100">
              <template #default="scope">
                <el-tag :type="getExtensionDaysType(scope.row.extensionDays)" size="small">
                  +{{ scope.row.extensionDays }}天
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="reason" label="延期原因" min-width="200" show-overflow-tooltip />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="scope">
                <el-tag :type="getStatusType(scope.row.status)" size="small">
                  {{ getStatusLabel(scope.row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="applyTime" label="申请时间" width="150" />
            <el-table-column prop="reviewTime" label="审核时间" width="150" />
            <el-table-column prop="reviewer" label="审核人" width="100" />
            <el-table-column label="操作" width="200" fixed="right">
              <template #default="scope">
                <el-button size="small" @click="viewExtension(scope.row)">查看</el-button>
                <el-button 
                  v-if="scope.row.status === 'pending'"
                  size="small" 
                  type="warning" 
                  @click="editExtension(scope.row)"
                >
                  编辑
                </el-button>
                <el-button 
                  v-if="scope.row.status === 'pending'"
                  size="small" 
                  type="danger" 
                  @click="cancelExtension(scope.row)"
                >
                  取消
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <!-- 空状态 -->
        <el-empty
          v-else
          description="暂无延期申请记录"
        >
          <el-button type="primary" @click="showApplyDialog">申请第一个延期</el-button>
        </el-empty>
      </div>

      <!-- 项目选择提示 -->
      <el-empty
        v-else
        description="请先选择一个项目"
      />
    </el-card>

    <!-- 申请延期对话框 -->
    <el-dialog
      v-model="applyDialogVisible"
      :title="isEditing ? '编辑延期申请' : '申请项目延期'"
      width="60%"
    >
      <el-form :model="extensionForm" :rules="extensionRules" ref="extensionFormRef" label-width="120px">
        <el-form-item label="项目名称" prop="projectId">
          <el-select v-model="extensionForm.projectId" placeholder="选择项目" :disabled="isEditing">
            <el-option
              v-for="project in projectList"
              :key="project.id"
              :label="project.name"
              :value="project.id"
            />
          </el-select>
        </el-form-item>
        
        <el-form-item label="原截止时间" prop="originalDeadline">
          <el-date-picker
            v-model="extensionForm.originalDeadline"
            type="date"
            placeholder="选择原截止时间"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            :disabled="isEditing"
          />
        </el-form-item>
        
        <el-form-item label="申请截止时间" prop="requestedDeadline">
          <el-date-picker
            v-model="extensionForm.requestedDeadline"
            type="date"
            placeholder="选择申请截止时间"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        
        <el-form-item label="延期天数">
          <el-input-number
            v-model="extensionForm.extensionDays"
            :min="1"
            :max="365"
            :disabled="true"
            style="width: 200px"
          />
          <span style="margin-left: 10px; color: #909399;">天</span>
        </el-form-item>
        
        <el-form-item label="延期原因" prop="reason">
          <el-input
            v-model="extensionForm.reason"
            type="textarea"
            :rows="4"
            placeholder="请详细说明申请延期的原因和必要性"
          />
        </el-form-item>
        
        <el-form-item label="影响分析" prop="impactAnalysis">
          <el-input
            v-model="extensionForm.impactAnalysis"
            type="textarea"
            :rows="3"
            placeholder="请分析延期对项目进度、质量等方面的影响"
          />
        </el-form-item>
        
        <el-form-item label="补救措施" prop="remedialMeasures">
          <el-input
            v-model="extensionForm.remedialMeasures"
            type="textarea"
            :rows="3"
            placeholder="请说明将采取哪些措施来减少延期的影响"
          />
        </el-form-item>
        
        <el-form-item label="附件材料">
          <el-upload
            ref="uploadRef"
            :auto-upload="false"
            :on-change="handleFileChange"
            :on-remove="handleFileRemove"
            :file-list="extensionForm.attachments"
            multiple
          >
            <el-button type="primary">选择文件</el-button>
            <template #tip>
              <div class="el-upload__tip">
                支持任意类型文件，单个文件不超过50MB
              </div>
            </template>
          </el-upload>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="applyDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitExtension" :loading="submitting">提交申请</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 延期申请详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      title="延期申请详情"
      width="70%"
    >
      <div v-if="currentExtension" class="extension-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="项目名称">{{ currentExtension.projectName }}</el-descriptions-item>
          <el-descriptions-item label="申请状态">
            <el-tag :type="getStatusType(currentExtension.status)" size="small">
              {{ getStatusLabel(currentExtension.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="原截止时间">{{ currentExtension.originalDeadline }}</el-descriptions-item>
          <el-descriptions-item label="申请截止时间">{{ currentExtension.requestedDeadline }}</el-descriptions-item>
          <el-descriptions-item label="延期天数">
            <el-tag :type="getExtensionDaysType(currentExtension.extensionDays)" size="small">
              +{{ currentExtension.extensionDays }}天
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="申请时间">{{ currentExtension.applyTime }}</el-descriptions-item>
          <el-descriptions-item label="审核人">{{ currentExtension.reviewer || '待审核' }}</el-descriptions-item>
          <el-descriptions-item label="审核时间">{{ currentExtension.reviewTime || '待审核' }}</el-descriptions-item>
        </el-descriptions>
        
        <el-divider />
        
        <h4>延期原因</h4>
        <p>{{ currentExtension.reason }}</p>
        
        <el-divider />
        
        <h4>影响分析</h4>
        <p>{{ currentExtension.impactAnalysis || '未提供' }}</p>
        
        <el-divider />
        
        <h4>补救措施</h4>
        <p>{{ currentExtension.remedialMeasures || '未提供' }}</p>
        
        <el-divider v-if="currentExtension.attachments && currentExtension.attachments.length > 0" />
        
        <div v-if="currentExtension.attachments && currentExtension.attachments.length > 0">
          <h4>附件材料</h4>
          <ul>
            <li v-for="attachment in currentExtension.attachments" :key="attachment.name">
              <el-link type="primary" @click="downloadAttachment(attachment)">
                {{ attachment.name }}
              </el-link>
            </li>
          </ul>
        </div>
        
        <el-divider v-if="currentExtension.comments" />
        
        <div v-if="currentExtension.comments">
          <h4>审核意见</h4>
          <p>{{ currentExtension.comments }}</p>
        </div>
      </div>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="detailDialogVisible = false">关闭</el-button>
          <el-button 
            v-if="currentExtension && currentExtension.status === 'pending'"
            type="warning" 
            @click="editExtension(currentExtension)"
          >
            编辑申请
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Clock, Check, Warning, Close, Document } from '@element-plus/icons-vue'
import { studentService } from '../../services/studentService'

// 响应式数据
const selectedProjectId = ref('')
const projectList = ref([])
const extensions = ref([])
const loading = ref(false)
const applyDialogVisible = ref(false)
const detailDialogVisible = ref(false)
const currentExtension = ref(null)
const isEditing = ref(false)
const submitting = ref(false)

const extensionForm = ref({
  projectId: '',
  originalDeadline: '',
  requestedDeadline: '',
  extensionDays: 0,
  reason: '',
  impactAnalysis: '',
  remedialMeasures: '',
  attachments: []
})

const extensionRules = {
  projectId: [{ required: true, message: '请选择项目', trigger: 'change' }],
  originalDeadline: [{ required: true, message: '请选择原截止时间', trigger: 'change' }],
  requestedDeadline: [{ required: true, message: '请选择申请截止时间', trigger: 'change' }],
  reason: [{ required: true, message: '请输入延期原因', trigger: 'blur' }],
  impactAnalysis: [{ required: true, message: '请输入影响分析', trigger: 'blur' }],
  remedialMeasures: [{ required: true, message: '请输入补救措施', trigger: 'blur' }]
}

const extensionFormRef = ref()
const uploadRef = ref()

// 延期申请统计
const extensionStats = computed(() => [
  {
    label: '总申请数',
    value: extensions.value.length,
    description: '所有申请',
    icon: Document,
    type: 'total'
  },
  {
    label: '待审核',
    value: extensions.value.filter(e => e.status === 'pending').length,
    description: '等待审核',
    icon: Clock,
    type: 'pending'
  },
  {
    label: '已通过',
    value: extensions.value.filter(e => e.status === 'approved').length,
    description: '审核通过',
    icon: Check,
    type: 'approved'
  },
  {
    label: '已拒绝',
    value: extensions.value.filter(e => e.status === 'rejected').length,
    description: '审核拒绝',
    icon: Close,
    type: 'rejected'
  }
])

// 监听截止时间变化，自动计算延期天数
watch(() => [extensionForm.value.originalDeadline, extensionForm.value.requestedDeadline], ([original, requested]) => {
  if (original && requested) {
    const originalDate = new Date(original)
    const requestedDate = new Date(requested)
    const diffTime = requestedDate.getTime() - originalDate.getTime()
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
    extensionForm.value.extensionDays = diffDays > 0 ? diffDays : 0
  }
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
        { id: 1, name: '智能校园系统', deadline: '2024-07-15' },
        { id: 2, name: '数据分析平台', deadline: '2024-06-14' },
        { id: 3, name: '在线教育平台', deadline: '2024-05-10' }
      ]
    }
  } catch (error) {
    console.error('加载项目列表失败:', error)
    // 使用模拟数据
    projectList.value = [
      { id: 1, name: '智能校园系统', deadline: '2024-07-15' },
      { id: 2, name: '数据分析平台', deadline: '2024-06-14' },
      { id: 3, name: '在线教育平台', deadline: '2024-05-10' }
    ]
  }
}

// 加载延期申请
const loadExtensions = async () => {
  if (!selectedProjectId.value) return
  
  loading.value = true
  try {
    // 这里应该调用实际的API
    // const response = await studentService.getProjectExtensions(selectedProjectId.value)
    
    // 使用模拟数据
    extensions.value = [
      {
        id: 1,
        projectId: 1,
        projectName: '智能校园系统',
        originalDeadline: '2024-07-15',
        requestedDeadline: '2024-08-15',
        extensionDays: 31,
        reason: '由于系统架构调整和新增功能需求，原计划时间不足以完成所有开发工作。需要额外时间进行充分测试和优化。',
        impactAnalysis: '延期将影响项目交付时间，但能确保系统质量和稳定性。',
        remedialMeasures: '增加开发人员，优化开发流程，采用敏捷开发方法提高效率。',
        status: 'pending',
        applyTime: '2024-06-01 10:30:00',
        attachments: [
          { name: '延期申请说明.docx', url: '/api/files/1' },
          { name: '项目进度报告.pdf', url: '/api/files/2' }
        ]
      },
      {
        id: 2,
        projectId: 2,
        projectName: '数据分析平台',
        originalDeadline: '2024-06-14',
        requestedDeadline: '2024-07-14',
        extensionDays: 30,
        reason: '数据采集模块遇到技术难题，需要重新设计数据管道架构。',
        impactAnalysis: '延期将影响后续模块开发，但能解决技术瓶颈问题。',
        remedialMeasures: '引入专家顾问，重新评估技术方案，优化开发计划。',
        status: 'approved',
        applyTime: '2024-05-15 14:20:00',
        reviewTime: '2024-05-20 09:15:00',
        reviewer: '李教授',
        comments: '延期理由充分，技术方案合理，同意延期申请。',
        attachments: [
          { name: '技术方案说明.pdf', url: '/api/files/3' }
        ]
      },
      {
        id: 3,
        projectId: 3,
        projectName: '在线教育平台',
        originalDeadline: '2024-05-10',
        requestedDeadline: '2024-06-10',
        extensionDays: 31,
        reason: '用户界面设计需要重新优化，用户体验测试反馈较多。',
        impactAnalysis: '延期将影响项目上线时间，但能提升用户体验质量。',
        remedialMeasures: '重新设计UI/UX，增加用户测试环节，优化交互流程。',
        status: 'rejected',
        applyTime: '2024-04-20 16:45:00',
        reviewTime: '2024-04-25 11:30:00',
        reviewer: '王教授',
        comments: '延期理由不够充分，建议在现有时间内优化完成。',
        attachments: [
          { name: 'UI设计稿.pdf', url: '/api/files/4' }
        ]
      }
    ]
  } catch (error) {
    console.error('加载延期申请失败:', error)
    ElMessage.error('加载延期申请失败')
  } finally {
    loading.value = false
  }
}

// 显示申请对话框
const showApplyDialog = () => {
  if (!selectedProjectId.value) {
    ElMessage.warning('请先选择一个项目')
    return
  }
  
  isEditing.value = false
  currentExtension.value = null
  resetExtensionForm()
  applyDialogVisible.value = true
}

// 编辑延期申请
const editExtension = (extension) => {
  isEditing.value = true
  currentExtension.value = extension
  extensionForm.value = {
    projectId: extension.projectId,
    originalDeadline: extension.originalDeadline,
    requestedDeadline: extension.requestedDeadline,
    extensionDays: extension.extensionDays,
    reason: extension.reason,
    impactAnalysis: extension.impactAnalysis,
    remedialMeasures: extension.remedialMeasures,
    attachments: extension.attachments || []
  }
  applyDialogVisible.value = true
}

// 查看延期申请详情
const viewExtension = (extension) => {
  currentExtension.value = extension
  detailDialogVisible.value = true
}

// 取消延期申请
const cancelExtension = async (extension) => {
  try {
    await ElMessageBox.confirm('确定要取消这个延期申请吗？', '确认取消', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    // 这里应该调用实际的API
    // await studentService.cancelExtension(extension.id)
    
    const index = extensions.value.findIndex(e => e.id === extension.id)
    if (index > -1) {
      extensions.value.splice(index, 1)
    }
    
    ElMessage.success('延期申请已取消')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('取消延期申请失败:', error)
      ElMessage.error('操作失败')
    }
  }
}

// 文件选择变化
const handleFileChange = (file, fileList) => {
  extensionForm.value.attachments = fileList
}

// 文件移除
const handleFileRemove = (file, fileList) => {
  extensionForm.value.attachments = fileList
}

// 下载附件
const downloadAttachment = (attachment) => {
  // 这里应该调用实际的下载API
  ElMessage.success(`开始下载附件: ${attachment.name}`)
  
  // 模拟下载
  const link = document.createElement('a')
  link.href = attachment.url || '#'
  link.download = attachment.name
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
}

// 重置延期申请表单
const resetExtensionForm = () => {
  extensionForm.value = {
    projectId: selectedProjectId.value,
    originalDeadline: '',
    requestedDeadline: '',
    extensionDays: 0,
    reason: '',
    impactAnalysis: '',
    remedialMeasures: '',
    attachments: []
  }
}

// 提交延期申请
const submitExtension = async () => {
  try {
    await extensionFormRef.value.validate()
    
    if (extensionForm.value.extensionDays <= 0) {
      ElMessage.warning('申请截止时间必须晚于原截止时间')
      return
    }
    
    submitting.value = true
    
    // 这里应该调用实际的API
    // await studentService.submitExtension(extensionForm.value)
    
    // 模拟提交过程
    await new Promise(resolve => setTimeout(resolve, 1500))
    
    if (isEditing.value) {
      // 编辑现有申请
      Object.assign(currentExtension.value, {
        originalDeadline: extensionForm.value.originalDeadline,
        requestedDeadline: extensionForm.value.requestedDeadline,
        extensionDays: extensionForm.value.extensionDays,
        reason: extensionForm.value.reason,
        impactAnalysis: extensionForm.value.impactAnalysis,
        remedialMeasures: extensionForm.value.remedialMeasures,
        attachments: extensionForm.value.attachments
      })
      
      ElMessage.success('延期申请更新成功')
    } else {
      // 创建新申请
      const newExtension = {
        id: Date.now(),
        projectId: extensionForm.value.projectId,
        projectName: projectList.value.find(p => p.id === extensionForm.value.projectId)?.name,
        originalDeadline: extensionForm.value.originalDeadline,
        requestedDeadline: extensionForm.value.requestedDeadline,
        extensionDays: extensionForm.value.extensionDays,
        reason: extensionForm.value.reason,
        impactAnalysis: extensionForm.value.impactAnalysis,
        remedialMeasures: extensionForm.value.remedialMeasures,
        status: 'pending',
        applyTime: new Date().toLocaleString('zh-CN'),
        attachments: extensionForm.value.attachments
      }
      
      extensions.value.unshift(newExtension)
      ElMessage.success('延期申请提交成功')
    }
    
    applyDialogVisible.value = false
  } catch (error) {
    console.error('提交延期申请失败:', error)
    ElMessage.error('提交失败')
  } finally {
    submitting.value = false
  }
}

// 获取状态类型
const getStatusType = (status) => {
  const typeMap = {
    'pending': 'warning',
    'approved': 'success',
    'rejected': 'danger'
  }
  return typeMap[status] || 'info'
}

// 获取状态标签
const getStatusLabel = (status) => {
  const labelMap = {
    'pending': '待审核',
    'approved': '已通过',
    'rejected': '已拒绝'
  }
  return labelMap[status] || '未知'
}

// 获取延期天数类型
const getExtensionDaysType = (days) => {
  if (days <= 7) return 'success'
  if (days <= 30) return 'warning'
  return 'danger'
}

// 组件挂载时加载数据
onMounted(() => {
  loadProjects()
})
</script>

<style scoped>
.project-extensions {
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

.extensions-container {
  margin-top: 20px;
}

.extension-stats {
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

.extensions-content {
  margin-top: 20px;
}

.extension-detail {
  padding: 20px;
}

.extension-detail h4 {
  margin: 20px 0 10px 0;
  color: #2c3e50;
  font-size: 16px;
}

.extension-detail p {
  margin: 10px 0;
  color: #606266;
  line-height: 1.6;
}

.extension-detail ul {
  margin: 10px 0;
  padding-left: 20px;
}

.extension-detail li {
  margin: 5px 0;
  color: #606266;
}

.dialog-footer {
  text-align: right;
}

:deep(.el-upload__tip) {
  color: #909399;
  font-size: 12px;
  margin-top: 5px;
}
</style> 
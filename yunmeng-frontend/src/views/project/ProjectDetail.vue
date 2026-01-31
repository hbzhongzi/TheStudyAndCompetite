<template>
  <div class="project-detail">
    <el-card>
      <template #header>
        <div class="page-header">
          <el-button @click="goBack" icon="ArrowLeft">返回</el-button>
          <span class="page-title">项目详情</span>
          <div class="header-actions">
            <el-button 
              v-if="canEdit" 
              type="primary" 
              @click="editProject"
              icon="Edit"
            >
              编辑项目
            </el-button>
            <el-button 
              v-if="canSubmit" 
              type="success" 
              @click="submitProject"
              icon="Upload"
            >
              提交审核
            </el-button>
          </div>
        </div>
      </template>

      <div v-loading="loading" class="detail-content">
        <!-- 项目基本信息 -->
        <el-descriptions title="项目基本信息" :column="2" border>
          <el-descriptions-item label="项目标题">{{ project.title }}</el-descriptions-item>
          <el-descriptions-item label="项目类型">{{ project.type }}</el-descriptions-item>
          <el-descriptions-item label="项目状态">
            <el-tag :type="getStatusType(project.status)">{{ getStatusText(project.status) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ formatDate(project.createdAt) }}</el-descriptions-item>
          <el-descriptions-item label="项目描述" :span="2">{{ project.description }}</el-descriptions-item>
        </el-descriptions>

        <!-- 项目进度和生命周期信息 -->
        <el-card class="progress-card" style="margin-top: 20px;">
          <template #header>
            <div class="card-header">
              <span>项目进度</span>
              <el-button 
                v-if="canUpdateProgress" 
                type="primary" 
                size="small" 
                @click="showProgressDialog = true"
                icon="Edit"
              >
                更新进度
              </el-button>
            </div>
          </template>
          
          <el-progress 
            :percentage="project.progress || 0" 
            :status="getProgressStatus(project.progress)"
            :stroke-width="20"
          />
          
          <div class="progress-info">
            <el-row :gutter="20">
              <el-col :span="8">
                <div class="info-item">
                  <span class="label">开始时间:</span>
                  <span class="value">{{ formatDate(project.startDate) || '未设置' }}</span>
                </div>
              </el-col>
              <el-col :span="8">
                <div class="info-item">
                  <span class="label">预计完成:</span>
                  <span class="value">{{ formatDate(project.expectedEndDate) || '未设置' }}</span>
                </div>
              </el-col>
              <el-col :span="8">
                <div class="info-item">
                  <span class="label">实际完成:</span>
                  <span class="value">{{ formatDate(project.actualEndDate) || '未完成' }}</span>
                </div>
              </el-col>
            </el-row>
          </div>
        </el-card>

        <!-- 项目里程碑 -->
        <el-card class="milestone-card" style="margin-top: 20px;">
          <template #header>
            <div class="card-header">
              <span>项目里程碑</span>
              <el-button 
                v-if="canManageMilestone" 
                type="primary" 
                size="small" 
                @click="showMilestoneDialog = true"
                icon="Plus"
              >
                添加里程碑
              </el-button>
            </div>
          </template>
          
          <div v-if="milestones.length === 0" class="empty-milestones">
            <el-empty description="暂无里程碑" :image-size="60" />
          </div>
          
          <el-timeline v-else>
            <el-timeline-item
              v-for="milestone in milestones"
              :key="milestone.id"
              :timestamp="formatDate(milestone.dueDate)"
              :type="getMilestoneType(milestone.status)"
            >
              <el-card class="milestone-item">
                <div class="milestone-header">
                  <h4>{{ milestone.title }}</h4>
                  <el-tag :type="getMilestoneStatusType(milestone.status)">
                    {{ getMilestoneStatusText(milestone.status) }}
                  </el-tag>
                </div>
                <p class="milestone-description">{{ milestone.description }}</p>
                <div class="milestone-progress">
                  <el-progress 
                    :percentage="milestone.progress" 
                    :status="getProgressStatus(milestone.progress)"
                    :stroke-width="8"
                  />
                </div>
                <div class="milestone-actions" v-if="canManageMilestone">
                  <el-button 
                    size="small" 
                    @click="editMilestone(milestone)"
                    icon="Edit"
                  >
                    编辑
                  </el-button>
                  <el-button 
                    size="small" 
                    type="success" 
                    @click="completeMilestone(milestone)"
                    v-if="milestone.status === 'in_progress'"
                    icon="Check"
                  >
                    完成
                  </el-button>
                </div>
              </el-card>
            </el-timeline-item>
          </el-timeline>
        </el-card>

        <!-- 项目文件 -->
        <el-card class="file-card" style="margin-top: 20px;">
          <template #header>
            <div class="card-header">
              <span>项目文件</span>
              <el-button 
                v-if="canUploadFile" 
                type="primary" 
                size="small" 
                @click="showFileDialog = true"
                icon="Upload"
              >
                上传文件
              </el-button>
            </div>
          </template>
          
          <el-tabs v-model="activeFileTab" @tab-click="handleFileTabClick">
            <el-tab-pane label="全部文件" name="all">
              <ProjectFileList 
                :files="allFiles" 
                :canReview="canReviewFile"
                @review="reviewFile"
                @download="downloadFile"
              />
            </el-tab-pane>
            <el-tab-pane label="开题报告" name="proposal">
              <ProjectFileList 
                :files="getFilesByType('proposal')" 
                :canReview="canReviewFile"
                @review="reviewFile"
                @download="downloadFile"
              />
            </el-tab-pane>
            <el-tab-pane label="中期报告" name="midterm">
              <ProjectFileList 
                :files="getFilesByType('midterm')" 
                :canReview="canReviewFile"
                @review="reviewFile"
                @download="downloadFile"
              />
            </el-tab-pane>
            <el-tab-pane label="结题报告" name="final">
              <ProjectFileList 
                :files="getFilesByType('final')" 
                :canReview="canReviewFile"
                @review="reviewFile"
                @download="downloadFile"
              />
            </el-tab-pane>
            <el-tab-pane label="成果展示" name="achievement">
              <ProjectFileList 
                :files="getFilesByType('achievement')" 
                :canReview="canReviewFile"
                @review="reviewFile"
                @download="downloadFile"
              />
            </el-tab-pane>
          </el-tabs>
        </el-card>

        <!-- 项目延期申请 -->
        <el-card class="extension-card" style="margin-top: 20px;" v-if="canApplyExtension">
          <template #header>
            <div class="card-header">
              <span>延期申请</span>
              <el-button 
                type="warning" 
                size="small" 
                @click="showExtensionDialog = true"
                icon="Clock"
              >
                申请延期
              </el-button>
            </div>
          </template>
          
          <div v-if="extensions.length === 0" class="empty-extensions">
            <el-empty description="暂无延期申请" :image-size="60" />
          </div>
          
          <el-timeline v-else>
            <el-timeline-item
              v-for="extension in extensions"
              :key="extension.id"
              :timestamp="formatDate(extension.createdAt)"
              :type="getExtensionType(extension.status)"
            >
              <el-card class="extension-item">
                <div class="extension-header">
                  <h4>延期申请</h4>
                  <el-tag :type="getExtensionStatusType(extension.status)">
                    {{ getExtensionStatusText(extension.status) }}
                  </el-tag>
                </div>
                <div class="extension-content">
                  <p><strong>延期原因:</strong> {{ extension.reason }}</p>
                  <p><strong>原定完成时间:</strong> {{ formatDate(extension.originalEndDate) }}</p>
                  <p><strong>申请延期到:</strong> {{ formatDate(extension.requestedEndDate) }}</p>
                  <p v-if="extension.reviewComments">
                    <strong>审核意见:</strong> {{ extension.reviewComments }}
                  </p>
                </div>
              </el-card>
            </el-timeline-item>
          </el-timeline>
        </el-card>

        <!-- 项目状态历史 -->
        <el-card class="status-history-card" style="margin-top: 20px;">
          <template #header>
            <span>状态变更历史</span>
          </template>
          
          <el-timeline>
            <el-timeline-item
              v-for="history in statusHistory"
              :key="history.id"
              :timestamp="formatDate(history.changedAt)"
              :type="getHistoryType(history.newStatus)"
            >
              <el-card class="history-item">
                <div class="history-header">
                  <span class="status-change">
                    {{ getStatusText(history.oldStatus) }} → {{ getStatusText(history.newStatus) }}
                  </span>
                  <span class="operator">操作人: {{ history.operatorName }}</span>
                </div>
                <p class="change-reason">{{ history.changeReason }}</p>
              </el-card>
            </el-timeline-item>
          </el-timeline>
        </el-card>

        <!-- 审核记录 -->
        <el-card class="review-card" style="margin-top: 20px;">
          <template #header>
            <span>审核记录</span>
          </template>
          
          <div v-if="reviews.length === 0" class="empty-reviews">
            <el-empty description="暂无审核记录" :image-size="60" />
          </div>
          
          <el-timeline v-else>
            <el-timeline-item
              v-for="review in reviews"
              :key="review.id"
              :timestamp="formatDate(review.reviewTime)"
              :type="getReviewType(review.status)"
            >
              <el-card class="review-item">
                <div class="review-header">
                  <span class="reviewer">审核人: {{ review.reviewer }}</span>
                  <el-tag :type="getReviewStatusType(review.status)">
                    {{ getReviewStatusText(review.status) }}
                  </el-tag>
                </div>
                <p class="review-comments">{{ review.comments }}</p>
              </el-card>
            </el-timeline-item>
          </el-timeline>
        </el-card>
      </div>
    </el-card>

    <!-- 进度更新对话框 -->
    <el-dialog v-model="showProgressDialog" title="更新项目进度" width="500px">
      <el-form :model="progressForm" label-width="100px">
        <el-form-item label="项目进度">
          <el-slider
            v-model="progressForm.progress"
            :min="0"
            :max="100"
            :step="5"
            show-input
            show-stops
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showProgressDialog = false">取消</el-button>
        <el-button type="primary" @click="updateProgress">确定</el-button>
      </template>
    </el-dialog>

    <!-- 里程碑对话框 -->
    <el-dialog v-model="showMilestoneDialog" title="添加里程碑" width="600px">
      <el-form :model="milestoneForm" :rules="milestoneRules" label-width="100px" ref="milestoneFormRef">
        <el-form-item label="里程碑标题" prop="title">
          <el-input v-model="milestoneForm.title" placeholder="请输入里程碑标题" />
        </el-form-item>
        <el-form-item label="里程碑描述" prop="description">
          <el-input
            v-model="milestoneForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入里程碑描述"
          />
        </el-form-item>
        <el-form-item label="预计完成时间" prop="dueDate">
          <el-date-picker
            v-model="milestoneForm.dueDate"
            type="datetime"
            placeholder="选择预计完成时间"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showMilestoneDialog = false">取消</el-button>
        <el-button type="primary" @click="createMilestone">确定</el-button>
      </template>
    </el-dialog>

    <!-- 文件上传对话框 -->
    <el-dialog v-model="showFileDialog" title="上传项目文件" width="600px">
      <el-form :model="fileForm" :rules="fileRules" label-width="100px" ref="fileFormRef">
        <el-form-item label="文件名称" prop="fileName">
          <el-input v-model="fileForm.fileName" placeholder="请输入文件名称" />
        </el-form-item>
        <el-form-item label="文件类型" prop="fileType">
          <el-select v-model="fileForm.fileType" placeholder="请选择文件类型" style="width: 100%">
            <el-option label="开题报告" value="proposal" />
            <el-option label="中期报告" value="midterm" />
            <el-option label="结题报告" value="final" />
            <el-option label="成果展示" value="achievement" />
            <el-option label="其他材料" value="other" />
          </el-select>
        </el-form-item>
        <el-form-item label="文件版本" prop="fileVersion">
          <el-input v-model="fileForm.fileVersion" placeholder="请输入文件版本，如：1.0" />
        </el-form-item>
        <el-form-item label="文件上传" prop="file">
          <el-upload
            ref="uploadRef"
            :auto-upload="false"
            :on-change="handleFileChange"
            :limit="1"
            accept=".pdf,.doc,.docx,.ppt,.pptx,.zip,.rar,.jpg,.jpeg,.png"
          >
            <el-button type="primary">选择文件</el-button>
            <template #tip>
              <div class="el-upload__tip">
                支持 PDF、Word、PPT、压缩包、图片等格式，文件大小不超过50MB
              </div>
            </template>
          </el-upload>
        </el-form-item>
        <el-form-item label="是否公开">
          <el-switch v-model="fileForm.isPublic" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showFileDialog = false">取消</el-button>
        <el-button type="primary" @click="uploadFile">确定</el-button>
      </template>
    </el-dialog>

    <!-- 延期申请对话框 -->
    <el-dialog v-model="showExtensionDialog" title="申请项目延期" width="600px">
      <el-form :model="extensionForm" :rules="extensionRules" label-width="120px" ref="extensionFormRef">
        <el-form-item label="延期原因" prop="reason">
          <el-input
            v-model="extensionForm.reason"
            type="textarea"
            :rows="4"
            placeholder="请详细说明申请延期的原因"
          />
        </el-form-item>
        <el-form-item label="申请延期到" prop="requestedEndDate">
          <el-date-picker
            v-model="extensionForm.requestedEndDate"
            type="datetime"
            placeholder="选择新的完成时间"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showExtensionDialog = false">取消</el-button>
        <el-button type="primary" @click="applyExtension">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { projectService } from '@/services/projectService'
import ProjectFileList from '@/components/ProjectFileList.vue'

export default {
  name: 'ProjectDetail',
  components: {
    ProjectFileList
  },
  setup() {
    const route = useRoute()
    const router = useRouter()
    
    // 响应式数据
    const loading = ref(false)
    const project = ref({})
    const milestones = ref([])
    const allFiles = ref([])
    const extensions = ref([])
    const statusHistory = ref([])
    const reviews = ref([])
    
    // 对话框控制
    const showProgressDialog = ref(false)
    const showMilestoneDialog = ref(false)
    const showFileDialog = ref(false)
    const showExtensionDialog = ref(false)
    
    // 表单数据
    const progressForm = reactive({
      progress: 0
    })
    
    const milestoneForm = reactive({
      title: '',
      description: '',
      dueDate: null
    })
    
    const fileForm = reactive({
      fileName: '',
      fileType: '',
      fileVersion: '1.0',
      file: null,
      isPublic: false
    })
    
    const extensionForm = reactive({
      reason: '',
      requestedEndDate: null
    })
    
    // 表单验证规则
    const milestoneRules = {
      title: [{ required: true, message: '请输入里程碑标题', trigger: 'blur' }],
      dueDate: [{ required: true, message: '请选择预计完成时间', trigger: 'change' }]
    }
    
    const fileRules = {
      fileName: [{ required: true, message: '请输入文件名称', trigger: 'blur' }],
      fileType: [{ required: true, message: '请选择文件类型', trigger: 'change' }],
      file: [{ required: true, message: '请选择要上传的文件', trigger: 'change' }]
    }
    
    const extensionRules = {
      reason: [{ required: true, message: '请输入延期原因', trigger: 'blur' }],
      requestedEndDate: [{ required: true, message: '请选择新的完成时间', trigger: 'change' }]
    }
    
    // 计算属性
    const canEdit = computed(() => {
      return project.value.status === 'draft' || project.value.status === 'rejected'
    })
    
    const canSubmit = computed(() => {
      return project.value.status === 'draft'
    })
    
    const canUpdateProgress = computed(() => {
      return project.value.status === 'in_progress' || project.value.status === 'approved'
    })
    
    const canManageMilestone = computed(() => {
      return project.value.status !== 'draft'
    })
    
    const canUploadFile = computed(() => {
      return project.value.status !== 'draft'
    })
    
    const canReviewFile = computed(() => {
      // 教师或管理员可以审核文件
      return true // 这里需要根据用户角色判断
    })
    
    const canApplyExtension = computed(() => {
      return project.value.status === 'in_progress'
    })
    
    // 方法
    const goBack = () => {
      router.go(-1)
    }
    
    const loadProjectDetail = async () => {
      try {
        loading.value = true
        const projectId = route.params.id
        
        // 验证项目ID是否存在
        if (!projectId) {
          ElMessage.error('项目ID无效')
          return
        }
        
        const response = await projectService.getProjectDetail(projectId)
        project.value = response.data
        
        // 加载相关数据
        await Promise.all([
          loadMilestones(projectId),
          loadFiles(projectId),
          loadExtensions(projectId),
          loadStatusHistory(projectId),
          loadReviews(projectId)
        ])
      } catch (error) {
        ElMessage.error(error.message)
      } finally {
        loading.value = false
      }
    }
    
    const loadMilestones = async (projectId) => {
      try {
        if (!projectId) {
          console.warn('项目ID无效，跳过里程碑加载')
          return
        }
        const response = await projectService.getProjectMilestones(projectId)
        milestones.value = response.data || []
      } catch (error) {
        console.error('加载里程碑失败:', error)
      }
    }
    
    const loadFiles = async (projectId) => {
      try {
        if (!projectId) {
          console.warn('项目ID无效，跳过文件加载')
          return
        }
        const response = await projectService.getProjectFilesByType(projectId)
        allFiles.value = response.data || []
      } catch (error) {
        console.error('加载文件失败:', error)
      }
    }
    
    const loadExtensions = async (projectId) => {
      try {
        if (!projectId) {
          console.warn('项目ID无效，跳过延期申请加载')
          return
        }
        // 这里需要后端提供获取延期申请的API
        extensions.value = []
      } catch (error) {
        console.error('加载延期申请失败:', error)
      }
    }
    
    const loadStatusHistory = async (projectId) => {
      try {
        if (!projectId) {
          console.warn('项目ID无效，跳过状态历史加载')
          return
        }
        const response = await projectService.getProjectStatusHistory(projectId)
        statusHistory.value = response.data || []
      } catch (error) {
        console.error('加载状态历史失败:', error)
      }
    }
    
    const loadReviews = async (projectId) => {
      try {
        if (!projectId) {
          console.warn('项目ID无效，跳过审核记录加载')
          return
        }
        const response = await projectService.getProjectReviews(projectId)
        reviews.value = response.data || []
      } catch (error) {
        console.error('加载审核记录失败:', error)
      }
    }
    
    const editProject = () => {
      router.push(`/project/edit/${project.value.id}`)
    }
    
    const submitProject = async () => {
      try {
        await ElMessageBox.confirm('确定要提交项目进行审核吗？提交后将无法修改核心信息。', '确认提交', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        
        await projectService.submitProject(project.value.id)
        ElMessage.success('项目提交成功')
        await loadProjectDetail()
      } catch (error) {
        if (error !== 'cancel') {
          ElMessage.error(error.message)
        }
      }
    }
    
    const updateProgress = async () => {
      try {
        await projectService.updateProjectProgress(project.value.id, progressForm)
        ElMessage.success('项目进度更新成功')
        showProgressDialog.value = false
        await loadProjectDetail()
      } catch (error) {
        ElMessage.error(error.message)
      }
    }
    
    const createMilestone = async () => {
      try {
        await projectService.createProjectMilestone(project.value.id, milestoneForm)
        ElMessage.success('里程碑创建成功')
        showMilestoneDialog.value = false
        await loadMilestones(project.value.id)
        // 重置表单
        Object.assign(milestoneForm, {
          title: '',
          description: '',
          dueDate: null
        })
      } catch (error) {
        ElMessage.error(error.message)
      }
    }
    
    const editMilestone = (milestone) => {
      // 实现编辑里程碑逻辑
      console.log('编辑里程碑:', milestone)
    }
    
    const completeMilestone = async (milestone) => {
      try {
        await projectService.updateProjectMilestone(milestone.id, {
          status: 'completed',
          progress: 100,
          completedDate: new Date()
        })
        ElMessage.success('里程碑完成')
        await loadMilestones(project.value.id)
      } catch (error) {
        ElMessage.error(error.message)
      }
    }
    
    const handleFileChange = (file) => {
      fileForm.file = file.raw
    }
    
    const uploadFile = async () => {
      try {
        // 这里需要实现文件上传逻辑
        const formData = new FormData()
        formData.append('file', fileForm.file)
        formData.append('fileName', fileForm.fileName)
        formData.append('fileType', fileForm.fileType)
        formData.append('fileVersion', fileForm.fileVersion)
        formData.append('isPublic', fileForm.isPublic)
        
        await projectService.uploadProjectFile(project.value.id, formData)
        ElMessage.success('文件上传成功')
        showFileDialog.value = false
        await loadFiles(project.value.id)
        // 重置表单
        Object.assign(fileForm, {
          fileName: '',
          fileType: '',
          fileVersion: '1.0',
          file: null,
          isPublic: false
        })
      } catch (error) {
        ElMessage.error(error.message)
      }
    }
    
    const applyExtension = async () => {
      try {
        await projectService.applyProjectExtension(project.value.id, extensionForm)
        ElMessage.success('延期申请提交成功')
        showExtensionDialog.value = false
        await loadExtensions(project.value.id)
        // 重置表单
        Object.assign(extensionForm, {
          reason: '',
          requestedEndDate: null
        })
      } catch (error) {
        ElMessage.error(error.message)
      }
    }
    
    const reviewFile = async (fileId, reviewData) => {
      try {
        await projectService.reviewProjectFile(fileId, reviewData)
        ElMessage.success('文件审核完成')
        await loadFiles(project.value.id)
      } catch (error) {
        ElMessage.error(error.message)
      }
    }
    
    const downloadFile = (file) => {
      // 实现文件下载逻辑
      console.log('下载文件:', file)
    }
    
    const getFilesByType = (type) => {
      return allFiles.value.filter(file => file.fileType === type)
    }
    
    const handleFileTabClick = (tab) => {
      console.log('切换到文件标签:', tab.props.name)
    }
    
    // 工具方法
    const formatDate = (date) => {
      if (!date) return ''
      return new Date(date).toLocaleString('zh-CN')
    }
    
    const getStatusType = (status) => {
      const statusMap = {
        draft: 'info',
        pending: 'warning',
        approved: 'success',
        rejected: 'danger',
        in_progress: 'primary',
        completed: 'success',
        archived: 'info',
        suspended: 'warning',
        need_revision: 'warning'
      }
      return statusMap[status] || 'info'
    }
    
    const getStatusText = (status) => {
      const statusMap = {
        draft: '草稿',
        pending: '待审核',
        approved: '已通过',
        rejected: '已驳回',
        in_progress: '进行中',
        completed: '已完成',
        archived: '已归档',
        suspended: '已暂停',
        need_revision: '需修改'
      }
      return statusMap[status] || status
    }
    
    const getProgressStatus = (progress) => {
      if (progress >= 100) return 'success'
      if (progress >= 80) return 'warning'
      if (progress >= 50) return ''
      return 'exception'
    }
    
    const getMilestoneType = (status) => {
      const typeMap = {
        pending: 'primary',
        in_progress: 'warning',
        completed: 'success',
        overdue: 'danger'
      }
      return typeMap[status] || 'info'
    }
    
    const getMilestoneStatusType = (status) => {
      const typeMap = {
        pending: 'info',
        in_progress: 'warning',
        completed: 'success',
        overdue: 'danger'
      }
      return typeMap[status] || 'info'
    }
    
    const getMilestoneStatusText = (status) => {
      const statusMap = {
        pending: '待开始',
        in_progress: '进行中',
        completed: '已完成',
        overdue: '已逾期'
      }
      return statusMap[status] || status
    }
    
    const getExtensionType = (status) => {
      const typeMap = {
        pending: 'warning',
        approved: 'success',
        rejected: 'danger'
      }
      return typeMap[status] || 'info'
    }
    
    const getExtensionStatusType = (status) => {
      const typeMap = {
        pending: 'warning',
        approved: 'success',
        rejected: 'danger'
      }
      return typeMap[status] || 'info'
    }
    
    const getExtensionStatusText = (status) => {
      const statusMap = {
        pending: '待审核',
        approved: '已通过',
        rejected: '已驳回'
      }
      return statusMap[status] || status
    }
    
    const getHistoryType = (status) => {
      return getStatusType(status)
    }
    
    const getReviewType = (status) => {
      const typeMap = {
        approved: 'success',
        rejected: 'danger',
        pending: 'warning'
      }
      return typeMap[status] || 'info'
    }
    
    const getReviewStatusType = (status) => {
      const typeMap = {
        approved: 'success',
        rejected: 'danger',
        pending: 'warning'
      }
      return typeMap[status] || 'info'
    }
    
    const getReviewStatusText = (status) => {
      const statusMap = {
        approved: '通过',
        rejected: '驳回',
        pending: '待审核'
      }
      return statusMap[status] || status
    }
    
    // 生命周期
    onMounted(() => {
      loadProjectDetail()
    })
    
    return {
      // 响应式数据
      loading,
      project,
      milestones,
      allFiles,
      extensions,
      statusHistory,
      reviews,
      
      // 对话框控制
      showProgressDialog,
      showMilestoneDialog,
      showFileDialog,
      showExtensionDialog,
      
      // 表单数据
      progressForm,
      milestoneForm,
      fileForm,
      extensionForm,
      
      // 表单验证规则
      milestoneRules,
      fileRules,
      extensionRules,
      
      // 计算属性
      canEdit,
      canSubmit,
      canUpdateProgress,
      canManageMilestone,
      canUploadFile,
      canReviewFile,
      canApplyExtension,
      
      // 方法
      goBack,
      editProject,
      submitProject,
      updateProgress,
      createMilestone,
      editMilestone,
      completeMilestone,
      handleFileChange,
      uploadFile,
      applyExtension,
      reviewFile,
      downloadFile,
      getFilesByType,
      handleFileTabClick,
      
      // 工具方法
      formatDate,
      getStatusType,
      getStatusText,
      getProgressStatus,
      getMilestoneType,
      getMilestoneStatusType,
      getMilestoneStatusText,
      getExtensionType,
      getExtensionStatusType,
      getExtensionStatusText,
      getHistoryType,
      getReviewType,
      getReviewStatusType,
      getReviewStatusText
    }
  }
}
</script>

<style scoped>
.project-detail {
  padding: 20px;
}

.page-header {
  display: flex;
  align-items: center;
  gap: 15px;
}

.page-title {
  font-size: 18px;
  font-weight: bold;
  flex: 1;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.detail-content {
  margin-top: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.progress-info {
  margin-top: 20px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 10px;
}

.info-item .label {
  font-weight: bold;
  color: #606266;
}

.info-item .value {
  color: #303133;
}

.empty-milestones,
.empty-extensions,
.empty-reviews {
  text-align: center;
  padding: 40px 0;
}

.milestone-item,
.extension-item,
.history-item,
.review-item {
  margin-bottom: 10px;
}

.milestone-header,
.extension-header,
.history-header,
.review-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.milestone-description,
.extension-content,
.change-reason,
.review-comments {
  margin: 10px 0;
  color: #606266;
}

.milestone-progress {
  margin: 15px 0;
}

.milestone-actions {
  margin-top: 15px;
  display: flex;
  gap: 10px;
}

.status-change {
  font-weight: bold;
  color: #303133;
}

.operator,
.reviewer {
  color: #909399;
  font-size: 14px;
}

.el-timeline-item {
  padding-bottom: 20px;
}

.el-timeline-item:last-child {
  padding-bottom: 0;
}
</style> 
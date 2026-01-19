<template>
  <div class="teacher-project-milestones">
    <el-card>
      <template #header>
        <div class="header-content">
          <span>项目里程碑管理</span>
          <el-button type="primary" @click="refreshMilestones">刷新</el-button>
        </div>
      </template>

      <!-- 项目选择 -->
      <el-form :inline="true" class="project-selector">
        <el-form-item label="选择项目">
          <el-select v-model="selectedProjectId" placeholder="请选择项目" @change="loadMilestones">
            <el-option
              v-for="project in projectList"
              :key="project.id"
              :label="project.name"
              :value="project.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="success" @click="loadMilestones" :disabled="!selectedProjectId">
            加载里程碑
          </el-button>
        </el-form-item>
      </el-form>

      <!-- 里程碑管理区域 -->
      <div v-if="selectedProjectId && milestones.length > 0" class="milestones-container">
        <!-- 里程碑统计 -->
        <el-row :gutter="20" class="milestone-stats">
          <el-col :span="6" v-for="stat in milestoneStats" :key="stat.label">
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

        <!-- 里程碑时间线 -->
        <el-card class="timeline-card">
          <template #header>
            <span>里程碑时间线</span>
          </template>
          
          <el-timeline>
            <el-timeline-item
              v-for="milestone in milestones"
              :key="milestone.id"
              :timestamp="milestone.date"
              :type="milestone.completed ? 'success' : 'primary'"
              :color="milestone.completed ? '#67C23A' : '#409EFF'"
            >
              <el-card class="milestone-card" :class="{ completed: milestone.completed }">
                <template #header>
                  <div class="milestone-header">
                    <h4 class="milestone-title">{{ milestone.title }}</h4>
                    <div class="milestone-actions">
                      <el-tag :type="milestone.completed ? 'success' : 'warning'" size="small">
                        {{ milestone.completed ? '已完成' : '进行中' }}
                      </el-tag>
                      <el-button
                        v-if="!milestone.completed"
                        size="small"
                        type="success"
                        @click="approveMilestone(milestone)"
                      >
                        审核通过
                      </el-button>
                      <el-button
                        v-if="!milestone.completed"
                        size="small"
                        type="warning"
                        @click="requestRevision(milestone)"
                      >
                        要求修改
                      </el-button>
                      <el-button size="small" type="primary" @click="provideGuidance(milestone)">
                        提供指导
                      </el-button>
                      <el-button size="small" type="info" @click="viewMilestoneDetail(milestone)">
                        查看详情
                      </el-button>
                    </div>
                  </div>
                </template>
                
                <div class="milestone-content">
                  <p class="milestone-description">{{ milestone.description }}</p>
                  
                  <div class="milestone-details">
                    <p><strong>计划日期:</strong> {{ milestone.date }}</p>
                    <p><strong>负责人:</strong> {{ milestone.assignee || '未分配' }}</p>
                    <p><strong>优先级:</strong> 
                      <el-tag :type="getPriorityType(milestone.priority)" size="small">
                        {{ getPriorityLabel(milestone.priority) }}
                      </el-tag>
                    </p>
                    <p v-if="milestone.completed"><strong>完成时间:</strong> {{ milestone.completedDate }}</p>
                    <p v-if="milestone.completed"><strong>审核人:</strong> {{ milestone.approver || '系统自动' }}</p>
                  </div>

                  <!-- 子任务 -->
                  <div v-if="milestone.subtasks && milestone.subtasks.length > 0" class="subtasks">
                    <h5>子任务</h5>
                    <el-table :data="milestone.subtasks" style="width: 100%" size="small">
                      <el-table-column prop="title" label="任务名称" />
                      <el-table-column prop="status" label="状态" width="100">
                        <template #default="scope">
                          <el-tag :type="getSubtaskStatusType(scope.row.status)" size="small">
                            {{ getSubtaskStatusLabel(scope.row.status) }}
                          </el-tag>
                        </template>
                      </el-table-column>
                      <el-table-column prop="assignee" label="负责人" width="100" />
                      <el-table-column prop="dueDate" label="截止日期" width="120" />
                      <el-table-column label="操作" width="150">
                        <template #default="scope">
                          <el-button size="small" @click="reviewSubtask(scope.row)">审核</el-button>
                          <el-button size="small" type="primary" @click="provideSubtaskGuidance(scope.row)">指导</el-button>
                        </template>
                      </el-table-column>
                    </el-table>
                  </div>

                  <!-- 指导记录 -->
                  <div v-if="milestone.guidanceRecords && milestone.guidanceRecords.length > 0" class="guidance-records">
                    <h5>指导记录</h5>
                    <el-timeline>
                      <el-timeline-item
                        v-for="record in milestone.guidanceRecords"
                        :key="record.id"
                        :timestamp="record.date"
                        type="success"
                      >
                        <h6>{{ record.title }}</h6>
                        <p>{{ record.content }}</p>
                        <p><strong>指导类型:</strong> {{ record.type }}</p>
                        <p><strong>指导人:</strong> {{ record.teacherName }}</p>
                      </el-timeline-item>
                    </el-timeline>
                  </div>
                </div>
              </el-card>
            </el-timeline-item>
          </el-timeline>
        </el-card>
      </div>

      <!-- 空状态 -->
      <el-empty
        v-else-if="selectedProjectId && milestones.length === 0"
        description="该项目暂无里程碑"
      >
        <el-button type="primary" @click="createMilestone">创建第一个里程碑</el-button>
      </el-empty>

      <!-- 项目选择提示 -->
      <el-empty
        v-else
        description="请先选择一个项目"
      />
    </el-card>

    <!-- 里程碑详情对话框 -->
    <el-dialog
      v-model="milestoneDetailVisible"
      title="里程碑详情"
      width="70%"
    >
      <div v-if="currentMilestone" class="milestone-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="里程碑标题">{{ currentMilestone.title }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="currentMilestone.completed ? 'success' : 'warning'" size="small">
              {{ currentMilestone.completed ? '已完成' : '进行中' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="计划日期">{{ currentMilestone.date }}</el-descriptions-item>
          <el-descriptions-item label="负责人">{{ currentMilestone.assignee || '未分配' }}</el-descriptions-item>
          <el-descriptions-item label="优先级">
            <el-tag :type="getPriorityType(currentMilestone.priority)" size="small">
              {{ getPriorityLabel(currentMilestone.priority) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="完成时间">{{ currentMilestone.completedDate || '未完成' }}</el-descriptions-item>
        </el-descriptions>
        
        <el-divider />
        
        <h4>里程碑描述</h4>
        <p>{{ currentMilestone.description }}</p>
        
        <el-divider />
        
        <h4>验收标准</h4>
        <div v-if="currentMilestone.acceptanceCriteria && currentMilestone.acceptanceCriteria.length > 0">
          <ul>
            <li v-for="criteria in currentMilestone.acceptanceCriteria" :key="criteria.id">
              {{ criteria.description }}
              <el-tag 
                v-if="criteria.completed" 
                type="success" 
                size="small" 
                style="margin-left: 10px"
              >
                已完成
              </el-tag>
            </li>
          </ul>
        </div>
        <p v-else>暂无验收标准</p>
        
        <el-divider />
        
        <h4>相关文档</h4>
        <div v-if="currentMilestone.documents && currentMilestone.documents.length > 0">
          <el-table :data="currentMilestone.documents" style="width: 100%">
            <el-table-column prop="name" label="文档名称" />
            <el-table-column prop="type" label="类型" />
            <el-table-column prop="uploadTime" label="上传时间" />
            <el-table-column label="操作" width="150">
              <template #default="scope">
                <el-button size="small" @click="downloadDocument(scope.row)">下载</el-button>
                <el-button size="small" type="primary" @click="previewDocument(scope.row)">预览</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
        <el-empty v-else description="暂无相关文档" />
      </div>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="milestoneDetailVisible = false">关闭</el-button>
          <el-button 
            v-if="currentMilestone && !currentMilestone.completed"
            type="success" 
            @click="approveMilestone(currentMilestone)"
          >
            审核通过
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 指导建议对话框 -->
    <el-dialog
      v-model="guidanceDialogVisible"
      title="提供指导建议"
      width="60%"
    >
      <el-form :model="guidanceForm" :rules="guidanceRules" ref="guidanceFormRef" label-width="100px">
        <el-form-item label="指导类型" prop="type">
          <el-select v-model="guidanceForm.type" placeholder="选择指导类型">
            <el-option label="技术指导" value="technical" />
            <el-option label="进度指导" value="progress" />
            <el-option label="方法指导" value="methodology" />
            <el-option label="问题解决" value="problem_solving" />
            <el-option label="其他" value="other" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="指导标题" prop="title">
          <el-input v-model="guidanceForm.title" placeholder="请输入指导标题" />
        </el-form-item>
        
        <el-form-item label="指导内容" prop="content">
          <el-input
            v-model="guidanceForm.content"
            type="textarea"
            :rows="5"
            placeholder="请详细描述指导内容"
          />
        </el-form-item>
        
        <el-form-item label="具体建议">
          <el-input
            v-model="guidanceForm.suggestions"
            type="textarea"
            :rows="3"
            placeholder="请提供具体的建议措施"
          />
        </el-form-item>
        
        <el-form-item label="预期效果">
          <el-input
            v-model="guidanceForm.expectedOutcome"
            type="textarea"
            :rows="2"
            placeholder="请描述预期达到的效果"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="guidanceDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitGuidance">提交指导</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 审核里程碑对话框 -->
    <el-dialog
      v-model="approvalDialogVisible"
      title="里程碑审核"
      width="50%"
    >
      <el-form :model="approvalForm" :rules="approvalRules" ref="approvalFormRef" label-width="100px">
        <el-form-item label="审核结果" prop="result">
          <el-radio-group v-model="approvalForm.result">
            <el-radio label="approved">通过</el-radio>
            <el-radio label="revision">需要修改</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item label="审核意见" prop="comments">
          <el-input
            v-model="approvalForm.comments"
            type="textarea"
            :rows="4"
            placeholder="请详细说明审核意见"
          />
        </el-form-item>
        
        <el-form-item label="修改要求" v-if="approvalForm.result === 'revision'">
          <el-input
            v-model="approvalForm.revisionRequirements"
            type="textarea"
            :rows="3"
            placeholder="请说明需要修改的具体内容"
          />
        </el-form-item>
        
        <el-form-item label="修改期限" v-if="approvalForm.result === 'revision'">
          <el-date-picker
            v-model="approvalForm.revisionDeadline"
            type="date"
            placeholder="选择修改截止日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        
        <el-form-item label="评分" v-if="approvalForm.result === 'approved'">
          <el-rate
            v-model="approvalForm.score"
            :max="10"
            show-score
            :texts="['很差', '差', '一般', '好', '很好']"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="approvalDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitApproval">提交审核</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Clock, Check, Warning, Document } from '@element-plus/icons-vue'
import { teacherService } from '../../services/teacherService'

// 响应式数据
const selectedProjectId = ref('')
const projectList = ref([])
const milestones = ref([])
const loading = ref(false)
const milestoneDetailVisible = ref(false)
const guidanceDialogVisible = ref(false)
const approvalDialogVisible = ref(false)
const currentMilestone = ref(null)

const guidanceForm = ref({
  type: '',
  title: '',
  content: '',
  suggestions: '',
  expectedOutcome: ''
})

const approvalForm = ref({
  result: '',
  comments: '',
  revisionRequirements: '',
  revisionDeadline: '',
  score: 8
})

const guidanceRules = {
  type: [{ required: true, message: '请选择指导类型', trigger: 'change' }],
  title: [{ required: true, message: '请输入指导标题', trigger: 'blur' }],
  content: [{ required: true, message: '请输入指导内容', trigger: 'blur' }]
}

const approvalRules = {
  result: [{ required: true, message: '请选择审核结果', trigger: 'change' }],
  comments: [{ required: true, message: '请输入审核意见', trigger: 'blur' }]
}

const guidanceFormRef = ref()
const approvalFormRef = ref()

// 里程碑统计
const milestoneStats = computed(() => [
  {
    label: '总里程碑',
    value: milestones.value.length,
    description: '所有里程碑',
    icon: Document,
    type: 'total'
  },
  {
    label: '已完成',
    value: milestones.value.filter(m => m.completed).length,
    description: '审核通过',
    icon: Check,
    type: 'completed'
  },
  {
    label: '进行中',
    value: milestones.value.filter(m => !m.completed).length,
    description: '正在执行',
    icon: Clock,
    type: 'ongoing'
  },
  {
    label: '需要修改',
    value: milestones.value.filter(m => m.status === 'revision').length,
    description: '要求修改',
    icon: Warning,
    type: 'revision'
  }
])

// 获取优先级类型
const getPriorityType = (priority) => {
  const typeMap = {
    low: 'info',
    medium: 'warning',
    high: 'danger',
    urgent: 'danger'
  }
  return typeMap[priority] || 'info'
}

// 获取优先级标签
const getPriorityLabel = (priority) => {
  const labelMap = {
    low: '低',
    medium: '中',
    high: '高',
    urgent: '紧急'
  }
  return labelMap[priority] || '中'
}

// 获取子任务状态类型
const getSubtaskStatusType = (status) => {
  const typeMap = {
    'not_started': 'info',
    'in_progress': 'warning',
    'completed': 'success',
    'on_hold': 'danger'
  }
  return typeMap[status] || 'info'
}

// 获取子任务状态标签
const getSubtaskStatusLabel = (status) => {
  const labelMap = {
    'not_started': '未开始',
    'in_progress': '进行中',
    'completed': '已完成',
    'on_hold': '暂停'
  }
  return labelMap[status] || '未知'
}

// 加载项目列表
const loadProjects = async () => {
  try {
    const response = await teacherService.getGuidedProjects()
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

// 加载里程碑
const loadMilestones = async () => {
  if (!selectedProjectId.value) return
  
  loading.value = true
  try {
    // 这里应该调用实际的API
    // const response = await teacherService.getProjectMilestones(selectedProjectId.value)
    
    // 使用模拟数据
    milestones.value = [
      {
        id: 1,
        title: '需求分析',
        description: '完成用户需求调研和分析，确定系统功能范围',
        date: '2024-01-20',
        assignee: '张三',
        priority: 'high',
        completed: true,
        completedDate: '2024-01-25',
        approver: '李教授',
        acceptanceCriteria: [
          { id: 1, description: '完成用户访谈报告', completed: true },
          { id: 2, description: '完成需求规格说明书', completed: true },
          { id: 3, description: '完成需求评审', completed: true }
        ],
        subtasks: [
          { id: 1, title: '用户访谈', status: 'completed', assignee: '张三', dueDate: '2024-01-18' },
          { id: 2, title: '需求文档编写', status: 'completed', assignee: '张三', dueDate: '2024-01-22' },
          { id: 3, title: '需求评审', status: 'completed', assignee: '张三', dueDate: '2024-01-25' }
        ],
        documents: [
          { name: '需求分析报告.pdf', type: 'PDF', uploadTime: '2024-01-25 15:30:00' },
          { name: '用户访谈记录.docx', type: 'Word', uploadTime: '2024-01-25 15:30:00' }
        ],
        guidanceRecords: [
          { id: 1, title: '需求分析方法指导', content: '建议采用结构化分析方法', type: '方法指导', date: '2024-01-15', teacherName: '李教授' }
        ]
      },
      {
        id: 2,
        title: '系统设计',
        description: '完成系统架构设计、数据库设计和界面设计',
        date: '2024-02-15',
        assignee: '李四',
        priority: 'high',
        completed: true,
        completedDate: '2024-02-20',
        approver: '李教授',
        acceptanceCriteria: [
          { id: 1, description: '完成系统架构设计', completed: true },
          { id: 2, description: '完成数据库设计', completed: true },
          { id: 3, description: '完成界面设计', completed: true }
        ],
        subtasks: [
          { id: 4, title: '架构设计', status: 'completed', assignee: '李四', dueDate: '2024-02-10' },
          { id: 5, title: '数据库设计', status: 'completed', assignee: '李四', dueDate: '2024-02-12' },
          { id: 6, title: '界面设计', status: 'completed', assignee: '李四', dueDate: '2024-02-15' }
        ],
        documents: [
          { name: '系统架构设计.pdf', type: 'PDF', uploadTime: '2024-02-20 14:20:00' },
          { name: '数据库设计.sql', type: 'SQL', uploadTime: '2024-02-20 14:20:00' }
        ],
        guidanceRecords: [
          { id: 1, title: '架构设计指导', content: '建议采用微服务架构', type: '技术指导', date: '2024-02-05', teacherName: '李教授' }
        ]
      },
      {
        id: 3,
        title: '功能开发',
        description: '核心功能模块的编码实现',
        date: '2024-05-15',
        assignee: '王五',
        priority: 'medium',
        completed: false,
        acceptanceCriteria: [
          { id: 1, description: '完成用户管理模块', completed: false },
          { id: 2, description: '完成核心业务模块', completed: false },
          { id: 3, description: '完成数据统计模块', completed: false }
        ],
        subtasks: [
          { id: 7, title: '用户管理模块', status: 'in_progress', assignee: '王五', dueDate: '2024-05-10' },
          { id: 8, title: '核心业务模块', status: 'not_started', assignee: '王五', dueDate: '2024-05-20' },
          { id: 9, title: '数据统计模块', status: 'not_started', assignee: '王五', dueDate: '2024-05-25' }
        ],
        documents: [
          { name: '开发计划.docx', type: 'Word', uploadTime: '2024-05-01 10:00:00' }
        ],
        guidanceRecords: [
          { id: 1, title: '开发方法指导', content: '建议采用敏捷开发方法', type: '方法指导', date: '2024-05-05', teacherName: '李教授' }
        ]
      },
      {
        id: 4,
        title: '系统测试',
        description: '功能测试、性能测试和用户验收测试',
        date: '2024-06-15',
        assignee: '赵六',
        priority: 'medium',
        completed: false,
        acceptanceCriteria: [
          { id: 1, description: '完成功能测试', completed: false },
          { id: 2, description: '完成性能测试', completed: false },
          { id: 3, description: '完成用户验收测试', completed: false }
        ],
        subtasks: [
          { id: 10, title: '功能测试', status: 'not_started', assignee: '赵六', dueDate: '2024-06-10' },
          { id: 11, title: '性能测试', status: 'not_started', assignee: '赵六', dueDate: '2024-06-12' },
          { id: 12, title: '用户验收测试', status: 'not_started', assignee: '赵六', dueDate: '2024-06-15' }
        ],
        documents: [],
        guidanceRecords: []
      }
    ]
  } catch (error) {
    console.error('加载里程碑失败:', error)
    ElMessage.error('加载里程碑失败')
  } finally {
    loading.value = false
  }
}

// 查看里程碑详情
const viewMilestoneDetail = (milestone) => {
  currentMilestone.value = milestone
  milestoneDetailVisible.value = true
}

// 提供指导建议
const provideGuidance = (milestone) => {
  currentMilestone.value = milestone
  guidanceForm.value = {
    type: '',
    title: '',
    content: '',
    suggestions: '',
    expectedOutcome: ''
  }
  guidanceDialogVisible.value = true
}

// 提交指导建议
const submitGuidance = async () => {
  try {
    await guidanceFormRef.value.validate()
    
    // 这里应该调用实际的API
    // await teacherService.submitMilestoneGuidance(currentMilestone.value.id, guidanceForm.value)
    
    // 添加指导记录到里程碑
    const newGuidance = {
      id: Date.now(),
      title: guidanceForm.value.title,
      content: guidanceForm.value.content,
      type: getGuidanceTypeLabel(guidanceForm.value.type),
      date: new Date().toLocaleString('zh-CN'),
      teacherName: '当前用户'
    }
    
    if (!currentMilestone.value.guidanceRecords) {
      currentMilestone.value.guidanceRecords = []
    }
    currentMilestone.value.guidanceRecords.push(newGuidance)
    
    ElMessage.success('指导建议提交成功')
    guidanceDialogVisible.value = false
  } catch (error) {
    console.error('提交指导建议失败:', error)
    ElMessage.error('提交失败')
  }
}

// 审核里程碑
const approveMilestone = (milestone) => {
  currentMilestone.value = milestone
  approvalForm.value = {
    result: 'approved',
    comments: '',
    revisionRequirements: '',
    revisionDeadline: '',
    score: 8
  }
  approvalDialogVisible.value = true
}

// 要求修改
const requestRevision = (milestone) => {
  currentMilestone.value = milestone
  approvalForm.value = {
    result: 'revision',
    comments: '',
    revisionRequirements: '',
    revisionDeadline: '',
    score: 0
  }
  approvalDialogVisible.value = true
}

// 提交审核
const submitApproval = async () => {
  try {
    await approvalFormRef.value.validate()
    
    // 这里应该调用实际的API
    // await teacherService.submitMilestoneApproval(currentMilestone.value.id, approvalForm.value)
    
    if (approvalForm.value.result === 'approved') {
      currentMilestone.value.completed = true
      currentMilestone.value.completedDate = new Date().toLocaleString('zh-CN')
      currentMilestone.value.approver = '当前用户'
      ElMessage.success('里程碑审核通过')
    } else {
      currentMilestone.value.status = 'revision'
      ElMessage.success('里程碑需要修改')
    }
    
    approvalDialogVisible.value = false
  } catch (error) {
    console.error('提交审核失败:', error)
    ElMessage.error('提交失败')
  }
}

// 审核子任务
const reviewSubtask = (subtask) => {
  ElMessage.info(`审核子任务: ${subtask.title}`)
}

// 为子任务提供指导
const provideSubtaskGuidance = (subtask) => {
  ElMessage.info(`为子任务提供指导: ${subtask.title}`)
}

// 下载文档
const downloadDocument = (document) => {
  ElMessage.success(`开始下载文档: ${document.name}`)
}

// 预览文档
const previewDocument = (document) => {
  ElMessage.info(`预览文档: ${document.name}`)
}

// 创建里程碑
const createMilestone = () => {
  ElMessage.info('跳转到里程碑创建页面')
}

// 获取指导类型标签
const getGuidanceTypeLabel = (type) => {
  const labelMap = {
    'technical': '技术指导',
    'progress': '进度指导',
    'methodology': '方法指导',
    'problem_solving': '问题解决',
    'other': '其他'
  }
  return labelMap[type] || '其他'
}

// 刷新里程碑
const refreshMilestones = () => {
  loadMilestones()
  ElMessage.success('里程碑列表已刷新')
}

// 组件挂载时加载数据
onMounted(() => {
  loadProjects()
})
</script>

<style scoped>
.teacher-project-milestones {
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

.milestones-container {
  margin-top: 20px;
}

.milestone-stats {
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

.stat-card.completed {
  border-left: 4px solid #4facfe;
}

.stat-card.ongoing {
  border-left: 4px solid #f093fb;
}

.stat-card.revision {
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

.stat-icon.completed {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.stat-icon.ongoing {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.stat-icon.revision {
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

.timeline-card {
  margin-top: 20px;
}

.milestone-card {
  margin-bottom: 10px;
  transition: all 0.3s ease;
}

.milestone-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.milestone-card.completed {
  opacity: 0.8;
  background-color: #f0f9ff;
}

.milestone-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.milestone-title {
  margin: 0;
  color: #2c3e50;
}

.milestone-actions {
  display: flex;
  gap: 8px;
  align-items: center;
}

.milestone-content {
  padding: 10px 0;
}

.milestone-description {
  color: #606266;
  margin-bottom: 15px;
  line-height: 1.6;
}

.milestone-details p {
  margin: 8px 0;
  color: #606266;
}

.subtasks {
  margin-top: 15px;
  padding-top: 15px;
  border-top: 1px solid #ebeef5;
}

.subtasks h5 {
  margin: 0 0 10px 0;
  color: #2c3e50;
}

.guidance-records {
  margin-top: 15px;
  padding-top: 15px;
  border-top: 1px solid #ebeef5;
}

.guidance-records h5 {
  margin: 0 0 10px 0;
  color: #2c3e50;
}

.guidance-records h6 {
  margin: 0 0 5px 0;
  color: #2c3e50;
  font-size: 14px;
}

.milestone-detail {
  padding: 20px;
}

.milestone-detail h4 {
  margin: 20px 0 10px 0;
  color: #2c3e50;
  font-size: 16px;
}

.milestone-detail p {
  margin: 10px 0;
  color: #606266;
  line-height: 1.6;
}

.milestone-detail ul {
  margin: 10px 0;
  padding-left: 20px;
}

.milestone-detail li {
  margin: 5px 0;
  color: #606266;
}

.dialog-footer {
  text-align: right;
}

:deep(.el-timeline-item__node) {
  background-color: #409EFF;
}

:deep(.el-timeline-item__node--success) {
  background-color: #67C23A;
}

:deep(.el-timeline-item__wrapper) {
  padding-left: 20px;
}
</style> 
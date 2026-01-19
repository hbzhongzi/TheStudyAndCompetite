<template>
  <div class="project-milestones">
    <el-card>
      <template #header>
        <div class="header-content">
          <span>项目里程碑管理</span>
          <el-button type="primary" @click="showAddMilestoneDialog">添加里程碑</el-button>
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

      <!-- 里程碑列表 -->
      <div v-if="selectedProjectId && milestones.length > 0" class="milestones-container">
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
                      @click="completeMilestone(milestone)"
                    >
                      标记完成
                    </el-button>
                    <el-button size="small" type="primary" @click="editMilestone(milestone)">
                      编辑
                    </el-button>
                    <el-button size="small" type="danger" @click="deleteMilestone(milestone)">
                      删除
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
                </div>

                <!-- 子任务 -->
                <div v-if="milestone.subtasks && milestone.subtasks.length > 0" class="subtasks">
                  <h5>子任务</h5>
                  <el-checkbox-group v-model="milestone.completedSubtasks">
                    <div v-for="subtask in milestone.subtasks" :key="subtask.id" class="subtask-item">
                      <el-checkbox 
                        :label="subtask.id"
                        :disabled="milestone.completed"
                        @change="updateSubtaskProgress(milestone)"
                      >
                        {{ subtask.title }}
                      </el-checkbox>
                    </div>
                  </el-checkbox-group>
                  <p class="subtask-progress">
                    子任务进度: {{ milestone.completedSubtasks.length }}/{{ milestone.subtasks.length }}
                  </p>
                </div>
              </div>
            </el-card>
          </el-timeline-item>
        </el-timeline>
      </div>

      <!-- 空状态 -->
      <el-empty
        v-else-if="selectedProjectId && milestones.length === 0"
        description="该项目暂无里程碑"
      >
        <el-button type="primary" @click="showAddMilestoneDialog">添加第一个里程碑</el-button>
      </el-empty>

      <!-- 项目选择提示 -->
      <el-empty
        v-else
        description="请先选择一个项目"
      />
    </el-card>

    <!-- 添加/编辑里程碑对话框 -->
    <el-dialog
      v-model="milestoneDialogVisible"
      :title="isEditing ? '编辑里程碑' : '添加里程碑'"
      width="50%"
    >
      <el-form :model="milestoneForm" :rules="milestoneRules" ref="milestoneFormRef" label-width="100px">
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
        
        <el-form-item label="计划日期" prop="date">
          <el-date-picker
            v-model="milestoneForm.date"
            type="date"
            placeholder="选择计划日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        
        <el-form-item label="负责人" prop="assignee">
          <el-input v-model="milestoneForm.assignee" placeholder="请输入负责人姓名" />
        </el-form-item>
        
        <el-form-item label="优先级" prop="priority">
          <el-select v-model="milestoneForm.priority" placeholder="选择优先级">
            <el-option label="低" value="low" />
            <el-option label="中" value="medium" />
            <el-option label="高" value="high" />
            <el-option label="紧急" value="urgent" />
          </el-select>
        </el-form-item>

        <!-- 子任务 -->
        <el-form-item label="子任务">
          <div v-for="(subtask, index) in milestoneForm.subtasks" :key="index" class="subtask-input">
            <el-input
              v-model="milestoneForm.subtasks[index]"
              placeholder="请输入子任务"
              style="width: 80%"
            />
            <el-button
              type="danger"
              size="small"
              @click="removeSubtask(index)"
              style="margin-left: 10px"
            >
              删除
            </el-button>
          </div>
          <el-button type="primary" size="small" @click="addSubtask">添加子任务</el-button>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="milestoneDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitMilestone">确定</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { studentService } from '../../services/studentService'

// 响应式数据
const selectedProjectId = ref('')
const projectList = ref([])
const milestones = ref([])
const milestoneDialogVisible = ref(false)
const isEditing = ref(false)
const currentMilestone = ref(null)

const milestoneForm = ref({
  title: '',
  description: '',
  date: '',
  assignee: '',
  priority: 'medium',
  subtasks: ['']
})

const milestoneRules = {
  title: [{ required: true, message: '请输入里程碑标题', trigger: 'blur' }],
  description: [{ required: true, message: '请输入里程碑描述', trigger: 'blur' }],
  date: [{ required: true, message: '请选择计划日期', trigger: 'change' }],
  priority: [{ required: true, message: '请选择优先级', trigger: 'change' }]
}

const milestoneFormRef = ref()

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

// 加载里程碑
const loadMilestones = async () => {
  if (!selectedProjectId.value) return
  
  try {
    // 这里应该调用实际的API
    // const response = await studentService.getProjectMilestones(selectedProjectId.value)
    
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
        subtasks: [
          { id: 1, title: '用户访谈' },
          { id: 2, title: '需求文档编写' },
          { id: 3, title: '需求评审' }
        ],
        completedSubtasks: [1, 2, 3]
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
        subtasks: [
          { id: 4, title: '架构设计' },
          { id: 5, title: '数据库设计' },
          { id: 6, title: '界面设计' }
        ],
        completedSubtasks: [4, 5, 6]
      },
      {
        id: 3,
        title: '功能开发',
        description: '核心功能模块的编码实现',
        date: '2024-05-15',
        assignee: '王五',
        priority: 'medium',
        completed: false,
        subtasks: [
          { id: 7, title: '用户管理模块' },
          { id: 8, title: '核心业务模块' },
          { id: 9, title: '数据统计模块' }
        ],
        completedSubtasks: [7]
      },
      {
        id: 4,
        title: '系统测试',
        description: '功能测试、性能测试和用户验收测试',
        date: '2024-06-15',
        assignee: '赵六',
        priority: 'medium',
        completed: false,
        subtasks: [
          { id: 10, title: '单元测试' },
          { id: 11, title: '集成测试' },
          { id: 12, title: '用户验收测试' }
        ],
        completedSubtasks: []
      }
    ]
  } catch (error) {
    console.error('加载里程碑失败:', error)
    ElMessage.error('加载里程碑失败')
  }
}

// 显示添加里程碑对话框
const showAddMilestoneDialog = () => {
  if (!selectedProjectId.value) {
    ElMessage.warning('请先选择一个项目')
    return
  }
  
  isEditing.value = false
  currentMilestone.value = null
  resetMilestoneForm()
  milestoneDialogVisible.value = true
}

// 编辑里程碑
const editMilestone = (milestone) => {
  isEditing.value = true
  currentMilestone.value = milestone
  milestoneForm.value = {
    title: milestone.title,
    description: milestone.description,
    date: milestone.date,
    assignee: milestone.assignee,
    priority: milestone.priority,
    subtasks: milestone.subtasks ? milestone.subtasks.map(s => s.title) : ['']
  }
  milestoneDialogVisible.value = true
}

// 删除里程碑
const deleteMilestone = async (milestone) => {
  try {
    await ElMessageBox.confirm('确定要删除这个里程碑吗？', '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    // 这里应该调用实际的API
    // await studentService.deleteMilestone(milestone.id)
    
    const index = milestones.value.findIndex(m => m.id === milestone.id)
    if (index > -1) {
      milestones.value.splice(index, 1)
    }
    
    ElMessage.success('里程碑删除成功')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除里程碑失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 完成里程碑
const completeMilestone = async (milestone) => {
  try {
    await ElMessageBox.confirm('确定要将此里程碑标记为完成吗？', '确认完成', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'info'
    })
    
    // 这里应该调用实际的API
    // await studentService.completeMilestone(milestone.id)
    
    milestone.completed = true
    milestone.completedDate = new Date().toISOString().split('T')[0]
    
    ElMessage.success('里程碑已完成')
  } catch (error) {
    if (error !== 'cancel') {
      console.error('完成里程碑失败:', error)
      ElMessage.error('操作失败')
    }
  }
}

// 更新子任务进度
const updateSubtaskProgress = (milestone) => {
  // 这里可以调用API更新子任务进度
  console.log('子任务进度更新:', milestone.completedSubtasks)
}

// 添加子任务
const addSubtask = () => {
  milestoneForm.value.subtasks.push('')
}

// 删除子任务
const removeSubtask = (index) => {
  milestoneForm.value.subtasks.splice(index, 1)
}

// 重置里程碑表单
const resetMilestoneForm = () => {
  milestoneForm.value = {
    title: '',
    description: '',
    date: '',
    assignee: '',
    priority: 'medium',
    subtasks: ['']
  }
}

// 提交里程碑
const submitMilestone = async () => {
  try {
    await milestoneFormRef.value.validate()
    
    if (isEditing.value) {
      // 编辑里程碑
      Object.assign(currentMilestone.value, {
        title: milestoneForm.value.title,
        description: milestoneForm.value.description,
        date: milestoneForm.value.date,
        assignee: milestoneForm.value.assignee,
        priority: milestoneForm.value.priority,
        subtasks: milestoneForm.value.subtasks.filter(s => s.trim()).map((title, index) => ({
          id: currentMilestone.value.subtasks[index]?.id || Date.now() + index,
          title: title.trim()
        }))
      })
      
      ElMessage.success('里程碑更新成功')
    } else {
      // 添加里程碑
      const newMilestone = {
        id: Date.now(),
        title: milestoneForm.value.title,
        description: milestoneForm.value.description,
        date: milestoneForm.value.date,
        assignee: milestoneForm.value.assignee,
        priority: milestoneForm.value.priority,
        completed: false,
        subtasks: milestoneForm.value.subtasks.filter(s => s.trim()).map((title, index) => ({
          id: Date.now() + index,
          title: title.trim()
        })),
        completedSubtasks: []
      }
      
      milestones.value.push(newMilestone)
      ElMessage.success('里程碑添加成功')
    }
    
    milestoneDialogVisible.value = false
  } catch (error) {
    console.error('提交里程碑失败:', error)
    ElMessage.error('操作失败')
  }
}

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

// 组件挂载时加载数据
onMounted(() => {
  loadProjects()
})
</script>

<style scoped>
.project-milestones {
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

.subtask-item {
  margin: 8px 0;
}

.subtask-progress {
  margin-top: 10px;
  color: #909399;
  font-size: 12px;
}

.subtask-input {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
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
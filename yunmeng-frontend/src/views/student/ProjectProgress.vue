<template>
  <div class="project-progress">
    <el-card>
      <template #header>
        <div class="header-content">
          <span>项目进度跟踪</span>
          <el-button type="primary" @click="showProgressReportDialog">生成进度报告</el-button>
        </div>
      </template>

      <!-- 项目选择 -->
      <el-form :inline="true" class="project-selector">
        <el-form-item label="选择项目">
          <el-select v-model="selectedProjectId" placeholder="请选择项目" @change="loadProgress">
            <el-option
              v-for="project in projectList"
              :key="project.id"
              :label="project.name"
              :value="project.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="success" @click="loadProgress" :disabled="!selectedProjectId">
            加载进度数据
          </el-button>
        </el-form-item>
      </el-form>

      <!-- 进度跟踪区域 -->
      <div v-if="selectedProjectId && currentProgress" class="progress-container">
        <!-- 总体进度概览 -->
        <el-row :gutter="20" class="progress-overview">
          <el-col :span="6">
            <el-card class="overview-card">
              <div class="overview-content">
                <h3>总体进度</h3>
                <div class="progress-circle">
                  <el-progress
                    type="circle"
                    :percentage="currentProgress.overallProgress"
                    :color="getProgressColor(currentProgress.overallProgress)"
                    :width="80"
                  />
                </div>
                <p class="progress-text">{{ currentProgress.overallProgress }}%</p>
              </div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card class="overview-card">
              <div class="overview-content">
                <h3>计划进度</h3>
                <div class="progress-circle">
                  <el-progress
                    type="circle"
                    :percentage="currentProgress.plannedProgress"
                    :color="getProgressColor(currentProgress.plannedProgress)"
                    :width="80"
                  />
                </div>
                <p class="progress-text">{{ currentProgress.plannedProgress }}%</p>
              </div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card class="overview-card">
              <div class="overview-content">
                <h3>剩余天数</h3>
                <div class="days-remaining">
                  <span class="days-number">{{ currentProgress.daysRemaining }}</span>
                  <span class="days-label">天</span>
                </div>
                <p class="deadline-text">截止: {{ currentProgress.deadline }}</p>
              </div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card class="overview-card">
              <div class="overview-content">
                <h3>状态评估</h3>
                <div class="status-assessment">
                  <el-tag :type="getStatusType(currentProgress.status)" size="large">
                    {{ getStatusLabel(currentProgress.status) }}
                  </el-tag>
                </div>
                <p class="status-desc">{{ getStatusDescription(currentProgress.status) }}</p>
              </div>
            </el-card>
          </el-col>
        </el-row>

        <!-- 进度趋势图 -->
        <el-card class="trend-chart-card">
          <template #header>
            <span>进度趋势分析</span>
          </template>
          <div class="chart-container">
            <div ref="trendChartRef" class="trend-chart"></div>
          </div>
        </el-card>

        <!-- 模块进度详情 -->
        <el-card class="module-progress-card">
          <template #header>
            <span>模块进度详情</span>
          </template>
          <el-table :data="currentProgress.modules" style="width: 100%">
            <el-table-column prop="name" label="模块名称" min-width="150" />
            <el-table-column prop="plannedProgress" label="计划进度" width="120">
              <template #default="scope">
                <el-progress 
                  :percentage="scope.row.plannedProgress" 
                  :status="getProgressStatus(scope.row.plannedProgress)"
                />
              </template>
            </el-table-column>
            <el-table-column prop="actualProgress" label="实际进度" width="120">
              <template #default="scope">
                <el-progress 
                  :percentage="scope.row.actualProgress" 
                  :status="getProgressStatus(scope.row.actualProgress)"
                />
              </template>
            </el-table-column>
            <el-table-column prop="deviation" label="偏差" width="100">
              <template #default="scope">
                <el-tag :type="getDeviationType(scope.row.deviation)" size="small">
                  {{ scope.row.deviation > 0 ? '+' : '' }}{{ scope.row.deviation }}%
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="scope">
                <el-tag :type="getModuleStatusType(scope.row.status)" size="small">
                  {{ getModuleStatusLabel(scope.row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="startDate" label="开始时间" width="120" />
            <el-table-column prop="endDate" label="结束时间" width="120" />
            <el-table-column label="操作" width="150" fixed="right">
              <template #default="scope">
                <el-button size="small" @click="viewModuleDetail(scope.row)">查看详情</el-button>
                <el-button size="small" type="primary" @click="updateModuleProgress(scope.row)">更新进度</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>

        <!-- 里程碑进度 -->
        <el-card class="milestone-progress-card">
          <template #header>
            <span>里程碑进度</span>
          </template>
          <el-timeline>
            <el-timeline-item
              v-for="milestone in currentProgress.milestones"
              :key="milestone.id"
              :timestamp="milestone.date"
              :type="milestone.completed ? 'success' : 'primary'"
              :color="milestone.completed ? '#67C23A' : '#409EFF'"
            >
              <el-card class="milestone-card" :class="{ completed: milestone.completed }">
                <div class="milestone-header">
                  <h4>{{ milestone.title }}</h4>
                  <el-tag :type="milestone.completed ? 'success' : 'warning'" size="small">
                    {{ milestone.completed ? '已完成' : '进行中' }}
                  </el-tag>
                </div>
                <p class="milestone-description">{{ milestone.description }}</p>
                <div class="milestone-progress">
                  <el-progress 
                    :percentage="milestone.progress" 
                    :status="milestone.completed ? 'success' : ''"
                  />
                  <span class="progress-text">{{ milestone.progress }}%</span>
                </div>
                <div class="milestone-details">
                  <p><strong>负责人:</strong> {{ milestone.assignee }}</p>
                  <p><strong>计划完成:</strong> {{ milestone.date }}</p>
                  <p v-if="milestone.completed"><strong>实际完成:</strong> {{ milestone.completedDate }}</p>
                </div>
              </el-card>
            </el-timeline-item>
          </el-timeline>
        </el-card>

        <!-- 风险与问题 -->
        <el-card class="risks-issues-card">
          <template #header>
            <span>风险与问题跟踪</span>
          </template>
          <el-row :gutter="20">
            <el-col :span="12">
              <h4>风险项</h4>
              <div v-if="currentProgress.risks && currentProgress.risks.length > 0">
                <div v-for="risk in currentProgress.risks" :key="risk.id" class="risk-item">
                  <el-tag :type="getRiskLevelType(risk.level)" size="small">{{ getRiskLevelLabel(risk.level) }}</el-tag>
                  <span class="risk-title">{{ risk.title }}</span>
                  <p class="risk-description">{{ risk.description }}</p>
                  <p class="risk-mitigation"><strong>缓解措施:</strong> {{ risk.mitigation }}</p>
                </div>
              </div>
              <el-empty v-else description="暂无风险项" />
            </el-col>
            <el-col :span="12">
              <h4>问题项</h4>
              <div v-if="currentProgress.issues && currentProgress.issues.length > 0">
                <div v-for="issue in currentProgress.issues" :key="issue.id" class="issue-item">
                  <el-tag :type="getIssuePriorityType(issue.priority)" size="small">{{ getIssuePriorityLabel(issue.priority) }}</el-tag>
                  <span class="issue-title">{{ issue.title }}</span>
                  <p class="issue-description">{{ issue.description }}</p>
                  <p class="issue-solution"><strong>解决方案:</strong> {{ issue.solution }}</p>
                </div>
              </div>
              <el-empty v-else description="暂无问题项" />
            </el-col>
          </el-row>
        </el-card>
      </div>

      <!-- 项目选择提示 -->
      <el-empty
        v-else-if="!selectedProjectId"
        description="请先选择一个项目"
      />
      <el-empty
        v-else
        description="加载进度数据中..."
      />
    </el-card>

    <!-- 更新模块进度对话框 -->
    <el-dialog
      v-model="progressUpdateDialogVisible"
      title="更新模块进度"
      width="50%"
    >
      <el-form :model="progressUpdateForm" :rules="progressUpdateRules" ref="progressUpdateFormRef" label-width="100px">
        <el-form-item label="模块名称">
          <el-input v-model="progressUpdateForm.moduleName" disabled />
        </el-form-item>
        
        <el-form-item label="当前进度" prop="progress">
          <el-slider v-model="progressUpdateForm.progress" :min="0" :max="100" :step="5" show-input />
        </el-form-item>
        
        <el-form-item label="进度说明" prop="description">
          <el-input
            v-model="progressUpdateForm.description"
            type="textarea"
            :rows="3"
            placeholder="请描述当前进度情况"
          />
        </el-form-item>
        
        <el-form-item label="遇到的问题">
          <el-input
            v-model="progressUpdateForm.issues"
            type="textarea"
            :rows="2"
            placeholder="请描述遇到的问题（可选）"
          />
        </el-form-item>
        
        <el-form-item label="下一步计划">
          <el-input
            v-model="progressUpdateForm.nextSteps"
            type="textarea"
            :rows="2"
            placeholder="请描述下一步计划（可选）"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="progressUpdateDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="submitProgressUpdate">提交更新</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 进度报告对话框 -->
    <el-dialog
      v-model="progressReportDialogVisible"
      title="生成进度报告"
      width="60%"
    >
      <el-form :model="reportForm" label-width="120px">
        <el-form-item label="报告类型">
          <el-radio-group v-model="reportForm.type">
            <el-radio label="weekly">周报</el-radio>
            <el-radio label="monthly">月报</el-radio>
            <el-radio label="milestone">里程碑报告</el-radio>
            <el-radio label="custom">自定义</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item label="报告时间范围" v-if="reportForm.type === 'custom'">
          <el-date-picker
            v-model="reportForm.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        
        <el-form-item label="包含内容">
          <el-checkbox-group v-model="reportForm.content">
            <el-checkbox label="overview">总体进度概览</el-checkbox>
            <el-checkbox label="modules">模块进度详情</el-checkbox>
            <el-checkbox label="milestones">里程碑完成情况</el-checkbox>
            <el-checkbox label="risks">风险与问题</el-checkbox>
            <el-checkbox label="trends">进度趋势分析</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        
        <el-form-item label="报告格式">
          <el-radio-group v-model="reportForm.format">
            <el-radio label="pdf">PDF</el-radio>
            <el-radio label="excel">Excel</el-radio>
            <el-radio label="word">Word</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="progressReportDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="generateProgressReport" :loading="generatingReport">生成报告</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import * as echarts from 'echarts'
import { studentService } from '../../services/studentService'

// 响应式数据
const selectedProjectId = ref('')
const projectList = ref([])
const currentProgress = ref(null)
const loading = ref(false)
const progressUpdateDialogVisible = ref(false)
const progressReportDialogVisible = ref(false)
const currentModule = ref(null)
const generatingReport = ref(false)

const progressUpdateForm = ref({
  moduleName: '',
  progress: 0,
  description: '',
  issues: '',
  nextSteps: ''
})

const progressUpdateRules = {
  progress: [{ required: true, message: '请设置进度', trigger: 'change' }],
  description: [{ required: true, message: '请输入进度说明', trigger: 'blur' }]
}

const reportForm = ref({
  type: 'weekly',
  dateRange: [],
  content: ['overview', 'modules', 'milestones'],
  format: 'pdf'
})

const progressUpdateFormRef = ref()
const trendChartRef = ref()
let trendChart = null

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

// 加载进度数据
const loadProgress = async () => {
  if (!selectedProjectId.value) return
  
  loading.value = true
  try {
    // 这里应该调用实际的API
    // const response = await studentService.getProjectProgress(selectedProjectId.value)
    
    // 使用模拟数据
    currentProgress.value = {
      projectId: selectedProjectId.value,
      projectName: projectList.value.find(p => p.id === selectedProjectId.value)?.name,
      overallProgress: 65,
      plannedProgress: 70,
      daysRemaining: 45,
      deadline: '2024-07-15',
      status: 'on_track',
      modules: [
        {
          id: 1,
          name: '需求分析',
          plannedProgress: 100,
          actualProgress: 100,
          deviation: 0,
          status: 'completed',
          startDate: '2024-01-15',
          endDate: '2024-01-30'
        },
        {
          id: 2,
          name: '系统设计',
          plannedProgress: 100,
          actualProgress: 100,
          deviation: 0,
          status: 'completed',
          startDate: '2024-02-01',
          endDate: '2024-02-28'
        },
        {
          id: 3,
          name: '数据库开发',
          plannedProgress: 80,
          actualProgress: 75,
          deviation: -5,
          status: 'in_progress',
          startDate: '2024-03-01',
          endDate: '2024-04-15'
        },
        {
          id: 4,
          name: '后端开发',
          plannedProgress: 60,
          actualProgress: 50,
          deviation: -10,
          status: 'in_progress',
          startDate: '2024-03-15',
          endDate: '2024-05-30'
        },
        {
          id: 5,
          name: '前端开发',
          plannedProgress: 40,
          actualProgress: 35,
          deviation: -5,
          status: 'in_progress',
          startDate: '2024-04-01',
          endDate: '2024-06-15'
        },
        {
          id: 6,
          name: '系统测试',
          plannedProgress: 20,
          actualProgress: 0,
          deviation: -20,
          status: 'not_started',
          startDate: '2024-06-01',
          endDate: '2024-07-01'
        }
      ],
      milestones: [
        {
          id: 1,
          title: '需求确认',
          description: '完成用户需求调研和分析',
          date: '2024-01-30',
          progress: 100,
          completed: true,
          completedDate: '2024-01-28',
          assignee: '张三'
        },
        {
          id: 2,
          title: '架构设计',
          description: '完成系统架构和数据库设计',
          date: '2024-02-28',
          progress: 100,
          completed: true,
          completedDate: '2024-02-25',
          assignee: '李四'
        },
        {
          id: 3,
          title: '核心开发',
          description: '完成核心功能模块开发',
          date: '2024-05-30',
          progress: 60,
          completed: false,
          assignee: '王五'
        },
        {
          id: 4,
          title: '系统集成',
          description: '完成各模块集成测试',
          date: '2024-06-30',
          progress: 0,
          completed: false,
          assignee: '赵六'
        }
      ],
      risks: [
        {
          id: 1,
          title: '技术难点',
          description: '某些技术模块实现复杂度较高',
          level: 'medium',
          mitigation: '引入专家顾问，增加技术调研时间'
        },
        {
          id: 2,
          title: '人员变动',
          description: '关键开发人员可能离职',
          level: 'low',
          mitigation: '做好知识转移，建立备份人员'
        }
      ],
      issues: [
        {
          id: 1,
          title: '数据库性能问题',
          description: '大数据量查询性能不理想',
          priority: 'high',
          solution: '优化查询语句，添加索引，考虑分库分表'
        },
        {
          id: 2,
          title: '接口兼容性',
          description: '新旧接口版本兼容性问题',
          priority: 'medium',
          solution: '制定接口版本管理策略，逐步迁移'
        }
      ]
    }
    
    // 初始化图表
    await nextTick()
    initTrendChart()
  } catch (error) {
    console.error('加载进度数据失败:', error)
    ElMessage.error('加载进度数据失败')
  } finally {
    loading.value = false
  }
}

// 初始化趋势图表
const initTrendChart = () => {
  if (!trendChartRef.value || !currentProgress.value) return
  
  trendChart = echarts.init(trendChartRef.value)
  
  const option = {
    title: {
      text: '项目进度趋势',
      left: 'center'
    },
    tooltip: {
      trigger: 'axis',
      formatter: function(params) {
        let result = params[0].name + '<br/>'
        params.forEach(param => {
          result += param.marker + param.seriesName + ': ' + param.value + '%<br/>'
        })
        return result
      }
    },
    legend: {
      data: ['计划进度', '实际进度'],
      bottom: 10
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '15%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: ['1月', '2月', '3月', '4月', '5月', '6月', '7月']
    },
    yAxis: {
      type: 'value',
      max: 100,
      axisLabel: {
        formatter: '{value}%'
      }
    },
    series: [
      {
        name: '计划进度',
        type: 'line',
        smooth: true,
        data: [15, 30, 45, 60, 75, 85, 100],
        itemStyle: {
          color: '#409EFF'
        }
      },
      {
        name: '实际进度',
        type: 'line',
        smooth: true,
        data: [15, 30, 45, 55, 65, 75, 85],
        itemStyle: {
          color: '#67C23A'
        }
      }
    ]
  }
  
  trendChart.setOption(option)
}

// 查看模块详情
const viewModuleDetail = (module) => {
  ElMessage.info(`查看模块详情: ${module.name}`)
}

// 更新模块进度
const updateModuleProgress = (module) => {
  currentModule.value = module
  progressUpdateForm.value = {
    moduleName: module.name,
    progress: module.actualProgress,
    description: '',
    issues: '',
    nextSteps: ''
  }
  progressUpdateDialogVisible.value = true
}

// 提交进度更新
const submitProgressUpdate = async () => {
  try {
    await progressUpdateFormRef.value.validate()
    
    // 这里应该调用实际的API
    // await studentService.updateModuleProgress(currentModule.value.id, progressUpdateForm.value)
    
    // 更新本地数据
    currentModule.value.actualProgress = progressUpdateForm.value.progress
    currentModule.value.deviation = progressUpdateForm.value.progress - currentModule.value.plannedProgress
    
    // 重新计算总体进度
    const totalProgress = currentProgress.value.modules.reduce((sum, module) => sum + module.actualProgress, 0)
    currentProgress.value.overallProgress = Math.round(totalProgress / currentProgress.value.modules.length)
    
    ElMessage.success('进度更新成功')
    progressUpdateDialogVisible.value = false
    
    // 重新渲染图表
    if (trendChart) {
      initTrendChart()
    }
  } catch (error) {
    console.error('更新进度失败:', error)
    ElMessage.error('更新失败')
  }
}

// 显示进度报告对话框
const showProgressReportDialog = () => {
  if (!selectedProjectId.value) {
    ElMessage.warning('请先选择一个项目')
    return
  }
  progressReportDialogVisible.value = true
}

// 生成进度报告
const generateProgressReport = async () => {
  try {
    generatingReport.value = true
    
    // 这里应该调用实际的API
    // await studentService.generateProgressReport(reportForm.value)
    
    // 模拟生成过程
    await new Promise(resolve => setTimeout(resolve, 2000))
    
    ElMessage.success('进度报告生成成功')
    progressReportDialogVisible.value = false
  } catch (error) {
    console.error('生成报告失败:', error)
    ElMessage.error('生成失败')
  } finally {
    generatingReport.value = false
  }
}

// 获取进度颜色
const getProgressColor = (progress) => {
  if (progress >= 80) return '#67C23A'
  if (progress >= 60) return '#E6A23C'
  if (progress >= 40) return '#F56C6C'
  return '#909399'
}

// 获取状态类型
const getStatusType = (status) => {
  const typeMap = {
    'on_track': 'success',
    'behind_schedule': 'warning',
    'at_risk': 'danger',
    'completed': 'info'
  }
  return typeMap[status] || 'info'
}

// 获取状态标签
const getStatusLabel = (status) => {
  const labelMap = {
    'on_track': '按计划进行',
    'behind_schedule': '进度滞后',
    'at_risk': '存在风险',
    'completed': '已完成'
  }
  return labelMap[status] || '未知'
}

// 获取状态描述
const getStatusDescription = (status) => {
  const descMap = {
    'on_track': '项目进度符合预期',
    'behind_schedule': '需要加快进度',
    'at_risk': '需要关注风险',
    'completed': '项目已完成'
  }
  return descMap[status] || '状态未知'
}

// 获取进度状态
const getProgressStatus = (progress) => {
  if (progress >= 100) return 'success'
  if (progress >= 80) return 'warning'
  if (progress >= 50) return ''
  return 'exception'
}

// 获取偏差类型
const getDeviationType = (deviation) => {
  if (deviation >= 0) return 'success'
  if (deviation >= -10) return 'warning'
  return 'danger'
}

// 获取模块状态类型
const getModuleStatusType = (status) => {
  const typeMap = {
    'completed': 'success',
    'in_progress': 'primary',
    'not_started': 'info',
    'on_hold': 'warning'
  }
  return typeMap[status] || 'info'
}

// 获取模块状态标签
const getModuleStatusLabel = (status) => {
  const labelMap = {
    'completed': '已完成',
    'in_progress': '进行中',
    'not_started': '未开始',
    'on_hold': '暂停'
  }
  return labelMap[status] || '未知'
}

// 获取风险等级类型
const getRiskLevelType = (level) => {
  const typeMap = {
    'high': 'danger',
    'medium': 'warning',
    'low': 'info'
  }
  return typeMap[level] || 'info'
}

// 获取风险等级标签
const getRiskLevelLabel = (level) => {
  const labelMap = {
    'high': '高',
    'medium': '中',
    'low': '低'
  }
  return labelMap[level] || '未知'
}

// 获取问题优先级类型
const getIssuePriorityType = (priority) => {
  const typeMap = {
    'high': 'danger',
    'medium': 'warning',
    'low': 'info'
  }
  return typeMap[priority] || 'info'
}

// 获取问题优先级标签
const getIssuePriorityLabel = (priority) => {
  const labelMap = {
    'high': '高',
    'medium': '中',
    'low': '低'
  }
  return labelMap[priority] || '未知'
}

// 组件挂载时加载数据
onMounted(() => {
  loadProjects()
})
</script>

<style scoped>
.project-progress {
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

.progress-container {
  margin-top: 20px;
}

.progress-overview {
  margin-bottom: 20px;
}

.overview-card {
  height: 200px;
  transition: all 0.3s ease;
}

.overview-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}

.overview-content {
  text-align: center;
  padding: 20px 0;
}

.overview-content h3 {
  margin: 0 0 20px 0;
  color: #2c3e50;
  font-size: 16px;
}

.progress-circle {
  margin-bottom: 15px;
}

.progress-text {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #409EFF;
}

.days-remaining {
  margin-bottom: 15px;
}

.days-number {
  font-size: 48px;
  font-weight: 700;
  color: #E6A23C;
}

.days-label {
  font-size: 18px;
  color: #909399;
  margin-left: 5px;
}

.deadline-text {
  margin: 0;
  color: #606266;
  font-size: 12px;
}

.status-assessment {
  margin-bottom: 15px;
}

.status-desc {
  margin: 0;
  color: #606266;
  font-size: 12px;
}

.trend-chart-card {
  margin-bottom: 20px;
}

.chart-container {
  height: 400px;
}

.trend-chart {
  width: 100%;
  height: 100%;
}

.module-progress-card {
  margin-bottom: 20px;
}

.milestone-progress-card {
  margin-bottom: 20px;
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
  margin-bottom: 10px;
}

.milestone-header h4 {
  margin: 0;
  color: #2c3e50;
}

.milestone-description {
  margin: 0 0 15px 0;
  color: #606266;
  line-height: 1.6;
}

.milestone-progress {
  display: flex;
  align-items: center;
  gap: 15px;
  margin-bottom: 15px;
}

.milestone-progress .el-progress {
  flex: 1;
}

.progress-text {
  color: #409EFF;
  font-weight: 600;
  min-width: 50px;
}

.milestone-details p {
  margin: 5px 0;
  color: #606266;
  font-size: 12px;
}

.risks-issues-card {
  margin-bottom: 20px;
}

.risks-issues-card h4 {
  margin: 0 0 15px 0;
  color: #2c3e50;
  font-size: 16px;
}

.risk-item, .issue-item {
  margin-bottom: 20px;
  padding: 15px;
  background-color: #f5f7fa;
  border-radius: 8px;
  border-left: 4px solid #409EFF;
}

.risk-item {
  border-left-color: #E6A23C;
}

.issue-item {
  border-left-color: #F56C6C;
}

.risk-title, .issue-title {
  margin-left: 10px;
  font-weight: 600;
  color: #2c3e50;
}

.risk-description, .issue-description {
  margin: 10px 0;
  color: #606266;
  line-height: 1.6;
}

.risk-mitigation, .issue-solution {
  margin: 10px 0 0 0;
  color: #606266;
  font-size: 12px;
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
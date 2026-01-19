<template>
  <div class="student-management">
    <!-- 搜索和筛选 -->
    <el-card class="filter-card">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-input
            v-model="searchQuery"
            placeholder="搜索学生姓名或学号"
            clearable
            @input="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-col>
        <el-col :span="4">
          <el-select v-model="gradeFilter" placeholder="年级筛选" clearable @change="handleFilter">
            <el-option label="全部" value="" />
            <el-option label="大一" value="大一" />
            <el-option label="大二" value="大二" />
            <el-option label="大三" value="大三" />
            <el-option label="大四" value="大四" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-select v-model="majorFilter" placeholder="专业筛选" clearable @change="handleFilter">
            <el-option label="全部" value="" />
            <el-option label="计算机科学" value="计算机科学" />
            <el-option label="软件工程" value="软件工程" />
            <el-option label="信息安全" value="信息安全" />
            <el-option label="人工智能" value="人工智能" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-button type="primary" @click="loadStudents">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- 学生列表 -->
    <el-card class="student-list">
      <template #header>
        <div class="card-header">
          <span>学生列表 ({{ filteredStudents.length }})</span>
        </div>
      </template>
      
      <el-table
        :data="filteredStudents"
        v-loading="loading"
        style="width: 100%"
        @row-click="handleRowClick"
      >
        <el-table-column prop="name" label="姓名" width="100" />
        <el-table-column prop="studentNumber" label="学号" width="120" />
        <el-table-column prop="grade" label="年级" width="80" />
        <el-table-column prop="major" label="专业" width="120" />
        <el-table-column prop="email" label="邮箱" width="180" />
        <el-table-column prop="phone" label="电话" width="120" />
        <el-table-column prop="projectCount" label="项目数" width="80" />
        <el-table-column prop="competitionCount" label="竞赛数" width="80" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click.stop="viewStudent(row)">查看详情</el-button>
            <el-button size="small" type="primary" @click.stop="viewGuidance(row)">指导记录</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 学生详情对话框 -->
    <el-dialog
      v-model="detailVisible"
      title="学生详情"
      width="60%"
      :before-close="handleCloseDetail"
    >
      <div v-if="currentStudent" class="student-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="姓名">{{ currentStudent.name }}</el-descriptions-item>
          <el-descriptions-item label="学号">{{ currentStudent.studentNumber }}</el-descriptions-item>
          <el-descriptions-item label="年级">{{ currentStudent.grade }}</el-descriptions-item>
          <el-descriptions-item label="专业">{{ currentStudent.major }}</el-descriptions-item>
          <el-descriptions-item label="邮箱">{{ currentStudent.email }}</el-descriptions-item>
          <el-descriptions-item label="电话">{{ currentStudent.phone }}</el-descriptions-item>
          <el-descriptions-item label="入学时间">{{ currentStudent.enrollmentDate }}</el-descriptions-item>
          <el-descriptions-item label="指导老师">{{ currentStudent.advisor }}</el-descriptions-item>
        </el-descriptions>

        <!-- 学生项目 -->
        <div class="section">
          <h4>参与项目</h4>
          <el-table :data="currentStudent.projects || []" style="width: 100%">
            <el-table-column prop="title" label="项目标题" />
            <el-table-column prop="type" label="类型" width="100" />
            <el-table-column prop="role" label="角色" width="100" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)">
                  {{ getStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <!-- 学生竞赛 -->
        <div class="section">
          <h4>参与竞赛</h4>
          <el-table :data="currentStudent.competitions || []" style="width: 100%">
            <el-table-column prop="name" label="竞赛名称" />
            <el-table-column prop="level" label="级别" width="100" />
            <el-table-column prop="result" label="成绩" width="100" />
            <el-table-column prop="date" label="参赛时间" width="120" />
          </el-table>
        </div>
      </div>
    </el-dialog>

    <!-- 指导记录对话框 -->
    <el-dialog
      v-model="guidanceVisible"
      title="指导记录"
      width="70%"
      :before-close="handleCloseGuidance"
    >
      <div v-if="currentStudent" class="guidance-record">
        <div class="student-info">
          <h4>{{ currentStudent.name }} - 指导记录</h4>
        </div>

        <!-- 指导记录列表 -->
        <el-table :data="guidanceRecords" style="width: 100%">
          <el-table-column prop="date" label="指导时间" width="160" />
          <el-table-column prop="type" label="指导类型" width="120" />
          <el-table-column prop="content" label="指导内容" />
          <el-table-column prop="duration" label="时长" width="100" />
          <el-table-column label="操作" width="120">
            <template #default="{ row }">
              <el-button size="small" @click="viewGuidanceDetail(row)">查看详情</el-button>
            </template>
          </el-table-column>
        </el-table>

        <!-- 添加指导记录 -->
        <div class="add-guidance">
          <el-button type="primary" @click="showAddGuidance = true">
            <el-icon><Plus /></el-icon>
            添加指导记录
          </el-button>
        </div>
      </div>
    </el-dialog>

    <!-- 添加指导记录对话框 -->
    <el-dialog
      v-model="showAddGuidance"
      title="添加指导记录"
      width="50%"
    >
      <el-form :model="guidanceForm" label-width="100px">
        <el-form-item label="指导类型">
          <el-select v-model="guidanceForm.type" placeholder="请选择指导类型">
            <el-option label="项目指导" value="项目指导" />
            <el-option label="竞赛指导" value="竞赛指导" />
            <el-option label="学术指导" value="学术指导" />
            <el-option label="其他" value="其他" />
          </el-select>
        </el-form-item>
        <el-form-item label="指导内容">
          <el-input
            v-model="guidanceForm.content"
            type="textarea"
            :rows="4"
            placeholder="请输入指导内容"
          />
        </el-form-item>
        <el-form-item label="指导时长">
          <el-input-number
            v-model="guidanceForm.duration"
            :min="0.5"
            :max="8"
            :step="0.5"
            placeholder="小时"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showAddGuidance = false">取消</el-button>
          <el-button type="primary" @click="addGuidanceRecord" :loading="submitting">
            确认
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, Refresh, Plus } from '@element-plus/icons-vue'

// 响应式数据
const loading = ref(false)
const searchQuery = ref('')
const gradeFilter = ref('')
const majorFilter = ref('')
const students = ref([])
const detailVisible = ref(false)
const guidanceVisible = ref(false)
const showAddGuidance = ref(false)
const currentStudent = ref(null)
const submitting = ref(false)

const guidanceForm = ref({
  type: '',
  content: '',
  duration: 1
})

// 模拟学生数据
const mockStudents = [
  {
    id: 1,
    name: '张三',
    studentNumber: '2021001',
    grade: '大三',
    major: '计算机科学',
    email: 'zhangsan@example.com',
    phone: '13800138001',
    enrollmentDate: '2021-09-01',
    advisor: '李教授',
    projectCount: 3,
    competitionCount: 2,
    projects: [
      { title: '智能校园管理系统', type: '科研', role: '负责人', status: 'approved' },
      { title: '在线学习平台', type: '创新', role: '成员', status: 'pending' }
    ],
    competitions: [
      { name: '全国大学生程序设计竞赛', level: '国家级', result: '二等奖', date: '2023-10-15' },
      { name: '蓝桥杯程序设计大赛', level: '省级', result: '一等奖', date: '2023-04-20' }
    ]
  },
  {
    id: 2,
    name: '李四',
    studentNumber: '2021002',
    grade: '大二',
    major: '软件工程',
    email: 'lisi@example.com',
    phone: '13800138002',
    enrollmentDate: '2021-09-01',
    advisor: '王教授',
    projectCount: 1,
    competitionCount: 1,
    projects: [
      { title: '智能校园管理系统', type: '科研', role: '成员', status: 'approved' }
    ],
    competitions: [
      { name: '蓝桥杯程序设计大赛', level: '省级', result: '三等奖', date: '2023-04-20' }
    ]
  },
  {
    id: 3,
    name: '王五',
    studentNumber: '2021003',
    grade: '大四',
    major: '信息安全',
    email: 'wangwu@example.com',
    phone: '13800138003',
    enrollmentDate: '2020-09-01',
    advisor: '张教授',
    projectCount: 5,
    competitionCount: 3,
    projects: [
      { title: '在线学习平台', type: '创新', role: '负责人', status: 'approved' },
      { title: '网络安全防护系统', type: '科研', role: '负责人', status: 'approved' }
    ],
    competitions: [
      { name: '全国大学生信息安全竞赛', level: '国家级', result: '一等奖', date: '2023-11-10' },
      { name: '全国大学生程序设计竞赛', level: '国家级', result: '三等奖', date: '2023-10-15' }
    ]
  }
]

// 模拟指导记录数据
const guidanceRecords = ref([
  {
    id: 1,
    date: '2024-01-15 14:30:00',
    type: '项目指导',
    content: '讨论智能校园管理系统的技术架构和实现方案',
    duration: 2
  },
  {
    id: 2,
    date: '2024-01-10 10:00:00',
    type: '竞赛指导',
    content: '分析程序设计竞赛的解题思路和算法优化',
    duration: 1.5
  },
  {
    id: 3,
    date: '2024-01-05 16:00:00',
    type: '学术指导',
    content: '指导论文写作和学术研究方法',
    duration: 1
  }
])

// 计算属性
const filteredStudents = computed(() => {
  let result = students.value

  // 年级筛选
  if (gradeFilter.value) {
    result = result.filter(s => s.grade === gradeFilter.value)
  }

  // 专业筛选
  if (majorFilter.value) {
    result = result.filter(s => s.major === majorFilter.value)
  }

  // 搜索筛选
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(s => 
      s.name.toLowerCase().includes(query) ||
      s.studentNumber.toLowerCase().includes(query)
    )
  }

  return result
})

// 方法
const loadStudents = async () => {
  loading.value = true
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    students.value = [...mockStudents]
    ElMessage.success('学生列表加载成功')
  } catch (error) {
    ElMessage.error('加载学生列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  // 搜索逻辑已在计算属性中处理
}

const handleFilter = () => {
  // 筛选逻辑已在计算属性中处理
}

const handleRowClick = (row) => {
  viewStudent(row)
}

const viewStudent = (student) => {
  currentStudent.value = student
  detailVisible.value = true
}

const handleCloseDetail = () => {
  detailVisible.value = false
  currentStudent.value = null
}

const viewGuidance = (student) => {
  currentStudent.value = student
  guidanceVisible.value = true
}

const handleCloseGuidance = () => {
  guidanceVisible.value = false
  currentStudent.value = null
}

const viewGuidanceDetail = (record) => {
  ElMessage.info(`查看指导记录详情: ${record.content}`)
}

const addGuidanceRecord = async () => {
  if (!guidanceForm.value.type || !guidanceForm.value.content) {
    ElMessage.warning('请填写完整的指导记录信息')
    return
  }

  submitting.value = true
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // 添加新记录
    const newRecord = {
      id: guidanceRecords.value.length + 1,
      date: new Date().toLocaleString(),
      type: guidanceForm.value.type,
      content: guidanceForm.value.content,
      duration: guidanceForm.value.duration
    }
    
    guidanceRecords.value.unshift(newRecord)
    
    ElMessage.success('指导记录添加成功')
    showAddGuidance.value = false
    
    // 重置表单
    guidanceForm.value = {
      type: '',
      content: '',
      duration: 1
    }
  } catch (error) {
    ElMessage.error('添加指导记录失败')
  } finally {
    submitting.value = false
  }
}

const getStatusType = (status) => {
  const statusMap = {
    draft: 'info',
    pending: 'warning',
    approved: 'success',
    rejected: 'danger'
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status) => {
  const statusMap = {
    draft: '草稿',
    pending: '待审核',
    approved: '已通过',
    rejected: '已拒绝'
  }
  return statusMap[status] || status
}

// 组件挂载时加载数据
onMounted(() => {
  loadStudents()
})
</script>

<style scoped>
.student-management {
  padding: 20px;
}

.filter-card {
  margin-bottom: 20px;
}

.student-list {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.student-detail {
  max-height: 60vh;
  overflow-y: auto;
}

.guidance-record {
  max-height: 60vh;
  overflow-y: auto;
}

.student-info {
  margin-bottom: 20px;
}

.student-info h4 {
  margin: 0;
  color: #2c3e50;
  font-size: 18px;
  font-weight: 600;
}

.section {
  margin-top: 20px;
}

.section h4 {
  margin-bottom: 10px;
  color: #2c3e50;
  font-size: 16px;
  font-weight: 600;
}

.add-guidance {
  margin-top: 20px;
  text-align: center;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style> 
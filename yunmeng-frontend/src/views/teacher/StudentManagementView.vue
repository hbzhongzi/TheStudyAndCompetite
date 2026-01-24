<template>
  <div class="student-management">
    <!-- 搜索和筛选 -->
    <el-card class="filter-card">
      <el-row :gutter="20" align="middle">
        <el-col :xs="24" :sm="12" :md="8" :lg="6">
          <el-input
            v-model="searchQuery"
            placeholder="搜索学生姓名或用户名"
            clearable
            @input="handleSearch"
            @clear="handleSearch"
            @keyup.enter="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-col>
        <el-col :xs="24" :sm="12" :md="8" :lg="4">
          <el-select
            v-model="gradeFilter"
            placeholder="年级筛选"
            clearable
            filterable
            @change="handleFilter"
          >
            <el-option label="全部" value="" />
            <el-option 
              v-for="grade in gradeOptions"
              :key="grade"
              :label="grade"
              :value="grade"
            />
          </el-select>
        </el-col>
        <el-col :xs="24" :sm="12" :md="8" :lg="4">
          <el-select
            v-model="majorFilter"
            placeholder="专业筛选"
            clearable
            filterable
            @change="handleFilter"
          >
            <el-option label="全部" value="" />
            <el-option 
              v-for="major in majorOptions"
              :key="major"
              :label="major"
              :value="major"
            />
          </el-select>
        </el-col>
        <el-col :xs="24" :sm="12" :md="8" :lg="4">
          <el-button-group>
            <el-button type="primary" @click="loadStudents" :loading="loading">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
            <el-button type="success" @click="exportData">
              <el-icon><Download /></el-icon>
              导出
            </el-button>
          </el-button-group>
        </el-col>
      </el-row>
    </el-card>

    <!-- 学生列表 -->
    <el-card class="student-list">
      <template #header>
        <div class="card-header">
          <span>学生列表 ({{ pagination.total || 0 }})</span>
        </div>
      </template>
      
      <el-table
        :data="filteredStudents"
        v-loading="loading"
        style="width: 100%"
        @selection-change="handleSelectionChange"
        @row-click="handleRowClick"
        stripe
        border
      >
        <el-table-column type="selection" width="55" />
        <el-table-column 
          prop="realName" 
          label="姓名" 
          width="100"
          sortable
        >
          <template #default="{ row }">
            <el-tag v-if="row.status === 'inactive'" type="info" size="small">已禁用</el-tag>
            {{ row.realName || row.name }}
          </template>
        </el-table-column>
        <el-table-column prop="name" label="用户名" width="120" sortable />
        <el-table-column prop="grade" label="年级" width="80" sortable>
          <template #default="{ row }">
            <span v-if="row.grade">{{ row.grade }}</span>
            <span v-else class="text-gray">未设置</span>
          </template>
        </el-table-column>
        <el-table-column prop="major" label="专业" width="120" sortable>
          <template #default="{ row }">
            <span v-if="row.major">{{ row.major }}</span>
            <span v-else class="text-gray">未设置</span>
          </template>
        </el-table-column>
        <el-table-column prop="email" label="邮箱" width="180">
          <template #default="{ row }">
            <el-link :href="`mailto:${row.email}`" type="primary">
              {{ row.email }}
            </el-link>
          </template>
        </el-table-column>
        <el-table-column prop="phone" label="电话" width="120">
          <template #default="{ row }">
            <span v-if="row.phone">{{ row.phone }}</span>
            <span v-else class="text-gray">未设置</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag 
              :type="row.status === 'active' ? 'success' : 'info'"
              size="small"
            >
              {{ row.status === 'active' ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="注册时间" width="150" sortable>
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button-group>
              <el-button 
                size="small" 
                @click.stop="viewStudent(row)"
              >
                详情
              </el-button>
              <el-button 
                size="small" 
                type="primary" 
                @click.stop="viewGuidance(row)"
              >
                指导
              </el-button>
              <el-dropdown @command="(command) => handleMoreAction(command, row)">
                <el-button size="small">
                  更多
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="status">
                      {{ row.status === 'active' ? '禁用账户' : '启用账户' }}
                    </el-dropdown-item>
                    <el-dropdown-item command="reset">重置密码</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页组件 -->
      <div class="pagination-container" v-if="pagination.total > 0">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.size"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 学生详情对话框 -->
    <el-dialog
      v-model="detailVisible"
      title="学生详情"
      width="800px"
    >
      <div class="student-detail" v-if="currentStudent">
        <div class="student-info">
          <h4>{{ currentStudent.realName || currentStudent.name }}</h4>
          <p>用户名: {{ currentStudent.name }}</p>
          <p>邮箱: {{ currentStudent.email }}</p>
          <p>电话: {{ currentStudent.phone || '未设置' }}</p>
          <p>年级: {{ currentStudent.grade || '未设置' }}</p>
          <p>专业: {{ currentStudent.major || '未设置' }}</p>
        </div>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="handleCloseDetail">关闭</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 指导记录对话框 -->
    <el-dialog
      v-model="guidanceVisible"
      :title="`${currentStudent?.realName || currentStudent?.name} 的指导记录`"
      width="900px"
    >
      <div class="guidance-record" v-if="currentStudent">
        <el-table :data="guidanceRecords" style="width: 100%">
          <el-table-column prop="date" label="时间" width="150" />
          <el-table-column prop="type" label="类型" width="100" />
          <el-table-column prop="content" label="内容" />
          <el-table-column prop="duration" label="时长(小时)" width="100" />
          
        </el-table>

        <div class="add-guidance">
          <el-button 
            type="primary" 
            @click="showAddGuidance = true"
            v-if="!showAddGuidance"
          >
            添加指导记录
          </el-button>
          
          <el-form 
            v-else
            :model="guidanceForm" 
            label-width="80px" 
            style="margin-top: 20px;"
          >
            <el-form-item label="指导类型">
              <el-select v-model="guidanceForm.type" placeholder="请选择指导类型">
                <el-option label="项目指导" value="项目指导" />
                <el-option label="竞赛指导" value="竞赛指导" />
                <el-option label="学术指导" value="学术指导" />
                <el-option label="就业指导" value="就业指导" />
              </el-select>
            </el-form-item>
            <el-form-item label="指导内容">
              <el-input
                v-model="guidanceForm.content"
                type="textarea"
                :rows="3"
                placeholder="请输入指导内容"
              />
            </el-form-item>
            <el-form-item label="时长(小时)">
              <el-input-number 
                v-model="guidanceForm.duration" 
                :min="0.5" 
                :step="0.5" 
                :precision="1"
              />
            </el-form-item>
            <el-form-item>
              <el-button @click="showAddGuidance = false">取消</el-button>
              <el-button 
                type="primary" 
                @click="addGuidanceRecord"
                :loading="submitting"
              >
                保存
              </el-button>
            </el-form-item>
          </el-form>
        </div>
      </div>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="handleCloseGuidance">关闭</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Search, Refresh, Download, Plus } from '@element-plus/icons-vue'
import { teacherService } from '../../services/teacherService'
import { ElMessage } from 'element-plus'
import { debounce } from 'lodash-es'

// 响应式数据
const searchQuery = ref('')
const gradeFilter = ref('')
const majorFilter = ref('')
const loading = ref(false)
const students = ref([])
const selectedStudents = ref([])
const detailVisible = ref(false)
const guidanceVisible = ref(false)
const showAddGuidance = ref(false)
const currentStudent = ref(null)
const submitting = ref(false)

// 分页信息
const pagination = reactive({
  page: 1,
  size: 20,
  total: 0
})

// 选项数据
const gradeOptions = ref(['大一', '大二', '大三', '大四', '研一', '研二', '研三'])
const majorOptions = ref(['计算机科学', '软件工程', '信息安全', '人工智能', '数据科学', '网络工程'])

// 指导表单
const guidanceForm = ref({
  type: '',
  content: '',
  duration: 1
})

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

// 计算属性 - 本地筛选（当需要前端筛选时使用）
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
      (s.realName && s.realName.toLowerCase().includes(query)) ||
      (s.name && s.name.toLowerCase().includes(query)) ||
      (s.email && s.email.toLowerCase().includes(query))
    )
  }

  return result
})

// 从数据中提取筛选选项
const extractOptionsFromData = () => {
  const grades = new Set()
  const majors = new Set()
  
  students.value.forEach(student => {
    if (student.grade) grades.add(student.grade)
    if (student.major) majors.add(student.major)
  })
  
  // 合并现有选项和新选项
  gradeOptions.value = Array.from(new Set([...gradeOptions.value, ...grades]))
  majorOptions.value = Array.from(new Set([...majorOptions.value, ...majors]))
}

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return ''
  try {
    const date = new Date(dateString)
    return date.toLocaleDateString('zh-CN') + ' ' + date.toLocaleTimeString('zh-CN', {
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch (error) {
    return dateString
  }
}

// 防抖搜索函数
const debouncedSearch = debounce(() => {
  pagination.page = 1
  loadStudents()
}, 500)

// 搜索处理
const handleSearch = () => {
  debouncedSearch()
}

// 筛选处理
const handleFilter = () => {
  pagination.page = 1
  loadStudents()
}

// 分页处理
const handleSizeChange = (size) => {
  pagination.size = size
  pagination.page = 1
  loadStudents()
}

const handleCurrentChange = (page) => {
  pagination.page = page
  loadStudents()
}

// 行选择
const handleSelectionChange = (selection) => {
  selectedStudents.value = selection
}

// 加载学生数据
const loadStudents = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      size: pagination.size,
      search: searchQuery.value,
      grade: gradeFilter.value,
      major: majorFilter.value
    }
    
    const response = await teacherService.getStudents(params)
    if (response.code === 200) {
      students.value = response.data.list || []
      pagination.total = response.data.total || 0
      pagination.page = response.data.page || 1
      pagination.size = response.data.size || 20
      
      // 提取年级和专业选项
      extractOptionsFromData()
      
      ElMessage.success(`加载成功，共 ${pagination.total} 名学生`)
    } else {
      ElMessage.error(response.message || '加载失败')
    }
  } catch (error) {
    console.error('加载学生列表失败:', error)
    ElMessage.error('加载学生列表失败')
  } finally {
    loading.value = false
  }
}

// 操作函数
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
  showAddGuidance.value = false
}

const handleCloseGuidance = () => {
  guidanceVisible.value = false
  currentStudent.value = null
  showAddGuidance.value = false
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

const handleMoreAction = (command, row) => {
  switch (command) {
    case 'edit':
      editStudent(row)
      break
    case 'status':
      toggleStatus(row)
      break
    case 'reset':
      resetPassword(row)
      break
    case 'delete':
      deleteStudent(row)
      break
  }
}

const editStudent = (row) => {
  console.log('编辑学生:', row)
  // 实现编辑功能
}

const toggleStatus = async (row) => {
  try {
    const newStatus = row.status === 'active' ? 'inactive' : 'active'
    // 这里调用您的接口
    // await teacherService.updateStudentStatus(row.id, newStatus)
    row.status = newStatus
    ElMessage.success('状态更新成功')
  } catch (error) {
    ElMessage.error('状态更新失败')
  }
}

const resetPassword = async (row) => {
  try {
    // 这里调用您的接口
    // await teacherService.resetStudentPassword(row.id)
    ElMessage.success('密码重置成功')
  } catch (error) {
    ElMessage.error('密码重置失败')
  }
}

const deleteStudent = async (row) => {
  try {
    // 这里调用您的接口
    // await teacherService.deleteStudent(row.id)
    loadStudents()
    ElMessage.success('删除成功')
  } catch (error) {
    ElMessage.error('删除失败')
  }
}

const exportData = () => {
  console.log('导出数据')
  // 实现导出功能
}

const handleBatchAction = () => {
  console.log('批量操作:', selectedStudents.value)
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

.header-actions {
  display: flex;
  gap: 10px;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
  padding-top: 20px;
  border-top: 1px solid #f0f0f0;
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

.add-guidance {
  margin-top: 20px;
  text-align: center;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.text-gray {
  color: #999;
}

/* 响应式调整 */
@media (max-width: 768px) {
  .filter-card .el-col {
    margin-bottom: 10px;
  }
  
  .card-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 10px;
  }
  
  .header-actions {
    width: 100%;
    justify-content: flex-end;
  }
}
</style>
<template>
  <div class="project-type-management">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-left">
        <h2>项目分类管理</h2>
        <p class="header-desc">管理系统内的项目分类，支持创建、编辑、删除等操作</p>
      </div>
      <div class="header-right">
        <el-button type="primary" @click="showCreateDialog">
          <i class="el-icon-plus"></i>
          新建分类
        </el-button>
        <el-button type="success" @click="showStats">
          <i class="el-icon-s-data"></i>
          分类统计
        </el-button>
      </div>
    </div>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row" v-if="typeStats.length > 0">
      <el-col :span="8" v-for="stat in typeStats.slice(0, 3)" :key="stat.id">
        <el-card class="stat-card">
          <div class="stat-content">
            <div class="stat-icon">
              <i class="el-icon-folder"></i>
            </div>
            <div class="stat-info">
              <h4>{{ stat.name }}</h4>
              <p class="stat-number">{{ stat.projectCount }}</p>
              <p class="stat-desc">个项目</p>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 分类列表 -->
    <el-card class="type-list-card">
      <template #header>
        <div class="list-header">
          <span>分类列表</span>
          <div class="list-actions">
            <el-button size="small" @click="refreshList">
              <i class="el-icon-refresh"></i>
              刷新
            </el-button>
          </div>
        </div>
      </template>

      <el-table 
        :data="projectTypes" 
        style="width: 100%" 
        v-loading="loading"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="分类名称" min-width="150" />
        <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="projectCount" label="项目数量" width="100">
          <template #default="{ row }">
            <el-tag :type="row.projectCount > 0 ? 'success' : 'info'">
              {{ row.projectCount }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column prop="updatedAt" label="更新时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.updatedAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="viewTypeDetail(row)">查看详情</el-button>
            <el-button size="small" type="primary" @click="editType(row)">编辑</el-button>
            <el-button 
              size="small" 
              type="danger" 
              @click="deleteType(row)"
              :disabled="row.projectCount > 0"
            >
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="totalTypes"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 创建/编辑分类对话框 -->
    <el-dialog
      v-model="showTypeDialog"
      :title="isEdit ? '编辑分类' : '新建分类'"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form :model="typeForm" :rules="typeRules" ref="typeFormRef" label-width="100px">
        <el-form-item label="分类名称" prop="name">
          <el-input
            v-model="typeForm.name"
            placeholder="请输入分类名称"
            maxlength="50"
            show-word-limit
          />
        </el-form-item>
        <el-form-item label="分类描述" prop="description">
          <el-input
            v-model="typeForm.description"
            type="textarea"
            :rows="4"
            placeholder="请输入分类描述"
            maxlength="200"
            show-word-limit
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showTypeDialog = false">取消</el-button>
        <el-button type="primary" @click="submitType" :loading="submitting">
          {{ isEdit ? '更新' : '创建' }}
        </el-button>
      </template>
    </el-dialog>

    <!-- 分类详情对话框 -->
    <el-dialog
      v-model="showDetailDialog"
      title="分类详情"
      width="600px"
      :close-on-click-modal="false"
    >
      <div v-if="selectedType" class="type-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="分类ID">{{ selectedType.id }}</el-descriptions-item>
          <el-descriptions-item label="分类名称">{{ selectedType.name }}</el-descriptions-item>
          <el-descriptions-item label="项目数量">{{ selectedType.projectCount }}</el-descriptions-item>
          <el-descriptions-item label="创建时间">{{ formatDate(selectedType.createdAt) }}</el-descriptions-item>
          <el-descriptions-item label="更新时间">{{ formatDate(selectedType.updatedAt) }}</el-descriptions-item>
          <el-descriptions-item label="分类描述" :span="2">{{ selectedType.description || '暂无描述' }}</el-descriptions-item>
        </el-descriptions>
      </div>
    </el-dialog>

    <!-- 分类统计对话框 -->
    <el-dialog
      v-model="showStatsDialog"
      title="分类统计"
      width="70%"
      :close-on-click-modal="false"
    >
      <div class="stats-content">
        <el-row :gutter="20">
          <el-col :span="12">
            <h4>分类项目分布</h4>
            <div class="chart-container">
              <!-- 这里可以添加图表组件 -->
              <el-table :data="typeStats" style="width: 100%">
                <el-table-column prop="name" label="分类名称" />
                <el-table-column prop="projectCount" label="项目数量" />
                <el-table-column label="占比">
                  <template #default="{ row }">
                    {{ ((row.projectCount / totalProjects) * 100).toFixed(1) }}%
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </el-col>
          <el-col :span="12">
            <h4>统计信息</h4>
            <el-card class="stats-summary">
              <p><strong>总分类数：</strong>{{ typeStats.length }}</p>
              <p><strong>总项目数：</strong>{{ totalProjects }}</p>
              <p><strong>有项目的分类：</strong>{{ activeTypes }}</p>
              <p><strong>空分类：</strong>{{ emptyTypes }}</p>
            </el-card>
          </el-col>
        </el-row>
      </div>
    </el-dialog>

    <!-- 删除确认对话框 -->
    <el-dialog
      v-model="showDeleteDialog"
      title="确认删除"
      width="400px"
    >
      <p>确定要删除分类"{{ selectedType?.name }}"吗？此操作不可恢复。</p>
      <template #footer>
        <el-button @click="showDeleteDialog = false">取消</el-button>
        <el-button type="danger" @click="confirmDelete" :loading="deleting">
          确认删除
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { projectTypeService } from '../../services/projectTypeService'

// 响应式数据
const loading = ref(false)
const submitting = ref(false)
const deleting = ref(false)
const currentPage = ref(1)
const pageSize = ref(20)
const totalTypes = ref(0)
const showTypeDialog = ref(false)
const showDetailDialog = ref(false)
const showStatsDialog = ref(false)
const showDeleteDialog = ref(false)
const isEdit = ref(false)
const selectedType = ref(null)
const typeFormRef = ref(null)

// 分类列表
const projectTypes = ref([])
const typeStats = ref([])

// 分类表单
const typeForm = reactive({
  name: '',
  description: ''
})

// 表单验证规则
const typeRules = {
  name: [
    { required: true, message: '请输入分类名称', trigger: 'blur' },
    { min: 2, max: 50, message: '分类名称长度在 2 到 50 个字符', trigger: 'blur' }
  ],
  description: [
    { max: 200, message: '描述长度不能超过 200 个字符', trigger: 'blur' }
  ]
}

// 计算属性
const totalProjects = computed(() => {
  return typeStats.value.reduce((sum, stat) => sum + stat.projectCount, 0)
})

const activeTypes = computed(() => {
  return typeStats.value.filter(stat => stat.projectCount > 0).length
})

const emptyTypes = computed(() => {
  return typeStats.value.filter(stat => stat.projectCount === 0).length
})

// 方法
const loadProjectTypes = async () => {
  loading.value = true
  try {
    const response = await projectTypeService.getProjectTypeList()
    if (response && response.code === 200) {
      projectTypes.value = response.data.list || []
      totalTypes.value = response.data.total || 0
    }
  } catch (error) {
    console.error('加载项目分类列表失败:', error)
    ElMessage.error('加载项目分类列表失败')
  } finally {
    loading.value = false
  }
}

const loadTypeStats = async () => {
  try {
    const response = await projectTypeService.getProjectTypeStats()
    if (response && response.code === 200) {
      typeStats.value = response.data || []
    }
  } catch (error) {
    console.error('加载分类统计失败:', error)
  }
}

const refreshList = () => {
  loadProjectTypes()
  loadTypeStats()
}

const handleSizeChange = (size) => {
  pageSize.value = size
  loadProjectTypes()
}

const handleCurrentChange = (page) => {
  currentPage.value = page
  loadProjectTypes()
}

const showCreateDialog = () => {
  isEdit.value = false
  typeForm.name = ''
  typeForm.description = ''
  showTypeDialog.value = true
}

const editType = (type) => {
  isEdit.value = true
  selectedType.value = type
  typeForm.name = type.name
  typeForm.description = type.description || ''
  showTypeDialog.value = true
}

const submitType = async () => {
  if (!typeFormRef.value) return
  
  try {
    await typeFormRef.value.validate()
  } catch (error) {
    return
  }

  submitting.value = true
  try {
    if (isEdit.value) {
      await projectTypeService.updateProjectType(selectedType.value.id, typeForm)
      ElMessage.success('分类更新成功')
    } else {
      await projectTypeService.createProjectType(typeForm)
      ElMessage.success('分类创建成功')
    }
    
    showTypeDialog.value = false
    refreshList()
  } catch (error) {
    ElMessage.error(error.message)
  } finally {
    submitting.value = false
  }
}

const viewTypeDetail = async (type) => {
  try {
    const response = await projectTypeService.getProjectTypeDetail(type.id)
    if (response && response.code === 200) {
      selectedType.value = response.data
      showDetailDialog.value = true
    }
  } catch (error) {
    ElMessage.error('获取分类详情失败')
  }
}

const deleteType = (type) => {
  if (type.projectCount > 0) {
    ElMessage.warning('该分类下还有项目，无法删除')
    return
  }
  
  selectedType.value = type
  showDeleteDialog.value = true
}

const confirmDelete = async () => {
  deleting.value = true
  try {
    await projectTypeService.deleteProjectType(selectedType.value.id)
    ElMessage.success('分类删除成功')
    showDeleteDialog.value = false
    refreshList()
  } catch (error) {
    ElMessage.error(error.message)
  } finally {
    deleting.value = false
  }
}

const showStats = () => {
  showStatsDialog.value = true
}

const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

// 组件挂载时加载数据
onMounted(() => {
  loadProjectTypes()
  loadTypeStats()
})
</script>

<style scoped>
.project-type-management {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
}

.header-left h2 {
  margin: 0 0 8px 0;
  color: #2c3e50;
  font-size: 24px;
  font-weight: 600;
}

.header-desc {
  margin: 0;
  color: #7f8c8d;
  font-size: 14px;
}

.header-right {
  display: flex;
  gap: 10px;
}

.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
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
  background: linear-gradient(135deg, #3498db, #2980b9);
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
  margin: 0;
  font-size: 28px;
  font-weight: 600;
  color: #2c3e50;
}

.stat-desc {
  margin: 5px 0 0 0;
  color: #7f8c8d;
  font-size: 12px;
}

.type-list-card {
  border-radius: 8px;
}

.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.list-actions {
  display: flex;
  gap: 10px;
}

.pagination-wrapper {
  margin-top: 20px;
  text-align: right;
}

.el-table {
  border-radius: 8px;
  overflow: hidden;
}

.el-table :deep(.el-table__row) {
  cursor: pointer;
}

.el-table :deep(.el-table__row:hover) {
  background-color: #f5f7fa;
}

.type-detail {
  padding: 10px 0;
}

.stats-content {
  padding: 10px 0;
}

.stats-content h4 {
  margin: 0 0 15px 0;
  color: #2c3e50;
  font-size: 16px;
  font-weight: 600;
}

.chart-container {
  margin-bottom: 20px;
}

.stats-summary {
  padding: 15px;
}

.stats-summary p {
  margin: 8px 0;
  color: #2c3e50;
}
</style> 
<template>
  <div class="project-overview">

    <!-- 搜索与筛选 -->
    <el-card class="filter-card">
      <el-row :gutter="20">
        <el-col :span="4">
          <el-select v-model="query.status" placeholder="项目状态" clearable @change="reload">
            <el-option label="进行中" value="ongoing" />
            <el-option label="已完成" value="completed" />
            <el-option label="待审核" value="pending" />
            <el-option label="已拒绝" value="rejected" />
          </el-select>
        </el-col>

        <el-col :span="4">
          <el-select v-model="query.type" placeholder="项目类型" clearable @change="reload">
            <el-option label="科研项目" value="科研项目" />
            <el-option label="创新项目" value="创新项目" />
            <el-option label="竞赛项目" value="竞赛项目" />
          </el-select>
        </el-col>

      <el-col :span="6">
          <el-button type="primary" @click="openCreateDialog">创建项目</el-button>
          <el-button @click="reload">刷新</el-button>
      </el-col>
      </el-row>
    </el-card>

    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stat-row">
      <el-col :span="6" v-for="item in stats" :key="item.title">
        <el-card class="stat-card">
          <div class="stat-value">{{ item.value }}</div>
          <div class="stat-title">{{ item.title }}</div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 项目列表 -->
    <el-card>
      <template #header>
        <span>我的项目</span>
      </template>

      <el-table :data="projects" v-loading="loading">
        <el-table-column label="项目名称" min-width="200">
          <template #default="{ row }">
            <el-link type="primary" @click="openDetail(row)">
              {{ row.name }}
            </el-link>
          </template>
        </el-table-column>

        <el-table-column prop="type" label="类型" width="120" />

        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusTag(row.status)">
              {{ statusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="进度" width="160">
          <template #default="{ row }">
            <el-progress :percentage="row.progress" />
          </template>
        </el-table-column>
          <el-table-column prop="createdAt" label="创建时间" width="140">
              <template #default="scope">
                    <span>{{ formatDateTime(scope.row.createdAt) }}</span>
              </template>
          </el-table-column>
<el-table-column label="操作" width="320" fixed="right">
  <template #default="{ row }">
    <el-button size="small" @click="openDetail(row)">查看</el-button>

    <!-- 提交审核：仅草稿可见 -->
    <el-button
      v-if="row.status === 'draft'"
      size="small"
      type="warning"
      @click="submitForReview(row)"
    >
      提交审核
    </el-button>

    <!-- 更新进度：仅审核通过 -->
    <el-button
      v-if="row.status === 'approved'"
      size="small"
      type="success"
      @click="openProgress(row)"
    >
      更新进度
    </el-button>

    <el-button size="small" type="danger" @click="remove(row)">
      删除
    </el-button>
  </template>
</el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="page.page"
          v-model:page-size="page.size"
          :total="page.total"
          layout="total, prev, pager, next"
          @current-change="reload"
        />
      </div>
    </el-card>

    <!-- 项目详情 -->
    <el-dialog v-model="detailVisible" title="项目详情" width="70%">
      <el-descriptions border :column="2">
        <el-descriptions-item label="项目名称">{{ current?.name }}</el-descriptions-item>
        <el-descriptions-item label="类型">{{ current?.type }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ statusText(current?.status) }}</el-descriptions-item>
        <el-descriptions-item label="进度">{{ current?.progress }}%</el-descriptions-item>
      </el-descriptions>

      <el-divider />
      <p><strong>项目描述：</strong>{{ current?.description }}</p >
      <p><strong>项目计划：</strong>{{ current?.plan }}</p >
    </el-dialog>

    <!-- 进度更新 -->
    <el-dialog v-model="progressVisible" title="更新进度" width="40%">
      <el-form>
        <el-form-item label="进度">
          <el-slider v-model="progressForm.progress" />
        </el-form-item>
        <el-form-item label="说明">
          <el-input v-model="progressForm.description" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="progressVisible=false">取消</el-button>
        <el-button type="primary" @click="submitProgress">提交</el-button>
      </template>
    </el-dialog>

  </div>

<el-dialog
  v-model="createVisible"
  title="创建项目"
  width="600px"
>
  <el-form
    ref="createFormRef"
    :model="createForm"
    :rules="createRules"
    label-width="100px"
  >
    <!-- 项目标题 -->
    <el-form-item label="项目标题" prop="title">
      <el-input
        v-model="createForm.title"
        placeholder="请输入项目标题（不超过100字）"
        maxlength="100"
        show-word-limit
      />
    </el-form-item>

    <!-- 项目类型 -->
    <el-form-item label="项目类型" prop="type">
      <el-radio-group v-model="createForm.type">
        <el-radio label="科研">科研项目</el-radio>
        <el-radio label="竞赛">竞赛项目</el-radio>
      </el-radio-group>
    </el-form-item>

    <!-- 指导教师 -->
    <el-form-item label="指导教师" prop="teacherId">
      <el-select
        v-model="createForm.teacherId"
        placeholder="请选择指导教师"
        filterable
      >
        <el-option
            v-for="t in teacherList"
            :key="t.id"
            :label="`${t.realName}（${t.department}）`"
            :value="t.id"
        />
      </el-select>
    </el-form-item>

    <!-- 项目描述 -->
    <el-form-item label="项目描述">
      <el-input
        v-model="createForm.description"
        type="textarea"
        rows="3"
        placeholder="项目背景、研究内容说明"
      />
    </el-form-item>

    <!-- 预期成果 -->
    <el-form-item label="预期成果">
      <el-input
        v-model="createForm.plan"
        type="textarea"
        rows="3"
        placeholder="论文、系统、竞赛成果等"
      />
    </el-form-item>

    <!-- 预计完成时间 -->
    <el-form-item label="预计完成时间">
      <el-date-picker
        v-model="createForm.FinishedAt"
        type="date"
        placeholder="选择日期"
        value-format="YYYY-MM-DDTHH:mm:ssZ"
        format="YYYY-MM-DD"
      />
    </el-form-item>
  </el-form>

  <template #footer>
    <el-button @click="createVisible = false">取消</el-button>
    <el-button type="primary" @click="submitCreate">提交</el-button>
  </template>
</el-dialog>

</template>

<script setup>
/* ==================== 基础依赖 ==================== */
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

import studentService from '@/services/studentService'
import { projectService } from '@/services/projectService'
import teacherService  from '@/services/teacherService'

/* ==================== 创建项目相关 ==================== */
const createVisible = ref(false)
const createFormRef = ref(null)

const createForm = ref({
  title: '',
  description: '',
  type: '科研',
  status: 'draft',
  teacherId: '',
  plan: '',
  FinishedAt: null
})

const createRules = {
  title: [
    { required: true, message: '请输入项目标题', trigger: 'blur' },
    { max: 100, message: '标题不超过100字', trigger: 'blur' }
  ],
  teacherId: [
    { required: true, message: '请选择指导教师', trigger: 'change' }
  ],
  type: [
    { required: true, message: '请选择项目类型', trigger: 'change' }
  ]
}

const teacherList = ref([])

const openCreateDialog = () => {
  createVisible.value = true
}

const submitCreate = () => {
  createFormRef.value.validate(async valid => {
    if (!valid) return

    await projectService.createProject(createForm.value)
    ElMessage.success('项目创建成功')

    createVisible.value = false
    reloadProjects()
    loadStats()
  })
}

const submitForReview = row => {
  ElMessageBox.confirm(
    '提交后将进入审核流程，确认提交？',
    '提示',
    { type: 'warning' }
  ).then(async () => {
    await projectService.submitProject(row.id)
    ElMessage.success('项目已提交审核')
    reloadProjects()
    loadStats()
  })
}

const openProgress = row => {
  if (row.status !== 'approved') {
    ElMessage.warning('项目审核通过后才能更新进度')
    return
  }

  current.value = row
  progressForm.value.progress = row.progress
  progressForm.value.description = ''
  progressVisible.value = true
}

const loadTeachers = async () => {
  const res = await teacherService.getTeacherList()
  teacherList.value = res.data || []
}

/* ==================== 项目列表 & 分页 ==================== */
const loading = ref(false)
const projects = ref([])
const current = ref(null)

const query = ref({
  keyword: '',
  status: '',
  type: ''
})

const page = ref({
  page: 1,
  size: 10,
  total: 0
})

const reloadProjects = async () => {
  loading.value = true

  const res = await studentService.getMyProjects({
    ...query.value,
    page: page.value.page,
    size: page.value.size
  })

  projects.value = (res.data.list || []).map(p => ({
    ...p,
    name: p.title
  }))

  page.value.total = res.data.total
  loading.value = false
}

/* ==================== 统计数据 ==================== */
const stats = ref([])

const loadStats = async () => {
  const { data } = await studentService.getProjectStats()

  stats.value = [
    { title: '项目总数', value: data.totalProjects },
    { title: '进行中', value: data.ongoingProjects },
    { title: '已完成', value: data.completedProjects },
    { title: '待审核', value: data.pendingProjects }
  ]
}

/* ==================== 项目详情 & 进度 ==================== */
const detailVisible = ref(false)
const progressVisible = ref(false)

const progressForm = ref({
  progress: 0,
  description: ''
})

const openDetail = row => {
  current.value = row
  detailVisible.value = true
}


const submitProgress = async () => {
  await studentService.updateProjectProgress(
    current.value.id,
    progressForm.value
  )

  ElMessage.success('进度更新成功')
  progressVisible.value = false
  reloadProjects()
}

/* ==================== 删除项目 ==================== */
const remove = row => {
  ElMessageBox.confirm('确认删除该项目？', '提示', { type: 'warning' })
    .then(async () => {
      await studentService.deleteProject(row.id)
      ElMessage.success('删除成功')
      reloadProjects()
      loadStats()
    })
}

/* ==================== 工具函数 ==================== */
const formatDateTime = dateString => {
  if (!dateString) return '--'
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', { hour12: false })
}

const statusText = s =>
  ({
    draft: '草稿',
    pending: '待审核',
    approved: '审核通过',
    ongoing: '进行中',
    completed: '已完成',
    rejected: '已拒绝'
  }[s] || s)

const statusTag = s =>
  ({ draft: '', ongoing: 'success', completed: 'success', approved: 'warning', rejected: 'danger' }[s] || '')



/* ==================== 生命周期 ==================== */
onMounted(() => {
  reloadProjects()
  loadStats()
  loadTeachers()
})
</script> 
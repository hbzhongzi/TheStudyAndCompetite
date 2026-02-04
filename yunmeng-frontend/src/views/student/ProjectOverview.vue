<template>
  <div class="project-overview">

    <!-- 筛选栏 -->
    <el-card class="filter-card">
      <el-row :gutter="20">
        <el-col :span="4">
          <el-select v-model="query.status" placeholder="项目状态" clearable @change="reloadProjects">
            <el-option label="草稿" value="draft" />
            <el-option label="待审核" value="pending" />
            <el-option label="已通过" value="approved" />
            <el-option label="已拒绝" value="rejected" />
          </el-select>
        </el-col>

        <el-col :span="4">
          <el-select v-model="query.type" placeholder="项目类型" clearable @change="reloadProjects">
            <el-option label="科研项目" value="科研项目" />
            <el-option label="创新项目" value="创新项目" />
            <el-option label="竞赛项目" value="竞赛项目" />
          </el-select>
        </el-col>

        <el-col :span="6">
          <el-button type="primary" @click="openCreateDialog">创建项目</el-button>
          <el-button @click="reloadProjects">刷新</el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- 统计 -->
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
      <template #header>我的项目</template>

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
            <el-progress :percentage="row.progress || 0" />
          </template>
        </el-table-column>

        <el-table-column label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatDateTime(row.createdAt) }}
          </template>
        </el-table-column>

        <!-- 操作列 -->
        <el-table-column label="操作" width="340" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="openDetail(row)">查看</el-button>

            <el-button
              size="small"
              type="success"
              :disabled="row.status !== 'approved'"
              @click="openProgress(row)"
            >
              进度
            </el-button>

            <el-button
              v-if="row.status === 'draft'"
              size="small"
              type="warning"
              @click="submitForReview(row)"
            >
              提交审核
            </el-button>

            <el-button size="small" type="danger" @click="remove(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        class="pagination"
        v-model:current-page="page.page"
        v-model:page-size="page.size"
        :total="page.total"
        layout="total, prev, pager, next"
        @current-change="reloadProjects"
      />
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

    <!-- 更新进度 -->
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
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

import studentService from '@/services/studentService'
import { projectService } from '@/services/projectService'
import { teacherService } from '@/services/teacherService'

/* ================= 项目列表 ================= */
const loading = ref(false)
const projects = ref([])
const current = ref(null)

const query = ref({ status: '', type: '' })
const page = ref({ page: 1, size: 10, total: 0 })

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

/* ================= 统计 ================= */
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

/* ================= 项目详情 ================= */
const detailVisible = ref(false)
const openDetail = row => {
  current.value = row
  detailVisible.value = true
}

/* ================= 项目进度 ================= */
const progressVisible = ref(false)
const progressForm = ref({ progress: 0, description: '' })

const openProgress = row => {
  if (row.status !== 'approved') {
    ElMessage.warning('项目通过审核后才能更新进度')
    return
  }
  current.value = row
  progressForm.value.progress = row.progress
  progressVisible.value = true
}

const submitProgress = async () => {
  await studentService.updateProjectProgress(current.value.id, progressForm.value)
  ElMessage.success('进度更新成功')
  progressVisible.value = false
  reloadProjects()
}

/* ================= 提交审核 ================= */
const submitForReview = row => {
  ElMessageBox.confirm('确认提交项目审核？', '提示', { type: 'warning' })
    .then(async () => {
      await projectService.submitProject(row.id)
      ElMessage.success('已提交审核')
      reloadProjects()
      loadStats()
    })
}

/* ================= 删除 ================= */
const remove = row => {
  ElMessageBox.confirm('确认删除该项目？', '提示', { type: 'warning' })
    .then(async () => {
      await studentService.deleteProject(row.id)
      ElMessage.success('删除成功')
      reloadProjects()
      loadStats()
    })
}

/* ================= 工具 ================= */
const statusText = s =>
  ({ draft: '草稿', submitted: '待审核', approved: '已通过', rejected: '已拒绝' }[s] || s)

const statusTag = s =>
  ({ draft: '', submitted: 'warning', approved: 'success', rejected: 'danger' }[s] || '')

const formatDateTime = d => d ? new Date(d).toLocaleString('zh-CN') : '--'

/* ================= 生命周期 ================= */
onMounted(() => {
  reloadProjects()
  loadStats()
})
</script>
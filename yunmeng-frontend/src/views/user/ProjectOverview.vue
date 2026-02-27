<template>
  <div class="admin-project-list">

    <h2>项目管理</h2>

    <el-table
      :data="projectList"
      border
      stripe
      v-loading="loading"
    >

      <el-table-column prop="id" label="ID" width="80" />

      <el-table-column prop="title" label="项目名称" />

      <el-table-column prop="type" label="类型" width="120" />

      <el-table-column label="状态" width="120">
        <template #default="{ row }">
          <el-tag
            :type="row.status === 'approved'
              ? 'success'
              : row.status === 'submitted'
              ? 'warning'
              : 'info'"
          >
            {{ row.status }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column label="学生" width="120">
        <template #default="{ row }">
          {{ row.student?.name }}
        </template>
      </el-table-column>

      <el-table-column label="学号" width="120">
        <template #default="{ row }">
          {{ row.student?.studentId }}
        </template>
      </el-table-column>

      <el-table-column label="操作" width="120">
        <template #default="{ row }">
          <el-button
            size="small"
            type="primary"
            @click="openDetail(row.id)"
          >
            查看
          </el-button>
        </template>
      </el-table-column>

    </el-table>

    <!-- 项目详情弹窗 -->
    <el-dialog
      v-model="detailVisible"
      title="项目详情"
      width="800px"
    >

      <div v-if="projectDetail">

        <!-- 基本信息 -->
        <el-divider content-position="left">项目信息</el-divider>

        <el-descriptions :column="2" border>

          <el-descriptions-item label="项目名称">
            {{ projectDetail.title }}
          </el-descriptions-item>

          <el-descriptions-item label="类型">
            {{ projectDetail.type }}
          </el-descriptions-item>

          <el-descriptions-item label="状态">
            {{ projectDetail.status }}
          </el-descriptions-item>

          <el-descriptions-item label="是否通过">
            <el-tag :type="projectDetail.isApproved ? 'success' : 'danger'">
              {{ projectDetail.isApproved ? '已通过' : '未通过' }}
            </el-tag>
          </el-descriptions-item>

          <el-descriptions-item label="创建时间">
            {{ formatDate(projectDetail.createdAt) }}
          </el-descriptions-item>

          <el-descriptions-item label="更新时间">
            {{ formatDate(projectDetail.updatedAt) }}
          </el-descriptions-item>

        </el-descriptions>

        <!-- 项目描述 -->
        <el-divider content-position="left">项目描述</el-divider>
        <p>{{ projectDetail.description }}</p >

        <!-- 项目计划 -->
        <el-divider content-position="left">项目计划</el-divider>
        <p>{{ projectDetail.plan }}</p >

        <!-- 学生信息 -->
        <el-divider content-position="left">学生信息</el-divider>

        <el-descriptions :column="2" border>
          <el-descriptions-item label="姓名">
            {{ projectDetail.student?.realName }}
          </el-descriptions-item>
          <el-descriptions-item label="学号">
            {{ projectDetail.student?.studentId }}
          </el-descriptions-item>
          <el-descriptions-item label="邮箱">
            {{ projectDetail.student?.email }}
          </el-descriptions-item>
          <el-descriptions-item label="电话">
            {{ projectDetail.student?.phone }}
          </el-descriptions-item>
          <el-descriptions-item label="学院">
            {{ projectDetail.student?.department }}
          </el-descriptions-item>
        </el-descriptions>

        <!-- 指导教师 -->
        <el-divider content-position="left">指导教师</el-divider>

        <el-descriptions :column="2" border>
          <el-descriptions-item label="姓名">
            {{ projectDetail.teacher?.realName }}
          </el-descriptions-item>
          <el-descriptions-item label="邮箱">
            {{ projectDetail.teacher?.email }}
          </el-descriptions-item>
          <el-descriptions-item label="电话">
            {{ projectDetail.teacher?.phone }}
          </el-descriptions-item>
          <el-descriptions-item label="学院">
            {{ projectDetail.teacher?.department }}
          </el-descriptions-item>
        </el-descriptions>

        <!-- 附件 -->
        <el-divider content-position="left">附件列表</el-divider>

        <el-table
          :data="projectDetail.files"
          border
          size="small"
        >

          <el-table-column prop="originalName" label="文件名" />

          <el-table-column label="大小(KB)" width="120">
            <template #default="{ row }">
              {{ (row.size / 1024).toFixed(2) }}
            </template>
          </el-table-column>

          <el-table-column label="上传时间" width="180">
            <template #default="{ row }">
              {{ formatDate(row.uploadTime) }}
            </template>
          </el-table-column>

          <el-table-column label="操作" width="100">
            <template #default="{ row }">
              <el-button
                size="small"
                type="primary"
                @click="downloadFile(row.fileUrl)"
              >
                下载
              </el-button>
            </template>
          </el-table-column>

        </el-table>

      </div>

    </el-dialog>

  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import adminService from '@/services/adminService'

const loading = ref(false)
const projectList = ref([])

const detailVisible = ref(false)
const projectDetail = ref(null)

// 加载项目列表
const loadProjects = async () => {
  loading.value = true
  try {
    const res = await adminService.getAllProjects()
    projectList.value = res.data.list
  } catch (error) {
    ElMessage.error('加载项目失败')
  }
  loading.value = false
}

// 打开详情
const openDetail = async (id) => {
  try {
    const res = await adminService.getProjectDetail(id)
    projectDetail.value = res.data
    detailVisible.value = true
  } catch (error) {
    ElMessage.error('获取项目详情失败')
  }
}

// 下载文件
const downloadFile = (url) => {
  window.open(`http://localhost:8080${url}`)
}

// 时间格式
const formatDate = (time) => {
  if (!time) return '-'
  return time.replace('T', ' ').substring(0, 16)
}

onMounted(loadProjects)
</script>

<style scoped>
.admin-project-list {
  background: #fff;
  padding: 20px;
  border-radius: 8px;
}
</style>
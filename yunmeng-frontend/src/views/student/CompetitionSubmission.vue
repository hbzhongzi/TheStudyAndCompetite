<template>
  <div class="competition-submission">
    <h2>我的竞赛报名</h2>

    <el-table
      :data="competitionList"
      border
      stripe
      v-loading="loading"
    >

      <el-table-column label="竞赛名称">
        <template #default="{ row }">
          {{ row.competition?.title }}
        </template>
      </el-table-column>

      <el-table-column label="报名时间" width="180">
        <template #default="{ row }">
          {{ formatDate(row.register_time) }}
        </template>
      </el-table-column>

      <el-table-column label="状态" width="120">
        <template #default="{ row }">
          <el-tag
            :type="row.status === 'approved'
              ? 'success'
              : row.status === 'pending'
              ? 'warning'
              : 'danger'"
          >
            {{ getStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column label="操作" width="180">
        <template #default="{ row }">

          <el-button
            v-if="row.status === 'approved'"
            size="small"
            type="primary"
            @click="openUploadDialog(row)"
          >
            提交作品
          </el-button>

          <span v-else>-</span>

        </template>
      </el-table-column>

    </el-table>

    <!-- 上传作品弹窗 -->
    <el-dialog
      v-model="uploadDialogVisible"
      title="提交作品"
      width="500px"
    >

      <el-form :model="uploadForm" label-width="80px">

        <el-form-item label="竞赛">
          {{ currentCompetition?.competition?.title }}
        </el-form-item>

        <el-form-item label="标题">
          <el-input v-model="uploadForm.title" />
        </el-form-item>

        <el-form-item label="描述">
          <el-input
            type="textarea"
            v-model="uploadForm.description"
            :rows="4"
          />
        </el-form-item>

        <el-form-item label="文件">
          <el-upload
            :auto-upload="false"
            :on-change="handleFileChange"
            :limit="1"
          >
            <el-button type="primary">选择文件</el-button>
          </el-upload>
        </el-form-item>

      </el-form>

      <template #footer>
        <el-button @click="uploadDialogVisible=false">
          取消
        </el-button>
        <el-button type="primary" @click="submitFile">
          提交
        </el-button>
      </template>

    </el-dialog>

  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import CompetitionService from '@/services/CompetitionService'

const loading = ref(false)
const competitionList = ref([])

const uploadDialogVisible = ref(false)
const currentCompetition = ref(null)

const uploadForm = ref({
  title: '',
  description: '',
  file: null
})


// 加载我的报名记录
const loadMyCompetitions = async () => {
  loading.value = true
  try {
    const res = await CompetitionService.getMyRegistrations()
    competitionList.value = res.data.list
  } catch (error) {
    ElMessage.error('加载报名记录失败')
  }
  loading.value = false
}


// 打开上传弹窗
const openUploadDialog = (row) => {
  currentCompetition.value = row
  uploadForm.value = {
    title: '',
    description: '',
    file: null
  }
  uploadDialogVisible.value = true
}


// 选择文件
const handleFileChange = (file) => {
  uploadForm.value.file = file.raw
}


// 提交作品
const submitFile = async () => {
  if (!uploadForm.value.title ||
      !uploadForm.value.file) {
    ElMessage.warning('请填写完整信息')
    return
  }

  const formData = new FormData()
  formData.append('CompetitionID', currentCompetition.value.competition_id)
  formData.append('Title', uploadForm.value.title)
  formData.append('Description', uploadForm.value.description)
  formData.append('file', uploadForm.value.file)

  try {
    await CompetitionService.submitCompetitionResult(formData)

    ElMessage.success('作品提交成功')
    uploadDialogVisible.value = false
    loadMyCompetitions()

  } catch (error) {
    ElMessage.error('提交失败')
  }
}


// 状态文字
const getStatusText = (status) => {
  if (status === 'approved') return '已通过'
  if (status === 'pending') return '审核中'
  return '已拒绝'
}


// 时间格式
const formatDate = (time) => {
  if (!time) return '-'
  return time.replace('T', ' ').substring(0, 16)
}

onMounted(loadMyCompetitions)
</script>

<style scoped>
.student-submission {
  background: #fff;
  padding: 20px;
  border-radius: 8px;
}
</style>
<template>
  <div class="project-file-list">
    <div v-if="files.length === 0" class="empty-files">
      <el-empty description="暂无文件" :image-size="60" />
    </div>
    
    <div v-else class="files-container">
      <el-table :data="files" style="width: 100%">
        <el-table-column prop="fileName" label="文件名称" min-width="200">
          <template #default="{ row }">
            <div class="file-name">
              <el-icon class="file-icon">
                <Document v-if="row.fileType === 'proposal' || row.fileType === 'midterm' || row.fileType === 'final'" />
                <Picture v-else-if="row.fileType === 'achievement'" />
                <Files v-else />
              </el-icon>
              <span>{{ row.fileName }}</span>
            </div>
          </template>
        </el-table-column>
        
        <el-table-column prop="fileType" label="文件类型" width="120">
          <template #default="{ row }">
            <el-tag :type="getFileTypeTagType(row.fileType)" size="small">
              {{ getFileTypeText(row.fileType) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="fileVersion" label="版本" width="80" />
        
        <el-table-column prop="reviewStatus" label="审核状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getReviewStatusType(row.reviewStatus)" size="small">
              {{ getReviewStatusText(row.reviewStatus) }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column prop="uploadTime" label="上传时间" width="160">
          <template #default="{ row }">
            {{ formatDate(row.uploadTime) }}
          </template>
        </el-table-column>
        
        <el-table-column prop="fileSize" label="文件大小" width="100">
          <template #default="{ row }">
            {{ formatFileSize(row.fileSize) }}
          </template>
        </el-table-column>
        
        <el-table-column prop="downloadCount" label="下载次数" width="100" />
        
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button 
              size="small" 
              type="primary" 
              @click="downloadFile(row)"
              icon="Download"
            >
              下载
            </el-button>
            
            <el-button 
              v-if="canReview && row.reviewStatus === 'pending'"
              size="small" 
              type="warning" 
              @click="showReviewDialog(row)"
              icon="Edit"
            >
              审核
            </el-button>
            
            <el-button 
              v-if="canReview && row.reviewStatus !== 'pending'"
              size="small" 
              type="info" 
              @click="viewReviewDetails(row)"
              icon="View"
            >
              查看审核
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
    
    <!-- 文件审核对话框 -->
    <el-dialog v-model="showReviewDialog" title="审核文件" width="500px">
      <el-form :model="reviewForm" :rules="reviewRules" label-width="100px" ref="reviewFormRef">
        <el-form-item label="审核结果" prop="reviewStatus">
          <el-radio-group v-model="reviewForm.reviewStatus">
            <el-radio label="approved">通过</el-radio>
            <el-radio label="rejected">驳回</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item label="审核意见" prop="reviewComments">
          <el-input
            v-model="reviewForm.reviewComments"
            type="textarea"
            :rows="4"
            placeholder="请输入审核意见"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="showReviewDialog = false">取消</el-button>
        <el-button type="primary" @click="submitReview">确定</el-button>
      </template>
    </el-dialog>
    
    <!-- 审核详情对话框 -->
    <el-dialog v-model="showReviewDetailsDialog" title="审核详情" width="500px">
      <el-descriptions :column="1" border>
        <el-descriptions-item label="审核状态">
          <el-tag :type="getReviewStatusType(currentFile.reviewStatus)">
            {{ getReviewStatusText(currentFile.reviewStatus) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="审核意见">
          {{ currentFile.reviewComments || '无' }}
        </el-descriptions-item>
        <el-descriptions-item label="审核时间">
          {{ formatDate(currentFile.reviewedAt) }}
        </el-descriptions-item>
        <el-descriptions-item label="审核人">
          {{ currentFile.reviewedBy || '未知' }}
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script>
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { Document, Picture, Files } from '@element-plus/icons-vue'

export default {
  name: 'ProjectFileList',
  components: {
    Document,
    Picture,
    Files
  },
  props: {
    files: {
      type: Array,
      default: () => []
    },
    canReview: {
      type: Boolean,
      default: false
    }
  },
  emits: ['review', 'download'],
  setup(props, { emit }) {
    // 响应式数据
    const showReviewDialog = ref(false)
    const showReviewDetailsDialog = ref(false)
    const currentFile = ref({})
    
    const reviewForm = reactive({
      reviewStatus: 'approved',
      reviewComments: ''
    })
    
    const reviewRules = {
      reviewStatus: [
        { required: true, message: '请选择审核结果', trigger: 'change' }
      ],
      reviewComments: [
        { required: true, message: '请输入审核意见', trigger: 'blur' }
      ]
    }
    
    // 方法
    const downloadFile = (file) => {
      emit('download', file)
    }
    
    const openReviewDialog = (file) => {
      currentFile.value = file
      showReviewDialog.value = true
      // 重置表单
      Object.assign(reviewForm, {
        reviewStatus: 'approved',
        reviewComments: ''
      })
    }
    
    const submitReview = async () => {
      try {
        const reviewData = {
          reviewStatus: reviewForm.reviewStatus,
          reviewComments: reviewForm.reviewComments
        }
        
        emit('review', currentFile.value.id, reviewData)
        showReviewDialog.value = false
      } catch (error) {
        ElMessage.error('审核提交失败')
      }
    }
    
    const viewReviewDetails = (file) => {
      currentFile.value = file
      showReviewDetailsDialog.value = true
    }
    
    // 工具方法
    const formatDate = (date) => {
      if (!date) return ''
      return new Date(date).toLocaleString('zh-CN')
    }
    
    const formatFileSize = (bytes) => {
      if (!bytes) return '0 B'
      const k = 1024
      const sizes = ['B', 'KB', 'MB', 'GB']
      const i = Math.floor(Math.log(bytes) / Math.log(k))
      return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
    }
    
    const getFileTypeTagType = (fileType) => {
      const typeMap = {
        proposal: 'primary',
        midterm: 'warning',
        final: 'success',
        achievement: 'info',
        other: 'default'
      }
      return typeMap[fileType] || 'default'
    }
    
    const getFileTypeText = (fileType) => {
      const typeMap = {
        proposal: '开题报告',
        midterm: '中期报告',
        final: '结题报告',
        achievement: '成果展示',
        other: '其他材料'
      }
      return typeMap[fileType] || fileType
    }
    
    const getReviewStatusType = (status) => {
      const typeMap = {
        pending: 'warning',
        approved: 'success',
        rejected: 'danger'
      }
      return typeMap[status] || 'info'
    }
    
    const getReviewStatusText = (status) => {
      const statusMap = {
        pending: '待审核',
        approved: '已通过',
        rejected: '已驳回'
      }
      return statusMap[status] || status
    }
    
    return {
      // 响应式数据
      showReviewDialog,
      showReviewDetailsDialog,
      currentFile,
      reviewForm,
      reviewRules,
      
      // 方法
      downloadFile,
      openReviewDialog,
      submitReview,
      viewReviewDetails,
      
      // 工具方法
      formatDate,
      formatFileSize,
      getFileTypeTagType,
      getFileTypeText,
      getReviewStatusType,
      getReviewStatusText
    }
  }
}
</script>

<style scoped>
.project-file-list {
  width: 100%;
}

.empty-files {
  text-align: center;
  padding: 40px 0;
}

.files-container {
  margin-top: 20px;
}

.file-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.file-icon {
  font-size: 16px;
  color: #409eff;
}

.el-table {
  margin-top: 10px;
}

.el-table .el-button {
  margin-right: 5px;
}

.el-table .el-button:last-child {
  margin-right: 0;
}
</style> 
<template>
  <div class="project-management">
    <!-- 页面标题和统计信息 -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">项目管理</h1>
        <p class="page-subtitle">管理您指导的学生项目</p>
      </div>
      <div class="header-stats">
        <div class="stat-item">
          <div class="stat-number">{{ projectStats.total || 0 }}</div>
          <div class="stat-label">总项目数</div>
        </div>
        <div class="stat-item">
          <div class="stat-number">{{ projectStats.pending || 0 }}</div>
          <div class="stat-label">待审核</div>
        </div>
        <div class="stat-item">
          <div class="stat-number">{{ projectStats.approved || 0 }}</div>
          <div class="stat-label">已通过</div>
        </div>
        <div class="stat-item">
          <div class="stat-number">{{ projectStats.inProgress || 0 }}</div>
          <div class="stat-label">进行中</div>
        </div>
      </div>
    </div>

    <!-- 新增：项目统计图表 -->
    <el-card class="stats-chart-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <span class="header-title">项目统计概览</span>
          <div class="header-actions">
            <el-button size="small" @click="refreshStats" :loading="statsLoading">
              <el-icon><Refresh /></el-icon>
              刷新统计
            </el-button>
          </div>
        </div>
      </template>
      
      <el-row :gutter="20">
        <el-col :span="12">
          <div class="chart-container">
            <h4>项目状态分布</h4>
            <div class="pie-chart" ref="statusChartRef"></div>
          </div>
        </el-col>
        <el-col :span="12">
          <div class="chart-container">
            <h4>项目类型分布</h4>
            <div class="pie-chart" ref="typeChartRef"></div>
          </div>
        </el-col>
      </el-row>
      
      <el-row :gutter="20" style="margin-top: 20px;">
        <el-col :span="24">
          <div class="chart-container">
            <h4>项目进度趋势</h4>
            <div class="line-chart" ref="progressChartRef"></div>
          </div>
        </el-col>
      </el-row>
    </el-card>

    <!-- 新增：快速操作面板 -->
    <el-card class="quick-actions-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <span class="header-title">快速操作</span>
        </div>
      </template>
      
      <el-row :gutter="20">
        <el-col :span="6">
          <el-button type="primary" @click="batchApprove" :disabled="!hasSelectedProjects" class="quick-action-btn">
            <el-icon><Check /></el-icon>
            批量通过
          </el-button>
        </el-col>
        <el-col :span="6">
          <el-button type="warning" @click="batchReview" :disabled="!hasSelectedProjects" class="quick-action-btn">
            <el-icon><Edit /></el-icon>
            批量审核
          </el-button>
        </el-col>
        <el-col :span="6">
          <el-button type="info" @click="exportProjects" class="quick-action-btn">
            <el-icon><Download /></el-icon>
            导出数据
          </el-button>
        </el-col>
        <el-col :span="6">
          <el-button type="success" @click="generateReport" class="quick-action-btn">
            <el-icon><Document /></el-icon>
            生成报告
          </el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- 搜索和筛选 -->
    <el-card class="filter-card" shadow="hover">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-input
            v-model="searchQuery"
            placeholder="搜索项目标题、学生姓名或描述"
            clearable
            @input="handleSearch"
            class="search-input"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-col>
        <el-col :span="4">
          <el-select v-model="statusFilter" placeholder="状态筛选" clearable @change="handleFilter" class="filter-select">
            <el-option label="全部状态" value="" />
            <el-option label="草稿" value="draft" />
            <el-option label="已提交" value="submitted" />
            <el-option label="待审核" value="pending" />
            <el-option label="已通过" value="approved" />
            <el-option label="已拒绝" value="rejected" />
            <el-option label="进行中" value="in_progress" />
            <el-option label="已完成" value="completed" />
            <el-option label="已暂停" value="suspended" />
            <el-option label="需修改" value="need_revision" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-select v-model="typeFilter" placeholder="类型筛选" clearable @change="handleFilter" class="filter-select">
            <el-option label="全部类型" value="" />
            <el-option label="科研" value="科研" />
            <el-option label="竞赛" value="竞赛" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-select v-model="levelFilter" placeholder="级别筛选" clearable @change="handleFilter" class="filter-select">
            <el-option label="全部级别" value="" />
            <el-option label="校级" value="校级" />
            <el-option label="省级" value="省级" />
            <el-option label="国家级" value="国家级" />
            <el-option label="国际级" value="国际级" />
          </el-select>
        </el-col>
        <el-col :span="6">
          <el-button type="primary" @click="loadProjects" :loading="loading" class="refresh-btn">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
          <el-button @click="resetFilters" class="reset-btn">
            <el-icon><Delete /></el-icon>
            重置
          </el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- 项目列表 -->
    <el-card class="project-list" shadow="hover">
      <template #header>
        <div class="card-header">
          <span class="header-title">项目列表</span>
          <div class="header-actions">
            <span class="project-count">共 {{ filteredProjects.length }} 个项目</span>
          </div>
        </div>
      </template>
      
      <el-table
        :data="filteredProjects"
        v-loading="loading"
        style="width: 100%"
        @row-click="handleRowClick"
        :row-class-name="getRowClassName"
        stripe
        highlight-current-row
        @selection-change="handleSelectionChange"
      >
        <!-- 新增：多选列 -->
        <el-table-column type="selection" width="55" />
        
        <el-table-column prop="title" label="项目标题" min-width="200" show-overflow-tooltip>
          <template #default="{ row }">
            <div class="project-title">
              <span class="title-text">{{ row.title }}</span>
              <el-tag v-if="row.isExtended" size="small" type="warning" class="extension-tag">延期</el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="studentName" label="学生" width="120" show-overflow-tooltip>
          <template #default="{ row }">
            <div class="student-info">
              <el-avatar :size="24" class="student-avatar">{{ row.studentName?.charAt(0) }}</el-avatar>
              <span class="student-name">{{ row.studentName }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.type === '科研' ? 'primary' : 'success'" size="small">
              {{ row.type }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="level" label="级别" width="100">
          <template #default="{ row }">
            <el-tag v-if="row.level" :type="getLevelType(row.level)" size="small">
              {{ row.level }}
            </el-tag>
            <span v-else class="no-data">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="progress" label="进度" width="120">
          <template #default="{ row }">
            <div class="progress-container">
              <el-progress 
                :percentage="row.progress || 0" 
                :status="getProgressStatus(row.progress)"
                :stroke-width="8"
                class="progress-bar"
              />
              <span class="progress-text">{{ row.progress || 0 }}%</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="expectedEndDate" label="预计完成" width="120">
          <template #default="{ row }">
            <span v-if="row.expectedEndDate" class="date-text">
              {{ formatDate(row.expectedEndDate) }}
            </span>
            <span v-else class="no-data">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="160">
          <template #default="{ row }">
            <span class="date-text">{{ formatDate(row.createdAt) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <div class="action-buttons">
              <el-button size="small" @click.stop="viewProject(row)" class="view-btn">
                <el-icon><View /></el-icon>
                查看
              </el-button>
            <el-button 
                v-if="row.status === 'submitted' || row.status === 'pending'"
              size="small" 
              type="success" 
              @click.stop="approveProject(row)"
                class="approve-btn"
            >
                <el-icon><Check /></el-icon>
              通过
            </el-button>
            <el-button 
                v-if="row.status === 'submitted' || row.status === 'pending'"
              size="small" 
              type="danger" 
              @click.stop="rejectProject(row)"
                class="reject-btn"
            >
                <el-icon><Close /></el-icon>
              驳回
            </el-button>
              <el-button 
                v-if="row.status === 'approved' || row.status === 'in_progress'"
                size="small" 
                type="warning" 
                @click.stop="updateProgress(row)"
                class="progress-btn"
              >
                <el-icon><Edit /></el-icon>
                进度
              </el-button>
              <!-- 新增：质量评估按钮 -->
              <el-button 
                v-if="row.status === 'in_progress' || row.status === 'completed'"
                size="small" 
                type="info" 
                @click.stop="assessQuality(row)"
                class="quality-btn"
              >
                <el-icon><Star /></el-icon>
                评估
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container" v-if="totalProjects > 0">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="totalProjects"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
          background
        />
      </div>
    </el-card>

    <!-- 项目详情对话框 -->
    <el-dialog
      v-model="detailVisible"
      title="项目详情"
      width="70%"
      :before-close="handleCloseDetail"
      class="project-detail-dialog"
    >
      <div v-if="currentProject" class="project-detail">
        <!-- 基本信息 -->
        <el-card class="detail-section" shadow="never">
          <template #header>
            <div class="section-header">
              <el-icon><InfoFilled /></el-icon>
              <span>基本信息</span>
            </div>
          </template>
          <el-descriptions :column="3" border>
            <el-descriptions-item label="项目标题" :span="3">
              <span class="project-title-large">{{ currentProject.title }}</span>
            </el-descriptions-item>
          <el-descriptions-item label="学生姓名">{{ currentProject.studentName }}</el-descriptions-item>
            <el-descriptions-item label="项目类型">
              <el-tag :type="currentProject.type === '科研' ? 'primary' : 'success'">
                {{ currentProject.type }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="项目级别">
              <el-tag v-if="currentProject.level" :type="getLevelType(currentProject.level)">
                {{ currentProject.level }}
              </el-tag>
              <span v-else class="no-data">-</span>
            </el-descriptions-item>
          <el-descriptions-item label="项目状态">
            <el-tag :type="getStatusType(currentProject.status)">
              {{ getStatusText(currentProject.status) }}
            </el-tag>
          </el-descriptions-item>
            <el-descriptions-item label="项目进度">
              <el-progress 
                :percentage="currentProject.progress || 0" 
                :status="getProgressStatus(currentProject.progress)"
                :stroke-width="12"
              />
            </el-descriptions-item>
            <el-descriptions-item label="预计完成时间">
              <span v-if="currentProject.expectedEndDate">{{ formatDate(currentProject.expectedEndDate) }}</span>
              <span v-else class="no-data">-</span>
            </el-descriptions-item>
            <el-descriptions-item label="实际完成时间">
              <span v-if="currentProject.actualEndDate">{{ formatDate(currentProject.actualEndDate) }}</span>
              <span v-else class="no-data">-</span>
            </el-descriptions-item>
            <el-descriptions-item label="是否延期">
              <el-tag v-if="currentProject.isExtended" type="warning">是</el-tag>
              <el-tag v-else type="info">否</el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="项目描述" :span="3">
              <div class="description-content">{{ currentProject.description || '暂无描述' }}</div>
            </el-descriptions-item>
        </el-descriptions>
        </el-card>

        <!-- 项目成员 -->
        <el-card class="detail-section" shadow="never" v-if="currentProject.members && currentProject.members.length > 0">
          <template #header>
            <div class="section-header">
              <el-icon><User /></el-icon>
              <span>项目成员</span>
        </div>
          </template>
          <el-table :data="currentProject.members" style="width: 100%" stripe>
            <el-table-column prop="name" label="姓名" width="120" />
            <el-table-column prop="studentNumber" label="学号" width="120" />
            <el-table-column prop="role" label="角色" width="120" />
          </el-table>
        </el-card>

                 <!-- 项目附件 -->
        <el-card class="detail-section" shadow="never" v-if="currentProject.files && currentProject.files.length > 0">
          <template #header>
            <div class="section-header">
              <el-icon><Document /></el-icon>
              <span>项目附件</span>
            </div>
          </template>
          <el-table :data="currentProject.files" style="width: 100%" stripe>
            <el-table-column prop="fileName" label="文件名" min-width="200" />
            <el-table-column prop="fileType" label="类型" width="100" />
            <el-table-column prop="fileSize" label="大小" width="100">
               <template #default="{ row }">
                <span>{{ formatFileSize(row.fileSize) }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="uploadTime" label="上传时间" width="160">
              <template #default="{ row }">
                <span>{{ formatDate(row.uploadTime) }}</span>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="120">
              <template #default="{ row }">
                <el-button size="small" @click="downloadFile(row)" type="primary">
                  <el-icon><Download /></el-icon>
                  下载
                </el-button>
               </template>
             </el-table-column>
           </el-table>
        </el-card>

         <!-- 审核记录 -->
        <el-card class="detail-section" shadow="never" v-if="currentProject.reviewRecords && currentProject.reviewRecords.length > 0">
          <template #header>
            <div class="section-header">
              <el-icon><Clock /></el-icon>
              <span>审核记录</span>
            </div>
          </template>
           <el-timeline>
             <el-timeline-item
               v-for="record in currentProject.reviewRecords"
               :key="record.id"
               :timestamp="formatDate(record.reviewTime)"
               :type="record.status === 'approved' ? 'success' : 'danger'"
              :color="record.status === 'approved' ? '#67C23A' : '#F56C6C'"
             >
              <el-card class="timeline-card" shadow="hover">
                <div class="review-header">
                  <span class="reviewer-name">{{ record.reviewer || record.reviewerName }}</span>
                  <el-tag :type="record.status === 'approved' ? 'success' : 'danger'" size="small">
                    {{ record.status === 'approved' ? '通过' : '驳回' }}
                  </el-tag>
                </div>
                <div class="review-comments">{{ record.comments || record.comment }}</div>
               </el-card>
             </el-timeline-item>
           </el-timeline>
        </el-card>
      </div>
    </el-dialog>

    <!-- 审核对话框 -->
    <el-dialog
      v-model="reviewVisible"
      :title="reviewType === 'approve' ? '通过项目' : '驳回项目'"
      width="50%"
      class="review-dialog"
    >
      <el-form :model="reviewForm" label-width="80px" :rules="reviewRules" ref="reviewFormRef">
        <el-form-item label="审核意见" prop="comment">
          <el-input
            v-model="reviewForm.comment"
            type="textarea"
            :rows="4"
            :placeholder="reviewType === 'approve' ? '请输入通过意见（可选）' : '请输入驳回原因'"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="reviewVisible = false">取消</el-button>
          <el-button type="primary" @click="submitReview" :loading="submitting">
            {{ reviewType === 'approve' ? '确认通过' : '确认驳回' }}
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 进度更新对话框 -->
    <el-dialog
      v-model="progressVisible"
      title="更新项目进度"
      width="50%"
      class="progress-dialog"
    >
      <el-form :model="progressForm" label-width="100px" :rules="progressRules" ref="progressFormRef">
        <el-form-item label="当前进度" prop="progress">
          <el-slider
            v-model="progressForm.progress"
            :min="0"
            :max="100"
            :step="5"
            show-input
            input-size="large"
            :format-tooltip="(val) => val + '%'"
          />
        </el-form-item>
        <el-form-item label="进度说明" prop="description">
          <el-input
            v-model="progressForm.description"
            type="textarea"
            :rows="3"
            placeholder="请描述项目当前进展情况和遇到的问题"
            maxlength="300"
            show-word-limit
          />
        </el-form-item>
        <el-form-item label="预计完成" prop="expectedEndDate">
          <el-date-picker
            v-model="progressForm.expectedEndDate"
            type="datetime"
            placeholder="选择预计完成时间"
            format="YYYY-MM-DD HH:mm"
            value-format="YYYY-MM-DD HH:mm:ss"
            :disabled-date="disabledDate"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="progressVisible = false">取消</el-button>
          <el-button type="primary" @click="submitProgress" :loading="submitting">
            更新进度
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 新增：批量操作对话框 -->
    <el-dialog
      v-model="batchOperationVisible"
      title="批量操作"
      width="60%"
      class="batch-operation-dialog"
    >
      <div class="batch-info">
        <el-alert
          title="批量操作提示"
          type="info"
          :closable="false"
          show-icon
        >
          <template #default>
            已选择 <strong>{{ selectedProjects.length }}</strong> 个项目进行批量操作
          </template>
        </el-alert>
      </div>
      
      <el-form :model="batchForm" label-width="120px" :rules="batchRules" ref="batchFormRef">
        <el-form-item label="操作类型" prop="operationType">
          <el-radio-group v-model="batchForm.operationType">
            <el-radio label="approve">批量通过</el-radio>
            <el-radio label="reject">批量驳回</el-radio>
            <el-radio label="updateProgress">批量更新进度</el-radio>
            <el-radio label="addComment">批量添加评语</el-radio>
          </el-radio-group>
        </el-form-item>
        
        <el-form-item label="操作内容" prop="content" v-if="batchForm.operationType === 'addComment'">
          <el-input
            v-model="batchForm.content"
            type="textarea"
            :rows="4"
            placeholder="请输入要添加的评语内容"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>
        
        <el-form-item label="进度值" prop="progress" v-if="batchForm.operationType === 'updateProgress'">
          <el-slider
            v-model="batchForm.progress"
            :min="0"
            :max="100"
            :step="5"
            show-input
            input-size="large"
            :format-tooltip="(val) => val + '%'"
          />
        </el-form-item>
        
        <el-form-item label="驳回原因" prop="rejectReason" v-if="batchForm.operationType === 'reject'">
          <el-input
            v-model="batchForm.rejectReason"
            type="textarea"
            :rows="4"
            placeholder="请输入驳回原因"
            maxlength="500"
            show-word-limit
            required
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="batchOperationVisible = false">取消</el-button>
          <el-button type="primary" @click="submitBatchOperation" :loading="submitting">
            确认执行
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 新增：项目质量评估对话框 -->
    <el-dialog
      v-model="qualityAssessmentVisible"
      title="项目质量评估"
      width="70%"
      class="quality-assessment-dialog"
    >
      <div v-if="currentProject" class="quality-assessment">
        <el-form :model="qualityForm" label-width="120px" :rules="qualityRules" ref="qualityFormRef">
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="项目创新性" prop="innovation">
                <el-rate
                  v-model="qualityForm.innovation"
                  :max="5"
                  show-text
                  :texts="['很差', '较差', '一般', '较好', '很好']"
                />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="技术可行性" prop="feasibility">
                <el-rate
                  v-model="qualityForm.feasibility"
                  :max="5"
                  show-text
                  :texts="['很差', '较差', '一般', '较好', '很好']"
                />
              </el-form-item>
            </el-col>
          </el-row>
          
          <el-row :gutter="20">
            <el-col :span="12">
              <el-form-item label="完成质量" prop="completion">
                <el-rate
                  v-model="qualityForm.completion"
                  :max="5"
                  show-text
                  :texts="['很差', '较差', '一般', '较好', '很好']"
                />
              </el-form-item>
            </el-col>
            <el-col :span="12">
              <el-form-item label="团队协作" prop="collaboration">
                <el-rate
                  v-model="qualityForm.collaboration"
                  :max="5"
                  show-text
                  :texts="['很差', '较差', '一般', '较好', '很好']"
                />
              </el-form-item>
            </el-col>
          </el-row>
          
          <el-form-item label="综合评语" prop="overallComment">
            <el-input
              v-model="qualityForm.overallComment"
              type="textarea"
              :rows="4"
              placeholder="请输入综合评语和建议"
              maxlength="1000"
              show-word-limit
            />
          </el-form-item>
          
          <el-form-item label="改进建议" prop="suggestions">
            <el-input
              v-model="qualityForm.suggestions"
              type="textarea"
              :rows="3"
              placeholder="请输入具体的改进建议"
              maxlength="500"
              show-word-limit
            />
          </el-form-item>
        </el-form>
      </div>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="qualityAssessmentVisible = false">取消</el-button>
          <el-button type="primary" @click="submitQualityAssessment" :loading="submitting">
            提交评估
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Delete, User, Document, Download, Clock, InfoFilled, Edit, View, Check, Close, Star } from '@element-plus/icons-vue'
import { teacherService, fileService } from '../../services/teacherService'
import { validateApiResponse, safeArrayFilter } from '../../utils/dataValidator'

// 响应式数据
const loading = ref(false)
const statsLoading = ref(false)
const searchQuery = ref('')
const statusFilter = ref('')
const typeFilter = ref('')
const levelFilter = ref('')
const projects = ref([]) // 确保初始化为空数组
const detailVisible = ref(false)
const reviewVisible = ref(false)
const progressVisible = ref(false)
const batchOperationVisible = ref(false)
const qualityAssessmentVisible = ref(false)
const currentProject = ref(null)
const reviewType = ref('approve')
const submitting = ref(false)
const selectedProjects = ref([])

const reviewForm = ref({
  comment: ''
})

const progressForm = ref({
  progress: 0,
  description: '',
  expectedEndDate: null
})

const batchForm = ref({
  operationType: 'approve',
  content: '',
  progress: 50,
  rejectReason: ''
})

const qualityForm = ref({
  innovation: 0,
  feasibility: 0,
  completion: 0,
  collaboration: 0,
  overallComment: '',
  suggestions: ''
})

// 表单引用
const reviewFormRef = ref()
const progressFormRef = ref()
const batchFormRef = ref()
const qualityFormRef = ref()

// 表单验证规则
const reviewRules = {
  comment: [
    { required: true, message: '请输入审核意见', trigger: 'blur' },
    { max: 500, message: '审核意见不能超过500个字符', trigger: 'blur' }
  ]
}

const progressRules = {
  progress: [
    { required: true, message: '请选择当前进度', trigger: 'change' }
  ],
  expectedEndDate: [
    { required: true, message: '请选择预计完成时间', trigger: 'change' }
  ]
}

const batchRules = {
  operationType: [
    { required: true, message: '请选择操作类型', trigger: 'change' }
  ],
  content: [
    { required: true, message: '请输入评语内容', trigger: 'blur' },
    { max: 500, message: '评语内容不能超过500个字符', trigger: 'blur' }
  ],
  rejectReason: [
    { required: true, message: '请输入驳回原因', trigger: 'blur' },
    { max: 500, message: '驳回原因不能超过500个字符', trigger: 'blur' }
  ]
}

const qualityRules = {
  innovation: [
    { required: true, message: '请评估项目创新性', trigger: 'change' }
  ],
  feasibility: [
    { required: true, message: '请评估技术可行性', trigger: 'change' }
  ],
  completion: [
    { required: true, message: '请评估完成质量', trigger: 'change' }
  ],
  collaboration: [
    { required: true, message: '请评估团队协作', trigger: 'change' }
  ],
  overallComment: [
    { required: true, message: '请输入综合评语', trigger: 'blur' },
    { max: 1000, message: '综合评语不能超过1000个字符', trigger: 'blur' }
  ]
}

// 模拟项目数据
const mockProjects = [
  {
    id: 1,
    title: '智能校园管理系统',
    studentName: '张三',
    type: '科研',
    status: 'pending',
    description: '基于物联网技术的智能校园管理系统，实现校园设施的智能化管理。',
    createTime: '2024-01-15 10:30:00',
    updateTime: '2024-01-15 14:20:00',
    members: [
      { name: '张三', studentNumber: '2021001', role: '负责人' },
      { name: '李四', studentNumber: '2021002', role: '成员' }
    ],
    attachments: [
      { fileName: '项目申请书.pdf', fileSize: '2.5MB' },
      { fileName: '技术方案.docx', fileSize: '1.8MB' }
    ],
    reviewRecords: []
  },
  {
    id: 2,
    title: '在线学习平台',
    studentName: '王五',
    type: '创新',
    status: 'approved',
    description: '基于Web技术的在线学习平台，支持课程管理和在线学习。',
    createTime: '2024-01-10 09:15:00',
    updateTime: '2024-01-12 16:45:00',
    members: [
      { name: '王五', studentNumber: '2021003', role: '负责人' }
    ],
    attachments: [
      { fileName: '项目计划书.pdf', fileSize: '3.2MB' }
    ],
    reviewRecords: [
      {
        id: 1,
                  reviewer: '李教授',
        status: 'approved',
                  comments: '项目方案合理，技术可行，同意通过。',
        reviewTime: '2024-01-12 16:45:00'
      }
    ]
  },
  {
    id: 3,
    title: '全国大学生程序设计竞赛',
    studentName: '赵六',
    type: '竞赛',
    status: 'rejected',
    description: '参加全国大学生程序设计竞赛，提升算法和编程能力。',
    createTime: '2024-01-08 11:20:00',
    updateTime: '2024-01-09 15:30:00',
    members: [
      { name: '赵六', studentNumber: '2021004', role: '参赛者' }
    ],
    attachments: [
      { fileName: '竞赛报名表.pdf', fileSize: '0.8MB' }
    ],
    reviewRecords: [
      {
        id: 2,
        reviewerName: '王教授',
        result: 'rejected',
        comment: '竞赛报名材料不完整，请补充相关证明材料后重新提交。',
        reviewTime: '2024-01-09 15:30:00'
      }
    ]
  }
]

// 计算属性
const filteredProjects = computed(() => {
  // 确保projects.value是数组
  let result = Array.isArray(projects.value) ? projects.value : []

  // 状态筛选
  if (statusFilter.value) {
    result = result.filter(p => p && p.status === statusFilter.value)
  }

  // 类型筛选
  if (typeFilter.value) {
    result = result.filter(p => p && p.type === typeFilter.value)
  }

  // 级别筛选
  if (levelFilter.value) {
    result = result.filter(p => p && p.level === levelFilter.value)
  }

  // 搜索筛选
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(p => 
      p && p.title && p.title.toLowerCase().includes(query) ||
      p && p.studentName && p.studentName.toLowerCase().includes(query) ||
      p && p.description && p.description.toLowerCase().includes(query)
    )
  }

  return result
})

const projectStats = computed(() => {
  const total = projects.value.length;
  const pending = projects.value.filter(p => p.status === 'pending').length;
  const approved = projects.value.filter(p => p.status === 'approved').length;
  const inProgress = projects.value.filter(p => p.status === 'in_progress').length;
  return { total, pending, approved, inProgress };
});

const totalProjects = computed(() => {
  return projects.value.length;
});

const hasSelectedProjects = computed(() => {
  return selectedProjects.value.length > 0;
});

const selectedProjectsCount = computed(() => {
  return selectedProjects.value.length;
});

const currentPage = ref(1);
const pageSize = ref(10);

// 方法
const loadProjects = async () => {
  loading.value = true
  try {
    const params = {
      page: currentPage.value,
      pageSize: pageSize.value,
    }
    if (statusFilter.value) params.status = statusFilter.value
    if (typeFilter.value) params.type = typeFilter.value
    if (levelFilter.value) params.level = levelFilter.value
    if (searchQuery.value) params.search = searchQuery.value
    
    // 使用getMyProjects获取我指导的项目列表
    const response = await teacherService.getMyProjects(params)
    const validation = validateApiResponse(response)
    
    if (validation.isValid) {
      projects.value = safeArrayFilter(validation.data)
      ElMessage.success('项目列表加载成功')
    } else {
      projects.value = []
      ElMessage.error(validation.error || '加载项目列表失败')
    }
  } catch (error) {
    console.error('加载项目列表失败:', error)
    ElMessage.error(error.message || '加载项目列表失败')
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

const resetFilters = () => {
  searchQuery.value = ''
  statusFilter.value = ''
  typeFilter.value = ''
  levelFilter.value = ''
  currentPage.value = 1
  loadProjects()
}

const handleRowClick = (row) => {
  viewProject(row)
}

const handleSelectionChange = (selection) => {
  selectedProjects.value = selection;
};

const handleCurrentChange = (page) => {
  currentPage.value = page
  loadProjects()
}

const handleSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
  loadProjects()
}

const viewProject = async (project) => {
  try {
    // 获取项目详情和审核记录
    const [detailResponse, reviewsResponse] = await Promise.all([
      teacherService.getProjectDetail(project.id),
      teacherService.getProjectReviews(project.id)
    ])
    
    if (detailResponse.code === 200) {
      const projectData = detailResponse.data
      const reviewData = reviewsResponse.code === 200 ? reviewsResponse.data : []
      
      // 确保reviewRecords是数组
      const reviewRecords = Array.isArray(reviewData) ? reviewData : []
      
      currentProject.value = {
        ...projectData,
        reviewRecords: reviewRecords
      }
      detailVisible.value = true
    } else {
      ElMessage.error(detailResponse.message || '获取项目详情失败')
    }
  } catch (error) {
    console.error('获取项目详情失败:', error)
    ElMessage.error(error.message || '获取项目详情失败')
  }
}

const handleCloseDetail = () => {
  detailVisible.value = false
  currentProject.value = null
}

const approveProject = (project) => {
  reviewType.value = 'approve'
  currentProject.value = project
  reviewForm.value.comment = ''
  reviewVisible.value = true
}

const rejectProject = (project) => {
  reviewType.value = 'reject'
  currentProject.value = project
  reviewForm.value.comment = ''
  reviewVisible.value = true
}

const submitReview = async () => {
  try {
    await reviewFormRef.value.validate();
  } catch (error) {
    return;
  }

  if (reviewType.value === 'reject' && !reviewForm.value.comment.trim()) {
    ElMessage.warning('请输入驳回原因')
    return
  }

  submitting.value = true
  try {
    const reviewData = {
      status: reviewType.value === 'approve' ? 'approved' : 'rejected',
      comments: reviewForm.value.comment || (reviewType.value === 'approve' ? '审核通过' : '审核不通过')
    }
    
    const response = await teacherService.reviewProject(currentProject.value.id, reviewData)
    
    if (response.code === 200) {
      // 更新本地数据
      const project = projects.value.find(p => p.id === currentProject.value.id)
      if (project) {
        project.status = reviewType.value === 'approve' ? 'approved' : 'rejected'
        project.updateTime = new Date().toLocaleString()
      }

      ElMessage.success(reviewType.value === 'approve' ? '项目审核通过' : '项目已驳回')
      reviewVisible.value = false
      currentProject.value = null
      
      // 重新加载项目列表
      loadProjects()
    } else {
      ElMessage.error(response.message || '审核操作失败')
    }
  } catch (error) {
    console.error('审核操作失败:', error)
    ElMessage.error(error.message || '审核操作失败')
  } finally {
    submitting.value = false
  }
}

const updateProgress = (project) => {
  progressForm.value.progress = project.progress || 0;
  progressForm.value.description = project.progressDescription || '';
  progressForm.value.expectedEndDate = project.expectedEndDate ? new Date(project.expectedEndDate) : null;
  progressVisible.value = true;
  currentProject.value = project; // 确保当前项目是正在更新的项目
};

const submitProgress = async () => {
  try {
    await progressFormRef.value.validate();
  } catch (error) {
    return;
  }

  submitting.value = true;
  try {
    const progressData = {
      progress: progressForm.value.progress,
      progressDescription: progressForm.value.description,
      expectedEndDate: progressForm.value.expectedEndDate ? progressForm.value.expectedEndDate.toISOString() : null,
      comment: progressForm.value.description // 进度更新时，审核意见即为进度说明
    };

    const response = await teacherService.updateProjectProgress(currentProject.value.id, progressData);

    if (response.code === 200) {
      const project = projects.value.find(p => p.id === currentProject.value.id);
      if (project) {
        project.progress = progressForm.value.progress;
        project.progressDescription = progressForm.value.description;
        project.expectedEndDate = progressForm.value.expectedEndDate ? progressForm.value.expectedEndDate.toISOString() : null;
        project.updateTime = new Date().toLocaleString();
      }
      ElMessage.success('项目进度更新成功');
      progressVisible.value = false;
      currentProject.value = null;
      loadProjects(); // 重新加载项目列表以更新进度
    } else {
      ElMessage.error(response.message || '更新项目进度失败');
    }
  } catch (error) {
    console.error('更新项目进度失败:', error);
    ElMessage.error(error.message || '更新项目进度失败');
  } finally {
    submitting.value = false;
  }
};

const assessQuality = (project) => {
  currentProject.value = project;
  qualityAssessmentVisible.value = true;
  qualityForm.value = {
    innovation: project.innovationRating || 0,
    feasibility: project.feasibilityRating || 0,
    completion: project.completionRating || 0,
    collaboration: project.collaborationRating || 0,
    overallComment: project.overallComment || '',
    suggestions: project.suggestions || ''
  };
};

const submitQualityAssessment = async () => {
  try {
    await qualityFormRef.value.validate();
  } catch (error) {
    return;
  }

  submitting.value = true;
  try {
    const qualityData = {
      innovationRating: qualityForm.value.innovation,
      feasibilityRating: qualityForm.value.feasibility,
      completionRating: qualityForm.value.completion,
      collaborationRating: qualityForm.value.collaboration,
      overallComment: qualityForm.value.overallComment,
      suggestions: qualityForm.value.suggestions
    };

    const response = await teacherService.assessProjectQuality(currentProject.value.id, qualityData);

    if (response.code === 200) {
      const project = projects.value.find(p => p.id === currentProject.value.id);
      if (project) {
        project.innovationRating = qualityForm.value.innovation;
        project.feasibilityRating = qualityForm.value.feasibility;
        project.completionRating = qualityForm.value.completion;
        project.collaborationRating = qualityForm.value.collaboration;
        project.overallComment = qualityForm.value.overallComment;
        project.suggestions = qualityForm.value.suggestions;
        project.updateTime = new Date().toLocaleString();
      }
      ElMessage.success('项目质量评估成功');
      qualityAssessmentVisible.value = false;
      currentProject.value = null;
      loadProjects(); // 重新加载项目列表以更新评估
    } else {
      ElMessage.error(response.message || '评估失败');
    }
  } catch (error) {
    console.error('评估失败:', error);
    ElMessage.error(error.message || '评估失败');
  } finally {
    submitting.value = false;
  }
};

const downloadFile = async (file) => {
  try {
    await fileService.downloadFile(file.fileUrl, file.fileName)
    ElMessage.success('文件下载成功')
  } catch (error) {
    console.error('文件下载失败:', error)
    ElMessage.error(error.message || '文件下载失败')
  }
}

const getStatusType = (status) => {
  const statusMap = {
    draft: 'info',
    submitted: 'info', // 新增
    pending: 'warning',
    approved: 'success',
    rejected: 'danger',
    in_progress: 'primary', // 新增
    completed: 'success', // 新增
    suspended: 'info', // 新增
    need_revision: 'danger' // 新增
  }
  return statusMap[status] || 'info'
}

const getStatusText = (status) => {
  const statusMap = {
    draft: '草稿',
    submitted: '已提交', // 新增
    pending: '待审核',
    approved: '已通过',
    rejected: '已拒绝',
    in_progress: '进行中', // 新增
    completed: '已完成', // 新增
    suspended: '已暂停', // 新增
    need_revision: '需修改' // 新增
  }
  return statusMap[status] || status
}

const getLevelType = (level) => {
  const levelMap = {
    '校级': 'info',
    '省级': 'success',
    '国家级': 'warning',
    '国际级': 'danger'
  }
  return levelMap[level] || 'info'
}

const getProgressStatus = (progress) => {
  if (progress >= 90) return 'success';
  if (progress >= 70) return 'exception';
  if (progress >= 50) return 'warning';
  return 'info';
};

const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN')
}

const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 Bytes';
  const k = 1024;
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
};

const disabledDate = (time) => {
  return time.getTime() < Date.now() - 8.64e7; // 不能选择过去的时间
};

const getRowClassName = ({ row }) => {
  if (row.status === 'pending' || row.status === 'submitted') {
    return 'pending-row';
  }
  if (row.status === 'rejected') {
    return 'rejected-row';
  }
  if (row.isExtended) {
    return 'extended-row';
  }
  return '';
};

const batchApprove = () => {
  if (selectedProjects.value.length === 0) {
    ElMessage.warning('请至少选择一个项目进行批量通过');
    return;
  }
  batchForm.value.operationType = 'approve';
  batchForm.value.rejectReason = ''; // 清空驳回原因
  batchOperationVisible.value = true;
};

const batchReview = () => {
  if (selectedProjects.value.length === 0) {
    ElMessage.warning('请至少选择一个项目进行批量审核');
    return;
  }
  batchForm.value.operationType = 'reject';
  batchForm.value.content = ''; // 清空评语
  batchOperationVisible.value = true;
};

const submitBatchOperation = async () => {
  try {
    await batchFormRef.value.validate();
  } catch (error) {
    return;
  }

  submitting.value = true;
  try {
    const operationData = {
      projectIds: selectedProjects.value.map(p => p.id),
      operationType: batchForm.value.operationType,
      content: batchForm.value.operationType === 'addComment' ? batchForm.value.content : undefined,
      progress: batchForm.value.operationType === 'updateProgress' ? batchForm.value.progress : undefined,
      rejectReason: batchForm.value.operationType === 'reject' ? batchForm.value.rejectReason : undefined
    };

    let response;
    if (batchForm.value.operationType === 'approve') {
      response = await teacherService.batchApproveProjects(operationData);
    } else if (batchForm.value.operationType === 'reject') {
      response = await teacherService.batchRejectProjects(operationData);
    } else if (batchForm.value.operationType === 'updateProgress') {
      response = await teacherService.batchUpdateProjectProgress(operationData);
    } else if (batchForm.value.operationType === 'addComment') {
      response = await teacherService.batchAddComments(operationData);
    }

    if (response.code === 200) {
      ElMessage.success('批量操作成功');
      batchOperationVisible.value = false;
      selectedProjects.value = []; // 清空选中项目
      loadProjects(); // 重新加载项目列表
    } else {
      ElMessage.error(response.message || '批量操作失败');
    }
  } catch (error) {
    console.error('批量操作失败:', error);
    ElMessage.error(error.message || '批量操作失败');
  } finally {
    submitting.value = false;
  }
};

const exportProjects = () => {
  ElMessageBox.confirm('确定要导出选中的项目数据吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(async () => {
    try {
      const projectIds = selectedProjects.value.map(p => p.id);
      const response = await teacherService.exportProjects(projectIds);
      if (response.code === 200) {
        const blob = new Blob([response.data], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' });
        const fileName = `项目数据_${new Date().toISOString().slice(0, 10)}.xlsx`;
        const link = document.createElement('a');
        link.href = window.URL.createObjectURL(blob);
        link.download = fileName;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        ElMessage.success('项目数据导出成功');
      } else {
        ElMessage.error(response.message || '导出项目数据失败');
      }
    } catch (error) {
      console.error('导出项目数据失败:', error);
      ElMessage.error(error.message || '导出项目数据失败');
    }
  }).catch(() => {
    // 用户取消
  });
};

const generateReport = () => {
  ElMessageBox.confirm('确定要生成项目报告吗？', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(async () => {
    try {
      const projectIds = selectedProjects.value.map(p => p.id);
      const response = await teacherService.generateProjectReport(projectIds);
      if (response.code === 200) {
        const blob = new Blob([response.data], { type: 'application/pdf' });
        const fileName = `项目报告_${new Date().toISOString().slice(0, 10)}.pdf`;
        const link = document.createElement('a');
        link.href = window.URL.createObjectURL(blob);
        link.download = fileName;
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
        ElMessage.success('项目报告生成成功');
      } else {
        ElMessage.error(response.message || '生成项目报告失败');
      }
    } catch (error) {
      console.error('生成项目报告失败:', error);
      ElMessage.error(error.message || '生成项目报告失败');
    }
  }).catch(() => {
    // 用户取消
  });
};

// 组件挂载时加载数据
onMounted(() => {
  loadProjects()
  loadProjectStats()
})

// 新增：加载项目统计
const loadProjectStats = async () => {
  statsLoading.value = true
  try {
    const response = await teacherService.getProjectStats()
    if (response.code === 200) {
      // 这里可以处理统计数据，用于图表显示
      console.log('项目统计数据:', response.data)
    }
  } catch (error) {
    console.error('加载项目统计失败:', error)
  } finally {
    statsLoading.value = false
  }
}

// 新增：刷新统计
const refreshStats = () => {
  loadProjectStats()
}

// 新增：初始化图表
const initCharts = () => {
  // 这里可以初始化ECharts图表
  // 由于需要引入ECharts库，这里先提供框架
  console.log('初始化图表')
}

// 新增：更新图表数据
const updateCharts = () => {
  // 更新图表数据
  console.log('更新图表数据')
}
</script>

<style scoped>
.project-management {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 20px;
  background-color: #f5f7fa;
  border-radius: 8px;
}

.header-content {
  flex-grow: 1;
  margin-right: 20px;
}

.page-title {
  font-size: 24px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 5px;
}

.page-subtitle {
  font-size: 16px;
  color: #606266;
}

.header-stats {
  display: flex;
  gap: 20px;
}

.stat-item {
  text-align: center;
}

.stat-number {
  font-size: 24px;
  font-weight: bold;
  color: #409eff;
  margin-bottom: 5px;
}

.stat-label {
  font-size: 14px;
  color: #909399;
}

.filter-card {
  margin-bottom: 20px;
}

.search-input {
  width: 100%;
}

.filter-select {
  width: 100%;
}

.refresh-btn, .reset-btn {
  margin-right: 10px;
}

.project-list {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-title {
  font-size: 18px;
  font-weight: bold;
  color: #303133;
}

.header-actions {
  display: flex;
  align-items: center;
}

.project-count {
  font-size: 14px;
  color: #909399;
}

.project-detail-dialog .el-dialog__header {
  background-color: #f5f7fa;
  border-bottom: 1px solid #ebeef5;
}

.project-detail-dialog .el-dialog__body {
  padding: 20px;
  max-height: 60vh;
  overflow-y: auto;
}

.project-detail {
  max-height: 60vh;
  overflow-y: auto;
}

.detail-section {
  margin-bottom: 20px;
}

.section-header {
  display: flex;
  align-items: center;
  font-size: 16px;
  font-weight: bold;
  color: #303133;
  margin-bottom: 15px;
}

.section-header .el-icon {
  margin-right: 8px;
  font-size: 18px;
}

.project-title-large {
  font-size: 18px;
  font-weight: bold;
  color: #303133;
}

.project-title {
  display: flex;
  align-items: center;
}

.title-text {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  margin-right: 5px;
}

.extension-tag {
  margin-left: 10px;
}

.student-info {
  display: flex;
  align-items: center;
}

.student-avatar {
  margin-right: 8px;
}

.student-name {
  font-size: 14px;
  color: #606266;
}

.progress-container {
  display: flex;
  align-items: center;
  margin-top: 10px;
}

.progress-bar {
  flex-grow: 1;
  margin-right: 10px;
}

.progress-text {
  font-size: 14px;
  color: #909399;
}

.date-text {
  font-size: 14px;
  color: #606266;
}

.no-data {
  color: #909399;
}

.action-buttons {
  display: flex;
  gap: 10px;
}

.view-btn, .approve-btn, .reject-btn, .progress-btn, .quality-btn {
  flex-shrink: 0;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.review-dialog .el-dialog__header {
  background-color: #f5f7fa;
  border-bottom: 1px solid #ebeef5;
}

.review-dialog .el-dialog__body {
  padding: 20px;
}

.review-dialog .el-form-item {
  margin-bottom: 15px;
}

.review-dialog .el-form-item__label {
  font-weight: bold;
}

.review-dialog .el-textarea {
  margin-top: 5px;
}

.review-dialog .el-input__inner {
  border-radius: 4px;
}

.review-dialog .el-button {
  border-radius: 4px;
}

.progress-dialog .el-dialog__header {
  background-color: #f5f7fa;
  border-bottom: 1px solid #ebeef5;
}

.progress-dialog .el-dialog__body {
  padding: 20px;
}

.progress-dialog .el-form-item {
  margin-bottom: 15px;
}

.progress-dialog .el-form-item__label {
  font-weight: bold;
}

.progress-dialog .el-slider {
  margin-top: 5px;
}

.progress-dialog .el-input__inner {
  border-radius: 4px;
}

.progress-dialog .el-button {
  border-radius: 4px;
}

.review-dialog .dialog-footer, .progress-dialog .dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.timeline-card {
  margin-bottom: 10px;
}

.review-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 5px;
}

.reviewer-name {
  font-size: 14px;
  font-weight: bold;
  color: #303133;
}

.review-comments {
  font-size: 14px;
  color: #606266;
  line-height: 1.6;
}

.description-content {
  font-size: 14px;
  color: #606266;
  line-height: 1.8;
}

/* 行样式 */
.pending-row {
  background-color: #fdf6ec !important;
}

.rejected-row {
  background-color: #fef0f0 !important;
}

.extended-row {
  background-color: #f0f9ff !important;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
  }
  
  .header-stats {
    margin-top: 20px;
    width: 100%;
    justify-content: space-between;
  }
  
  .filter-card .el-row {
    margin: 0;
  }
  
  .filter-card .el-col {
    margin-bottom: 15px;
  }
  
  .action-buttons {
    flex-direction: column;
    gap: 5px;
  }
  
  .action-buttons .el-button {
    width: 100%;
  }
}

/* 新增样式 */
.stats-chart-card {
  margin-bottom: 20px;
}

.quick-actions-card {
  margin-bottom: 20px;
}

.chart-container {
  text-align: center;
  padding: 20px;
}

.chart-container h4 {
  margin-bottom: 15px;
  color: #303133;
  font-size: 16px;
}

.pie-chart {
  height: 200px;
  background-color: #f5f7fa;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #909399;
}

.line-chart {
  height: 250px;
  background-color: #f5f7fa;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #909399;
}

.quick-action-btn {
  width: 100%;
  height: 60px;
  font-size: 14px;
}

.batch-operation-dialog .el-dialog__header {
  background-color: #f5f7fa;
  border-bottom: 1px solid #ebeef5;
}

.batch-operation-dialog .el-dialog__body {
  padding: 20px;
}

.batch-info {
  margin-bottom: 20px;
}

.quality-assessment-dialog .el-dialog__header {
  background-color: #f5f7fa;
  border-bottom: 1px solid #ebeef5;
}

.quality-assessment-dialog .el-dialog__body {
  padding: 20px;
}

.quality-assessment .el-form-item {
  margin-bottom: 20px;
}

.quality-assessment .el-rate {
  margin-top: 5px;
}

.quality-btn {
  margin-left: 5px;
}
</style> 
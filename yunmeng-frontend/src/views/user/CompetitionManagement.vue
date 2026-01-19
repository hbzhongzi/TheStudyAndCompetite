<template>
  <div class="competition-management">
    <div class="page-header">
      <h2>竞赛管理</h2>
      <el-button type="primary" @click="showAddCompetitionDialog = true">
        <el-icon><Plus /></el-icon>
        发布竞赛
      </el-button>
    </div>

    <!-- 搜索和筛选 -->
    <el-card class="search-card">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-input
            v-model="searchQuery"
            placeholder="搜索竞赛名称或主办方"
            clearable
            @input="handleSearch"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
        </el-col>
        <el-col :span="4">
          <el-select v-model="statusFilter" placeholder="状态筛选" clearable @change="handleSearch">
            <el-option label="全部状态" value="" />
            <el-option label="草稿" value="draft" />
            <el-option label="报名中" value="registration" />
            <el-option label="提交中" value="submission" />
            <el-option label="评审中" value="review" />
            <el-option label="已完成" value="completed" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-select v-model="typeFilter" placeholder="类型筛选" clearable @change="handleSearch">
            <el-option label="全部类型" value="" />
            <el-option label="校级" value="校级" />
            <el-option label="省级" value="省级" />
            <el-option label="国家级" value="国家级" />
            <el-option label="国际级" value="国际级" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-select v-model="isOpenFilter" placeholder="开放状态" clearable @change="handleSearch">
            <el-option label="全部" value="" />
            <el-option label="开放报名" :value="true" />
            <el-option label="关闭报名" :value="false" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
        </el-col>
      </el-row>
      
      <!-- 高级筛选 -->
      <el-row :gutter="20" style="margin-top: 15px;">
        <el-col :span="4">
          <el-select v-model="teacherLimitFilter" placeholder="教师指导" clearable @change="handleSearch">
            <el-option label="全部" value="" />
            <el-option label="需要教师指导" :value="true" />
            <el-option label="不需要教师指导" :value="false" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-select v-model="locationFilter" placeholder="地点筛选" clearable @change="handleSearch">
            <el-option label="全部地点" value="" />
            <el-option label="线上" value="线上" />
            <el-option label="线下" value="线下" />
            <el-option label="混合" value="混合" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-date-picker
            v-model="dateRangeFilter"
            type="datetimerange"
            range-separator="至"
            start-placeholder="开始时间"
            end-placeholder="结束时间"
            format="YYYY-MM-DD HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
            @change="handleSearch"
            style="width: 100%"
          />
        </el-col>
        <el-col :span="4">
          <el-button @click="clearFilters" type="info">
            <el-icon><Refresh /></el-icon>
            清除筛选
          </el-button>
        </el-col>
      </el-row>
    </el-card>

    <!-- 竞赛列表 -->
    <el-card class="competition-list-card">
      <el-table :data="safeCompetitions" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="竞赛ID" width="80" />
        <el-table-column prop="title" label="竞赛名称" width="250" />
        <el-table-column prop="organizer" label="主办方" width="150" />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="scope">
            <el-tag :type="getTypeTagType(scope.row.type)">
              {{ scope.row.type }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getStatusTagType(scope.row.status)">
              {{ getStatusLabel(scope.row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="is_open" label="报名状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.is_open ? 'success' : 'danger'">
              {{ scope.row.is_open ? '开放' : '关闭' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="start_time" label="比赛开始时间" width="120">
          <template #default="scope">
            {{ formatDate(scope.row.start_time) }}
          </template>
        </el-table-column>
        <el-table-column prop="end_time" label="比赛结束时间" width="120">
          <template #default="scope">
            {{ formatDate(scope.row.end_time) }}
          </template>
        </el-table-column>
        <el-table-column prop="registration_start" label="报名开始时间" width="120">
          <template #default="scope">
            {{ formatDate(scope.row.registration_start) || '未设置' }}
          </template>
        </el-table-column>
        <el-table-column prop="registration_end" label="报名截止时间" width="120">
          <template #default="scope">
            {{ formatDate(scope.row.registration_end) || '未设置' }}
          </template>
        </el-table-column>
        <el-table-column prop="submission_deadline" label="作品提交截止" width="120">
          <template #default="scope">
            {{ formatDate(scope.row.submission_deadline) || '未设置' }}
          </template>
        </el-table-column>
        <el-table-column prop="location" label="竞赛地点" width="150">
          <template #default="scope">
            {{ scope.row.location || '未设置' }}
          </template>
        </el-table-column>
        <el-table-column prop="contact" label="联系方式" width="120">
          <template #default="scope">
            {{ scope.row.contact || '未设置' }}
          </template>
        </el-table-column>
        <el-table-column prop="department_limit" label="院系限制" width="120">
          <template #default="scope">
            {{ scope.row.department_limit || '不限' }}
          </template>
        </el-table-column>
        <el-table-column prop="teacher_limit" label="教师指导" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.teacher_limit ? 'warning' : 'info'" size="small">
              {{ scope.row.teacher_limit ? '需要' : '不需要' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="current_participants" label="参赛人数" width="100" />
        <el-table-column prop="max_participants" label="最大人数" width="100" />
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="scope">
            <div class="action-buttons">
              <!-- 主要操作按钮组 -->
              <el-button-group class="main-actions">
                <el-button 
                  size="small" 
                  type="primary" 
                  @click="editCompetition(scope.row)"
                  :icon="Edit"
                  title="编辑竞赛"
                >
                  编辑
                </el-button>
                <el-button 
                  size="small" 
                  type="info" 
                  @click="manageRegistrations(scope.row)"
                  :icon="User"
                  title="报名管理"
                >
                  报名
                </el-button>
                <el-button 
                  size="small" 
                  type="success" 
                  @click="manageSubmissions(scope.row)"
                  :icon="Document"
                  title="作品管理"
                >
                  作品
                </el-button>
              </el-button-group>
              
              <!-- 次要操作下拉菜单 -->
              <el-dropdown @command="handleCommand" trigger="click" class="secondary-actions">
                <el-button size="small" type="default">
                  更多<el-icon class="el-icon--right"><ArrowDown /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <!-- 基础操作 -->
                    <el-dropdown-item :command="{action: 'view', id: scope.row.id}">
                      <el-icon><View /></el-icon>查看详情
                    </el-dropdown-item>
                    <el-dropdown-item :command="{action: 'toggleStatus', id: scope.row.id}">
                      <el-icon><Switch /></el-icon>{{ scope.row.is_open ? '关闭报名' : '开放报名' }}
                    </el-dropdown-item>
                    
                    <!-- 评审相关 -->
                    <el-dropdown-item :command="{action: 'assignJudges', id: scope.row.id}">
                      <el-icon><UserFilled /></el-icon>分配评审
                    </el-dropdown-item>
                    <el-dropdown-item :command="{action: 'judgingProgress', id: scope.row.id}">
                      <el-icon><TrendCharts /></el-icon>评审进度
                    </el-dropdown-item>
                    
                    <!-- 数据管理 -->
                    <el-dropdown-item :command="{action: 'exportRegistrations', id: scope.row.id}">
                      <el-icon><Download /></el-icon>导出报名
                    </el-dropdown-item>
                    <el-dropdown-item :command="{action: 'exportResults', id: scope.row.id}">
                      <el-icon><Download /></el-icon>导出结果
                    </el-dropdown-item>
                    
                    <!-- 成绩管理 -->
                    <el-dropdown-item :command="{action: 'manageResults', id: scope.row.id}">
                      <el-icon><Trophy /></el-icon>成绩管理
                    </el-dropdown-item>
                    
                    <!-- 高级操作 -->
                    <el-dropdown-item :command="{action: 'finalize', id: scope.row.id}" divided>
                      <el-icon><Check /></el-icon>最终确认
                    </el-dropdown-item>
                    <el-dropdown-item :command="{action: 'delete', id: scope.row.id}" divided>
                      <el-icon><Delete /></el-icon>删除竞赛
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handleCurrentChange"
        />
      </div>
    </el-card>

    <!-- 创建/编辑竞赛对话框 -->
    <el-dialog
      v-model="showAddCompetitionDialog"
      :title="editingCompetition ? '编辑竞赛' : '发布竞赛'"
      width="60%"
    >
      <el-form
        ref="competitionFormRef"
        :model="competitionForm"
        :rules="competitionRules"
        label-width="120px"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="竞赛名称" prop="title">
              <el-input v-model="competitionForm.title" placeholder="请输入竞赛名称" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="竞赛类型" prop="type">
              <el-select v-model="competitionForm.type" placeholder="请选择竞赛类型">
                <el-option label="校级" value="校级" />
                <el-option label="省级" value="省级" />
                <el-option label="国家级" value="国家级" />
                <el-option label="国际级" value="国际级" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="主办方" prop="organizer">
              <el-input v-model="competitionForm.organizer" placeholder="请输入主办方" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="最大参与人数" prop="max_participants">
              <el-input-number v-model="competitionForm.max_participants" :min="1" :max="10000" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="报名开始时间" prop="registration_start">
              <el-date-picker
                v-model="competitionForm.registration_start"
                type="datetime"
                placeholder="选择报名开始时间"
                format="YYYY-MM-DD HH:mm:ss"
                value-format="YYYY-MM-DD HH:mm:ss"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="报名截止时间" prop="registration_end">
              <el-date-picker
                v-model="competitionForm.registration_end"
                type="datetime"
                placeholder="选择报名截止时间"
                format="YYYY-MM-DD HH:mm:ss"
                value-format="YYYY-MM-DD HH:mm:ss"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="比赛开始时间" prop="start_time">
              <el-date-picker
                v-model="competitionForm.start_time"
                type="datetime"
                placeholder="选择比赛开始时间"
                format="YYYY-MM-DD HH:mm:ss"
                value-format="YYYY-MM-DD HH:mm:ss"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="比赛结束时间" prop="end_time">
              <el-date-picker
                v-model="competitionForm.end_time"
                type="datetime"
                placeholder="选择比赛结束时间"
                format="YYYY-MM-DD HH:mm:ss"
                value-format="YYYY-MM-DD HH:mm:ss"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="作品提交截止" prop="submission_deadline">
              <el-date-picker
                v-model="competitionForm.submission_deadline"
                type="datetime"
                placeholder="选择作品提交截止时间"
                format="YYYY-MM-DD HH:mm:ss"
                value-format="YYYY-MM-DD HH:mm:ss"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="竞赛地点" prop="location">
              <el-input v-model="competitionForm.location" placeholder="请输入竞赛地点" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="联系方式" prop="contact">
              <el-input v-model="competitionForm.contact" placeholder="请输入联系方式" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="院系限制" prop="department_limit">
              <el-input v-model="competitionForm.department_limit" placeholder="请输入院系限制，留空表示不限" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="教师指导" prop="teacher_limit">
              <el-switch 
                v-model="competitionForm.teacher_limit" 
                active-text="需要教师指导"
                inactive-text="不需要教师指导"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="竞赛描述" prop="description">
          <el-input
            v-model="competitionForm.description"
            type="textarea"
            :rows="4"
            placeholder="请输入竞赛描述"
          />
        </el-form-item>

        <el-form-item label="竞赛规则" prop="rules">
          <el-input
            v-model="competitionForm.rules"
            type="textarea"
            :rows="3"
            placeholder="请输入竞赛规则"
          />
        </el-form-item>

        <el-form-item label="报名要求" prop="requirements">
          <el-input
            v-model="competitionForm.requirements"
            type="textarea"
            :rows="3"
            placeholder="请输入报名要求"
          />
        </el-form-item>

        <el-form-item label="评审方式" prop="judging_method">
          <el-input
            v-model="competitionForm.judging_method"
            type="textarea"
            :rows="3"
            placeholder="请输入评审方式"
          />
        </el-form-item>

        <el-form-item label="重要注意事项" prop="important_notes">
          <el-input
            v-model="competitionForm.important_notes"
            type="textarea"
            :rows="3"
            placeholder="请输入重要注意事项"
          />
        </el-form-item>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="相关链接" prop="website">
              <el-input v-model="competitionForm.website" placeholder="请输入相关链接" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="QQ群" prop="qq_group">
              <el-input v-model="competitionForm.qq_group" placeholder="请输入QQ群号" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="支持文件格式" prop="file_formats">
              <el-input v-model="competitionForm.file_formats" placeholder="如：PDF、PPT、DOC、ZIP" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="文件大小限制" prop="file_size_limit">
              <el-input v-model="competitionForm.file_size_limit" placeholder="如：单个文件不超过50MB" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="附件URL" prop="attachment">
          <el-input v-model="competitionForm.attachment" placeholder="请输入附件URL" />
        </el-form-item>

        <el-form-item label="开放报名" prop="is_open">
          <el-switch v-model="competitionForm.is_open" />
        </el-form-item>

        <el-form-item label="竞赛状态" prop="status">
          <el-select v-model="competitionForm.status" placeholder="请选择竞赛状态">
            <el-option label="草稿" value="draft" />
            <el-option label="报名中" value="registration" />
            <el-option label="提交中" value="submission" />
            <el-option label="评审中" value="review" />
            <el-option label="已完成" value="completed" />
          </el-select>
        </el-form-item>
      </el-form>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="closeCompetitionDialog">取消</el-button>
          <el-button type="primary" @click="submitCompetition" :loading="submitting">
            {{ editingCompetition ? '更新' : '创建' }}
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 报名管理对话框 -->
    <el-dialog
      v-model="showRegistrationsDialog"
      title="报名管理"
      width="80%"
    >
      <div v-if="currentCompetition">
        <h3>{{ currentCompetition.title }} - 报名记录</h3>
        <el-table :data="registrations" style="width: 100%" v-loading="registrationsLoading">
          <el-table-column prop="id" label="报名ID" width="80" />
          <el-table-column prop="student.username" label="学生姓名" width="120" />
          <el-table-column prop="student.email" label="邮箱" width="180" />
          <el-table-column prop="team_name" label="团队名称" width="150" />
          <el-table-column prop="status" label="状态" width="100">
            <template #default="scope">
              <el-tag :type="getRegistrationStatusType(scope.row.status)">
                {{ getRegistrationStatusLabel(scope.row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="register_time" label="报名时间" width="150">
            <template #default="scope">
              {{ formatDate(scope.row.register_time) }}
            </template>
          </el-table-column>
          <el-table-column prop="contact_phone" label="联系电话" width="120" />
          <el-table-column prop="contact_email" label="联系邮箱" width="180" />
        </el-table>
      </div>
    </el-dialog>

    <!-- 作品管理对话框 -->
    <el-dialog
      v-model="showSubmissionsDialog"
      title="作品管理"
      width="80%"
    >
      <div v-if="currentCompetition">
        <h3>{{ currentCompetition.title }} - 提交作品</h3>
        <el-table :data="submissions" style="width: 100%" v-loading="submissionsLoading">
          <el-table-column prop="id" label="作品ID" width="80" />
          <el-table-column prop="student.username" label="学生姓名" width="120" />
          <el-table-column prop="file_name" label="文件名" width="200" />
          <el-table-column prop="file_size" label="文件大小" width="100">
            <template #default="scope">
              {{ formatFileSize(scope.row.file_size) }}
            </template>
          </el-table-column>
          <el-table-column prop="status" label="状态" width="100">
            <template #default="scope">
              <el-tag :type="getSubmissionStatusType(scope.row.status)">
                {{ getSubmissionStatusLabel(scope.row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="submit_time" label="提交时间" width="150">
            <template #default="scope">
              {{ formatDate(scope.row.submit_time) }}
            </template>
          </el-table-column>
          <el-table-column prop="description" label="描述" />
          <el-table-column label="操作" width="150">
            <template #default="scope">
              <el-button size="small" @click="downloadSubmission(scope.row)">下载</el-button>
              <el-button size="small" type="primary" @click="viewSubmission(scope.row)">查看</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-dialog>

    <!-- 竞赛详情查看对话框 -->
    <el-dialog
      v-model="showCompetitionDetailDialog"
      title="竞赛详情"
      width="70%"
    >
      <div v-if="currentCompetition" class="competition-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="竞赛名称" :span="2">
            {{ currentCompetition.title }}
          </el-descriptions-item>
          <el-descriptions-item label="竞赛类型">
            <el-tag :type="getTypeTagType(currentCompetition.type)">
              {{ currentCompetition.type }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="主办方">
            {{ currentCompetition.organizer }}
          </el-descriptions-item>
          <el-descriptions-item label="竞赛状态">
            <el-tag :type="getStatusTagType(currentCompetition.status)">
              {{ getStatusLabel(currentCompetition.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="报名状态">
            <el-tag :type="currentCompetition.is_open ? 'success' : 'danger'">
              {{ currentCompetition.is_open ? '开放' : '关闭' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="报名开始时间">
            {{ formatDate(currentCompetition.registration_start) || '未设置' }}
          </el-descriptions-item>
          <el-descriptions-item label="报名截止时间">
            {{ formatDate(currentCompetition.registration_end) || '未设置' }}
          </el-descriptions-item>
          <el-descriptions-item label="作品提交截止">
            {{ formatDate(currentCompetition.submission_deadline) || '未设置' }}
          </el-descriptions-item>
          <el-descriptions-item label="比赛开始时间">
            {{ formatDate(currentCompetition.start_time) }}
          </el-descriptions-item>
          <el-descriptions-item label="比赛结束时间">
            {{ formatDate(currentCompetition.end_time) }}
          </el-descriptions-item>
          <el-descriptions-item label="竞赛地点">
            {{ currentCompetition.location || '未设置' }}
          </el-descriptions-item>
          <el-descriptions-item label="联系方式">
            {{ currentCompetition.contact || '未设置' }}
          </el-descriptions-item>
          <el-descriptions-item label="院系限制">
            {{ currentCompetition.department_limit || '不限' }}
          </el-descriptions-item>
          <el-descriptions-item label="教师指导">
            <el-tag :type="currentCompetition.teacher_limit ? 'warning' : 'info'">
              {{ currentCompetition.teacher_limit ? '需要' : '不需要' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="最大参与人数">
            {{ currentCompetition.max_participants }}
          </el-descriptions-item>
          <el-descriptions-item label="当前参与人数">
            {{ currentCompetition.current_participants }}
          </el-descriptions-item>
          <el-descriptions-item label="支持文件格式" :span="2">
            {{ currentCompetition.file_formats || '未设置' }}
          </el-descriptions-item>
          <el-descriptions-item label="文件大小限制" :span="2">
            {{ currentCompetition.file_size_limit || '未设置' }}
          </el-descriptions-item>
          <el-descriptions-item label="竞赛描述" :span="2">
            {{ currentCompetition.description }}
          </el-descriptions-item>
          <el-descriptions-item label="竞赛规则" :span="2">
            {{ currentCompetition.rules || '未设置' }}
          </el-descriptions-item>
          <el-descriptions-item label="报名要求" :span="2">
            {{ currentCompetition.requirements || '未设置' }}
          </el-descriptions-item>
          <el-descriptions-item label="评审方式" :span="2">
            {{ currentCompetition.judging_method || '未设置' }}
          </el-descriptions-item>
          <el-descriptions-item label="重要注意事项" :span="2">
            {{ currentCompetition.important_notes || '未设置' }}
          </el-descriptions-item>
          <el-descriptions-item label="相关链接" :span="2">
            <el-link v-if="currentCompetition.website" :href="currentCompetition.website" target="_blank" type="primary">
              {{ currentCompetition.website }}
            </el-link>
            <span v-else>未设置</span>
          </el-descriptions-item>
          <el-descriptions-item label="QQ群" :span="2">
            {{ currentCompetition.qq_group || '未设置' }}
          </el-descriptions-item>
          <el-descriptions-item label="附件" :span="2">
            <el-link v-if="currentCompetition.attachment" :href="currentCompetition.attachment" target="_blank" type="primary">
              查看附件
            </el-link>
            <span v-else>无附件</span>
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </el-dialog>

    <!-- 成绩管理对话框 -->
    <el-dialog
      v-model="showResultsDialog"
      title="成绩管理"
      width="80%"
    >
      <div v-if="currentCompetition">
        <div class="results-header">
          <h3>{{ currentCompetition.title }} - 获奖结果</h3>
          <el-button type="primary" @click="showAddResultDialog = true">
            <el-icon><Plus /></el-icon>
            添加获奖结果
          </el-button>
        </div>
        <el-table :data="results" style="width: 100%" v-loading="resultsLoading">
          <el-table-column prop="id" label="结果ID" width="80" />
          <el-table-column prop="student.username" label="学生姓名" width="120" />
          <el-table-column prop="award_level" label="获奖等级" width="100">
            <template #default="scope">
              <el-tag :type="getAwardLevelType(scope.row.award_level)">
                {{ scope.row.award_level }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="final_score" label="最终得分" width="100" />
          <el-table-column prop="publish_time" label="公布时间" width="150">
            <template #default="scope">
              {{ formatDate(scope.row.publish_time) }}
            </template>
          </el-table-column>
          <el-table-column prop="certificate_url" label="证书" width="100">
            <template #default="scope">
              <el-button size="small" v-if="scope.row.certificate_url" @click="downloadCertificate(scope.row)">
                下载证书
              </el-button>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="150">
            <template #default="scope">
              <el-button size="small" type="primary" @click="editResult(scope.row)">编辑</el-button>
              <el-button size="small" type="danger" @click="deleteResult(scope.row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-dialog>

    <!-- 添加获奖结果对话框 -->
    <el-dialog
      v-model="showAddResultDialog"
      title="添加获奖结果"
      width="50%"
    >
      <el-form
        ref="resultFormRef"
        :model="resultForm"
        :rules="resultRules"
        label-width="120px"
      >
        <el-form-item label="学生" prop="student_id">
          <el-select v-model="resultForm.student_id" placeholder="请选择学生" filterable>
            <el-option
              v-for="student in availableStudents"
              :key="student.id"
              :label="student.username"
              :value="student.id"
            />
          </el-select>
        </el-form-item>
        
        <el-form-item label="获奖等级" prop="award_level">
          <el-select v-model="resultForm.award_level" placeholder="请选择获奖等级">
            <el-option label="一等奖" value="一等奖" />
            <el-option label="二等奖" value="二等奖" />
            <el-option label="三等奖" value="三等奖" />
            <el-option label="优秀奖" value="优秀奖" />
            <el-option label="参与奖" value="参与奖" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="最终得分" prop="final_score">
          <el-input-number v-model="resultForm.final_score" :min="0" :max="100" />
        </el-form-item>
        
        <el-form-item label="证书URL" prop="certificate_url">
          <el-input v-model="resultForm.certificate_url" placeholder="请输入证书URL" />
        </el-form-item>
      </el-form>

      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showAddResultDialog = false">取消</el-button>
          <el-button type="primary" @click="submitResult" :loading="submittingResult">
            添加
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script>
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { 
  Plus, 
  Search, 
  ArrowDown, 
  Edit, 
  User, 
  Document, 
  View, 
  Switch, 
  UserFilled, 
  TrendCharts, 
  Download, 
  Trophy, 
  Check, 
  Delete,
  Refresh
} from '@element-plus/icons-vue'
import competitionService from '@/services/competitionService'

export default {
  name: 'CompetitionManagement',
  components: {
    Plus,
    Search,
    ArrowDown,
    Edit,
    User,
    Document,
    View,
    Switch,
    UserFilled,
    TrendCharts,
    Download,
    Trophy,
    Check,
    Delete,
    Refresh
  },
  setup() {
    // 响应式数据
    const loading = ref(false)
    const competitions = ref([])
    
    // 确保竞赛列表始终是一个数组
    const safeCompetitions = computed(() => {
      return Array.isArray(competitions.value) ? competitions.value : []
    })
    
    const currentPage = ref(1)
    const pageSize = ref(20)
    const total = ref(0)
    
    // 搜索和筛选
    const searchQuery = ref('')
    const statusFilter = ref('')
    const typeFilter = ref('')
    const isOpenFilter = ref('')
    const teacherLimitFilter = ref('') // 新增教师指导筛选
    const locationFilter = ref('') // 新增地点筛选
    const dateRangeFilter = ref([]) // 新增日期范围筛选
    
    // 对话框状态
    const showAddCompetitionDialog = ref(false)
    const showRegistrationsDialog = ref(false)
    const showSubmissionsDialog = ref(false)
    const showResultsDialog = ref(false)
    const showAddResultDialog = ref(false)
    const showCompetitionDetailDialog = ref(false) // 新增竞赛详情对话框
    
    // 表单数据
    const editingCompetition = ref(null)
    const currentCompetition = ref(null)
    const submitting = ref(false)
    const submittingResult = ref(false)
    
    // 列表数据
    const registrations = ref([])
    const submissions = ref([])
    const results = ref([])
    const availableStudents = ref([])
    const registrationsLoading = ref(false)
    const submissionsLoading = ref(false)
    const resultsLoading = ref(false)
    
    // 表单引用
    const competitionFormRef = ref()
    const resultFormRef = ref()
    
    // 竞赛表单
    const competitionForm = reactive({
      title: '',
      type: '',
      organizer: '',
      start_time: '',
      end_time: '',
      description: '',
      attachment: '',
      is_open: true,
      status: 'draft',
      max_participants: 100,
      registration_start: '', // 新增报名开始时间
      registration_end: '', // 新增报名截止时间
      submission_deadline: '', // 新增作品提交截止时间
      location: '', // 新增竞赛地点
      contact: '', // 新增联系方式
      department_limit: '', // 新增院系限制
      teacher_limit: false, // 新增教师指导
      rules: '', // 新增竞赛规则
      requirements: '', // 新增报名要求
      judging_method: '', // 新增评审方式
      important_notes: '', // 新增重要注意事项
      website: '', // 新增相关链接
      qq_group: '', // 新增QQ群
      file_formats: '', // 新增支持文件格式
      file_size_limit: '' // 新增文件大小限制
    })
    
    // 结果表单
    const resultForm = reactive({
      student_id: null,
      award_level: '',
      final_score: null,
      certificate_url: ''
    })
    
    // 表单验证规则
    const competitionRules = {
      title: [{ required: true, message: '请输入竞赛名称', trigger: 'blur' }],
      type: [{ required: true, message: '请选择竞赛类型', trigger: 'change' }],
      organizer: [{ required: true, message: '请输入主办方', trigger: 'blur' }],
      registration_start: [{ required: false, message: '请选择报名开始时间', trigger: 'change' }],
      registration_end: [{ required: false, message: '请选择报名截止时间', trigger: 'change' }],
      submission_deadline: [{ required: false, message: '请选择作品提交截止时间', trigger: 'change' }],
      start_time: [{ required: true, message: '请选择比赛开始时间', trigger: 'change' }],
      end_time: [{ required: true, message: '请选择比赛结束时间', trigger: 'change' }],
      description: [{ required: true, message: '请输入竞赛描述', trigger: 'blur' }],
      location: [{ required: false, message: '请输入竞赛地点', trigger: 'blur' }],
      contact: [{ required: false, message: '请输入联系方式', trigger: 'blur' }],
      department_limit: [{ required: false, message: '请输入院系限制', trigger: 'blur' }],
      teacher_limit: [{ required: false, message: '请选择是否需要教师指导', trigger: 'change' }],
      rules: [{ required: false, message: '请输入竞赛规则', trigger: 'blur' }],
      requirements: [{ required: false, message: '请输入报名要求', trigger: 'blur' }],
      judging_method: [{ required: false, message: '请输入评审方式', trigger: 'blur' }],
      important_notes: [{ required: false, message: '请输入重要注意事项', trigger: 'blur' }],
      website: [{ required: false, message: '请输入相关链接', trigger: 'blur' }],
      qq_group: [{ required: false, message: '请输入QQ群号', trigger: 'blur' }],
      file_formats: [{ required: false, message: '请输入支持的文件格式', trigger: 'blur' }],
      file_size_limit: [{ required: false, message: '请输入文件大小限制', trigger: 'blur' }]
    }
    
    const resultRules = {
      student_id: [{ required: true, message: '请选择学生', trigger: 'change' }],
      award_level: [{ required: true, message: '请选择获奖等级', trigger: 'change' }]
    }
    
    // 获取竞赛列表
    const fetchCompetitions = async () => {
      loading.value = true
      try {
        const params = {
          page: currentPage.value,
          size: pageSize.value,
          search: searchQuery.value,
          status: statusFilter.value,
          type: typeFilter.value,
          is_open: isOpenFilter.value,
          teacher_limit: teacherLimitFilter.value, // 添加教师指导筛选参数
          location: locationFilter.value, // 添加地点筛选参数
          start_time: dateRangeFilter.value ? dateRangeFilter.value[0] : undefined, // 添加开始时间筛选参数
          end_time: dateRangeFilter.value ? dateRangeFilter.value[1] : undefined // 添加结束时间筛选参数
        }
        
        const response = await competitionService.getCompetitions(params)
        console.log('竞赛列表响应:', response) // 调试日志
        
        // 根据后端响应结构调整数据获取
        if (response && response.data) {
          competitions.value = response.data.list || []
          total.value = response.data.total || 0
        } else if (response && Array.isArray(response)) {
          // 如果响应直接是数组
          competitions.value = response
          total.value = response.length
        } else {
          competitions.value = []
          total.value = 0
        }
        
        console.log('处理后的竞赛数据:', competitions.value) // 调试日志
      } catch (error) {
        console.error('获取竞赛列表失败:', error)
        ElMessage.error('获取竞赛列表失败: ' + (error.response?.data?.message || error.message))
        competitions.value = []
        total.value = 0
      } finally {
        loading.value = false
      }
    }
    
    // 搜索处理
    const handleSearch = () => {
      currentPage.value = 1
      fetchCompetitions()
    }
    
    // 分页处理
    const handleSizeChange = (val) => {
      pageSize.value = val
      fetchCompetitions()
    }
    
    const handleCurrentChange = (val) => {
      currentPage.value = val
      fetchCompetitions()
    }
    
    // 清除所有筛选条件
    const clearFilters = () => {
      searchQuery.value = ''
      statusFilter.value = ''
      typeFilter.value = ''
      isOpenFilter.value = ''
      teacherLimitFilter.value = ''
      locationFilter.value = ''
      dateRangeFilter.value = []
      handleSearch() // 重新搜索以应用新的筛选条件
    }
    
    // 查看竞赛详情
    const viewCompetition = (competition) => {
      currentCompetition.value = competition
      showCompetitionDetailDialog.value = true
    }
    
    // 编辑竞赛
    const editCompetition = (competition) => {
      editingCompetition.value = competition
      
      // 重置表单，然后加载竞赛数据
      resetCompetitionForm()
      Object.assign(competitionForm, competition)
      
      // 设置默认值（如果字段为空）
      if (!competitionForm.registration_start) {
        const now = new Date()
        competitionForm.registration_start = now.toISOString().slice(0, 19).replace('T', ' ')
      }
      
      if (!competitionForm.registration_end && competitionForm.start_time) {
        // 如果报名截止时间未设置，默认设置为比赛开始时间
        competitionForm.registration_end = competitionForm.start_time
      }
      
      if (!competitionForm.submission_deadline && competitionForm.end_time) {
        // 如果作品提交截止时间未设置，默认设置为比赛结束时间
        competitionForm.submission_deadline = competitionForm.end_time
      }
      
      if (!competitionForm.location) {
        competitionForm.location = '待定'
      }
      
      if (!competitionForm.contact) {
        competitionForm.contact = '待定'
      }
      
      if (!competitionForm.department_limit) {
        competitionForm.department_limit = '不限'
      }
      
      if (!competitionForm.file_formats) {
        competitionForm.file_formats = 'PDF、PPT、DOC、ZIP'
      }
      
      if (!competitionForm.file_size_limit) {
        competitionForm.file_size_limit = '单个文件不超过50MB'
      }
      
      showAddCompetitionDialog.value = true
    }
    
    // 提交竞赛表单
    const submitCompetition = async () => {
      if (!competitionFormRef.value) return
      
      try {
        await competitionFormRef.value.validate()
        submitting.value = true
        
        // 准备提交的数据，确保格式正确
        const submitData = {
          title: competitionForm.title,
          type: competitionForm.type,
          organizer: competitionForm.organizer,
          description: competitionForm.description,
          attachment: competitionForm.attachment,
          is_open: competitionForm.is_open,
          status: competitionForm.status,
          max_participants: competitionForm.max_participants,
          registration_start: competitionForm.registration_start, // 新增报名开始时间
          registration_end: competitionForm.registration_end, // 新增报名截止时间
          submission_deadline: competitionForm.submission_deadline, // 新增作品提交截止时间
          location: competitionForm.location, // 新增竞赛地点
          contact: competitionForm.contact, // 新增联系方式
          department_limit: competitionForm.department_limit, // 新增院系限制
          teacher_limit: competitionForm.teacher_limit, // 新增教师指导
          rules: competitionForm.rules, // 新增竞赛规则
          requirements: competitionForm.requirements, // 新增报名要求
          judging_method: competitionForm.judging_method, // 新增评审方式
          important_notes: competitionForm.important_notes, // 新增重要注意事项
          website: competitionForm.website, // 新增相关链接
          qq_group: competitionForm.qq_group, // 新增QQ群
          file_formats: competitionForm.file_formats, // 新增支持文件格式
          file_size_limit: competitionForm.file_size_limit // 新增文件大小限制
        }
        
        // 处理时间字段
        if (competitionForm.registration_start) {
          submitData.registration_start = new Date(competitionForm.registration_start)
        }
        if (competitionForm.registration_end) {
          submitData.registration_end = new Date(competitionForm.registration_end)
        }
        if (competitionForm.start_time) {
          submitData.start_time = new Date(competitionForm.start_time)
        }
        if (competitionForm.end_time) {
          submitData.end_time = new Date(competitionForm.end_time)
        }
        if (competitionForm.submission_deadline) {
          submitData.submission_deadline = new Date(competitionForm.submission_deadline)
        }
        
        if (editingCompetition.value) {
          await competitionService.updateCompetition(editingCompetition.value.id, submitData)
          ElMessage.success('竞赛更新成功')
        } else {
          await competitionService.createCompetition(submitData)
          ElMessage.success('竞赛创建成功')
        }
        
        showAddCompetitionDialog.value = false
        fetchCompetitions()
        resetCompetitionForm()
      } catch (error) {
        console.error('提交竞赛表单失败:', error)
        
        // 根据错误状态码提供具体的错误信息
        if (error.response?.status === 403) {
          ElMessage.error('权限不足，需要管理员权限才能创建竞赛。请使用管理员账号登录。')
        } else if (error.response?.status === 401) {
          ElMessage.error('登录已过期，请重新登录。正在跳转到登录页面...')
          // 清除过期的认证信息
          localStorage.removeItem('token')
          localStorage.removeItem('userInfo')
          localStorage.removeItem('userRole')
          // 延迟跳转到登录页面
          setTimeout(() => {
            window.location.href = '/login'
          }, 2000)
        } else {
          ElMessage.error(error.response?.data?.message || error.message || '操作失败')
        }
      } finally {
        submitting.value = false
      }
    }
    
    // 重置竞赛表单
    const resetCompetitionForm = () => {
      editingCompetition.value = null
      const now = new Date()
      const currentDateTime = now.toISOString().slice(0, 19).replace('T', ' ')
      
      Object.assign(competitionForm, {
        title: '',
        type: '',
        organizer: '',
        start_time: currentDateTime, // 默认设置为当前日期时间
        end_time: '',
        description: '',
        attachment: '',
        is_open: true,
        status: 'draft',
        max_participants: 100,
        registration_start: '', // 重置报名开始时间
        registration_end: '', // 重置报名截止时间
        submission_deadline: '', // 重置作品提交截止时间
        location: '', // 重置竞赛地点
        contact: '', // 重置联系方式
        department_limit: '', // 重置院系限制
        teacher_limit: false, // 重置教师指导
        rules: '', // 重置竞赛规则
        requirements: '', // 重置报名要求
        judging_method: '', // 重置评审方式
        important_notes: '', // 重置重要注意事项
        website: '', // 重置相关链接
        qq_group: '', // 重置QQ群
        file_formats: '', // 重置支持文件格式
        file_size_limit: '' // 重置文件大小限制
      })
    }
    
    // 关闭竞赛对话框
    const closeCompetitionDialog = () => {
      showAddCompetitionDialog.value = false
      resetCompetitionForm()
    }
    
    // 管理报名
    const manageRegistrations = async (competition) => {
      currentCompetition.value = competition
      showRegistrationsDialog.value = true
      await fetchRegistrations(competition.id)
    }
    
    // 获取报名记录
    const fetchRegistrations = async (competitionId) => {
      registrationsLoading.value = true
      try {
        const response = await competitionService.getCompetitionRegistrations(competitionId)
        registrations.value = response.data.list || []
      } catch (error) {
        ElMessage.error('获取报名记录失败')
      } finally {
        registrationsLoading.value = false
      }
    }
    
    // 管理作品
    const manageSubmissions = async (competition) => {
      currentCompetition.value = competition
      showSubmissionsDialog.value = true
      await fetchSubmissions(competition.id)
    }
    
    // 获取提交作品
    const fetchSubmissions = async (competitionId) => {
      submissionsLoading.value = true
      try {
        const response = await competitionService.getCompetitionSubmissions(competitionId)
        submissions.value = response.data.list || []
      } catch (error) {
        ElMessage.error('获取提交作品失败')
      } finally {
        submissionsLoading.value = false
      }
    }
    
    // 管理结果
    const manageResults = async (competition) => {
      currentCompetition.value = competition
      showResultsDialog.value = true
      await fetchResults(competition.id)
    }
    
    // 获取获奖结果
    const fetchResults = async (competitionId) => {
      resultsLoading.value = true
      try {
        // 这里需要根据实际API调整
        const response = await competitionService.getCompetitionRegistrations(competitionId)
        results.value = response.data.list || []
        availableStudents.value = response.data.list || []
      } catch (error) {
        ElMessage.error('获取获奖结果失败')
      } finally {
        resultsLoading.value = false
      }
    }
    
    // 提交获奖结果
    const submitResult = async () => {
      if (!resultFormRef.value) return
      
      try {
        await resultFormRef.value.validate()
        submittingResult.value = true
        
        await competitionService.submitResult(currentCompetition.value.id, resultForm)
        ElMessage.success('获奖结果添加成功')
        
        showAddResultDialog.value = false
        fetchResults(currentCompetition.value.id)
        resetResultForm()
      } catch (error) {
        ElMessage.error(error.message || '添加失败')
      } finally {
        submittingResult.value = false
      }
    }
    
    // 重置结果表单
    const resetResultForm = () => {
      Object.assign(resultForm, {
        student_id: null,
        award_level: '',
        final_score: null,
        certificate_url: ''
      })
    }
    
    // 下拉菜单命令处理
    const handleCommand = async (command) => {
      const { action, id } = command
      
      switch (action) {
        case 'view':
          await viewCompetition(competitions.value.find(c => c.id === id))
          break
        case 'toggleStatus':
          await toggleCompetitionStatus(id)
          break
        case 'assignJudges':
          await assignJudges(id)
          break
        case 'judgingProgress':
          await viewJudgingProgress(id)
          break
        case 'exportRegistrations':
          await exportData(id, 'registrations')
          break
        case 'exportResults':
          await exportData(id, 'results')
          break
        case 'manageResults':
          await manageResults(competitions.value.find(c => c.id === id))
          break
        case 'finalize':
          await finalizeResults(id)
          break
        case 'delete':
          await deleteCompetition(id)
          break
      }
    }
    
    // 切换竞赛状态
    const toggleCompetitionStatus = async (id) => {
      try {
        const competition = competitions.value.find(c => c.id === id)
        if (!competition) return
        
        await competitionService.updateCompetition(id, {
          is_open: !competition.is_open
        })
        
        ElMessage.success('状态更新成功')
        fetchCompetitions()
      } catch (error) {
        ElMessage.error('状态更新失败')
      }
    }
    
    // 分配评审教师
    const assignJudges = async (id) => {
      ElMessage.info('分配评审教师功能待实现')
    }
    
    // 查看评审进度
    const viewJudgingProgress = async (id) => {
      try {
        const response = await competitionService.getJudgingProgress(id)
        ElMessage.info('评审进度功能待实现')
      } catch (error) {
        ElMessage.error('获取评审进度失败')
      }
    }
    
    // 导出数据
    const exportData = async (id, type) => {
      try {
        const response = await competitionService.exportCompetitionData(id, type)
        
        // 创建下载链接
        const blob = new Blob([response.data], { type: 'application/vnd.ms-excel' })
        const url = window.URL.createObjectURL(blob)
        const link = document.createElement('a')
        link.href = url
        link.download = `competition_${id}_${type}.xlsx`
        link.click()
        window.URL.revokeObjectURL(url)
        
        ElMessage.success('导出成功')
      } catch (error) {
        ElMessage.error('导出失败')
      }
    }
    
    // 最终确认成绩
    const finalizeResults = async (id) => {
      try {
        await ElMessageBox.confirm('确认要最终确认成绩吗？此操作不可撤销。', '确认操作', {
          confirmButtonText: '确认',
          cancelButtonText: '取消',
          type: 'warning'
        })
        
        await competitionService.finalizeResults(id)
        ElMessage.success('成绩最终确认成功')
        fetchCompetitions()
      } catch (error) {
        if (error !== 'cancel') {
          ElMessage.error('最终确认失败')
        }
      }
    }
    
    // 删除竞赛
    const deleteCompetition = async (id) => {
      try {
        await ElMessageBox.confirm('确认要删除这个竞赛吗？此操作不可撤销。', '确认删除', {
          confirmButtonText: '删除',
          cancelButtonText: '取消',
          type: 'warning'
        })
        
        await competitionService.deleteCompetition(id)
        ElMessage.success('竞赛删除成功')
        fetchCompetitions()
      } catch (error) {
        if (error !== 'cancel') {
          ElMessage.error('删除失败')
        }
      }
    }
    
    // 工具函数
    const formatDate = (date) => {
      if (!date) return ''
      return new Date(date).toLocaleString()
    }
    
    const formatFileSize = (bytes) => {
      if (!bytes) return '0 B'
      const k = 1024
      const sizes = ['B', 'KB', 'MB', 'GB']
      const i = Math.floor(Math.log(bytes) / Math.log(k))
      return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
    }
    
    const getTypeTagType = (type) => {
      const types = {
        '校级': 'info',
        '省级': 'warning',
        '国家级': 'danger',
        '国际级': 'success'
      }
      return types[type] || 'info'
    }
    
    const getStatusTagType = (status) => {
      const types = {
        'draft': 'info',
        'registration': 'warning',
        'submission': 'primary',
        'review': 'success',
        'completed': 'success'
      }
      return types[status] || 'info'
    }
    
    const getStatusLabel = (status) => {
      const labels = {
        'draft': '草稿',
        'registration': '报名中',
        'submission': '提交中',
        'review': '评审中',
        'completed': '已完成'
      }
      return labels[status] || status
    }
    
    const getRegistrationStatusType = (status) => {
      const types = {
        'registered': 'info',
        'approved': 'success',
        'rejected': 'danger',
        'withdrawn': 'warning'
      }
      return types[status] || 'info'
    }
    
    const getRegistrationStatusLabel = (status) => {
      const labels = {
        'registered': '已报名',
        'approved': '已通过',
        'rejected': '已拒绝',
        'withdrawn': '已撤回'
      }
      return labels[status] || status
    }
    
    const getSubmissionStatusType = (status) => {
      const types = {
        'submitted': 'info',
        'reviewing': 'warning',
        'approved': 'success',
        'rejected': 'danger'
      }
      return types[status] || 'info'
    }
    
    const getSubmissionStatusLabel = (status) => {
      const labels = {
        'submitted': '已提交',
        'reviewing': '评审中',
        'approved': '已通过',
        'rejected': '已拒绝'
      }
      return labels[status] || status
    }
    
    const getAwardLevelType = (level) => {
      const types = {
        '一等奖': 'danger',
        '二等奖': 'warning',
        '三等奖': 'success',
        '优秀奖': 'info',
        '参与奖': 'info'
      }
      return types[level] || 'info'
    }
    
    // 下载作品
    const downloadSubmission = (submission) => {
      if (submission.file_url) {
        window.open(submission.file_url, '_blank')
      } else {
        ElMessage.warning('文件链接不存在')
      }
    }
    
    // 查看作品
    const viewSubmission = (submission) => {
      ElMessage.info('查看作品详情功能待实现')
    }
    
    // 下载证书
    const downloadCertificate = (result) => {
      if (result.certificate_url) {
        window.open(result.certificate_url, '_blank')
      } else {
        ElMessage.warning('证书链接不存在')
      }
    }
    
    // 编辑结果
    const editResult = (result) => {
      ElMessage.info('编辑获奖结果功能待实现')
    }
    
    // 删除结果
    const deleteResult = async (result) => {
      try {
        await ElMessageBox.confirm('确认要删除这个获奖结果吗？', '确认删除', {
          confirmButtonText: '删除',
          cancelButtonText: '取消',
          type: 'warning'
        })
        
        ElMessage.success('删除成功')
        fetchResults(currentCompetition.value.id)
      } catch (error) {
        if (error !== 'cancel') {
          ElMessage.error('删除失败')
        }
      }
    }
    
    // 初始化
    onMounted(() => {
      fetchCompetitions()
      resetCompetitionForm() // 初始化表单，设置默认开始时间
    })
    
    return {
      // 响应式数据
      loading,
      competitions,
      safeCompetitions,
      currentPage,
      pageSize,
      total,
      searchQuery,
      statusFilter,
      typeFilter,
      isOpenFilter,
      teacherLimitFilter, // 新增教师指导筛选
      locationFilter, // 新增地点筛选
      dateRangeFilter, // 新增日期范围筛选
      
      // 对话框状态
      showAddCompetitionDialog,
      showRegistrationsDialog,
      showSubmissionsDialog,
      showResultsDialog,
      showAddResultDialog,
      showCompetitionDetailDialog, // 新增竞赛详情对话框
      
      // 表单数据
      editingCompetition,
      currentCompetition,
      submitting,
      submittingResult,
      competitionForm,
      resultForm,
      
      // 列表数据
      registrations,
      submissions,
      results,
      availableStudents,
      registrationsLoading,
      submissionsLoading,
      resultsLoading,
      
      // 表单引用
      competitionFormRef,
      resultFormRef,
      
      // 验证规则
      competitionRules,
      resultRules,
      
      // 方法
      handleSearch,
      handleSizeChange,
      handleCurrentChange,
      viewCompetition,
      editCompetition,
      submitCompetition,
      closeCompetitionDialog,
      manageRegistrations,
      manageSubmissions,
      manageResults,
      submitResult,
      handleCommand,
      clearFilters, // 新增清除筛选方法
      
      // 工具函数
      formatDate,
      formatFileSize,
      getTypeTagType,
      getStatusTagType,
      getStatusLabel,
      getRegistrationStatusType,
      getRegistrationStatusLabel,
      getSubmissionStatusType,
      getSubmissionStatusLabel,
      getAwardLevelType,
      downloadSubmission,
      viewSubmission,
      downloadCertificate,
      editResult,
      deleteResult
    }
  }
}
</script>

<style scoped>
.competition-management {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
  color: #303133;
}

.search-card {
  margin-bottom: 20px;
}

.competition-list-card {
  margin-bottom: 20px;
}

.pagination-container {
  margin-top: 20px;
  text-align: right;
}

.results-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.results-header h3 {
  margin: 0;
}

.dialog-footer {
  text-align: right;
}

/* 操作按钮样式 */
.action-buttons {
  display: flex;
  align-items: center;
  gap: 8px;
}

.main-actions {
  display: flex;
  align-items: center;
}

.main-actions .el-button {
  margin-right: 0;
  border-radius: 0;
}

.main-actions .el-button:first-child {
  border-top-left-radius: 4px;
  border-bottom-left-radius: 4px;
}

.main-actions .el-button:last-child {
  border-top-right-radius: 4px;
  border-bottom-right-radius: 4px;
  margin-right: 8px;
}

.secondary-actions {
  margin-left: auto;
}

/* 下拉菜单图标样式 */
.el-dropdown-menu__item .el-icon {
  margin-right: 8px;
  width: 16px;
  height: 16px;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .action-buttons {
    flex-direction: column;
    gap: 4px;
  }
  
  .main-actions {
    width: 100%;
  }
  
  .main-actions .el-button {
    flex: 1;
    min-width: 60px;
  }
  
  .secondary-actions {
    margin-left: 0;
    width: 100%;
  }
  
  .secondary-actions .el-button {
    width: 100%;
  }
}
</style> 
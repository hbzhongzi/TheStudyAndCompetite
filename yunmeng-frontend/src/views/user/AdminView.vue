<template>
  <div class="admin-container">
    <el-container>
      <!-- 顶部导航 -->
      <el-header class="header">
        <div class="header-left">
          <h2>管理员端 - 云梦高校科研竞赛管理系统</h2>
        </div>
        <div class="header-right">
          <el-dropdown @command="handleCommand">
            <span class="user-info">
              <i class="el-icon-setting"></i>
              系统管理员
              <i class="el-icon-arrow-down"></i>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">个人信息</el-dropdown-item>
                <el-dropdown-item command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      
      <el-container>
        <!-- 侧边栏 -->
        <el-aside width="250px" class="sidebar">
          <el-menu
            :default-active="activeMenu"
            class="sidebar-menu"
            @select="handleMenuSelect"
          >
            <el-menu-item index="dashboard">
              <i class="el-icon-s-home"></i>
              <span>系统概览</span>
            </el-menu-item>
            <el-menu-item index="users">
              <i class="el-icon-user"></i>
              <span>用户管理</span>
            </el-menu-item>
            <el-sub-menu index="project-management">
              <template #title>
                <i class="el-icon-s-operation"></i>
                <span>项目管理</span>
              </template>
              <el-menu-item index="project-overview">项目概览</el-menu-item>
              <el-menu-item index="project-review">项目审核</el-menu-item>
              <el-menu-item index="project-types">分类管理</el-menu-item>
              <el-menu-item index="review-flows">审核流程</el-menu-item>
              <el-menu-item index="project-stats">项目统计</el-menu-item>
            </el-sub-menu>
            <el-sub-menu index="competitions">
              <template #title>
                <i class="el-icon-trophy"></i>
                <span>竞赛管理</span>
              </template>
              <el-menu-item index="competition-basic">基础竞赛管理</el-menu-item>
              <el-menu-item index="competition-results">审核报名</el-menu-item>
              <el-menu-item index="competition-registrations">竞赛任务分配</el-menu-item>
              <el-menu-item index="audit-logs">审计日志</el-menu-item>
            </el-sub-menu>
            <el-menu-item index="system">
              <i class="el-icon-setting"></i>
              <span>系统设置</span>
            </el-menu-item>
            <el-menu-item index="reports">
              <i class="el-icon-s-data"></i>
              <span>数据统计</span>
            </el-menu-item>
            <el-menu-item index="logs">
              <i class="el-icon-document-copy"></i>
              <span>系统日志</span>
            </el-menu-item>
          </el-menu>
        </el-aside>
        
        <!-- 主要内容区域 -->
        <el-main class="main-content">
          <div class="content-header">
            <h3>{{ pageTitle }}</h3>
          </div>
          
          <div class="content-body">
            <!-- 系统概览页面 -->
            <div v-if="activeMenu === 'dashboard'" class="welcome-content">
              <DashboardView />
            </div>
            
            <!-- 用户管理页面 -->
            <div v-else-if="activeMenu === 'users'" class="welcome-content">
              <UserManagement />
            </div>
            
            <!-- 项目管理页面 -->
            <div v-else-if="activeMenu === 'project-overview'" class="welcome-content">
              <ProjectOverview />
            </div>
            <div v-else-if="activeMenu === 'project-review'" class="welcome-content">
              <ProjectReview />
            </div>
            <div v-else-if="activeMenu === 'project-types'" class="welcome-content">
              <ProjectTypeManagement />
            </div>
            <div v-else-if="activeMenu === 'review-flows'" class="welcome-content">
              <ReviewFlows />
            </div>
            <div v-else-if="activeMenu === 'project-stats'" class="welcome-content">
              <ProjectStats />
            </div>
            
            <!-- 竞赛管理页面 -->
            <div v-else-if="activeMenu === 'competition-basic'" class="welcome-content">
              <CompetitionManagement />
            </div>
            <div v-else-if="activeMenu === 'competition-registrations'" class="welcome-content">
              <CompetitionRegistrationReview />
            </div>
            <div v-else-if="activeMenu === 'competition-results'" class="welcome-content">
              <CompetitionResults />
            </div>
            <div v-else-if="activeMenu === 'audit-logs'" class="welcome-content">
              <CompetitionAuditLogs />
            </div>
            
            <!-- 系统设置页面 -->
            <div v-else-if="activeMenu === 'system'" class="welcome-content">
              <SystemSettings />
            </div>
            
            <!-- 数据统计页面 -->
            <div v-else-if="activeMenu === 'reports'" class="welcome-content">
              <DataReports />
            </div>
            
            <!-- 系统日志页面 -->
            <div v-else-if="activeMenu === 'logs'" class="welcome-content">
              <SystemLogs />
            </div>
            
            <!-- 默认显示系统概览 -->
            <div v-else class="welcome-content">
              <DashboardView />
            </div>
          </div>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import UserManagement from './UserManagement.vue'
import ProjectManagement from './ProjectManagement.vue'
import ProjectTypeManagement from './ProjectTypeManagement.vue'
import CompetitionManagement from './CompetitionManagement.vue'
import CompetitionRegistrationReview from './CompetitionRegistrationReview.vue'
import CompetitionResults from './CompetitionResults.vue'
import CompetitionAuditLogs from './CompetitionAuditLogs.vue'
import SystemSettings from './SystemSettings.vue'
import DataReports from './DataReports.vue'
import SystemLogs from './SystemLogs.vue'
import DashboardView from './DashboardView.vue'
import ProjectOverview from './ProjectOverview.vue'
import ProjectReview from './ProjectReview.vue'
import ReviewFlows from './ReviewFlows.vue'
import ProjectStats from './ProjectStats.vue'

const router = useRouter()
const activeMenu = ref('dashboard')

const pageTitle = computed(() => {
  const menuMap = {
    dashboard: '系统概览',
    users: '用户管理',
    'project-overview': '项目概览',
    'project-review': '项目审核',
    'project-types': '项目分类管理',
    'review-flows': '审核流程',
    'project-stats': '项目统计',
    'competition-basic': '基础竞赛管理',
    'competition-registrations': '报名审核',
    'competition-results': '结果管理',
    'audit-logs': '审计日志',
    system: '系统设置',
    reports: '数据统计',
    logs: '系统日志'
  }
  return menuMap[activeMenu.value] || '系统概览'
})

const handleMenuSelect = (index) => {
  activeMenu.value = index
}

const navigateToMenu = (menu) => {
  activeMenu.value = menu
}

const handleCommand = async (command) => {
  if (command === 'logout') {
    try {
      // 清除本地存储
      localStorage.removeItem('token')
      localStorage.removeItem('userInfo')
      localStorage.removeItem('userRole')
      ElMessage.success('退出登录成功')
      router.push('/login')
    } catch (error) {
      ElMessage.error('退出登录失败')
    }
  } else if (command === 'profile') {
    activeMenu.value = 'profile'
  }
}
</script>

<style scoped>
.admin-container {
  height: 100vh;
}

.header {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.header h2 {
  margin: 0;
  font-size: 20px;
  font-weight: 600;
}

.user-info {
  color: white;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
}

.sidebar {
  background: #f8f9fa;
  border-right: 1px solid #e9ecef;
}

.sidebar-menu {
  border-right: none;
  height: 100%;
}

.main-content {
  background: #f5f7fa;
  padding: 20px;
}

.content-header {
  margin-bottom: 20px;
}

.content-header h3 {
  margin: 0;
  color: #2c3e50;
  font-size: 24px;
  font-weight: 600;
}

.welcome-content {
  background: white;
  border-radius: 10px;
  padding: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.stat-card {
  border-radius: 10px;
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
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
  display: flex;
  align-items: center;
  justify-content: center;
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

.recent-activity {
  border-radius: 10px;
}

.el-timeline-item__content {
  color: #2c3e50;
}

.el-timeline-item__content h4 {
  margin: 0 0 8px 0;
  color: #2c3e50;
}

.el-timeline-item__content p {
  margin: 0;
  color: #7f8c8d;
  line-height: 1.6;
}

.el-button {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.el-button i {
  font-size: 20px;
}
</style> 
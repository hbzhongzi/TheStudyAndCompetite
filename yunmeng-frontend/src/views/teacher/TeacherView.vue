<template>
  <div class="teacher-container">
    <el-container>
      <!-- 顶部导航 -->
      <el-header class="header">
        <div class="header-left">
          <h2>教师端 - 云梦高校科研竞赛管理系统</h2>
        </div>
        <div class="header-right">
          <el-dropdown @command="handleCommand">
            <span class="user-info">
              <i class="el-icon-s-custom"></i>
              教师用户
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
              <el-icon><House /></el-icon>
              <span>首页概览</span>
            </el-menu-item>
            <el-sub-menu index="project-management">
              <template #title>
                <el-icon><Operation /></el-icon>
                <span>项目管理</span>
              </template>
              <el-menu-item index="project-overview">项目概览</el-menu-item>
              <el-menu-item index="project-detail">项目详情</el-menu-item>
              <el-menu-item index="project-review">延期审核</el-menu-item>
              <el-menu-item index="project-files">文件审核</el-menu-item>

            </el-sub-menu>
            <el-menu-item index="students">
              <el-icon><User /></el-icon>
              <span>学生管理</span>
            </el-menu-item>
            <el-sub-menu index="competitions">
              <template #title>
                <el-icon><Trophy /></el-icon>
                <span>竞赛指导</span>
              </template>
              <el-menu-item index="competition-guidance">竞赛指导</el-menu-item>
              <el-menu-item index="competition-judging">作品评审</el-menu-item>
            </el-sub-menu>
            <el-menu-item index="applications">
              <el-icon><EditPen /></el-icon>
              <span>申请审核</span>
            </el-menu-item>
            <el-menu-item index="reports">
              <el-icon><DataAnalysis /></el-icon>
              <span>统计报告</span>
            </el-menu-item>
          </el-menu>
        </el-aside>
        
        <!-- 主要内容区域 -->
        <el-main class="main-content">
          <div class="content-header">
            <h3>{{ pageTitle }}</h3>
          </div>
          
          <div class="content-body">
            <!-- 首页概览 -->
            <div v-if="activeMenu === 'dashboard'" class="welcome-content">
              <el-row :gutter="20">
                <el-col :span="6">
                  <el-card class="stat-card">
                    <div class="stat-content">
                      <div class="stat-icon">
                        <el-icon><House /></el-icon>
                      </div>
                      <div class="stat-info">
                        <h4>指导项目</h4>
                        <p class="stat-number">5</p>
                      </div>
                    </div>
                  </el-card>
                </el-col>
                
                <el-col :span="6">
                  <el-card class="stat-card">
                    <div class="stat-content">
                      <div class="stat-icon">
                        <el-icon><User /></el-icon>
                      </div>
                      <div class="stat-info">
                        <h4>指导学生</h4>
                        <p class="stat-number">12</p>
                      </div>
                    </div>
                  </el-card>
                </el-col>
                
                <el-col :span="6">
                  <el-card class="stat-card">
                    <div class="stat-content">
                      <div class="stat-icon">
                        <el-icon><Trophy /></el-icon>
                      </div>
                      <div class="stat-info">
                        <h4>竞赛指导</h4>
                        <p class="stat-number">3</p>
                      </div>
                    </div>
                  </el-card>
                </el-col>
                
                <el-col :span="6">
                  <el-card class="stat-card">
                    <div class="stat-content">
                      <div class="stat-icon">
                        <el-icon><EditPen /></el-icon>
                      </div>
                      <div class="stat-info">
                        <h4>待审核</h4>
                        <p class="stat-number">8</p>
                      </div>
                    </div>
                  </el-card>
                </el-col>
              </el-row>
              
              <el-row :gutter="20" style="margin-top: 20px;">
                <el-col :span="12">
                  <el-card class="recent-activity">
                    <template #header>
                      <span>待审核申请</span>
                    </template>
                    <el-table :data="pendingApplications" style="width: 100%">
                      <el-table-column prop="student" label="学生" width="100" />
                      <el-table-column prop="type" label="类型" width="100" />
                      <el-table-column prop="title" label="标题" />
                      <el-table-column prop="date" label="申请时间" width="120" />
                      <el-table-column label="操作" width="120">
                        <template #default>
                          <el-button size="small" type="primary">审核</el-button>
                        </template>
                      </el-table-column>
                    </el-table>
                  </el-card>
                </el-col>
                
                <el-col :span="12">
                  <el-card class="recent-activity">
                    <template #header>
                      <span>最近活动</span>
                    </template>
                    <el-timeline>
                      <el-timeline-item timestamp="2024-01-15" placement="top">
                        <el-card>
                          <h4>项目申请审核</h4>
                          <p>审核了学生张三的"智能校园管理系统"项目申请。</p>
                        </el-card>
                      </el-timeline-item>
                      <el-timeline-item timestamp="2024-01-12" placement="top">
                        <el-card>
                          <h4>竞赛指导</h4>
                          <p>为"全国大学生程序设计竞赛"提供了指导。</p>
                        </el-card>
                      </el-timeline-item>
                    </el-timeline>
                  </el-card>
                </el-col>
              </el-row>
            </div>
            
            <!-- 项目管理子界面 -->
            <div v-else-if="activeMenu === 'project-overview'" class="welcome-content">

              <TeacherProjectOverview />
            </div>
            <div v-else-if="activeMenu === 'project-review'" class="welcome-content">

              <TeacherProjectReview />
            </div>

            <div v-else-if="activeMenu === 'project-files'" class="welcome-content">

              <TeacherProjectFiles />
            </div>
            <div v-else-if="activeMenu === 'project-detail'" class="welcome-content">

              <TeacherProjectDetail />
            </div>
            
            <!-- 其他主要界面 -->
            <div v-else-if="activeMenu === 'projects'" class="welcome-content">
              <ProjectManagementView />
            </div>
            <div v-else-if="activeMenu === 'students'" class="welcome-content">
              <StudentManagementView />
            </div>
            <div v-else-if="activeMenu === 'competition-guidance'" class="welcome-content">
              <CompetitionGuidanceView />
            </div>
            <div v-else-if="activeMenu === 'competition-judging'" class="welcome-content">
              <CompetitionJudging />
            </div>
            <div v-else-if="activeMenu === 'applications'" class="welcome-content">
              <ApplicationReviewView />
            </div>
            <div v-else-if="activeMenu === 'reports'" class="welcome-content">
              <ReportView />
            </div>
            <div v-else-if="activeMenu === 'profile'" class="welcome-content">
              <TeacherProfileView />
            </div>
            
            <!-- 默认显示欢迎信息 -->
            <div v-else class="welcome-content">
              <h2>欢迎使用教师端管理系统</h2>
              <p>请从左侧菜单选择要使用的功能</p>
              
            </div>
          </div>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { 
  House, 
  Document, 
  Operation, 
  User, 
  Trophy, 
  EditPen, 
  DataAnalysis 
} from '@element-plus/icons-vue'
import ProjectManagementView from './ProjectManagementView.vue'
import StudentManagementView from './StudentManagementView.vue'
import CompetitionGuidanceView from './CompetitionGuidanceView.vue'
import CompetitionJudging from './CompetitionJudging.vue'
import ApplicationReviewView from './ApplicationReviewView.vue'
import ReportView from './ReportView.vue'
import TeacherProfileView from './TeacherProfileView.vue'
import TeacherProjectOverview from './TeacherProjectOverview.vue'
import TeacherProjectReview from './TeacherProjectReview.vue'
import TeacherProjectFiles from './TeacherProjectFiles.vue'
import TeacherProjectDetail from './TeacherProjectDetail.vue'
import TeacherProjectExtensions from './TeacherProjectExtensions.vue'


const router = useRouter()
const activeMenu = ref('dashboard')

const pageTitle = computed(() => {
  const menuMap = {
    dashboard: '首页概览',
    projects: '项目管理',
    'project-overview': '项目概览',
    'project-review': '项目审核',
    'project-files': '文件审核',
    'review-tasks': '审核任务',
    students: '学生管理',
    competitions: '竞赛指导',
    applications: '申请审核',
    reports: '统计报告',
    profile: '个人信息'
  }
  return menuMap[activeMenu.value] || '首页概览'
})

// 模拟待审核申请数据
const pendingApplications = ref([
  {
    student: '张三',
    type: '项目申请',
    title: '智能校园管理系统',
    date: '2024-01-15'
  },
  {
    student: '李四',
    type: '竞赛报名',
    title: '全国大学生程序设计竞赛',
    date: '2024-01-14'
  },
  {
    student: '王五',
    type: '项目申请',
    title: '在线学习平台',
    date: '2024-01-13'
  }
])

const handleMenuSelect = async (index) => {
  try {
    console.log('菜单选择:', index)
    activeMenu.value = index
    
    // 使用更安全的 nextTick
    await nextTick()
    
    // 等待 DOM 完全渲染
    setTimeout(() => {
      // 在这里执行需要 DOM 的操作
      console.log('当前活动菜单:', activeMenu.value)
      
      // 如果有特定操作需要 DOM，添加检查
      const element = document.querySelector('.your-element')
      if (element && element.nextSibling) {
        // 安全地操作
      }
    }, 100) // 等待 100ms 确保 DOM 更新完成
  } catch (error) {
    console.error('菜单切换错误:', error)
    ElMessage.error('页面切换失败，请重试')
  }
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
.teacher-container {
  height: 100vh;
}

.header {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
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
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
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
</style> 
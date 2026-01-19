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
            <el-menu-item index="projects">
              <el-icon><Document /></el-icon>
              <span>我的项目</span>
            </el-menu-item>
            <el-sub-menu index="project-management">
              <template #title>
                <el-icon><Operation /></el-icon>
                <span>项目管理</span>
              </template>
              <el-menu-item index="project-overview">项目概览</el-menu-item>
              <el-menu-item index="project-review">项目审核</el-menu-item>
              <el-menu-item index="project-milestones">里程碑管理</el-menu-item>
              <el-menu-item index="project-files">文件审核</el-menu-item>
              <el-menu-item index="project-extensions">延期审核</el-menu-item>
              <el-menu-item index="review-tasks">审核任务</el-menu-item>
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
            <!-- 添加调试信息 -->
            <p style="color: #999; font-size: 12px; margin: 5px 0 0 0;">
              当前菜单: {{ activeMenu }} | 页面标题: {{ pageTitle }}
            </p>
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
              <div style="background: #e8f4fd; padding: 10px; margin-bottom: 10px; border-radius: 4px;">
                <p><strong>调试信息:</strong> 显示项目概览组件 (activeMenu: {{ activeMenu }})</p>
              </div>
              <TeacherProjectOverview />
            </div>
            <div v-else-if="activeMenu === 'project-review'" class="welcome-content">
              <div style="background: #e8f4fd; padding: 10px; margin-bottom: 10px; border-radius: 4px;">
                <p><strong>调试信息:</strong> 显示项目审核组件 (activeMenu: {{ activeMenu }})</p>
              </div>
              <TeacherProjectReview />
            </div>
            <div v-else-if="activeMenu === 'project-milestones'" class="welcome-content">
              <div style="background: #e8f4fd; padding: 10px; margin-bottom: 10px; border-radius: 4px;">
                <p><strong>调试信息:</strong> 显示里程碑管理组件 (activeMenu: {{ activeMenu }})</p>
              </div>
              <TeacherProjectMilestones />
            </div>
            <div v-else-if="activeMenu === 'project-files'" class="welcome-content">
              <div style="background: #e8f4fd; padding: 10px; margin-bottom: 10px; border-radius: 4px;">
                <p><strong>调试信息:</strong> 显示文件审核组件 (activeMenu: {{ activeMenu }})</p>
              </div>
              <TeacherProjectFiles />
            </div>
            <div v-else-if="activeMenu === 'project-extensions'" class="welcome-content">
              <div style="background: #e8f4fd; padding: 10px; margin-bottom: 10px; border-radius: 4px;">
                <p><strong>调试信息:</strong> 显示延期审核组件 (activeMenu: {{ activeMenu }})</p>
              </div>
              <TeacherProjectExtensions />
            </div>
            <div v-else-if="activeMenu === 'review-tasks'" class="welcome-content">
              <div style="background: #e8f4fd; padding: 10px; margin-bottom: 10px; border-radius: 4px;">
                <p><strong>调试信息:</strong> 显示审核任务组件 (activeMenu: {{ activeMenu }})</p>
              </div>
              <TeacherReviewTasks />
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
              
              <!-- 添加测试信息 -->
              <el-card style="margin-top: 20px; background: #f8f9fa;">
                <template #header>
                  <span>调试信息</span>
                </template>
                <div>
                  <p><strong>当前选中的菜单:</strong> {{ activeMenu }}</p>
                  <p><strong>页面标题:</strong> {{ pageTitle }}</p>
                  <p><strong>可用菜单项:</strong></p>
                  <ul>
                    <li>dashboard - 首页概览</li>
                    <li>projects - 项目管理</li>
                    <li>project-overview - 项目概览</li>
                    <li>project-review - 项目审核</li>
                    <li>project-milestones - 里程碑管理</li>
                    <li>project-files - 文件审核</li>
                    <li>project-extensions - 延期审核</li>
                    <li>review-tasks - 审核任务</li>
                    <li>students - 学生管理</li>
                    <li>competition-guidance - 竞赛指导</li>
                    <li>competition-judging - 作品评审</li>
                    <li>applications - 申请审核</li>
                    <li>reports - 统计报告</li>
                    <li>profile - 个人信息</li>
                  </ul>
                </div>
              </el-card>
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
import TeacherProjectMilestones from './TeacherProjectMilestones.vue'
import TeacherProjectFiles from './TeacherProjectFiles.vue'
import TeacherProjectExtensions from './TeacherProjectExtensions.vue'
import TeacherReviewTasks from './TeacherReviewTasks.vue'
import ErrorBoundary from '../../components/ErrorBoundary.vue'

const router = useRouter()
const activeMenu = ref('dashboard')

const pageTitle = computed(() => {
  const menuMap = {
    dashboard: '首页概览',
    projects: '项目管理',
    'project-overview': '项目概览',
    'project-review': '项目审核',
    'project-milestones': '里程碑管理',
    'project-files': '文件审核',
    'project-extensions': '延期审核',
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
    console.log('菜单选择:', index) // 添加调试日志
    activeMenu.value = index
    
    // 等待DOM更新完成
    await nextTick()
    
    // 验证组件是否正确加载
    console.log('当前活动菜单:', activeMenu.value)
    
    // 强制触发响应式更新
    if (index.startsWith('project-') || index === 'review-tasks') {
      console.log('项目管理子菜单选中:', index)
    }
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
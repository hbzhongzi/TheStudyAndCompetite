<template>
  <div class="student-container">
    <el-container>
      <!-- 顶部导航 -->
      <el-header class="header">
        <div class="header-left">
          <h2>学生端 - 云梦高校科研竞赛管理系统</h2>
        </div>
        <div class="header-right">
          <el-dropdown @command="handleCommand">
            <span class="user-info">
              <i class="el-icon-user"></i>
              学生用户
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
              <span>首页概览</span>
            </el-menu-item>
            <el-sub-menu index="project-management">
              <template #title>
                <i class="el-icon-s-operation"></i>
                <span>项目管理</span>
              </template>
              <el-menu-item index="project-overview">项目概览</el-menu-item>
              <el-menu-item index="project-files">项目详情</el-menu-item>
              <el-menu-item index="project-extensions">延期申请</el-menu-item>
              <el-menu-item index="project-progress">进度跟踪</el-menu-item>
            </el-sub-menu>
            <el-sub-menu index="competitions">
              <template #title>
                <i class="el-icon-trophy"></i>
                <span>竞赛管理</span>
              </template>
              <el-menu-item index="competition-info">竞赛信息</el-menu-item>
              <el-menu-item index="competition-submission">作品提交</el-menu-item>
            </el-sub-menu>
            <el-menu-item index="my-competitions">
              <i class="el-icon-medal"></i>
              <span>我的竞赛</span>
            </el-menu-item>
            <el-menu-item index="applications">
              <i class="el-icon-edit-outline"></i>
              <span>申请管理</span>
            </el-menu-item>
            <el-menu-item index="profile">
              <i class="el-icon-user"></i>
              <span>个人信息</span>
            </el-menu-item>
          </el-menu>
        </el-aside>
        
        <!-- 主要内容区域 -->
        <el-main class="main-content">
          <div class="content-header">
            <h3>{{ pageTitle }}</h3>
          </div>
          
          <div class="content-body">
            <!-- 根据选中的菜单显示不同的内容 -->
            <div v-if="activeMenu === 'dashboard'" class="welcome-content">
              <el-row :gutter="20">
                <el-col :span="8">
                  <el-card class="stat-card">
                    <div class="stat-content">
                      <div class="stat-icon">
                        <i class="el-icon-document"></i>
                      </div>
                      <div class="stat-info">
                        <h4>参与项目</h4>
                        <p class="stat-number">{{ stats.projectCount }}</p>
                      </div>
                    </div>
                  </el-card>
                </el-col>
                
                <el-col :span="8">
                  <el-card class="stat-card">
                    <div class="stat-content">
                      <div class="stat-icon">
                        <i class="el-icon-clock"></i>
                      </div>
                      <div class="stat-info">
                        <h4>待处理</h4>
                        <p class="stat-number">{{ stats.pendingCount }}</p>
                      </div>
                    </div>
                  </el-card>
                </el-col>
                
                <el-col :span="8">
                  <el-card class="stat-card">
                    <div class="stat-content">
                      <div class="stat-icon">
                        <i class="el-icon-trophy"></i>
                      </div>
                      <div class="stat-info">
                        <h4>竞赛项目</h4>
                        <p class="stat-number">{{ stats.competitionCount }}</p>
                      </div>
                    </div>
                  </el-card>
                </el-col>
              </el-row>
              
              <!-- 最近活动 -->
              <el-row style="margin-top: 20px;">
                <el-col :span="24">
                  <el-card>
                    <template #header>
                      <span>最近活动</span>
                    </template>
                    <el-timeline>
                      <el-timeline-item 
                        v-for="activity in recentActivity" 
                        :key="activity.id"
                        :timestamp="activity.timestamp"
                        :type="activity.type"
                      >
                        <div class="activity-content">
                          <h4>{{ activity.title }}</h4>
                          <p>{{ activity.description }}</p>
                        </div>
                      </el-timeline-item>
                    </el-timeline>
                  </el-card>
                </el-col>
              </el-row>
            </div>
            
            <!-- 项目管理子界面 -->
            <div v-else-if="activeMenu === 'project-overview'" class="welcome-content">
              <ProjectOverview />
            </div>
            <div v-else-if="activeMenu === 'project-milestones'" class="welcome-content">
              <ProjectMilestones />
            </div>
            <div v-else-if="activeMenu === 'project-files'" class="welcome-content">
              <ProjectFiles />
            </div>
            <div v-else-if="activeMenu === 'project-extensions'" class="welcome-content">
              <ProjectExtensions />
            </div>
            <div v-else-if="activeMenu === 'project-progress'" class="welcome-content">
              <ProjectProgress />
            </div>
            
            <!-- 其他界面 -->
            <div v-else-if="currentComponent" class="welcome-content">
              <component :is="currentComponent" />
            </div>
            
            <!-- 默认显示欢迎信息 -->
            <div v-else class="welcome-content">
              <h2>欢迎使用学生端管理系统</h2>
              <p>请从左侧菜单选择要使用的功能</p>
            </div>
          </div>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import ProjectList from './ProjectManagement.vue'
import CompetitionView from './CompetitionView.vue'
import CompetitionSubmission from './CompetitionSubmission.vue'
import MyCompetitions from './MyCompetitions.vue'
import ApplicationView from './ApplicationView.vue'
import ProfileView from './ProfileView.vue'
import ProjectOverview from './ProjectOverview.vue'
import ProjectFiles from './ProjectFiles.vue'
import ProjectExtensions from './ProjectExtensions.vue'
import ProjectProgress from './ProjectProgress.vue'
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'  // 添加 ElMessageBox







const router = useRouter()
const activeMenu = ref('dashboard')

// 统计数据
const stats = ref({
  projectCount: 0,
  pendingCount: 0,
  competitionCount: 0
})

// 心跳检测
let heartbeatInterval

const startHeartbeat = () => {
  heartbeatInterval = setInterval(() => {
    const token = localStorage.getItem('token')
    if (!token) {
      console.log('心跳检测: token丢失')
      clearInterval(heartbeatInterval)
      ElMessage.warning('登录已过期，请重新登录')
      router.push('/login')
    }
  }, 30000)
}

// 检查登录状态
const checkLoginStatus = () => {
  const token = localStorage.getItem('token')
  const userRole = localStorage.getItem('userRole')
  
  if (!token || userRole !== 'student') {
    return false
  }
  return true
}

// 组件挂载
onMounted(() => {
  console.log('学生端组件初始化')
  
  // 检查登录状态
  if (!checkLoginStatus()) {
    ElMessage.warning('请先登录')
    router.push('/login')
    return
  }
  
  // 启动心跳检测
  startHeartbeat()
  
  // 加载数据
  loadStats()
})

// 组件卸载
onUnmounted(() => {
  if (heartbeatInterval) {
    clearInterval(heartbeatInterval)
  }
})

// 最近活动数据
const recentActivity = ref([
  {
    id: 1,
    title: '项目申请已提交',
    description: '您提交的"智能校园管理系统"项目申请已成功提交，等待审核。',
    timestamp: '2024-01-15',
    type: 'primary'
  },
  {
    id: 2,
    title: '竞赛报名成功',
    description: '您已成功报名参加"全国大学生程序设计竞赛"。',
    timestamp: '2024-01-10',
    type: 'success'
  },
  {
    id: 3,
    title: '项目进度更新',
    description: '您的"数据分析平台"项目进度已更新至75%。',
    timestamp: '2024-01-08',
    type: 'info'
  }
])

const pageTitle = computed(() => {
  const menuMap = {
    dashboard: '首页概览',
    projects: '我的项目',
    'competition-info': '竞赛信息',
    'competition-submission': '作品提交',
    'my-competitions': '我的竞赛',
    applications: '申请管理',
    profile: '个人信息'
  }
  return menuMap[activeMenu.value] || '首页概览'
})

// 动态组件映射
const currentComponent = computed(() => {
  const componentMap = {
    projects: ProjectList,
    'competition-info': CompetitionView,
    'competition-submission': CompetitionSubmission,
    'my-competitions': MyCompetitions,
    applications: ApplicationView,
    profile: ProfileView
  }
  return componentMap[activeMenu.value] || null
})

// 加载统计数据
const loadStats = async () => {
  try {
    const response = await projectService.getMyProjects()
    
    // 检查响应格式
    if (response && response.code === 200) {
      const projects = response.data || []
      
      stats.value.projectCount = projects.length
      stats.value.pendingCount = projects.filter(p => p.status === 'pending').length
      stats.value.competitionCount = projects.filter(p => p.type === '竞赛').length
    } else {
      // 如果响应格式不符合预期，使用默认值
      loadDefaultStats()
    }
  } catch (error) {
    console.error('加载统计数据失败:', error)
    // 设置默认值，避免页面显示异常
    loadDefaultStats()
  }
}

// 加载默认统计数据
const loadDefaultStats = () => {
  stats.value.projectCount = 3
  stats.value.pendingCount = 1
  stats.value.competitionCount = 1
}

const handleMenuSelect = async (index) => {
  try {
    activeMenu.value = index
    
    // 等待DOM更新完成
    await nextTick()
    
    // 如果切换到项目页面，重新加载统计数据
    if (index === 'projects') {
      await loadStats()
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

// 组件挂载时加载统计数据
onMounted(() => {
  loadStats()
})
</script>

<style scoped>
.student-container {
  height: 100vh;
}

.header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
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
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
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
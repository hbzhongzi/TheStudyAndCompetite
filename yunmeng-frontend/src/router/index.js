import { createRouter, createWebHistory } from 'vue-router'
import LoginView from '../views/login/LoginView.vue'

const routes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/login',
    name: 'Login',
    component: LoginView,
    meta: { requiresAuth: false }
  },
  {
    path: '/dashboard',
    name: 'Dashboard',
    component: () => import('../views/dashboard/DashboardView.vue'),
    meta: { requiresAuth: true }
  },
  {
    path: '/student',
    name: 'Student',
    component: () => import('../views/student/StudentView.vue'),
    meta: { requiresAuth: true, role: 'student' }
  },
  {
    path: '/teacher',
    name: 'Teacher',
    component: () => import('../views/teacher/TeacherView.vue'),
    meta: { requiresAuth: true, role: 'teacher' }
  },
  {
    path: '/teacher/projects',
    name: 'TeacherProjects',
    component: () => import('../views/teacher/ProjectManagementView.vue'),
    meta: { requiresAuth: true, role: 'teacher' }
  },
  {
    path: '/teacher/students',
    name: 'TeacherStudents',
    component: () => import('../views/teacher/StudentManagementView.vue'),
    meta: { requiresAuth: true, role: 'teacher' }
  },
  {
    path: '/teacher/competitions',
    name: 'TeacherCompetitions',
    component: () => import('../views/teacher/CompetitionGuidanceView.vue'),
    meta: { requiresAuth: true, role: 'teacher' }
  },
  {
    path: '/teacher/applications',
    name: 'TeacherApplications',
    component: () => import('../views/teacher/ApplicationReviewView.vue'),
    meta: { requiresAuth: true, role: 'teacher' }
  },
  {
    path: '/teacher/reports',
    name: 'TeacherReports',
    component: () => import('../views/teacher/ReportView.vue'),
    meta: { requiresAuth: true, role: 'teacher' }
  },
  {
    path: '/teacher/profile',
    name: 'TeacherProfile',
    component: () => import('../views/teacher/TeacherProfileView.vue'),
    meta: { requiresAuth: true, role: 'teacher' }
  },
  {
    path: '/admin',
    name: 'Admin',
    component: () => import('../views/user/AdminView.vue'),
    meta: { requiresAuth: true, role: 'admin' }
  },
  {
    path: '/admin/projects',
    name: 'AdminProjects',
    component: () => import('../views/user/ProjectManagement.vue'),
    meta: { requiresAuth: true, role: 'admin' }
  },
  {
    path: '/admin/project-types',
    name: 'AdminProjectTypes',
    component: () => import('../views/user/ProjectTypeManagement.vue'),
    meta: { requiresAuth: true, role: 'admin' }
  },
  {
    path: '/project/edit/:id',
    name: 'ProjectEdit',
    component: () => import('../views/project/ProjectEdit.vue'),
    meta: { requiresAuth: true, role: 'student' }
  },
  {
    path: '/student/competitions',
    name: 'StudentCompetitions',
    component: () => import('../views/student/CompetitionView.vue'),
    meta: { requiresAuth: true, role: 'student' }
  },
  {
    path: '/student/applications',
    name: 'StudentApplications',
    component: () => import('../views/student/ApplicationView.vue'),
    meta: { requiresAuth: true, role: 'student' }
  },
  {
    path: '/student/profile',
    name: 'StudentProfile',
    component: () => import('../views/student/ProfileView.vue'),
    meta: { requiresAuth: true, role: 'student' }
  },
  {
    path: '/student/notifications',
    name: 'StudentNotifications',
    component: () => import('../views/student/NotificationView.vue'),
    meta: { requiresAuth: true, role: 'student' }
  },
  {
    path: '/student/files',
    name: 'StudentFiles',
    component: () => import('../views/student/FileView.vue'),
    meta: { requiresAuth: true, role: 'student' }
  },
  // 新增竞赛管理模块路由
  {
    path: '/student/submission',
    name: 'StudentSubmission',
    component: () => import('../views/student/CompetitionSubmission.vue'),
    meta: { requiresAuth: true, role: 'student' }
  },
  {
    path: '/admin/competition-registrations',
    name: 'CompetitionRegistrationReview',
    component: () => import('../views/user/CompetitionRegistrationReview.vue'),
    meta: { requiresAuth: true, role: 'admin' }
  },
  {
    path: '/admin/competition-results',
    name: 'CompetitionResults',
    component: () => import('../views/user/CompetitionResults.vue'),
    meta: { requiresAuth: true, role: 'admin' }
  },
  {
    path: '/admin/audit-logs',
    name: 'CompetitionAuditLogs',
    component: () => import('../views/user/CompetitionAuditLogs.vue'),
    meta: { requiresAuth: true, role: 'admin' }
  },
  {
    path: '/teacher/judging',
    name: 'CompetitionJudging',
    component: () => import('../views/teacher/CompetitionJudging.vue'),
    meta: { requiresAuth: true, role: 'teacher' }
  },

]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 路由守卫 - 移除自动跳转到登录页面的功能
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('token')
  const userRole = localStorage.getItem('userRole')
  
  // 如果访问需要认证的页面但没有token，不自动跳转，让用户手动访问登录页面
  if (to.meta.requiresAuth && !token) {
    // 不自动跳转到登录页面，而是显示错误信息或保持当前页面
    console.log('需要登录才能访问此页面，请手动访问登录页面')
    // 可以选择显示一个提示信息，但不强制跳转
    next(false) // 阻止导航
  } else if (to.meta.role && to.meta.role !== userRole) {
    // 角色不匹配，也不自动跳转
    console.log('权限不足，请使用正确的角色登录')
    next(false) // 阻止导航
  } else {
    next()
  }
})

export default router 
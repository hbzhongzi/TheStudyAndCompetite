<template>
  <div class="login-container">
    <!-- 背景装饰 -->
    <div class="background-decoration">
      <div class="circle circle-1"></div>
      <div class="circle circle-2"></div>
      <div class="circle circle-3"></div>
    </div>
    
    <!-- 登录卡片 -->
    <div class="login-card">
      <!-- 标题区域 -->
      <div class="login-header">
        <div class="logo-container">
          <div class="logo-icon">
            <i class="el-icon-school"></i>
          </div>
        </div>
        <h1 class="system-title">云梦高校科研竞赛管理系统</h1>
        <p class="system-subtitle">Student Research & Competition Management System</p>
      </div>
      
      <!-- 登录表单 -->
      <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        class="login-form"
        @keydown.enter="handleLogin"
      >
        <!-- 用户名输入框 -->
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="请输入用户名/学号/工号"
            size="large"
            prefix-icon="el-icon-user"
            clearable
          />
        </el-form-item>
        
        <!-- 密码输入框 -->
        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="请输入密码"
            size="large"
            prefix-icon="el-icon-lock"
            show-password
            clearable
          />
        </el-form-item>
        
        <!-- 角色选择 -->
        <el-form-item prop="role">
          <el-select
            v-model="loginForm.role"
            placeholder="请选择角色"
            size="large"
            style="width: 100%"
          >
            <el-option label="学生" value="student" />
            <el-option label="教师" value="teacher" />
            <el-option label="管理员" value="admin" />
          </el-select>
        </el-form-item>
        
        <!-- 记住角色选项 -->
        <el-form-item>
          <el-checkbox v-model="rememberRole">记住角色选择</el-checkbox>
        </el-form-item>
        
        <!-- 登录按钮 -->
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            class="login-button"
            @click="handleLogin"
          >
            {{ loading ? '登录中...' : '登录' }}
          </el-button>
        </el-form-item>
        

      </el-form>
      
      <!-- 底部信息 -->
      <div class="login-footer">
        <p class="copyright">© 2024 云梦高校. All rights reserved.</p>
      </div>
    </div>
    

  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { login } from '../../services/auth'

const router = useRouter()
const loginFormRef = ref()
const loading = ref(false)
const rememberRole = ref(false)

const loginForm = reactive({
  username: 'admin', // 测试用默认值
  password: '123456', // 测试用默认值
  role: 'admin' // 测试用默认值
})

const loginRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  role: [
    { required: true, message: '请选择角色', trigger: 'change' }
  ]
}

const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  try {
    // 表单验证
    await loginFormRef.value.validate()
    
    loading.value = true
    
    // 调用登录接口
    const res = await login(loginForm)
    
    if (res.success) {
      // 登录成功
      localStorage.setItem('token', res.token)
      localStorage.setItem('userInfo', JSON.stringify(res.data))
      localStorage.setItem('userRole', loginForm.role)
      
      // 记住角色选择
      if (rememberRole.value) {
        localStorage.setItem('rememberedRole', loginForm.role)
      } else {
        localStorage.removeItem('rememberedRole')
      }
      
      ElMessage.success(res.message || '登录成功')
      
      // ✅ 添加调试信息，查看跳转前状态
      console.log('登录成功，开始跳转:', {
        role: loginForm.role,
        token: res.token,
        userInfo: res.data
      })
      
      // ✅ 等待一小段时间让消息提示显示
      setTimeout(async () => {
        try {
          // 根据角色跳转到对应页面
          let targetRoute = '/dashboard'
          switch (loginForm.role) {
            case 'student':
              targetRoute = '/student'
              break
            case 'teacher':
              targetRoute = '/teacher'
              break
            case 'admin':
              targetRoute = '/admin'
              break
            default:
              targetRoute = '/dashboard'
          }
          
          console.log('跳转到:', targetRoute)
          
          // ✅ 使用 replace 而不是 push，避免历史记录问题
          await router.replace(targetRoute)
          
          // ✅ 强制刷新页面以确保路由生效
          window.location.reload()
          
        } catch (routerError) {
          console.error('路由跳转失败:', routerError)
          ElMessage.error('页面跳转失败，请手动刷新页面')
        }
      }, 100) // 等待100ms
      
    } else {
      // 登录失败，显示后端返回的错误信息
      ElMessage.error(res.message || '登录失败')
    }
  } catch (error) {
    // 网络异常或其他错误
    console.error('登录异常:', error)
    ElMessage.error('网络异常或服务器无响应，请稍后重试')
  } finally {
    loading.value = false
  }
}

// 验证登录状态
const validateLoginState = () => {
  const token = localStorage.getItem('token')
  const userRole = localStorage.getItem('userRole')
  const userInfo = localStorage.getItem('userInfo')
  
  return { token, userRole, userInfo }
}

onMounted(() => {
  // 恢复记住的角色选择
  const rememberedRole = localStorage.getItem('rememberedRole')
  if (rememberedRole) {
    loginForm.role = rememberedRole
    rememberRole.value = true
  }
  
  // 移除自动登录功能 - 不再检查是否已经登录并自动跳转
  // 用户每次访问登录页面都需要重新登录
})
</script>

<style scoped>
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  position: relative;
  overflow: hidden;
}

/* 背景装饰 */
.background-decoration {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
}

.circle {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  animation: float 6s ease-in-out infinite;
}

.circle-1 {
  width: 200px;
  height: 200px;
  top: 10%;
  left: 10%;
  animation-delay: 0s;
}

.circle-2 {
  width: 150px;
  height: 150px;
  top: 60%;
  right: 15%;
  animation-delay: 2s;
}

.circle-3 {
  width: 100px;
  height: 100px;
  bottom: 20%;
  left: 20%;
  animation-delay: 4s;
}

@keyframes float {
  0%, 100% { transform: translateY(0px); }
  50% { transform: translateY(-20px); }
}

/* 登录卡片 */
.login-card {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 20px;
  padding: 40px;
  width: 100%;
  max-width: 420px;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

/* 标题区域 */
.login-header {
  text-align: center;
  margin-bottom: 30px;
}

.logo-container {
  margin-bottom: 20px;
}

.logo-icon {
  width: 80px;
  height: 80px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto;
  box-shadow: 0 10px 20px rgba(102, 126, 234, 0.3);
}

.logo-icon i {
  font-size: 40px;
  color: white;
}

.system-title {
  font-size: 24px;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 8px 0;
  font-family: 'Microsoft YaHei', 'PingFang SC', sans-serif;
}

.system-subtitle {
  font-size: 14px;
  color: #7f8c8d;
  margin: 0;
  font-family: 'Roboto', sans-serif;
}

/* 登录表单 */
.login-form {
  margin-bottom: 20px;
}

.login-form :deep(.el-input__wrapper) {
  border-radius: 10px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  border: 1px solid #e1e8ed;
}

.login-form :deep(.el-input__wrapper:hover) {
  border-color: #667eea;
}

.login-form :deep(.el-input__wrapper.is-focus) {
  border-color: #667eea;
  box-shadow: 0 0 0 2px rgba(102, 126, 234, 0.2);
}

.login-form :deep(.el-select .el-input__wrapper) {
  border-radius: 10px;
}

.login-button {
  width: 100%;
  height: 48px;
  border-radius: 10px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;
  font-size: 16px;
  font-weight: 600;
  transition: all 0.3s ease;
}

.login-button:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 20px rgba(102, 126, 234, 0.3);
}

.login-button:active {
  transform: translateY(0);
}

/* 底部信息 */
.login-footer {
  text-align: center;
  margin-top: 20px;
}

.copyright {
  font-size: 12px;
  color: #95a5a6;
  margin: 0;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .login-card {
    margin: 20px;
    padding: 30px 20px;
  }
  
  .system-title {
    font-size: 20px;
  }
  
  .logo-icon {
    width: 60px;
    height: 60px;
  }
  
  .logo-icon i {
    font-size: 30px;
  }
}

@media (max-width: 480px) {
  .login-card {
    padding: 20px 15px;
  }
  
  .system-title {
    font-size: 18px;
  }
  
  .system-subtitle {
    font-size: 12px;
  }
}


</style> 
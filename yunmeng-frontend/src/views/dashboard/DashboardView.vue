<template>
  <div class="dashboard-container">
    <div class="dashboard-header">
      <h1>欢迎使用云梦高校科研竞赛管理系统</h1>
      <p>请选择您的角色进入相应功能模块</p>
    </div>
    
    <div class="dashboard-content">
      <el-row :gutter="20">
        <el-col :span="8">
          <el-card class="role-card" @click="navigateToRole('student')">
            <div class="card-content">
              <div class="icon-wrapper student-icon">
                <i class="el-icon-user"></i>
              </div>
              <h3>学生端</h3>
              <p>查看项目、参与竞赛、提交申请</p>
            </div>
          </el-card>
        </el-col>
        
        <el-col :span="8">
          <el-card class="role-card" @click="navigateToRole('teacher')">
            <div class="card-content">
              <div class="icon-wrapper teacher-icon">
                <i class="el-icon-s-custom"></i>
              </div>
              <h3>教师端</h3>
              <p>管理项目、指导竞赛、审核申请</p>
            </div>
          </el-card>
        </el-col>
        
        <el-col :span="8">
          <el-card class="role-card" @click="navigateToRole('admin')">
            <div class="card-content">
              <div class="icon-wrapper admin-icon">
                <i class="el-icon-setting"></i>
              </div>
              <h3>管理员端</h3>
              <p>系统管理、用户管理、数据统计</p>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()

const navigateToRole = (role) => {
  router.push(`/${role}`)
}
</script>

<style scoped>
.dashboard-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  padding: 40px 20px;
}

.dashboard-header {
  text-align: center;
  margin-bottom: 60px;
}

.dashboard-header h1 {
  font-size: 32px;
  color: #2c3e50;
  margin-bottom: 10px;
  font-weight: 600;
}

.dashboard-header p {
  font-size: 16px;
  color: #7f8c8d;
  margin: 0;
}

.dashboard-content {
  max-width: 1200px;
  margin: 0 auto;
}

.role-card {
  cursor: pointer;
  transition: all 0.3s ease;
  border-radius: 15px;
  overflow: hidden;
  height: 300px;
  display: flex;
  align-items: center;
}

.role-card:hover {
  transform: translateY(-10px);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
}

.card-content {
  text-align: center;
  padding: 20px;
  width: 100%;
}

.icon-wrapper {
  width: 80px;
  height: 80px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 20px;
  transition: all 0.3s ease;
}

.icon-wrapper i {
  font-size: 40px;
  color: white;
}

.student-icon {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.teacher-icon {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.admin-icon {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.role-card:hover .icon-wrapper {
  transform: scale(1.1);
}

.card-content h3 {
  font-size: 24px;
  color: #2c3e50;
  margin-bottom: 15px;
  font-weight: 600;
}

.card-content p {
  font-size: 14px;
  color: #7f8c8d;
  line-height: 1.6;
  margin: 0;
}

@media (max-width: 768px) {
  .dashboard-header h1 {
    font-size: 24px;
  }
  
  .dashboard-content {
    padding: 0 10px;
  }
  
  .role-card {
    height: 250px;
    margin-bottom: 20px;
  }
}
</style> 
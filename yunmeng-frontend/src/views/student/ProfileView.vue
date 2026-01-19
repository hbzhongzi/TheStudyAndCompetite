<template>
  <div class="profile-view">
    <!-- 页面标题 -->
    <div class="page-header">
      <h2>个人信息</h2>
      <p>管理您的个人资料和学术成果</p>
    </div>

    <el-row :gutter="20">
      <!-- 左侧：基本信息 -->
      <el-col :span="16">
        <el-card class="profile-card">
          <template #header>
            <div class="card-header">
              <span>基本信息</span>
              <el-button 
                v-if="!isEditing" 
                type="primary" 
                size="small"
                @click="startEdit"
              >
                编辑信息
              </el-button>
              <div v-else>
                <el-button size="small" @click="cancelEdit">取消</el-button>
                <el-button type="primary" size="small" @click="saveProfile" :loading="saving">
                  保存
                </el-button>
              </div>
            </div>
          </template>
          
          <el-form 
            :model="profileForm" 
            :rules="profileRules"
            ref="profileFormRef"
            label-width="100px"
            :disabled="!isEditing"
          >
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="姓名" prop="name">
                  <el-input v-model="profileForm.name" placeholder="请输入姓名" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="学号" prop="studentNumber">
                  <el-input v-model="profileForm.studentNumber" placeholder="请输入学号" disabled />
                </el-form-item>
              </el-col>
            </el-row>
            
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="性别" prop="gender">
                  <el-select v-model="profileForm.gender" placeholder="请选择性别">
                    <el-option label="男" value="male" />
                    <el-option label="女" value="female" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="出生日期" prop="birthDate">
                  <el-date-picker
                    v-model="profileForm.birthDate"
                    type="date"
                    placeholder="请选择出生日期"
                    style="width: 100%"
                  />
                </el-form-item>
              </el-col>
            </el-row>
            
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="学院" prop="college">
                  <el-input v-model="profileForm.college" placeholder="请输入学院" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="专业" prop="major">
                  <el-input v-model="profileForm.major" placeholder="请输入专业" />
                </el-form-item>
              </el-col>
            </el-row>
            
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="年级" prop="grade">
                  <el-select v-model="profileForm.grade" placeholder="请选择年级">
                    <el-option label="大一" value="1" />
                    <el-option label="大二" value="2" />
                    <el-option label="大三" value="3" />
                    <el-option label="大四" value="4" />
                    <el-option label="研究生" value="graduate" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="班级" prop="className">
                  <el-input v-model="profileForm.className" placeholder="请输入班级" />
                </el-form-item>
              </el-col>
            </el-row>
            
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="联系电话" prop="phone">
                  <el-input v-model="profileForm.phone" placeholder="请输入联系电话" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="邮箱" prop="email">
                  <el-input v-model="profileForm.email" placeholder="请输入邮箱" />
                </el-form-item>
              </el-col>
            </el-row>
            
            <el-form-item label="个人简介" prop="bio">
              <el-input 
                v-model="profileForm.bio" 
                type="textarea" 
                placeholder="请输入个人简介"
                :rows="4"
              />
            </el-form-item>
          </el-form>
        </el-card>

        <!-- 学术成果 -->
        <el-card class="achievement-card" style="margin-top: 20px;">
          <template #header>
            <div class="card-header">
              <span>学术成果</span>
              <el-button type="primary" size="small" @click="showAddAchievement = true">
                添加成果
              </el-button>
            </div>
          </template>
          
          <div class="achievement-list">
            <el-empty v-if="achievements.length === 0" description="暂无学术成果" />
            <div v-else class="achievement-items">
              <div 
                v-for="achievement in achievements" 
                :key="achievement.id"
                class="achievement-item"
              >
                <div class="achievement-header">
                  <h4>{{ achievement.title }}</h4>
                  <div class="achievement-actions">
                    <el-button size="small" @click="editAchievement(achievement)">编辑</el-button>
                    <el-button size="small" type="danger" @click="deleteAchievement(achievement)">删除</el-button>
                  </div>
                </div>
                <div class="achievement-content">
                  <p><strong>类型：</strong>{{ achievement.type }}</p>
                  <p><strong>时间：</strong>{{ formatDate(achievement.date) }}</p>
                  <p><strong>描述：</strong>{{ achievement.description }}</p>
                  <p v-if="achievement.award"><strong>奖项：</strong>{{ achievement.award }}</p>
                </div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>

      <!-- 右侧：统计信息和快捷操作 -->
      <el-col :span="8">
        <!-- 统计信息 -->
        <el-card class="stats-card">
          <template #header>
            <span>统计信息</span>
          </template>
          
          <div class="stats-content">
            <div class="stat-item">
              <div class="stat-icon">
                <el-icon><Document /></el-icon>
              </div>
              <div class="stat-info">
                <h4>参与项目</h4>
                <p>{{ stats.projectCount }}</p>
              </div>
            </div>
            
            <div class="stat-item">
              <div class="stat-icon">
                <el-icon><Trophy /></el-icon>
              </div>
              <div class="stat-info">
                <h4>参与竞赛</h4>
                <p>{{ stats.competitionCount }}</p>
              </div>
            </div>
            
            <div class="stat-item">
              <div class="stat-icon">
                <el-icon><Star /></el-icon>
              </div>
              <div class="stat-info">
                <h4>学术成果</h4>
                <p>{{ stats.achievementCount }}</p>
              </div>
            </div>
            
            <div class="stat-item">
              <div class="stat-icon">
                <el-icon><Trophy /></el-icon>
              </div>
              <div class="stat-info">
                <h4>获得奖项</h4>
                <p>{{ stats.awardCount }}</p>
              </div>
            </div>
          </div>
        </el-card>

        <!-- 快捷操作 -->
        <el-card class="quick-actions-card" style="margin-top: 20px;">
          <template #header>
            <span>快捷操作</span>
          </template>
          
          <div class="quick-actions">
            <el-button type="primary" @click="changePassword" style="width: 100%; margin-bottom: 10px;">
              修改密码
            </el-button>
            <el-button @click="exportProfile" style="width: 100%; margin-bottom: 10px;">
              导出个人档案
            </el-button>
            <el-button @click="uploadAvatar" style="width: 100%;">
              上传头像
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <!-- 添加/编辑学术成果对话框 -->
    <el-dialog
      v-model="showAddAchievement"
      :title="editingAchievement ? '编辑学术成果' : '添加学术成果'"
      width="50%"
    >
      <el-form :model="achievementForm" :rules="achievementRules" ref="achievementFormRef" label-width="100px">
        <el-form-item label="成果标题" prop="title">
          <el-input v-model="achievementForm.title" placeholder="请输入成果标题" />
        </el-form-item>
        
        <el-form-item label="成果类型" prop="type">
          <el-select v-model="achievementForm.type" placeholder="请选择成果类型" style="width: 100%">
            <el-option label="论文发表" value="paper" />
            <el-option label="专利" value="patent" />
            <el-option label="竞赛获奖" value="competition" />
            <el-option label="项目成果" value="project" />
            <el-option label="其他" value="other" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="获得时间" prop="date">
          <el-date-picker
            v-model="achievementForm.date"
            type="date"
            placeholder="请选择获得时间"
            style="width: 100%"
          />
        </el-form-item>
        
        <el-form-item label="成果描述" prop="description">
          <el-input 
            v-model="achievementForm.description" 
            type="textarea" 
            placeholder="请输入成果描述"
            :rows="3"
          />
        </el-form-item>
        
        <el-form-item label="获得奖项" prop="award">
          <el-input v-model="achievementForm.award" placeholder="请输入获得奖项（可选）" />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showAddAchievement = false">取消</el-button>
          <el-button type="primary" @click="saveAchievement" :loading="savingAchievement">
            保存
          </el-button>
        </span>
      </template>
    </el-dialog>

    <!-- 修改密码对话框 -->
    <el-dialog
      v-model="showPasswordDialog"
      title="修改密码"
      width="40%"
    >
      <el-form :model="passwordForm" :rules="passwordRules" ref="passwordFormRef" label-width="100px">
        <el-form-item label="原密码" prop="oldPassword">
          <el-input 
            v-model="passwordForm.oldPassword" 
            type="password" 
            placeholder="请输入原密码"
            show-password
          />
        </el-form-item>
        
        <el-form-item label="新密码" prop="newPassword">
          <el-input 
            v-model="passwordForm.newPassword" 
            type="password" 
            placeholder="请输入新密码"
            show-password
          />
        </el-form-item>
        
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input 
            v-model="passwordForm.confirmPassword" 
            type="password" 
            placeholder="请再次输入新密码"
            show-password
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showPasswordDialog = false">取消</el-button>
          <el-button type="primary" @click="confirmChangePassword" :loading="changingPassword">
            确认修改
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Document, Trophy, Star } from '@element-plus/icons-vue'

// 响应式数据
const isEditing = ref(false)
const saving = ref(false)
const showAddAchievement = ref(false)
const showPasswordDialog = ref(false)
const editingAchievement = ref(null)
const savingAchievement = ref(false)
const changingPassword = ref(false)

// 表单引用
const profileFormRef = ref()
const achievementFormRef = ref()
const passwordFormRef = ref()

// 个人信息表单
const profileForm = reactive({
  name: '张三',
  studentNumber: '2021001001',
  gender: 'male',
  birthDate: '2000-01-01',
  college: '计算机学院',
  major: '软件工程',
  grade: '3',
  className: '软工2101班',
  phone: '13800138000',
  email: 'zhangsan@example.com',
  bio: '热爱编程，专注于Web开发和人工智能技术研究。'
})

// 个人信息验证规则
const profileRules = {
  name: [
    { required: true, message: '请输入姓名', trigger: 'blur' }
  ],
  phone: [
    { required: true, message: '请输入联系电话', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号码', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ]
}

// 学术成果列表
const achievements = ref([
  {
    id: 1,
    title: '基于深度学习的图像识别系统',
    type: '项目成果',
    date: '2024-01-15',
    description: '开发了一个基于深度学习的图像识别系统，用于医疗影像分析',
    award: '校级优秀项目奖'
  },
  {
    id: 2,
    title: '全国大学生程序设计竞赛',
    type: '竞赛获奖',
    date: '2023-12-01',
    description: '参加全国大学生程序设计竞赛，获得省级二等奖',
    award: '省级二等奖'
  }
])

// 学术成果表单
const achievementForm = reactive({
  title: '',
  type: '',
  date: '',
  description: '',
  award: ''
})

// 学术成果验证规则
const achievementRules = {
  title: [
    { required: true, message: '请输入成果标题', trigger: 'blur' }
  ],
  type: [
    { required: true, message: '请选择成果类型', trigger: 'change' }
  ],
  date: [
    { required: true, message: '请选择获得时间', trigger: 'change' }
  ],
  description: [
    { required: true, message: '请输入成果描述', trigger: 'blur' }
  ]
}

// 密码表单
const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 密码验证规则
const passwordRules = {
  oldPassword: [
    { required: true, message: '请输入原密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== passwordForm.newPassword) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

// 统计数据
const stats = ref({
  projectCount: 3,
  competitionCount: 2,
  achievementCount: 2,
  awardCount: 1
})

// 开始编辑
const startEdit = () => {
  isEditing.value = true
}

// 取消编辑
const cancelEdit = () => {
  isEditing.value = false
  // 重置表单数据
  Object.assign(profileForm, {
    name: '张三',
    studentNumber: '2021001001',
    gender: 'male',
    birthDate: '2000-01-01',
    college: '计算机学院',
    major: '软件工程',
    grade: '3',
    className: '软工2101班',
    phone: '13800138000',
    email: 'zhangsan@example.com',
    bio: '热爱编程，专注于Web开发和人工智能技术研究。'
  })
}

// 保存个人信息
const saveProfile = async () => {
  try {
    await profileFormRef.value.validate()
    saving.value = true
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    isEditing.value = false
    ElMessage.success('个人信息保存成功')
  } catch (error) {
    console.error('保存失败:', error)
    ElMessage.error('保存失败，请检查输入信息')
  } finally {
    saving.value = false
  }
}

// 编辑学术成果
const editAchievement = (achievement) => {
  editingAchievement.value = achievement
  Object.assign(achievementForm, {
    title: achievement.title,
    type: achievement.type,
    date: achievement.date,
    description: achievement.description,
    award: achievement.award
  })
  showAddAchievement.value = true
}

// 删除学术成果
const deleteAchievement = async (achievement) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除学术成果"${achievement.title}"吗？`,
      '确认删除',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    const index = achievements.value.findIndex(a => a.id === achievement.id)
    if (index !== -1) {
      achievements.value.splice(index, 1)
      ElMessage.success('删除成功')
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

// 保存学术成果
const saveAchievement = async () => {
  try {
    await achievementFormRef.value.validate()
    savingAchievement.value = true
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 800))
    
    if (editingAchievement.value) {
      // 编辑现有成果
      Object.assign(editingAchievement.value, {
        title: achievementForm.title,
        type: achievementForm.type,
        date: achievementForm.date,
        description: achievementForm.description,
        award: achievementForm.award
      })
    } else {
      // 添加新成果
      const newAchievement = {
        id: Date.now(),
        title: achievementForm.title,
        type: achievementForm.type,
        date: achievementForm.date,
        description: achievementForm.description,
        award: achievementForm.award
      }
      achievements.value.push(newAchievement)
    }
    
    showAddAchievement.value = false
    editingAchievement.value = null
    Object.assign(achievementForm, {
      title: '',
      type: '',
      date: '',
      description: '',
      award: ''
    })
    ElMessage.success(editingAchievement.value ? '编辑成功' : '添加成功')
  } catch (error) {
    console.error('保存失败:', error)
    ElMessage.error('保存失败，请检查输入信息')
  } finally {
    savingAchievement.value = false
  }
}

// 修改密码
const changePassword = () => {
  Object.assign(passwordForm, {
    oldPassword: '',
    newPassword: '',
    confirmPassword: ''
  })
  showPasswordDialog.value = true
}

// 确认修改密码
const confirmChangePassword = async () => {
  try {
    await passwordFormRef.value.validate()
    changingPassword.value = true
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    showPasswordDialog.value = false
    ElMessage.success('密码修改成功')
  } catch (error) {
    console.error('修改失败:', error)
    ElMessage.error('修改失败，请检查输入信息')
  } finally {
    changingPassword.value = false
  }
}

// 导出个人档案
const exportProfile = () => {
  ElMessage.info('导出功能开发中...')
}

// 上传头像
const uploadAvatar = () => {
  ElMessage.info('上传头像功能开发中...')
}

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN')
}

// 监听对话框关闭
const handleDialogClose = () => {
  editingAchievement.value = null
  Object.assign(achievementForm, {
    title: '',
    type: '',
    date: '',
    description: '',
    award: ''
  })
}

// 组件挂载时加载数据
onMounted(() => {
  // 这里可以加载用户数据
})
</script>

<style scoped>
.profile-view {
  padding: 20px;
}

.page-header {
  margin-bottom: 30px;
  text-align: center;
}

.page-header h2 {
  margin: 0 0 10px 0;
  color: #2c3e50;
  font-size: 28px;
  font-weight: 600;
}

.page-header p {
  margin: 0;
  color: #7f8c8d;
  font-size: 16px;
}

.profile-card,
.achievement-card,
.stats-card,
.quick-actions-card {
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.achievement-list {
  margin-top: 10px;
}

.achievement-items {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.achievement-item {
  border: 1px solid #e9ecef;
  border-radius: 8px;
  padding: 15px;
  background: #f8f9fa;
}

.achievement-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.achievement-header h4 {
  margin: 0;
  color: #2c3e50;
  font-size: 16px;
  font-weight: 600;
}

.achievement-actions {
  display: flex;
  gap: 8px;
}

.achievement-content p {
  margin: 0 0 5px 0;
  color: #5a6c7d;
  font-size: 14px;
}

.stats-content {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 15px;
  padding: 15px;
  background: #f8f9fa;
  border-radius: 8px;
}

.stat-icon {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-icon .el-icon {
  font-size: 20px;
  color: white;
}

.stat-info h4 {
  margin: 0 0 5px 0;
  color: #7f8c8d;
  font-size: 14px;
}

.stat-info p {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
  color: #2c3e50;
}

.quick-actions {
  display: flex;
  flex-direction: column;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style> 
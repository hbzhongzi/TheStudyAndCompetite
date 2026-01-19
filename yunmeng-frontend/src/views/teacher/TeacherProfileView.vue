<template>
  <div class="teacher-profile">
    <el-row :gutter="20">
      <!-- 个人信息卡片 -->
      <el-col :span="8">
        <el-card class="profile-card">
          <div class="profile-header">
            <div class="avatar-container">
              <el-avatar :size="100" :src="profile.avatar" />
              <el-button 
                size="small" 
                type="primary" 
                class="upload-btn"
                @click="handleAvatarUpload"
              >
                更换头像
              </el-button>
            </div>
            <h3>{{ profile.name }}</h3>
            <p class="title">{{ profile.title }}</p>
            <p class="department">{{ profile.department }}</p>
          </div>
          
          <el-divider />
          
          <div class="profile-info">
            <div class="info-item">
              <span class="label">工号：</span>
              <span class="value">{{ profile.employeeId }}</span>
            </div>
            <div class="info-item">
              <span class="label">邮箱：</span>
              <span class="value">{{ profile.email }}</span>
            </div>
            <div class="info-item">
              <span class="label">电话：</span>
              <span class="value">{{ profile.phone }}</span>
            </div>
            <div class="info-item">
              <span class="label">办公室：</span>
              <span class="value">{{ profile.office }}</span>
            </div>
            <div class="info-item">
              <span class="label">入职时间：</span>
              <span class="value">{{ profile.hireDate }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <!-- 详细信息 -->
      <el-col :span="16">
        <el-card class="detail-card">
          <template #header>
            <div class="card-header">
              <span>详细信息</span>
              <el-button type="primary" @click="editMode = true" v-if="!editMode">
                编辑信息
              </el-button>
              <div v-else>
                <el-button @click="cancelEdit">取消</el-button>
                <el-button type="primary" @click="saveProfile" :loading="saving">保存</el-button>
              </div>
            </div>
          </template>
          
          <el-form 
            :model="profile" 
            :disabled="!editMode"
            label-width="100px"
            class="profile-form"
          >
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="姓名">
                  <el-input v-model="profile.name" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="工号">
                  <el-input v-model="profile.employeeId" disabled />
                </el-form-item>
              </el-col>
            </el-row>
            
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="职称">
                  <el-select v-model="profile.title" placeholder="请选择职称">
                    <el-option label="教授" value="教授" />
                    <el-option label="副教授" value="副教授" />
                    <el-option label="讲师" value="讲师" />
                    <el-option label="助教" value="助教" />
                  </el-select>
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="所属院系">
                  <el-input v-model="profile.department" />
                </el-form-item>
              </el-col>
            </el-row>
            
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="邮箱">
                  <el-input v-model="profile.email" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="电话">
                  <el-input v-model="profile.phone" />
                </el-form-item>
              </el-col>
            </el-row>
            
            <el-row :gutter="20">
              <el-col :span="12">
                <el-form-item label="办公室">
                  <el-input v-model="profile.office" />
                </el-form-item>
              </el-col>
              <el-col :span="12">
                <el-form-item label="入职时间">
                  <el-date-picker
                    v-model="profile.hireDate"
                    type="date"
                    placeholder="选择入职时间"
                    format="YYYY-MM-DD"
                    value-format="YYYY-MM-DD"
                  />
                </el-form-item>
              </el-col>
            </el-row>
            
            <el-form-item label="个人简介">
              <el-input
                v-model="profile.bio"
                type="textarea"
                :rows="4"
                placeholder="请输入个人简介"
              />
            </el-form-item>
            
            <el-form-item label="研究方向">
              <el-input
                v-model="profile.researchAreas"
                type="textarea"
                :rows="3"
                placeholder="请输入研究方向"
              />
            </el-form-item>
          </el-form>
        </el-card>
        
        <!-- 修改密码 -->
        <el-card class="password-card">
          <template #header>
            <span>修改密码</span>
          </template>
          
          <el-form :model="passwordForm" label-width="100px" class="password-form">
            <el-form-item label="当前密码">
              <el-input 
                v-model="passwordForm.oldPassword" 
                type="password" 
                placeholder="请输入当前密码"
                show-password
              />
            </el-form-item>
            <el-form-item label="新密码">
              <el-input 
                v-model="passwordForm.newPassword" 
                type="password" 
                placeholder="请输入新密码"
                show-password
              />
            </el-form-item>
            <el-form-item label="确认密码">
              <el-input 
                v-model="passwordForm.confirmPassword" 
                type="password" 
                placeholder="请再次输入新密码"
                show-password
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="changePassword" :loading="changingPassword">
                修改密码
              </el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
    
    <!-- 头像上传对话框 -->
    <el-dialog
      v-model="avatarDialogVisible"
      title="上传头像"
      width="400px"
    >
      <el-upload
        class="avatar-uploader"
        :show-file-list="false"
        :before-upload="beforeAvatarUpload"
        :http-request="uploadAvatar"
      >
        <img v-if="avatarUrl" :src="avatarUrl" class="avatar" />
        <el-icon v-else class="avatar-uploader-icon"><Plus /></el-icon>
      </el-upload>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="avatarDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="confirmAvatarUpload" :loading="uploading">
            确认上传
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'

// 响应式数据
const editMode = ref(false)
const saving = ref(false)
const changingPassword = ref(false)
const avatarDialogVisible = ref(false)
const avatarUrl = ref('')
const uploading = ref(false)

// 个人信息
const profile = reactive({
  name: '李教授',
  employeeId: 'T2021001',
  title: '教授',
  department: '计算机科学与技术学院',
  email: 'li.professor@university.edu.cn',
  phone: '13800138000',
  office: 'A栋301室',
  hireDate: '2015-09-01',
  bio: '从事计算机科学教学和研究工作多年，主要研究方向包括人工智能、机器学习、数据挖掘等。在国内外学术期刊发表论文50余篇，主持国家级科研项目3项。',
  researchAreas: '人工智能、机器学习、数据挖掘、计算机视觉',
  avatar: ''
})

// 密码表单
const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 方法
const handleAvatarUpload = () => {
  avatarDialogVisible.value = true
}

const beforeAvatarUpload = (file) => {
  const isJPG = file.type === 'image/jpeg' || file.type === 'image/png'
  const isLt2M = file.size / 1024 / 1024 < 2

  if (!isJPG) {
    ElMessage.error('头像只能是 JPG 或 PNG 格式!')
    return false
  }
  if (!isLt2M) {
    ElMessage.error('头像大小不能超过 2MB!')
    return false
  }
  
  // 预览图片
  const reader = new FileReader()
  reader.onload = (e) => {
    avatarUrl.value = e.target.result
  }
  reader.readAsDataURL(file)
  
  return false // 阻止自动上传
}

const uploadAvatar = (options) => {
  // 这里应该调用实际的上传API
  console.log('上传文件:', options.file)
}

const confirmAvatarUpload = async () => {
  if (!avatarUrl.value) {
    ElMessage.warning('请先选择头像文件')
    return
  }
  
  uploading.value = true
  try {
    // 模拟上传
    await new Promise(resolve => setTimeout(resolve, 1000))
    profile.avatar = avatarUrl.value
    ElMessage.success('头像上传成功')
    avatarDialogVisible.value = false
  } catch (error) {
    ElMessage.error('头像上传失败')
  } finally {
    uploading.value = false
  }
}

const saveProfile = async () => {
  saving.value = true
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    editMode.value = false
    ElMessage.success('个人信息保存成功')
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

const cancelEdit = () => {
  editMode.value = false
  // 这里可以重置表单数据
}

const changePassword = async () => {
  if (!passwordForm.oldPassword || !passwordForm.newPassword || !passwordForm.confirmPassword) {
    ElMessage.warning('请填写完整的密码信息')
    return
  }
  
  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    ElMessage.error('两次输入的新密码不一致')
    return
  }
  
  if (passwordForm.newPassword.length < 6) {
    ElMessage.error('新密码长度不能少于6位')
    return
  }
  
  changingPassword.value = true
  try {
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // 清空表单
    passwordForm.oldPassword = ''
    passwordForm.newPassword = ''
    passwordForm.confirmPassword = ''
    
    ElMessage.success('密码修改成功')
  } catch (error) {
    ElMessage.error('密码修改失败')
  } finally {
    changingPassword.value = false
  }
}
</script>

<style scoped>
.teacher-profile {
  padding: 20px;
}

.profile-card {
  border-radius: 10px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.profile-header {
  text-align: center;
  padding: 20px 0;
}

.avatar-container {
  position: relative;
  display: inline-block;
  margin-bottom: 15px;
}

.upload-btn {
  position: absolute;
  bottom: 0;
  right: 0;
  transform: translate(50%, 50%);
}

.profile-header h3 {
  margin: 10px 0 5px 0;
  color: #2c3e50;
  font-size: 20px;
}

.title {
  margin: 5px 0;
  color: #7f8c8d;
  font-size: 14px;
}

.department {
  margin: 5px 0;
  color: #95a5a6;
  font-size: 12px;
}

.profile-info {
  padding: 0 20px;
}

.info-item {
  display: flex;
  margin-bottom: 10px;
}

.info-item .label {
  font-weight: 600;
  color: #7f8c8d;
  width: 80px;
}

.info-item .value {
  color: #2c3e50;
  flex: 1;
}

.detail-card {
  border-radius: 10px;
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.profile-form {
  padding: 20px 0;
}

.password-card {
  border-radius: 10px;
}

.password-form {
  padding: 20px 0;
}

.avatar-uploader {
  text-align: center;
}

.avatar-uploader .el-upload {
  border: 1px dashed #d9d9d9;
  border-radius: 6px;
  cursor: pointer;
  position: relative;
  overflow: hidden;
  transition: var(--el-transition-duration-fast);
}

.avatar-uploader .el-upload:hover {
  border-color: var(--el-color-primary);
}

.avatar-uploader-icon {
  font-size: 28px;
  color: #8c939d;
  width: 178px;
  height: 178px;
  text-align: center;
  line-height: 178px;
}

.avatar {
  width: 178px;
  height: 178px;
  display: block;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style> 
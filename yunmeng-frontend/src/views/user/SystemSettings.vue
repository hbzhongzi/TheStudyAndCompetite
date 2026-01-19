<template>
  <div class="system-settings">
    <div class="page-header">
      <h2>系统设置</h2>
    </div>

    <el-tabs v-model="activeTab" type="border-card">
      <!-- 基本设置 -->
      <el-tab-pane label="基本设置" name="basic">
        <el-card class="setting-card">
          <template #header>
            <span>系统基本信息</span>
          </template>
          <el-form :model="basicSettings" :rules="basicRules" ref="basicFormRef" label-width="120px">
            <el-form-item label="系统名称" prop="systemName">
              <el-input v-model="basicSettings.systemName" placeholder="请输入系统名称" />
            </el-form-item>
            <el-form-item label="系统版本" prop="version">
              <el-input v-model="basicSettings.version" placeholder="请输入系统版本" />
            </el-form-item>
            <el-form-item label="管理员邮箱" prop="adminEmail">
              <el-input v-model="basicSettings.adminEmail" placeholder="请输入管理员邮箱" />
            </el-form-item>
            <el-form-item label="系统描述" prop="description">
              <el-input 
                v-model="basicSettings.description" 
                type="textarea" 
                :rows="3"
                placeholder="请输入系统描述" 
              />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveBasicSettings" :loading="saving">保存设置</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-tab-pane>

      <!-- 权限管理 -->
      <el-tab-pane label="权限管理" name="permissions">
        <el-card class="setting-card">
          <template #header>
            <span>角色权限配置</span>
          </template>
          <el-table :data="rolePermissions" style="width: 100%">
            <el-table-column prop="role" label="角色" width="120" />
            <el-table-column prop="description" label="描述" width="200" />
            <el-table-column label="权限" min-width="400">
              <template #default="scope">
                <el-checkbox-group v-model="scope.row.permissions">
                  <el-checkbox label="user_manage">用户管理</el-checkbox>
                  <el-checkbox label="project_manage">项目管理</el-checkbox>
                  <el-checkbox label="competition_manage">竞赛管理</el-checkbox>
                  <el-checkbox label="system_config">系统配置</el-checkbox>
                  <el-checkbox label="data_export">数据导出</el-checkbox>
                  <el-checkbox label="log_view">日志查看</el-checkbox>
                </el-checkbox-group>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="120">
              <template #default="scope">
                <el-button size="small" @click="saveRolePermissions(scope.row)">保存</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-tab-pane>

      <!-- 安全设置 -->
      <el-tab-pane label="安全设置" name="security">
        <el-card class="setting-card">
          <template #header>
            <span>安全配置</span>
          </template>
          <el-form :model="securitySettings" :rules="securityRules" ref="securityFormRef" label-width="150px">
            <el-form-item label="密码最小长度" prop="minPasswordLength">
              <el-input-number 
                v-model="securitySettings.minPasswordLength" 
                :min="6" 
                :max="20"
                style="width: 200px"
              />
            </el-form-item>
            <el-form-item label="密码复杂度要求" prop="passwordComplexity">
              <el-checkbox-group v-model="securitySettings.passwordComplexity">
                <el-checkbox label="uppercase">必须包含大写字母</el-checkbox>
                <el-checkbox label="lowercase">必须包含小写字母</el-checkbox>
                <el-checkbox label="number">必须包含数字</el-checkbox>
                <el-checkbox label="special">必须包含特殊字符</el-checkbox>
              </el-checkbox-group>
            </el-form-item>
            <el-form-item label="登录失败锁定" prop="loginLockEnabled">
              <el-switch v-model="securitySettings.loginLockEnabled" />
            </el-form-item>
            <el-form-item label="失败次数阈值" prop="maxLoginAttempts" v-if="securitySettings.loginLockEnabled">
              <el-input-number 
                v-model="securitySettings.maxLoginAttempts" 
                :min="3" 
                :max="10"
                style="width: 200px"
              />
            </el-form-item>
            <el-form-item label="锁定时间(分钟)" prop="lockDuration" v-if="securitySettings.loginLockEnabled">
              <el-input-number 
                v-model="securitySettings.lockDuration" 
                :min="5" 
                :max="1440"
                style="width: 200px"
              />
            </el-form-item>
            <el-form-item label="会话超时时间" prop="sessionTimeout">
              <el-input-number 
                v-model="securitySettings.sessionTimeout" 
                :min="15" 
                :max="480"
                style="width: 200px"
              />
              <span style="margin-left: 10px;">分钟</span>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveSecuritySettings" :loading="saving">保存设置</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-tab-pane>

      <!-- 邮件设置 -->
      <el-tab-pane label="邮件设置" name="email">
        <el-card class="setting-card">
          <template #header>
            <span>邮件服务配置</span>
          </template>
          <el-form :model="emailSettings" :rules="emailRules" ref="emailFormRef" label-width="120px">
            <el-form-item label="SMTP服务器" prop="smtpServer">
              <el-input v-model="emailSettings.smtpServer" placeholder="如：smtp.qq.com" />
            </el-form-item>
            <el-form-item label="SMTP端口" prop="smtpPort">
              <el-input-number 
                v-model="emailSettings.smtpPort" 
                :min="1" 
                :max="65535"
                style="width: 200px"
              />
            </el-form-item>
            <el-form-item label="发件人邮箱" prop="senderEmail">
              <el-input v-model="emailSettings.senderEmail" placeholder="请输入发件人邮箱" />
            </el-form-item>
            <el-form-item label="发件人名称" prop="senderName">
              <el-input v-model="emailSettings.senderName" placeholder="请输入发件人名称" />
            </el-form-item>
            <el-form-item label="邮箱密码" prop="emailPassword">
              <el-input 
                v-model="emailSettings.emailPassword" 
                type="password" 
                placeholder="请输入邮箱密码或授权码"
                show-password
              />
            </el-form-item>
            <el-form-item label="启用SSL" prop="enableSSL">
              <el-switch v-model="emailSettings.enableSSL" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveEmailSettings" :loading="saving">保存设置</el-button>
              <el-button @click="testEmailSettings" :loading="testing">测试连接</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-tab-pane>

      <!-- 备份设置 -->
      <el-tab-pane label="备份设置" name="backup">
        <el-card class="setting-card">
          <template #header>
            <span>数据备份配置</span>
          </template>
          <el-form :model="backupSettings" :rules="backupRules" ref="backupFormRef" label-width="150px">
            <el-form-item label="自动备份" prop="autoBackup">
              <el-switch v-model="backupSettings.autoBackup" />
            </el-form-item>
            <el-form-item label="备份频率" prop="backupFrequency" v-if="backupSettings.autoBackup">
              <el-select v-model="backupSettings.backupFrequency" style="width: 200px">
                <el-option label="每天" value="daily" />
                <el-option label="每周" value="weekly" />
                <el-option label="每月" value="monthly" />
              </el-select>
            </el-form-item>
            <el-form-item label="备份时间" prop="backupTime" v-if="backupSettings.autoBackup">
              <el-time-picker 
                v-model="backupSettings.backupTime" 
                format="HH:mm"
                style="width: 200px"
              />
            </el-form-item>
            <el-form-item label="保留备份数量" prop="backupRetention">
              <el-input-number 
                v-model="backupSettings.backupRetention" 
                :min="1" 
                :max="100"
                style="width: 200px"
              />
            </el-form-item>
            <el-form-item label="备份路径" prop="backupPath">
              <el-input v-model="backupSettings.backupPath" placeholder="请输入备份文件保存路径" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" @click="saveBackupSettings" :loading="saving">保存设置</el-button>
              <el-button type="success" @click="createBackup" :loading="backingUp">立即备份</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'

// 响应式数据
const activeTab = ref('basic')
const saving = ref(false)
const testing = ref(false)
const backingUp = ref(false)

// 表单引用
const basicFormRef = ref()
const securityFormRef = ref()
const emailFormRef = ref()
const backupFormRef = ref()

// 基本设置
const basicSettings = reactive({
  systemName: '云梦高校科研竞赛管理系统',
  version: 'v1.0.0',
  adminEmail: 'admin@yunmeng.edu.cn',
  description: '云梦高校科研竞赛管理系统是一个综合性的科研项目管理平台，支持项目申报、竞赛管理、成果展示等功能。'
})

const basicRules = {
  systemName: [
    { required: true, message: '请输入系统名称', trigger: 'blur' }
  ],
  version: [
    { required: true, message: '请输入系统版本', trigger: 'blur' }
  ],
  adminEmail: [
    { required: true, message: '请输入管理员邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ]
}

// 角色权限
const rolePermissions = ref([
  {
    role: 'admin',
    description: '系统管理员',
    permissions: ['user_manage', 'project_manage', 'competition_manage', 'system_config', 'data_export', 'log_view']
  },
  {
    role: 'teacher',
    description: '教师',
    permissions: ['project_manage', 'competition_manage', 'data_export']
  },
  {
    role: 'student',
    description: '学生',
    permissions: []
  }
])

// 安全设置
const securitySettings = reactive({
  minPasswordLength: 8,
  passwordComplexity: ['uppercase', 'lowercase', 'number'],
  loginLockEnabled: true,
  maxLoginAttempts: 5,
  lockDuration: 30,
  sessionTimeout: 120
})

const securityRules = {
  minPasswordLength: [
    { required: true, message: '请设置密码最小长度', trigger: 'blur' }
  ],
  maxLoginAttempts: [
    { required: true, message: '请设置失败次数阈值', trigger: 'blur' }
  ],
  lockDuration: [
    { required: true, message: '请设置锁定时间', trigger: 'blur' }
  ],
  sessionTimeout: [
    { required: true, message: '请设置会话超时时间', trigger: 'blur' }
  ]
}

// 邮件设置
const emailSettings = reactive({
  smtpServer: 'smtp.qq.com',
  smtpPort: 587,
  senderEmail: 'system@yunmeng.edu.cn',
  senderName: '云梦大学系统',
  emailPassword: '',
  enableSSL: true
})

const emailRules = {
  smtpServer: [
    { required: true, message: '请输入SMTP服务器地址', trigger: 'blur' }
  ],
  smtpPort: [
    { required: true, message: '请输入SMTP端口', trigger: 'blur' }
  ],
  senderEmail: [
    { required: true, message: '请输入发件人邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  senderName: [
    { required: true, message: '请输入发件人名称', trigger: 'blur' }
  ],
  emailPassword: [
    { required: true, message: '请输入邮箱密码', trigger: 'blur' }
  ]
}

// 备份设置
const backupSettings = reactive({
  autoBackup: true,
  backupFrequency: 'daily',
  backupTime: new Date(2024, 0, 1, 2, 0),
  backupRetention: 30,
  backupPath: '/backup'
})

const backupRules = {
  backupRetention: [
    { required: true, message: '请设置保留备份数量', trigger: 'blur' }
  ],
  backupPath: [
    { required: true, message: '请输入备份路径', trigger: 'blur' }
  ]
}

// 方法
const saveBasicSettings = async () => {
  try {
    await basicFormRef.value.validate()
    saving.value = true
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    ElMessage.success('基本设置保存成功')
  } catch (error) {
    ElMessage.error('保存失败：' + error.message)
  } finally {
    saving.value = false
  }
}

const saveRolePermissions = async (role) => {
  try {
    saving.value = true
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 500))
    
    ElMessage.success(`${role.description}权限保存成功`)
  } catch (error) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

const saveSecuritySettings = async () => {
  try {
    await securityFormRef.value.validate()
    saving.value = true
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    ElMessage.success('安全设置保存成功')
  } catch (error) {
    ElMessage.error('保存失败：' + error.message)
  } finally {
    saving.value = false
  }
}

const saveEmailSettings = async () => {
  try {
    await emailFormRef.value.validate()
    saving.value = true
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    ElMessage.success('邮件设置保存成功')
  } catch (error) {
    ElMessage.error('保存失败：' + error.message)
  } finally {
    saving.value = false
  }
}

const testEmailSettings = async () => {
  try {
    testing.value = true
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 2000))
    
    ElMessage.success('邮件连接测试成功')
  } catch (error) {
    ElMessage.error('邮件连接测试失败')
  } finally {
    testing.value = false
  }
}

const saveBackupSettings = async () => {
  try {
    await backupFormRef.value.validate()
    saving.value = true
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    ElMessage.success('备份设置保存成功')
  } catch (error) {
    ElMessage.error('保存失败：' + error.message)
  } finally {
    saving.value = false
  }
}

const createBackup = async () => {
  try {
    backingUp.value = true
    
    // 模拟API调用
    await new Promise(resolve => setTimeout(resolve, 3000))
    
    ElMessage.success('数据备份创建成功')
  } catch (error) {
    ElMessage.error('数据备份创建失败')
  } finally {
    backingUp.value = false
  }
}
</script>

<style scoped>
.system-settings {
  padding: 20px;
}

.page-header {
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
  color: #2c3e50;
}

.setting-card {
  margin-bottom: 20px;
}

.el-checkbox-group {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}
</style> 
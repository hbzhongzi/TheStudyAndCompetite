<template>
  <div class="permission-container">
    <!-- ================= 搜索区域 ================= -->
    <el-card shadow="never" class="search-card">
      <el-form :inline="true" :model="searchForm">

        <el-form-item label="用户名">
          <el-input
            v-model="searchForm.username"
            placeholder="请输入用户名"
            clearable
          />
        </el-form-item>

        <el-form-item label="角色">
          <el-select v-model="searchForm.role" placeholder="请选择角色" clearable>
            <el-option label="管理员" value="admin" />
            <el-option label="教师" value="teacher" />
            <el-option label="学生" value="student" />
          </el-select>
        </el-form-item>

        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="请选择状态" clearable>
            <el-option label="启用" value="active" />
            <el-option label="禁用" value="inactive" />
          </el-select>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="loadUsers">
            搜索
          </el-button>
          <el-button @click="resetSearch">
            重置
          </el-button>
        </el-form-item>

      </el-form>
    </el-card>

    <!-- ================= 用户表格 ================= -->
    <el-card shadow="never" class="table-card">

      <el-table
        :data="filteredUsers"
        border
        stripe
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="80" />

        <el-table-column prop="username" label="用户名" />

        <el-table-column prop="realName" label="真实姓名" />

        <el-table-column prop="email" label="邮箱" />

        <el-table-column prop="department" label="学院" />

        <el-table-column label="角色">
          <template #default="{ row }">
            <el-tag
              v-for="role in row.roleNames"
              :key="role"
              type="info"
              class="role-tag"
            >
              {{ role }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="状态">
          <template #default="{ row }">
            <el-tag
              :type="row.status === 'active' ? 'success' : 'danger'"
            >
              {{ row.status === 'active' ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column prop="createTime" label="创建时间" width="180" />

        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button
              size="small"
              type="primary"
              @click="openStatusDialog(row)"
            >
              修改权限
            </el-button>
          </template>
        </el-table-column>

      </el-table>

      <!-- 分页 -->
      <el-pagination
        class="pagination"
        background
        layout="total, prev, pager, next"
        :current-page="page"
        :page-size="size"
        :total="total"
        @current-change="handlePageChange"
      />

    </el-card>

    <!-- ================= 修改权限弹窗 ================= -->
    <el-dialog
      v-model="statusDialogVisible"
      title="修改用户权限"
      width="400px"
    >
      <div>
        <p><strong>用户名：</strong>{{ currentUser.username }}</p >
        <p><strong>真实姓名：</strong>{{ currentUser.realName }}</p >
      </div>

      <el-divider />

      <el-radio-group v-model="statusForm.status">
        <el-radio label="active">启用</el-radio>
        <el-radio label="inactive">禁用</el-radio>
      </el-radio-group>

      <template #footer>
        <el-button @click="statusDialogVisible = false">
          取消
        </el-button>
        <el-button type="primary" @click="submitStatus">
          确认修改
        </el-button>
      </template>
    </el-dialog>

  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import userService from '@/services/userService'

// ================= 数据 =================
const users = ref([])
const page = ref(1)
const size = ref(20)
const total = ref(0)

const searchForm = ref({
  username: '',
  role: '',
  status: ''
})

const statusDialogVisible = ref(false)
const currentUser = ref({})
const statusForm = ref({
  status: ''
})

// ================= 获取用户列表 =================
const loadUsers = async () => {
  try {
    const res = await userService.getUserList({
      page: page.value,
      size: size.value
    })

    users.value = res.data.list
    total.value = res.data.total

  } catch (error) {
    ElMessage.error('获取用户列表失败')
  }
}

// ================= 前端过滤 =================
const filteredUsers = computed(() => {
  return users.value.filter(user => {

    if (searchForm.value.username &&
        !user.username.includes(searchForm.value.username)) {
      return false
    }

    if (searchForm.value.role &&
        !user.roleNames.includes(searchForm.value.role)) {
      return false
    }

    if (searchForm.value.status &&
        user.status !== searchForm.value.status) {
      return false
    }

    return true
  })
})

// ================= 分页 =================
const handlePageChange = (newPage) => {
  page.value = newPage
  loadUsers()
}

// ================= 重置搜索 =================
const resetSearch = () => {
  searchForm.value = {
    username: '',
    role: '',
    status: ''
  }
}

// ================= 打开修改权限弹窗 =================
const openStatusDialog = (row) => {
  currentUser.value = row
  statusForm.value.status = row.status
  statusDialogVisible.value = true
}

// ================= 提交权限修改 =================
const submitStatus = async () => {
  try {
    await userService.toggleUserStatus(currentUser.value.id, statusForm.value.status)

    ElMessage.success('权限修改成功')
    statusDialogVisible.value = false
    loadUsers()

  } catch (error) {
    ElMessage.error('权限修改失败')
  }
}

onMounted(() => {
  loadUsers()
})
</script>

<style scoped>
.permission-container {
  padding: 20px;
}

.search-card {
  margin-bottom: 20px;
}

.table-card {
  margin-top: 10px;
}

.role-tag {
  margin-right: 5px;
}

.pagination {
  margin-top: 20px;
  text-align: right;
}
</style>
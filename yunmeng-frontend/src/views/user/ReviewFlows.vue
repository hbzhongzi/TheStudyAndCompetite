<template>
  <div class="review-flows">
    <!-- 审核流程配置 -->
    <el-card style="margin-bottom: 20px;">
      <template #header>
        <span>审核流程配置</span>
        <el-button type="primary" @click="addFlow" style="float: right;">
          新增流程
        </el-button>
      </template>
      
      <el-table :data="flowList" style="width: 100%">
        <el-table-column prop="name" label="流程名称" />
        <el-table-column prop="type" label="适用类型" />
        <el-table-column prop="steps" label="审核步骤">
          <template #default="scope">
            <el-tag v-for="step in scope.row.steps" :key="step.id" style="margin-right: 5px;">
              {{ step.name }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态">
          <template #default="scope">
            <el-tag :type="scope.row.status === 'active' ? 'success' : 'info'">
              {{ scope.row.status === 'active' ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="scope">
            <el-button size="small" @click="editFlow(scope.row)">编辑</el-button>
            <el-button 
              size="small" 
              :type="scope.row.status === 'active' ? 'warning' : 'success'"
              @click="toggleFlowStatus(scope.row)"
            >
              {{ scope.row.status === 'active' ? '禁用' : '启用' }}
            </el-button>
            <el-button size="small" type="danger" @click="deleteFlow(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 流程编辑对话框 -->
    <el-dialog
      v-model="flowDialogVisible"
      :title="flowDialogTitle"
      width="60%"
    >
      <el-form :model="flowForm" :rules="flowRules" ref="flowFormRef" label-width="100px">
        <el-form-item label="流程名称" prop="name">
          <el-input v-model="flowForm.name" placeholder="请输入流程名称" />
        </el-form-item>
        
        <el-form-item label="适用类型" prop="type">
          <el-select v-model="flowForm.type" placeholder="请选择适用类型">
            <el-option label="软件开发" value="software" />
            <el-option label="科研项目" value="research" />
            <el-option label="创新项目" value="innovation" />
            <el-option label="竞赛项目" value="competition" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="审核步骤" prop="steps">
          <div v-for="(step, index) in flowForm.steps" :key="index" class="step-item">
            <el-row :gutter="10">
              <el-col :span="8">
                <el-input v-model="step.name" placeholder="步骤名称" />
              </el-col>
              <el-col :span="6">
                <el-select v-model="step.role" placeholder="审核角色">
                  <el-option label="指导教师" value="teacher" />
                  <el-option label="系主任" value="department_head" />
                  <el-option label="教务处长" value="academic_dean" />
                  <el-option label="校长" value="president" />
                </el-select>
              </el-col>
              <el-col :span="6">
                <el-input-number v-model="step.duration" :min="1" :max="30" placeholder="期限(天)" />
              </el-col>
              <el-col :span="4">
                <el-button type="danger" @click="removeStep(index)" :disabled="flowForm.steps.length <= 1">
                  删除
                </el-button>
              </el-col>
            </el-row>
          </div>
          <el-button type="primary" @click="addStep">添加步骤</el-button>
        </el-form-item>
        
        <el-form-item label="流程描述">
          <el-input 
            v-model="flowForm.description" 
            type="textarea" 
            :rows="3" 
            placeholder="请输入流程描述"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="flowDialogVisible = false">取消</el-button>
          <el-button type="primary" @click="saveFlow">保存</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

// 响应式数据
const flowList = ref([])
const flowDialogVisible = ref(false)
const flowDialogTitle = ref('')
const flowFormRef = ref(null)

const flowForm = reactive({
  id: null,
  name: '',
  type: '',
  steps: [{ name: '', role: '', duration: 7 }],
  description: '',
  status: 'active'
})

const flowRules = {
  name: [{ required: true, message: '请输入流程名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择适用类型', trigger: 'change' }],
  steps: [{ required: true, message: '请配置审核步骤', trigger: 'change' }]
}

// 加载流程列表
const loadFlows = async () => {
  try {
    // 模拟数据
    flowList.value = [
      {
        id: 1,
        name: '标准项目审核流程',
        type: '软件开发',
        steps: [
          { id: 1, name: '指导教师审核', role: 'teacher', duration: 3 },
          { id: 2, name: '系主任审核', role: 'department_head', duration: 5 },
          { id: 3, name: '教务处长审核', role: 'academic_dean', duration: 7 }
        ],
        status: 'active'
      },
      {
        id: 2,
        name: '竞赛项目审核流程',
        type: '竞赛项目',
        steps: [
          { id: 1, name: '指导教师审核', role: 'teacher', duration: 2 },
          { id: 2, name: '竞赛委员会审核', role: 'competition_committee', duration: 3 }
        ],
        status: 'active'
      }
    ]
  } catch (error) {
    console.error('加载流程列表失败:', error)
    ElMessage.error('加载流程列表失败')
  }
}

// 新增流程
const addFlow = () => {
  flowForm.id = null
  flowForm.name = ''
  flowForm.type = ''
  flowForm.steps = [{ name: '', role: '', duration: 7 }]
  flowForm.description = ''
  flowForm.status = 'active'
  flowDialogTitle.value = '新增审核流程'
  flowDialogVisible.value = true
}

// 编辑流程
const editFlow = (flow) => {
  Object.assign(flowForm, {
    id: flow.id,
    name: flow.name,
    type: flow.type,
    steps: [...flow.steps],
    description: flow.description || '',
    status: flow.status
  })
  flowDialogTitle.value = '编辑审核流程'
  flowDialogVisible.value = true
}

// 添加步骤
const addStep = () => {
  flowForm.steps.push({ name: '', role: '', duration: 7 })
}

// 删除步骤
const removeStep = (index) => {
  if (flowForm.steps.length > 1) {
    flowForm.steps.splice(index, 1)
  }
}

// 保存流程
const saveFlow = async () => {
  try {
    await flowFormRef.value.validate()
    
    if (flowForm.id) {
      // 更新流程
      const index = flowList.value.findIndex(f => f.id === flowForm.id)
      if (index !== -1) {
        flowList.value[index] = { ...flowForm }
      }
      ElMessage.success('流程更新成功')
    } else {
      // 新增流程
      const newFlow = {
        ...flowForm,
        id: Date.now()
      }
      flowList.value.push(newFlow)
      ElMessage.success('流程创建成功')
    }
    
    flowDialogVisible.value = false
  } catch (error) {
    console.error('保存流程失败:', error)
    ElMessage.error('保存失败，请检查表单')
  }
}

// 切换流程状态
const toggleFlowStatus = async (flow) => {
  try {
    const newStatus = flow.status === 'active' ? 'inactive' : 'active'
    flow.status = newStatus
    ElMessage.success(`流程已${newStatus === 'active' ? '启用' : '禁用'}`)
  } catch (error) {
    console.error('切换流程状态失败:', error)
    ElMessage.error('操作失败')
  }
}

// 删除流程
const deleteFlow = async (flow) => {
  try {
    await ElMessageBox.confirm('确定要删除这个审核流程吗？', '确认删除', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    const index = flowList.value.findIndex(f => f.id === flow.id)
    if (index !== -1) {
      flowList.value.splice(index, 1)
      ElMessage.success('流程删除成功')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除流程失败:', error)
      ElMessage.error('删除失败')
    }
  }
}

// 组件挂载时加载数据
onMounted(() => {
  loadFlows()
})
</script>

<style scoped>
.review-flows {
  padding: 20px;
}

.step-item {
  margin-bottom: 10px;
  padding: 10px;
  border: 1px solid #e4e7ed;
  border-radius: 4px;
}

.dialog-footer {
  text-align: right;
}
</style> 
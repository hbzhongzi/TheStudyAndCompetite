<template>
  <div class="competition-management">

    <div class="header">
      <h2>竞赛管理</h2>
      <el-button type="primary" @click="openDialog">
        <el-icon><Plus /></el-icon>
        新建竞赛
      </el-button>
    </div>

    <!-- 表格 -->
    <el-table :data="competitions" v-loading="loading" border>

      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="title" label="竞赛名称" min-width="200" />
      <el-table-column prop="category" label="类别" width="120" />

      <el-table-column label="级别" width="100">
        <template #default="scope">
          <el-tag :type="scope.row.level === 'national' ? 'danger' : 'success'">
            {{ scope.row.level === 'national' ? '国家级' : '校级' }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column label="报名时间" width="260">
        <template #default="scope">
          <div class="time-range">
            <div>{{ formatTime(scope.row.registrationStart) }}</div>
            <div class="divider">至</div>
            <div>{{ formatTime(scope.row.registrationEnd) }}</div>
          </div>
        </template>
      </el-table-column>

      <el-table-column label="报名状态" width="100">
        <template #default="scope">
          <el-tag :type="scope.row.isOpen ? 'success' : 'info'">
            {{ scope.row.isOpen ? '开放' : '关闭' }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column label="当前人数" width="120">
        <template #default="scope">
          {{ scope.row.currentParticipants || 0 }} / {{ scope.row.maxParticipants }}
        </template>
      </el-table-column>

      <!-- 操作列 -->
      <el-table-column label="操作" width="300" fixed="right">
        <template #default="scope">

          <el-button
            size="small"
            @click="viewDetail(scope.row.id)">
            详情
          </el-button>

          <el-button
            size="small"
            :type="scope.row.isOpen ? 'warning' : 'success'"
            @click="toggleOpen(scope.row)">
            {{ scope.row.isOpen ? '关闭报名' : '开放报名' }}
          </el-button>

          <el-button
            size="small"
            type="danger"
            @click="handleDelete(scope.row.id)">
            删除
          </el-button>

        </template>
      </el-table-column>

    </el-table>

    <!-- 分页 -->
    <div style="margin-top:20px;text-align:right">
      <el-pagination
        background
        layout="total, prev, pager, next"
        :total="total"
        :page-size="size"
        :current-page="page"
        @current-change="handlePageChange"
      />
    </div>

    <!-- 创建竞赛 -->
    <el-dialog v-model="dialogVisible" title="创建竞赛" width="600px">
      <el-form :model="form" label-width="110px">

        <el-form-item label="竞赛名称">
          <el-input v-model="form.title" />
        </el-form-item>

        <el-form-item label="竞赛类别">
          <el-input v-model="form.type" />
        </el-form-item>

        <el-form-item label="主办单位">
          <el-input v-model="form.organizer" />
        </el-form-item>

        <el-form-item label="竞赛级别">
          <el-select v-model="form.level">
            <el-option label="校级" value="school" />
            <el-option label="国家级" value="national" />
          </el-select>
        </el-form-item>

        <el-form-item label="报名时间">
          <el-date-picker
            v-model="registrationRange"
            type="datetimerange"
            format="YYYY-MM-DD HH:mm"
            value-format="YYYY-MM-DDTHH:mm:ss[Z]"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item label="竞赛时间">
          <el-date-picker
            v-model="competitionRange"
            type="datetimerange"
            format="YYYY-MM-DD HH:mm"
            value-format="YYYY-MM-DDTHH:mm:ss[Z]"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item label="最大人数">
          <el-input-number v-model="form.maxParticipants" :min="1" />
        </el-form-item>

        <el-form-item label="开放报名">
          <el-switch v-model="form.isOpen" />
        </el-form-item>

        <el-form-item label="竞赛说明">
          <el-input type="textarea" rows="3" v-model="form.description" />
        </el-form-item>

        <el-form-item label="附件地址">
          <el-input v-model="form.attachment" />
        </el-form-item>

      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="createCompetition">
          创建
        </el-button>
      </template>
    </el-dialog>

  </div>

  <!-- ===================== 详情弹窗 ===================== -->
    <el-dialog v-model="detailVisible" title="竞赛详情" width="750px">

      <div v-if="detailData">

        <!-- 基本信息 -->
        <el-divider>基本信息</el-divider>
        <el-descriptions :column="2" border>

          <el-descriptions-item label="竞赛名称">
            {{ detailData.title }}
          </el-descriptions-item>

          <el-descriptions-item label="类别">
            {{ detailData.category }}
          </el-descriptions-item>

          <el-descriptions-item label="级别">
            <el-tag :type="detailData.level === 'national' ? 'danger' : 'success'">
              {{ detailData.level === 'national' ? '国家级' : '校级' }}
            </el-tag>
          </el-descriptions-item>

          <el-descriptions-item label="系统状态">
            <el-tag type="warning">
              {{ formatStatus(detailData.status) }}
            </el-tag>
          </el-descriptions-item>

        </el-descriptions>

        <!-- 时间 -->
        <el-divider>时间安排</el-divider>
        <el-descriptions :column="1" border>

          <el-descriptions-item label="报名时间">
            {{ formatTime(detailData.registrationStart) }}
            至
            {{ formatTime(detailData.registrationEnd) }}
          </el-descriptions-item>

          <el-descriptions-item label="提交时间">
            {{ formatTime(detailData.submissionStart) }}
            至
            {{ formatTime(detailData.submissionEnd) }}
          </el-descriptions-item>

        </el-descriptions>

        <!-- 报名情况 -->
        <el-divider>报名情况</el-divider>
        <el-descriptions :column="2" border>

          <el-descriptions-item label="报名状态">
            <el-tag :type="detailData.isOpen ? 'success' : 'info'">
              {{ detailData.isOpen ? '开放报名' : '关闭报名' }}
            </el-tag>
          </el-descriptions-item>

          <el-descriptions-item label="报名人数">
            {{ detailData.currentParticipants }} /
            {{ detailData.maxParticipants }}
          </el-descriptions-item>

        </el-descriptions>

        <!-- 奖项 -->
        <el-divider>奖项配置</el-divider>
        <el-descriptions :column="3" border>

          <el-descriptions-item label="一等奖">
            {{ detailData.awardConfig?.first_prize || 0 }}
          </el-descriptions-item>

          <el-descriptions-item label="二等奖">
            {{ detailData.awardConfig?.second_prize || 0 }}
          </el-descriptions-item>

          <el-descriptions-item label="三等奖">
            {{ detailData.awardConfig?.third_prize || 0 }}
          </el-descriptions-item>

        </el-descriptions>

        <!-- 系统信息 -->
        <el-divider>系统信息</el-divider>
        <el-descriptions :column="1" border>

          <el-descriptions-item label="创建时间">
            {{ formatTime(detailData.createdAt) }}
          </el-descriptions-item>

          <el-descriptions-item label="更新时间">
            {{ formatTime(detailData.updatedAt) }}
          </el-descriptions-item>

        </el-descriptions>

        <!-- 说明 -->
        <el-divider>竞赛说明</el-divider>
        <div class="description-box">
          {{ detailData.description }}
        </div>

      </div>

    </el-dialog>

</template>

<script>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import competitionService from '@/services/competitionService'
import { ElMessageBox } from 'element-plus'

export default {
  name: 'CompetitionManagement',
  components: { Plus },

  setup() {

    const loading = ref(false)
    const competitions = ref([])
    const page = ref(1)
    const size = ref(10)
    const total = ref(0)
    const dialogVisible = ref(false)
    const detailVisible = ref(false)
    const detailData = ref(null)

    const registrationRange = ref([])
    const competitionRange = ref([])

    const form = reactive({
      title: '',
      type: '',
      organizer: '',
      level: 'school',
      description: '',
      attachment: '',
      isOpen: true,
      maxParticipants: 100
    })

    const formatTime = (time) => {
      if (!time) return ''
      return time.replace('T', ' ').replace('Z', '')
    }

    const fetchCompetitions = async () => {
      try {
        loading.value = true
        const res = await competitionService.getCompetitions({
          page: page.value,
          size: size.value
        })
        competitions.value = res.data.list
        total.value = res.data.total
      } catch {
        ElMessage.error('获取失败')
      } finally {
        loading.value = false
      }
    }

const viewDetail = async (id) => {
  try {
    const res = await competitionService.getCompetitionDetail(id)
    detailData.value = res.data
    detailVisible.value = true
  } catch (error) {
    console.log(error)
    ElMessage.error(error.response?.data?.message || '获取详情失败')
  }
}

const toggleOpen = async (row) => {
  try {
    await competitionService.toggleCompetitionOpen(row.id)

    ElMessage.success('状态更新成功')

    // 更新当前行状态
    row.isOpen = !row.isOpen

  } catch {
    ElMessage.error('更新失败')
  }
}

const handleDelete = async (id) => {
  try {
    await ElMessageBox.confirm(
      '确定要删除该竞赛吗？此操作不可恢复！',
      '危险操作',
      {
        confirmButtonText: '确认删除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    await competitionService.deleteCompetition(id)

    ElMessage.success('删除成功')

    fetchCompetitions()

  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}


    const handlePageChange = (newPage) => {
      page.value = newPage
      fetchCompetitions()
    }

    const openDialog = () => {
      dialogVisible.value = true
    }
    const formatStatus = (status) => {
  const map = {
    draft: '草稿',
    registration: '报名中',
    ongoing: '进行中',
    finished: '已结束'
  }
  return map[status] || status
}

const createCompetition = async () => {
  try {
    const payload = {
      ...form,
      RegistrationStart: registrationRange.value?.[0] || null,
      RegistrationEnd: registrationRange.value?.[1] || null,
      StartTime: competitionRange.value?.[0] || null,
      EndTime: competitionRange.value?.[1] || null
    }

    console.log("最终提交:", payload)

    await competitionService.createCompetition(payload)

    ElMessage.success('创建成功')
    dialogVisible.value = false
    fetchCompetitions()

  } catch (e) {
    ElMessage.error('创建失败')
  }
}

    onMounted(fetchCompetitions)

    return {
      competitions,
      loading,
      page,
      size,
      total,
      dialogVisible,
      form,
      registrationRange,
      competitionRange,
      openDialog,
      createCompetition,
      handlePageChange,
      formatTime,
      formatStatus,
      viewDetail,
      toggleOpen,
      handleDelete,
      detailVisible,
      detailData,
      viewDetail,
    }
  }
}
</script>

<style scoped>
.header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 20px;
}

.time-range {
  line-height: 1.5;
  font-size: 13px;
}

.divider {
  color: #999;
  font-size: 12px;
}
</style>
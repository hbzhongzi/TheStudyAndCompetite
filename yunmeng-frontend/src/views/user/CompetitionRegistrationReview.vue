<template>
  <div class="judge-assign">

    <h2>评审任务分配</h2>

    <!-- 选择竞赛 -->
    <el-card style="margin-bottom:20px">
      <el-form label-width="100px">
        <el-form-item label="选择竞赛">
          <el-select
            v-model="selectedCompetitionId"
            placeholder="请选择开放报名的竞赛"
            @change="loadJudges"
          >
            <el-option
              v-for="item in openCompetitions"
              :key="item.id"
              :label="item.title"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 已分配教师 -->
    <el-card v-if="selectedCompetitionId" style="margin-bottom:20px">
      <template #header>已分配评审教师</template>

      <el-table :data="judgeList" border>
        <el-table-column prop="teacher.realName" label="姓名" />
        <el-table-column prop="teacher.department" label="学院" />
        <el-table-column prop="teacher.email" label="邮箱" />
        <el-table-column prop="status" label="状态" />
      </el-table>
    </el-card>

    <!-- 分配教师 -->
    <el-card v-if="selectedCompetitionId">
      <template #header>分配新评审教师</template>

      <el-form label-width="100px">

        <el-form-item label="选择教师">
          <el-select
            v-model="assignForm.teacherID"
            placeholder="请选择教师"
            filterable
          >
            <el-option
              v-for="teacher in teacherList"
              :key="teacher.id"
              :label="teacher.realName + ' - ' + teacher.department"
              :value="teacher.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" @click="assignJudge">
            分配评审
          </el-button>
        </el-form-item>

      </el-form>

    </el-card>

  </div>
</template>

<script>
import { ref, onMounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import competitionService from '@/services/competitionService'
import userService from '@/services/userService'

export default {
  name: 'AdminJudgeAssign',

  setup() {

    const competitions = ref([])
    const openCompetitions = ref([])
    const selectedCompetitionId = ref(null)

    const judgeList = ref([])
    const teacherList = ref([])

    const assignForm = ref({
      teacherID: null
    })

    // 加载竞赛
    const loadCompetitions = async () => {
      const res = await competitionService.getCompetitions()

      competitions.value = res.data.list

      openCompetitions.value = competitions.value.filter(c => c.isOpen)

      // 默认选中第一个开放竞赛
      if (openCompetitions.value.length > 0) {
        selectedCompetitionId.value = openCompetitions.value[0].id
        loadJudges()
      }
    }

    // 加载评审教师
    const loadJudges = async () => {
      if (!selectedCompetitionId.value) return
      const res = await competitionService.getJudgesByCompetition(selectedCompetitionId.value)
      judgeList.value = res.data
    }

    // 加载教师列表
    const loadTeachers = async () => {
      const res = await userService.getUserList({ role: 'teacher' })
      teacherList.value = res.data.list
    }

    // 分配教师
    const assignJudge = async () => {

      if (!assignForm.value.teacherID) {
        ElMessage.warning('请选择教师')
        return
      }

      await competitionService.distributeJudge({
        competition_id: selectedCompetitionId.value,
        teacher_id: assignForm.value.teacherID,
        status: 'active'
      })

      ElMessage.success('分配成功')

      assignForm.value.teacherID = null

      loadJudges()
    }

    onMounted(() => {
      loadCompetitions()
      loadTeachers()
    })

    return {
      openCompetitions,
      selectedCompetitionId,
      judgeList,
      teacherList,
      assignForm,
      loadJudges,
      assignJudge
    }
  }
}
</script>

<style scoped>
.judge-assign {
  padding: 20px;
}
</style>
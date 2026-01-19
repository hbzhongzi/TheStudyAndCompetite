<template>
  <div class="project-templates">
    <!-- 模板筛选 -->
    <div class="template-filter">
      <el-row :gutter="20">
        <el-col :span="6">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索模板名称"
            clearable
            @input="filterTemplates"
          >
            <template #prefix>
              <i class="el-icon-search"></i>
            </template>
          </el-input>
        </el-col>
        <el-col :span="4">
          <el-select v-model="selectedType" placeholder="项目类型" clearable @change="filterTemplates">
            <el-option label="全部类型" value="" />
            <el-option label="科研" value="科研" />
            <el-option label="竞赛" value="竞赛" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-select v-model="selectedLevel" placeholder="项目级别" clearable @change="filterTemplates">
            <el-option label="全部级别" value="" />
            <el-option label="校级" value="校级" />
            <el-option label="省级" value="省级" />
            <el-option label="国家级" value="国家级" />
            <el-option label="国际级" value="国际级" />
          </el-select>
        </el-col>
        <el-col :span="4">
          <el-select v-model="selectedCategory" placeholder="项目分类" clearable @change="filterTemplates">
            <el-option label="全部分类" value="" />
            <el-option label="计算机科学" value="计算机科学" />
            <el-option label="机械工程" value="机械工程" />
            <el-option label="经济管理" value="经济管理" />
            <el-option label="人文社科" value="人文社科" />
          </el-select>
        </el-col>
        <el-col :span="6">
          <el-button type="primary" @click="filterTemplates">筛选</el-button>
          <el-button @click="resetFilter">重置</el-button>
        </el-col>
      </el-row>
    </div>

    <!-- 模板列表 -->
    <div class="template-list">
      <el-row :gutter="20">
        <el-col 
          v-for="template in filteredTemplates" 
          :key="template.id" 
          :span="8"
          class="template-item"
        >
          <el-card class="template-card" shadow="hover">
            <div class="template-header">
              <div class="template-type">
                <el-tag :type="template.type === '科研' ? 'primary' : 'success'" size="small">
                  {{ template.type }}
                </el-tag>
                <el-tag :type="getLevelType(template.level)" size="small">
                  {{ template.level }}
                </el-tag>
              </div>
              <div class="template-rating">
                <el-rate 
                  v-model="template.rating" 
                  disabled 
                  show-score 
                  text-color="#ff9900"
                  score-template="{value}"
                />
              </div>
            </div>
            
            <div class="template-title">
              <h3>{{ template.name }}</h3>
            </div>
            
            <div class="template-description">
              <p>{{ template.description }}</p>
            </div>
            
            <div class="template-meta">
              <div class="meta-item">
                <i class="el-icon-user"></i>
                <span>适用人数: {{ template.memberCount }}人</span>
              </div>
              <div class="meta-item">
                <i class="el-icon-time"></i>
                <span>预计周期: {{ template.duration }}个月</span>
              </div>
              <div class="meta-item">
                <i class="el-icon-s-data"></i>
                <span>使用次数: {{ template.usageCount }}次</span>
              </div>
            </div>
            
            <div class="template-tags">
              <el-tag 
                v-for="tag in template.tags" 
                :key="tag" 
                size="small" 
                class="template-tag"
              >
                {{ tag }}
              </el-tag>
            </div>
            
            <div class="template-actions">
              <el-button type="primary" @click="previewTemplate(template)">
                <i class="el-icon-view"></i>
                预览
              </el-button>
              <el-button type="success" @click="selectTemplate(template)">
                <i class="el-icon-check"></i>
                使用模板
              </el-button>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>

    <!-- 分页 -->
    <div class="template-pagination" v-if="totalTemplates > 0">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[9, 18, 36]"
        :total="totalTemplates"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handleCurrentChange"
      />
    </div>

    <!-- 模板预览对话框 -->
    <el-dialog
      v-model="showPreviewDialog"
      title="模板预览"
      width="70%"
      :close-on-click-modal="false"
    >
      <div v-if="selectedTemplate" class="template-preview">
        <div class="preview-header">
          <h2>{{ selectedTemplate.name }}</h2>
          <div class="preview-meta">
            <el-tag :type="selectedTemplate.type === '科研' ? 'primary' : 'success'">
              {{ selectedTemplate.type }}
            </el-tag>
            <el-tag :type="getLevelType(selectedTemplate.level)">
              {{ selectedTemplate.level }}
            </el-tag>
            <el-tag type="info">{{ selectedTemplate.category }}</el-tag>
          </div>
        </div>
        
        <div class="preview-content">
          <el-row :gutter="20">
            <el-col :span="16">
              <div class="preview-section">
                <h3>项目描述</h3>
                <p>{{ selectedTemplate.description }}</p>
              </div>
              
              <div class="preview-section">
                <h3>项目目标</h3>
                <ul>
                  <li v-for="(goal, index) in selectedTemplate.goals" :key="index">
                    {{ goal }}
                  </li>
                </ul>
              </div>
              
              <div class="preview-section">
                <h3>技术要求</h3>
                <ul>
                  <li v-for="(requirement, index) in selectedTemplate.requirements" :key="index">
                    {{ requirement }}
                  </li>
                </ul>
              </div>
              
              <div class="preview-section">
                <h3>预期成果</h3>
                <ul>
                  <li v-for="(outcome, index) in selectedTemplate.outcomes" :key="index">
                    {{ outcome }}
                  </li>
                </ul>
              </div>
            </el-col>
            
            <el-col :span="8">
              <div class="preview-sidebar">
                <div class="preview-card">
                  <h4>基本信息</h4>
                  <div class="info-item">
                    <span class="label">适用人数:</span>
                    <span class="value">{{ selectedTemplate.memberCount }}人</span>
                  </div>
                  <div class="info-item">
                    <span class="label">预计周期:</span>
                    <span class="value">{{ selectedTemplate.duration }}个月</span>
                  </div>
                  <div class="info-item">
                    <span class="label">难度等级:</span>
                    <span class="value">{{ selectedTemplate.difficulty }}</span>
                  </div>
                  <div class="info-item">
                    <span class="label">使用次数:</span>
                    <span class="value">{{ selectedTemplate.usageCount }}次</span>
                  </div>
                  <div class="info-item">
                    <span class="label">评分:</span>
                    <span class="value">
                      <el-rate 
                        v-model="selectedTemplate.rating" 
                        disabled 
                        show-score 
                        text-color="#ff9900"
                        score-template="{value}"
                      />
                    </span>
                  </div>
                </div>
                
                <div class="preview-card">
                  <h4>标签</h4>
                  <div class="preview-tags">
                    <el-tag 
                      v-for="tag in selectedTemplate.tags" 
                      :key="tag" 
                      size="small"
                    >
                      {{ tag }}
                    </el-tag>
                  </div>
                </div>
                
                <div class="preview-card">
                  <h4>操作</h4>
                  <el-button type="success" @click="selectTemplate(selectedTemplate)" block>
                    使用此模板
                  </el-button>
                </div>
              </div>
            </el-col>
          </el-row>
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { studentService } from '../services/studentService'

// 响应式数据
const searchKeyword = ref('')
const selectedType = ref('')
const selectedLevel = ref('')
const selectedCategory = ref('')
const currentPage = ref(1)
const pageSize = ref(9)
const totalTemplates = ref(0)
const showPreviewDialog = ref(false)
const selectedTemplate = ref(null)

// 模板列表
const templates = ref([])

// 计算属性
const filteredTemplates = computed(() => {
  let filtered = templates.value

  if (searchKeyword.value) {
    filtered = filtered.filter(template => 
      template.name.toLowerCase().includes(searchKeyword.value.toLowerCase()) ||
      template.description.toLowerCase().includes(searchKeyword.value.toLowerCase())
    )
  }

  if (selectedType.value) {
    filtered = filtered.filter(template => template.type === selectedType.value)
  }

  if (selectedLevel.value) {
    filtered = filtered.filter(template => template.level === selectedLevel.value)
  }

  if (selectedCategory.value) {
    filtered = filtered.filter(template => template.category === selectedCategory.value)
  }

  return filtered
})

// 方法
const loadTemplates = async () => {
  try {
    const response = await studentService.getProjectTemplates()
    if (response && response.code === 200) {
      templates.value = response.data.list || []
      totalTemplates.value = response.data.total || 0
    } else {
      // 使用模拟数据
      templates.value = getMockTemplates()
      totalTemplates.value = templates.value.length
    }
  } catch (error) {
    console.error('加载模板失败:', error)
    // 使用模拟数据
    templates.value = getMockTemplates()
    totalTemplates.value = templates.value.length
  }
}

const getMockTemplates = () => {
  return [
    {
      id: 1,
      name: '智能校园导航系统',
      type: '科研',
      level: '校级',
      category: '计算机科学',
      description: '基于移动端的智能校园导航系统，集成地图服务、路径规划、实时定位等功能。',
      memberCount: 3,
      duration: 6,
      usageCount: 15,
      rating: 4.5,
      tags: ['移动开发', '地图服务', '路径规划'],
      goals: [
        '实现校园地图的数字化展示',
        '提供最优路径规划算法',
        '集成实时定位和导航功能'
      ],
      requirements: [
        '熟悉移动端开发技术',
        '了解地图服务API使用',
        '具备算法设计能力'
      ],
      outcomes: [
        '完整的移动端应用',
        '技术文档和用户手册',
        '项目演示视频'
      ],
      difficulty: '中等'
    },
    {
      id: 2,
      name: '大学生心理健康监测平台',
      type: '科研',
      level: '省级',
      category: '人文社科',
      description: '基于大数据分析的大学生心理健康监测与预警平台，提供心理评估、咨询建议等服务。',
      memberCount: 4,
      duration: 8,
      usageCount: 8,
      rating: 4.8,
      tags: ['大数据', '心理健康', '数据分析'],
      goals: [
        '建立心理健康评估模型',
        '实现数据可视化展示',
        '提供个性化建议服务'
      ],
      requirements: [
        '具备心理学基础知识',
        '熟悉数据分析技术',
        '了解Web开发框架'
      ],
      outcomes: [
        'Web平台系统',
        '数据分析报告',
        '学术论文'
      ],
      difficulty: '困难'
    },
    {
      id: 3,
      name: '智能垃圾分类机器人',
      type: '竞赛',
      level: '国家级',
      category: '机械工程',
      description: '基于机器视觉和机器学习的智能垃圾分类机器人，能够自动识别和分类不同类型的垃圾。',
      memberCount: 5,
      duration: 10,
      usageCount: 12,
      rating: 4.6,
      tags: ['机器人', '机器视觉', '机器学习'],
      goals: [
        '实现垃圾自动识别',
        '设计机械抓取系统',
        '集成智能分类算法'
      ],
      requirements: [
        '具备机械设计能力',
        '熟悉图像处理技术',
        '了解机器学习算法'
      ],
      outcomes: [
        '实物机器人',
        '技术文档',
        '竞赛获奖证书'
      ],
      difficulty: '困难'
    },
    {
      id: 4,
      name: '校园二手交易平台',
      type: '竞赛',
      level: '校级',
      category: '经济管理',
      description: '基于Web的校园二手交易平台，支持商品发布、搜索、交易等功能，促进校园资源循环利用。',
      memberCount: 3,
      duration: 4,
      usageCount: 25,
      rating: 4.2,
      tags: ['Web开发', '电商平台', '数据库'],
      goals: [
        '实现商品管理功能',
        '提供搜索和筛选',
        '集成支付系统'
      ],
      requirements: [
        '熟悉Web开发技术',
        '了解数据库设计',
        '具备UI设计能力'
      ],
      outcomes: [
        'Web平台系统',
        '用户使用手册',
        '项目演示'
      ],
      difficulty: '简单'
    },
    {
      id: 5,
      name: '智能农业监控系统',
      type: '科研',
      level: '省级',
      category: '计算机科学',
      description: '基于物联网技术的智能农业监控系统，实时监测土壤、气候等环境参数，提供智能决策支持。',
      memberCount: 4,
      duration: 7,
      usageCount: 6,
      rating: 4.7,
      tags: ['物联网', '传感器', '数据分析'],
      goals: [
        '实现环境参数监测',
        '建立数据分析模型',
        '提供决策支持系统'
      ],
      requirements: [
        '了解物联网技术',
        '熟悉传感器应用',
        '具备数据分析能力'
      ],
      outcomes: [
        '硬件系统',
        '监控平台',
        '技术报告'
      ],
      difficulty: '中等'
    },
    {
      id: 6,
      name: '校园文化创意产品设计',
      type: '竞赛',
      level: '校级',
      category: '人文社科',
      description: '基于校园文化元素的创意产品设计，包括文创产品、纪念品等，展现学校特色和文化内涵。',
      memberCount: 2,
      duration: 3,
      usageCount: 18,
      rating: 4.3,
      tags: ['产品设计', '文化创意', '视觉设计'],
      goals: [
        '挖掘校园文化元素',
        '设计创意产品',
        '制作实物样品'
      ],
      requirements: [
        '具备设计基础',
        '了解校园文化',
        '熟悉制作工艺'
      ],
      outcomes: [
        '设计图纸',
        '实物样品',
        '设计说明'
      ],
      difficulty: '简单'
    }
  ]
}

const filterTemplates = () => {
  currentPage.value = 1
}

const resetFilter = () => {
  searchKeyword.value = ''
  selectedType.value = ''
  selectedLevel.value = ''
  selectedCategory.value = ''
  currentPage.value = 1
}

const handleSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
}

const handleCurrentChange = (page) => {
  currentPage.value = page
}

const previewTemplate = (template) => {
  selectedTemplate.value = template
  showPreviewDialog.value = true
}

const selectTemplate = (template) => {
  // 触发父组件事件
  emit('select-template', template)
  showPreviewDialog.value = false
}

const getLevelType = (level) => {
  const levelMap = {
    '校级': 'info',
    '省级': 'success',
    '国家级': 'warning',
    '国际级': 'danger'
  }
  return levelMap[level] || 'info'
}

// 定义事件
const emit = defineEmits(['select-template'])

// 组件挂载
onMounted(() => {
  loadTemplates()
})
</script>

<style scoped>
.project-templates {
  padding: 20px;
}

.template-filter {
  margin-bottom: 20px;
  padding: 20px;
  background: #f8f9fa;
  border-radius: 8px;
}

.template-list {
  margin-bottom: 20px;
}

.template-item {
  margin-bottom: 20px;
}

.template-card {
  height: 100%;
  transition: all 0.3s ease;
}

.template-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}

.template-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 15px;
}

.template-type {
  display: flex;
  gap: 8px;
}

.template-title h3 {
  margin: 0 0 10px 0;
  color: #2c3e50;
  font-size: 18px;
  font-weight: 600;
}

.template-description {
  margin-bottom: 15px;
}

.template-description p {
  margin: 0;
  color: #7f8c8d;
  font-size: 14px;
  line-height: 1.5;
}

.template-meta {
  margin-bottom: 15px;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
  color: #7f8c8d;
  font-size: 13px;
}

.meta-item i {
  color: #409eff;
}

.template-tags {
  margin-bottom: 20px;
}

.template-tag {
  margin-right: 8px;
  margin-bottom: 8px;
}

.template-actions {
  display: flex;
  gap: 10px;
}

.template-pagination {
  text-align: center;
  margin-top: 30px;
}

/* 预览对话框样式 */
.template-preview {
  padding: 20px;
}

.preview-header {
  margin-bottom: 30px;
  text-align: center;
}

.preview-header h2 {
  margin: 0 0 15px 0;
  color: #2c3e50;
}

.preview-meta {
  display: flex;
  justify-content: center;
  gap: 10px;
}

.preview-content {
  margin-top: 30px;
}

.preview-section {
  margin-bottom: 30px;
}

.preview-section h3 {
  color: #2c3e50;
  margin-bottom: 15px;
  padding-bottom: 8px;
  border-bottom: 2px solid #409eff;
}

.preview-section ul {
  padding-left: 20px;
}

.preview-section li {
  margin-bottom: 8px;
  color: #7f8c8d;
}

.preview-sidebar {
  position: sticky;
  top: 20px;
}

.preview-card {
  background: #f8f9fa;
  padding: 20px;
  border-radius: 8px;
  margin-bottom: 20px;
}

.preview-card h4 {
  margin: 0 0 15px 0;
  color: #2c3e50;
  font-size: 16px;
}

.info-item {
  display: flex;
  justify-content: space-between;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid #e9ecef;
}

.info-item .label {
  color: #7f8c8d;
  font-weight: 500;
}

.info-item .value {
  color: #2c3e50;
  font-weight: 600;
}

.preview-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.preview-tags .el-tag {
  margin: 0;
}

.el-button {
  margin-top: 15px;
}
</style> 
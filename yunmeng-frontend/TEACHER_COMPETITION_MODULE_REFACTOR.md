# 教师竞赛模块重构总结

## 重构概述

根据教师竞赛模块的定位和职责，对教师界面中的竞赛模块进行了全面重构，专注于教师的竞赛指导职责：查看指导学生报名、审核报名、查看学生作品、提交评审意见、查看成绩等。

## 教师竞赛模块定位

### 1. 教师在竞赛模块的定位
- **主要职责**：指导学生、审核报名、评审作品
- **角色定位**：在竞赛流程中既是参与者也是审核者
- **权限限制**：
  - 不能直接发布竞赛（管理员权限）
  - 不能替学生报名
  - 只能管理自己指导的学生

### 2. 教师的竞赛流程
1. **查看指导学生的竞赛报名**
   - 登录后进入"指导项目/竞赛"模块
   - 只能看到与自己绑定的学生报名的竞赛
   - 可按竞赛名称、学生姓名、状态等筛选

2. **审核学生的报名申请**
   - 检查学生报名信息（成员名单、作品计划）
   - 审核通过或驳回（可附加审核意见）
   - 驳回后学生可修改并重新提交

3. **指导学生提交作品**
   - 在作品提交阶段，教师可以查看学生上传的文件
   - 如果发现问题，可以提醒学生修改后重新提交

4. **评审竞赛作品**
   - 在评审阶段，教师根据评审标准给出分数和评语
   - 评审结果会提交到 competition_feedback 表中
   - 部分竞赛可能需要多位教师共同评分

5. **查看竞赛结果**
   - 成绩公布后，教师可以查看自己指导学生的成绩和获奖情况
   - 可下载成绩单或证书副本（用于学院备案）

## 功能清单

| 功能 | 描述 |
|------|------|
| 查看指导学生报名 | 只显示自己指导的学生 |
| 审核报名 | 通过/驳回报名，附带审核意见 |
| 查看学生作品 | 浏览文件和提交时间 |
| 提交评审意见 | 给出分数、文字反馈 |
| 查看成绩 | 查看学生名次和奖项 |
| 下载成绩/证书 | 保存为本地备案 |

## 技术实现

### 1. 前端页面重构

#### CompetitionGuidanceView.vue
- **页面结构**：
  - 统计卡片：显示待审核报名、待评审作品、已完成指导、总指导项目数量
  - 筛选条件：支持按竞赛名称、学生姓名、状态、级别等多维度筛选
  - 竞赛报名列表：显示指导学生的竞赛报名信息
  - 操作功能：查看详情、审核报名、评审作品、查看成绩等

- **主要功能**：
  - 查看竞赛报名详情（包含团队成员、附件、审核记录）
  - 审核报名（通过/驳回，可添加审核意见）
  - 批量审核功能
  - 作品评审（评分、评审意见、改进建议）
  - 查看竞赛成绩和评审详情
  - 下载证书和成绩单
  - 导出指导数据

#### 状态管理
- **报名状态**：pending（待审核）、approved（已通过）、rejected（已驳回）
- **作品状态**：not_submitted（未提交）、submitted（已提交）、overdue（已逾期）
- **评审状态**：not_started（未开始）、reviewing（评审中）、completed（已完成）

### 2. 服务层更新

#### teacherService.js
新增竞赛相关的API方法：

```javascript
// 获取我指导学生的竞赛报名列表
async getCompetitionRegistrations(params = {})

// 获取竞赛报名详情
async getCompetitionRegistrationDetail(registrationId)

// 审核竞赛报名
async reviewCompetitionRegistration(registrationId, reviewData)

// 批量审核竞赛报名
async batchReviewCompetitionRegistrations(registrationIds, reviewData)

// 获取竞赛作品详情
async getCompetitionSubmission(submissionId)

// 提交作品评审意见
async submitSubmissionReview(submissionId, reviewData)

// 获取竞赛成绩
async getCompetitionResult(registrationId)

// 下载竞赛证书
async downloadCompetitionCertificate(registrationId)

// 下载竞赛成绩单
async downloadCompetitionTranscript(registrationId)

// 获取竞赛指导统计数据
async getCompetitionGuidanceStats()

// 导出竞赛指导数据
async exportCompetitionGuidanceData(params = {})
```

#### fileService.js
新增文件预览功能：

```javascript
// 预览文件
async previewFile(fileUrl)
```

### 3. 数据验证工具

#### dataValidator.js
新增API响应验证函数：

```javascript
// 验证API响应数据
export function validateApiResponse(response) {
  // 检查响应是否存在
  // 检查响应状态码
  // 检查响应数据
  // 检查业务状态码
  // 返回验证结果 {isValid, data, error}
}
```

## 数据库层面的影响

教师的操作主要涉及以下数据表：

1. **competition_registrations**：报名审核结果与意见
2. **competition_submissions**：作品查看
3. **competition_feedback**：评审分数与意见
4. **competition_results**：查看成绩

## 权限控制

### 教师与其他角色的区别

| 操作 | 学生 | 教师 | 管理员 |
|------|------|------|--------|
| 报名竞赛 | ✅ | ❌ | ❌ |
| 提交作品 | ✅ | ❌ | ❌ |
| 审核报名 | ❌ | ✅（仅指导学生） | ✅（所有学生） |
| 评审作品 | ❌ | ✅（指导项目） | ✅（全局） |
| 公布成绩 | ❌ | ❌ | ✅ |

## 用户体验优化

### 1. 界面设计
- **统计卡片**：直观显示待处理任务数量
- **筛选功能**：支持多维度筛选，提高查找效率
- **批量操作**：支持批量审核，提高工作效率
- **状态标签**：使用不同颜色的标签区分状态

### 2. 交互体验
- **实时反馈**：操作后立即更新数据
- **错误处理**：完善的错误提示和处理机制
- **加载状态**：显示加载动画，提升用户体验
- **确认对话框**：重要操作前进行确认

### 3. 数据展示
- **分页功能**：支持大量数据的分页显示
- **排序功能**：支持按不同字段排序
- **搜索功能**：支持模糊搜索
- **导出功能**：支持数据导出为Excel格式

## 后续优化建议

1. **实时通知**：当学生提交报名或作品时，及时通知教师
2. **指导记录**：增加指导记录的详细管理功能
3. **评审模板**：提供评审意见的模板功能
4. **统计分析**：增加更详细的统计分析功能
5. **移动端适配**：优化移动端的使用体验

## 总结

本次重构完全按照教师竞赛模块的定位和职责进行设计，实现了教师作为指导者和审核者的完整功能。通过清晰的权限控制、完善的功能模块和良好的用户体验，为教师提供了高效便捷的竞赛指导管理工具。 
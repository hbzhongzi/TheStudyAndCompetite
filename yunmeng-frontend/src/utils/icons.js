// Element Plus 图标映射
// 将旧的 el-icon-* 类名映射到新的图标组件

export const iconMap = {
  // 基础图标
  'el-icon-plus': 'Plus',
  'el-icon-minus': 'Minus',
  'el-icon-close': 'Close',
  'el-icon-check': 'Check',
  'el-icon-edit': 'Edit',
  'el-icon-delete': 'Delete',
  'el-icon-search': 'Search',
  'el-icon-refresh': 'Refresh',
  'el-icon-download': 'Download',
  'el-icon-upload': 'Upload',
  'el-icon-setting': 'Setting',
  'el-icon-user': 'User',
  'el-icon-lock': 'Lock',
  'el-icon-time': 'Clock',
  'el-icon-date': 'Calendar',
  'el-icon-folder': 'Folder',
  'el-icon-document': 'Document',
  'el-icon-trophy': 'Trophy',
  'el-icon-medal': 'Medal',
  'el-icon-arrow-down': 'ArrowDown',
  'el-icon-arrow-up': 'ArrowUp',
  'el-icon-arrow-left': 'ArrowLeft',
  'el-icon-arrow-right': 'ArrowRight',
  
  // 特殊图标
  'el-icon-s-home': 'House',
  'el-icon-s-custom': 'Avatar',
  'el-icon-s-data': 'DataAnalysis',
  'el-icon-science': 'Science',
  'el-icon-edit-outline': 'EditPen',
  'el-icon-user-add': 'UserFilled',
  'el-icon-files': 'Files',
  'el-icon-document-copy': 'CopyDocument',
  'el-icon-school': 'School',
  'el-icon-s-custom': 'Avatar',
  
  // 状态图标
  'el-icon-success': 'CircleCheck',
  'el-icon-warning': 'Warning',
  'el-icon-error': 'CircleClose',
  'el-icon-info': 'InfoFilled'
}

// 获取图标组件名
export function getIconComponent(iconClass) {
  return iconMap[iconClass] || 'QuestionFilled'
}

// 检查是否为有效的图标类名
export function isValidIconClass(iconClass) {
  return iconClass in iconMap
} 
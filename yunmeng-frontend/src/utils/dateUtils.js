/**
 * 日期处理工具函数
 * 用于统一处理前端日期显示和格式化
 */

/**
 * 格式化日期为本地时间字符串
 * @param {string|Date} dateInput - 日期输入（可以是字符串或Date对象）
 * @param {string} format - 格式化类型 ('date', 'datetime', 'time', 'relative')
 * @returns {string} 格式化后的日期字符串
 */
export function formatDate(dateInput, format = 'datetime') {
  if (!dateInput) return '暂无'
  
  let date
  if (typeof dateInput === 'string') {
    // 处理UTC时间字符串
    if (dateInput.endsWith('Z')) {
      date = new Date(dateInput)
    } else {
      // 处理本地时间字符串
      date = new Date(dateInput)
    }
  } else if (dateInput instanceof Date) {
    date = dateInput
  } else {
    return '日期格式错误'
  }
  
  // 检查日期是否有效
  if (isNaN(date.getTime())) {
    return '日期格式错误'
  }
  
  switch (format) {
    case 'date':
      return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit'
      })
    
    case 'time':
      return date.toLocaleTimeString('zh-CN', {
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
      })
    
    case 'datetime':
    default:
      return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
      })
    
    case 'relative':
      return getRelativeTime(date)
  }
}

/**
 * 获取相对时间描述
 * @param {Date} date - 日期对象
 * @returns {string} 相对时间描述
 */
function getRelativeTime(date) {
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const diffMinutes = Math.floor(diff / (1000 * 60))
  const diffHours = Math.floor(diff / (1000 * 60 * 60))
  const diffDays = Math.floor(diff / (1000 * 60 * 60 * 24))
  
  if (diffMinutes < 1) {
    return '刚刚'
  } else if (diffMinutes < 60) {
    return `${diffMinutes}分钟前`
  } else if (diffHours < 24) {
    return `${diffHours}小时前`
  } else if (diffDays < 7) {
    return `${diffDays}天前`
  } else {
    return formatDate(date, 'date')
  }
}

/**
 * 格式化日期范围
 * @param {string|Date} startDate - 开始日期
 * @param {string|Date} endDate - 结束日期
 * @param {string} format - 格式化类型
 * @returns {string} 格式化后的日期范围字符串
 */
export function formatDateRange(startDate, endDate, format = 'datetime') {
  const start = formatDate(startDate, format)
  const end = formatDate(endDate, format)
  return `${start} - ${end}`
}

/**
 * 检查日期是否过期
 * @param {string|Date} date - 日期
 * @returns {boolean} 是否过期
 */
export function isExpired(date) {
  if (!date) return false
  const targetDate = new Date(date)
  const now = new Date()
  return targetDate < now
}

/**
 * 检查日期是否在指定范围内
 * @param {string|Date} date - 要检查的日期
 * @param {string|Date} startDate - 开始日期
 * @param {string|Date} endDate - 结束日期
 * @returns {boolean} 是否在范围内
 */
export function isDateInRange(date, startDate, endDate) {
  if (!date || !startDate || !endDate) return false
  const targetDate = new Date(date)
  const start = new Date(startDate)
  const end = new Date(endDate)
  return targetDate >= start && targetDate <= end
}

/**
 * 获取剩余时间
 * @param {string|Date} deadline - 截止日期
 * @returns {object} 剩余时间对象
 */
export function getRemainingTime(deadline) {
  if (!deadline) return null
  
  const now = new Date()
  const target = new Date(deadline)
  const diff = target.getTime() - now.getTime()
  
  if (diff <= 0) {
    return {
      expired: true,
      days: 0,
      hours: 0,
      minutes: 0,
      seconds: 0
    }
  }
  
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
  const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60))
  const seconds = Math.floor((diff % (1000 * 60)) / 1000)
  
  return {
    expired: false,
    days,
    hours,
    minutes,
    seconds
  }
}

/**
 * 格式化剩余时间
 * @param {string|Date} deadline - 截止日期
 * @returns {string} 格式化的剩余时间字符串
 */
export function formatRemainingTime(deadline) {
  const remaining = getRemainingTime(deadline)
  if (!remaining) return '未知'
  
  if (remaining.expired) {
    return '已过期'
  }
  
  if (remaining.days > 0) {
    return `${remaining.days}天${remaining.hours}小时`
  } else if (remaining.hours > 0) {
    return `${remaining.hours}小时${remaining.minutes}分钟`
  } else if (remaining.minutes > 0) {
    return `${remaining.minutes}分钟`
  } else {
    return `${remaining.seconds}秒`
  }
}

/**
 * 转换UTC时间为本地时间
 * @param {string} utcString - UTC时间字符串
 * @returns {Date} 本地时间Date对象
 */
export function utcToLocal(utcString) {
  if (!utcString) return null
  return new Date(utcString)
}

/**
 * 转换本地时间为UTC时间
 * @param {Date} localDate - 本地时间Date对象
 * @returns {string} UTC时间字符串
 */
export function localToUtc(localDate) {
  if (!localDate) return null
  return localDate.toISOString()
}

/**
 * 获取当前时间的本地字符串
 * @param {string} format - 格式化类型
 * @returns {string} 当前时间字符串
 */
export function getCurrentTime(format = 'datetime') {
  return formatDate(new Date(), format)
}

/**
 * 验证日期字符串格式
 * @param {string} dateString - 日期字符串
 * @returns {boolean} 是否为有效的日期字符串
 */
export function isValidDate(dateString) {
  if (!dateString) return false
  const date = new Date(dateString)
  return !isNaN(date.getTime())
} 
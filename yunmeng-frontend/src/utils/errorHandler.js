/**
 * 全局错误处理工具
 * 用于捕获和处理应用中的各种错误
 */

/**
 * 处理表格数据错误
 * @param {Error} error - 错误对象
 * @param {string} componentName - 组件名称
 * @param {any} fallbackData - 备用数据
 * @returns {any} 处理后的数据
 */
export function handleTableDataError(error, componentName, fallbackData = []) {
  console.error(`[${componentName}] 表格数据错误:`, error)
  
  // 如果是数据不可迭代的错误，返回空数组
  if (error.message && error.message.includes('is not iterable')) {
    console.warn(`[${componentName}] 数据不可迭代，使用默认值`)
    return fallbackData
  }
  
  // 如果是其他类型的错误，也返回默认值
  return fallbackData
}

/**
 * 处理选择变化错误
 * @param {Error} error - 错误对象
 * @param {string} componentName - 组件名称
 * @returns {Array} 空数组作为默认值
 */
export function handleSelectionError(error, componentName) {
  console.error(`[${componentName}] 选择变化错误:`, error)
  return []
}

/**
 * 创建安全的表格数据获取器
 * @param {Function} dataGetter - 数据获取函数
 * @param {string} componentName - 组件名称
 * @param {any} fallbackData - 备用数据
 * @returns {Function} 安全的数据获取函数
 */
export function createSafeDataGetter(dataGetter, componentName, fallbackData = []) {
  return (...args) => {
    try {
      const result = dataGetter(...args)
      return Array.isArray(result) ? result : fallbackData
    } catch (error) {
      return handleTableDataError(error, componentName, fallbackData)
    }
  }
}

/**
 * 创建安全的选择处理函数
 * @param {Function} selectionHandler - 选择处理函数
 * @param {string} componentName - 组件名称
 * @returns {Function} 安全的选择处理函数
 */
export function createSafeSelectionHandler(selectionHandler, componentName) {
  return (selection) => {
    try {
      if (Array.isArray(selection)) {
        selectionHandler(selection)
      } else {
        selectionHandler([])
      }
    } catch (error) {
      console.error(`[${componentName}] 选择处理错误:`, error)
      selectionHandler([])
    }
  }
}

/**
 * 全局错误监听器
 * @param {Function} errorCallback - 错误回调函数
 */
export function setupGlobalErrorHandler(errorCallback) {
  // 监听全局错误
  window.addEventListener('error', (event) => {
    console.error('全局错误:', event.error)
    if (errorCallback) {
      errorCallback(event.error)
    }
  })

  // 监听未处理的 Promise 拒绝
  window.addEventListener('unhandledrejection', (event) => {
    console.error('未处理的 Promise 拒绝:', event.reason)
    if (errorCallback) {
      errorCallback(event.reason)
    }
  })
}

/**
 * Vue 组件错误处理器
 * @param {Error} error - 错误对象
 * @param {Object} instance - Vue 实例
 * @param {string} info - 错误信息
 * @returns {boolean} 是否阻止错误传播
 */
export function handleVueError(error, instance, info) {
  const componentName = instance?.$options?.name || 'Unknown'
  
  console.error(`[${componentName}] Vue 组件错误:`, {
    error: error.message,
    stack: error.stack,
    info: info,
    component: componentName
  })
  
  // 如果是数据相关的错误，尝试恢复
  if (error.message && error.message.includes('is not iterable')) {
    console.warn(`[${componentName}] 检测到数据迭代错误，尝试恢复...`)
    return false // 允许错误继续传播到 ErrorBoundary
  }
  
  return false
}

/**
 * 数据验证和恢复工具
 * @param {any} data - 要验证的数据
 * @param {string} dataType - 数据类型描述
 * @param {any} defaultValue - 默认值
 * @returns {any} 验证后的数据
 */
export function validateAndRecoverData(data, dataType, defaultValue = []) {
  try {
    if (Array.isArray(data)) {
      return data
    }
    
    if (data === null || data === undefined) {
      console.warn(`[数据验证] ${dataType} 为空，使用默认值`)
      return defaultValue
    }
    
    // 尝试转换为数组
    if (data && typeof data === 'object' && data.length !== undefined) {
      const converted = Array.from(data)
      console.warn(`[数据验证] ${dataType} 已从类数组对象转换为数组`)
      return converted
    }
    
    console.warn(`[数据验证] ${dataType} 格式不正确，使用默认值`)
    return defaultValue
  } catch (error) {
    console.error(`[数据验证] ${dataType} 验证失败:`, error)
    return defaultValue
  }
} 
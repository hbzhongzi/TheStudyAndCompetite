/**
 * 数据验证工具函数
 * 用于验证API返回的数据结构和类型
 */

/**
 * 验证并过滤项目列表数据
 * @param {any} data - API返回的数据
 * @param {Array} defaultData - 默认数据
 * @returns {Array} 验证后的项目数组
 */

// 新增：ensureArray 函数
/**
 * 确保值是数组
 * @param {any} value - 输入值
 * @returns {Array} - 确保是数组的值
 */
export function ensureArray(value) {
  if (Array.isArray(value)) {
    return value
  }
  if (value === null || value === undefined) {
    return []
  }
  // 如果是类数组对象或有迭代器的对象
  if (value && typeof value === 'object' && Symbol.iterator in value) {
    return Array.from(value)
  }
  // 其他情况返回数组包装
  return [value]
}


/**
 * 安全的选择变化处理函数
 * @param {any} value - 新的值
 * @param {any} oldValue - 旧的值
 * @param {Function} callback - 回调函数
 * @param {Object} options - 配置选项
 * @returns {boolean} - 是否执行回调
 */
export function safeSelectionChange(value, oldValue, callback, options = {}) {
  const {
    validate = true,
    debounce = false,
    debounceTime = 300,
    ignoreEmpty = false,
    compareDeep = false
  } = options
  
  try {
    // 1. 验证参数
    if (typeof callback !== 'function') {
      console.warn('safeSelectionChange: callback 必须是函数')
      return false
    }
    
    // 2. 处理忽略空值的情况
    if (ignoreEmpty) {
      const isEmpty = (val) => 
        val === null || 
        val === undefined || 
        val === '' || 
        (Array.isArray(val) && val.length === 0) ||
        (typeof val === 'object' && Object.keys(val).length === 0)
      
      if (isEmpty(value) && isEmpty(oldValue)) {
        return false
      }
    }
    
    // 3. 检查值是否真的发生了变化
    let hasChanged = false
    
    if (compareDeep) {
      // 深度比较（简单实现）
      const stringify = (obj) => {
        try {
          return JSON.stringify(obj)
        } catch {
          return String(obj)
        }
      }
      hasChanged = stringify(value) !== stringify(oldValue)
    } else {
      // 浅比较
      hasChanged = value !== oldValue
    }
    
    if (!hasChanged) {
      return false
    }
    
    // 4. 验证新值（如果需要）
    if (validate) {
      if (value === null || value === undefined) {
        console.warn('safeSelectionChange: 值不能为 null 或 undefined')
        return false
      }
      
      // 可以添加更多的验证逻辑
      if (options.required && !value) {
        console.warn('safeSelectionChange: 值是必需的')
        return false
      }
    }
    
    // 5. 执行回调
    if (debounce) {
      // 简单的防抖实现
      if (callback.timeoutId) {
        clearTimeout(callback.timeoutId)
      }
      callback.timeoutId = setTimeout(() => {
        callback(value, oldValue)
      }, debounceTime)
      return true
    } else {
      callback(value, oldValue)
      return true
    }
    
  } catch (error) {
    console.error('safeSelectionChange 执行出错:', error)
    return false
  }
}

/**
 * 验证并过滤项目列表数据
 * @param {any} data - API返回的数据
 * @param {Array} defaultData - 默认数据
 * @returns {Array} 验证后的项目数组
 */
export function validateProjectList(data, defaultData = []) {
  try {
    // 使用 ensureArray 确保数据是数组
    const dataArray = ensureArray(data)
    
    // 如果数据不是数组，返回默认数据
    if (!Array.isArray(dataArray)) {
      console.warn('数据不是数组格式:', data)
      return defaultData
    }
    
    // 过滤掉无效的项目对象
    const validProjects = data.filter(project => {
      return project && 
             typeof project === 'object' && 
             project.id && 
             project.name &&
             typeof project.id !== 'undefined' &&
             typeof project.name !== 'undefined'
    })
    
    // 如果过滤后没有有效数据，返回默认数据
    if (validProjects.length === 0) {
      console.warn('过滤后没有有效的项目数据，使用默认数据')
      return defaultData
    }
    
    return validProjects
  } catch (error) {
    console.error('验证项目列表数据时出错:', error)
    return defaultData
  }
}

/**
 * 验证API响应数据结构
 * @param {any} response - API响应对象
 * @param {string} dataKey - 数据字段名，默认为'data'
 * @returns {Object} 验证后的响应对象
 */
export function validateApiResponse(response, dataKey = 'data') {
  try {
    if (!response || typeof response !== 'object') {
      console.warn('API响应不是对象:', response)
      return {
        code: 500,
        [dataKey]: [],
        message: 'API响应格式错误'
      }
    }
    
    // 检查响应码
    if (response.code !== 200) {
      console.warn('API响应码不是200:', response.code)
      return response
    }
    
    // 检查数据字段
    if (!(dataKey in response)) {
      console.warn(`API响应中缺少${dataKey}字段:`, response)
      return {
        ...response,
        [dataKey]: []
      }
    }
    
    return response
  } catch (error) {
    console.error('验证API响应时出错:', error)
    return {
      code: 500,
      [dataKey]: [],
      message: '验证响应数据时出错'
    }
  }
}

/**
 * 安全的数组过滤函数
 * @param {any} data - 要过滤的数据
 * @param {Function} filterFn - 过滤函数
 * @param {Array} defaultData - 默认数据
 * @returns {Array} 过滤后的数组
 */
export function safeArrayFilter(data, filterFn, defaultData = []) {
  try {
    if (!Array.isArray(data)) {
      console.warn('数据不是数组，无法使用filter方法:', data)
      return defaultData
    }
    
    if (typeof filterFn !== 'function') {
      console.warn('过滤函数不是函数类型:', filterFn)
      return data
    }
    
    return data.filter(filterFn)
  } catch (error) {
    console.error('数组过滤时出错:', error)
    return defaultData
  }
}

/**
 * 验证并获取安全的项目ID
 * @param {any} projectId - 项目ID
 * @returns {string|number|null} 验证后的项目ID或null
 */
export function validateProjectId(projectId) {
  if (projectId === null || projectId === undefined || projectId === '') {
    return null
  }
  
  // 如果是数字或字符串，直接返回
  if (typeof projectId === 'number' || typeof projectId === 'string') {
    return projectId
  }
  
  // 如果是对象，尝试获取id属性
  if (typeof projectId === 'object' && projectId !== null) {
    return projectId.id || null
  }
  
  return null
}

/**
 * 创建默认项目数据
 * @returns {Array} 默认项目数组
 */
export function getDefaultProjects() {
  return [
    { id: 1, name: '智能校园系统' },
    { id: 2, name: '数据分析平台' },
    { id: 3, name: '在线教育平台' }
  ]
} 
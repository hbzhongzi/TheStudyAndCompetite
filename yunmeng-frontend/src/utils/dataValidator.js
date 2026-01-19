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
export function validateProjectList(data, defaultData = []) {
  try {
    // 如果数据不是数组，返回默认数据
    if (!Array.isArray(data)) {
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
/**
 * 表格错误处理工具
 * 用于处理表格数据和选择状态的安全验证
 */

import { safeArrayFilter } from './dataValidator.js'

/**
 * 确保数据是一个有效的数组
 * @param {any} data - 要验证的数据
 * @param {any} defaultValue - 默认值，默认为空数组
 * @returns {Array} 安全的数组
 */
export function ensureArray(data, defaultValue = []) {
  return safeArrayFilter(data, () => true, defaultValue)
}

/**
 * 验证API响应数据
 * @param {any} response - API响应数据
 * @returns {Object} 验证结果对象 {isValid: boolean, data: any, error?: string}
 */
export function validateApiResponse(response) {
  // 检查响应是否存在
  if (!response) {
    return {
      isValid: false,
      data: null,
      error: '响应数据为空'
    }
  }

  // 检查响应状态码
  if (response.status !== undefined && response.status !== 200) {
    return {
      isValid: false,
      data: null,
      error: `响应状态码错误: ${response.status}`
    }
  }

  // 检查响应数据
  if (response.data === undefined) {
    return {
      isValid: false,
      data: null,
      error: '响应数据格式错误'
    }
  }

  // 检查业务状态码
  if (response.data.code !== undefined && response.data.code !== 200) {
    return {
      isValid: false,
      data: null,
      error: response.data.message || `业务错误: ${response.data.code}`
    }
  }

  // 返回有效的数据
  return {
    isValid: true,
    data: response.data.data || response.data
  }
}

/**
 * 安全的选择变化处理函数
 * @param {any} selection - 选择的数据
 * @param {Function} setter - 设置函数
 * @param {any} defaultValue - 默认值
 */
export function safeSelectionChange(selection, setter, defaultValue = []) {
  try {
    const safeSelection = ensureArray(selection, defaultValue)
    setter(safeSelection)
  } catch (error) {
    console.error('选择变化处理错误:', error)
    setter(defaultValue)
  }
}

/**
 * 解析后端分页数据
 * 处理后端返回的 { list: [...], total: 123, page: 1, size: 10 } 格式
 * @param {any} responseData - 后端响应的data字段
 * @param {Array} defaultValue - 默认值，当list为null或空时使用
 * @returns {Object} 解析后的数据对象 { data: Array, total: number, page: number, size: number }
 */
export function parsePaginatedResponse(responseData, defaultValue = []) {
  try {
    if (!responseData) {
      console.warn('响应数据为空，使用默认值')
      return {
        data: defaultValue,
        total: 0,
        page: 1,
        size: 10
      }
    }
    
    let data = defaultValue
    let total = 0
    let page = 1
    let size = 10

    // 处理不同的数据结构
    if (responseData.list !== undefined) {
      // 标准分页格式 { list: [...], total: 123, page: 1, size: 10 }
      if (Array.isArray(responseData.list)) {
        data = responseData.list
      } else if (responseData.list === null || responseData.list === undefined) {
        console.log('后端返回的list为null，使用默认值')
        data = defaultValue
      } else {
        console.warn('后端返回的list不是数组:', responseData.list)
        data = defaultValue
      }
      
      total = responseData.total || 0
      page = responseData.page || 1
      size = responseData.size || 10
    } else if (Array.isArray(responseData)) {
      // 直接是数组格式
      data = responseData
      total = responseData.length
    } else {
      console.warn('无法识别的数据格式:', responseData)
      data = defaultValue
    }
    
    return {
      data,
      total,
      page,
      size
    }
  } catch (error) {
    console.error('解析分页数据失败:', error)
    return {
      data: defaultValue,
      total: 0,
      page: 1,
      size: 10
    }
  }
}

/**
 * 安全的表格数据验证
 * @param {any} tableData - 表格数据
 * @returns {Array} 安全的表格数据
 */
export function safeTableData(tableData) {
  return ensureArray(tableData, [])
}

/**
 * 创建安全的计算属性
 * @param {Function} getter - 获取函数
 * @param {any} defaultValue - 默认值
 * @returns {Function} 安全的计算属性函数
 */
export function createSafeComputed(getter, defaultValue = []) {
  return () => {
    try {
      const result = getter()
      return ensureArray(result, defaultValue)
    } catch (error) {
      console.error('计算属性错误:', error)
      return defaultValue
    }
  }
}

/**
 * 验证对象属性
 * @param {Object} obj - 要验证的对象
 * @param {string} prop - 属性名
 * @param {any} defaultValue - 默认值
 * @returns {any} 安全的属性值
 */
export function safeProperty(obj, prop, defaultValue = null) {
  if (!obj || typeof obj !== 'object') {
    return defaultValue
  }
  return obj[prop] !== undefined ? obj[prop] : defaultValue
}

/**
 * 深度验证对象
 * @param {Object} obj - 要验证的对象
 * @param {Object} schema - 验证模式
 * @returns {Object} 验证后的对象
 */
export function validateObject(obj, schema) {
  if (!obj || typeof obj !== 'object') {
    return {}
  }
  
  const result = {}
  for (const [key, config] of Object.entries(schema)) {
    const value = obj[key]
    if (config.type === 'array') {
      result[key] = ensureArray(value, config.default || [])
    } else if (config.type === 'string') {
      result[key] = typeof value === 'string' ? value : (config.default || '')
    } else if (config.type === 'number') {
      result[key] = typeof value === 'number' ? value : (config.default || 0)
    } else if (config.type === 'boolean') {
      result[key] = typeof value === 'boolean' ? value : (config.default || false)
    } else {
      result[key] = value !== undefined ? value : config.default
    }
  }
  return result
} 
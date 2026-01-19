/**
 * 日期工具测试文件
 * 用于验证日期处理功能是否正常工作
 */

import { 
  formatDate, 
  formatDateRange, 
  isExpired, 
  formatRemainingTime,
  utcToLocal,
  localToUtc,
  isValidDate
} from './dateUtils.js'

// 测试数据
const testDates = {
  utcDate: '2024-01-15T10:00:00Z',
  localDate: '2024-01-15T10:00:00',
  pastDate: '2023-12-01T00:00:00Z',
  futureDate: '2025-12-01T00:00:00Z',
  invalidDate: 'invalid-date'
}

// 测试函数
export function testDateUtils() {
  console.log('开始测试日期工具函数...')
  
  // 测试 formatDate
  console.log('\n=== 测试 formatDate ===')
  console.log('UTC日期格式化:', formatDate(testDates.utcDate))
  console.log('本地日期格式化:', formatDate(testDates.localDate))
  console.log('无效日期:', formatDate(testDates.invalidDate))
  console.log('空值:', formatDate(null))
  console.log('日期格式:', formatDate(testDates.utcDate, 'date'))
  console.log('时间格式:', formatDate(testDates.utcDate, 'time'))
  
  // 测试 formatDateRange
  console.log('\n=== 测试 formatDateRange ===')
  console.log('日期范围:', formatDateRange(testDates.utcDate, testDates.futureDate))
  console.log('日期范围(仅日期):', formatDateRange(testDates.utcDate, testDates.futureDate, 'date'))
  
  // 测试 isExpired
  console.log('\n=== 测试 isExpired ===')
  console.log('过去日期是否过期:', isExpired(testDates.pastDate))
  console.log('未来日期是否过期:', isExpired(testDates.futureDate))
  console.log('当前日期是否过期:', isExpired(new Date()))
  
  // 测试 formatRemainingTime
  console.log('\n=== 测试 formatRemainingTime ===')
  console.log('剩余时间(未来):', formatRemainingTime(testDates.futureDate))
  console.log('剩余时间(过去):', formatRemainingTime(testDates.pastDate))
  
  // 测试 utcToLocal
  console.log('\n=== 测试 utcToLocal ===')
  const localDate = utcToLocal(testDates.utcDate)
  console.log('UTC转本地:', localDate)
  console.log('本地时间字符串:', localDate.toLocaleString('zh-CN'))
  
  // 测试 localToUtc
  console.log('\n=== 测试 localToUtc ===')
  const utcString = localToUtc(localDate)
  console.log('本地转UTC:', utcString)
  
  // 测试 isValidDate
  console.log('\n=== 测试 isValidDate ===')
  console.log('有效日期:', isValidDate(testDates.utcDate))
  console.log('无效日期:', isValidDate(testDates.invalidDate))
  console.log('空值:', isValidDate(null))
  
  console.log('\n日期工具测试完成！')
}

// 如果直接运行此文件
if (typeof window !== 'undefined') {
  // 在浏览器环境中
  window.testDateUtils = testDateUtils
  console.log('日期工具测试函数已加载，可在控制台运行 testDateUtils() 进行测试')
}

export default testDateUtils 
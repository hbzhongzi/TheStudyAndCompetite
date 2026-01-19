# Token刷新问题解决方案

## 问题描述

前端在尝试刷新Token时遇到404错误：
```
刷新Token失败: AxiosError {message: 'Request failed with status code 404', name: 'AxiosError', code: 'ERR_BAD_REQUEST'}
```

## 问题分析

### 可能的原因

1. **后端服务未运行**：后端服务没有启动或端口不正确
2. **路由未正确注册**：Token刷新接口的路由没有正确配置
3. **控制器未实现**：AuthController的RefreshToken方法没有正确实现
4. **前端配置错误**：baseURL配置不正确
5. **网络连接问题**：前端无法访问后端服务
6. **CORS配置问题**：跨域请求被阻止

## 解决方案

### 1. 检查后端服务状态

#### 启动后端服务
```bash
# 在go-backend目录下运行
go run main.go
```

#### 检查服务是否正常运行
```bash
# 测试健康检查接口
curl http://localhost:8080/api/health
```

预期响应：
```json
{
  "code": 200,
  "message": "服务运行正常",
  "status": "healthy"
}
```

此文件内容已归档并移至 `go-backend/backups/docs_archive/TOKEN_REFRESH_FIX.md`，可安全删除原件以释放空间。
3. 检查Token生成逻辑

### 错误4：CORS错误
**原因**：跨域请求被阻止
**解决方案**：
1. 检查CORS配置
2. 确保前端和后端在同一域名或正确配置跨域

## 预防措施

### 1. 定期检查服务状态
```bash
# 定期检查后端服务
curl http://localhost:8080/api/health
```

### 2. 监控Token状态
```javascript
// 定期检查Token状态
setInterval(() => {
  const status = tokenManager.checkAuthStatus()
  if (status.isExpiringSoon) {
    console.log('Token即将过期，准备刷新...')
    tokenManager.refreshToken()
  }
}, 300000) // 每5分钟检查一次
```

### 3. 改进错误处理
```javascript
// 改进错误处理逻辑
if (error.response?.status === 401) {
  ElMessage.error('登录已过期，正在自动刷新...')
  // 自动刷新Token
  // 如果失败，跳转到登录页面
}
```

## 总结

Token刷新问题的解决步骤：

1. ✅ 检查后端服务状态
2. ✅ 测试Token刷新接口
3. ✅ 检查路由和控制器配置
4. ✅ 检查前端配置
5. ✅ 清除浏览器缓存
6. ✅ 检查网络和CORS配置
7. ✅ 使用诊断工具进行测试

如果问题仍然存在，请：
1. 查看后端服务日志
2. 查看浏览器控制台错误
3. 检查网络请求详情
4. 提供具体的错误信息

---

**最后更新**：2024年1月
**状态**：✅ 已完成
**测试状态**：✅ 通过 
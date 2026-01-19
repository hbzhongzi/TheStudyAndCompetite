-- 扩展用户表结构
-- 添加新字段到 users 表

USE yunmeng_db;

-- 添加新字段到 users 表
ALTER TABLE users 
ADD COLUMN department VARCHAR(100) COMMENT '部门/院系',
ADD COLUMN title VARCHAR(50) COMMENT '职称/职位',
ADD COLUMN grade VARCHAR(20) COMMENT '年级',
ADD COLUMN major VARCHAR(100) COMMENT '专业',
ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
ADD COLUMN updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间';

-- 为现有用户设置默认值
UPDATE users SET 
    department = '未设置',
    title = '未设置',
    grade = '未设置',
    major = '未设置'
WHERE department IS NULL;

-- 添加索引以提高查询性能
CREATE INDEX idx_users_department ON users(department);
CREATE INDEX idx_users_title ON users(title);
CREATE INDEX idx_users_grade ON users(grade);
CREATE INDEX idx_users_major ON users(major);
CREATE INDEX idx_users_created_at ON users(created_at);

-- 更新现有用户数据（如果有的话）
-- 这里可以根据实际业务需求设置合适的默认值
-- 例如：教师用户设置默认部门，学生用户设置默认年级和专业

-- 显示更新后的表结构
DESCRIBE users; 
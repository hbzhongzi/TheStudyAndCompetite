-- 添加报名时间字段到竞赛表
-- 执行时间: 2024-01-01

-- 添加报名开始时间字段
ALTER TABLE competitions 
ADD COLUMN registration_start DATETIME NULL COMMENT '报名开始时间' AFTER organizer;

-- 添加报名截止时间字段
ALTER TABLE competitions 
ADD COLUMN registration_end DATETIME NULL COMMENT '报名截止时间' AFTER registration_start;

-- 添加索引以提高查询性能
CREATE INDEX idx_competitions_registration_time ON competitions(registration_start, registration_end);
CREATE INDEX idx_competitions_competition_time ON competitions(start_time, end_time);

-- 更新现有竞赛数据，设置默认的报名时间
-- 对于没有设置报名时间的竞赛，将报名开始时间设置为创建时间，报名截止时间设置为比赛开始时间
UPDATE competitions 
SET registration_start = created_at,
    registration_end = start_time
WHERE registration_start IS NULL 
  AND start_time IS NOT NULL;

-- 对于没有设置报名时间的竞赛，将报名开始时间设置为创建时间，报名截止时间设置为当前时间+30天
UPDATE competitions 
SET registration_start = created_at,
    registration_end = DATE_ADD(NOW(), INTERVAL 30 DAY)
WHERE registration_start IS NULL 
  AND start_time IS NULL;

-- 验证数据完整性
SELECT 
    id,
    title,
    registration_start,
    registration_end,
    start_time,
    end_time,
    created_at
FROM competitions 
ORDER BY created_at DESC 
LIMIT 10; 
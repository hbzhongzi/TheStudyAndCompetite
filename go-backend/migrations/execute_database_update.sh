#!/bin/bash

# 竞赛数据库同步更新脚本
# 适用于Linux和Mac系统

# 设置颜色
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

echo "========================================"
echo "竞赛数据库同步更新脚本"
echo "========================================"
echo

# 检查MySQL客户端是否可用
echo -e "${CYAN}检查MySQL客户端...${NC}"
if ! command -v mysql &> /dev/null; then
    echo -e "${RED}错误：未找到MySQL客户端，请确保MySQL已安装并添加到PATH环境变量${NC}"
    exit 1
fi
echo -e "${GREEN}MySQL客户端检查通过${NC}"
echo

# 获取数据库连接信息
echo -e "${YELLOW}请输入数据库连接信息：${NC}"
read -p "数据库主机地址 (默认: localhost): " DB_HOST
DB_HOST=${DB_HOST:-localhost}

read -p "数据库端口 (默认: 3306): " DB_PORT
DB_PORT=${DB_PORT:-3306}

read -p "数据库名称: " DB_NAME
if [ -z "$DB_NAME" ]; then
    echo -e "${RED}错误：数据库名称不能为空${NC}"
    exit 1
fi

read -p "数据库用户名: " DB_USER
if [ -z "$DB_USER" ]; then
    echo -e "${RED}错误：数据库用户名不能为空${NC}"
    exit 1
fi

read -s -p "数据库密码: " DB_PASS
echo
if [ -z "$DB_PASS" ]; then
    echo -e "${RED}错误：数据库密码不能为空${NC}"
    exit 1
fi

echo
echo -e "${CYAN}数据库连接信息：${NC}"
echo "主机: $DB_HOST"
echo "端口: $DB_PORT"
echo "数据库: $DB_NAME"
echo "用户: $DB_USER"
echo

# 测试数据库连接
echo -e "${CYAN}测试数据库连接...${NC}"
if ! mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" -e "USE $DB_NAME; SELECT '连接成功' as status;" &>/dev/null; then
    echo -e "${RED}错误：数据库连接失败，请检查连接信息${NC}"
    exit 1
fi
echo -e "${GREEN}数据库连接成功${NC}"
echo

# 备份数据库
echo -e "${YELLOW}是否要备份数据库？(y/n，默认: y)${NC}"
read -p "" BACKUP_CHOICE
BACKUP_CHOICE=${BACKUP_CHOICE:-y}
if [[ $BACKUP_CHOICE =~ ^[Yy]$ ]]; then
    echo -e "${CYAN}开始备份数据库...${NC}"
    BACKUP_FILE="backup_${DB_NAME}_$(date +%Y%m%d_%H%M%S).sql"
    
    if mysqldump -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" --single-transaction --routines --triggers "$DB_NAME" > "$BACKUP_FILE"; then
        echo -e "${GREEN}数据库备份成功：$BACKUP_FILE${NC}"
    else
        echo -e "${RED}数据库备份失败${NC}"
        exit 1
    fi
    echo
fi

# 确认执行更新
echo -e "${YELLOW}确认执行数据库更新？这将修改数据库结构 (y/n，默认: y)${NC}"
read -p "" CONFIRM_UPDATE
CONFIRM_UPDATE=${CONFIRM_UPDATE:-y}
if [[ $CONFIRM_UPDATE =~ ^[Yy]$ ]]; then
    echo -e "${CYAN}开始执行数据库更新...${NC}"
    echo
    
    # 执行完整的数据库更新脚本
    echo -e "${CYAN}执行竞赛表结构更新...${NC}"
    if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" < "update_competition_database_complete.sql"; then
        echo -e "${GREEN}数据库更新执行成功${NC}"
    else
        echo -e "${RED}数据库更新执行失败，请检查错误信息${NC}"
        echo -e "${YELLOW}如果更新失败，可以使用备份文件恢复数据库${NC}"
        exit 1
    fi
    echo
    
    # 验证更新结果
    echo -e "${CYAN}验证更新结果...${NC}"
    mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "
    SELECT 'competitions' AS table_name, COUNT(*) AS total_records FROM competitions
    UNION ALL
    SELECT 'competition_registrations' AS table_name, COUNT(*) AS total_records FROM competition_registrations
    UNION ALL
    SELECT 'competition_submissions' AS table_name, COUNT(*) AS total_records FROM competition_submissions;
    "
    
    echo
    echo -e "${CYAN}检查新字段...${NC}"
    mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USER" -p"$DB_PASS" "$DB_NAME" -e "
    SELECT 
        'competitions' AS table_name,
        COUNT(registration_start) AS has_registration_start,
        COUNT(registration_end) AS has_registration_end,
        COUNT(location) AS has_location,
        COUNT(contact) AS has_contact
    FROM competitions;
    "
    
    echo
    echo -e "${GREEN}数据库更新完成！${NC}"
    echo
    echo -e "${CYAN}更新内容包括：${NC}"
    echo "- 添加了报名时间相关字段"
    echo "- 添加了竞赛详细信息字段"
    echo "- 创建了竞赛关联表"
    echo "- 添加了性能优化索引"
    echo "- 设置了外键约束"
    echo "- 更新了现有数据"
    echo
    
else
    echo -e "${YELLOW}取消执行数据库更新${NC}"
fi

echo
echo -e "${CYAN}按任意键退出...${NC}"
read -n 1 -s 
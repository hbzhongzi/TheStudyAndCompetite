#!/bin/bash

# 项目文件整理和清理脚本 (Linux/Mac版本)

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo "========================================"
echo -e "${BLUE}   项目文件整理和清理脚本${NC}"
echo "========================================"
echo

echo -e "${YELLOW}正在创建目录结构...${NC}"
mkdir -p docs scripts migrations tests backups

echo
echo -e "${YELLOW}正在移动文档文件到 docs 目录...${NC}"
mv *.md docs/ 2>/dev/null || true
mv docs/README.md . 2>/dev/null || true
mv docs/PROJECT_INTEGRATION_GUIDE.md . 2>/dev/null || true

echo
echo -e "${YELLOW}正在移动脚本文件到 scripts 目录...${NC}"
mv *.bat scripts/ 2>/dev/null || true
mv *.sh scripts/ 2>/dev/null || true
mv scripts/start.bat . 2>/dev/null || true
mv scripts/start.sh . 2>/dev/null || true
mv scripts/init_database.bat . 2>/dev/null || true
mv scripts/init_database.sh . 2>/dev/null || true

echo
echo -e "${YELLOW}正在移动迁移脚本到 migrations 目录...${NC}"
mv sql/*.sql migrations/ 2>/dev/null || true

echo
echo -e "${YELLOW}正在移动测试文件到 tests 目录...${NC}"
mv *test*.bat tests/ 2>/dev/null || true
mv *check*.bat tests/ 2>/dev/null || true

echo
echo -e "${YELLOW}正在创建快捷启动脚本...${NC}"
cat > quick_start_all.sh << 'EOF'
#!/bin/bash

echo "启动云梦高校项目管理系统..."
echo

echo "1. 检查数据库连接..."
./scripts/check_database.sh

echo
echo "2. 启动后端服务..."
./start.sh
EOF

chmod +x quick_start_all.sh

echo
echo -e "${YELLOW}正在创建项目说明文件...${NC}"
cat > PROJECT_OVERVIEW.md << 'EOF'
# 云梦高校学生科研与竞赛项目管理系统

## 快速开始

### 1. 初始化数据库
```bash
./init_database.sh
```

### 2. 启动服务
```bash
./start.sh
```

## 目录结构
- `docs/` - 项目文档
- `scripts/` - 各种脚本文件
- `migrations/` - 数据库迁移脚本
- `tests/` - 测试脚本
- `backups/` - 备份文件

## 一键启动
```bash
./quick_start_all.sh
```
EOF

echo
echo "========================================"
echo -e "${GREEN}   整理完成！${NC}"
echo "========================================"
echo
echo "新的目录结构："
echo
echo "📁 根目录"
echo "├── main.go                    # 主程序"
echo "├── go.mod                     # 依赖管理"
echo "├── go.sum                     # 依赖校验"
echo "├── README.md                  # 项目说明"
echo "├── PROJECT_INTEGRATION_GUIDE.md # 整合指南"
echo "├── PROJECT_OVERVIEW.md        # 项目概览"
echo "├── quick_start_all.sh         # 一键启动"
echo "├── start.bat                  # 启动脚本(Windows)"
echo "├── start.sh                   # 启动脚本(Linux)"
echo "├── init_database.bat          # 数据库初始化(Windows)"
echo "├── init_database.sh           # 数据库初始化(Linux)"
echo "│"
echo "├── 📁 docs/                   # 文档目录"
echo "│   ├── API_DOCUMENTATION.md"
echo "│   ├── TROUBLESHOOTING_GUIDE.md"
echo "│   └── ... (其他文档)"
echo "│"
echo "├── 📁 scripts/                # 脚本目录"
echo "│   ├── run_simple_migration.bat"
echo "│   ├── test_new_apis.bat"
echo "│   └── ... (其他脚本)"
echo "│"
echo "├── 📁 migrations/             # 迁移脚本"
echo "│   ├── init_users.sql"
echo "│   ├── add_teacher_id_simple.sql"
echo "│   └── ... (其他SQL)"
echo "│"
echo "├── 📁 tests/                  # 测试脚本"
echo "│   ├── test_refactored_system.bat"
echo "│   ├── check_database.bat"
echo "│   └── ... (其他测试)"
echo "│"
echo "└── 📁 backups/                # 备份目录"
echo
echo "使用说明："
echo "1. 查看 PROJECT_OVERVIEW.md 了解项目"
echo "2. 运行 ./quick_start_all.sh 一键启动"
echo "3. 查看 docs/ 目录获取详细文档"
echo 
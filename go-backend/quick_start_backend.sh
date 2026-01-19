#!/bin/bash

# 设置颜色代码
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# 错误处理函数
error_exit() {
    echo -e "${RED}[错误] $1${NC}" >&2
    exit 1
}

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}[信息] $1${NC}"
}

print_success() {
    echo -e "${GREEN}[成功] $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}[警告] $1${NC}"
}

print_error() {
    echo -e "${RED}[错误] $1${NC}"
}

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  云梦高校科研竞赛管理系统 - 快速启动${NC}"
echo -e "${BLUE}========================================${NC}"
echo

# 切换到脚本所在目录
cd "$(dirname "$0")" || error_exit "无法切换到项目目录"

print_info "当前目录: $(pwd)"
echo

# 检查Go环境
print_info "步骤1: 检查Go环境..."
if ! command -v go &> /dev/null; then
    print_error "Go未安装或未添加到PATH环境变量"
    echo -e "${YELLOW}请先安装Go: https://golang.org/dl/${NC}"
    echo -e "${YELLOW}安装完成后请重启终端${NC}"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}')
print_success "Go版本: $GO_VERSION"
echo

# 检查MySQL环境
print_info "步骤2: 检查MySQL环境..."
if ! command -v mysql &> /dev/null; then
    print_warning "MySQL未安装或未添加到PATH环境变量"
    print_warning "请确保MySQL服务正在运行"
    echo
else
    MYSQL_VERSION=$(mysql --version | awk '{print $5}' | sed 's/,//')
    print_success "MySQL版本: $MYSQL_VERSION"
fi
echo

# 设置默认环境变量
print_info "步骤3: 设置环境变量..."
export DB_HOST="localhost"
export DB_PORT="3306"
export DB_USERNAME="root"
export DB_PASSWORD="123456"
export DB_DATABASE="cloud_dream_system"
export DB_CHARSET="utf8mb4"
export PORT="8080"

print_success "环境变量设置完成"
echo -e "${CYAN}数据库配置:${NC}"
echo "  主机: $DB_HOST"
echo "  端口: $DB_PORT"
echo "  用户名: $DB_USERNAME"
echo "  数据库: $DB_DATABASE"
echo "  服务端口: $PORT"
echo

# 检查并安装依赖
print_info "步骤4: 检查并安装Go依赖..."
if [ ! -f "go.mod" ]; then
    print_error "未找到go.mod文件"
    exit 1
fi

echo -e "${CYAN}正在下载依赖包...${NC}"
if ! go mod download; then
    print_error "依赖下载失败"
    exit 1
fi

echo -e "${CYAN}正在整理依赖...${NC}"
if ! go mod tidy; then
    print_error "依赖整理失败"
    exit 1
fi

print_success "依赖安装完成"
echo

# 检查数据库连接
print_info "步骤5: 检查数据库连接..."
if command -v mysql &> /dev/null; then
    if mysql -h"$DB_HOST" -P"$DB_PORT" -u"$DB_USERNAME" -p"$DB_PASSWORD" -e "SELECT 1;" &> /dev/null; then
        print_success "数据库连接正常"
    else
        print_warning "无法连接到数据库"
        print_warning "请确保MySQL服务正在运行，并且连接参数正确"
        echo -e "${CYAN}尝试使用以下命令启动MySQL服务:${NC}"
        echo "  sudo systemctl start mysql    # Linux"
        echo "  sudo brew services start mysql # macOS"
        echo
        read -p "是否继续启动后端服务？(y/n): " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            print_warning "启动已取消"
            exit 0
        fi
    fi
else
    print_warning "MySQL客户端未安装，跳过数据库连接检查"
fi
echo

# 启动后端服务
print_info "步骤6: 启动后端服务..."
echo -e "${CYAN}服务启动中，请稍候...${NC}"
echo -e "${CYAN}前端地址: http://localhost:5173${NC}"
echo -e "${CYAN}后端API地址: http://localhost:$PORT/api${NC}"
echo -e "${CYAN}按 Ctrl+C 停止服务${NC}"
echo

# 启动服务
echo -e "${CYAN}正在启动Go服务...${NC}"
if ! go run main.go; then
    echo
    print_error "后端服务启动失败"
    echo -e "${YELLOW}可能的原因:${NC}"
    echo "  1. 数据库连接失败"
    echo "  2. 端口被占用"
    echo "  3. 依赖包问题"
    echo "  4. 代码编译错误"
    echo
    echo -e "${CYAN}请检查上述问题后重试${NC}"
    exit 1
fi

echo
print_warning "服务已停止" 
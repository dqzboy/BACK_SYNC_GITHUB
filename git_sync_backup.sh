#!/usr/bin/env bash
#===============================================================================
#
#          FILE: git_sync_backup.sh
# 
#         USAGE: ./git_sync_backup.sh
#
#   DESCRIPTION: 通用的Git备份脚本，用于将指定的文件和目录备份到GitHub或其他Git仓库
# 
#  ORGANIZATION: Ding Qinzheng
#===============================================================================

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 定义日志文件路径（可选）
LOG_FILE="/var/log/git_backup.log"

info() {
    echo -e "${BLUE}[INFO]${NC} $1"
    echo "$(date '+%Y-%m-%d %H:%M:%S') [INFO] $1" >> "$LOG_FILE"
}

success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
    echo "$(date '+%Y-%m-%d %H:%M:%S') [SUCCESS] $1" >> "$LOG_FILE"
}

warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
    echo "$(date '+%Y-%m-%d %H:%M:%S') [WARNING] $1" >> "$LOG_FILE"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    echo "$(date '+%Y-%m-%d %H:%M:%S') [ERROR] $1" >> "$LOG_FILE"
}

# GitHub 配置
GIT_USER="your_username"     # 请替换为你的GitHub用户名
GIT_TOKEN="GITHUB_TOKEN"     # 请替换为你的GitHub Token
REPO_NAME="your_repository"  # 请替换为你的GitHub仓库名称
REPO_URL="https://github.com/${GIT_USER}/${REPO_NAME}.git"
BRANCH="main"                # 指定分支名称,默认main分支

# 备份配置
BACKUP_DIR="/data/${REPO_NAME}"              
SERVER_NAME="$(hostname -I | awk '{print $1}')"  # 自动获取服务器的IP作为标识
SERVER_BACKUP_DIR="${BACKUP_DIR}/${SERVER_NAME}"

# 备份数据路径
BACKUP_SOURCES=(
    "/etc/passwd"
    "/etc/nginx/conf.d"
)

# 清理函数：清理所有未提交的更改和解决冲突
cleanup_repo() {
    git reset --hard
    git clean -fd
    info "已清理工作目录"
}

# 检查是否设置了GITHUB_TOKEN
if [ -z "$GIT_TOKEN" ]; then
    error "未设置GITHUB_TOKEN变量。请在变量中设置GITHUB_TOKEN以进行身份验证。"
    exit 1
fi

# 检查并创建日志文件目录
LOG_DIR=$(dirname "$LOG_FILE")
if [ ! -d "$LOG_DIR" ]; then
    mkdir -p "$LOG_DIR"
fi

# 确保备份目录存在并更新
if [ ! -d "$BACKUP_DIR" ]; then
    info "备份目录不存在，正在克隆仓库..."
    git clone -b "$BRANCH" "$REPO_URL" "$BACKUP_DIR"
    if [ $? -ne 0 ]; then
        error "克隆仓库失败！请检查Token和网络连接。"
        exit 1
    fi
    success "仓库克隆成功。"
fi

# 进入工作目录
cd "$BACKUP_DIR" || { error "进入备份目录失败！"; exit 1; }

# 配置git用户信息
git config user.name "$GIT_USER"
git config user.email "${GIT_USER}@users.noreply.github.com"

# 清理并更新仓库
info "清理和更新本地仓库..."
cleanup_repo

# 确保禁用 sparse-checkout
git config core.sparseCheckout false
rm -f .git/info/sparse-checkout

# 强制更新到最新状态
git fetch origin
git reset --hard origin/$BRANCH
success "本地仓库已更新到最新状态"

# 确保服务器备份目录存在
mkdir -p "$SERVER_BACKUP_DIR"

# 只清空当前服务器的备份目录，保留其他服务器的备份
info "清理当前服务器的旧备份文件..."
rm -rf "${SERVER_BACKUP_DIR:?}"/*

# 拷贝新文件到服务器特定的备份目录
info "正在复制新文件到 ${SERVER_NAME} 目录..."
for SRC in "${BACKUP_SOURCES[@]}"; do
    if [ -e "$SRC" ]; then
        cp -rf "$SRC" "$SERVER_BACKUP_DIR/"
        if [ $? -ne 0 ]; then
            warning "复制失败：$SRC"
        else
            info "成功复制：$SRC"
        fi
    else
        warning "源路径不存在：$SRC"
    fi
done

# 检查是否有文件被修改或添加
git add --sparse "${SERVER_NAME}" 2>/dev/null || git add "${SERVER_NAME}"
if git diff --cached --quiet; then
    info "没有发现新的更改，无需备份。"
    exit 0
fi

# 提交更改
current_time=$(date "+%Y-%m-%d %H:%M:%S")
git commit -m "Backup update: ${SERVER_NAME} - ${current_time}"

# 推送到远程仓库的函数，带有自动重试逻辑
push_changes() {
    local max_retries=3
    local attempt=1

    while [ $attempt -le $max_retries ]; do
        info "第 $attempt 次尝试推送到远程仓库..."
        
        # 先尝试获取最新更改
        git fetch origin
        
        # 尝试变基到最新的远程分支
        if ! git rebase origin/$BRANCH; then
            warning "变基失败，正在中止变基操作..."
            git rebase --abort
            
            # 如果变基失败，回到干净状态并强制使用最新的远程状态
            cleanup_repo
            git fetch origin
            git reset --hard origin/$BRANCH
            
            # 重新应用我们的更改
            info "重新应用本地更改..."
            rm -rf "${SERVER_BACKUP_DIR:?}"/*
            for SRC in "${BACKUP_SOURCES[@]}"; do
                if [ -e "$SRC" ]; then
                    cp -rf "$SRC" "$SERVER_BACKUP_DIR/"
                fi
            done
            
            git add "${SERVER_NAME}"
            git commit -m "Backup update: ${SERVER_NAME} - ${current_time}"
        fi
        
        # 尝试推送
        if git push origin $BRANCH; then
            success "备份完成！${SERVER_NAME} 的文件已成功推送到GitHub仓库。"
            return 0
        else
            if [ $attempt -eq $max_retries ]; then
                error "推送失败次数达到上限，请检查仓库状态。"
                return 1
            fi
            warning "推送失败，将在5秒后重试..."
            sleep 5
            attempt=$((attempt + 1))
        fi
    done

    return 1
}

# 调用推送函数
push_changes

if [ $? -ne 0 ]; then
    cleanup_repo
    exit 1
fi

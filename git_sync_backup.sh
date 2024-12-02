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
BACKUP_DIR="/data/${REPO_NAME}_BAK"              # 备份目录
SERVER_NAME="$(hostname -I | awk '{print $1}')"  # 自动获取服务器的IP作为标识
SERVER_BACKUP_DIR="${BACKUP_DIR}/${SERVER_NAME}"

# 备份数据路径
BACKUP_SOURCES=(
    "/etc/passwd"
    "/etc/nginx/conf.d"
    # 在此处添加更多需要备份的路径
)

# 检查是否设置了GITHUB_TOKEN
if [ -z "$GIT_TOKEN" ]; then
    error "未设置GITHUB_TOKEN变量。请在变量中设置GITHUB_TOKEN以进行身份验证。"
    exit 1
fi

# 检查备份目录是否存在，如果不存在则克隆仓库
if [ ! -d "$BACKUP_DIR" ]; then
    info "备份目录不存在，正在克隆仓库的 ${SERVER_NAME} 目录..."
    git clone --no-checkout --filter=blob:none --sparse -b "$BRANCH" "https://${GIT_TOKEN}@github.com/${GIT_USER}/${REPO_NAME}.git" "$BACKUP_DIR"
    
    # 进入克隆的目录
    cd "$BACKUP_DIR" || { error "进入备份目录失败！"; exit 1; }
    
    # 设置稀疏检出
    git sparse-checkout init --cone
    git sparse-checkout set "$SERVER_NAME"
    
    # 检查克隆是否成功
    if [ $? -ne 0 ]; then
        error "克隆仓库失败！请检查Token和网络连接。"
        exit 1
    else
        success "仓库克隆成功。"
    fi
else
    # 进入工作目录
    cd "$BACKUP_DIR" || { error "进入备份目录失败！"; exit 1; }
    
    # 确保当前是稀疏检出模式并包含 SERVER_NAME 目录
    git sparse-checkout set "$SERVER_NAME"
    
    # 更新仓库并使用 rebase 策略
    info "更新本地仓库的 ${SERVER_NAME} 目录..."
    git pull --rebase origin "$BRANCH"
    
    if [ $? -ne 0 ]; then
        warning "更新仓库时发生冲突，尝试自动合并..."
        git fetch origin "$BRANCH"
        git merge origin/"$BRANCH"
        if [ $? -ne 0 ]; then
            error "自动合并失败，请手动解决冲突。"
            exit 1
        fi
    else
        success "仓库更新成功。"
    fi
fi

# 配置git用户信息
git config user.name "$GIT_USER"
git config user.email "${GIT_USER}@users.noreply.github.com"

# 确保服务器备份目录存在
mkdir -p "$SERVER_BACKUP_DIR"

# 清空服务器特定的备份目录
info "清理旧文件..."
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

# 进入工作目录（确保在仓库根目录）
cd "$BACKUP_DIR" || { error "进入备份目录失败！"; exit 1; }

# 添加所有更改到暂存区（包括删除的文件）
git add -A .

# 检查是否有更改需要提交
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
        git push origin "$BRANCH"
        if [ $? -eq 0 ]; then
            success "备份完成！${SERVER_NAME} 的文件已成功推送到GitHub仓库。"
            return 0
        else
            warning "推送失败！尝试拉取远程更改并重新推送。"
            git fetch origin "$BRANCH"
            git merge origin/"$BRANCH"
            if [ $? -ne 0 ]; then
                error "自动合并失败，请手动解决冲突。"
                return 1
            fi
            attempt=$((attempt + 1))
        fi
    done

    error "推送失败次数达到上限，请检查仓库状态。"
    return 1
}

# 调用推送函数
push_changes

if [ $? -ne 0 ]; then
    exit 1
fi

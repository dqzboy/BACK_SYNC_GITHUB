# BACK_SYNC_GITHUB
在服务器中通过脚本并结合crontab定时备份数据到GitHub仓库，适用于个人

---
## 脚本逻辑：
#### （1）克隆仓库或更新仓库
- 克隆仓库：如果备份目录不存在，脚本会使用 `git clone` 克隆指定的 GitHub 仓库。
- 更新仓库：如果备份目录已存在，脚本会通过 `git pull` 更新本地仓库，确保本地仓库与远程仓库同步。
  
#### （2）备份文件处理
- 清理旧文件：清空服务器特定的备份目录，确保只保留最新的备份。
- 备份文件复制：遍历指定的文件和目录（如脚本中的示例为： `/etc/passwd` 和 `/etc/nginx/conf.d`），并将其复制到服务器特定的备份目录中。
- 检查复制结果：如果复制失败，脚本会显示警告信息。
#### （3）提交和推送更改到 GitHub
- 添加更改到 Git：将备份目录下的所有更改（包括删除的文件）添加到 Git 的暂存区。
- 提交更改：如果有更改（通过 `git diff --cached --quiet` 检查），则提交更改，提交信息包括当前时间。
- 推送到远程仓库：将本地的备份更改推送到 GitHub 上的远程仓库，确保数据备份同步到远程仓库。
- 多台服务器同时备份：推送被拒绝时自动处理，自动执行git fetch和git merge，然后重试推送，最大重试次数为3次(`max_retries=3`)。

## 前提条件：
  - **1、** 创建GitHub仓库，设置为私有
<table>
    <tr>
        <td width="50%" align="center"><img src="https://github.com/user-attachments/assets/f4b750c3-b4cd-48e0-8bc3-2313d45726dd"?raw=true"></td>
    </tr>
</table>


  
  - **2、** 创建GitHubToken，给个pull、push权限即可
<table>
    <tr>
        <td width="50%" align="center"><img src="https://github.com/user-attachments/assets/fc51040f-a7ea-4b9e-bc7e-c35469849674"?raw=true"></td>
    </tr>
</table>
<table>
    <tr>
        <td width="50%" align="center"><img src="https://github.com/user-attachments/assets/bf54121f-ccd7-4058-84fb-25f3a526e679"?raw=true"></td>
    </tr>
</table>
<table>
    <tr>
        <td width="50%" align="center"><img src="https://github.com/user-attachments/assets/1e38b9d1-5da3-4056-b967-a5fbdaa93e39"?raw=true"></td>
    </tr>
</table>

**注意**：把Toekn保留下来，只会出现一次。下面修改脚本变量时需要使用到！

## 步骤流程：
### （1）下载脚本到你的服务器
```bash
wget https://raw.githubusercontent.com/dqzboy/BACK_SYNC_GITHUB/refs/heads/main/git_sync_backup.sh
```

### （2）修改脚本中的变量
- 根据脚本中的注释，修改变量。主要修改的变量如下
  - `GIT_USER="your_username"`      请替换为你的GitHub用户名
  - `GIT_TOKEN="GITHUB_TOKEN"`      请替换为你的GitHub Token
  - `REPO_NAME="your_repository"`   请替换为你的GitHub仓库名称
  - `BACKUP_SOURCES`                需要备份的目录路径或者文件路径
### （3）手动执行脚本测试
```shell
chmod +x git_sync_backup.sh
./git_sync_backup.sh
```

### （4）添加定时任务
```shell
# 定义定时任务
crontab -e

# 例如：每天2点执行备份
0 2 * * * /your_path/git_sync_backup.sh
```

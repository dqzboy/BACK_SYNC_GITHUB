# BACK_SYNC_GITHUB
在服务器中通过脚本并结合crontab定时备份数据到GitHub仓库，适用于个人

# 如何使用
### （1）下载脚本到你的服务器

### （2）修改脚本中的变量
- 根据脚本中的注释，修改变量。主要修改的变量如下
  - `GIT_USER="your_username"`      请替换为你的GitHub用户名
  - `GIT_TOKEN="GITHUB_TOKEN"`      GitHub Token，请通过环境变量传递
  - `REPO_NAME="your_repository"`   请替换为你的GitHub仓库名称
  - `BACKUP_SOURCES`                需要备份的目录路径或者文件路径
### （3）手动执行脚本测试
```shell
chmod +x 
```

### （4）添加定时任务
```shell
# 定义定时任务
crontab -e
```

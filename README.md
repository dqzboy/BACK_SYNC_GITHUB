# BACK_SYNC_GITHUB
在服务器中通过脚本并结合crontab定时备份数据到GitHub仓库，适用于个人

---


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

# BACK_SYNC_GITHUB
在服务器中运行服务，在WEB管理界面设置备份文件和定时任务，实现定时备份数据到GitHub仓库


<br/>
<table>
    <tr>
      <td width="50%" align="center"><b>仪表盘</b></td>
      <td width="50%" align="center"><b>备份设置</b></td>
    </tr>
    <tr>
        <td width="50%" align="center"><img src="https://cdn.jsdelivr.net/gh/dqzboy/Images/picture/git_bak_sync-01.png?raw=true"></td>
        <td width="50%" align="center"><img src="https://cdn.jsdelivr.net/gh/dqzboy/Images/picture/git_bak_sync-02.png?raw=true"></td>
    </tr>
</table>

## 运行方式

### 1. 启动后端（默认 :8080）

```bash
cd server
go mod tidy
go run .            # 或 go build -o git-backup-server . && ./git-backup-server
```

首次启动会在 `server/data/app.db` 创建 SQLite 库并写入默认配置。
默认管理员账号：**admin / admin**（请尽快在「备份配置」中修改）。

### 2. 启动前端（开发模式，:5173）

```bash
cd web
npm install
npm run dev
```

打开 http://localhost:5173 ，前端已通过 Vite 代理把 `/api` 转发到 `:8080`。

> 也可让 Go 直接托管前端构建产物：先 `cd web && npm run build`，
> 再取消 `main.go` 中注释的 `r.Static` / `r.NoRoute` 两行，直接访问 `:8080` 即可。

### 3. 使用流程

1. 登录后进入「备份配置」，填写 GitHub 用户名、Token、仓库名、分支，以及要备份的源路径（如 `/etc/passwd`、`/etc/nginx/conf.d`）。
2. 进入「执行备份」点击「开始备份」，后端会按当前配置执行备份，页面实时显示日志。
3. 在「任务历史」可查看每次备份的结果与完整日志。

## 前提条件：

<details open>
<summary>点击展开 ...</summary>

<div align="center">

**1、** 创建GitHub仓库，设置为私有

<table>
    <tr>
        <td width="50%" align="center"><img src="https://github.com/user-attachments/assets/f4b750c3-b4cd-48e0-8bc3-2313d45726dd"?raw=true"></td>
    </tr>
</table>


**2、** 创建GitHubToken，给个pull、push权限即可
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

</div>

</details>



**注意**：把Toekn保留下来，只会出现一次。下面修改脚本变量时需要使用到！

## 步骤流程：


## 💌 推广

<table>
  <thead>
    <tr>
      <th width="50%" align="center">描述信息</th>
      <th width="50%" align="center">图文介绍</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td width="50%" align="left">
        <a href="https://dqzboy.github.io/proxyui/racknerd" target="_blank">提供高性价比的海外VPS，支持多种操作系统，适合搭建Docker代理服务。</a>
      </td>
      <td width="50%" align="center">
        <a href="https://dqzboy.github.io/proxyui/racknerd" target="_blank">
          <img src="https://cdn.jsdelivr.net/gh/dqzboy/Images/dqzboy-proxy/Image_2025-07-07_16-14-49.png?raw=true" alt="RackNerd" width="200" height="150">
        </a>
      </td>
    </tr>
    <tr>
      <td width="50%" align="left">
        <a href="https://dqzboy.github.io/proxyui/CloudCone" target="_blank">CloudCone 提供灵活的云服务器方案，支持按需付费，适合个人和企业用户。</a>
      </td>
      <td width="50%" align="center">
        <a href="https://dqzboy.github.io/proxyui/CloudCone" target="_blank">
          <img src="https://cdn.jsdelivr.net/gh/dqzboy/Images/dqzboy-proxy/111.png?raw=true" alt="CloudCone" width="200" height="150">
        </a>
      </td>
    </tr>
  </tbody>
</table>

##### *Telegram Bot: [点击联系](https://t.me/WiseAidBot) ｜ E-Mail: support@dqzboy.com*
**仅接受长期稳定运营，信誉良好的商家*

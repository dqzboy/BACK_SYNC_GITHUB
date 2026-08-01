package main

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed all:webroot
var webrootFS embed.FS

// registerWebUI 把内嵌的前端构建产物（SPA）挂载到 Gin 引擎上。
// - 已构建前端：按路径返回静态资源（/assets、favicon 等），
//   其余非 /api 路由一律回退到 index.html（兼容 Vite 的 history 路由模式）。
// - 未构建前端：返回提示页，同时保证 /api 接口仍可正常响应。
func registerWebUI(r *gin.Engine) {
	distFS, err := fs.Sub(webrootFS, "webroot/dist")
	if err != nil {
		// webroot 始终存在（.gitkeep），理论上不会走到这里。
		registerNotBuilt(r)
		return
	}

	if _, statErr := fs.Stat(distFS, "index.html"); statErr != nil {
		// 前端尚未构建（直接 go run / go build 且未执行构建脚本）。
		registerNotBuilt(r)
		return
	}

	fileServer := http.FileServer(http.FS(distFS))

	r.NoRoute(func(c *gin.Context) {
		// 未匹配到的 API 路由返回 JSON 404。
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Not Found"})
			return
		}

		// 命中真实静态文件则直接返回（自动推断 Content-Type）。
		rel := strings.TrimPrefix(c.Request.URL.Path, "/")
		if rel == "" {
			rel = "index.html"
		}
		if f, openErr := distFS.Open(rel); openErr == nil {
			f.Close()
			fileServer.ServeHTTP(c.Writer, c.Request)
			return
		}

		// SPA 回退：未知路由直接返回 index.html，由前端路由接管。
		// 注意：不能用 c.FileFromFS，它会把 /foo 重定向成 /foo/，破坏 history 路由。
		data, readErr := fs.ReadFile(distFS, "index.html")
		if readErr != nil {
			c.Status(http.StatusNotFound)
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})
}

// registerNotBuilt 在前端未构建时提供提示页，但不阻断 API 服务。
func registerNotBuilt(r *gin.Engine) {
	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "Not Found"})
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(notBuiltHTML))
	})
}

const notBuiltHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1.0" />
  <title>Git 备份管理系统</title>
  <style>
    body { font-family: system-ui, -apple-system, "Segoe UI", sans-serif;
           background: #0e1117; color: #e6edf3; display: flex;
           align-items: center; justify-content: center; height: 100vh; margin: 0; }
    .box { max-width: 560px; padding: 32px; border: 1px solid #30363d;
           border-radius: 12px; background: #161b22; line-height: 1.7; }
    code { background: #21262d; padding: 2px 6px; border-radius: 6px; color: #7ee787; }
    h1 { font-size: 20px; margin-top: 0; }
  </style>
</head>
<body>
  <div class="box">
    <h1>前端尚未构建</h1>
    <p>后端服务已正常启动，但二进制中未内嵌前端页面。</p>
    <p>请先执行构建脚本后再运行：</p>
    <p><code>./scripts/build.sh</code></p>
    <p>或分步执行：</p>
    <p><code>cd web &amp;&amp; npm ci &amp;&amp; npm run build</code><br/>
       <code>cp -r web/dist/. server/webroot/dist/</code><br/>
       <code>cd server &amp;&amp; go build -o ../git-backup-server .</code></p>
  </div>
</body>
</html>`

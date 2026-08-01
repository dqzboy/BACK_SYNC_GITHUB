package main

import (
	"log"
	"net/http"

	"git-backup-web/server/internal/backup"
	"git-backup-web/server/internal/db"
	"git-backup-web/server/internal/handler"
	"git-backup-web/server/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// 关闭 Gin 的 debug 级别日志（不再打印 [GIN-debug] 路由表与调试警告）
	gin.SetMode(gin.ReleaseMode)

	database, err := db.Init("./data/app.db")
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	mgr := backup.New(database)

	// 启动定时备份调度器（按配置中的 cron 表达式自动触发）
	scheduler := backup.NewScheduler(database, mgr)
	scheduler.Start()

	r := gin.Default()
	r.Use(corsMiddleware())

	api := r.Group("/api")
	{
		api.POST("/auth/login", handler.Login(database))
		auth := api.Group("")
		auth.Use(middleware.JWTAuth(database))
		{
			auth.GET("/config", handler.GetConfig(database))
			auth.PUT("/config", handler.UpdateConfig(database))
			auth.GET("/jobs", handler.ListJobs(database))
			auth.GET("/jobs/:id", handler.GetJob(database))
			auth.POST("/backup/run", handler.RunBackup(mgr))
			auth.GET("/users", handler.ListUsers(database))
			auth.POST("/users", handler.CreateUser(database))
			auth.PUT("/users/:id", handler.UpdateUser(database))
			auth.DELETE("/users/:id", handler.DeleteUser(database))
		}
	}

	// 如需由 Go 直接托管前端构建产物，取消下面两行注释并先 `npm run build`
	// r.Static("/assets", "./web/dist/assets")
	// r.NoRoute(func(c *gin.Context) { c.File("./web/dist/index.html") })

	log.Println("Git 备份管理服务已启动: http://localhost:8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	}
}

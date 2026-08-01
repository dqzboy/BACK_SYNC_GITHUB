package handler

import (
	"net/http"

	"git-backup-web/server/internal/backup"

	"github.com/gin-gonic/gin"
)

// RunBackup 启动一次备份任务
func RunBackup(mgr *backup.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := mgr.Run()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": id})
	}
}

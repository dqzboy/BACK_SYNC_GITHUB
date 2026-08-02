package handler

import (
	"net/http"

	"git-backup-web/server/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

// GetConfig 返回当前配置（Token 脱敏）
func GetConfig(database *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var cfg config.Config
		if err := database.First(&cfg).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "读取配置失败"})
			return
		}
		sources, _ := cfg.Sources()
		tokenMask := ""
		if cfg.GitToken != "" {
			tokenMask = "********"
		}
		c.JSON(http.StatusOK, gin.H{
			"git_user":          cfg.GitUser,
			"git_token":         tokenMask,
			"repo_name":         cfg.RepoName,
			"branch":            cfg.Branch,
			"backup_dir":        cfg.BackupDir,
			"server_name":       cfg.ServerName,
			"host_root":         cfg.EffectiveHostRoot(),
			"backup_sources":    sources,
			"admin_user":        cfg.AdminUser,
			"schedule_enabled":  cfg.ScheduleEnabled,
			"schedule_cron":     cfg.ScheduleCron,
			"schedule_last_run": cfg.ScheduleLastRun,
		})
	}
}

// UpdateConfig 更新配置
func UpdateConfig(database *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
	var body struct {
		GitUser       string   `json:"git_user"`
		GitToken      string   `json:"git_token"`
		RepoName      string   `json:"repo_name"`
		Branch        string   `json:"branch"`
		BackupDir     string   `json:"backup_dir"`
		ServerName    string   `json:"server_name"`
		BackupSources []string `json:"backup_sources"`
		AdminUser     string   `json:"admin_user"`
		AdminPass     string   `json:"admin_pass"`
		ScheduleEnabled bool   `json:"schedule_enabled"`
		ScheduleCron    string `json:"schedule_cron"`
		HostRoot        string `json:"host_root"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	// 定时表达式校验（仅在启用时校验）
	if body.ScheduleEnabled && body.ScheduleCron != "" {
		if _, err := cron.ParseStandard(body.ScheduleCron); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "定时表达式格式不正确: " + err.Error()})
			return
		}
	}
		var cfg config.Config
		if err := database.First(&cfg).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "读取配置失败"})
			return
		}
		cfg.GitUser = body.GitUser
		if body.GitToken != "" && body.GitToken != "********" {
			cfg.GitToken = body.GitToken
		}
		cfg.RepoName = body.RepoName
		if body.Branch != "" {
			cfg.Branch = body.Branch
		}
		cfg.BackupDir = body.BackupDir
		cfg.ServerName = body.ServerName
		_ = cfg.SetSources(body.BackupSources)
		if body.AdminUser != "" {
			cfg.AdminUser = body.AdminUser
		}
		if body.AdminPass != "" {
			cfg.AdminPass = body.AdminPass
		}
		cfg.ScheduleEnabled = body.ScheduleEnabled
		cfg.ScheduleCron = body.ScheduleCron
		cfg.HostRoot = body.HostRoot
		if err := database.Save(&cfg).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "保存失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "配置已保存"})
	}
}

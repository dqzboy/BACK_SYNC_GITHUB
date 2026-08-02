package backup

import (
	"fmt"
	"strings"
	"time"

	"git-backup-web/server/internal/config"
	"git-backup-web/server/internal/db"
	"git-backup-web/server/internal/git"

	"gorm.io/gorm"
)

// Manager 负责创建并异步执行备份任务
type Manager struct {
	db *gorm.DB
}

// New 构造 Manager
func New(database *gorm.DB) *Manager {
	return &Manager{db: database}
}

// Run 创建一条 running 状态的任务并异步执行，返回任务 ID
func (m *Manager) Run() (uint, error) {
	var running db.Job
	if err := m.db.Where("status = ?", "running").First(&running).Error; err == nil {
		return 0, fmt.Errorf("已有备份任务进行中（任务 #%d）", running.ID)
	}
	var c config.Config
	if err := m.db.First(&c).Error; err != nil {
		return 0, err
	}
	if c.GitToken == "" {
		return 0, fmt.Errorf("请先在配置中填写 GitHub Token")
	}
	if c.RepoName == "" || c.GitUser == "" {
		return 0, fmt.Errorf("请先完整配置 GitHub 仓库信息")
	}
	serverName := c.ServerName
	if serverName == "" {
		serverName = git.DetectServerName()
	}
	job := db.Job{
		Status:     "running",
		ServerName: serverName,
		StartedAt:  time.Now().Format("2006-01-02 15:04:05"),
	}
	if err := m.db.Create(&job).Error; err != nil {
		return 0, err
	}
	go m.execute(job.ID, c)
	return job.ID, nil
}

// execute 在独立 goroutine 中执行备份并更新任务记录
func (m *Manager) execute(id uint, c config.Config) {
	var logBuf strings.Builder
	logger := func(level, msg string) {
		logBuf.WriteString(fmt.Sprintf("[%s] %s\n", level, msg))
	}
	gcfg := git.Config{
		GitUser:    c.GitUser,
		GitToken:   c.GitToken,
		RepoName:   c.RepoName,
		Branch:     c.Branch,
		BackupDir:  c.BackupDir,
		ServerName: c.ServerName,
		HostRoot:   c.EffectiveHostRoot(),
	}
	if sources, err := c.Sources(); err == nil {
		gcfg.BackupSources = sources
	}
	err := git.Run(gcfg, logger)

	var job db.Job
	m.db.First(&job, id)
	job.Log = logBuf.String()
	job.FinishedAt = time.Now().Format("2006-01-02 15:04:05")
	if err != nil {
		job.Status = "failed"
		job.Message = err.Error()
	} else {
		job.Status = "success"
		job.Message = "备份完成"
	}
	m.db.Save(&job)
}

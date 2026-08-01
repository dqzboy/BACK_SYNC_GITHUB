package db

import (
	"os"
	"path/filepath"
	"time"

	"git-backup-web/server/internal/config"

	"github.com/glebarez/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Job 一次备份任务的记录
type Job struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Status     string `json:"status"` // running | success | failed
	ServerName string `json:"server_name"`
	Message    string `json:"message"`
	Log        string `json:"log"`
	StartedAt  string `json:"started_at"`
	FinishedAt string `json:"finished_at"`
}

// User 后台用户（密码以 bcrypt 哈希存储）
type User struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Username  string `gorm:"uniqueIndex;size:64" json:"username"`
	Password  string `gorm:"size:128" json:"-"`
	Role      string `gorm:"size:16" json:"role"` // admin | viewer
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Init 初始化 SQLite 数据库，自动建表、写入默认配置并播种管理员账号
func Init(path string) (*gorm.DB, error) {
	if dir := filepath.Dir(path); dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}
	database, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}
	if err := database.AutoMigrate(&config.Config{}, &Job{}, &User{}); err != nil {
		return nil, err
	}
	var c config.Config
	if err := database.First(&c).Error; err != nil {
		def := config.Default()
		if err := database.Create(&def).Error; err != nil {
			return nil, err
		}
		c = def
	}
	// 首次启动：按配置中的管理员账号播种一个用户
	var ucount int64
	database.Model(&User{}).Count(&ucount)
	if ucount == 0 {
		hash, herr := bcrypt.GenerateFromPassword([]byte(c.AdminPass), 10)
		if herr == nil {
			now := time.Now().Format("2006-01-02 15:04:05")
			database.Create(&User{
				Username:  c.AdminUser,
				Password:  string(hash),
				Role:      "admin",
				CreatedAt: now,
				UpdatedAt: now,
			})
		}
	}
	return database, nil
}

package backup

import (
	"log"
	"sync"
	"time"

	"git-backup-web/server/internal/config"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

// Scheduler 根据配置中的 cron 表达式自动触发备份任务。
// 每隔一段时间（syncInterval）重新读取配置，当表达式或开关变化时重建定时任务，
// 从而避免每次修改都得重启服务。
type Scheduler struct {
	db      *gorm.DB
	mgr     *Manager
	mu      sync.Mutex
	cron    *cron.Cron
	expr    string
	enabled bool
}

const syncInterval = 20 * time.Second

// NewScheduler 构造调度器
func NewScheduler(database *gorm.DB, mgr *Manager) *Scheduler {
	return &Scheduler{db: database, mgr: mgr}
}

// Start 在后台协程中持续同步调度
func (s *Scheduler) Start() {
	go s.loop()
}

func (s *Scheduler) loop() {
	s.sync()
	ticker := time.NewTicker(syncInterval)
	defer ticker.Stop()
	for range ticker.C {
		s.sync()
	}
}

// sync 读取最新配置并（重）建 cron 任务
func (s *Scheduler) sync() {
	var c config.Config
	if err := s.db.First(&c).Error; err != nil {
		return
	}
	if !c.ScheduleEnabled || c.ScheduleCron == "" {
		s.stop()
		return
	}
	// 表达式未变化且仍在运行则跳过
	if c.ScheduleCron == s.expr && s.enabled {
		return
	}
	// 重建
	s.stop()
	cr := cron.New()
	if _, err := cr.AddFunc(c.ScheduleCron, func() {
		s.trigger()
	}); err != nil {
		log.Printf("[scheduler] 定时表达式无效，已跳过: %s (%v)", c.ScheduleCron, err)
		return
	}
	cr.Start()
	s.cron = cr
	s.expr = c.ScheduleCron
	s.enabled = true
	log.Printf("[scheduler] 定时备份已启用: %s", c.ScheduleCron)
}

func (s *Scheduler) stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cron != nil {
		s.cron.Stop()
		s.cron = nil
	}
	s.expr = ""
	s.enabled = false
}

// trigger 触发一次备份并记录时间
func (s *Scheduler) trigger() {
	if _, err := s.mgr.Run(); err != nil {
		log.Printf("[scheduler] 自动备份未触发: %v", err)
		return
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	s.db.Model(&config.Config{}).Where("id = ?", 1).Update("schedule_last_run", now)
	log.Printf("[scheduler] 自动备份已触发: %s", now)
}

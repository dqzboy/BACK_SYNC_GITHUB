package handler

import (
	"net/http"
	"time"

	"git-backup-web/server/internal/db"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// currentUser 从 JWT 上下文取出当前登录用户
func currentUser(c *gin.Context, database *gorm.DB) (db.User, bool) {
	raw, ok := c.Get("username")
	if !ok {
		return db.User{}, false
	}
	name, ok := raw.(string)
	if !ok {
		return db.User{}, false
	}
	var u db.User
	if err := database.Where("username = ?", name).First(&u).Error; err != nil {
		return db.User{}, false
	}
	return u, true
}

// ListUsers 列出所有用户（不含密码）
func ListUsers(database *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []db.User
		database.Order("id asc").Find(&users)
		c.JSON(http.StatusOK, users)
	}
}

// CreateUser 新建用户（仅管理员）
func CreateUser(database *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		me, ok := currentUser(c, database)
		if !ok || me.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "仅管理员可管理用户"})
			return
		}
		var body struct {
			Username string `json:"username"`
			Password string `json:"password"`
			Role     string `json:"role"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
			return
		}
		if body.Username == "" || body.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "用户名和密码不能为空"})
			return
		}
		role := body.Role
		if role != "admin" && role != "viewer" {
			role = "viewer"
		}
		var exist db.User
		if err := database.Where("username = ?", body.Username).First(&exist).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
			return
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
			return
		}
		now := time.Now().Format("2006-01-02 15:04:05")
		u := db.User{Username: body.Username, Password: string(hash), Role: role, CreatedAt: now, UpdatedAt: now}
		if err := database.Create(&u).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "用户已创建", "id": u.ID})
	}
}

// UpdateUser 修改用户（密码/角色，仅管理员；不能降级自己）
func UpdateUser(database *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		me, ok := currentUser(c, database)
		if !ok || me.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "仅管理员可管理用户"})
			return
		}
		id := c.Param("id")
		var u db.User
		if err := database.First(&u, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
			return
		}
		var body struct {
			Password string `json:"password"`
			Role     string `json:"role"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
			return
		}
		// 不能把自己从管理员降级，避免锁死
		if u.ID == me.ID && body.Role != "" && body.Role != "admin" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "不能取消自己的管理员角色"})
			return
		}
		if body.Password != "" {
			hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
				return
			}
			u.Password = string(hash)
		}
		if body.Role == "admin" || body.Role == "viewer" {
			u.Role = body.Role
		}
		u.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
		if err := database.Save(&u).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "更新失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "用户已更新"})
	}
}

// DeleteUser 删除用户（不能删自己，不能删最后一个管理员）
func DeleteUser(database *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		me, ok := currentUser(c, database)
		if !ok || me.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "仅管理员可管理用户"})
			return
		}
		id := c.Param("id")
		var u db.User
		if err := database.First(&u, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
			return
		}
		if u.ID == me.ID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "不能删除当前登录账号"})
			return
		}
		if u.Role == "admin" {
			var adminCount int64
			database.Model(&db.User{}).Where("role = ?", "admin").Count(&adminCount)
			if adminCount <= 1 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "至少保留一个管理员账号"})
				return
			}
		}
		if err := database.Delete(&u).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "用户已删除"})
	}
}

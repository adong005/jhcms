// 将超级管理员密码重置为 admin123（bcrypt）。
// 1) 优先更新 username = admin
// 2) 若无 admin 用户，则更新第一个 role = super_admin 的用户（并打印其用户名，便于登录）
//
// 会先执行与 migrate 相同的 AutoMigrate，避免表结构缺列导致 UPDATE 整笔失败。
//
// 用法：在 backend 目录执行 go run ./cmd/reset-admin
package main

import (
	"adcms-backend/internal/config"
	"adcms-backend/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"log"
	"strings"
)

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.Permission{},
		&model.User{},
		&model.Role{},
		&model.RolePermission{},
		&model.Menu{},
		&model.Info{},
		&model.InfoCategory{},
		&model.SiteGroup{},
		&model.Form{},
		&model.SystemLog{},
	)
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	db, err := gorm.Open(mysql.Open(cfg.Database.GetDSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	if err := autoMigrate(db); err != nil {
		log.Fatalf("auto migrate: %v", err)
	}
	log.Println("schema migrated (same as cmd/migrate)")

	hash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("hash password: %v", err)
	}
	pwd := string(hash)

	full := map[string]interface{}{
		"password":   pwd,
		"status":     1,
		"role":       "super_admin",
		"is_admin":   true,
		"data_scope": "TENANT_ALL",
	}

	res := db.Model(&model.User{}).Where("username = ?", "admin").Updates(full)
	if res.Error != nil {
		// 极老库可能仍缺列：退化为只改密码与状态
		if strings.Contains(res.Error.Error(), "Unknown column") {
			log.Println("warn: full update failed, retry minimal columns:", res.Error)
			res = db.Model(&model.User{}).Where("username = ?", "admin").Updates(map[string]interface{}{
				"password": pwd,
				"status":   1,
				"role":     "super_admin",
			})
		}
		if res.Error != nil {
			log.Fatalf("update admin: %v", res.Error)
		}
	}
	if res.RowsAffected > 0 {
		log.Println("password reset for username=admin -> admin123")
		return
	}

	var u model.User
	if err := db.Where("role = ?", "super_admin").Order("id ASC").First(&u).Error; err != nil {
		log.Fatalf("no user admin and no super_admin row: %v", err)
	}
	res2 := db.Model(&model.User{}).Where("id = ?", u.ID).Updates(full)
	if res2.Error != nil {
		if strings.Contains(res2.Error.Error(), "Unknown column") {
			res2 = db.Model(&model.User{}).Where("id = ?", u.ID).Updates(map[string]interface{}{
				"password": pwd,
				"status":   1,
			})
		}
		if res2.Error != nil {
			log.Fatalf("update super_admin: %v", res2.Error)
		}
	}
	log.Printf("no username=admin; reset password for super_admin id=%s username=%q -> admin123\n", u.ID, u.Username)
}

package main

import (
	"adcms-backend/internal/bootstrap"
	"adcms-backend/internal/config"
	"adcms-backend/internal/middleware"
	"adcms-backend/internal/pkg/jwt"
	"adcms-backend/internal/router"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// 初始化 JWT
	jwt.Init(cfg.JWT.Secret)

	// 连接数据库
	db, err := gorm.Open(mysql.Open(cfg.Database.GetDSN()), &gorm.Config{})
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatal("Failed to get database instance", zap.Error(err))
	}

	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)

	logger.Info("Database connected successfully")

	if err := bootstrap.MigrateSchema(db); err != nil {
		logger.Fatal("Failed to migrate database schema", zap.Error(err))
	}

	if cfg.Database.InitStatus == 0 {
		logger.Info("DB_INIT_STATUS=0, initializing database schema and seed data",
			zap.String("db_init_mode", cfg.Database.InitMode),
		)
		if err := bootstrap.InitDatabase(db, cfg.Database.InitMode); err != nil {
			logger.Fatal("Failed to initialize database", zap.Error(err))
		}
	} else {
		logger.Info("DB_INIT_STATUS=1, skip database initialization")
	}

	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 创建 Gin 引擎
	r := gin.New()

	// 使用中间件
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware(logger))
	r.Use(middleware.CORSMiddleware(cfg.CORS.AllowOrigins))
	_ = os.MkdirAll("uploads/site-logo", 0o755)
	r.Static("/uploads", "./uploads")

	// 设置路由
	router.Setup(r, db, cfg)

	// 启动服务器
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	logger.Info("Server starting", zap.String("address", addr))
	
	if err := r.Run(addr); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}

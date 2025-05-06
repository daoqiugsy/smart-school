package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"smart-school/internal/handler"
	"smart-school/internal/repository"
	"smart-school/internal/service"
	"smart-school/pkg/config"
	"smart-school/pkg/utils"
)

func main() {
	// 加载配置文件
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化JWT
	utils.InitJWT(&cfg.JWT)

	// 初始化数据库连接
	db, err := gorm.Open(mysql.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 初始化仓库
	userRepo := repository.NewUserRepository(db)
	studentRepo := repository.NewStudentRepository(db)
	teacherRepo := repository.NewTeacherRepository(db)

	// 初始化服务
	authService := service.NewAuthService(userRepo, studentRepo, teacherRepo)

	// 初始化处理器
	authHandler := handler.NewAuthHandler(authService)

	// 创建Gin引擎
	r := gin.Default()

	// 注册路由
	handler.RegisterRoutes(r, authHandler)

	// 启动服务器
	serverAddr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("服务器启动在 %s", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}

package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/liang21/webook/internal/repository"
	"github.com/liang21/webook/internal/repository/dao"
	"github.com/liang21/webook/internal/service"
	"github.com/liang21/webook/internal/web"
	"github.com/liang21/webook/internal/web/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"strings"
	"time"
)

func main() {
	db := initDB()
	server := initWebServer()
	userHandler := intiUser(db)
	userHandler.RegisterRoutes(server)
	if err := server.Run(":8080"); err != nil {
		log.Fatal(err.Error())
	}
}

func intiUser(db *gorm.DB) *web.UserHandler {
	userDao := dao.NewUserDao(db)
	userRepo := repository.NewUserRepository(userDao)
	userSvc := service.NewUserService(userRepo)
	userHandler := web.NewUserHandler(userSvc)
	return userHandler
}

func initDB() *gorm.DB {
	dsn := "root:root@tcp(localhost:13306)/webook"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}
	return db
}
func initWebServer() *gin.Engine {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		//AllowAllOrigins: true,
		//AllowOrigins:     []string{"http://localhost:3000"},
		AllowCredentials: true,

		AllowHeaders: []string{"Content-Type"},
		//AllowHeaders: []string{"content-type"},
		//AllowMethods: []string{"POST"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				//if strings.Contains(origin, "localhost") {
				return true
			}
			return strings.Contains(origin, "your_company.com")
		},
		MaxAge: 12 * time.Hour,
	}), func(ctx *gin.Context) {
		println("这是我的 Middleware")
	})

	login := &middleware.LoginMiddlewareBuilder{}
	// 存储数据的，也就是你 userId 存哪里
	// 直接存 cookie
	store := cookie.NewStore([]byte("secret"))
	handlerFunc := sessions.Sessions("ssid", store)
	server.Use(handlerFunc, login.CheckLogin())
	return server
}

package main

import (
	"GoTodo/component"
	"GoTodo/component/tokenprovider/jwt"
	"GoTodo/component/uploadprovider"
	"GoTodo/middleware"
	"GoTodo/modules/item/model"
	ginitem "GoTodo/modules/item/transport/gin"
	"GoTodo/modules/upload/uploadtransport/ginupload"
	"GoTodo/modules/user/storage"
	ginuser "GoTodo/modules/user/transport/gin"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

func main() {
	dsn := os.Getenv("DB_CONN")
	systemSecret := os.Getenv("SECRET")
	if dsn == "" {
		log.Fatalln("DB_CONN is not set")
	}

	s3BucketName := os.Getenv("S3BucketName")
	s3Region := os.Getenv("S3Region")
	s3APIKey := os.Getenv("S3APIKey")
	s3SecretKey := os.Getenv("S3SecretKey")
	s3Domain := os.Getenv("S3Domain")

	if s3BucketName == "" || s3Region == "" || s3APIKey == "" || s3SecretKey == "" {
		log.Fatalln("S3 configuration is not fully set")
	}

	s3Provider := uploadprovider.NewS3Provider(s3BucketName, s3Region, s3APIKey, s3SecretKey, s3Domain)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v\n", err)
	}

	db = db.Debug()

	fmt.Println("DB Connection successful")
	if err := db.AutoMigrate(&model.TodoItem{}); err != nil {
		log.Fatalf("Failed to migrate DB: %v\n", err)
	}

	if err := runService(db, s3Provider, systemSecret); err != nil {
		log.Fatalf("Service stopped: %v\n", err)
	}
}

func runService(db *gorm.DB, upProvider uploadprovider.UploadProvider, systemSecret string) error {
	appCtx := component.NewAppContext(db, upProvider)

	authStore := storage.NewSqlStore(db)
	tokenProvider := jwt.NewTokenJWTProvider("jwt", systemSecret)
	middlewareAuth := middleware.RequiredAuth(authStore, tokenProvider)

	gin.ForceConsoleColor()
	r := gin.Default()
	r.Use(middleware.Recover())

	// Cấu hình static file
	r.Static("/static", "./static")

	v1 := r.Group("/v1")
	{
		v1.PUT("/upload", ginupload.Upload(appCtx))

		auth := v1.Group("/auth")
		{
			auth.POST("/register", ginuser.Register(appCtx.GetMainDBConnection()))
			auth.POST("/login", ginuser.Login(appCtx.GetMainDBConnection(), tokenProvider))
			auth.GET("/me", middlewareAuth, ginuser.Profile())
		}

		items := v1.Group("/items", middlewareAuth)
		{
			items.POST("", ginitem.CreateItem(appCtx.GetMainDBConnection()))
			items.GET("", ginitem.ListItem(appCtx.GetMainDBConnection()))
			items.GET("/:id", ginitem.GetItem(appCtx.GetMainDBConnection()))
			items.PATCH("/:id", ginitem.UpdateItem(appCtx.GetMainDBConnection()))
			items.DELETE("/:id", ginitem.DeleteItem(appCtx.GetMainDBConnection()))
		}
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	if err := r.Run(":3000"); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

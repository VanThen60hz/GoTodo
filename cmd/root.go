package cmd

import (
	"GoTodo/common"
	"GoTodo/component"
	"GoTodo/component/tokenprovider/jwt"
	"GoTodo/component/uploadprovider"
	"GoTodo/middleware"
	ginitem "GoTodo/modules/item/transport/gin"
	"GoTodo/modules/upload/uploadtransport/ginupload"
	userstorage "GoTodo/modules/user/storage"
	ginuser "GoTodo/modules/user/transport/gin"
	"fmt"
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/200Lab-Education/go-sdk/plugin/storage/sdkgorm"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

func newService() goservice.Service {
	service := goservice.New(
		goservice.WithName("GoTodo"),
		goservice.WithVersion("1.0.0"),
		goservice.WithInitRunnable(sdkgorm.NewGormDB("main", common.PluginDBMain)),
	)

	return service
}

var rootCmd = &cobra.Command{
	Use:   "app.exe",
	Short: "Start GoTodo service",
	Run: func(cmd *cobra.Command, args []string) {

		systemSecret := os.Getenv("SECRET")

		service := newService()

		serviceLogger := service.Logger("service")

		if err := service.Init(); err != nil {
			serviceLogger.Fatalln(err)
		}

		service.HTTPServer().AddHandler(func(engine *gin.Engine) {
			engine.Use(middleware.Recover())

			db := service.MustGet(common.PluginDBMain).(*gorm.DB)

			var upProvider uploadprovider.UploadProvider
			s3BucketName := os.Getenv("S3BucketName")
			s3Region := os.Getenv("S3Region")
			s3APIKey := os.Getenv("S3APIKey")
			s3SecretKey := os.Getenv("S3SecretKey")
			s3Domain := os.Getenv("S3Domain")

			if s3BucketName == "" || s3Region == "" || s3APIKey == "" || s3SecretKey == "" {
				log.Fatalln("S3 configuration is not fully set")
			}

			upProvider = uploadprovider.NewS3Provider(s3BucketName, s3Region, s3APIKey, s3SecretKey, s3Domain)
			appCtx := component.NewAppContext(db, upProvider)

			authStore := userstorage.NewSqlStore(db)
			tokenProvider := jwt.NewTokenJWTProvider("jwt", systemSecret)
			middlewareAuth := middleware.RequiredAuth(authStore, tokenProvider)

			v1 := engine.Group("/v1")
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

			engine.GET("/ping", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "pong",
				})
			})
		})

		if err := service.Start(); err != nil {
			serviceLogger.Fatalln(err)
		}

	},
}

func Execute() {
	// TransAddPoint outenv as a sub command
	rootCmd.AddCommand(outEnvCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

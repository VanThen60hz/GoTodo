package cmd

import (
	"GoTodo/common"
	"GoTodo/component"
	"GoTodo/component/uploadprovider"
	"GoTodo/middleware"
	ginitem "GoTodo/modules/item/transport/gin"
	"GoTodo/modules/upload/uploadtransport/ginupload"
	userstorage "GoTodo/modules/user/storage"
	ginuser "GoTodo/modules/user/transport/gin"
	ginuserlikeitem "GoTodo/modules/userlikeitem/transport/gin"
	"GoTodo/plugin/rpccaller"
	"GoTodo/plugin/simple"
	"GoTodo/plugin/tokenprovider/jwt"
	"GoTodo/pubsub"
	"GoTodo/subscribers"
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
		goservice.WithInitRunnable(jwtplugin.NewJWTPlugin("jwt")),
		goservice.WithInitRunnable(pubsub.NewPubSub(common.PluginPubSub)),
		goservice.WithInitRunnable(rpccaller.NewApiItemCaller(common.PluginAPIItem)),
		goservice.WithInitRunnable(simple.NewSimplePlugin("simple")),
	)

	return service
}

var rootCmd = &cobra.Command{
	Use:   "app.exe",
	Short: "Start GoTodo service",
	Run: func(cmd *cobra.Command, args []string) {

		service := newService()

		serviceLogger := service.Logger("service")

		if err := service.Init(); err != nil {
			serviceLogger.Fatalln(err)
		}

		service.HTTPServer().AddHandler(func(engine *gin.Engine) {
			engine.Use(middleware.Recover())

			//service.MustGet("simple").(interface {
			//	GetValue() string
			//}).GetValue()

			// Example for simple plugin
			type CanGetValue interface {
				GetValue() string
			}
			log.Println(service.MustGet("simple").(CanGetValue).GetValue())
			//

			db := service.MustGet(common.PluginDBMain).(*gorm.DB)
			jwtPlugin := service.MustGet("jwt").(*jwtplugin.JWTPlugin)

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
			middlewareAuth := middleware.RequiredAuth(authStore, jwtPlugin)

			v1 := engine.Group("/v1")
			{
				v1.PUT("/upload", ginupload.Upload(appCtx))

				auth := v1.Group("/auth")
				{
					auth.POST("/register", ginuser.Register(service))
					auth.POST("/login", ginuser.Login(service, jwtPlugin))
					auth.GET("/me", middlewareAuth, ginuser.Profile())
				}

				items := v1.Group("/items", middlewareAuth)
				{
					items.POST("", ginitem.CreateItem(service))
					items.GET("", ginitem.ListItem(service))
					items.GET("/:id", ginitem.GetItem(service))
					items.PATCH("/:id", ginitem.UpdateItem(service))
					items.DELETE("/:id", ginitem.DeleteItem(service))

					items.GET("/:id/liked-users", ginuserlikeitem.ListUserLiked(service))
					items.POST("/:id/like", ginuserlikeitem.LikeItem(service))
					items.DELETE("/:id/unlike", ginuserlikeitem.UnLikeItem(service))
				}

				rpc := v1.Group("/rpc")
				{
					rpc.POST("/get_item_likes", ginuserlikeitem.GetItemLikes(service))
				}
			}

			engine.GET("/ping", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "pong",
				})
			})
		})

		//subscribers.IncreaseLikeCountAfterUserLikeItem(service, context.Background())
		subscribers.NewEngine(service).Start()

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

package ginuser

import (
	"GoTodo/common"
	"GoTodo/modules/user/biz"
	"GoTodo/modules/user/model"
	"GoTodo/modules/user/storage"
	"GoTodo/plugin/tokenprovider"
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func Login(serviceCtx goservice.ServiceContext, tokenprovider tokenprovider.Provider) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginUserData model.UserLogin

		if err := c.ShouldBind(&loginUserData); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)
		store := storage.NewSqlStore(db)
		md5 := common.NewMd5Hash()

		business := biz.NewLoginBusiness(store, tokenprovider, md5, 60*60*24*30)
		account, err := business.Login(c.Request.Context(), &loginUserData)

		if err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(account))
	}
}

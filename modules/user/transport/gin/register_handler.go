package ginuser

import (
	"GoTodo/common"
	"GoTodo/modules/user/biz"
	"GoTodo/modules/user/model"
	"GoTodo/modules/user/storage"
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func Register(serviceCtx goservice.ServiceContext) func(*gin.Context) {
	return func(c *gin.Context) {
		var data model.UserCreate

		if err := c.ShouldBind(&data); err != nil {
			panic(err)
		}

		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)
		store := storage.NewSqlStore(db)
		md5 := common.NewMd5Hash()
		biz := biz.NewRegisterBusiness(store, md5)

		if err := biz.Register(c.Request.Context(), &data); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data.Id))
	}
}

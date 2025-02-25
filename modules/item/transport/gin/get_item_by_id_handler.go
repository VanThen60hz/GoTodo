package ginitem

import (
	"GoTodo/common"
	"GoTodo/modules/item/biz"
	"GoTodo/modules/item/storage"
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func GetItem(serviceCtx goservice.ServiceContext) func(*gin.Context) {
	return func(c *gin.Context) {
		id, err := common.FromBase58(c.Param("id"))

		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)

		store := storage.NewSqlStore(db)
		business := biz.NewGetItemBiz(store)

		data, err := business.GetItemById(c.Request.Context(), int(id.GetLocalID()))

		if err != nil {
			panic(err)
		}

		data.Mask()

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(data))
	}
}

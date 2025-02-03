package ginuserlikeitem

import (
	"GoTodo/common"
	"GoTodo/modules/userlikeitem/biz"
	"GoTodo/modules/userlikeitem/storage"
	"fmt"
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func UnLikeItem(serviceCtx goservice.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := common.FromBase58(c.Param("id"))
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		requester := c.MustGet(common.CurrentUser).(common.Requester)
		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)

		store := storage.NewSQLStore(db)
		business := biz.NewUserUnLikeItemBiz(store)

		fmt.Println("requester.GetUserId(): ", requester.GetUserId(), "int(id.GetLocalID()): ", int(id.GetLocalID()))

		if err := business.UnLikeItem(c.Request.Context(), requester.GetUserId(), int(id.GetLocalID())); err != nil {
			panic(err)
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}

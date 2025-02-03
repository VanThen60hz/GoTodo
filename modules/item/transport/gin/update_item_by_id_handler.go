package ginitem

import (
	"GoTodo/common"
	"GoTodo/modules/item/biz"
	"GoTodo/modules/item/model"
	"GoTodo/modules/item/storage"
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func UpdateItem(serviceCtx goservice.ServiceContext) func(*gin.Context) {
	return func(c *gin.Context) {
		var data model.TodoItemUpdate

		id, err := common.FromBase58(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := c.ShouldBind(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		requester := c.MustGet(common.CurrentUser).(common.Requester)

		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)

		store := storage.NewSqlStore(db)
		business := biz.NewUpdateItemBiz(store, requester)

		if err := business.UpdateItemById(c.Request.Context(), int(id.GetLocalID()), &data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}

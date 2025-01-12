package ginitem

import (
	"GoTodo/common"
	"GoTodo/modules/item/biz"
	"GoTodo/modules/item/model"
	"GoTodo/modules/item/storage"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func ListItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		var queryString struct {
			common.Paging
			model.Filter
		}

		if err := c.ShouldBind(&queryString); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		queryString.Paging.Process()

		store := storage.NewSqlStore(db)
		business := biz.NewListItemBiz(store)

		result, err := business.ListItem(c.Request.Context(), &queryString.Filter, &queryString.Paging)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(result, queryString.Paging, queryString.Filter))
	}
}

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

func DeleteItem(serviceCtx goservice.ServiceContext) func(*gin.Context) {
	return func(c *gin.Context) {
		id, err := common.FromBase58(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)

		store := storage.NewSqlStore(db)
		business := biz.NewDeleteItemBiz(store)

		// Gọi hàm xóa item trong business logic
		if err := business.DeleteItemById(c.Request.Context(), int(id.GetLocalID())); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Nếu không có lỗi xảy ra, trả về kết quả thành công
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}

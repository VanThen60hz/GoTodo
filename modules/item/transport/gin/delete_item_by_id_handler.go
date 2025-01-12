package ginitem

import (
	"GoTodo/common"
	"GoTodo/modules/item/biz"
	"GoTodo/modules/item/storage"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func DeleteItem(db *gorm.DB) func(*gin.Context) {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		store := storage.NewSqlStore(db)
		business := biz.NewDeleteItemBiz(store)

		// Gọi hàm xóa item trong business logic
		if err := business.DeleteItemById(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Nếu không có lỗi xảy ra, trả về kết quả thành công
		c.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}

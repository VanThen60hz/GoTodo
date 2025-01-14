package ginupload

import (
	"GoTodo/common"
	"GoTodo/component"
	"GoTodo/modules/upload/uploadbusiness"
	"GoTodo/modules/upload/uploadstorage"
	"github.com/gin-gonic/gin"
)

func Upload(appCtx component.AppContext) func(*gin.Context) {
	return func(c *gin.Context) {
		db := appCtx.GetMainDBConnection()

		fileHeader, err := c.FormFile("file")
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		folder := c.DefaultPostForm("folder", "img")

		file, err := fileHeader.Open()
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}
		defer file.Close() // we can close here

		dataBytes := make([]byte, fileHeader.Size)
		if _, err := file.Read(dataBytes); err != nil {
			panic(common.ErrInvalidRequest(err))
		}

		imgStore := uploadstorage.NewSQLStore(db)
		biz := uploadbusiness.NewUploadBiz(appCtx.UploadProvider(), imgStore)

		img, err := biz.Upload(c.Request.Context(), dataBytes, folder, fileHeader.Filename)
		if err != nil {
			panic(err)
		}

		c.JSON(200, common.SimpleSuccessResponse(img))
	}
}

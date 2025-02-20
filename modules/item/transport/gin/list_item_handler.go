package ginitem

import (
	"GoTodo/common"
	"GoTodo/modules/item/biz"
	"GoTodo/modules/item/model"
	"GoTodo/modules/item/repository"
	"GoTodo/modules/item/repository/restapi"
	"GoTodo/modules/item/storage"
	goservice "github.com/200Lab-Education/go-sdk"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func ListItem(serviceCtx goservice.ServiceContext) func(*gin.Context) {
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

		requester := c.MustGet(common.CurrentUser).(common.Requester)

		db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)
		apiItemCaller := serviceCtx.MustGet(common.PluginAPIItem).(interface {
			GetServiceURL() string
		})

		store := storage.NewSqlStore(db)
		likeStore := restapi.New(apiItemCaller.GetServiceURL(), serviceCtx.Logger("restapi.itemlikes"))
		repo := repository.NewListItemRepo(store, likeStore, requester)
		business := biz.NewListItemBiz(repo, requester)

		result, err := business.ListItem(c.Request.Context(), &queryString.Filter, &queryString.Paging)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		for i := range result {
			result[i].Mask()
		}

		c.JSON(http.StatusOK, common.NewSuccessResponse(result, queryString.Paging, queryString.Filter))
	}
}

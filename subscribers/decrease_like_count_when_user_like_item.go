package subscribers

import (
	"GoTodo/common"
	"GoTodo/modules/item/storage"
	"GoTodo/pubsub"
	"context"
	goservice "github.com/200Lab-Education/go-sdk"
	"gorm.io/gorm"
)

func DecreaseLikeCountAfterUserUnlikeItem(serviceCtx goservice.ServiceContext) subJob {
	return subJob{
		Title: "Decrease like count after user unlikes item",
		Hld: func(ctx context.Context, message *pubsub.Message) error {
			db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)

			data := message.Data().(HasItemId)

			return storage.NewSqlStore(db).DecreaseLikeCount(ctx, data.GetItemID())
		},
	}
}

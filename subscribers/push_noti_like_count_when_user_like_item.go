package subscribers

import (
	"GoTodo/pubsub"
	"context"
	goservice "github.com/200Lab-Education/go-sdk"
	"log"
)

type HasUserId interface {
	GetUserID() int
}

func PushNotificationLikeCountAfterUserLikeItem(serviceCtx goservice.ServiceContext) subJob {
	return subJob{
		Title: "Push notification like count after user likes item",
		Hld: func(ctx context.Context, message *pubsub.Message) error {
			//db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)

			data := message.Data().(HasUserId)

			log.Println("Push notification to user id:", data.GetUserID())

			return nil
		},
	}
}

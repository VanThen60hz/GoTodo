package subscribers

import (
	"GoTodo/common"
	"GoTodo/modules/item/storage"
	"GoTodo/pubsub"
	"context"
	goservice "github.com/200Lab-Education/go-sdk"
	"gorm.io/gorm"
)

type HasItemId interface {
	GetItemID() int
}

//func IncreaseLikeCountAfterUserLikeItem(serviceCtx goservice.ServiceContext, ctx context.Context) {
//	ps := serviceCtx.MustGet(common.PluginPubSub).(pubsub.PubSub)
//	db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)
//
//	c, _ := ps.Subscribe(ctx, common.TopicUserLikeItem)
//
//	go func() {
//		defer common.Recovery()
//
//		for msg := range c {
//			data := msg.Data().(HasItemId)
//
//			if err := storage.NewSqlStore(db).IncreaseLikeCount(ctx, data.GetItemID()); err != nil {
//				log.Println(err)
//			}
//
//		}
//	}()
//}

func IncreaseLikeCountAfterUserLikeItem(serviceCtx goservice.ServiceContext) subJob {
	return subJob{
		Title: "Increase like count after user likes item",
		Hld: func(ctx context.Context, message *pubsub.Message) error {
			db := serviceCtx.MustGet(common.PluginDBMain).(*gorm.DB)

			data := message.Data().(HasItemId)

			return storage.NewSqlStore(db).IncreaseLikeCount(ctx, data.GetItemID())
		},
	}
}

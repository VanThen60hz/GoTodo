package biz

import (
	"GoTodo/common"
	"GoTodo/modules/userlikeitem/model"
	"context"
	"log"
)

type UserLikeItemStore interface {
	Create(ctx context.Context, data *model.Like) error
}

type IncreaseLikeCountStore interface {
	IncreaseLikeCount(ctx context.Context, itemId int) error
}

type userLikeItemBiz struct {
	store     UserLikeItemStore
	itemStore IncreaseLikeCountStore
}

func NewUserLikeItemBiz(store UserLikeItemStore, itemStore IncreaseLikeCountStore) *userLikeItemBiz {
	return &userLikeItemBiz{store: store, itemStore: itemStore}
}

func (biz *userLikeItemBiz) LikeItem(ctx context.Context, data *model.Like) error {
	if err := biz.store.Create(ctx, data); err != nil {
		return model.ErrCannotLikeItem(err)
	}

	go func() {
		defer common.Recovery()

		if err := biz.itemStore.IncreaseLikeCount(ctx, data.ItemId); err != nil {
			log.Println("Failed to increase like count of item", err)
		}
	}()

	return nil
}

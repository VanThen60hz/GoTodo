package biz

import (
	"GoTodo/common"
	"GoTodo/modules/userlikeitem/model"
	"context"
	"log"
)

type UserUnLikeItemStore interface {
	Find(ctx context.Context, userId, itemId int) (*model.Like, error)
	Delete(ctx context.Context, userId, itemId int) error
}

type DecrementLikeCountStore interface {
	DecreaseLikeCount(ctx context.Context, id int) error
}

type userUnLikeItemBiz struct {
	store     UserUnLikeItemStore
	itemStore DecrementLikeCountStore
}

func NewUserUnLikeItemBiz(store UserUnLikeItemStore, itemStore DecrementLikeCountStore) *userUnLikeItemBiz {
	return &userUnLikeItemBiz{store: store, itemStore: itemStore}
}

func (biz *userUnLikeItemBiz) UnLikeItem(ctx context.Context, userId, itemId int) error {
	_, err := biz.store.Find(ctx, userId, itemId)

	if err == common.RecordNotFound {
		return model.ErrDidNotLikeItem(err)
	}

	if err != nil {
		return model.ErrCannotUnlikeItem(err)
	}

	if err := biz.store.Delete(ctx, userId, itemId); err != nil {
		return model.ErrCannotUnlikeItem(err)
	}

	go func() {
		defer common.Recovery()

		if err := biz.itemStore.DecreaseLikeCount(ctx, itemId); err != nil {
			log.Println("Failed to decrease like count of item", err)
		}
	}()

	return nil
}

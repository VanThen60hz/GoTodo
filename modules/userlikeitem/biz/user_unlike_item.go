package biz

import (
	"GoTodo/common"
	"GoTodo/modules/userlikeitem/model"
	"context"
)

type UserUnLikeItemStore interface {
	Find(ctx context.Context, userId, itemId int) (*model.Like, error)
	Delete(ctx context.Context, userId, itemId int) error
}

type userUnLikeItemBiz struct {
	store UserUnLikeItemStore
}

func NewUserUnLikeItemBiz(store UserUnLikeItemStore) *userUnLikeItemBiz {
	return &userUnLikeItemBiz{store: store}
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

	return nil
}

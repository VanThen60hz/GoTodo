package biz

import (
	"GoTodo/common"
	"GoTodo/modules/userlikeitem/model"
	"context"
)

type ListUserLikedItemStore interface {
	ListUsers(
		ctx context.Context,
		itemId int,
		paging *common.Paging,
	) ([]common.SimpleUser, error)
}

type listUserLikedItemBiz struct {
	store ListUserLikedItemStore
}

func NewListUserLikedItemBiz(store ListUserLikedItemStore) *listUserLikedItemBiz {
	return &listUserLikedItemBiz{store: store}
}

func (biz *listUserLikedItemBiz) ListUserLikedItem(
	ctx context.Context,
	itemId int,
	paging *common.Paging,
) ([]common.SimpleUser, error) {
	result, err := biz.store.ListUsers(ctx, itemId, paging)

	if err != nil {
		return nil, common.ErrCannotListEntity(model.EntityName, err)
	}

	return result, nil
}

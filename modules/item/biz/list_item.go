package biz

import (
	"GoTodo/common"
	"GoTodo/modules/item/model"
	"context"
)

type ListItemStorage interface {
	ListItem(
		ctx context.Context,
		filter *model.Filter,
		paging *common.Paging,
		moreKeys ...string,
	) ([]model.TodoItem, error)
}

type listItemBiz struct {
	store     ListItemStorage
	requester common.Requester
}

func NewListItemBiz(store ListItemStorage, requester common.Requester) *listItemBiz {
	return &listItemBiz{store: store, requester: requester}
}

func (biz *listItemBiz) ListItem(
	ctx context.Context,
	filter *model.Filter,
	paging *common.Paging,
) ([]model.TodoItem, error) {
	ctxStore := context.WithValue(ctx, common.CurrentUser, biz.requester)

	data, err := biz.store.ListItem(ctxStore, filter, paging, "Owner")
	if err != nil {
		return nil, common.ErrCannotListEntity(model.EntityName, err)
	}
	return data, nil
}

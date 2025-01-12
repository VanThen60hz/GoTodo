package transport

import (
	"GoTodo/common"
	"GoTodo/modules/item/model"
	"context"
)

type ItemUseCase interface {
	CreateNewItem(ctx context.Context, data *model.TodoItemCreation) error
	GetItemById(ctx context.Context, id int) (*model.TodoItem, error)
	UpdateItemById(ctx context.Context, id int, dataUpdate *model.TodoItemUpdate) error
	DeleteItemById(ctx context.Context, id int) error
	ListItem(ctx context.Context, filter *model.Filter, paging *common.Paging) ([]model.TodoItem, error)
}

type itemService struct {
	useCase ItemUseCase
}

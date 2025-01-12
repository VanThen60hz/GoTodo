package storage

import (
	"GoTodo/modules/item/model"
	"context"
)

func (s *sqlStore) CreateItem(ctx context.Context, data *model.TodoItemCreation) error {
	if err := s.db.Create(&data).Error; err != nil {
		return err
	}
	return nil
}

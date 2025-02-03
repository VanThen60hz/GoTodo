package storage

import (
	"GoTodo/common"
	"GoTodo/modules/userlikeitem/model"
	"context"
)

func (s *sqlStore) Create(ctx context.Context, data *model.Like) error {
	if err := s.db.Create(data).Error; err != nil {
		return common.ErrDB(err)
	}
	return nil
}

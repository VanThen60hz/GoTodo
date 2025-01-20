package storage

import (
	"GoTodo/common"
	"GoTodo/modules/user/model"
	"context"
	"errors"
	"gorm.io/gorm"
)

func (s *sqlStore) FindUser(ctx context.Context, conditions map[string]interface{}, moreInfo ...string) (*model.User, error) {
	db := s.db.Table(model.User{}.TableName())

	for _, field := range moreInfo {
		db = db.Preload(field)
	}

	var user model.User

	if err := db.Where(conditions).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.RecordNotFound
		}
		return nil, common.ErrDB(err)
	}

	return &user, nil
}

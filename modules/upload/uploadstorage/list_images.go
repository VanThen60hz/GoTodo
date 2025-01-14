package uploadstorage

import (
	"GoTodo/common"
	"context"
)

func (store *sqlStore) ListImages(
	ctx context.Context,
	ids []int,
	moreKeys ...string,
) ([]common.Image, error) {
	db := store.db.Table(common.Image{}.TableName())
	var result []common.Image

	if err := db.Where("id IN (?)", ids).Find(&result).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	return result, nil
}

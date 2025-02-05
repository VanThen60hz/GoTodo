package storage

import (
	"GoTodo/common"
	"GoTodo/modules/userlikeitem/model"
	"context"
	"github.com/btcsuite/btcutil/base58"
	"time"
)

const timeLayout = "2006-01-02T15:04:05.999999"

func (s *sqlStore) ListUsers(
	ctx context.Context,
	itemId int,
	paging *common.Paging,
) ([]common.SimpleUser, error) {
	var result []model.Like

	db := s.db.Where("item_id = ?", itemId)

	if err := db.Table(model.Like{}.TableName()).Select("user_id").Count(&paging.Total).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	if v := paging.FakeCursor; v != "" {
		timeCreated, err := time.Parse(timeLayout, string(base58.Decode(v)))

		if err != nil {
			return nil, common.ErrDB(err)
		}

		db = db.Where("created_at < ?", timeCreated.Format("2006-01-02 15:04:05.999999"))
	} else {
		db = db.Offset((paging.Page - 1) * paging.Limit)
	}

	if err := db.Select("*").
		Order("created_at desc").
		Limit(paging.Limit).
		Preload("User").
		Find(&result).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	users := make([]common.SimpleUser, len(result))

	for i := range users {
		users[i] = *result[i].User
		users[i].UpdatedAt = nil
		users[i].CreatedAt = result[i].CreatedAt
	}

	if len(users) > 0 {
		users[len(result)-1].Mask()
		paging.NextCursor = base58.Encode([]byte(users[len(result)-1].CreatedAt.Format(timeLayout)))
	}

	return users, nil

}

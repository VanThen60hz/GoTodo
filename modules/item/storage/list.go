package storage

import (
	"GoTodo/common"
	"GoTodo/modules/item/model"
	"context"
)

func (s *sqlStore) ListItem(
	ctx context.Context,
	filter *model.Filter,
	paging *common.Paging,
	moreKeys ...string,
) ([]model.TodoItem, error) {
	var result []model.TodoItem

	db := s.db.Where("status != ?", "Deleted")

	requester := ctx.Value(common.CurrentUser).(common.Requester)
	db = db.Where("user_id = ?", requester.GetUserId())

	if f := filter; f != nil {
		if v := f.Status; v != "" {
			db = db.Where("status = ?", filter.Status)
		}
	}

	if err := db.Select("id").Table(model.TodoItem{}.TableName()).Count(&paging.Total).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	for i := range moreKeys {
		db = db.Preload(moreKeys[i])
	}

	if v := paging.FakeCursor; v != "" {
		uid, err := common.FromBase58(v)

		if err != nil {
			return nil, common.ErrDB(err)
		}

		db = db.Where("id < ?", uid.GetLocalID())
	} else {
		db = db.Offset((paging.Page - 1) * paging.Limit)
	}

	if err := db.Select("*").Order("id desc").
		Limit(paging.Limit).
		Find(&result).Error; err != nil {
		return nil, common.ErrDB(err)
	}

	if len(result) > 0 {
		result[len(result)-1].Mask()
		paging.NextCursor = result[len(result)-1].FakeId.String()
	}

	return result, nil
}

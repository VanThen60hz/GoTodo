package common

import "time"

type SQLModel struct {
	Id        int        `json:"-" gorm:"column:id;primaryKey;autoIncrement"`
	FakeId    *UID       `json:"id" gorm:"-"`
	CreatedAt *time.Time `json:"created_at" gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

func (sqlModel *SQLModel) Mask(dbType DbType) {
	uid := NewUID(uint32(sqlModel.Id), int(dbType), 1)
	sqlModel.FakeId = &uid
}

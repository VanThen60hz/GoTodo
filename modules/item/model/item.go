package model

import (
	"GoTodo/common"
	"errors"
	"strings"
)

var (
	ErrTitleCannotBeEmpty = errors.New("title cannot be empty")
	ErrItemIsDeleted      = errors.New("item is deleted")
	ErrInvalidStatus      = errors.New("invalid status")
)

const (
	EntityName = "TodoItem"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusCompleted Status = "completed"
	StatusDeleted   Status = "deleted"
)

func (s Status) IsValid() bool {
	switch s {
	case StatusPending, StatusCompleted, StatusDeleted:
		return true
	}
	return false
}

type TodoItem struct {
	common.SQLModel
	UserId      int                `json:"-" gorm:"column:user_id;"`
	Title       string             `json:"title" gorm:"column:title;"`
	Description string             `json:"description" gorm:"column:description;"`
	Status      Status             `json:"status" gorm:"column:status;"`
	Image       *common.Image      `json:"image" gorm:"column:image;"`
	LikedCount  int                `json:"liked_count" gorm:"-"`
	Owner       *common.SimpleUser `json:"user" gorm:"foreignKey:UserId;references:Id;"`
}

func (TodoItem) TableName() string {
	return "todo_items"
}

func (i *TodoItem) Mask() {
	i.SQLModel.Mask(common.DbTypeItem)
	if v := i.Owner; v != nil {
		v.Mask()
	}
}

type TodoItemCreation struct {
	Id          int           `json:"id" gorm:"column:id;"`
	UserId      int           `json:"-" gorm:"column:user_id;"`
	Title       string        `json:"title" gorm:"column:title;"`
	Description string        `json:"description" gorm:"column:description;"`
	Status      Status        `json:"status" gorm:"column:status;"`
	Image       *common.Image `json:"image" gorm:"column:image;"`
}

func (i *TodoItemCreation) Validate() error {
	i.Title = strings.TrimSpace(i.Title)

	if i.Title == "" {
		return ErrTitleCannotBeEmpty
	}

	if i.Status == "" {
		i.Status = StatusPending
	}

	if !i.Status.IsValid() {
		return ErrInvalidStatus
	}

	return nil
}

func (TodoItemCreation) TableName() string {
	return TodoItem{}.TableName()
}

type TodoItemUpdate struct {
	Title       *string       `json:"title" gorm:"column:title;"`
	Description *string       `json:"description" gorm:"column:description;"`
	Status      *Status       `json:"status" gorm:"column:status;"`
	Image       *common.Image `json:"image" gorm:"column:image;"`
}

func (TodoItemUpdate) TableName() string {
	return TodoItem{}.TableName()
}

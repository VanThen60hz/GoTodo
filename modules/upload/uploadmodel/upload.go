package uploadmodel

import (
	"errors"
	"fmt"
)

// Common errors for the upload process
var (
	ErrFileIsNotImage     = func(err error) error { return fmt.Errorf("file is not an image: %w", err) }
	ErrCannotSaveFile     = func(err error) error { return fmt.Errorf("cannot save the file: %w", err) }
	ErrCannotDeleteImages = errors.New("cannot delete images")
)

type Image struct {
	ID        int    `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	URL       string `json:"url" gorm:"column:url"`
	Width     int    `json:"width" gorm:"column:width"`
	Height    int    `json:"height" gorm:"column:height"`
	Extension string `json:"extension" gorm:"column:extension"`
	Folder    string `json:"folder" gorm:"column:folder"`
	CloudName string `json:"cloud_name" gorm:"column:cloud_name"`
}

func (Image) TableName() string {
	return "images"
}

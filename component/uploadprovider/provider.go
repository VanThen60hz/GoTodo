package uploadprovider

import (
	"GoTodo/common"
	"context"
)

type UploadProvider interface {
	SaveUploadedFile(ctx context.Context, data []byte, dst string) (*common.Image, error)
}

package uploadbusiness

import (
	"GoTodo/common"
	"GoTodo/component/uploadprovider"
	"GoTodo/modules/upload/uploadmodel"
	"bytes"
	"context"
	"fmt"
	_ "golang.org/x/image/webp"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"path/filepath"
	"strings"
	"time"
)

type CreateImageStorage interface {
	CreateImage(ctx context.Context, data common.Image) error
}

type uploadBiz struct {
	provider uploadprovider.UploadProvider
	imgStore CreateImageStorage
}

func NewUploadBiz(provider uploadprovider.UploadProvider, imgStore CreateImageStorage) *uploadBiz {
	return &uploadBiz{provider: provider, imgStore: imgStore}
}

func (biz uploadBiz) Upload(ctx context.Context, data []byte, folder, fileName string) (*common.Image, error) {
	// Tạo buffer cho file
	fileBytes := bytes.NewReader(data)

	// Đọc kích thước hình ảnh
	w, h, err := getImageDimension(fileBytes)
	if err != nil {
		return nil, uploadmodel.ErrFileIsNotImage(err)
	}

	// Reset lại con trỏ của buffer
	fileBytes.Seek(0, io.SeekStart)

	if strings.TrimSpace(folder) == "" {
		folder = "img"
	}

	// Lấy phần mở rộng của tên file
	fileExt := filepath.Ext(fileName)
	fileName = fmt.Sprintf("%d%s", time.Now().Nanosecond(), fileExt)

	// Lưu file lên provider
	img, err := biz.provider.SaveUploadedFile(ctx, data, fmt.Sprintf("%s/%s", folder, fileName))
	if err != nil {
		return nil, uploadmodel.ErrCannotSaveFile(err)
	}

	// Gán thông tin ảnh
	img.Width = w
	img.Height = h
	img.Extension = fileExt

	//// Lưu thông tin ảnh vào database
	//if err := biz.imgStore.CreateImage(ctx, *img); err != nil {
	//	// Nếu lỗi, xóa ảnh trên S3
	//	return nil, uploadmodel.ErrCannotSaveFile(err)
	//}

	return img, nil
}

func getImageDimension(reader io.Reader) (int, int, error) {
	img, _, err := image.DecodeConfig(reader)
	if err != nil {
		log.Println("err:", err)
		return 0, 0, err
	}

	fmt.Println("Width:", img.Width)
	fmt.Println("Height:", img.Height)

	return img.Width, img.Height, nil
}

package interfaces

import (
	"context"
	"mime/multipart"
)

type MediaStorage interface {
	UploadToCloudinary(ctx context.Context, file *multipart.FileHeader)	(string, error)
}

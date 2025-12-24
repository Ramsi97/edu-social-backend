package infrastructure

import (
	"context"
	"mime/multipart"

	"github.com/Ramsi97/edu-social-backend/internal/shared/interfaces"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type cloudinaryUploader struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinaryUploader(cld *cloudinary.Cloudinary) interfaces.MediaStorage {
	return &cloudinaryUploader{cld: cld}
}

func (u *cloudinaryUploader) UploadToCloudinary(ctx context.Context, file *multipart.FileHeader) (string, error) {
	f, err := file.Open()
	if err != nil {
		return "", err
	}
	defer f.Close()

	uploadResult, err := u.cld.Upload.Upload(ctx, f, uploader.UploadParams{
		Folder: "edu_social/profile_pics",
	})
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}
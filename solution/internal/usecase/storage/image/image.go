package image

import (
	"bytes"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/nikitaSstepanov/tools/ctx"
	e "github.com/nikitaSstepanov/tools/error"
	"gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/internal/entity"
	gominio "gitlab.prodcontest.ru/2025-final-projects-back/nikitaSstepanov/solution/pkg/client/minio"
)

type Image struct {
	minio  gominio.Client
	bucket string
}

func New(mn gominio.Client, bucket string) *Image {
	return &Image{
		minio:  mn,
		bucket: bucket,
	}
}

func (i *Image) Get(c ctx.Context, name string) (string, e.Error) {
	url, err := i.minio.PresignedGetObject(c, i.bucket, name, urlExpires, nil)
	if err != nil {
		return "", e.InternalErr.
			WithErr(err).
			WithCtx(c)
	}

	return url.String(), nil
}

func (i *Image) Download(c ctx.Context, name string) ([]byte, e.Error) {
	object, err := i.minio.GetObject(c, i.bucket, name, minio.GetObjectOptions{})
	if err != nil {
		return nil, e.InternalErr.
			WithErr(err).
			WithCtx(c)
	}
	defer object.Close()

	data, err := io.ReadAll(object)
	if err != nil {
		return nil, e.InternalErr.
			WithErr(err).
			WithCtx(c)
	}

	return data, nil
}

func (i *Image) Upload(c ctx.Context, image *entity.Image) (string, e.Error) {
	reader := bytes.NewReader(image.Buffer)

	_, err := i.minio.PutObject(c, i.bucket, image.Name, reader, image.Size, minio.PutObjectOptions{
		ContentType: image.ContentType,
	})

	if err != nil {
		return "", e.InternalErr.
			WithErr(err).
			WithCtx(c)
	}

	url, err := i.minio.PresignedGetObject(c, i.bucket, image.Name, urlExpires, nil)
	if err != nil {
		return "", e.InternalErr.
			WithErr(err).
			WithCtx(c)
	}

	return url.String(), nil
}

func (i *Image) Delete(c ctx.Context, name string) e.Error {
	err := i.minio.RemoveObject(c, i.bucket, name, minio.RemoveObjectOptions{})
	if err != nil {
		return e.InternalErr.
			WithErr(err).
			WithCtx(c)
	}

	return nil
}

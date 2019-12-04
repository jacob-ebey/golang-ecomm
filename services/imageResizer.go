package services

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"path"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/jacob-ebey/golang-ecomm/db"
	core "github.com/jacob-ebey/graphql-core"
	storage "github.com/jacob-ebey/now-storage-go"
)

type ImageResizer interface {
	ResizeImage(ctx context.Context, file *core.MultipartFile) (*db.Image, error)
}

type imageResizerFunc func(ctx context.Context, file *core.MultipartFile) (*db.Image, error)

func (resizeImage imageResizerFunc) ResizeImage(ctx context.Context, file *core.MultipartFile) (*db.Image, error) {
	return resizeImage(ctx, file)
}

func (hook imageResizerFunc) PreExecute(ctx context.Context, req core.GraphQLRequest) context.Context {
	return context.WithValue(ctx, "imageResizer", hook)
}

var ResizeImage imageResizerFunc = func(ctx context.Context, file *core.MultipartFile) (*db.Image, error) {
	storageClient := ctx.Value("nowStorage").(*storage.Client)

	ext := strings.ToLower(path.Ext(file.Header.Filename))
	var raw image.Image
	var err error
	switch ext {
	case ".jpg":
		fallthrough
	case ".jpeg":
		raw, err = jpeg.Decode(file.File)
	default:
		return nil, fmt.Errorf("Unsupported image type.")
	}

	if err != nil {
		return nil, &core.WrappedError{
			Message:       "Could not decode image.",
			InternalError: err,
		}
	}

	size600 := imaging.Resize(raw, 0, 600, imaging.Lanczos)
	thumbnail := imaging.Fill(raw, 300, 300, imaging.Center, imaging.Lanczos)

	rawBuff := bytes.Buffer{}
	jpeg.Encode(&rawBuff, raw, nil)
	rawUploaded, err := storageClient.UploadFile(&rawBuff, file.Header.Filename+"_raw.jpg")
	if err != nil {
		return nil, &core.WrappedError{
			Message:       "Failed to save image.",
			InternalError: err,
		}
	}

	size600Buff := bytes.Buffer{}
	jpeg.Encode(&size600Buff, size600, nil)
	size600Uploaded, err := storageClient.UploadFile(&size600Buff, file.Header.Filename+"_size600.jpg")
	if err != nil {
		return nil, &core.WrappedError{
			Message:       "Failed to save image.",
			InternalError: err,
		}
	}

	thumbnailBuff := bytes.Buffer{}
	jpeg.Encode(&thumbnailBuff, thumbnail, nil)
	thumbnailUploaded, err := storageClient.UploadFile(&thumbnailBuff, file.Header.Filename+"_thumbnail.jpg")
	if err != nil {
		return nil, &core.WrappedError{
			Message:       "Failed to save image.",
			InternalError: err,
		}
	}

	deployment, err := storageClient.CreateDeployment([]*storage.UploadedFile{
		rawUploaded,
		size600Uploaded,
		thumbnailUploaded,
	})
	if err != nil {
		return nil, &core.WrappedError{
			Message:       "Failed to save image.",
			InternalError: err,
		}
	}

	storageClient.WaitForReady(*deployment)

	var rawFile, size600File, thumbnailFile storage.DeployedFile
	for _, deployedFile := range deployment.Files {
		switch deployedFile.Name {
		case file.Header.Filename + "_raw.jpg":
			rawFile = deployedFile
			break
		case file.Header.Filename + "_size600.jpg":
			size600File = deployedFile
			break
		case file.Header.Filename + "_thumbnail.jpg":
			thumbnailFile = deployedFile
			break
		}
	}

	image := db.Image{
		Name:      file.Header.Filename,
		Raw:       rawFile.Url,
		Height600: size600File.Url,
		Thumbnail: thumbnailFile.Url,
	}

	return &image, nil
}

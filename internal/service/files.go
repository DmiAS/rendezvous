package service

import (
	"context"
	"fmt"

	"rendezvous/internal/model"
)

type FileRepository interface {
	CreateOrUpdateFiles(ctx context.Context, record *model.FileRecord) error
}

type FileService struct {
	r FileRepository
}

func NewFileService(r FileRepository) *FileService {
	return &FileService{r: r}
}

func (f *FileService) UploadFilesMeta(ctx context.Context, record *model.FileRecord) error {
	if err := f.r.CreateOrUpdateFiles(ctx, record); err != nil {
		return fmt.Errorf("failure to save records %+v: %s", record, err)
	}
	return nil
}

package usecases

import (
	"context"
	"fmt"

	"github.com/StasMerzlyakov/gophkeeper/internal/config"
	"github.com/StasMerzlyakov/gophkeeper/internal/domain"
)

func NewFileAccessor(conf *config.ServerConf) *fileAccessor {
	return &fileAccessor{
		fileStoragePath: conf.FStoragePath,
	}
}

func (fa *fileAccessor) StateFullStorage(stflStorage StateFullStorage) *fileAccessor {
	fa.stflStorage = stflStorage
	return fa
}

func (fa *fileAccessor) FileStorage(fileStorage FileStorage) *fileAccessor {
	fa.fileStorage = fileStorage
	return fa
}

type fileAccessor struct {
	fileStoragePath string
	stflStorage     StateFullStorage
	fileStorage     FileStorage
}

func (fa *fileAccessor) GetFileInfoList(ctx context.Context) ([]domain.FileInfo, error) {
	log := domain.GetCtxLogger(ctx)
	action := domain.GetAction(1)

	log.Debugw(action, "msg", "start")

	bucket, err := fa.stflStorage.GetUserFilesBucket(ctx)
	if err != nil {
		log.Infow(action, "err", err.Error())
		return nil, err
	}

	lst, err := fa.fileStorage.GetFileInfoList(ctx, bucket)
	if err != nil {
		log.Infow(action, "list err", err.Error())
		return nil, fmt.Errorf("%w list err %s", domain.ErrServerInternal, err.Error())
	}

	log.Debugw(action, "msg", "complete")
	return lst, nil
}

func (fa *fileAccessor) DeleteFileInfo(ctx context.Context, name string) error {
	log := domain.GetCtxLogger(ctx)
	action := domain.GetAction(1)

	log.Debugw(action, "msg", "start")

	bucket, err := fa.stflStorage.GetUserFilesBucket(ctx)
	if err != nil {
		log.Infow(action, "err", err.Error())
		return err
	}

	if err = fa.fileStorage.DeleteFileInfo(ctx, bucket, name); err != nil {
		log.Infow(action, "delete err", err.Error())
		return err
	}

	log.Debugw(action, "msg", "complete")
	return nil
}

func (fa *fileAccessor) CreateStreamFileWriter(ctx context.Context) (domain.StreamFileWriter, error) {
	log := domain.GetCtxLogger(ctx)
	action := domain.GetAction(1)

	log.Debugw(action, "msg", "start")

	bucket, err := fa.stflStorage.GetUserFilesBucket(ctx)
	if err != nil {
		log.Infow(action, "err", err.Error())
		return nil, err
	}

	wrt, err := fa.fileStorage.CreateStreamFileWriter(ctx, bucket)
	if err != nil {
		log.Infow(action, "crete err", err.Error())
		return nil, err
	}

	log.Debugw(action, "msg", "complete")
	return wrt, nil
}

func (fa *fileAccessor) CreateStreamReader(ctx context.Context, info *domain.FileInfo) (domain.StreamFileReader, error) {
	log := domain.GetCtxLogger(ctx)
	action := domain.GetAction(1)

	log.Debugw(action, "msg", "start")

	bucket, err := fa.stflStorage.GetUserFilesBucket(ctx)
	if err != nil {
		log.Infow(action, "err", err.Error())
		return nil, err
	}

	rdr, err := fa.fileStorage.CreateStreamFileReader(ctx, bucket, info.Name)
	if err != nil {
		log.Infow(action, "crete err", err.Error())
		return nil, err
	}

	log.Debugw(action, "msg", "complete")
	return rdr, nil
}

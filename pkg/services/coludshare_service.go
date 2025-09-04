package services

import (
	"github.com/banyar/go-packages/pkg/frontlog"
	"github.com/banyar/go-packages/pkg/repositories"
	"go.uber.org/zap"
)

type CloudShareService struct {
	CloudShareRepo *repositories.CloudShareRepository
}

func NewCloudShareService(cloudeshareRepo *repositories.CloudShareRepository) *CloudShareService {
	return &CloudShareService{
		CloudShareRepo: cloudeshareRepo,
	}
}

func (c *CloudShareService) Get() ([]string, error) {
	resp, err := c.CloudShareRepo.FindCSVFiles()
	if err != nil {
		frontlog.Logger.Error("Error CloudShare connection:", zap.Any("error=", err))
		return nil, err
	}
	return resp, nil
}

func (c *CloudShareService) Download() ([]string, error) {
	resp, err := c.CloudShareRepo.DownloadCloudShareFile()
	if err != nil {
		frontlog.Logger.Error("Error CloudShare connection :", zap.Any("error=", err))
		return nil, err
	}
	return resp, nil
}

func (c *CloudShareService) Upload(fileNameList []string, cloudShareUploadFolder string) ([]string, []string, error) {
	successList, failedList, err := c.CloudShareRepo.UploadToCloudShare(fileNameList, cloudShareUploadFolder)
	if err != nil {
		frontlog.Logger.Error("Error CloudShare connection :", zap.Any("error=", err))
		return successList, failedList, err
	}
	return successList, failedList, nil
}

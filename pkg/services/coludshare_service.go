package services

import (
	"log"

	"github.com/banyar/go-packages/pkg/repositories"
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
		log.Fatal("ERROR : ", err)
		return nil, err
	}
	return resp, nil
}

func (c *CloudShareService) Download() ([]string, error) {
	resp, err := c.CloudShareRepo.DownloadCloudShareFile()
	if err != nil {
		log.Fatal("ERROR : ", err)
		return nil, err
	}
	return resp, nil
}

func (c *CloudShareService) Upload(fileNameList []string, cloudShareUploadFolder string) ([]string, []string, error) {
	successList, failedList, err := c.CloudShareRepo.UploadToCloudShare(fileNameList, cloudShareUploadFolder)
	if err != nil {
		log.Fatal("ERROR : ", err)
		return successList, failedList, err
	}
	return successList, failedList, nil
}

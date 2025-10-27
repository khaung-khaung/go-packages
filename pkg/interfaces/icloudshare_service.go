package interfaces

type ICloudShareService interface {
	Get() ([]string, error)
	Download() ([]string, error)
	Upload(fileNameList []string, cloudShareUploadFolder string) ([]string, []string, error)
}

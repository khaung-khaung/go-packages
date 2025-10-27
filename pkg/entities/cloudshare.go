package entities

import "time"

type DSNCloudShare struct {
	CloudShareBaseUrl          string //cloudshare base url
	CloudShareUserName         string //cloudshare login username
	CloudSharePassword         string //cloudshare login password
	CloudShareDownLoadFilePath string //cloudshare download file path
	FileAge                    time.Duration
	OutputFilePath             string //file data store path
	InputFilePath              string
	CloudShareUploadFilePath   string //cloudshare upload file path
	FileAgeOn                  bool
}

type FileDownloadResponse struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

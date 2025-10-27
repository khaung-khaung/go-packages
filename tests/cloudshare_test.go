package tests

import (
	"fmt"

	"os"
	"strconv"
	"testing"
	"time"

	"github.com/banyar/go-packages/pkg/adapters"
	"github.com/banyar/go-packages/pkg/common"
	"github.com/banyar/go-packages/pkg/entities"
	"github.com/banyar/go-packages/pkg/frontlog"
	"go.uber.org/zap"

	"github.com/gocarina/gocsv"
	"github.com/joho/godotenv"
)

type Test struct {
	Id     string // Ticket ID
	Reason string // Reason skipped
}

var (
	TestList          []Test
	UploadCSVFileList []string
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		frontlog.Logger.Error("Error loading .env file :", zap.Any("error=", err))
	}
	frontlog.Logger.Info("config.init()")
}

func TestCSVDownload(t *testing.T) {
	DSNCloudShare := GetDSNCloudShare()
	common.DisplayJsonFormat("DS", &DSNCloudShare)
	cloudShareAdapter := adapters.NewCloudShareAdapter(&DSNCloudShare)
	resp, err := cloudShareAdapter.CloudShareService.Download()

	if err != nil {
		frontlog.Logger.Error("Error :", zap.Any("", err))
	}

	csv, err := cloudShareAdapter.CloudShareService.Get()
	if err != nil {
		frontlog.Logger.Error("Error :", zap.Any("", err))
	}
	common.DisplayJsonFormat("cloudShareAdapter ", resp)
	common.DisplayJsonFormat("csv files ", csv)
}

func TestGetCSVFile(t *testing.T) {
	DSNCloudShare := GetDSNCloudShare()
	common.DisplayJsonFormat("DS", &DSNCloudShare)
	cloudShareAdapter := adapters.NewCloudShareAdapter(&DSNCloudShare)
	csv, err := cloudShareAdapter.CloudShareService.Get()
	if err != nil {
		frontlog.Logger.Error("Error :", zap.Any("", err))
	}
	common.DisplayJsonFormat("csv files ", csv)
}

func TestCSVUpload(t *testing.T) {
	recordID := "1"
	reason := "Test Reason"
	TestList = append(TestList, Test{Id: recordID, Reason: reason}) // Add record to skipped list
	fmt.Println("Full File Name", TestList)

	if len(TestList) > 0 {
		if fileName, err := DumpCSVfile("list.csv", TestList); err != nil {
			fmt.Println("csv_write", err.Error())
		} else {
			UploadCSVFileList = append(UploadCSVFileList, fileName)
		}
	}

	if len(UploadCSVFileList) > 0 {
		UploadCsvToCloudShare(UploadCSVFileList)
	}
}

func GetDSNCloudShare() entities.DSNCloudShare {
	fileAge, err := time.ParseDuration(os.Getenv("FILE_AGE"))
	if err != nil {

		frontlog.Logger.Error("converting FILE_AGE to time duration :", zap.Any("", err))
	}
	fileAgeOn, err := strconv.ParseBool(os.Getenv("FILE_AGE_ON"))

	if err != nil {

		frontlog.Logger.Error("Error converting FILE_AGE to bool :", zap.Any("", err))
	}

	return entities.DSNCloudShare{
		CloudShareBaseUrl:          os.Getenv("CLOUDSHARE_BASE_URL"),
		CloudShareUserName:         os.Getenv("CLOUDSHARE_USERNAME"),
		CloudSharePassword:         os.Getenv("CLOUDSHARE_PASSWORD"),
		CloudShareDownLoadFilePath: os.Getenv("CLOUDSHARE_DOWNLOAD_FILEPATH"),
		OutputFilePath:             os.Getenv("OUTPUT_FILEPATH"),
		FileAge:                    fileAge,
		FileAgeOn:                  fileAgeOn,
		CloudShareUploadFilePath:   os.Getenv("CLOUDSHARE_UPLOAD_FILEPATH"),
	}

}

func DumpCSVfile[T any](fileName string, itemList []T) (string, error) {
	DSNCloudShare := GetDSNCloudShare()
	startTime := time.Now()
	DTFormatFilename := "2006-01-02"
	appStartTime := startTime.Format(DTFormatFilename)
	fullFileName := fmt.Sprintf("%s-%s", appStartTime, fileName)

	frontlog.Logger.Info("Full File Name :", zap.Any("", fullFileName))

	outputPath := fmt.Sprintf("%s/dump", DSNCloudShare.OutputFilePath)
	if err := common.EnsureOutputFolderExists(outputPath); err != nil {
		frontlog.Logger.Error("csv_write error :", zap.Any("", err.Error()))
		return "", err
	}
	csvFileFullPath := fmt.Sprintf("%s/%s", outputPath, fullFileName)
	csvContent, err := gocsv.MarshalString(&itemList)
	if err != nil {

		frontlog.Logger.Error("csv_write error :", zap.Any("", err.Error()))
		return "", err
	}
	// Write the CSV content to a file
	file, err := os.Create(csvFileFullPath)
	if err != nil {
		return "", err
	}

	defer file.Close()
	_, err = file.WriteString(csvContent)
	if err != nil {
		frontlog.Logger.Error("csv_write error :", zap.Any("", err.Error()))
		return "", err
	}
	return csvFileFullPath, nil
}

func UploadCsvToCloudShare(fileNameList []string) {
	AppStartTime := time.Now()
	DTFormatDirName := "2006-01-02"
	cloudShareFolderName := AppStartTime.Format(DTFormatDirName)
	DSNCloudShare := GetDSNCloudShare()
	cloudShareUploadFullPath := fmt.Sprintf("%s/%s", DSNCloudShare.CloudShareUploadFilePath, cloudShareFolderName)
	fmt.Println("cloudShareUploadFullPath", cloudShareUploadFullPath)

	cloudShareAdapter := adapters.NewCloudShareAdapter(&DSNCloudShare)

	//upload to cloudshare
	successList, failList, err := cloudShareAdapter.CloudShareService.Upload(fileNameList, cloudShareUploadFullPath)

	if err != nil {
		frontlog.Logger.Error("csv upload error :", zap.Any("", err.Error()))
	}

	frontlog.Logger.Info("Success List :", zap.Any("", successList))
	frontlog.Logger.Info("Failed List :", zap.Any("", failList))

}

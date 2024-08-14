package tests

import (
	"fmt"
	"log"

	"os"
	"strconv"
	"testing"
	"time"

	"github.com/banyar/go-packages/pkg/adapters"
	"github.com/banyar/go-packages/pkg/common"
	"github.com/banyar/go-packages/pkg/entities"

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
		log.Fatal("Error loading .env file")
	}
	fmt.Println("config.init()")
}

func TestCSVDownload(t *testing.T) {
	DSNCloudShare := GetDSNCloudShare()
	common.DisplayJsonFormat("DS", &DSNCloudShare)
	cloudShareAdapter := adapters.NewCloudShareAdapter(&DSNCloudShare)
	resp, err := cloudShareAdapter.CloudShareService.Download()

	if err != nil {
		log.Fatal("ERROR : ", err)
	}

	csv, err := cloudShareAdapter.CloudShareService.Get()
	if err != nil {
		log.Fatal("ERROR : ", err)
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
		log.Fatal("ERROR : ", err)
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
		log.Fatalf("Error converting FILE_AGE to time duration: %v", err)
	}
	fileAgeOn, err := strconv.ParseBool(os.Getenv("FILE_AGE_ON"))

	if err != nil {
		log.Fatalf("Error converting FILE_AGE to bool: %v", err)
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
	fmt.Println("Full File Name", fullFileName)

	outputPath := fmt.Sprintf("%s/dump", DSNCloudShare.OutputFilePath)
	if err := common.EnsureOutputFolderExists(outputPath); err != nil {
		fmt.Println("csv_write error", err.Error())
		return "", err
	}
	csvFileFullPath := fmt.Sprintf("%s/%s", outputPath, fullFileName)
	csvContent, err := gocsv.MarshalString(&itemList)
	if err != nil {
		fmt.Println("csv_write error", err.Error())
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
		fmt.Println("csv_write error", err.Error())
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
		fmt.Println("csv upload error", err.Error())
	}
	fmt.Println("Success List", successList)
	fmt.Println("Failed List", failList)
}

package repositories

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	entities "github.com/banyar/go-packages/pkg/entities"

	"github.com/studio-b12/gowebdav"
)

type CloudShareRepository struct {
	client *gowebdav.Client
	dsn    *entities.DSNCloudShare
}

func ConnectCloudShare(DSNCloudShare *entities.DSNCloudShare) *CloudShareRepository {
	return &CloudShareRepository{
		client: gowebdav.NewClient(DSNCloudShare.CloudShareBaseUrl, DSNCloudShare.CloudShareUserName, DSNCloudShare.CloudSharePassword),
		dsn:    DSNCloudShare,
	}
}

func (r *CloudShareRepository) FindCSVFiles() ([]string, error) {
	dir := r.dsn.OutputFilePath
	// fmt.Println("csv dir ===> ", dir)
	fileAge := r.dsn.FileAge
	// fmt.Println("fileAge ====> ", fileAge)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var csvFiles []string
	now := time.Now()
	lastModifiedTime := now.Add(-fileAge)
	// fmt.Println("fileAge ====> ", fileAge)
	for _, entry := range entries {

		// check if the entry is a regular file and has a .csv extension
		if entry.Type().IsRegular() && filepath.Ext(entry.Name()) == ".csv" {
			// check if the entry was modified within the given time
			info, err := entry.Info()
			if err != nil {
				return nil, err
			}

			if r.dsn.FileAgeOn {
				modTime := info.ModTime()
				if modTime.After(lastModifiedTime) {
					csvFiles = append(csvFiles, entry.Name())
				} else {
					duration := lastModifiedTime.Sub(modTime)
					// Print the duration
					// fmt.Printf("The difference between %v and %v is %v\n", lastModifiedTime, modTime, duration)
					// Get the difference in hours, minutes, and seconds
					hours := int(duration.Hours())
					minutes := int(duration.Minutes()) % 60
					seconds := int(duration.Seconds()) % 60
					fmt.Printf("Difference: %d hours, %d minutes, %d seconds = %s \n", hours, minutes, seconds, entry.Name())
				}
			} else {
				csvFiles = append(csvFiles, entry.Name())
			}
		}
	}
	return csvFiles, nil
}

func (r *CloudShareRepository) DownloadCloudShareFile() ([]string, error) {
	csInputFilePath := r.dsn.CloudShareDownLoadFilePath
	fmt.Println("csInputFilePath ===> ", csInputFilePath)
	fileAge := r.dsn.FileAge
	fmt.Println("fileAge ====> ", fileAge)

	// Get the files in the root path
	files, err := r.client.ReadDir(csInputFilePath)
	if err != nil {
		fmt.Println("Error reading CloudShare folder", err.Error())
	}
	csvFiles := []string{"No csv file found"}
	if len(files) != 0 {
		now := time.Now()
		lastModifiedTime := now.Add(-fileAge)
		// loop over the files and filter by extension and modification time
		for _, file := range files {
			// check if the file is a regular file and has a .csv extension
			if file.Mode().IsRegular() && filepath.Ext(file.Name()) == ".csv" {
				// check if the file was modified within the last 15 minutes
				if file.ModTime().After(lastModifiedTime) {
					// CloudShare file name with path
					csFileName := fmt.Sprintf("%s/%s", csInputFilePath, file.Name())
					// read the file content as a byte slice
					data, err := r.client.Read(csFileName)
					if err != nil {
						fmt.Println("Error reading CloudShare file", err.Error())
						continue
					}
					// write the file content to a local file with the same name
					fileName := fmt.Sprintf("%s/%s", r.dsn.OutputFilePath, file.Name())
					err = os.WriteFile(fileName, data, 0644)
					if err != nil {
						fmt.Println("Error writing file", err.Error())
						continue
					}
					csvFiles = nil
					csvFiles = append(csvFiles, fmt.Sprintf("%s %s", "File downloaded", file.Name()))
				}
			}
		}
	}
	return csvFiles, nil
}

// UploadFile check and if file not existing create folder to WebDAV (CloudShare) connection
func (r *CloudShareRepository) UploadToCloudShare(fileNameList []string, cloudShareUploadFolder string) ([]string, []string, error) {
	var failList []string
	var successList []string
	folderExist, err := r.checkAndFolderCreate(cloudShareUploadFolder)
	if err != nil {
		fmt.Println("Error creating CloudShare folder", err.Error())
		return successList, failList, err
	}

	if folderExist {
		// upload to cloudshare
		for _, fileName := range fileNameList {
			err := r.uploadFile(fileName, cloudShareUploadFolder)
			if err != nil {
				failList = append(failList, fileName)
				fmt.Println("Error uploading to cloud share", fileName)
			} else {
				successList = append(successList, fileName)
				fmt.Println("Success uploading to cloud share", fileName)
			}
		}
	}
	return successList, failList, nil
}

// UploadFile check and if file not existing create folder to WebDAV (CloudShare) connection
func (r *CloudShareRepository) checkAndFolderCreate(cloudShareUploadFolder string) (bool, error) {
	folderExist := true
	_, isFound := r.client.Stat(cloudShareUploadFolder)
	if isFound != nil {
		fmt.Println("CloudShare folder does not exist. ", isFound.Error())
		err := r.client.Mkdir(cloudShareUploadFolder, 0755) // create upload folder to cloudshare
		if err != nil {
			folderExist = false
			return folderExist, err
		}
		fmt.Println("CloudShare folder successfully created", cloudShareUploadFolder)
		folderExist = true
	}
	return folderExist, nil
}

// UploadFile Upload individual file to existing WebDAV (CloudShare) connection
func (r *CloudShareRepository) uploadFile(fileName string, remotePath string) error {
	// Open the local file for reading
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file", err.Error())
		return err
	}
	defer file.Close()
	// Read the content of the local file into a []byte slice
	fileStat, _ := file.Stat()
	fileSize := fileStat.Size()
	fileContent := make([]byte, fileSize)
	_, err = file.Read(fileContent)
	if err != nil {
		fmt.Println("Error reading local file", fileStat.Name())
		return err
	}

	parts := strings.Split(fileName, "/")
	remotePath = fmt.Sprintf("%s/%s", remotePath, parts[len(parts)-1])
	// Upload the file to the WebDAV server
	err = r.client.Write(remotePath, fileContent, 0644)
	if err != nil {
		fmt.Println("Error uploading file", err.Error())
	}
	return nil
}

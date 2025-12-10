package handlers

import (
	"sccsmsserver/db/pg"
	"sccsmsserver/i18n"
	"sccsmsserver/pkg/aws"
	"sccsmsserver/pkg/mysf"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Recieve client upload multiple files handler
func RecieveFilesHandler(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		zap.L().Error("RecieveFilesHandler invalid param:", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get Operator ID
	operatorID, resStatus := GetOperatorID(c)
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, nil)
		return
	}
	// Extract files and other info from formdata
	var fileNumber int
	files := form.File["files"]
	fileNumber = len(files)

	if fileNumber == 0 {
		zap.L().Error("RecieveFilesHandler no files")
		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
		return
	}
	// Get file key
	filekeyArr, ok := c.GetPostFormArray("fileKey")
	if !ok {
		zap.L().Error("RecieveFilesHandler  c.GetPostFormArray(filekey) not ok")
		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
		return
	}
	if len(filekeyArr) != fileNumber {
		zap.L().Error("RecieveFilesHandler filekey number not enough")
		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
		return
	}
	// Get file type
	fileTypeArr, ok := c.GetPostFormArray("fileType")
	if !ok {
		zap.L().Error("RecieveFilesHandler c.GetPostFormArray(fileType) not ok")
		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
		return
	}
	if len(fileTypeArr) != fileNumber {
		zap.L().Error("RecieveFilesHandler filetype number not enough")
		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
		return
	}
	// Get file name
	fileNameArr, ok := c.GetPostFormArray("fileName")
	if !ok {
		zap.L().Error("RecieveFilesHandler c.GetPostFormArray(filename) not ok")
		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
		return
	}
	if len(fileNameArr) != fileNumber {
		zap.L().Error("RecieveFilesHandler filetype filename not enough")
		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
		return
	}
	// Get isImage
	fileIsImageArr, ok := c.GetPostFormArray("isImage")
	if !ok {
		zap.L().Error("RecieveFilesHandler c.GetPostFormArray(isimage) not ok")
		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
		return
	}
	if len(fileIsImageArr) != fileNumber {
		zap.L().Error("RecieveFilesHandler isimage number not enough")
		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
		return
	}
	// Get file hash
	fileHashArr, ok := c.GetPostFormArray("hash")
	if !ok {
		zap.L().Error("RecieveFilesHandler c.GetPostFormArray(filehash) not ok")
		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
		return
	}
	if len(fileHashArr) != fileNumber {
		zap.L().Error("RecieveFilesHandler filehash number not enough")
		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
		return
	}
	// Get camera model
	modelArr, ok := c.GetPostFormArray("model")
	if !ok {
		zap.L().Error("RecieveFilesHandler c.GetPostFormArray(model) not ok")
		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
		return
	}
	if len(modelArr) != fileNumber {
		zap.L().Error("RecieveFilesHandler model number not enough")
		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
	}
	// Get latitude
	latitudeArr, ok := c.GetPostFormArray("latitude")
	if !ok {
		zap.L().Error("RecieveFilesHandler c.GetPostFormArray(latitude) not ok")
		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
		return
	}
	if len(latitudeArr) != fileNumber {
		zap.L().Error("RecieveFilesHandler latitude number not enough")
		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
	}
	// Get longitude
	longitudeArr, ok := c.GetPostFormArray("longitude")
	if !ok {
		zap.L().Error("RecieveFilesHandler c.GetPostFormArray(longitude) not ok")
		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
		return
	}
	if len(longitudeArr) != fileNumber {
		zap.L().Error("RecieveFilesHandler longitude number not enough")
		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
	}
	// Get original date time
	dateTimeOriginal, ok := c.GetPostFormArray("DateTimeOriginal")
	if !ok {
		zap.L().Error("RecieveFilesHandler c.GetPostFormArray(DateTimeOriginal) not ok")
		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
		return
	}
	if len(dateTimeOriginal) != fileNumber {
		zap.L().Error("RecieveFilesHandler dateTimeOriginal number not enough")
		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
	}

	// Get file source
	source, ok := c.GetPostFormArray("source")
	if !ok {
		zap.L().Error("RecieveFilesHandler c.GetPostFormArray(source) not ok")
		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
		return
	}
	if len(source) != fileNumber {
		zap.L().Error("RecieveFilesHandler source number not enough")
		ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
	}

	var fileRespArr []pg.File
	var build strings.Builder
	for index, file := range files {
		var fileInfo pg.File
		fileInfo.OriginFileName = fileNameArr[index]
		fileInfo.FileKey, _ = strconv.Atoi(filekeyArr[index])
		fileInfo.Hash = fileHashArr[index]
		fileInfo.FileType = fileTypeArr[index]
		fileInfo.IsImage, _ = strconv.Atoi(fileIsImageArr[index])
		fileInfo.Size = file.Size
		fileInfo.Model = modelArr[index]
		fileInfo.Longitude, _ = strconv.ParseFloat(longitudeArr[index], 64)
		fileInfo.Latitude, _ = strconv.ParseFloat(latitudeArr[index], 64)
		fileInfo.DateTimeOriginal = dateTimeOriginal[index]
		fileInfo.Source = source[index]

		fileInfo.CreatorID = operatorID
		fileInfo.UpLoadDate = time.Now()

		// Generate upload file name using Snowflake algorithm
		fileName := strconv.FormatInt(mysf.GenID(), 10)
		build.Reset()
		build.WriteString(fileName)
		build.WriteString(fileTypeArr[index])
		fileInfo.MinioFileName = build.String()

		fileObj, err := file.Open()
		if err != nil {
			zap.L().Error("RecieveFilesHandler File Open failed:", zap.Error(err))
			ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
			return
		}
		// Upload file to minio server
		_, err = aws.UploadFile(fileInfo.MinioFileName, fileObj, file.Size)
		if err != nil {
			zap.L().Error("RecieveFilesHandler File Upload failed:", zap.Error(err))
			ResponseWithMsg(c, i18n.CodeInvalidParm, nil)
			return
		}
		// Write file information to database
		resStatus, _ = fileInfo.Add()
		if resStatus != i18n.StatusOK {
			ResponseWithMsg(c, resStatus, err)
			return
		}
		fileRespArr = append(fileRespArr, fileInfo)
	}
	// Response
	ResponseWithMsg(c, i18n.StatusOK, fileRespArr)
}

// Get File information by file hash handler
func GetFileInfoByHashHandler(c *gin.Context) {
	f := new(pg.File)
	err := c.ShouldBind(f)
	if err != nil {
		zap.L().Error("GetFileInfoByHash invalid param:", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get file information
	resStatus, _ := f.GetFileInfoByHash()
	if resStatus != i18n.StatusOK {
		ResponseWithMsg(c, resStatus, f)
		return
	}
	// Response
	ResponseWithMsg(c, i18n.StatusOK, f)
}

// Get file array by file hash array handler
func GetFilesByHashHandler(c *gin.Context) {
	fs := new([]pg.File)
	err := c.ShouldBind(fs)
	if err != nil {
		zap.L().Error("GetFilesByHashHandler invalid param:", zap.Error(err))
		ResponseWithMsg(c, i18n.CodeInvalidParm, err)
		return
	}
	// Get files information
	fileArr, resStatus, _ := pg.GetFilesByHash(*fs)
	// Response
	ResponseWithMsg(c, resStatus, fileArr)
}

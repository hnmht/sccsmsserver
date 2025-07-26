package pg

import (
	"database/sql"
	"encoding/json"
	"sccsmsserver/cache"
	"sccsmsserver/i18n"
	"sccsmsserver/pkg/minio"
	"sccsmsserver/pub"
	"time"

	"go.uber.org/zap"
)

// File details
type File struct {
	ID               int32     `db:"id" json:"ID"`
	Hash             string    `db:"hash" json:"hash"`
	MinioFileName    string    `db:"miniofilename" json:"minioFileName"`
	OriginFileName   string    `db:"originfilename" json:"originFileName"`
	FileKey          int       `db:"filekey" json:"fileKey"`
	FilePath         string    `json:"filePath"`
	FileUri          string    `json:"fileUri"`
	Mime             string    `json:"mime"`
	FileType         string    `db:"filetype" json:"fileType"`
	IsImage          int       `db:"isimage" json:"isImage"`
	Model            string    `db:"model" json:"model"`
	Longitude        float64   `db:"longitude" json:"longitude"`
	Latitude         float64   `db:"latitude" json:"latitude"`
	Size             int64     `db:"size" json:"size"`
	FileUrl          string    `db:"fileurl" json:"fileUrl"`
	DateTimeOriginal string    `db:"datetimeoriginal" json:"dateTimeOriginal"`
	UpLoadDate       time.Time `db:"uploaddate" json:"uploadTime"`
	Source           string    `db:"source" json:"source"`
	CreatorID        int32     `db:" creatorid" json:"creatorid"`
	CreatorName      string    `json:"creatorName"`
	Dr               int16     `db:"dr" json:"dr"`
	Ts               time.Time `db:"ts" json:"ts"`
}

// Get File information by file ID.
func (file *File) GetFileInfoByID() (resStatus i18n.ResKey, err error) {
	// Get file information from cache
	number, fb, _ := cache.Get(pub.File, file.ID)
	if number > 0 {
		err = json.Unmarshal(fb, &file)
		if err != nil {
			resStatus = i18n.StatusInternalError
			zap.L().Error("GetPersonInfoByID json.Unmarshal failed", zap.Error(err))
			return
		}
		resStatus = i18n.StatusOK
		return
	}
	// If file information isn't in cache, retrieve it from databases
	sqlStr := `select miniofilename,originfilename,filekey,filetype,isimage,
	model,longitude,latitude,size,datetimeoriginal,
	uploaddate,creatorid,filehash,source,ts 
	from filelist where id=$1`
	err = db.QueryRow(sqlStr, file.ID).Scan(&file.MinioFileName, &file.OriginFileName, &file.FileKey, &file.FileType, &file.IsImage,
		&file.Model, &file.Longitude, &file.Latitude, &file.Size, &file.DateTimeOriginal,
		&file.UpLoadDate, &file.CreatorName, &file.Hash, &file.Source, &file.Ts)
	if err != nil {
		if err == sql.ErrNoRows {
			resStatus = i18n.StatusFileNotExist
			file.CreatorID = 0
			file.FileUrl = ""
			return
		}
		resStatus = i18n.StatusInternalError
		zap.L().Error("file getFileInfoByID db.queryrow failed")
		return
	}
	// Get File URL
	fileUrl, err := minio.GetFileUrl(file.MinioFileName, pub.FileURLExpireTime)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("file.getFileInfoByID minio.GetFileUrl failed", zap.Error(err))
		return
	}
	file.FileUrl = fileUrl
	// Write into cache
	jsonB, _ := json.Marshal(file)
	err = cache.Set(pub.File, file.ID, jsonB)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("file getFileInfoByID cache.Set failed", zap.Error(err))
		return
	}
	return i18n.StatusOK, nil
}

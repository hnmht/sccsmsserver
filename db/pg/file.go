package pg

import (
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"sccsmsserver/cache"
	"sccsmsserver/i18n"
	"sccsmsserver/pkg/minio"
	"sccsmsserver/pub"
	"time"

	"go.uber.org/zap"
)

// File details
type File struct {
	ID               int32     `db:"id" json:"id"`
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
	CreatorID        int32     `db:" creatorid" json:"creatorID"`
	CreatorName      string    `json:"creatorName"`
	Dr               int16     `db:"dr" json:"dr"`
	Ts               time.Time `db:"ts" json:"ts"`
}

// Voucher File details
type VoucherFile struct {
	ID         int32     `db:"id" json:"id"`
	BillBID    int32     `db:"billbid" json:"billBID"`
	BillHID    int32     `db:"billhid" json:"billHID"`
	File       File      `db:"fileid" json:"file"`
	CreateDate time.Time `db:"createtime" json:"createDate"`
	Creator    Person    `db:"creatorid" json:"creator"`
	ModifyDate time.Time `db:"modifytime" json:"modifyDate"`
	Modifier   Person    `db:"modifierid" json:"modifier"`
	Ts         time.Time `db:"ts" json:"ts"`
	Dr         int16     `db:"dr" json:"dr"`
}

// Get File information by file ID.
func (file *File) GetFileInfoByID() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
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
	// Write File into cache
	jsonB, _ := json.Marshal(file)
	err = cache.Set(pub.File, file.ID, jsonB)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("file getFileInfoByID cache.Set failed", zap.Error(err))
		return
	}
	// Write FileHash into cache
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(file.ID))
	hashString := fmt.Sprintf("%s%s%s", pub.FileHash, ":", file.Hash)
	err = cache.SetOther(hashString, buf)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("file getFileInfoByID cache.SetOther failed", zap.Error(err))
		return
	}

	return
}

// Add File record to database
func (file *File) Add() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Add Record to files table
	sqlStr := `insert into filelist(miniofilename,originfilename,filekey,filetype,isimage,
	model,longitude,latitude,size,datetimeoriginal,
	uploaddate,creatorid,filehash,source)
	values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) returning id`
	err = db.QueryRow(sqlStr, file.MinioFileName, file.OriginFileName, file.FileKey, file.FileType, file.IsImage,
		file.Model, file.Longitude, file.Latitude, file.Size, file.DateTimeOriginal,
		file.UpLoadDate, file.CreatorID, file.Hash, file.Source).Scan(&file.ID)
	if err != nil {
		zap.L().Error("file.Add stmt.QueryRow failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	// Request Minio server to get file URL
	fileUrl, err := minio.GetFileUrl(file.MinioFileName, pub.FileURLExpireTime)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("file.Add minio.GetFileUrl failed", zap.Error(err))
		return
	}
	file.FileUrl = fileUrl
	return
}

// Get File information by file hash.
func (file *File) GetFileInfoByHash() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	// Get file ID from cache
	hashString := fmt.Sprintf("%s%s%s", pub.FileHash, ":", file.Hash)
	number, fb, _ := cache.GetOther(hashString)
	if number > 0 {
		// If file ID is in cache, get file information by file ID
		file.ID = int32(binary.BigEndian.Uint32(fb))
		resStatus, err = file.GetFileInfoByID()
		return
	}
	// If file ID isn't in cache, retrieve it from databases
	// Get file information from filelist table
	sqlStr := `select id,miniofilename,originfilename,filetype,isimage,
	model,longitude,latitude,size,datetimeoriginal,
	source,uploaddate,creatorid,ts
	from filelist where filehash=$1 limit 1`
	err = db.QueryRow(sqlStr, file.Hash).Scan(&file.ID, &file.MinioFileName, &file.OriginFileName, &file.FileType, &file.IsImage,
		&file.Model, &file.Longitude, &file.Latitude, &file.Size, &file.DateTimeOriginal,
		&file.Source, &file.UpLoadDate, &file.CreatorID, &file.Ts)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
			file.ID = 0
			file.FileUrl = ""
			return
		}
		resStatus = i18n.StatusInternalError
		zap.L().Error("file GetFileInfoByHash db.queryrow failed")
		return
	}
	// Request Minio server to get file URL
	fileUrl, err := minio.GetFileUrl(file.MinioFileName, pub.FileURLExpireTime)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("file.GetFileInfoByHash minio.GetFileUrl failed", zap.Error(err))
		return
	}
	file.FileUrl = fileUrl

	// Write File into cache
	jsonB, _ := json.Marshal(file)
	err = cache.Set(pub.File, file.ID, jsonB)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("file GetFileInfoByHash cache.Set failed", zap.Error(err))
		return
	}
	// Write FileHash into cache
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(file.ID))
	err = cache.SetOther(hashString, buf)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("file GetFileInfoByHash cache.SetOther failed", zap.Error(err))
		return
	}

	return
}

// Get file array by file hash array
func GetFilesByHash(files []File) (fileArr []File, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	fileArr = make([]File, 0)
	for _, file := range files {
		resStatus, err = file.GetFileInfoByHash()
		if err != nil || resStatus != i18n.StatusOK {
			return
		}
		fileArr = append(fileArr, file)
	}
	return
}

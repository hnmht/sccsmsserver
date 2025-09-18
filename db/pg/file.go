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

// Voucher File details
type VoucherFile struct {
	ID         int32     `db:"id" json:"id"`                   //行id
	BillBid    int32     `db:"billbid" json:"billbid"`         //单据表体id
	BIllHid    int32     `db:"billhid" json:"billhid"`         //单据表头id
	File       File      `db:"file_id" json:"file"`            //文件id
	CreateDate time.Time `db:"create_time" json:"createdate"`  //创建日期
	CreateUser Person    `db:"createuserid" json:"createuser"` //创建人
	ModifyDate time.Time `db:"modify_time" json:"modifydate"`  //更新日期
	ModifyUser Person    `db:"modifyuserid" json:"modifyuser"` //更新人
	Ts         time.Time `db:"ts" json:"ts"`                   //时间戳
	Dr         int16     `db:"dr" json:"dr"`                   //删除标志
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

/*
// Add 将文件信息写入filelist表
func (file *File) Add() (err error) {
	//向数据库filelist表中写入记录预处理
	sqlStr := `insert into filelist(miniofilename,originfilename,filekey,filetype,isimage,
	model,longitude,latitude,size,datetimeoriginal,
	uploaddate,createuserid,filehash,source)
	values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14) returning id`
	stmt, err := db.Prepare(sqlStr)
	if err != nil {
		zap.L().Error("File.Add db.Perpare failed", zap.Error(err))
		return
	}
	defer stmt.Close()

	//向数据库中写入记录
	err = stmt.QueryRow(file.MinioFileName, file.OriginFileName, file.FileKey, file.FileType, file.IsImage,
		file.Model, file.Longitude, file.Latitude, file.Size, file.DateTimeOriginal,
		file.UpLoadDate, file.CreateUserID, file.FileHash, file.Source).Scan(&file.FileId)
	if err != nil {
		zap.L().Error("file.Add stmt.QueryRow failed", zap.Error(err))
		return
	}
	//获取文件url
	fileUrl, err := minio.GetFileUrl(file.MinioFileName, durtion)
	file.FileUrl = fileUrl
	return
}

// GetFileInfoByHash 根据文件hash获取文件信息
func (file *File) GetFileInfoByHash() (err error) {
	//获取文件信息
	sqlStr := `select id,miniofilename,originfilename,filetype,isimage,
	model,longitude,latitude,size,datetimeoriginal,
	source,uploaddate,createuserid,ts
	from filelist where filehash=$1 limit 1`
	err = db.QueryRow(sqlStr, file.FileHash).Scan(&file.FileId, &file.MinioFileName, &file.OriginFileName, &file.FileType, &file.IsImage,
		&file.Model, &file.Longitude, &file.Latitude, &file.Size, &file.DateTimeOriginal,
		&file.Source, &file.UpLoadDate, &file.CreateUserID, &file.Ts)

	if err != nil {
		if err == sql.ErrNoRows {
			file.FileId = 0
			file.FileUrl = ""
			return nil
		}
		zap.L().Error("file GetFileInfoByHash db.queryrow failed")
		return err
	}
	//获取文件url
	fileUrl, err := minio.GetFileUrl(file.MinioFileName, durtion)
	file.FileUrl = fileUrl
	return
}

// GetFilesByHash 根据hash批量获取文件信息
func GetFilesByHash(files []File) (fileArr []File, resStatus pub.ResStatus, err error) {
	for _, file := range files {
		file.GetFileInfoByHash()
		fileArr = append(fileArr, file)
	}
	resStatus = pub.StatusOK
	return
}
*/

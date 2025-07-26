package pg

import "time"

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
	CreateUserName   string    `json:"creatorName"`
	Dr               int16     `db:"dr" json:"dr"`
	Ts               time.Time `db:"ts" json:"ts"`
}

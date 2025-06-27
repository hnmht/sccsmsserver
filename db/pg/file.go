package pg

import "time"

//File 文件详情
type File struct {
	FileId           int32     `db:"id" json:"fileID"`
	FileHash         string    `db:"filehash" json:"filehash"`
	MinioFileName    string    `db:"miniofilename" json:"miniofilename"`
	OriginFileName   string    `db:"originfilename" json:"originfilename"`
	FileKey          int       `db:"filekey" json:"filekey"`
	FilePath         string    `json:"filepath"`
	FileUri          string    `json:"fileuri"`
	Mime             string    `json:"mime"`
	FileType         string    `db:"filetype" json:"filetype"`
	IsImage          int       `db:"isimage" json:"isimage"`
	Model            string    `db:"model" json:"model"`
	Longitude        float64   `db:"longitude" json:"longitude"` //经度
	Latitude         float64   `db:"latitude" json:"latitude"`   //纬度
	Size             int64     `db:"size" json:"size"`
	FileUrl          string    `db:"fileurl" json:"fileurl"`
	DateTimeOriginal string    `db:"datetimeoriginal" json:"datetimeoriginal"`
	UpLoadDate       time.Time `db:"uploaddate" json:"uploadtime"`
	Source           string    `db:"source" json:"source"`
	CreateUserID     int32     `db:" createuserid" json:"createuserid"`    //创建用户id
	CreateUserName   string    `db:"createusername" json:"createusername"` //创建用户名称
	Ts               time.Time `db:"ts" json:"ts"`
}

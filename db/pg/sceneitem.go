package pg

/* //SceneItemOption 现场档案自定义项模型
type SceneItemOption struct {
	ID           int32           `db:"id" json:"id"`
	Code         string          `db:"code" json:"code"`
	Name         string          `db:"name" json:"name"`
	DisplayName  string          `db:"displayname" json:"displayname"`
	UDC          UserDefineClass `db:"udc_id" json:"udc"`
	DefaultValue UserDefineDoc   `db:"defaultvalue_id" json:"defaultvalue"`
	Enable       int16           `db:"enable" json:"enable"`
	IsModify     int16           `json:"ismodify"`
	CreateDate   time.Time       `db:"create_time" json:"createdate"`
	CreateUser   Person          `db:"createuserid" json:"createuser"`
	ModifyDate   time.Time       `db:"modify_time" json:"modifydate"`
	ModifyUser   Person          `db:"modifyuserid" json:"modifyuser"`
	Ts           time.Time       `db:"ts" json:"ts"`
	Dr           int16           `db:"dr" json:"dr"` //删除标志
} */

// 现场档案选项初始化表
func initSceneItemOption() (isFinish bool, err error) {
	//检查是否存在记录
	sqlStr := "select count(id) from sceneitemoption"
	hasRecord, isFinish, err := checkRecord("sceneitemoption", sqlStr)
	if hasRecord || !isFinish || err != nil { //有数据 或 没有完成 或有错误
		return
	}
	//没有数据继续执行初始化
	/* 	var options = []SceneItemOption{
	   		{ID: 1, Code: "udf1", Name: "自定义项1", DisplayName: "自定义项1"},
	   		{ID: 2, Code: "udf2", Name: "自定义项2", DisplayName: "自定义项2"},
	   		{ID: 3, Code: "udf3", Name: "自定义项3", DisplayName: "自定义项3"},
	   		{ID: 4, Code: "udf4", Name: "自定义项4", DisplayName: "自定义项4"},
	   		{ID: 5, Code: "udf5", Name: "自定义项5", DisplayName: "自定义项5"},
	   		{ID: 6, Code: "udf6", Name: "自定义项6", DisplayName: "自定义项6"},
	   		{ID: 7, Code: "udf7", Name: "自定义项7", DisplayName: "自定义项7"},
	   		{ID: 8, Code: "udf8", Name: "自定义项8", DisplayName: "自定义项8"},
	   		{ID: 9, Code: "udf9", Name: "自定义项9", DisplayName: "自定义项9"},
	   		{ID: 10, Code: "udf10", Name: "自定义项10", DisplayName: "自定义项10"},
	   	}

	   	sqlStr = "insert into sceneitemoption(id,code,name,displayname) values($1,$2,$3,$4)"
	   	for _, option := range options {
	   		_, err = db.Exec(sqlStr, option.ID, option.Code, option.Name, option.DisplayName)
	   		if err != nil {
	   			isFinish = false
	   			zap.L().Error("initSceneItemOption insert initvalues failed", zap.Error(err))
	   			return
	   		}
	   	} */
	return
}

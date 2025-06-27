package pg

import "time"

// 人员档案（简化版用户档案）
type Person struct {
	ID          int32     `db:"id" json:"id"`
	Code        string    `db:"usercode" json:"code"`
	Name        string    `db:"username" json:"name"`
	Avatar      File      `db:"file_id" json:"avatar"`
	DeptID      int32     `db:"dept_id" json:"deptid"`
	DeptCode    string    `json:"deptcode"`
	DeptName    string    `json:"deptname"`
	IsOperator  int16     `json:"isoperator"`
	OpID        int32     `json:"op_id"`
	OpName      string    `json:"opname"` //岗位名称
	Description string    `db:"description" json:"description"`
	Mobile      string    `db:"mobile" json:"mobile"`
	Email       string    `db:"email" json:"email"`
	Gender      int16     `db:"gender" json:"gender"`
	SystemFlag  int16     `db:"systemflag" json:"systemflag"`
	Status      int16     `db:"status" json:"status"`
	CreateDate  time.Time `db:"create_time" json:"createdate"`
	Ts          time.Time `db:"ts" json:"ts"`
	Dr          int16     `db:"dr" json:"dr"`
}

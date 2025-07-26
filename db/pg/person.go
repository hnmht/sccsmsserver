package pg

import "time"

// Person Master Data (simplify User)
type Person struct {
	ID           int32     `db:"id" json:"id"`
	Code         string    `db:"code" json:"code"`
	Name         string    `db:"name" json:"name"`
	Avatar       File      `db:"fileid" json:"avatar"`
	DeptID       int32     `db:"deptid" json:"deptID"`
	DeptCode     string    `json:"deptCode"`
	DeptName     string    `json:"deptName"`
	IsOperator   int16     `json:"isOperator"`
	PositionID   int32     `db:"positionid" json:"positionID"`
	PositionName string    `json:"positionName"`
	Description  string    `db:"description" json:"description"`
	Mobile       string    `db:"mobile" json:"mobile"`
	Email        string    `db:"email" json:"email"`
	Gender       int16     `db:"gender" json:"gender"`
	SystemFlag   int16     `db:"systemflag" json:"systemflag"`
	Status       int16     `db:"status" json:"status"`
	CreateDate   time.Time `db:"createtime" json:"createDate"`
	Ts           time.Time `db:"ts" json:"ts"`
	Dr           int16     `db:"dr" json:"dr"`
}

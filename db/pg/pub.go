package pg

import (
	"sccsmsserver/i18n"
)

// SeaCloud Date Type struct
type ScDataType struct {
	ID        int32  `json:"id"`
	TypeCode  string `json:"code"`
	TypeName  string `json:"name"`
	DataType  string `json:"dataType"`
	FrontDb   string `json:"frontDb"`
	InputMode string `json:"inputMode"`
}

// Data references check struct
type DataReferenceCheck struct {
	Description    string
	SqlStr         string
	UsedReturnCode i18n.ResKey
}

// Data query conditions
type QueryParams struct {
	QueryString string `json:"queryString"`
}

// Data query pagination
type PagingQueryParams struct {
	QueryString string `json:"queryString"`
	Page        int32  `json:"page"`
	PerPage     int32  `json:"perPage"`
}

// the struct Check the archive is refreneced
type ArchiveCheckUsed struct {
	Description    string
	SqlStr         string
	UsedReturnCode i18n.ResKey
}

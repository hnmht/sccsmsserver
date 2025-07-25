package pg

import (
	"sccsmsserver/i18n"
)

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

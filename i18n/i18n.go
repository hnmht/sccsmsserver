package i18n

import (
	"embed"
)

//go:embed translations/*
var transFile embed.FS

// System messages to users
type ResKey string

// plural.Selectf function parameters
type SelectfParms struct {
	Arg    int           `json:"arg"`
	Format string        `json:"format"`
	Cases  []interface{} `json:"cases"`
}

type TransItem struct {
	Key          string       `json:"key"`
	Type         string       `json:"type"` // string || plural
	Message      string       `json:"message"`
	SelectfParms SelectfParms `json:"selectParms"`
}

// Language translation file structor
type LangTransFile struct {
	Language   string      `json:"language"`
	TransItems []TransItem `json:"transItems"`
}

const (
	// System Message
	CodeSuccess       ResKey = "CodeSuccess"
	CodeInvalidParm   ResKey = "CodeInvalidParm"
	CodeServerBusy    ResKey = "CodeServerBusy"
	CodeInternalError ResKey = "CodeInternalError"
	CodeInvalidToken  ResKey = "CodeInvalidToken"
	CodeNeedLogin     ResKey = "CodeNeedLogin"
	CodeClientUnknown ResKey = "CodeClientUnknown"
	CodeClientEmpty   ResKey = "CodeClientEmpty"
	CodeTokenDestroy  ResKey = "CodeTokenDestroy"
	CodeLoginOther    ResKey = "CodeLoginOther"
	// Logic Message

	// Menu Name

)

func InitTranslators() (err error) {
	return
}

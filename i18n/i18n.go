package i18n

import (
	"embed"
	"encoding/json"
	"errors"

	"go.uber.org/zap"
	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
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
	Language string      `json:"language"`
	Messages []TransItem `json:"messages"`
}

type SystemLanguage struct {
	Language string       `json:"language"`
	Tag      language.Tag `json:"tag"`
	FileName string       `json:"fileName"`
	Default  bool         `json:"default"`
}

// List of languages the system wil support
var SupportLanguages []SystemLanguage = []SystemLanguage{
	{Language: "en-US", Tag: language.AmericanEnglish, FileName: "en-US.json", Default: true},
	{Language: "zh-Hans", Tag: language.SimplifiedChinese, FileName: "zh-Hans.json", Default: false},
}

// Default language message printer
var defaultPrinter *message.Printer

// Language message printers map
var printers map[string]*message.Printer

func InitTranslators() (err error) {
	//Initialize printerMap
	printers = make(map[string]*message.Printer, len(SupportLanguages))

	for _, lang := range SupportLanguages {
		// Read translation file content
		file, err := transFile.Open("translations/" + lang.FileName)
		if err != nil {
			zap.L().Error("InitTranslators  transFile.Open failed:", zap.Error(err))
			return err
		}

		// Parse JSON file
		decoder := json.NewDecoder(file)
		var ltf LangTransFile
		err = decoder.Decode(&ltf)
		if err != nil {
			zap.L().Error("InitTranslators  decoder.Decode failed:", zap.Error(err))
			return err
		}
		// Write messages
		for _, msg := range ltf.Messages {
			switch msg.Type {
			case "string":
				err = message.SetString(lang.Tag, msg.Key, msg.Message)
				if err != nil {
					zap.L().Error("InitTranslators message.SetString failed:", zap.Error(err))
					return err
				}
			case "plural":
				cm := plural.Selectf(msg.SelectfParms.Arg, msg.SelectfParms.Format, msg.SelectfParms.Cases...)
				err = message.Set(lang.Tag, msg.Key, cm)
				if err != nil {
					zap.L().Error("InitTranslators message.Set failed:", zap.Error(err))
					return err
				}
			default:
				zap.L().Error("InitTranslators msg.Type undefined")
				return errors.New("InitTranslators TransItem type undefined")
			}
		}

		p := message.NewPrinter(lang.Tag)
		if lang.Default {
			defaultPrinter = p
		}
		// assign a value to the printers
		printers[lang.Language] = p
		// Close the file
		file.Close()
	}
	return
}

// Type Reskey to string
func (r ResKey) String() string {
	return string(r)
}

// Type Reskey to Msg
func (r ResKey) Msg(lang string, params ...interface{}) (result string) {
	rs := r.String()
	p, ok := printers[lang]

	if !ok {
		result = defaultPrinter.Sprintf(rs, params...)
	} else {
		result = p.Sprintf(rs, params...)
	}
	if result == rs {
		result = defaultPrinter.Sprintf(rs, params...)
	}
	return
}

package pg

import (
	"database/sql"
	"sccsmsserver/i18n"
	"time"

	"go.uber.org/zap"
)

type LandingPageInfo struct {
	SysNameDisp string    `db:"sysnamedisp" json:"sysNameDisp"`
	IntroText   string    `db:"introtext" json:"introText"`
	File        File      `db:"fileid" json:"file"`
	ModifyDate  time.Time `db:"modifytime" json:"modifyDate"`
	Modifier    Person    `db:"modifierid" json:"modifier"`
	Ts          time.Time `db:"ts" json:"ts"`
}

// Initialize the landing page table
func initLandingPage() (isFinish bool, err error) {
	// Step 1: Check if a record exists for default value in the landingpage table
	sqlStr := "select count(fileid) as rownum from landingpage"
	hasRecord, isFinish, err := genericCheckRecord("landingpage", sqlStr)
	if hasRecord || !isFinish || err != nil {
		return
	}
	// Step 2: Insert a default record into the landingpage table
	sqlStr = `insert into landingpage(sysnamedisp,introtext,fileid,modifierid) 
		values('SeaCloud Construction Site Managemnet System',
		'An open-source construction site management system that helps managers effectively implement on-site management measures.'
		,0,10000);`
	_, err = db.Exec(sqlStr)
	if err != nil {
		isFinish = false
		zap.L().Error("landingpage insert initvalue failed", zap.Error(err))
		return isFinish, err
	}
	return
}

// Get Landing Page Info
func GetLandingPageInfo() (info LandingPageInfo, resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	sqlStr := "select sysnamedisp,introtext,fileid,modifytime,modifierid,ts from landingpage limit 1"
	err = db.QueryRow(sqlStr).Scan(&info.SysNameDisp, &info.IntroText, &info.File.ID, &info.ModifyDate, &info.Modifier.ID, &info.Ts)
	if err != nil {
		resStatus = i18n.StatusInternalError
		zap.L().Error("GetLandingPageInfo db.QueryRow failed ", zap.Error(err))
		return
	}
	// Get File details
	if info.File.ID > 0 {
		resStatus, err = info.File.GetFileInfoByID()
		if err != nil || resStatus != i18n.StatusOK {
			return
		}
	} else {
		info.File.FileUrl = "/static/img/brands/introduce.jpg"
	}
	// Get Modifier details
	if info.Modifier.ID > 0 {
		resStatus, err = info.Modifier.GetPersonInfoByID()
		if err != nil || resStatus != i18n.StatusOK {
			return
		}
	}

	return
}

// Modify Landing Page Info
func (info *LandingPageInfo) Modify() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	sqlStr := `update landingpage set sysnamedisp=$1,introtext=$2,fileid=$3,modifytime=current_timestamp,modifierid=$4,ts=current_timestamp where ts=$5
	 returning modifytime,ts,modifierid`
	err = db.QueryRow(sqlStr, info.SysNameDisp, info.IntroText, info.File.ID, info.Modifier.ID, info.Ts).Scan(&info.ModifyDate, &info.Ts, &info.Modifier.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			resStatus = i18n.StatusOtherEdit
			return
		}
		resStatus = i18n.StatusInternalError
		zap.L().Error("LandingPageInfo.Modify db.QueryRow failed", zap.Error(err))
		return
	}
	// Get File info
	if info.File.ID > 0 {
		resStatus, err = info.File.GetFileInfoByID()
		if err != nil || resStatus != i18n.StatusOK {
			return
		}
	}

	// Get Modifier info
	if info.Modifier.ID > 0 {
		resStatus, err = info.Modifier.GetPersonInfoByID()
		if err != nil || resStatus != i18n.StatusOK {
			return
		}
	}
	return
}

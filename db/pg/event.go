package pg

import (
	"sccsmsserver/i18n"
	"sccsmsserver/pub"
	"time"

	"go.uber.org/zap"
)

// User Event struct
type Event struct {
	ID              int32            `json:"id"`
	CSA             ConstructionSite `json:"csa"`
	EPT             EPT              `json:"ept"`
	Start           time.Time        `json:"start"`
	End             time.Time        `json:"end"`
	Status          int16            `json:"status"`
	Editable        bool             `json:"editable"`
	AllDay          bool             `json:"allDay"`
	BackgroundColor string           `json:"backgroundColor"`
	BillType        string           `json:"billType"`
	HID             int32            `json:"hid"`
	BillNumber      string           `json:"billNumber"`
	RowNumber       int32            `json:"rowNumber"`
	HDescription    string           `json:"hDescription"`
	BDescription    string           `json:"bDescription"`
	EpaName         string           `json:"epaName"`
	EpaValueDIsp    string           `json:"epaValueDisp"`
	Files           []VoucherFile    `json:"files"`
	Creator         Person           `json:"creator"`
}

// User Event Params
type UserEvents struct {
	ID           int32     `json:"userID"`
	Start        time.Time `json:"start"`
	End          time.Time `json:"end"`
	ResultNumber int32     `json:"resultNumber"`
	Events       []Event   `json:"events"`
}

// Retrieve User Events
func (ue *UserEvents) GetEvents() (resStatus i18n.ResKey, err error) {
	resStatus = i18n.StatusOK
	ue.Events = make([]Event, 0)
	// Retrieve Events from Work Order
	sqlStr := `select 
	b.id as bid,
	b.csaid,
	b.eptid,
	ept_h.name,
	b.starttime as starttime,
	b.endtime as endtime,
	b.status,
	false as editable,
	false as allday,
	'WO' as billtype,
	h.id as hid,
	h.billnumber,
	b.rownumber,
	h.description as hdescription,
	b.description as bdescription,
	b.creatorid
	from workorder_b as b
	left join workorder_h as h on b.hid=h.id
	left join ept_h on b.eptid = ept_h.id
	left join csa on b.csaid = csa.id
	where b.dr=0 and h.dr=0 and b.executorid=$1 and b.starttime >= $2 and starttime<=$3`
	rows, err := db.Query(sqlStr, ue.ID, ue.Start, ue.End)
	if err != nil {
		zap.L().Error("UserEvents.GetEvents db.Query(sqlStr) failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer rows.Close()
	// Extract Event row by row
	for rows.Next() {
		var e Event
		err = rows.Scan(&e.ID, &e.CSA.ID, &e.EPT.HID, &e.EPT.Name, &e.Start,
			&e.End, &e.Status, &e.Editable, &e.AllDay, &e.BillType,
			&e.HID, &e.BillNumber, &e.RowNumber, &e.HDescription, &e.BDescription,
			&e.Creator.ID)
		if err != nil {
			zap.L().Error("UserEvents.GetEvents rows.Next() failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}
		ue.ResultNumber++
		// Get Construction Site details
		if e.CSA.ID > 0 {
			resStatus, err = e.CSA.GetInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Creator Details
		if e.Creator.ID > 0 {
			resStatus, err = e.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Color
		e.BackgroundColor = pub.EventBackgroundColors[e.Status]
		ue.Events = append(ue.Events, e)
	}

	// Retrieve Events form Execution Order
	edSqlStr := `select 
	b.id, 
	h.csaid,
	h.eptid,
	ept_h.name,
	b.handlestarttime as starttime,
	b.handleendtime as endtime,
	b.status,
	false as editable,
	false as allday,
	'EO' as billtype,
	h.id as hid,
	h.billnumber,
	b.rownumber,
	h.description as hdescription,
	b.description as bdescription,
	b.creatorid,
	epa.name,
	b.executionvaluedisp 
	from executionorder_b as b 
	left join executionorder_h as h on b.hid = h.id
	left join epa on b.epaid = epa.id
	left join csa as csa on h.csaid = csa.id
	left join ept_h as ept_h on h.eptid = ept_h.id
	where b.ishandle=1 and b.dr=0 and h.dr=0 and b.issueownerid=$1 and b.handlestarttime >= $2 and b.handlestarttime<=$3 ;`
	edRows, err := db.Query(edSqlStr, ue.ID, ue.Start, ue.End)
	if err != nil {
		zap.L().Error("UserEvents.GetEvents db.Query(edSqlStr) failed", zap.Error(err))
		resStatus = i18n.StatusInternalError
		return
	}
	defer edRows.Close()
	// Extract event row by row
	for edRows.Next() {
		var e Event
		err = edRows.Scan(&e.ID, &e.CSA.ID, &e.EPT.HID, &e.EPT.Name, &e.Start,
			&e.End, &e.Status, &e.Editable, &e.AllDay, &e.BillType,
			&e.HID, &e.BillNumber, &e.RowNumber, &e.HDescription, &e.BDescription,
			&e.Creator.ID, &e.EpaName, &e.EpaValueDIsp)
		if err != nil {
			zap.L().Error("UserEvents.GetEvents edRows.Next() failed", zap.Error(err))
			resStatus = i18n.StatusInternalError
			return
		}
		ue.ResultNumber++
		// Get Constrution Site details
		if e.CSA.ID > 0 {
			resStatus, err = e.CSA.GetInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Creator details
		if e.Creator.ID > 0 {
			resStatus, err = e.Creator.GetPersonInfoByID()
			if resStatus != i18n.StatusOK || err != nil {
				return
			}
		}
		// Get Attachments
		e.Files, resStatus, err = GetEORowFiles(e.ID)
		if resStatus != i18n.StatusOK || err != nil {
			return
		}

		// Get Color
		e.BackgroundColor = pub.EventBackgroundColors[e.Status]
		ue.Events = append(ue.Events, e)
	}
	return
}

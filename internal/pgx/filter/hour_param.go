package filter

import "time"

type HourParam struct {
	ID        int64     `json:"ID"`
	Value     int64     `json:"Value"`
	Timestamp time.Time `json:"Timestamp"`
	DateFrom  time.Time `json:"DateFrom"`
	DateTo    time.Time `json:"DateTo"`
	ParamID   int64     `json:"ParamID"`
	XmlCreate bool      `json:"XmlCreate"`
	Manual    bool      `json:"Manual"`
	ChangeBy  string    `json:"ChangeBy"`
	Comment   string    `json:"Comment"`
	Limit     int64     `json:"Limit"`
	Offset    int64     `json:"Offset"`
}

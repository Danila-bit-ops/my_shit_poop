package filter

import "time"

type HourParam struct {
	ID        int64
	Value     int64
	Timestamp time.Time
	ParamID   int64
	XmlCreate bool
	Manual    bool
	ChangeBy  string
	Comment   string
	Limit     int64
}

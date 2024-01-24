package model

import (
	"fmt"
	"time"
)

type HourParamList []HourParam

func (l HourParamList) ToHTMLTable() string {
	var table string

	for _, v := range l {
		id := fmt.Sprintf("<td>%d</td>", v.ID)
		paramID := fmt.Sprintf("<td>%d</td>", v.ParamID)
		val := fmt.Sprintf("<td>%g</td>", v.Val)
		timestamp := fmt.Sprintf("<td>%s</td>", v.Timestamp.Format(time.DateTime))
		table += "<tr>" + id + paramID + val + timestamp + "</tr>"
	}

	return table
}

type HourParam struct {
	Timestamp time.Time `json:"Timestamp"`
	ChangeBy  string    `json:"ChangeBy"`
	Comment   string    `json:"Comment"`
	ID        int64     `json:"ID"`
	Val       float64   `json:"Val"`
	ParamID   int64     `json:"ParamID"`
	XMLCreate bool      `json:"XMLCreate"`
	Manual    bool      `json:"Manual"`
}

type Column struct {
	Name     string `toml:"name"`
	DataType string `toml:"data_type"`
}

type Table struct {
	TableName string   `toml:"TableName"`
	Columns   []Column `toml:"Columns"`
}

type TableConfig struct {
	IP       string
	DBName   string
	Login    string
	Password string
	Sslmode  string
	Tables   []Table `toml:"tables"`
}

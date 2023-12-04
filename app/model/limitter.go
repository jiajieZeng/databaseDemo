package model

import "time"

// EventLog 搜集日志的结构
type EventLog struct {
	EventTime time.Time `json:"time"`
	Log       string    `json:"info"`
}

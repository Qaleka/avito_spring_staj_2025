package models

import "time"

const (
	STATUS_ACTIVE = "in_progress"
	STATUS_CLOSED = "closed"
)

type Reception struct {
	Id       string
	DateTime time.Time
	PvzId    string
	Status   string
}

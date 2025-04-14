package models

import "time"

type Pvz struct {
	Id               string
	RegistrationDate time.Time
	City             string
	Receptions       []Reception `json:"-"`
}

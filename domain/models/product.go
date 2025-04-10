package models

import "time"

const (
	ELECTRONIC_TYPE = "электроника"
	CLOTHES_TYPE    = "одежда"
	BOOTS_TYPE      = "обувь"
)

type Product struct {
	Id          string
	DateTime    time.Time
	Type        string
	ReceptionId string
}

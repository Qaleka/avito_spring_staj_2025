package responses

import "time"

type RegisterResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type CreatePvzResponse struct {
	Id               string    `json:"id"`
	RegistrationDate time.Time `json:"registrationDate"`
	City             string    `json:"city"`
}

type CreateReceptionResponse struct {
	Id       string    `json:"id"`
	DateTime time.Time `json:"dateTime"`
	PvzId    string    `json:"pvzId"`
	Status   string    `json:"status"`
}

type AddProductResponse struct {
	Id          string    `json:"id"`
	DateTime    time.Time `json:"dateTime"`
	Type        string    `json:"type"`
	ReceptionId string    `json:"receptionId"`
}

type CloseReceptionResponse struct {
	Id       string    `json:"id"`
	DateTime time.Time `json:"dateTime"`
	PvzId    string    `json:"pvzId"`
	Status   string    `json:"status"`
}

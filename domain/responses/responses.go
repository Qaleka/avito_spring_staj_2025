package responses

import (
	"avito_spring_staj_2025/domain/models"
	"time"
)

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

type GetPvzsInformationResponse struct {
	Pvz        models.Pvz                 `json:"pvz"`
	Receptions []GetReceptionWithProducts `json:"receptions"`
}

type GetReceptionWithProducts struct {
	Reception models.Reception `json:"reception"`
	Products  []models.Product `json:"products"`
}

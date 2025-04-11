package requests

import "time"

type DummyLoginRequest struct {
	Role string `json:"role"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreatePvzRequest struct {
	Id               string    `json:"id"`
	RegistrationDate time.Time `json:"registrationDate"`
	City             string    `json:"city"`
}

type CreateReceptionRequest struct {
	PvzId string `json:"pvzId"`
}

type GetPvzInfoQueries struct {
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
	Page      int       `json:"page"`
	Limit     int       `json:"limit"`
}

type AddProductRequest struct {
	Type  string `json:"type"`
	PvzId string `json:"pvzId"`
}

type GetPvzListRequest struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
}

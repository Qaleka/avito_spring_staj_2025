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
	Id string `json:"id"`
}

type AddProductRequest struct {
	Type  string `json:"type"`
	PvzId string `json:"pvzId"`
}

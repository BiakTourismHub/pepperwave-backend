package models

type AuthRegisterRequest struct {
	Fullname  string `json:"fullname"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
	Role      string `json:"role"`
	AccountId string `json:"account_id"`
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

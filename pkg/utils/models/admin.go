package models

type AdminLogin struct {
	Email    string `json:"email,omitempty" validate:"required"`
	Password string `json:"password" validate:"min=8,max=20"`
}
type AdminDetailsResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Email string `json:"email" validate:"email"`
}

type AdminResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`

}
type UserDetailsAtAdmin struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email" validate:"email"`
	Phone   string `json:"phone"`
	Blocked bool   `json:"blocked"`
}
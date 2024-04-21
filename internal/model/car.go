package model

type Car struct {
	Id     int    `json:"id" example:"1"`
	RegNum string `json:"regNum" example:"X123XX150"`
	Mark   string `json:"mark" example:"Lada"`
	Model  string `json:"model" example:"Vesta"`
	Year   int32  `json:"year,omitempty" example:"2002"`
	Owner  People `json:"owner"`
}

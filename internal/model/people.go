package model

type People struct {
	Id         int    `json:"id" example:"1"`
	Name       string `json:"name" example:"string"`
	Surname    string `json:"surname" example:"string"`
	Patronymic string `json:"patronymic,omitempty" example:"string"`
}

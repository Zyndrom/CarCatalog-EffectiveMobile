package repository

import (
	"CarsCatalog/internal/model"
	"CarsCatalog/internal/repository/postgres"
	"CarsCatalog/internal/repository/utils"
)

type repository struct {
	db dbStorage
}

type dbStorage interface {
	AddNewCar(car model.Car) error
	AddNewPeople(people model.People) (int, error)
	GetCarsWithOwners(filters []utils.Filter, search utils.SearchOptions) ([]model.Car, error)
	DeleteCarById(id int) error
	UpdateCarById(id int, car model.Car) error
	UpdatePeopleById(id int, people model.People) error
	ContainsCarById(id int) (bool, error)
	ContainsPeopleById(id int) (bool, error)
}

func New() *repository {
	return &repository{db: postgres.New()}
}

func (r *repository) AddNewCar(car model.Car) error {
	return r.db.AddNewCar(car)
}
func (r *repository) AddNewPeople(people model.People) (int, error) {
	return r.db.AddNewPeople(people)
}
func (r *repository) GetCarsWithOwners(query map[string][]string) ([]model.Car, error) {
	filters := utils.GetFilterOptions(query)
	options := utils.GetSearchOptions(query)
	return r.db.GetCarsWithOwners(filters, options)
}
func (r *repository) DeleteCarById(id int) error {
	return r.db.DeleteCarById(id)
}
func (r *repository) UpdateCarById(id int, car model.Car) error {
	return r.db.UpdateCarById(id, car)
}
func (r *repository) UpdatePeopleById(id int, people model.People) error {
	return r.db.UpdatePeopleById(id, people)
}
func (r *repository) ContainsCarById(id int) (bool, error) {
	return r.db.ContainsCarById(id)
}
func (r *repository) ContainsPeopleById(id int) (bool, error) {
	return r.db.ContainsPeopleById(id)
}

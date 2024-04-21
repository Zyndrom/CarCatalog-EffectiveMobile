package car

import (
	"CarsCatalog/internal/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

type carService struct {
	storage storage
}
type storage interface {
	AddNewCar(car model.Car) error
	AddNewPeople(people model.People) (int, error)
	GetCarsWithOwners(query map[string][]string) ([]model.Car, error)
	DeleteCarById(id int) error
	UpdateCarById(id int, car model.Car) error
	ContainsCarById(id int) (bool, error)
	UpdatePeopleById(id int, people model.People) error
	ContainsPeopleById(id int) (bool, error)
}

var (
	externalApi string
)

func New(storage storage) *carService {
	externalApi = os.Getenv("EXTERNAL_CAR_API")
	return &carService{storage: storage}
}

func (c *carService) AddNewCars(regNums []string) error {
	newCars := []model.Car{}
	for _, regNum := range regNums {
		query := fmt.Sprintf("%s/info?regNum=%s", externalApi, regNum)
		resp, err := http.Get(query)
		if err != nil {
			logrus.Debug(err)
			return err
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil || body == nil {
			logrus.Debug(err)
			return err
		}
		var car model.Car
		if err := json.Unmarshal(body, &car); err != nil {
			logrus.Debugf("Can not unmarshal JSON, %s", err.Error())
			return err
		}
		newCars = append(newCars, car)
	}

	for _, car := range newCars {
		id, err := c.AddNewPeople(car.Owner)
		if err != nil {
			return err
		}
		car.Owner.Id = id

		err = c.addNewCar(car)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *carService) ContainsCarById(id int) (bool, error) {
	return c.storage.ContainsCarById(id)
}

func (c *carService) addNewCar(car model.Car) error {
	return c.storage.AddNewCar(car)
}
func (c *carService) AddNewPeople(people model.People) (int, error) {
	return c.storage.AddNewPeople(people)
}
func (c *carService) GetCarsWithOwners(query map[string][]string) ([]model.Car, error) {
	return c.storage.GetCarsWithOwners(query)
}
func (c *carService) DeleteCarById(id int) error {
	return c.storage.DeleteCarById(id)
}
func (c *carService) UpdateCarById(id int, car model.Car) error {
	return c.storage.UpdateCarById(id, car)
}
func (c *carService) UpdatePeopleById(id int, people model.People) error {
	return c.storage.UpdatePeopleById(id, people)
}

func (c *carService) ContainsPeopleById(id int) (bool, error) {
	return c.storage.ContainsPeopleById(id)
}

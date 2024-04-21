package postgres

import (
	"CarsCatalog/internal/model"
	"CarsCatalog/internal/repository/utils"
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type postgresql struct {
	db *sql.DB
}

func New() *postgresql {
	usr := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	dbName := os.Getenv("POSTGRES_DATABASE")
	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
		usr, pass, host, port, dbName)

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		logrus.Fatalf(err.Error())
	}
	err = db.Ping()
	if err != nil {
		logrus.Fatalf(err.Error())
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logrus.Fatalf(err.Error())
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)

	if err != nil {
		logrus.Fatalf(err.Error())
	}
	m.Up()
	psg := &postgresql{db: db}
	return psg

}

func (p *postgresql) AddNewCar(car model.Car) error {
	query := `INSERT INTO cars (reg_num, mark, model, year, owner_id) VALUES ($1, $2, $3, $4, $5)`
	_, err := p.db.Exec(query, car.RegNum, car.Mark, car.Model, car.Year, car.Owner.Id)
	return err
}

func (p *postgresql) AddNewPeople(people model.People) (int, error) {
	query := `INSERT INTO peoples (name, surname,patronymic) VALUES ($1, $2, $3) returning id`
	var id int
	err := p.db.QueryRow(query, people.Name, people.Surname, people.Patronymic).Scan(&id)
	return id, err
}

func (p *postgresql) ContainsCarById(id int) (bool, error) {
	query := `SELECT COUNT(*) FROM cars WHERE id=$1`
	row := p.db.QueryRow(query, id)
	var count int
	err := row.Scan(&count)
	if err != nil {
		logrus.Fatalln(err)
		return true, err
	}
	return count > 0, nil

}

func (p *postgresql) ContainsPeopleById(id int) (bool, error) {
	query := `SELECT COUNT(*) FROM peoples WHERE id=$1`
	row := p.db.QueryRow(query, id)
	var count int
	err := row.Scan(&count)
	if err != nil {
		logrus.Fatalln(err)
		return true, err
	}
	return count > 0, nil

}

func (p *postgresql) GetCarsWithOwners(filters []utils.Filter, search utils.SearchOptions) ([]model.Car, error) {
	filterString := generateFilter(filters)
	var query string
	if filterString != "" {
		query = fmt.Sprintf(`SELECT cars.id, cars.reg_num, cars.mark, cars.model, cars.year, cars.owner_id, peoples.name, peoples.surname, peoples.patronymic 
		FROM cars JOIN peoples on peoples.id = cars.owner_id WHERE %s ORDER BY id %s LIMIT %s OFFSET %s;`,
			filterString, search.Order, search.Limit, search.Offset)
	} else {
		query = fmt.Sprintf(`SELECT cars.id, cars.reg_num, cars.mark, cars.model, cars.year, cars.owner_id, peoples.name, peoples.surname, peoples.patronymic 
		FROM cars JOIN peoples on peoples.id = cars.owner_id ORDER BY id %s LIMIT %s OFFSET %s;`,
			search.Order, search.Limit, search.Offset)
	}
	rows, err := p.db.Query(query)

	if err != nil {
		return nil, err
	}
	cars := []model.Car{}
	for rows.Next() {
		car := model.Car{
			Owner: model.People{},
		}
		year := sql.NullInt32{}
		patronymic := sql.NullString{}
		err := rows.Scan(&car.Id, &car.RegNum, &car.Mark, &car.Model, &year, &car.Owner.Id, &car.Owner.Name, &car.Owner.Surname, &patronymic)
		car.Year = year.Int32
		car.Owner.Patronymic = patronymic.String
		if err != nil {
			continue
		}
		cars = append(cars, car)
	}
	return cars, nil
}

func (p *postgresql) DeleteCarById(id int) error {
	query := `DELETE FROM cars WHERE id=$1;`
	_, err := p.db.Exec(query, id)
	return err
}

func (p *postgresql) UpdateCarById(id int, car model.Car) error {
	query := `UPDATE cars SET `
	args := []interface{}{}
	updateFields := []string{}
	i := 1
	if car.RegNum != "" {
		updateFields = append(updateFields, fmt.Sprintf("reg_num = $%d", i))
		args = append(args, car.RegNum)
		i++
	}

	if car.Model != "" {
		updateFields = append(updateFields, fmt.Sprintf("model = $%d", i))
		args = append(args, car.Model)
		i++
	}

	if car.Mark != "" {
		updateFields = append(updateFields, fmt.Sprintf("mark = $%d", i))
		args = append(args, car.Mark)
		i++
	}
	if car.Year != 0 {
		updateFields = append(updateFields, fmt.Sprintf("year = $%d", i))
		args = append(args, car.Year)
		i++
	}
	if car.Owner.Id != 0 {
		updateFields = append(updateFields, fmt.Sprintf("owner_id = $%d", i))
		args = append(args, car.Owner.Id)
		i++
	}

	args = append(args, car.Id)

	query += strings.Join(updateFields, ", ") + fmt.Sprintf(" WHERE id = $%d;", i)
	_, err := p.db.Exec(query, args...)
	return err
}

func (p *postgresql) UpdatePeopleById(id int, people model.People) error {
	query := `UPDATE peoples SET `
	args := []interface{}{}
	updateFields := []string{}
	i := 1
	if people.Name != "" {
		updateFields = append(updateFields, fmt.Sprintf("name = $%d", i))
		args = append(args, people.Name)
		i++
	}

	if people.Surname != "" {
		updateFields = append(updateFields, fmt.Sprintf("surname = $%d", i))
		args = append(args, people.Surname)
		i++
	}

	if people.Patronymic != "" {
		updateFields = append(updateFields, fmt.Sprintf("patronymic = $%d", i))
		args = append(args, people.Patronymic)
		i++
	}

	args = append(args, people.Id)

	query += strings.Join(updateFields, ", ") + fmt.Sprintf(" WHERE id = $%d;", i)

	_, err := p.db.Exec(query, args...)
	return err
}

func generateFilter(filters []utils.Filter) string {
	if len(filters) == 0 {
		return ""
	}
	var builder bytes.Buffer

	for i, val := range filters {
		if val.Operator != string(utils.Between) {

			builder.WriteString(fmt.Sprintf("%s %s '%s'", val.Field, transformOperator(val.Operator), val.Value))

		} else {
			if strings.Contains(val.Value, ":") {
				_, err := strconv.Atoi(val.Value[:strings.Index(val.Value, ":")])
				if err != nil {
					continue
				}
				_, err = strconv.Atoi(val.Value[strings.Index(val.Value, ":")+1:])
				if err != nil {
					continue
				}
				builder.WriteString(val.Field + "BETWEEN" + val.Value[:strings.Index(val.Value, ":")] + "AND" +
					val.Value[strings.Index(val.Value, ":")+1:])
			} else {
				continue
			}
		}
		if i+1 < len(filters) {
			builder.WriteString(" AND ")
		}
	}
	return builder.String()
}

func transformOperator(operator string) string {
	switch operator {
	case string(utils.LowerThan):
		return "<"
	case string(utils.LowerThanEq):
		return "<="
	case string(utils.GreaterThan):
		return ">"
	case string(utils.GreaterThanEq):
		return ">="
	default:
		return "="
	}
}

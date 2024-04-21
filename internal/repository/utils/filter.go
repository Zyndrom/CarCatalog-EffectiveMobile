package utils

import "strings"

const (
	carId    = "id"
	regNum   = "reg_num"
	carMark  = "mark"
	carModel = "model"
	year     = "year"
)

const (
	LowerThan     = "lt"
	LowerThanEq   = "ltq"
	GreaterThan   = "gt"
	GreaterThanEq = "gtq"
	Equal         = "eq"
	Between       = "between"
)

type Filter struct {
	Field    string
	Operator string
	Value    string
}

func GetFilterOptions(query map[string][]string) []Filter {
	filters := []Filter{}
	carId, ok := query[carId]
	if ok {
		filters = append(filters, Filter{
			Field:    "cars.id",
			Operator: string(Equal),
			Value:    carId[0],
		})
	}
	regNum, ok := query[regNum]
	if ok {
		filters = append(filters, Filter{
			Field:    "cars.reg_num",
			Operator: string(Equal),
			Value:    regNum[0],
		})
	}
	carMark, ok := query[carMark]
	if ok {
		filters = append(filters, Filter{
			Field:    "cars.mark",
			Operator: string(Equal),
			Value:    carMark[0],
		})
	}
	carModel, ok := query[carModel]
	if ok {
		filters = append(filters, Filter{
			Field:    "cars.model",
			Operator: string(Equal),
			Value:    carModel[0],
		})

	}
	yearQuery, ok := query[year]
	if ok {
		var year string
		operator := string(Equal)
		if strings.Contains(yearQuery[0], ":") {
			operator = yearQuery[0][:strings.Index(yearQuery[0], ":")]
			year = yearQuery[0][strings.Index(yearQuery[0], ":")+1:]
		} else {
			year = yearQuery[0]
		}

		filters = append(filters, Filter{
			Field:    "cars.year",
			Operator: operator,
			Value:    year,
		})
	}

	return filters
}

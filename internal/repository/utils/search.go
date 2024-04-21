package utils

import "strconv"

const (
	Asc  = "asc"
	Desc = "desc"
)

const (
	Order  = "order"
	Limit  = "limit"
	Offset = "offset"
)

const (
	defaultOffset = "0"
	defaultLimit  = "10"
	defaultOrder  = Desc
)

type SearchOptions struct {
	Order  string
	Limit  string
	Offset string
}

func GetSearchOptions(query map[string][]string) SearchOptions {
	searchOptions := SearchOptions{
		Order:  defaultOrder,
		Limit:  defaultLimit,
		Offset: defaultOffset,
	}
	order, ok := query[Order]
	if ok {
		if order[0] == Asc || order[0] == Desc {
			searchOptions.Order = order[0]
		} else {
			searchOptions.Order = Desc
		}
	}
	limit, ok := query[Limit]
	limitAsNum := 10
	if ok {
		var err error
		limitAsNum, err = strconv.Atoi(limit[0])
		if err != nil {
			searchOptions.Limit = defaultLimit
		} else {
			searchOptions.Limit = limit[0]
		}
	}
	offset, ok := query[Offset]
	if ok {
		offsetAsNum, err := strconv.Atoi(offset[0])
		if err != nil {
			searchOptions.Offset = defaultOffset
		} else {
			offsetVal := strconv.Itoa(offsetAsNum * limitAsNum)
			searchOptions.Offset = offsetVal
		}
	}
	return searchOptions
}

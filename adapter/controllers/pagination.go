package controllers

import (
	"net/url"
	"strconv"
)

const (
	defaultPageLimit = 100
	defaultStartPage = 1
	defaultSort      = ""
)

type Pagination struct {
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
	Sort  string `json:"sort"`
}

func GeneratePaginationFromRequest(query url.Values) Pagination {
	// Initializing default
	limit := defaultPageLimit
	page := defaultStartPage
	sort := defaultSort

	for key, value := range query {
		queryValue := value[len(value)-1]
		switch key {
		case "limit":
			limit, _ = strconv.Atoi(queryValue)
			break
		case "page":
			page, _ = strconv.Atoi(queryValue)
			break
		case "sort":
			sort = queryValue
			break

		}
	}
	return Pagination{
		Limit: limit,
		Page:  page,
		Sort:  sort,
	}

}

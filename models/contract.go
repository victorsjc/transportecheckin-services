package models

type Contract struct {
	Id                  string `json:"id"`
	CompanyId           string `json:"companyId"`
	Type                string `json:"type"`
	Name                string `json:"name"`
	Contact             string `json:"contact"`
	Local               string `json:"local"`
	DaysOfWeek          string `json:"daysOfWeek"`
	HR			        string `json:"hr"`
	Status				string `json:"status"`
}
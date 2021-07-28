package rabotaua

import (
	"fmt"
)

//easyjson:json
type Cities []City

type City struct {
	ID   int    `json:"id"`
	City string `json:"ua"`
}

type SearchResult struct {
	Took         int       `json:"took"`
	Start        int       `json:"start"`
	Count        int       `json:"count"`
	Total        int       `json:"total"`
	ErrorMessage string    `json:"errorMessage"`
	Vacancy      []Vacancy `json:"documents"`
}

type ScheduleMap map[int]string

//easyjson:json
type schedules []schedulesStruct

type schedulesStruct struct {
	ID int    `json:"id"`
	Ua string `json:"ua"`
}

type Vacancy struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Date             string `json:"date"`
	DateTxt          string `json:"dateTxt"`
	CityName         string `json:"cityName"`
	NotebookID       int    `json:"notebookId"`
	CompanyName      string `json:"companyName"`
	ShortDescription string `json:"shortDescription"`
}

func (v Vacancy) GetURL() string {
	return fmt.Sprintf("https://rabota.ua/ua/company%d/vacancy%d", v.NotebookID, v.ID)
}

type VacancyParameters struct {
	Keywords   string `json:"k"`
	CityID     int    `json:"ci"`
	ScheduleID int    `json:"si"`
}

type VacancyParametersPage struct {
	VacancyParameters VacancyParameters `json:"vacancy_parameters"`
	Page              int               `json:"page"`
	Count             int               `json:"count"`
}

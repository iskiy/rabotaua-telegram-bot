package rabotaua

import (
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
	"net/http"
	"net/url"
	"strconv"
)

const (
	homePage                      = "https://api.rabota.ua"
	autocompleteCityEndPoint      = "/autocomplete/city"
	vacancySearchEndPoint         = "/vacancy/search"
	getScheduleDictionaryEndPoint = "/dictionary/schedule"
	dictionaryCityEndPoint        = "/dictionary/city"
)

var (
	CantFindCityError = errors.New("can`t find city: ")
)

type RabotaClient struct {
	client fasthttp.Client
}

func NewRabotaClient() *RabotaClient {
	return &RabotaClient{client: fasthttp.Client{}}
}

func (c *RabotaClient) GetSearchResultFromParameters(p VacancyParameters) (SearchResult, error) {
	requestURL := getVacanciesRequestURLFromParams(p)
	return c.getSearchResultFromURL(requestURL)
}

func (c *RabotaClient) GetSearchResultFromParametersPage(p VacancyParametersPage) (SearchResult, error) {
	requestURL := getVacanciesRequestURLFromParamsPage(p)
	return c.getSearchResultFromURL(requestURL)
}

func getVacanciesRequestURLFromParams(p VacancyParameters) string {
	params := getURLValuesFromVacancyParameters(p)
	return getVacancySearchURLFromURLValues(params)
}

func getVacanciesRequestURLFromParamsPage(p VacancyParametersPage) string {
	params := getURLValuesFromVacancyParametersPage(p)
	return getVacancySearchURLFromURLValues(params)
}

func getVacancySearchURLFromURLValues(params url.Values) string {
	return homePage + vacancySearchEndPoint + "?" + params.Encode()
}

func getURLValuesFromVacancyParameters(p VacancyParameters) url.Values {
	params := url.Values{}
	params.Add("keyWords", p.Keywords)
	params.Add("ukrainian", "true")
	params.Add("scheduleId", strconv.Itoa(p.ScheduleID))
	params.Add("cityId", strconv.Itoa(p.CityID))
	return params
}

func getURLValuesFromVacancyParametersPage(p VacancyParametersPage) url.Values {
	params := getURLValuesFromVacancyParameters(p.VacancyParameters)
	params.Add("count", strconv.Itoa(p.Count))
	params.Add("page", strconv.Itoa(p.Page))
	return params
}

func (c *RabotaClient) getSearchResultFromURL(URL string) (SearchResult, error) {
	status, body, err := c.client.Get(nil, URL)
	if err != nil {
		return SearchResult{}, err
	}
	if status != http.StatusOK {
		return SearchResult{}, fmt.Errorf("bad status code: %d", status)
	}
	return getSearchResultFromBytes(body)
}

func getSearchResultFromBytes(body []byte) (SearchResult, error) {
	searchResult := &SearchResult{}
	err := searchResult.UnmarshalJSON(body)
	if err != nil {
		return SearchResult{}, err
	}
	return *searchResult, nil
}

func (c *RabotaClient) GetCityFromName(cityName string) (City, error) {
	URL := getCityIDURL(cityName)
	return c.getCityFromURL(URL)
}

func (c *RabotaClient) GetCities() (Cities, error) {
	return c.getCitiesFromURL(homePage + dictionaryCityEndPoint)
}

func (c *RabotaClient) getCitiesFromURL(URL string) (Cities, error) {
	status, body, err := c.client.Get(nil, URL)
	if err != nil {
		return Cities{}, err
	}
	if status != http.StatusOK {
		return Cities{}, fmt.Errorf("bad http status: %d", status)
	}
	return getCitiesFromBytes(body)
}

func (c *RabotaClient) getCityFromURL(URL string) (City, error) {
	cities, err := c.getCitiesFromURL(URL)
	if err != nil {
		return City{}, err
	}
	if len(cities) == 0 {
		return City{}, CantFindCityError
	}
	return cities[0], nil
}

func getCityIDURL(cityName string) string {
	par := url.Values{}
	par.Add("term", cityName)
	return homePage + autocompleteCityEndPoint + "?" + par.Encode()
}

func getCitiesFromBytes(body []byte) ([]City, error) {
	cities := &Cities{}
	err := cities.UnmarshalJSON(body)
	if err != nil {
		return Cities{}, err
	}
	return *cities, nil
}

func (c *RabotaClient) GetSchedulesMap() (ScheduleMap, error) {
	URL := homePage + getScheduleDictionaryEndPoint
	return c.getSchedulesMapFromURL(URL)
}

func (c *RabotaClient) getSchedulesMapFromURL(URL string) (ScheduleMap, error) {
	status, body, err := c.client.Get(nil, URL)
	if err != nil {
		return ScheduleMap{}, err
	}

	if status != http.StatusOK {
		return ScheduleMap{}, fmt.Errorf("bad http status: %d", status)
	}
	return getSchedulesFromBytes(body)
}

func getSchedulesFromBytes(body []byte) (ScheduleMap, error) {
	schedules := &schedules{}
	err := schedules.UnmarshalJSON(body)
	if err != nil {
		return ScheduleMap{}, err
	}
	res := make(ScheduleMap)
	for _, s := range *schedules {
		res[s.ID] = s.Ua
	}
	return res, nil
}

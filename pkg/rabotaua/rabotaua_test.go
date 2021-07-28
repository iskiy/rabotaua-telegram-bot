package rabotaua

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// TODO: Errors tests
const (
	dataSourceVacancies = "test_vacancies_data.json"
	dataSourceCities    = "test_cities_data.json"
	dataSourceSchedules = "test_schedules_data.json"
)

var (
	rabotaClient = NewRabotaClient()
)

func MockGetVacanciesServerHandler(w http.ResponseWriter, r *http.Request) {
	mockWriteResponse(w, dataSourceVacancies)
}

func MockGetCityIDServerHandler(w http.ResponseWriter, r *http.Request) {
	mockWriteResponse(w, dataSourceCities)
}

func MockGetSchedulesServerHandler(w http.ResponseWriter, r *http.Request) {
	mockWriteResponse(w, dataSourceSchedules)
}

func mockWriteResponse(w http.ResponseWriter, dataSource string) {
	responseBody, err := readBytes(dataSource)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-type", "application/json")
	_, err = w.Write(responseBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func readBytes(dataSource string) ([]byte, error) {
	file, err := os.Open(dataSource)
	if err != nil {
		return []byte{}, err
	}
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return []byte{}, err
	}
	return bytes, nil
}

func TestMockGetVacancies(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(MockGetVacanciesServerHandler))
	defer mockServer.Close()
	res, err := rabotaClient.getSearchResultFromURL(mockServer.URL)
	if err != nil {
		t.Errorf(err.Error())
	}
	if &res == nil || res.Count != 20 {
		t.Errorf("wrong result count - %d != 20", res.Count)
	}
}

func TestGetVacanciesFromParameters(t *testing.T) {
	parameters := VacancyParameters{"golang", 1, 1}
	_, err := rabotaClient.GetSearchResultFromParameters(parameters)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}

func BenchmarkMockGetVacancies(b *testing.B) {
	mockServer := httptest.NewServer(http.HandlerFunc(MockGetVacanciesServerHandler))
	defer mockServer.Close()
	for n := 0; n < b.N; n++ {
		_, err := rabotaClient.getSearchResultFromURL(mockServer.URL)
		if err != nil {
			b.Errorf(err.Error())
		}
	}
}

func TestMockGetCityID(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(MockGetCityIDServerHandler))
	defer mockServer.Close()
	city, err := rabotaClient.getCityFromURL(mockServer.URL)
	if err != nil {
		t.Errorf(err.Error())
	}
	if &city == nil || city.ID != 1 {
		t.Errorf("wrong result: %d != 1", city.ID)
	}
}

func TestCityIDFromName(t *testing.T) {
	cityName := "Ки"
	_, err := rabotaClient.GetCityFromName(cityName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func BenchmarkMockGetCityID(b *testing.B) {
	mockServer := httptest.NewServer(http.HandlerFunc(MockGetCityIDServerHandler))
	defer mockServer.Close()
	for n := 0; n < b.N; n++ {
		_, err := rabotaClient.getCityFromURL(mockServer.URL)
		if err != nil {
			b.Errorf(err.Error())
		}
	}
}

func TestRabotaClient_GetSchedulesMapFromURL(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(MockGetSchedulesServerHandler))
	defer mockServer.Close()
	schedules, err := rabotaClient.getSchedulesMapFromURL(mockServer.URL)
	if err != nil {
		t.Errorf(err.Error())
	}
	fullTime := schedules[1]
	if fullTime != "повна зайнятість" {
		t.Errorf("wrong result: %s != \"повна зайнятість\"", fullTime)
	}
}

func TestRabotaClient_GetSchedulesMap(t *testing.T) {
	_, err := rabotaClient.GetSchedulesMap()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
}

func BenchmarkMockGetSchedulesFromURL(b *testing.B) {
	mockServer := httptest.NewServer(http.HandlerFunc(MockGetSchedulesServerHandler))
	defer mockServer.Close()
	for n := 0; n < b.N; n++ {
		_, err := rabotaClient.getSchedulesMapFromURL(mockServer.URL)
		if err != nil {
			b.Errorf(err.Error())
		}
	}
}

func TestVacancyParametersPage_JSONParsing(t *testing.T) {
	params := VacancyParameters{"dasdsa", 3, 5}
	pageParams := VacancyParametersPage{params, 5, 6}
	jsonBody, err := pageParams.MarshalJSON()
	if err != nil {
		t.Errorf(err.Error())
	}
	pageParamsUnmarshal := &VacancyParametersPage{}
	err = pageParamsUnmarshal.UnmarshalJSON(jsonBody)
	if err != nil {
		t.Errorf(err.Error())
	}
	if pageParamsUnmarshal.Page != pageParams.Page {
		t.Errorf("unexpected unmarshal result")
	}
}

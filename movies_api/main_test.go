package main_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	request "movies/Contract/Request"
	response "movies/Contract/Response"
	"net/http"
	"os"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err.Error())
	}

	testEnvironment := os.Getenv("TEST_ENVIRONMENT")

	if !strings.Contains(testEnvironment, "staging") {
		fmt.Println("To run this file, set the test environment to staging first, then rebuild the image")
		os.Exit(2)
	}

	code := m.Run()
	DeleteData()
	os.Exit(code)
}

func InitializeDb() *sql.DB {
	driver := os.Getenv("DB_DRIVER")
	username := os.Getenv("DB_USER_NAME")
	password := os.Getenv("DB_PASSWORD")
	address := os.Getenv("DB_ADDRESS_TEST")
	database := os.Getenv("DB_DATABASE_TEST")
	dbase, err := sql.Open(driver, username+":"+password+"@tcp("+address+")"+"/"+database+"?parseTime=true")
	if err != nil {
		panic(err.Error())
	}

	return dbase
}

func TestGetNonExistentMovies(t *testing.T) {

	DeleteData()

	apiCall, err := http.NewRequest("GET", "http://localhost:8080/Movies", nil)
	if err != nil {
		panic(err.Error())
	}

	c := &http.Client{}

	r, err := c.Do(apiCall)
	if err != nil {
		panic(err.Error())
	}

	assert.Equal(t, 404, r.StatusCode)

}

func TestGetMovies(t *testing.T) {

	DeleteData()

	var requestData []request.Movie

	var responseData []response.Movie

	mockData := request.Movie{
		Id:          1,
		Title:       "Pengabdi Setan 2 Comunion",
		Description: "dalah sebuah film horor Indonesia tahun 2022 yang disutradarai dan ditulis oleh Joko Anwar sebagai sekuel dari film tahun 2017, Pengabdi Setan.",
		Rating:      7,
		Image:       "",
		Created_at:  "2022-08-01 10:56:31",
		Updated_at:  "2022-08-13 09:30:23",
	}

	requestData = append(requestData, mockData)

	InsertData(requestData)

	apiCall, err := http.NewRequest("GET", "http://localhost:8080/Movies", nil)
	if err != nil {
		panic(err.Error())
	}

	c := &http.Client{}

	r, err := c.Do(apiCall)
	if err != nil {
		panic(err.Error())
	}

	defer r.Body.Close()

	bodyResponse, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	json.Unmarshal(bodyResponse, &responseData)

	assert.Equal(t, 200, r.StatusCode)
	assert.NotEmpty(t, r.Body)

	for i := 0; i < len(responseData); i++ {
		assert.Equal(t, responseData[i].Id, requestData[i].Id)
		assert.Equal(t, responseData[i].Title, requestData[i].Title)
		assert.Equal(t, responseData[i].Description, requestData[i].Description)
		assert.Equal(t, responseData[i].Rating, requestData[i].Rating)
		assert.Equal(t, responseData[i].Image, requestData[i].Image)
		assert.Equal(t, responseData[i].Created_at, requestData[i].Created_at)
		assert.Equal(t, responseData[i].Updated_at, requestData[i].Updated_at)
	}

}

func TestGetNonExistentMovie(t *testing.T) {

	DeleteData()

	apiCall, err := http.NewRequest("GET", "http://localhost:8080/Movies/1", nil)
	if err != nil {
		panic(err.Error())
	}

	c := &http.Client{}

	r, err := c.Do(apiCall)
	if err != nil {
		panic(err.Error())
	}

	assert.Equal(t, 404, r.StatusCode)

}

func TestGetMovie(t *testing.T) {

	DeleteData()

	var requestData []request.Movie

	var response response.Movie

	mockData := request.Movie{
		Id:          1,
		Title:       "Pengabdi Setan 2 Comunion",
		Description: "dalah sebuah film horor Indonesia tahun 2022 yang disutradarai dan ditulis oleh Joko Anwar sebagai sekuel dari film tahun 2017, Pengabdi Setan.",
		Rating:      7,
		Image:       "",
		Created_at:  "2022-08-01 10:56:31",
		Updated_at:  "2022-08-13 09:30:23",
	}

	requestData = append(requestData, mockData)

	InsertData(requestData)

	apiCall, err := http.NewRequest("GET", "http://localhost:8080/Movies/1", nil)
	if err != nil {
		panic(err.Error())
	}

	c := &http.Client{}

	r, err := c.Do(apiCall)
	if err != nil {
		panic(err.Error())
	}

	defer r.Body.Close()

	bodyResponse, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	json.Unmarshal(bodyResponse, &response)

	fmt.Println(response.Id)

	assert.Equal(t, 200, r.StatusCode)
	assert.NotEmpty(t, r.Body)

	assert.Equal(t, mockData.Id, response.Id)
	assert.Equal(t, mockData.Title, response.Title)
	assert.Equal(t, mockData.Description, response.Description)
	assert.Equal(t, mockData.Rating, response.Rating)
	assert.Equal(t, mockData.Image, response.Image)
	assert.Equal(t, mockData.Created_at, response.Created_at)
	assert.Equal(t, mockData.Updated_at, response.Updated_at)

}

func TestInsertMovie(t *testing.T) {

	DeleteData()

	mockData := request.Movie{
		Id:          1,
		Title:       "Pengabdi Setan 2 Comunion",
		Description: "dalah sebuah film horor Indonesia tahun 2022 yang disutradarai dan ditulis oleh Joko Anwar sebagai sekuel dari film tahun 2017, Pengabdi Setan.",
		Rating:      7,
		Image:       "",
		Created_at:  "2022-08-01 10:56:31",
		Updated_at:  "2022-08-13 09:30:23",
	}

	data, err := json.Marshal(mockData)
	if err != nil {
		panic(err.Error())
	}
	reader := bytes.NewReader(data)

	apiCall, err := http.NewRequest("POST", "http://localhost:8080/Movies", reader)
	if err != nil {
		panic(err.Error())
	}

	c := &http.Client{}

	r, err := c.Do(apiCall)
	if err != nil {
		panic(err.Error())
	}

	assert.Equal(t, 200, r.StatusCode)

}

func TestInsertMovieFailValidation(t *testing.T) {

	DeleteData()

	var errorResponse response.ErrorResponse

	mockData := request.Movie{
		Id:          1,
		Title:       "",
		Description: "",
		Rating:      0,
		Image:       "",
		Created_at:  "",
		Updated_at:  "",
	}

	data, err := json.Marshal(mockData)
	if err != nil {
		panic(err.Error())
	}
	reader := bytes.NewReader(data)

	apiCall, err := http.NewRequest("POST", "http://localhost:8080/Movies", reader)
	if err != nil {
		panic(err.Error())
	}

	c := &http.Client{}

	r, err := c.Do(apiCall)
	if err != nil {
		panic(err.Error())
	}

	defer r.Body.Close()

	bodyResponse, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	json.Unmarshal(bodyResponse, &errorResponse)

	assert.Equal(t, 400, r.StatusCode)
	assert.NotEmpty(t, errorResponse.Message)
}

func TestUpdateMovie(t *testing.T) {

	DeleteData()

	var requestData []request.Movie

	mockData := request.Movie{
		Id:          1,
		Title:       "Pengabdi Setan 2 Comunion",
		Description: "dalah sebuah film horor Indonesia tahun 2022 yang disutradarai dan ditulis oleh Joko Anwar sebagai sekuel dari film tahun 2017, Pengabdi Setan.",
		Rating:      7,
		Image:       "",
		Created_at:  "2022-08-01 10:56:31",
		Updated_at:  "2022-08-13 09:30:23",
	}

	requestData = append(requestData, mockData)

	InsertData(requestData)

	data, err := json.Marshal(mockData)
	if err != nil {
		panic(err.Error())
	}
	reader := bytes.NewReader(data)

	apiCall, err := http.NewRequest("PATCH", "http://localhost:8080/Movies/1", reader)
	if err != nil {
		panic(err.Error())
	}

	c := &http.Client{}

	r, err := c.Do(apiCall)
	if err != nil {
		panic(err.Error())
	}

	assert.Equal(t, 200, r.StatusCode)

}

func TestUpdateMovieFailValidation(t *testing.T) {

	DeleteData()

	var errorResponse response.ErrorResponse

	mockData := request.Movie{
		Id:          1,
		Title:       "",
		Description: "",
		Rating:      0,
		Image:       "",
		Created_at:  "",
		Updated_at:  "",
	}

	data, err := json.Marshal(mockData)
	if err != nil {
		panic(err.Error())
	}
	reader := bytes.NewReader(data)

	apiCall, err := http.NewRequest("PATCH", "http://localhost:8080/Movies/1", reader)
	if err != nil {
		panic(err.Error())
	}

	c := &http.Client{}

	r, err := c.Do(apiCall)
	if err != nil {
		panic(err.Error())
	}

	defer r.Body.Close()

	bodyResponse, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	json.Unmarshal(bodyResponse, &errorResponse)

	assert.Equal(t, 400, r.StatusCode)
	assert.NotEmpty(t, errorResponse.Message)
}

func TestDeleteMovie(t *testing.T) {

	DeleteData()

	var requestData []request.Movie

	mockData := request.Movie{
		Id:          1,
		Title:       "Pengabdi Setan 2 Comunion",
		Description: "dalah sebuah film horor Indonesia tahun 2022 yang disutradarai dan ditulis oleh Joko Anwar sebagai sekuel dari film tahun 2017, Pengabdi Setan.",
		Rating:      7,
		Image:       "",
		Created_at:  "2022-08-01 10:56:31",
		Updated_at:  "2022-08-13 09:30:23",
	}

	requestData = append(requestData, mockData)

	InsertData(requestData)

	apiCall, err := http.NewRequest("DELETE", "http://localhost:8080/Movies/1", nil)
	if err != nil {
		panic(err.Error())
	}

	c := &http.Client{}

	r, err := c.Do(apiCall)
	if err != nil {
		panic(err.Error())
	}

	assert.Equal(t, 200, r.StatusCode)

}

func DeleteData() {
	db := InitializeDb()
	stmt, err := db.Prepare("DELETE FROM movies;")
	if err != nil {
		panic(err.Error())
	}

	_, err = stmt.Exec()
	if err != nil {
		panic(err.Error())
	}

	defer stmt.Close()

	fmt.Println("Data deleted")
}

func InsertData(r []request.Movie) {

	db := InitializeDb()

	stmt, err := db.Prepare("INSERT INTO movies(id, title, description, rating, image, created_at, updated_at) VALUES (?,?,?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}

	for i := 0; i < len(r); i++ {
		_, err = stmt.Exec(r[i].Id, r[i].Title, r[i].Description, r[i].Rating, r[i].Image, r[i].Created_at, r[i].Updated_at)
		if err != nil {
			panic(err.Error())
		}
	}

	defer stmt.Close()

	fmt.Println("Data inserted")
}

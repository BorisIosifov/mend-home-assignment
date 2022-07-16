package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/BorisIosifov/mend-home-assignment/storage"
)

func TestServeHTTP(t *testing.T) {
	var (
		params           map[string]string
		expectedResponse string
	)
	strg, err := storage.PrepareLocalMemory()
	if err != nil {
		t.Fatal(err)
	}
	api := API{
		Storage: strg,
	}

	// Testing of getting an empty list of objects
	checkRequest(t, api, "GET", "/books/", nil, 200, `[]`)

	// Testing of getting an empty list of objects
	checkRequest(t, api, "GET", "/cars/", nil, 200, `[]`)

	// Testing of adding an object
	params = map[string]string{
		"author": "Ernest Hamingway",
		"title":  "For Whom the Bell Tolls",
	}
	expectedResponse = `{"ID":1,"Author":"Ernest Hamingway","Title":"For Whom the Bell Tolls"}`
	checkRequest(t, api, "POST", "/books/", params, 200, expectedResponse)

	// Testing of adding an object
	params = map[string]string{
		"author": "Fedor Dostoevsky",
		"title":  "Crime and Punishment",
	}
	expectedResponse = `{"ID":2,"Author":"Fedor Dostoevsky","Title":"Crime and Punishment"}`
	checkRequest(t, api, "POST", "/books/", params, 200, expectedResponse)

	// Testing of getting a list of objects
	expectedResponse = `[{"ID":1,"Author":"Ernest Hamingway","Title":"For Whom the Bell Tolls"},{"ID":2,"Author":"Fedor Dostoevsky","Title":"Crime and Punishment"}]`
	checkRequest(t, api, "GET", "/books/", nil, 200, expectedResponse)

	// Testing of getting an object
	expectedResponse = `{"ID":1,"Author":"Ernest Hamingway","Title":"For Whom the Bell Tolls"}`
	checkRequest(t, api, "GET", "/books/1/", nil, 200, expectedResponse)

	// Testing of getting a non-existent object
	expectedResponse = `{"error":"Object books with id 12 not found"}`
	checkRequest(t, api, "GET", "/books/12/", nil, 404, expectedResponse)

	// Testing of updating an object
	params = map[string]string{
		"author": "Fedor Dostoevsky",
		"title":  "Crime and Punishment!!!",
	}
	expectedResponse = `{"ID":2,"Author":"Fedor Dostoevsky","Title":"Crime and Punishment!!!"}`
	checkRequest(t, api, "PUT", "/books/2/", params, 200, expectedResponse)

	// Testing of updating a non-existent object
	params = map[string]string{
		"author": "Sholem Aleichem",
		"title":  "Wandering Stars",
	}
	expectedResponse = `{"error":"Object books with id 123 not found"}`
	checkRequest(t, api, "PUT", "/books/123/", params, 404, expectedResponse)

	// Testing of getting a list of objects after updating one of them
	expectedResponse = `[{"ID":1,"Author":"Ernest Hamingway","Title":"For Whom the Bell Tolls"},{"ID":2,"Author":"Fedor Dostoevsky","Title":"Crime and Punishment!!!"}]`
	checkRequest(t, api, "GET", "/books/", nil, 200, expectedResponse)

	// Testing of deleting an object
	checkRequest(t, api, "DELETE", "/books/2/", params, 200, `{}`)

	// Testing of deleting a non-existent object
	expectedResponse = `{"error":"Object books with id 1234 not found"}`
	checkRequest(t, api, "DELETE", "/books/1234/", params, 404, expectedResponse)

	// Testing of getting a deleted object
	expectedResponse = `{"error":"Object books with id 2 not found"}`
	checkRequest(t, api, "GET", "/books/2/", nil, 404, expectedResponse)

	// Testing of getting a list of objects after deleting one of them
	expectedResponse = `[{"ID":1,"Author":"Ernest Hamingway","Title":"For Whom the Bell Tolls"}]`
	checkRequest(t, api, "GET", "/books/", nil, 200, expectedResponse)

	// Testing of getting a list of non-existent objects
	expectedResponse = `{"error":"Page not found"}`
	checkRequest(t, api, "GET", "/goods/", nil, 404, expectedResponse)

	// Testing of a bad request
	expectedResponse = `{"error":"Bad request"}`
	checkRequest(t, api, "GET", "/books/123/456/", nil, 400, expectedResponse)
}

func checkRequest(t *testing.T, api API, method string, path string, params map[string]string, expectedStatus int, expectedResponse string) {
	var reqBody io.Reader
	v := url.Values{}
	if len(params) > 0 {
		for param, value := range params {
			v.Set(param, value)
		}
		reqBody = strings.NewReader(v.Encode())
	}

	req, err := http.NewRequest(method, path, reqBody)
	if err != nil {
		t.Fatal(err)
	}
	if reqBody != nil {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	}

	w := httptest.NewRecorder()
	api.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	assert.Equal(t, expectedStatus, resp.StatusCode)
	assert.JSONEq(t, expectedResponse, string(body))
}

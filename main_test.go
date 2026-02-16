package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeCount(t *testing.T) {
	var actualCount int
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		count int
		wait  int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, len(cafeList["tula"])},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/cafe?city=tula&count=%d", v.count), nil)
		handler.ServeHTTP(response, req)
		responseString := strings.TrimSpace(response.Body.String())
		if responseString == "" {
			actualCount = 0
		} else {
			responseArray := strings.Split(responseString, ",")
			actualCount = len(responseArray)
		}
		require.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, v.wait, actualCount)
	}
}
func TestCafeSearch(t *testing.T) {
	var actualCount int
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		search    string
		wantCount int
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", fmt.Sprintf("/cafe?city=moscow&search=%s", v.search), nil)
		handler.ServeHTTP(response, req)
		responseString := strings.TrimSpace(response.Body.String())
		if responseString == "" {
			actualCount = 0
		} else {
			responseArray := strings.Split(responseString, ",")
			result := []string{}
			for _, j := range responseArray {
				if strings.Contains(strings.ToLower(j), v.search) {
					result = append(result, j)
				}
			}
			actualCount = len(result)
		}
		require.Equal(t, http.StatusOK, response.Code)
		assert.Equal(t, v.wantCount, actualCount)
	}
}
func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

func TestConvert(t *testing.T) {
	t.Log("hell")
}

func TestRedicCreation(t *testing.T) {
	client := createRedisClient()
	resp, err := client.Ping().Result()
	if err != nil {
		t.Fail()
	}
	if resp != "PONG" {
		t.Fail()
	}
}

func TestHomepage(t *testing.T) {
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()
	router := Router()
	router.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")
	resp, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fail()
	}
	assert.Equalf(t, "welcome to my program", string(resp), "it should return a welcome")

}

func TestAddTask(t *testing.T) {
	request, _ := http.NewRequest("GET", "/addTask?Command=ls&Argument=/tmp", nil)
	response := httptest.NewRecorder()
	router := Router()
	router.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")
}

func TestExecuteTask(t *testing.T) {
	s := miniredis.RunT(t)
	s.Set("tasks", "{\"Command\":\"ls\",\"Argument\":\"/tmp\"}")
	request, _ := http.NewRequest("GET", "/executeTask", nil)
	response := httptest.NewRecorder()
	router := Router()
	router.ServeHTTP(response, request)
	assert.Equal(t, 200, response.Code, "OK response is expected")
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fail()
	}
	assert.Contains(t, string(body), "tmux-1000")

}

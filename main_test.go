package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-redis/redismock/v8"
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
	client, mock := redismock.NewClientMock()

	// TODO 5 - Stub redis responses

	// TODO 6 - Stub executeTask responses

	// TODO 7 - Run test

	// TODO 8 - Compare test response with expected stubbed responses
}

package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainHandlerCorrectRequest(t *testing.T) {
	req, err := http.NewRequest("GET", "/cafe?city=moscow&count=2", nil)
	require.NoError(t, err, "Создание запроса не должно приводить к ошибке")

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusOK, responseRecorder.Code, "Код ответа должен быть 200")
	body := responseRecorder.Body.String()
	assert.NotEmpty(t, body, "Тело ответа не должно быть пустым")

	expectedCafe := "Мир кофе,Сладкоежка"
	assert.Equal(t, expectedCafe, body, "Ответ должен содержать корректный список кафе")
}

func TestMainHandlerUnsupportedCity(t *testing.T) {
	req, err := http.NewRequest("GET", "/cafe?city=paris&count=2", nil)
	require.NoError(t, err, "Создание запроса не должно приводить к ошибке")

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code, "Код ответа должен быть 400")
	body := responseRecorder.Body.String()
	assert.Equal(t, "wrong city value", body, "Ответ должен содержать ошибку 'wrong city value'")
}

func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4 // Количество доступных кафе в Москве

	req, err := http.NewRequest("GET", "/cafe?city=moscow&count=10", nil)
	require.NoError(t, err, "Создание запроса не должно приводить к ошибке")

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusOK, responseRecorder.Code, "Код ответа должен быть 200")
	body := responseRecorder.Body.String()
	assert.NotEmpty(t, body, "Тело ответа не должно быть пустым")

	actualCafeList := strings.Split(body, ",")
	assert.Len(t, actualCafeList, totalCount, "Должны вернуться все доступные кафе")
}

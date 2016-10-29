package handler

import (
	"net/http/httptest"

	"net/http"

	"bytes"
	"encoding/json"

	"github.com/mrkaspa/geoserver/models"
	"os"
	"github.com/mrkaspa/geoserver/utils"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	utils.LoadEnv("../.env_test")
	utils.InitLogger()
	retCode := m.Run()
	os.Exit(retCode)
}

func TestController_storeHandler(t *testing.T) {
	testController := controller{persistorCreator: models.NewMockPersistor}
	stroke := models.Stroke{
		UserID:   "a1",
		Info:     "a1",
		Location: []float64{-79.38066843, 43.65483486},
	}

	data, _ := json.Marshal(stroke)
	req, _ := http.NewRequest(http.MethodPost, "", bytes.NewReader(data))
	w := httptest.NewRecorder()

	testController.storeHandler(w, req)
	assert.Equal(t, w.Code, http.StatusOK)
}

func TestController_nearHandler(t *testing.T) {
	testController := controller{persistorCreator: models.NewMockPersistor}
	stroke := models.Stroke{
		UserID:   "a1",
		Info:     "a1",
		Location: []float64{-79.38066843, 43.65483486},
	}

	data, _ := json.Marshal(stroke)
	req, _ := http.NewRequest(http.MethodPost, "", bytes.NewReader(data))
	w := httptest.NewRecorder()

	testController.nearHandler(w, req)
	assert.Equal(t, w.Code, http.StatusOK)
	var matches []models.Stroke
	json.Unmarshal(w.Body.Bytes(), &matches)
	assert.NotEmpty(t, matches)
}

func TestController_recentHandler(t *testing.T) {
	testController := controller{persistorCreator: models.NewMockPersistor}
	req, _ := http.NewRequest(http.MethodGet, "/recent/a1", nil)
	w := httptest.NewRecorder()
	testController.recentStrokes(w, req)
	assert.Equal(t, w.Code, http.StatusOK)
	var matches []models.Stroke
	json.Unmarshal(w.Body.Bytes(), &matches)
	assert.NotEmpty(t, matches)
}

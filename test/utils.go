package test

import (
	"encoding/json"

	"github.com/mrkaspa/geoserver/models"
)

func createStrokeNear(userID string, loc []float64) (*models.StrokeNear, []byte) {
	stroke := models.Stroke{
		UserID:   userID,
		Info:     userID,
		Location: loc,
	}
	strokeNear := models.StrokeNear{
		TimeRange:   3,
		MaxDistance: 5,
		Stroke:      stroke,
	}
	json, _ := json.Marshal(&strokeNear)
	return &strokeNear, json
}

func createStroke(userID string, loc []float64) (*models.Stroke, []byte) {
	stroke := models.Stroke{
		UserID:   userID,
		Info:     userID,
		Location: loc,
	}
	json, _ := json.Marshal(&stroke)
	return &stroke, json
}

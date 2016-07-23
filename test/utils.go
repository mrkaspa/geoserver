package test

import (
	"encoding/json"

	"github.com/mrkaspa/geoserver/models"
	. "github.com/onsi/gomega"
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

func BeIn(arr []interface{}, val interface{}) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

func matchTwo(err1, err2 error, resp1, resp2, infoA, infoB string) {
	posibilities := []interface{}{infoA, infoB}
	Expect(err1).To(BeNil())
	Expect(err2).To(BeNil())
	Expect(resp1).NotTo(BeEquivalentTo(resp2))
	Expect(BeIn(posibilities, resp1)).To(BeTrue())
	Expect(BeIn(posibilities, resp2)).To(BeTrue())
}

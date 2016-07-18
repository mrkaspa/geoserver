package test

import (
	"bytes"
	"net/http"

	"encoding/json"

	"github.com/mrkaspa/geoserver/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("/near", func() {

	It("should get OK", func() {
		_, data1 := createStroke("a1", []float64{-79.38066843, 43.65483486})
		http.Post(ts.URL+"/store", "application/json; charset=utf-8", bytes.NewBuffer(data1))

		_, data2 := createStrokeNear("a2", []float64{-79.38066843, 43.65483486})
		res, _ :=
			http.Post(ts.URL+"/near", "application/json; charset=utf-8", bytes.NewBuffer(data2))
		Expect(res.StatusCode).To(BeEquivalentTo(http.StatusOK))

		var strokes []models.Stroke
		json.NewDecoder(res.Body).Decode(&strokes)
		Expect(len(strokes)).To(BeEquivalentTo(1))
		Expect(strokes[0].UserID).To(BeEquivalentTo("a1"))
	})

})

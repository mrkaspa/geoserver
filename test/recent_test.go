package test

import (
	"net/http"

	"encoding/json"

	"bytes"

	"github.com/mrkaspa/geoserver/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("/recent", func() {

	It("should get OK", func() {
		_, data := createStroke("a1", []float64{-79.38066843, 43.65483486})
		http.Post(ts.URL+"/store", "application/json; charset=utf-8", bytes.NewBuffer(data))
		res, _ := http.Get(ts.URL + "/recent/a1")
		createStroke("a1", []float64{-79.38066843, 43.65483486})

		Expect(res.StatusCode).To(Equal(http.StatusOK))

		var strokes []models.Stroke
		json.NewDecoder(res.Body).Decode(&strokes)
		Expect(len(strokes)).To(Equal(1))
		Expect(strokes[0].UserID).To(Equal("a1"))
	})

})

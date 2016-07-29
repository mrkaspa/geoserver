package handler

import (
	"net/http/httptest"

	"net/http"

	"bytes"
	"encoding/json"

	"github.com/mrkaspa/geoserver/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("controller", func() {

	testController := controller{persistorCreator: newMockPersistor}

	It("tests storeHandler", func() {
		stroke := models.Stroke{
			UserID:   "a1",
			Info:     "a1",
			Location: []float64{-79.38066843, 43.65483486},
		}

		data, _ := json.Marshal(stroke)
		req, _ := http.NewRequest(http.MethodPost, "", bytes.NewReader(data))
		w := httptest.NewRecorder()

		testController.storeHandler(w, req)
		Expect(w.Code).To(Equal(http.StatusOK))
	})

	It("tests nearHandler", func() {
		stroke := models.Stroke{
			UserID:   "a1",
			Info:     "a1",
			Location: []float64{-79.38066843, 43.65483486},
		}

		data, _ := json.Marshal(stroke)
		req, _ := http.NewRequest(http.MethodPost, "", bytes.NewReader(data))
		w := httptest.NewRecorder()

		testController.nearHandler(w, req)
		Expect(w.Code).To(Equal(http.StatusOK))
		var matches []models.Stroke
		json.Unmarshal(w.Body.Bytes(), &matches)
		Expect(matches).ToNot(BeEmpty())
	})

})

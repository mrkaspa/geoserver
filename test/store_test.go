package test

import (
	"bytes"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("/store", func() {

	It("should get OK", func() {
		_, data := createStroke("a1", []float64{-79.38066843, 43.65483486})
		res, _ :=
			http.Post(ts.URL+"/store", "application/json; charset=utf-8", bytes.NewBuffer(data))
		Expect(res.StatusCode).To(Equal(http.StatusOK))
	})

})

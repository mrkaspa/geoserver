package handler

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"

	"github.com/mrkaspa/geoserver/utils"
)

func TestHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Handler Suite")
}

var _ = BeforeSuite(func() {
	utils.LoadEnv("../.env_test")
	utils.InitLogger()
})

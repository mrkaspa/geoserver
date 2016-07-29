package models

import (
	"testing"

	"github.com/mrkaspa/geoserver/utils"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestPersistor(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Models Suite")
}

var _ = BeforeSuite(func() {
	utils.LoadEnv("../.env_test")
	utils.InitLogger()
	InitDB()
})

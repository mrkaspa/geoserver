package test

import (
	"net/http/httptest"
	"testing"

	"gopkg.in/mgo.v2/bson"

	"github.com/mrkaspa/geoserver/handler"
	"github.com/mrkaspa/geoserver/models"
	"github.com/mrkaspa/geoserver/utils"

	"time"

	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	ts *httptest.Server
)

func TestBlackBox(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Api WS")
}

var _ = BeforeSuite(func() {
	utils.LoadEnv("../.env_test")
	utils.InitLogger()
	models.InitDB()
	handler.InitSearcher(models.NewPersistor)
	ts = httptest.NewServer(handler.NewRouter())
})

var _ = AfterSuite(func() {
	ts.Close()
	models.Session.Close()
	time.Sleep(1 * time.Second)
})

var _ = BeforeEach(func() {
	fmt.Println("----------------------------------------------------------------")
	cleanDB()
})

var _ = AfterEach(func() {
	fmt.Println("****************************************************************")
})

func cleanDB() {
	models.StrokesCollection.RemoveAll(bson.M{})
}

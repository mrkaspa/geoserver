package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gopkg.in/mgo.v2/bson"

	"github.com/mrkaspa/geoserver/handler"
	"github.com/mrkaspa/geoserver/models"
	"github.com/mrkaspa/geoserver/utils"

	"time"

	"fmt"

	"github.com/joho/godotenv"
	"github.com/mrkaspa/geoserver/ws"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	ts *httptest.Server
)

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Api WS")
}

var _ = BeforeSuite(func() {
	router := handler.NewRouter()
	http.Handle("/", router)
	ts = httptest.NewServer(router)
})

var _ = AfterSuite(func() {
	ts.Close()
	models.Session.Close()
	time.Sleep(1 * time.Second)
})

var _ = BeforeEach(func() {
	fmt.Println("----------------------------------------------------------------")
	cleanDB()
	wsConnUser1 = createClient(ts.URL, username1)
	wsConnUser2 = createClient(ts.URL, username2)
	wsConnUser3 = createClient(ts.URL, username3)
})

var _ = AfterEach(func() {
	if wsConnUser1.IsServerConn() {
		wsConnUser1.Close()
	}
	if wsConnUser2.IsServerConn() {
		wsConnUser2.Close()
	}
	if wsConnUser3.IsServerConn() {
		wsConnUser3.Close()
	}
	ws.SearcherVar.Clean <- true
	fmt.Println("****************************************************************")
})

func cleanDB() {
	models.StrokesCollection.RemoveAll(bson.M{})
}

func init() {
	path := ".env_test"
	for i := 1; ; i++ {
		if err := godotenv.Load(path); err != nil {
			if i > 3 {
				panic("Error loading .env_test file")
			} else {
				path = "../" + path
			}
		} else {
			break
		}
	}
	utils.Init()
	models.Init()
}

package test

import (
	"time"

	"sync"

	"github.com/mrkaspa/geoserver/models"
	"github.com/mrkaspa/geoserver/utils"
	"github.com/mrkaspa/geoserver/ws"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/websocket"
)

var (
	username1           = "a1"
	wsConnUser1         *websocket.Conn
	postStrokeUser1Byte []byte
	postStrokeUser1     *models.StrokeNear
	username2           = "a2"
	wsConnUser2         *websocket.Conn
	postStrokeUser2Byte []byte
	postStrokeUser2     *models.StrokeNear
	username3           = "a3"
	wsConnUser3         *websocket.Conn
	postStrokeUser3Byte []byte
	postStrokeUser3     *models.StrokeNear
	resp1               string
	resp2               string
	resp11              string
	resp12              string
	resp21              string
	resp22              string
	resp31              string
	resp32              string
	err1                error
	err2                error
	err11               error
	err12               error
	err21               error
	err22               error
	err31               error
	err32               error
)

func withTwoUsers(sleepTime int) {
	websocket.Message.Send(wsConnUser1, postStrokeUser1Byte)
	if sleepTime > 0 {
		time.Sleep(time.Duration(sleepTime) * time.Second)
	}
	websocket.Message.Send(wsConnUser2, postStrokeUser2Byte)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		err1 = websocket.Message.Receive(wsConnUser1, &resp1)
	}()
	go func() {
		defer wg.Done()
		err2 = websocket.Message.Receive(wsConnUser2, &resp2)
	}()
	wg.Wait()
}

func withThreeUsers(sleepTime1, sleepTime2 int) {
	websocket.Message.Send(wsConnUser1, postStrokeUser1Byte)
	if sleepTime1 > 0 {
		time.Sleep(time.Duration(sleepTime1) * time.Second)
	}
	websocket.Message.Send(wsConnUser2, postStrokeUser2Byte)
	if sleepTime2 > 0 {
		time.Sleep(time.Duration(sleepTime2) * time.Second)
	}
	websocket.Message.Send(wsConnUser3, postStrokeUser3Byte)
	wg := sync.WaitGroup{}
	wg.Add(6)
	go func() {
		defer wg.Done()
		err11 = websocket.Message.Receive(wsConnUser1, &resp11)
	}()
	go func() {
		defer wg.Done()
		err12 = websocket.Message.Receive(wsConnUser1, &resp12)
	}()
	go func() {
		defer wg.Done()
		err21 = websocket.Message.Receive(wsConnUser2, &resp21)
	}()
	go func() {
		defer wg.Done()
		err22 = websocket.Message.Receive(wsConnUser2, &resp22)
	}()
	go func() {
		defer wg.Done()
		err31 = websocket.Message.Receive(wsConnUser3, &resp31)
	}()
	go func() {
		defer wg.Done()
		err32 = websocket.Message.Receive(wsConnUser3, &resp32)
	}()
	wg.Wait()
}

var _ = Describe("WS Behavior", func() {

	BeforeEach(func() {
		wsConnUser1 = createClient(ts.URL, username1)
		wsConnUser2 = createClient(ts.URL, username2)
		wsConnUser3 = createClient(ts.URL, username3)
	})

	AfterEach(func() {
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
	})

	Context("with two users", func() {

		BeforeEach(func() {
			postStrokeUser1, postStrokeUser1Byte = createStrokeNear(username1, []float64{-79.38066843, 43.65483486})
			postStrokeUser2, postStrokeUser2Byte = createStrokeNear(username2, []float64{-79.38066843, 43.65483486})
		})

		It("should do match", func() {
			withTwoUsers(0)
			Expect(err1).To(BeNil())
			Expect(resp1).To(BeEquivalentTo(postStrokeUser2.Stroke.Info))
			Expect(err2).To(BeNil())
			Expect(resp2).To(BeEquivalentTo(postStrokeUser1.Stroke.Info))
		})

		It("should do match after 1 second", func() {
			withTwoUsers(1)
			Expect(err1).To(BeNil())
			Expect(resp1).To(BeEquivalentTo(postStrokeUser2.Stroke.Info))
			Expect(err2).To(BeNil())
			Expect(resp2).To(BeEquivalentTo(postStrokeUser1.Stroke.Info))
		})

		It("should do not match", func() {
			withTwoUsers(4)
			Expect(err1).NotTo(BeNil())
			Expect(err2).NotTo(BeNil())
		})

	})

	Context("with two users far away", func() {

		BeforeEach(func() {
			postStrokeUser1, postStrokeUser1Byte = createStrokeNear(username1, []float64{-79.38066843, 43.65483486})
			postStrokeUser2, postStrokeUser2Byte = createStrokeNear(username2, []float64{-49.38066843, 43.65483486})
		})

		It("should do not match", func() {
			withTwoUsers(0)
			Expect(err1).NotTo(BeNil())
			Expect(err2).NotTo(BeNil())
		})

	})

	Context("with three users", func() {

		BeforeEach(func() {
			postStrokeUser1, postStrokeUser1Byte = createStrokeNear(username1, []float64{-79.38066843, 43.65483486})
			postStrokeUser2, postStrokeUser2Byte = createStrokeNear(username2, []float64{-79.38066843, 43.65483486})
			postStrokeUser3, postStrokeUser3Byte = createStrokeNear(username3, []float64{-79.38066843, 43.65483486})
		})

		It("should do match", func() {
			withThreeUsers(0, 0)
			utils.Log.Infof("matchOtherTwo a1")
			matchTwo(err11, err12, resp11, resp12, postStrokeUser2.Stroke.Info, postStrokeUser3.Stroke.Info)
			utils.Log.Infof("matchOtherTwo a2")
			matchTwo(err21, err22, resp21, resp22, postStrokeUser1.Stroke.Info, postStrokeUser3.Stroke.Info)
			utils.Log.Infof("matchOtherTwo a3")
			matchTwo(err31, err32, resp31, resp32, postStrokeUser1.Stroke.Info, postStrokeUser2.Stroke.Info)
		})

		It("should do match after 1 - 2 seconds", func() {
			withThreeUsers(2, 1)
			utils.Log.Infof("matchOtherTwo a1")
			matchTwo(err11, err12, resp11, resp12, postStrokeUser2.Stroke.Info, postStrokeUser3.Stroke.Info)
			utils.Log.Infof("matchOtherTwo a2")
			matchTwo(err21, err22, resp21, resp22, postStrokeUser1.Stroke.Info, postStrokeUser3.Stroke.Info)
			utils.Log.Infof("matchOtherTwo a3")
			matchTwo(err31, err32, resp31, resp32, postStrokeUser1.Stroke.Info, postStrokeUser2.Stroke.Info)
		})

		It("a1 and a2 should match, a3 shouldn't match", func() {
			withThreeUsers(0, 4)
			Expect((err11 == nil && err12 != nil) || (err11 != nil && err12 == nil)).To(BeTrue())
			Expect((err21 == nil && err22 != nil) || (err21 != nil && err22 == nil)).To(BeTrue())
			Expect(err31).NotTo(BeNil())
			Expect(err32).NotTo(BeNil())
		})

	})

})

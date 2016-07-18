package test

import (
	"time"

	"github.com/mrkaspa/geoserver/utils"
	"github.com/mrkaspa/geoserver/ws"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/net/websocket"
)

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
			postStrokeUser1, postStrokeUser1Byte = createStroke(username1, []float64{-79.38066843, 43.65483486})
			postStrokeUser2, postStrokeUser2Byte = createStroke(username2, []float64{-79.38066843, 43.65483486})
		})

		It("should do match", func() {
			websocket.Message.Send(wsConnUser1, postStrokeUser1Byte)
			websocket.Message.Send(wsConnUser2, postStrokeUser2Byte)
			resp1 := new(string)
			resp2 := new(string)
			err1 := websocket.Message.Receive(wsConnUser1, resp1)
			err2 := websocket.Message.Receive(wsConnUser2, resp2)
			Expect(err1).To(BeNil())
			Expect(*resp1).To(BeEquivalentTo(postStrokeUser2.Info))
			Expect(err2).To(BeNil())
			Expect(*resp2).To(BeEquivalentTo(postStrokeUser1.Info))
		})

		It("should do match after 1 second", func() {
			websocket.Message.Send(wsConnUser1, postStrokeUser1Byte)
			time.Sleep(1 * time.Second)
			websocket.Message.Send(wsConnUser2, postStrokeUser2Byte)
			resp1 := new(string)
			resp2 := new(string)
			err1 := websocket.Message.Receive(wsConnUser1, resp1)
			err2 := websocket.Message.Receive(wsConnUser2, resp2)
			Expect(err1).To(BeNil())
			Expect(*resp1).To(BeEquivalentTo(postStrokeUser2.Info))
			Expect(err2).To(BeNil())
			Expect(*resp2).To(BeEquivalentTo(postStrokeUser1.Info))
		})

		It("should do not match", func() {
			websocket.Message.Send(wsConnUser1, postStrokeUser1Byte)
			time.Sleep(3 * time.Second)
			websocket.Message.Send(wsConnUser2, postStrokeUser2Byte)
			resp1 := new(string)
			resp2 := new(string)
			err1 := websocket.Message.Receive(wsConnUser1, resp1)
			err2 := websocket.Message.Receive(wsConnUser2, resp2)
			Expect(err1).NotTo(BeNil())
			Expect(err2).NotTo(BeNil())
		})

	})

	Context("with two users far away", func() {

		BeforeEach(func() {
			postStrokeUser1, postStrokeUser1Byte = createStroke(username1, []float64{-79.38066843, 43.65483486})
			postStrokeUser2, postStrokeUser2Byte = createStroke(username2, []float64{-49.38066843, 43.65483486})
		})

		It("should do not match", func() {
			websocket.Message.Send(wsConnUser1, postStrokeUser1Byte)
			time.Sleep(3 * time.Second)
			websocket.Message.Send(wsConnUser2, postStrokeUser2Byte)
			resp1 := new(string)
			resp2 := new(string)
			err1 := websocket.Message.Receive(wsConnUser1, resp1)
			err2 := websocket.Message.Receive(wsConnUser2, resp2)
			Expect(err1).NotTo(BeNil())
			Expect(err2).NotTo(BeNil())
		})

	})

	Context("with three users", func() {

		BeforeEach(func() {
			postStrokeUser1, postStrokeUser1Byte = createStroke(username1, []float64{-79.38066843, 43.65483486})
			postStrokeUser2, postStrokeUser2Byte = createStroke(username2, []float64{-79.38066843, 43.65483486})
			postStrokeUser3, postStrokeUser3Byte = createStroke(username3, []float64{-79.38066843, 43.65483486})
		})

		It("should do match", func() {
			websocket.Message.Send(wsConnUser1, postStrokeUser1Byte)
			websocket.Message.Send(wsConnUser2, postStrokeUser2Byte)
			websocket.Message.Send(wsConnUser3, postStrokeUser3Byte)
			utils.Log.Infof("matchOtherTwo a1")
			matchOtherTwo(wsConnUser1, postStrokeUser2.Info, postStrokeUser3.Info)
			utils.Log.Infof("matchOtherTwo a2")
			matchOtherTwo(wsConnUser2, postStrokeUser1.Info, postStrokeUser3.Info)
			utils.Log.Infof("matchOtherTwo a3")
			matchOtherTwo(wsConnUser3, postStrokeUser2.Info, postStrokeUser1.Info)
		})

		It("a1 and a2 should match, a3 shouldn't match", func() {
			websocket.Message.Send(wsConnUser1, postStrokeUser1Byte)
			websocket.Message.Send(wsConnUser2, postStrokeUser2Byte)
			time.Sleep(3 * time.Second)
			websocket.Message.Send(wsConnUser3, postStrokeUser3Byte)
			resp1 := new(string)
			resp2 := new(string)
			resp3 := new(string)
			err1 := websocket.Message.Receive(wsConnUser1, resp1)
			err2 := websocket.Message.Receive(wsConnUser2, resp2)
			err3 := websocket.Message.Receive(wsConnUser3, resp3)
			Expect(err1).To(BeNil())
			Expect(*resp1).To(BeEquivalentTo(postStrokeUser2.Info))
			Expect(err2).To(BeNil())
			Expect(*resp2).To(BeEquivalentTo(postStrokeUser1.Info))
			Expect(err3).NotTo(BeNil())
		})

	})

})

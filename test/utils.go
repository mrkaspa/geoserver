package test

import (
	"encoding/json"

	. "github.com/onsi/gomega"
	"golang.org/x/net/websocket"
)

type postStroke struct {
	Info string
	Loc  []float64
}

func createPostStroke(info string, loc []float64) (*postStroke, []byte) {
	stroke := postStroke{
		Info: info,
		Loc:  loc,
	}
	json, _ := json.Marshal(stroke)
	return &stroke, json
}

func BeIn(arr []interface{}, val interface{}) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

func matchOtherTwo(wsConn *websocket.Conn, infoA, infoB string) {
	resp12 := new(string)
	resp13 := new(string)
	err12 := websocket.Message.Receive(wsConn, resp12)
	err13 := websocket.Message.Receive(wsConn, resp13)
	posibilities := []interface{}{infoA, infoB}
	Expect(err12).To(BeNil())
	Expect(err13).To(BeNil())
	Expect(*resp12).NotTo(BeEquivalentTo(resp13))
	Expect(BeIn(posibilities, *resp12)).To(BeTrue())
	Expect(BeIn(posibilities, *resp13)).To(BeTrue())
}

package test

import (
	"encoding/json"

	"github.com/mrkaspa/geoserver/models"
	. "github.com/onsi/gomega"
	"golang.org/x/net/websocket"
)

func createStroke(userID string, loc []float64) (*models.Stroke, []byte) {
	stroke := models.Stroke{
		UserID:   userID,
		Info:     userID,
		Location: loc,
	}
	json, _ := json.Marshal(&stroke)
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

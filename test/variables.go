package test

import (
	"github.com/mrkaspa/geoserver/models"
	"golang.org/x/net/websocket"
)

var (
	username1           = "a1"
	wsConnUser1         *websocket.Conn
	postStrokeUser1Byte []byte
	postStrokeUser1     *models.Stroke
	username2           = "a2"
	wsConnUser2         *websocket.Conn
	postStrokeUser2Byte []byte
	postStrokeUser2     *models.Stroke
	username3           = "a3"
	wsConnUser3         *websocket.Conn
	postStrokeUser3Byte []byte
	postStrokeUser3     *models.Stroke
)

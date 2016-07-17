package test

import (
	"fmt"
	"net/url"

	"golang.org/x/net/websocket"
)

func createClient(urlString, username string) *websocket.Conn {
	url, err := url.Parse(urlString)
	if err != nil {
		panic("malformed url")
	}
	newURL := fmt.Sprintf("ws://%s/ws/%s", url.Host, username)
	conn, err := websocket.Dial(newURL, "", urlString)
	if err != nil {
		panic(err.Error())
	}
	return conn
}

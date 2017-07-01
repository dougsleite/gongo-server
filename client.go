package main

import (
	"github.com/gorilla/websocket"
	"gopkg.in/mgo.v2"
)

type FindHandler func(string) (Handler, bool)

type Client struct {
	sendChannel chan Message
	websocket   *websocket.Conn
	findHandler FindHandler
	session     *mgo.Session
}

func (client *Client) Read() {
	var message Message
	for {
		if err := client.websocket.ReadJSON(&message); err != nil {
			break
		}
		if handler, found := client.findHandler(message.Name); found {
			handler(client, message.Data)
		}
	}
	client.websocket.Close()
}

func (client *Client) Write() {
	for msg := range client.sendChannel {
		if err := client.websocket.WriteJSON(msg); err != nil {
			break
		}
	}
	client.websocket.Close()
}

func (c *Client) Close() {
	close(c.sendChannel)
	c.session.Close()
}

func NewClient(websocket *websocket.Conn, findHandler FindHandler, session *mgo.Session) *Client {
	return &Client{
		sendChannel: make(chan Message),
		websocket:   websocket,
		findHandler: findHandler,
		session:     session,
	}
}

package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
}

// NewClient crea una nueva conexión WebSocket al servidor
func NewClient(url string) *Client {
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatalf("Error al conectarse al servidor WebSocket: %v", err)
	}
	return &Client{conn: conn}
}

// SendMessage envía un mensaje al servidor
func (c *Client) SendMessage(msg interface{}) {
	err := c.conn.WriteJSON(msg)
	if err != nil {
		log.Printf("Error al enviar mensaje: %v", err)
	}
}

// Listen escucha mensajes del servidor
func (c *Client) Listen(handleMessage func(msg map[string]interface{})) {
	for {
		var msg map[string]interface{}
		err := c.conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error al recibir mensaje: %v", err)
			break
		}
		handleMessage(msg)
	}
}

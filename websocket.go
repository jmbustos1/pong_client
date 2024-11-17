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
	log.Println("Conexión establecida con el servidor WebSocket:", url)
	return &Client{conn: conn}
}

// SendMessage envía un mensaje al servidor
func (c *Client) SendMessage(msg interface{}) {
	log.Printf("Enviando mensaje al servidor: %+v\n", msg)
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
		log.Printf("Mensaje recibido del servidor: %+v\n", msg)
		handleMessage(msg)
	}
}

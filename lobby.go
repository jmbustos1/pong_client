package main

import "log"

type Lobby struct {
	ID   string
	Name string
}

func (g *Game) createLobby(lobbyName string) {
	g.client.SendMessage(map[string]interface{}{
		"event":      "create_lobby",
		"player_id":  g.playerID,
		"lobby_name": lobbyName,
	})
}

func (g *Game) leaveLobby() {
	g.client.SendMessage(map[string]interface{}{
		"event":     "leave_lobby",
		"player_id": g.playerID,
		"lobby_id":  g.currentLobbyID, // Asegúrate de tener el ID del lobby actual
	})
	g.currentLobbyID = "" // Limpia el ID del lobby actual
	log.Println("Saliendo del lobby...")
}

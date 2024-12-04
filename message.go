package main

import "log"

type Message struct {
	Event    string   `json:"event"`
	PlayerID string   `json:"player_id,omitempty"`
	LobbyID  string   `json:"lobby_id,omitempty"`
	Lobbies  []string `json:"lobbies,omitempty"`
	Data     string   `json:"data,omitempty"`
}

func (g *Game) handleServerMessage(msg Message) {
	switch msg.Event {
	case "lobby_players":
		log.Printf("Jugadores en el lobby: %+v\n", msg.Lobbies)
		g.lobbyPlayers = msg.Lobbies // Actualiza la lista de jugadores en el lobby actual

	case "lobbies_list":
		log.Printf("Lista de lobbies recibida: %+v\n", msg.Lobbies)
		// Convertir []string en []Lobby
		var lobbies []Lobby
		for _, lobbyName := range msg.Lobbies {
			lobbies = append(lobbies, Lobby{
				ID:   "",        // Si no se recibe un ID, déjalo vacío
				Name: lobbyName, // Usa el nombre del lobby
			})
		}
		g.lobbies = lobbies // Asignar la lista convertida
	case "lobby_joined":
		log.Printf("Unido al lobby: %s\n", msg.LobbyID)
		g.state = LobbyIn
	case "player_joined":
		log.Printf("Nuevo jugador conectado: %s\n", msg.PlayerID)
	case "game_start":
		log.Println("El juego ha comenzado.")
		g.state = Playing
	default:
		log.Printf("Evento desconocido recibido: %s\n", msg.Event)
	}
}

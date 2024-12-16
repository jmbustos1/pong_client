package main

import "log"

type Message struct {
	Event    string      `json:"event"`
	PlayerID string      `json:"player_id,omitempty"`
	LobbyID  string      `json:"lobby_id,omitempty"`
	Lobbies  interface{} `json:"lobbies,omitempty"` // Cambia esto a interface{}
	Data     string      `json:"data,omitempty"`
}

func (g *Game) HandleServerMessage(msg Message) {
	log.Printf("Evento recibido: %s, Detalles: %+v\n", msg.Event, msg)
	if msg.Event == "test_message" {
		log.Println("Mensaje de prueba recibido:", msg.Data)
	}
	log.Println("MENSAJEEEEEE", msg.Event)
	switch msg.Event {
	case "lobby_players":
		// Actualizar la lista de jugadores en el lobby
		log.Printf("Jugadores en el lobby actualizados: %+v\n", msg.Lobbies)
		if players, ok := msg.Lobbies.([]interface{}); ok {
			var updatedPlayers []string
			for _, player := range players {
				if playerName, ok := player.(string); ok {
					updatedPlayers = append(updatedPlayers, playerName)
				}
			}
			g.lobbyPlayers = updatedPlayers
		}
		log.Printf("Jugadores actualizados: %+v\n", g.lobbyPlayers)

	case "lobbies_list":
		log.Printf("Lista de lobbies recibida: %+v\n", msg.Lobbies)

		// Convertir la lista de lobbies en []Lobby
		var lobbies []Lobby
		if list, ok := msg.Lobbies.([]interface{}); ok {
			for _, lobbyData := range list {
				if lobbyMap, ok := lobbyData.(map[string]interface{}); ok {
					// Extraer los campos ID y Name del mapa
					id, _ := lobbyMap["id"].(string)
					name, _ := lobbyMap["name"].(string)

					lobbies = append(lobbies, Lobby{
						ID:   id,
						Name: name,
					})
				}
			}
		}
		g.lobbies = lobbies // Actualiza la lista de lobbies

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

package main

func (g *Game) createLobby(lobbyName string) {
	g.client.SendMessage(map[string]interface{}{
		"event":      "create_lobby",
		"player_id":  g.playerID,
		"lobby_name": lobbyName,
	})
}

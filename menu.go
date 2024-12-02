package main

import (
	"image/color"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const (
	inputCooldown = 200 * time.Millisecond
	stateCooldown = 1000 * time.Millisecond
)

func (g *Game) startNewGame() {
	g.state = Playing
	g.player1Y = screenHeight/2 - paddleHeight/2
	g.player2Y = screenHeight/2 - paddleHeight/2
	g.player1X = 20
	g.player2X = screenWidth - 30
	g.ballX = screenWidth / 2
	g.ballY = screenHeight / 2
	g.ballDirection = Vector{X: 1, Y: 1}
}

func (g *Game) updateMenu() {
	if time.Since(g.lastInputTime) < inputCooldown {
		return
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.menuSelection = (g.menuSelection + 2) % 3
		g.lastInputTime = time.Now()
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.menuSelection = (g.menuSelection + 1) % 3
		g.lastInputTime = time.Now()
	}

	// Solo permite el cambio de estado si ha pasado el cooldown
	if ebiten.IsKeyPressed(ebiten.KeyEnter) && time.Since(g.stateChangeTime) > stateCooldown {
		switch g.menuSelection {
		case 0:
			g.startNewGame()
		case 1:
			g.state = LobbyMenu
			g.lobbyMenuSelection = 0       // Reinicia la selección del menú de lobbies
			g.stateChangeTime = time.Now() // Marca el tiempo del cambio de estado
			g.client.SendMessage(Message{
				Event: "get_lobbies",
			})
		case 2:
			log.Println("Exiting game.")
			os.Exit(0)
		}
	}
}

func (g *Game) drawMenu(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 50, G: 50, B: 50, A: 255})
	menuItems := []string{"New Game", "Lobby", "Quit"}
	for i, item := range menuItems {
		yPos := screenHeight/2 + (i * fontSize) - 20
		if i == g.menuSelection {
			text.Draw(screen, "> "+item, g.font, screenWidth/2-60, yPos, color.RGBA{R: 255, G: 100, B: 100, A: 255})
		} else {
			text.Draw(screen, item, g.font, screenWidth/2-30, yPos, color.White)
		}
	}
}

func (g *Game) drawLobbyMenu(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 30, G: 30, B: 30, A: 255}) // Fondo gris oscuro

	text.Draw(screen, "Available Lobbies:", g.font, screenWidth/2-100, 50, color.White)

	// Dibujar la lista de lobbies dinámica
	for i, lobby := range g.lobbies {
		yPos := 100 + (i * fontSize)   // Espaciado vertical
		if g.lobbyMenuSelection == i { // Compara con el índice actual
			text.Draw(screen, "> "+lobby.Name, g.font, screenWidth/2-60, yPos, color.RGBA{R: 255, G: 100, B: 100, A: 255})
		} else {
			text.Draw(screen, lobby.Name, g.font, screenWidth/2-30, yPos, color.White)
		}
	}

	// Dibujar las opciones adicionales
	createLobbyYPos := 100 + (len(g.lobbies) * fontSize) // "Create Lobby" justo después de los lobbies
	backYPos := createLobbyYPos + fontSize + 10          // "Back" después de "Create Lobby"

	if g.lobbyMenuSelection == len(g.lobbies) { // Selección en "Create Lobby"
		text.Draw(screen, "> Create Lobby", g.font, screenWidth/2-60, createLobbyYPos, color.RGBA{R: 100, G: 255, B: 100, A: 255})
	} else {
		text.Draw(screen, "Create Lobby", g.font, screenWidth/2-30, createLobbyYPos, color.White)
	}

	if g.lobbyMenuSelection == len(g.lobbies)+1 { // Selección en "Back"
		text.Draw(screen, "> Back", g.font, screenWidth/2-60, backYPos, color.RGBA{R: 255, G: 100, B: 100, A: 255})
	} else {
		text.Draw(screen, "Back", g.font, screenWidth/2-30, backYPos, color.White)
	}
}

func (g *Game) drawLobby(screen *ebiten.Image) {
	screen.Fill(color.RGBA{R: 50, G: 50, B: 50, A: 255}) // Fondo gris más claro

	// Jugadores conectados (simulados aquí)
	players := []string{"Player 1", "Player 2"} // Cambiar por datos reales del servidor

	text.Draw(screen, "Lobby Players:", g.font, screenWidth/2-100, 50, color.White)

	// Dibujar la lista de jugadores
	for i, player := range players {
		yPos := 100 + (i * fontSize) // Espaciado vertical
		text.Draw(screen, player, g.font, screenWidth/2-30, yPos, color.White)
	}

	// Opciones: "Start Game" y "Back"
	yPos := screenHeight - 80
	text.Draw(screen, "Start Game", g.font, screenWidth/2-60, yPos, color.RGBA{R: 100, G: 255, B: 100, A: 255})
	yPos += fontSize + 10
	text.Draw(screen, "Back", g.font, screenWidth/2-60, yPos, color.RGBA{R: 255, G: 100, B: 100, A: 255})
}

func (g *Game) updateLobbyMenu() {
	// Evita procesar entradas si el cooldown no ha terminado
	if time.Since(g.lastInputTime) < inputCooldown || time.Since(g.stateChangeTime) < stateCooldown {
		return
	}

	totalOptions := len(g.lobbies) + 2 // Lobbies dinámicos + "Create Lobby" + "Back"

	// Navegación en el menú de lobbies
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.lobbyMenuSelection = (g.lobbyMenuSelection - 1 + totalOptions) % totalOptions
		g.lastInputTime = time.Now()
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.lobbyMenuSelection = (g.lobbyMenuSelection + 1) % totalOptions
		g.lastInputTime = time.Now()
	}
	log.Print(g.lobbyMenuSelection)
	// Selección de opciones
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		switch {
		case g.lobbyMenuSelection == 0: // Crear un lobby
			g.createLobby("My Awesome Lobby")
			g.state = LobbyIn
			g.stateChangeTime = time.Now()
		case g.lobbyMenuSelection == len(g.lobbies)+1: // Volver al menú principal
			g.state = Menu
			g.stateChangeTime = time.Now()
		default: // Unirse a un lobby existente
			selectedLobby := g.lobbies[g.lobbyMenuSelection-1] // Ajusta selección dinámica
			g.joinLobby(selectedLobby.ID)
			g.state = LobbyIn
			g.stateChangeTime = time.Now()
		}
	}
}

func (g *Game) joinLobby(lobbyID string) {
	g.client.SendMessage(map[string]interface{}{
		"event":     "join_lobby",
		"player_id": g.playerID,
		"lobby_id":  lobbyID,
	})
}

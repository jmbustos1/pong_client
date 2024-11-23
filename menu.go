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

	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		switch g.menuSelection {
		case 0:
			g.startNewGame()
		case 1:
			g.state = LobbyMenu
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

func (g *Game) updateLobbyMenu() {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		switch g.menuSelection {
		case 0: // Crear un lobby
			g.createLobby("My Awesome Lobby")
			g.state = Lobby // Cambia al estado Lobby para esperar al jugador 2
		case 1: // Unirse a un lobby existente
			g.joinLobby("selected_lobby_id") // Usa el ID del lobby seleccionado
			g.state = Lobby
		case 2: // Volver al menú principal
			g.state = Menu
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

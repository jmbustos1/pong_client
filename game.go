package main

import (
	"crypto/sha1"
	"fmt"
	"image/color"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

type Vector struct {
	X, Y float64
}

type GameState int

const (
	Menu      GameState = iota
	LobbyMenu           // Menú para ver o crear lobbies
	Playing
	LobbyIn
)

type Game struct {
	state              GameState
	menuSelection      int       // Para el menú principal
	lobbyMenuSelection int       // Para el menú de lobbies
	stateChangeTime    time.Time // Marca cuándo se cambió el estado
	lastInputTime      time.Time
	font               font.Face
	player1Y           float64
	playerID           string
	client             *Client
	player2Y           float64
	player1X           float64
	player2X           float64
	ballX              float64
	ballY              float64
	ballDirection      Vector
	lastHitPaddle      bool
	lobbies            []Lobby
	selectedLobby      int
	currentLobbyID     string   // Almacena el ID del lobby actual
	lobbyPlayers       []string // Lista de jugadores en el lobby actual
	serverError        string
}

func generatePlayerID() string {
	// Usa el hostname como base para generar el ID
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	// Hashea el hostname para obtener un ID único
	hash := sha1.Sum([]byte(hostname))
	return fmt.Sprintf("%x", hash[:4]) // Usa los primeros 4 bytes como ID
}

// NewGame inicializa el juego y carga los recursos necesarios.
func NewGame() *Game {
	fontFace := loadFont() // Asume que tienes una función `loadFont()` en `assets.go`
	pID := generatePlayerID()
	fmt.Printf("Generated Player ID: %s\n", pID)

	client, err := NewClient("ws://172.17.0.1:8088/ws")
	if err != nil {
		log.Println("No se pudo conectar al servidor:", err)
		client = nil // Permitir que el juego continúe sin servidor
	}
	game := &Game{
		state:         Menu,
		menuSelection: 0,
		lastInputTime: time.Now(),
		font:          fontFace,
		player1X:      20,
		player2X:      screenWidth - 30,
		playerID:      pID,
		client:        client,
		ballX:         screenWidth / 2,
		ballY:         screenHeight / 2,
		ballDirection: Vector{X: 1, Y: 1},
	}
	// Define el comportamiento al perder la conexión
	onDisconnect := func() {
		log.Println("Servidor caído. Volviendo al menú principal...")
		game.state = Menu
		game.serverError = "Servidor desconectado"
	}
	// Escuchar mensajes si hay conexión
	if client != nil {
		go client.Listen(game.HandleServerMessage, onDisconnect)
	}
	return game
}

func (g *Game) updateGame() {
	// Actualiza la posición de las palas
	g.updatePaddles(g.client, g.playerID)

	// Actualiza la posición de la pelota y maneja colisiones
	g.updateBall()
}

func (g *Game) Update() error {
	g.ProcessIncomingMessages()
	switch g.state {
	case Menu:
		g.updateMenu()
	case LobbyMenu:
		g.updateLobbyMenu()
	case LobbyIn:
		g.updateLobby()
	case Playing:
		g.updateGame()
	}
	return nil
}

// Draw renderiza la pantalla del juego según el estado actual.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black) // Fondo negro para el juego
	//log.Printf("Estado actual: %v\n", g.state, g.menuSelection, g.lobbyMenuSelection, g.lobbies)
	switch g.state {
	case Menu:
		g.drawMenu(screen)
	case Playing:
		g.drawGame(screen)
	case LobbyMenu:
		g.drawLobbyMenu(screen)
	case LobbyIn:
		g.drawLobby(screen)
	}
}

// drawGame renderiza los elementos del juego en curso
func (g *Game) drawGame(screen *ebiten.Image) {
	// Dibujar palas
	player1Paddle := ebiten.NewImage(paddleWidth, paddleHeight)
	player1Paddle.Fill(color.White)
	op1 := &ebiten.DrawImageOptions{}
	op1.GeoM.Translate(g.player1X, g.player1Y)
	screen.DrawImage(player1Paddle, op1)

	player2Paddle := ebiten.NewImage(paddleWidth, paddleHeight)
	player2Paddle.Fill(color.White)
	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Translate(g.player2X, g.player2Y)
	screen.DrawImage(player2Paddle, op2)

	// Dibujar bola
	ball := ebiten.NewImage(ballSize, ballSize)
	ball.Fill(color.White)
	op3 := &ebiten.DrawImageOptions{}
	op3.GeoM.Translate(g.ballX, g.ballY)
	screen.DrawImage(ball, op3)
}

// Layout ajusta el tamaño de la ventana
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

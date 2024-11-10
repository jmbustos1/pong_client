package main

import (
	// "fmt"
	"image/color"
	"io"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	screenWidth   = 640
	screenHeight  = 480
	paddleWidth   = 10
	paddleHeight  = 80
	ballSize      = 10
	paddleSpeed   = 5
	ballSpeed     = 4
	inputCooldown = 200 * time.Millisecond
	fontSize      = 24 // Tamaño de la fuente para el menú
)

type GameState int

const (
	Menu GameState = iota
	Playing
	Lobby
)

// Vector define la dirección de la bola
type Vector struct {
	X, Y float64
}

type Game struct {
	state              GameState
	menuSelection      int
	lastInputTime      time.Time
	player1Y, player2Y float64
	player1X, player2X float64
	ballX, ballY       float64
	ballDirection      Vector
	lastHitPaddle      bool
	font               font.Face
}

func NewGame() *Game {
	// Cargar la fuente desde el archivo TTF
	fontFile, err := os.Open("assets/Roboto-Regular.ttf")
	if err != nil {
		log.Fatal("Error abriendo el archivo de fuente:", err)
	}
	defer fontFile.Close()

	// Lee el archivo completo
	fontBytes, err := io.ReadAll(fontFile)
	if err != nil {
		log.Fatal("Error leyendo los bytes del archivo de fuente:", err)
	}

	// Verifica el tamaño del archivo
	// fmt.Println("Tamaño de la fuente en bytes:", len(fontBytes))
	// Parsear la fuente y crear un `font.Face` para el tamaño deseado
	ttf, err := opentype.Parse(fontBytes)
	// fmt.Println(ttf)
	if err != nil {
		log.Fatal(err)
	}

	fontFace, err := opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	return &Game{
		state:         Menu,
		menuSelection: 0,
		lastInputTime: time.Now(),
		font:          fontFace,
	}
}

func (g *Game) Update() error {
	switch g.state {
	case Menu:
		g.updateMenu()
	case Playing:
		g.updateGame()
	case Lobby:
	}
	return nil
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
			g.state = Lobby
		case 2:
			log.Println("Exiting game.")
			os.Exit(0)
		}
	}
}

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

func (g *Game) updateGame() {
	// Controles del Jugador 1 (W y S)
	if ebiten.IsKeyPressed(ebiten.KeyW) && g.player1Y > 0 {
		g.player1Y -= paddleSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) && g.player1Y < screenHeight-paddleHeight {
		g.player1Y += paddleSpeed
	}

	// Controles del Jugador 2 (Flechas ↑ y ↓)
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) && g.player2Y > 0 {
		g.player2Y -= paddleSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) && g.player2Y < screenHeight-paddleHeight {
		g.player2Y += paddleSpeed
	}

	g.ballX += g.ballDirection.X * ballSpeed
	g.ballY += g.ballDirection.Y * ballSpeed

	if g.ballY <= 0 || g.ballY >= screenHeight-ballSize {
		g.ballDirection.Y *= -1
	}

	g.checkPaddleCollision(g.player1X, g.player1Y, true)
	g.checkPaddleCollision(g.player2X, g.player2Y, false)
}

func (g *Game) checkPaddleCollision(paddleX, paddleY float64, isLeftPaddle bool) {
	if (isLeftPaddle && g.ballX <= paddleX+paddleWidth) || (!isLeftPaddle && g.ballX+ballSize >= paddleX) {
		if g.ballY+ballSize >= paddleY && g.ballY <= paddleY+paddleHeight {
			if g.lastHitPaddle {
				return
			}
			g.lastHitPaddle = true

			collisionPoint := (g.ballY - paddleY) / paddleHeight

			if collisionPoint <= 0.1 && g.ballDirection.Y > 0 {
				g.ballDirection.X *= -1
				g.ballDirection.Y *= -1
			} else if collisionPoint >= 0.9 && g.ballDirection.Y < 0 {
				g.ballDirection.X *= -1
				g.ballDirection.Y *= -1
			} else {
				g.ballDirection.X *= -1
			}

			if isLeftPaddle {
				g.ballX = paddleX + paddleWidth
			} else {
				g.ballX = paddleX - ballSize
			}
		} else {
			g.lastHitPaddle = false
		}
	} else {
		g.lastHitPaddle = false
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case Menu:
		g.drawMenu(screen)
	case Playing:
		g.drawGame(screen)
	case Lobby:
		screen.Fill(color.RGBA{R: 30, G: 30, B: 30, A: 255})
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

func (g *Game) drawGame(screen *ebiten.Image) {
	screen.Fill(color.Black)
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

	bal := ebiten.NewImage(ballSize, ballSize)
	bal.Fill(color.White)
	op3 := &ebiten.DrawImageOptions{}
	op3.GeoM.Translate(g.ballX, g.ballY)
	screen.DrawImage(bal, op3)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	game := NewGame()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Pong Game with Custom Font")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

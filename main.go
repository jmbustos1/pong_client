package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	// "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
)

const (
	screenWidth  = 640
	screenHeight = 480
	paddleWidth  = 10
	paddleHeight = 80
	ballSize     = 10
	paddleSpeed  = 5
	ballSpeed    = 4
)

var ballDirection int = 1

type Game struct {
	player1Y, player2Y float64
	player1X, player2X float64
	ballX, ballY       float64
}

func (g *Game) Update() error {
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
	// Movimiento de la pelota basado en la dirección
	switch ballDirection {
	case 0: // derecha-abajo
		g.ballX += ballSpeed
		g.ballY += ballSpeed
	case 1: // izquierda-abajo
		g.ballX -= ballSpeed
		g.ballY += ballSpeed
	case 2: // derecha-arriba
		g.ballX += ballSpeed
		g.ballY -= ballSpeed
	case 3: // izquierda-arriba
		g.ballX -= ballSpeed
		g.ballY -= ballSpeed
	}

	// Verificación de colisiones con los bordes
	if g.ballY <= 0 { // Rebote en el borde superior
		if ballDirection == 2 {
			ballDirection = 0
		} else if ballDirection == 3 {
			ballDirection = 1
		}
	}
	if g.ballY >= screenHeight-ballSize { // Rebote en el borde inferior
		if ballDirection == 0 {
			ballDirection = 2
		} else if ballDirection == 1 {
			ballDirection = 3
		}
	}
	if g.ballX <= 0 { // Rebote en el borde izquierdo
		if ballDirection == 3 {
			ballDirection = 2
		} else if ballDirection == 1 {
			ballDirection = 0
		}
	}
	if g.ballX >= screenWidth-ballSize { // Rebote en el borde derecho
		if ballDirection == 0 {
			ballDirection = 1
		} else if ballDirection == 2 {
			ballDirection = 3
		}
	}

	// Verificación de colisiones con las palas
	// Rebote en la pala izquierda (Jugador 1)
	if g.ballX <= g.player1X+paddleWidth && // Colisión en el borde de la pala
		g.ballY+ballSize >= g.player1Y && // Dentro del rango Y de la pala
		g.ballY <= g.player1Y+paddleHeight {

		if ballDirection == 1 {
			ballDirection = 0
		} else if ballDirection == 3 {
			ballDirection = 2
		}
	}

	// Rebote en la pala derecha (Jugador 2)
	if g.ballX+ballSize >= g.player2X && // Colisión en el borde de la pala
		g.ballY+ballSize >= g.player2Y && // Dentro del rango Y de la pala
		g.ballY <= g.player2Y+paddleHeight {

		if ballDirection == 0 {
			ballDirection = 1
		} else if ballDirection == 2 {
			ballDirection = 3
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Fondo
	screen.Fill(color.Black)

	// Dibujar palas de los jugadores
	player1Paddle := ebiten.NewImage(paddleWidth, paddleHeight)
	player1Paddle.Fill(color.White)
	op1 := &ebiten.DrawImageOptions{}
	op1.GeoM.Translate(20, g.player1Y)
	screen.DrawImage(player1Paddle, op1)

	player2Paddle := ebiten.NewImage(paddleWidth, paddleHeight)
	player2Paddle.Fill(color.White)
	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Translate(screenWidth-30, g.player2Y)
	screen.DrawImage(player2Paddle, op2)

	// bola
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
	game := &Game{
		player1Y: screenHeight/2 - paddleHeight/2,
		player2Y: screenHeight/2 - paddleHeight/2,
		player1X: 20,
		player2X: screenWidth - 30,
		ballX:    screenWidth / 2,
		ballY:    screenHeight / 2,
	}
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Pong Game")

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}

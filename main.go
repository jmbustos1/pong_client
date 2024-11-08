package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
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

// Vector define la dirección de la bola
type Vector struct {
	X, Y float64
}

type Game struct {
	player1Y, player2Y float64
	player1X, player2X float64
	ballX, ballY       float64
	ballDirection      Vector
	lastHitPaddle      bool
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

	// Movimiento de la bola
	g.ballX += g.ballDirection.X * ballSpeed
	g.ballY += g.ballDirection.Y * ballSpeed

	// Verificación de colisiones con los bordes superior e inferior
	if g.ballY <= 0 || g.ballY >= screenHeight-ballSize {
		g.ballDirection.Y *= -1 // Invertir dirección vertical
	}

	// Verificación de colisiones con las palas
	g.checkPaddleCollision(g.player1X, g.player1Y, true)  // Pala izquierda
	g.checkPaddleCollision(g.player2X, g.player2Y, false) // Pala derecha

	return nil
}

// checkPaddleCollision calcula el rebote dependiendo de la posición de colisión en la pala
func (g *Game) checkPaddleCollision(paddleX, paddleY float64, isLeftPaddle bool) {
	// Verificar si la bola está en la misma coordenada X que la pala
	if (isLeftPaddle && g.ballX <= paddleX+paddleWidth) || (!isLeftPaddle && g.ballX+ballSize >= paddleX) {
		// Verificar si la bola está dentro de la altura de la pala
		if g.ballY+ballSize >= paddleY && g.ballY <= paddleY+paddleHeight {
			// Evitar rebotes consecutivos en el mismo frame
			if g.lastHitPaddle {
				return
			}
			g.lastHitPaddle = true // Marcar que la pala fue golpeada en este frame

			// Calcular la posición de colisión en la pala (de 0.0 a 1.0)
			collisionPoint := (g.ballY - paddleY) / paddleHeight

			// Si la bola viene hacia abajo y choca con la parte superior de la pala (10% superior)
			if collisionPoint <= 0.1 && g.ballDirection.Y > 0 {
				// Rebote de 180 grados: invertir X y Y
				g.ballDirection.X *= -1
				g.ballDirection.Y *= -1
			} else if collisionPoint >= 0.9 && g.ballDirection.Y < 0 { // Si viene hacia arriba y choca con la parte inferior (10% inferior)
				// Rebote de 180 grados: invertir X y Y
				g.ballDirection.X *= -1
				g.ballDirection.Y *= -1
			} else {
				// Rebote normal: solo invertimos X
				g.ballDirection.X *= -1
			}

			// Separar la bola de la pala para evitar rebotes infinitos
			if isLeftPaddle {
				g.ballX = paddleX + paddleWidth // Posicionar justo a la derecha de la pala
			} else {
				g.ballX = paddleX - ballSize // Posicionar justo a la izquierda de la pala
			}
		} else {
			g.lastHitPaddle = false // Resetear el estado si no hay colisión en este frame
		}
	} else {
		g.lastHitPaddle = false // Resetear el estado si no hay colisión en este frame
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Fondo
	screen.Fill(color.Black)

	// Dibujar palas de los jugadores
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

	// Bola
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
		player1Y:      screenHeight/2 - paddleHeight/2,
		player2Y:      screenHeight/2 - paddleHeight/2,
		player1X:      20,
		player2X:      screenWidth - 30,
		ballX:         screenWidth / 2,
		ballY:         screenHeight / 2,
		ballDirection: Vector{X: 1, Y: 1}, // Dirección inicial
	}
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Pong Game")

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}

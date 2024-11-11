package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	paddleWidth  = 10
	paddleHeight = 80
	paddleSpeed  = 5
)

// Actualiza la posición de las palas con las teclas W/S y las flechas ↑/↓.
func (g *Game) updatePaddles() {
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
}

// Verifica y maneja las colisiones entre las palas y la pelota.
func (g *Game) checkPaddleCollision(paddleX, paddleY float64, isLeftPaddle bool) {
	if (isLeftPaddle && g.ballX <= paddleX+paddleWidth) || (!isLeftPaddle && g.ballX+ballSize >= paddleX) {
		if g.ballY+ballSize >= paddleY && g.ballY <= paddleY+paddleHeight {
			if g.lastHitPaddle {
				return
			}
			g.lastHitPaddle = true

			collisionPoint := (g.ballY - paddleY) / paddleHeight

			// Rebote basado en la posición de colisión
			if collisionPoint <= 0.1 && g.ballDirection.Y > 0 {
				g.ballDirection.X *= -1
				g.ballDirection.Y *= -1
			} else if collisionPoint >= 0.9 && g.ballDirection.Y < 0 {
				g.ballDirection.X *= -1
				g.ballDirection.Y *= -1
			} else {
				g.ballDirection.X *= -1
			}

			// Ajusta la posición de la pelota para evitar múltiples colisiones
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

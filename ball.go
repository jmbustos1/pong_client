package main

const (
	ballSize  = 10
	ballSpeed = 4
)

// Actualiza la posición de la pelota y verifica colisiones con los bordes.
func (g *Game) updateBall() {
	// Movimiento de la pelota
	g.ballX += g.ballDirection.X * ballSpeed
	g.ballY += g.ballDirection.Y * ballSpeed

	// Rebote en los bordes superior e inferior
	if g.ballY <= 0 || g.ballY >= screenHeight-ballSize {
		g.ballDirection.Y *= -1
	}

	// Colisiones con las palas
	g.checkPaddleCollision(g.player1X, g.player1Y, true)  // Pala izquierda
	g.checkPaddleCollision(g.player2X, g.player2Y, false) // Pala derecha
}

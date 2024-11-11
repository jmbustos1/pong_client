package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// Crear una nueva instancia del juego
	game := NewGame()

	// Configuración de la ventana
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Pong Game with Custom Font")

	// Ejecutar el juego
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

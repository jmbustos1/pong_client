package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// Crear una nueva instancia del juego
	game := NewGame()

	// Configuración de la ventana
	// Hacer que la ventana sea resizable y a pantalla completa
	// ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	// ebiten.SetFullscreen(false) // Cambia a true si quieres comenzar en pantalla completa

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Pong Game with Custom Font")

	// Ejecutar el juego
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

package main

// PENDIENTES:
// REVISAR QUE SE PUEDA CORRER JUEGO SERVIDOR ESTA CAIDO
// REVISAR UN BUEN DEBOUCE DE BOTONES
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
	// CREAR ESTRUCTURA PARA MEESAGE SEND

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Pong Game with Custom Font")

	// Ejecutar el juego
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

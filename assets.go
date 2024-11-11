package main

import (
	"io"
	"log"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const fontSize = 24

// loadFont carga la fuente desde un archivo TTF.
func loadFont() font.Face {
	fontFile, err := os.Open("assets/Roboto-Regular.ttf")
	if err != nil {
		log.Fatal("Error abriendo el archivo de fuente:", err)
	}
	defer fontFile.Close()

	fontBytes, err := io.ReadAll(fontFile)
	if err != nil {
		log.Fatal("Error leyendo los bytes del archivo de fuente:", err)
	}

	ttf, err := opentype.Parse(fontBytes)
	if err != nil {
		log.Fatal("Error parseando la fuente:", err)
	}

	fontFace, err := opentype.NewFace(ttf, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal("Error creando `font.Face`:", err)
	}

	return fontFace
}

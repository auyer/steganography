package main

import (
	"image"
	"log"
	"testing"
)

var rawInputFile = "./stegosaurus.png"
var encodedInputFile = "./encoded_stegosaurus.png"

func TestOpenImageFromPath(t *testing.T) {
	img, err := OpenImageFromPath(rawInputFile)
	if err != nil {
		log.Printf("Error opening or Decoding file %s: %v", rawInputFile, err)
		t.FailNow()
	}
	if (img.Bounds().Bounds() != image.Rectangle{image.Point{0, 0}, image.Point{1195, 642}}) {
		log.Printf("Image has wrong size")
		t.FailNow()

	}
}

func TestEmptyPathHelperFunction(t *testing.T) {
	_, err := OpenImageFromPath(" ")
	if err == nil {
		log.Print("Empty path given, err could not be nil.")
		t.FailNow()
	}
}

package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
)

var inputFile string
var outputFile string

func main() {

	inputFile = "input.png"
	outputFile = "output.png"

	rgbIm := imageToRGBA(decodeImage(inputFile))

	//c := rgbIm.At(0, 0)

	var c color.RGBA

	for x := 0; x < 20; x++ {
		for y := 0; y < 20; y++ {
			rgbIm.SetRGBA(x, y, c)
		}
	}

	encodePNG(outputFile, rgbIm)

}

func imageToRGBA(src image.Image) *image.RGBA {
	b := src.Bounds()

	var m *image.RGBA
	var width, height int

	width = b.Dx()
	height = b.Dy()

	m = image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(m, m.Bounds(), src, b.Min, draw.Src)
	return m
}

func decodeImage(filename string) image.Image {
	inFile, err := os.Open(filename)

	if err != nil {
		log.Fatalf("Error opening file %s: %v", filename, err)
	}

	defer inFile.Close()

	reader := bufio.NewReader(inFile)

	img, _, err := image.Decode(reader)

	fmt.Println("Read", filename)
	return img
}

func encodePNG(filename string, img image.Image) {
	fo, err := os.Create(filename)

	if err != nil {
		log.Fatalf("Error creating file %s: %v", filename, err)
	}

	defer fo.Close()
	defer fo.Sync()

	writer := bufio.NewWriter(fo)
	defer writer.Flush()

	err = png.Encode(writer, img)

	fmt.Println("Wrote to", filename)
}

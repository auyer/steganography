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

	var b byte = 58

	print("LSB: ", getLSB(b), "\n")
	setLSB(&b, true)
	print("LSB: ", getLSB(b), "\n")

	print("MAX: ", maxEncodeSize(rgbIm), "\n")

	var c color.RGBA

	for x := 0; x < 20; x++ {
		for y := 0; y < 20; y++ {
			rgbIm.SetRGBA(x, y, c)
		}
	}

	encodePNG(outputFile, rgbIm)

}

// Convert given image to RGBA image
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

// Read and return an image at the given path
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

// Will write out a given image to a given path in filename
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
}

// Given an image will find how many bytes can be stored in that image using least significant bit encoding
func maxEncodeSize(img image.Image) int {

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	return int(((width * height * 3) / 8))
}

// Given a byte, will return the least significant bit of that byte
func getLSB(b byte) byte {
	b &= 1
	return b
}

// Given a byte will set that byte's least significant bit to a given value (where true is 1 and false is 0)
func setLSB(b *byte, bit bool) {
	if bit == true {
		*b = *b | 1
	} else if bit == false {
		var mask byte = 0xFE
		*b = *b & mask
	}
}

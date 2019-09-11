package steganography

import (
	"bufio"
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
	"testing"
)

var rawInputFilePng = "./examples/stegosaurus.png"
var rawInputFileJpg = "./examples/stegosaurus.jpg"
var encodedInputFilePng = "./examples/encoded_stegosaurus.png"
var encodedInputFileJpg = "./examples/encoded_stegosaurus.jpg"

var bitmessage = []uint8{84, 104, 101, 113, 117, 97, 100, 114, 117, 112, 101, 100, 97, 108, 83, 116, 101, 103, 111, 115, 97, 117, 114, 117, 115, 105, 115, 111, 110, 101, 111, 102, 116, 104, 101, 109, 111, 115, 116, 101, 97, 115, 105, 108, 121, 105, 100, 101, 110, 116, 105, 102, 105, 97, 98, 108, 101, 100, 105, 110, 111, 115, 97, 117, 114, 103, 101, 110, 101, 114, 97, 44, 100, 117, 101, 116, 111, 116, 104, 101, 100, 105, 115, 116, 105, 110, 99, 116, 105, 118, 101, 100, 111, 117, 98, 108, 101, 114, 111, 119, 111, 102, 107, 105, 116, 101, 45, 115, 104, 97, 112, 101, 100, 112, 108, 97, 116, 101, 115, 114, 105, 115, 105, 110, 103, 118, 101, 114, 116, 105, 99, 97, 108, 108, 121, 97, 108, 111, 110, 103, 116, 104, 101, 114, 111, 117, 110, 100, 101, 100, 98, 97, 99, 107, 97, 110, 100, 116, 104, 101, 116, 119, 111, 112, 97, 105, 114, 115, 111, 102, 108, 111, 110, 103, 115, 112, 105, 107, 101, 115, 101, 120, 116, 101, 110, 100, 105, 110, 103, 104, 111, 114, 105, 122, 111, 110, 116, 97, 108, 108, 121, 110, 101, 97, 114, 116, 104, 101, 101, 110, 100, 111, 102, 116, 104, 101, 116, 97, 105, 108, 46, 65, 108, 116, 104, 111, 117, 103, 104, 108, 97, 114, 103, 101, 105, 110, 100, 105, 118, 105, 100, 117, 97, 108, 115, 99, 111, 117, 108, 100, 103, 114, 111, 119, 117, 112, 116, 111, 57, 109, 40, 50, 57, 46, 53, 102, 116, 41, 105, 110, 108, 101, 110, 103, 116, 104, 91, 52, 93, 97, 110, 100, 53, 46, 51, 116, 111, 55, 109, 101, 116, 114, 105, 99, 116, 111, 110, 115, 40, 53, 46, 56, 116, 111, 55, 46, 55, 115, 104, 111, 114, 116, 116, 111, 110, 115, 41, 105, 110, 119, 101, 105, 103, 104, 116, 44, 91, 53, 93, 91, 54, 93, 116, 104, 101, 118, 97, 114, 105, 111, 117, 115, 115, 112, 101, 99, 105, 101, 115, 111, 102, 83, 116, 101, 103, 111, 115, 97, 117, 114, 117, 115, 119, 101, 114, 101, 100, 119, 97, 114, 102, 101, 100, 98, 121, 99, 111, 110, 116, 101, 109, 112, 111, 114, 97, 114, 105, 101, 115, 44, 116, 104, 101, 103, 105, 97, 110, 116, 115, 97, 117, 114, 111, 112, 111, 100, 115, 46, 83, 111, 109, 101, 102, 111, 114, 109, 111, 102, 97, 114, 109, 111, 114, 97, 112, 112, 101, 97, 114, 115, 116, 111, 104, 97, 118, 101, 98, 101, 101, 110, 110, 101, 99, 101, 115, 115, 97, 114, 121, 44, 97, 115, 83, 116, 101, 103, 111, 115, 97, 117, 114, 117, 115, 115, 112, 101, 99, 105, 101, 115, 99, 111, 101, 120, 105, 115, 116, 101, 100, 119, 105, 116, 104, 108, 97, 114, 103, 101, 112, 114, 101, 100, 97, 116, 111, 114, 121, 116, 104, 101, 114, 111, 112, 111, 100, 100, 105, 110, 111, 115, 97, 117, 114, 115, 44, 115, 117, 99, 104, 97, 115, 65, 108, 108, 111, 115, 97, 117, 114, 117, 115, 97, 110, 100, 67, 101, 114, 97, 116, 111, 115, 97, 117, 114, 117, 115, 46}

func TestEncodeFromPngFile(t *testing.T) {

	inFile, err := os.Open(rawInputFilePng)
	if err != nil {
		log.Printf("Error opening file %s: %v", rawInputFilePng, err)
		t.FailNow()

	}
	defer inFile.Close()

	reader := bufio.NewReader(inFile)
	img, _, err := image.Decode(reader)
	if err != nil {
		log.Printf("Error decoding. %v", err)
		t.FailNow()
	}
	w := new(bytes.Buffer)
	err = Encode(w, img, bitmessage) // Encode the message into the image file
	if err != nil {
		log.Printf("Error Encoding file %v", err)
		t.FailNow()

	}
	outFile, err := os.Create(encodedInputFilePng)
	if err != nil {
		log.Printf("Error creating file %s: %v", encodedInputFilePng, err)
		t.FailNow()

	}
	w.WriteTo(outFile)
	defer outFile.Close()
}

func TestEncodeFromJpgFile(t *testing.T) {

	inFile, err := os.Open(rawInputFileJpg)
	if err != nil {
		log.Printf("Error opening file %s: %v", rawInputFileJpg, err)
		t.FailNow()

	}
	defer inFile.Close()

	reader := bufio.NewReader(inFile)
	img, err := jpeg.Decode(reader)
	if err != nil {
		log.Printf("Error decoding. %v", err)
		t.FailNow()
	}
	w := new(bytes.Buffer)
	err = Encode(w, img, bitmessage) // Encode the message into the image file
	if err != nil {
		log.Printf("Error Encoding file %v", err)
		t.FailNow()

	}
	outFile, err := os.Create(encodedInputFileJpg)
	if err != nil {
		log.Printf("Error creating file %s: %v", encodedInputFileJpg, err)
		t.FailNow()

	}
	w.WriteTo(outFile)
	defer outFile.Close()
}

func TestDecodeFromPngFile(t *testing.T) {
	inFile, err := os.Open(encodedInputFilePng)
	if err != nil {
		log.Printf("Error opening file %s: %v", encodedInputFilePng, err)
		t.FailNow()
	}
	defer inFile.Close()

	reader := bufio.NewReader(inFile)
	img, _, err := image.Decode(reader)
	if err != nil {
		log.Print("Error decoding file")
		t.FailNow()
	}

	sizeOfMessage := GetMessageSizeFromImage(img)

	msg := Decode(sizeOfMessage, img) // Read the message from the picture file

	if !bytes.Equal(msg, bitmessage) {
		log.Println(string(msg))
		log.Print("messages dont match")
		t.FailNow()
	}
}

func TestDecodeFromJpgFile(t *testing.T) {
	inFile, err := os.Open(encodedInputFileJpg)
	if err != nil {
		log.Printf("Error opening file %s: %v", encodedInputFileJpg, err)
		t.FailNow()
	}
	defer inFile.Close()

	reader := bufio.NewReader(inFile)
	img, _, err := image.Decode(reader)
	if err != nil {
		log.Print("Error decoding file")
		t.FailNow()
	}

	sizeOfMessage := GetMessageSizeFromImage(img)

	msg := Decode(sizeOfMessage, img) // Read the message from the picture file

	if !bytes.Equal(msg, bitmessage) {
		log.Println(string(msg))
		log.Print("messages dont match")
		t.FailNow()
	}
}

func TestEncodeDecodeGeneratedSmallImage(t *testing.T) {
	// Creating image
	width := 30
	height := 1

	upLeft := image.Point{0, 0}
	lowRight := image.Point{width, height}

	newimg := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// Set color for each pixel.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			newimg.Set(x, y, color.White)
		}
	}

	w := new(bytes.Buffer)
	err := EncodeRGBA(w, newimg, []uint8{84, 84, 84}) // Encode the message into the image file
	if err != nil {
		log.Printf("Error Encoding file %v", err)
		t.FailNow()

	}
	decodeImg, _, err := image.Decode(w)
	if err != nil {
		log.Println("Failed to Decode Image")
		t.FailNow()
	}

	sizeOfMessage := GetMessageSizeFromImage(decodeImg)

	msg := Decode(sizeOfMessage, decodeImg) // Read the message from the picture file

	// otherwise, print the message to STDOUT

	if !bytes.Equal(msg, []uint8{84, 84, 84}) {
		log.Println(string(msg))
		log.Print("messages dont match")
		t.FailNow()
	}
}
func TestSmalImage(t *testing.T) {

	miniImage := image.Image(image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{18, 1}}))
	log.Print(MaxEncodeSize(miniImage))
	if MaxEncodeSize(miniImage) > 0 {
		log.Printf("Uncaught small image size")
		t.FailNow()
	}

	miniImage = image.Image(image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{23, 1}}))

	if MaxEncodeSize(miniImage) != 4 {
		log.Printf("Uncaught minimal image size")
		t.FailNow()
	}
}

func TestMessageTooLarge(t *testing.T) {

	miniImage := image.Image(image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{24, 1}}))
	w := new(bytes.Buffer)
	err := Encode(w, miniImage, bitmessage) // Encode the message into the image file
	if err == nil {
		log.Printf("Uncaught error: message too large for image")
		t.FailNow()

	}

}

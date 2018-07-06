package steganography

import (
	"bufio"
	"bytes"
	"image"
	"log"
	"os"
	"testing"
)

var rawInputFile = "./examples/monalisa.png"
var encodedInputFile = "./examples/encoded.png"

// var message = "Hello Steganography !"
var bitmessage = []uint8{72, 101, 108, 108, 111, 32, 83, 116, 101, 103, 97, 110, 111, 103, 114, 97, 112, 104, 121, 32, 33, 10}

func TestEncode(t *testing.T) {

	inFile, err := os.Open(rawInputFile)
	if err != nil {
		log.Printf("Error opening file %s: %v", rawInputFile, err)
		t.FailNow()

	}
	defer inFile.Close()

	reader := bufio.NewReader(inFile)
	img, _, err := image.Decode(reader)
	// println(name)
	encodedImg := EncodeString(bitmessage, img) // Encode the message into the image file
	outFile, err := os.Create(encodedInputFile)
	if err != nil {
		log.Printf("Error creating file %s: %v", encodedInputFile, err)
		t.FailNow()

	}
	w := bufio.NewWriter(outFile)
	w.Write(encodedImg.Bytes())
}

func TestDecode(t *testing.T) {
	inFile, err := os.Open(encodedInputFile)
	if err != nil {
		log.Printf("Error opening file %s: %v", encodedInputFile, err)
		t.FailNow()
	}
	defer inFile.Close()

	reader := bufio.NewReader(inFile)
	img, _, err := image.Decode(reader)
	if err != nil {
		log.Print("Error decoding file")
		t.FailNow()
	}

	sizeOfMessage := GetSizeOfMessageFromImage(img)

	msg := DecodeMessageFromPicture(4, sizeOfMessage, img) // Read the message from the picture file

	// otherwise, print the message to STDOUT

	if !bytes.Equal(msg, bitmessage) {
		log.Print("messages dont match")
		t.FailNow()
	}
}

func TestEncodeDecode(t *testing.T) {
	inFile, err := os.Open(rawInputFile)
	if err != nil {
		log.Printf("Error opening file %s: %v", rawInputFile, err)
		t.FailNow()

	}
	defer inFile.Close()

	reader := bufio.NewReader(inFile)
	img, _, err := image.Decode(reader)
	// println(name)
	encodedImg := EncodeString(bitmessage, img) // Encode the message into the image file

	img, _, err = image.Decode(bytes.NewReader(encodedImg.Bytes()))
	if err != nil {
		log.Print("Error decoding file")
		t.FailNow()
	}

	sizeOfMessage := GetSizeOfMessageFromImage(img)

	msg := DecodeMessageFromPicture(4, sizeOfMessage, img) // Read the message from the picture file

	if !bytes.Equal(msg, bitmessage) {
		log.Print("messages dont match")
		t.FailNow()
	}
}

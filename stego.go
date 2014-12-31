package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"unicode/utf8"
)

var pictureInputFile string
var pictureOutputFile string
var messageInputFile string
var messageOutputFile string

var read bool
var write bool
var help bool

var ascii bool

func init() {

	flag.BoolVar(&read, "r", false, "Specifies if you would like to read a message from a given PNG file")
	flag.BoolVar(&write, "w", false, "Specifies if you would like to write a message to a given PNG file")

	flag.StringVar(&pictureInputFile, "imgi", "input.png", "Path to the the input image")
	flag.StringVar(&pictureOutputFile, "imgo", "output.png", "Path to the the output image")

	flag.StringVar(&messageInputFile, "msgi", "message.txt", "Path to the message input file")
	flag.StringVar(&messageOutputFile, "msgo", "", "Path to the message output file")

	flag.BoolVar(&help, "help", false, "Help")

	flag.BoolVar(&ascii, "ascii", false, "For use in read mode. Specifies if the anticipated message is in textual form.")

	flag.Parse()
}

func main() {

	if (!read && !write) || help {
		if help {
			fmt.Println("go-steg has two modes: write and read:")

			fmt.Println("- Write: take a message and write it into a specified location")
			fmt.Println("\t+ EX: ./stego -w -msgi message.txt -imgi plain.png -imgo secret.png")

			fmt.Println("- Read: take a picture and read the message from it")
			fmt.Println("\t+ EX: ./stego -r -imgi secret.png -msgo secret.txt")
		} else if !read || !write {
			fmt.Println("You must specify either the read or write flag. See -help for more information\n")
		}
		return
	}

	if write {
		message, err := ioutil.ReadFile(messageInputFile) // Read the message from the message file
		if err != nil {
			print("Error reading from file!!!")
			return
		}
		encodeString(string(message)) // Encode the message into the image file
	}

	if read {
		msg := decodeMessageFromPicture() // Read the message from the picture file

		// if the message is textual in nature eliminate excess non-ascii characters from the message
		if ascii == true {
			var lastIndexOfMsg int = len(msg)

			// iterate through every character in the message
			for i := range msg {
				// once a non-ascii character has been detected, set this as the last index of the message
				if msg[i] < 32 || 127 < msg[i] {
					lastIndexOfMsg = i
					break
				}
			}
			msg = msg[:lastIndexOfMsg] // truncate the message to eliminate all the garbage values off the end
		}

		// if the user specifies a location to write the message to...
		if messageOutputFile != "" {

			// write the message to the given output file
			err := ioutil.WriteFile(messageOutputFile, msg, 0644)

			if err != nil {
				fmt.Println("There was an error writing to file: ", messageOutputFile)
			}

		} else { // otherwise, print the message to STDOUT
			for i := range msg {
				fmt.Printf("%c", msg[i])
			}
		}

	}

}

// using LSB steganography, decode the message from the picture and return it as a sequence of bytes
func decodeMessageFromPicture() (message []byte) {

	var byteIndex int = 0
	var bitIndex int = 0

	rgbIm := imageToRGBA(decodeImage(pictureInputFile))

	width := rgbIm.Bounds().Dx()
	height := rgbIm.Bounds().Dy()

	var c color.RGBA
	var lsb byte

	message = append(message, 0)

	// iterate through every pixel in the image and stitch together the message bit by bit
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {

			c = rgbIm.RGBAAt(x, y) // get the color of the pixel

			/*  RED  */
			lsb = getLSB(c.R)                                                    // get the least significant bit from the red component of this pixel
			message[byteIndex] = setBitInByte(message[byteIndex], bitIndex, lsb) // add this bit to the message
			bitIndex++

			if bitIndex > 7 { // when we have filled up a byte, move on to the next byte
				bitIndex = 0
				byteIndex++
				message = append(message, 0)
			}

			/*  GREEN  */
			lsb = getLSB(c.G)
			message[byteIndex] = setBitInByte(message[byteIndex], bitIndex, lsb)
			bitIndex++

			if bitIndex > 7 {
				bitIndex = 0
				byteIndex++
				message = append(message, 0)
			}

			/*  BLUE  */
			lsb = getLSB(c.B)
			message[byteIndex] = setBitInByte(message[byteIndex], bitIndex, lsb)
			bitIndex++

			if bitIndex > 7 {
				bitIndex = 0
				byteIndex++
				message = append(message, 0)
			}
		}
	}
	return
}

// encodes a given string into the input image using least significant bit encryption
func encodeString(message string) {

	rgbIm := imageToRGBA(decodeImage(pictureInputFile))

	var messageLength int = utf8.RuneCountInString(message)
	var width int = rgbIm.Bounds().Dx()
	var height int = rgbIm.Bounds().Dy()

	if maxEncodeSize(rgbIm) < messageLength {
		print("Error! The message you are trying to encode is too large.")
		return
	}

	var c color.RGBA
	var offsetIntoMessage int = 0
	var bit byte
	var err error

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {

			c = rgbIm.RGBAAt(x, y) // get the color at this pixel

			/*  RED  */
			bit, err = getNextBitFromString(message)
			if err != nil { // if we don't have any more bits left in our message
				rgbIm.SetRGBA(x, y, c)
				encodePNG(pictureOutputFile, rgbIm) // write the encoded file out
				return
			}
			setLSB(&c.R, bit)

			/*  GREEN  */
			bit, err = getNextBitFromString(message)
			if err != nil {
				rgbIm.SetRGBA(x, y, c)
				encodePNG(pictureOutputFile, rgbIm)
				return
			}
			setLSB(&c.G, bit)

			/*  BLUE  */
			bit, err = getNextBitFromString(message)
			if err != nil {
				rgbIm.SetRGBA(x, y, c)
				encodePNG(pictureOutputFile, rgbIm)
				return
			}
			setLSB(&c.B, bit)

			rgbIm.SetRGBA(x, y, c)

			offsetIntoMessage++
		}
	}

	encodePNG(pictureOutputFile, rgbIm)
}

// convert given image to RGBA image
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

// read and return an image at the given path
func decodeImage(filename string) image.Image {
	inFile, err := os.Open(filename)

	if err != nil {
		log.Fatalf("Error opening file %s: %v", filename, err)
	}

	defer inFile.Close()
	reader := bufio.NewReader(inFile)
	img, _, err := image.Decode(reader)
	return img
}

// will write out a given image to a given path in filename
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

// given an image will find how many bytes can be stored in that image using least significant bit encoding
func maxEncodeSize(img image.Image) int {

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	return int(((width * height * 3) / 8))
}

// given a byte, will return the least significant bit of that byte
func getLSB(b byte) byte {
	if b%2 == 0 {
		return 0
	} else {
		return 1
	}

	return b
}

// given a byte will set that byte's least significant bit to a given value (where true is 1 and false is 0)
func setLSB(b *byte, bit byte) {
	if bit == 1 {
		*b = *b | 1
	} else if bit == 0 {
		var mask byte = 0xFE
		*b = *b & mask
	}
}

// given a bit will return a bit from that byte
func getBitFromByte(b byte, indexInByte int) byte {
	b = b << uint(indexInByte)
	var mask byte = 0x80

	var bit byte = mask & b

	if bit == 128 {
		return 1
	}
	return 0
}

// sets a specific bit in a byte to a given value and returns the new byte
func setBitInByte(b byte, indexInByte int, bit byte) byte {
	var mask byte = 0x80
	mask = mask >> uint(indexInByte)

	if bit == 0 {
		mask = ^mask
		b = b & mask
	} else if bit == 1 {
		b = b | mask
	}
	return b
}

var offsetInBytes int = 0
var offsetInBitsIntoByte int = 0

// each call will return the next subsequent bit in the string
func getNextBitFromString(s string) (byte, error) {

	lenOfString := len(s)

	if offsetInBytes >= lenOfString {
		return 0, errors.New("Error! Can't offset that far into the string.")
	}

	byteArray := []byte(s)
	choiceByte := byteArray[offsetInBytes]
	choiceBit := getBitFromByte(choiceByte, offsetInBitsIntoByte)

	offsetInBitsIntoByte++

	if offsetInBitsIntoByte >= 8 {
		offsetInBitsIntoByte = 0
		offsetInBytes++
	}

	return choiceBit, nil
}
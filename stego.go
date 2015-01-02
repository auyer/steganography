package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"
)

var pictureInputFile string
var pictureOutputFile string
var messageInputFile string
var messageOutputFile string
var printLen bool
var read bool
var write bool
var help bool

func init() {

	flag.BoolVar(&read, "r", false, "Specifies if you would like to read a message from a given PNG file")
	flag.BoolVar(&write, "w", false, "Specifies if you would like to write a message to a given PNG file")
	flag.BoolVar(&printLen, "length", false, "When set, will print out the max message size that can fit into given image.")

	flag.StringVar(&pictureInputFile, "imgi", "", "Path to the the input image")
	flag.StringVar(&pictureOutputFile, "imgo", "", "Path to the the output image")

	flag.StringVar(&messageInputFile, "msgi", "", "Path to the message input file")
	flag.StringVar(&messageOutputFile, "msgo", "", "Path to the message output file")

	flag.BoolVar(&help, "help", false, "Help")

	flag.Parse()
}

func main() {
	if (!read && !write && !printLen) || help {
		if help {
			fmt.Println("go-steg has two modes: write and read:")

			fmt.Println("- Write: take a message and write it into a specified location")
			fmt.Println("\t+ EX: ./stego -w -msgi message.txt -imgi plain.png -imgo secret.png")

			fmt.Println("- Read: take a picture and read the message from it")
			fmt.Println("\t+ EX: ./stego -r -imgi secret.png -msgo secret.txt")
		} else if !read || !write {
			fmt.Println("You must specify either the read, write, or length flag. See -help for more information\n")
		}
		return
	}

	if printLen {

		if pictureInputFile == "" {
			fmt.Println("Error: In order to run stego in length mode, you must specify: ")
			fmt.Println("-imgi: the image that you would like to check the maximum encoding length of")
			return
		}

		rgbIm := imageToRGBA(decodeImage(pictureInputFile))

		var sizeInBytes uint32 = maxEncodeSize(rgbIm)

		fmt.Println("B:\t", sizeInBytes)
		fmt.Println("KB:\t", float64(sizeInBytes)/1000)
		fmt.Println("MB:\t", (float64(sizeInBytes)/1000)/1000)
	}

	if write {

		if messageInputFile == "" || pictureInputFile == "" || pictureOutputFile == "" {
			fmt.Println("Error: In order to run stego in write mode, you must specify: ")
			fmt.Println("-imgi: the plain image that you would like to encode with")
			fmt.Println("-imgo: where you would like to store the encoded image")
			fmt.Println("-msgi: the message that you would like to embed in the image")
			return
		}

		message, err := ioutil.ReadFile(messageInputFile) // Read the message from the message file
		if err != nil {
			print("Error reading from file!!!")
			return
		}

		encodeString(message) // Encode the message into the image file
	}

	if read {

		if pictureInputFile == "" {
			fmt.Println("Error: In order to run stego in read mode, you must specify: ")
			fmt.Println("-imgi: the image with the embeded message")
			return
		}

		sizeOfMessage := getSizeOfMessageFromImage()

		msg := decodeMessageFromPicture(4, sizeOfMessage) // Read the message from the picture file

		// if the user specifies a location to write the message to...
		if messageOutputFile != "" {
			err := ioutil.WriteFile(messageOutputFile, msg, 0644) // write the message to the given output file

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

// encodes a given string into the input image using least significant bit encryption
func encodeString(message []byte) {

	rgbIm := imageToRGBA(decodeImage(pictureInputFile))

	var messageLength uint32 = uint32(len(message))

	var width int = rgbIm.Bounds().Dx()
	var height int = rgbIm.Bounds().Dy()
	var c color.RGBA
	var bit byte
	var ok bool

	if maxEncodeSize(rgbIm) < messageLength+4 {
		print("Error! The message you are trying to encode is too large.")
		return
	}

	one, two, three, four := splitToBytes(messageLength)

	message = append([]byte{four}, message...)
	message = append([]byte{three}, message...)
	message = append([]byte{two}, message...)
	message = append([]byte{one}, message...)

	ch := make(chan byte, 100)

	go getNextBitFromString(message, ch)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {

			c = rgbIm.RGBAAt(x, y) // get the color at this pixel

			/*  RED  */
			bit, ok = <-ch
			if !ok { // if we don't have any more bits left in our message
				rgbIm.SetRGBA(x, y, c)
				encodePNG(pictureOutputFile, rgbIm) // write the encoded file out
				return
			}
			setLSB(&c.R, bit)

			/*  GREEN  */
			bit, ok = <-ch
			if !ok {
				rgbIm.SetRGBA(x, y, c)
				encodePNG(pictureOutputFile, rgbIm)
				return
			}
			setLSB(&c.G, bit)

			/*  BLUE  */
			bit, ok = <-ch
			if !ok {
				rgbIm.SetRGBA(x, y, c)
				encodePNG(pictureOutputFile, rgbIm)
				return
			}
			setLSB(&c.B, bit)

			rgbIm.SetRGBA(x, y, c)
		}
	}

	encodePNG(pictureOutputFile, rgbIm)
}

// using LSB steganography, decode the message from the picture and return it as a sequence of bytes
func decodeMessageFromPicture(startOffset uint32, msgLen uint32) (message []byte) {

	var byteIndex uint32 = 0
	var bitIndex uint32 = 0

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

				if byteIndex >= msgLen+startOffset {
					return message[startOffset : msgLen+startOffset]
				}

				message = append(message, 0)
			}

			/*  GREEN  */
			lsb = getLSB(c.G)
			message[byteIndex] = setBitInByte(message[byteIndex], bitIndex, lsb)
			bitIndex++

			if bitIndex > 7 {

				bitIndex = 0
				byteIndex++

				if byteIndex >= msgLen+startOffset {
					return message[startOffset : msgLen+startOffset]
				}

				message = append(message, 0)
			}

			/*  BLUE  */
			lsb = getLSB(c.B)
			message[byteIndex] = setBitInByte(message[byteIndex], bitIndex, lsb)
			bitIndex++

			if bitIndex > 7 {
				bitIndex = 0
				byteIndex++

				if byteIndex >= msgLen+startOffset {
					return message[startOffset : msgLen+startOffset]
				}

				message = append(message, 0)
			}
		}
	}
	return
}

// gets the size of the message from the first four bytes encoded in the image
func getSizeOfMessageFromImage() (size uint32) {

	sizeAsByteArray := decodeMessageFromPicture(0, 4)
	size = combineToInt(sizeAsByteArray[0], sizeAsByteArray[1], sizeAsByteArray[2], sizeAsByteArray[3])
	return
}

// given an image will find how many bytes can be stored in that image using least significant bit encoding
func maxEncodeSize(img image.Image) uint32 {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	return uint32(((width * height * 3) / 8)) - 4
}

// each call will return the next subsequent bit in the string
func getNextBitFromString(byteArray []byte, ch chan byte) {

	var offsetInBytes int = 0
	var offsetInBitsIntoByte int = 0
	var choiceByte byte

	lenOfString := len(byteArray)

	for {
		if offsetInBytes >= lenOfString {
			close(ch)
			return
		}

		choiceByte = byteArray[offsetInBytes]
		ch <- getBitFromByte(choiceByte, offsetInBitsIntoByte)

		offsetInBitsIntoByte++

		if offsetInBitsIntoByte >= 8 {
			offsetInBitsIntoByte = 0
			offsetInBytes++
		}
	}
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
func setBitInByte(b byte, indexInByte uint32, bit byte) byte {
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

// given four bytes, will return the 32 bit unsigned integer which is the composition of those four bytes (one is MSB)
func combineToInt(one, two, three, four byte) (ret uint32) {
	ret = uint32(one)
	ret = ret << 8
	ret = ret | uint32(two)
	ret = ret << 8
	ret = ret | uint32(three)
	ret = ret << 8
	ret = ret | uint32(four)
	return
}

// given an unsigned integer, will split this integer into its four bytes
func splitToBytes(x uint32) (one, two, three, four byte) {
	one = byte(x >> 24)
	var mask uint32 = 255

	two = byte((x >> 16) & mask)
	three = byte((x >> 8) & mask)
	four = byte(x & mask)
	return
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

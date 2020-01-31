// Package steganography is a library that provides functions to execute steganography encoding and decoding in a given image. It is also able to check the maximum encoding size, and the size of an encoded message.
package steganography

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
)

// EncodeRGBA encodes a given string into the input image using least significant bit encryption (LSB steganography)
// The minnimum image size is 24 pixels for one byte. For each additional byte, it is necessary 3 more pixels.
/*
	Input:
		writeBuffer *bytes.Buffer : the destination of the encoded image bytes
		pictureInputFile image.RGBA : image data used in encoding
		message []byte : byte slice of the message to be encoded
	Output:
		bytes buffer ( io.writter ) to create file, or send data.
*/
func EncodeRGBA(writeBuffer *bytes.Buffer, rgbImage *image.RGBA, message []byte) error {

	var messageLength = uint32(len(message))

	var width = rgbImage.Bounds().Dx()
	var height = rgbImage.Bounds().Dy()
	var c color.RGBA
	var bit byte
	var ok bool
	//var encodedImage image.Image
	if MaxEncodeSize(rgbImage) < messageLength+4 {
		return errors.New("message too large for image")
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

			c = rgbImage.RGBAAt(x, y) // get the color at this pixel

			/*  RED  */
			bit, ok = <-ch
			if !ok { // if we don't have any more bits left in our message
				rgbImage.SetRGBA(x, y, c)
				png.Encode(writeBuffer, rgbImage)
				// return *writeBuffer, nil
			}
			setLSB(&c.R, bit)

			/*  GREEN  */
			bit, ok = <-ch
			if !ok {
				rgbImage.SetRGBA(x, y, c)
				png.Encode(writeBuffer, rgbImage)
				return nil
			}
			setLSB(&c.G, bit)

			/*  BLUE  */
			bit, ok = <-ch
			if !ok {
				rgbImage.SetRGBA(x, y, c)
				png.Encode(writeBuffer, rgbImage)
				return nil
			}
			setLSB(&c.B, bit)

			rgbImage.SetRGBA(x, y, c)
		}
	}

	err := png.Encode(writeBuffer, rgbImage)
	fmt.Println("err")
	return err
}

// Encode encodes a given string into the input image using least significant bit encryption (LSB steganography)
// The minnimum image size is 23 pixels
// It wraps EncodeRGBA making the conversion from image.Image to image.RGBA
/*
	Input:
		writeBuffer *bytes.Buffer : the destination of the encoded image bytes
		message []byte : byte slice of the message to be encoded
		pictureInputFile image.Image : image data used in encoding
	Output:
		bytes buffer ( io.writter ) to create file, or send data.
*/
func Encode(writeBuffer *bytes.Buffer, pictureInputFile image.Image, message []byte) error {

	rgbImage := imageToRGBA(pictureInputFile)

	return EncodeRGBA(writeBuffer, rgbImage, message)

}

// decodeRGBA gets messages from pictures using LSB steganography, decode the message from the picture and return it as a sequence of bytes
/*
	Input:
		startOffset uint32 : number of bytes used to declare size of message
		msgLen uint32 : size of the message to be decoded
		pictureInputFile image.RGBA : image data used in decoding
	Output:
		message []byte decoded from image
*/
func decodeRGBA(startOffset uint32, msgLen uint32, rgbImage *image.RGBA) (message []byte) {

	var byteIndex uint32
	var bitIndex uint32

	width := rgbImage.Bounds().Dx()
	height := rgbImage.Bounds().Dy()

	var c color.RGBA
	var lsb byte

	message = append(message, 0)

	// iterate through every pixel in the image and stitch together the message bit by bit
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {

			c = rgbImage.RGBAAt(x, y) // get the color of the pixel

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

// decode gets messages from pictures using LSB steganography, decode the message from the picture and return it as a sequence of bytes
// It wraps EncodeRGBA making the conversion from image.Image to image.RGBA
/*
	Input:
		startOffset uint32 : number of bytes used to declare size of message
		msgLen uint32 : size of the message to be decoded
		pictureInputFile image.Image : image data used in decoding
	Output:
		message []byte decoded from image
*/
func decode(startOffset uint32, msgLen uint32, pictureInputFile image.Image) (message []byte) {

	rgbImage := imageToRGBA(pictureInputFile)
	return decodeRGBA(startOffset, msgLen, rgbImage)

}

// Decode gets messages from pictures using LSB steganography, decode the message from the picture and return it as a sequence of bytes
// It wraps EncodeRGBA making the conversion from image.Image to image.RGBA
/*
	Input:
		msgLen uint32 : size of the message to be decoded
		pictureInputFile image.Image : image data used in decoding
	Output:
		message []byte decoded from image
*/
func Decode(msgLen uint32, pictureInputFile image.Image) (message []byte) {
	return decode(4, msgLen, pictureInputFile) // the offset of 4 skips the "header" where message length is defined

}

// MaxEncodeSize given an image will find how many bytes can be stored in that image using least significant bit encoding
// ((width * height * 3) / 8 ) - 4
// The result must be at least 4,
func MaxEncodeSize(img image.Image) uint32 {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	eval := ((width * height * 3) / 8) - 4
	if eval < 4 {
		eval = 0
	}
	return uint32(eval)
}

// GetMessageSizeFromImage gets the size of the message from the first four bytes encoded in the image
func GetMessageSizeFromImage(pictureInputFile image.Image) (size uint32) {

	sizeAsByteArray := decode(0, 4, pictureInputFile)
	size = combineToInt(sizeAsByteArray[0], sizeAsByteArray[1], sizeAsByteArray[2], sizeAsByteArray[3])
	return
}

// getNextBitFromString each call will return the next subsequent bit in the string
func getNextBitFromString(byteArray []byte, ch chan byte) {

	var offsetInBytes int
	var offsetInBitsIntoByte int
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

// getLSB given a byte, will return the least significant bit of that byte
func getLSB(b byte) byte {
	if b%2 == 0 {
		return 0
	}
	return 1
}

// setLSB given a byte will set that byte's least significant bit to a given value (where true is 1 and false is 0)
func setLSB(b *byte, bit byte) {
	if bit == 1 {
		*b = *b | 1
	} else if bit == 0 {
		var mask byte = 0xFE
		*b = *b & mask
	}
}

// getBitFromByte given a bit will return a bit from that byte
func getBitFromByte(b byte, indexInByte int) byte {
	b = b << uint(indexInByte)
	var mask byte = 0x80

	var bit = mask & b

	if bit == 128 {
		return 1
	}
	return 0
}

// setBitInByte sets a specific bit in a byte to a given value and returns the new byte
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

// combineToInt given four bytes, will return the 32 bit unsigned integer which is the composition of those four bytes (one is MSB)
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

// splitToBytes given an unsigned integer, will split this integer into its four bytes
func splitToBytes(x uint32) (one, two, three, four byte) {
	one = byte(x >> 24)
	var mask uint32 = 255

	two = byte((x >> 16) & mask)
	three = byte((x >> 8) & mask)
	four = byte(x & mask)
	return
}

// imageToRGBA converts image.Image to image.RGBA
func imageToRGBA(src image.Image) *image.RGBA {
	bounds := src.Bounds()

	var m *image.RGBA
	var width, height int

	width = bounds.Dx()
	height = bounds.Dy()

	m = image.NewRGBA(image.Rect(0, 0, width, height))

	draw.Draw(m, m.Bounds(), src, bounds.Min, draw.Src)
	return m
}

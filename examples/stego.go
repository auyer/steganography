package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"os"

	"github.com/auyer/steganography"
)

var pictureInputFile string
var pictureOutputFile string
var messageInputFile string
var messageOutputFile string
var decode bool
var encode bool
var help bool

func init() {

	flag.BoolVar(&decode, "d", false, "Specifies if you would like to decode a message from a given PNG file")
	flag.BoolVar(&encode, "e", false, "Specifies if you would like to encode a message to a given PNG file")

	flag.StringVar(&pictureInputFile, "i", "", "Path to the the input image")
	flag.StringVar(&pictureOutputFile, "o", "", "Path to the the output image")

	flag.StringVar(&messageInputFile, "mi", "", "Path to the message input file")
	flag.StringVar(&messageOutputFile, "mo", "", "Path to the message output file")

	flag.BoolVar(&help, "help", false, "Help")

	flag.Parse()
}

func main() {
	if (!decode && !encode) || help {
		if help {

			fmt.Println("-e: take a message and encodes it into a specified location")
			fmt.Println("\t+ EX: ./main -e -mi message.txt -i plain.png -o secret.png")

			fmt.Println("-d: take a picture and decodes the message from it")
			fmt.Println("\t+ EX: ./stego -d -i secret.png -mo secret.txt")
		} else if !decode || !encode {
			fmt.Println("You must specify either the decode, write, flag. See -help for more information\n")
		}
		return
	}

	if encode {
		if messageInputFile == "" || pictureInputFile == "" || pictureOutputFile == "" {
			fmt.Println("Error: In order to run in encode mode, you must specify: ")
			fmt.Println("-i: the plain image that you would like to encode with")
			fmt.Println("-o: where you would like to store the encodeded image")
			fmt.Println("-mi: the message that you would like to embed in the image")
			return
		}
		message, err := ioutil.ReadFile(messageInputFile) // Read the message from the message file
		if err != nil {
			print("Error reading from file!!!")
			return
		}

		inFile, err := os.Open(pictureInputFile)
		if err != nil {
			log.Fatalf("Error opening file %s: %v", pictureInputFile, err)
		}
		defer inFile.Close()

		reader := bufio.NewReader(inFile)
		img, name, err := image.Decode(reader)
		println(name)
		encodedImg := steganography.EncodeString(message, img) // Encode the message into the image file
		outFile, err := os.Create(pictureOutputFile)
		if err != nil {
			log.Fatalf("Error creating file %s: %v", pictureOutputFile, err)
		}
		w := bufio.NewWriter(outFile)
		w.Write(encodedImg.Bytes())

	}

	if decode {

		if pictureInputFile == "" {
			fmt.Println("Error: In order to run stego in read mode, you must specify: -im: the image with the embedded message")
			return
		}

		inFile, err := os.Open(pictureInputFile)
		if err != nil {
			log.Fatalf("Error opening file %s: %v", pictureInputFile, err)
		}
		defer inFile.Close()

		reader := bufio.NewReader(inFile)
		img, _, err := image.Decode(reader)
		if err != nil {
			log.Fatalf("error decoding file", img)
		}

		sizeOfMessage := steganography.GetSizeOfMessageFromImage(img)

		msg := steganography.DecodeMessageFromPicture(4, sizeOfMessage, img) // Read the message from the picture file

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

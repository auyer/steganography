# steganography Lib

[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](http://godoc.org/github.com/auyer/steganography) [![Go Report Card](https://goreportcard.com/badge/github.com/auyer/steganography)](https://goreportcard.com/report/github.com/auyer/steganography) [![LICENSE MIT](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://img.shields.io/badge/license-MIT-brightgreen.svg) [![Build Status](https://travis-ci.org/auyer/steganography.svg?branch=master)](https://travis-ci.org/auyer/steganography)

Steganography is a library written in Pure go to allow simple LSB steganography on images. It is capable of both encoding and decoding images. It can store files of any format.
This librery is inspired by Stego, a command line utility with the same purpose.

## Demonstration

| Original        | Encoded           |
| -------------------- | ------------------|
| ![Original File](examples/stegosaurus.png) | ![Encoded File](examples/encoded_stegosaurus.png)
|   79,955 bytes       |   80,629 bytes   |

The second image contains the first paragaph of the description of a stegosaurus on [Wikipidia](https://en.wikipedia.org/wiki/Stegosaurus), also available in [examples/message.txt](examples/message.txt) as an example.

Encode
------
Write mode is used to take a message and embed it into an image file using LSB steganography in order to produce a secret image file that will contain your message.
```go
encodedImg := steganography.EncodeString(message, img, pictureOutputFile) // Encode the message into the image file
outFile, err := os.Create(pictureOutputFile)
if err != nil {
    log.Fatalf("Error creating file %s: %v", pictureOutputFile, err)
}
w := bufio.NewWriter(outFile)
w.Write(encodedImg.Bytes())
```

Decode
-----
Read mode is used to read an image that has been encoded using LSB steganography, and extract the hidden message from that image.
```go
reader := bufio.NewReader(inFile)
img, _, err := image.Decode(reader)
sizeOfMessage := steganography.GetSizeOfMessageFromImage(img)

msg := steganography.DecodeMessageFromPicture(4, sizeOfMessage, img) // Read the message from the picture file
```

Size of Message
------
Length mode can be used in order to preform a preliminary check on the carrier image in order to deduce how large of a file it can store. Length is given in bytes, kilobytes, and megabytes.

```go
img, _, err := image.Decode(reader)
sizeOfMessage := steganography.GetSizeOfMessageFromImage(img)
```
-----
### Attributions
 - mStegossaurus Picture By Matt Martyniuk - Own work, CC BY-SA 3.0, https://commons.wikimedia.org/w/index.php?curid=42215661

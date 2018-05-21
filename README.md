steganography Lib
=====

Stego is a command line utility written in go to allow simple LSB steganography on PNG images. It is capable of both encoding and decoding images. It can store files of any format.

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
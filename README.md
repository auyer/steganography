Stego
=====

Stego is a command line utility written in go to allow simple LSB steganography on PNG images. It is capable of both encoding and decoding images:



Write
------
Write mode is used to take a message and embed it into an image file using LSB steganography in order to produce a secret image file that will contain your message.
```bash
./stego -w -msgi message.txt -imgi plain.png -imgo secret.png
```
Read
-----
Read mode is used to read an image that has been encoded using LSB steganography, and extract the hidden message from that image.
```bash
./stego -r -imgi secret.png -msgo secret.txt
```
If no 'msgo' is specified, the message will be printed to stdout:
```bash
./stego -r -imgi secret.png
```
If the encoded message is comprised of ASCII characters, you may use the ASCII flag in order to discard all non-ascii characters at the end of the message:
```bash
./stego -r -imgi secret.png -msgo secret.txt -ascii
```

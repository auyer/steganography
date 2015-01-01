Stego
=====

Stego is a command line utility written in go to allow simple LSB steganography on PNG images. It is capable of both encoding and decoding images. It can store files of any format.

Length
------
Length mode can be used in order to preform a preliminary check on the carrier image in order to deduce how large of a file it can store. Length is given in bytes, kilobytes, and megabytes.

```bash
› ./stego -length -imgi plain.png
B:	 863996
KB:	 863.996
MB:	 0.863996
```

Write
------
Write mode is used to take a message and embed it into an image file using LSB steganography in order to produce a secret image file that will contain your message.
```bash
› ./stego -w -msgi message.txt -imgi plain.png -imgo steged.png
› ./stego -w -msgi hideme.png -imgi plain.png -imgo steged.png
```
Read
-----
Read mode is used to read an image that has been encoded using LSB steganography, and extract the hidden message from that image.
```bash
› ./stego -r -imgi steged.png -msgo secret.txt
› ./stego -r -imgi steged.png -msgo secret.png
```
If no 'msgo' is specified, the message will be printed to stdout:
```bash
› ./stego -r -imgi secret.png
```

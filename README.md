Stego
=====

Stego is a command line utility written in go to allow simple LSB steganography on PNG images. It is capable of both encoding and decoding images.

Write:
------
./stego -w -msgi message.txt -imgi plain.png -imgo secret.png

Read:
------
./stego -r -imgi secret.png -msgo secret.txt

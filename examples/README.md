Example usage: 

    Decoding message: go run stego.go -d -i encoded_stegosaurus.png

    Encoding message: go run stego.go -e -i stegosaurus.png -mi message.txt -o encoded_stegosaurus.png

Usage stego.go

    -help 	Will show this message below

    -d	Specifies if you would like to decode a message from a given PNG file

    -e 	Specifies if you would like to encode a message to a given PNG file

    -i string Path to the the input image

    -mi string Path to the message input file
    
    -mo string Path to the message output file
    
    -o string Path to the the output image (default "encoded.png")
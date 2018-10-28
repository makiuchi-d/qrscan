package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

func main() {
	flag.Parse()
	args := flag.Args()

	var reader io.Reader

	if len(args) <= 0 {
		reader = os.Stdin
	} else {
		filename := args[0]
		file, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		reader = file
	}

	img, _, err := image.Decode(reader)
	if err != nil {
		panic(err)
	}

	qrreader := qrcode.NewQRCodeReader()
	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		panic(err)
	}
	result, err := qrreader.Decode(bmp, nil)
	if err != nil {
		panic(err)
	}

	fmt.Print(result)
}

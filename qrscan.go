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
	"github.com/makiuchi-d/gozxing/multi/qrcode"
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

	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}

	qrreader := qrcode.NewQRCodeMultiReader()
	results, err := qrreader.DecodeMultipleWithoutHint(bmp)
	if err != nil {
		panic(fmt.Sprintf("%+v", err))
	}

	for _, result := range results {
		fmt.Println(result)
	}
}

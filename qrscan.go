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
	"github.com/makiuchi-d/gozxing/common"
	"github.com/makiuchi-d/gozxing/qrcode"
)

type ImageSource struct {
	gozxing.LuminanceSourceBase
	img       image.Image
	top, left int
}

func NewImageSource(img image.Image) gozxing.LuminanceSource {
	rect := img.Bounds()
	top := rect.Min.Y
	left := rect.Min.X
	width := rect.Max.X - rect.Min.X
	height := rect.Max.Y - rect.Min.Y
	return &ImageSource{
		gozxing.LuminanceSourceBase{width, height},
		img,
		top,
		left,
	}
}
func (this *ImageSource) GetRow(y int, row []byte) ([]byte, error) {
	for x := 0; x < this.GetWidth(); x++ {
		r, g, b, _ := this.img.At(this.left+x, this.top+y).RGBA()
		row[x] = byte((r + 2*g + b) * 255 / (4 * 0xffff))
	}
	return row, nil
}
func (this *ImageSource) GetMatrix() []byte {
	width := this.GetWidth()
	height := this.GetHeight()
	matrix := make([]byte, width*height)
	for y := 0; y < height; y++ {
		offset := y * width
		for x := 0; x < width; x++ {
			r, g, b, _ := this.img.At(this.left+x, this.top+y).RGBA()
			matrix[offset+x] = byte((r + 2*g + b) * 255 / (4 * 0xffff))
		}
	}
	return matrix
}
func (this *ImageSource) Invert() gozxing.LuminanceSource {
	return gozxing.LuminanceSourceInvert(this)
}
func (this *ImageSource) String() string {
	return gozxing.LuminanceSourceString(this)
}

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
	bmp, err := gozxing.NewBinaryBitmap(common.NewHybridBinarizer(NewImageSource(img)))
	if err != nil {
		panic(err)
	}
	result, err := qrreader.Decode(bmp, nil)
	if err != nil {
		panic(err)
	}

	fmt.Print(result)
}

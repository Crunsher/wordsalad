package canvas

import (
	"bufio"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/golang/freetype"
)

var size = 12.0

func PaintImage(field Field, fontPath string) (*image.RGBA, error) {
	fontBytes, err := ioutil.ReadFile(fontPath)
	if err != nil {
		return nil, err
	}
	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}

	rgba := image.NewRGBA(image.Rect(0, 0, int(size)*field.width+40, int(size)*field.height+40))
	draw.Draw(rgba, rgba.Bounds(), image.White, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(f)
	c.SetFontSize(size)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(image.Black)

	for x, line := range field.field {
		for y, elem := range line {
			pt := freetype.Pt(int(size) * y + 20, x * int(size) + 20)
			_, err = c.DrawString(string(elem.c), pt)
			if err != nil {
				return nil, err
			}
		}
	}

	return rgba, nil
}

func WriteToFile(image *image.RGBA, filename string) error {
	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	err = png.Encode(b, image)
	if err != nil {
		return err
	}
	err = b.Flush()
	if err != nil {
		return err
	}
	fmt.Printf("Wrote \"%s\" OK.\n", filename)
	return nil
}

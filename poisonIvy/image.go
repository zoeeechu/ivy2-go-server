// Special Thanks: https://github.com/dtgreene/ivy2
package poisonIvy

import (
	"bytes"
	"image"
	"image/jpeg"
	_ "image/png"
	"os"

	"github.com/nfnt/resize"
)

const (
	PrintFinalWidth  = 640
	PrintFinalHeight = 1616
)

func PrepareImage(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	resized := resize.Resize(PrintFinalWidth, PrintFinalHeight, img, resize.Lanczos3)

	buf := new(bytes.Buffer)

	err = jpeg.Encode(buf, resized, &jpeg.Options{Quality: 100})
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

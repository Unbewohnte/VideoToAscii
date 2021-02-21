package processor

import (
	"image"
	_ "image/jpeg"
	"image/png"

	_ "image/png"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
)

type DataForAscii struct {
	Img      *image.Image
	Width    uint
	Height   uint
	Filename string
}

// GetImage returns image.Image from a filepath
func GetImage(pathToFile string) (*image.Image, error) {
	file, err := os.Open(filepath.Join(pathToFile))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	image, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return &image, nil
}

// SaveImage takes an image.Image and saves it to a file
func SaveImage(filename string, img *image.Image) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	err = png.Encode(f, *img)
	if err != nil {
		return err
	}
	return nil
}

// ResizeImage takes an image.Image and returns a resized one using https://github.com/nfnt/resize
func ResizeImage(img image.Image, newWidth uint, newHeight uint) image.Image {
	resizedImage := resize.Resize(newWidth, newHeight, img, resize.Lanczos3)
	return resizedImage
}

// GetChar returns a char from chars corresponding to the pixel brightness
func GetChar(chars []string, pixelBrightness int) string {
	charsLen := len(chars)
	return chars[int((charsLen*pixelBrightness)/256)]
}

// ASCIIfy converts and image.Image into ASCII art
func ASCIIfy(ASCIIchars []string, img *image.Image, cols, rows uint, filename string) {

	var resized image.Image
	if cols == uint(0) || rows == uint(0) {
		resized = *img
	} else {
		resized = ResizeImage(*img, cols, rows)
	}

	imgBounds := resized.Bounds()

	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	for y := 0; y < imgBounds.Max.Y; y++ {
		for x := 0; x < imgBounds.Max.X; x++ {
			r, g, b, _ := resized.At(x, y).RGBA()
			r = r / 257
			g = g / 257
			b = b / 257
			currentPixelBrightness := int((float64(0.2126)*float64(r) + float64(0.7152)*float64(g) + float64(0.0722)*float64(b)))
			f.Write([]byte(GetChar(ASCIIchars, currentPixelBrightness)))
		}
		f.Write([]byte("\n"))
	}
}

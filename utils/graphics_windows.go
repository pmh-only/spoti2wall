package utils

import (
	"image"
	"image/jpeg"
	"os"

	"github.com/disintegration/gift"
)

func BlurImage(path string, blurOption int) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	g := gift.New(gift.GaussianBlur(float32(blurOption)))
	dst := image.NewRGBA(g.Bounds(img.Bounds()))
	g.Draw(dst, img)

	outFile, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	jpeg.Encode(outFile, dst, nil)
}

func DarkImage(path string, darkOption int) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	g := gift.New(gift.Brightness(-float32(darkOption)))

	dst := image.NewRGBA(g.Bounds(img.Bounds()))
	g.Draw(dst, img)

	outFile, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	jpeg.Encode(outFile, dst, nil)
}

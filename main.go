package main

import (
	"image"
	"image/draw"
	"image/jpeg"
	_ "image/png"
	"log"
	"math"
	"os"
)

const (
	xInterval = 500
	yInterval = 500
)

func main() {
	originalFile, err := os.Open("./assets/torres.jpg")
	if err != nil {
		log.Fatalf("failed to open original: %s", err)
	}
	defer originalFile.Close()

	stackFile, err := os.Open("./assets/liverpool.jpg")
	if err != nil {
		log.Fatalf("failed to open stack: %s", err)
	}
	defer stackFile.Close()

	originalImg, _, err := image.Decode(originalFile)
	if err != nil {
		log.Fatalf("failed to decode original: %s", err)
	}
	stackImg, _, err := image.Decode(stackFile)
	if err != nil {
		log.Fatalf("failed to decode stack: %s", err)
	}

	rcts := make([]image.Rectangle, 0)
	originalWidth := originalImg.Bounds().Dx()
	originalHeight := originalImg.Bounds().Dy()
	blocks := calcBlocks(originalWidth, originalHeight)
	column := 0
	x := 0
	y := 0
	for i := 0; i < int(blocks); i++ {
		x = column * xInterval
		if x > originalWidth {
			column = 0
			x = 0
			y += yInterval
		}
		startPoint := image.Point{X: x, Y: y}
		rcts = append(rcts, image.Rectangle{
			Min: startPoint,
			Max: startPoint.Add(stackImg.Bounds().Size()),
		})
		column++
	}

	originalRct := image.Rectangle{
		Min: image.ZP,
		Max: originalImg.Bounds().Size(),
	}
	rgba := image.NewRGBA(originalRct)
	draw.Draw(rgba, originalRct, originalImg, image.ZP, draw.Src)
	for _, v := range rcts {
		draw.Draw(rgba, v, stackImg, image.ZP, draw.Over)
	}

	output, err := os.Create("./output.jpg")
	if err != nil {
		log.Fatalf("failed to create result image: %s", err)
	}
	defer output.Close()

	jpeg.Encode(output, rgba, &jpeg.Options{
		Quality: jpeg.DefaultQuality,
	})
}

func calcBlocks(w, h int) float64 {
	return math.Ceil(float64(w)/float64(xInterval)) * math.Ceil(float64(h)/float64(yInterval))
}

package main

import (
	"flag"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"
)

func main() {
	small := flag.String("small", "", "64x64 PNG output path")
	large := flag.String("large", "", "256x256 PNG output path")
	flag.Parse()
	if *small == "" || *large == "" {
		log.Fatal("both -small and -large are required")
	}
	writeIcon(*small, 64)
	writeIcon(*large, 256)
}

func writeIcon(path string, size int) {
	canvas := image.NewRGBA(image.Rect(0, 0, size, size))
	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{C: color.RGBA{R: 19, G: 78, B: 130, A: 255}}, image.Point{}, draw.Src)
	margin := size / 8
	panel := image.Rect(margin, margin, size-margin, size-margin)
	draw.Draw(canvas, panel, &image.Uniform{C: color.RGBA{R: 45, G: 180, B: 176, A: 255}}, image.Point{}, draw.Src)
	stripe := image.Rect(margin, size/2-size/16, size-margin, size/2+size/16)
	draw.Draw(canvas, stripe, &image.Uniform{C: color.RGBA{R: 255, G: 255, B: 255, A: 255}}, image.Point{}, draw.Src)

	file, err := os.Create(path)
	if err != nil {
		log.Fatalf("create %s: %v", path, err)
	}
	defer file.Close()
	if err := png.Encode(file, canvas); err != nil {
		log.Fatalf("encode %s: %v", path, err)
	}
}

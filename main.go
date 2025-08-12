package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
)

var asciiChars = []byte(" .:-=+*#%@")

func main() {
	filenameWithoutExt := "input"
	dir := "."

	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	var foundFile string
	for _, entry := range entries {
		if !entry.IsDir() {
			name := entry.Name()
			ext := filepath.Ext(name)
			base := name[:len(name)-len(ext)]
			if base == filenameWithoutExt {
				foundFile = filepath.Join(dir, name)
				break
			}
		}
	}

	if foundFile == "" {
		fmt.Println("Файл не найден")
		return
	}

	fmt.Println("Нашли файл:", foundFile)

	file, err := os.Open(foundFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	newWidth := 80
	ratio := float64(img.Bounds().Dy()) / float64(img.Bounds().Dx())
	newHeight := uint(float64(newWidth) * ratio * 0.55)

	resized := resize.Resize(uint(newWidth), newHeight, img, resize.Lanczos3)

	generateASCII(resized)
}

func generateASCII(img image.Image) {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			brightness := c.Y
			index := int(brightness) * (len(asciiChars) - 1) / 255
			fmt.Printf("%c", asciiChars[index])
		}
		fmt.Println()
	}
}

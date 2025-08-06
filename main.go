package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"

	"github.com/nfnt/resize"
)

var asciiChars = []byte(" .:-=+*#%@")

func main() {
	// Открываем изображение
	file, err := os.Open("input.jpg")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {
		panic(err)
	}

	// Меняем размер
	newWidth := 80
	ratio := float64(img.Bounds().Dy()) / float64(img.Bounds().Dx())
	newHeight := uint(float64(newWidth) * ratio * 0.55) // коррекция пропорций
	resized := resize.Resize(uint(newWidth), newHeight, img, resize.Lanczos3)

	// Генерируем
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

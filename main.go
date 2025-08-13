package main

import (
	"flag"
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
	filenamePtr := flag.String("file", "input", "Имя файла без расширения")
	pathPtr := flag.String("path", ".", "Путь к директории с изображением")
	widthPtr := flag.Int("width", 80, "Ширина ASCII-арт")
	outputPtr := flag.String("output", "", "Путь к файлу для сохранения результата (если пусто — вывод в консоль)")

	flag.Parse()

	filenameWithoutExt := *filenamePtr
	dir := *pathPtr
	newWidth := *widthPtr
	outputFile := *outputPtr

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

	ratio := float64(img.Bounds().Dy()) / float64(img.Bounds().Dx())
	newHeight := uint(float64(newWidth) * ratio * 0.55)

	resized := resize.Resize(uint(newWidth), newHeight, img, resize.Lanczos3)

	if outputFile != "" {
		outF, err := os.Create(outputFile)
		if err != nil {
			log.Fatal(err)
		}
		defer outF.Close()

		generateASCII(resized, outF)
		fmt.Println("Результат сохранен в файл:", outputFile)
	} else {
		generateASCII(resized, os.Stdout)
	}
}

func generateASCII(img image.Image, out *os.File) {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			brightness := c.Y
			index := int(brightness) * (len(asciiChars) - 1) / 255
			fmt.Fprintf(out, "%c", asciiChars[index])
		}
		fmt.Fprintln(out)
	}
}

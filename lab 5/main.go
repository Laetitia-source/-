package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"sync"
	"time"
)

// Функция для обработки чисел через канал (Задача 1)
func count(ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for num := range ch {
		fmt.Printf("Число %d в квадрате: %d\n", num, num*num)
	}
}

// Функция для последовательной обработки изображения (Задача 2)
func filterSequential(img draw.Image) {
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			originalColor := img.At(x, y)
			r, g, b, _ := originalColor.RGBA()

			// Перевод в оттенки серого
			gray := uint16((r + g + b) / 3)
			grayColor := color.RGBA64{R: gray, G: gray, B: gray, A: 65535}
			img.Set(x, y, grayColor)
		}
	}
}

// Функция для обработки одной строки изображения (Задача 3)
func filterParallel(img draw.Image, y int, wg *sync.WaitGroup) {
	defer wg.Done()
	bounds := img.Bounds()
	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		originalColor := img.At(x, y)
		r, g, b, _ := originalColor.RGBA()

		// Перевод в оттенки серого
		gray := uint16((r + g + b) / 3)
		grayColor := color.RGBA64{R: gray, G: gray, B: gray, A: 65535}
		img.Set(x, y, grayColor)
	}
}

// Функция для применения матричного фильтра (Задача 4)
func filterWithKernel(img image.Image, output draw.Image, kernel [][]float64, wg *sync.WaitGroup, yStart, yEnd int) {
	defer wg.Done()

	bounds := img.Bounds()
	kernelSize := len(kernel)
	offset := kernelSize / 2
	adjustmentFactor := 1.1 // Уменьшаем эффект размытия ещё больше

	for y := yStart; y < yEnd; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var rSum, gSum, bSum float64

			for ky := 0; ky < kernelSize; ky++ {
				for kx := 0; kx < kernelSize; kx++ {
					ix := x + kx - offset
					iy := y + ky - offset

					// Проверка границ изображения
					if ix >= bounds.Min.X && ix < bounds.Max.X && iy >= bounds.Min.Y && iy < bounds.Max.Y {
						r, g, b, _ := img.At(ix, iy).RGBA()
						weight := kernel[ky][kx]

						rSum += float64(r) * weight
						gSum += float64(g) * weight
						bSum += float64(b) * weight
					}
				}
			}

			output.Set(x, y, color.RGBA64{
				R: uint16(clamp(rSum*adjustmentFactor/16, 0, 65535)),
				G: uint16(clamp(gSum*adjustmentFactor/16, 0, 65535)),
				B: uint16(clamp(bSum*adjustmentFactor/16, 0, 65535)),
				A: 65535,
			})
		}
	}
}

// Функция для ограничения значений
func clamp(value float64, min float64, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}



func main() {

	// --- Задача 1: Горутины и каналы ---
	ch := make(chan int)
	var wg sync.WaitGroup

	wg.Add(1)
	go count(ch, &wg)

	for i := 1; i <= 5; i++ {
		ch <- i
	}
	close(ch)
	wg.Wait()
	fmt.Println("Задача 1 завершена.")


	// --- Задачи 2 и 3: Обработка изображений ---
	inputFile, err := os.Open("image.png")
	if err != nil {
		fmt.Println("Ошибка при открытии изображения:", err)
		return
	}
	defer inputFile.Close()

	srcImage, _, err := image.Decode(inputFile)
	if err != nil {
		fmt.Println("Ошибка при декодировании изображения:", err)
		return
	}

	bounds := srcImage.Bounds()
	dstImage := image.NewRGBA(bounds)
	draw.Draw(dstImage, bounds, srcImage, bounds.Min, draw.Src)

	// --- Задача 2: Последовательная обработка ---
	start := time.Now()
	filterSequential(dstImage)
	duration := time.Since(start)
	fmt.Println("Время последовательной обработки:", duration)

	outputFileSeq, err := os.Create("output_sequential.png")
	if err != nil {
		fmt.Println("Ошибка при создании выходного файла:", err)
		return
	}
	defer outputFileSeq.Close()

	err = png.Encode(outputFileSeq, dstImage)
	if err != nil {
		fmt.Println("Ошибка при сохранении изображения:", err)
		return
	}
	fmt.Println("Задача 2 завершена. Последовательное изображение сохранено.")

	// --- Задача 3: Параллельная обработка ---
	dstImage = image.NewRGBA(bounds)
	draw.Draw(dstImage, bounds, srcImage, bounds.Min, draw.Src)

	start = time.Now()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		wg.Add(1)
		go filterParallel(dstImage, y, &wg)
	}
	wg.Wait()
	duration = time.Since(start)
	fmt.Println("Время параллельной обработки:", duration)

	outputFilePar, err := os.Create("output_parallel.png")
	if err != nil {
		fmt.Println("Ошибка при создании выходного файла:", err)
		return
	}
	defer outputFilePar.Close()

	err = png.Encode(outputFilePar, dstImage)
	if err != nil {
		fmt.Println("Ошибка при сохранении изображения:", err)
		return
	}
	fmt.Println("Задача 3 завершена. Параллельное изображение сохранено.")

	
	// --- Задача 4: Матричный фильтр ---
	kernel := [][]float64{
		{1, 2, 1},
		{2, 4, 2},
		{1, 2, 1},
	}
	for i := range kernel {
		for j := range kernel[i] {
			kernel[i][j] /= 16.0
		}
	}
	dstImage = image.NewRGBA(bounds)
	draw.Draw(dstImage, bounds, srcImage, bounds.Min, draw.Src)

	start = time.Now()
	rows := 4 // Number of rows to split the image into
	rowHeight := bounds.Dy() / rows

	for i := 0; i < rows; i++ {
		yStart := bounds.Min.Y + i*rowHeight
		yEnd := yStart + rowHeight
		if i == rows-1 {
			yEnd = bounds.Max.Y
		}
		wg.Add(1)
		go filterWithKernel(srcImage, dstImage, kernel, &wg, yStart, yEnd)
	}
	wg.Wait()
	duration = time.Since(start)
	fmt.Println("Время обработки с матричным фильтром:", duration)

	outputFileKernel, err := os.Create("output_kernel.png")
	if err != nil {
		fmt.Println("Ошибка при создании выходного файла:", err)
		return
	}
	defer outputFileKernel.Close()

	err = png.Encode(outputFileKernel, dstImage)
	if err != nil {
		fmt.Println("Ошибка при сохранении изображения:", err)
		return
	}
	fmt.Println("Задача 4 завершена. Изображение с матричным фильтром сохранено.")
}            
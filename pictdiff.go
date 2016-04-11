// run as standalone binary for best performance (20x difference)

package main

import "image"
import "image/draw"
import _ "image/jpeg"
import "image/png"
import "os"
import "log"
import "fmt"
import "runtime"

func AsList(img *image.RGBA, x int, y int) ([]int) {
	off := img.PixOffset(x, y)
	r := int(img.Pix[off])
	g := int(img.Pix[off + 1])
	b := int(img.Pix[off + 2])
	a := int(img.Pix[off + 3])
	return []int{r, g, b, a}
}

func Abs(x int) (r int) {
	r = x
	if r < 0 {
		r = -r
	}
	return 
}

func Max(x int, y int) (r int) {
	r = x
	if y > x {
		r = y
	}
	return
}

type calcrowret struct {
	Y int
	Diff int
	Pixels *[]uint8
}

func calcrow(c *chan calcrowret, img1 *image.RGBA, img2 *image.RGBA, y int, width int) {
	totaldiff := 0
	pixel_list := make([]uint8, width * 4, width * 4)

	for x := 0; x < width; x += 1 {
		p1 := AsList(img1, x, y)
		p2 := AsList(img2, x, y)

		totplus := 0
		absdiff := Abs(p2[3] - p1[3])
		diffpixel := []int{255, 255, 255}

		for i := 0; i < 3; i += 1 {
			diff := p2[i] - p1[i]
			absdiff += Abs(diff)
			totplus += Max(0, diff)
			diffpixel[i] += diff
		}
		totaldiff += absdiff

		for i := 0; i < 3; i += 1 {
			diffpixel[i] -= totplus
			if absdiff > 0 && absdiff < 5 {
				diffpixel[i] -= 5
			}
			diffpixel[i] = Max(0, diffpixel[i])
		}

		pixel_list[x * 4 + 0] = uint8(diffpixel[0])
		pixel_list[x * 4 + 1] = uint8(diffpixel[1])
		pixel_list[x * 4 + 2] = uint8(diffpixel[2])
		pixel_list[x * 4 + 3] = 255
	}

	*c <- calcrowret{Y: y, Diff: totaldiff, Pixels: &pixel_list}
}

func Load(c *chan *image.RGBA, name string) {
	f, err := os.Open(name)
	if err != nil {
		log.Fatal("Image could not be opened")
	}
	rimg, _, err := image.Decode(f)
	if err != nil {
		log.Fatal("Image could not be decoded")
	}
	width := rimg.Bounds().Dx()
	height := rimg.Bounds().Dy()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), rimg, rimg.Bounds().Min, draw.Src)
	*c <- img
}

func main() {
	cimg1 := make(chan *image.RGBA)
	cimg2 := make(chan *image.RGBA)

	go Load(&cimg1, os.Args[1])
	go Load(&cimg2, os.Args[2])

	img1 := <- cimg1
	img2 := <- cimg2

	if img1.Bounds() != img2.Bounds() {
		log.Fatal("Images don't have the same size")
	}

	width := img1.Bounds().Dx()
	height := img1.Bounds().Dy()

	totaldiff := 0
	mapimg := image.NewNRGBA(image.Rect(0, 0, width, height))

	runtime.GOMAXPROCS(runtime.NumCPU())
	diffmeasurements := make(chan calcrowret, height)
	for y := 0; y < height; y += 1 {
		go calcrow(&diffmeasurements, img1, img2, y, width)
	}
	for y := 0; y < height; y += 1 {
		result := <- diffmeasurements
		totaldiff += result.Diff
		for x := 0; x < width; x += 1 {
			off := mapimg.PixOffset(x, result.Y)
			off2 := x * 4
			mapimg.Pix[off + 0] = (*result.Pixels)[off2 + 0]
			mapimg.Pix[off + 1] = (*result.Pixels)[off2 + 1]
			mapimg.Pix[off + 2] = (*result.Pixels)[off2 + 2]
			mapimg.Pix[off + 3] = (*result.Pixels)[off2 + 3]
		}
	}

	mapfile, _ := os.Create(os.Args[3])
	defer mapfile.Close()

    	png.Encode(mapfile, mapimg)
	fmt.Printf("Difference: %v\n", totaldiff)
}

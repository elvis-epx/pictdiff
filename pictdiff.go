// run as standalone binary for best performance (20x difference)

package main

import "image"
import "image/draw"
import _ "image/jpeg"
import "image/png"
import "os"
import "log"
import "fmt"

func Abs(x int) (r int) {
    r = x
    if r < 0 {
        r = -r
    }
    return
}

type calcrowret struct {
	Y int
	Diff int
	Pixels *[]uint8
}

func calcrow(c chan calcrowret, img1 *image.RGBA, img2 *image.RGBA, y int, width int) {
	totaldiff := 0
	pixel_list := make([]uint8, width * 4, width * 4)
    p1 := make([]uint8, 4, 4)
    p2 := make([]uint8, 4, 4)

	for x := 0; x < width; x += 1 {
        off1 := img1.PixOffset(x, y)
        off2 := img2.PixOffset(x, y)
        copy(p1, img1.Pix[off1:])
        copy(p2, img2.Pix[off2:])

		totplus := 0
		absdiff := Abs(int(p2[3]) - int(p1[3]))
		diffpixel := [3]int{255, 255, 255}

		for i := 0; i < 3; i += 1 {
			diff := int(p2[i]) - int(p1[i])
			absdiff += Abs(diff)
			totplus += max(0, diff)
			diffpixel[i] += diff
		}
		totaldiff += absdiff

		for i := 0; i < 3; i += 1 {
			diffpixel[i] -= totplus
			if absdiff > 0 && absdiff < 5 {
				diffpixel[i] -= 2
			}
			diffpixel[i] = max(0, diffpixel[i])
		}

		pixel_list[x * 4 + 0] = uint8(diffpixel[0])
		pixel_list[x * 4 + 1] = uint8(diffpixel[1])
		pixel_list[x * 4 + 2] = uint8(diffpixel[2])
		pixel_list[x * 4 + 3] = 255
	}

	c <- calcrowret{Y: y, Diff: totaldiff, Pixels: &pixel_list}
}

func Load(c chan *image.RGBA, name string) {
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
	c <- img
}

func main() {
	cimg1 := make(chan *image.RGBA)
	cimg2 := make(chan *image.RGBA)

	if len(os.Args) < 4 {
		os.Stderr.WriteString("\n")
		os.Stderr.WriteString("Usage: pictdiff <picture A> <picture B> <diff map>\n")
		os.Stderr.WriteString("\n")
		os.Stderr.WriteString("Example: pictdiff a.png b.png diff.png\n")
		os.Stderr.WriteString("\n")
		os.Exit(2)
	}

	go Load(cimg1, os.Args[1])
	go Load(cimg2, os.Args[2])

	img1 := <-cimg1
	img2 := <-cimg2

	if img1.Bounds() != img2.Bounds() {
		log.Fatal("Images don't have the same size")
	}

	width := img1.Bounds().Dx()
	height := img1.Bounds().Dy()

	totaldiff := 0
	mapimg := image.NewNRGBA(image.Rect(0, 0, width, height))

	diffmeasurements := make(chan calcrowret, height)
	for y := 0; y < height; y += 1 {
		go calcrow(diffmeasurements, img1, img2, y, width)
	}
	for y := 0; y < height; y += 1 {
		result := <-diffmeasurements
		totaldiff += result.Diff
        copy(mapimg.Pix[mapimg.PixOffset(0, result.Y):], *result.Pixels)
	}

	mapfile, err := os.Create(os.Args[3])
	if err != nil {
		os.Stderr.WriteString("Cannot open diff map file for writing\n")
	} else {
		defer mapfile.Close()
		err := (&png.Encoder{CompressionLevel: png.BestSpeed}).Encode(mapfile, mapimg)
		if err != nil {
			os.Stderr.WriteString("Cannot write to diff map file\n")
		}
	}
	fmt.Printf("%v\n", totaldiff)
}

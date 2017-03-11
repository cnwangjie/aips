package aips

import (
	"image"
	"image/color"
	"math"
)

func ColorDiff(a, b color.Color) float64 {
	ra, ga, ba, _ := a.RGBA()
	rb, gb, bb, _ := b.RGBA()
	return math.Sqrt((30*float64(ra-rb)*float64(ra-rb) + 59*float64(ga-gb)*float64(ga-gb) + 11*float64(ba-bb)*float64(ba-bb)) / 100)
}

func Gaussian(x, y int, sigma float64) float64 {
	return (1 / (2 * math.Pi * sigma * sigma)) * math.Pow(math.E, -(float64(x)*float64(x)+float64(y)*float64(y))/(4*sigma*sigma))
}

func Gray(c color.Color) uint16 {
	r, g, b, _ := c.RGBA()
	return uint16((r*30 + g*59 + b*11) / 100)
}

func Binarization(src image.Image) image.Image {
	rect := src.Bounds()
	width := rect.Dx()
	height := rect.Dy()
	dst := image.NewRGBA(rect)
	for x := 0; x < width; x += 1 {
		for y := 0; y < height; y += 1 {
			c := color.Black
			if Gray(src.At(x, y)) > 32768 {
				c = color.White
			}
			dst.Set(x, y, c)
		}
	}
	return dst
}

func GrayImg(src image.Image) image.Image {
	rect := src.Bounds()
	width := rect.Dx()
	height := rect.Dy()
	dst := image.NewRGBA(rect)
	for x := 0; x < width; x += 1 {
		for y := 0; y < height; y += 1 {
			gray := Gray(src.At(x, y))
			dst.Set(x, y, color.Gray16{gray})
		}
	}
	return dst
}

func Blur(src image.Image, ra int) image.Image {
	srcRect := src.Bounds()
	width := srcRect.Dx()
	height := srcRect.Dy()
	dst := image.NewRGBA(srcRect)

	weight := make([][]float64, ra*2+1)
	var sum float64
	sigma := (float64(ra)*2 + 1) / 2
	for x := 0; x < ra*2+1; x += 1 {
		col := make([]float64, ra*2+1)
		for y := 0; y < ra*2+1; y += 1 {
			col[y] = Gaussian(x-ra, y-ra, sigma)
			sum += col[y]
		}
		weight[x] = col
	}
	for x := 0; x < ra*2+1; x += 1 {
		for y := 0; y < ra*2+1; y += 1 {
			weight[x][y] = weight[x][y] / sum
		}
	}

	for x := 0; x < width; x += 1 {
		for y := 0; y < height; y += 1 {
			var r, g, b, a float64 = 0, 0, 0, 0
			for dx := 0; dx < ra*2+1; dx += 1 {
				for dy := 0; dy < ra*2+1; dy += 1 {
					tx, ty := dx, dy
					if dx < 0 || dx > width {
						tx = x
					}
					if dy < 0 || dy > height {
						ty = y
					}
					dr, dg, db, da := src.At(x-ra+tx, y-ra+ty).RGBA()
					r += float64(dr) * weight[dx][dy]
					g += float64(dg) * weight[dx][dy]
					b += float64(db) * weight[dx][dy]
					a += float64(da) * weight[dx][dy]
				}

			}
			a = -1
			dst.Set(x, y, color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)})
		}
	}
	return dst
}

package aips

import (
	"fmt"
	"image"
	"image/color"
)

var (
	NCTemp = []float64{
		100, 0, 0, 0, 100,
		0, 50, 0, 50, 0,
		0, 0, 0, 0, 0,
		0, 50, 0, 50, 0,
		100, 0, 0, 0, 100}
	SharpTemp = []float64{
		-1, -4, -7, -4, -1,
		-4, -16, -26, -16, -4,
		-7, -26, 505, -26, -7,
		-4, -16, -26, -16, -4,
		-1, -4, -7, -4, -1}
	GaussSmoothTemp = []float64{
		1, 4, 7, 4, 1,
		4, 16, 26, 16, 4,
		7, 26, 41, 26, 7,
		4, 16, 26, 16, 4,
		1, 4, 7, 4, 1}
)

// 八通道过滤
func AreaFilter(src image.Image, diff float64, sum int) image.Image {
	srcRect := src.Bounds()
	width := srcRect.Dx()
	height := srcRect.Dy()
	dst := image.NewRGBA(srcRect)
	for x := 0; x < width; x += 1 {
		for y := 0; y < height; y += 1 {
			diffSum := 0          // 周围色差大于标准的个数
			diffN := [8]float64{} // 色差从大到小的排序
			c := [8]color.Color{}
			for dx, dy, n := -1, -1, 0; n < 8; dx += 1 {
				if dx == 2 {
					dx = -1
					dy += 1
				}
				if dx == 0 && dy == 0 {
					continue
				}
				// 当前扫描的点与中间点的色差
				diffC := ColorDiff(src.At(x, y), src.At(x+dx, y+dy))
				i := n
				// 将当前颜色插入
				for i > 0 && diffN[i-1] < diffC {
					diffN[i] = diffN[i-1]
					c[i] = c[i-1]
					i -= 1
				}
				diffN[i] = diffC
				c[i] = src.At(x+dx, y+dy)
				if diffC >= diff {
					diffSum += 1
				}

				n += 1
			}

			fmt.Println(diffN, c[0])
			if diffSum >= sum {
				dst.Set(x, y, c[0])
			}
		}
	}
	return dst
}

func MidianFilter(src image.Image) image.Image {
	srcRect := src.Bounds()
	width := srcRect.Dx()
	height := srcRect.Dy()
	dst := image.NewRGBA(srcRect)
	var ra, ga, ba, aa, sa [9]uint32
	for x := 0; x < width; x += 1 {
		for y := 0; y < height; y += 1 {
			n := 9
			if x == 0 || x == width-1 {
				n = 6
			}
			if y == 0 || y == height-1 {
				if n == 6 {
					n = 4
				} else {
					n = 6
				}
			}
			ra[0], ga[0], ba[0], aa[0] = src.At(x-1, y-1).RGBA()
			ra[1], ga[1], ba[1], aa[1] = src.At(x, y-1).RGBA()
			ra[2], ga[2], ba[2], aa[2] = src.At(x+1, y-1).RGBA()
			ra[3], ga[3], ba[3], aa[3] = src.At(x+1, y).RGBA()
			ra[4], ga[4], ba[4], aa[4] = src.At(x+1, y+1).RGBA()
			ra[5], ga[5], ba[5], aa[5] = src.At(x, y+1).RGBA()
			ra[6], ga[6], ba[6], aa[6] = src.At(x-1, y+1).RGBA()
			ra[7], ga[7], ba[7], aa[7] = src.At(x-1, y).RGBA()
			ra[8], ga[8], ba[8], aa[8] = src.At(x, y).RGBA()
			for i, _ := range sa {
				sa[i] = ra[i]
			}
			for i := 0; i < 5; i += 1 {
				max := i
				for j := i + 1; j < 9; j += 1 {
					if sa[j] > sa[max] {
						max = j
					}
				}
				sa[i] = sa[i] + sa[max]
				sa[max] = sa[i] - sa[max]
				sa[i] = sa[i] - sa[max]
				ra[i] = ra[i] + ra[max]
				ra[max] = ra[i] - ra[max]
				ra[i] = ra[i] - ra[max]
				ga[i] = ga[i] + ga[max]
				ga[max] = ga[i] - ga[max]
				ga[i] = ga[i] - ga[max]
				ba[i] = ba[i] + ba[max]
				ba[max] = ba[i] - ba[max]
				ba[i] = ba[i] - ba[max]
				aa[i] = aa[i] + aa[max]
				aa[max] = aa[i] - aa[max]
				aa[i] = aa[i] - aa[max]
			}
			var r, g, b, a uint32
			if true {
				r, g, b, a = ra[4], ga[4], ba[4], aa[4]
			} else if n == 6 {
				r = (ra[3] + ra[2]) / 2
				g = (ga[3] + ga[2]) / 2
				b = (ba[3] + ba[2]) / 2
				a = (aa[3] + aa[2]) / 2
			} else {
				r = (ra[1] + ra[2]) / 2
				g = (ga[1] + ga[2]) / 2
				b = (ba[1] + ba[2]) / 2
				a = (aa[1] + aa[2]) / 2
			}
			dst.Set(x, y, color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)})

		}
	}
	return dst
}

func TemplateFilter(src image.Image, temp []float64) image.Image {
	srcRect := src.Bounds()
	width := srcRect.Dx()
	height := srcRect.Dy()
	size := len(temp)
	min := 0
	if width < height {
		min = width
	} else {
		min = height
	}
	filterWidth := 1
	// 计算模板的宽度 模板宽度不能比图片的窄边还宽
	for i := 3; i < min; i += 2 {
		if i*i >= size {
			filterWidth = i
			break
		}
	}
	if filterWidth == 1 {
		return src
	}
	dst := image.NewRGBA(srcRect)
	var sum float64
	for i := 0; i < size; i += 1 {
		sum += temp[i]
	}
	for i := 0; i < size; i += 1 {
		temp[i] = temp[i] / sum
	}
	for x := 0; x < width; x += 1 {
		for y := 0; y < height; y += 1 {
			var r, g, b, a float64
			for dx, dy, n := 0, 0, 0; n < size; dx += 1 {
				if dx == filterWidth {
					dx = 0
					dy += 1
				}
				dr, dg, db, da := src.At(x+dx-filterWidth/2, y+dy-filterWidth/2).RGBA()
				r += float64(dr) * temp[n]
				g += float64(dg) * temp[n]
				b += float64(db) * temp[n]
				a += float64(da) * temp[n]
				n += 1
			}
			dst.Set(x, y, color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)})
		}
	}
	return dst
}

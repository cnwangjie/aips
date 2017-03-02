package aips

import (
	"image"
	"image/color"
	"math"
)

func Rotate(src image.Image, ang float64) image.Image {
	// 计算旋转后保留完整图像的宽和高
	cosA, sinA := math.Cos(-ang), math.Sin(-ang)
	srcRect := src.Bounds()
	srcWidth := srcRect.Dx()
	srcHeight := srcRect.Dy()
	dstWidth := int(math.Abs(float64(srcWidth)*cosA - float64(srcHeight)*sinA + 0.5))
	dstHeight := int(math.Abs(float64(-srcWidth)*cosA + float64(srcHeight)*sinA + 0.5))
	dstRect := image.Rect(0, 0, dstWidth, dstHeight)
	dst := image.NewRGBA(dstRect)

	// 计算新图片
	rx, ry := float64(srcWidth)/2, float64(srcHeight)/2
	dx, dy := float64(dstWidth)/2, float64(dstHeight)/2

	for x := 0; x < dstWidth; x += 1 {
		for y := 0; y < dstHeight; y += 1 {
			// 反向映射
			sx := (float64(x)-dx)*cosA - (float64(y)-dy)*(-sinA) + rx
			sy := (float64(x)-dx)*(-sinA) + (float64(y)-dy)*cosA + ry
			Lf, Rf, Tf, Bf := math.Floor(sx), math.Ceil(sx), math.Floor(sy), math.Floor(sy)
			Li, Ri, Ti, Bi := int(Lf), int(Rf), int(Tf), int(Bf)
			r1, g1, b1, a1 := src.At(Li, Ti).RGBA()
			r2, g2, b2, a2 := src.At(Ri, Ti).RGBA()
			r3, g3, b3, a3 := src.At(Li, Bi).RGBA()
			r4, g4, b4, a4 := src.At(Ri, Bi).RGBA()
			r1f, g1f, b1f, a1f := float64(r1), float64(g1), float64(b1), float64(a1)
			r2f, g2f, b2f, a2f := float64(r2), float64(g2), float64(b2), float64(a2)
			r3f, g3f, b3f, a3f := float64(r3), float64(g3), float64(b3), float64(a3)
			r4f, g4f, b4f, a4f := float64(r4), float64(g4), float64(b4), float64(a4)

			p1 := (1 - sx + Lf) * (1 - sy + Tf)
			p2 := (sx - Lf) * (1 - sy + Tf)
			p3 := (1 - sx + Lf) * (sy - Tf)
			p4 := (sx - Lf) * (sy - Tf)
			r := p1*r1f + p2*r2f + p3*r3f + p4*r4f
			g := p1*g1f + p2*g2f + p3*g3f + p4*g4f
			b := p1*b1f + p2*b2f + p3*b3f + p4*b4f
			a := p1*a1f + p2*a2f + p3*a3f + p4*a4f

			dst.Set(x, y, color.RGBA64{uint16(r + 0.5), uint16(g + 0.5), uint16(b + 0.5), uint16(a + 0.5)})
		}
	}

	return dst
}

func Scale(src image.Image, zoomx float64, zoomy float64) image.Image {
	srcRect := src.Bounds()
	srcWidth := srcRect.Dx()
	srcHeight := srcRect.Dy()
	dstWidth := int(float64(srcWidth) * zoomx)
	dstHeight := int(float64(srcHeight) * zoomy)
	dstRect := image.Rect(0, 0, dstWidth, dstHeight)
	dst := image.NewRGBA(dstRect)

	for x := 0; x < dstWidth; x += 1 {
		for y := 0; y < dstHeight; y += 1 {
			// 反向映射
			sx := float64(x) / zoomx
			sy := float64(y) / zoomy
			Lf, Rf, Tf, Bf := math.Floor(sx), math.Ceil(sx), math.Floor(sy), math.Floor(sy)
			Li, Ri, Ti, Bi := int(Lf), int(Rf), int(Tf), int(Bf)
			r1, g1, b1, a1 := src.At(Li, Ti).RGBA()
			r2, g2, b2, a2 := src.At(Ri, Ti).RGBA()
			r3, g3, b3, a3 := src.At(Li, Bi).RGBA()
			r4, g4, b4, a4 := src.At(Ri, Bi).RGBA()
			r1f, g1f, b1f, a1f := float64(r1), float64(g1), float64(b1), float64(a1)
			r2f, g2f, b2f, a2f := float64(r2), float64(g2), float64(b2), float64(a2)
			r3f, g3f, b3f, a3f := float64(r3), float64(g3), float64(b3), float64(a3)
			r4f, g4f, b4f, a4f := float64(r4), float64(g4), float64(b4), float64(a4)

			p1 := (1 - sx + Lf) * (1 - sy + Tf)
			p2 := (sx - Lf) * (1 - sy + Tf)
			p3 := (1 - sx + Lf) * (sy - Tf)
			p4 := (sx - Lf) * (sy - Tf)
			r := p1*r1f + p2*r2f + p3*r3f + p4*r4f
			g := p1*g1f + p2*g2f + p3*g3f + p4*g4f
			b := p1*b1f + p2*b2f + p3*b3f + p4*b4f
			a := p1*a1f + p2*a2f + p3*a3f + p4*a4f

			dst.Set(x, y, color.RGBA64{uint16(r + 0.5), uint16(g + 0.5), uint16(b + 0.5), uint16(a + 0.5)})
		}
	}

	return dst
}

func Resize(src image.Image, dstWidth int, dstHeight int) image.Image {
	srcRect := src.Bounds()
	srcWidth := srcRect.Dx()
	srcHeight := srcRect.Dy()
	dstRect := image.Rect(0, 0, dstWidth, dstHeight)
	dst := image.NewRGBA(dstRect)
	zoomx := float64(dstWidth) / float64(srcWidth)
	zoomy := float64(dstHeight) / float64(srcHeight)
	for x := 0; x < dstWidth; x += 1 {
		for y := 0; y < dstHeight; y += 1 {
			// 反向映射
			sx := float64(x) / zoomx
			sy := float64(y) / zoomy
			Lf, Rf, Tf, Bf := math.Floor(sx), math.Ceil(sx), math.Floor(sy), math.Floor(sy)
			Li, Ri, Ti, Bi := int(Lf), int(Rf), int(Tf), int(Bf)
			r1, g1, b1, a1 := src.At(Li, Ti).RGBA()
			r2, g2, b2, a2 := src.At(Ri, Ti).RGBA()
			r3, g3, b3, a3 := src.At(Li, Bi).RGBA()
			r4, g4, b4, a4 := src.At(Ri, Bi).RGBA()
			r1f, g1f, b1f, a1f := float64(r1), float64(g1), float64(b1), float64(a1)
			r2f, g2f, b2f, a2f := float64(r2), float64(g2), float64(b2), float64(a2)
			r3f, g3f, b3f, a3f := float64(r3), float64(g3), float64(b3), float64(a3)
			r4f, g4f, b4f, a4f := float64(r4), float64(g4), float64(b4), float64(a4)

			p1 := (1 - sx + Lf) * (1 - sy + Tf)
			p2 := (sx - Lf) * (1 - sy + Tf)
			p3 := (1 - sx + Lf) * (sy - Tf)
			p4 := (sx - Lf) * (sy - Tf)
			r := p1*r1f + p2*r2f + p3*r3f + p4*r4f
			g := p1*g1f + p2*g2f + p3*g3f + p4*g4f
			b := p1*b1f + p2*b2f + p3*b3f + p4*b4f
			a := p1*a1f + p2*a2f + p3*a3f + p4*a4f

			dst.Set(x, y, color.RGBA64{uint16(r + 0.5), uint16(g + 0.5), uint16(b + 0.5), uint16(a + 0.5)})
		}
	}

	return dst
}

func Cut(src image.Image, rect image.Rectangle) image.Image {
	dstWidth := rect.Dx()
	dstHeight := rect.Dy()
	dstRect := image.Rect(0, 0, dstWidth, dstHeight)
	dst := image.NewRGBA(dstRect)
	for x := 0; x < dstWidth; x += 1 {
		for y := 0; y < dstHeight; y += 1 {
			dst.Set(x, y, src.At(x+rect.Min.X, y+rect.Min.Y))
		}
	}
	return dst
}

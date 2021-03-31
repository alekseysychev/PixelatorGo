package pixelator

import (
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"log"

	_ "image/jpeg"
	_ "image/png"
)

type clusterRGBA struct {
	R     uint32
	G     uint32
	B     uint32
	A     uint32
	Count uint32
}

func (c *clusterRGBA) Append(r, g, b, a uint32) {
	c.R += r >> 8
	c.G += g >> 8
	c.B += b >> 8
	c.A += a >> 8
	c.Count++
}

func (c *clusterRGBA) Bit8() color.RGBA {
	var bit uint32 = 8
	return color.RGBA{
		R: uint8(uint32(c.R*bit/c.Count/255) * 255 / bit),
		G: uint8(uint32(c.G*bit/c.Count/255) * 255 / bit),
		B: uint8(uint32(c.B*bit/c.Count/255) * 255 / bit),
		A: uint8(uint32(c.A*bit/c.Count/255) * 255 / bit),
	}
}

func (c *clusterRGBA) Avg() color.RGBA {
	log.Println(c.R/c.Count, int(c.R*8/c.Count/255)*255/8)
	return color.RGBA{
		R: uint8(c.R / c.Count),
		G: uint8(c.G / c.Count),
		B: uint8(c.B / c.Count),
		A: uint8(c.A / c.Count),
	}
}

func Compile(r io.Reader, w io.Writer, s int, q int) error {
	if s <= 0 {
		return errors.New("cluster size must be more than 0")
	}
	if q <= 0 {
		return errors.New("quality must be more than 0")
	}

	imageData, _, err := image.Decode(r)
	if err != nil {
		return err
	}

	clusters := make(map[image.Point]*clusterRGBA)

	bounds := imageData.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if clusters[image.Point{X: x / s, Y: y / s}] == nil {
				clusters[image.Point{X: x / s, Y: y / s}] = &clusterRGBA{}
			}
			clusters[image.Point{X: x / s, Y: y / s}].Append(imageData.At(x, y).RGBA())
		}
	}

	newImage := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// newImage.SetRGBA(x, y, clusters[image.Point{X: x / s, Y: y / s}].Avg())
			newImage.SetRGBA(x, y, clusters[image.Point{X: x / s, Y: y / s}].Bit8())
		}
	}

	log.Println(q)
	options := jpeg.Options{
		Quality: q,
	}
	err = jpeg.Encode(w, newImage, &options)
	if err != nil {
		return err
	}

	return nil
}

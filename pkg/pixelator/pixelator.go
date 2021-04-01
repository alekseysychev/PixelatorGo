package pixelator

import (
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"io"
	"log"

	_ "image/png"
)

type clusterRGBA struct {
	R     float64
	G     float64
	B     float64
	A     float64
	Count float64
}

func (c *clusterRGBA) Append(r, g, b, a uint32) {
	c.Count++

	ca := (c.Count - 1) / c.Count
	cb := 1 / c.Count

	c.R = c.R*ca + float64(r)*cb
	c.G = c.G*ca + float64(g)*cb
	c.B = c.B*ca + float64(b)*cb
	c.A = c.A*ca + float64(a)*cb

	// log.Print(r, c.R)
}

func (c *clusterRGBA) Bit4() color.RGBA {
	a := 65536
	log.Fatalln(a, a>>14)

	return color.RGBA{

		R: uint8(int(c.R) >> 15 * 128),
		G: uint8(int(c.G) >> 14 * 64),
		B: uint8(int(c.B) >> 15 * 128),
		A: uint8(int(c.A) >> 15 * 128),
	}
}

func (c *clusterRGBA) Color() color.RGBA {
	return color.RGBA{
		R: uint8(int(c.R) >> 8),
		G: uint8(int(c.G) >> 8),
		B: uint8(int(c.B) >> 8),
		A: uint8(int(c.A) >> 8),
	}
}

type Settings struct {
	ClusterSize int
	Quality     int
	Colors      []color.Color
}

func Compile(r io.Reader, w io.Writer, s Settings) error {
	err := settingsCheck(s)
	if err != nil {
		return err
	}

	imageData, _, err := image.Decode(r)
	if err != nil {
		return err
	}

	clusters := make(map[image.Point]*clusterRGBA)

	bounds := imageData.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if clusters[image.Point{X: x / s.ClusterSize, Y: y / s.ClusterSize}] == nil {
				clusters[image.Point{X: x / s.ClusterSize, Y: y / s.ClusterSize}] = &clusterRGBA{}
			}
			clusters[image.Point{X: x / s.ClusterSize, Y: y / s.ClusterSize}].Append(imageData.At(x, y).RGBA())
		}
	}

	newImage := image.NewRGBA(bounds)
	var clr color.RGBA
	var cluster *clusterRGBA
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			cluster = clusters[image.Point{X: x / s.ClusterSize, Y: y / s.ClusterSize}]

			clr = cluster.Color()

			newImage.SetRGBA(x, y, clr)
		}
	}

	err = jpeg.Encode(w, newImage, &jpeg.Options{
		Quality: s.Quality,
	})

	if err != nil {
		return err
	}

	return nil
}

var (
	ErrWrongClusterSize error = errors.New("wrong cluster size")
	ErrWrongQuality     error = errors.New("wrong quality")
)

func settingsCheck(s Settings) error {
	if s.ClusterSize < 1 {
		return ErrWrongClusterSize
	}
	if s.Quality < 0 || s.Quality > 100 {
		return ErrWrongQuality
	}
	if len(s.Colors) == 0 {
		s.Colors = append(s.Colors, color.Black, color.White)
	}
	return nil
}

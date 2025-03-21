package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/verdantco/flatpack/geometry"
)

func main() {
	if len(os.Args) < 2 {
		panic("no .stl file provided")
	}

	file, err := os.Open("simple.stl")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	s := bufio.NewScanner(file)
	o := []*geometry.F{}

	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(line, "facet normal") {
			face := readFace(s, line)
			o = append(o, face)
		}
	}

	if err := s.Err(); err != nil {
		panic(err)
	}

	d := 200.0
	n := geometry.Normalize(o, d)
	up := &geometry.V3{X: 0.0, Y: 0.0, Z: 1.0}
	r := make([]*geometry.R, len(n))
	f := make([]*geometry.F, len(n))

	for i := range n {
		r[i] = geometry.RotationBetween(n[i].N, up)

		f[i] = n[i].Rotate(r[i])
	}

	writeMarkup(d, n, r)
}

func readFace(s *bufio.Scanner, normalLine string) *geometry.F {
	n := readVertex(strings.Fields(normalLine)[2:])

	if !s.Scan() || s.Text() != "outer loop" {
		panic("invalid facet: missing outer loop")
	}

	v0 := readVertexFromScanner(s, "v0")
	v1 := readVertexFromScanner(s, "v1")
	v2 := readVertexFromScanner(s, "v2")

	if !s.Scan() || s.Text() != "endloop" {
		panic("invalid facet: missing endloop")
	}
	if !s.Scan() || s.Text() != "endfacet" {
		panic("invalid facet: missing endfacet")
	}

	return &geometry.F{
		N: n,
		V: [3]*geometry.V3{v0, v1, v2},
	}
}

func readVertexFromScanner(s *bufio.Scanner, label string) *geometry.V3 {
	if !s.Scan() || !strings.HasPrefix(s.Text(), "vertex ") {
		panic(fmt.Sprintf("invalid facet: missing %s", label))
	}
	return readVertex(strings.Fields(s.Text())[1:])
}

func readVertex(str []string) *geometry.V3 {
	if len(str) != 3 {
		panic("invalid vertex")
	}

	x, err := strconv.ParseFloat(str[0], 64)
	if err != nil {
		panic("invalid vertex X")
	}
	y, err := strconv.ParseFloat(str[1], 64)
	if err != nil {
		panic("invalid vertex Y")
	}
	z, err := strconv.ParseFloat(str[2], 64)
	if err != nil {
		panic("invalid vertex Z")
	}

	return &geometry.V3{X: x, Y: y, Z: z}
}

func writeMarkup(s float64, o []*geometry.F, r []*geometry.R) {
	svgOpen := fmt.Sprintf("<svg viewBox=\"0 0 %d %d\" style=\"background:black\">", int(s), int(s))
	fmt.Println(svgOpen)

	format := `<polygon class="face" points="%.5f,%.5f %.5f,%.5f %.5f,%.5f" />`

	for i := range o {
		x0, x1, x2 := o[i].V[0].X, o[i].V[1].X, o[i].V[2].X
		y0, y1, y2 := s-o[i].V[0].Y, s-o[i].V[1].Y, s-o[i].V[2].Y

		triangle := fmt.Sprintf(format, x0, y0, x1, y1, x2, y2)
		fmt.Println(triangle)
	}

	svgClose := "</svg>"
	fmt.Println(svgClose)
}

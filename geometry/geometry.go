package geometry

import (
	"encoding/json"
	"math"
)

type V3 struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

func (v *V3) String() string {
	bytes, _ := json.Marshal(v)

	return string(bytes)
}

func (v *V3) Len() float64 {
	x := math.Pow(v.X, 2)
	y := math.Pow(v.Y, 2)
	z := math.Pow(v.Z, 2)

	return math.Sqrt(x + y + z)
}

func (v *V3) Scale(s float64) *V3 {
	return &V3{
		X: v.X * s,
		Y: v.Y * s,
		Z: v.Z * s,
	}
}

func (v *V3) Unit() *V3 {
	if v.Len() == 0.0 {
		return v
	}

	return v.Scale(1 / v.Len())
}

func (v *V3) Rotate(r *R) *V3 {
	cosQ := math.Cos(r.Q)
	sinQ := math.Sin(r.Q)
	dot := Dot(r.A, v)

	term1 := v.Scale(cosQ)
	term2 := Cross(r.A, v).Scale(sinQ)
	term3 := r.A.Scale(dot * (1 - cosQ))

	return Add(term1, Add(term2, term3))
}

type F struct {
	N *V3    `json:"n"`
	V [3]*V3 `json:"v"`
}

func (f *F) String() string {
	bytes, _ := json.Marshal(f)

	return string(bytes)
}

func (f *F) Rotate(r *R) *F {
	return &F{
		N: f.N.Rotate(r),
		V: [3]*V3{
			f.V[0].Rotate(r),
			f.V[1].Rotate(r),
			f.V[2].Rotate(r),
		},
	}
}

type R struct {
	Q float64 `json:"q"`
	A *V3     `json:"a"`
}

func (r *R) String() string {
	bytes, _ := json.Marshal(r)

	return string(bytes)
}

func Add(v1, v2 *V3) *V3 {
	return &V3{
		X: v1.X + v2.X,
		Y: v1.Y + v2.Y,
		Z: v1.Z + v2.Z,
	}
}

func Subtract(v1, v2 *V3) *V3 {
	return &V3{
		X: v1.X - v2.X,
		Y: v1.Y - v2.Y,
		Z: v1.Z - v2.Z,
	}
}

func Dot(v1, v2 *V3) float64 {
	return (v1.X * v2.X) + (v1.Y * v2.Y) + (v1.Z * v2.Z)
}

func Cross(v1, v2 *V3) *V3 {
	return &V3{
		X: v1.Y*v2.Z - v1.Z*v2.Y,
		Y: v1.Z*v2.X - v1.X*v2.Z,
		Z: v1.X*v2.Y - v1.Y*v2.X,
	}
}

func RotationBetween(v1, v2 *V3) *R {
	cos := Dot(v1, v2) / (v1.Len() * v2.Len())

	return &R{
		Q: math.Acos(cos),
		A: Cross(v1, v2).Unit(),
	}
}

func Normalize(o []*F, s float64) []*F {
	n := make([]*F, len(o))

	minX, minY, minZ := math.Inf(1), math.Inf(1), math.Inf(1)
	maxX, maxY, maxZ := math.Inf(-1), math.Inf(-1), math.Inf(-1)

	for i := range o {
		for j := range o[i].V {
			minX, maxX = math.Min(minX, o[i].V[j].X), math.Max(maxX, o[i].V[j].X)
			minY, maxY = math.Min(minY, o[i].V[j].Y), math.Max(maxY, o[i].V[j].Y)
			minZ, maxZ = math.Min(minZ, o[i].V[j].Z), math.Max(maxZ, o[i].V[j].Z)
		}
	}

	xRange, yRange, zRange := maxX-minX, maxY-minY, maxZ-minZ
	maxRange := math.Max(xRange, math.Max(yRange, zRange))
	minimums := &V3{X: minX, Y: minY, Z: minZ}

	for i := range o {
		n[i] = &F{
			N: &V3{X: o[i].N.X, Y: o[i].N.Y, Z: o[i].N.Z},
			V: [3]*V3{
				Subtract(o[i].V[0], minimums).Scale(s / maxRange),
				Subtract(o[i].V[1], minimums).Scale(s / maxRange),
				Subtract(o[i].V[2], minimums).Scale(s / maxRange),
			},
		}
	}

	return n
}

package geometry

import (
	"fmt"
	"math"
	"testing"
)

func Test_Len(t *testing.T) {
	v := &V3{X: 3.0, Y: 4.0, Z: 0.0}

	if v.Len() != 5.0 {
		t.Fail()
	}
}

func Test_Scale(t *testing.T) {
	v := &V3{X: 1.0, Y: 2.0, Z: 3.0}

	l := v.Len()
	v = v.Scale(2)

	if v.Len() != l*2 {
		t.Fail()
	}
}

func Test_Unit(t *testing.T) {
	v := &V3{X: 3.0, Y: 4.0, Z: 0.0}

	u := v.Unit()

	if u.Len() != 1.0 {
		t.Fail()
	}
}

func Test_Rotate_V(t *testing.T) {
	r := &R{
		Q: math.Pi,
		A: &V3{X: 0.0, Y: 0.0, Z: 1.0},
	}

	v := &V3{X: 1.0, Y: 0.0, Z: 0.0}
	u := v.Rotate(r)

	if u.X != -1 {
		t.Fail()
	}
}

func Test_Add(t *testing.T) {
	v1 := &V3{X: 1.0, Y: 1.0, Z: 1.0}
	v2 := &V3{X: 1.0, Y: 1.0, Z: 1.0}

	v := Add(v1, v2)

	if v.X != 2 || v.Y != 2 || v.Z != 2 {
		t.Fail()
	}
}

func Test_Subtract(t *testing.T) {
	v1 := &V3{X: 1.0, Y: 1.0, Z: 1.0}
	v2 := &V3{X: 1.0, Y: 1.0, Z: 1.0}

	v := Subtract(v1, v2)

	if v.Len() != 0.0 {
		t.Fail()
	}
}

func Test_Dot(t *testing.T) {
	v := &V3{X: 1.0, Y: 2.0, Z: 3.0}

	if Dot(v, v) != 14.0 {
		t.Fail()
	}
}

func Test_Cross(t *testing.T) {
	v := &V3{X: 1.0, Y: 2.0, Z: 3.0}

	if Cross(v, v).Len() != 0.0 {
		t.Fail()
	}
}

func Test_RotationBetween(t *testing.T) {
	v1 := &V3{1.0, 0.0, 0.0}
	v2 := &V3{-1.0, 0.0, 0.0}

	r := RotationBetween(v1, v2)

	fmt.Println(r.String())

	if r.Q != math.Pi || r.A.Z < 0.999 {
		t.Fail()
	}
}

func Test_Normalize(t *testing.T) {
	o := []*F{
		{
			N: &V3{X: 0, Y: 0, Z: 1},
			V: [3]*V3{
				{X: -1, Y: -1, Z: 0},
				{X: -1, Y: 1, Z: 0},
				{X: 2, Y: 2, Z: 0},
			},
		},
	}

	n := Normalize(o, 1.0)

	fmt.Print(n)

	for i := range n[0].V {
		outsideXBounds := n[0].V[i].X < 0 || n[0].V[i].Y < 0
		outsideYBounds := n[0].V[i].X > 1 || n[0].V[i].Y > 1

		if outsideXBounds || outsideYBounds {
			t.Fail()
		}
	}
}

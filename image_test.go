package main

import (
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func isNilOrObject(e any, a any) bool {
	if reflect.ValueOf(a).IsNil() && reflect.ValueOf(e).IsNil() {
		return true
	}
	if !reflect.ValueOf(a).IsNil() && !reflect.ValueOf(e).IsNil() {
		return true
	}

	return false
}

func asserPointsAreEqual(t *testing.T, e, a Point, msg any) {
	assert.InDelta(t, e.X, a.X, 0.1, msg)
	assert.InDelta(t, e.Y, a.Y, 0.1, msg)
}

func TestRead(t *testing.T) {
	testdata := []struct {
		file  string
		image *Image
	}{
		{
			"testdata/ime.hvif",
			&Image{
				styles: []Style{
					&Color{Red: 1, Green: 1, Blue: 1, Alpha: 116},
					&Color{Red: 1, Green: 1, Blue: 1, Alpha: 255},
					&Gradient{Type: GradientLinear, Offsets: []uint8{0, 255},
						Colors: []Color{
							{Red: 246, Green: 197, Blue: 79, Alpha: 255},
							{Red: 168, Green: 120, Blue: 4, Alpha: 255},
						},
						Transformable: &TransformerAffine{
							Matrix: [6]float32{0.25, -0.125, 0.25, 0.5, 16.0, 8.0},
						},
					},
					&Gradient{Type: GradientLinear, Offsets: []uint8{0, 255},
						Colors: []Color{
							{Red: 255, Green: 238, Blue: 199, Alpha: 255},
							{Red: 246, Green: 197, Blue: 79, Alpha: 255},
						},
						Transformable: &TransformerAffine{
							Matrix: [6]float32{0.25, -0.125, 0.25, 0.5, 16.0, 8.0},
						},
					},
					&Gradient{Type: GradientLinear, Offsets: []uint8{0, 254},
						Colors: []Color{
							{Red: 168, Green: 120, Blue: 4, Alpha: 255},
							{Red: 202, Green: 154, Blue: 37, Alpha: 255},
						},
						Transformable: &TransformerAffine{
							Matrix: [6]float32{0.25, -0.0625, 0.25, 0.5, 16.0, 4.0},
						},
					},
					&Gradient{Type: GradientCircular, Offsets: []uint8{0, 185, 255},
						Colors: []Color{
							{Red: 235, Green: 137, Blue: 255, Alpha: 255},
							{Red: 195, Green: 4, Blue: 233, Alpha: 255},
							{Red: 151, Green: 8, Blue: 179, Alpha: 255},
						},
						Transformable: &TransformerAffine{
							Matrix: [6]float32{0.125, 0.0, 0.0, 0.0625, 32.0, 16.0},
						},
					},
					&Gradient{Type: GradientCircular, Offsets: []uint8{0, 185, 255},
						Colors: []Color{
							{Red: 189, Green: 244, Blue: 178, Alpha: 255},
							{Red: 37, Green: 198, Blue: 5, Alpha: 255},
							{Red: 20, Green: 107, Blue: 2, Alpha: 255},
						},
						Transformable: &TransformerAffine{
							Matrix: [6]float32{0.125, 0.0, 0.0, 0.0625, 32.0, 16.0},
						},
					},
					&Gradient{Type: GradientCircular, Offsets: []uint8{0, 185, 255},
						Colors: []Color{
							{Red: 137, Green: 176, Blue: 255, Alpha: 255},
							{Red: 4, Green: 79, Blue: 233, Alpha: 255},
							{Red: 8, Green: 64, Blue: 179, Alpha: 255},
						},
						Transformable: &TransformerAffine{
							Matrix: [6]float32{0.125, 0.0, 0.0, 0.0625, 32.0, 16.0},
						},
					},
					&Gradient{Type: GradientCircular, Offsets: []uint8{0, 185, 255},
						Colors: []Color{
							{Red: 255, Green: 137, Blue: 137, Alpha: 255},
							{Red: 233, Green: 6, Blue: 6, Alpha: 255},
							{Red: 179, Green: 9, Blue: 9, Alpha: 255},
						},
						Transformable: &TransformerAffine{
							Matrix: [6]float32{0.125, 0.0, 0.0, 0.0625, 32.0, 16.0},
						},
					},
					&Gradient{Type: GradientCircular, Offsets: []uint8{0, 255},
						Colors: []Color{
							{Red: 237, Green: 237, Blue: 237, Alpha: 255},
							{Red: 53, Green: 53, Blue: 53, Alpha: 255},
						},
						Transformable: &TransformerAffine{
							Matrix: [6]float32{0.03, 0.0, 0.0, 0.25, 32.0, 16.0},
						},
					},
					&Gradient{Type: GradientLinear, Offsets: []uint8{0, 255},
						Colors: []Color{
							{Red: 255, Green: 255, Blue: 255, Alpha: 255},
							{Red: 124, Green: 147, Blue: 177, Alpha: 255},
						},
						Transformable: &TransformerAffine{
							Matrix: [6]float32{-0.03, 0.0, 0.0, 1.0, 32.0, 0.0},
						},
					},
					&Gradient{Type: GradientCircular, Offsets: []uint8{0, 255},
						Colors: []Color{
							{Red: 221, Green: 5, Blue: 5, Alpha: 255},
							{Red: 255, Green: 5, Blue: 5, Alpha: 0},
						},
						Transformable: &TransformerAffine{
							Matrix: [6]float32{0.25, 0.0, 0.0, 0.125, 32.0, 0.0},
						},
					},
					&Gradient{Type: GradientCircular, Offsets: []uint8{0, 255},
						Colors: []Color{
							{Red: 255, Green: 255, Blue: 255, Alpha: 255},
							{Red: 160, Green: 109, Blue: 30, Alpha: 255},
						},
						Transformable: &TransformerAffine{
							Matrix: [6]float32{0.06, -0.06, 0.06, 0.06, 32.0, 4.0},
						},
					},
				},
				pathes: []*Path{
					{
						isClosed: true, Elements: []PathElement{
							&Curve{PointIn: Point{18, 22}, Point: Point{18, 22}, PointOut: Point{18, 22}},
							&Curve{PointIn: Point{18, 56}, Point: Point{18, 56}, PointOut: Point{34, 56}},
							&Curve{PointIn: Point{38, 44}, Point: Point{34, 48}, PointOut: Point{38, 44}},
							&Curve{PointIn: Point{44, 46}, Point: Point{40, 46}, PointOut: Point{48, 46}},
							&Curve{PointIn: Point{55, 45}, Point: Point{54, 46}, PointOut: Point{56, 44}},
							&Curve{PointIn: Point{64, 42}, Point: Point{64, 45}, PointOut: Point{64, 40}},
							&Curve{PointIn: Point{61, 39}, Point: Point{59.9, 39.9}, PointOut: Point{61, 39}},
							&Curve{PointIn: Point{62, 34}, Point: Point{62, 37}, PointOut: Point{62, 28}},
							&Curve{PointIn: Point{50, 22}, Point: Point{58, 26}, PointOut: Point{50, 22}},
						},
					},
					{
						isClosed: true, Elements: []PathElement{
							&Curve{PointIn: Point{2, 38}, Point: Point{2, 24}, PointOut: Point{2, 48}},
							&Curve{PointIn: Point{18, 52}, Point: Point{12, 52}, PointOut: Point{30, 52}},
							&Curve{PointIn: Point{33, 41}, Point: Point{29, 45}, PointOut: Point{37, 37}},
							&Curve{PointIn: Point{44, 42}, Point: Point{39, 42}, PointOut: Point{54, 42}},
							&Curve{PointIn: Point{58, 28}, Point: Point{58, 36}, PointOut: Point{58, 20}},
							&Curve{PointIn: Point{38, 14}, Point: Point{48, 14}, PointOut: Point{20, 14}},
						},
					},
				},
			},
		},
	}

	for _, tc := range testdata {
		file, err := os.Open(tc.file)
		if err != nil {
			t.Errorf("open file: %e", err)
		}
		defer file.Close()

		img, err := ReadImage(file)
		if err != nil {
			t.Errorf("read image: %e", err)
		}

		// Styles
		aStyles := img.GetStyles()
		eStyles := tc.image.styles
		assert.Len(t, aStyles, len(eStyles))
		for i := range len(aStyles) {
			assert.Equal(t, reflect.TypeOf(aStyles[i]), reflect.TypeOf(eStyles[i]), i)

			switch actual := aStyles[i].(type) {
			case *Color:
				expected, ok := eStyles[i].(*Color)
				assert.True(t, ok, i)
				assert.Equal(t, expected, actual, i)
			case *Gradient:
				expected, ok := eStyles[i].(*Gradient)
				assert.True(t, ok)
				assert.Equal(t, expected.Type, actual.Type, i)
				assert.Equal(t, expected.Colors, actual.Colors, i)
				assert.Equal(t, expected.Offsets, actual.Offsets, i)
				assert.True(t, isNilOrObject(expected.Transformable, actual.Transformable), i)
				assert.InDeltaSlice(t, expected.Transformable.Matrix[:], actual.Transformable.Matrix[:], 0.01, i)
			default:
				assert.Failf(t, "unrecognized style", "[%d] type: %s", i, reflect.TypeOf(aStyles[i]))
			}
		}

		// Pathes
		aPathes := img.GetPathes()
		ePathes := tc.image.pathes
		// assert.Len(t, aPathes, len(ePathes))
		for i := range len(aPathes[:2]) {
			assert.Equal(t, reflect.TypeOf(aPathes[i]), reflect.TypeOf(ePathes[i]), i)

			assert.Equal(t, aPathes[i].isClosed, ePathes[i].isClosed, i)
			for pathID := range len(aPathes[i].Elements) {
				switch actual := aPathes[i].Elements[pathID].(type) {
				case *Point:
					expected, ok := ePathes[i].Elements[pathID].(*Point)
					assert.True(t, ok)
					asserPointsAreEqual(t, *expected, *actual, i)
				case *HLine:
					expected, ok := ePathes[i].Elements[pathID].(*HLine)
					assert.True(t, ok)
					assert.InDeltaSlice(t, expected.X, actual.X, 0.1, i)
				case *VLine:
					expected, ok := ePathes[i].Elements[pathID].(*VLine)
					assert.True(t, ok)
					assert.InDeltaSlice(t, expected.Y, actual.Y, 0.1, i)
				case *Curve:
					expected, ok := ePathes[i].Elements[pathID].(*Curve)
					assert.True(t, ok)
					asserPointsAreEqual(t, expected.PointIn, actual.PointIn, i)
					asserPointsAreEqual(t, expected.PointOut, actual.PointOut, i)
					asserPointsAreEqual(t, expected.Point, actual.Point, i)
				default:
					assert.Failf(t, "unrecognized path", "[%d] type: %s", i, reflect.TypeOf(aPathes[i].Elements[pathID]))
				}
			}
		}
	}

}

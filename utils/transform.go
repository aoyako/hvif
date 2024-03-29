package utils

import (
	"encoding/binary"
	"io"
	"math"
)

func ReadFloatCoord(r io.Reader) float32 {
	var val uint16
	var x uint8
	binary.Read(r, binary.LittleEndian, &x)
	val = uint16(x)
	if val&0x80 != 0 {
		var xlow uint8
		binary.Read(r, binary.LittleEndian, &xlow)
		val = val & (0x80 - 1)
		val = (val << 8) | uint16(xlow)

		return float32(val)/102.0 - 128.0
	}

	return float32(val) - 32.0
}

func ReadFloatTrans(r io.Reader) float32 {
	var b1 uint8
	var b2 uint8
	var b3 uint8
	binary.Read(r, binary.LittleEndian, &b1)
	binary.Read(r, binary.LittleEndian, &b2)
	binary.Read(r, binary.LittleEndian, &b3)

	value := uint32(uint32(b1)<<16 | uint32(b2)<<8 | uint32(b3))
	if value == 0 {
		return 0.0
	}

	sign := ((value & 0b100000000000000000000000) >> 23)
	expo := ((value & 0b011111100000000000000000) >> 17) - 32
	mant := ((value & 0b000000011111111111111111) >> 6)

	bits := (sign << 31) | ((expo + 127) << 23) | mant
	return math.Float32frombits(bits)
}

const TransformMatrixSize = 6
const PerspectiveMatrixSize = 9

type TransformerType uint8

const (
	TransformerTypeAffine TransformerType = 20 + iota
	TransformerTypeContour
	TransformerTypePerspective
	TransformerTypeStroke
)

type Transformable struct {
	Matrix [TransformMatrixSize]float32
}

type Translation struct {
	X float32
	Y float32
}

type LodScale struct {
	MinS float32
	MaxS float32
}

func ReadMatrix(r io.Reader, size int) []float32 {
	res := make([]float32, size)
	for i := 0; i < size; i++ {
		res[i] = ReadFloatTrans(r)
	}
	return res
}

func ReadTransformable(r io.Reader) Transformable {
	var t Transformable
	copy(t.Matrix[:], ReadMatrix(r, TransformMatrixSize))
	return t
}

type TransformerAffine struct {
	Matrix [TransformMatrixSize]float32
}

func ReadAffine(r io.Reader) TransformerAffine {
	var t TransformerAffine
	copy(t.Matrix[:], ReadMatrix(r, TransformMatrixSize))
	return t
}

type LineJoinOptions uint8

const (
	MiterJoin LineJoinOptions = iota
	MiterJoinRevert
	RoundJoin
	BevelJoin
	MiterJoinRound
)

type TransformerContour struct {
	Width      float64
	LineJoin   LineJoinOptions
	MiterLimit float64
}

func ReadCountour(r io.Reader) TransformerContour {
	var t TransformerContour
	var width uint8
	var lineJoin uint8
	var miterLimit uint8
	binary.Read(r, binary.LittleEndian, &width)
	binary.Read(r, binary.LittleEndian, &lineJoin)
	binary.Read(r, binary.LittleEndian, &miterLimit)

	t.Width = (float64(width) - 128.0)
	t.LineJoin = LineJoinOptions(lineJoin)
	t.MiterLimit = float64(miterLimit)
	return t
}

type TransformerPerspective struct {
	Matrix [PerspectiveMatrixSize]float32
}

func ReadTransformerPerspective(r io.Reader) TransformerPerspective {
	var t TransformerPerspective
	copy(t.Matrix[:], ReadMatrix(r, TransformMatrixSize))
	return t
}

type LineCapOptions uint8

const (
	ButtCap LineCapOptions = iota
	SquareCap
	RoundCap
)

type TransformerStroke struct {
	Width      float32
	LineJoin   LineJoinOptions
	LineCap    LineCapOptions
	MiterLimit float32
}

func ReadTransformerStroke(r io.Reader) TransformerStroke {
	var t TransformerStroke
	var width uint8
	var lineOptions uint8
	var miterLimit uint8
	binary.Read(r, binary.LittleEndian, &width)
	binary.Read(r, binary.LittleEndian, &lineOptions)
	binary.Read(r, binary.LittleEndian, &miterLimit)

	t.Width = (float32(width) - 128.0)
	t.LineJoin = LineJoinOptions(lineOptions & 15)
	t.LineCap = LineCapOptions(lineOptions >> 4)
	t.MiterLimit = float32(miterLimit)

	return t
}

type Transformer any

func ReadTransformer(r io.Reader) Transformer {
	var ttype TransformerType
	binary.Read(r, binary.LittleEndian, &ttype)
	switch ttype {
	case TransformerTypeAffine:
		return ReadAffine(r)
	case TransformerTypeContour:
		return ReadCountour(r)
	case TransformerTypePerspective:
		return ReadTransformerPerspective(r)
	case TransformerTypeStroke:
		return ReadTransformerStroke(r)
	}
	return nil
}

func ReadTranslation(r io.Reader) Translation {
	var t Translation
	t.X = ReadFloatCoord(r)
	t.Y = ReadFloatCoord(r)
	return t
}

func ReadLodScale(r io.Reader) LodScale {
	var ls LodScale
	var scale uint8
	binary.Read(r, binary.LittleEndian, &scale)
	ls.MinS = float32(scale) / 63.75
	binary.Read(r, binary.LittleEndian, &scale)
	ls.MaxS = float32(scale) / 63.75

	return ls
}

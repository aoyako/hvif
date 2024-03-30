package utils

import (
	"encoding/binary"
	"io"
)

const TransformMatrixSize = 6
const PerspectiveMatrixSize = 9

type LineJoinOptions uint8

const (
	MiterJoin LineJoinOptions = iota
	MiterJoinRevert
	RoundJoin
	BevelJoin
	MiterJoinRound
)

type LineCapOptions uint8

const (
	ButtCap LineCapOptions = iota
	SquareCap
	RoundCap
)

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

type Transformer any

type TransformerAffine struct {
	Matrix [TransformMatrixSize]float32
}

type TransformerPerspective struct {
	Matrix [PerspectiveMatrixSize]float32
}

type TransformerContour struct {
	Width      float32
	LineJoin   LineJoinOptions
	MiterLimit float32
}

type TransformerStroke struct {
	Width      float32
	LineJoin   LineJoinOptions
	LineCap    LineCapOptions
	MiterLimit float32
}

func ReadTransformable(r io.Reader) Transformable {
	var t Transformable
	copy(t.Matrix[:], ReadMatrix(r, TransformMatrixSize))
	return t
}

func ReadAffine(r io.Reader) TransformerAffine {
	var t TransformerAffine
	copy(t.Matrix[:], ReadMatrix(r, TransformMatrixSize))
	return t
}

func ReadCountour(r io.Reader) TransformerContour {
	var t TransformerContour
	var width uint8
	var lineJoin uint8
	var miterLimit uint8
	binary.Read(r, binary.LittleEndian, &width)
	binary.Read(r, binary.LittleEndian, &lineJoin)
	binary.Read(r, binary.LittleEndian, &miterLimit)

	t.Width = (float32(width) - 128.0)
	t.LineJoin = LineJoinOptions(lineJoin)
	t.MiterLimit = float32(miterLimit)
	return t
}

func ReadTransformerPerspective(r io.Reader) TransformerPerspective {
	var t TransformerPerspective
	copy(t.Matrix[:], ReadMatrix(r, TransformMatrixSize))
	return t
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

package utils

import (
	"encoding/binary"
	"fmt"
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

func ReadCountour(r io.Reader) (TransformerContour, error) {
	var t TransformerContour
	var width uint8
	var lineJoin uint8
	var miterLimit uint8

	if err := binary.Read(r, binary.LittleEndian, &width); err != nil {
		return t, fmt.Errorf("reading width: %w", err)
	}
	if err := binary.Read(r, binary.LittleEndian, &lineJoin); err != nil {
		return t, fmt.Errorf("reading line join options: %w", err)
	}
	if err := binary.Read(r, binary.LittleEndian, &miterLimit); err != nil {
		return t, fmt.Errorf("reading miter limit: %w", err)
	}

	t.Width = (float32(width) - 128.0)
	t.LineJoin = LineJoinOptions(lineJoin)
	t.MiterLimit = float32(miterLimit)

	return t, nil
}

func ReadTransformerPerspective(r io.Reader) TransformerPerspective {
	var t TransformerPerspective
	copy(t.Matrix[:], ReadMatrix(r, TransformMatrixSize))
	return t
}

func ReadTransformerStroke(r io.Reader) (TransformerStroke, error) {
	var t TransformerStroke
	var width uint8
	var lineOptions uint8
	var miterLimit uint8
	if err := binary.Read(r, binary.LittleEndian, &width); err != nil {
		return t, fmt.Errorf("reading width: %w", err)
	}
	if err := binary.Read(r, binary.LittleEndian, &lineOptions); err != nil {
		return t, fmt.Errorf("reading line options: %w", err)
	}
	if err := binary.Read(r, binary.LittleEndian, &miterLimit); err != nil {
		return t, fmt.Errorf("reading miter limit: %w", err)
	}

	t.Width = (float32(width) - 128.0)
	t.LineJoin = LineJoinOptions(lineOptions & 15)
	t.LineCap = LineCapOptions(lineOptions >> 4)
	t.MiterLimit = float32(miterLimit)

	return t, nil
}

func ReadTransformer(r io.Reader) (Transformer, error) {
	var ttype TransformerType
	binary.Read(r, binary.LittleEndian, &ttype)
	switch ttype {
	case TransformerTypeAffine:
		return ReadAffine(r), nil
	case TransformerTypeContour:
		if t, err := ReadCountour(r); err != nil {
			return nil, fmt.Errorf("reading countour: %w", err)
		} else {
			return t, nil
		}
	case TransformerTypePerspective:
		return ReadTransformerPerspective(r), nil
	case TransformerTypeStroke:
		t, err := ReadTransformerStroke(r)
		if err != nil {
			return t, fmt.Errorf("read stroke transformer: %w", err)
		}
		return t, nil
	}
	return nil, nil
}

func ReadTranslation(r io.Reader) (Translation, error) {
	var t Translation
	x, err := ReadFloatCoord(r)
	if err != nil {
		return t, fmt.Errorf("read x coord: %w", err)
	}
	y, err := ReadFloatCoord(r)
	if err != nil {
		return t, fmt.Errorf("read y coord: %w", err)
	}
	t.X = x
	t.Y = y
	return t, nil
}

func ReadLodScale(r io.Reader) (LodScale, error) {
	var ls LodScale
	var minScale uint8
	var maxScale uint8
	if err := binary.Read(r, binary.LittleEndian, &minScale); err != nil {
		return ls, fmt.Errorf("reading min scale: %w", err)
	}
	if err := binary.Read(r, binary.LittleEndian, &maxScale); err != nil {
		return ls, fmt.Errorf("reading max scale: %w", err)
	}
	ls.MinS = float32(minScale) / 63.75
	ls.MaxS = float32(maxScale) / 63.75

	return ls, nil
}

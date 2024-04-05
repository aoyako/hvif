package main

import (
	"encoding/binary"
	"fmt"
	"io"
)

const (
	transformMatrixSize   = 6
	perspectiveMatrixSize = 9
)

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

type transformerType uint8

const (
	transformerTypeAffine transformerType = 20 + iota
	transformerTypeContour
	transformerTypePerspective
	transformerTypeStroke
)

// TransformerTranslation | TransformerLodScale | TransformerAffine | TransformerPerspective | TransformerContour | TransformerStroke
type Transformer any

type TransformerTranslation struct {
	X float32
	Y float32
}

type TransformerLodScale struct {
	MinS float32
	MaxS float32
}

type TransformerAffine struct {
	Matrix [transformMatrixSize]float32
}

type TransformerPerspective struct {
	Matrix [perspectiveMatrixSize]float32
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

func readAffine(r io.Reader) (TransformerAffine, error) {
	var t TransformerAffine
	mx, err := ReadMatrix(r, transformMatrixSize)
	if err != nil {
		return t, fmt.Errorf("reading matrix: %w", err)
	}
	copy(t.Matrix[:], mx)
	return t, nil
}

func readCountour(r io.Reader) (TransformerContour, error) {
	var t TransformerContour

	var width uint8
	if err := binary.Read(r, binary.LittleEndian, &width); err != nil {
		return t, fmt.Errorf("reading width: %w", err)
	}
	t.Width = (float32(width) - 128.0)

	var lineJoin uint8
	if err := binary.Read(r, binary.LittleEndian, &lineJoin); err != nil {
		return t, fmt.Errorf("reading line join options: %w", err)
	}
	t.LineJoin = LineJoinOptions(lineJoin)

	var miterLimit uint8
	if err := binary.Read(r, binary.LittleEndian, &miterLimit); err != nil {
		return t, fmt.Errorf("reading miter limit: %w", err)
	}
	t.MiterLimit = float32(miterLimit)

	return t, nil
}

func readTransformerPerspective(r io.Reader) (TransformerPerspective, error) {
	var t TransformerPerspective
	mx, err := ReadMatrix(r, perspectiveMatrixSize)
	if err != nil {
		return t, fmt.Errorf("reading matrix: %w", err)
	}
	copy(t.Matrix[:], mx)
	return t, nil
}

func readTransformerStroke(r io.Reader) (TransformerStroke, error) {
	var t TransformerStroke

	var width uint8
	if err := binary.Read(r, binary.LittleEndian, &width); err != nil {
		return t, fmt.Errorf("reading width: %w", err)
	}
	t.Width = (float32(width) - 128.0)

	var lineOptions uint8
	if err := binary.Read(r, binary.LittleEndian, &lineOptions); err != nil {
		return t, fmt.Errorf("reading line options: %w", err)
	}
	t.LineJoin = LineJoinOptions(lineOptions & 15)
	t.LineCap = LineCapOptions(lineOptions >> 4)

	var miterLimit uint8
	if err := binary.Read(r, binary.LittleEndian, &miterLimit); err != nil {
		return t, fmt.Errorf("reading miter limit: %w", err)
	}
	t.MiterLimit = float32(miterLimit)

	return t, nil
}

func readTransformer(r io.Reader) (Transformer, error) {
	var ttype transformerType
	err := binary.Read(r, binary.LittleEndian, &ttype)
	if err != nil {
		return nil, fmt.Errorf("reading type: %w", err)
	}

	switch ttype {
	case transformerTypeAffine:
		t, err := readAffine(r)
		if err != nil {
			return nil, fmt.Errorf("reading affine transformer: %w", err)
		}
		return t, nil
	case transformerTypeContour:
		if t, err := readCountour(r); err != nil {
			return nil, fmt.Errorf("reading countour: %w", err)
		} else {
			return t, nil
		}
	case transformerTypePerspective:
		t, err := readTransformerPerspective(r)
		if err != nil {
			return nil, fmt.Errorf("reading perspective transformer: %w", err)
		}
		return t, nil
	case transformerTypeStroke:
		t, err := readTransformerStroke(r)
		if err != nil {
			return t, fmt.Errorf("read stroke transformer: %w", err)
		}
		return t, nil
	}
	return nil, nil
}

func readTranslation(r io.Reader) (TransformerTranslation, error) {
	var t TransformerTranslation
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

func readLodScale(r io.Reader) (TransformerLodScale, error) {
	var ls TransformerLodScale
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
